package anomaly

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"
	"time"
)

// recoveredEvent holds a complete interval observed between periodic scans.
type recoveredEvent struct {
	event      Event
	endedAt    time.Time
	lastValue  *string
	inProgress bool
}

// The caller holds d.mu. Pending closed intervals use the same per-device
// bound as persisted history; an interval being written is never discarded.
func (d *Detector) queueRecoveredLocked(event Event, endedAt time.Time, lastValue *string) {
	queue := d.recovered[event.DeviceID]
	if len(queue) >= DefaultRetention {
		oldest := -1
		for i, item := range queue {
			if item.inProgress || item.event.ID != 0 {
				continue
			}
			if oldest < 0 || item.event.StartedAt.Before(queue[oldest].event.StartedAt) {
				oldest = i
			}
		}
		if oldest < 0 {
			return
		}
		queue = slices.Delete(queue, oldest, oldest+1)
	}
	d.recovered[event.DeviceID] = append(queue, &recoveredEvent{event: event, endedAt: endedAt, lastValue: lastValue})
}

// A recovery observed before Load can already be covered by an adopted event.
func (d *Detector) discardCoveredRecoveredLocked(event Event, recoveredAt time.Time) {
	if recoveredAt.IsZero() {
		return
	}
	queue := slices.DeleteFunc(d.recovered[event.DeviceID], func(item *recoveredEvent) bool {
		return !item.inProgress && item.event.ID == 0 && item.event.Kind == event.Kind && item.event.PointID == event.PointID &&
			!item.event.StartedAt.Before(event.StartedAt) && !item.endedAt.After(recoveredAt)
	})
	if len(queue) == 0 {
		delete(d.recovered, event.DeviceID)
	} else {
		d.recovered[event.DeviceID] = queue
	}
}

// flushRecovered persists complete intervals before opening newer active ones.
// A partially persisted event stays queued so its end can be retried.
func (d *Detector) flushRecovered(ctx context.Context) error {
	d.mu.Lock()
	devices := slices.Sorted(maps.Keys(d.recovered))
	d.mu.Unlock()
	var failures []error
	for _, deviceID := range devices {
		for range DefaultRetention {
			d.mu.Lock()
			queue := d.recovered[deviceID]
			if len(queue) == 0 {
				d.mu.Unlock()
				break
			}
			entry := queue[0]
			entry.inProgress = true
			d.mu.Unlock()

			var err error
			if entry.event.ID == 0 {
				var persisted Event
				persisted, err = d.repo.Open(ctx, entry.event)
				if err == nil {
					d.mu.Lock()
					entry.event = persisted
					d.mu.Unlock()
				}
			}
			if err == nil {
				err = d.repo.End(ctx, entry.event.ID, entry.endedAt, entry.lastValue, "")
			}

			d.mu.Lock()
			entry.inProgress = false
			if err == nil || errors.Is(err, ErrNotActive) {
				queue = d.recovered[deviceID]
				if len(queue) > 0 && queue[0] == entry {
					queue = queue[1:]
					if len(queue) == 0 {
						delete(d.recovered, deviceID)
					} else {
						d.recovered[deviceID] = queue
					}
				}
			}
			d.mu.Unlock()
			if err != nil {
				failures = append(failures, fmt.Errorf("record recovered %s anomaly: %w", entry.event.Kind, err))
				break
			}
		}
	}
	return errors.Join(failures...)
}
