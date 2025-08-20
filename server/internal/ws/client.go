package ws

import (
	"time"

	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	ID       string `json:"id"`
	RoomID   string `json:"room_id"`
	Username string `json:"username"`
}

type Message struct {
	Content   string `json:"content"`
	RoomID    string `json:"room_id"`
	Username  string `json:"username"`
	UserID    string `json:"user_id,omitempty"`
	System    bool   `json:"system"`
	Timestamp string `json:"timestamp,omitempty"`
}

func (c *Client) ReadMessage(core *Core) {
	defer func() {
		core.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg := &Message{
			Content:   string(m),
			RoomID:    c.RoomID,
			Username:  c.Username,
			UserID:    c.ID,
			Timestamp: time.Now().Format("2006-01-02T15:04:05Z07:00"),
		}

		core.Broadcast <- msg
	}
}

func (c *Client) WriteMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.Message
		if !ok {
			return
		}

		c.Conn.WriteJSON(message)
	}
}
