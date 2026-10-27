package grouppipeline

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

// SQLServerNativeSplit is the typed dbtarget SPLIT executor. The root wiring
// supplies it; it never accepts caller SQL.
type SQLServerNativeSplit func(context.Context, *sql.DB, *dbtarget.SQLServerRecordingCapability,
	*dbtarget.SQLServerManagedInspection, *dbtarget.SQLServerManagedPartitionEvidence) error

// MaintainSQLServerPartitions runs one bounded native maintenance pass over
// applied SQL Server scopes using the existing pipeline lifetime and ledger.
func (p *Pipeline) MaintainSQLServerPartitions(ctx context.Context) error {
	if p.recovery == nil || p.sqlServerSplit == nil || p.deps.Groups == nil || p.deps.Destinations == nil || !p.partitionMu.TryLock() {
		return nil
	}
	defer p.partitionMu.Unlock()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	listed, err := p.deps.Groups.List(ctx)
	if err != nil {
		return err
	}
	seen := map[recordingplan.StorageScope]bool{}
	var failures []error
	for _, group := range listed.Groups {
		if err := ctx.Err(); err != nil {
			return errors.Join(append(failures, err)...)
		}
		if group == nil || group.AppliedRevision == "" || group.StorageMode != workspace.WriteGroupStorageModeNativeMonthly {
			continue
		}
		snapshot, err := p.deps.Groups.ResolveAppliedAt(ctx, group.ID, p.config.Now())
		if err != nil || snapshot == nil {
			failures = append(failures, errors.New("sqlserver partition maintenance: applied snapshot unavailable"))
			continue
		}
		frozen, err := p.loadRuntimeSnapshot(ctx, groupKey(snapshot))
		if err != nil || frozen == nil {
			failures = append(failures, errors.New("sqlserver partition maintenance: frozen proof unavailable"))
			continue
		}
		if frozen.Group.StorageMode != workspace.WriteGroupStorageModeNativeMonthly ||
			frozen.RuntimeLayout.Dialect != dbtarget.SQLDialectSQLServer {
			continue
		}
		scope := workspace.StorageScopeOf(frozen.Group)
		if seen[scope] {
			continue
		}
		seen[scope] = true
		if err := p.maintainSQLServerPartitionScope(ctx, frozen); err != nil {
			failures = append(failures, err)
		}
	}
	return errors.Join(failures...)
}

func (p *Pipeline) maintainSQLServerPartitionScope(ctx context.Context, frozen *workspace.WriteGroupAppliedSnapshot) (err error) {
	group := frozen.Group
	if group == nil || group.Destination.StorageStrategy != workspace.WriteGroupStorageStrategyManaged ||
		workspace.StorageScopeOf(group) == (recordingplan.StorageScope{}) {
		return recordingplan.ErrSQLServerPartitionMaintenanceRequest
	}
	opened, err := p.deps.Destinations.OpenDestination(ctx, group.Destination.ConnectorID, group.Destination.ConnectorRevision)
	if err != nil {
		return err
	}
	if opened == nil {
		return dbtarget.ErrSQLServerManagedTableUnconfirmed
	}
	if opened.Close != nil {
		defer func() { err = errors.Join(err, opened.Close()) }()
	}
	if opened.Kind != schema.DatabaseConnectorKindSQLServer || opened.DB == nil || opened.SQLServer == nil {
		return nil
	}
	var maintenanceUnavailable error
	if !opened.SQLServer.SupportsNativePartition {
		if !opened.SQLServer.SupportsManagedLong {
			return nil
		}
		maintenanceUnavailable = recordingplan.ErrTargetInspectionUnconfirmed
		if !opened.SQLServer.Permissions.AlterAnyDataspace {
			maintenanceUnavailable = recordingplan.ErrTargetPermissionDenied
		}
	}
	target := recordingplan.SQLServerPartitionMaintenanceTarget{
		Inspect: func(ctx context.Context, table string) (recordingplan.TargetTableInspection, error) {
			if maintenanceUnavailable != nil {
				return recordingplan.TargetTableInspection{}, maintenanceUnavailable
			}
			physical, inspectErr := dbtarget.InspectOpenSQLServerManagedTable(ctx, opened.DB, opened.SQLServer, table)
			return sqlServerRestoreTargetInspection(physical), mapSQLServerStorageRestoreError(inspectErr)
		},
		Split: func(ctx context.Context, before, after *recordingplan.SQLServerPartitionEvidence) error {
			if maintenanceUnavailable != nil {
				return maintenanceUnavailable
			}
			beforeGroup, inspectErr := dbtarget.InspectOpenSQLServerManagedTable(ctx, opened.DB, opened.SQLServer, frozen.Group.Destination.TableName)
			if inspectErr != nil {
				return mapSQLServerStorageRestoreError(inspectErr)
			}
			actualDefinition := beforeGroup.Evidence.NativeMonthly
			if !beforeGroup.Evidence.PartitionInspected || !sqlServerPartitionDefinitionMatches(actualDefinition, before) {
				return fmt.Errorf("%w: sqlserver catalog changed between preflight and split", recordingplan.ErrIncompatibleExistingTable)
			}
			return p.sqlServerSplit(ctx, opened.DB, opened.SQLServer, beforeGroup, sqlServerPartitionEvidenceForDB(after))
		},
	}
	result, err := p.recovery.ledger.MaintainSQLServerPartitions(ctx, recordingplan.SQLServerPartitionMaintenanceRequest{
		Scope: workspace.StorageScopeOf(frozen.Group), GroupID: frozen.Group.ID,
		ConnectorRevision: frozen.Group.Destination.ConnectorRevision, At: p.config.Now(),
	}, target)
	if err != nil {
		return err
	}
	if result.Code != "" {
		return fmt.Errorf("%s: %s", result.Code, result.Reason)
	}
	return nil
}

func sqlServerPartitionEvidenceForDB(source *recordingplan.SQLServerPartitionEvidence) *dbtarget.SQLServerManagedPartitionEvidence {
	if source == nil {
		return nil
	}
	partitions := make([]dbtarget.SQLServerPartitionBoundary, len(source.Partitions))
	for i, partition := range source.Partitions {
		partitions[i] = dbtarget.SQLServerPartitionBoundary{Name: partition.Name, LessThan: partition.LessThan, Ordinal: partition.Ordinal}
	}
	return &dbtarget.SQLServerManagedPartitionEvidence{
		Method: source.Method, Column: source.Column, FunctionName: source.FunctionName, SchemeName: source.SchemeName,
		Filegroup: source.Filegroup, Partitions: partitions,
	}
}

func sqlServerPartitionDefinitionMatches(actual *dbtarget.SQLServerManagedPartitionEvidence, expected *recordingplan.SQLServerPartitionEvidence) bool {
	if actual == nil || expected == nil {
		return false
	}
	if actual.Method != expected.Method || actual.Column != expected.Column || actual.FunctionName != expected.FunctionName ||
		actual.SchemeName != expected.SchemeName || actual.Filegroup != expected.Filegroup ||
		len(actual.Partitions) != len(expected.Partitions) {
		return false
	}
	return slices.EqualFunc(actual.Partitions, expected.Partitions, func(a dbtarget.SQLServerPartitionBoundary, b recordingplan.SQLServerPartitionBoundary) bool {
		return a.Name == b.Name && a.LessThan == b.LessThan && a.Ordinal == b.Ordinal
	})
}
