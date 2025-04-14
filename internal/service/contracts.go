package service

import (
	"context"
	"errors"
	"github.com/igorgrichanov/toDoList/internal/models"
)

var ErrInternal = errors.New("internal server error")
var ErrNotFound = errors.New("not found")
var ErrConflict = errors.New("already exists")

type TasksService interface {
	CreateTask(ctx context.Context, task *models.Task) (*models.Task, error)
	ListTasks(ctx context.Context) ([]*models.Task, error)
	UpdateTask(ctx context.Context, task *models.Task) error
	DeleteTask(ctx context.Context, id int) (*models.Task, error)
}
