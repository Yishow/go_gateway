package grouppipeline

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/recordingplan"
)

type sqlServerStorageRestoreTargetOps struct {
	Inspect func(context.Context, *sql.DB, *dbtarget.SQLServerRecordingCapability, string) (*dbtarget.SQLServerManagedInspection, error)
	Execute func(context.Context, *sql.DB, *dbtarget.SQLServerRecordingCapability, dbtarget.SQLServerManagedRepairRequest) (*dbtarget.SchemaStatementExecution, error)
}

func sqlServerStorageRestoreTarget(target groupdelivery.Target, schemaName string) recordingplan.StorageRestoreTarget {
	return sqlServerStorageRestoreTargetWithOps(target, schemaName, sqlServerStorageRestoreTargetOps{
		Inspect: dbtarget.InspectOpenSQLServerManagedTable,
		Execute: dbtarget.RunSQLServerManagedRepairBatch,
	})
}

func sqlServerStorageRestoreTargetWithOps(target groupdelivery.Target, schemaName string, ops sqlServerStorageRestoreTargetOps) recordingplan.StorageRestoreTarget {
	bound := func() bool {
		return target.DB != nil && target.SQLServerCapability != nil && target.SQLServerCapability.Schema == schemaName
	}
	return recordingplan.StorageRestoreTarget{
		Inspect: func(ctx context.Context, table string) (recordingplan.TargetTableInspection, error) {
			if !bound() || ops.Inspect == nil {
				return recordingplan.TargetTableInspection{}, recordingplan.ErrTargetInspectionUnconfirmed
			}
			physical, err := ops.Inspect(ctx, target.DB, target.SQLServerCapability, table)
			state := sqlServerRestoreTargetInspection(physical)
			if err != nil {
				return state, mapSQLServerStorageRestoreError(err)
			}
			if physical == nil || physical.Inspection == nil || physical.Evidence == nil {
				return state, recordingplan.ErrTargetInspectionUnconfirmed
			}
			switch physical.Inspection.Status {
			case dbtarget.TableInspectionExists, dbtarget.TableInspectionMissing:
				return state, nil
			case dbtarget.TableInspectionForbidden:
				return state, recordingplan.ErrTargetPermissionDenied
			default:
				return state, recordingplan.ErrTargetInspectionUnconfirmed
			}
		},
		Execute: func(context.Context, []string) (recordingplan.SchemaExecution, error) {
			return recordingplan.SchemaExecution{}, recordingplan.ErrSchemaAdapterUnverified
		},
		ExecuteSQLServerRepair: func(ctx context.Context, tableName, operationID string, statements []string) (recordingplan.SchemaExecution, error) {
			if !bound() || ops.Execute == nil || !target.SQLServerCapability.SupportsManagedLong || !target.SQLServerCapability.SupportsMissingRepair ||
				tableName == "" || tableName == dbtarget.EffectReceiptTable || operationID == "" {
				return recordingplan.SchemaExecution{}, recordingplan.ErrSchemaAdapterUnverified
			}
			result, err := ops.Execute(ctx, target.DB, target.SQLServerCapability, dbtarget.SQLServerManagedRepairRequest{
				TableName: tableName, OperationID: operationID, Statements: statements,
			})
			if result == nil {
				return recordingplan.SchemaExecution{}, mapSQLServerStorageRestoreError(err)
			}
			return recordingplan.SchemaExecution{Committed: result.Committed, RolledBack: result.RolledBack}, mapSQLServerStorageRestoreError(err)
		},
		CanCreate: func(context.Context) (bool, error) {
			if !bound() {
				return false, nil
			}
			c := target.SQLServerCapability
			return c.SupportsManagedLong && c.SupportsMissingRepair && c.Permissions.Verified && c.Permissions.CreateTable && c.Permissions.AlterSchema, nil
		},
	}
}

func sqlServerRestoreTargetInspection(physical *dbtarget.SQLServerManagedInspection) recordingplan.TargetTableInspection {
	if physical == nil || physical.Inspection == nil {
		return recordingplan.TargetTableInspection{}
	}
	i := physical.Inspection
	state := recordingplan.TargetTableInspection{
		Status: string(i.Status), ColumnTypes: map[string]string{}, ColumnNullable: map[string]bool{}, ColumnPrimaryKey: map[string]bool{},
		IndexesInspected: physical.Indexes != nil,
	}
	for _, column := range i.Columns {
		state.Columns = append(state.Columns, column.Name)
		state.ColumnTypes[column.Name], state.ColumnNullable[column.Name], state.ColumnPrimaryKey[column.Name] = column.DataType, column.Nullable, column.PrimaryKey
	}
	for _, index := range physical.Indexes {
		state.Indexes = append(state.Indexes, recordingplan.TargetIndexInspection{Name: index.Name, Columns: index.Columns, Unique: index.Unique})
	}
	if e := physical.Evidence; e != nil {
		state.SQLServer = &recordingplan.TargetSQLServerTableEvidence{
			ServerVersion: e.ServerVersion, ServerEdition: e.ServerEdition, CapabilityRevision: e.CapabilityRevision,
			Database: e.Database, Schema: e.Schema, EngineEdition: e.EngineEdition, CompatibilityLevel: e.CompatibilityLevel,
			MetadataVerified: e.MetadataVerified, TableType: e.TableType, ObjectID: e.ObjectID, Plain: e.Plain,
			PrimaryKey: e.PrimaryKey, ColumnCollations: e.ColumnCollations, ColumnDefaults: e.ColumnDefaults, ProvenanceJSONChecked: e.ProvenanceJSONChecked,
			RepairOperationID: e.RepairOperationID,
			NativeMonthly:     sqlServerRestoreTargetPartition(e.NativeMonthly), PartitionInspected: e.PartitionInspected,
		}
	}
	return state
}

func sqlServerRestoreTargetPartition(evidence *dbtarget.SQLServerManagedPartitionEvidence) *recordingplan.TargetSQLServerPartitionEvidence {
	if evidence == nil {
		return nil
	}
	partitions := make([]recordingplan.SQLServerPartitionBoundary, len(evidence.Partitions))
	for i, partition := range evidence.Partitions {
		partitions[i] = recordingplan.SQLServerPartitionBoundary{Name: partition.Name, LessThan: partition.LessThan, Ordinal: partition.Ordinal}
	}
	return &recordingplan.TargetSQLServerPartitionEvidence{
		Method: evidence.Method, Column: evidence.Column, FunctionName: evidence.FunctionName, SchemeName: evidence.SchemeName,
		Filegroup: evidence.Filegroup, Partitions: partitions, FunctionExists: evidence.FunctionExists, SchemeExists: evidence.SchemeExists,
		Exclusive: evidence.Exclusive, IndexesAligned: evidence.IndexesAligned,
	}
}

func mapSQLServerStorageRestoreError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, dbtarget.ErrSQLServerManagedTableForbidden), errors.Is(err, dbtarget.ErrSchemaExecutionPermissionDenied):
		return fmt.Errorf("%w: %w", recordingplan.ErrTargetPermissionDenied, err)
	case errors.Is(err, dbtarget.ErrSQLServerManagedScopeChanged):
		return fmt.Errorf("%w: %w", recordingplan.ErrPreviewTokenStale, err)
	case errors.Is(err, dbtarget.ErrSQLServerManagedIndexUnusable), errors.Is(err, dbtarget.ErrSQLServerManagedPhysicalMismatch):
		return fmt.Errorf("%w: %w", recordingplan.ErrIncompatibleExistingTable, err)
	default:
		return err
	}
}
