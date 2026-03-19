package models

import "time"

// Todo represents a todo item in the application
type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title" binding:"required"`
	IsCompleted bool      `json:"is_completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateTodoRequest represents the request body for creating a todo
type CreateTodoRequest struct {
	Title string `json:"title" binding:"required,min=1,max=200"`
}

// UpdateTodoRequest represents the request body for updating a todo
type UpdateTodoRequest struct {
	Title       *string `json:"title,omitempty" binding:"omitempty,min=1,max=200"`
	IsCompleted *bool   `json:"is_completed,omitempty"`
}
