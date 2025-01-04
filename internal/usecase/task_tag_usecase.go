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
