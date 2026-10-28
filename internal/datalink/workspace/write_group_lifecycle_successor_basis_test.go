package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLifecycleSuccessorSeedUsesImmutableAppliedMeaning(t *testing.T) {
	for _, kind := range []string{"basic", "profile"} {
		t.Run(kind, func(t *testing.T) {
			var h *circuitProfileHarness
			var req writeGroupRecreateReservationRequest
			if kind == "basic" {
				h, req = recreateBasicFixture(t, "sqlite")
			} else {
				h, req = recreateRepositoryFixture(t, "profile")
			}
			applied := seedRecreateEngineEvidence(t, h, req)
			h.service.WithBacklogOwnershipGuard(writeGroupBacklogGuardFunc(func(context.Context, *sql.Tx, *WriteGroup) error { return nil }))
			_, err := h.db.Exec(`UPDATE write_groups SET name='unapplied draft',row_policy=json_set(row_policy,'$.interval_seconds',77),revision='draft-R2' WHERE id=?`, req.PredecessorGroupID)
			require.NoError(t, err)
			_, err = h.db.Exec(`UPDATE write_group_versions SET payload=json_set(payload,'$.destination.schema_revision','private-schema-proof','$.destination.schema_digest','private-schema-digest','$.migration.review_result','private-review') WHERE group_id=?`, req.PredecessorGroupID)
			require.NoError(t, err)
			before := recreateRepositorySnapshot(t, h)
			view, err := h.service.Lifecycle(t.Context(), req.PredecessorGroupID)
			require.NoError(t, err)
			require.True(t, view.Actions[3].Allowed, "%+v", view.Actions[3])
			require.Equal(t, "draft-R2", view.GroupRevision)
			seed := view.SuccessorSeed
			require.NotNil(t, seed)
			require.Equal(t, "applied", seed.BasisKind)
			require.Equal(t, applied.Group.Revision, seed.BasisRevision)
			require.Equal(t, applied.Group.RowPolicy, seed.Group.RowPolicy)
			require.Equal(t, applied.Group.WritePolicy, seed.Group.WritePolicy)
			require.Equal(t, applied.Group.Members, seed.Group.Members)
			require.Equal(t, applied.Group.RecordingProfile, seed.Group.RecordingProfile)
			require.Equal(t, applied.Group.Name, seed.Group.Name)
			raw, err := json.Marshal(seed)
			require.NoError(t, err)
			for _, private := range []string{"private-schema", "private-review", "schema_token", "schema_revision", "schema_digest", "runtime_layout", `"capability"`, "connection_config", "ledger", "raw_sql"} {
				require.NotContains(t, string(raw), private)
			}
			reviewReq := WriteGroupSuccessorReviewRequest{WorkspaceID: view.WorkspaceID, ExpectedWorkspaceRevision: view.WorkspaceRevision, LineageID: view.History.LineageID, ExpectedLineageRevision: view.History.LineageRevision, ExpectedPredecessorRevision: view.GroupRevision, ExpectedPredecessorConnectorRevision: view.ConnectorRevision, ExpectedTargetConnectorRevision: seed.Group.Destination.ConnectorRevision, OperationID: "seed-review", Boundary: time.Date(2026, 10, 28, 12, 0, 0, 0, time.UTC), Group: seed.Group}
			review, err := h.service.SuccessorReview(t.Context(), view.GroupID, reviewReq)
			require.NoError(t, err)
			require.Empty(t, review.Blockers, "%+v", review.Blockers)
			for _, mutation := range []string{"foreign", "meaning", "stale"} {
				bad := reviewReq
				bad.Group = cloneWriteGroup(reviewReq.Group)
				switch mutation {
				case "foreign":
					bad.Group.WorkspaceID = "foreign"
				case "meaning":
					bad.Group.Members[0].SourceRevision = "tampered"
				case "stale":
					bad.ExpectedPredecessorRevision = "stale"
				}
				rejected, err := h.service.SuccessorReview(t.Context(), view.GroupID, bad)
				require.True(t, err != nil || len(rejected.Blockers) != 0, "seed never grants authority: %s", mutation)
			}
			require.Equal(t, before, recreateRepositorySnapshot(t, h))
		})
	}
}

func TestLifecycleSuccessorSeedProvesNeverApplied(t *testing.T) {
	h, req := recreateBasicFixture(t, "sqlite")
	view, err := h.service.Lifecycle(t.Context(), req.PredecessorGroupID)
	require.NoError(t, err)
	require.True(t, view.Actions[3].Allowed)
	require.Equal(t, "saved", view.SuccessorSeed.BasisKind)
	require.Equal(t, req.ExpectedPredecessorRevision, view.SuccessorSeed.BasisRevision)
	seedRecreateEngineEvidence(t, h, req)
	_, err = h.db.Exec(`UPDATE write_groups SET applied_revision='' WHERE id=?`, req.PredecessorGroupID)
	require.NoError(t, err)
	view, err = h.service.Lifecycle(t.Context(), req.PredecessorGroupID)
	require.NoError(t, err)
	require.False(t, view.Actions[3].Allowed)
	require.Equal(t, "source_definition_unavailable", view.Actions[3].ReasonCode)
	require.Nil(t, view.SuccessorSeed, "empty applied pointer cannot disguise historical evidence as a saved-only basis")
}

func TestLifecycleSuccessorRealUpdateKeepsDraftTransitionBlocked(t *testing.T) {
	h, id, _ := successorActiveWideFixture(t)
	original, err := h.service.Get(t.Context(), id)
	require.NoError(t, err)
	draft := cloneWriteGroup(original.Group)
	draft.Name = "Saved changes are not active"
	request := lifecycleReadMutation(t, h.service, original.Group)
	request.Group = draft
	saved, err := h.service.Update(t.Context(), id, request)
	require.NoError(t, err)
	require.Equal(t, WriteGroupStatusDraft, saved.Group.Status)
	require.Equal(t, original.Group.AppliedRevision, saved.Group.AppliedRevision)
	view, err := h.service.Lifecycle(t.Context(), id)
	require.NoError(t, err)
	require.Nil(t, view.SuccessorSeed)
	for _, action := range view.Actions[3:] {
		require.False(t, action.Allowed)
		require.Equal(t, "source_draft_not_active", action.ReasonCode)
	}
	require.True(t, view.Actions[0].Allowed)
	require.True(t, view.Actions[1].Allowed, "ordinary safe stop is not a successor action")
}
