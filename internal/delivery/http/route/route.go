package route

import (
	"github.com/abdisetiakawan/go-clean-arch/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App               *fiber.App
	UserController    *http.UserController
	TaskController 	*http.TaskController
	TagsController *http.TagsController
	TaskTagController *http.TaskTagController
	AuthMiddleware    fiber.Handler
}

func (c *RouteConfig) Setup() {
	c.SetupAuthRoute()
	c.SetupUserRoute()
}

func (c *RouteConfig) SetupAuthRoute() {
	c.App.Post("/api/users", c.UserController.Register)
	c.App.Post("/api/users/_login", c.UserController.Login)
}

func (c *RouteConfig) SetupUserRoute() {
	c.App.Use(c.AuthMiddleware)
	c.App.Patch("/api/users/_current", c.UserController.Update)
	c.App.Get("/api/users/_current", c.UserController.Current)

	c.App.Get("/api/tasks", c.TaskController.List)
	c.App.Post("/api/tasks", c.TaskController.Create)
	c.App.Put("/api/tasks/:taskId", c.TaskController.Update)
	c.App.Get("/api/tasks/:taskId", c.TaskController.Get)
	c.App.Delete("/api/tasks/:taskId", c.TaskController.Delete)

	c.App.Post("/api/tags", c.TagsController.Create)
	c.App.Get("/api/tags", c.TagsController.List)
	c.App.Get("/api/tags/:tagId", c.TagsController.Get)
	c.App.Put("/api/tags/:tagId", c.TagsController.Update)
	c.App.Delete("/api/tags/:tagId", c.TagsController.Delete)

	c.App.Post("/api/tasks/:taskId/tags", c.TaskTagController.Create)
	c.App.Get("/api/taskswithtags", c.TaskTagController.List)
	c.App.Get("/api/tags/:tagId/tasks", c.TaskTagController.ListByTagId)
	c.App.Delete("/api/tasks/:taskId/tags/:tagId", c.TaskTagController.Delete)
}
