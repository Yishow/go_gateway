package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"testing"

	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func successorAPIConfirmRequest(t *testing.T, f successorAPIFixture, review workspace.WriteGroupSuccessorReview, prepared workspace.WriteGroupSuccessorPreparation, token recordingplan.SchemaPreviewToken) map[string]any {
	t.Helper()
	request := successorAPIStage(t, f, review, prepared)
	request["action"], request["token"], request["digest"], request["schema_operation_id"] = "schema-confirm", token.Token, token.Digest, token.OperationID
	return request
}

func successorAPIRequireUnprepared(t *testing.T, f successorAPIFixture, prepared workspace.WriteGroupSuccessorPreparation, predecessor *workspace.WriteGroup) {
	t.Helper()
	saved, err := f.f.groups.Get(t.Context(), prepared.Group.ID)
	require.NoError(t, err)
	require.Equal(t, workspace.WriteGroupStatusDraft, saved.Group.Status)
	require.Empty(t, saved.Group.AppliedRevision)
	require.Empty(t, saved.Group.Destination.SchemaDigest)
	current, err := f.f.groups.Get(t.Context(), f.groupID)
	require.NoError(t, err)
	require.Equal(t, predecessor, current.Group)
	var head, pending string
	require.NoError(t, f.f.db.QueryRow(`SELECT current_group_id,reserved_group_id FROM write_group_lineages WHERE lineage_id=?`, prepared.LineageID).Scan(&head, &pending))
	require.Equal(t, f.groupID, head)
	require.Equal(t, prepared.Group.ID, pending)
}

func TestWriteGroupSuccessorRecoversProofFailureAfterRealDDL(t *testing.T) {
	f := newSuccessorAPIFixture(t)
	review := f.review(t)
	prepared := f.reserve(t, review)
	token := successorAPIPreview(t, f, review, prepared)
	old, err := f.f.groups.Get(t.Context(), f.groupID)
	require.NoError(t, err)
	request := successorAPIConfirmRequest(t, f, review, prepared, token)
	_, err = f.f.db.Exec(`CREATE TRIGGER fail_successor_proof BEFORE UPDATE ON system_settings WHEN NEW.key='studio_v2_workspace' BEGIN SELECT RAISE(ABORT,'proof-save-failed'); END`)
	require.NoError(t, err)
	response := performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+f.groupID+"/successor", request)
	require.GreaterOrEqual(t, response.Code, 400, response.Body.String())
	require.NotContains(t, response.Body.String(), "proof-save-failed")
	ledger := recordingplan.NewService(recordingplan.NewSQLRepository(f.f.db))
	committed, err := ledger.GetSchemaOperation(t.Context(), f.f.record.ID, token.OperationID)
	require.NoError(t, err)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, committed.Status)
	require.NotEmpty(t, committed.VerifiedDigest)
	successorAPIRequireUnprepared(t, f, prepared, old.Group)
	target, err := sql.Open("sqlite", f.f.file)
	require.NoError(t, err)
	defer target.Close()
	var count int
	require.NoError(t, target.QueryRow(`SELECT COUNT(*) FROM sqlite_schema WHERE type='table' AND name=?`, prepared.Group.Destination.TableName).Scan(&count))
	require.Equal(t, 1, count, "DDL happened despite local proof failure")
	targetBefore := successorAPIAllState(t, target)
	_, err = f.f.db.Exec(`DROP TRIGGER fail_successor_proof`)
	require.NoError(t, err)
	response = performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+f.groupID+"/successor", request)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	recovered, err := ledger.GetSchemaOperation(t.Context(), f.f.record.ID, token.OperationID)
	require.NoError(t, err)
	require.Equal(t, committed, recovered, "readback reuses original completed operation, never claims a replacement")
	require.Equal(t, targetBefore, successorAPIAllState(t, target))
	saved, err := f.f.groups.Get(t.Context(), prepared.Group.ID)
	require.NoError(t, err)
	require.Equal(t, committed.VerifiedDigest, saved.Group.Destination.SchemaDigest)
	require.Equal(t, workspace.WriteGroupStatusDraft, saved.Group.Status)
	require.Empty(t, saved.Group.AppliedRevision)
}

func TestWriteGroupSuccessorUnknownSchemaResultDoesNotReplayDDL(t *testing.T) {
	f := newSuccessorAPIFixture(t)
	review := f.review(t)
	prepared := f.reserve(t, review)
	token := successorAPIPreview(t, f, review, prepared)
	old, err := f.f.groups.Get(t.Context(), f.groupID)
	require.NoError(t, err)
	ledger := recordingplan.NewService(recordingplan.NewSQLRepository(f.f.db))
	operation, outcome, err := ledger.ClaimSchemaApply(t.Context(), &token)
	require.NoError(t, err)
	require.Equal(t, recordingplan.ClaimAcquired, outcome)
	// A prior executor performed an actual target statement but could not prove
	// the whole batch. Save that uncertainty through the real existing ledger.
	target, err := sql.Open("sqlite", f.f.file)
	require.NoError(t, err)
	defer target.Close()
	require.NotEmpty(t, token.Statements)
	_, err = target.Exec(token.Statements[0])
	require.NoError(t, err)
	unknown, err := ledger.FinishSchemaOperation(t.Context(), operation.OperationID, operation.Owner, recordingplan.SchemaOperationResult{Status: recordingplan.SchemaOperationUnknown, ExecutedStatements: 1, Reason: "outcome_unknown", NextAction: "inspect_original_operation"})
	require.NoError(t, err)
	before := successorAPIAllState(t, f.f.db)
	targetBefore := successorAPIAllState(t, target)
	request := successorAPIConfirmRequest(t, f, review, prepared, token)
	for i := 0; i < 2; i++ {
		response := performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+f.groupID+"/successor", request)
		require.Equal(t, http.StatusOK, response.Code, response.Body.String())
		var body struct {
			Data recordingplan.SchemaOperation `json:"data"`
		}
		require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
		require.Equal(t, unknown.OperationID, body.Data.OperationID)
		require.Equal(t, recordingplan.SchemaOperationUnknown, body.Data.Status)
		require.Empty(t, body.Data.VerifiedDigest)
		require.Equal(t, before, successorAPIAllState(t, f.f.db))
		require.Equal(t, targetBefore, successorAPIAllState(t, target))
		successorAPIRequireUnprepared(t, f, prepared, old.Group)
	}
}

func TestWriteGroupSuccessorPostDDLConnectorCASRetainsTrueOutcome(t *testing.T) {
	f := newSuccessorAPIFixture(t)
	review := f.review(t)
	prepared := f.reserve(t, review)
	token := successorAPIPreview(t, f, review, prepared)
	old, err := f.f.groups.Get(t.Context(), f.groupID)
	require.NoError(t, err)
	request := successorAPIConfirmRequest(t, f, review, prepared, token)
	// The target transaction has committed before the schema ledger records
	// succeeded. Inject a concurrent connector edit at exactly that boundary.
	_, err = f.f.db.Exec(`CREATE TRIGGER change_connector_after_successor_ddl AFTER UPDATE ON managed_schema_operations WHEN NEW.operation_id='` + token.OperationID + `' AND NEW.status='succeeded' BEGIN UPDATE database_connectors SET identity_revision='connector-after-ddl' WHERE id='connector-A'; END`)
	require.NoError(t, err)
	response := performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+f.groupID+"/successor", request)
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
	ledger := recordingplan.NewService(recordingplan.NewSQLRepository(f.f.db))
	committed, err := ledger.GetSchemaOperation(t.Context(), f.f.record.ID, token.OperationID)
	require.NoError(t, err)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, committed.Status)
	successorAPIRequireUnprepared(t, f, prepared, old.Group)
	target, err := sql.Open("sqlite", f.f.file)
	require.NoError(t, err)
	defer target.Close()
	var count int
	require.NoError(t, target.QueryRow(`SELECT COUNT(*) FROM sqlite_schema WHERE type='table' AND name=?`, prepared.Group.Destination.TableName).Scan(&count))
	require.Equal(t, 1, count)
	before := successorAPIAllState(t, f.f.db)
	targetBefore := successorAPIAllState(t, target)
	response = performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+f.groupID+"/successor", request)
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
	require.Equal(t, before, successorAPIAllState(t, f.f.db))
	require.Equal(t, targetBefore, successorAPIAllState(t, target))
}
