package http

import (
	"github.com/abdisetiakawan/go-clean-arch/internal/delivery/http/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"github.com/abdisetiakawan/go-clean-arch/internal/model"
	"github.com/abdisetiakawan/go-clean-arch/internal/usecase"
)

type TaskController struct {
	UseCase *usecase.TaskUseCase
	Log     *logrus.Logger
}

func NewTaskController(useCase *usecase.TaskUseCase, logger *logrus.Logger) *TaskController {
	return &TaskController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *TaskController) List(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := &model.ListTaskRequest{
		UserID: auth.ID,
		Title: ctx.Query("title", ""),
		Description: ctx.Query("description", ""),
		Page: ctx.QueryInt("page", 1),
		Size: ctx.QueryInt("size", 10),
	}

	responses, total, err := c.UseCase.List(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to list tasks : %+v", err)
		return err
	}
	paging := &model.PageMetaData{
		Page: request.Page,
		Size: request.Size,
		Total: total,
		TotalPage: int(total / int64(request.Size)) + 1,
	}
	return ctx.JSON(model.WebResponse[[]model.ListTaskResponse]{Paging: responses, Meta: paging})
}

func (c *TaskController) Create(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := new(model.CreateTaskRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}
	request.UserID = auth.ID
	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create task : %+v", err)
	}
	return ctx.JSON(model.WebResponse[*model.CreateTaskResponse]{Data: response})
}

func (c *TaskController) Update(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := new(model.UpdateTaskRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}
	request.UserID = auth.ID
	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to update task : %+v", err)
	}
	return ctx.JSON(model.WebResponse[*model.UpdateTaskResponse]{Data: response})
}

func (c *TaskController) Get(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := &model.GetTaskRequest{
		ID: ctx.Params("taskId"),
		UserID: auth.ID,
	}

	response, err := c.UseCase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to get task : %+v", err)
	}
	return ctx.JSON(model.WebResponse[*model.GetTaskResponse]{Data: response})
}

func (c *TaskController) Delete(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := &model.DeleteTaskRequest{
		ID: ctx.Params("taskId"),
		UserID: auth.ID,
	}
	response, err := c.UseCase.Delete(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to delete task : %+v", err)
	}
	return ctx.JSON(model.WebResponse[*model.DeleteTaskResponse]{Data: response})
}