package workspace

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSuccessorSQLServerUnitActivationEngineGate(t *testing.T) {
	group, _ := validSQLServerTestFixture(t)
	for _, mode := range []WriteGroupStorageMode{WriteGroupStorageModePlain, WriteGroupStorageModeNativeMonthly} {
		group.StorageMode = mode
		require.True(t, successorActivationEngineSupported(group, "sqlserver"), "%s", mode)
	}
	group.StorageMode = WriteGroupStorageModeTimescale
	require.False(t, successorActivationEngineSupported(group, "sqlserver"))
	group.StorageMode, group.LayoutVersion = WriteGroupStorageModePlain, WriteGroupLayoutManagedWideProfileV1
	require.False(t, successorActivationEngineSupported(group, "sqlserver"))
}

func TestSuccessorSQLServerUnitReadinessUsesGuardedReload(t *testing.T) {
	f := newSQLServerApplyCompatibilityFixture(t, WriteGroupStorageModePlain)
	snapshot, err := f.h.service.readinessSnapshot(t.Context(), f.group.ID, nil)
	require.NoError(t, err)
	called := false
	stale := errors.New("reservation changed while inspecting")
	_, err = f.h.service.evaluateReadinessSnapshot(t.Context(), f.group.ID, nil, snapshot, func(context.Context) (*writeGroupReadinessSnapshot, error) {
		called = true
		return nil, stale
	})
	require.True(t, called, "SQLServer must use caller's operation-bound reloader")
	require.ErrorIs(t, err, stale)
}

func TestSuccessorSQLServerUnitLocalReviewUsesSucceededToken(t *testing.T) {
	for _, mode := range []WriteGroupStorageMode{WriteGroupStorageModePlain, WriteGroupStorageModeNativeMonthly} {
		t.Run(string(mode), func(t *testing.T) {
			h, id, req, p := successorSQLServerUnitConfirmed(t, mode)
			calls, catalogs := p.calls, p.catalogCalls
			p.err = errors.New("offline: local review may not probe")
			before := lineageReadSnapshot(t, h)
			view, err := h.service.SuccessorActivationReview(t.Context(), id, req)
			require.NoError(t, err)
			require.Empty(t, view.Blockers)
			require.Equal(t, calls, p.calls)
			require.Equal(t, catalogs, p.catalogCalls)
			require.Equal(t, before, lineageReadSnapshot(t, h))
		})
	}
}

func TestSuccessorSQLServerUnitSchemaConfirmationRequiresFreshFacts(t *testing.T) {
	for _, mode := range []WriteGroupStorageMode{WriteGroupStorageModePlain, WriteGroupStorageModeNativeMonthly} {
		t.Run(string(mode), func(t *testing.T) {
			h, id, stage, p := successorSQLServerUnitReserved(t, mode)
			scope, err := h.service.SuccessorSchemaScope(t.Context(), id, stage)
			require.NoError(t, err)
			token, digest := successorSQLServerUnitSucceeded(t, h, scope, p)
			p.capability.CapabilityRevision = "changed-after-success"
			before := lineageReadSnapshot(t, h)
			_, err = h.service.ConfirmSuccessorSchema(t.Context(), id, stage, token.OperationID, digest, &StorageEvidence{Layout: *token.GroupLayout, CreatedTable: true})
			require.Error(t, err, "a succeeded local token cannot replace fresh target verification")
			require.Equal(t, before, lineageReadSnapshot(t, h))
		})
	}
}

func TestSuccessorSQLServerUnitSourceDescriptorMustValidate(t *testing.T) {
	h, id, req, p := successorSQLServerUnitFixture(t, WriteGroupStorageModePlain)
	_, err := h.db.Exec(`UPDATE wg_runtime_versions SET payload=json_remove(payload,'$.runtime_layout.sqlserver') WHERE group_id=?`, id)
	require.NoError(t, err)
	before := lineageReadSnapshot(t, h)
	view, err := h.service.SuccessorReview(t.Context(), id, req)
	require.NoError(t, err)
	require.Contains(t, view.Blockers, WriteGroupSuccessorBlocker{ReasonCode: "source-engine-evidence-unavailable", NextAction: "review_lifecycle"}, "present corrupt descriptor must not fall back to SQLServer driver label")
	require.Zero(t, p.calls)
	require.Equal(t, before, lineageReadSnapshot(t, h))
}
