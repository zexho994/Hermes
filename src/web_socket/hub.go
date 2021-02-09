package main

type Hub struct{
	users map[*User]bool
	broadcast chan []byte
	register chan *User
	unregister chan *User
}

func newHub() *Hub {
	return &Hub{
		broadcast: make(chan []byte),
		register: make(chan *User),
		unregister: make(chan *User),
		users: make(map[*User]bool),
	}
}

func (h *Hub) run(){
	for {
		select {
		case user:= <-h.register:
			h.users[user] = true

		case user:= <-h.unregister:
			if _,ok := h.users[user]; ok {
				delete(h.users,user)
				close(user.send)
			}
		case message := <-h.broadcast:
			for user := range h.users {
				select {
				case user.send <- message :
				default:
					close(user.send)
					delete(h.users,user)
				}
			}
		}
	}
}
