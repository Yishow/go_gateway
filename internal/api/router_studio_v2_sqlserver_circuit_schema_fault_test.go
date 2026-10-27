package api

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
	"testing"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"

	"github.com/stretchr/testify/require"
)

// TestSQLServerCircuitSchemaFaultIntegration keeps the actual public preview,
// operation ledger, SQL Server transaction and live post-catalog verification.
// Only the transport fault is injected; no catalog result is manufactured.
func TestSQLServerCircuitSchemaFaultIntegration(t *testing.T) {
	fixture := SQLServerCircuitFixtureOrSkip(t)
	for _, loseACK := range []bool{false, true} {
		name := "partial_ddl_rollback"
		if loseACK {
			name = "committed_ddl_lost_ack"
		}
		t.Run(name, func(t *testing.T) {
			f := newSQLServerCircuitRawFixture(t, fixture, "bool", true)
			capability, err := f.targetService.ProbeSQLServerRecordingAtRevision(t.Context(), "connector-A", "connector-1")
			require.NoError(t, err)
			if !capability.SupportsMissingRepair {
				t.Skip("NOT_RUN: SQL Server scoped repair transaction gate is closed")
			}
			token := previewLong(t, f.managedGroupFixture, f.instantID)
			require.NotEmpty(t, token.Statements)
			require.NotNil(t, token.GroupLayout)
			// Never inject DDL faults into a previously existing table.
			before, err := dbtarget.InspectOpenSQLServerManagedTable(t.Context(), f.target, capability, token.GroupLayout.TableName)
			require.NoError(t, err)
			require.Equal(t, dbtarget.TableInspectionMissing, before.Inspection.Status)
			receiptBefore, err := dbtarget.InspectOpenSQLServerManagedTable(t.Context(), f.target, capability, dbtarget.EffectReceiptTable)
			require.NoError(t, err)
			faultDB, faults := sqlServerCircuitFaultDB(t, fixture)
			t.Cleanup(func() { require.NoError(t, faultDB.Close()) })
			var rejected, committed atomic.Int64
			if loseACK {
				faults.Lose.Store(true)
				hook := func() { committed.Add(1) }
				faults.AfterCommit.Store(&hook)
			} else {
				reject := func(query string) error {
					if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(query)), "CREATE INDEX") {
						rejected.Add(1)
						return errors.New("fixture interrupted schema batch after CREATE TABLE")
					}
					return nil
				}
				faults.Reject.Store(&reject)
			}
			target := recordingplan.SQLServerSchemaApplyTarget{
				Inspect: func(ctx context.Context, name string) (recordingplan.TargetTableInspection, error) {
					physical, err := dbtarget.InspectOpenSQLServerManagedTable(ctx, f.target, capability, name)
					return sqlServerCircuitTargetInspection(physical), err
				},
				Execute: func(ctx context.Context, saved *recordingplan.SchemaPreviewToken, statements []string) (recordingplan.SchemaExecution, error) {
					execution, err := dbtarget.RunSQLServerManagedRepairBatch(ctx, faultDB, capability, dbtarget.SQLServerManagedRepairRequest{
						TableName: saved.GroupLayout.TableName, OperationID: saved.OperationID, Statements: statements,
					})
					if execution == nil {
						return recordingplan.SchemaExecution{}, err
					}
					return recordingplan.SchemaExecution{Committed: execution.Committed, RolledBack: execution.RolledBack}, err
				},
			}
			repo := recordingplan.NewSQLRepository(f.db)
			ledger := recordingplan.NewService(repo)
			op, _, err := ledger.ApplySQLServerSchemaPreview(t.Context(), &token, target)
			faults.Lose.Store(false)
			require.NoError(t, err)
			require.Equal(t, token.OperationID, op.OperationID)
			if loseACK {
				require.EqualValues(t, 1, committed.Load(), "the real SQL transaction committed before losing its ACK")
				require.Equal(t, recordingplan.SchemaOperationSucceeded, op.Status)
				require.NotEmpty(t, op.VerifiedDigest)
				registerSQLServerCircuitCreatedTableCleanup(t, f.managedGroupFixture, f.targetService, token)
			} else {
				require.EqualValues(t, 1, rejected.Load(), "the partial DDL seam must actually execute")
				require.Equal(t, recordingplan.SchemaOperationFailed, op.Status)
				missing, err := dbtarget.InspectOpenSQLServerManagedTable(t.Context(), f.target, capability, token.GroupLayout.TableName)
				require.NoError(t, err)
				require.Equal(t, dbtarget.TableInspectionMissing, missing.Inspection.Status, "CREATE TABLE must roll back with the failed index")
			}
			stored, err := repo.GetSchemaOperation(t.Context(), f.record.ID, token.OperationID)
			require.NoError(t, err)
			require.Equal(t, op.Status, stored.Status)
			// Restart and replay retain the same operation; no second DDL batch.
			target.Execute = func(context.Context, *recordingplan.SchemaPreviewToken, []string) (recordingplan.SchemaExecution, error) {
				t.Fatal("a completed operation replay must not execute DDL")
				return recordingplan.SchemaExecution{}, nil
			}
			replay, _, err := recordingplan.NewService(recordingplan.NewSQLRepository(f.db)).ApplySQLServerSchemaPreview(t.Context(), &token, target)
			require.NoError(t, err)
			require.Equal(t, op.OperationID, replay.OperationID)
			require.Equal(t, op.Status, replay.Status)
			receiptAfter, err := dbtarget.InspectOpenSQLServerManagedTable(t.Context(), f.target, capability, dbtarget.EffectReceiptTable)
			require.NoError(t, err)
			if receiptBefore.Inspection.Status == dbtarget.TableInspectionExists {
				require.Equal(t, receiptBefore.Evidence.ObjectID, receiptAfter.Evidence.ObjectID, "existing receipts must retain their object identity")
			}
		})
	}
}

func sqlServerCircuitTargetInspection(physical *dbtarget.SQLServerManagedInspection) recordingplan.TargetTableInspection {
	if physical == nil || physical.Inspection == nil || physical.Evidence == nil {
		return recordingplan.TargetTableInspection{}
	}
	i, e := physical.Inspection, physical.Evidence
	state := recordingplan.TargetTableInspection{
		Status: string(i.Status), IndexesInspected: physical.Indexes != nil,
		ColumnTypes: map[string]string{}, ColumnNullable: map[string]bool{}, ColumnPrimaryKey: map[string]bool{},
		SQLServer: &recordingplan.TargetSQLServerTableEvidence{
			Database: e.Database, Schema: e.Schema, ServerVersion: e.ServerVersion, ServerEdition: e.ServerEdition,
			CapabilityRevision: e.CapabilityRevision, EngineEdition: e.EngineEdition, CompatibilityLevel: e.CompatibilityLevel,
			MetadataVerified: e.MetadataVerified, TableType: e.TableType, ObjectID: e.ObjectID, Plain: e.Plain,
			PrimaryKey: e.PrimaryKey, ColumnCollations: e.ColumnCollations, ColumnDefaults: e.ColumnDefaults,
			ProvenanceJSONChecked: e.ProvenanceJSONChecked, RepairOperationID: e.RepairOperationID, PartitionInspected: e.PartitionInspected,
		},
	}
	for _, column := range i.Columns {
		state.Columns = append(state.Columns, column.Name)
		state.ColumnTypes[column.Name], state.ColumnNullable[column.Name], state.ColumnPrimaryKey[column.Name] = column.DataType, column.Nullable, column.PrimaryKey
	}
	for _, index := range physical.Indexes {
		state.Indexes = append(state.Indexes, recordingplan.TargetIndexInspection{Name: index.Name, Columns: index.Columns, Unique: index.Unique})
	}
	return state
}
