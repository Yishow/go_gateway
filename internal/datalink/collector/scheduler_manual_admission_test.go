package collector

import (
	"context"
	"testing"
	"testing/synctest"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/wirebudget"

	"github.com/stretchr/testify/require"
)

// Both entrypoints use the real wirebudget boundary so baseline scalar behavior
// can be compared without pretending that a missing fake method is a RED.
type manualAdmissionProtocol struct{ scheduledExecutionProtocol }

func (p *manualAdmissionProtocol) Read(ctx context.Context, _ connector.ReadRequest) (connector.ReadResult, error) {
	if err := wirebudget.Take(ctx); err != nil {
		return connector.ReadResult{}, err
	}
	p.sends.Add(1)
	return connector.ReadResult{Quality: schema.QualityGood, Value: int16(10)}, nil
}

func manualAdmissionFixture(t *testing.T) (*Scheduler, *manualAdmissionProtocol) {
	t.Helper()
	p := &manualAdmissionProtocol{scheduledExecutionProtocol: scheduledExecutionProtocol{timeout: time.Second}}
	previous, exists := connector.GetFactory(schema.ProtocolMC3E)
	connector.Register(schema.ProtocolMC3E, func() connector.Protocol { return p })
	t.Cleanup(func() {
		if exists {
			connector.Register(schema.ProtocolMC3E, previous)
		} else {
			connector.Unregister(schema.ProtocolMC3E)
		}
	})
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	t.Cleanup(func() { require.NoError(t, manager.CloseAll()) })
	config := DefaultSchedulerConfig()
	config.DefaultRetryCount = 0
	s := NewScheduler(config, manager)
	s.AddDevice(&schema.Device{ID: "device", Protocol: schema.ProtocolMC3E, ConnectionConfig: "{}"})
	s.AddPoint(&schema.Point{ID: "point", DeviceID: "device", Address: "D0", DataType: schema.DataTypeInt16, RevisionID: "old"})
	return s, p
}

func TestManualAdmissionRejectsCapturedInstallationAfterGateWait(t *testing.T) {
	for _, mutation := range []string{"point", "device"} {
		t.Run(mutation, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				s, p := manualAdmissionFixture(t)
				gate := s.deviceLocks["device"]
				gate.Lock()
				done := make(chan []CollectedValue, 1)
				go func() { done <- s.PollNowContext(t.Context(), []string{"point"}) }()
				synctest.Wait()
				if mutation == "point" {
					s.AddPoint(&schema.Point{ID: "point", DeviceID: "device", Address: "D1", DataType: schema.DataTypeInt16, RevisionID: "new"})
				} else {
					s.AddDevice(&schema.Device{ID: "device", Protocol: schema.ProtocolMC3E, ConnectionConfig: `{"timeout":2}`})
				}
				gate.Unlock()
				values := <-done
				require.Zero(t, p.sends.Load(), "a stale captured intent cannot initiate its first send after replacement")
				require.Len(t, values, 1)
				require.Equal(t, "old", values[0].PointRevision)
				require.Equal(t, schema.QualityBad, values[0].Quality)
				require.Empty(t, collectedExecutionValues(s))
			})
		})
	}
}

func TestManualInvalidCapacityFailsClosedBeforeStart(t *testing.T) {
	typ := schema.ProtocolType("manual-invalid-capacity")
	p := &deviceGateProtocol{}
	connector.Register(typ, func() connector.Protocol { return p })
	defer connector.Unregister(typ)
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	defer manager.CloseAll()
	for _, active := range []bool{true, false} {
		config := DefaultSchedulerConfig()
		zero := 0
		if active {
			config.MaxActiveScans = &zero
		} else {
			config.PendingScanCapacity = &zero
		}
		s := NewScheduler(config, manager)
		s.AddDevice(&schema.Device{ID: "device", Protocol: typ, ConnectionConfig: "{}"})
		s.AddPoint(&schema.Point{ID: "point", DeviceID: "device", Address: "D0", DataType: schema.DataTypeInt16})
		ctx, cancel := context.WithTimeout(t.Context(), time.Second)
		values := s.PollNowContext(ctx, []string{"point"})
		cancel()
		require.Len(t, values, 1)
		require.Equal(t, "capacity-invalid", values[0].QualityReason)
		require.Equal(t, schema.QualityBad, values[0].Quality)
	}
}

func TestManualRequestedPointCannotOverrideForbiddenCoverage(t *testing.T) {
	s, p := manualAdmissionFixture(t)
	s.AddDevice(&schema.Device{ID: "device", Protocol: schema.ProtocolMC3E, ConnectionConfig: `{"acquisition":{"version":1,"profile":"safe-v1","coverage":{"revision":"deny","forbidden":[{"area":"D","start":0,"end":1}]}}}`})
	s.AddPoint(&schema.Point{ID: "point", DeviceID: "device", Address: "D0", DataType: schema.DataTypeInt16, RevisionID: "denied"})
	values := s.PollNowContext(t.Context(), []string{"point", "point"})
	require.Zero(t, p.sends.Load(), "explicit intent does not grant a forbidden cell")
	require.Len(t, values, 1)
	require.Equal(t, schema.QualityBad, values[0].Quality)
	require.Contains(t, values[0].Error, "forbidden-address")
}

func TestManualPlannedCallerCancelReturnsEachKnownPointOnce(t *testing.T) {
	for _, preCanceled := range []bool{true, false} {
		s, p := manualAdmissionFixture(t)
		s.AddPoint(&schema.Point{ID: "second", DeviceID: "device", Address: "D10", DataType: schema.DataTypeInt16, RevisionID: "old-second"})
		ctx, cancel := context.WithCancel(t.Context())
		p.hook = func(ctx context.Context, req connector.PlannedReadRequest, _ int) ([]connector.PlannedPointResult, error) {
			cancel()
			return executionResults(req, time.Now(), ctx.Err()), ctx.Err()
		}
		if preCanceled {
			cancel()
		}
		values := s.PollNowContext(ctx, []string{"point", "second", "point", "unknown"})
		cancel()
		require.Len(t, values, 2)
		if preCanceled {
			require.Zero(t, p.sends.Load())
		} else {
			require.EqualValues(t, 1, p.sends.Load())
		}
		seen := map[string]bool{}
		for _, value := range values {
			require.False(t, seen[value.PointID])
			seen[value.PointID] = true
			require.Equal(t, schema.QualityBad, value.Quality)
			require.Contains(t, value.Error, context.Canceled.Error())
		}
		require.Empty(t, collectedExecutionValues(s))
	}
}
