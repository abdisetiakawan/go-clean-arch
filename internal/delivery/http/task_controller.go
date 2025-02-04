package http

import (
	"math"

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
	request := &model.SearchTaskRequest{
		Email: auth.Email,
		Title: ctx.Query("title", ""),
		Description: ctx.Query("description", ""),
		Status: ctx.Query("status", ""),
		Page: ctx.QueryInt("page", 1),
		Size: ctx.QueryInt("size", 10),
	}

	responses, total, err := c.UseCase.Search(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to list tasks : %+v", err)
		return err
	}
	paging := &model.PageMetadata{
		Page: request.Page,
		Size: request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(responses, "Tasks fetched successfully", fiber.StatusOK, paging))
}

func (c *TaskController) Create(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := new(model.CreateTaskRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return model.ErrBadRequest
	}
	request.Email = auth.Email
	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create task : %+v", err)
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.NewWebResponse(response, "Successfully created task", fiber.StatusCreated, nil))
}

func (c *TaskController) Update(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := new(model.UpdateTaskRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}
	request.Email = auth.Email
	request.ID = ctx.Params("taskId")
	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to update task : %+v", err)
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(response, "Successfully updated task", fiber.StatusOK, nil))
}

func (c *TaskController) Get(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := &model.GetTaskRequest{
		ID: ctx.Params("taskId"),
		Email: auth.Email,
	}

	response, err := c.UseCase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to get task : %+v", err)
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(response, "Successfully get task", fiber.StatusOK, nil))
}

func (c *TaskController) Delete(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := &model.GetTaskRequest{
		ID: ctx.Params("taskId"),
		Email: auth.Email,
	}

	if err := c.UseCase.Delete(ctx.UserContext(), request); err != nil {
		c.Log.Warnf("Failed to delete task : %+v", err)
		return err
	}
	
	return ctx.SendStatus(fiber.StatusNoContent)
}