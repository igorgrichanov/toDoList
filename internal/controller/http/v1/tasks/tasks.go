package tasks

import (
	"github.com/gofiber/fiber/v2"
	"log/slog"
)

type Tasker interface {
	Create(c *fiber.Ctx) error
	List(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type TaskController struct {
	log *slog.Logger
}
