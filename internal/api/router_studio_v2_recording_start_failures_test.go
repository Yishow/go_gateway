package api

import (
	"context"
	"database/sql"
	"net/http"
	"strings"
	"testing"

	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

type recordingStartBlockedDeviceStub struct{}

func (recordingStartBlockedDeviceStub) GetByID(_ context.Context, id string) (*schema.Device, error) {
	return &schema.Device{ID: id}, nil
}

func (recordingStartBlockedDeviceStub) CheckReadiness(_ context.Context, id string) (*schema.DeviceReadiness, error) {
	return &schema.DeviceReadiness{
		DeviceID: id, BlockingReasons: []string{"test readiness block"},
	}, nil
}

func TestRecordingStartSaveFailureDoesNotAdvanceGroupOrApply(t *testing.T) {
	f := newRecordingStartFixture(t, recordingStartSuccess())
	current, err := f.groups.Get(t.Context(), f.group["id"].(string))
	require.NoError(t, err)
	request := f.request("save-failure")
	request.Groups[0].Draft = current.Group
	_, err = f.db.ExecContext(t.Context(), `CREATE TRIGGER fail_recording_start_save
		BEFORE UPDATE ON system_settings
		WHEN NEW.key = 'studio_v2_workspace'
		BEGIN SELECT RAISE(ABORT, 'save failure'); END`)
	require.NoError(t, err)

	failed, err := f.start.Start(t.Context(), request)
	require.NoError(t, err)
	require.Equal(t, recordingplan.SchemaOperationFailed, failed.Status)
	require.Equal(t, "retry_same_request", failed.NextAction)
	require.False(t, failed.Groups[0].Applied)
	require.Zero(t, f.activation.Calls())
	after, err := f.groups.Get(t.Context(), f.group["id"].(string))
	require.NoError(t, err)
	require.Equal(t, current.Group.Revision, after.Group.Revision)
	require.Empty(t, after.Group.AppliedRevision)

	_, err = f.db.ExecContext(t.Context(), `DROP TRIGGER fail_recording_start_save`)
	require.NoError(t, err)
	retried, err := f.start.Start(t.Context(), request)
	require.NoError(t, err)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, retried.Status)
	require.True(t, retried.Groups[0].Applied)
	require.True(t, retried.Devices[0].Activated)
}

func TestRecordingStartReadinessFailureHasNoApplyEffect(t *testing.T) {
	managed := newManagedGroupFixture(t, "sqlite")
	token := managed.preview(t)
	confirm := managed.schemaRequest()
	confirm["token"], confirm["operation_id"] = token.Token, token.OperationID
	response := performJSONRequest(t, managed.router, http.MethodPost, managed.confirmPath(), confirm)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	reloaded := performJSONRequest(t, managed.router, http.MethodGet, writeGroupsPath+"/"+managed.group["id"].(string), nil)
	require.Equal(t, http.StatusOK, reloaded.Code, reloaded.Body.String())
	managed.data, managed.group = groupSaveData(t, decodeJSONBody(t, reloaded))
	f := newRecordingStartServiceFixture(t, managed, recordingStartSuccess())
	f.workspace.WithReadinessServices(recordingStartBlockedDeviceStub{}, nil, nil, nil).WithWriteGroupReadiness(f.groups)

	failed, err := f.start.Start(t.Context(), f.request("readiness-failure"))
	require.NoError(t, err)
	require.Equal(t, recordingplan.SchemaOperationFailed, failed.Status)
	require.Equal(t, "scope_not_ready", failed.Reason)
	require.Equal(t, "review_device", failed.NextAction)
	require.False(t, failed.Groups[0].Applied)
	require.Zero(t, f.activation.Calls())
	current, err := f.groups.Get(t.Context(), f.group["id"].(string))
	require.NoError(t, err)
	require.Empty(t, current.Group.AppliedRevision)
}

func TestRecordingStartGetEnforcesWorkspaceAndActionScope(t *testing.T) {
	f := newManagedGroupFixture(t, "sqlite")
	token := f.preview(t)
	path := strings.TrimSuffix(writeGroupsPath, "/write-groups") + "/recording-start/operations/"
	wrongAction := performJSONRequest(t, f.router, http.MethodGet, path+token.OperationID, nil)
	require.Equal(t, http.StatusNotFound, wrongAction.Code, wrongAction.Body.String())

	ledger := recordingplan.NewService(recordingplan.NewSQLRepository(f.db))
	foreign, _, err := ledger.ClaimRecordingStart(t.Context(), "foreign-workspace", "foreign-request", strings.Repeat("b", 64), "{}")
	require.NoError(t, err)
	foreignScope := performJSONRequest(t, f.router, http.MethodGet, path+foreign.OperationID, nil)
	require.Equal(t, http.StatusNotFound, foreignScope.Code, foreignScope.Body.String())
}

func TestRecordingStartRejectsChangedDigestWithoutExtraEffects(t *testing.T) {
	f := newRecordingStartFixture(t, recordingStartSuccess())
	request := f.request("changed-digest")
	completed, err := f.start.Start(t.Context(), request)
	require.NoError(t, err)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, completed.Status)
	changed := request
	changed.Groups = append([]workspace.RecordingStartGroupIntent(nil), request.Groups...)
	changed.Groups[0].ExpectedConnectorRevision = "connector-changed"

	_, err = f.start.Start(t.Context(), changed)
	require.ErrorIs(t, err, recordingplan.ErrSchemaOperationMismatch)
	require.Equal(t, 1, f.activation.Calls())
	var versions int
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT count(*) FROM write_group_versions WHERE group_id=?`, f.group["id"]).Scan(&versions))
	require.Equal(t, 1, versions)
}

func TestRecordingStartRejectsForeignDeviceWithoutClaimingOperation(t *testing.T) {
	f := newRecordingStartFixture(t, recordingStartSuccess())
	request := f.request("foreign-device")
	request.DeviceIDs = []string{"device-foreign"}

	_, err := f.start.Start(t.Context(), request)
	require.ErrorIs(t, err, workspace.ErrRecordingStartInvalid)
	require.Zero(t, f.activation.Calls())
	var operations, versions int
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT count(*) FROM managed_schema_operations WHERE action='recording_start'`).Scan(&operations))
	require.Zero(t, operations)
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT count(*) FROM write_group_versions WHERE group_id=?`, f.group["id"]).Scan(&versions))
	require.Zero(t, versions)
}

func TestRecordingStartRejectsSourceTagConnectorAndSchemaDrift(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*testing.T, *recordingStartFixture)
	}{
		{
			name: "source",
			mutate: func(t *testing.T, f *recordingStartFixture) {
				_, err := f.db.ExecContext(t.Context(), `UPDATE points SET address='40002' WHERE id='point-A'`)
				require.NoError(t, err)
			},
		},
		{
			name: "tagtype",
			mutate: func(t *testing.T, f *recordingStartFixture) {
				_, err := f.db.ExecContext(t.Context(), `UPDATE tags SET data_type='int64' WHERE id='tag-A'`)
				require.NoError(t, err)
			},
		},
		{
			name: "connector",
			mutate: func(t *testing.T, f *recordingStartFixture) {
				_, err := f.db.ExecContext(t.Context(), `UPDATE database_connectors SET identity_revision='connector-2' WHERE id='connector-A'`)
				require.NoError(t, err)
			},
		},
		{
			name: "schema",
			mutate: func(t *testing.T, f *recordingStartFixture) {
				destination := f.group["destination"].(map[string]any)
				target, err := sql.Open("sqlite", f.file)
				require.NoError(t, err)
				target.SetMaxOpenConns(1)
				t.Cleanup(func() { require.NoError(t, target.Close()) })
				_, err = target.ExecContext(t.Context(), `DROP TABLE "`+destination["table_name"].(string)+`"`)
				require.NoError(t, err)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := newRecordingStartFixture(t, recordingStartFailure())
			request := f.request("drift-" + test.name)
			partial, err := f.start.Start(t.Context(), request)
			require.NoError(t, err)
			require.Equal(t, recordingplan.SchemaOperationPartial, partial.Status)
			require.True(t, partial.Groups[0].Applied)
			appliedRevision := partial.Groups[0].AppliedRevision
			test.mutate(t, &f)

			retried, err := f.start.Start(t.Context(), request)
			require.NoError(t, err)
			require.NotEqual(t, recordingplan.SchemaOperationSucceeded, retried.Status)
			require.Equal(t, partial.OperationID, retried.OperationID)
			require.True(t, retried.Groups[0].Applied)
			require.Equal(t, appliedRevision, retried.Groups[0].AppliedRevision)
			require.Equal(t, 1, f.activation.Calls(), "drift must not activate again")
			var versions int
			require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT count(*) FROM write_group_versions WHERE group_id=?`, f.group["id"]).Scan(&versions))
			require.Equal(t, 1, versions, "drift must not apply a second version")
		})
	}
}
