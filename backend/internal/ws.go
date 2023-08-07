package internal

import (
	"fmt"
	"github.com/gofiber/contrib/websocket"
)

type Client struct {
	name string
	conn *websocket.Conn
	send chan []byte
}

func (c *Client) read(room *Room) {
	defer func() {
		room.unregister <- c
		c.conn.Close()
	}()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			room.unregister <- c
			c.conn.Close()
			break
		}

		prepareMessage := fmt.Sprintf("[%s]: %s", c.name, msg)
		room.broadcast <- []byte(prepareMessage)
	}
}

func NewClient(name string, conn *websocket.Conn) *Client {
	return &Client{
		name: name,
		conn: conn,
		send: make(chan []byte),
	}
}

type Room struct {
	name       string
	members    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func (r *Room) run() {
	// running this in a goroutine so it doesn't block the main thread
	go func() {
		for {
			select {
			case client := <-r.register:
				r.members[client] = true
			case client := <-r.unregister:
				fmt.Println("unregistering client: ", client)
				// check if client is connected
				if _, ok := r.members[client]; ok {
					// delete client from members
					delete(r.members, client)
					// close channel send for client
					close(client.send)
				}
			case message := <-r.broadcast:
				for client := range r.members {
					// this is a goroutine, so it doesn't block the main thread
					// each write is done in a separate goroutine to avoid blocking
					go func(client *Client) {
						client.conn.WriteMessage(1, message)
					}(client)

				}

			}
		}
	}()
}

func NewRoom(name string) *Room {
	room := &Room{
		name:    name,
		members: make(map[*Client]bool),
		// channels can be any type, here I'm creating a channel of type slice of bytes
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	return room
}

type RoomStore struct {
	rooms map[string]*Room
}

func (rs *RoomStore) GetRoom(name string) (*Room, bool) {
	room, exists := rs.rooms[name]
	return room, exists
}

func (rs *RoomStore) CreateRoom(name string) (*Room, error) {
	if _, exists := rs.rooms[name]; exists {
		return nil, fmt.Errorf("room %s already exists", name)
	}
	room := NewRoom(name)
	rs.rooms[name] = room
	room.run()
	return room, nil
}

func NewRoomStore() *RoomStore {
	return &RoomStore{
		rooms: make(map[string]*Room),
	}
}

func (rs *RoomStore) WsHandler(conn *websocket.Conn) {
	roomNum := conn.Params("room")
	user := conn.Params("name")

	room, exists := rs.rooms[roomNum]
	if !exists {
		room, _ = rs.CreateRoom(roomNum)
	}

	conn.WriteMessage(1, []byte("you just joined room: "+roomNum)) // welcome user
	client := NewClient(user, conn)
	room.broadcast <- []byte(user + " joined the room!") // broadcast to all users that a new user joined
	room.register <- client                              // we're sending the client to the register channel of the room
	client.read(room)                                    // this is a goroutine, so it doesn't block the main thread
}
