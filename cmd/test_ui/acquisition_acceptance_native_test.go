package main

import (
	"fmt"
	"sync"
	"time"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/schema"
)

// Bounded synchronous test-only observer: no lossy channel, no runtime history
// registry. Quality/identity exist here; decoded values are NOT exposed by the
// native observer and are checked only at SampleSink and in paired adapter tests.
type acceptanceNativeCapture struct {
	mu           sync.Mutex
	want         map[string]map[string]int16
	current      map[string]string
	seen         map[string]map[string]bool
	completed    map[string][]string
	limit        int
	bad          error
	total        uint64
	lastReceived map[string]time.Time
}

func newAcceptanceNativeCapture(load []acceptanceScheduledDevice, limit int) *acceptanceNativeCapture {
	n := &acceptanceNativeCapture{want: map[string]map[string]int16{}, current: map[string]string{}, seen: map[string]map[string]bool{}, completed: map[string][]string{}, limit: limit, lastReceived: map[string]time.Time{}}
	for _, d := range load {
		n.want[d.id] = d.want
	}
	return n
}

func (n *acceptanceNativeCapture) accept(o collector.AcquisitionObservation) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.total++
	if len(n.completed[o.DeviceID]) >= n.limit && n.current[o.DeviceID] != o.AcquisitionID {
		return
	}
	_, known := n.want[o.DeviceID][o.PointID]
	if !known || o.Quality != schema.QualityGood || o.AcquisitionID == "" || o.RequestID == "" || o.ConfigFingerprint == "" || o.PlanFingerprint == "" || o.ReceivedAt.IsZero() || o.ReceivedAt.Before(n.lastReceived[o.PointID]) {
		n.bad = fmt.Errorf("invalid native quality/identity device=%s point=%s", o.DeviceID, o.PointID)
		return
	}
	n.lastReceived[o.PointID] = o.ReceivedAt
	if n.current[o.DeviceID] != o.AcquisitionID {
		if prior := n.seen[o.DeviceID]; len(prior) > 0 && len(prior) != len(n.want[o.DeviceID]) {
			n.bad = fmt.Errorf("incomplete native scan device=%s: %d/%d", o.DeviceID, len(prior), len(n.want[o.DeviceID]))
		}
		n.current[o.DeviceID] = o.AcquisitionID
		n.seen[o.DeviceID] = map[string]bool{}
	}
	if n.seen[o.DeviceID][o.PointID] {
		n.bad = fmt.Errorf("duplicate native point %s/%s", o.AcquisitionID, o.PointID)
	}
	n.seen[o.DeviceID][o.PointID] = true
	if len(n.seen[o.DeviceID]) == len(n.want[o.DeviceID]) && len(n.completed[o.DeviceID]) < n.limit {
		n.completed[o.DeviceID] = append(n.completed[o.DeviceID], o.AcquisitionID)
	}
}

func (n *acceptanceNativeCapture) next(device string, ordinal int) (acquisitionID string, available bool, err error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	ids := n.completed[device]
	if ordinal >= len(ids) {
		return "", false, n.bad
	}
	return ids[ordinal], true, n.bad
}

// Detach bounded per-scan delivery counts at the same cut as the global total.
func (d *acceptanceSampleDrain) deliverySnapshot() (total uint64, deliveries map[string]map[string]int, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	result := make(map[string]map[string]int, len(d.delivered))
	for device, scans := range d.delivered {
		result[device] = make(map[string]int, len(scans))
		for id, count := range scans {
			result[device][id] = count
		}
	}
	return d.total, result, d.bad
}

func (n *acceptanceNativeCapture) count() uint64 { n.mu.Lock(); defer n.mu.Unlock(); return n.total }

func (n *acceptanceNativeCapture) validationError() error {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.bad
}
