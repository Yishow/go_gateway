package workspace

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSuccessorSQLServerUnitUnsupportedSourceAndTargetHaveNoEffects(t *testing.T) {
	for _, mutation := range []string{"wide", "cross-engine", "cross-mode", "partial", "meaning", "never-applied"} {
		t.Run(mutation, func(t *testing.T) {
			h, id, req, p := successorSQLServerUnitFixture(t, WriteGroupStorageModePlain)
			switch mutation {
			case "wide":
				req.Group.LayoutVersion, req.Group.StorageMode = WriteGroupLayoutManagedWideProfileV1, ""
			case "cross-engine":
				_, err := h.db.Exec(`UPDATE database_connectors SET kind='postgres' WHERE id=?`, req.Group.Destination.ConnectorID)
				require.NoError(t, err)
			case "cross-mode":
				req.Group.StorageMode = WriteGroupStorageModeNativeMonthly
			case "partial":
				req.Group.Members = nil
			case "meaning":
				req.Group.RecordingProfile.Interval++
			case "never-applied":
				_, err := h.db.Exec(`UPDATE write_groups SET applied_revision='' WHERE id=?`, id)
				require.NoError(t, err)
			}
			before := lineageReadSnapshot(t, h)
			view, err := h.service.SuccessorReview(t.Context(), id, req)
			if err == nil {
				require.NotEmpty(t, view.Blockers)
			}
			require.Zero(t, p.calls)
			require.Equal(t, before, lineageReadSnapshot(t, h))
		})
	}
}

func TestSuccessorSQLServerUnitActivationRollbackKeepsAllOwners(t *testing.T) {
	for _, mode := range []WriteGroupStorageMode{WriteGroupStorageModePlain, WriteGroupStorageModeNativeMonthly} {
		t.Run(string(mode), func(t *testing.T) {
			h, id, req, _ := successorSQLServerUnitConfirmed(t, mode)
			view, err := h.service.SuccessorActivationReview(t.Context(), id, req)
			require.NoError(t, err)
			confirmation := WriteGroupSuccessorActivationConfirmation{Envelope: view.Envelope, ReviewDigest: view.ReviewDigest}
			ready, err := h.service.SuccessorActivationReadiness(t.Context(), id, confirmation)
			require.NoError(t, err)
			require.True(t, ready.Readiness.Ready)
			confirmation.ReadinessDigest = ready.ReadinessDigest
			_, err = h.db.Exec(`CREATE TRIGGER reject_sqlserver_activation BEFORE UPDATE ON write_group_lineages BEGIN SELECT RAISE(ABORT,'unit rollback'); END`)
			require.NoError(t, err)
			before := recreateRepositorySnapshot(t, h)
			_, err = h.service.ConfirmSuccessorActivation(t.Context(), id, confirmation)
			require.Error(t, err)
			require.Equal(t, before, recreateRepositorySnapshot(t, h))
			require.Zero(t, h.countRows(`SELECT COUNT(*) FROM write_group_cutovers`))
		})
	}
}

func TestSuccessorSQLServerUnitStaleWorkspaceRejectsBeforeProbe(t *testing.T) {
	h, id, stage, p := successorSQLServerUnitReserved(t, WriteGroupStorageModePlain)
	scope, err := h.service.SuccessorSchemaScope(t.Context(), id, stage)
	require.NoError(t, err)
	token, digest := successorSQLServerUnitSucceeded(t, h, scope, p)
	stage.Current.WorkspaceRevision = "stale-workspace"
	before, calls := lineageReadSnapshot(t, h), p.calls
	_, err = h.service.ConfirmSuccessorSchema(t.Context(), id, stage, token.OperationID, digest, &StorageEvidence{Layout: *token.GroupLayout, CreatedTable: true})
	require.Error(t, err)
	require.Equal(t, calls, p.calls)
	require.Equal(t, before, lineageReadSnapshot(t, h))
}

func TestSuccessorSQLServerUnitFailedProbeStillReloadsReservation(t *testing.T) {
	f := newSQLServerApplyCompatibilityFixture(t, WriteGroupStorageModePlain)
	snapshot, err := f.h.service.readinessSnapshot(t.Context(), f.group.ID, nil)
	require.NoError(t, err)
	f.h.service.WithManagedTableInspector(&successorSQLServerUnitInspector{err: errors.New("unit offline")})
	stale := errors.New("stale reservation")
	called := false
	_, err = f.h.service.evaluateReadinessSnapshot(t.Context(), f.group.ID, nil, snapshot, func(context.Context) (*writeGroupReadinessSnapshot, error) {
		called = true
		return nil, stale
	})
	require.True(t, called)
	require.ErrorIs(t, err, stale)
}
