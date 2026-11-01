package collector

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func scheduledPlanFixture(t *testing.T, config string, addresses ...string) (*Scheduler, scheduledSnapshot) {
	t.Helper()
	s := NewScheduler(DefaultSchedulerConfig(), nil)
	s.AddDevice(&schema.Device{ID: "device", Protocol: schema.ProtocolMC3E, ConnectionConfig: config})
	groupID := "source-cadence"
	for i, address := range addresses {
		p := &schema.Point{ID: fmt.Sprintf("point-%d", i), DeviceID: "device", Address: address, DataType: schema.DataTypeInt16, RevisionID: "point-v1", PollingGroupID: &groupID, FormatMode: schema.PointFormatInherit, FormatVersion: 1}
		profile, err := connector.ResolveReadProfile(*p, decoding.FormatContext{Protocol: schema.ProtocolMC3E, RuleID: "source-rule", RuleOwned: true, RuleRevision: "rule-v1", ConnectionRevision: connector.ReadConnectionRevision(config)})
		require.NoError(t, err)
		require.NotZero(t, s.AddPointWithObservationLease(p, profile))
	}
	snapshots := s.captureScheduledGroup(groupID)
	require.Len(t, snapshots, 1)
	return s, snapshots[0]
}

func TestScheduledPlanUsesInstalledSourceSnapshotNotObservationOrLeaseIdentity(t *testing.T) {
	s, snapshot := scheduledPlanFixture(t, `{}`, "D0", "D1", "D2")
	first, err := buildScheduledPlan(snapshot)
	require.NoError(t, err)
	require.Empty(t, first.failures)
	require.Len(t, first.plan.Requests(), 1)
	require.Equal(t, 3, first.plan.Requests()[0].Length)
	for _, binding := range first.bindings {
		require.Equal(t, "source-rule", binding.Point.SourceIdentity)
		require.Equal(t, "rule-v1", binding.Point.SourceRevision)
		require.Equal(t, "point-v1", binding.PointRevision)
		require.NotZero(t, binding.InstallationLease)
		require.Equal(t, s.InstallationScope(), binding.InstallationScope)
	}
	snapshot.capturedAt = snapshot.capturedAt.Add(time.Hour)
	for i := range snapshot.points {
		snapshot.points[i].observationLease += 100
		snapshot.points[i].installationScope = "another-installation"
		snapshot.points[i].observationBinding = PointObservationBinding{WorkspaceID: "diagnostic-only"}
	}
	second, err := buildScheduledPlan(snapshot)
	require.NoError(t, err)
	require.Equal(t, first.plan.Fingerprint(), second.plan.Fingerprint(), "ephemeral execution/observer metadata cannot become configuration authority")
	require.Equal(t, first.deadline.Add(time.Hour), second.deadline)
	for _, binding := range second.bindings {
		require.NotEqual(t, first.bindings[binding.Point.ID].InstallationLease, binding.InstallationLease)
	}
}

func TestScheduledPlanLegacyProfileIsRealAndFingerprintUnchanged(t *testing.T) {
	s := NewScheduler(DefaultSchedulerConfig(), nil)
	d := schema.Device{ID: "device", Protocol: schema.ProtocolMC3E, ConnectionConfig: `{"data_format":"CDAB"}`}
	p := schema.Point{ID: "point", DeviceID: d.ID, Address: "D0", DataType: schema.DataTypeInt32, RevisionID: "legacy-v1"}
	s.AddDevice(&d)
	s.AddPoint(&p)
	plan, err := buildScheduledPlan(s.captureScheduledGroup("")[0])
	require.NoError(t, err)
	require.Empty(t, plan.failures)
	binding := plan.bindings[p.ID]
	require.True(t, binding.Profile.Present())
	require.Equal(t, schema.PointFormatLegacy, binding.Profile.Resolution().Mode)
	require.Equal(t, "CDAB", binding.Profile.Resolution().EffectiveFormat)
	want, err := measurement.AcquisitionConfigFingerprint(d, p)
	require.NoError(t, err)
	require.Equal(t, want, binding.ConfigFingerprint)
	_, _, err = (connector.PlannedReadRequest{Plan: plan.plan, AcquisitionID: "scan", RequestID: "scan/0", Bindings: []connector.PlannedPointBinding{binding}}).Snapshot()
	require.NoError(t, err, "legacy binding must satisfy ordinary profile checks, not bypass Present")
}

func TestScheduledPlanDeadlineReservesOneFallbackWithoutGrantingAttempts(t *testing.T) {
	_, snapshot := scheduledPlanFixture(t, `{"timeout":2,"acquisition":{"version":1,"profile":"mitsubishi-fx5-cpu-3e-binary-v1"}}`, "D0", "D10", "D20")
	plan, err := buildScheduledPlan(snapshot)
	require.NoError(t, err)
	require.Len(t, plan.plan.Requests(), 1)
	require.Equal(t, "random", plan.plan.Requests()[0].Kind)
	require.Equal(t, snapshot.capturedAt.Add(8*time.Second), plan.deadline, "one initial plus three precomputed safe requests, no retry multiplier")
	safe := safeAcquisitionInput(plan.input, map[string]bool{"point-0": true})
	for _, source := range safe.Sources {
		require.Equal(t, "unsupported", source.Capability.RandomSupport)
		for _, point := range source.Points {
			require.NotEqual(t, "point-0", point.ID)
		}
	}
	_, snapshot = scheduledPlanFixture(t, `{"timeout":9223372036}`, "D0", "D10")
	_, err = buildScheduledPlan(snapshot)
	require.ErrorContains(t, err, "duration")
}

func TestScheduledPlanUnverifiedSourceGeometryCannotAuthorizeGap(t *testing.T) {
	_, snapshot := scheduledPlanFixture(t, `{"acquisition":{"version":1,"profile":"safe-v1","coverage":{"revision":"declared-v1","readable":[{"area":"D","start":0,"end":21}]}}}`, "D0", "D10", "D20")
	plan, err := buildScheduledPlan(snapshot)
	require.NoError(t, err)
	require.Empty(t, plan.failures)
	require.Len(t, plan.plan.Requests(), 3, "installed active points alone cannot prove there are no disabled or skipped source cells in a declared readable gap")
	require.Equal(t, planReasonSourceCoverageUnverified, plan.coverageReason)
}

func TestScheduledPlanUnverifiedGeometryPreservesRestrictionsAndExactReads(t *testing.T) {
	declaration := `{"acquisition":{"version":1,"profile":"safe-v1","coverage":{"revision":"declared-v1","readable":[{"area":"D","start":0,"end":21}],"forbidden":[{"area":"D","start":10,"end":11}],"excluded":[{"area":"D","start":1,"end":10}]}}}`
	_, snapshot := scheduledPlanFixture(t, declaration, "D0", "D10", "D20")
	plan, err := buildScheduledPlan(snapshot)
	require.NoError(t, err)
	require.Contains(t, plan.failures, "point-1", "forbidden remains absolute even for an active point")
	for _, source := range plan.input.Sources {
		require.Empty(t, source.Coverage.Readable)
		require.Len(t, source.Coverage.Forbidden, 1)
		require.Len(t, source.Coverage.Excluded, 1)
	}
	for _, tc := range []struct {
		profile, kind string
		addresses     []string
	}{
		{"safe-v1", "contiguous", []string{"D0", "D1", "D2"}},
		{"mitsubishi-fx5-cpu-3e-binary-v1", "random", []string{"D0", "D10", "D20"}},
	} {
		config := fmt.Sprintf(`{"acquisition":{"version":1,"profile":%q,"coverage":{"revision":"declared-v1","readable":[{"area":"D","start":0,"end":21}]}}}`, tc.profile)
		_, snapshot = scheduledPlanFixture(t, config, tc.addresses...)
		plan, err = buildScheduledPlan(snapshot)
		require.NoError(t, err)
		require.Empty(t, plan.failures)
		require.Len(t, plan.plan.Requests(), 1)
		require.Equal(t, tc.kind, plan.plan.Requests()[0].Kind)
	}
	_, snapshot = scheduledPlanFixture(t, `{}`, "D0")
	plan, err = buildScheduledPlan(snapshot)
	require.NoError(t, err)
	require.Empty(t, plan.coverageReason, "ordinary config without declared gaps is not a coverage error")
}

func TestScheduledPlanReasonIsClosedAndIndependentOfGoodQuality(t *testing.T) {
	s, snapshot := scheduledPlanFixture(t, `{}`, "D0")
	var observations []AcquisitionObservation
	s.config.AcquisitionObserver = func(o AcquisitionObservation) { observations = append(observations, o) }
	for _, reason := range []string{planReasonSourceCoverageUnverified, "private unsafe text", ""} {
		value := s.collectedValueWithReason(snapshot.device, snapshot.points[0], connector.ReadResult{Quality: schema.QualityGood, Value: int16(7)}, nil, "scan", time.Now(), "", acquisitionTrace{planReason: reason})
		require.Equal(t, schema.QualityGood, value.Quality)
		require.Empty(t, value.QualityReason)
	}
	require.Len(t, observations, 3)
	require.Equal(t, planReasonSourceCoverageUnverified, observations[0].PlanReason)
	require.Empty(t, observations[0].QualityReason)
	require.Empty(t, observations[1].PlanReason)
	require.Empty(t, observations[2].PlanReason)
}

func TestScheduledAdmissionRejectsReinstallConfigAndStoppedTicker(t *testing.T) {
	s, snapshot := scheduledPlanFixture(t, `{}`, "D0")
	require.NoError(t, s.admitScheduled(snapshot, snapshot.points))
	p := snapshot.points[0].identityPoint()
	s.RemovePoint(p.ID)
	require.NotZero(t, s.AddPointWithObservationLease(&p, snapshot.points[0].readProfile))
	require.ErrorIs(t, s.admitScheduled(snapshot, snapshot.points), errAcquisitionSuperseded)
	s, snapshot = scheduledPlanFixture(t, `{}`, "D0")
	s.AddDevice(&schema.Device{ID: "device", Protocol: schema.ProtocolMC3E, ConnectionConfig: `{"timeout":2}`})
	require.ErrorIs(t, s.admitScheduled(snapshot, snapshot.points), errAcquisitionSuperseded)
	s, snapshot = scheduledPlanFixture(t, `{}`, "D0")
	snapshot.running, snapshot.ctx = true, context.Background()
	require.ErrorIs(t, s.admitScheduled(snapshot, snapshot.points), errAcquisitionSuperseded)
}

func TestScheduledScatterRejectsForeignTimeAndDuplicatePoint(t *testing.T) {
	_, snapshot := scheduledPlanFixture(t, `{}`, "D0")
	plan, err := buildScheduledPlan(snapshot)
	require.NoError(t, err)
	read := connector.PlannedReadRequest{Plan: plan.plan, AcquisitionID: "scan", RequestID: "request", Bindings: []connector.PlannedPointBinding{plan.bindings["point-0"]}}
	request, bindings, err := read.Snapshot()
	require.NoError(t, err)
	foreignTime, completed := time.Unix(1, 0), time.Unix(20, 0)
	foreign := connector.PlannedPointResult{Binding: bindings[0], Group: request.Group, PlanFingerprint: plan.plan.Fingerprint(), AcquisitionID: "other-scan", RequestID: "request", ReceivedAt: foreignTime, ReadResult: connector.ReadResult{Quality: schema.QualityGood, Value: int16(99)}}
	for _, malformed := range [][]connector.PlannedPointResult{{foreign}, {foreign, foreign}} {
		results, err := validateScheduledResults(read, malformed, nil, completed)
		require.ErrorIs(t, err, connector.ErrPlannedRead)
		require.Len(t, results, 1)
		require.Equal(t, completed, results[0].ReceivedAt)
		require.Equal(t, schema.QualityBad, results[0].Quality)
		require.Nil(t, results[0].Value)
		require.Equal(t, "scan", results[0].AcquisitionID)
	}
}
