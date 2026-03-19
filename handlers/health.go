package handlers

import (
	"net/http"
	"time"

	"github.com/WarunaUdara/jenkins-e2e-cicd-pipeline-google-kubernetes-engine/repository"
	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	repo      *repository.TodoRepository
	startTime time.Time
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(repo *repository.TodoRepository) *HealthHandler {
	return &HealthHandler{
		repo:      repo,
		startTime: time.Now(),
	}
}

// Health returns the health status of the application
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "todo-app-go",
		"uptime":  time.Since(h.startTime).String(),
	})
}

// Ready returns the readiness status of the application
func (h *HealthHandler) Ready(c *gin.Context) {
	// Check if the application is ready to serve requests
	// For this simple app, we just check if the repository is initialized
	if h.repo == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not ready",
			"reason": "repository not initialized",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "ready",
		"service":    "todo-app-go",
		"todo_count": h.repo.Count(),
	})
}
