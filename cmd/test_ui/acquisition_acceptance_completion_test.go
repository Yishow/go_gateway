package main

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/collector"

	"github.com/stretchr/testify/require"
)

// Test-only storage: fixed first-N completion events, never a product registry.
// Callback delivery order is not scan order; matching uses native AcquisitionID.
// First-N admission is callback order. If overtaking at the bound omits an
// expected native ID, the strict join fails; later facts never replace it.
type acceptanceCompletionSlot struct {
	events                   chan collector.AcquisitionCompletionObservation
	attempted, lost, surplus atomic.Uint64
	pending                  map[string]collector.AcquisitionCompletionObservation // consumer-owned
	seen                     map[string]bool                                       // consumer-owned, bounded by limit
}
type acceptanceCompletionCapture struct {
	slots         map[string]*acceptanceCompletionSlot // immutable after construction
	limit         uint64
	unknown       atomic.Uint64
	callbackNanos atomic.Int64
	callbackCalls atomic.Uint64
}

func newAcceptanceCompletionCapture(load []acceptanceScheduledDevice, limit int) *acceptanceCompletionCapture {
	c := &acceptanceCompletionCapture{slots: map[string]*acceptanceCompletionSlot{}, limit: uint64(limit)}
	for _, d := range load {
		c.slots[d.id] = &acceptanceCompletionSlot{events: make(chan collector.AcquisitionCompletionObservation, limit), pending: map[string]collector.AcquisitionCompletionObservation{}, seen: map[string]bool{}}
	}
	return c
}
func (c *acceptanceCompletionCapture) accept(o collector.AcquisitionCompletionObservation) {
	started := time.Now()
	defer func() { c.callbackNanos.Add(int64(time.Since(started))); c.callbackCalls.Add(1) }()
	slot := c.slots[o.DeviceID]
	if slot == nil {
		c.unknown.Add(1)
		return
	}
	if slot.attempted.Add(1) > c.limit {
		slot.surplus.Add(1)
		return
	}
	select {
	case slot.events <- o:
	default:
		slot.lost.Add(1)
	}
}
func (c *acceptanceCompletionCapture) next(device, id string) (collector.AcquisitionCompletionObservation, bool, error) {
	var zero collector.AcquisitionCompletionObservation
	if c.unknown.Load() != 0 {
		return zero, false, fmt.Errorf("unknown completion device")
	}
	slot := c.slots[device]
	if slot == nil {
		return zero, false, fmt.Errorf("unknown requested device %s", device)
	}
	if slot.lost.Load() != 0 {
		return zero, false, fmt.Errorf("completion observer lost first required events for %s: %d", device, slot.lost.Load())
	}
	for {
		select {
		case event := <-slot.events:
			key := event.Metrics.AcquisitionID
			if key == "" || slot.seen[key] {
				return zero, false, fmt.Errorf("empty/duplicate completion ID for %s", device)
			}
			slot.seen[key] = true
			slot.pending[key] = event
		default:
			event, ok := slot.pending[id]
			if ok {
				delete(slot.pending, id)
			}
			return event, ok, nil
		}
	}
}
func (c *acceptanceCompletionCapture) counts() (lost, surplus uint64) {
	for _, slot := range c.slots {
		lost += slot.lost.Load()
		surplus += slot.surplus.Load()
	}
	return lost, surplus
}

func TestAcceptanceCompletionBoundedIdentityJoin(t *testing.T) {
	c := newAcceptanceCompletionCapture([]acceptanceScheduledDevice{{id: "d"}}, 101)
	event := func(id string) collector.AcquisitionCompletionObservation {
		return collector.AcquisitionCompletionObservation{DeviceID: "d", Metrics: collector.AcquisitionActualMetrics{AcquisitionID: id}}
	}
	// Even reversed publication callback delivery cannot substitute native order.
	c.accept(event("second"))
	c.accept(event("first"))
	first, ok, err := c.next("d", "first")
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, "first", first.Metrics.AcquisitionID)
	second, ok, err := c.next("d", "second")
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, "second", second.Metrics.AcquisitionID)
	for i := 2; i < 101; i++ {
		c.accept(event(fmt.Sprintf("scan-%d", i)))
	}
	c.accept(event("surplus-not-a-replacement"))
	_, ok, err = c.next("d", "missing-first-required")
	require.NoError(t, err)
	require.False(t, ok)
	loss, surplus := c.counts()
	require.Zero(t, loss)
	require.EqualValues(t, 1, surplus)
	require.Len(t, c.slots["d"].seen, 101)
	require.Len(t, c.slots["d"].pending, 99)
	require.Equal(t, 101, cap(c.slots["d"].events))
}
func TestAcceptanceCompletionLossIsNotProductDrop(t *testing.T) {
	c := newAcceptanceCompletionCapture([]acceptanceScheduledDevice{{id: "d"}}, 2)
	// Deliberately undersized test channel proves explicit observation failure.
	c.slots["d"].events = make(chan collector.AcquisitionCompletionObservation, 1)
	for _, id := range []string{"a", "b"} {
		c.accept(collector.AcquisitionCompletionObservation{DeviceID: "d", Metrics: collector.AcquisitionActualMetrics{AcquisitionID: id}})
	}
	_, ok, err := c.next("d", "a")
	require.False(t, ok)
	require.ErrorContains(t, err, "observer lost")
	loss, surplus := c.counts()
	require.EqualValues(t, 1, loss)
	require.Zero(t, surplus)
}
