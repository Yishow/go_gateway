package workspace

import (
	"context"
	"strings"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
)

// MySQLTestWriteLayout verifies the saved group's confirmed schema for an
// explicit test write. It does not apply the group or authorize runtime intake.
// Native proof targets only the ledger-authorized current physical partitions;
// the immutable recording confirmation and runtime descriptor stay unchanged.
func (s *WriteGroupService) MySQLTestWriteLayout(ctx context.Context, groupID, expectedGroupRevision string) (*WriteGroupRuntimeLayout, error) {
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
	if !ready.ConfigReady || !ready.SchemaReady || frozen == nil || frozen.Dialect != dbtarget.SQLDialectMySQL || frozen.MySQL == nil {
		return nil, ErrWriteGroupValidation
	}
	saved, err := s.Get(ctx, groupID)
	if err != nil {
		return nil, err
	}
	group := saved.Group
	if group.Revision != expectedGroupRevision {
		return nil, ErrWriteGroupRevisionConflict
	}
	if err := ValidateMySQLRuntimeLayout(group, frozen); err != nil {
		return nil, err
	}
	if group.StorageMode == WriteGroupStorageModeNativeMonthly {
		if err := s.mysqlTestWriteCurrentPartitions(ctx, group, frozen); err != nil {
			return nil, err
		}
	}
	return frozen, nil
}

func (s *WriteGroupService) mysqlTestWriteCurrentPartitions(ctx context.Context, group *WriteGroup, frozen *WriteGroupRuntimeLayout) error {
	scope := StorageScopeOf(group)
	lineage, err := recordingplan.ReadStorageLineage(ctx, s.repo.db, scope)
	if err != nil {
		return err
	}
	if !lineage.Confirmed() || lineage.ConnectorRevision != group.Destination.ConnectorRevision ||
		lineage.ConfirmedDigest != group.Destination.SchemaDigest {
		return ErrWriteGroupValidation
	}
	state, err := recordingplan.ReadMySQLPartitionMaintenanceState(ctx, s.repo.db, scope)
	if err != nil {
		return err
	}
	actual := recordingPlanMySQLInspection(frozen.MySQL.Group)
	if state != nil {
		actual.MySQL.Partitions = state.Current
	}
	normalized, err := recordingplan.NormalizeMySQLPartitionCatalog(ctx, s.repo.db, scope, *lineage.Layout, actual)
	if err != nil {
		return err
	}
	digest, err := recordingplan.VerifyManagedMySQLCatalog(*lineage.Layout, group.Destination.Database,
		normalized, recordingPlanMySQLInspection(frozen.MySQL.Receipt))
	if err != nil || digest != frozen.SchemaDigest {
		return ErrWriteGroupValidation
	}
	// These values belong to this fresh readiness read. Clone the proof nodes
	// before replacing partitions so the normalized runtime proof is not edited.
	proof := *frozen.MySQL
	physical := *proof.Group
	evidence := *physical.Evidence
	evidence.Partitions = dbtargetMySQLPartitions(actual.MySQL.Partitions)
	physical.Evidence = &evidence
	proof.Group = &physical
	frozen.MySQL = &proof
	return nil
}
