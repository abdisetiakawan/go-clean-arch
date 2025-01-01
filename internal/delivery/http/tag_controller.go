package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"github.com/abdisetiakawan/go-clean-arch/internal/model"
	"github.com/abdisetiakawan/go-clean-arch/internal/usecase"
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

func (c *TagsController) List(ctx *fiber.Ctx) error {
	responses, err := c.UseCase.List(ctx.UserContext())
	if err != nil {
		c.Log.Warnf("Failed to list tags : %+v", err)
		return err
	}
	return ctx.JSON(model.WebResponse[[]model.TagResponse]{Data: responses})
}

func (c *TagsController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateTagRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}
	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create tag : %+v", err)
	}
	return ctx.JSON(model.WebResponse[*model.TagResponse]{Data: response})
}

func (c *TagsController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateTagRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}
	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to update tag : %+v", err)
	}
	return ctx.JSON(model.WebResponse[*model.TagResponse]{Data: response})
}

func (c *TagsController) Delete(ctx *fiber.Ctx) error {
	request := new(model.DeleteTagRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}
	response, err := c.UseCase.Delete(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to delete tag : %+v", err)
	}
	return ctx.JSON(model.WebResponse[*model.TagResponse]{Data: response})
}

func (c *TagsController) Get(ctx *fiber.Ctx) error {
	request := new(model.GetTagRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}
	response, err := c.UseCase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to get tag : %+v", err)
	}
	return ctx.JSON(model.WebResponse[*model.TagResponse]{Data: response})
}