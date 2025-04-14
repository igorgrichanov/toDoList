package app

import (
	"context"
	"github.com/go-playground/validator/v10"
	_ "github.com/igorgrichanov/toDoList/docs"
	"github.com/igorgrichanov/toDoList/internal/config"
	"github.com/igorgrichanov/toDoList/internal/controller"
	httpRouter "github.com/igorgrichanov/toDoList/internal/controller/http"
	tasksController "github.com/igorgrichanov/toDoList/internal/controller/http/v1/tasks"
	"github.com/igorgrichanov/toDoList/internal/repository/tasks"
	"github.com/igorgrichanov/toDoList/internal/service/tasksService"
	"github.com/igorgrichanov/toDoList/pkg/logger/sl"
	"github.com/igorgrichanov/toDoList/pkg/postgres"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func Run(log *slog.Logger, conf *config.Config) {
	// db
	db, err := postgres.New(&conf.DB, log)
	if err != nil {
		log.Error("Error connecting to database", sl.Err(err))
		os.Exit(1)
	}
	defer db.Close()

	// infrastructure
	repo := tasks.NewTasksRepository(log, db)
	validate := validator.New()

	// service
	uc := tasksService.NewUseCase(log, repo)

	// controller
	tasksCtrl := tasksController.NewTaskController(log, uc, validate)
	ctrl := controller.New(tasksCtrl)

	// router
	app := httpRouter.NewRouter(log, &conf.Server, ctrl)

	addr := conf.Server.Host + ":" + conf.Server.Port
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info("starting server at http://" + addr)
		if err := app.Listen(addr); err != nil {
			log.Error(err.Error())
		}
	}()

	<-done
	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), conf.Server.ShutdownTimeout)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Error("error while shutting down server", sl.Err(err))
	} else {
		log.Info("server shutdown complete")
	}
}
