package internal

import (
	"log"

	"github.com/gorilla/websocket"
)

type GorillaClientList map[*GorillaClient]bool

type GorillaClient struct {
	conn    *websocket.Conn
	manager *Manager

	// egress is used to avoid concurrent writes on the websocket connection
	egress chan []byte
}

func GorillaNewClient(conn *websocket.Conn, manager *Manager) *GorillaClient {
	return &GorillaClient{
		conn:    conn,
		manager: manager,
		egress:  make(chan []byte),
	}
}

func (c *GorillaClient) readMessage() {
	defer func() {
		c.manager.removeClient(c)
	}()

	for {
		msgType, payload, err := c.conn.ReadMessage()
		if err != nil {
			// if an unexpected close error happens we log below.
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}

		for client := range c.manager.clients {
			client.egress <- payload
		}

		log.Println(msgType)
		log.Println("message received: ", string(payload))
	}
}

func (c *GorillaClient) writeMessage() {
	defer func() {
		c.manager.removeClient(c)
	}()

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed: ", err)
				}
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("failed to send message: %v", err)
			}
			log.Println("message sent!")
		}
	}

}
