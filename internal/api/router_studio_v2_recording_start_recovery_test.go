package api

import (
	"context"
	"database/sql"
	"net/http"
	"sync"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

type recordingStartGatedActivation struct {
	delegate *recordingStartActivationStub
	entered  chan struct{}
	release  chan struct{}
	once     sync.Once
}

func (s *recordingStartGatedActivation) ActivateScope(ctx context.Context, deviceIDs, groupIDs []string) (*workspace.ActivationResponse, error) {
	s.once.Do(func() { close(s.entered) })
	select {
	case <-s.release:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	return s.delegate.ActivateScope(ctx, deviceIDs, groupIDs)
}

func (s *recordingStartGatedActivation) Calls() int {
	return s.delegate.Calls()
}

func TestRecordingStartConcurrentSameRequestHasOneOperationAndAppliedVersion(t *testing.T) {
	f := newRecordingStartFixture(t)
	request := f.request("concurrent-same-request")
	entered, release := make(chan struct{}), make(chan struct{})
	firstActivation := &recordingStartGatedActivation{
		delegate: &recordingStartActivationStub{workspace: f.record.ID}, entered: entered, release: release,
	}
	secondActivation := &recordingStartActivationStub{workspace: f.record.ID}
	firstStart := workspace.NewRecordingStartService(f.workspace, f.groups, firstActivation, f.ledger)
	secondLedger := recordingplan.NewService(recordingplan.NewSQLRepository(f.db))
	secondStart := workspace.NewRecordingStartService(f.workspace, f.groups, secondActivation, secondLedger)
	type startResult struct {
		operation *workspace.RecordingStartOperation
		err       error
	}
	firstDone := make(chan startResult, 1)
	go func() {
		operation, err := firstStart.Start(t.Context(), request)
		firstDone <- startResult{operation: operation, err: err}
	}()
	<-entered

	inFlight, err := secondStart.Start(t.Context(), request)
	require.NoError(t, err)
	require.Equal(t, recordingplan.SchemaOperationRunning, inFlight.Status)
	require.Equal(t, 0, secondActivation.Calls(), "live lease must not be taken over")

	close(release)
	first := <-firstDone
	require.NoError(t, first.err)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, first.operation.Status)
	require.Equal(t, first.operation.OperationID, inFlight.OperationID)

	replayed, err := secondStart.Start(t.Context(), request)
	require.NoError(t, err)
	require.Equal(t, first.operation.OperationID, replayed.OperationID)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, replayed.Status)
	require.Zero(t, secondActivation.Calls())
	var operations, versions int
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT count(*) FROM managed_schema_operations WHERE action='recording_start'`).Scan(&operations))
	require.Equal(t, 1, operations)
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT count(*) FROM write_group_versions WHERE group_id=?`, f.group["id"]).Scan(&versions))
	require.Equal(t, 1, versions)
}

func TestRecordingStartStaleLeaseIsTakenOverByRestartedService(t *testing.T) {
	f := newRecordingStartFixture(t)
	request := f.request("stale-lease-restart")
	entered, release := make(chan struct{}), make(chan struct{})
	firstActivation := &recordingStartGatedActivation{
		delegate: &recordingStartActivationStub{workspace: f.record.ID}, entered: entered, release: release,
	}
	firstStart := workspace.NewRecordingStartService(f.workspace, f.groups, firstActivation, f.ledger)
	secondActivation := &recordingStartActivationStub{workspace: f.record.ID}
	secondLedger := recordingplan.NewService(recordingplan.NewSQLRepository(f.db))
	secondStart := workspace.NewRecordingStartService(f.workspace, f.groups, secondActivation, secondLedger)
	type startResult struct {
		operation *workspace.RecordingStartOperation
		err       error
	}
	firstDone := make(chan startResult, 1)
	go func() {
		operation, err := firstStart.Start(t.Context(), request)
		firstDone <- startResult{operation: operation, err: err}
	}()
	<-entered

	var owner string
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT owner FROM managed_schema_operations WHERE action='recording_start'`).Scan(&owner))
	_, err := f.db.ExecContext(t.Context(), `UPDATE managed_schema_operations SET updated_at=? WHERE action='recording_start'`, time.Now().UTC().Add(-recordingplan.SchemaOperationLease-time.Minute))
	require.NoError(t, err)

	restarted, err := secondStart.Start(t.Context(), request)
	require.NoError(t, err)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, restarted.Status)
	require.Equal(t, 1, secondActivation.Calls())
	var newOwner string
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT owner FROM managed_schema_operations WHERE action='recording_start'`).Scan(&newOwner))
	require.NotEqual(t, owner, newOwner, "a restarted process must fence the stale owner")

	close(release)
	first := <-firstDone
	require.NoError(t, first.err)
	require.Equal(t, restarted.OperationID, first.operation.OperationID)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, first.operation.Status)
	var versions int
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT count(*) FROM write_group_versions WHERE group_id=?`, f.group["id"]).Scan(&versions))
	require.Equal(t, 1, versions)
}

func TestRecordingStartPartialRetryReusesAppliedRevision(t *testing.T) {
	f := newRecordingStartFixture(t, recordingStartFailure(), recordingStartSuccess())
	request := f.request("partial-retry")

	first, err := f.start.Start(t.Context(), request)
	require.NoError(t, err)
	require.Equal(t, recordingplan.SchemaOperationPartial, first.Status)
	require.True(t, first.Groups[0].Applied)
	require.NotEmpty(t, first.Groups[0].AppliedRevision)
	require.False(t, first.Devices[0].Activated)
	appliedRevision := first.Groups[0].AppliedRevision
	var versions int
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT count(*) FROM write_group_versions WHERE group_id=?`, f.group["id"]).Scan(&versions))
	require.Equal(t, 1, versions)

	second, err := f.start.Start(t.Context(), request)
	require.NoError(t, err)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, second.Status)
	require.True(t, second.Groups[0].Applied)
	require.Equal(t, appliedRevision, second.Groups[0].AppliedRevision)
	require.True(t, second.Devices[0].Activated)
	require.Equal(t, 2, f.activation.Calls())
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT count(*) FROM write_group_versions WHERE group_id=?`, f.group["id"]).Scan(&versions))
	require.Equal(t, 1, versions, "retry must not create another applied version")
}

func TestRecordingStartApplyCheckpointRollbackDoesNotClaimApplied(t *testing.T) {
	f := newRecordingStartFixture(t, recordingStartSuccess())
	request := f.request("checkpoint-rollback")
	_, err := f.db.ExecContext(t.Context(), `CREATE TRIGGER fail_recording_start_applied_checkpoint
		BEFORE UPDATE OF detail ON managed_schema_operations
		WHEN NEW.action = 'recording_start' AND instr(NEW.detail, '"applied":true') > 0
		BEGIN SELECT RAISE(ABORT, 'checkpoint failure'); END`)
	require.NoError(t, err)

	failed, err := f.start.Start(t.Context(), request)
	require.NoError(t, err)
	require.Equal(t, recordingplan.SchemaOperationFailed, failed.Status)
	require.Equal(t, "retry_same_request", failed.NextAction)
	require.False(t, failed.Groups[0].Applied, "a rolled back checkpoint must not report applied")
	require.Empty(t, failed.Groups[0].AppliedRevision)
	var versions int
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT count(*) FROM write_group_versions WHERE group_id=?`, f.group["id"]).Scan(&versions))
	require.Zero(t, versions, "the group Apply and its checkpoint must share one transaction")
	var detail string
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT detail FROM managed_schema_operations WHERE operation_id=?`, failed.OperationID).Scan(&detail))
	require.NotContains(t, detail, `"applied":true`, "persisted progress must not claim the rolled back Apply")

	_, err = f.db.ExecContext(t.Context(), `DROP TRIGGER fail_recording_start_applied_checkpoint`)
	require.NoError(t, err)
	retried, err := f.start.Start(t.Context(), request)
	require.NoError(t, err)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, retried.Status)
	require.True(t, retried.Groups[0].Applied)
	require.NotEmpty(t, retried.Groups[0].AppliedRevision)
	require.True(t, retried.Devices[0].Activated)
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT count(*) FROM write_group_versions WHERE group_id=?`, f.group["id"]).Scan(&versions))
	require.Equal(t, 1, versions)
}

func TestRecordingStartCompletedReplayAfterSQLConfigRestartIsReadOnly(t *testing.T) {
	f := newRecordingStartFixture(t, recordingStartSuccess())
	request := f.request("restart-replay")
	_, err := f.start.Start(t.Context(), request)
	require.NoError(t, err, "the original response is intentionally discarded to model a lost reply")
	var operationID string
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT operation_id FROM managed_schema_operations WHERE action='recording_start'`).Scan(&operationID))

	configPath := sqliteConfigPath(t, f.db)
	require.NoError(t, f.db.Close())
	reopened, err := sql.Open("sqlite", configPath)
	require.NoError(t, err)
	reopened.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, reopened.Close()) })
	f.db = reopened
	f.workspace = workspace.NewService(workspace.NewSQLRepository(reopened))
	targets := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(reopened))
	f.groups = workspace.NewWriteGroupService(f.workspace, workspace.NewSQLWriteGroupRepository(reopened)).
		WithTableInspector(dbtarget.NewReadOnlyTableInspector(targets))
	f.groups.WithManagedTableInspector(f.groups.ManagedTableInspector(targets))
	f.workspace.WithReadinessServices(recordingStartReadinessDeviceStub{}, nil, nil, nil).WithWriteGroupReadiness(f.groups)
	f.activation = &recordingStartActivationStub{workspace: f.record.ID}
	f.ledger = recordingplan.NewService(recordingplan.NewSQLRepository(reopened))
	f.start = workspace.NewRecordingStartService(f.workspace, f.groups, f.activation, f.ledger)
	f.router = NewRouter(&DatalinkServices{
		Workspace: f.workspace, WriteGroups: f.groups, DBTarget: targets, RecordingPlan: f.ledger,
	})

	replayed, err := f.start.Start(t.Context(), request)
	require.NoError(t, err)
	require.Equal(t, operationID, replayed.OperationID)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, replayed.Status)
	require.Zero(t, f.activation.Calls(), "completed history must be returned without reapplying or reactivating")

	operation := performJSONRequest(t, f.router, http.MethodGet, "/api/v1/datalink/studio-v2/workspace/recording-start/operations/"+operationID, nil)
	require.Equal(t, http.StatusOK, operation.Code, operation.Body.String())
}

func sqliteConfigPath(t *testing.T, db *sql.DB) string {
	t.Helper()
	var sequence int
	var name, path string
	require.NoError(t, db.QueryRowContext(t.Context(), `PRAGMA database_list`).Scan(&sequence, &name, &path))
	require.NotEmpty(t, path)
	return path
}
