package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID               uuid.UUID  `json:"id"`
	Name             string     `json:"name"`
	CreatorID        *uuid.UUID `json:"creator_id"`
	CreatedAt        time.Time  `json:"created_at"`
	ExpiresAt        time.Time  `json:"expires_at"`
	IsPinned         bool       `json:"is_pinned"`
	TopicTitle       *string    `json:"topic_title,omitempty"`
	TopicDescription *string    `json:"topic_description,omitempty"`
	TopicURL         *string    `json:"topic_url,omitempty"`
	TopicSource      *string    `json:"topic_source,omitempty"`
	TopicUpdatedAt   *time.Time `json:"topic_updated_at,omitempty"`
}

type Message struct {
	ID        uuid.UUID  `json:"id"`
	RoomID    uuid.UUID  `json:"room_id"`
	UserID    *uuid.UUID `json:"user_id,omitempty"`
	Username  string     `json:"username"`
	Content   string     `json:"content"`
	IsSystem  bool       `json:"is_system"`
	CreatedAt time.Time  `json:"created_at"`
}

type RoomRepository struct {
	db *sql.DB
}

func NewRoomRepository(db *sql.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) CreateRoom(ctx context.Context, room *Room) (*Room, error) {
	var query string
	var err error
	
	if room.IsPinned {
		// For pinned rooms, we can set a custom expires_at time
		query = `
			INSERT INTO rooms (name, creator_id, is_pinned, topic_title, topic_description, topic_url, topic_source, topic_updated_at, expires_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING id, created_at, expires_at
		`
		err = r.db.QueryRowContext(ctx, query, 
			room.Name, room.CreatorID, room.IsPinned, 
			room.TopicTitle, room.TopicDescription, room.TopicURL, 
			room.TopicSource, room.TopicUpdatedAt, room.ExpiresAt,
		).Scan(
			&room.ID,
			&room.CreatedAt,
			&room.ExpiresAt,
		)
	} else {
		// Regular rooms get default 24-hour expiration
		query = `
			INSERT INTO rooms (name, creator_id)
			VALUES ($1, $2)
			RETURNING id, created_at, expires_at
		`
		err = r.db.QueryRowContext(ctx, query, room.Name, room.CreatorID).Scan(
			&room.ID,
			&room.CreatedAt,
			&room.ExpiresAt,
		)
	}

	if err != nil {
		return nil, fmt.Errorf("insert room: %w", err)
	}

	return room, nil
}

func (r *RoomRepository) GetRoomByID(ctx context.Context, id uuid.UUID) (*Room, error) {
	query := `
		SELECT id, name, creator_id, created_at, expires_at, is_pinned, 
		       topic_title, topic_description, topic_url, topic_source, topic_updated_at
		FROM rooms
		WHERE id = $1 AND expires_at > NOW()
	`

	var room Room
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&room.ID,
		&room.Name,
		&room.CreatorID,
		&room.CreatedAt,
		&room.ExpiresAt,
		&room.IsPinned,
		&room.TopicTitle,
		&room.TopicDescription,
		&room.TopicURL,
		&room.TopicSource,
		&room.TopicUpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Room not found or expired
		}
		return nil, fmt.Errorf("query room by id: %w", err)
	}

	return &room, nil
}

func (r *RoomRepository) GetAllActiveRooms(ctx context.Context) ([]*Room, error) {
	query := `
		SELECT id, name, creator_id, created_at, expires_at, is_pinned,
		       topic_title, topic_description, topic_url, topic_source, topic_updated_at
		FROM rooms
		WHERE expires_at > NOW()
		ORDER BY is_pinned DESC, created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query all rooms: %w", err)
	}
	defer rows.Close()

	var rooms []*Room
	for rows.Next() {
		var room Room
		err := rows.Scan(
			&room.ID,
			&room.Name,
			&room.CreatorID,
			&room.CreatedAt,
			&room.ExpiresAt,
			&room.IsPinned,
			&room.TopicTitle,
			&room.TopicDescription,
			&room.TopicURL,
			&room.TopicSource,
			&room.TopicUpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan room: %w", err)
		}
		rooms = append(rooms, &room)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rooms: %w", err)
	}

	return rooms, nil
}

func (r *RoomRepository) CountActiveRooms(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM rooms WHERE expires_at > NOW()`
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count active rooms: %w", err)
	}
	return count, nil
}

func (r *RoomRepository) CreateMessage(ctx context.Context, msg *Message) (*Message, error) {
	query := `
		INSERT INTO messages (room_id, user_id, username, content, is_system)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	err := r.db.QueryRowContext(
		ctx, query,
		msg.RoomID, msg.UserID, msg.Username, msg.Content, msg.IsSystem,
	).Scan(&msg.ID, &msg.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("insert message: %w", err)
	}

	return msg, nil
}

func (r *RoomRepository) GetRoomMessages(ctx context.Context, roomID uuid.UUID, limit int) ([]*Message, error) {
	query := `
		SELECT m.id, m.room_id, m.user_id, m.username, m.content, m.is_system, m.created_at
		FROM messages m
		INNER JOIN rooms r ON m.room_id = r.id
		WHERE m.room_id = $1 AND r.expires_at > NOW()
		ORDER BY m.created_at DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, roomID, limit)
	if err != nil {
		return nil, fmt.Errorf("query room messages: %w", err)
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(
			&msg.ID,
			&msg.RoomID,
			&msg.UserID,
			&msg.Username,
			&msg.Content,
			&msg.IsSystem,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan message: %w", err)
		}
		messages = append(messages, &msg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate messages: %w", err)
	}

	// Reverse the messages to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func (r *RoomRepository) DeleteExpiredRooms(ctx context.Context) (int, error) {
	query := `DELETE FROM rooms WHERE expires_at <= NOW()`
	
	result, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("delete expired rooms: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("get rows affected: %w", err)
	}

	return int(rowsAffected), nil
}

func (r *RoomRepository) HasActiveRoom(ctx context.Context, userID uuid.UUID) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM rooms WHERE creator_id = $1 AND expires_at > NOW()`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("check active room: %w", err)
	}
	return count > 0, nil
}

func (r *RoomRepository) CountPinnedRooms(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM rooms WHERE is_pinned = true AND expires_at > NOW()`
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count pinned rooms: %w", err)
	}
	return count, nil
}