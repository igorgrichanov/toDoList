package response

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strings"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func ValidationError(errs validator.ValidationErrors) string {
	var errMsgs []string

	for _, err := range errs {
		field := err.Field()
		switch err.Tag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field '%s' is required", field))
		case "oneof":
			errMsgs = append(errMsgs, fmt.Sprintf("field '%s' must be one of [%s]", field, err.Param()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", field))
		}
	}

	return strings.Join(errMsgs, ", ")
}

func ErrorBadRequest(c *fiber.Ctx, err string) error {
	resp := Response{
		Success: false,
		Message: err,
	}
	return c.Status(http.StatusBadRequest).JSON(resp)
}

func ErrorNotFound(c *fiber.Ctx) error {
	resp := Response{
		Success: false,
		Message: http.StatusText(http.StatusNotFound),
	}
	return c.Status(http.StatusNotFound).JSON(resp)
}

func ErrorInternal(c *fiber.Ctx) error {
	resp := Response{
		Success: false,
		Message: http.StatusText(http.StatusInternalServerError),
	}
	return c.Status(http.StatusInternalServerError).JSON(resp)
}

func ErrorConflict(c *fiber.Ctx, err string) error {
	resp := Response{
		Success: false,
		Message: http.StatusText(http.StatusConflict),
	}
	return c.Status(http.StatusConflict).JSON(resp)
}
