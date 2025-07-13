package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/melkeydev/chat-go/internal/api/model"
	roomRepo "github.com/melkeydev/chat-go/internal/repo/room"
	"github.com/melkeydev/chat-go/internal/ws"
	"github.com/melkeydev/chat-go/util"
)

type CoreHandler struct {
	core      *ws.Core
	roomRepo  *roomRepo.RoomRepository
	roomLimit int
}

func NewCoreHandler(c *ws.Core) *CoreHandler {
	// Default room limit is 100, can be overridden by MAX_ROOMS env var
	roomLimit := 50
	if maxRoomsStr := os.Getenv("MAX_ROOMS"); maxRoomsStr != "" {
		if limit, err := strconv.Atoi(maxRoomsStr); err == nil && limit > 0 {
			roomLimit = limit
		}
	}

	return &CoreHandler{
		core:      c,
		roomRepo:  roomRepo.NewRoomRepository(c.GetDB()),
		roomLimit: roomLimit,
	}
}

func (h *CoreHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var req model.CreateRoomReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	ctx := r.Context()

	// Check room limit
	activeRooms, err := h.roomRepo.CountActiveRooms(ctx)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "failed to check room limit")
		return
	}

	if activeRooms >= h.roomLimit {
		util.WriteError(w, http.StatusTooManyRequests, "maximum number of rooms reached")
		return
	}

	// Create room in database
	room := &roomRepo.Room{
		Name: req.Name,
	}
	room, err = h.roomRepo.CreateRoom(ctx, room)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "failed to create room")
		return
	}

	// Add to in-memory map
	h.core.Rooms[room.ID.String()] = &ws.Room{
		ID:      room.ID.String(),
		Name:    room.Name,
		Clients: make(map[string]*ws.Client),
	}

	// Return the room with the database-generated ID
	resp := model.CreateRoomReq{
		ID:   room.ID.String(),
		Name: room.Name,
	}
	util.WriteJSON(w, http.StatusOK, resp)
}

func (h *CoreHandler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomId") // from /ws/{roomId}

	// Verify room exists in database
	roomUUID, err := uuid.Parse(roomID)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid room ID")
		return
	}

	ctx := r.Context()
	dbRoom, err := h.roomRepo.GetRoomByID(ctx, roomUUID)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "failed to verify room")
		return
	}
	if dbRoom == nil {
		util.WriteError(w, http.StatusNotFound, "room not found or expired")
		return
	}

	// Ensure room exists in memory map
	if _, exists := h.core.Rooms[roomID]; !exists {
		h.core.Rooms[roomID] = &ws.Room{
			ID:      roomID,
			Name:    dbRoom.Name,
			Clients: make(map[string]*ws.Client),
		}
	}

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// TODO:tighten this check!
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid connection upgrade")
		return
	}

	q := r.URL.Query()
	clientID := q.Get("userId")
	username := q.Get("username")

	cl := &ws.Client{
		Conn:     conn,
		Message:  make(chan *ws.Message, 10),
		ID:       clientID,
		RoomID:   roomID,
		Username: username,
	}

	h.core.Register <- cl

	go cl.WriteMessage()
	cl.ReadMessage(h.core)

}

func (h *CoreHandler) GetRooms(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Fetch active rooms from database
	dbRooms, err := h.roomRepo.GetAllActiveRooms(ctx)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "failed to fetch rooms")
		return
	}

	rooms := make([]model.RoomRes, 0, len(dbRooms))
	for _, room := range dbRooms {
		rooms = append(rooms, model.RoomRes{
			ID:   room.ID.String(),
			Name: room.Name,
		})

		// Ensure room exists in memory map
		if _, exists := h.core.Rooms[room.ID.String()]; !exists {
			h.core.Rooms[room.ID.String()] = &ws.Room{
				ID:      room.ID.String(),
				Name:    room.Name,
				Clients: make(map[string]*ws.Client),
			}
		}
	}

	util.WriteJSON(w, http.StatusOK, rooms)
}

func (h *CoreHandler) GetClients(w http.ResponseWriter, r *http.Request) {
	var clients []model.ClientRes
	roomID := chi.URLParam(r, "roomId") // from /ws/{roomId}

	if _, ok := h.core.Rooms[roomID]; !ok {
		clients = make([]model.ClientRes, 0)
		util.WriteJSON(w, http.StatusOK, clients)
		return
	}

	for _, c := range h.core.Rooms[roomID].Clients {
		clients = append(clients, model.ClientRes{
			ID:       c.ID,
			Username: c.Username,
		})
	}

	util.WriteJSON(w, http.StatusOK, clients)
}
