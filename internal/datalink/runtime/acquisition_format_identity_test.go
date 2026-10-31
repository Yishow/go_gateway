package runtime

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"testing"
)

func TestModernProjectionObservationUsesInstalledFormatIdentity(t *testing.T) {
	f := newAcquisitionInstallFixture(t, true)
	for i, p := range f.projection.Points {
		p.DataType = schema.DataTypeUint32
		p.FormatMode = schema.PointFormatInherit
		p.FormatVersion = 1
		p.RevisionID = "same-point-revision"
		f.projection.RuleLinks[i].Address = p.Address
		f.projection.Rules[i].DataFormat = []string{"ABCD", "DCBA"}[i]
	}
	require.NoError(t, f.service.ApplyWorkspaceProjection(t.Context(), f.projection))
	first := f.poll(t, "P1")
	second := f.poll(t, "P2")
	require.NotEmpty(t, first.Binding.ConfigFingerprint)
	require.NotEqual(t, first.ConfigFingerprint, second.ConfigFingerprint)
	profile, err := f.service.ResolvePointReadProfile(t.Context(), f.projection.Points[0])
	require.NoError(t, err)
	expected, err := measurement.AcquisitionConfigFingerprint(*f.projection.Devices[0], *f.projection.Points[0], profile)
	require.NoError(t, err)
	require.Equal(t, expected, first.ConfigFingerprint)
	require.Equal(t, expected, first.Binding.ConfigFingerprint)
	require.Equal(t, "R1", profile.IdentityView().RuleID)
	// A late result cannot be rebound to the next rule/default generation even
	// when the point revision and old projected DataFormat are unchanged.
	oldLease := f.service.scheduler.CapturePointLease("P1", "same-point-revision")
	f.projection.Rules[0].DataFormat = "CDAB"
	f.projection.Rules[0].RevisionID = "r-next"
	f.projection.Version = "next"
	require.NoError(t, f.service.ApplyWorkspaceProjection(t.Context(), f.projection))
	next := f.poll(t, "P1")
	require.NotEqual(t, first.ConfigFingerprint, next.ConfigFingerprint)
	require.False(t, f.service.scheduler.AcceptsInstallationSample("P1", oldLease, first.Binding.InstallationScope))
	require.False(t, f.service.scheduler.BindPointObservation("P1", oldLease, first.Binding))
	_, err = measurement.AcquisitionConfigFingerprint(*f.projection.Devices[0], *f.projection.Points[0], connector.ReadProfile{})
	require.Error(t, err)
}

func TestLegacyRawCannotOverwriteNewModernPointState(t *testing.T) {
	f := newAcquisitionInstallFixture(t, false)
	p := f.projection.Points[0]
	repo := point.NewMemoryRepository()
	require.NoError(t, repo.Create(t.Context(), p))
	svc := point.NewService(repo, nil)
	f.service.pointSvc = svc
	f.service.config.UpdatePointState = true
	f.service.writer = &mockWriter{}
	require.NoError(t, f.service.ApplyWorkspaceProjection(t.Context(), f.projection))
	old := f.service.scheduler.PollNowContext(t.Context(), []string{p.ID})
	require.Len(t, old, 1)
	require.Equal(t, schema.QualityGood, old[0].Quality)
	p.DataType = schema.DataTypeUint32
	p.FormatMode = schema.PointFormatInherit
	p.FormatVersion = 1
	p.RevisionID = "modern-next"
	f.projection.RuleLinks[0].Address = p.Address
	require.NoError(t, svc.Restore(t.Context(), p))
	require.NoError(t, f.service.ApplyWorkspaceProjection(t.Context(), f.projection))
	f.service.handleCollectedValue(t.Context(), old[0])
	old[0].InstallationLease = 0
	old[0].InstallationScope = ""
	f.service.handleCollectedValue(t.Context(), old[0])
	current, err := svc.GetByID(t.Context(), p.ID)
	require.NoError(t, err)
	require.Nil(t, current.LastValue, "old legacy raw result must not overwrite the newly installed modern point state")
}

func TestLegacyInstalledRawStillAcceptsCurrentLegacyBinding(t *testing.T) {
	f := newAcquisitionInstallFixture(t, false)
	p := f.projection.Points[0]
	repo := point.NewMemoryRepository()
	require.NoError(t, repo.Create(t.Context(), p))
	f.service.pointSvc = point.NewService(repo, nil)
	f.service.config.UpdatePointState = true
	f.service.writer = &mockWriter{}
	require.NoError(t, f.service.ApplyWorkspaceProjection(t.Context(), f.projection))
	values := f.service.scheduler.PollNowContext(t.Context(), []string{p.ID})
	require.Len(t, values, 1)
	require.NotZero(t, values[0].InstallationLease)
	f.service.handleCollectedValue(t.Context(), values[0])
	current, err := f.service.pointSvc.GetByID(t.Context(), p.ID)
	require.NoError(t, err)
	require.NotNil(t, current.LastValue)
	values[0].InstallationLease = 0
	values[0].InstallationScope = ""
	f.service.handleCollectedValue(t.Context(), values[0])
	current, err = f.service.pointSvc.GetByID(t.Context(), p.ID)
	require.NoError(t, err)
	require.NotNil(t, current.LastValue)
}
