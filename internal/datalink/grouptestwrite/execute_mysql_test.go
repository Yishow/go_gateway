package grouptestwrite

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

const mysqlReadbackOwner = "_gw_owner_0123456789abcdef01234567"

func mysqlReadbackPlan(t *testing.T, expected measurement.ExactValue, proof *dbtarget.MySQLManagedRowProof) *plan {
	t.Helper()
	return &plan{
		dialect: dbtarget.SQLDialectMySQL,
		columns: []dbtarget.ColumnInfo{
			{Name: "record_id", DataType: "CHAR(64)"},
			{Name: "value_decimal", DataType: "DECIMAL(65,30)"},
			{Name: mysqlReadbackOwner, DataType: "CHAR(16)"},
			{Name: "provenance", DataType: "JSON"},
		},
		members:     []memberFixture{{column: "value_decimal", kind: measurement.ExactDecimal, value: expected}},
		ownerColumn: mysqlReadbackOwner,
		mysql:       &mysqlTestPlan{proof: proof, members: []mysqlTestMember{{kind: measurement.ExactDecimal, value: expected}}},
		group: &workspace.WriteGroup{RowPolicy: workspace.WriteGroupRowPolicy{
			ProvenanceColumn: "provenance",
		}},
	}
}

func mysqlReadbackState(recordID string) progress {
	return progress{
		Schema: "main", Table: "gw_long_test", OwnerColumn: mysqlReadbackOwner, OwnerValue: "gw-test-op1",
		TableOwnerValue: "gateway", ProvenanceColumn: "provenance",
		WorkspaceID: "ws-1", GroupID: "group-1", GroupRevision: "rev-1",
		RecordIDs: []string{recordID}, SampleIDs: []string{"gw-test-op1-0"},
		EffectKeys: []string{"effect-1"}, EffectDigests: map[string]string{"effect-1": "digest-1"},
	}
}

func openMySQLReadbackDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	for _, statement := range []string{
		"CREATE TABLE `gw_long_test` (`record_id` TEXT PRIMARY KEY, `workspace_id` TEXT, `group_id` TEXT, `group_revision` TEXT, `provenance` TEXT, `" + mysqlReadbackOwner + "` TEXT, `value_decimal` TEXT)",
		"CREATE TABLE `gw_effect_receipts` (`effect_key` TEXT PRIMARY KEY, `payload_digest` TEXT NOT NULL, `committed_at` TEXT NOT NULL)",
	} {
		_, err := db.ExecContext(t.Context(), statement)
		require.NoError(t, err)
	}
	return db
}

func TestMySQLVerifyNormalizesPaddedDecimalForPlainAndNative(t *testing.T) {
	expected, err := measurement.NewDecimal("12000.4")
	require.NoError(t, err)
	plainProof := &dbtarget.MySQLManagedRowProof{ConnectorID: "conn", ConnectorRevision: "crev"}
	nativeProof := &dbtarget.MySQLManagedRowProof{ConnectorID: "conn", ConnectorRevision: "crev", Group: &dbtarget.MySQLManagedInspection{
		Evidence: &dbtarget.MySQLManagedTableEvidence{Partitions: &dbtarget.MySQLPartitionEvidence{
			Method: "RANGE COLUMNS", Expression: "bucket_start",
		}},
	}}
	for name, proof := range map[string]*dbtarget.MySQLManagedRowProof{"plain": plainProof, "native": nativeProof} {
		t.Run(name, func(t *testing.T) {
			db := openMySQLReadbackDB(t)
			recordID := strings.Repeat("a", 64)
			state := mysqlReadbackState(recordID)
			provenance := `{"sample_id":"gw-test-op1-0"}`
			_, err := db.ExecContext(t.Context(), "INSERT INTO `gw_long_test` (`record_id`,`workspace_id`,`group_id`,`group_revision`,`provenance`,`"+mysqlReadbackOwner+"`,`value_decimal`) VALUES (?,?,?,?,?,?,?)",
				recordID, state.WorkspaceID, state.GroupID, state.GroupRevision, provenance, "gateway", "12000.400000000000000000000000000000")
			require.NoError(t, err)
			_, err = db.ExecContext(t.Context(), "INSERT INTO `gw_effect_receipts` (`effect_key`,`payload_digest`,`committed_at`) VALUES (?,?,?)", "effect-1", "digest-1", "2026-01-01T00:00:00.000000Z")
			require.NoError(t, err)
			expectedRow := dbtarget.EncodedRow{RecordID: recordID, EffectKey: "effect-1", Cells: []dbtarget.EncodedCell{
				{Column: "record_id", Value: recordID},
				{Column: "value_decimal", Value: "12000.4"},
				{Column: mysqlReadbackOwner, Value: "gateway"},
				{Column: "provenance", Value: provenance},
			}}
			p := mysqlReadbackPlan(t, expected, proof)
			svc := New(Dependencies{}, Config{StepTimeout: time.Second, Now: func() time.Time { return time.Unix(0, 0) }})
			result := svc.verifyMySQL(t.Context(), &dbtarget.OpenedDestination{DB: db, Kind: schema.DatabaseConnectorKindMySQL}, p, state, []dbtarget.EncodedRow{expectedRow})
			require.Equal(t, WriteVerified, result.outcome, "padded DECIMAL must compare as the same exact rational, got %s", result.reason)
		})
	}
}

func TestMySQLCommittedAtHasExactlySixFractionalDigits(t *testing.T) {
	formatted := time.Date(2026, 1, 2, 3, 4, 5, 123456789, time.UTC).Format(mysqlTestCommittedAtLayout)
	require.Equal(t, "2026-01-02T03:04:05.123456Z", formatted)
	require.True(t, strings.HasSuffix(formatted, ".123456Z"))
}

func TestMySQLSecondInsertFailureCleansCommittedFirstRow(t *testing.T) {
	group := mysqlTestGroup()
	group.Destination.TableSchema = "main"
	layout := mysqlTestRuntimeLayout(t)
	svc, _ := mysqlTestService(t, group, layout)
	p, err := svc.buildMySQLTestPlan(t.Context(), group, &schema.DatabaseConnector{ID: "conn", Kind: schema.DatabaseConnectorKindMySQL, IdentityRevision: "crev"})
	require.NoError(t, err)
	ledger := recordingplan.NewService(recordingplan.NewMemoryRepository())
	svc.deps.Ledger = ledger
	scope := recordingplan.TestWriteScope{
		WorkspaceID: group.WorkspaceID, WorkspaceRevision: "setup-1", GroupID: group.ID, GroupRevision: group.Revision,
		ConnectorID: "conn", ConnectorRevision: "crev", Dialect: "mysql", Database: "db",
		Schema: group.Destination.TableSchema, Table: group.Destination.TableName,
	}
	token, err := ledger.PrepareTestWritePreview(t.Context(), scope, p.digest)
	require.NoError(t, err)
	op, _, err := ledger.ClaimTestWrite(t.Context(), token)
	require.NoError(t, err)

	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	columns := append([]string{"record_id", "workspace_id", "group_id", "group_revision"}, p.readColumns()...)
	quoted := make([]string, len(columns))
	for i, column := range columns {
		quoted[i] = "`" + strings.ReplaceAll(column, "`", "``") + "`"
	}
	_, err = db.ExecContext(t.Context(), "CREATE TABLE `"+p.group.Destination.TableName+"` ("+strings.Join(quoted, " TEXT,")+" TEXT)")
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), "CREATE TABLE `gw_effect_receipts` (`effect_key` TEXT PRIMARY KEY, `payload_digest` TEXT NOT NULL, `committed_at` TEXT NOT NULL)")
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), "INSERT INTO `"+p.group.Destination.TableName+"` (`record_id`,`workspace_id`,`group_id`,`group_revision`,`provenance`,`"+p.ownerColumn+"`) VALUES (?,?,?,?,?,?)",
		strings.Repeat("c", 64), "production", "production", "rev", `{"sample_id":"production"}`, "gateway")
	require.NoError(t, err)

	inserts := 0
	svc.insertRow = func(ctx context.Context, target *sql.DB, req dbtarget.GroupInsertRequest) (dbtarget.GroupInsertResult, error) {
		inserts++
		if inserts == 2 {
			return dbtarget.GroupInsertResult{}, &dbtarget.GroupInsertError{Phase: dbtarget.InsertPhasePreCommit, Cause: errors.New("second row failed")}
		}
		values := make([]any, len(columns))
		for i, column := range columns {
			for _, cell := range req.Row.Cells {
				if strings.EqualFold(cell.Column, column) {
					values[i] = cell.Value
					break
				}
			}
		}
		placeholders := strings.TrimSuffix(strings.Repeat("?,", len(columns)), ",")
		if _, err := target.ExecContext(ctx, "INSERT INTO `"+p.group.Destination.TableName+"` ("+strings.Join(quoted, ",")+") VALUES ("+placeholders+")", values...); err != nil {
			return dbtarget.GroupInsertResult{}, err
		}
		_, err := target.ExecContext(ctx, "INSERT INTO `gw_effect_receipts` (`effect_key`,`payload_digest`,`committed_at`) VALUES (?,?,?)", req.Row.EffectKey, req.PayloadDigest, req.CommittedAt)
		return dbtarget.GroupInsertResult{}, err
	}
	state := progress{Phase: phaseClaimed, OwnerValue: ownerValue(op.OperationID)}
	written, finished, err := svc.writeMySQLPhase(t.Context(), op, p, token.CreatedAt, &dbtarget.OpenedDestination{DB: db, Kind: schema.DatabaseConnectorKindMySQL}, &state)
	require.NoError(t, err)
	require.Nil(t, finished)
	require.Equal(t, WriteUnverified, written.outcome)
	require.Equal(t, 2, inserts)
	cleanup := svc.cleanupPhase(t.Context(), op, p, &dbtarget.OpenedDestination{DB: db, Kind: schema.DatabaseConnectorKindMySQL}, &state)
	require.Equal(t, CleanupCleaned, cleanup.status, "%+v", state)
	var production int
	require.NoError(t, db.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM `"+p.group.Destination.TableName+"` WHERE `record_id` = ?", strings.Repeat("c", 64)).Scan(&production))
	require.Equal(t, 1, production)
	var testRows int
	require.NoError(t, db.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM `"+p.group.Destination.TableName+"` WHERE `record_id` IN (?,?)", state.RecordIDs[0], state.RecordIDs[1]).Scan(&testRows))
	require.Zero(t, testRows)
}

func TestMySQLFirstInsertPreCommitFailureIsTerminalFailed(t *testing.T) {
	group := mysqlTestGroup()
	group.Destination.TableSchema = "main"
	layout := mysqlTestRuntimeLayout(t)
	svc, _ := mysqlTestService(t, group, layout)
	p, err := svc.buildMySQLTestPlan(t.Context(), group, &schema.DatabaseConnector{ID: "conn", Kind: schema.DatabaseConnectorKindMySQL, IdentityRevision: "crev"})
	require.NoError(t, err)
	ledger := recordingplan.NewService(recordingplan.NewMemoryRepository())
	svc.deps.Ledger = ledger
	scope := recordingplan.TestWriteScope{
		WorkspaceID: group.WorkspaceID, WorkspaceRevision: "setup-1", GroupID: group.ID, GroupRevision: group.Revision,
		ConnectorID: "conn", ConnectorRevision: "crev", Dialect: "mysql", Database: "db",
		Schema: group.Destination.TableSchema, Table: group.Destination.TableName,
	}
	token, err := ledger.PrepareTestWritePreview(t.Context(), scope, p.digest)
	require.NoError(t, err)
	op, _, err := ledger.ClaimTestWrite(t.Context(), token)
	require.NoError(t, err)

	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	columns := append([]string{"record_id", "workspace_id", "group_id", "group_revision"}, p.readColumns()...)
	quoted := make([]string, len(columns))
	for i, column := range columns {
		quoted[i] = "`" + strings.ReplaceAll(column, "`", "``") + "`"
	}
	_, err = db.ExecContext(t.Context(), "CREATE TABLE `"+p.group.Destination.TableName+"` ("+strings.Join(quoted, " TEXT,")+" TEXT)")
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), "CREATE TABLE `gw_effect_receipts` (`effect_key` TEXT PRIMARY KEY, `payload_digest` TEXT NOT NULL, `committed_at` TEXT NOT NULL)")
	require.NoError(t, err)

	inserts := 0
	svc.insertRow = func(context.Context, *sql.DB, dbtarget.GroupInsertRequest) (dbtarget.GroupInsertResult, error) {
		inserts++
		return dbtarget.GroupInsertResult{}, &dbtarget.GroupInsertError{Phase: dbtarget.InsertPhasePreCommit, Cause: dbtarget.ErrExistingRowDiffers}
	}
	state := progress{Phase: phaseClaimed, OwnerValue: ownerValue(op.OperationID)}
	written, finished, err := svc.writeMySQLPhase(t.Context(), op, p, token.CreatedAt, &dbtarget.OpenedDestination{DB: db, Kind: schema.DatabaseConnectorKindMySQL}, &state)
	require.NoError(t, err)
	require.Equal(t, 1, inserts)
	require.Zero(t, written.outcome)
	require.NotNil(t, finished, "a known-empty failed target must terminate before cleanup")
	require.Equal(t, recordingplan.SchemaOperationFailed, finished.Status)
	require.Equal(t, WriteFailed, finished.WriteOutcome)
	require.Equal(t, reasonRowRejected, finished.Reason)
	require.Equal(t, CleanupNotAttempted, finished.CleanupStatus)
	require.Equal(t, op.PayloadDigest, finished.PayloadDigest)
	require.NotEmpty(t, finished.Detail)
}
