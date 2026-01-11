package service

import (
	"context"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/yourname/company-superapp/internal/domain"
	"github.com/yourname/company-superapp/internal/pkg/fcm"
)

type NotificationService struct {
	pushTokenRepo domain.PushTokenRepository
	fcmClient     *fcm.FCMClient
}

func NewNotificationService(pushTokenRepo domain.PushTokenRepository, fcmClient *fcm.FCMClient) *NotificationService {
	return &NotificationService{
		pushTokenRepo: pushTokenRepo,
		fcmClient:     fcmClient,
	}
}

type RegisterTokenRequest struct {
	Token      string  `json:"token" binding:"required"`
	DeviceInfo *string `json:"device_info,omitempty"`
}

func (s *NotificationService) RegisterToken(ctx context.Context, userID uuid.UUID, req RegisterTokenRequest) error {
	pushToken := &domain.PushToken{
		ID:         uuid.New(),
		UserID:     userID,
		Token:      req.Token,
		DeviceInfo: req.DeviceInfo,
	}

	return s.pushTokenRepo.Create(ctx, pushToken)
}

func (s *NotificationService) UnregisterToken(ctx context.Context, token string) error {
	return s.pushTokenRepo.Delete(ctx, token)
}

func (s *NotificationService) SendToUser(ctx context.Context, userID uuid.UUID, title string, body string, data map[string]string) error {
	tokens, err := s.pushTokenRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if len(tokens) == 0 {
		log.Printf("No push tokens found for user %s", userID)
		return nil
	}

	for _, token := range tokens {
		var sendErr error
		
		// Проверка, является ли токен Expo токеном
		if strings.HasPrefix(token.Token, "ExponentPushToken") {
			sendErr = s.fcmClient.SendToExpoToken(ctx, token.Token, title, body, data)
		} else {
			sendErr = s.fcmClient.SendNotification(ctx, token.Token, title, body, data)
		}

		if sendErr != nil {
			log.Printf("Failed to send notification to token %s: %v", token.Token, sendErr)
			// Продолжаем отправку на другие токены
		}
	}

	return nil
}

func (s *NotificationService) GetUserTokens(ctx context.Context, userID uuid.UUID) ([]domain.PushToken, error) {
	return s.pushTokenRepo.GetByUserID(ctx, userID)
}
