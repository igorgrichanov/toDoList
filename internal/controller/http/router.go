package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/igorgrichanov/toDoList/internal/config"
	"github.com/igorgrichanov/toDoList/internal/controller"
	"github.com/igorgrichanov/toDoList/internal/controller/http/middleware/logger"
	"github.com/igorgrichanov/toDoList/internal/controller/http/middleware/request_id"
	swagger "github.com/swaggo/fiber-swagger"
	"log/slog"
)

//	@Title			To-do list API
//	@Version		1.0
//	@Description	To-do list API allows to perform operations with the list of tasks.

//	@Host		localhost:8080
//	@BasePath	/
//	@Schemes	http
//	@Accept		json
//	@Produce	json

// @Tag.name			tasks
// @Tag.description	operations with the list of tasks
func NewRouter(log *slog.Logger, cfg *config.Server, ctrl *controller.Controllers) *fiber.App {
	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	})
	app.Use(recover.New())
	app.Use(request_id.NewRequestIDMiddleware())
	app.Use(logger.NewLoggerMiddleware(log))

	tasks := app.Group("/tasks")
	tasks.Post("/", ctrl.Tasks.Create)
	tasks.Get("/", ctrl.Tasks.List)
	tasks.Put("/:id", ctrl.Tasks.Update)
	tasks.Delete("/:id", ctrl.Tasks.Delete)

	sw := app.Group("/swagger")
	sw.Use(func(c *fiber.Ctx) error {
		c.Set("Cache-Control", "no-store, no-cache, must-revalidate")
		return c.Next()
	})
	sw.Get("/*", swagger.WrapHandler)
	return app
}
