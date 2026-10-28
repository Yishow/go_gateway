package workspace

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRecreateIntentIgnoresOnlyServerOwnedNoise(t *testing.T) {
	h, req := recreateRepositoryFixture(t, "profile")
	original := cloneWriteGroup(req.Group)
	normalized, digest, err := normalizeRecreateReservationRequest(req)
	require.NoError(t, err)
	noisy := req
	noisy.Group = cloneWriteGroup(req.Group)
	noisy.Group.ID = "caller-id"
	noisy.Group.Revision = "caller-revision"
	noisy.Group.AppliedRevision = "caller-applied"
	noisy.Group.BasicManagedDeviceID = "caller-device"
	noisy.Group.Status = WriteGroupStatusRunning
	noisy.Group.CreatedAt = time.Now()
	noisy.Group.UpdatedAt = time.Now()
	noisy.Group.Destination.SchemaRevision = "proof"
	noisy.Group.Destination.SchemaDigest = "proof-digest"
	noisy.Group.Migration = WriteGroupMigration{SourceKind: "caller", SourceIDs: []string{"old"}}
	noisy.Group.Name = "  " + noisy.Group.Name + "  "
	_, same, err := normalizeRecreateReservationRequest(noisy)
	require.NoError(t, err)
	require.Equal(t, digest, same)
	require.Equal(t, original, req.Group, "normalization must not mutate frozen caller request")
	require.Empty(t, normalized.Group.ID)
	require.Empty(t, normalized.Group.AppliedRevision)
	first, err := h.service.reserveRecreatedDraft(t.Context(), req)
	require.NoError(t, err)
	before := lineageReadSnapshot(t, h)
	replay, err := h.service.reserveRecreatedDraft(t.Context(), noisy)
	require.NoError(t, err)
	require.Equal(t, first, replay)
	require.Equal(t, before, lineageReadSnapshot(t, h))
}

func TestRecreateServiceFreshCASAndScopeRefusalsHaveNoEffects(t *testing.T) {
	cases := map[string]func(*writeGroupRecreateReservationRequest){
		"workspace":              func(r *writeGroupRecreateReservationRequest) { r.ExpectedWorkspaceRevision = "stale" },
		"lineage":                func(r *writeGroupRecreateReservationRequest) { r.ExpectedLineageRevision++ },
		"predecessor":            func(r *writeGroupRecreateReservationRequest) { r.ExpectedPredecessorRevision = "stale" },
		"saved_source_connector": func(r *writeGroupRecreateReservationRequest) { r.ExpectedPredecessorConnectorRevision = "stale" },
		"target_connector": func(r *writeGroupRecreateReservationRequest) {
			r.ExpectedTargetConnectorRevision = "stale"
			r.Group.Destination.ConnectorRevision = "stale"
		},
		"foreign_workspace": func(r *writeGroupRecreateReservationRequest) {
			r.WorkspaceID = "foreign"
			r.Group.WorkspaceID = "foreign"
		},
		"foreign_lineage": func(r *writeGroupRecreateReservationRequest) { r.LineageID = "foreign" },
	}
	for name, change := range cases {
		t.Run(name, func(t *testing.T) {
			h, req := recreateRepositoryFixture(t, "profile")
			change(&req)
			before := lineageReadSnapshot(t, h)
			_, err := h.service.reserveRecreatedDraft(t.Context(), req)
			require.Error(t, err)
			require.True(t, errors.Is(err, ErrWriteGroupRevisionConflict) || errors.Is(err, ErrWriteGroupConnectorRevisionConflict))
			require.Equal(t, before, lineageReadSnapshot(t, h))
		})
	}
}
