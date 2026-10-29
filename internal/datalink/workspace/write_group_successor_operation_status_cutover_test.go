package workspace

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/recordingplan"
)

// This fixture seeds local authority records only. It performs no target SQL,
// target probes or runtime activation and does not claim end-to-end cutover proof.
func statusCutoverFixture(t *testing.T, disablePredecessor ...bool) (*circuitProfileHarness, string, WriteGroupSuccessorStageRequest, *WriteGroupCutover) {
	t.Helper()
	h, id, request := successorActiveWideFixture(t)
	_, err := h.db.Exec(`UPDATE database_connectors SET connection_config=json_set(connection_config,'$.schema','main') WHERE id=?`, request.Group.Destination.ConnectorID)
	require.NoError(t, err)
	request.Group.Destination.TableSchema = "main"
	initial, err := h.service.SuccessorReview(t.Context(), id, request)
	require.NoError(t, err)
	require.Empty(t, initial.Blockers)
	reserved, err := h.service.ReserveSuccessor(t.Context(), id, WriteGroupSuccessorReservationRequest{Envelope: initial.Envelope, ReviewDigest: initial.ReviewDigest})
	require.NoError(t, err)
	stage := WriteGroupSuccessorStageRequest{Envelope: initial.Envelope, ReviewDigest: initial.ReviewDigest, Current: WriteGroupSuccessorCurrentCAS{
		WorkspaceRevision: reserved.WorkspaceRevision, LineageRevision: reserved.LineageRevision, PredecessorRevision: request.ExpectedPredecessorRevision,
		PredecessorConnectorRevision: request.ExpectedPredecessorConnectorRevision, GroupID: reserved.Group.ID, GroupRevision: reserved.Group.Revision, TargetConnectorRevision: request.ExpectedTargetConnectorRevision}}
	token, repo, ledger := statusPreview(t, h, id, stage)
	op, _, err := ledger.ClaimSchemaApply(t.Context(), token)
	require.NoError(t, err)
	digest := strings.Repeat("a", 64)
	_, err = repo.FinishSchemaOperation(t.Context(), op.OperationID, op.Owner, recordingplan.SchemaOperationResult{Status: recordingplan.SchemaOperationSucceeded, VerifiedDigest: digest})
	require.NoError(t, err)
	confirmed, err := h.service.ConfirmSuccessorSchema(t.Context(), id, stage, token.OperationID, digest, &StorageEvidence{Layout: *token.GroupLayout, CreatedTable: true})
	require.NoError(t, err)
	stage.Current.WorkspaceRevision = confirmed.WorkspaceRevision
	if len(disablePredecessor) > 0 && disablePredecessor[0] {
		stopped, err := h.service.Disable(t.Context(), id, WriteGroupMutation{WorkspaceID: request.WorkspaceID, ExpectedWorkspaceRevision: stage.Current.WorkspaceRevision, ExpectedGroupRevision: stage.Current.PredecessorRevision, ExpectedConnectorRevision: stage.Current.PredecessorConnectorRevision})
		require.NoError(t, err)
		stage.Current.WorkspaceRevision, stage.Current.PredecessorRevision = stopped.WorkspaceRevision, stopped.Group.Revision
		require.NotEqual(t, request.ExpectedPredecessorRevision, stage.Current.PredecessorRevision)
	}
	h.service.now = func() time.Time { return stage.Envelope.Request.Boundary.Add(-time.Hour) }
	review, err := h.service.SuccessorActivationReview(t.Context(), id, WriteGroupSuccessorActivationRequest{Stage: stage, Boundary: stage.Envelope.Request.Boundary})
	require.NoError(t, err)
	c := WriteGroupSuccessorActivationConfirmation{Envelope: review.Envelope, ReviewDigest: review.ReviewDigest, ReadinessDigest: strings.Repeat("b", 64)}
	_, intent, err := successorActivationIntent(id, c)
	require.NoError(t, err)
	event := &WriteGroupCutover{ID: "status-cutover", WorkspaceID: stage.Envelope.Request.WorkspaceID, Kind: WriteGroupCutoverForward,
		SourceGroupID: id, Boundary: review.Envelope.Request.Boundary, ReviewDigest: review.ReviewDigest,
		Ended:     []CutoverEnd{{GroupID: id, Revision: stage.Envelope.SourceAppliedRevision}},
		Started:   []CutoverStart{{GroupID: stage.Current.GroupID, Revision: stage.Current.GroupRevision}},
		CreatedAt: h.service.clockNow(), OwnershipFenceRevision: "status-cutover", intent: intent}
	for _, member := range stage.Envelope.Request.Group.Members {
		event.Moved = append(event.Moved, memberRefOf(member))
	}
	sortMemberRefs(event.Moved)
	tx, err := h.db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	group := cloneWriteGroup(confirmed.Group)
	group.Status, group.AppliedRevision = WriteGroupStatusReady, group.Revision
	require.NoError(t, h.service.repo.insertVersionInTx(t.Context(), tx, group, event.Boundary))
	require.NoError(t, h.service.repo.insertCutoverInTx(t.Context(), tx, event))
	require.NoError(t, h.service.repo.insertRevisionEndInTx(t.Context(), tx, event.WorkspaceID, id, event.Ended[0].Revision, event.ID, event.Boundary))
	require.NoError(t, h.service.repo.insertOwnershipFenceInTx(t.Context(), tx, event))
	require.NoError(t, tx.Commit())
	return h, id, stage, event
}

func TestSuccessorOperationStatusAllowsReviewedLaterPredecessorCAS(t *testing.T) {
	h, id, stage, event := statusCutoverFixture(t, true)
	status, err := h.service.SuccessorOperationStatus(t.Context(), id, statusRequest(stage))
	require.NoError(t, err)
	require.Equal(t, event, status.Cutover.Result)
	require.Equal(t, "succeeded", status.Cutover.Status)
}

func TestSuccessorOperationStatusReadsOriginalCutoverAndRejectsCorruption(t *testing.T) {
	for _, mutation := range []string{"none", "moved", "boundary", "intent", "activation-kind", "schema", "fence", "ended"} {
		t.Run(mutation, func(t *testing.T) {
			h, id, stage, event := statusCutoverFixture(t)
			var err error
			switch mutation {
			case "moved":
				detail, e := json.Marshal(cutoverDetail{Ended: event.Ended, Started: event.Started, OwnershipFenceRevision: event.ID})
				require.NoError(t, e)
				_, err = h.db.Exec(`UPDATE write_group_cutovers SET detail=? WHERE id=?`, string(detail), event.ID)
			case "boundary":
				_, err = h.db.Exec(`UPDATE write_group_cutovers SET boundary=? WHERE id=?`, event.Boundary.Add(time.Second).Format(time.RFC3339Nano), event.ID)
			case "intent":
				_, err = h.db.Exec(`UPDATE write_group_cutovers SET intent=json_set(intent,'$.confirmation.readiness_digest','changed') WHERE id=?`, event.ID)
			case "activation-kind":
				var stored struct {
					Confirmation WriteGroupSuccessorActivationConfirmation `json:"confirmation"`
				}
				require.NoError(t, json.Unmarshal([]byte(event.intent), &stored))
				stored.Confirmation.Envelope.ActivationKind = successorActivationFirstStart
				_, digest, e := normalizeSuccessorActivationEnvelope(stored.Confirmation.Envelope)
				require.NoError(t, e)
				stored.Confirmation.ReviewDigest = digest
				_, intent, e := successorActivationIntent(id, stored.Confirmation)
				require.NoError(t, e)
				_, err = h.db.Exec(`UPDATE write_group_cutovers SET intent=?,review_digest=? WHERE id=?`, intent, digest, event.ID)
			case "schema":
				_, err = h.db.Exec(`UPDATE managed_schema_operations SET status='unknown' WHERE action='schema_apply'`)
			case "fence":
				_, err = h.db.Exec(`DELETE FROM write_group_ownership_fences WHERE fence_revision=?`, event.ID)
			case "ended":
				_, err = h.db.Exec(`DELETE FROM write_group_revision_ends WHERE cutover_id=?`, event.ID)
			}
			require.NoError(t, err)
			before := lineageReadSnapshot(t, h)
			status, err := h.service.SuccessorOperationStatus(t.Context(), id, statusRequest(stage))
			if mutation == "none" {
				require.NoError(t, err)
				require.Equal(t, "succeeded", status.Cutover.Status)
				require.Equal(t, event, status.Cutover.Result)
			} else {
				require.Error(t, err)
				require.Nil(t, status)
			}
			require.Equal(t, before, lineageReadSnapshot(t, h))
		})
	}
}
