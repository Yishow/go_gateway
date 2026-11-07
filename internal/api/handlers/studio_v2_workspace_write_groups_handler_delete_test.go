package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestWriteGroupDeleteResponseContract(t *testing.T) {
	for _, tc := range []struct {
		name, code, action string
		status             int
	}{
		{"success", "", "", http.StatusOK},
		{"backlog", "WRITE_GROUP_DELETE_BACKLOG", "stop_and_wait", http.StatusConflict},
		{"active", "WRITE_GROUP_DELETE_ACTIVE", "stop_first", http.StatusConflict},
		{"pending", "WRITE_GROUP_DELETE_OPERATION_PENDING", "open_recording", http.StatusConflict},
		{"successor", "WRITE_GROUP_DELETE_HAS_SUCCESSOR", "open_recording", http.StatusConflict},
	} {
		t.Run(tc.name, func(t *testing.T) {
			db, ws, _ := recordingDeviceFixture(t)
			record, err := ws.GetOrCreate(t.Context())
			require.NoError(t, err)
			_, err = db.ExecContext(t.Context(), `INSERT INTO database_connectors(id,name,kind,connection_config,identity_revision,status,enabled) VALUES('connector','target','sqlite','{}','connector-revision-1','ready',1)`)
			require.NoError(t, err)
			seedDeviceRecording(t, db, record.ID, "R1", "Test", "disabled")
			seedDeleteResponseCondition(t, db, record.ID, tc.name)
			groups := workspace.NewWriteGroupService(ws, workspace.NewSQLWriteGroupRepository(db))
			handler := NewStudioV2WorkspaceWriteGroupsHandler(groups)
			payload, err := json.Marshal(map[string]string{"workspace_id": record.ID, "expected_workspace_revision": record.DatabaseSetupRevision, "expected_group_revision": "rev-1", "expected_connector_revision": "connector-revision-1"})
			require.NoError(t, err)
			response := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(response)
			c.Request = httptest.NewRequestWithContext(t.Context(), http.MethodDelete, "/write-groups/R1", bytes.NewReader(payload))
			c.Params = gin.Params{{Key: "id", Value: "R1"}}
			handler.Delete(c)
			require.Equal(t, tc.status, response.Code, response.Body.String())
			var body map[string]json.RawMessage
			require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
			if tc.name == "success" {
				var data struct {
					Deleted            bool   `json:"deleted"`
					GroupID            string `json:"group_id"`
					LineageID          string `json:"lineage_id"`
					DeletedGenerations int    `json:"deleted_generations"`
				}
				require.NoError(t, json.Unmarshal(body["data"], &data))
				require.True(t, data.Deleted)
				require.Equal(t, "R1", data.GroupID)
				require.Empty(t, data.LineageID)
				require.Equal(t, 1, data.DeletedGenerations)
				return
			}
			var failure struct {
				Code, Action string
				Details      json.RawMessage
			}
			require.NoError(t, json.Unmarshal(body["error"], &failure))
			require.Equal(t, tc.code, failure.Code)
			require.Equal(t, tc.action, failure.Action)
			if tc.name == "backlog" {
				require.JSONEq(t, `{"samples_unconsumed":1,"outbox":{"pending":2,"future_vendor_state":1}}`, string(failure.Details))
			} else if tc.name == "successor" {
				require.JSONEq(t, `{"referenced_by":[{"id":"R2","name":"other"}],"total":1}`, string(failure.Details))
			} else {
				require.Empty(t, failure.Details)
			}
		})
	}
}

func seedDeleteResponseCondition(t *testing.T, db *sql.DB, workspaceID, condition string) {
	t.Helper()
	var err error
	switch condition {
	case "backlog":
		_, err = db.ExecContext(t.Context(), `INSERT INTO wg_delivery_samples(workspace_id,group_id,group_revision,sample_id,member_key,observed_at,bucket_start,payload,payload_digest,created_at) VALUES(?,'R1','rev-1','S1','member','2026-01-01','2026-01-01','{}','digest','2026-01-01')`, workspaceID)
		require.NoError(t, err)
		for id, state := range map[string]string{"Q1": "pending", "Q2": "pending", "Q3": "future_vendor_state"} {
			_, err = db.ExecContext(t.Context(), `INSERT INTO wg_delivery_outbox(effect_key,record_id,workspace_id,group_id,group_revision,bucket_start,partition_key,destination_scope,connector_id,connector_revision,table_name,payload,payload_digest,state,next_retry_at,created_at,updated_at) VALUES(?,?,?,'R1','rev-1','2026-01-01','partition','scope','connector','revision','target','{}','digest',?,'2026-01-01','2026-01-01','2026-01-01')`, id, id, workspaceID, state)
			require.NoError(t, err)
		}
	case "active":
		_, err = db.ExecContext(t.Context(), `UPDATE write_groups SET status='ready',applied_revision='rev-1' WHERE id='R1'`)
		require.NoError(t, err)
		group, getErr := workspace.NewSQLWriteGroupRepository(db).Get(t.Context(), workspaceID, "R1")
		require.NoError(t, getErr)
		payload, encodeErr := json.Marshal(group)
		require.NoError(t, encodeErr)
		_, err = db.ExecContext(t.Context(), `INSERT INTO write_group_versions(group_id,group_revision,effective_at,payload) VALUES('R1','rev-1','2026-01-01T00:00:00Z',?)`, string(payload))
	case "pending":
		detail := `{"scope":{"workspace_id":"` + workspaceID + `","lineage_id":"L1","predecessor_group_id":"R1","successor_group_id":"R2"}}`
		_, err = db.ExecContext(t.Context(), `INSERT INTO managed_schema_operations(operation_id,token,workspace_id,scope_key,owner,action,status,created_at,updated_at,detail) VALUES('O2','T2',?,'scope','owner','successor_first_start','running',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,?)`, workspaceID, detail)
	case "successor":
		seedDeviceRecording(t, db, workspaceID, "R2", "other", "draft")
		_, err = db.ExecContext(t.Context(), `INSERT INTO wg_delivery_sample_handoffs(transfer_id,workspace_id,source_group_id,source_revision,source_seq,target_group_id,target_revision,target_seq,origin_seq,state,record) VALUES('H1',?,'R1','rev-1',1,'R2','rev-1',2,1,'finalized','{}')`, workspaceID)
	}
	require.NoError(t, err)
}
