package grouppipeline

import (
	"context"
	"errors"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	datalinksnapshot "go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/workspace"
)

const sqlServerMetadataVerified = "verified"

func (r destinationResolver) resolveSQLServer(ctx context.Context, item groupdelivery.OutboxItem, opened *dbtarget.OpenedDestination) (groupdelivery.Target, error) {
	if opened.DB == nil || r.loadRuntimeSnapshot == nil || !sqlServerCapabilityMatchesDestination(opened.SQLServer, item.Destination) {
		return rejectOpenedDestination(opened, "SQL Server destination does not match a verified frozen recording scope")
	}
	frozen, err := r.loadRuntimeSnapshot(ctx, item.Key)
	if err != nil || !sqlServerFrozenScopeMatches(item, frozen) || !sqlServerCapabilityMatchesProof(opened.SQLServer, frozen.RuntimeLayout.SQLServer) {
		return rejectOpenedDestination(opened, "SQL Server frozen runtime proof is unavailable or no longer matches the destination")
	}
	proof := dbtarget.CloneSQLServerManagedRowProof(frozen.RuntimeLayout.SQLServer)
	if r.loadSQLServerProof != nil {
		proof, err = r.loadSQLServerProof(ctx, frozen, opened)
		if err != nil || proof == nil {
			if errors.Is(err, errSQLServerProofReadTransient) {
				if opened.Close != nil {
					err = errors.Join(err, opened.Close())
				}
				return groupdelivery.Target{}, err
			}
			return rejectOpenedDestination(opened, "SQL Server physical generation is not authorized by the confirmed repair ledger")
		}
	}
	return groupdelivery.Target{
		DB: opened.DB, Kind: opened.Kind, SQLServer: dbtarget.CloneSQLServerManagedRowProof(proof),
		SQLServerCapability: opened.SQLServer, Close: opened.Close,
	}, nil
}

func sqlServerCapabilityMatchesDestination(capability *dbtarget.SQLServerRecordingCapability, destination groupdelivery.FrozenDestination) bool {
	if capability == nil || !capability.Connected || capability.MetadataStatus != sqlServerMetadataVerified || !capability.SupportsManagedLong ||
		capability.ConnectorID != destination.ConnectorID || capability.ConnectorRevision != destination.ConnectorRevision ||
		capability.Database != destination.Database || capability.Schema != destination.TableSchema {
		return false
	}
	p := capability.Permissions
	return p.Verified && (p.DatabaseViewDefinition || p.SchemaViewDefinition) && p.CreateTable && p.AlterSchema && p.SelectSchema && p.InsertSchema && p.DeleteSchema
}

func sqlServerFrozenScopeMatches(item groupdelivery.OutboxItem, frozen *workspace.WriteGroupAppliedSnapshot) bool {
	if frozen == nil || frozen.Group == nil || frozen.RuntimeLayout == nil || frozen.EffectiveAt.IsZero() ||
		frozen.WorkspaceID != item.Key.WorkspaceID || frozen.GroupID != item.Key.GroupID ||
		frozen.GroupRevision != item.Key.GroupRevision || frozen.AppliedRevision != item.Key.GroupRevision ||
		frozen.Group.WorkspaceID != item.Key.WorkspaceID || frozen.Group.ID != item.Key.GroupID ||
		frozen.Group.Revision != item.Key.GroupRevision || frozen.Group.AppliedRevision != item.Key.GroupRevision ||
		!longGroup(frozen.Group) || !validLongRuntimeProof(frozen.Group, frozen.RuntimeLayout) ||
		frozen.RuntimeLayout.Dialect != dbtarget.SQLDialectSQLServer {
		return false
	}
	destination := frozen.Group.Destination
	outbox := item.Destination
	scope, err := datalinksnapshot.DestinationScope(destination.ConnectorID, destination.ConnectorRevision, destination.Database, destination.TableSchema, destination.TableName)
	return err == nil && scope == outbox.Scope && destination.ConnectorID == outbox.ConnectorID &&
		destination.ConnectorRevision == outbox.ConnectorRevision && destination.Database == outbox.Database &&
		destination.TableSchema == outbox.TableSchema && destination.TableName == outbox.TableName &&
		frozen.Group.WritePolicy.DedupeCapability == outbox.DedupeCapability && frozen.Group.RowPolicy.RecordKeyColumn == outbox.RecordKeyColumn
}

func sqlServerCapabilityMatchesProof(capability *dbtarget.SQLServerRecordingCapability, proof *dbtarget.SQLServerManagedRowProof) bool {
	if capability == nil || proof == nil || proof.Group == nil || proof.Group.Evidence == nil || proof.Receipt == nil || proof.Receipt.Evidence == nil ||
		capability.ConnectorID != proof.ConnectorID || capability.ConnectorRevision != proof.ConnectorRevision {
		return false
	}
	native := proof.Group.Evidence.NativeMonthly != nil
	for _, evidence := range []*dbtarget.SQLServerManagedTableEvidence{proof.Group.Evidence, proof.Receipt.Evidence} {
		if capability.Database != evidence.Database || capability.Schema != evidence.Schema ||
			capability.ServerVersion != evidence.ServerVersion || capability.ServerEdition != evidence.ServerEdition ||
			capability.EngineEdition != evidence.EngineEdition || capability.CompatibilityLevel != evidence.CompatibilityLevel {
			return false
		}
		if capability.CapabilityRevision != evidence.CapabilityRevision &&
			(!native || !dbtarget.SQLServerNativeDeliveryCapabilityRevisionMatches(capability, evidence.CapabilityRevision)) {
			return false
		}
	}
	return true
}
