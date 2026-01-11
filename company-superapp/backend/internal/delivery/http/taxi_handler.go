package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourname/company-superapp/internal/service"
)

type TaxiHandler struct {
	taxiService *service.TaxiService
}

func NewTaxiHandler(taxiService *service.TaxiService) *TaxiHandler {
	return &TaxiHandler{taxiService: taxiService}
}

func (h *TaxiHandler) RegisterRoutes(rg *gin.RouterGroup) {
	taxi := rg.Group("/taxi")
	taxi.Use(AuthMiddleware())
	{
		taxi.POST("/generate-upload-url", h.GenerateUploadURL)
		taxi.POST("/confirm-upload", h.ConfirmUpload)
		taxi.GET("/requests", h.GetUserRequests)
	}
}

func (h *TaxiHandler) GenerateUploadURL(c *gin.Context) {
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

	var req service.GenerateUploadURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.ContentType == "" {
		req.ContentType = "image/jpeg"
	}

	response, err := h.taxiService.GenerateUploadURL(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate upload URL"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *TaxiHandler) ConfirmUpload(c *gin.Context) {
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

	var req service.ConfirmUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	taxiRequest, err := h.taxiService.ConfirmUpload(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to confirm upload"})
		return
	}

	c.JSON(http.StatusCreated, taxiRequest)
}

func (h *TaxiHandler) GetUserRequests(c *gin.Context) {
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

	requests, err := h.taxiService.GetUserRequestsWithURLs(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get requests"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"requests": requests})
}
