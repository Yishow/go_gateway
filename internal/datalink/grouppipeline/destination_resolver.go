package grouppipeline

import (
	"context"
	"errors"
	"fmt"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/schema"
	datalinksnapshot "go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/workspace"
)

// Resolve opens the outbox row's frozen destination; MySQL additionally needs
// the applied runtime proof that was frozen for the same group revision. A
// native_monthly group resolves only when the canonical workspace validator
// accepts its frozen proof.
func (r destinationResolver) Resolve(ctx context.Context, item groupdelivery.OutboxItem) (groupdelivery.Target, error) {
	if r.destinations == nil {
		return groupdelivery.Target{}, fmt.Errorf("%w: destination service unavailable", groupdelivery.ErrTargetBlocked)
	}
	opened, err := r.destinations.OpenDestination(ctx, item.Destination.ConnectorID, item.Destination.ConnectorRevision)
	if err != nil {
		if errors.Is(err, dbtarget.ErrDestinationBlocked) {
			return groupdelivery.Target{}, fmt.Errorf("%w: %w", groupdelivery.ErrTargetBlocked, err)
		}
		return groupdelivery.Target{}, err
	}
	if opened == nil {
		return groupdelivery.Target{}, fmt.Errorf("%w: destination opener returned no connection", groupdelivery.ErrTargetBlocked)
	}
	if opened.Kind == schema.DatabaseConnectorKindSQLServer {
		return r.resolveSQLServer(ctx, item, opened)
	}
	if opened.Kind != schema.DatabaseConnectorKindMySQL {
		return groupdelivery.Target{DB: opened.DB, Kind: opened.Kind, Close: opened.Close}, nil
	}
	return r.resolveMySQL(ctx, item, opened)
}

func (r destinationResolver) resolveMySQL(
	ctx context.Context,
	item groupdelivery.OutboxItem,
	opened *dbtarget.OpenedDestination,
) (groupdelivery.Target, error) {
	if opened.DB == nil || r.loadRuntimeSnapshot == nil || !mysqlCapabilityMatchesDestination(opened.MySQL, item.Destination) {
		return rejectOpenedDestination(opened, "MySQL destination does not match a verified frozen recording scope")
	}
	frozen, err := r.loadRuntimeSnapshot(ctx, item.Key)
	if err != nil || !mysqlFrozenScopeMatches(item, frozen) ||
		!mysqlCapabilityMatchesProof(opened.MySQL, frozen.RuntimeLayout.MySQL) {
		return rejectOpenedDestination(opened, "MySQL frozen runtime proof is unavailable or no longer matches the destination")
	}
	proof := frozen.RuntimeLayout.MySQL
	if r.loadMySQLProof != nil {
		proof, err = r.loadMySQLProof(ctx, frozen)
		if err != nil || proof == nil {
			if errors.Is(err, dbtarget.ErrMySQLManagedPartitionChanged) || errors.Is(err, errMySQLPartitionProofReadTransient) {
				if opened.Close != nil {
					err = errors.Join(err, opened.Close())
				}
				return groupdelivery.Target{}, err
			}
			return rejectOpenedDestination(opened, "MySQL authorized partition proof is unavailable")
		}
	}
	return groupdelivery.Target{
		DB: opened.DB, Kind: opened.Kind, MySQL: proof,
		MySQLCapability: opened.MySQL, Close: opened.Close,
	}, nil
}

func mysqlCapabilityMatchesDestination(capability *dbtarget.MySQLRecordingCapability, destination groupdelivery.FrozenDestination) bool {
	return capability != nil && capability.Connected && capability.MetadataStatus == "verified" && capability.TransactionalEngine &&
		capability.ConnectorID == destination.ConnectorID && capability.ConnectorRevision == destination.ConnectorRevision &&
		capability.Database == destination.Database
}

func mysqlFrozenScopeMatches(item groupdelivery.OutboxItem, frozen *workspace.WriteGroupAppliedSnapshot) bool {
	if frozen == nil || frozen.Group == nil || frozen.RuntimeLayout == nil || frozen.EffectiveAt.IsZero() ||
		frozen.WorkspaceID != item.Key.WorkspaceID || frozen.GroupID != item.Key.GroupID ||
		frozen.GroupRevision != item.Key.GroupRevision || frozen.AppliedRevision != item.Key.GroupRevision ||
		frozen.Group.WorkspaceID != item.Key.WorkspaceID || frozen.Group.ID != item.Key.GroupID ||
		frozen.Group.Revision != item.Key.GroupRevision || frozen.Group.AppliedRevision != item.Key.GroupRevision ||
		!longGroup(frozen.Group) ||
		frozen.Group.StorageMode != workspace.WriteGroupStorageModePlain && frozen.Group.StorageMode != workspace.WriteGroupStorageModeNativeMonthly ||
		frozen.Group.Destination.StorageStrategy != workspace.WriteGroupStorageStrategyManaged ||
		frozen.Group.WritePolicy.DedupeCapability != groupdelivery.DedupeReceipt ||
		frozen.RuntimeLayout.Dialect != dbtarget.SQLDialectMySQL || frozen.RuntimeLayout.MySQL == nil ||
		workspace.ValidateMySQLRuntimeLayout(frozen.Group, frozen.RuntimeLayout) != nil {
		return false
	}
	destination := frozen.Group.Destination
	itemDestination := item.Destination
	scope, err := datalinksnapshot.DestinationScope(
		destination.ConnectorID, destination.ConnectorRevision, destination.Database, destination.TableSchema, destination.TableName,
	)
	if err != nil || scope != itemDestination.Scope || destination.ConnectorID != itemDestination.ConnectorID ||
		destination.ConnectorRevision != itemDestination.ConnectorRevision || destination.Database != itemDestination.Database ||
		destination.TableSchema != itemDestination.TableSchema || destination.TableName != itemDestination.TableName ||
		frozen.Group.WritePolicy.DedupeCapability != itemDestination.DedupeCapability ||
		frozen.Group.RowPolicy.RecordKeyColumn != itemDestination.RecordKeyColumn {
		return false
	}
	proof := frozen.RuntimeLayout.MySQL
	return proof.ConnectorID == destination.ConnectorID && proof.ConnectorRevision == destination.ConnectorRevision &&
		proof.Group != nil && proof.Group.Inspection != nil && proof.Group.Evidence != nil &&
		proof.Receipt != nil && proof.Receipt.Inspection != nil && proof.Receipt.Evidence != nil
}

func mysqlCapabilityMatchesProof(
	capability *dbtarget.MySQLRecordingCapability,
	proof *dbtarget.MySQLManagedRowProof,
) bool {
	if capability == nil || proof == nil || proof.Group == nil || proof.Group.Evidence == nil ||
		proof.Receipt == nil || proof.Receipt.Evidence == nil {
		return false
	}
	group, receipt := proof.Group.Evidence, proof.Receipt.Evidence
	return capability.ConnectorID == proof.ConnectorID && capability.ConnectorRevision == proof.ConnectorRevision &&
		capability.Database == group.Database && capability.ServerFamily == group.ServerFamily &&
		capability.ServerVersion == group.ServerVersion && capability.CapabilityRevision == group.CapabilityRevision &&
		group.Database == receipt.Database && group.ServerFamily == receipt.ServerFamily &&
		group.ServerVersion == receipt.ServerVersion && group.CapabilityRevision == receipt.CapabilityRevision
}

func rejectOpenedDestination(opened *dbtarget.OpenedDestination, reason string) (groupdelivery.Target, error) {
	blocked := fmt.Errorf("%w: %s", groupdelivery.ErrTargetBlocked, reason)
	if opened != nil && opened.Close != nil {
		if err := opened.Close(); err != nil {
			blocked = errors.Join(blocked, fmt.Errorf("close rejected destination: %w", err))
		}
	}
	return groupdelivery.Target{}, blocked
}
