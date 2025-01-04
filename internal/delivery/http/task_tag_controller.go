package http

import (
	"strconv"

	"github.com/abdisetiakawan/go-clean-arch/internal/delivery/http/middleware"
	"github.com/abdisetiakawan/go-clean-arch/internal/model"
	"github.com/abdisetiakawan/go-clean-arch/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type TaskTagController struct {
	UseCase *usecase.TaskTagUseCase
	Log     *logrus.Logger
}

func NewTaskTagController(useCase *usecase.TaskTagUseCase, logger *logrus.Logger) *TaskTagController {
	return &TaskTagController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *TaskTagController) Create(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := new(model.CreateTaskTagRequest)
	if err := ctx.BodyParser(&request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return model.ErrBadRequest
	}
	taskIdStr := ctx.Params("taskId")
	taskId, err := strconv.ParseUint(taskIdStr, 10, 32)
    if err != nil {
        c.Log.Warnf("Invalid task ID : %+v", err)
        return model.ErrBadRequest
    }
	request.TaskId = uint(taskId)
	response, err := c.UseCase.Create(ctx.UserContext(), request, auth.Email)
	if err != nil {
		c.Log.Warnf("Failed to create task tag : %+v", err)
		return err
	}
	return ctx.JSON(model.WebResponse[*model.TaskTagResponse]{Data: response})
}