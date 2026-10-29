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

func successorAPIPreview(t *testing.T, f successorAPIFixture, review workspace.WriteGroupSuccessorReview, prepared workspace.WriteGroupSuccessorPreparation) recordingplan.SchemaPreviewToken {
	t.Helper()
	request := successorAPIStage(t, f, review, prepared)
	request["action"] = "schema-preview"
	response := performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+f.groupID+"/successor", request)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	var body struct {
		Data recordingplan.SchemaPreviewToken `json:"data"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
	return body.Data
}

func TestWriteGroupSuccessorSchemaPreparationDoesNotActivate(t *testing.T) {
	f := newSuccessorAPIFixture(t)
	predecessor, err := f.f.groups.Get(t.Context(), f.groupID)
	require.NoError(t, err)
	review := f.review(t)
	prepared := f.reserve(t, review)
	target, err := sql.Open("sqlite", f.f.file)
	require.NoError(t, err)
	defer target.Close()
	beforeTarget := successorAPIAllState(t, target)
	token := successorAPIPreview(t, f, review, prepared)
	require.Equal(t, prepared.Group.ID, token.PlanID)
	require.NotEqual(t, prepared.OperationID, token.OperationID, "schema operation and transition reservation identities are distinct")
	require.NotEmpty(t, token.SourceDigest)
	require.NotEmpty(t, token.Statements)
	require.Equal(t, beforeTarget, successorAPIAllState(t, target), "schema preview must not execute DDL")
	request := successorAPIStage(t, f, review, prepared)
	request["action"], request["token"], request["digest"], request["schema_operation_id"] = "schema-confirm", token.Token, token.Digest, token.OperationID
	response := performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+f.groupID+"/successor", request)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	var body struct {
		Data recordingplan.SchemaOperation `json:"data"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
	require.Equal(t, recordingplan.SchemaOperationSucceeded, body.Data.Status)
	require.Equal(t, token.OperationID, body.Data.OperationID)
	require.NotEmpty(t, body.Data.VerifiedDigest)
	var count int
	require.NoError(t, target.QueryRow(`SELECT COUNT(*) FROM sqlite_schema WHERE type='table' AND name=?`, prepared.Group.Destination.TableName).Scan(&count))
	require.Equal(t, 1, count)
	require.Equal(t, beforeTarget[f.table], successorAPIAllState(t, target)[f.table], "old SQL history is not rewritten")
	saved, err := f.f.groups.Get(t.Context(), prepared.Group.ID)
	require.NoError(t, err)
	require.Equal(t, workspace.WriteGroupStatusDraft, saved.Group.Status)
	require.Empty(t, saved.Group.AppliedRevision)
	require.NotEmpty(t, saved.Group.Destination.SchemaDigest)
	current, err := f.f.groups.Get(t.Context(), f.groupID)
	require.NoError(t, err)
	require.Equal(t, predecessor.Group, current.Group)
	var head string
	require.NoError(t, f.f.db.QueryRow(`SELECT current_group_id FROM write_group_lineages WHERE lineage_id=?`, prepared.LineageID).Scan(&head))
	require.Equal(t, f.groupID, head)
	var newVersions, newRuntime, newFences int
	require.NoError(t, f.f.db.QueryRow(`SELECT count(*) FROM write_group_versions WHERE group_id=?`, prepared.Group.ID).Scan(&newVersions))
	require.NoError(t, f.f.db.QueryRow(`SELECT count(*) FROM wg_runtime_versions WHERE group_id=?`, prepared.Group.ID).Scan(&newRuntime))
	require.NoError(t, f.f.db.QueryRow(`SELECT count(*) FROM write_group_ownership_fences WHERE group_id=?`, prepared.Group.ID).Scan(&newFences))
	require.Zero(t, newVersions)
	require.Zero(t, newRuntime)
	require.Zero(t, newFences)
	after := successorAPIAllState(t, f.f.db)
	replay := f.reserve(t, review)
	require.Equal(t, prepared.Group.ID, replay.Group.ID)
	require.Equal(t, after, successorAPIAllState(t, f.f.db), "original frozen reserve retries survive schema CAS updates")
	response = performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+f.groupID+"/successor", request)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	var replayed struct {
		Data recordingplan.SchemaOperation `json:"data"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &replayed))
	require.Equal(t, token.OperationID, replayed.Data.OperationID)
	require.Equal(t, after, successorAPIAllState(t, f.f.db), "same schema token replay must not run a new operation or rewrite proof")
}

func TestWriteGroupSuccessorSchemaStagesRefuseStaleOrForeignScope(t *testing.T) {
	f := newSuccessorAPIFixture(t)
	review := f.review(t)
	prepared := f.reserve(t, review)
	before := successorAPIAllState(t, f.f.db)
	request := successorAPIStage(t, f, review, prepared)
	stale := request["current"].(workspace.WriteGroupSuccessorCurrentCAS)
	stale.WorkspaceRevision = "stale"
	request["action"], request["current"] = "schema-preview", stale
	response := performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+f.groupID+"/successor", request)
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
	require.Equal(t, before, successorAPIAllState(t, f.f.db), "stale stage must not create preview tokens")
	token := successorAPIPreview(t, f, review, prepared)
	before = successorAPIAllState(t, f.f.db)
	target, err := sql.Open("sqlite", f.f.file)
	require.NoError(t, err)
	defer target.Close()
	beforeTarget := successorAPIAllState(t, target)
	for _, field := range []string{"token", "digest", "schema_operation_id"} {
		request = successorAPIStage(t, f, review, prepared)
		request["action"], request["token"], request["digest"], request["schema_operation_id"] = "schema-confirm", token.Token, token.Digest, token.OperationID
		request[field] = "foreign"
		response = performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+f.groupID+"/successor", request)
		require.GreaterOrEqual(t, response.Code, 400, response.Body.String())
		require.Less(t, response.Code, 500, response.Body.String())
		require.Equal(t, before, successorAPIAllState(t, f.f.db))
		require.Equal(t, beforeTarget, successorAPIAllState(t, target))
	}
}

func TestWriteGroupSuccessorSchemaConfirmRechecksConnectorBeforeDDL(t *testing.T) {
	f := newSuccessorAPIFixture(t)
	review := f.review(t)
	prepared := f.reserve(t, review)
	token := successorAPIPreview(t, f, review, prepared)
	request := successorAPIStage(t, f, review, prepared)
	request["action"], request["token"], request["digest"], request["schema_operation_id"] = "schema-confirm", token.Token, token.Digest, token.OperationID
	_, err := f.f.db.Exec(`UPDATE database_connectors SET identity_revision='connector-new' WHERE id='connector-A'`)
	require.NoError(t, err)
	before := successorAPIAllState(t, f.f.db)
	target, err := sql.Open("sqlite", f.f.file)
	require.NoError(t, err)
	defer target.Close()
	targetBefore := successorAPIAllState(t, target)
	response := performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+f.groupID+"/successor", request)
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
	require.Equal(t, before, successorAPIAllState(t, f.f.db))
	require.Equal(t, targetBefore, successorAPIAllState(t, target))
}

func TestWriteGroupSuccessorPreparationFailureKeepsPredecessor(t *testing.T) {
	f := newSuccessorAPIFixture(t)
	review := f.review(t)
	prepared := f.reserve(t, review)
	token := successorAPIPreview(t, f, review, prepared)
	old, err := f.f.groups.Get(t.Context(), f.groupID)
	require.NoError(t, err)
	target, err := sql.Open("sqlite", f.f.file)
	require.NoError(t, err)
	defer target.Close()
	// A real target changes after preview: the proposed table name is occupied
	// by an incompatible object. Confirmation must not treat preview as proof.
	_, err = target.Exec(`CREATE VIEW "` + prepared.Group.Destination.TableName + `" AS SELECT 1 AS incompatible`)
	require.NoError(t, err)
	request := successorAPIStage(t, f, review, prepared)
	request["action"], request["token"], request["digest"], request["schema_operation_id"] = "schema-confirm", token.Token, token.Digest, token.OperationID
	response := performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+f.groupID+"/successor", request)
	// The existing ledger renderer successfully returns a failed operation;
	// HTTP success is not schema success.
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	var outcome struct {
		Data recordingplan.SchemaOperation `json:"data"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &outcome))
	require.Equal(t, recordingplan.SchemaOperationFailed, outcome.Data.Status)
	require.Equal(t, token.OperationID, outcome.Data.OperationID)
	require.Zero(t, outcome.Data.ExecutedStatements)
	saved, err := f.f.groups.Get(t.Context(), prepared.Group.ID)
	require.NoError(t, err)
	require.Equal(t, workspace.WriteGroupStatusDraft, saved.Group.Status)
	require.Empty(t, saved.Group.AppliedRevision)
	require.Empty(t, saved.Group.Destination.SchemaDigest)
	current, err := f.f.groups.Get(t.Context(), f.groupID)
	require.NoError(t, err)
	require.Equal(t, old.Group, current.Group)
	var head string
	require.NoError(t, f.f.db.QueryRow(`SELECT current_group_id FROM write_group_lineages WHERE lineage_id=?`, prepared.LineageID).Scan(&head))
	require.Equal(t, f.groupID, head)
	var count int
	require.NoError(t, target.QueryRow(`SELECT COUNT(*) FROM sqlite_schema WHERE type='view' AND name=?`, prepared.Group.Destination.TableName).Scan(&count))
	require.Equal(t, 1, count)
}

func TestWriteGroupSuccessorPreparesLongWithoutMovingWideOwner(t *testing.T) {
	f := newSuccessorAPIFixture(t)
	f.request.Group.LayoutVersion = workspace.WriteGroupLayoutManagedLongV1
	f.request.Group.StorageMode = workspace.WriteGroupStorageModePlain
	f.request.Group.Destination.TableName = ""
	review := f.review(t)
	prepared := f.reserve(t, review)
	require.Equal(t, workspace.WriteGroupLayoutManagedLongV1, prepared.Group.LayoutVersion)
	token := successorAPIPreview(t, f, review, prepared)
	response := performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+f.groupID+"/successor", successorAPIConfirmRequest(t, f, review, prepared, token))
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	var body struct {
		Data recordingplan.SchemaOperation `json:"data"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
	require.Equal(t, recordingplan.SchemaOperationSucceeded, body.Data.Status)
	saved, err := f.f.groups.Get(t.Context(), prepared.Group.ID)
	require.NoError(t, err)
	require.Equal(t, workspace.WriteGroupStatusDraft, saved.Group.Status)
	require.Empty(t, saved.Group.AppliedRevision)
	require.Equal(t, body.Data.VerifiedDigest, saved.Group.Destination.SchemaDigest)
	var current, reserved string
	require.NoError(t, f.f.db.QueryRow(`SELECT current_group_id,reserved_group_id FROM write_group_lineages WHERE lineage_id=?`, prepared.LineageID).Scan(&current, &reserved))
	require.Equal(t, f.groupID, current)
	require.Equal(t, prepared.Group.ID, reserved)
	var confirmedOperation, confirmedDigest string
	require.NoError(t, f.f.db.QueryRow(`SELECT confirmed_operation_id,confirmed_digest FROM write_group_storage_scopes WHERE table_name=?`, prepared.Group.Destination.TableName).Scan(&confirmedOperation, &confirmedDigest))
	require.Equal(t, token.OperationID, confirmedOperation)
	require.Equal(t, body.Data.VerifiedDigest, confirmedDigest)
}
