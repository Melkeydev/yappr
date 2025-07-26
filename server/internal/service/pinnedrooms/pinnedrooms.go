package pinnedrooms

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	roomRepo "github.com/melkeydev/chat-go/internal/repo/room"
	"github.com/melkeydev/chat-go/internal/service/topics"
	"github.com/melkeydev/chat-go/internal/ws"
)

type PinnedRoomsService struct {
	roomRepo     *roomRepo.RoomRepository
	topicService *topics.TopicService
	wsCore       *ws.Core
}

func NewPinnedRoomsService(db *sql.DB, wsCore *ws.Core) *PinnedRoomsService {
	return &PinnedRoomsService{
		roomRepo:     roomRepo.NewRoomRepository(db),
		topicService: topics.NewTopicService(),
		wsCore:       wsCore,
	}
}

// getNextMidnightUTC returns the next midnight UTC time
func getNextMidnightUTC() time.Time {
	// TEMPORARY: Changed to 2 minutes for debugging
	return time.Now().Add(2 * time.Minute)
	// Original code:
	// now := time.Now().UTC()
	// midnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
	// return midnight
}

// RefreshPinnedRooms creates new pinned rooms with fresh topics
func (s *PinnedRoomsService) RefreshPinnedRooms(ctx context.Context) error {
	log.Println("Refreshing pinned rooms...")

	// Fetch fresh topics
	topics, err := s.topicService.FetchAllTopics(ctx)
	if err != nil {
		return fmt.Errorf("fetch topics: %w", err)
	}

	// Get the next midnight UTC for all pinned rooms to expire at the same time
	expiresAt := getNextMidnightUTC()
	now := time.Now()

	// Create rooms for each topic
	roomNames := []string{"Tech Talk", "World News", "Fun Facts"}
	
	for i, topic := range topics {
		if i >= len(roomNames) {
			break
		}

		room := &roomRepo.Room{
			Name:             roomNames[i],
			IsPinned:         true,
			TopicTitle:       &topic.Title,
			TopicDescription: &topic.Description,
			TopicURL:         &topic.URL,
			TopicSource:      &topic.Source,
			TopicUpdatedAt:   &now,
			ExpiresAt:        expiresAt,
		}

		log.Printf("Creating pinned room %s with topic data:", roomNames[i])
		log.Printf("  Title: %v", room.TopicTitle)
		log.Printf("  Description: %v", room.TopicDescription)
		log.Printf("  URL: %v", room.TopicURL)
		log.Printf("  Source: %v", room.TopicSource)
		
		createdRoom, err := s.roomRepo.CreateRoom(ctx, room)
		if err != nil {
			log.Printf("Failed to create pinned room %s: %v", roomNames[i], err)
			continue
		}
		
		log.Printf("Created room in DB with ID: %s", createdRoom.ID.String())
		log.Printf("  DB Title: %v", createdRoom.TopicTitle)
		log.Printf("  DB Description: %v", createdRoom.TopicDescription)
		log.Printf("  DB URL: %v", createdRoom.TopicURL)
		log.Printf("  DB Source: %v", createdRoom.TopicSource)

		// Add to WebSocket core's in-memory map
		s.wsCore.Rooms[createdRoom.ID.String()] = &ws.Room{
			ID:               createdRoom.ID.String(),
			Name:             createdRoom.Name,
			Clients:          make(map[string]*ws.Client),
			IsPinned:         createdRoom.IsPinned,
			TopicTitle:       createdRoom.TopicTitle,
			TopicDescription: createdRoom.TopicDescription,
			TopicURL:         createdRoom.TopicURL,
			TopicSource:      createdRoom.TopicSource,
		}

		log.Printf("Created pinned room: %s with topic: %s", createdRoom.Name, topic.Title)
	}

	return nil
}

// CheckAndRefreshPinnedRooms checks if pinned rooms need to be created
func (s *PinnedRoomsService) CheckAndRefreshPinnedRooms(ctx context.Context) error {
	count, err := s.roomRepo.CountPinnedRooms(ctx)
	if err != nil {
		return fmt.Errorf("count pinned rooms: %w", err)
	}

	// If we don't have 3 pinned rooms, create them
	if count < 3 {
		log.Printf("Only %d pinned rooms found, creating new ones...", count)
		return s.RefreshPinnedRooms(ctx)
	}

	return nil
}