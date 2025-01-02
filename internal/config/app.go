package config

import (
	"github.com/abdisetiakawan/go-clean-arch/internal/delivery/http"
	"github.com/abdisetiakawan/go-clean-arch/internal/delivery/http/middleware"
	"github.com/abdisetiakawan/go-clean-arch/internal/delivery/http/route"
	"github.com/abdisetiakawan/go-clean-arch/internal/helper"
	"github.com/abdisetiakawan/go-clean-arch/internal/repository"
	"github.com/abdisetiakawan/go-clean-arch/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
    DB       *gorm.DB
    App      *fiber.App
    Log      *logrus.Logger
    Validate *validator.Validate
    Config   *viper.Viper
    Help     *helper.Helper
}

func Bootstrap(config *BootstrapConfig) {
    userRepository := repository.NewUserRepository(config.DB, config.Log)

    userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository, config.Help)

    userController := http.NewUserController(userUseCase, config.Log)

    authMiddleware := middleware.NewAuth(userUseCase)
    routeConfig := route.RouteConfig{
        App:            config.App,
        UserController: userController,
        AuthMiddleware: authMiddleware,
    }
    routeConfig.Setup()
}