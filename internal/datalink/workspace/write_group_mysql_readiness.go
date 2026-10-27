package workspace

import (
	"context"
	"strings"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
)

type managedMySQLRuntimeInspector interface {
	InspectMySQLManagedTableAtRevision(context.Context, string, string, string, string, string) (*dbtarget.MySQLManagedInspection, error)
}

func (s *WriteGroupService) mysqlReadiness(ctx context.Context, id string, content *WriteGroup, before *writeGroupReadinessSnapshot, result *WriteGroupReadiness) (*WriteGroupReadiness, error) {
	s.populateMySQLPartitionMaintenanceReadiness(ctx, before.validated, result)
	frozen, capability, inspectErr := s.readMySQLRuntimeLayoutWithCapability(ctx, before)
	after, err := s.readinessSnapshot(ctx, id, content)
	if err != nil {
		return nil, err
	}
	if !sameWriteGroupReadinessSnapshot(before, after) || after.validateErr != nil {
		result = s.newSnapshotReadiness(ctx, after)
		s.populateMySQLPartitionMaintenanceReadiness(ctx, after.validated, result)
		result.Issues = append(result.Issues, writeGroupReadinessIssue("readiness-snapshot-changed", "configuration changed during schema verification", id))
		return writeGroupReadinessResult(finalizeWriteGroupReadiness(result))
	}
	if inspectErr != nil {
		result.Issues = append(result.Issues, writeGroupReadinessIssue("managed-schema-unverified", "managed destination ownership and canonical schema could not be verified", id))
		return writeGroupReadinessResult(finalizeWriteGroupReadiness(result))
	}
	result.SchemaReady, result.SchemaDigest, result.runtimeLayout = true, frozen.SchemaDigest, frozen
	if !mysqlRecordingOperationGate(capability, before.validated) {
		message := "MySQL recording operations are not verified for this destination"
		if capability != nil && len(capability.Reasons) > 0 {
			message += ": " + strings.Join(capability.Reasons, ", ")
		}
		result.Issues = append(result.Issues, writeGroupReadinessIssue("mysql-recording-unverified", message, id))
	}
	return finalizeWriteGroupReadiness(result), nil
}

func mysqlRecordingOperationGate(capability *dbtarget.MySQLRecordingCapability, group *WriteGroup) bool {
	if capability == nil || group == nil || !capability.SupportsManagedLong {
		return false
	}
	if group.StorageMode == WriteGroupStorageModeNativeMonthly && !capability.SupportsNativePartition {
		return false
	}
	return true
}

func (s *WriteGroupService) readMySQLRuntimeLayout(ctx context.Context, snapshot *writeGroupReadinessSnapshot) (*WriteGroupRuntimeLayout, error) {
	frozen, _, err := s.readMySQLRuntimeLayoutWithCapability(ctx, snapshot)
	return frozen, err
}

func (s *WriteGroupService) readMySQLRuntimeLayoutWithCapability(ctx context.Context, snapshot *writeGroupReadinessSnapshot) (*WriteGroupRuntimeLayout, *dbtarget.MySQLRecordingCapability, error) {
	layout, capability, err := s.managedMySQLSchemaLayoutWithCapability(ctx, snapshot.validated)
	if err != nil {
		return nil, nil, err
	}
	reader, ok := s.managedTableInspector.(managedMySQLRuntimeInspector)
	if !ok {
		return nil, nil, ErrWriteGroupReadinessUnavailable
	}
	group := snapshot.validated
	read := func(table string) (*dbtarget.MySQLManagedInspection, error) {
		return reader.InspectMySQLManagedTableAtRevision(ctx, group.Destination.ConnectorID, group.Destination.ConnectorRevision,
			group.Destination.Database, layout.MySQL.CapabilityRevision, table)
	}
	physical, err := read(group.Destination.TableName)
	if err != nil {
		return nil, nil, err
	}
	receipt, err := read(dbtarget.EffectReceiptTable)
	if err != nil {
		return nil, nil, err
	}
	groupInspection := recordingPlanMySQLInspection(physical)
	groupPhysical := physical
	if group.StorageMode == WriteGroupStorageModeNativeMonthly && layout.MySQL != nil && layout.MySQL.NativeMonthly != nil &&
		physical != nil && physical.Inspection != nil && physical.Inspection.Status == dbtarget.TableInspectionExists {
		var err error
		groupInspection, err = recordingplan.NormalizeMySQLPartitionCatalog(ctx, s.repo.db, StorageScopeOf(group), *layout, groupInspection)
		if err != nil {
			return nil, nil, err
		}
		if physical.Evidence != nil {
			clonedEvidence := *physical.Evidence
			clonedEvidence.Partitions = dbtargetMySQLPartitions(layout.MySQL.NativeMonthly)
			clonedPhysical := *physical
			clonedPhysical.Evidence = &clonedEvidence
			groupPhysical = &clonedPhysical
		}
	}
	digest, err := recordingplan.VerifyManagedMySQLCatalog(*layout, group.Destination.Database,
		groupInspection, recordingPlanMySQLInspection(receipt))
	if err != nil || physical == nil || physical.Inspection == nil {
		return nil, nil, ErrWriteGroupValidation
	}
	frozen := &WriteGroupRuntimeLayout{
		Dialect: dbtarget.SQLDialectMySQL, Columns: physical.Inspection.Columns, TagTypes: snapshot.tagTypes,
		ReceiptTableReady: true, SchemaDigest: digest,
		MySQL: &dbtarget.MySQLManagedRowProof{ConnectorID: group.Destination.ConnectorID, ConnectorRevision: group.Destination.ConnectorRevision, Group: groupPhysical, Receipt: receipt},
	}
	if err := ValidateMySQLRuntimeLayout(group, frozen); err != nil {
		return nil, nil, err
	}
	return frozen, capability, nil
}
