package dashboard_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type fakeDevice struct {
	id       string
	lastSeen time.Time
	deleted  bool
}

type fakeMessage struct {
	id        string
	deviceID  string
	state     string
	createdAt time.Time
}

type fakeUpstream struct {
	server *httptest.Server

	devices  []fakeDevice
	messages []fakeMessage

	mu        sync.Mutex
	msgCounts map[string]int
	msgHits   atomic.Int64
	failMsgs  atomic.Bool
}

func newFakeUpstream() *fakeUpstream {
	f := &fakeUpstream{
		msgCounts: make(map[string]int),
	}
	f.server = httptest.NewServer(http.HandlerFunc(f.handle))

	return f
}

func (f *fakeUpstream) close() {
	f.server.Close()
}

func (f *fakeUpstream) setCount(state, date string, count int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.msgCounts[state+"|"+date] = count
}

func (f *fakeUpstream) handle(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && r.URL.Path == "/messages":
		f.handleMessages(w, r)
	case r.Method == http.MethodGet && r.URL.Path == "/devices":
		f.handleDevices(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (f *fakeUpstream) handleMessages(w http.ResponseWriter, r *http.Request) {
	f.msgHits.Add(1)

	if f.failMsgs.Load() {
		http.Error(w, "upstream error", http.StatusInternalServerError)
		return
	}

	query := r.URL.Query()

	f.mu.Lock()
	total, items := f.pickMessages(query, query.Get("state"))
	f.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	_ = json.NewEncoder(w).Encode(items)
}

// pickMessages returns the X-Total-Count and page items for a /messages request.
// When message fixtures are set, they are filtered and paginated; otherwise the
// legacy count fixtures (state+date keys) are used.
func (f *fakeUpstream) pickMessages(query url.Values, state string) (int, []map[string]any) {
	if len(f.messages) > 0 {
		return f.pageMessages(query, state)
	}

	date := ""
	if from := query.Get("from"); from != "" {
		parsed, err := time.Parse(time.RFC3339, from)
		if err == nil {
			date = parsed.UTC().Format("2006-01-02")
		}
	}

	return f.msgCounts[state+"|"+date], []map[string]any{}
}

func (f *fakeUpstream) pageMessages(query url.Values, state string) (int, []map[string]any) {
	now := time.Now().UTC()

	filtered := make([]fakeMessage, 0, len(f.messages))
	for _, m := range f.messages {
		if state != "" && m.state != state {
			continue
		}
		if from := query.Get("from"); from != "" {
			if parsed, err := time.Parse(time.RFC3339, from); err == nil && m.createdAt.Before(parsed) {
				continue
			}
		}
		if to := query.Get("to"); to != "" {
			if parsed, err := time.Parse(time.RFC3339, to); err == nil && !m.createdAt.Before(parsed) {
				continue
			}
		}
		filtered = append(filtered, m)
	}

	limit := 100
	if v := query.Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}
	offset := 0
	if v := query.Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			offset = n
		}
	}
	if offset < 0 {
		offset = 0
	}
	if offset > len(filtered) {
		offset = len(filtered)
	}
	end := min(offset+limit, len(filtered))

	items := make([]map[string]any, 0, end-offset)
	for _, m := range filtered[offset:end] {
		items = append(items, map[string]any{
			"id":         m.id,
			"deviceId":   m.deviceID,
			"state":      m.state,
			"recipients": []map[string]any{},
			"states":     map[string]any{m.state: now.Format(time.RFC3339)},
		})
	}

	return len(filtered), items
}

func (f *fakeUpstream) handleDevices(w http.ResponseWriter, _ *http.Request) {
	now := time.Now().UTC()

	items := make([]map[string]any, 0, len(f.devices))
	for _, d := range f.devices {
		var deletedAt *time.Time
		if d.deleted {
			t := now
			deletedAt = &t
		}

		items = append(items, map[string]any{
			"id":        d.id,
			"name":      d.id,
			"createdAt": now.Format(time.RFC3339),
			"updatedAt": now.Format(time.RFC3339),
			"lastSeen":  d.lastSeen.UTC().Format(time.RFC3339),
			"deletedAt": deletedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(items)
}
