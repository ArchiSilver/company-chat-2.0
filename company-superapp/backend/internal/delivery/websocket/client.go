package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Разрешаем все origins для разработки
	},
}

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	userID uuid.UUID
	chatID uuid.UUID
}

type IncomingMessage struct {
	Type    string `json:"type"`
	ChatID  string `json:"chat_id"`
	Content string `json:"content"`
}

type OutgoingMessage struct {
	Type      string    `json:"type"`
	ID        int64     `json:"id,omitempty"`
	ChatID    string    `json:"chat_id"`
	SenderID  string    `json:"sender_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var incoming IncomingMessage
		if err := json.Unmarshal(message, &incoming); err != nil {
			log.Printf("error unmarshalling message: %v", err)
			continue
		}

		if incoming.Type == "message" {
			chatID, _ := uuid.Parse(incoming.ChatID)
			outgoing := OutgoingMessage{
				Type:      "message",
				ChatID:    incoming.ChatID,
				SenderID:  c.userID.String(),
				Content:   incoming.Content,
				CreatedAt: time.Now(),
			}

			outgoingBytes, _ := json.Marshal(outgoing)

			// Publish to Redis for other server instances
			c.hub.redis.Publish(c.hub.ctx, "messages:"+chatID.String(), outgoingBytes)

			// Send to all clients in the room on this instance
			c.hub.mu.RLock()
			if room, ok := c.hub.rooms[chatID]; ok {
				for client := range room {
					select {
					case client.send <- outgoingBytes:
					default:
						close(client.send)
						delete(c.hub.clients, client)
					}
				}
			}
			c.hub.mu.RUnlock()
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, userID uuid.UUID, chatID uuid.UUID) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		userID: userID,
		chatID: chatID,
	}
	client.hub.register <- client

	// Add client to the chat room
	hub.mu.Lock()
	if hub.rooms[chatID] == nil {
		hub.rooms[chatID] = make(map[*Client]bool)
		go hub.SubscribeToChat(chatID)
	}
	hub.rooms[chatID][client] = true
	hub.mu.Unlock()

	go client.writePump()
	go client.readPump()
}
