package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/WarunaUdara/jenkins-e2e-cicd-pipeline-google-kubernetes-engine/models"
	"github.com/WarunaUdara/jenkins-e2e-cicd-pipeline-google-kubernetes-engine/repository"
	"github.com/gin-gonic/gin"
)

// TodoHandler handles todo-related requests
type TodoHandler struct {
	repo *repository.TodoRepository
}

// NewTodoHandler creates a new todo handler
func NewTodoHandler(repo *repository.TodoRepository) *TodoHandler {
	return &TodoHandler{repo: repo}
}

// ListTodos returns all todos
func (h *TodoHandler) ListTodos(c *gin.Context) {
	todos := h.repo.GetAll()
	c.JSON(http.StatusOK, gin.H{
		"todos": todos,
		"count": len(todos),
	})
}

// GetTodo returns a single todo by ID
func (h *TodoHandler) GetTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid todo ID",
		})
		return
	}

	todo, err := h.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrTodoNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "todo not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrieve todo",
		})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// CreateTodo creates a new todo
func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var req models.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	todo := h.repo.Create(req.Title)
	c.JSON(http.StatusCreated, todo)
}

// UpdateTodo updates an existing todo
func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid todo ID",
		})
		return
	}

	var req models.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	todo, err := h.repo.Update(id, req.Title, req.IsCompleted)
	if err != nil {
		if errors.Is(err, repository.ErrTodoNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "todo not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update todo",
		})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// DeleteTodo deletes a todo by ID
func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid todo ID",
		})
		return
	}

	if err := h.repo.Delete(id); err != nil {
		if errors.Is(err, repository.ErrTodoNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "todo not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete todo",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "todo deleted successfully",
	})
}
