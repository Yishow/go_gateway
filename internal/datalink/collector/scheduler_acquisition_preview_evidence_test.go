package collector

import (
	"context"
	"testing"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/mcprotocol"

	"github.com/stretchr/testify/require"
)

func previewExecutionFixture(t *testing.T, config string, addresses ...string) (*Scheduler, *scheduledExecutionProtocol) {
	t.Helper()
	s, snapshot := scheduledPlanFixture(t, config, addresses...)
	s.config.DefaultRetryCount, s.config.DefaultRetryDelay = 0, 0
	p := &scheduledExecutionProtocol{timeout: time.Second}
	old, exists := connector.GetFactory(schema.ProtocolMC3E)
	connector.Register(schema.ProtocolMC3E, func() connector.Protocol { return p })
	t.Cleanup(func() {
		if exists {
			connector.Register(schema.ProtocolMC3E, old)
		} else {
			connector.Unregister(schema.ProtocolMC3E)
		}
	})
	s.connMgr = connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	t.Cleanup(func() { require.NoError(t, s.connMgr.CloseAll()) })
	ticker := time.NewTicker(time.Hour)
	t.Cleanup(ticker.Stop)
	s.groupTickers[snapshot.groupID] = &groupTicker{group: &schema.PollingGroup{ID: snapshot.groupID, IntervalMs: 1000, Enabled: true}, ticker: ticker, stopCh: make(chan struct{})}
	receipt := s.CaptureAcquisitionGeometry("device")
	s.acquisitionGeometry = map[string]*acquisitionGeometryState{"device": {publication: &acquisitionGeometry{workspaceID: "w", sourceRevision: "source", revision: "source", device: receipt.device, slots: receipt.slots}}}
	return s, p
}
func TestAcquisitionPreviewActualOnlyScheduledMatchingLatest(t *testing.T) {
	s, p := previewExecutionFixture(t, `{}`, "D0", "D1")
	before, _, err := s.PreviewAcquisition(t.Context(), "w", "device", "source")
	require.NoError(t, err)
	require.Equal(t, "available", before.Status)
	require.EqualValues(t, 0, p.connects.Load())
	require.EqualValues(t, 0, p.sends.Load())
	snapshot := s.captureScheduledGroup("source-cadence")[0]
	s.pollAcquisitionSnapshot(snapshot, "scheduled-1", false, func(CollectedValue) {})
	view, proof, err := s.PreviewAcquisition(t.Context(), "w", "device", "source")
	require.NoError(t, err)
	actual := view.Groups[0].Actual
	require.Equal(t, "available", actual.Status)
	require.Equal(t, "scheduled-1", actual.AcquisitionID)
	require.EqualValues(t, 1, actual.WireAttempts)
	require.EqualValues(t, 1, actual.LogicalRequests.Succeeded)
	s.pollAcquisitionSnapshot(s.captureScheduledGroup("source-cadence")[0], "manual-2", true, func(CollectedValue) {})
	again, _, err := s.PreviewAcquisition(t.Context(), "w", "device", "source")
	require.NoError(t, err)
	require.Equal(t, "scheduled-1", again.Groups[0].Actual.AcquisitionID)
	require.Equal(t, view.PreviewRevision, again.PreviewRevision)
	require.True(t, s.AcquisitionPreviewCurrent(proof))
	p.generation.Add(1)
	require.False(t, s.AcquisitionPreviewCurrent(proof))
	changed, _, err := s.PreviewAcquisition(t.Context(), "w", "device", "source")
	require.NoError(t, err)
	require.Equal(t, "unavailable", changed.Groups[0].Actual.Status)
}
func TestAcquisitionPreviewWarmSafeInitialIsNotFallback(t *testing.T) {
	s, p := previewExecutionFixture(t, `{"acquisition":{"version":1,"profile":"mitsubishi-fx5-cpu-3e-binary-v1"}}`, "D0", "D10", "D20")
	p.hook = func(_ context.Context, req connector.PlannedReadRequest, _ int) ([]connector.PlannedPointResult, error) {
		request, _ := req.Plan.Request(req.RequestIndex)
		if request.Kind == "random" {
			err := &mcprotocol.MCError{Code: 0xC059}
			return executionResults(req, time.Now(), err), err
		}
		return executionResults(req, time.Now(), nil), nil
	}
	s.pollAcquisitionSnapshot(s.captureScheduledGroup("source-cadence")[0], "fallback-scan", false, func(CollectedValue) {})
	first, _, err := s.PreviewAcquisition(t.Context(), "w", "device", "source")
	require.NoError(t, err)
	require.Equal(t, "safe_current_connection", first.Groups[0].Basis)
	require.Equal(t, "unavailable", first.Groups[0].Actual.Status, "fallback changed context, not a matching warm-safe initial scan")
	s.pollAcquisitionSnapshot(s.captureScheduledGroup("source-cadence")[0], "warm-scan", false, func(CollectedValue) {})
	second, _, err := s.PreviewAcquisition(t.Context(), "w", "device", "source")
	require.NoError(t, err)
	require.Equal(t, first.PreviewRevision, second.PreviewRevision)
	require.Equal(t, "available", second.Groups[0].Actual.Status)
	require.EqualValues(t, 0, second.Groups[0].Actual.FallbackWireAttempts)
	require.EqualValues(t, 3, second.Groups[0].Actual.WireAttempts)
	require.Equal(t, "safe", second.Groups[0].Actual.Phases[0].Phase)
}
func TestAcquisitionPreviewEvidenceLateAndLossIsolation(t *testing.T) {
	s, _ := previewExecutionFixture(t, `{}`, "D0")
	snapshot := s.captureScheduledGroup("source-cadence")[0]
	scan, err := buildScheduledPlan(snapshot)
	require.NoError(t, err)
	_, old := s.beginAcquisitionEvidence(t.Context(), scan, "old")
	_, newer := s.beginAcquisitionEvidence(t.Context(), scan, "new")
	newer.snapshot.capturedAt = old.snapshot.capturedAt.Add(time.Nanosecond)
	s.finishAcquisitionEvidence(newer)
	s.finishAcquisitionEvidence(old)
	view, _, err := s.PreviewAcquisition(t.Context(), "w", "device", "source")
	require.NoError(t, err)
	require.Equal(t, "new", view.Groups[0].Actual.AcquisitionID)
	require.Equal(t, "not_started", view.Groups[0].Actual.LogicalRequests.Status)
	require.Nil(t, view.Groups[0].Actual.LogicalRequests.AcquisitionLogicalCounts)
	require.EqualValues(t, 0, view.Groups[0].Actual.WireAttempts)
	s.outputLoss.record(CollectedValue{PointID: "foreign", PointRevision: "v", InstallationScope: "other"})
	require.Equal(t, "unavailable", s.previewOutputLoss(s.CaptureAcquisitionGeometry("device")).Status)
	pt := snapshot.points[0]
	s.outputLoss.record(CollectedValue{PointID: pt.ID, PointRevision: pt.configRevision, InstallationScope: pt.installationScope})
	loss := s.previewOutputLoss(s.CaptureAcquisitionGeometry("device"))
	require.Equal(t, "partial", loss.Status)
	require.Len(t, loss.RetainedDetail, 1)
	require.EqualValues(t, 1, loss.RetainedDetail[0].Losses)
	s.RemovePollingGroup("source-cadence")
	require.Empty(t, s.acquisitionEvidence)
}

func TestAcquisitionEvidenceCompletedRejectionBeforeCancelStaysFailed(t *testing.T) {
	s, scan, p, conn := scheduledExecutionFixture(t, `{"acquisition":{"version":1,"profile":"mitsubishi-fx5-cpu-3e-binary-v1"}}`, "D0", "D10", "D20")
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	scan.snapshot.ctx = ctx
	p.hook = func(_ context.Context, req connector.PlannedReadRequest, _ int) ([]connector.PlannedPointResult, error) {
		err := &mcprotocol.MCError{Code: 0xC059}
		results := executionResults(req, time.Now(), err)
		cancel()
		return results, err
	}
	ctx, evidence := s.beginAcquisitionEvidence(ctx, scan, "canceled-after-rejection")
	s.executeAcquisitionPlan(ctx, conn, scan, "canceled-after-rejection", nil, func(CollectedValue) {})
	require.EqualValues(t, 1, evidence.metrics.LogicalRequests.Failed)
	require.EqualValues(t, 0, evidence.metrics.LogicalRequests.Canceled)
	require.EqualValues(t, 1, evidence.metrics.WireAttempts)
	require.EqualValues(t, 0, evidence.metrics.FallbackWireAttempts, "canceled scan cannot send fallback but retains finished rejection")
}

func TestAcquisitionPreviewConcurrentCompletionRemovalAndRestart(t *testing.T) {
	s, _ := previewExecutionFixture(t, `{}`, "D0")
	snapshot := s.captureScheduledGroup("source-cadence")[0]
	scan, err := buildScheduledPlan(snapshot)
	require.NoError(t, err)
	_, e := s.beginAcquisitionEvidence(t.Context(), scan, "completed")
	done := make(chan struct{})
	go func() {
		defer close(done)
		for range 30 {
			completed := *e
			s.finishAcquisitionEvidence(&completed)
		}
	}()
	for range 30 {
		_, _, _ = s.PreviewAcquisition(t.Context(), "w", "device", "source")
	}
	<-done
	s.RemovePollingGroup("source-cadence")
	require.Empty(t, s.acquisitionEvidence)
	s.finishAcquisitionEvidence(e)
	require.Empty(t, s.acquisitionEvidence, "late removed tuple cannot republish")
	require.NoError(t, s.Start(nil))
	require.NoError(t, s.Stop())
	s.finishAcquisitionEvidence(e)
	require.Empty(t, s.acquisitionEvidence, "restart cannot borrow old lifetime evidence")
}
