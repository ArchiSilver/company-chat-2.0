package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourname/company-superapp/internal/service"
)

type SearchHandler struct {
	searchService *service.GlobalSearchService
}

func NewSearchHandler(searchService *service.GlobalSearchService) *SearchHandler {
	return &SearchHandler{searchService: searchService}
}

func (h *SearchHandler) RegisterRoutes(rg *gin.RouterGroup) {
	search := rg.Group("/search")
	search.Use(AuthMiddleware())
	{
		search.GET("", h.Search)
	}
}

func (h *SearchHandler) Search(c *gin.Context) {
	query := c.Query("q")

	results, err := h.searchService.Search(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
		return
	}

	c.JSON(http.StatusOK, results)
}
