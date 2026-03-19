package repository

import (
	"errors"
	"sync"
	"time"

	"github.com/WarunaUdara/jenkins-e2e-cicd-pipeline-google-kubernetes-engine/models"
)

var (
	// ErrTodoNotFound is returned when a todo is not found
	ErrTodoNotFound = errors.New("todo not found")
)

// TodoRepository provides thread-safe access to todo data
type TodoRepository struct {
	mu     sync.RWMutex
	todos  map[int]*models.Todo
	nextID int
}

// NewTodoRepository creates a new todo repository
func NewTodoRepository() *TodoRepository {
	return &TodoRepository{
		todos:  make(map[int]*models.Todo),
		nextID: 1,
	}
}

// GetAll returns all todos
func (r *TodoRepository) GetAll() []*models.Todo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	todos := make([]*models.Todo, 0, len(r.todos))
	for _, todo := range r.todos {
		todos = append(todos, todo)
	}
	return todos
}

// GetByID returns a todo by ID
func (r *TodoRepository) GetByID(id int) (*models.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	todo, exists := r.todos[id]
	if !exists {
		return nil, ErrTodoNotFound
	}
	return todo, nil
}

// Create adds a new todo
func (r *TodoRepository) Create(title string) *models.Todo {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	todo := &models.Todo{
		ID:          r.nextID,
		Title:       title,
		IsCompleted: false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	r.todos[r.nextID] = todo
	r.nextID++

	return todo
}

// Update modifies an existing todo
func (r *TodoRepository) Update(id int, title *string, isCompleted *bool) (*models.Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	todo, exists := r.todos[id]
	if !exists {
		return nil, ErrTodoNotFound
	}

	if title != nil {
		todo.Title = *title
	}
	if isCompleted != nil {
		todo.IsCompleted = *isCompleted
	}
	todo.UpdatedAt = time.Now()

	return todo, nil
}

// Delete removes a todo by ID
func (r *TodoRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.todos[id]; !exists {
		return ErrTodoNotFound
	}

	delete(r.todos, id)
	return nil
}

// Count returns the total number of todos
func (r *TodoRepository) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.todos)
}
