package collector

import (
	"context"
	"errors"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/wirebudget"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAcquisitionPreviewMissingInstallIsUnavailable(t *testing.T) {
	s := NewScheduler(DefaultSchedulerConfig(), nil)
	view, _, err := s.PreviewAcquisition(t.Context(), "w", "device", "source")
	require.NoError(t, err)
	require.Equal(t, "unavailable", view.Status)
	require.Empty(t, view.Groups)
}

func TestAcquisitionPreviewInstalledParityAndHead(t *testing.T) {
	s, snapshot := scheduledPlanFixture(t, `{}`, "D0", "D1", "D20")
	s.connMgr = connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	s.mu.Lock()
	s.groupTickers[snapshot.groupID] = &groupTicker{group: &schema.PollingGroup{ID: snapshot.groupID, IntervalMs: 1000, Enabled: true}}
	s.acquisitionGeometry = map[string]*acquisitionGeometryState{"device": {publication: &acquisitionGeometry{workspaceID: "w", sourceRevision: "source", revision: "source", device: snapshot.device, slots: s.captureAcquisitionGeometryLocked("device").slots}}}
	s.mu.Unlock()
	view, proof, err := s.PreviewAcquisition(t.Context(), "w", "device", "source")
	require.NoError(t, err)
	require.Equal(t, "available", view.Status)
	require.Len(t, view.Groups, 1)
	require.Equal(t, 3, view.Groups[0].AcquisitionPointCount)
	require.Equal(t, 2, view.Groups[0].PlannedRequests)
	require.Equal(t, "declared_no_connection", view.Groups[0].Basis)
	require.Equal(t, "unavailable", view.Groups[0].Actual.Status)
	require.True(t, s.AcquisitionPreviewCurrent(proof))
	second, _, err := s.PreviewAcquisition(t.Context(), "w", "device", "source")
	require.NoError(t, err)
	require.Equal(t, view.PreviewRevision, second.PreviewRevision)
	s.RemovePoint("point-0")
	require.False(t, s.AcquisitionPreviewCurrent(proof))
}

func TestAcquisitionPreviewDifferentSourceAndCadenceUnavailable(t *testing.T) {
	s, _ := scheduledPlanFixture(t, `{}`, "D0")
	view, _, err := s.PreviewAcquisition(t.Context(), "w", "device", "other")
	require.NoError(t, err)
	require.Equal(t, "unavailable", view.Status)
	require.Empty(t, view.Groups)
}

func TestAcquisitionEvidenceWireLogicalUnitsAndCancellation(t *testing.T) {
	e := &acquisitionScanEvidence{metrics: AcquisitionActualMetrics{LogicalRequests: AcquisitionLogicalOutcomes{Status: "not_started"}}}
	ctx := context.WithValue(t.Context(), acquisitionEvidenceKey{}, e)
	first, finish := beginEvidenceRequest(ctx, "initial")
	require.NoError(t, wirebudget.Take(first))
	require.NoError(t, wirebudget.Take(first))
	finish(nil, errors.New("unsupported"))
	fallback, done := beginEvidenceRequest(ctx, "fallback")
	require.NoError(t, wirebudget.Take(fallback))
	require.NoError(t, wirebudget.Take(fallback))
	done(nil, context.Canceled)
	require.Equal(t, uint64(4), e.metrics.WireAttempts)
	require.Equal(t, uint64(2), e.metrics.RetryWireAttempts)
	require.Equal(t, uint64(2), e.metrics.FallbackWireAttempts)
	require.Equal(t, uint64(1), e.metrics.LogicalRequests.Failed)
	require.Equal(t, uint64(1), e.metrics.LogicalRequests.Canceled)
	_, completed := beginEvidenceRequest(ctx, "initial")
	completed([]connector.PlannedPointResult{{ReadResult: connector.ReadResult{Quality: schema.QualityGood}, ReceivedAt: time.Now()}}, nil)
	require.Equal(t, uint64(1), e.metrics.LogicalRequests.Succeeded)
	require.Equal(t, uint64(4), e.metrics.WireAttempts, "logical zero-send result is not an invented wire send")
}

func TestAcquisitionPreviewGroupsRemainSeparateAndRejectedMembershipVisible(t *testing.T) {
	s, _ := previewExecutionFixture(t, `{}`, "D0", "invalid")
	secondGroup := "slow"
	point := &schema.Point{ID: "slow-point", DeviceID: "device", Address: "D0", DataType: schema.DataTypeUint16, PollingGroupID: &secondGroup, RevisionID: "slow-v1"}
	s.AddPoint(point)
	ticker := time.NewTicker(time.Hour)
	t.Cleanup(ticker.Stop)
	s.groupTickers[secondGroup] = &groupTicker{group: &schema.PollingGroup{ID: secondGroup, IntervalMs: 5000, Enabled: true}, ticker: ticker, stopCh: make(chan struct{})}
	receipt := s.CaptureAcquisitionGeometry("device")
	s.acquisitionGeometry["device"].publication = &acquisitionGeometry{workspaceID: "w", sourceRevision: "source", revision: "source", device: receipt.device, slots: receipt.slots}
	view, _, err := s.PreviewAcquisition(t.Context(), "w", "device", "source")
	require.NoError(t, err)
	require.Len(t, view.Groups, 2)
	require.Equal(t, int64(5000), view.Groups[0].PollIntervalMS)
	require.Equal(t, int64(1000), view.Groups[1].PollIntervalMS)
	require.Equal(t, 2, view.Groups[1].AcquisitionPointCount)
	require.Len(t, view.Groups[1].PointErrors, 1)
	require.Equal(t, 1, view.Groups[1].UniqueWireUnits[0].UniqueAddressUnits)
}

func TestAcquisitionPreviewNeverBorrowsIdenticalForeignWorkspaceGeometry(t *testing.T) {
	s, _ := previewExecutionFixture(t, `{}`, "D0")
	view, _, err := s.PreviewAcquisition(t.Context(), "other-workspace", "device", "source")
	require.NoError(t, err)
	require.Equal(t, "unavailable", view.Status)
	require.Empty(t, view.Groups)
}

func TestAcquisitionPreviewActualRejectsForeignWorkspaceWithSameSlots(t *testing.T) {
	s, _ := previewExecutionFixture(t, `{}`, "D0")
	snapshot := s.captureScheduledGroup("source-cadence")[0]
	scan, err := buildScheduledPlan(snapshot)
	require.NoError(t, err)
	_, e := s.beginAcquisitionEvidence(t.Context(), scan, "old-workspace")
	s.finishAcquisitionEvidence(e)
	old := s.acquisitionGeometry["device"].publication
	replacement := *old
	replacement.workspaceID = "other-workspace"
	s.acquisitionGeometry["device"].publication = &replacement
	view, _, err := s.PreviewAcquisition(t.Context(), "other-workspace", "device", "source")
	require.NoError(t, err)
	require.Equal(t, "available", view.Status)
	require.Equal(t, "unavailable", view.Groups[0].Actual.Status, "source semantic and point IDs cannot relabel a prior workspace scan")
}

func TestAcquisitionPreviewCanceledReadDoesNotPublish(t *testing.T) {
	s, p := previewExecutionFixture(t, `{}`, "D0")
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	view, _, err := s.PreviewAcquisition(ctx, "w", "device", "source")
	require.ErrorIs(t, err, context.Canceled)
	require.Equal(t, "unavailable", view.Status)
	require.Empty(t, view.Groups)
	require.Zero(t, p.connects.Load())
	require.Zero(t, p.sends.Load())
}
