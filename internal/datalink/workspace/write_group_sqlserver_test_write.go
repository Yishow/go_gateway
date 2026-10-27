package workspace

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
)

const sqlServerTestWriteMetadataVerified = "verified"

type managedSQLServerTestWriteInspector interface {
	InspectSQLServerManagedTableAtRevision(context.Context, string, string, string, string, string, string) (*dbtarget.SQLServerManagedInspection, error)
}

// SQLServerTestWriteLayout verifies the saved SQL Server group proof for one
// explicit test write. It never applies the group or mutates its frozen digest.
func (s *WriteGroupService) SQLServerTestWriteLayout(ctx context.Context, groupID, expectedGroupRevision string) (*WriteGroupRuntimeLayout, error) {
	if strings.TrimSpace(expectedGroupRevision) == "" {
		return nil, ErrWriteGroupRevisionConflict
	}
	ready, err := s.Readiness(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if ready.GroupRevision != expectedGroupRevision {
		return nil, ErrWriteGroupRevisionConflict
	}
	frozen := ready.runtimeLayout
	if !ready.ConfigReady || !ready.SchemaReady || frozen == nil || frozen.Dialect != dbtarget.SQLDialectSQLServer || frozen.SQLServer == nil {
		return nil, ErrWriteGroupValidation
	}
	saved, err := s.Get(ctx, groupID)
	if err != nil {
		return nil, err
	}
	group := saved.Group
	if group == nil || group.Revision != expectedGroupRevision {
		return nil, ErrWriteGroupRevisionConflict
	}
	if err := ValidateSQLServerRuntimeLayout(group, frozen); err != nil {
		return nil, err
	}
	prober, ok := s.managedTableInspector.(ManagedWriteGroupSQLServerProbeInspector)
	if !ok {
		return nil, ErrWriteGroupValidation
	}
	capability, err := prober.ProbeSQLServerRecordingAtRevision(ctx, group.Destination.ConnectorID, group.Destination.ConnectorRevision)
	if err != nil || !sqlServerTestWriteCapabilityMatches(group, frozen, capability) {
		return nil, ErrWriteGroupValidation
	}
	if group.StorageMode == WriteGroupStorageModeNativeMonthly {
		return s.sqlServerTestWriteCurrentPartitions(ctx, ready, group, frozen)
	}
	return frozen, nil
}

func sqlServerTestWriteCapabilityMatches(group *WriteGroup, frozen *WriteGroupRuntimeLayout, capability *dbtarget.SQLServerRecordingCapability) bool {
	if group == nil || frozen == nil || frozen.SQLServer == nil || frozen.SQLServer.Group == nil ||
		frozen.SQLServer.Group.Evidence == nil || capability == nil || !capability.Connected ||
		capability.MetadataStatus != sqlServerTestWriteMetadataVerified || !capability.SupportsManagedLong || !capability.SupportsGroupTestWrite ||
		capability.ConnectorID != group.Destination.ConnectorID || capability.ConnectorRevision != group.Destination.ConnectorRevision ||
		capability.Database != group.Destination.Database || capability.Schema != group.Destination.TableSchema {
		return false
	}
	evidence := frozen.SQLServer.Group.Evidence
	return capability.CapabilityRevision == evidence.CapabilityRevision && capability.ServerVersion == evidence.ServerVersion &&
		capability.ServerEdition == evidence.ServerEdition && capability.EngineEdition == evidence.EngineEdition &&
		capability.CompatibilityLevel == evidence.CompatibilityLevel &&
		(group.StorageMode != WriteGroupStorageModeNativeMonthly || capability.SupportsNativePartition)
}

func (s *WriteGroupService) sqlServerTestWriteCurrentPartitions(
	ctx context.Context,
	ready *WriteGroupReadiness,
	group *WriteGroup,
	frozen *WriteGroupRuntimeLayout,
) (*WriteGroupRuntimeLayout, error) {
	if s == nil || s.repo == nil || frozen == nil || frozen.SQLServer == nil || frozen.SQLServer.Group == nil ||
		frozen.SQLServer.Group.Evidence == nil || frozen.SQLServer.Group.Evidence.CapabilityRevision == "" {
		return nil, ErrWriteGroupValidation
	}
	mutation := WriteGroupMutation{
		WorkspaceID: ready.WorkspaceID, ExpectedWorkspaceRevision: ready.WorkspaceRevision,
		ExpectedGroupRevision: group.Revision, ExpectedConnectorRevision: group.Destination.ConnectorRevision,
	}
	scope, err := s.ManagedSchemaScope(ctx, group.ID, mutation)
	if err != nil || scope.GroupLayout == nil || scope.SchemaDigest != frozen.SchemaDigest ||
		scope.ConnectorRevision != group.Destination.ConnectorRevision {
		return nil, ErrWriteGroupValidation
	}
	inspector, ok := s.managedTableInspector.(managedSQLServerTestWriteInspector)
	if !ok {
		return nil, ErrWriteGroupValidation
	}
	inspect := func(table string) (*dbtarget.SQLServerManagedInspection, error) {
		return inspector.InspectSQLServerManagedTableAtRevision(ctx,
			group.Destination.ConnectorID, group.Destination.ConnectorRevision,
			group.Destination.Database, group.Destination.TableSchema,
			frozen.SQLServer.Group.Evidence.CapabilityRevision, table)
	}
	physical, err := inspect(group.Destination.TableName)
	if err != nil || physical == nil || physical.Inspection == nil || physical.Evidence == nil ||
		physical.Inspection.Status != dbtarget.TableInspectionExists {
		return nil, ErrWriteGroupValidation
	}
	receipt, err := inspect(dbtarget.EffectReceiptTable)
	if err != nil || receipt == nil || receipt.Inspection == nil || receipt.Evidence == nil ||
		receipt.Inspection.Status != dbtarget.TableInspectionExists {
		return nil, ErrWriteGroupValidation
	}
	actual := recordingPlanSQLServerInspection(physical)
	actualReceipt := recordingPlanSQLServerInspection(receipt)
	normalized, err := s.NormalizeSQLServerPartitionInspection(ctx, scope, actual, actualReceipt)
	if err != nil {
		return nil, fmt.Errorf("verify SQL Server current partition catalog: %w", err)
	}
	verifiedDigest, err := recordingplan.VerifyManagedSQLServerCatalog(*scope.GroupLayout, scope.Database, normalized, actualReceipt)
	if err != nil || verifiedDigest != frozen.SchemaDigest {
		return nil, ErrWriteGroupValidation
	}
	if physical.Evidence.NativeMonthly == nil || !physical.Evidence.PartitionInspected {
		return nil, ErrWriteGroupValidation
	}

	// The physical proof carries the authorized current partition boundaries to
	// the bounded insert path. The schema digest remains the immutable confirmed
	// digest, whose native layout is the original partition prefix.
	layout := *frozen
	proof := dbtarget.CloneSQLServerManagedRowProof(frozen.SQLServer)
	groupProof := *proof.Group
	groupEvidence := *groupProof.Evidence
	groupEvidence.NativeMonthly = cloneSQLServerTestWritePartition(physical.Evidence.NativeMonthly)
	groupEvidence.PartitionInspected = physical.Evidence.PartitionInspected
	groupProof.Evidence = &groupEvidence
	proof.Group = &groupProof
	layout.SQLServer = proof
	return &layout, nil
}

func cloneSQLServerTestWritePartition(evidence *dbtarget.SQLServerManagedPartitionEvidence) *dbtarget.SQLServerManagedPartitionEvidence {
	if evidence == nil {
		return nil
	}
	cloned := *evidence
	cloned.Partitions = slices.Clone(evidence.Partitions)
	return &cloned
}
