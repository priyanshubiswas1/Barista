package ws

import (
	"encoding/json"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, m *Manager, to string) error

const (
	EventOrderPlaced = "order_placed"
	EventSendMessage = "send_message" //send
	EventNewMessage  = "new_message"  //receive
	EventChangeRoom  = "change_room"  //receive
)

type OrderPlacedEvent struct {
	To string `json:"to"`
}
