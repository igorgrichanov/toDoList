package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/igorgrichanov/toDoList/internal/config"
	"github.com/igorgrichanov/toDoList/internal/controller/http/middleware/request_id"
	"github.com/igorgrichanov/toDoList/internal/models"
	"github.com/igorgrichanov/toDoList/internal/repository/tasks"
	"github.com/igorgrichanov/toDoList/internal/service/tasksService"
	"github.com/igorgrichanov/toDoList/pkg/logger/sl"
	"github.com/igorgrichanov/toDoList/pkg/postgres"
	"log/slog"
	"math/rand"
	"os"
)

func main() {
	count := flag.Int("count", 10, "number of fake tasks to insert")
	flag.Parse()

	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)
	conf, err := config.New()
	if err != nil {
		log.Error("Error loading config", sl.Err(err))
		os.Exit(1)
	}

	// db
	db, err := postgres.New(&conf.DB, log)
	if err != nil {
		log.Error("Error connecting to database", sl.Err(err))
		os.Exit(1)
	}
	defer db.Close()

	// repository
	repo := tasks.NewTasksRepository(log, db)

	// service
	ts := tasksService.NewUseCase(log, repo)

	// check if database is empty
	var countTasks int
	err = db.Pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM tasks").Scan(&countTasks)
	if err != nil {
		log.Error("failed to check existing tasks", sl.Err(err))
		os.Exit(1)
	}

	if countTasks > 0 {
		log.Info("skipping pre-population: tasks table already has data", slog.Int("count", countTasks))
		return
	}

	statuses := []string{"new", "in_progress", "done"}
	inserted := 0
	for i := 0; i < *count; i++ {
		task := &models.Task{
			Title:       gofakeit.Name(),
			Description: gofakeit.Sentence(5),
			Status:      statuses[rand.Intn(len(statuses))],
		}
		ctx := context.WithValue(context.Background(), request_id.RequestIDKey, uuid.New().String())
		_, err = ts.CreateTask(ctx, task)
		if err != nil {
			log.Error("Error creating task", sl.Err(err))
			continue
		}
		inserted++
	}

	successMsg := fmt.Sprintf("%d tasks inserted successfully", *count)
	log.Info(successMsg)
}
