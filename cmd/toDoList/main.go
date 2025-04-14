package main

import (
	"github.com/igorgrichanov/toDoList/internal/app"
	"github.com/igorgrichanov/toDoList/internal/config"
	"github.com/igorgrichanov/toDoList/pkg/logger/sl"
	"log/slog"
	"os"
)

func main() {
	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)
	conf, err := config.New()
	if err != nil {
		log.Error("Error loading config", sl.Err(err))
		os.Exit(1)
	}

	app.Run(log, conf)
}
