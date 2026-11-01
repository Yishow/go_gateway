package main

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
)

// Bounded per-device/Point state; never retain 100*devices*Points envelopes.
// Both fixture channels are continuously drained, so their fixed 256 capacity
// cannot itself block the production ingestor or manufacture output loss.
type acceptanceSampleDrain struct {
	mu           sync.Mutex
	ids          map[string]string
	seen         map[string]map[string]bool
	want         map[string]map[string]int16
	lastReceived map[string]time.Time
	limit        int
	bad          error
	total        uint64
	stopping     bool
	delivered    map[string]map[string]int
}

func startAcceptanceSampleDrain(t *testing.T, capture *acquisitionWiringCapture, load []acceptanceScheduledDevice, limit int) *acceptanceSampleDrain {
	t.Helper()
	d := &acceptanceSampleDrain{ids: map[string]string{}, seen: map[string]map[string]bool{}, want: map[string]map[string]int16{}, lastReceived: map[string]time.Time{}, limit: limit, delivered: map[string]map[string]int{}}
	for _, device := range load {
		d.want[device.id] = device.want
		d.delivered[device.id] = map[string]int{}
	}
	// Observe synchronously at the existing test SampleSink decorator so queued
	// fixture-channel items cannot masquerade as product ingestion backlog.
	capture.acceptanceSampleObserver = d.accept
	ctx, cancel := context.WithCancel(t.Context())
	finished := make(chan struct{})
	go func() {
		defer close(finished)
		for {
			select {
			case <-capture.samples:
			case <-capture.observations:
			case <-ctx.Done():
				return
			}
		}
	}()
	t.Cleanup(func() { cancel(); <-finished })
	return d
}

func (d *acceptanceSampleDrain) accept(sample measurement.SampleEnvelope) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.stopping {
		return
	}
	d.total++
	want, known := d.want[sample.DeviceID][sample.PointID]
	if !known || sample.Value != want || sample.Quality != schema.QualityGood || sample.AcquisitionID == "" || sample.ConfigFingerprint == "" || sample.SourceRevision == "" || sample.ReceivedAt.IsZero() || sample.ObservedAt != sample.ReceivedAt || sample.ReceivedAt.Before(d.lastReceived[sample.PointID]) {
		d.bad = fmt.Errorf("invalid production sample for device=%s point=%s quality=%s value=%v", sample.DeviceID, sample.PointID, sample.Quality, sample.Value)
		return
	}
	d.lastReceived[sample.PointID] = sample.ReceivedAt
	if d.ids[sample.DeviceID] != sample.AcquisitionID {

		d.ids[sample.DeviceID] = sample.AcquisitionID
		d.seen[sample.DeviceID] = map[string]bool{}
	}
	if d.seen[sample.DeviceID][sample.PointID] {
		d.bad = fmt.Errorf("duplicate production sample %s/%s", sample.AcquisitionID, sample.PointID)
	}
	d.seen[sample.DeviceID][sample.PointID] = true
	byScan := d.delivered[sample.DeviceID]
	if _, exists := byScan[sample.AcquisitionID]; exists || len(byScan) < d.limit {
		byScan[sample.AcquisitionID]++
	}
}

func (d *acceptanceSampleDrain) stopValidation() {
	d.mu.Lock()
	d.stopping = true
	d.mu.Unlock()
}

func (d *acceptanceSampleDrain) count() uint64 { d.mu.Lock(); defer d.mu.Unlock(); return d.total }
