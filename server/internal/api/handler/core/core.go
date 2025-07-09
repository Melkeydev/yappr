package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/melkeydev/chat-go/internal/api/model"
	"github.com/melkeydev/chat-go/internal/ws"
	"github.com/melkeydev/chat-go/util"
)

type CoreHandler struct {
	core *ws.Core
}

func NewCoreHandler(c *ws.Core) *CoreHandler {
	return &CoreHandler{
		core: c,
	}
}

func (h *CoreHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var req model.CreateRoomReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	h.core.Rooms[req.ID] = &ws.Room{
		ID:      req.ID,
		Name:    req.Name,
		Clients: make(map[string]*ws.Client),
	}

	util.WriteJSON(w, http.StatusOK, req)
}

func (h *CoreHandler) JoinRoom(w http.ResponseWriter, r *http.Request) {
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

	roomID := chi.URLParam(r, "roomId") // from /ws/{roomId}
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
	rooms := make([]model.RoomRes, 0)

	for _, r := range h.core.Rooms {
		rooms = append(rooms, model.RoomRes{
			ID:   r.ID,
			Name: r.Name,
		})
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
