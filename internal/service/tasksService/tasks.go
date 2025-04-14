package tasksService

import (
	"context"
	"errors"
	"github.com/igorgrichanov/toDoList/internal/controller/http/middleware/request_id"
	"github.com/igorgrichanov/toDoList/internal/models"
	"github.com/igorgrichanov/toDoList/internal/repository"
	"github.com/igorgrichanov/toDoList/internal/service"
	"github.com/igorgrichanov/toDoList/pkg/logger/sl"
	"log/slog"
)

type TasksRepository interface {
	Create(ctx context.Context, task *models.Task) (*models.Task, error)
	List(ctx context.Context) ([]*models.Task, error)
	Update(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, id int) (*models.Task, error)
}

type UseCase struct {
	log  *slog.Logger
	repo TasksRepository
}

func NewUseCase(log *slog.Logger, repo TasksRepository) *UseCase {
	return &UseCase{
		log:  log,
		repo: repo,
	}
}

func (u *UseCase) CreateTask(ctx context.Context, task *models.Task) (*models.Task, error) {
	const op = "service.tasks.CreateTask"
	requestID := ctx.Value(request_id.RequestIDKey).(string)
	log := u.log.With(
		slog.String("op", op),
		slog.String("request_id", requestID),
	)
	res, err := u.repo.Create(ctx, task)
	if errors.Is(err, repository.ErrInvalidInput) {
		log.Error("invalid input", sl.Err(err))
		return nil, service.ErrInternal
	} else if err != nil {
		log.Error("failed to create task", sl.Err(err))
		return nil, service.ErrInternal
	}

	log.Info("task created", slog.Any("id", res.ID))
	return res, nil
}

func (u *UseCase) ListTasks(ctx context.Context) ([]*models.Task, error) {
	const op = "service.tasks.ListTasks"
	requestID := ctx.Value(request_id.RequestIDKey).(string)
	log := u.log.With(
		slog.String("op", op),
		slog.String("request_id", requestID),
	)
	res, err := u.repo.List(ctx)
	if err != nil {
		log.Error("failed to get list of tasks", sl.Err(err))
		return nil, service.ErrInternal
	}
	log.Info("tasks list received")
	return res, nil
}

func (u *UseCase) UpdateTask(ctx context.Context, task *models.Task) error {
	const op = "service.tasks.UpdateTask"
	requestID := ctx.Value(request_id.RequestIDKey).(string)
	log := u.log.With(
		slog.String("op", op),
		slog.String("request_id", requestID),
	)
	err := u.repo.Update(ctx, task)
	if errors.Is(repository.ErrNotFound, err) {
		log.Error("Task not found", sl.Err(err))
		return service.ErrNotFound
	} else if errors.Is(repository.ErrConcurrentUpdate, err) {
		log.Error("concurrent update", sl.Err(err))
		return service.ErrConflict
	} else if errors.Is(repository.ErrInvalidInput, err) {
		log.Error("invalid input", sl.Err(err))
		return service.ErrInternal
	} else if err != nil {
		log.Error("failed to update task", sl.Err(err))
		return service.ErrInternal
	}
	log.Info("task updated", slog.Any("id", task.ID))
	return nil
}

func (u *UseCase) DeleteTask(ctx context.Context, id int) (*models.Task, error) {
	const op = "service.tasks.DeleteTask"
	requestID := ctx.Value(request_id.RequestIDKey).(string)
	log := u.log.With(
		slog.String("op", op),
		slog.String("request_id", requestID),
	)
	res, err := u.repo.Delete(ctx, id)
	if errors.Is(repository.ErrNotFound, err) {
		log.Error("task not found", sl.Err(err))
		return nil, service.ErrNotFound
	} else if err != nil {
		log.Error("failed to delete task", sl.Err(err))
		return nil, service.ErrInternal
	}
	log.Info("task deleted", slog.Any("id", id))
	return res, nil
}
