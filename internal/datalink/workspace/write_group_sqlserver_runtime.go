package workspace

import (
	"fmt"
	"reflect"
	"slices"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
)

// ValidateSQLServerRuntimeLayout verifies the frozen SQL Server physical proof against
// the canonical group without reading an editable destination or the network.
func ValidateSQLServerRuntimeLayout(group *WriteGroup, frozen *WriteGroupRuntimeLayout) error {
	if group == nil || frozen == nil || frozen.Dialect != dbtarget.SQLDialectSQLServer ||
		!group.isLongLayout() || group.StorageMode != WriteGroupStorageModePlain && group.StorageMode != WriteGroupStorageModeNativeMonthly ||
		group.Destination.StorageStrategy != WriteGroupStorageStrategyManaged ||
		group.Destination.ConnectorID == "" || group.Destination.ConnectorRevision == "" ||
		group.Destination.Database == "" || group.Destination.TableSchema == "" ||
		group.Destination.SchemaRevision == "" || len(group.Members) == 0 ||
		frozen.Hypertable != nil || frozen.MySQL != nil ||
		group.WritePolicy.DedupeCapability != dbtarget.GroupEffectReceipt || !frozen.ReceiptTableReady ||
		!LongRecordingShapeSupported(group, frozen.TagTypes) || !longRuntimeSupported(group) {
		return ErrWriteGroupValidation
	}
	proof := frozen.SQLServer
	if proof == nil || proof.ConnectorID != group.Destination.ConnectorID || proof.ConnectorRevision != group.Destination.ConnectorRevision ||
		proof.Group == nil || proof.Group.Inspection == nil || proof.Group.Evidence == nil ||
		proof.Receipt == nil || proof.Receipt.Inspection == nil || proof.Receipt.Evidence == nil {
		return ErrWriteGroupValidation
	}
	for _, member := range group.Members {
		if value, ok := frozen.TagTypes[member.TagID]; !ok || !value.IsValid() {
			return ErrWriteGroupValidation
		}
	}
	for name, physical := range map[string]*dbtarget.SQLServerManagedInspection{
		group.Destination.TableName: proof.Group,
		dbtarget.EffectReceiptTable: proof.Receipt,
	} {
		if physical.Inspection.Table != name ||
			physical.Inspection.Schema != group.Destination.TableSchema ||
			physical.Evidence.Database != group.Destination.Database ||
			physical.Evidence.Schema != group.Destination.TableSchema ||
			physical.Indexes == nil ||
			!slices.EqualFunc(physical.Indexes, physical.Inspection.Indexes, func(a, b dbtarget.IndexInfo) bool {
				return a.Name == b.Name && a.Unique == b.Unique && slices.Equal(a.Columns, b.Columns)
			}) {
			return ErrWriteGroupValidation
		}
	}
	table, owner, ok := LongStorageIdentity(group)
	if !ok || table != group.Destination.TableName {
		return ErrWriteGroupValidation
	}
	evidence := proof.Group.Evidence
	receiptEvidence := proof.Receipt.Evidence
	if !evidence.MetadataVerified || !receiptEvidence.Plain || !receiptEvidence.MetadataVerified ||
		evidence.TableType != "USER_TABLE" || receiptEvidence.TableType != "USER_TABLE" ||
		evidence.ServerVersion == "" || evidence.ServerEdition == "" || evidence.CapabilityRevision == "" || evidence.CompatibilityLevel == 0 ||
		receiptEvidence.ServerVersion != evidence.ServerVersion || receiptEvidence.ServerEdition != evidence.ServerEdition ||
		receiptEvidence.CapabilityRevision != evidence.CapabilityRevision || receiptEvidence.CompatibilityLevel != evidence.CompatibilityLevel ||
		receiptEvidence.EngineEdition != evidence.EngineEdition {
		return ErrWriteGroupValidation
	}
	if (group.StorageMode == WriteGroupStorageModeNativeMonthly) != (evidence.NativeMonthly != nil) ||
		receiptEvidence.NativeMonthly != nil || (evidence.NativeMonthly != nil) == evidence.Plain {
		return ErrWriteGroupValidation
	}
	if evidence.NativeMonthly != nil && (!evidence.PartitionInspected || !evidence.NativeMonthly.Exclusive ||
		!evidence.NativeMonthly.IndexesAligned || !evidence.NativeMonthly.FunctionExists || !evidence.NativeMonthly.SchemeExists) {
		return ErrWriteGroupValidation
	}

	options := recordingplan.SQLServerGroupSchemaOptions{
		Database:           group.Destination.Database,
		Schema:             group.Destination.TableSchema,
		Collation:          evidence.ColumnCollations["record_id"],
		ServerVersion:      evidence.ServerVersion,
		ServerEdition:      evidence.ServerEdition,
		CapabilityRevision: evidence.CapabilityRevision,
		EngineEdition:      evidence.EngineEdition,
		CompatibilityLevel: evidence.CompatibilityLevel,
	}
	if evidence.NativeMonthly != nil {
		options.NativeMonthly = recordingPlanSQLServerLayoutPartition(evidence.NativeMonthly)
	}
	layout, err := recordingplan.ManagedSQLServerLongGroupLayout(table, owner, options)
	if err != nil {
		return fmt.Errorf("%w: frozen sqlserver layout: %w", ErrWriteGroupValidation, err)
	}
	digest, err := recordingplan.VerifyManagedSQLServerCatalog(layout, group.Destination.Database,
		recordingPlanSQLServerInspection(proof.Group), recordingPlanSQLServerInspection(proof.Receipt))
	if err != nil || digest == "" || digest != group.Destination.SchemaDigest || digest != frozen.SchemaDigest ||
		!reflect.DeepEqual(frozen.Columns, proof.Group.Inspection.Columns) {
		return ErrWriteGroupValidation
	}
	return nil
}

func recordingPlanSQLServerInspection(physical *dbtarget.SQLServerManagedInspection) recordingplan.TargetTableInspection {
	if physical == nil || physical.Inspection == nil || physical.Evidence == nil {
		return recordingplan.TargetTableInspection{}
	}
	state := recordingPlanTableInspection(physical.Inspection)
	state.IndexesInspected = physical.Indexes != nil
	state.Indexes = nil
	for _, index := range physical.Indexes {
		state.Indexes = append(state.Indexes, recordingplan.TargetIndexInspection{
			Name:    index.Name,
			Columns: index.Columns,
			Unique:  index.Unique,
		})
	}
	e := physical.Evidence
	state.SQLServer = &recordingplan.TargetSQLServerTableEvidence{
		ServerVersion:         e.ServerVersion,
		ServerEdition:         e.ServerEdition,
		CapabilityRevision:    e.CapabilityRevision,
		Database:              e.Database,
		Schema:                e.Schema,
		EngineEdition:         e.EngineEdition,
		CompatibilityLevel:    e.CompatibilityLevel,
		MetadataVerified:      e.MetadataVerified,
		TableType:             e.TableType,
		ObjectID:              e.ObjectID,
		Plain:                 e.Plain,
		PrimaryKey:            e.PrimaryKey,
		ColumnCollations:      e.ColumnCollations,
		ColumnDefaults:        e.ColumnDefaults,
		ProvenanceJSONChecked: e.ProvenanceJSONChecked,
		RepairOperationID:     e.RepairOperationID,
		NativeMonthly:         recordingPlanSQLServerTargetPartition(e.NativeMonthly),
		PartitionInspected:    e.PartitionInspected,
	}
	return state
}

func recordingPlanSQLServerLayoutPartition(evidence *dbtarget.SQLServerManagedPartitionEvidence) *recordingplan.SQLServerPartitionEvidence {
	if evidence == nil {
		return nil
	}
	partitions := make([]recordingplan.SQLServerPartitionBoundary, len(evidence.Partitions))
	for i, partition := range evidence.Partitions {
		partitions[i] = recordingplan.SQLServerPartitionBoundary{Name: partition.Name, LessThan: partition.LessThan, Ordinal: partition.Ordinal}
	}
	return &recordingplan.SQLServerPartitionEvidence{
		Method: evidence.Method, Column: evidence.Column, FunctionName: evidence.FunctionName,
		SchemeName: evidence.SchemeName, Filegroup: evidence.Filegroup, Partitions: partitions,
	}
}

func recordingPlanSQLServerTargetPartition(evidence *dbtarget.SQLServerManagedPartitionEvidence) *recordingplan.TargetSQLServerPartitionEvidence {
	if evidence == nil {
		return nil
	}
	partitions := make([]recordingplan.SQLServerPartitionBoundary, len(evidence.Partitions))
	for i, partition := range evidence.Partitions {
		partitions[i] = recordingplan.SQLServerPartitionBoundary{Name: partition.Name, LessThan: partition.LessThan, Ordinal: partition.Ordinal}
	}
	return &recordingplan.TargetSQLServerPartitionEvidence{
		Method: evidence.Method, Column: evidence.Column, FunctionName: evidence.FunctionName,
		SchemeName: evidence.SchemeName, Filegroup: evidence.Filegroup, Partitions: partitions,
		FunctionExists: evidence.FunctionExists, SchemeExists: evidence.SchemeExists,
		Exclusive: evidence.Exclusive, IndexesAligned: evidence.IndexesAligned,
	}
}
