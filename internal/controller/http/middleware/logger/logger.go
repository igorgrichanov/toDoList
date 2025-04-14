package logger

import (
	"github.com/gofiber/fiber/v2"
	"github.com/igorgrichanov/toDoList/internal/controller/http/middleware/request_id"
	"log/slog"
	"time"
)

func NewLoggerMiddleware(log *slog.Logger) fiber.Handler {
	log = log.With(
		slog.String("component", "middleware/logger"),
	)
	log.Info("logger middleware enabled")
	return func(c *fiber.Ctx) error {
		requestID, ok := c.UserContext().Value(request_id.RequestIDKey).(string)
		if !ok {
			log.Error("missing request id", slog.String("request_id", requestID))
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		entry := log.With(
			slog.String("method", c.Method()),
			slog.String("path", c.OriginalURL()),
			slog.String("remote_addr", c.IP()),
			slog.String("user_agent", c.Get("User-Agent")),
			slog.String("request_id", requestID),
		)

		start := time.Now()
		defer func() {
			entry.Info("request completed",
				slog.Int("status", c.Response().StatusCode()),
				slog.Int("bytes", len(c.Response().Body())),
				slog.String("duration", time.Since(start).String()),
			)
		}()
		return c.Next()
	}
}
