package anomaly

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"
)

const (
	defaultStaleMultiplier = 5
	minimumStaleThreshold  = 10 * time.Second
	maxLastValueLength     = 256
)

// StaleThreshold resolves a device's stale threshold: the configured value, or
// five times the acquisition interval and at least 10 seconds.
func StaleThreshold(configuredMs *int, intervalMs int) time.Duration {
	if configuredMs != nil && *configuredMs > 0 {
		return time.Duration(*configuredMs) * time.Millisecond
	}
	return max(time.Duration(intervalMs)*time.Millisecond*defaultStaleMultiplier, minimumStaleThreshold)
}

// FrozenThreshold resolves a device's frozen threshold; zero disables it.
func FrozenThreshold(configuredMs int) time.Duration {
	if configuredMs <= 0 {
		return 0
	}
	return time.Duration(configuredMs) * time.Millisecond
}

// Thresholds are one device's resolved detection thresholds. A zero Stale
// skips stale evaluation; a zero Frozen disables frozen detection.
type Thresholds struct {
	Stale  time.Duration
	Frozen time.Duration
}

// Observation is one acquisition result already received by the runtime.
type Observation struct {
	DeviceID string
	PointID  string
	At       time.Time
	OK       bool
	Value    any
	// Reason is the collector quality reason of a failed read.
	Reason string
}

// ScanInput is the runtime's current view used by one scan.
type ScanInput struct {
	// Monitored maps currently scheduled point IDs to their device ID.
	Monitored map[string]string
	// Thresholds maps device ID to its resolved thresholds.
	Thresholds map[string]Thresholds
	// BreakerOpen marks devices whose circuit breaker is open.
	BreakerOpen map[string]bool
}

type eventKey struct {
	deviceID, pointID string
	kind              Kind
}

// activeEvent is an event that is persisted (ID > 0) or being opened (ID 0).
// recoveredAt is the first clearing observation since it became active.
type activeEvent struct {
	event       Event
	recoveredAt time.Time
	ending      bool
}

type pointState struct {
	deviceID    string
	firstSeen   time.Time
	lastSuccess time.Time
	lastChange  time.Time
	hasValue    bool
	value       any
	valueText   *string
}

type deviceState struct {
	failingSince time.Time
	failReason   string
	breakerSince time.Time
	lastSuccess  time.Time
}

// Detector keeps the latest observation per point in memory and turns it into
// persisted anomaly events on Scan. Observe never performs I/O, so it is safe
// on the acquisition path; Scan is the only writer to the repository.
type Detector struct {
	repo Repository
	now  func() time.Time

	scanMu     sync.Mutex
	mu         sync.Mutex
	points     map[string]*pointState
	devices    map[string]*deviceState
	active     map[eventKey]*activeEvent
	thresholds map[string]Thresholds
	recovered  map[string][]*recoveredEvent
	needReload bool
}

// NewDetector creates a detector. now is injectable for tests.
func NewDetector(repo Repository, now func() time.Time) *Detector {
	if now == nil {
		now = time.Now
	}
	return &Detector{repo: repo, now: now, points: map[string]*pointState{}, devices: map[string]*deviceState{}, active: map[eventKey]*activeEvent{}, recovered: map[string][]*recoveredEvent{}}
}

// Load adopts events left active by a previous run so restarts never create a
// second active event for the same device, point and kind.
func (d *Detector) Load(ctx context.Context) error {
	events, err := d.repo.ListActive(ctx)
	if err != nil {
		return fmt.Errorf("load active acquisition anomalies: %w", err)
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, ev := range events {
		key := eventKey{ev.DeviceID, ev.PointID, ev.Kind}
		if _, known := d.active[key]; !known {
			entry := &activeEvent{event: ev}
			if ev.Kind == KindDisconnected {
				if ds := d.devices[ev.DeviceID]; ds != nil && ds.lastSuccess.After(ev.StartedAt) {
					entry.recoveredAt = ds.lastSuccess
				}
			} else if ps := d.points[ev.PointID]; ps != nil && ps.deviceID == ev.DeviceID && ps.lastSuccess.After(ev.StartedAt) {
				if ev.Kind == KindStale || (ev.Kind == KindFrozen && frozenValueChanged(ev.LastValue, ps.valueText)) {
					entry.recoveredAt = ps.lastSuccess
				}
			}
			d.discardCoveredRecoveredLocked(ev, entry.recoveredAt)
			d.active[key] = entry
		}
	}
	d.needReload = false
	return nil
}

// Observe records one acquisition result. It only updates memory.
func (d *Detector) Observe(o Observation) {
	if o.DeviceID == "" || o.PointID == "" {
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	if o.At.IsZero() {
		o.At = d.now()
	}
	ps := d.point(o.PointID, o.DeviceID, o.At)
	ds := d.device(o.DeviceID)
	if !o.OK {
		if (o.Reason == ReasonConnectionFailed || o.Reason == ReasonCircuitOpen) && ds.failingSince.IsZero() {
			ds.failingSince, ds.failReason = o.At, o.Reason
		}
		return
	}
	valueText := formatValue(o.Value)
	disconnectedKey := eventKey{o.DeviceID, "", KindDisconnected}
	if active := d.active[disconnectedKey]; !ds.failingSince.IsZero() && o.At.After(ds.failingSince) && (active == nil || !active.recoveredAt.IsZero()) {
		d.queueRecoveredLocked(Event{DeviceID: o.DeviceID, Kind: KindDisconnected, StartedAt: ds.failingSince, Detail: ds.failReason}, o.At, nil)
	}
	ds.failingSince, ds.failReason = time.Time{}, ""
	if o.At.After(ds.lastSuccess) {
		ds.lastSuccess = o.At
	}
	d.markRecovered(disconnectedKey, o.At)
	staleKey := eventKey{o.DeviceID, o.PointID, KindStale}
	threshold := d.thresholds[o.DeviceID].Stale
	reference := ps.firstSeen
	if !ps.lastSuccess.IsZero() {
		reference = ps.lastSuccess
	}
	if active := d.active[staleKey]; threshold > 0 && o.At.Sub(reference) > threshold && (active == nil || !active.recoveredAt.IsZero()) {
		d.queueRecoveredLocked(Event{DeviceID: o.DeviceID, PointID: o.PointID, Kind: KindStale, StartedAt: reference.Add(threshold), LastValue: ps.valueText, Detail: DetailNoSuccessfulRead}, o.At, valueText)
	}
	d.markRecovered(staleKey, o.At)
	if !ps.hasValue || !reflect.DeepEqual(ps.value, o.Value) {
		frozenKey := eventKey{o.DeviceID, o.PointID, KindFrozen}
		frozen := d.active[frozenKey]
		thresholds := d.thresholds[o.DeviceID]
		fresh := thresholds.Stale <= 0 || o.At.Sub(ps.lastSuccess) <= thresholds.Stale
		if ps.hasValue && thresholds.Frozen > 0 && fresh && o.At.Sub(ps.lastChange) > thresholds.Frozen && (frozen == nil || !frozen.recoveredAt.IsZero()) {
			d.queueRecoveredLocked(Event{DeviceID: o.DeviceID, PointID: o.PointID, Kind: KindFrozen, StartedAt: ps.lastChange.Add(thresholds.Frozen), LastValue: ps.valueText, Detail: DetailValueUnchanged}, o.At, valueText)
		}
		if ps.hasValue || (frozen != nil && frozenValueChanged(frozen.event.LastValue, valueText)) {
			d.markRecovered(frozenKey, o.At)
		}
		ps.lastChange = o.At
	}
	ps.hasValue, ps.value, ps.valueText = true, o.Value, valueText
	if o.At.After(ps.lastSuccess) {
		ps.lastSuccess = o.At
	}
}

func frozenValueChanged(previous, current *string) bool {
	return previous != nil && current != nil && *previous != *current
}

func (d *Detector) point(pointID, deviceID string, at time.Time) *pointState {
	ps := d.points[pointID]
	if ps == nil || ps.deviceID != deviceID {
		ps = &pointState{deviceID: deviceID, firstSeen: at}
		d.points[pointID] = ps
	}
	return ps
}

func (d *Detector) device(deviceID string) *deviceState {
	ds := d.devices[deviceID]
	if ds == nil {
		ds = &deviceState{}
		d.devices[deviceID] = ds
	}
	return ds
}

func (d *Detector) markRecovered(key eventKey, at time.Time) {
	if ae := d.active[key]; ae != nil && ae.recoveredAt.IsZero() {
		ae.recoveredAt = at
	}
}

func formatValue(v any) *string {
	if v == nil {
		return nil
	}
	var text string
	if raw, err := json.Marshal(v); err == nil {
		text = string(raw)
	} else {
		text = fmt.Sprint(v)
	}
	if len(text) > maxLastValueLength {
		text = strings.ToValidUTF8(text[:maxLastValueLength], "")
	}
	return &text
}
