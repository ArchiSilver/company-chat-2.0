package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourname/company-superapp/internal/service"
)

type ReportHandler struct {
	reportService *service.ReportService
}

func NewReportHandler(reportService *service.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

func (h *ReportHandler) RegisterRoutes(rg *gin.RouterGroup) {
	reports := rg.Group("/reports")
	reports.Use(AuthMiddleware())
	{
		reports.GET("/tasks", h.GenerateTasksReport)
	}
}

// GenerateTasksReport generates a PDF report of tasks
// GET /api/v1/reports/tasks?from=2024-01-01&to=2024-01-31
func (h *ReportHandler) GenerateTasksReport(c *gin.Context) {
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

	// Parse date parameters
	fromStr := c.Query("from")
	toStr := c.Query("to")

	if fromStr == "" || toStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from and to date parameters are required"})
		return
	}

	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from date format, use YYYY-MM-DD"})
		return
	}

	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to date format, use YYYY-MM-DD"})
		return
	}

	// Set to end of day for the 'to' date
	to = to.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	// Generate PDF
	pdfBytes, err := h.reportService.GenerateTasksReport(c.Request.Context(), userID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate report"})
		return
	}

	// Set headers for PDF download
	filename := fmt.Sprintf("tasks_report_%s_%s.pdf", fromStr, toStr)
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))

	// Send PDF bytes
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}
