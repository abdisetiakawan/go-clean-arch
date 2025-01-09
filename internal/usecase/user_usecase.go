package usecase

import (
	"context"
	"encoding/json"
	"time"

	"github.com/abdisetiakawan/go-clean-arch/internal/entity"
	"github.com/abdisetiakawan/go-clean-arch/internal/helper"
	"github.com/abdisetiakawan/go-clean-arch/internal/model"
	"github.com/abdisetiakawan/go-clean-arch/internal/model/converter"
	"github.com/abdisetiakawan/go-clean-arch/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase struct {
    DB             *gorm.DB
    Log            *logrus.Logger
    Validate       *validator.Validate
    UserRepository *repository.UserRepository
    Jwt            *helper.JwtHelper
    Cache          *helper.CacheHelper
}

func NewUserUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, userRepository *repository.UserRepository, jwt *helper.JwtHelper, cache *helper.CacheHelper) *UserUseCase {
    return &UserUseCase{
        DB:             db,
        Log:            log,
        Validate:       validate,
        UserRepository: userRepository,
        Jwt:            jwt,
        Cache:          cache,
    }
}

func (c *UserUseCase) Create(ctx context.Context, request *model.CreateUserRequest) (*model.UserResponse, error) {
    tx := c.DB.WithContext(ctx).Begin()
    defer tx.Rollback()

    err := c.Validate.Struct(request)
    if err != nil {
        c.Log.Warnf("Failed to validate request body : %+v", err)
        return nil, model.ErrBadRequest
    }

    total, err := c.UserRepository.CountByEmail(tx, request.Email)
    if err != nil {
        c.Log.Warnf("Failed to count user : %+v", err)
        return nil, model.ErrInternalServer
    }
    if total > 0 {
        return nil, model.ErrUserAlreadyExists
    }

    password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
    if err != nil {
        c.Log.Warnf("Failed to hash password : %+v", err)
        return nil, model.ErrInternalServer
    }
    
    accessToken, refreshToken, err := c.Jwt.GenerateTokenUser(model.UserResponse{
        Name:  request.Name,
        Email: request.Email,
    })
    if err != nil {
        c.Log.Warnf("Failed to generate tokens : %+v", err)
        return nil, model.ErrInternalServer
    }

    user := &entity.User{
        Name:     request.Name,
        Email:    request.Email,
        Password: string(password),
        AccessToken: accessToken,
    }

    err = c.UserRepository.Create(tx, user)
    if err != nil {
        c.Log.Warnf("Failed to create user : %+v", err)
        return nil, model.ErrInternalServer
    }


    err = tx.Commit().Error
    if err != nil {
        c.Log.Warnf("Failed to commit transaction : %+v", err)
        return nil, model.ErrInternalServer
    }

    sessionKey := "session:" + user.Email
    sessionData := map[string]string{
        "accessToken":  accessToken,
        "refreshToken": refreshToken,
    }
    sessionDataJSON, _ := json.Marshal(sessionData)
    
    err = c.Cache.Set(ctx, sessionKey, sessionDataJSON, 1*time.Hour)
    if err != nil {
        c.Log.Warnf("Failed to store session : %+v", err)
        return nil, model.ErrInternalServer
    }

    return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Login(ctx context.Context, request *model.LoginUserRequest) (*model.UserResponse, error) {
    tx := c.DB.WithContext(ctx).Begin()
    defer tx.Rollback()

    err := c.Validate.Struct(request)
    if err != nil {
        c.Log.Warnf("Failed to validate request body : %+v", err)
        return nil, model.ErrBadRequest
    }

    user := new(entity.User)
    err = c.UserRepository.FindByEmail(tx, user, request.Email)
    if err != nil {
        c.Log.Warnf("Failed to find user : %+v", err)
        return nil, model.ErrInvalidCredentials
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
    if err != nil {
        c.Log.Warnf("Failed to compare password : %+v", err)
        return nil, model.ErrInvalidCredentials
    }

    accessToken, refreshToken, err := c.Jwt.GenerateTokenUser(model.UserResponse{
        Name:  user.Name,
        Email: user.Email,
    })
    if err != nil {
        c.Log.Warnf("Failed to generate tokens : %+v", err)
        return nil, model.ErrInternalServer
    }
    user.AccessToken = accessToken

    err = tx.Commit().Error
    if err != nil {
        c.Log.Warnf("Failed to commit transaction : %+v", err)
        return nil, model.ErrInternalServer
    }

    sessionKey := "session:" + user.Email
    sessionData := map[string]string{
        "accessToken":  accessToken,
        "refreshToken": refreshToken,
    }
    sessionDataJSON, _ := json.Marshal(sessionData)
    
    err = c.Cache.Set(ctx, sessionKey, sessionDataJSON, 1*time.Hour)
    if err != nil {
        c.Log.Warnf("Failed to store session : %+v", err)
        return nil, model.ErrInternalServer
    }

    return converter.UserToResponse(user), nil
}


func (c *UserUseCase) Update(ctx context.Context, request *model.UpdateUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Failed to validate request body : %+v", err)
		return nil, model.ErrBadRequest
	}
	user := new(entity.User)

	err = c.UserRepository.FindByEmail(tx, user, request.Email)
	if err != nil {
		c.Log.Warnf("Failed to find user : %+v", err)
		return nil, model.ErrInternalServer
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
			return nil, model.ErrInternalServer
		}
		user.Password = string(password)
	}

	err = c.UserRepository.Update(tx, user)
	if err != nil {
		c.Log.Warnf("Failed to update user : %+v", err)
		return nil, model.ErrInternalServer
	}

	err = tx.Commit().Error
	if err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, model.ErrInternalServer
	}

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Current(ctx context.Context, request *model.GetUserRequest) (*model.UserResponse, error) {
    var userResponse model.UserResponse
    cacheKey := "current-user:" + request.Email
    if err := c.Cache.GetAndUnmarshal(ctx, cacheKey, &userResponse); err == nil {
        return &userResponse, nil
    }
    
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Failed to validate request body : %+v", err)
		return nil, model.ErrBadRequest
	}
	user := new(entity.User)

	err = c.UserRepository.FindByEmail(tx, user, request.Email)
	if err != nil {
		c.Log.Warnf("Failed to find user : %+v", err)
		return nil, model.ErrInternalServer
	}

	err = tx.Commit().Error
	if err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, model.ErrInternalServer
	}

    userResponse = *converter.UserToResponse(user)
    userResponseJSON, _ := json.Marshal(userResponse)
    c.Cache.Set(ctx, cacheKey, userResponseJSON, 1*time.Minute)

	return &userResponse, nil
}