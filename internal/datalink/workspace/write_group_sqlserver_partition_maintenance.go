package workspace

import (
	"context"
	"errors"
	"fmt"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
)

// NormalizeSQLServerPartitionInspection projects an existing native-monthly
// group table back to the immutable confirmed layout only when the actual
// receipt and the ledger-authorized current projection both prove the same
// catalog. Unrecorded drift remains an error.
func (s *WriteGroupService) NormalizeSQLServerPartitionInspection(
	ctx context.Context,
	scope recordingplan.SchemaPreviewScope,
	state, receipt recordingplan.TargetTableInspection,
) (recordingplan.TargetTableInspection, error) {
	if state.Status != string(dbtarget.TableInspectionExists) || state.SQLServer == nil ||
		scope.GroupLayout == nil || scope.GroupLayout.SQLServer == nil || scope.GroupLayout.SQLServer.NativeMonthly == nil {
		return state, nil
	}
	if scope.GroupLayout.OwnerColumn != "" {
		hasOwner := false
		for _, column := range state.Columns {
			if column == scope.GroupLayout.OwnerColumn {
				hasOwner = true
				break
			}
		}
		if !hasOwner {
			return state, nil
		}
	}
	if s == nil || s.repo == nil || s.repo.db == nil {
		return state, fmt.Errorf("%w: repository unavailable", ErrWriteGroupValidation)
	}
	storageScope := recordingplan.StorageScope{
		WorkspaceID: scope.WorkspaceID, ConnectorID: scope.ConnectorID, Database: scope.Database,
		TableSchema: scope.Schema, StorageMode: string(WriteGroupStorageModeNativeMonthly),
		LayoutVersion: string(WriteGroupLayoutManagedLongV1),
	}
	lineage, err := recordingplan.ReadStorageLineage(ctx, s.repo.db, storageScope)
	if err != nil && !errors.Is(err, recordingplan.ErrStorageLineageNotFound) {
		return state, err
	}
	if !lineage.Confirmed() {
		// Initial post-DDL inspection precedes the first recorded confirmation.
		// Validate the exact initial catalog; no maintenance projection is trusted.
		_, err := recordingplan.VerifyManagedSQLServerCatalog(*scope.GroupLayout, scope.Database, state, receipt)
		return state, err
	}
	return recordingplan.NormalizeSQLServerPartitionCatalog(ctx, s.repo.db, storageScope, *scope.GroupLayout, state, receipt)
}

// populateSQLServerPartitionMaintenanceReadiness exposes the durable current
// SQL Server partition projection without changing the frozen initial layout.
func (s *WriteGroupService) populateSQLServerPartitionMaintenanceReadiness(ctx context.Context, connectorKind string, group *WriteGroup, result *WriteGroupReadiness) {
	if s == nil || s.repo == nil || s.repo.db == nil || group == nil || result == nil ||
		connectorKind != writeGroupConnectorKindSQLServer || group.StorageMode != WriteGroupStorageModeNativeMonthly {
		return
	}
	state, err := recordingplan.ReadSQLServerPartitionMaintenanceState(ctx, s.repo.db, StorageScopeOf(group))
	if err != nil {
		result.PartitionMaintenance = &WriteGroupPartitionMaintenanceReadiness{
			Status: WriteGroupMigrationPreviewStatusBlocked, Code: partitionMaintenanceBlockedCode, Reason: "state_unreadable",
		}
		return
	}
	if state == nil {
		return
	}
	result.SQLServerCurrent = recordingplan.CloneSQLServerPartitionEvidence(state.Current)
	result.PartitionMaintenance = &WriteGroupPartitionMaintenanceReadiness{
		Status: state.Status, Reason: state.Reason, UpdatedAt: state.UpdatedAt,
	}
	if state.Status == WriteGroupMigrationPreviewStatusBlocked {
		result.PartitionMaintenance.Code = partitionMaintenanceBlockedCode
	}
}
