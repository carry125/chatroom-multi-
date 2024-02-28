package main

import "fmt"

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	//identity of room
	roomId string
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients. 把消息存到這個channel 之後會有其他goroutine遍歷client 把消息發給客戶端
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub(roomId string) *Hub { //實體化
	return &Hub{
		roomId:     roomId,
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	defer func() {
		close(h.register)
		close(h.broadcast)
		close(h.unregister)
	}()
	for { //不斷從channel讀取數據
		select {
		case client := <-h.register: //註冊客戶端
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {

				delete(h.clients, client) //刪除引用

				close(client.send) //切斷連線
			}
			if len(h.clients) == 0 {
				delete(house, h.roomId)
				fmt.Println("logout")
				return //如果只用break 只會結束select不會結束for
			}
		case message := <-h.broadcast:
			for client := range h.clients { //對每個客戶端
				select {
				case client.send <- message: //把消息寫入每個客戶端的client.send channel中，實現廣播
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
