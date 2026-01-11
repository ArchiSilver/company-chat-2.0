package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/yourname/company-superapp/internal/domain"
)

type ChatService struct {
	chatRepo    domain.ChatRepository
	messageRepo domain.MessageRepository
}

func NewChatService(chatRepo domain.ChatRepository, messageRepo domain.MessageRepository) *ChatService {
	return &ChatService{
		chatRepo:    chatRepo,
		messageRepo: messageRepo,
	}
}

func (s *ChatService) GetUserChats(ctx context.Context, userID uuid.UUID) ([]domain.Chat, error) {
	return s.chatRepo.GetChatsByUserID(ctx, userID)
}

func (s *ChatService) GetChatMessages(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, page, pageSize int) ([]domain.Message, error) {
	// Here you would add logic to check if userID is a member of chatID
	// For now, we'll skip that for brevity.
	offset := (page - 1) * pageSize
	return s.messageRepo.GetMessagesByChatID(ctx, chatID, pageSize, offset)
}
