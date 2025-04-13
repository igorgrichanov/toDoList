package models

import "time"

type Task struct {
	ID          int       `json:"id,omitempty"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}
