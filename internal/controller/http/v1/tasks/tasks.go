package tasks

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/igorgrichanov/toDoList/internal/controller/http/middleware/request_id"
	"github.com/igorgrichanov/toDoList/internal/models"
	"github.com/igorgrichanov/toDoList/internal/service"
	"github.com/igorgrichanov/toDoList/pkg/api/response"
	"github.com/igorgrichanov/toDoList/pkg/logger/sl"
	"log/slog"
	"strconv"
	"strings"
)

type Tasker interface {
	Create(c *fiber.Ctx) error
	List(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type TaskController struct {
	log       *slog.Logger
	uc        service.Tasks
	validator *validator.Validate
}

func NewTaskController(log *slog.Logger, uc service.Tasks, v *validator.Validate) *TaskController {
	return &TaskController{log: log, uc: uc, validator: v}
}

type CreateRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty" validate:"omitempty,oneof=new in_progress done"`
}

// @Summary	Create a new task
// @Tags		tasks
// @Param		Task	body		CreateRequest	true	"Specify task title. Description and status are optional"
// @Success	201		{object}	models.Task
// @Failure	400		{object}	response.Response	"invalid request body"
// @Failure	500		{object}	response.Response	"internal server error"
// @Router		/tasks [post]
func (tc *TaskController) Create(c *fiber.Ctx) error {
	const op = "controller.tasks.Create"
	requestID := c.UserContext().Value(request_id.RequestIDKey).(string)
	log := tc.log.With(
		slog.String("op", op),
		slog.String("request_id", requestID),
	)

	req := &CreateRequest{}
	if err := c.BodyParser(req); err != nil {
		log.Error("failed to parse request body", sl.Err(err))
		return response.ErrorBadRequest(c, "invalid request body")
	}
	req.Status = strings.ToLower(req.Status)
	if err := tc.validator.Struct(req); err != nil {
		log.Error("validation failed", sl.Err(err))
		var validateErr validator.ValidationErrors
		if errors.As(err, &validateErr) {
			resp := response.ValidationError(validateErr)
			return response.ErrorBadRequest(c, resp)
		}
		return response.ErrorBadRequest(c, "request body validation failed")

	}
	log.Info("request received", slog.Any("data", req))

	task, err := tc.uc.CreateTask(c.UserContext(), &models.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
	})
	if err != nil {
		log.Error("failed to create task", sl.Err(err))
		return response.ErrorInternal(c)
	}
	log.Info("task created", slog.Any("data", task))

	return c.Status(fiber.StatusCreated).JSON(task)
}

// @Summary	Get list of existing tasks
// @Tags		tasks
// @Success	200	{object}	[]models.Task
// @Failure	500	{object}	response.Response	"internal server error"
// @Router		/tasks [get]
func (tc *TaskController) List(c *fiber.Ctx) error {
	const op = "controller.tasks.List"
	requestID := c.UserContext().Value(request_id.RequestIDKey).(string)
	log := tc.log.With(
		slog.String("op", op),
		slog.String("request_id", requestID),
	)
	log.Info("request received")

	tasks, err := tc.uc.ListTasks(c.UserContext())
	if err != nil {
		log.Error("failed to list tasks", sl.Err(err))
		return response.ErrorInternal(c)
	}

	log.Info("tasks received", slog.Any("data", tasks))
	return c.Status(fiber.StatusOK).JSON(tasks)
}

type UpdateRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status" validate:"oneof=new in_progress done"`
}

// @Summary	Update task
// @Tags		tasks
// @Param		id		path	int				true	"Task ID"
// @Param		Task	body	UpdateRequest	true	"Specify fields to update"
// @Success	204
// @Failure	400	{object}	response.Response	"invalid request body or task ID"
// @Failure	404	{object}	response.Response	"task not found"
// @Failure	409	{object}	response.Response	"task has already been updated, try again"
// @Failure	500	{object}	response.Response	"internal server error"
// @Router		/tasks/{id} [put]
func (tc *TaskController) Update(c *fiber.Ctx) error {
	const op = "controller.tasks.Update"
	requestID := c.UserContext().Value(request_id.RequestIDKey).(string)
	log := tc.log.With(
		slog.String("op", op),
		slog.String("request_id", requestID),
	)
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Error("failed to parse id param", sl.Err(err))
		return response.ErrorBadRequest(c, "invalid task ID")
	}
	if id < 0 {
		log.Error("invalid id param", slog.Int("id", id))
		return response.ErrorBadRequest(c, "invalid task ID")
	}
	req := &UpdateRequest{}
	if err := c.BodyParser(req); err != nil {
		log.Error("failed to parse request body", sl.Err(err))
		return response.ErrorBadRequest(c, "invalid request body")
	}
	req.Status = strings.ToLower(req.Status)
	if err = tc.validator.Struct(req); err != nil {
		log.Error("validation failed", sl.Err(err))
		var validateErr validator.ValidationErrors
		if errors.As(err, &validateErr) {
			resp := response.ValidationError(validateErr)
			return response.ErrorBadRequest(c, resp)
		}
		return response.ErrorBadRequest(c, "request body validation failed")
	}
	log.Info("request received", slog.Any("data", req))

	err = tc.uc.UpdateTask(c.UserContext(), &models.Task{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
	})
	if errors.Is(err, service.ErrNotFound) {
		log.Error("task not found", sl.Err(err))
		return response.ErrorNotFound(c)
	} else if errors.Is(err, service.ErrConflict) {
		log.Error("task has already been updated", sl.Err(err))
		return response.ErrorConflict(c, "task has already been updated, try again")
	} else if err != nil {
		log.Error("failed to update task", sl.Err(err))
		return response.ErrorInternal(c)
	}
	log.Info("task updated", slog.Any("data", req))

	return c.SendStatus(fiber.StatusNoContent)
}

// @Summary	Delete task
// @Tags		tasks
// @Param		id	path		int	true	"Task ID"
// @Success	200	{object}	models.Task
// @Failure	400	{object}	response.Response	"invalid task ID"
// @Failure	404	{object}	response.Response	"task not found"
// @Failure	500	{object}	response.Response	"internal server error"
// @Router		/tasks/{id} [delete]
func (tc *TaskController) Delete(c *fiber.Ctx) error {
	const op = "controller.tasks.Delete"
	requestID := c.UserContext().Value(request_id.RequestIDKey).(string)
	log := tc.log.With(
		slog.String("op", op),
		slog.String("request_id", requestID),
	)
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Error("failed to parse id param", slog.String("error", err.Error()))
		return response.ErrorBadRequest(c, "invalid task ID")
	}
	if id < 0 {
		log.Error("invalid id param", slog.Int("id", id))
		return response.ErrorBadRequest(c, "invalid task ID")
	}

	task, err := tc.uc.DeleteTask(c.UserContext(), id)
	if errors.Is(err, service.ErrNotFound) {
		log.Error("task not found", sl.Err(err))
		return response.ErrorNotFound(c)
	} else if err != nil {
		log.Error("failed to delete task", sl.Err(err))
		return response.ErrorInternal(c)
	}
	log.Info("task deleted", slog.Any("data", task))

	return c.Status(fiber.StatusOK).JSON(task)
}
