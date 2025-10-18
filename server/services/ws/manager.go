package ws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		CheckOrigin:     checkOrigin,
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Manager struct {
	clients ClientList
	sync.RWMutex

	handlers map[string]EventHandler
}

func NewManager() *Manager {
	m := &Manager{
		clients:  make(ClientList),
		handlers: make(map[string]EventHandler),
	}
	m.setupEventHandlers()
	return m
}

func (m *Manager) setupEventHandlers() {
	m.handlers[EventOrderPlaced] = OrderPlaced
}

func OrderPlaced(event Event, m *Manager, to string) error {
	var chatevent OrderPlacedEvent

	if err := json.Unmarshal(event.Payload, &chatevent); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	data, err := json.Marshal(chatevent)
	if err != nil {
		return fmt.Errorf("failed to marshal chatevent message: %v", err)
	}

	outgoingEvent := Event{
		Payload: data,
		Type:    EventOrderPlaced,
	}

	for client := range m.clients {
		if client.id == to {
			client.egress <- outgoingEvent
			break
		}
	}

	return nil
}

// call it directly to send
// to and event type send here inside event
func (m *Manager) RouteEvent(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return fmt.Errorf("invalid http method")
	}

	event := Event{}
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		return fmt.Errorf("invalid req payload, err: %v", err)
	}
	if handler, ok := m.handlers[event.Type]; ok {
		var reqPayload struct {
			ID string `json:"to"`
		}
		if err := json.Unmarshal(event.Payload, &reqPayload); err != nil {
			return err
		}
		if err := handler(event, m, reqPayload.ID); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("invalid event type")
	}
}

func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		log.Println("Missing id query param")
		return
	}

	log.Println("new connection")

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(id, conn, m)
	m.addClient(client)

	go client.readMessages()
	go client.writeMessages()
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	m.clients[client] = true
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		client.connection.Close()
		delete(m.clients, client)
	}
}

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	switch origin {
	case "http://localhost:5173":
		return true
	default:
		return false
	}
}
