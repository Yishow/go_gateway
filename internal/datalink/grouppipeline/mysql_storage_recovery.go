package grouppipeline

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	mysqldriver "github.com/go-sql-driver/mysql"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"
)

type mysqlStorageRestoreTargetOps struct {
	Inspect func(context.Context, *sql.DB, *dbtarget.MySQLRecordingCapability, string, string) (*dbtarget.MySQLManagedInspection, error)
	Execute func(context.Context, *sql.DB, *dbtarget.MySQLRecordingCapability, string, string) error
}

func mysqlStorageRestoreTarget(target groupdelivery.Target, database string) recordingplan.StorageRestoreTarget {
	return mysqlStorageRestoreTargetWithOps(target, database, mysqlStorageRestoreTargetOps{
		Inspect: dbtarget.InspectOpenMySQLManagedTable,
		Execute: dbtarget.ExecuteOpenMySQLManagedStatement,
	})
}

func mysqlStorageRestoreTargetWithOps(
	target groupdelivery.Target,
	database string,
	ops mysqlStorageRestoreTargetOps,
) recordingplan.StorageRestoreTarget {
	return recordingplan.StorageRestoreTarget{
		Inspect: func(ctx context.Context, table string) (recordingplan.TargetTableInspection, error) {
			if target.DB == nil || target.MySQLCapability == nil {
				return recordingplan.TargetTableInspection{}, dbtarget.ErrMySQLManagedServerMismatch
			}
			if ops.Inspect == nil {
				return recordingplan.TargetTableInspection{}, recordingplan.ErrTargetInspectionUnconfirmed
			}
			inspected, err := ops.Inspect(ctx, target.DB, target.MySQLCapability, database, table)
			state := mysqlRestoreTargetInspection(inspected)
			if err != nil {
				return state, mapMySQLStorageRestoreError(err)
			}
			if inspected == nil || inspected.Inspection == nil || inspected.Evidence == nil {
				return state, fmt.Errorf("%w: mysql catalog evidence is incomplete", recordingplan.ErrPreviewTokenStale)
			}
			switch inspected.Inspection.Status {
			case dbtarget.TableInspectionExists, dbtarget.TableInspectionMissing:
			case dbtarget.TableInspectionForbidden:
				return state, recordingplan.ErrTargetPermissionDenied
			default:
				return state, fmt.Errorf("%w: mysql catalog state is unconfirmed", recordingplan.ErrPreviewTokenStale)
			}
			return state, nil
		},
		Execute: func(ctx context.Context, statements []string) (recordingplan.SchemaExecution, error) {
			if len(statements) != 1 {
				return recordingplan.SchemaExecution{}, dbtarget.ErrMySQLManagedStatementRefused
			}
			if target.DB == nil || target.MySQLCapability == nil {
				return recordingplan.SchemaExecution{}, dbtarget.ErrMySQLManagedServerMismatch
			}
			if ops.Execute == nil {
				return recordingplan.SchemaExecution{}, recordingplan.ErrSchemaAdapterUnverified
			}
			if err := ops.Execute(ctx, target.DB, target.MySQLCapability, database, statements[0]); err != nil {
				return recordingplan.SchemaExecution{}, mapMySQLStorageRestoreError(err)
			}
			return recordingplan.SchemaExecution{Committed: 1}, nil
		},
	}
}

func mysqlRestoreTargetInspection(inspected *dbtarget.MySQLManagedInspection) recordingplan.TargetTableInspection {
	if inspected == nil || inspected.Inspection == nil {
		return recordingplan.TargetTableInspection{}
	}
	info := inspected.Inspection
	state := recordingplan.TargetTableInspection{
		Status: string(info.Status), Columns: make([]string, 0, len(info.Columns)),
		ColumnTypes: map[string]string{}, ColumnNullable: map[string]bool{}, ColumnPrimaryKey: map[string]bool{},
		IndexesInspected: info.IndexesInspected,
	}
	for _, column := range info.Columns {
		state.Columns = append(state.Columns, column.Name)
		state.ColumnTypes[column.Name] = column.DataType
		state.ColumnNullable[column.Name] = column.Nullable
		state.ColumnPrimaryKey[column.Name] = column.PrimaryKey
	}
	for _, index := range info.Indexes {
		state.Indexes = append(state.Indexes, recordingplan.TargetIndexInspection{
			Name: index.Name, Columns: index.Columns, Unique: index.Unique,
		})
	}
	if evidence := inspected.Evidence; evidence != nil {
		state.MySQL = &recordingplan.TargetMySQLTableEvidence{
			ServerFamily: evidence.ServerFamily, ServerVersion: evidence.ServerVersion,
			CapabilityRevision: evidence.CapabilityRevision, Database: evidence.Database,
			Engine: evidence.Engine, Charset: evidence.Charset, Collation: evidence.Collation, TableType: evidence.TableType,
			PrimaryKey: evidence.PrimaryKey, ColumnCollations: evidence.ColumnCollations,
			ColumnDefaults: evidence.ColumnDefaults, JSONColumns: evidence.JSONColumns,
			Partitions: mirrorMySQLPartitionEvidence(evidence.Partitions),
		}
	}
	return state
}

func mirrorMySQLPartitionEvidence(evidence *dbtarget.MySQLPartitionEvidence) *recordingplan.MySQLPartitionEvidence {
	if evidence == nil {
		return nil
	}
	var partitions []recordingplan.MySQLPartitionBoundary
	if evidence.Partitions != nil {
		partitions = make([]recordingplan.MySQLPartitionBoundary, len(evidence.Partitions))
		for i, boundary := range evidence.Partitions {
			partitions[i] = recordingplan.MySQLPartitionBoundary{
				Name:     boundary.Name,
				LessThan: boundary.LessThan,
				Ordinal:  boundary.Ordinal,
			}
		}
	}
	return &recordingplan.MySQLPartitionEvidence{
		Method:     evidence.Method,
		Expression: evidence.Expression,
		Partitions: partitions,
	}
}

func mapMySQLStorageRestoreError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, dbtarget.ErrMySQLManagedTableForbidden) || errors.Is(err, dbtarget.ErrSchemaExecutionPermissionDenied) || mysqlPermissionDenied(err) {
		return fmt.Errorf("%w: %w", recordingplan.ErrTargetPermissionDenied, err)
	}
	if errors.Is(err, dbtarget.ErrMySQLManagedTableUnconfirmed) || errors.Is(err, dbtarget.ErrMySQLManagedCapabilityChanged) ||
		errors.Is(err, dbtarget.ErrMySQLManagedServerMismatch) {
		return fmt.Errorf("%w: %w", recordingplan.ErrPreviewTokenStale, err)
	}
	if errors.Is(err, dbtarget.ErrMySQLManagedIndexUnusable) {
		return fmt.Errorf("%w: %w", recordingplan.ErrIncompatibleExistingTable, err)
	}
	return err
}

func mysqlPermissionDenied(err error) bool {
	var mysqlErr *mysqldriver.MySQLError
	if !errors.As(err, &mysqlErr) {
		return false
	}
	switch mysqlErr.Number {
	case 1044, 1142, 1143, 1227:
		return true
	default:
		return false
	}
}

func (r *StorageRecovery) mysqlConfirmedFallbackLayout(
	ctx context.Context,
	scope recordingplan.StorageScope,
	table, owner string,
) (*recordingplan.GroupSchemaLayout, *recordingplan.StorageLineage, error) {
	lineage, err := recordingplan.ReadStorageLineage(ctx, r.store.DB(), scope)
	if errors.Is(err, recordingplan.ErrStorageLineageNotFound) {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}
	if lineage.StorageScope != scope || !lineage.Confirmed() || lineage.Generation < 1 || lineage.Layout == nil || lineage.Layout.MySQL == nil ||
		lineage.TableName != table || lineage.OwnerColumn != owner || lineage.Layout.TableName != table || lineage.Layout.OwnerColumn != owner ||
		scope.Database != scope.TableSchema ||
		scope.StorageMode != string(workspace.WriteGroupStorageModePlain) && scope.StorageMode != string(workspace.WriteGroupStorageModeNativeMonthly) ||
		(scope.StorageMode == string(workspace.WriteGroupStorageModeNativeMonthly)) != (lineage.Layout.MySQL.NativeMonthly != nil) ||
		scope.LayoutVersion != string(workspace.WriteGroupLayoutManagedLongV1) {
		return nil, lineage, nil
	}
	canonical, err := recordingplan.ManagedMySQLLongGroupLayout(table, owner, *lineage.Layout.MySQL)
	if err != nil {
		return nil, lineage, fmt.Errorf("build canonical mysql recovery layout: %w", err)
	}
	if !reflect.DeepEqual(canonical, *lineage.Layout) {
		return nil, lineage, nil
	}
	encoded, err := json.Marshal(*lineage.Layout)
	if err != nil {
		return nil, lineage, fmt.Errorf("encode confirmed mysql storage layout: %w", err)
	}
	digest := sha256.Sum256(encoded)
	if lineage.LayoutDigest != hex.EncodeToString(digest[:]) {
		return nil, lineage, nil
	}
	layout := *lineage.Layout
	return &layout, lineage, nil
}

func mysqlCapabilityMatchesLayout(capability *dbtarget.MySQLRecordingCapability, layout *recordingplan.GroupSchemaLayout, database string) bool {
	if capability == nil || layout == nil || layout.MySQL == nil {
		return false
	}
	options := layout.MySQL
	return capability.Connected && capability.MetadataStatus == "verified" && capability.TransactionalEngine &&
		capability.Database == database && capability.ServerFamily == options.ServerFamily &&
		capability.ServerVersion == options.ServerVersion && capability.CapabilityRevision == options.CapabilityRevision
}
