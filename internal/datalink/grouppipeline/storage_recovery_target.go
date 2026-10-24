package grouppipeline

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
)

// restoreTarget adapts the destination delivery opened to the restore engine.
func restoreTarget(target groupdelivery.Target, schemaName string, hypertable bool) recordingplan.StorageRestoreTarget {
	return recordingplan.StorageRestoreTarget{
		Inspect: func(ctx context.Context, table string) (recordingplan.TargetTableInspection, error) {
			return inspectOpenTarget(ctx, target, schemaName, table, hypertable)
		},
		Execute: func(ctx context.Context, statements []string) (recordingplan.SchemaExecution, error) {
			execution, err := dbtarget.RunManagedSchemaBatch(ctx, target.DB, target.Kind, schemaName, statements)
			if errors.Is(err, dbtarget.ErrSchemaExecutionPermissionDenied) {
				err = fmt.Errorf("%w: %w", recordingplan.ErrTargetPermissionDenied, err)
			}
			if execution == nil {
				return recordingplan.SchemaExecution{RolledBack: true}, err
			}
			return recordingplan.SchemaExecution{Committed: execution.Committed, RolledBack: execution.RolledBack}, err
		},
		CanCreate: func(ctx context.Context) (bool, error) {
			return dbtarget.CanCreateManagedTable(ctx, target.DB, target.Kind, schemaName)
		},
	}
}

// inspectOpenTarget reads one table through the opened destination. The shared
// long table also gets its indexes and, for a hypertable layout, the read-only
// TimescaleDB capability and hypertable evidence.
func inspectOpenTarget(ctx context.Context, target groupdelivery.Target, schemaName, table string, hypertable bool) (recordingplan.TargetTableInspection, error) {
	long := table != receiptTableName
	info, err := dbtarget.InspectOpenManagedTable(ctx, target.DB, target.Kind, schemaName, table, long)
	if err != nil {
		return recordingplan.TargetTableInspection{}, err
	}
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
		state.Indexes = append(state.Indexes, recordingplan.TargetIndexInspection{Name: index.Name, Columns: index.Columns, Unique: index.Unique})
	}
	if long && hypertable && target.Kind == schema.DatabaseConnectorKindPostgres {
		evidence, err := timescaleEvidence(ctx, target.DB, schemaName, table, info.Status == dbtarget.TableInspectionExists)
		if err != nil {
			return recordingplan.TargetTableInspection{}, err
		}
		state.Timescale = evidence
	}
	return state, nil
}

func timescaleEvidence(ctx context.Context, db *sql.DB, schemaName, table string, exists bool) (*recordingplan.TargetTimescaleEvidence, error) {
	capability, err := dbtarget.ProbeTimescale(ctx, db, schemaName)
	if err != nil {
		return nil, err
	}
	evidence := &recordingplan.TargetTimescaleEvidence{
		Available: capability.Available, Reasons: capability.Reasons, ExtensionSchema: capability.ExtensionSchema, Version: capability.Version,
	}
	if !exists {
		return evidence, nil
	}
	hypertable, err := dbtarget.InspectHypertable(ctx, db, schemaName, table)
	if err != nil {
		return nil, err
	}
	if hypertable != nil {
		evidence.Hypertable = &recordingplan.TargetHypertable{
			TimeColumn: hypertable.TimeColumn, ChunkInterval: hypertable.ChunkInterval, Dimensions: hypertable.Dimensions,
			CompressionEnabled: hypertable.CompressionEnabled, Jobs: hypertable.Jobs,
		}
	}
	return evidence, nil
}

// ResolveStale ends restores whose owner never came back, from the catalog of
// their destination alone, so a gateway that restarted mid-restore records
// what the table proves. It runs no statement and starts no restore.
func (r *StorageRecovery) ResolveStale(ctx context.Context) error {
	operations, err := r.ledger.StaleStorageRepairs(ctx)
	if err != nil {
		return err
	}
	var errs []error
	for _, operation := range operations {
		if err := r.resolveStaleOperation(ctx, operation); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (r *StorageRecovery) resolveStaleOperation(ctx context.Context, operation recordingplan.SchemaOperation) error {
	request, ok := recordingplan.StorageRepairRequestOf(operation, "")
	if !ok {
		return nil
	}
	lineage, err := recordingplan.ReadStorageLineage(ctx, r.store.DB(), request.Scope)
	if errors.Is(err, recordingplan.ErrStorageLineageNotFound) {
		// The operation names a scope this gateway never saved; there is no
		// confirmed table to judge it against.
		return nil
	}
	if err != nil || !lineage.Confirmed() {
		return err
	}
	request.ConnectorRevision = lineage.ConnectorRevision
	if r.destinations == nil {
		return nil
	}
	opened, err := r.destinations.OpenDestination(ctx, request.Scope.ConnectorID, request.ConnectorRevision)
	if err != nil {
		// An unreachable or blocked destination keeps the operation as it is.
		return nil
	}
	defer func() {
		if opened.Close != nil {
			_ = opened.Close() //nolint:errcheck // releasing a finished connection; nothing to recover
		}
	}()
	dialect, ok := dialectName(opened.Kind)
	if !ok {
		return nil
	}
	request.Dialect = dialect
	target := groupdelivery.Target{DB: opened.DB, Kind: opened.Kind}
	release := r.lock(request.Scope)
	defer release()
	_, err = r.ledger.ResolveStaleStorageRepair(ctx, request, restoreTarget(target, request.Scope.TableSchema, lineage.Layout.Partition != nil))
	return err
}
