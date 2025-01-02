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

type TaskUseCase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	TaskRepository *repository.TaskRepository
}

func NewTaskUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate, taskRepository *repository.TaskRepository) *TaskUseCase {
	return &TaskUseCase{
		DB: db,
		Log: logger,
		Validate: validate,
		TaskRepository: taskRepository,
	}
}

func (c *TaskUseCase) Create(ctx context.Context, request *model.CreateTaskRequest) (*model.TaskResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validate request body")
		return nil, model.ErrBadRequest
	}
	task := &entity.Task{
		Email: request.Email,
		Title: request.Title,
		Description: request.Description,
		Status: request.Status,
		DueDate: request.DueDate,
	}
	if err := c.TaskRepository.Create(tx, task); err != nil {
		c.Log.WithError(err).Error("error create task")
		return nil, model.ErrInternalServer
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error create task")
		return nil, model.ErrInternalServer
	}

	return converter.TaskToResponse(task), nil
}

func (c *TaskUseCase) Search(ctx context.Context, request *model.SearchTaskRequest) ([]model.TaskResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validate request body")
		return nil, 0, model.ErrBadRequest
	}
	tasks, total, err := c.TaskRepository.Search(tx, request)
	if err != nil {
		c.Log.WithError(err).Error("error search task")
		return nil, 0, model.ErrInternalServer
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error search task")
		return nil, 0, model.ErrInternalServer
	}
	responses := make([]model.TaskResponse, len(tasks))
	for i, task := range tasks {
		task.Email = ""
		responses[i] = *converter.TaskToResponse(&task)
	}
	return responses, total, nil
}

func (c *TaskUseCase) Get(ctx context.Context, request *model.GetTaskRequest) (*model.TaskResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validate request query")
		return nil, model.ErrBadRequest
	}
	task := new(entity.Task)
	if err := c.TaskRepository.FindByEmailAndId(tx, task, request.ID, request.Email); err != nil {
		c.Log.WithError(err).Error("error search task")
		return nil, model.ErrNotFound
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error search task")
		return nil, model.ErrInternalServer
	}
	return converter.TaskToResponse(task), nil
}