package grouppipeline

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"
	"reflect"
	"strings"
	"time"
)

// StorageObservationOutcome describes a proven schema-operation result, never
// recovered historical data or SQL delivery of an outbox effect.
type StorageObservationOutcome string

const (
	StorageObservationCompleted   StorageObservationOutcome = "completed"
	StorageObservationBlocked     StorageObservationOutcome = "blocked"
	StorageObservationUnavailable StorageObservationOutcome = "unavailable"
)

// StorageOperationObservation contains only original frozen identity and the
// real settled ledger ID. No new operation or observation ID is generated.
type StorageOperationObservation struct {
	Available                                                      bool
	Key                                                            groupdelivery.GroupKey
	WorkspaceRevision, ConnectorID, ConnectorRevision, OperationID string
	Outcome                                                        StorageObservationOutcome
	Reason                                                         string
}

// StorageObservationConfig is immutable after construction. Observer callbacks
// must be nonblocking and cannot authorize repair or modify native outcomes.
type StorageObservationConfig struct {
	RecordLimit              int
	OnDrop                   func(uint64)
	InfoEnabled, WarnEnabled bool
	ScopeRevision            func(context.Context, groupdelivery.GroupKey) (string, bool)
	Observe                  func(StorageOperationObservation)
}

// NewStorageRecoveryWithObservation installs optional capture before use. It
// never copies the recovery's mutex, maps, or an already-running instance.
func NewStorageRecoveryWithObservation(store *groupdelivery.Store, ledger *recordingplan.Service, destinations Destinations, cfg StorageObservationConfig) *StorageRecovery {
	r := NewStorageRecovery(store, ledger, destinations)
	r.observation = cfg
	return r
}

// Only scalar identity and fixed-size hashes cross the native operation. The
// original verified input is not a current draft or a newly reconstructed ID.
type storageOriginalProof struct {
	available                            bool
	key                                  groupdelivery.GroupKey
	effect, connector, connectorRevision string
	frozenHash, itemHash                 [32]byte
}

func originalStorageProof(item groupdelivery.OutboxItem, snap *workspace.WriteGroupAppliedSnapshot) (out storageOriginalProof) {
	defer func() {
		if recover() != nil {
			out = storageOriginalProof{}
		}
	}()
	if snap == nil || snap.Group == nil || snap.RuntimeLayout == nil || !finalizedDestinationMatches(snap, item) {
		return out
	}
	for _, id := range []string{item.Key.WorkspaceID, item.Key.GroupID, item.Key.GroupRevision, item.EffectKey, item.RecordID, item.Destination.ConnectorID, item.Destination.ConnectorRevision} {
		if !safeOwnerID(id) {
			return out
		}
	}
	itemHash, ok := storageItemHash(item)
	if !ok {
		return out
	}
	frozenHash, ok := storageProofHash(snap)
	if !ok {
		return out
	}
	return storageOriginalProof{available: true, key: item.Key, effect: strings.Clone(item.EffectKey), connector: strings.Clone(item.Destination.ConnectorID), connectorRevision: strings.Clone(item.Destination.ConnectorRevision), frozenHash: frozenHash, itemHash: itemHash}
}

func storageProofHash(value any) ([32]byte, bool) {
	data, err := json.Marshal(value)
	if err != nil {
		return [32]byte{}, false
	}
	return sha256.Sum256(data), true
}

// Mutable delivery state is deliberately excluded: schema recreation proves
// neither receipt commitment nor the current send/retry claim.
func storageItemHash(item groupdelivery.OutboxItem) ([32]byte, bool) {
	sum := sha256.Sum256(item.Payload)
	if hex.EncodeToString(sum[:]) != item.PayloadDigest {
		return [32]byte{}, false
	}
	row, err := groupdelivery.DecodeRowPayload(item.Payload)
	if err != nil || row.EffectKey != item.EffectKey || row.RecordID != item.RecordID || row.EntityKey != item.EntityKey || !row.BucketStart.Equal(item.BucketStart) {
		return [32]byte{}, false
	}
	item.Payload = nil
	item.State, item.LastErrorCode = "", ""
	item.ClaimEpoch, item.RetryCount = 0, 0
	item.NextRetryAt, item.CommittedAt = time.Time{}, time.Time{}
	return storageProofHash(item)
}

// Defer finish before native lock/close defers. The single existing A capture
// owns descriptor admission and loss; callbacks run after every outer lock.
func (r *StorageRecovery) captureStorageObservation(ctx context.Context) (context.Context, *storageOriginalProof, func(bool)) {
	cfg := r.observation
	if cfg.Observe == nil || !cfg.InfoEnabled && !cfg.WarnEnabled {
		return ctx, nil, func(bool) {}
	}
	original := &storageOriginalProof{}
	seen := false
	captured, flush := recordingplan.CaptureStorageSettlements(ctx, recordingplan.StorageObservationConfig{
		RecordLimit: cfg.RecordLimit, InfoEnabled: cfg.InfoEnabled, WarnEnabled: cfg.WarnEnabled,
		OnDrop: func(n uint64) {
			seen = true
			if cfg.OnDrop != nil {
				cfg.OnDrop(n)
			}
		},
		Observe: func(s recordingplan.StorageSettlement) {
			seen = true
			out := r.verifyStorageObservation(ctx, *original, s)
			if !out.Available && !cfg.WarnEnabled {
				return
			}
			cfg.Observe(out)
		},
	})
	return captured, original, func(refused bool) {
		flush()
		// Before a claim there is no real ledger ID. Capture only unavailable,
		// once, and never duplicate a captured or admission-dropped settlement.
		if !seen && refused && cfg.WarnEnabled {
			if cfg.RecordLimit <= 0 {
				storageObservationDrop(cfg)
				return
			}
			func() {
				defer func() {
					if recover() != nil {
						storageObservationDrop(cfg)
					}
				}()
				cfg.Observe(StorageOperationObservation{})
			}()
		}
	}
}
func storageObservationDrop(cfg StorageObservationConfig) {
	defer func() { _ = recover() }()
	if cfg.OnDrop != nil {
		cfg.OnDrop(1)
	}
}

// This brackets observation-time metadata only. It cannot assert historical
// workspace configuration, restore data, authorize DDL, or prove SQL delivery.
func (r *StorageRecovery) verifyStorageObservation(ctx context.Context, p storageOriginalProof, s recordingplan.StorageSettlement) (out StorageOperationObservation) {
	defer func() {
		if recover() != nil {
			out = StorageOperationObservation{}
		}
	}()
	cfg := r.observation
	if !p.available || !s.Available || !s.OwnedExecution || r.store == nil || r.ledger == nil || cfg.ScopeRevision == nil || ctx.Err() != nil {
		return out
	}
	if s.WorkspaceID != p.key.WorkspaceID || !safeOwnerID(s.OperationID) || s.Status != recordingplan.SchemaOperationSucceeded && s.Status != recordingplan.SchemaOperationFailed {
		return out
	}
	revision, ok := cfg.ScopeRevision(ctx, p.key)
	if !ok || !safeOwnerID(revision) {
		return out
	}
	before, ok := r.storageSettledProof(ctx, p, s)
	if !ok {
		return out
	}
	afterRevision, ok := cfg.ScopeRevision(ctx, p.key)
	if !ok || afterRevision != revision {
		return out
	}
	// Re-read after the second provider callback; it may have raced with a
	// frozen row, lineage, or operation replacement. Volatile outbox counters
	// are not compared and no transaction-snapshot claim is made.
	after, ok := r.storageSettledProof(ctx, p, s)
	if !ok || !reflect.DeepEqual(before, after) {
		return out
	}
	out = StorageOperationObservation{Available: true, Key: p.key, WorkspaceRevision: revision, ConnectorID: p.connector, ConnectorRevision: p.connectorRevision, OperationID: s.OperationID, Outcome: StorageObservationCompleted, Reason: "none"}
	if s.Status == recordingplan.SchemaOperationFailed {
		out.Outcome = StorageObservationBlocked
		out.Reason = storageFailureReason(before.operation.Reason)
	}
	return out
}

type storageSettledEvidence struct {
	operation *recordingplan.SchemaOperation
	lineage   *recordingplan.StorageLineage
}

func (r *StorageRecovery) storageSettledProof(ctx context.Context, p storageOriginalProof, s recordingplan.StorageSettlement) (storageSettledEvidence, bool) {
	var empty storageSettledEvidence
	snap, err := loadFinalizedStorageSnapshot(ctx, r.store, p.key)
	if err != nil || snap == nil || snap.Group == nil || snap.RuntimeLayout == nil {
		return empty, false
	}
	hash, ok := storageProofHash(snap)
	if !ok || hash != p.frozenHash {
		return empty, false
	}
	item, err := r.store.GetOutbox(ctx, p.effect)
	if err != nil || !finalizedDestinationMatches(snap, item) {
		return empty, false
	}
	hash, ok = storageItemHash(item)
	if !ok || hash != p.itemHash {
		return empty, false
	}
	op, err := r.ledger.GetSchemaOperation(ctx, p.key.WorkspaceID, s.OperationID)
	if err != nil || op == nil || op.OperationID != s.OperationID || op.WorkspaceID != s.WorkspaceID || op.Action != recordingplan.StorageRepairAction || op.Status != s.Status || op.VerifiedDigest != s.VerifiedDigest || !op.UpdatedAt.Equal(s.UpdatedAt) || op.CompletedAt == nil || !op.CompletedAt.Equal(s.CompletedAt) {
		return empty, false
	}
	scope := workspace.StorageScopeOf(snap.Group)
	lineage, err := recordingplan.ReadStorageLineage(ctx, r.store.DB(), scope)
	frozen := snap.Group.Destination
	if err != nil || !lineage.Confirmed() || lineage.ConnectorRevision != p.connectorRevision || lineage.ConfirmedOperationID != frozen.SchemaRevision || lineage.ConfirmedDigest != frozen.SchemaDigest {
		return empty, false
	}
	table, owner, ok := workspace.LongStorageIdentity(snap.Group)
	if !ok || lineage.TableName != table || lineage.OwnerColumn != owner || lineage.Layout.TableName != table || lineage.Layout.OwnerColumn != owner {
		return empty, false
	}
	key := (&recordingplan.SchemaPreviewToken{WorkspaceID: scope.WorkspaceID, ConnectorID: scope.ConnectorID, Database: scope.Database, Schema: scope.TableSchema, GroupLayout: lineage.Layout}).ScopeKey()
	if key == "" || op.ScopeKey != key {
		return empty, false
	}
	if s.Status == recordingplan.SchemaOperationSucceeded {
		if op.VerifiedDigest == "" || op.VerifiedDigest != lineage.ConfirmedDigest || lineage.LastReconstructionOperationID != op.OperationID || !lineage.HistoryMayBeLost || lineage.Reconstructions <= 0 || lineage.LastReconstructedAt.IsZero() {
			return empty, false
		}
	} else if storageFailureReason(op.Reason) == "" {
		return empty, false
	}
	return storageSettledEvidence{operation: op, lineage: lineage}, true
}
func storageFailureReason(reason string) string {
	switch reason {
	case recordingplan.StorageReasonPermissionDenied, recordingplan.StorageReasonStatementFailed:
		return "unavailable"
	case recordingplan.StorageReasonConfirmationMissing, recordingplan.StorageReasonOwnershipMismatch, recordingplan.StorageReasonRepairNeedsConfirm, recordingplan.StorageReasonConnectorRevision, recordingplan.StorageReasonTableOccupied, recordingplan.StorageReasonRecoveryDisabled, recordingplan.StorageReasonPolicyUnconfirmed, recordingplan.StorageReasonReceiptsMissing, recordingplan.StorageReasonTimescaleUnavailable, recordingplan.StorageReasonTimescaleMoved, recordingplan.StorageReasonVerificationMismatch:
		return "effective-proof-missing"
	}
	return ""
}
