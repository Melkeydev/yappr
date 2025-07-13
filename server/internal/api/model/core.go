package model

type CreateRoomReq struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
}

type ClientRes struct {
	ID       string `json:"id"`
	Username string `json:"username" `
}

type RoomRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
