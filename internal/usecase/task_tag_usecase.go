package usecase

import (
	"context"

	"github.com/abdisetiakawan/go-clean-arch/internal/entity"
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
}

func NewTaskTagUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, taskTagRepository *repository.TaskTagRepository) *TaskTagUseCase {
	return &TaskTagUseCase{
		DB:            db,
		Log:           log,
		Validate:      validate,
		TaskTagRepository: taskTagRepository,
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
	if err := c.TaskTagRepository.CheckIsAdded(tx, taskTag); err != nil {
		c.Log.WithError(err).Error("error check availability task tag")
		return nil, err
	}

	if err := c.TaskTagRepository.CreateTaskTag(tx, taskTag, email); err != nil {
		c.Log.WithError(err).Error("error create task tag")
		return nil, err
	}
	
    if err := tx.Commit().Error; err != nil {
        c.Log.WithError(err).Error("error create task tag")
        return nil, model.ErrInternalServer
    }

	return converter.TaskTagToResponse(taskTag), nil
}

func (c *TaskTagUseCase) Search(ctx context.Context, request *model.SearchTaskTagRequest) ([]model.TaskTagResult, int64, error) {
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

	return responses, total, nil
}

func (c *TaskTagUseCase) SearchTaskTagRequestWithTagId(ctx context.Context, request *model.SearchTaskTagRequestWithTagId) ([]model.TaskTagResult, int64, error) {
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

	return responses, total, nil
}