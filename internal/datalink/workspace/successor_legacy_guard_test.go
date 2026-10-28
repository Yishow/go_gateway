package workspace

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func legacyHistoricalFixture(t *testing.T) (*circuitProfileHarness, string, CutoverRequest) {
	t.Helper()
	h, id, req := successorActiveWideFixture(t)
	review, err := h.service.SuccessorReview(t.Context(), id, req)
	require.NoError(t, err)
	require.Empty(t, review.Blockers)
	reserved, err := h.service.ReserveSuccessor(t.Context(), id, WriteGroupSuccessorReservationRequest{Envelope: review.Envelope, ReviewDigest: review.ReviewDigest})
	require.NoError(t, err)
	_, err = h.db.Exec(`UPDATE write_group_lineages SET current_group_id=?,reserved_group_id=NULL,reservation_operation_id=NULL,reservation_intent_digest=NULL,revision=revision+1 WHERE lineage_id=?`, reserved.Group.ID, reserved.LineageID)
	require.NoError(t, err)
	return h, id, CutoverRequest{Kind: WriteGroupCutoverForward, WorkspaceID: req.WorkspaceID, ExpectedWorkspaceRevision: reserved.WorkspaceRevision,
		SourceGroupID: id, ExpectedSourceRevision: req.ExpectedPredecessorRevision, ExpectedConnectorRevision: req.ExpectedPredecessorConnectorRevision,
		Boundary: req.Boundary}
}

func TestSuccessorLegacyGuardRejectsHistoricalGeneration(t *testing.T) {
	for _, kind := range []string{WriteGroupCutoverForward, WriteGroupCutoverRollback} {
		for _, status := range []string{"ready", "disabled", "deleted"} {
			t.Run(kind+"/"+status, func(t *testing.T) {
				h, id, req := legacyHistoricalFixture(t)
				_, err := h.db.Exec(`UPDATE write_groups SET status=?,revision='new-unapplied-draft' WHERE id=?`, status, id)
				require.NoError(t, err)
				req.Kind, req.ExpectedSourceRevision = kind, "new-unapplied-draft"
				if kind == WriteGroupCutoverRollback {
					req.CutoverID = "must-not-load-before-owner-guard"
				}
				inspector := &writeGroupReadinessInspector{}
				h.service.WithTableInspector(inspector)
				before := finalizedState(t, h.db)
				result, err := h.service.ReviewCutover(t.Context(), req)
				require.ErrorIs(t, err, ErrWriteGroupLifecycleBlocked)
				require.Nil(t, result)
				require.Zero(t, inspector.calls)
				require.Equal(t, before, finalizedState(t, h.db))
			})
		}
	}
}

func TestSuccessorLegacyGuardRejectsCorruptLineage(t *testing.T) {
	h, _, req := legacyHistoricalFixture(t)
	_, err := h.db.Exec(`UPDATE write_group_lineages SET workspace_id='foreign'`)
	require.NoError(t, err)
	result, err := h.service.ReviewCutover(t.Context(), req)
	require.ErrorIs(t, err, ErrWriteGroupLineageCorrupt)
	require.Nil(t, result)
}

func TestSuccessorLegacyGuardRejectsOriginBeforeProbe(t *testing.T) {
	for _, intent := range []string{`{"version":"go-gateway/successor-activation-intent/v1","predecessor_group_id":"unused","confirmation":{}}`, `{"version":"go-gateway/successor-activation-intent/v99"}`, `{"version":"go-gateway/successor-activation-intent/`} {
		t.Run(intent, func(t *testing.T) {
			f, _, req := counterCutoverFixture(t)
			forward, err := f.service.ApplyCutover(t.Context(), req)
			require.NoError(t, err)
			rollback := reviewCounterRollback(t, f, forward)
			_, err = f.db.Exec(`UPDATE write_group_cutovers SET intent=? WHERE id=?`, intent, forward.ID)
			require.NoError(t, err)
			probes := 0
			f.inspector.mutate = func() { probes++ }
			before := finalizedState(t, f.db)
			result, err := f.service.ReviewCutover(t.Context(), rollback)
			require.ErrorIs(t, err, ErrCutoverNotApplicable)
			require.Nil(t, result)
			require.Zero(t, probes)
			require.Equal(t, before, finalizedState(t, f.db))
		})
	}
}

func TestSuccessorLegacyGuardRechecksOriginInTransaction(t *testing.T) {
	f, _, req := counterCutoverFixture(t)
	forward, err := f.service.ApplyCutover(t.Context(), req)
	require.NoError(t, err)
	rollback := reviewCounterRollback(t, f, forward)
	plan, err := f.service.planRollback(t.Context(), rollback)
	require.NoError(t, err)
	require.True(t, plan.review.Applicable)
	_, err = f.db.Exec(`UPDATE write_group_cutovers SET intent=? WHERE id=?`, `{"version":"go-gateway/successor-activation-intent/v1"}`, forward.ID)
	require.NoError(t, err)
	before := finalizedState(t, f.db)
	tx, err := f.db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	record, err := readWorkspaceScope(t.Context(), tx, f.service.repo.query, req.WorkspaceID)
	require.NoError(t, err)
	result, err := f.service.applyCutoverInTx(t.Context(), tx, record, plan)
	require.ErrorIs(t, err, ErrCutoverNotApplicable)
	require.Nil(t, result)
	require.NoError(t, tx.Rollback())
	require.Equal(t, before, finalizedState(t, f.db))
}

func TestSuccessorLegacyGuardRechecksHeadInTransaction(t *testing.T) {
	h, id, req := legacyHistoricalFixture(t)
	source, err := h.service.repo.Get(t.Context(), req.WorkspaceID, id)
	require.NoError(t, err)
	h.service.now = func() time.Time { return req.Boundary.Add(-time.Minute) }
	plan := &cutoverPlan{req: req, source: source}
	tx, err := h.db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	record, err := readWorkspaceScope(t.Context(), tx, h.service.repo.query, req.WorkspaceID)
	require.NoError(t, err)
	result, err := h.service.applyCutoverInTx(t.Context(), tx, record, plan)
	require.ErrorIs(t, err, ErrWriteGroupLifecycleBlocked)
	require.Nil(t, result)
}

func TestSuccessorLegacyGuardPreservesCompletedLegacyReplay(t *testing.T) {
	f, _, req := counterCutoverFixture(t)
	completed, err := f.service.ApplyCutover(t.Context(), req)
	require.NoError(t, err)
	// Seed a later valid lineage generation. The completed legacy event remains
	// immutable; returning its exact result must not reacquire any ownership.
	tx, err := f.db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	meta, found, err := f.service.repo.loadWriteGroupLineage(t.Context(), tx, req.WorkspaceID, f.group.ID)
	require.NoError(t, err)
	require.True(t, found)
	next := cloneWriteGroup(f.group)
	next.ID, next.Revision, next.AppliedRevision, next.Status = "later-generation", "later-revision", "", WriteGroupStatusDraft
	require.NoError(t, f.service.repo.insertGroup(t.Context(), tx, next))
	_, err = tx.Exec(`INSERT INTO write_group_generations(group_id,lineage_id,generation,predecessor_group_id) VALUES(?,?,2,?)`, next.ID, meta.LineageID, f.group.ID)
	require.NoError(t, err)
	_, err = tx.Exec(`UPDATE write_group_lineages SET current_group_id=?,revision=revision+1 WHERE lineage_id=?`, next.ID, meta.LineageID)
	require.NoError(t, err)
	require.ErrorIs(t, f.service.repo.requireCurrentCanonicalOwnerInTx(t.Context(), tx, req.WorkspaceID, f.group.ID), ErrWriteGroupLifecycleBlocked)
	require.NoError(t, tx.Commit())
	before := finalizedState(t, f.db)
	replay, err := f.service.ApplyCutover(t.Context(), req)
	require.NoError(t, err)
	require.Equal(t, completed.ID, replay.ID)
	require.Equal(t, before, finalizedState(t, f.db))
}

func TestSuccessorLegacyGuardCannotReschedulePendingCurrentGeneration(t *testing.T) {
	for _, offset := range []time.Duration{-time.Second, 0, time.Second} {
		t.Run(offset.String(), func(t *testing.T) {
			h, _, req := legacyHistoricalFixture(t)
			var currentID string
			require.NoError(t, h.db.QueryRow(`SELECT current_group_id FROM write_group_lineages`).Scan(&currentID))
			current, err := h.service.repo.Get(t.Context(), req.WorkspaceID, currentID)
			require.NoError(t, err)
			current.AppliedRevision, current.Status = current.Revision, WriteGroupStatusReady
			tx, err := h.db.BeginTx(t.Context(), nil)
			require.NoError(t, err)
			require.NoError(t, h.service.repo.updateGroup(t.Context(), tx, current))
			require.NoError(t, h.service.repo.insertVersionInTx(t.Context(), tx, current, req.Boundary))
			require.NoError(t, tx.Commit())
			activation := req.Boundary
			h.service.now = func() time.Time { return activation.Add(offset) }
			req.SourceGroupID, req.ExpectedSourceRevision = current.ID, current.Revision
			req.Boundary = activation.Add(time.Minute)
			before := finalizedState(t, h.db)
			result, err := h.service.ReviewCutover(t.Context(), req)
			if offset < 0 {
				require.ErrorIs(t, err, ErrWriteGroupLifecycleBlocked)
				require.Nil(t, result)
			} else {
				require.NoError(t, err, "the pending guard ends at the immutable activation boundary")
				require.NotNil(t, result)
			}
			require.Equal(t, before, finalizedState(t, h.db))
		})
	}
}
