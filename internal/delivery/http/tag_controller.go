package http

import (
	"math"

	"github.com/abdisetiakawan/go-clean-arch/internal/delivery/http/middleware"
	"github.com/abdisetiakawan/go-clean-arch/internal/model"
	"github.com/abdisetiakawan/go-clean-arch/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type TagsController struct {
	UseCase *usecase.TagUseCase
	Log     *logrus.Logger
}

func NewTagsController(useCase *usecase.TagUseCase, logger *logrus.Logger) *TagsController {
	return &TagsController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *TagsController) Create(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := new(model.CreateTagRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("failed to parse request body: %+v", err)
		return model.ErrBadRequest
	}
	request.Email = auth.Email
	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("failed to create tag: %+v", err)
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[*model.TagResponse]{Data: response})
}

func (c *TagsController) List(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := &model.SearchTagRequest{
		Email: auth.Email,
		Name: ctx.Query("name", ""),
		Page: ctx.QueryInt("page", 1),
		Size: ctx.QueryInt("size", 10),
	}

	responses, total, err := c.UseCase.Search(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("failed to list tag: %+v", err)
		return err
	}
	paging := &model.PageMetadata{
		Page: request.Page,
		Size: request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[[]model.TagResponse]{Paging: paging, Data: responses})
}

func (c *TagsController) Get(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := &model.GetTagRequest{
		ID: ctx.Params("tagId"),
		Email: auth.Email,
	}

	response, err := c.UseCase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to get tag : %+v", err)
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[*model.TagResponse]{Data: response})
}

func (c *TagsController) Update(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := new(model.UpdateTagRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}
	request.Email = auth.Email
	request.ID = ctx.Params("tagId")
	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to update tag : %+v", err)
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[*model.TagResponse]{Data: response})
}


func (c *TagsController) Delete(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := &model.GetTagRequest{
		ID: ctx.Params("tagId"),
		Email: auth.Email,
	}

	if err := c.UseCase.Delete(ctx.UserContext(), request); err != nil {
		c.Log.Warnf("Failed to delete tag : %+v", err)
		return err
	}
	
	return ctx.Status(fiber.StatusNoContent).JSON(model.WebResponse[bool]{Data: true})
}