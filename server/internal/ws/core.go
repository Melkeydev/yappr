package ws

type Room struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"clients"`
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
			if _, ok := c.Rooms[cl.RoomID]; ok {
				r := c.Rooms[cl.RoomID]

				if _, ok := r.Clients[cl.ID]; !ok {
					r.Clients[cl.ID] = cl
				}
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
			if _, ok := c.Rooms[m.RoomID]; ok {

				for _, cl := range c.Rooms[m.RoomID].Clients {
					cl.Message <- m
				}
			}
		}
	}
}
