package workspace

import (
	"context"
	"fmt"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
)

const partitionMaintenanceBlockedCode = "partition_maintenance_blocked"

// WriteGroupPartitionMaintenanceReadiness exposes the current projection
// of native monthly partition maintenance for operator visibility.
type WriteGroupPartitionMaintenanceReadiness struct {
	Status    string                                `json:"status"`
	Code      string                                `json:"code,omitempty"`
	Reason    string                                `json:"reason,omitempty"`
	UpdatedAt time.Time                             `json:"updated_at"`
	Current   *recordingplan.MySQLPartitionEvidence `json:"current,omitempty"`
}

// NormalizeMySQLPartitionInspection projects an existing MySQL native-monthly
// group table catalog state back to its confirmed initial layout for schema
// verification. It skips missing tables, receipts, and plain storage.
func (s *WriteGroupService) NormalizeMySQLPartitionInspection(
	ctx context.Context,
	scope recordingplan.SchemaPreviewScope,
	state recordingplan.TargetTableInspection,
) (recordingplan.TargetTableInspection, error) {
	if state.Status != string(dbtarget.TableInspectionExists) ||
		state.MySQL == nil ||
		scope.GroupLayout == nil ||
		scope.GroupLayout.MySQL == nil ||
		scope.GroupLayout.MySQL.NativeMonthly == nil {
		return state, nil
	}
	if scope.GroupLayout.OwnerColumn != "" {
		hasOwner := false
		for _, col := range state.Columns {
			if col == scope.GroupLayout.OwnerColumn {
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
		WorkspaceID:   scope.WorkspaceID,
		ConnectorID:   scope.ConnectorID,
		Database:      scope.Database,
		TableSchema:   scope.Schema,
		StorageMode:   string(WriteGroupStorageModeNativeMonthly),
		LayoutVersion: string(WriteGroupLayoutManagedLongV1),
	}
	return recordingplan.NormalizeMySQLPartitionCatalogForExplicitReconfirmation(ctx, s.repo.db, storageScope, *scope.GroupLayout, state, scope.ConnectorRevision)
}

func (s *WriteGroupService) populateMySQLPartitionMaintenanceReadiness(ctx context.Context, group *WriteGroup, result *WriteGroupReadiness) {
	if s == nil || s.repo == nil || s.repo.db == nil || group == nil || result == nil ||
		group.StorageMode != WriteGroupStorageModeNativeMonthly {
		return
	}
	state, err := recordingplan.ReadMySQLPartitionMaintenanceState(ctx, s.repo.db, StorageScopeOf(group))
	if err != nil || state == nil {
		return
	}
	var code string
	if state.Status == "blocked" {
		code = partitionMaintenanceBlockedCode
	}
	result.PartitionMaintenance = &WriteGroupPartitionMaintenanceReadiness{
		Status:    state.Status,
		Code:      code,
		Reason:    state.Reason,
		UpdatedAt: state.UpdatedAt,
		Current:   state.Current,
	}
}

func dbtargetMySQLPartitions(evidence *recordingplan.MySQLPartitionEvidence) *dbtarget.MySQLPartitionEvidence {
	if evidence == nil {
		return nil
	}
	result := &dbtarget.MySQLPartitionEvidence{
		Method:     evidence.Method,
		Expression: evidence.Expression,
		Partitions: make([]dbtarget.MySQLPartitionBoundary, len(evidence.Partitions)),
	}
	for i, partition := range evidence.Partitions {
		result.Partitions[i] = dbtarget.MySQLPartitionBoundary{
			Name:     partition.Name,
			LessThan: partition.LessThan,
			Ordinal:  partition.Ordinal,
		}
	}
	return result
}
