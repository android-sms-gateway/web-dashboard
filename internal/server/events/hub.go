package events

import (
	"encoding/json"
	"fmt"
	"sync"
)

type Hub struct {
	mu    sync.RWMutex
	conns map[string]map[string]*Connection
}

func NewHub() *Hub {
	return &Hub{
		mu:    sync.RWMutex{},
		conns: make(map[string]map[string]*Connection),
	}
}

func (h *Hub) Add(conn *Connection) {
	if conn.User == "" {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if h.conns[conn.User] == nil {
		h.conns[conn.User] = make(map[string]*Connection)
	}

	if existing, ok := h.conns[conn.User][conn.ID]; ok {
		existing.Close()
	}

	h.conns[conn.User][conn.ID] = conn
}

func (h *Hub) Remove(id, userID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	conns, exists := h.conns[userID]
	if !exists {
		return
	}

	if conn, found := conns[id]; found {
		conn.Close()
		delete(conns, id)
	}

	if len(conns) == 0 {
		delete(h.conns, userID)
	}
}

func (h *Hub) SendToUser(userID string, event Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, conn := range h.conns[userID] {
		conn.Send(data)
	}

	return nil
}

func (h *Hub) Len() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	n := 0
	for _, conns := range h.conns {
		n += len(conns)
	}

	return n
}
