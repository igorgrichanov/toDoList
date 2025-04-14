package tasks

import (
	"context"
	"errors"
	"fmt"
	"github.com/igorgrichanov/toDoList/internal/models"
	"github.com/igorgrichanov/toDoList/internal/repository"
	"github.com/igorgrichanov/toDoList/pkg/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"log/slog"
	"time"
)

const (
	tasksTable = "tasks"
)

type Tasks struct {
	log *slog.Logger
	db  *postgres.Postgres
}

func NewTasksRepository(log *slog.Logger, db *postgres.Postgres) *Tasks {
	return &Tasks{log: log, db: db}
}

func (t *Tasks) Create(ctx context.Context, task *models.Task) (*models.Task, error) {
	createdAt := time.Now().UTC()
	updatedAt := createdAt
	sql, args, err := t.db.Builder.Insert(tasksTable).
		Columns("title", "description", "status", "created_at", "updated_at").
		Values(task.Title, task.Description, task.Status, createdAt, updatedAt).
		Suffix("RETURNING \"id\"").ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", repository.ErrBuildingSql, err)
	}
	var insertedID int
	err = t.db.Pool.QueryRow(ctx, sql, args...).Scan(&insertedID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23502": // not null violation
				return nil, fmt.Errorf("%w: missing required field", repository.ErrInvalidInput)
			case "23514": // check constraint
				return nil, fmt.Errorf("%w: invalid value in field with CHECK", repository.ErrInvalidInput)
			}
		}
		return nil, fmt.Errorf("%w: %s", repository.ErrExecutingSql, err)
	}
	task.ID = insertedID
	return task, nil
}

func (t *Tasks) List(ctx context.Context) ([]*models.Task, error) {
	// if table grows, SELECT "*" may cause errors when retrieving data
	sql, args, err := t.db.Builder.Select("id", "title", "description", "status", "created_at", "updated_at").
		From(tasksTable).OrderBy("id").ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", repository.ErrBuildingSql, err)
	}
	rows, err := t.db.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", repository.ErrExecutingSql, err)
	}
	defer rows.Close()

	tasks := make([]*models.Task, 0, 20) // Можно отправить select count(1) для cap, но это лишнее обращение в БД.
	// Для больших таблиц используют пагинацию
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", repository.ErrRetrievingData, err)
		}
		tasks = append(tasks, &task)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %s", repository.ErrRetrievingData, err)
	}
	return tasks, nil
}

func (t *Tasks) Update(ctx context.Context, task *models.Task) error {
	// get last update time to avoid data races
	var lastUpdatedAt time.Time
	sql, args, err := t.db.Builder.Select("updated_at").From(tasksTable).
		Where("id = ?", task.ID).ToSql()
	if err != nil {
		return fmt.Errorf("%w: %s", repository.ErrBuildingSql, err)
	}
	err = t.db.Pool.QueryRow(ctx, sql, args...).Scan(&lastUpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrNotFound
		}
		return fmt.Errorf("%w: %s", repository.ErrExecutingSql, err)
	}

	updatedAt := time.Now().UTC()
	sql, args, err = t.db.Builder.Update(tasksTable).
		Set("title", task.Title).
		Set("description", task.Description).
		Set("status", task.Status).
		Set("updated_at", updatedAt).
		Where("id = ?", task.ID).
		Where("updated_at = ?", lastUpdatedAt).ToSql()
	if err != nil {
		return fmt.Errorf("%w: %s", repository.ErrBuildingSql, err)
	}
	res, err := t.db.Pool.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23502": // not null violation
				return fmt.Errorf("%w: missing required field", repository.ErrInvalidInput)
			case "23514": // check constraint
				return fmt.Errorf("%w: invalid value in field with CHECK", repository.ErrInvalidInput)
			}
		}
		return fmt.Errorf("%w: %s", repository.ErrExecutingSql, err)
	}
	if res.RowsAffected() == 0 {
		return repository.ErrConcurrentUpdate
	}
	return nil
}

func (t *Tasks) Delete(ctx context.Context, id int) (*models.Task, error) {
	sql, args, err := t.db.Builder.Delete(tasksTable).Where("id = ?", id).
		Suffix("RETURNING id, title, description, status, created_at, updated_at").ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", repository.ErrBuildingSql, err)
	}
	var task models.Task
	err = t.db.Pool.QueryRow(ctx, sql, args...).
		Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %s", repository.ErrExecutingSql, err)
	}
	return &task, nil
}
