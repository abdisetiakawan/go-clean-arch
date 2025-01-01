package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type TaskTagController struct {
	Log     *logrus.Logger
	UseCase *usecase.TaskTagUseCase
}

func NewTaskTagController(useCase *usecase.TaskTagUseCase, logger *logrus.Logger) *TaskTagController {
	return &TaskTagController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *TaskTagController) List(ctx *fiber.Ctx) error {
	responses, err := c.UseCase.List(ctx.UserContext())
	if err != nil {
		c.Log.Warnf("Failed to list task tags : %+v", err)
		return err
	}
	return ctx.JSON(model.WebResponse[[]model.TaskTagResponse]{Data: responses})
}

func (c *TaskTagController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateTaskTagRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}
	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create task tag : %+v", err)
	}
	return ctx.JSON(model.WebResponse[*model.TaskTagResponse]{Data: response})
}