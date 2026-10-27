package workspace

import (
	"context"
	"fmt"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
)

// ManagedWriteGroupSQLServerProbeInspector reads live SQL Server recording
// capability facts at one saved connector revision.
type ManagedWriteGroupSQLServerProbeInspector interface {
	ProbeSQLServerRecordingAtRevision(ctx context.Context, connectorID, expectedRevision string) (*dbtarget.SQLServerRecordingCapability, error)
}

// managedSQLServerSchemaLayout builds a SQL Server long layout from live
// capability facts and the frozen native partition evidence. It never accepts
// caller-supplied server facts or a shared partition object.
func (s *WriteGroupService) managedSQLServerSchemaLayout(ctx context.Context, group *WriteGroup) (*recordingplan.GroupSchemaLayout, error) {
	if s == nil || s.managedTableInspector == nil {
		return nil, fmt.Errorf("%w: managed inspector is unavailable", ErrWriteGroupValidation)
	}
	prober, ok := s.managedTableInspector.(ManagedWriteGroupSQLServerProbeInspector)
	if !ok {
		return nil, fmt.Errorf("%w: sqlserver capability probe is not supported by inspector", ErrWriteGroupValidation)
	}
	capability, err := prober.ProbeSQLServerRecordingAtRevision(ctx, group.Destination.ConnectorID, group.Destination.ConnectorRevision)
	if err != nil {
		return nil, fmt.Errorf("%w: probe sqlserver capability failed: %w", ErrWriteGroupValidation, err)
	}
	if capability == nil || !capability.Connected {
		return nil, fmt.Errorf("%w: sqlserver destination is not connected", ErrWriteGroupValidation)
	}
	if capability.MetadataStatus != "verified" {
		return nil, fmt.Errorf("%w: sqlserver destination metadata is not verified (status=%s)", ErrWriteGroupValidation, capability.MetadataStatus)
	}
	if capability.ConnectorID != group.Destination.ConnectorID || capability.ConnectorRevision != group.Destination.ConnectorRevision ||
		capability.Database != group.Destination.Database || capability.Schema != group.Destination.TableSchema {
		return nil, fmt.Errorf("%w: sqlserver capability scope mismatch", ErrWriteGroupValidation)
	}
	if !capability.SupportsManagedLong {
		return nil, fmt.Errorf("%w: sqlserver managed long is not available for this destination", ErrWriteGroupValidation)
	}
	if group.StorageMode != WriteGroupStorageModePlain && group.StorageMode != WriteGroupStorageModeNativeMonthly {
		return nil, fmt.Errorf("%w: sqlserver requires plain or native monthly managed long schema", ErrWriteGroupValidation)
	}
	table, owner, ok := LongStorageIdentity(group)
	if !ok {
		return nil, fmt.Errorf("%w: group is not a long group", ErrWriteGroupValidation)
	}
	options := recordingplan.SQLServerGroupSchemaOptions{
		Database:           capability.Database,
		Schema:             capability.Schema,
		ServerVersion:      capability.ServerVersion,
		ServerEdition:      capability.ServerEdition,
		CapabilityRevision: capability.CapabilityRevision,
		EngineEdition:      capability.EngineEdition,
		CompatibilityLevel: capability.CompatibilityLevel,
	}
	if group.StorageMode == WriteGroupStorageModeNativeMonthly {
		if !capability.SupportsNativePartition {
			reasons := capability.Reasons
			if len(reasons) == 0 {
				reasons = []string{"native_partition_unverified"}
			}
			return nil, &recordingplan.SQLServerPartitionUnavailableError{Reasons: reasons}
		}
		options.NativeMonthly, err = s.sqlServerNativeMonthlyLayout(ctx, group)
		if err != nil {
			return nil, err
		}
	}
	layout, err := recordingplan.ManagedSQLServerLongGroupLayout(table, owner, options)
	if err != nil {
		return nil, fmt.Errorf("%w: build sqlserver long layout: %w", ErrWriteGroupValidation, err)
	}
	return &layout, nil
}
