package request_id

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Key to use when setting the request ID. Non-exported type is used to avoid collisions
type ctxKeyRequestID int

// RequestIDKey is the key that holds the unique request ID in a request context.
// TODO: use non-exported variable and FromContext() to retrieve value
const RequestIDKey ctxKeyRequestID = 0

// RequestIDHeader is the name of the HTTP Header which contains the request id.
// Exported so that it can be changed by developers
var RequestIDHeader = "X-Request-Id"

func NewRequestIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		reqID := c.Get(RequestIDHeader)
		if reqID == "" {
			reqID = uuid.New().String()
		}
		c.SetUserContext(context.WithValue(c.Context(), RequestIDKey, reqID))
		return c.Next()
	}
}
