package api

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouptestwrite"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"

	mssql "github.com/microsoft/go-mssqldb"
	"github.com/stretchr/testify/require"
)

// TestSQLServerGroupProbeLifecycleIntegration is the original 4.4 entrypoint:
// public preview/confirm token lifecycle, exact write/readback, operation-
// owned cleanup and honest replay semantics on the disposable SQL Server
// fixture. Both plain and native get a fresh fixture.
func TestSQLServerGroupProbeLifecycleIntegration(t *testing.T) {
	fixture := SQLServerCircuitFixtureOrSkip(t)
	for _, mode := range []workspace.WriteGroupStorageMode{workspace.WriteGroupStorageModePlain, workspace.WriteGroupStorageModeNativeMonthly} {
		t.Run(string(mode), func(t *testing.T) {
			f := NewSQLServerCircuitPlainFixture(t, fixture, SQLServerCircuitFixtureOptions{
				StorageMode: mode, TagType: "uint64",
			})
			runSQLServerGroupProbeLifecycle(t, f)
		})
	}
}

func runSQLServerGroupProbeLifecycle(t *testing.T, f *sqlServerCircuitPlainFixture) {
	t.Helper()
	ctx := t.Context()
	instantID, _, _ := f.GroupIDs()

	workspaceRecord, err := f.workspace.GetOrCreate(ctx)
	require.NoError(t, err)
	live := newSQLServerGroupProbeService(t, f, f.targetService)
	ws := grouptestwrite.Workspace{ID: workspaceRecord.ID, Revision: workspaceRecord.DatabaseSetupRevision}
	saved, err := f.groups.Get(ctx, instantID)
	require.NoError(t, err)
	dest := glmDestination(saved.Group)
	beforeReceipts := glmReceiptDigests(t, f.target, dest.TableSchema)

	preview, err := live.Preview(ctx, ws, instantID)
	require.NoError(t, err)
	require.NotEmpty(t, preview.Token)
	require.NotEmpty(t, preview.OperationID)
	require.Equal(t, "sqlserver", preview.Target.Dialect)

	outcome, err := live.Confirm(ctx, ws, grouptestwrite.Confirmation{
		Token: preview.Token, OperationID: preview.OperationID, GroupID: instantID,
	})
	require.NoError(t, err)
	require.NotNil(t, outcome.Operation)
	require.Equal(t, grouptestwrite.WriteVerified, outcome.Operation.WriteOutcome)
	require.Equal(t, grouptestwrite.CleanupCleaned, outcome.Operation.CleanupStatus)
	require.Equal(t, beforeReceipts, glmReceiptDigests(t, f.target, dest.TableSchema),
		"normal cleanup must not remove or add receipts outside the operation-owned IDs")
	require.Zero(t, sqlServerGroupProbeRowCount(t, f, dest, prefixIDs(preview.RecordIDs)))

	replay, err := live.Confirm(ctx, ws, grouptestwrite.Confirmation{
		Token: preview.Token, OperationID: preview.OperationID, GroupID: instantID,
	})
	require.NoError(t, err)
	require.Equal(t, outcome.Operation.WriteOutcome, replay.Operation.WriteOutcome)
	require.Equal(t, outcome.Operation.CleanupStatus, replay.Operation.CleanupStatus)
	require.Equal(t, beforeReceipts, glmReceiptDigests(t, f.target, dest.TableSchema))
}

// TestSQLServerGroupProbeLifecycleCleanupDenied uses an actual lostcommit
// Reject(DELETE) fault on the saved metadata proof. The write/readback must
// verify, cleanup must fail, and the exact operation rows/receipts must remain.
func TestSQLServerGroupProbeLifecycleCleanupDenied(t *testing.T) {
	fixture := SQLServerCircuitFixtureOrSkip(t)
	f := NewSQLServerCircuitPlainFixture(t, fixture)
	ctx := t.Context()
	instantID, _, _ := f.GroupIDs()
	saved, err := f.groups.Get(ctx, instantID)
	require.NoError(t, err)
	dest := glmDestination(saved.Group)
	beforeReceipts := glmReceiptDigests(t, f.target, dest.TableSchema)

	faultDB, faults := sqlServerCircuitFaultDB(t, fixture)
	t.Cleanup(func() { require.NoError(t, faultDB.Close()) })
	capability, err := f.targetService.ProbeSQLServerRecordingAtRevision(ctx, dest.ConnectorID, dest.ConnectorRevision)
	require.NoError(t, err)
	faultDestinations := sqlServerFaultDestinations{inner: f.targetService, db: faultDB, capability: capability}
	live := newSQLServerGroupProbeService(t, f, faultDestinations)
	reject := func(query string) error {
		if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(query)), "DELETE") {
			return mssql.Error{Number: 229, Class: 14, Message: "fixture denies operation-owned cleanup"}
		}
		return nil
	}
	faults.Reject.Store(&reject)

	workspaceRecord, err := f.workspace.GetOrCreate(ctx)
	require.NoError(t, err)
	ws := grouptestwrite.Workspace{ID: workspaceRecord.ID, Revision: workspaceRecord.DatabaseSetupRevision}
	preview, err := live.Preview(ctx, ws, instantID)
	require.NoError(t, err)
	require.True(t, f.createdTable, "cleanup fault must use a table this test created")
	outcome, err := live.Confirm(ctx, ws, grouptestwrite.Confirmation{
		Token: preview.Token, OperationID: preview.OperationID, GroupID: instantID,
	})
	require.NoError(t, err)
	require.NotNil(t, outcome.Operation)
	require.Equal(t, grouptestwrite.WriteVerified, outcome.Operation.WriteOutcome)
	require.Equal(t, grouptestwrite.CleanupFailed, outcome.Operation.CleanupStatus)
	detail := sqlServerGroupProbeDetail(t, outcome.Operation.Detail)
	require.NotEmpty(t, detail.EffectKeys)
	require.NotEmpty(t, detail.EffectDigests)
	require.Equal(t, len(preview.RecordIDs), sqlServerGroupProbeRowCount(t, f, dest, preview.RecordIDs),
		"failed cleanup must retain the exact operation-owned rows")
	for _, key := range detail.EffectKeys {
		require.Equal(t, detail.EffectDigests[key], glmReceiptDigests(t, f.target, dest.TableSchema)[key],
			"failed cleanup must retain the exact operation-owned receipt")
	}
	for key, digest := range beforeReceipts {
		require.Equal(t, digest, glmReceiptDigests(t, f.target, dest.TableSchema)[key],
			"failed cleanup must preserve pre-existing shared receipts")
	}

	replay, err := live.Confirm(ctx, ws, grouptestwrite.Confirmation{
		Token: preview.Token, OperationID: preview.OperationID, GroupID: instantID,
	})
	require.NoError(t, err)
	require.Equal(t, outcome.Operation.OperationID, replay.Operation.OperationID)
	require.Equal(t, outcome.Operation.WriteOutcome, replay.Operation.WriteOutcome)
	require.Equal(t, outcome.Operation.CleanupStatus, replay.Operation.CleanupStatus)
	require.Equal(t, len(preview.RecordIDs), sqlServerGroupProbeRowCount(t, f, dest, preview.RecordIDs))

	rejectInsert := func(query string) error {
		if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(query)), "INSERT") {
			return mssql.Error{Number: 229, Class: 14, Message: "fixture denies rewrite during cleanup repair"}
		}
		return nil
	}
	faults.Reject.Store(&rejectInsert)

	repaired, err := live.Confirm(ctx, ws, grouptestwrite.Confirmation{
		Token:       preview.Token,
		OperationID: preview.OperationID,
		GroupID:     instantID,
		Action:      "cleanup_repair",
	})
	require.NoError(t, err)
	require.NotNil(t, repaired)
	require.NotNil(t, repaired.Operation)
	require.Equal(t, outcome.Operation.OperationID, repaired.Operation.OperationID)
	require.Equal(t, grouptestwrite.WriteVerified, repaired.Operation.WriteOutcome)
	require.Equal(t, grouptestwrite.CleanupCleaned, repaired.Operation.CleanupStatus)
	require.Zero(t, sqlServerGroupProbeRowCount(t, f, dest, preview.RecordIDs), "all test rows must be removed")
	for _, key := range detail.EffectKeys {
		require.NotContains(t, glmReceiptDigests(t, f.target, dest.TableSchema), key, "operation receipts must be removed")
	}
	for key, digest := range beforeReceipts {
		require.Equal(t, digest, glmReceiptDigests(t, f.target, dest.TableSchema)[key], "pre-existing receipts must remain")
	}
}

// TestSQLServerGroupProbeLifecycleWriteDenialIsHonest covers a dropped owned
// table after proving this fixture created it and its object/owner identity
// still matches the frozen runtime proof.
func TestSQLServerGroupProbeLifecycleWriteDenialIsHonest(t *testing.T) {
	fixture := SQLServerCircuitFixtureOrSkip(t)
	f := NewSQLServerCircuitPlainFixture(t, fixture)
	ctx := t.Context()
	instantID, _, _ := f.GroupIDs()
	require.True(t, f.createdTable, "the dropped-table fault requires this test to have created the table")
	saved, err := f.groups.Get(ctx, instantID)
	require.NoError(t, err)
	dest := glmDestination(saved.Group)
	payload, err := f.store.RuntimeVersion(ctx, groupdelivery.GroupKey{
		WorkspaceID:   saved.Group.WorkspaceID,
		GroupID:       saved.Group.ID,
		GroupRevision: saved.Group.Revision,
	})
	require.NoError(t, err)
	var snapshot workspace.WriteGroupAppliedSnapshot
	require.NoError(t, json.Unmarshal(payload, &snapshot))
	require.NotNil(t, snapshot.RuntimeLayout)
	require.NotNil(t, snapshot.RuntimeLayout.SQLServer)
	capability, err := f.targetService.ProbeSQLServerRecordingAtRevision(ctx, dest.ConnectorID, dest.ConnectorRevision)
	require.NoError(t, err)
	current, err := dbtarget.InspectOpenSQLServerManagedTable(ctx, f.target, capability, dest.TableName)
	require.NoError(t, err)
	require.Equal(t, snapshot.RuntimeLayout.SQLServer.Group.Evidence.ObjectID, current.Evidence.ObjectID,
		"the actual table object must still be the frozen owned object")
	require.True(t, sqlServerGroupProbeOwnerDefault(current.Evidence.ColumnDefaults), "the actual table owner marker must remain gateway")
	var rows int
	require.NoError(t, f.target.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+sqlServerQualified(dest.TableSchema, dest.TableName)).Scan(&rows))
	require.Zero(t, rows, "only an empty fixture-created table may be removed for the denial case")
	_, err = f.target.ExecContext(ctx, "IF OBJECT_ID(@p1,N'U')=@p2 AND NOT EXISTS (SELECT 1 FROM "+
		sqlServerQualified(dest.TableSchema, dest.TableName)+") DROP TABLE "+
		sqlServerQualified(dest.TableSchema, dest.TableName), dest.TableSchema+"."+dest.TableName, current.Evidence.ObjectID)
	require.NoError(t, err)

	workspaceRecord, err := f.workspace.GetOrCreate(ctx)
	require.NoError(t, err)
	live := newSQLServerGroupProbeService(t, f, f.targetService)
	ws := grouptestwrite.Workspace{ID: workspaceRecord.ID, Revision: workspaceRecord.DatabaseSetupRevision}
	preview, previewErr := live.Preview(ctx, ws, instantID)
	if previewErr != nil {
		require.Nil(t, preview, "failed preview must yield no token or operation")
		var unsupported *grouptestwrite.UnsupportedError
		if errors.As(previewErr, &unsupported) {
			require.Equal(t, grouptestwrite.ReasonLayoutBlocked, unsupported.Reason, "dropped table must block layout honestly")
		} else {
			require.ErrorIs(t, previewErr, recordingplan.ErrTargetInspectionUnconfirmed)
		}
		return
	}
	outcome, err := live.Confirm(ctx, ws, grouptestwrite.Confirmation{
		Token: preview.Token, OperationID: preview.OperationID, GroupID: instantID,
	})
	require.Error(t, err)
	require.NotNil(t, outcome)
	require.NotNil(t, outcome.Operation)
	require.NotEqual(t, grouptestwrite.WriteVerified, outcome.Operation.WriteOutcome)
}

func newSQLServerGroupProbeService(t *testing.T, f *sqlServerCircuitPlainFixture, destinations grouptestwrite.Destinations) *grouptestwrite.Service {
	t.Helper()
	return grouptestwrite.New(grouptestwrite.Dependencies{
		Groups: f.groups, Tags: f.tags, Destinations: destinations,
		Inspector: dbtarget.NewReadOnlyTableInspector(f.targetService),
		Ledger:    recordingplan.NewService(recordingplan.NewSQLRepository(f.db)),
	}, grouptestwrite.Config{})
}

type sqlServerGroupProbeOperationDetail struct {
	Schema           string            `json:"schema"`
	Table            string            `json:"table"`
	OwnerColumn      string            `json:"owner_column"`
	OwnerValue       string            `json:"owner_value"`
	TableOwnerValue  string            `json:"table_owner_value"`
	ProvenanceColumn string            `json:"provenance_column"`
	WorkspaceID      string            `json:"workspace_id"`
	GroupID          string            `json:"group_id"`
	GroupRevision    string            `json:"group_revision"`
	SampleIDs        []string          `json:"sample_ids"`
	RecordIDs        []string          `json:"record_ids"`
	EffectKeys       []string          `json:"effect_keys"`
	EffectDigests    map[string]string `json:"effect_digests"`
}

func sqlServerGroupProbeDetail(t *testing.T, raw string) sqlServerGroupProbeOperationDetail {
	t.Helper()
	var detail sqlServerGroupProbeOperationDetail
	require.NoError(t, json.Unmarshal([]byte(raw), &detail))
	return detail
}

func sqlServerGroupProbeRowCount(t *testing.T, f *sqlServerCircuitPlainFixture, dest glmDestinationRef, recordIDs []string) int {
	t.Helper()
	if len(recordIDs) == 0 {
		return 0
	}
	query := "SELECT COUNT(*) FROM " + sqlServerQualified(dest.TableSchema, dest.TableName) +
		" WHERE record_id IN (" + sqlServerParameters(len(recordIDs)) + ")"
	var count int
	require.NoError(t, f.target.QueryRowContext(t.Context(), query, sqlServerArgs(recordIDs)...).Scan(&count))
	return count
}

func prefixIDs(ids []string) []string { return ids }

func sqlServerGroupProbeOwnerDefault(defaults map[string]*string) bool {
	for _, value := range defaults {
		if value != nil && *value == "gateway" {
			return true
		}
	}
	return false
}
