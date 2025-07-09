package ws

type Room struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"clients"`
	History []*Message
}

type Core struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewCore() *Core {
	return &Core{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
	}
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
				// replay
				go func(h []*Message) {
					for _, m := range h {
						cl.Message <- m
					}
				}(room.History)
			}

		case cl := <-c.Unregister:
			if _, ok := c.Rooms[cl.RoomID]; ok {
				if _, ok := c.Rooms[cl.RoomID].Clients[cl.ID]; ok {
					if len(c.Rooms[cl.RoomID].Clients) != 0 {
						c.Broadcast <- &Message{
							Content:  "user left the chat",
							RoomID:   cl.RoomID,
							Username: cl.Username,
						}
					}

					delete(c.Rooms[cl.RoomID].Clients, cl.ID)
					close(cl.Message)
				}
			}

			// FAN OUT
		case m := <-c.Broadcast:
			if room, ok := c.Rooms[m.RoomID]; ok {
				room.History = append(room.History, m) // NEW
				for _, cl := range room.Clients {
					cl.Message <- m
				}
			}
		}
	}
}
