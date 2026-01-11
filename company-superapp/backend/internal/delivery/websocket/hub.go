package websocket

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Hub поддерживает набор активных клиентов и транслирует сообщения клиентам.
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	rooms      map[uuid.UUID]map[*Client]bool
	mu         sync.RWMutex
	redis      *redis.Client
	ctx        context.Context
}

func NewHub(redis *redis.Client) *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		rooms:      make(map[uuid.UUID]map[*Client]bool),
		redis:      redis,
		ctx:        context.Background(),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				// Также удаляем из всех комнат
				for _, room := range h.rooms {
					delete(room, client)
				}
			}
			h.mu.Unlock()
		case <-h.broadcast:
			// Это для широковещательной рассылки всем клиентам, что нам может не понадобиться.
			// Фокусируемся на сообщениях для конкретных комнат.
			log.Println("Broadcast channel received a message, but it's not being sent to anyone.")
		}
	}
}

func (h *Hub) SubscribeToChat(chatId uuid.UUID) {
	pubsub := h.redis.Subscribe(h.ctx, "messages:"+chatId.String())
	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			var messageData map[string]interface{}
			if err := json.Unmarshal([]byte(msg.Payload), &messageData); err != nil {
				log.Printf("error unmarshalling message from redis: %v", err)
				continue
			}

			h.mu.RLock()
			room, ok := h.rooms[chatId]
			if !ok {
				h.mu.RUnlock()
				continue
			}

			for client := range room {
				select {
				case client.send <- []byte(msg.Payload):
				default:
					close(client.send)
					delete(h.clients, client)
					delete(room, client)
				}
			}
			h.mu.RUnlock()
		}
	}()
}
