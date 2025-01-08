package usecase

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/abdisetiakawan/go-clean-arch/internal/entity"
	"github.com/abdisetiakawan/go-clean-arch/internal/helper"
	"github.com/abdisetiakawan/go-clean-arch/internal/model"
	"github.com/abdisetiakawan/go-clean-arch/internal/model/converter"
	"github.com/abdisetiakawan/go-clean-arch/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TaskTagUseCase struct {
	DB            *gorm.DB
	Log           *logrus.Logger
	Validate      *validator.Validate
	TaskTagRepository *repository.TaskTagRepository
	Cache *helper.CacheHelper
}

func NewTaskTagUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, taskTagRepository *repository.TaskTagRepository, cache *helper.CacheHelper) *TaskTagUseCase {
	return &TaskTagUseCase{
		DB:            db,
		Log:           log,
		Validate:      validate,
		TaskTagRepository: taskTagRepository,
		Cache: cache,
	}
}

func (c *TaskTagUseCase) Create(ctx context.Context, request *model.CreateTaskTagRequest, email string) (*model.TaskTagResponse, error) {
    tx := c.DB.WithContext(ctx).Begin()
    defer tx.Rollback()
    if err := c.Validate.Struct(request); err != nil {
        c.Log.WithError(err).Error("error validate request body")
        return nil, model.ErrBadRequest
    }
    taskTag := &entity.TaskTag{
        TaskId: request.TaskId,
        TagId: request.TagId,
    }
    // Memastikan hanya satu tag untuk satu task
    if err := c.TaskTagRepository.CheckTaskTag(tx, int(request.TaskId), int(request.TagId)); err == nil {
        c.Log.Info("Task tag already exists")
        return nil, model.ErrConflict
    } else if err != gorm.ErrRecordNotFound {
        c.Log.WithError(err).Error("error checking task tag")
        return nil, model.ErrInternalServer
    }

    if err := c.TaskTagRepository.CreateTaskTag(tx, taskTag, email); err != nil {
        c.Log.WithError(err).Error("error create task tag")
        return nil, err
    }
    
    if err := tx.Commit().Error; err != nil {
        c.Log.WithError(err).Error("error create task tag")
        return nil, model.ErrInternalServer
    }

    // Clear all related caches
    c.Cache.Delete(ctx, "task_tags:" + email)
    c.Cache.Delete(ctx, "task_tags:search:" + email)
    c.Cache.Delete(ctx, "task_tags:" + strconv.Itoa(int(request.TagId)))

    return converter.TaskTagToResponse(taskTag), nil
}


func (c *TaskTagUseCase) Search(ctx context.Context, request *model.SearchTaskTagRequest) ([]model.TaskTagResult, int64, error) {
	cacheKey := "task_tags:search:" + request.Email
	var cachedData struct {
		Responses []model.TaskTagResult
		Total     int64
	}
	if err := c.Cache.GetAndUnmarshal(ctx, cacheKey, &cachedData); err == nil {
		return cachedData.Responses, cachedData.Total, nil
	}
	
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validate request query")
		return nil, 0, model.ErrBadRequest
	}

	taskTags, total, err := c.TaskTagRepository.SearchTaskTag(tx, request)
	if err != nil {
		c.Log.WithError(err).Error("error search task tag")
		return nil, 0, err
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error search task tag")
		return nil, 0, model.ErrInternalServer
	}

	responses := make([]model.TaskTagResult, len(taskTags))
	for i, taskTag := range taskTags {
		responses[i] = *converter.TaskWithTagsToResponse(&taskTag)
	}

	cachedData.Responses = responses
	cachedData.Total = total
	cachedDataJSON, _ := json.Marshal(cachedData)
	c.Cache.Set(ctx, cacheKey, cachedDataJSON, 1 * time.Minute)

	return responses, total, nil
}

func (c *TaskTagUseCase) SearchTaskTagRequestWithTagId(ctx context.Context, request *model.SearchTaskTagRequestWithTagId) ([]model.TaskTagResult, int64, error) {
	cacheKey := "task_tags:" + strconv.Itoa(int(request.TagId))
	var cachedData struct {
		Responses []model.TaskTagResult
		Total     int64
	}
	if err := c.Cache.GetAndUnmarshal(ctx, cacheKey, &cachedData); err == nil {
		return cachedData.Responses, cachedData.Total, nil
	}
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()
	
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validate request query")
		return nil, 0, model.ErrBadRequest
	}
	taskTags, total, err := c.TaskTagRepository.SearchTaskTagRequestWithTagId(tx, request)
	if err != nil {
		c.Log.WithError(err).Error("error search task tag request with tag id")
		return nil, 0, err
	}
	if err := tx.Commit().Error; err != nil {
	c.Log.WithError(err).Error("error search task tag request with tag id")
	return nil, 0, model.ErrInternalServer
	}

	responses := make([]model.TaskTagResult, len(taskTags))
	for i, taskTag := range taskTags {
		responses[i] = *converter.TaskWithTagsToResponse(&taskTag)
	}

	cachedData.Responses = responses
	cachedData.Total = total
	cachedDataJSON, _ := json.Marshal(cachedData)
	c.Cache.Set(ctx, cacheKey, cachedDataJSON, 1 * time.Minute)

	return responses, total, nil
}

func (c *TaskTagUseCase) Delete(ctx context.Context, request *model.GetTaskTagForDelete) error {
    tx := c.DB.WithContext(ctx).Begin()
    defer tx.Rollback()

    taskTag := new(entity.TaskTag)
    if err := c.Validate.Struct(request); err != nil {
        c.Log.WithError(err).Error("error validate request query")
        return model.ErrBadRequest
    }
    if err := c.TaskTagRepository.CheckIsAdded(tx, taskTag, request); err != nil {
        c.Log.WithError(err).Error("error check is added task tag")
        return err
    }
    if err := c.TaskTagRepository.Delete(tx, taskTag); err != nil {
        c.Log.WithError(err).Error("error delete task tag")
        return err
    }
    if err := tx.Commit().Error; err != nil {
        c.Log.WithError(err).Error("error delete task tag")
        return model.ErrInternalServer
    }

    // Clear all related caches
    c.Cache.Delete(ctx, "task_tags:" + request.Email)
    c.Cache.Delete(ctx, "task_tags:search:" + request.Email)
    c.Cache.Delete(ctx, "task_tags:" + strconv.Itoa(int(taskTag.TagId)))
    
    return nil
}