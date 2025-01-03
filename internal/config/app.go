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
    userRepository := repository.NewUserRepository(config.Log)
    userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository, config.Help)
    userController := http.NewUserController(userUseCase, config.Log)

    taskRepository := repository.NewTaskRepository(config.Log)
    taskUseCase := usecase.NewTaskUseCase(config.DB, config.Log, config.Validate, taskRepository)
    taskController := http.NewTaskController(taskUseCase, config.Log)

    tagRepository := repository.NewTagRepository(config.Log)
    tagUseCase := usecase.NewTagUseCase(config.DB, config.Log, config.Validate, tagRepository)
    tagController := http.NewTagsController(tagUseCase, config.Log)

    taskTagRepository := repository.NewtaskTagRepository(config.Log)
    taskTagUseCase := usecase.NewTaskTagUseCase(config.DB, config.Log, config.Validate, taskTagRepository)
    taskTagController := http.NewTaskTagController(taskTagUseCase, config.Log)
    
    authMiddleware := middleware.NewAuth(userUseCase, config.Config)
    routeConfig := route.RouteConfig{
        App:            config.App,
        UserController: userController,
        TaskController: taskController,
        TagsController: tagController,
        TaskTagController: taskTagController,
        AuthMiddleware: authMiddleware,
    }
    routeConfig.Setup()
}