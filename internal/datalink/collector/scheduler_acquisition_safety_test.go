package collector

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestScheduledSafetyInvalidationDuringBackoffDeniesActualSend(t *testing.T) {
	for _, change := range []string{"disable", "same-revision-reinstall", "configuration", "recovery-quarantine", "stop"} {
		t.Run(change, func(t *testing.T) {
			s, scan, p, conn := scheduledExecutionFixture(t, `{}`, "D0")
			s.config.DefaultRetryCount, s.config.DefaultRetryDelay = 3, 40*time.Millisecond
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			s.pollCtx, s.pollCancel, s.running, s.stopCh = ctx, cancel, true, make(chan struct{})
			scan.snapshot.ctx, scan.snapshot.running = ctx, true
			first := make(chan struct{}, 1)
			p.hook = func(_ context.Context, req connector.PlannedReadRequest, send int) ([]connector.PlannedPointResult, error) {
				if send == 1 {
					first <- struct{}{}
				}
				err := errors.New("retryable wire failure")
				return executionResults(req, time.Now(), err), err
			}
			done := make(chan struct{})
			go func() { defer close(done); s.executeScheduledPlan(ctx, conn, scan, "scan", s.deviceBreakers["device"]) }()
			select {
			case <-first:
			case <-time.After(time.Second):
				t.Fatal("first admitted send missing")
			}
			captured := scan.snapshot.points[0]
			point := captured.identityPoint()
			switch change {
			case "disable":
				proof, ok := s.CaptureOrdinaryPoint(point.ID)
				require.True(t, ok)
				point.Enabled = false
				require.True(t, s.ApplyOrdinaryPoint(&point, proof, captured.readProfile))
			case "same-revision-reinstall":
				s.RemovePoint(point.ID)
				require.NotZero(t, s.AddPointWithObservationLease(&point, captured.readProfile))
			case "configuration":
				s.AddDevice(&schema.Device{ID: "device", Protocol: schema.ProtocolMC3E, ConnectionConfig: `{"timeout":2}`})
			case "recovery-quarantine":
				proof, ok := s.CaptureOrdinaryPoint(point.ID)
				require.True(t, ok)
				point.Enabled = true
				_, ok = s.StagePointFormatRecovery(&point, proof, captured.readProfile)
				require.True(t, ok)
			case "stop":
				require.NoError(t, s.Stop())
			}
			select {
			case <-done:
			case <-time.After(time.Second):
				t.Fatal("invalidated retry did not stop")
			}
			require.EqualValues(t, 1, p.sends.Load(), "retry invocation is not a new admitted wire send")
			require.Equal(t, 1, p.budgets[0].Used(), "invalidation does not consume another attempt")
			values := collectedExecutionValues(s)
			if change == "stop" {
				require.Empty(t, values)
				return
			}
			require.Len(t, values, 1)
			require.Equal(t, schema.QualityBad, values[0].Quality)
			require.Equal(t, captured.observationLease, values[0].InstallationLease)
			require.Equal(t, captured.configRevision, values[0].PointRevision)
			require.Equal(t, scan.bindings[point.ID].ConfigFingerprint, values[0].ConfigFingerprint)
			require.Zero(t, s.deviceBreakers["device"].GetStats().TrackerStats.TotalRequests, "stale rejection is not the replacement installation's health")
		})
	}
}

func TestScheduledSafetyLateGoodRetainsOldIdentityWithoutNewHealth(t *testing.T) {
	s, scan, p, conn := scheduledExecutionFixture(t, `{}`, "D0")
	captured := scan.snapshot.points[0]
	p.hook = func(_ context.Context, req connector.PlannedReadRequest, _ int) ([]connector.PlannedPointResult, error) {
		point := captured.identityPoint()
		s.RemovePoint(point.ID)
		s.AddPointWithObservationLease(&point, captured.readProfile)
		return executionResults(req, time.Unix(100, 0), nil), nil
	}
	s.executeScheduledPlan(t.Context(), conn, scan, "old-scan", s.deviceBreakers["device"])
	values := collectedExecutionValues(s)
	require.Len(t, values, 1)
	require.Equal(t, captured.observationLease, values[0].InstallationLease)
	require.Equal(t, scan.bindings[captured.ID].ConfigFingerprint, values[0].ConfigFingerprint)
	require.False(t, s.AcceptsInstallationSample(captured.ID, values[0].InstallationLease, values[0].InstallationScope))
	require.Zero(t, s.deviceBreakers["device"].GetStats().TrackerStats.TotalRequests)
}

func TestScheduledSafetyScanDeadlineCannotBeExtendedByRetry(t *testing.T) {
	s, scan, p, conn := scheduledExecutionFixture(t, `{}`, "D0")
	s.config.DefaultRetryCount = 3
	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Millisecond)
	defer cancel()
	deadline, ok := ctx.Deadline()
	require.True(t, ok)
	p.hook = func(ctx context.Context, req connector.PlannedReadRequest, _ int) ([]connector.PlannedPointResult, error) {
		<-ctx.Done()
		return executionResults(req, time.Now(), ctx.Err()), ctx.Err()
	}
	s.executeScheduledPlan(ctx, conn, scan, "deadline", nil)
	require.EqualValues(t, 1, p.sends.Load())
	require.Equal(t, deadline, p.deadlines[0])
	values := collectedExecutionValues(s)
	require.Len(t, values, 1)
	require.Equal(t, schema.QualityBad, values[0].Quality)
	require.Contains(t, values[0].Error, context.DeadlineExceeded.Error())
}

func TestScheduledSafetyGateDeadlineEmitsOnlyCurrentSnapshot(t *testing.T) {
	for _, replace := range []bool{false, true} {
		t.Run(map[bool]string{false: "current", true: "replaced"}[replace], func(t *testing.T) {
			s, scan, _, _ := scheduledExecutionFixture(t, `{"timeout":1}`, "D0")
			scan.snapshot.capturedAt = time.Now().Add(-980 * time.Millisecond)
			gate := s.deviceLocks["device"]
			gate.Lock()
			defer gate.Unlock()
			if replace {
				point := scan.snapshot.points[0].identityPoint()
				s.RemovePoint(point.ID)
				s.AddPointWithObservationLease(&point, scan.snapshot.points[0].readProfile)
			}
			s.pollScheduledSnapshot(scan.snapshot)
			values := collectedExecutionValues(s)
			if replace {
				require.Empty(t, values)
				return
			}
			require.Len(t, values, 1)
			require.Equal(t, schema.QualityBad, values[0].Quality)
			require.Equal(t, "scan-deadline", values[0].QualityReason)
		})
	}
}

func TestScheduledSafetyOldTickerCannotCaptureReplacementGroup(t *testing.T) {
	s, scan, _, _ := scheduledExecutionFixture(t, `{}`, "D0")
	old := &groupTicker{group: &schema.PollingGroup{ID: scan.snapshot.groupID, IntervalMs: 100}}
	current := &groupTicker{group: &schema.PollingGroup{ID: scan.snapshot.groupID, IntervalMs: 1000}}
	s.groupTickers[scan.snapshot.groupID] = current
	require.Empty(t, s.captureScheduledGroup(scan.snapshot.groupID, old))
	require.Len(t, s.captureScheduledGroup(scan.snapshot.groupID, current), 1)
}
