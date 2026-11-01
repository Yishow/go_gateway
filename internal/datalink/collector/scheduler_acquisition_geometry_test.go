package collector

import (
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func seedAcquisitionGeometry(t *testing.T) *Scheduler {
	t.Helper()
	s := NewScheduler(DefaultSchedulerConfig(), nil)
	s.AddDevice(&schema.Device{ID: "d", Protocol: schema.ProtocolMC3E, ConnectionConfig: `{}`})
	s.AddPoint(&schema.Point{ID: "p", DeviceID: "d", Address: "D0", DataType: schema.DataTypeUint16})
	receipt := s.CaptureAcquisitionGeometry("d")
	s.acquisitionGeometry = map[string]*acquisitionGeometryState{"d": {publication: &acquisitionGeometry{revision: "source-semantic", device: receipt.device, slots: receipt.slots}, pending: map[*acquisitionGeometryMutation]struct{}{}}}
	require.NotNil(t, s.currentAcquisitionGeometryLocked("d"))
	return s
}

func TestAcquisitionGeometryMutationQuarantineDoesNotResumeOldProof(t *testing.T) {
	s := seedAcquisitionGeometry(t)
	before := s.CaptureAcquisitionGeometry("d")
	again := s.CaptureAcquisitionGeometry("d")
	require.Same(t, before.publication, again.publication, "pure capture cannot quarantine or publish")
	a := s.BeginAcquisitionGeometryMutation("d")
	require.Nil(t, s.currentAcquisitionGeometryLocked("d"))
	require.Equal(t, before.slots, s.CaptureAcquisitionGeometry("d").slots, "quarantine must not reinstall source points")
	b := s.BeginAcquisitionGeometryMutation("d")
	s.FinishAcquisitionGeometryMutation(a)
	require.Len(t, s.acquisitionGeometry["d"].pending, 1)
	s.FinishAcquisitionGeometryMutation(a)
	require.Len(t, s.acquisitionGeometry["d"].pending, 1, "duplicate finalization cannot release B")
	s.FinishAcquisitionGeometryMutation(AcquisitionGeometryMutation{})
	NewScheduler(DefaultSchedulerConfig(), nil).FinishAcquisitionGeometryMutation(b)
	require.Len(t, s.acquisitionGeometry["d"].pending, 1, "foreign scheduler cannot release B")
	s.FinishAcquisitionGeometryMutation(b)
	require.Empty(t, s.acquisitionGeometry["d"].pending)
	require.Nil(t, s.currentAcquisitionGeometryLocked("d"), "success, cancellation or rollback without accepted SQL proof must retain unknown")
	require.Equal(t, before.slots, s.CaptureAcquisitionGeometry("d").slots)
}

func TestAcquisitionGeometryFullSlotProofRejectsNonEpochRemovalAndForeignAdds(t *testing.T) {
	for _, change := range []string{"remove", "add", "same-revision-reinstall", "config"} {
		t.Run(change, func(t *testing.T) {
			s := seedAcquisitionGeometry(t)
			epoch := s.pointObservationLease
			switch change {
			case "remove":
				s.RemovePoint("p")
				require.Equal(t, epoch, s.pointObservationLease)
			case "add":
				s.AddPoint(&schema.Point{ID: "foreign", DeviceID: "d", Address: "D1", DataType: schema.DataTypeUint16})
			case "same-revision-reinstall":
				s.AddPoint(&schema.Point{ID: "p", DeviceID: "d", Address: "D0", DataType: schema.DataTypeUint16})
			case "config":
				s.AddDevice(&schema.Device{ID: "d", Protocol: schema.ProtocolMC3E, ConnectionConfig: `{"timeout":2}`})
			}
			require.Nil(t, s.currentAcquisitionGeometryLocked("d"))
		})
	}
}

func TestAcquisitionGeometryInstalledProofEnablesOnlyDeclaredReadableGaps(t *testing.T) {
	config := `{"acquisition":{"version":1,"profile":"safe-v1","coverage":{"revision":"operator-1","readable":[{"area":"D","start":0,"end":21}]}}}`
	s, _ := scheduledPlanFixture(t, config, "D0", "D10", "D20")
	before := s.CaptureAcquisitionGeometry("device")
	s.acquisitionGeometry = map[string]*acquisitionGeometryState{"device": {publication: &acquisitionGeometry{revision: "complete-source-1", device: before.device, slots: before.slots}, pending: map[*acquisitionGeometryMutation]struct{}{}}}
	captured := s.captureScheduledGroup("source-cadence")[0]
	planned, err := buildScheduledPlan(captured)
	require.NoError(t, err)
	require.Len(t, planned.plan.Requests(), 1, "declared gap coverage needs a captured, full installed source proof")
	require.Equal(t, 21, planned.plan.Requests()[0].Length)
	require.Empty(t, planned.coverageReason)
	token := s.BeginAcquisitionGeometryMutation("device")
	unknown, err := buildScheduledPlan(s.captureScheduledGroup("source-cadence")[0])
	require.NoError(t, err)
	require.Len(t, unknown.plan.Requests(), 3)
	require.Equal(t, planReasonSourceCoverageUnverified, unknown.coverageReason)
	require.NotEqual(t, planned.plan.Fingerprint(), unknown.plan.Fingerprint())
	s.FinishAcquisitionGeometryMutation(token)
	stillUnknown, err := buildScheduledPlan(s.captureScheduledGroup("source-cadence")[0])
	require.NoError(t, err)
	require.Equal(t, unknown.plan.Fingerprint(), stillUnknown.plan.Fingerprint(), "ending quarantine cannot restore stale complete geometry")
}

func TestAcquisitionGeometryUninstalledMutationDoesNotRetainUnknownEntries(t *testing.T) {
	s := NewScheduler(DefaultSchedulerConfig(), nil)
	a := s.BeginAcquisitionGeometryMutation("not-installed")
	b := s.BeginAcquisitionGeometryMutation("not-installed")
	s.FinishAcquisitionGeometryMutation(a)
	require.Contains(t, s.acquisitionGeometry, "not-installed")
	s.FinishAcquisitionGeometryMutation(b)
	require.NotContains(t, s.acquisitionGeometry, "not-installed")
}
