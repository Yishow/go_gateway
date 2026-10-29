package workspace

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/recordingplan"
)

// Seed only local immutable checkpoint records; no external readiness or DDL.
func statusFirstStartFixture(t *testing.T) (*circuitProfileHarness, string, WriteGroupSuccessorStageRequest, *WriteGroupSuccessorFirstStartResult) {
	t.Helper()
	h, id, stage, _ := successorSchemaFixture(t)
	token, repo, ledger := statusPreview(t, h, id, stage)
	op, _, err := ledger.ClaimSchemaApply(t.Context(), token)
	require.NoError(t, err)
	digest := strings.Repeat("a", 64)
	_, err = repo.FinishSchemaOperation(t.Context(), op.OperationID, op.Owner, recordingplan.SchemaOperationResult{Status: recordingplan.SchemaOperationSucceeded, VerifiedDigest: digest})
	require.NoError(t, err)
	confirmed, err := h.service.ConfirmSuccessorSchema(t.Context(), id, stage, token.OperationID, digest, &StorageEvidence{Layout: *token.GroupLayout, CreatedTable: true})
	require.NoError(t, err)
	stage.Current.WorkspaceRevision = confirmed.WorkspaceRevision
	h.service.now = func() time.Time { return stage.Envelope.Request.Boundary.Add(-time.Hour) }
	review, err := h.service.SuccessorActivationReview(t.Context(), id, WriteGroupSuccessorActivationRequest{Stage: stage, Boundary: stage.Envelope.Request.Boundary})
	require.NoError(t, err)
	confirmation := WriteGroupSuccessorActivationConfirmation{Envelope: review.Envelope, ReviewDigest: review.ReviewDigest, ReadinessDigest: strings.Repeat("b", 64)}
	_, intent, err := successorFirstStartIntent(id, confirmation)
	require.NoError(t, err)
	rootID := recordingplan.SuccessorReserveOperationID(stage.Envelope.Request.WorkspaceID, stage.Envelope.Request.OperationID)
	result := &WriteGroupSuccessorFirstStartResult{OperationID: recordingplan.SuccessorFirstStartOperationID(stage.Envelope.Request.WorkspaceID, rootID), RootOperationID: stage.Envelope.Request.OperationID,
		WorkspaceRevision: "committed-first-start-revision", LineageID: stage.Envelope.Request.LineageID, LineageRevision: stage.Current.LineageRevision + 1,
		PredecessorGroupID: id, GroupID: stage.Current.GroupID, AppliedRevision: stage.Current.GroupRevision, Boundary: review.Envelope.Request.Boundary, Generation: stage.Envelope.Generation + 1, Status: "succeeded"}
	payload, err := json.Marshal(result)
	require.NoError(t, err)
	tx, err := h.db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	group := cloneWriteGroup(confirmed.Group)
	group.Status, group.AppliedRevision = WriteGroupStatusReady, group.Revision
	require.NoError(t, h.service.repo.insertVersionInTx(t.Context(), tx, group, result.Boundary))
	_, err = recordingplan.RecordSuccessorFirstStartInTx(t.Context(), tx, &recordingplan.SuccessorFirstStartRecord{Scope: recordingplan.SuccessorOperationScope{WorkspaceID: stage.Envelope.Request.WorkspaceID, RequestID: stage.Envelope.Request.OperationID, LineageID: stage.Envelope.Request.LineageID, PredecessorGroupID: id, SuccessorGroupID: group.ID}, ReserveOperationID: rootID, IntentDigest: review.ReviewDigest, Intent: intent, Result: payload})
	require.NoError(t, err)
	require.NoError(t, tx.Commit())
	return h, id, stage, result
}

func TestSuccessorOperationStatusFirstStartHasOwnImmutableOutcome(t *testing.T) {
	h, id, stage, expected := statusFirstStartFixture(t)
	status, err := h.service.SuccessorOperationStatus(t.Context(), id, statusRequest(stage))
	require.NoError(t, err)
	require.Equal(t, expected, status.FirstStart.Result)
	require.Equal(t, "succeeded", status.FirstStart.Status)
	require.Equal(t, "not_submitted", status.Cutover.Status)
	require.Nil(t, status.Cutover.Result)
	_, err = h.db.Exec(`UPDATE write_group_lineages SET current_group_id=?,reserved_group_id=NULL,reservation_operation_id=NULL,reservation_intent_digest=NULL WHERE lineage_id=?`, expected.GroupID, expected.LineageID)
	require.NoError(t, err)
	h.service = NewWriteGroupService(NewService(NewSQLRepository(h.db)), NewSQLWriteGroupRepository(h.db))
	after, err := h.service.SuccessorOperationStatus(t.Context(), id, statusRequest(stage))
	require.NoError(t, err)
	require.Equal(t, status, after)
	for _, table := range []string{"write_group_cutovers", "write_group_revision_ends", "write_group_ownership_fences"} {
		require.Zero(t, h.countRows("SELECT COUNT(*) FROM "+table))
	}
}

func TestSuccessorOperationStatusRejectsFirstStartMissingImmutableVersion(t *testing.T) {
	h, id, stage, _ := statusFirstStartFixture(t)
	_, err := h.db.Exec(`DELETE FROM write_group_versions`)
	require.NoError(t, err)
	status, err := h.service.SuccessorOperationStatus(t.Context(), id, statusRequest(stage))
	require.ErrorIs(t, err, ErrWriteGroupLineageCorrupt)
	require.Nil(t, status)
}
