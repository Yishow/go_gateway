package dbtarget

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/recordingplan"

	"github.com/stretchr/testify/require"
)

type mysqlLiveLostFinish struct{ *recordingplan.SQLRepository }

func (*mysqlLiveLostFinish) FinishSchemaOperation(context.Context, string, string, recordingplan.SchemaOperationResult) (*recordingplan.SchemaOperation, error) {
	return nil, errors.New("injected local finish acknowledgement lost")
}

func TestMySQLManagedDDLRealRecovery(t *testing.T) {
	for _, family := range []string{mysqlFixtureFamily, mariaDBFixtureFamily} {
		t.Run(family, func(t *testing.T) {
			for _, fault := range []string{"partial_index", "finish_ack_lost", "foreign_table"} {
				t.Run(fault, func(t *testing.T) { runMySQLManagedLiveRecovery(t, family, fault) })
			}
		})
	}
}

func runMySQLManagedLiveRecovery(t *testing.T, family, fault string) {
	cfg := loadMySQLLongFixtureConfig(t, family)
	db := openMySQLLongFixtureDB(t, cfg)
	requireMySQLFixtureFamily(t, db, family)
	inspector, capability := mysqlManagedLiveInspector(t, cfg)
	require.Equal(t, family, capability.ServerFamily)
	table, owner := uniqueLongTestTable(t), uniqueLongTestOwnerColumn(t)
	layout, err := recordingplan.ManagedMySQLLongGroupLayout(table, owner, recordingplan.MySQLGroupSchemaOptions{
		ServerFamily: family, ServerVersion: capability.ServerVersion, CapabilityRevision: capability.CapabilityRevision,
		Engine: "InnoDB", Collation: mysqlFamilyCollationFor(family),
	})
	require.NoError(t, err)
	require.False(t, fixtureTableExists(t, db, longTestReceiptTable), "dedicated fixture must begin without an unowned receipt")
	var owned []string
	t.Cleanup(func() { dropFixtureTablesOwned(t, db, owned) })
	path := filepath.Join(t.TempDir(), "live-ddl-ledger.db")
	local := mysqlManagedLocalLedger(t, path, true)
	repo := recordingplan.NewSQLRepository(local)
	service := recordingplan.NewService(repo)
	if fault == "finish_ack_lost" {
		service = recordingplan.NewService(&mysqlLiveLostFinish{repo})
	}
	actualDDL := []string{}
	failIndex := fault == "partial_index"
	target := recordingplan.SchemaApplyTarget{
		Inspect: func(ctx context.Context, table string) (recordingplan.TargetTableInspection, error) {
			return mysqlManagedLiveState(ctx, inspector, capability, table)
		},
		Execute: func(ctx context.Context, statements []string) (recordingplan.SchemaExecution, error) {
			require.Len(t, statements, 1, "MySQL DDL must be dispatched one statement at a time")
			statement := statements[0]
			if failIndex && strings.HasPrefix(statement, "CREATE INDEX ") {
				return recordingplan.SchemaExecution{}, errors.New("injected index unavailable before execution")
			}
			err := inspector.ExecuteMySQLManagedStatementAtRevision(ctx, capability.ConnectorID, capability.ConnectorRevision, cfg.Database, capability.CapabilityRevision, statement)
			if err == nil {
				actualDDL = append(actualDDL, statement)
				if strings.HasPrefix(statement, "CREATE TABLE `"+table+"`") {
					owned = append(owned, table)
				} else if strings.HasPrefix(statement, "CREATE TABLE `"+longTestReceiptTable+"`") {
					owned = append(owned, longTestReceiptTable)
				}
				return recordingplan.SchemaExecution{Committed: 1}, nil
			}
			return recordingplan.SchemaExecution{}, err
		},
	}
	if fault == "foreign_table" {
		_, err := db.ExecContext(t.Context(), "CREATE TABLE `"+table+"` (foreign_value INT) ENGINE=InnoDB")
		require.NoError(t, err)
		owned = append(owned, table)
	}
	scope := recordingplan.SchemaPreviewScope{
		WorkspaceID: "ws-live", WorkspaceRevision: "ws-rev-1", PlanID: "group-live", PlanRevision: "group-rev-1",
		ConnectorID: capability.ConnectorID, ConnectorRevision: capability.ConnectorRevision, Dialect: "mysql",
		Database: cfg.Database, Schema: cfg.Database, GroupLayout: &layout, SourceDigest: "source-rev-1",
	}
	token, err := service.PrepareSchemaPreview(t.Context(), scope, target.Inspect)
	if fault == "foreign_table" {
		require.ErrorIs(t, err, recordingplan.ErrIncompatibleExistingTable)
		require.Empty(t, actualDDL)
		var count int
		require.NoError(t, db.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND COLUMN_NAME='foreign_value'", cfg.Database, table).Scan(&count))
		require.Equal(t, 1, count, "the foreign table is left unchanged")
		return
	}
	require.NoError(t, err)
	require.Len(t, token.Statements, 3)
	op, _, err := service.ApplySchemaPreview(t.Context(), token, target)
	if fault == "partial_index" {
		require.NoError(t, err)
		require.Equal(t, recordingplan.SchemaOperationPartial, op.Status)
		require.Len(t, actualDDL, 1)
	} else {
		require.ErrorIs(t, err, recordingplan.ErrSchemaOperationUnacknowledged)
		require.Len(t, actualDDL, 3)
		_, err = local.ExecContext(t.Context(), `UPDATE managed_schema_operations SET updated_at=? WHERE operation_id=?`, time.Now().UTC().Add(-11*time.Minute), op.OperationID)
		require.NoError(t, err)
	}
	stored, err := repo.GetSchemaOperation(t.Context(), scope.WorkspaceID, token.OperationID)
	require.NoError(t, err)
	require.Contains(t, stored.Detail, `"sql_hash"`)
	require.Contains(t, stored.Detail, `"verified"`)
	require.NoError(t, local.Close())
	restarted := recordingplan.NewService(recordingplan.NewSQLRepository(mysqlManagedLocalLedger(t, path, false)))
	loaded, err := restarted.ValidatePreviewTokenForApply(t.Context(), token.Token, scope)
	require.NoError(t, err)
	failIndex = false
	finished, _, err := restarted.ApplySchemaPreview(t.Context(), loaded, target)
	require.NoError(t, err)
	require.Equal(t, token.OperationID, finished.OperationID)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, finished.Status)
	require.Len(t, actualDDL, 3, "recovery never re-executes the already applied CREATE TABLE or index")
	require.Len(t, finished.VerifiedDigest, 64)
	prepared, err := restarted.PrepareSchemaPreview(t.Context(), scope, target.Inspect)
	require.NoError(t, err)
	require.Empty(t, prepared.Statements, "real catalog proves every table and index, not a batch rollback claim")
}
