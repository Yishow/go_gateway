package api

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestWriteGroupSuccessorReserveReplayAndNoActivation(t *testing.T) {
	f := newSuccessorAPIFixture(t)
	predecessor, err := f.f.groups.Get(t.Context(), f.groupID)
	require.NoError(t, err)
	before := successorAPIAllState(t, f.f.db)
	review := f.review(t)
	reserved := f.reserve(t, review)
	require.NotEqual(t, f.groupID, reserved.Group.ID)
	require.EqualValues(t, 2, reserved.Generation)
	require.Equal(t, f.request.OperationID, reserved.OperationID)
	require.Equal(t, workspace.WriteGroupStatusDraft, reserved.Group.Status)
	require.Empty(t, reserved.Group.AppliedRevision)
	require.Empty(t, reserved.Group.Destination.SchemaDigest)
	require.NotEqual(t, f.table, reserved.Group.Destination.TableName)
	current, err := f.f.groups.Get(t.Context(), f.groupID)
	require.NoError(t, err)
	require.Equal(t, predecessor.Group, current.Group)
	var head, pending string
	require.NoError(t, f.f.db.QueryRow(`SELECT current_group_id,reserved_group_id FROM write_group_lineages WHERE lineage_id=?`, reserved.LineageID).Scan(&head, &pending))
	require.Equal(t, f.groupID, head)
	require.Equal(t, reserved.Group.ID, pending)
	after := successorAPIAllState(t, f.f.db)
	for _, table := range []string{"write_group_versions", "wg_runtime_versions", "write_group_ownership_fences", "wg_delivery_samples", "wg_delivery_outbox", "wg_delivery_receipts", "write_group_profile_keys"} {
		require.Equal(t, before[table], after[table], table)
	}
	again := f.reserve(t, review)
	require.Equal(t, reserved.Group.ID, again.Group.ID)
	require.Equal(t, reserved.IntentDigest, again.IntentDigest)
	require.Equal(t, after, successorAPIAllState(t, f.f.db), "exact retry must not write")
	for _, action := range []string{"apply", "schema-preview"} {
		response := performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+reserved.Group.ID+"/"+action, map[string]any{
			"workspace_id": f.f.record.ID, "expected_workspace_revision": reserved.WorkspaceRevision,
			"expected_group_revision": reserved.Group.Revision, "expected_connector_revision": "connector-1",
		})
		if action == "apply" {
			require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
			require.Contains(t, response.Body.String(), "WRITE_GROUP_LIFECYCLE_BLOCKED")
		} else {
			// Existing ordinary schema renderer maps the reserved lifecycle refusal
			// to its safe generic unavailable response; this task does not change it.
			require.Equal(t, http.StatusServiceUnavailable, response.Code, response.Body.String())
			require.Contains(t, response.Body.String(), "RECORDING_PLAN_UNAVAILABLE")
		}
	}
	require.Equal(t, after, successorAPIAllState(t, f.f.db))
}

func TestWriteGroupSuccessorReserveRejectsChangedFrozenIntent(t *testing.T) {
	f := newSuccessorAPIFixture(t)
	review := f.review(t)
	before := successorAPIAllState(t, f.f.db)
	for _, mutation := range []struct {
		name   string
		change func(*workspace.WriteGroupSuccessorEnvelope)
	}{
		{"boundary", func(e *workspace.WriteGroupSuccessorEnvelope) {
			e.Request.Boundary = e.Request.Boundary.Add(time.Second)
		}},
		{"workspace", func(e *workspace.WriteGroupSuccessorEnvelope) { e.Request.WorkspaceID = "foreign" }},
		{"operation", func(e *workspace.WriteGroupSuccessorEnvelope) { e.Request.OperationID = "different" }},
		{"target", func(e *workspace.WriteGroupSuccessorEnvelope) { e.Request.Group.Name = "Changed intent" }},
		{"source_revision", func(e *workspace.WriteGroupSuccessorEnvelope) { e.SourceAppliedRevision = "foreign" }},
		{"workspace_cas", func(e *workspace.WriteGroupSuccessorEnvelope) { e.Request.ExpectedWorkspaceRevision = "stale" }},
		{"lineage_cas", func(e *workspace.WriteGroupSuccessorEnvelope) { e.Request.ExpectedLineageRevision++ }},
		{"predecessor_cas", func(e *workspace.WriteGroupSuccessorEnvelope) { e.Request.ExpectedPredecessorRevision = "stale" }},
		{"source_connector", func(e *workspace.WriteGroupSuccessorEnvelope) {
			e.Request.ExpectedPredecessorConnectorRevision = "stale"
		}},
		{"target_connector", func(e *workspace.WriteGroupSuccessorEnvelope) { e.Request.ExpectedTargetConnectorRevision = "stale" }},
		{"lineage", func(e *workspace.WriteGroupSuccessorEnvelope) { e.Request.LineageID = "foreign" }},
		{"generation", func(e *workspace.WriteGroupSuccessorEnvelope) { e.Generation++ }},
	} {
		t.Run(mutation.name, func(t *testing.T) {
			raw, err := json.Marshal(review.Envelope)
			require.NoError(t, err)
			var changed workspace.WriteGroupSuccessorEnvelope
			require.NoError(t, json.Unmarshal(raw, &changed))
			mutation.change(&changed)
			response := performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+f.groupID+"/successor", map[string]any{"action": "reserve", "envelope": changed, "review_digest": review.ReviewDigest})
			require.GreaterOrEqual(t, response.Code, 400, response.Body.String())
			require.Less(t, response.Code, 500, response.Body.String())
			require.Equal(t, before, successorAPIAllState(t, f.f.db))
		})
	}
}

func TestWriteGroupSuccessorStrictStagesCannotActivate(t *testing.T) {
	f := newSuccessorAPIFixture(t)
	review := f.review(t)
	before := successorAPIAllState(t, f.f.db)
	for _, payload := range []map[string]any{
		{"action": "activate", "envelope": review.Envelope, "review_digest": review.ReviewDigest},
		{"action": "reserve", "envelope": review.Envelope, "review_digest": review.ReviewDigest, "allow_reserved": true},
		{"action": "reserve", "envelope": review.Envelope, "review_digest": review.ReviewDigest, "token": "foreign"},
		{"action": "schema-preview", "envelope": review.Envelope, "review_digest": review.ReviewDigest, "schema_operation_id": "foreign"},
	} {
		response := performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+f.groupID+"/successor", payload)
		require.Equal(t, http.StatusBadRequest, response.Code, response.Body.String())
		require.Equal(t, before, successorAPIAllState(t, f.f.db))
	}
}

func TestWriteGroupSuccessorReserveRechecksLiveCAS(t *testing.T) {
	f := newSuccessorAPIFixture(t)
	review := f.review(t)
	_, err := f.f.db.Exec(`UPDATE database_connectors SET identity_revision='connector-new' WHERE id='connector-A'`)
	require.NoError(t, err)
	before := successorAPIAllState(t, f.f.db)
	response := performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+f.groupID+"/successor", map[string]any{"action": "reserve", "envelope": review.Envelope, "review_digest": review.ReviewDigest})
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
	require.Equal(t, before, successorAPIAllState(t, f.f.db))
}
