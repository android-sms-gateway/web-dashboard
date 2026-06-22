package events

import "encoding/json"

type EventType string

const (
	EventMessageReceived     EventType = "message.received"
	EventMessageStateChanged EventType = "message.state_changed"
	EventDeviceStatusChanged EventType = "device.status_changed"
	EventStatsUpdated        EventType = "stats.updated"
)

type Event struct {
	Type    EventType       `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

type CallbackPayload struct {
	Event   string          `json:"event"`
	Payload json.RawMessage `json:"payload"`
}
