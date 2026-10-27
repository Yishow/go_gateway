package grouppipeline

import (
	"context"
	"strings"

	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/workspace"
)

// loadFinalizedSnapshot is used only by delivery/storage, never by intake or
// replay. Finalized effects already contain their values and original identity.
func (p *Pipeline) loadFinalizedSnapshot(ctx context.Context, key groupdelivery.GroupKey) (*workspace.WriteGroupAppliedSnapshot, error) {
	return loadFinalizedStorageSnapshot(ctx, p.deps.Store, key)
}

func loadFinalizedStorageSnapshot(ctx context.Context, store *groupdelivery.Store, key groupdelivery.GroupKey) (*workspace.WriteGroupAppliedSnapshot, error) {
	payload, err := store.RuntimeVersion(ctx, key)
	if err != nil || payload == nil {
		return nil, err
	}
	projected, err := counterproof.SnapshotStorageProjection(payload)
	if err != nil {
		return nil, errRuntimeSnapshotInvalid
	}
	return decodeFrozenStorageSnapshot(projected, key)
}

// finalizedDestinationMatches retains the applied storage identity even when
// counter meaning cannot be replayed. Outbox metadata cannot redirect a row.
func finalizedDestinationMatches(snap *workspace.WriteGroupAppliedSnapshot, item groupdelivery.OutboxItem) bool {
	frozen, actual := snap.Group.Destination, item.Destination
	scope, err := snapshot.DestinationScope(frozen.ConnectorID, frozen.ConnectorRevision, frozen.Database, frozen.TableSchema, frozen.TableName)
	dedupe := strings.ToLower(strings.TrimSpace(snap.Group.WritePolicy.DedupeCapability))
	if dedupe == "" {
		dedupe = groupdelivery.DedupeNone
	}
	recordKey := ""
	if snap.RuntimeLayout.Dialect == dbtarget.SQLDialectMySQL || snap.RuntimeLayout.Dialect == dbtarget.SQLDialectSQLServer {
		recordKey = snap.Group.RowPolicy.RecordKeyColumn
	}
	return err == nil && actual.Scope == scope && frozen.ConnectorID == actual.ConnectorID && frozen.ConnectorRevision == actual.ConnectorRevision && frozen.Database == actual.Database && frozen.TableSchema == actual.TableSchema && frozen.TableName == actual.TableName && actual.DedupeCapability == dedupe && actual.RecordKeyColumn == recordKey
}
