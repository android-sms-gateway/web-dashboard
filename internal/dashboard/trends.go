package dashboard

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/android-sms-gateway/client-go/smsgateway"
)

const dateLayout = "2006-01-02"

// Stats represents aggregated dashboard statistics.
type Stats struct {
	MessagesSent    int `json:"messagesSent"`
	MessagesPending int `json:"messagesPending"`
	MessagesFailed  int `json:"messagesFailed"`
	DevicesActive   int `json:"devicesActive"`
	DevicesOnline   int `json:"devicesOnline"`
	DevicesTotal    int `json:"devicesTotal"`
}

// Trends represents per-day message volume and device activity.
type Trends struct {
	Days           int           `json:"days"`
	MessageVolume  []DayVolume   `json:"messageVolume"`
	DeviceActivity []DayActivity `json:"deviceActivity"`
}

// DayVolume represents message counts for a single day.
type DayVolume struct {
	Date    string `json:"date"`
	Sent    int    `json:"sent"`
	Pending int    `json:"pending"`
	Failed  int    `json:"failed"`
}

// DayActivity represents the number of distinct devices that sent messages on a single day.
type DayActivity struct {
	Date   string `json:"date"`
	Active int    `json:"active"`
}

const (
	trendsConcurrency = 6
	// activityPageSize is the page size used when listing messages per day
	// bucket to collect distinct device IDs. The upstream server accepts up
	// to 100.
	activityPageSize = 100
	// DeviceOnlineIn is the window within which a device is considered online.
	DeviceOnlineIn = 15 * time.Minute
)

type trendsBucket struct {
	start time.Time
	end   time.Time
	date  string
}

func trendsBuckets(now time.Time, days int) []trendsBucket {
	if days <= 0 {
		return nil
	}

	buckets := make([]trendsBucket, days)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	for i := range days {
		start := today.AddDate(0, 0, i-days+1)
		buckets[i] = trendsBucket{
			start: start,
			end:   start.AddDate(0, 0, 1),
			date:  start.Format(dateLayout),
		}
	}

	return buckets
}

func countMessagesRange(
	ctx context.Context,
	client *smsgateway.Client,
	from, to *time.Time,
) (int, int, int, error) {
	var wg sync.WaitGroup
	var t, p, f int
	var collectErr error
	var errMu sync.Mutex

	limit := 1
	pendingState := smsgateway.ProcessingStatePending
	failedState := smsgateway.ProcessingStateFailed

	recordErr := func(e error) {
		errMu.Lock()
		collectErr = errors.Join(collectErr, e)
		errMu.Unlock()
	}

	query := func(state *string, out *int) {
		_, n, err := client.ListMessages(ctx, smsgateway.ListMessagesOptions{
			State:          state,
			Limit:          &limit,
			From:           from,
			To:             to,
			DeviceID:       nil,
			Offset:         nil,
			IncludeContent: nil,
		})
		if err != nil {
			recordErr(err)
			return
		}
		*out = n
	}

	wg.Go(func() { query(nil, &t) })
	wg.Go(func() { query((*string)(&pendingState), &p) })
	wg.Go(func() { query((*string)(&failedState), &f) })

	wg.Wait()

	return t, p, f, collectErr
}

func sweepDays(
	ctx context.Context,
	client *smsgateway.Client,
	buckets []trendsBucket,
) ([]DayVolume, []DayActivity, error) {
	volumes := make([]DayVolume, len(buckets))
	activity := make([]DayActivity, len(buckets))
	for i, b := range buckets {
		volumes[i] = DayVolume{
			Date:    b.date,
			Sent:    0,
			Pending: 0,
			Failed:  0,
		}
		activity[i] = DayActivity{Date: b.date, Active: 0}
	}

	if len(buckets) == 0 {
		return volumes, activity, nil
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, trendsConcurrency)
	var collectErr error
	var errMu sync.Mutex

	recordErr := func(e error) {
		errMu.Lock()
		collectErr = errors.Join(collectErr, e)
		errMu.Unlock()
	}

	for i, b := range buckets {
		wg.Go(func() {
			sem <- struct{}{}
			defer func() { <-sem }()

			from, to := b.start, b.end
			stats, err := sweepDay(ctx, client, from, to)
			if err != nil {
				recordErr(err)
				return
			}

			volumes[i] = DayVolume{
				Date:    b.date,
				Sent:    max(stats.total-stats.pending-stats.failed, 0),
				Pending: stats.pending,
				Failed:  stats.failed,
			}
			activity[i] = DayActivity{Date: b.date, Active: stats.devices}
		})
	}

	wg.Wait()

	return volumes, activity, collectErr
}

// daySweepStats aggregates volume counts and distinct device activity for one
// day bucket, collected from the messages already fetched by the sweep.
type daySweepStats struct {
	total   int
	pending int
	failed  int
	devices int
}

// sweepDay returns per-state volume counts and the number of distinct device
// IDs that sent messages in [from, to). Volume counts come from
// countMessagesRange (O(1) per state); device IDs are collected by paging
// through messages.
func sweepDay(
	ctx context.Context,
	client *smsgateway.Client,
	from, to time.Time,
) (daySweepStats, error) {
	total, pending, failed := 0, 0, 0
	devices := make(map[string]struct{})
	limit := activityPageSize
	offset := 0

	for {
		msgs, _, listErr := client.ListMessages(ctx, smsgateway.ListMessagesOptions{
			From:           &from,
			To:             &to,
			State:          nil,
			DeviceID:       nil,
			Limit:          &limit,
			Offset:         &offset,
			IncludeContent: nil,
		})
		if listErr != nil {
			return daySweepStats{}, fmt.Errorf("failed to list messages: %w", listErr)
		}

		for _, m := range msgs {
			if m.DeviceID != "" {
				devices[m.DeviceID] = struct{}{}
			}

			switch m.State {
			case smsgateway.ProcessingStatePending:
				pending++
			case smsgateway.ProcessingStateFailed:
				failed++
			case smsgateway.ProcessingStateProcessed,
				smsgateway.ProcessingStateSent,
				smsgateway.ProcessingStateDelivered:
			}
		}

		total += len(msgs)

		offset += len(msgs)
		if len(msgs) < limit {
			break
		}
	}

	return daySweepStats{
		total:   total,
		pending: pending,
		failed:  failed,
		devices: len(devices),
	}, nil
}
