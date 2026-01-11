package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourname/company-superapp/internal/delivery/websocket"
	"github.com/yourname/company-superapp/internal/service"
)

type ChatHandler struct {
	service *service.ChatService
	hub     *websocket.Hub
}

func NewChatHandler(service *service.ChatService, hub *websocket.Hub) *ChatHandler {
	return &ChatHandler{
		service: service,
		hub:     hub,
	}
}

func (h *ChatHandler) RegisterRoutes(router *gin.RouterGroup) {
	chats := router.Group("/chats")
	{
		chats.GET("", h.getChats)
		chats.GET("/:id/messages", h.getMessages)
	}
	ws := router.Group("/ws")
	{
		ws.GET("/connect", h.handleWebSocket)
	}
}

func (h *ChatHandler) getChats(c *gin.Context) {
	// For now, we'll use a placeholder user ID. In a real app, this would come from the JWT.
	userID, _ := uuid.NewRandom() 
	chats, err := h.service.GetUserChats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get chats"})
		return
	}
	c.JSON(http.StatusOK, chats)
}

func (h *ChatHandler) getMessages(c *gin.Context) {
	chatID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat id"})
		return
	}
	// Placeholder user ID
	userID, _ := uuid.NewRandom()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "50"))

	messages, err := h.service.GetChatMessages(c.Request.Context(), chatID, userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get messages"})
		return
	}
	c.JSON(http.StatusOK, messages)
}

func (h *ChatHandler) handleWebSocket(c *gin.Context) {
	// In a real app, userID would come from JWT middleware
	userIDStr := c.Query("user_id")
	chatIDStr := c.Query("chat_id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat_id"})
		return
	}

	websocket.ServeWs(h.hub, c.Writer, c.Request, userID, chatID)
}
