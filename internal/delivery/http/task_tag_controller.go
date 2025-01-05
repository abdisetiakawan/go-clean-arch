package http

import (
	"math"
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

func (c *TaskTagController) List(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := &model.SearchTaskTagRequest{
		Email: auth.Email,
		Page: ctx.QueryInt("page", 1),
		Size: ctx.QueryInt("size", 10),
	}
	responses, total, err := c.UseCase.Search(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to list task tags : %+v", err)
		return err
	}
	paging := &model.PageMetadata{
		Page: request.Page,
		Size: request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
		}

	return ctx.JSON(model.WebResponse[[]model.TaskTagResult]{Paging: paging, Data: responses})
}

func (c *TaskTagController) ListByTagId(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	TagId := ctx.Params("tagId")
	tagId, err := strconv.ParseUint(TagId, 10, 32)
	if err != nil {
		c.Log.Warnf("Invalid tag ID : %+v", err)
		return model.ErrBadRequest
	}
	request := &model.SearchTaskTagRequestWithTagId{
		Email: auth.Email,
		TagId: uint(tagId),
		Page: ctx.QueryInt("page", 1),
		Size: ctx.QueryInt("size", 10),
	}
	responses, total, err := c.UseCase.SearchTaskTagRequestWithTagId(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to list task tags : %+v", err)
		return err
	}
	paging := &model.PageMetadata{
		Page: request.Page,
		Size: request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
		}

	return ctx.JSON(model.WebResponse[[]model.TaskTagResult]{Paging: paging, Data: responses})
}