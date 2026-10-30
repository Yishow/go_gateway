package runtime

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"testing"
)

func TestReadProfileProjectionMandatoryWithoutObserver(t *testing.T) {
	f := newAcquisitionInstallFixture(t, false)
	p := f.projection.Points[0]
	p.FormatMode = schema.PointFormatInherit
	p.FormatVersion = 1
	p.RevisionID = "point-1"
	p.DataType = schema.DataTypeUint32
	f.projection.Rules[0].DataFormat = "CDAB"
	f.projection.RuleLinks[0].Address = p.Address
	require.NoError(t, f.service.ApplyWorkspaceProjection(t.Context(), f.projection))
	lease := f.service.CapturePointLease(p.ID, p.RevisionID)
	require.NotZero(t, lease)
	profile, ok := f.service.CapturePointReadProfile(p.ID, lease)
	require.True(t, ok)
	require.True(t, profile.Present())
	require.Equal(t, "CDAB", profile.Resolution().EffectiveFormat)
	require.Equal(t, f.projection.Rules[0].RevisionID, profile.Resolution().RuleRevision)
	f.projection.Rules[0].DataFormat = "DCBA"
	old, ok := f.service.CapturePointReadProfile(p.ID, lease)
	require.True(t, ok)
	require.Equal(t, profile, old)
	values := f.service.scheduler.PollNowContext(t.Context(), []string{p.ID})
	require.Len(t, values, 1)
	require.Equal(t, schema.QualityGood, values[0].Quality)
}
func TestReadProfileRejectsAbsentAuthorityAndWrongRevision(t *testing.T) {
	f := newAcquisitionInstallFixture(t, false)
	p := *f.projection.Points[0]
	p.FormatMode = schema.PointFormatOverride
	p.FormatOverride = "CDAB"
	p.FormatVersion = 1
	p.DataType = schema.DataTypeUint32
	p.RevisionID = "not-persisted"
	require.Zero(t, f.service.InstallPointForMutation(&p))
	require.Zero(t, f.service.CapturePointLease(p.ID, p.RevisionID))
	_, err := f.service.pointReadProfile(t.Context(), &p)
	require.ErrorIs(t, err, connector.ErrReadProfile)
}
func TestReadProfileCompensationUsesCapturedPolicy(t *testing.T) {
	f := newAcquisitionInstallFixture(t, false)
	p := f.projection.Points[0]
	p.FormatMode = schema.PointFormatInherit
	p.FormatVersion = 1
	p.RevisionID = "point-1"
	p.DataType = schema.DataTypeUint32
	f.projection.Rules[0].DataFormat = "CDAB"
	f.projection.RuleLinks[0].Address = p.Address
	require.NoError(t, f.service.ApplyWorkspaceProjection(t.Context(), f.projection))
	original := *p
	lease := f.service.CapturePointLease(p.ID, p.RevisionID)
	profile, ok := f.service.CapturePointReadProfile(p.ID, lease)
	require.True(t, ok)
	p.RevisionID = "point-2"
	f.projection.Rules[0].DataFormat = "DCBA"
	f.projection.Rules[0].RevisionID = "rule-2"
	installed, ok := f.service.InstallPointForMutationAtLease(p, lease)
	require.True(t, ok)
	restored, ok := f.service.RestorePointProfileReceipt(&original, installed, profile)
	require.True(t, ok)
	actual, ok := f.service.CapturePointReadProfile(p.ID, restored)
	require.True(t, ok)
	require.Equal(t, profile, actual)
	foreign := f.service.InstallPointForMutation(p)
	require.NotZero(t, foreign)
	_, ok = f.service.RestorePointProfileReceipt(&original, restored, profile)
	require.False(t, ok)
	require.Equal(t, foreign, f.service.CapturePointLease(p.ID, p.RevisionID))
}

func TestReadProfileScopedDeviceIgnoresForeignInvalidAndDisabledPoints(t *testing.T) {
	f := newAcquisitionInstallFixture(t, false)
	p := f.projection.Points[0]
	p.FormatMode = schema.PointFormatOverride
	p.FormatOverride = "CDAB"
	p.FormatVersion = 1
	p.DataType = schema.DataTypeUint32
	p.RevisionID = "p1"
	f.projection.RuleLinks[0].Address = p.Address
	foreign := f.projection.Points[1]
	foreign.FormatMode = schema.PointFormatOverride
	foreign.FormatOverride = "UNKNOWN"
	foreign.FormatVersion = 1
	require.NoError(t, f.service.ApplyWorkspaceProjectionForDevices(t.Context(), f.projection, []string{p.DeviceID}))
	require.NotZero(t, f.service.CapturePointLease(p.ID, p.RevisionID))
	// Full apply also must not parse a device that is being disabled/removed.
	f.projection.Devices[1].Status = schema.DeviceStatusDisabled
	require.NoError(t, f.service.ApplyWorkspaceProjection(t.Context(), f.projection))
}
func TestReadProfileDisableAndMissingAuthorityStopOnlyOwnedSlot(t *testing.T) {
	f := newAcquisitionInstallFixture(t, false)
	p := f.projection.Points[0]
	p.FormatMode = schema.PointFormatOverride
	p.FormatOverride = "CDAB"
	p.FormatVersion = 1
	p.DataType = schema.DataTypeUint32
	p.RevisionID = "p1"
	f.projection.RuleLinks[0].Address = p.Address
	require.NoError(t, f.service.ApplyWorkspaceProjection(t.Context(), f.projection))
	proof, ok := f.service.CaptureOrdinaryPoint(p.ID)
	require.True(t, ok)
	f.service.workspace = nil
	changed := *p
	changed.RevisionID = "p2"
	require.False(t, f.service.ApplyOrdinaryPoint(&changed, proof))
	require.Zero(t, f.service.CapturePointLease(p.ID, p.RevisionID))
	// Disabled mutation never requires a decodable profile/provider.
	p.Enabled = false
	proof, ok = f.service.CaptureOrdinaryPoint(p.ID)
	require.True(t, ok)
	require.True(t, f.service.ApplyOrdinaryPoint(p, proof))
}

func TestReadProfileInvalidProjectionStopsMatchingOldInstall(t *testing.T) {
	f := newAcquisitionInstallFixture(t, false)
	p := f.projection.Points[0]
	p.FormatMode = schema.PointFormatOverride
	p.FormatOverride = "CDAB"
	p.FormatVersion = 1
	p.DataType = schema.DataTypeUint32
	p.RevisionID = "p1"
	f.projection.RuleLinks[0].Address = p.Address
	require.NoError(t, f.service.ApplyWorkspaceProjection(t.Context(), f.projection))
	lease := f.service.CapturePointLease(p.ID, p.RevisionID)
	require.NotZero(t, lease)
	p.FormatOverride = "UNKNOWN"
	p.RevisionID = "p2"
	require.Error(t, f.service.ApplyWorkspaceProjection(t.Context(), f.projection))
	require.Zero(t, f.service.CapturePointLease(p.ID, "p1"))
}

func TestReadProfileFailedOrdinaryResolutionPreservesForeignSameRevisionLease(t *testing.T) {
	f := newAcquisitionInstallFixture(t, false)
	p := f.projection.Points[0]
	p.FormatMode = schema.PointFormatOverride
	p.FormatOverride = "CDAB"
	p.FormatVersion = 1
	p.DataType = schema.DataTypeUint32
	p.RevisionID = "p1"
	f.projection.RuleLinks[0].Address = p.Address
	require.NoError(t, f.service.ApplyWorkspaceProjection(t.Context(), f.projection))
	proof, ok := f.service.CaptureOrdinaryPoint(p.ID)
	require.True(t, ok)
	foreign := f.service.InstallPointForMutation(p)
	require.NotZero(t, foreign)
	require.NotEqual(t, proof.Lease, foreign)
	f.service.workspace = nil
	changed := *p
	changed.RevisionID = "late-revision"
	require.False(t, f.service.ApplyOrdinaryPoint(&changed, proof))
	require.Equal(t, foreign, f.service.CapturePointLease(p.ID, p.RevisionID))
}

func TestReadProfileLookupConsumesOriginalInstallationProof(t *testing.T) {
	for _, conditional := range []bool{false, true} {
		for _, event := range []string{"normal", "foreign", "absent-invalidation"} {
			t.Run(fmt.Sprintf("conditional-%t/%s", conditional, event), func(t *testing.T) {
				f := newAcquisitionInstallFixture(t, false)
				p := f.projection.Points[0]
				p.FormatMode = schema.PointFormatInherit
				p.FormatVersion = 1
				p.RevisionID = "p1"
				p.DataType = schema.DataTypeUint32
				f.projection.RuleLinks[0].Address = p.Address
				require.NoError(t, f.service.ApplyWorkspaceProjection(t.Context(), f.projection))
				old := f.service.CapturePointLease(p.ID, p.RevisionID)
				profile, ok := f.service.CapturePointReadProfile(p.ID, old)
				require.True(t, ok)
				if event == "absent-invalidation" {
					require.True(t, f.service.SuspendPointLease(p.ID, old))
					old = 0
				}
				var foreign uint64
				f.service.workspace = &compensationProjectionReader{projection: f.projection, hook: func() {
					switch event {
					case "foreign":
						foreign = f.service.scheduler.AddPointWithObservationLease(p, profile)
					case "absent-invalidation":
						require.True(t, f.service.scheduler.InvalidatePointFormat(p.ID, p.RevisionID))
					}
				}}
				var got uint64
				if conditional {
					got, ok = f.service.InstallPointForMutationAtLease(p, old)
				} else {
					got = f.service.InstallPointForMutation(p)
					ok = got != 0
				}
				if event == "normal" {
					require.True(t, ok)
					require.NotZero(t, got)
					return
				}
				require.False(t, ok)
				require.Zero(t, got)
				require.Equal(t, foreign, f.service.CapturePointLease(p.ID, p.RevisionID))
			})
		}
	}
}
