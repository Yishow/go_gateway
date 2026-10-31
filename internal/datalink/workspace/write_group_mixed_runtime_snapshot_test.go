package workspace

import (
	"encoding/json"
	"go-gateway/internal/datalink/dbtarget"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMixedAppliedDescriptorRetainsMemberIntentAndOriginalBytes(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	ctx := t.Context()
	draft := mixedDraft(h, WriteGroupMemberCalculationMean, 15, "mixed_frozen")
	draft.LayoutVersion = WriteGroupLayoutManagedWideProfileV1
	saved, err := h.service.Create(ctx, h.mutation(ctx, draft))
	require.NoError(t, err)
	ready, err := h.service.readinessSnapshot(ctx, saved.Group.ID, nil)
	require.NoError(t, err)
	g := cloneWriteGroup(saved.Group)
	g.AppliedRevision = g.Revision
	require.NoError(t, ValidateMixedRecordingSnapshot(g, ready.tagTypes, nil))
	frozen := &WriteGroupAppliedSnapshot{WorkspaceID: g.WorkspaceID, GroupID: g.ID, GroupRevision: g.Revision, AppliedRevision: g.Revision, Group: g, RuntimeLayout: &WriteGroupRuntimeLayout{Dialect: dbtarget.SQLDialectSQLite, TagTypes: ready.tagTypes, ReceiptTableReady: true}}
	tx, err := h.db.BeginTx(ctx, nil)
	require.NoError(t, err)
	require.NoError(t, h.service.saveRuntimeVersionInTx(ctx, tx, frozen))
	require.NoError(t, tx.Commit())
	var original string
	require.NoError(t, h.db.QueryRowContext(t.Context(), `SELECT payload FROM wg_runtime_versions WHERE group_id=?`, g.ID).Scan(&original))
	var restored WriteGroupAppliedSnapshot
	require.NoError(t, json.Unmarshal([]byte(original), &restored))
	require.Equal(t, g.RecordID, restored.Group.RecordID)
	require.Equal(t, g.RecordingProfile, restored.Group.RecordingProfile)
	require.Equal(t, g.Members, restored.Group.Members)
	require.Zero(t, restored.Group.RowPolicy.IntervalSeconds)
	// Neither editing the caller nor an idempotent same-revision persistence call
	// rewrites accepted method/cadence under the old immutable identity.
	g.Members[0].Calculation = WriteGroupMemberCalculationLatest
	g.Members[0].Cadence.Interval = 60
	tx, err = h.db.BeginTx(ctx, nil)
	require.NoError(t, err)
	require.NoError(t, h.service.saveRuntimeVersionInTx(ctx, tx, frozen))
	require.NoError(t, tx.Commit())
	var after string
	require.NoError(t, h.db.QueryRowContext(t.Context(), `SELECT payload FROM wg_runtime_versions WHERE group_id=?`, g.ID).Scan(&after))
	require.Equal(t, original, after)
	// This private descriptor fixture does not make the public draft applied.
	public, err := h.service.Get(ctx, g.ID)
	require.NoError(t, err)
	require.Empty(t, public.Group.AppliedRevision)
	readiness, err := h.service.Readiness(ctx, g.ID)
	require.NoError(t, err)
	require.False(t, readiness.Ready)
}

func TestMixedHistoricalDeltaCannotManufactureOriginalProof(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	g := mixedDraft(h, WriteGroupMemberCalculationDelta, 15, "")
	g.LayoutVersion = WriteGroupLayoutManagedLongV1
	g.ID = "old-group"
	g.Revision = "old-revision"
	g.AppliedRevision = g.Revision
	require.True(t, requiresUnscaledCounterProof(g))
	tx, err := h.db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	err = h.service.saveHistoricalRuntimeVersionInTx(t.Context(), tx, &WriteGroupAppliedSnapshot{WorkspaceID: g.WorkspaceID, GroupID: g.ID, AppliedRevision: g.AppliedRevision, Group: g})
	require.ErrorIs(t, err, ErrWriteGroupApplyNotReady)
	require.NoError(t, tx.Rollback())
	g.Members[0].Calculation = WriteGroupMemberCalculationLatest
	require.False(t, requiresUnscaledCounterProof(g), "latest has no calculation baseline authority")
}
