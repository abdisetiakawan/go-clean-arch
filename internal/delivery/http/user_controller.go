package http

import (
	"github.com/abdisetiakawan/go-clean-arch/internal/delivery/http/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"github.com/abdisetiakawan/go-clean-arch/internal/model"
	"github.com/abdisetiakawan/go-clean-arch/internal/usecase"
)

type UserController struct {
	Log     *logrus.Logger
	UseCase *usecase.UserUseCase
}

func NewUserController(useCase *usecase.UserUseCase, logger *logrus.Logger) *UserController {
	return &UserController{
		Log: logger,
		UseCase: useCase,
	}
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	request := new(model.CreateUserRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return model.ErrBadRequest
	}
	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to register user : %+v", err)
		return err
	}
	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginUserRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return model.ErrBadRequest
	}
	response, err := c.UseCase.Login(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to login user : %+v", err)
		return err
	}
	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

func (c *UserController) Logout(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := &model.LogoutUserRequest{Email: auth.Email}
	response, err := c.UseCase.Logout(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to logout user : %+v", err)
		return model.ErrInternalServer
	}
	return ctx.JSON(model.WebResponse[bool]{Data: response})
}

func (c *UserController) Current(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := &model.GetUserRequest{Email: auth.Email}
	response, err := c.UseCase.Current(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to get current user : %+v", err)
		return model.ErrInternalServer
	}
	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

func (c *UserController) Update(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := new(model.UpdateUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return model.ErrBadRequest
	}
	request.Email = auth.Email
	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to update user : %+v", err)
		return model.ErrInternalServer
	}
	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}