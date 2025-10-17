package ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

type ClientList map[*Client]bool

type Client struct {
	id         string
	connection *websocket.Conn
	manager    *Manager
	egress     chan Event
}

func NewClient(id string, conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		id:         id,
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
	}
}

func (c *Client) readMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
		return
	}
	c.connection.SetReadLimit(512)
	c.connection.SetPongHandler(c.pongHandler)

	// for {
	// 	_, payload, err := c.connection.ReadMessage()
	// 	if err != nil {
	// 		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
	// 			log.Println("error reading message: ", err)
	// 		}
	// 		break
	// 	}
	// 	var request Event
	// 	if err := json.Unmarshal(payload, &request); err != nil {
	// 		log.Println("error marshalling event: ", err)
	// 	}

	// 	if err := c.manager.routeEvent(request, c); err != nil {
	// 		log.Println("error handling message: ", err)
	// 	}

	// 	// for wsclient := range c.manager.clients {
	// 	// 	wsclient.egress <- payload
	// 	// }
	// 	// log.Println(messageType)
	// 	// log.Println(string(payload))
	// }
}

func (c *Client) writeMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("error closing connection:", err)
				}
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Println("error marshaling message:", err)
				continue
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println("failed to send message:", err)
				return
			}
			log.Println("message sent")

		case <-ticker.C:
			log.Println("ping")
			if err := c.connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("ping failed:", err)
				return
			}
			log.Println("ping sent")
		}
	}
}

func (c *Client) pongHandler(pongMsg string) error {
	log.Println("pong")
	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}
