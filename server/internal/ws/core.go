package ws

import (
	"context"
	"database/sql"
	"log"

	"github.com/google/uuid"
	roomRepo "github.com/Melkeydev/yappr/internal/repo/room"
	statsRepo "github.com/Melkeydev/yappr/internal/repo/stats"
)

type Room struct {
	ID               string             `json:"id"`
	Name             string             `json:"name"`
	Clients          map[string]*Client `json:"clients"`
	History          []*Message
	IsPinned         bool    `json:"is_pinned"`
	TopicTitle       *string `json:"topic_title,omitempty"`
	TopicDescription *string `json:"topic_description,omitempty"`
	TopicURL         *string `json:"topic_url,omitempty"`
	TopicSource      *string `json:"topic_source,omitempty"`
}

type Core struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
	roomRepo   *roomRepo.RoomRepository
	statsRepo  *statsRepo.StatsRepository
	db         *sql.DB
}

func NewCore(db *sql.DB) *Core {
	return &Core{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
		roomRepo:   roomRepo.NewRoomRepository(db),
		statsRepo:  statsRepo.NewStatsRepository(db),
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
						userID := ""
						if msg.UserID != nil {
							userID = msg.UserID.String()
						}

						wsMsg := &Message{
							Content:   msg.Content,
							RoomID:    cl.RoomID,
							Username:  msg.Username,
							UserID:    userID,
							System:    msg.IsSystem,
							Timestamp: msg.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
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

				go func(msg *Message) {
					roomUUID, err := uuid.Parse(msg.RoomID)
					if err != nil {
						log.Printf("Invalid room ID: %v", err)
						return
					}

					var userID *uuid.UUID
					if msg.UserID != "" {
						if parsedUserID, err := uuid.Parse(msg.UserID); err == nil {
							userID = &parsedUserID
						}
					}

					dbMsg := &roomRepo.Message{
						RoomID:   roomUUID,
						UserID:   userID,
						Username: msg.Username,
						Content:  msg.Content,
						IsSystem: msg.System,
					}

					if _, err := c.roomRepo.CreateMessage(context.Background(), dbMsg); err != nil {
						log.Printf("Failed to persist message: %v", err)
					}

					if userID != nil {
						if err := c.statsRepo.IncrementMessageCount(context.Background(), *userID); err != nil {
							log.Printf("Failed to update message count for user %s: %v", userID.String(), err)
						} else {
							go func() {
								_, err := c.statsRepo.CheckAndAwardAchievements(context.Background(), *userID)
								if err != nil {
									log.Printf("Error checking achievements for message sender %s: %v", userID.String(), err)
								}
							}()
						}
					}
				}(m)

				for _, cl := range room.Clients {
					cl.Message <- m
				}
			}
		}
	}
}
