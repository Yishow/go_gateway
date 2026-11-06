package anomaly

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"time"
)

type scanAction struct {
	key       eventKey
	entry     *activeEvent
	open      bool
	endAt     time.Time
	detail    string
	lastValue *string
}

// Scan compares the latest observations with the thresholds and persists event
// starts and ends. Repository failures are returned and retried on the next
// scan with the same start time; nothing is reported as recorded unless the
// repository accepted it.
func (d *Detector) Scan(ctx context.Context, input ScanInput) error {
	d.scanMu.Lock()
	defer d.scanMu.Unlock()
	var errs []error
	d.mu.Lock()
	reload := d.needReload
	d.mu.Unlock()
	if reload {
		if err := d.Load(ctx); err != nil {
			errs = append(errs, err)
		}
	}
	actions := d.plan(d.now(), input)
	for _, a := range actions {
		if !a.open {
			errs = append(errs, d.apply(ctx, a))
		}
	}
	errs = append(errs, d.flushRecovered(ctx))
	for _, a := range actions {
		if a.open {
			errs = append(errs, d.apply(ctx, a))
		}
	}
	return errors.Join(errs...)
}

func (d *Detector) plan(now time.Time, input ScanInput) []scanAction {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.thresholds = maps.Clone(input.Thresholds)
	var actions []scanAction
	end := func(key eventKey, ae *activeEvent, at time.Time, detail string, value *string) {
		if ae.event.ID > 0 && !ae.ending {
			ae.ending = true
			actions = append(actions, scanAction{key: key, entry: ae, endAt: at, detail: detail, lastValue: value})
		}
	}
	open := func(ev Event) {
		key := eventKey{ev.DeviceID, ev.PointID, ev.Kind}
		ae := &activeEvent{event: ev}
		d.active[key] = ae
		actions = append(actions, scanAction{key: key, entry: ae, open: true})
	}

	monitoredDevices := map[string]bool{}
	for pointID, deviceID := range input.Monitored {
		if deviceID == "" || pointID == "" {
			continue
		}
		monitoredDevices[deviceID] = true
		ps := d.point(pointID, deviceID, now)
		th := input.Thresholds[deviceID]
		reference := ps.firstSeen
		if !ps.lastSuccess.IsZero() {
			reference = ps.lastSuccess
		}
		staleKey := eventKey{deviceID, pointID, KindStale}
		if ae := d.active[staleKey]; ae != nil {
			if !ae.recoveredAt.IsZero() {
				end(staleKey, ae, ae.recoveredAt, "", ps.valueText)
			}
		} else if th.Stale > 0 && now.Sub(reference) > th.Stale {
			open(Event{DeviceID: deviceID, PointID: pointID, Kind: KindStale, StartedAt: reference.Add(th.Stale), LastValue: ps.valueText, Detail: DetailNoSuccessfulRead})
		}
		frozenKey := eventKey{deviceID, pointID, KindFrozen}
		fresh := th.Stale <= 0 || now.Sub(reference) <= th.Stale
		if ae := d.active[frozenKey]; ae != nil {
			switch {
			case !ae.recoveredAt.IsZero():
				end(frozenKey, ae, ae.recoveredAt, "", ps.valueText)
			case th.Frozen <= 0:
				end(frozenKey, ae, now, DetailFrozenDisabled, ps.valueText)
			}
		} else if th.Frozen > 0 && ps.hasValue && fresh && now.Sub(ps.lastChange) > th.Frozen {
			open(Event{DeviceID: deviceID, PointID: pointID, Kind: KindFrozen, StartedAt: ps.lastChange.Add(th.Frozen), LastValue: ps.valueText, Detail: DetailValueUnchanged})
		}
	}

	for deviceID := range monitoredDevices {
		ds := d.device(deviceID)
		breakerOpen := input.BreakerOpen[deviceID]
		switch {
		case !breakerOpen:
			ds.breakerSince = time.Time{}
		case ds.breakerSince.IsZero():
			ds.breakerSince = now
		}
		key := eventKey{deviceID, "", KindDisconnected}
		if ae := d.active[key]; ae != nil {
			if !ae.recoveredAt.IsZero() && !breakerOpen {
				end(key, ae, ae.recoveredAt, "", nil)
			}
			continue
		}
		start, detail := ds.failingSince, ds.failReason
		if !ds.breakerSince.IsZero() && (start.IsZero() || ds.breakerSince.Before(start)) {
			start, detail = ds.breakerSince, ReasonBreakerOpen
		}
		if !start.IsZero() {
			open(Event{DeviceID: deviceID, Kind: KindDisconnected, StartedAt: start, Detail: detail})
		}
	}

	for key, ae := range d.active {
		stillMonitored := monitoredDevices[key.deviceID]
		if key.pointID != "" {
			stillMonitored = input.Monitored[key.pointID] == key.deviceID
		}
		if !stillMonitored {
			end(key, ae, now, DetailMonitoringStopped, nil)
		}
	}
	for pointID, ps := range d.points {
		if input.Monitored[pointID] != ps.deviceID {
			delete(d.points, pointID)
		}
	}
	for deviceID := range d.devices {
		if !monitoredDevices[deviceID] {
			delete(d.devices, deviceID)
		}
	}
	return actions
}

func (d *Detector) apply(ctx context.Context, a scanAction) error {
	if a.open {
		persisted, err := d.repo.Open(ctx, a.entry.event)
		d.mu.Lock()
		defer d.mu.Unlock()
		if err != nil {
			if d.active[a.key] == a.entry {
				delete(d.active, a.key)
			}
			if errors.Is(err, ErrAlreadyActive) {
				d.needReload = true
			}
			return fmt.Errorf("record %s anomaly start: %w", a.key.kind, err)
		}
		a.entry.event = persisted
		return nil
	}
	err := d.repo.End(ctx, a.entry.event.ID, a.endAt, a.lastValue, a.detail)
	d.mu.Lock()
	defer d.mu.Unlock()
	if err != nil && !errors.Is(err, ErrNotActive) {
		a.entry.ending = false
		return fmt.Errorf("record %s anomaly end: %w", a.key.kind, err)
	}
	if d.active[a.key] == a.entry {
		delete(d.active, a.key)
	}
	if err != nil {
		return fmt.Errorf("record %s anomaly end: %w", a.key.kind, err)
	}
	return nil
}
