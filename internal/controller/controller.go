package controller

import "github.com/igorgrichanov/toDoList/internal/controller/http/v1/tasks"

type Controllers struct {
	Tasks tasks.Tasker
}

func New(taskController tasks.Tasker) *Controllers {
	return &Controllers{
		Tasks: taskController,
	}
}
