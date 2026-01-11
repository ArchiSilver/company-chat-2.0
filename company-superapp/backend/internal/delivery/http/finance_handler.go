package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourname/company-superapp/internal/service"
)

type FinanceHandler struct {
	salaryService *service.SalaryService
}

func NewFinanceHandler(salaryService *service.SalaryService) *FinanceHandler {
	return &FinanceHandler{salaryService: salaryService}
}

func (h *FinanceHandler) RegisterRoutes(rg *gin.RouterGroup) {
	finance := rg.Group("/finance")
	finance.Use(AuthMiddleware())
	finance.Use(RBACMiddleware("admin", "manager")) // Only admin and manager can access finance
	{
		finance.GET("/salary", h.GetSalary)
		finance.PUT("/salary", h.UpdateSalary)
	}
}

func (h *FinanceHandler) GetSalary(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	salary, err := h.salaryService.GetSalary(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get salary"})
		return
	}

	c.JSON(http.StatusOK, salary)
}

type UpdateSalaryRequest struct {
	Amount float64 `json:"amount" binding:"required"`
}

func (h *FinanceHandler) UpdateSalary(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req UpdateSalaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.salaryService.UpdateSalary(c.Request.Context(), userID, req.Amount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update salary"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "salary updated successfully"})
}
