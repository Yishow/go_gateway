package runtime

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"slices"
	"testing"
)

func TestFormatRecoveryRefreshesMandatoryTypedBinding(t *testing.T) {
	f := newAcquisitionInstallFixture(t, false)
	p := f.projection.Points[0]
	p.FormatMode = schema.PointFormatOverride
	p.FormatOverride = "ABCD"
	p.DataFormat = "ABCD"
	p.FormatVersion = 1
	p.RevisionID = "old"
	p.DataType = schema.DataTypeUint32
	f.projection.RuleLinks[0].Address = p.Address
	f.projection.Tags = []*schema.Tag{{ID: "T", Key: "T", DataType: schema.DataTypeUint32}}
	f.projection.Mappings = []*schema.Mapping{{ID: "M", PointID: p.ID, TagID: "T", Enabled: true}}
	require.NoError(t, f.service.ApplyWorkspaceProjection(t.Context(), f.projection))
	require.Len(t, f.service.mappingIndex[p.ID], 1)
	old := f.service.mappingIndex[p.ID][0]
	require.True(t, f.service.SuspendPointRevision(p.ID, p.RevisionID))
	bindingProof, ok := f.service.CaptureFormatRecoveryBindings(p.ID)
	require.True(t, ok)
	p.RevisionID = "new"
	p.FormatOverride = "CDAB"
	p.DataFormat = "CDAB"
	profile, err := f.service.ResolvePointReadProfile(t.Context(), p)
	require.NoError(t, err)
	proof, ok := f.service.CaptureOrdinaryPoint(p.ID)
	require.True(t, ok)
	lease, ok := f.service.InstallPreparedFormatRecovery(p, proof, profile)
	require.True(t, ok)
	tags := map[string]*schema.Tag{"T": f.projection.Tags[0]}
	prepared, err := PrepareFormatRecoveryBindings(f.projection.WorkspaceID, *p, profile, *f.projection.Devices[0], f.projection.Mappings, tags)
	require.NoError(t, err)
	_, ok = f.service.ActivateFormatRecoveryBindings(p.ID, lease, profile, bindingProof, prepared)
	require.True(t, ok)
	values := f.service.scheduler.PollNowContext(t.Context(), []string{p.ID})
	require.Len(t, values, 1)
	require.NotEqual(t, old.ConfigFingerprint, values[0].ConfigFingerprint)
	sample, err := typedSampleFromValue(values[0], f.service.mappingIndex[p.ID][0], values[0].Value, values[0].Value)
	require.NoError(t, err, "recovered installation must supply current mandatory typed source binding")
	require.NotEqual(t, old.SourceRevision, sample.SourceRevision)
	require.Equal(t, values[0].ConfigFingerprint, f.service.mappingIndex[p.ID][0].ConfigFingerprint)
	sink := &sampleSinkCapture{}
	f.service.sampleSink = sink
	f.service.writer = &mockWriter{}
	f.service.handleCollectedValue(t.Context(), values[0])
	require.Len(t, sink.samples, 1)
	require.Equal(t, sample.SourceRevision, sink.samples[0].SourceRevision)
	require.Equal(t, values[0].ConfigFingerprint, sink.samples[0].ConfigFingerprint)
}

func recoveryBindingFixture(t *testing.T) (*Service, *schema.Point, connector.ReadProfile, FormatRecoveryBindingProof, PreparedFormatRecoveryBindings, uint64) {
	t.Helper()
	f := newAcquisitionInstallFixture(t, false)
	p := f.projection.Points[0]
	p.FormatMode = schema.PointFormatOverride
	p.FormatOverride = "ABCD"
	p.DataFormat = "ABCD"
	p.FormatVersion = 1
	p.RevisionID = "before"
	p.DataType = schema.DataTypeUint32
	f.projection.RuleLinks[0].Address = p.Address
	f.projection.Tags = []*schema.Tag{{ID: "T", Key: "T", DataType: schema.DataTypeUint32}}
	f.projection.Mappings = []*schema.Mapping{{ID: "M", PointID: p.ID, TagID: "T", Enabled: true}}
	require.NoError(t, f.service.ApplyWorkspaceProjection(t.Context(), f.projection))
	bindingProof, ok := f.service.CaptureFormatRecoveryBindings(p.ID)
	require.True(t, ok)
	require.True(t, f.service.SuspendPointRevision(p.ID, p.RevisionID))
	p.RevisionID = "after"
	p.FormatOverride = "CDAB"
	p.DataFormat = "CDAB"
	profile, err := f.service.ResolvePointReadProfile(t.Context(), p)
	require.NoError(t, err)
	proof, ok := f.service.CaptureOrdinaryPoint(p.ID)
	require.True(t, ok)
	prepared, err := PrepareFormatRecoveryBindings(f.projection.WorkspaceID, *p, profile, *f.projection.Devices[0], f.projection.Mappings, map[string]*schema.Tag{"T": f.projection.Tags[0]})
	require.NoError(t, err)
	lease, ok := f.service.InstallPreparedFormatRecovery(p, proof, profile)
	require.True(t, ok)
	return f.service, p, profile, bindingProof, prepared, lease
}

func TestFormatRecoveryBindingPublicationRejectsForeignInstall(t *testing.T) {
	for _, kind := range []string{"binding", "same-value-new-binding", "lease"} {
		t.Run(kind, func(t *testing.T) {
			s, p, profile, proof, prepared, lease := recoveryBindingFixture(t)
			foreignLease := uint64(0)
			s.mappingMu.Lock()
			foreign := slices.Clone(s.mappingIndex[p.ID])
			if kind == "binding" {
				foreign[0].MappingRevision = "foreign"
			}
			if kind != "lease" {
				s.mappingIndex[p.ID] = foreign
			}
			s.mappingMu.Unlock()
			if kind == "lease" {
				foreignLease = s.InstallPointForMutation(p)
				require.NotZero(t, foreignLease)
			}
			_, ok := s.ActivateFormatRecoveryBindings(p.ID, lease, profile, proof, prepared)
			require.False(t, ok)
			if kind == "lease" {
				require.Equal(t, foreignLease, s.CapturePointLease(p.ID, p.RevisionID))
			} else {
				require.Equal(t, foreign, s.mappingIndex[p.ID])
				require.Empty(t, s.scheduler.PollNowContext(t.Context(), []string{p.ID}))
			}
		})
	}
}

func TestFormatRecoveryBindingCleanupOnlyRemovesOwnInstall(t *testing.T) {
	for _, kind := range []string{"own", "foreign-binding", "foreign-lease"} {
		t.Run(kind, func(t *testing.T) {
			s, p, profile, proof, prepared, lease := recoveryBindingFixture(t)
			receipt, ok := s.ActivateFormatRecoveryBindings(p.ID, lease, profile, proof, prepared)
			require.True(t, ok)
			var foreign []mappingBinding
			var foreignLease uint64
			if kind != "own" {
				s.mappingMu.Lock()
				foreign = slices.Clone(s.mappingIndex[p.ID])
				foreign[0].MappingRevision = "foreign"
				s.mappingIndex[p.ID] = foreign
				s.mappingMu.Unlock()
			}
			if kind == "foreign-lease" {
				foreignLease = s.InstallPointForMutation(p)
				require.NotZero(t, foreignLease)
			}
			stopped := s.SuspendFormatRecovery(p.ID, lease, receipt)
			if kind == "foreign-lease" {
				require.False(t, stopped)
				require.Equal(t, foreignLease, s.CapturePointLease(p.ID, p.RevisionID))
			} else {
				require.True(t, stopped)
				require.Zero(t, s.CapturePointLease(p.ID, p.RevisionID))
			}
			if kind == "own" {
				require.Empty(t, s.mappingIndex[p.ID])
			} else {
				require.Equal(t, foreign, s.mappingIndex[p.ID])
			}
		})
	}
}

func TestFormatRecoveryAcceptedRawRetainsCapturedBinding(t *testing.T) {
	s, p, profile, proof, prepared, lease := recoveryBindingFixture(t)
	_, ok := s.ActivateFormatRecoveryBindings(p.ID, lease, profile, proof, prepared)
	require.True(t, ok)
	cv := s.scheduler.PollNowContext(t.Context(), []string{p.ID})[0]
	captured, accepted := s.captureCollectedBindings(cv)
	require.True(t, accepted)
	require.Len(t, captured, 1)
	s.mappingMu.Lock()
	next := slices.Clone(s.mappingIndex[p.ID])
	next[0].SourceRevision = "foreign"
	s.mappingIndex[p.ID] = next
	s.mappingMu.Unlock()
	sample, err := typedSampleFromValue(cv, captured[0], cv.Value, cv.Value)
	require.NoError(t, err)
	require.NotEqual(t, "foreign", sample.SourceRevision)
}
