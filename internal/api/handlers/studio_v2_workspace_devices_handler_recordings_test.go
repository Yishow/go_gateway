package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/audit"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestStudioV2DeviceDeleteRefusesRecordingReferences(t *testing.T) {
	db, ws, handler := recordingDeviceFixture(t)
	record, err := ws.GetOrCreate(t.Context())
	require.NoError(t, err)
	seedDeviceRecording(t, db, record.ID, "R1", "MC (Instantaneous)", "draft")
	seedDeviceRecording(t, db, record.ID, "R2", "100.76.76.1(mc)", "draft")
	seedSQLiteCascadeRule(t.Context(), t, sourcerule.NewSQLRepository(db), "rule-target", "device-target")
	response := deleteRecordingDevice(t, handler)
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
	var body struct {
		Error struct {
			Code, Action string
			Details      struct {
				Recordings []struct{ ID, Name string }
				Total      int
			}
		}
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
	require.Equal(t, "DEVICE_REFERENCED_BY_RECORDINGS", body.Error.Code)
	require.Equal(t, "remove_recordings_first", body.Error.Action)
	require.Equal(t, 2, body.Error.Details.Total)
	require.Len(t, body.Error.Details.Recordings, 2)
	var deviceCount, groupCount int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM devices WHERE id='device-target'`).Scan(&deviceCount))
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM write_groups WHERE id IN ('R1','R2')`).Scan(&groupCount))
	require.Equal(t, 1, deviceCount)
	require.Equal(t, 2, groupCount)
	var ruleCount int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM source_rules WHERE id='rule-target'`).Scan(&ruleCount))
	require.Equal(t, 1, ruleCount)
	current, err := ws.GetOrCreate(t.Context())
	require.NoError(t, err)
	require.Equal(t, record, current)
}

func TestStudioV2DeviceDeleteReportsAtMostTwentyReferences(t *testing.T) {
	db, ws, handler := recordingDeviceFixture(t)
	record, err := ws.GetOrCreate(t.Context())
	require.NoError(t, err)
	for i := range 21 {
		seedDeviceRecording(t, db, record.ID, fmt.Sprintf("R%02d", i), fmt.Sprintf("recording %02d", i), "draft")
	}
	response := deleteRecordingDevice(t, handler)
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
	var body struct {
		Error struct {
			Code    string
			Details struct {
				Recordings []struct{ ID, Name string }
				Total      int
			}
		}
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
	require.Equal(t, "DEVICE_REFERENCED_BY_RECORDINGS", body.Error.Code)
	require.Equal(t, 21, body.Error.Details.Total)
	require.Len(t, body.Error.Details.Recordings, 20)
	var count int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM write_groups`).Scan(&count))
	require.Equal(t, 21, count)
}

func TestStudioV2DeviceDeleteRefusesBasicKeyOnlyReference(t *testing.T) {
	db, ws, handler := recordingDeviceFixture(t)
	record, err := ws.GetOrCreate(t.Context())
	require.NoError(t, err)
	seedDeviceRecording(t, db, record.ID, "R1", "Basic", "draft")
	_, err = db.ExecContext(t.Context(), `DELETE FROM write_group_members WHERE group_id='R1'`)
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), `INSERT INTO write_group_basic_keys(workspace_id,device_id,canonical_role,group_id) VALUES(?,'device-target','basic','R1')`, record.ID)
	require.NoError(t, err)
	response := deleteRecordingDevice(t, handler)
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
	require.Contains(t, response.Body.String(), "DEVICE_REFERENCED_BY_RECORDINGS")
	var count int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM write_group_basic_keys WHERE group_id='R1'`).Scan(&count))
	require.Equal(t, 1, count)
}

type failingDeleteShareScope struct{ failedClosed bool }

func (s *failingDeleteShareScope) InvalidateWorkspaceMembership(context.Context, string) error {
	return errors.New("private backend detail")
}

func (s *failingDeleteShareScope) FailClosedMembership() { s.failedClosed = true }

func TestStudioV2DeviceDeleteShareRefreshFailureIsExplicitAndClosed(t *testing.T) {
	db, ws, handler := recordingDeviceFixture(t)
	share := &failingDeleteShareScope{}
	ws.SetMembershipScopeInvalidator(share)
	response := deleteRecordingDevice(t, handler)
	require.Equal(t, http.StatusInternalServerError, response.Code, response.Body.String())
	require.Contains(t, response.Body.String(), "DEVICE_DELETED_SHARE_SCOPE_REFRESH_FAILED")
	require.NotContains(t, response.Body.String(), "private backend detail")
	require.True(t, share.failedClosed)
	var count int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM devices WHERE id='device-target'`).Scan(&count))
	require.Zero(t, count, "the failure happens after the device transaction committed")
}

func TestStudioV2DeviceDeleteWritesAudit(t *testing.T) {
	db, ws, handler := recordingDeviceFixture(t)
	record, err := ws.GetOrCreate(t.Context())
	require.NoError(t, err)
	auditSvc := audit.NewService(audit.NewSQLRepository(db))
	handler.WithAudit(auditSvc)
	response := deleteRecordingDevice(t, handler)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	entries, err := auditSvc.List(t.Context(), audit.ListFilter{WorkspaceID: record.ID, Limit: 10})
	require.NoError(t, err)
	require.Len(t, entries, 1)
	require.Equal(t, audit.EventTypeDeviceDeleted, entries[0].EventType)
	require.Equal(t, "device-target", entries[0].ReferenceID)
	require.Contains(t, entries[0].Details, "device-target")
	require.NotContains(t, entries[0].Details, "127.0.0.1")
}

func TestStudioV2DeviceDeletePurgesTombstone(t *testing.T) {
	db, ws, handler := recordingDeviceFixture(t)
	record, err := ws.GetOrCreate(t.Context())
	require.NoError(t, err)
	seedDeviceRecording(t, db, record.ID, "R1", "Test", "deleted")
	response := deleteRecordingDevice(t, handler)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	var count int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM write_groups WHERE id='R1'`).Scan(&count))
	require.Zero(t, count)
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM devices WHERE id='device-target'`).Scan(&count))
	require.Zero(t, count)
}

func TestStudioV2DeviceDeleteRollsBackTombstoneBacklog(t *testing.T) {
	db, ws, handler := recordingDeviceFixture(t)
	record, err := ws.GetOrCreate(t.Context())
	require.NoError(t, err)
	seedDeviceRecording(t, db, record.ID, "R1", "Test", "deleted")
	_, err = db.ExecContext(t.Context(), `INSERT INTO wg_delivery_samples(workspace_id,group_id,group_revision,sample_id,member_key,observed_at,bucket_start,payload,payload_digest,created_at) VALUES(?,'R1','rev-1','S1','member','2026-01-01','2026-01-01','{}','digest','2026-01-01')`, record.ID)
	require.NoError(t, err)
	response := deleteRecordingDevice(t, handler)
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
	require.Contains(t, response.Body.String(), "WRITE_GROUP_DELETE_BACKLOG")
	for _, query := range []string{`SELECT COUNT(*) FROM devices WHERE id='device-target'`, `SELECT COUNT(*) FROM write_groups WHERE id='R1'`, `SELECT COUNT(*) FROM wg_delivery_samples WHERE group_id='R1' AND consumed=0`} {
		var count int
		require.NoError(t, db.QueryRowContext(t.Context(), query).Scan(&count))
		require.Equal(t, 1, count, query)
	}
	current, err := ws.GetOrCreate(t.Context())
	require.NoError(t, err)
	require.Equal(t, record, current)
}

func TestStudioV2DeviceDeleteRollsBackWhenDeviceDeleteFails(t *testing.T) {
	db, ws, handler := recordingDeviceFixture(t)
	record, err := ws.GetOrCreate(t.Context())
	require.NoError(t, err)
	seedDeviceRecording(t, db, record.ID, "R1", "Test", "deleted")
	_, err = db.ExecContext(t.Context(), `CREATE TRIGGER refuse_device_delete BEFORE DELETE ON devices BEGIN SELECT RAISE(ABORT,'private failure'); END`)
	require.NoError(t, err)
	response := deleteRecordingDevice(t, handler)
	require.Equal(t, http.StatusInternalServerError, response.Code, response.Body.String())
	require.NotContains(t, response.Body.String(), "private failure")
	for _, query := range []string{`SELECT COUNT(*) FROM devices WHERE id='device-target'`, `SELECT COUNT(*) FROM write_groups WHERE id='R1'`, `SELECT COUNT(*) FROM write_group_members WHERE group_id='R1'`} {
		var count int
		require.NoError(t, db.QueryRowContext(t.Context(), query).Scan(&count))
		require.Equal(t, 1, count, query)
	}
	current, err := ws.GetOrCreate(t.Context())
	require.NoError(t, err)
	require.Equal(t, record, current)
}

func TestStudioV2DeviceDeleteRefusesSharedTombstone(t *testing.T) {
	db, ws, handler := recordingDeviceFixture(t)
	record, err := ws.GetOrCreate(t.Context())
	require.NoError(t, err)
	seedDeviceRecording(t, db, record.ID, "R1", "Test", "deleted")
	_, err = db.ExecContext(t.Context(), `INSERT INTO write_group_members(group_id,member_index,device_id,point_id,tag_id,target_column) VALUES('R1',1,'other-device','point-2','tag-2','value2')`)
	require.NoError(t, err)
	response := deleteRecordingDevice(t, handler)
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
	require.Contains(t, response.Body.String(), "DEVICE_REFERENCED_BY_RECORDINGS")
	for _, query := range []string{`SELECT COUNT(*) FROM devices WHERE id='device-target'`, `SELECT COUNT(*) FROM write_groups WHERE id='R1'`} {
		var count int
		require.NoError(t, db.QueryRowContext(t.Context(), query).Scan(&count))
		require.Equal(t, 1, count, query)
	}
}

func recordingDeviceFixture(t *testing.T) (*sql.DB, *workspace.Service, *StudioV2WorkspaceDevicesHandler) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	db := openStudioV2CascadeTestDB(t)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	deviceRepo := device.NewSQLRepository(db)
	deviceSvc := device.NewService(deviceRepo, nil)
	ws := workspace.NewService(workspace.NewSQLRepository(db))
	seedSQLiteCascadeDevice(t.Context(), t, deviceRepo, "device-target")
	_, err := ws.AttachDevice(t.Context(), "device-target")
	require.NoError(t, err)
	return db, ws, NewStudioV2WorkspaceDevicesHandler(ws, deviceSvc).WithWriteGroups(workspace.NewWriteGroupService(ws, workspace.NewSQLWriteGroupRepository(db)))
}

func seedDeviceRecording(t *testing.T, db *sql.DB, workspaceID, id, name, status string) {
	t.Helper()
	_, err := db.ExecContext(t.Context(), `INSERT INTO write_groups(id,workspace_id,revision,name,status,destination_connector_id,destination_connector_revision,destination_table_name) VALUES(?,?,'rev-1',?,?, 'connector','connector-revision-1','target')`, id, workspaceID, name, status)
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), `INSERT INTO write_group_members(group_id,member_index,device_id,point_id,tag_id,target_column) VALUES(?,0,'device-target','point-1','tag-1','value')`, id)
	require.NoError(t, err)
}

func deleteRecordingDevice(t *testing.T, handler *StudioV2WorkspaceDevicesHandler) *httptest.ResponseRecorder {
	t.Helper()
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = httptest.NewRequestWithContext(t.Context(), http.MethodDelete, "/api/v1/datalink/studio-v2/workspace/devices/device-target", http.NoBody)
	c.Params = gin.Params{{Key: "id", Value: "device-target"}}
	handler.Delete(c)
	return resp
}
