package ws

import (
	"context"
	"database/sql"
	"log"

	"github.com/google/uuid"
	roomRepo "github.com/melkeydev/chat-go/internal/repo/room"
)

type Room struct {
	ID               string             `json:"id"`
	Name             string             `json:"name"`
	Clients          map[string]*Client `json:"clients"`
	History          []*Message
	IsPinned         bool               `json:"is_pinned"`
	TopicTitle       *string            `json:"topic_title,omitempty"`
	TopicDescription *string            `json:"topic_description,omitempty"`
	TopicURL         *string            `json:"topic_url,omitempty"`
	TopicSource      *string            `json:"topic_source,omitempty"`
}

type Core struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
	roomRepo   *roomRepo.RoomRepository
	db         *sql.DB
}

func NewCore(db *sql.DB) *Core {
	return &Core{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
		roomRepo:   roomRepo.NewRoomRepository(db),
		db:         db,
	}
}

func (c *Core) GetDB() *sql.DB {
	return c.db
}

// The core will be ran in a different go Routine
func (c *Core) Run() {
	for {
		select {
		case cl := <-c.Register:
			if room, ok := c.Rooms[cl.RoomID]; ok {
				if _, ok := room.Clients[cl.ID]; !ok {
					room.Clients[cl.ID] = cl
				}
				// Load and replay history from database
				go func() {
					roomUUID, err := uuid.Parse(cl.RoomID)
					if err != nil {
						log.Printf("Invalid room ID: %v", err)
						return
					}
					
					messages, err := c.roomRepo.GetRoomMessages(context.Background(), roomUUID, 100)
					if err != nil {
						log.Printf("Failed to load room messages: %v", err)
						return
					}
					
					for _, msg := range messages {
						wsMsg := &Message{
							Content:  msg.Content,
							RoomID:   cl.RoomID,
							Username: msg.Username,
							System:   msg.IsSystem,
						}
						cl.Message <- wsMsg
					}
				}()
			}

		case cl := <-c.Unregister:
			if _, ok := c.Rooms[cl.RoomID]; ok {
				if _, ok := c.Rooms[cl.RoomID].Clients[cl.ID]; ok {
					delete(c.Rooms[cl.RoomID].Clients, cl.ID)
					close(cl.Message)
				}
			}

			// FAN OUT
		case m := <-c.Broadcast:
			if room, ok := c.Rooms[m.RoomID]; ok {
				room.History = append(room.History, m)
				
				// Persist message to database
				go func(msg *Message) {
					roomUUID, err := uuid.Parse(msg.RoomID)
					if err != nil {
						log.Printf("Invalid room ID: %v", err)
						return
					}
					
					dbMsg := &roomRepo.Message{
						RoomID:   roomUUID,
						Username: msg.Username,
						Content:  msg.Content,
						IsSystem: msg.System,
					}
					
					// Try to parse user ID if it's in the username format
					// This is a placeholder - you might want to pass actual user ID
					if _, err := c.roomRepo.CreateMessage(context.Background(), dbMsg); err != nil {
						log.Printf("Failed to persist message: %v", err)
					}
				}(m)
				
				for _, cl := range room.Clients {
					cl.Message <- m
				}
			}
		}
	}
}
