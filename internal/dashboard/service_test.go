package dashboard_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/android-sms-gateway/web-dashboard/internal/dashboard"
	"github.com/android-sms-gateway/web-dashboard/internal/gateway"
	"github.com/go-core-fx/cachefx/cache"
	"go.uber.org/zap"
)

const (
	stateTotal   = ""
	statePending = "Pending"
	stateFailed  = "Failed"
	stateSent    = "Sent"
)

func newService(upstream *fakeUpstream) *dashboard.Service {
	return dashboard.NewService(gateway.NewFactory(upstream.server.URL), cache.NewMemory(0), zap.NewNop())
}

func bucketDate(now time.Time, offsetDays int) string {
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return today.AddDate(0, 0, offsetDays).Format("2006-01-02")
}

func TestStats(t *testing.T) {
	upstream := newFakeUpstream()
	t.Cleanup(upstream.close)

	now := time.Now().UTC()
	upstream.setCount(stateTotal, "", 10)
	upstream.setCount(statePending, "", 2)
	upstream.setCount(stateFailed, "", 1)
	upstream.devices = []fakeDevice{
		{id: "online", lastSeen: now},
		{id: "offline", lastSeen: now.Add(-time.Hour)},
		{id: "deleted", lastSeen: now, deleted: true},
	}

	svc := newService(upstream)
	stats, err := svc.Stats(context.Background(), "user", "pass")
	if err != nil {
		t.Fatalf("Stats() error: %v", err)
	}

	if stats.MessagesSent != 7 {
		t.Errorf("MessagesSent = %d, want 7", stats.MessagesSent)
	}
	if stats.MessagesPending != 2 {
		t.Errorf("MessagesPending = %d, want 2", stats.MessagesPending)
	}
	if stats.MessagesFailed != 1 {
		t.Errorf("MessagesFailed = %d, want 1", stats.MessagesFailed)
	}
	if stats.DevicesActive != 2 {
		t.Errorf("DevicesActive = %d, want 2", stats.DevicesActive)
	}
	if stats.DevicesOnline != 1 {
		t.Errorf("DevicesOnline = %d, want 1", stats.DevicesOnline)
	}
	if stats.DevicesTotal != 3 {
		t.Errorf("DevicesTotal = %d, want 3", stats.DevicesTotal)
	}
}

func TestTrendsDefaultDays(t *testing.T) {
	upstream := newFakeUpstream()
	t.Cleanup(upstream.close)

	svc := newService(upstream)
	trends, err := svc.Trends(context.Background(), "user", "pass", 7)
	if err != nil {
		t.Fatalf("Trends() error: %v", err)
	}

	if trends.Days != 7 {
		t.Errorf("Days = %d, want 7", trends.Days)
	}
	if len(trends.MessageVolume) != 7 {
		t.Errorf("len(MessageVolume) = %d, want 7", len(trends.MessageVolume))
	}
	if len(trends.DeviceActivity) != 7 {
		t.Errorf("len(DeviceActivity) = %d, want 7", len(trends.DeviceActivity))
	}

	data, err := json.Marshal(trends)
	if err != nil {
		t.Fatalf("json.Marshal() error: %v", err)
	}
	for _, key := range []string{"days", "messageVolume", "deviceActivity", "sent", "pending", "failed", "active"} {
		if !bytes.Contains(data, []byte(key)) {
			t.Errorf("marshaled Trends missing key %q", key)
		}
	}
}

func TestTrendsCounts(t *testing.T) {
	upstream := newFakeUpstream()
	t.Cleanup(upstream.close)

	now := time.Now().UTC()
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	// today: 10 messages (7 sent, 2 pending, 1 failed); yesterday: 5 failed.
	msgs := make([]fakeMessage, 0, 15)
	for i := range 7 {
		msgs = append(msgs, fakeMessage{
			id:        fmt.Sprintf("today-sent-%d", i),
			deviceID:  "dev-a",
			state:     stateSent,
			createdAt: dayStart.Add(time.Duration(10+i) * time.Minute),
		})
	}
	for i := range 2 {
		msgs = append(msgs, fakeMessage{
			id:        fmt.Sprintf("today-pending-%d", i),
			deviceID:  "dev-b",
			state:     statePending,
			createdAt: dayStart.Add(time.Duration(11+i) * time.Hour),
		})
	}
	msgs = append(msgs, fakeMessage{
		id:        "today-failed",
		deviceID:  "dev-c",
		state:     stateFailed,
		createdAt: dayStart.Add(13 * time.Hour),
	})
	for i := range 5 {
		msgs = append(msgs, fakeMessage{
			id:        fmt.Sprintf("yesterday-failed-%d", i),
			deviceID:  "dev-d",
			state:     stateFailed,
			createdAt: dayStart.Add(-24*time.Hour + time.Duration(i)*time.Minute),
		})
	}
	upstream.messages = msgs

	svc := newService(upstream)
	trends, err := svc.Trends(context.Background(), "user", "pass", 7)
	if err != nil {
		t.Fatalf("Trends() error: %v", err)
	}

	today := bucketDate(now, 0)
	yesterday := bucketDate(now, -1)

	got := volumeByDate(t, trends.MessageVolume, today)
	if got.Sent != 7 || got.Pending != 2 || got.Failed != 1 {
		t.Errorf("today volume = %+v, want {Sent:7 Pending:2 Failed:1}", got)
	}

	got = volumeByDate(t, trends.MessageVolume, yesterday)
	if got.Sent != 0 || got.Pending != 0 || got.Failed != 5 {
		t.Errorf("yesterday volume = %+v, want {Sent:0 Pending:0 Failed:5}", got)
	}
}

func TestTrendsVolumeStateDerivation(t *testing.T) {
	upstream := newFakeUpstream()
	t.Cleanup(upstream.close)

	now := time.Now().UTC()
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	// Raw-string state rule: only Pending and Failed are split out;
	// Processed/Delivered/Cancelling/Cancelled all count as sent.
	upstream.messages = []fakeMessage{
		{id: "m-sent", deviceID: "dev-a", state: stateSent, createdAt: dayStart.Add(10 * time.Hour)},
		{id: "m-processed", deviceID: "dev-a", state: "Processed", createdAt: dayStart.Add(10*time.Hour + time.Minute)},
		{id: "m-delivered", deviceID: "dev-b", state: "Delivered", createdAt: dayStart.Add(11 * time.Hour)},
		{
			id:        "m-cancelling",
			deviceID:  "dev-b",
			state:     "Cancelling",
			createdAt: dayStart.Add(11*time.Hour + time.Minute),
		},
		{id: "m-cancelled", deviceID: "dev-c", state: "Cancelled", createdAt: dayStart.Add(12 * time.Hour)},
		{id: "m-pending", deviceID: "dev-c", state: statePending, createdAt: dayStart.Add(12*time.Hour + time.Minute)},
		{id: "m-failed", deviceID: "dev-d", state: stateFailed, createdAt: dayStart.Add(13 * time.Hour)},
	}

	svc := newService(upstream)
	trends, err := svc.Trends(context.Background(), "user", "pass", 7)
	if err != nil {
		t.Fatalf("Trends() error: %v", err)
	}

	got := volumeByDate(t, trends.MessageVolume, bucketDate(now, 0))
	if got.Sent != 5 || got.Pending != 1 || got.Failed != 1 {
		t.Errorf("today volume = %+v, want {Sent:5 Pending:1 Failed:1}", got)
	}
}

func TestTrendsVolumeSentClamp(t *testing.T) {
	upstream := newFakeUpstream()
	t.Cleanup(upstream.close)

	now := time.Now().UTC()
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	// A day with only failed messages: sent = max(5-0-5, 0) must clamp to 0,
	// never go negative.
	msgs := make([]fakeMessage, 0, 5)
	for i := range 5 {
		msgs = append(msgs, fakeMessage{
			id:        fmt.Sprintf("fail-%d", i),
			deviceID:  "dev-f",
			state:     stateFailed,
			createdAt: dayStart.Add(time.Duration(10+i) * time.Hour),
		})
	}
	upstream.messages = msgs

	svc := newService(upstream)
	trends, err := svc.Trends(context.Background(), "user", "pass", 7)
	if err != nil {
		t.Fatalf("Trends() error: %v", err)
	}

	got := volumeByDate(t, trends.MessageVolume, bucketDate(now, 0))
	if got.Sent != 0 || got.Pending != 0 || got.Failed != 5 {
		t.Errorf("today volume = %+v, want {Sent:0 Pending:0 Failed:5}", got)
	}
}

func TestTrendsDeviceActivity(t *testing.T) {
	upstream := newFakeUpstream()
	t.Cleanup(upstream.close)

	now := time.Now().UTC()
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	// dev-b sends 3 messages today: distinct devices per day, not per message.
	// dev-a is active today and two days ago and must be counted in both days.
	upstream.messages = []fakeMessage{
		{id: "m1", deviceID: "dev-a", state: stateSent, createdAt: dayStart.Add(10 * time.Hour)},
		{id: "m2", deviceID: "dev-b", state: stateSent, createdAt: dayStart.Add(11 * time.Hour)},
		{id: "m3", deviceID: "dev-b", state: stateSent, createdAt: dayStart.Add(12 * time.Hour)},
		{id: "m4", deviceID: "dev-b", state: stateSent, createdAt: dayStart.Add(13 * time.Hour)},
		{id: "m5", deviceID: "dev-a", state: stateSent, createdAt: dayStart.Add(-48 * time.Hour)},
		{id: "m6", deviceID: "dev-c", state: stateSent, createdAt: dayStart.Add(-47 * time.Hour)},
		{id: "m7", deviceID: "dev-outside", state: stateSent, createdAt: dayStart.Add(-10 * 24 * time.Hour)},
	}

	svc := newService(upstream)
	trends, err := svc.Trends(context.Background(), "user", "pass", 7)
	if err != nil {
		t.Fatalf("Trends() error: %v", err)
	}

	today := bucketDate(now, 0)
	yesterday := bucketDate(now, -1)
	twoDaysAgo := bucketDate(now, -2)

	if got := activityByDate(t, trends.DeviceActivity, today); got != 2 {
		t.Errorf("active today = %d, want 2 (dev-a, dev-b; dev-b messages counted once)", got)
	}
	if got := activityByDate(t, trends.DeviceActivity, twoDaysAgo); got != 2 {
		t.Errorf("active two days ago = %d, want 2 (dev-a, dev-c)", got)
	}
	if got := activityByDate(t, trends.DeviceActivity, yesterday); got != 0 {
		t.Errorf("active yesterday = %d, want 0 (no messages)", got)
	}
}

func TestTrendsDeviceActivityPagination(t *testing.T) {
	upstream := newFakeUpstream()
	t.Cleanup(upstream.close)

	now := time.Now().UTC()
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	// 250 messages across 3 devices today: more than one page of results.
	// Plus 2 devices with messages two days ago.
	msgs := make([]fakeMessage, 0, 252)
	for i := range 250 {
		msgs = append(msgs, fakeMessage{
			id:        fmt.Sprintf("m-%d", i),
			deviceID:  fmt.Sprintf("dev-%d", i%3),
			state:     stateSent,
			createdAt: dayStart.Add(time.Duration(10+i/100) * time.Hour),
		})
	}
	msgs = append(msgs,
		fakeMessage{id: "old-1", deviceID: "dev-old-a", state: stateSent, createdAt: dayStart.Add(-48 * time.Hour)},
		fakeMessage{id: "old-2", deviceID: "dev-old-b", state: stateSent, createdAt: dayStart.Add(-47 * time.Hour)},
	)
	upstream.messages = msgs

	svc := newService(upstream)
	trends, err := svc.Trends(context.Background(), "user", "pass", 7)
	if err != nil {
		t.Fatalf("Trends() error: %v", err)
	}

	today := bucketDate(now, 0)
	twoDaysAgo := bucketDate(now, -2)

	if got := activityByDate(t, trends.DeviceActivity, today); got != 3 {
		t.Errorf("active today = %d, want 3 distinct devices across 250 messages", got)
	}
	if got := activityByDate(t, trends.DeviceActivity, twoDaysAgo); got != 2 {
		t.Errorf("active two days ago = %d, want 2", got)
	}
}

func TestTrendsCache(t *testing.T) {
	upstream := newFakeUpstream()
	t.Cleanup(upstream.close)

	svc := newService(upstream)

	for range 2 {
		if _, err := svc.Trends(context.Background(), "user", "pass", 7); err != nil {
			t.Fatalf("Trends() error: %v", err)
		}
	}

	// Per bucket: 3 countMessagesRange requests + 1 activity page request.
	wantHits := int64(7)
	if got := upstream.msgHits.Load(); got != wantHits {
		t.Errorf("upstream message hits = %d, want %d (cached)", got, wantHits)
	}
}

func TestTrendsCacheIsolatedPerUser(t *testing.T) {
	upstream := newFakeUpstream()
	t.Cleanup(upstream.close)

	svc := newService(upstream)

	if _, err := svc.Trends(context.Background(), "user-a", "pass", 7); err != nil {
		t.Fatalf("Trends(user-a) error: %v", err)
	}
	if _, err := svc.Trends(context.Background(), "user-b", "pass", 7); err != nil {
		t.Fatalf("Trends(user-b) error: %v", err)
	}
	if _, err := svc.Trends(context.Background(), "user-a", "pass", 7); err != nil {
		t.Fatalf("Trends(user-a) error: %v", err)
	}

	// 2 cache misses x 7 buckets x 4 requests per bucket.
	wantHits := int64(2 * 7 * 1)
	if got := upstream.msgHits.Load(); got != wantHits {
		t.Errorf("upstream message hits = %d, want %d", got, wantHits)
	}
}

func TestTrendsCacheIsolatedPerDays(t *testing.T) {
	upstream := newFakeUpstream()
	t.Cleanup(upstream.close)

	svc := newService(upstream)

	if _, err := svc.Trends(context.Background(), "user", "pass", 7); err != nil {
		t.Fatalf("Trends(days=7) error: %v", err)
	}
	if _, err := svc.Trends(context.Background(), "user", "pass", 30); err != nil {
		t.Fatalf("Trends(days=30) error: %v", err)
	}

	wantHits := int64(7 + 30)
	if got := upstream.msgHits.Load(); got != wantHits {
		t.Errorf("upstream message hits = %d, want %d", got, wantHits)
	}
}

func TestTrendsUpstreamError(t *testing.T) {
	upstream := newFakeUpstream()
	t.Cleanup(upstream.close)

	upstream.failMsgs.Store(true)

	svc := newService(upstream)
	if _, err := svc.Trends(context.Background(), "user", "pass", 7); err == nil {
		t.Error("Trends() error = nil, want error")
	}
	if _, err := svc.Stats(context.Background(), "user", "pass"); err == nil {
		t.Error("Stats() error = nil, want error")
	}
}

func volumeByDate(t *testing.T, volumes []dashboard.DayVolume, date string) dashboard.DayVolume {
	t.Helper()
	for _, v := range volumes {
		if v.Date == date {
			return v
		}
	}
	t.Fatalf("no volume for date %q", date)

	return dashboard.DayVolume{}
}

func activityByDate(t *testing.T, activity []dashboard.DayActivity, date string) int {
	t.Helper()
	for _, a := range activity {
		if a.Date == date {
			return a.Active
		}
	}
	t.Fatalf("no activity for date %q", date)

	return 0
}
