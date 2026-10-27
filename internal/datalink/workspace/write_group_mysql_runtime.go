package workspace

import (
	"fmt"
	"reflect"
	"slices"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
)

// ValidateMySQLRuntimeLayout verifies the frozen MySQL physical proof against
// the canonical group without reading an editable destination or the network.
func ValidateMySQLRuntimeLayout(group *WriteGroup, frozen *WriteGroupRuntimeLayout) error {
	if group == nil || frozen == nil || frozen.Dialect != dbtarget.SQLDialectMySQL ||
		!group.isLongLayout() || group.StorageMode != WriteGroupStorageModePlain && group.StorageMode != WriteGroupStorageModeNativeMonthly ||
		group.Destination.StorageStrategy != WriteGroupStorageStrategyManaged ||
		group.Destination.ConnectorID == "" || group.Destination.ConnectorRevision == "" ||
		group.Destination.Database == "" || group.Destination.Database != group.Destination.TableSchema ||
		group.Destination.SchemaRevision == "" || len(group.Members) == 0 || frozen.Hypertable != nil ||
		group.WritePolicy.DedupeCapability != "receipt" || !frozen.ReceiptTableReady ||
		!LongRecordingShapeSupported(group, frozen.TagTypes) || !longRuntimeSupported(group) {
		return ErrWriteGroupValidation
	}
	proof := frozen.MySQL
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
	for name, physical := range map[string]*dbtarget.MySQLManagedInspection{group.Destination.TableName: proof.Group, dbtarget.EffectReceiptTable: proof.Receipt} {
		if physical.Inspection.Table != name || physical.Inspection.Schema != group.Destination.Database ||
			physical.Indexes == nil || !slices.EqualFunc(physical.Indexes, physical.Inspection.Indexes, func(a, b dbtarget.IndexInfo) bool {
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
	if (group.StorageMode == WriteGroupStorageModeNativeMonthly) != (evidence.Partitions != nil) || proof.Receipt.Evidence.Partitions != nil {
		return ErrWriteGroupValidation
	}
	if !dbtarget.MySQLRecordingVersionVerified(evidence.ServerFamily, evidence.ServerVersion) {
		return ErrWriteGroupValidation
	}
	layout, err := recordingplan.ManagedMySQLLongGroupLayout(table, owner, recordingplan.MySQLGroupSchemaOptions{
		ServerFamily: evidence.ServerFamily, ServerVersion: evidence.ServerVersion, CapabilityRevision: evidence.CapabilityRevision,
		Engine: evidence.Engine, Charset: evidence.Charset, Collation: evidence.Collation,
		NativeMonthly: recordingPlanMySQLPartitions(evidence.Partitions),
	})
	if err != nil {
		return fmt.Errorf("%w: frozen mysql layout: %w", ErrWriteGroupValidation, err)
	}
	digest, err := recordingplan.VerifyManagedMySQLCatalog(layout, group.Destination.Database,
		recordingPlanMySQLInspection(proof.Group), recordingPlanMySQLInspection(proof.Receipt))
	if err != nil || digest == "" || digest != group.Destination.SchemaDigest || digest != frozen.SchemaDigest ||
		!reflect.DeepEqual(frozen.Columns, proof.Group.Inspection.Columns) {
		return ErrWriteGroupValidation
	}
	return nil
}

func recordingPlanMySQLInspection(physical *dbtarget.MySQLManagedInspection) recordingplan.TargetTableInspection {
	if physical == nil || physical.Inspection == nil || physical.Evidence == nil {
		return recordingplan.TargetTableInspection{}
	}
	state := recordingPlanTableInspection(physical.Inspection)
	// The dedicated MySQL proof persists the complete index slice explicitly:
	// nil is unverified, while [] is the verified absence of secondary indexes.
	// The generic inspection's IndexesInspected flag is deliberately not JSON.
	state.IndexesInspected = physical.Indexes != nil
	state.Indexes = nil
	for _, index := range physical.Indexes {
		state.Indexes = append(state.Indexes, recordingplan.TargetIndexInspection{Name: index.Name, Columns: index.Columns, Unique: index.Unique})
	}
	e := physical.Evidence
	state.MySQL = &recordingplan.TargetMySQLTableEvidence{
		ServerFamily: e.ServerFamily, ServerVersion: e.ServerVersion, CapabilityRevision: e.CapabilityRevision, Database: e.Database,
		Engine: e.Engine, Charset: e.Charset, Collation: e.Collation, TableType: e.TableType, PrimaryKey: e.PrimaryKey,
		ColumnCollations: e.ColumnCollations, ColumnDefaults: e.ColumnDefaults, JSONColumns: e.JSONColumns,
		Partitions: recordingPlanMySQLPartitions(e.Partitions),
	}
	return state
}
