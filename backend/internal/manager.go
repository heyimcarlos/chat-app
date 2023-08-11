package internal

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var allowOrigins = map[string]bool{
	"http://localhost:3000": true,
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			if origin == "" {
				return true
			}
			allow, ok := allowOrigins[origin]
			if ok && allow {
				return true
			}
			return false
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Manager struct {
	clients      GorillaClientList
	sync.RWMutex // provides a lock/unlock mechanism to protect the map, when clients are connecting and disconnecting.

	handlers map[string]EventHandler
}

func NewManager() *Manager {
	m := &Manager{
		clients:  make(GorillaClientList),
		handlers: make(map[string]EventHandler),
	}

	m.setupEventHandlers()
	return m
}

func (m *Manager) setupEventHandlers() {
	m.handlers[EventSendMessage] = SendMessage
}

func SendMessage(event Event, c *Client) error {
	fmt.Println(event)
	return nil
}

func (m *Manager) routeEvent(event Event, c *Client) error {
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("there is no handler for this event type")
	}
}

func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {
	log.Println("Websocket connection received")

	// upgrade regular HTTP connection to a websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := GorillaNewClient(conn, m)

	m.addClient(client)

	// Start client processses
	go client.readMessage()
	go client.writeMessage()
}

func (m *Manager) addClient(client *GorillaClient) {
	m.Lock()         // locking the map helps to lock the manager in case two clients are connecting at the same time.
	defer m.Unlock() // unlocks the map after the client is added.

	if _, ok := m.clients[client]; !ok {
		m.clients[client] = true
	}
}

func (m *Manager) removeClient(client *GorillaClient) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		client.conn.Close()
		delete(m.clients, client)
	}
}
