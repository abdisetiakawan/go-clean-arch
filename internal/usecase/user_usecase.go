package usecase

import (
	"context"

	"github.com/abdisetiakawan/go-clean-arch/internal/entity"
	"github.com/abdisetiakawan/go-clean-arch/internal/helper"
	"github.com/abdisetiakawan/go-clean-arch/internal/model"
	"github.com/abdisetiakawan/go-clean-arch/internal/model/converter"
	"github.com/abdisetiakawan/go-clean-arch/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase struct {
    DB             *gorm.DB
    Log            *logrus.Logger
    Validate       *validator.Validate
    UserRepository *repository.UserRepository
    Helper         *helper.Helper
}

func NewUserUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, userRepository *repository.UserRepository, helper *helper.Helper) *UserUseCase {
    return &UserUseCase{
        DB:             db,
        Log:            log,
        Validate:       validate,
        UserRepository: userRepository,
        Helper:         helper,
    }
}

func (c *UserUseCase) Create(ctx context.Context, request *model.CreateUserRequest) (*model.UserResponse, error) {
    tx := c.DB.WithContext(ctx).Begin()
    defer tx.Rollback()

    err := c.Validate.Struct(request)
    if err != nil {
        c.Log.Warnf("Failed to validate request body : %+v", err)
        return nil, fiber.ErrBadRequest
    }
    
    total, err := c.UserRepository.CountByEmail(tx, request.Email)
    if err != nil {
        c.Log.Warnf("Failed to count user : %+v", err)
        return nil, fiber.ErrInternalServerError
    }
    if total > 0 {
        c.Log.Warnf("User already exists : %+v", err)
        return nil, fiber.ErrConflict
    }
    password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
    if err != nil {
        c.Log.Warnf("Failed to hash password : %+v", err)
        return nil, fiber.ErrInternalServerError
    }
    access_token, refreshToken, err := (*c.Helper).GenerateTokenUser(model.UserResponse{
		Name:  request.Name,
		Email: request.Email,
	})
    if err != nil {
        c.Log.Warnf("Failed to generate tokens : %+v", err)
        return nil, fiber.ErrInternalServerError
    }

    user := &entity.User{
        Name:         request.Name,
        Email:        request.Email,
        Password:     string(password),
        Token: refreshToken,
		AccessToken: access_token,
    }

    err = c.UserRepository.Create(tx, user)
    if err != nil {
        c.Log.Warnf("Failed to create user : %+v", err)
        return nil, fiber.ErrInternalServerError
    }

    err = tx.Commit().Error
    if err != nil {
        c.Log.Warnf("Failed to commit transaction : %+v", err)
        return nil, fiber.ErrInternalServerError
    }

    return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Login(ctx context.Context, request *model.LoginUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Failed to validate request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	user := new(entity.User)

	err = c.UserRepository.FindByEmail(tx, user, request.Email)
	if err != nil {
		c.Log.Warnf("Failed to find user : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		c.Log.Warnf("Failed to compare password : %+v", err)
		return nil, fiber.ErrUnauthorized
	}

    accessToken, refreshToken, err := (*c.Helper).GenerateTokenUser(model.UserResponse{
		Name:  user.Name,
		Email: user.Email,
	})
	if err != nil {
		c.Log.Warnf("Failed to generate tokens : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	user.AccessToken = accessToken
	user.Token = refreshToken

	err = c.UserRepository.Update(tx, user)
	if err != nil {
		c.Log.Warnf("Failed to update user : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	err = tx.Commit().Error
	if err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Logout(ctx context.Context, request *model.LogoutUserRequest) (bool, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Failed to validate request body : %+v", err)
		return false, fiber.ErrBadRequest
	}

	user := new(entity.User)

	err = c.UserRepository.FindByEmail(tx, user, request.Email)
	if err != nil {
		c.Log.Warnf("Failed to find user : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	user.Token = ""

	err = c.UserRepository.Update(tx, user)
	if err != nil {
		c.Log.Warnf("Failed to update user : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	err = tx.Commit().Error
	if err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	return true, nil
}

func (c *UserUseCase) Update(ctx context.Context, request *model.UpdateUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Failed to validate request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}
	user := new(entity.User)

	err = c.UserRepository.FindByEmail(tx, user, request.Email)
	if err != nil {
		c.Log.Warnf("Failed to find user : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if request.Name != "" {
		user.Name = request.Name
	}

	if request.Email != "" {
		user.Email = request.Email
	}

	if request.Password != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			c.Log.Warnf("Failed to generate password : %+v", err)
			return nil, fiber.ErrInternalServerError
		}
		user.Password = string(password)
	}

	err = c.UserRepository.Update(tx, user)
	if err != nil {
		c.Log.Warnf("Failed to update user : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	err = tx.Commit().Error
	if err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Current(ctx context.Context, request *model.GetUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Failed to validate request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}
	user := new(entity.User)

	err = c.UserRepository.FindByEmail(tx, user, request.Email)
	if err != nil {
		c.Log.Warnf("Failed to find user : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	err = tx.Commit().Error
	if err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Verify(ctx context.Context, request *model.VerifyUserRequest) (*model.Auth, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Failed to validate request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}
	user := new(entity.User)
	if err := c.UserRepository.FindByToken(tx, user, request.Token); err != nil {
		c.Log.Warnf("Failed to find user : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	err = tx.Commit().Error
	if err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	return &model.Auth{Email: user.Email}, nil
}