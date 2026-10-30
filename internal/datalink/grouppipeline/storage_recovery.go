package grouppipeline

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

// StorageRecovery restores the confirmed shared table of a managed long group
// when it vanished from the destination, and only then. It adapts the delivery
// sender to the restore engine of the operation ledger: it gathers the local
// evidence frozen when the group revision was applied (policy, the group's own
// confirmation, ownership), hands over the destination delivery already
// opened, and turns the engine's answer into what the sender does with the row.
//
// It never opens a connection of its own to restore anything. The destination
// is the one delivery opened for the frozen connector revision, which cannot
// create a missing SQLite file, so a vanished database file is an outage, not
// permission to build an empty replacement.
type StorageRecovery struct {
	store        *groupdelivery.Store
	ledger       *recordingplan.Service
	destinations Destinations
	report       func(error)
	sqlServerOps *sqlServerStorageRestoreTargetOps
	observation  StorageObservationConfig

	snapshots sync.Map // groupdelivery.GroupKey -> *workspace.WriteGroupAppliedSnapshot (immutable per revision)

	mu    sync.Mutex
	locks map[string]*scopeLock
}

type scopeLock struct {
	mu   sync.Mutex
	refs int
}

// NewStorageRecovery builds the recovery over the local delivery store and the
// operation ledger, which share the gateway's local database.
func NewStorageRecovery(store *groupdelivery.Store, ledger *recordingplan.Service, destinations Destinations) *StorageRecovery {
	return &StorageRecovery{store: store, ledger: ledger, destinations: destinations, locks: map[string]*scopeLock{}}
}

var _ groupdelivery.StorageRecovery = (*StorageRecovery)(nil)

// snapshot returns the verified frozen snapshot of the item's group revision,
// or nil when the item does not belong to managed long storage. Counter replay
// proof is independent of finalized delivery. Invalid storage facts and read
// failures remain refusals, never permission to bypass lineage or ownership.
func (r *StorageRecovery) snapshot(ctx context.Context, key groupdelivery.GroupKey) (*workspace.WriteGroupAppliedSnapshot, error) {
	if cached, ok := r.snapshots.Load(key); ok {
		if snap, ok := cached.(*workspace.WriteGroupAppliedSnapshot); ok {
			return snap, nil
		}
		return nil, nil
	}
	snap, err := loadFinalizedStorageSnapshot(ctx, r.store, key)
	if err != nil {
		return nil, err
	}
	if snap == nil || !longGroup(snap.Group) || snap.Group.Destination.StorageStrategy != workspace.WriteGroupStorageStrategyManaged {
		r.snapshots.Store(key, (*workspace.WriteGroupAppliedSnapshot)(nil))
		return nil, nil
	}
	r.snapshots.Store(key, snap)
	return snap, nil
}

// Admit blocks delivery to a scope whose last automatic restore ended in a
// result nobody can vouch for, until an explicit confirmation checks it again.
func (r *StorageRecovery) Admit(ctx context.Context, item groupdelivery.OutboxItem) (string, error) {
	snap, err := r.snapshot(ctx, item.Key)
	if errors.Is(err, errRuntimeSnapshotInvalid) {
		return recordingplan.StorageCodeOwnerMismatch, nil
	}
	if err != nil {
		return "", err
	}
	if snap == nil {
		return "", nil
	}
	if !finalizedDestinationMatches(snap, item) {
		return recordingplan.StorageCodeOwnerMismatch, nil
	}
	lineage, err := recordingplan.ReadStorageLineage(ctx, r.store.DB(), workspace.StorageScopeOf(snap.Group))
	if errors.Is(err, recordingplan.ErrStorageLineageNotFound) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	if lineage.RepairState != "" {
		return recordingplan.StorageCodeRepairBlocked, nil
	}
	return "", nil
}

// BlockedCode names the stable code of managed long storage whose confirmed
// destination identity can no longer be opened.
func (r *StorageRecovery) BlockedCode(ctx context.Context, item groupdelivery.OutboxItem) string {
	snap, err := r.snapshot(ctx, item.Key)
	if err != nil {
		return recordingplan.StorageCodeRepairBlocked
	}
	if snap == nil {
		return ""
	}
	if _, err := recordingplan.ReadStorageLineage(ctx, r.store.DB(), workspace.StorageScopeOf(snap.Group)); err != nil {
		return ""
	}
	return recordingplan.StorageCodeRepairBlocked
}

// Restore decides what to do with a row the destination refused with a
// schema-shaped error. The error only starts the look; the catalog and the
// local evidence decide.
func (r *StorageRecovery) Restore(ctx context.Context, item groupdelivery.OutboxItem, target groupdelivery.Target, cause error) (decision groupdelivery.RestoreDecision, err error) {
	ctx, original, finish := r.captureStorageObservation(ctx)
	defer func() { finish(decision.Action == groupdelivery.RestoreBlock || err != nil) }()
	snap, err := r.snapshot(ctx, item.Key)
	if errors.Is(err, errRuntimeSnapshotInvalid) {
		return groupdelivery.RestoreDecision{Action: groupdelivery.RestoreBlock, Code: recordingplan.StorageCodeOwnerMismatch}, nil
	}
	if err != nil {
		return groupdelivery.RestoreDecision{}, err
	}
	if snap == nil {
		return groupdelivery.RestoreDecision{}, nil
	}
	if original != nil {
		*original = originalStorageProof(item, snap)
	}
	if target.Kind == schema.DatabaseConnectorKindSQLServer {
		return r.restoreSQLServer(ctx, item, target, cause, snap)
	}
	group := snap.Group
	dialect, ok := dialectName(target.Kind)
	if !ok {
		return groupdelivery.RestoreDecision{}, nil
	}
	table, owner, ok := workspace.LongStorageIdentity(group)
	if !ok {
		return groupdelivery.RestoreDecision{}, nil
	}
	scope := workspace.StorageScopeOf(group)
	var layout *recordingplan.GroupSchemaLayout
	var lineage *recordingplan.StorageLineage
	if target.Kind == schema.DatabaseConnectorKindMySQL {
		if target.DB == nil || !mysqlFrozenScopeMatches(item, snap) ||
			!mysqlCapabilityMatchesDestination(target.MySQLCapability, item.Destination) ||
			!mysqlCapabilityMatchesProof(target.MySQLCapability, snap.RuntimeLayout.MySQL) {
			return groupdelivery.RestoreDecision{Action: groupdelivery.RestoreBlock, Code: recordingplan.StorageCodeRepairBlocked}, nil
		}
		layout, lineage, err = r.mysqlConfirmedFallbackLayout(ctx, scope, table, owner)
		if err != nil {
			return groupdelivery.RestoreDecision{}, err
		}
	} else {
		resolved, layoutErr := workspace.LongStorageLayout(group, dialect)
		if layoutErr != nil {
			// The group's own mode cannot exist on this kind of destination.
			return groupdelivery.RestoreDecision{Action: groupdelivery.RestoreBlock, Code: recordingplan.StorageCodeOwnerMismatch}, nil
		}
		layout = resolved
	}
	denial := r.denial(ctx, snap, item, table, dialect)
	if lineage != nil && lineage.Confirmed() && lineage.ConfirmedDigest != group.Destination.SchemaDigest {
		denial = &recordingplan.StorageDenial{Code: recordingplan.StorageCodeOwnerMismatch, Reason: recordingplan.StorageReasonConfirmationMissing}
	}
	request := recordingplan.StorageRestoreRequest{
		Scope: scope, Dialect: dialect, ConnectorRevision: item.Destination.ConnectorRevision,
		TableName: table, OwnerColumn: owner, GroupID: group.ID, FallbackLayout: layout,
		Deny: denial,
	}
	release := r.lock(request.Scope)
	defer release()
	result, err := r.ledger.RestoreManagedStorage(ctx, request, restoreTarget(target, item.Destination.TableSchema, group.StorageMode == workspace.WriteGroupStorageModeTimescale))
	if err != nil {
		return groupdelivery.RestoreDecision{}, err
	}
	switch result.Outcome {
	case recordingplan.StorageRestored:
		return groupdelivery.RestoreDecision{Action: groupdelivery.RestoreRetryNow}, nil
	case recordingplan.StorageNotMissing:
		if dbtarget.IsMissingRelationError(cause) {
			// The table is there now (restored by another execution, or by a
			// person): the insert that found it gone is simply tried again.
			return groupdelivery.RestoreDecision{Action: groupdelivery.RestoreRetryNow}, nil
		}
		return groupdelivery.RestoreDecision{}, nil
	case recordingplan.StorageInProgress:
		return groupdelivery.RestoreDecision{Action: groupdelivery.RestoreWait}, nil
	default:
		return groupdelivery.RestoreDecision{Action: groupdelivery.RestoreBlock, Code: result.Code}, nil
	}
}

// denial turns what the frozen group says about itself into a reason not to
// restore its table. It matters only once the catalog shows the table missing.
func (r *StorageRecovery) denial(ctx context.Context, snap *workspace.WriteGroupAppliedSnapshot, item groupdelivery.OutboxItem, table, dialect string) *recordingplan.StorageDenial {
	group, frozen := snap.Group, snap.Group.Destination
	mismatch := &recordingplan.StorageDenial{Code: recordingplan.StorageCodeOwnerMismatch, Reason: recordingplan.StorageReasonOwnershipMismatch}
	switch {
	case frozen.ConnectorID != item.Destination.ConnectorID || frozen.ConnectorRevision != item.Destination.ConnectorRevision ||
		frozen.Database != item.Destination.Database || frozen.TableSchema != item.Destination.TableSchema ||
		frozen.TableName != item.Destination.TableName || frozen.TableName != table:
		return mismatch
	case snap.RuntimeLayout == nil || string(snap.RuntimeLayout.Dialect) != dialect:
		return mismatch
	case group.WritePolicy.RecreateMissing == nil:
		return &recordingplan.StorageDenial{Code: recordingplan.StorageCodeRepairBlocked, Reason: recordingplan.StorageReasonPolicyUnconfirmed}
	case !*group.WritePolicy.RecreateMissing:
		return &recordingplan.StorageDenial{Code: recordingplan.StorageCodeRepairBlocked, Reason: recordingplan.StorageReasonRecoveryDisabled}
	}
	if !r.groupConfirmed(ctx, snap) {
		return &recordingplan.StorageDenial{Code: recordingplan.StorageCodeOwnerMismatch, Reason: recordingplan.StorageReasonConfirmationMissing}
	}
	return nil
}

// groupConfirmed checks the group's own confirmation against the operation
// ledger: it names a succeeded operation that verified exactly that digest.
func (r *StorageRecovery) groupConfirmed(ctx context.Context, snap *workspace.WriteGroupAppliedSnapshot) bool {
	frozen := snap.Group.Destination
	if frozen.SchemaRevision == "" || len(frozen.SchemaDigest) != 64 {
		return false
	}
	operation, err := r.ledger.GetSchemaOperation(ctx, snap.Group.WorkspaceID, frozen.SchemaRevision)
	return err == nil && operation.Status == recordingplan.SchemaOperationSucceeded && operation.VerifiedDigest == frozen.SchemaDigest
}

// lock serializes restores of one storage scope within this process, so
// parallel groups sharing a table take turns: the first restores it, the others
// then find it there. Other processes meet in the operation ledger.
func (r *StorageRecovery) lock(scope recordingplan.StorageScope) (release func()) {
	key := fmt.Sprintf("%s\x1f%s\x1f%s\x1f%s\x1f%s\x1f%s", scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema, scope.StorageMode, scope.LayoutVersion)
	r.mu.Lock()
	entry, ok := r.locks[key]
	if !ok {
		entry = &scopeLock{}
		r.locks[key] = entry
	}
	entry.refs++
	r.mu.Unlock()
	entry.mu.Lock()
	return func() {
		entry.mu.Unlock()
		r.mu.Lock()
		entry.refs--
		if entry.refs == 0 {
			delete(r.locks, key)
		}
		r.mu.Unlock()
	}
}

func dialectName(kind schema.DatabaseConnectorKind) (string, bool) {
	switch kind {
	case schema.DatabaseConnectorKindSQLite:
		return "sqlite", true
	case schema.DatabaseConnectorKindPostgres:
		return "postgres", true
	case schema.DatabaseConnectorKindMySQL:
		return "mysql", true
	case schema.DatabaseConnectorKindSQLServer:
		return string(dbtarget.SQLDialectSQLServer), true
	}
	return "", false
}

// StorageView is the API-safe storage truth of one managed long group: which
// generation of the shared table it writes to and whether a rebuilt table means
// earlier history may be gone. It carries no values, DSNs or driver text.
type StorageView struct {
	Generation          int                          `json:"generation"`
	HistoryMayBeLost    bool                         `json:"history_may_be_lost"`
	Reconstructions     int                          `json:"reconstructions"`
	LastReconstructedAt *time.Time                   `json:"last_reconstructed_at,omitempty"`
	RepairState         string                       `json:"repair_state,omitempty"`
	RepairReason        string                       `json:"repair_reason,omitempty"`
	RecreateMissing     *bool                        `json:"recreate_missing,omitempty"`
	Events              []recordingplan.StorageEvent `json:"events"`
}

// View reads the storage diagnostics of a managed long group.
func (r *StorageRecovery) View(ctx context.Context, group *workspace.WriteGroup) (*StorageView, error) {
	if group == nil || !longGroup(group) || group.Destination.StorageStrategy != workspace.WriteGroupStorageStrategyManaged {
		return nil, nil
	}
	scope := workspace.StorageScopeOf(group)
	lineage, err := recordingplan.ReadStorageLineage(ctx, r.store.DB(), scope)
	if errors.Is(err, recordingplan.ErrStorageLineageNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	events, err := recordingplan.ReadStorageEvents(ctx, r.store.DB(), scope, group.ID)
	if err != nil {
		return nil, err
	}
	view := &StorageView{
		Generation: lineage.Generation, HistoryMayBeLost: lineage.HistoryMayBeLost, Reconstructions: lineage.Reconstructions,
		RepairState: lineage.RepairState, RepairReason: lineage.RepairReason, RecreateMissing: group.WritePolicy.RecreateMissing,
		Events: events,
	}
	if !lineage.LastReconstructedAt.IsZero() {
		at := lineage.LastReconstructedAt
		view.LastReconstructedAt = &at
	}
	return view, nil
}
