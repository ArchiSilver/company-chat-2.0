package fcm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type FCMClient struct {
	serverKey string
}

func NewFCMClient() (*FCMClient, error) {
	serverKey := os.Getenv("FCM_SERVER_KEY")
	if serverKey == "" {
		// Разрешаем запуск без FCM для разработки
		return &FCMClient{serverKey: ""}, nil
	}

	return &FCMClient{serverKey: serverKey}, nil
}

type FCMMessage struct {
	To           string            `json:"to,omitempty"`
	Notification *FCMNotification  `json:"notification,omitempty"`
	Data         map[string]string `json:"data,omitempty"`
}

type FCMNotification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Sound string `json:"sound,omitempty"`
	Badge string `json:"badge,omitempty"`
}

type FCMResponse struct {
	Success int `json:"success"`
	Failure int `json:"failure"`
}

func (c *FCMClient) SendNotification(ctx context.Context, token string, title string, body string, data map[string]string) error {
	if c.serverKey == "" {
		// Пропускаем отправку, если FCM не настроен
		fmt.Printf("[FCM MOCK] Would send to %s: %s - %s\n", token, title, body)
		return nil
	}

	message := FCMMessage{
		To: token,
		Notification: &FCMNotification{
			Title: title,
			Body:  body,
			Sound: "default",
		},
		Data: data,
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal FCM message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://fcm.googleapis.com/fcm/send", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create FCM request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "key="+c.serverKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send FCM request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("FCM returned non-200 status: %d", resp.StatusCode)
	}

	var fcmResp FCMResponse
	if err := json.NewDecoder(resp.Body).Decode(&fcmResp); err != nil {
		return fmt.Errorf("failed to decode FCM response: %w", err)
	}

	if fcmResp.Failure > 0 {
		return fmt.Errorf("FCM reported %d failures", fcmResp.Failure)
	}

	return nil
}

// SendToExpoToken отправляет уведомление через push-сервис Expo
func (c *FCMClient) SendToExpoToken(ctx context.Context, token string, title string, body string, data map[string]string) error {
	message := map[string]interface{}{
		"to":    token,
		"sound": "default",
		"title": title,
		"body":  body,
		"data":  data,
	}

	jsonData, err := json.Marshal([]interface{}{message})
	if err != nil {
		return fmt.Errorf("failed to marshal Expo message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://exp.host/--/api/v2/push/send", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create Expo request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send Expo request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Expo returned non-200 status: %d", resp.StatusCode)
	}

	return nil
}
