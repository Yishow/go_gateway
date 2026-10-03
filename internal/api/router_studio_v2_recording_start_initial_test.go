package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestRecordingStartRecoversClaimBeforeFirstProgressCheckpoint(t *testing.T) {
	f := newRecordingStartFixture(t)
	request := f.request("crash-before-initial-checkpoint")
	encoded, err := json.Marshal(request)
	require.NoError(t, err)
	hash := sha256.Sum256(encoded)
	digest := hex.EncodeToString(hash[:])
	view := workspace.RecordingStartOperation{
		WorkspaceID: request.WorkspaceID, DeviceIDs: request.DeviceIDs, IntentDigest: digest, Stage: "save",
		Groups:  []workspace.RecordingStartGroupProgress{{GroupID: request.Groups[0].GroupID, GroupRevision: request.Groups[0].ExpectedGroupRevision}},
		Devices: []workspace.RecordingStartDeviceProgress{{DeviceID: request.DeviceIDs[0]}},
	}
	initial, err := json.Marshal(map[string]any{
		"version": 1, "request": request, "workspace_revision": request.ExpectedWorkspaceRevision,
		"scopes":         []map[string]string{{"source_digest": "", "schema_revision": "", "schema_digest": ""}},
		"device_digests": []string{""}, "view": view,
	})
	require.NoError(t, err)
	claimed, outcome, err := f.ledger.ClaimRecordingStart(t.Context(), request.WorkspaceID, request.RequestID, digest, string(initial))
	require.NoError(t, err)
	require.Equal(t, recordingplan.ClaimAcquired, outcome)
	_, err = f.db.ExecContext(t.Context(), `UPDATE managed_schema_operations SET updated_at=? WHERE operation_id=?`, time.Now().UTC().Add(-20*time.Minute), claimed.OperationID)
	require.NoError(t, err)
	resumed, err := f.start.Start(t.Context(), request)
	require.NoError(t, err)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, resumed.Status, "%+v", resumed)
	require.Equal(t, claimed.OperationID, resumed.OperationID)
	require.Equal(t, 1, f.activation.Calls())
}

func TestRecordingStartGetDoesNotCreateWorkspace(t *testing.T) {
	f := newRecordingStartFixture(t)
	_, err := f.db.ExecContext(t.Context(), `DELETE FROM system_settings WHERE key='studio_v2_workspace'`)
	require.NoError(t, err)
	_, err = f.start.Get(t.Context(), "missing-operation")
	require.ErrorIs(t, err, recordingplan.ErrSchemaOperationNotFound)
	var count int
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM system_settings WHERE key='studio_v2_workspace'`).Scan(&count))
	require.Zero(t, count, "querying start progress must not create a workspace")
}
