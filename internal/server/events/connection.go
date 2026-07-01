package events

import (
	"bufio"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	sseKeepaliveInterval = 30 * time.Second
	sseChanCapacity      = 10
)

type Connection struct {
	ID   string
	User string

	send chan []byte
	done chan struct{}
}

func NewConnection(id, userID string) *Connection {
	return &Connection{
		ID:   id,
		User: userID,
		send: make(chan []byte, sseChanCapacity),
		done: make(chan struct{}),
	}
}

func (c *Connection) Send(data []byte) {
	select {
	case c.send <- data:
	default:
	}
}

func (c *Connection) Close() {
	select {
	case <-c.done:
	default:
		close(c.done)
	}
}

func (c *Connection) Done() <-chan struct{} {
	return c.done
}

func (c *Connection) Stream(ctx *fiber.Ctx, onDone func()) error {
	ctx.Set("Content-Type", "text/event-stream")
	ctx.Set("Cache-Control", "no-cache")
	ctx.Set("Connection", "keep-alive")

	ctx.Status(fiber.StatusOK).Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		defer c.Close()
		c.writeLoop(w)
		if onDone != nil {
			onDone()
		}
	})
	return nil
}

func (c *Connection) writeLoop(w *bufio.Writer) {
	keepalive := time.NewTicker(sseKeepaliveInterval)
	defer keepalive.Stop()

	for {
		select {
		case <-c.Done():
			return
		case <-keepalive.C:
			if !c.writeKeepalive(w) {
				return
			}
		case data, ok := <-c.send:
			if !ok {
				return
			}
			if !c.writeData(w, data) {
				return
			}
		}
	}
}

func (c *Connection) writeKeepalive(w *bufio.Writer) bool {
	if _, err := fmt.Fprintf(w, ": keepalive\n\n"); err != nil {
		return false
	}

	return w.Flush() == nil
}

func (c *Connection) writeData(w *bufio.Writer, data []byte) bool {
	if _, err := fmt.Fprintf(w, "event: message\ndata: %s\n\n", data); err != nil {
		return false
	}

	return w.Flush() == nil
}

func (c *Connection) SendEvent(event Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	c.Send(data)

	return nil
}
