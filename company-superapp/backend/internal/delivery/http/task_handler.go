package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourname/company-superapp/internal/domain"
	"github.com/yourname/company-superapp/internal/service"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) RegisterRoutes(router *gin.RouterGroup) {
	tasks := router.Group("/tasks")
	{
		tasks.POST("", h.createTask)
		tasks.POST("/from-message", h.createFromMessage)
		tasks.GET("", h.getTasks)
		tasks.GET("/:id", h.getTask)
		tasks.PUT("/:id", h.updateTask)
		tasks.PUT("/:id/status", h.updateStatus)
		tasks.DELETE("/:id", h.deleteTask)
	}
}

func (h *TaskHandler) createTask(c *gin.Context) {
	var input service.CreateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Placeholder user ID - in real app, this comes from JWT middleware
	creatorID, _ := uuid.NewRandom()

	task, err := h.service.Create(c.Request.Context(), creatorID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) createFromMessage(c *gin.Context) {
	var input service.CreateFromMessageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	creatorID, _ := uuid.NewRandom()

	task, err := h.service.CreateFromMessage(c.Request.Context(), creatorID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task from message"})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) getTasks(c *gin.Context) {
	var assigneeID *uuid.UUID
	var status *domain.TaskStatus

	if assigneeStr := c.Query("assignee_id"); assigneeStr != "" {
		if id, err := uuid.Parse(assigneeStr); err == nil {
			assigneeID = &id
		}
	}

	if statusStr := c.Query("status"); statusStr != "" {
		s := domain.TaskStatus(statusStr)
		status = &s
	}

	tasks, err := h.service.GetAll(c.Request.Context(), assigneeID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) getTask(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	task, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) updateTask(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	var input service.UpdateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.service.Update(c.Request.Context(), id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) updateStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	var input service.UpdateStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.UpdateStatus(c.Request.Context(), id, input.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *TaskHandler) deleteTask(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	err = h.service.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
