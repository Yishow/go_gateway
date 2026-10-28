package workspace

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/recordingplan"
)

func statusRequest(stage WriteGroupSuccessorStageRequest) WriteGroupSuccessorOperationStatusRequest {
	return WriteGroupSuccessorOperationStatusRequest{WorkspaceID: stage.Envelope.Request.WorkspaceID, LineageID: stage.Envelope.Request.LineageID, OperationID: stage.Envelope.Request.OperationID}
}

func TestSuccessorOperationStatusImmutableReserveAndForeignScope(t *testing.T) {
	h, id, stage, _ := successorSchemaFixture(t)
	req := statusRequest(stage)
	tx, err := h.db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	root, err := recordingplan.ReadSuccessorReserveInTx(t.Context(), tx, req.WorkspaceID, req.OperationID)
	require.NoError(t, err)
	_, reserved, err := decodeSuccessorReserveRecord(root)
	require.NoError(t, err, "validate original reserve payload")
	_, err = decodeSuccessorReserveSchemaBinding(root.SchemaSourceBinding, reserved)
	require.NoError(t, err, "validate original reserve schema binding")
	require.NoError(t, tx.Rollback())
	status, err := h.service.SuccessorOperationStatus(t.Context(), id, req)
	require.NoError(t, err)
	require.Equal(t, "succeeded", status.Reserve.Status)
	require.Equal(t, stage.Current.GroupID, status.Reserve.Result.Group.ID)
	require.Equal(t, "not_submitted", status.Cutover.Status)
	require.Equal(t, "not_submitted", status.FirstStart.Status)
	require.Empty(t, status.SchemaOperations)
	_, err = h.db.Exec(`UPDATE write_groups SET revision='later-revision',destination_schema_revision='later-proof' WHERE id=?`, stage.Current.GroupID)
	require.NoError(t, err)
	_, err = h.db.Exec(`UPDATE write_group_lineages SET current_group_id=?,reserved_group_id=NULL,reservation_operation_id=NULL,reservation_intent_digest=NULL WHERE lineage_id=?`, stage.Current.GroupID, req.LineageID)
	require.NoError(t, err)
	before := lineageReadSnapshot(t, h)
	again, err := h.service.SuccessorOperationStatus(t.Context(), id, req)
	require.NoError(t, err)
	require.Equal(t, status, again, "current head and draft revisions cannot replace the immutable reserve result")
	require.Equal(t, before, lineageReadSnapshot(t, h))
	for _, mode := range []string{"path", "workspace", "lineage", "operation"} {
		t.Run(mode, func(t *testing.T) {
			bad, path := req, id
			switch mode {
			case "path":
				path = "foreign"
			case "workspace":
				bad.WorkspaceID = "foreign"
			case "lineage":
				bad.LineageID = "foreign"
			case "operation":
				bad.OperationID = "foreign"
			}
			got, err := h.service.SuccessorOperationStatus(t.Context(), path, bad)
			require.ErrorIs(t, err, ErrWriteGroupNotFound)
			require.Nil(t, got)
		})
	}
}

func statusPreview(t *testing.T, h *circuitProfileHarness, id string, stage WriteGroupSuccessorStageRequest) (*recordingplan.SchemaPreviewToken, *recordingplan.SQLRepository, *recordingplan.Service) {
	t.Helper()
	scope, err := h.service.SuccessorSchemaScope(t.Context(), id, stage)
	require.NoError(t, err)
	repo := recordingplan.NewSQLRepository(h.db)
	ledger := recordingplan.NewService(repo)
	token, err := ledger.PrepareSchemaPreview(t.Context(), scope, func(context.Context, string) (recordingplan.TargetTableInspection, error) {
		return recordingplan.TargetTableInspection{Status: "missing"}, nil
	})
	require.NoError(t, err)
	return token, repo, ledger
}

func TestSuccessorOperationStatusRetainsUnknownBesideNewUnsubmittedPreview(t *testing.T) {
	h, id, stage, _ := successorSchemaFixture(t)
	first, repo, ledger := statusPreview(t, h, id, stage)
	op, _, err := ledger.ClaimSchemaApply(t.Context(), first)
	require.NoError(t, err)
	_, err = repo.FinishSchemaOperation(t.Context(), op.OperationID, op.Owner, recordingplan.SchemaOperationResult{Status: recordingplan.SchemaOperationUnknown, Reason: "unverifiable", NextAction: "postgres://private:password@private", Detail: `{"private":"secret"}`})
	require.NoError(t, err)
	second, _, _ := statusPreview(t, h, id, stage)
	before := lineageReadSnapshot(t, h)
	status, err := h.service.SuccessorOperationStatus(t.Context(), id, statusRequest(stage))
	require.NoError(t, err)
	require.Len(t, status.SchemaOperations, 2)
	states := map[string]string{}
	for _, item := range status.SchemaOperations {
		states[item.OperationID] = item.Status
	}
	require.Equal(t, "unknown", states[first.OperationID])
	require.Equal(t, "preview_not_submitted", states[second.OperationID])
	encoded, err := json.Marshal(status)
	require.NoError(t, err)
	for _, secret := range []string{first.Token, second.Token, "private", "password", "secret", "source_evidence", "schema_source_binding", "statements"} {
		require.NotContains(t, string(encoded), secret)
	}
	require.Equal(t, before, lineageReadSnapshot(t, h))
}
