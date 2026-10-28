package workspace

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func legacyAdvanceHolderHead(t *testing.T, f counterApplyFixture, holder *WriteGroup) {
	t.Helper()
	tx, err := f.db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	meta, found, err := f.service.repo.loadWriteGroupLineage(t.Context(), tx, holder.WorkspaceID, holder.ID)
	require.NoError(t, err)
	require.True(t, found)
	next := cloneWriteGroup(holder)
	next.ID, next.Revision, next.AppliedRevision, next.Status = "later-holder", "later-holder-revision", "", WriteGroupStatusDraft
	next.Destination.SchemaRevision, next.Destination.SchemaDigest = "", ""
	if next.isWideProfile() {
		next.Destination.TableName = managedGroupTableName(next.ID)
	}
	require.NoError(t, f.service.repo.insertGroup(t.Context(), tx, next))
	_, err = tx.Exec(`INSERT INTO write_group_generations(group_id,lineage_id,generation,predecessor_group_id) VALUES(?,?,2,?)`, next.ID, meta.LineageID, holder.ID)
	require.NoError(t, err)
	_, err = tx.Exec(`UPDATE write_group_lineages SET current_group_id=?,revision=revision+1 WHERE lineage_id=?`, next.ID, meta.LineageID)
	require.NoError(t, err)
	require.NoError(t, tx.Commit())
}

func TestSuccessorLegacyGuardRejectsHistoricalTargetBeforeProbe(t *testing.T) {
	f, target, req := counterCutoverFixture(t)
	legacyAdvanceHolderHead(t, f, target)
	probes := 0
	f.inspector.mutate = func() { probes++ }
	before := finalizedState(t, f.db)
	result, err := f.service.ReviewCutover(t.Context(), req)
	require.ErrorIs(t, err, ErrWriteGroupLifecycleBlocked)
	require.Nil(t, result)
	require.Zero(t, probes)
	require.Equal(t, before, finalizedState(t, f.db))
}

func TestSuccessorLegacyGuardRejectsHistoricalRollbackHolder(t *testing.T) {
	f, target, req := counterCutoverFixture(t)
	forward, err := f.service.ApplyCutover(t.Context(), req)
	require.NoError(t, err)
	rollback := reviewCounterRollback(t, f, forward)
	legacyAdvanceHolderHead(t, f, target)
	probes := 0
	f.inspector.mutate = func() { probes++ }
	before := finalizedState(t, f.db)
	result, err := f.service.ReviewCutover(t.Context(), rollback)
	require.ErrorIs(t, err, ErrWriteGroupLifecycleBlocked)
	require.Nil(t, result)
	require.Zero(t, probes)
	require.Equal(t, before, finalizedState(t, f.db))
}

func TestSuccessorLegacyGuardRechecksRollbackHoldersInTransaction(t *testing.T) {
	f, target, req := counterCutoverFixture(t)
	forward, err := f.service.ApplyCutover(t.Context(), req)
	require.NoError(t, err)
	rollback := reviewCounterRollback(t, f, forward)
	plan, err := f.service.planRollback(t.Context(), rollback)
	require.NoError(t, err)
	require.True(t, plan.review.Applicable)
	legacyAdvanceHolderHead(t, f, target)
	tx, err := f.db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	record, err := readWorkspaceScope(t.Context(), tx, f.service.repo.query, req.WorkspaceID)
	require.NoError(t, err)
	result, err := f.service.applyCutoverInTx(t.Context(), tx, record, plan)
	require.ErrorIs(t, err, ErrWriteGroupLifecycleBlocked)
	require.Nil(t, result)
}

func TestSuccessorLegacyGuardRechecksPriorSourceScopeInTransaction(t *testing.T) {
	f, target, req := counterCutoverFixture(t)
	forward, err := f.service.ApplyCutover(t.Context(), req)
	require.NoError(t, err)
	rollback := reviewCounterRollback(t, f, forward)
	plan, err := f.service.planRollback(t.Context(), rollback)
	require.NoError(t, err)
	_, err = f.db.Exec(`UPDATE write_group_cutovers SET source_group_id=? WHERE id=?`, target.ID, forward.ID)
	require.NoError(t, err)
	tx, err := f.db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	record, err := readWorkspaceScope(t.Context(), tx, f.service.repo.query, req.WorkspaceID)
	require.NoError(t, err)
	result, err := f.service.applyCutoverInTx(t.Context(), tx, record, plan)
	require.ErrorIs(t, err, ErrCutoverNotApplicable)
	require.Nil(t, result)
}
