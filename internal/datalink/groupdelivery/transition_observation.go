package groupdelivery

import (
	"bytes"
	"context"
	"math"
	"time"

	"go-gateway/internal/datalink/runtimeview"
)

// DeliveryObservationConfig optionally observes confirmed transitions. Both
// callbacks are immutable after construction; Observe must be nonblocking.
// ScopeRevision is a pure read of current observed workspace ownership/revision.
type DeliveryObservationConfig struct {
	ScopeRevision func(context.Context, GroupKey) (string, bool)
	Observe       func(DeliveryObservation)
}

// DeliveryObservation is a closed, ephemeral reference to one native transition.
// WorkspaceRevision is observed around this transition, not historical output
// metadata. Unavailable observations carry no identity, raw code or payload.
type DeliveryObservation struct {
	Available         bool
	Key               GroupKey
	WorkspaceRevision string
	ConnectorID       string
	ConnectorRevision string
	EffectKey         string
	RecordID          string
	PayloadDigest     string
	State             string
	Reason            runtimeview.Reason
	CommittedAt       time.Time
	ReceiptMatched    bool
}

// WithDeliveryObservation returns a configured clone without changing s.
func (s *Store) WithDeliveryObservation(config DeliveryObservationConfig) *Store {
	clone := *s
	clone.deliveryObservation = config
	return &clone
}

type deliveryObservationBefore struct {
	item     OutboxItem
	owner    string
	revision string
	valid    bool
}

func (s *Store) observationScope(ctx context.Context, key GroupKey) (revision string, ok bool) {
	defer func() {
		if recover() != nil {
			revision, ok = "", false
		}
	}()
	if s.deliveryObservation.ScopeRevision == nil {
		return "", false
	}
	revision, ok = s.deliveryObservation.ScopeRevision(ctx, key)
	return revision, ok && revision != ""
}

func (s *Store) prepareDeliveryObservation(ctx context.Context, effect string) (before deliveryObservationBefore) {
	defer func() {
		if recover() != nil {
			before = deliveryObservationBefore{}
		}
	}()
	if s.deliveryObservation.Observe == nil {
		return deliveryObservationBefore{}
	}
	item, owner, err := s.readDeliveryObservationRow(ctx, effect)
	if err != nil {
		return deliveryObservationBefore{}
	}
	revision, ok := s.observationScope(ctx, item.Key)
	if !ok {
		return deliveryObservationBefore{}
	}
	check, checkedOwner, err := s.readDeliveryObservationRow(ctx, effect)
	return deliveryObservationBefore{item: item, owner: owner, revision: revision, valid: err == nil && owner == checkedOwner && sameDeliveryObservationRow(item, check)}
}

func (s *Store) observeDeliveryTransition(ctx context.Context, effect string, claim *Claim, state, code string, interrupted bool, run func() error) error {
	if s.deliveryObservation.Observe == nil {
		return run()
	}
	before := s.prepareDeliveryObservation(ctx, effect)
	err := run()
	if err == nil {
		s.finishDeliveryObservation(ctx, effect, claim, state, code, interrupted, before)
	}
	return err
}

func (s *Store) finishDeliveryObservation(ctx context.Context, effect string, claim *Claim, state, code string, interrupted bool, before deliveryObservationBefore) {
	if s.deliveryObservation.Observe == nil {
		return
	}
	observation := s.proveDeliveryObservation(ctx, effect, claim, state, code, interrupted, before)
	// All native SQL statements/transactions have returned before this callback.
	// A logging panic cannot replace the already-established mutation outcome.
	defer func() { recover() }() //nolint:errcheck // Logging cannot change the settled native outcome.
	s.deliveryObservation.Observe(observation)
}

func (s *Store) proveDeliveryObservation(ctx context.Context, effect string, claim *Claim, state, code string, interrupted bool, before deliveryObservationBefore) (observation DeliveryObservation) {
	unavailable := DeliveryObservation{Reason: runtimeview.ReasonUnavailable}
	defer func() {
		if recover() != nil {
			observation = unavailable
		}
	}()
	if !before.valid || !validDeliveryIdentity(before.item) || before.item.State == StateCommitted {
		return unavailable
	}
	after, owner, err := s.readDeliveryObservationRow(ctx, effect)
	if err != nil || !sameDeliveryIdentity(before.item, after) || after.State != state {
		return unavailable
	}
	if claim != nil && (before.owner != claim.Owner || owner != claim.Owner || before.item.State != StateSending || before.item.ClaimEpoch != claim.Epoch || after.ClaimEpoch != claim.Epoch || claim.EffectKey != effect) {
		return unavailable
	}
	if claim == nil {
		expectedEpoch := before.item.ClaimEpoch
		if interrupted {
			if before.item.State != StateSending || expectedEpoch == math.MaxInt64 {
				return unavailable
			}
			expectedEpoch++
		} else if before.item.State != StatePending && before.item.State != StateRetrying {
			return unavailable
		}
		expectedOwner := before.owner
		if interrupted {
			expectedOwner = ""
		}
		if after.ClaimEpoch != expectedEpoch || owner != expectedOwner {
			return unavailable
		}
	}
	matched := false
	if state == StateCommitted {
		if after.CommittedAt.IsZero() {
			return unavailable
		}
		matched, err = s.receiptMatches(ctx, sourceRow{state: after.State, effectKey: after.EffectKey, digest: after.PayloadDigest, committedAt: timestamp(after.CommittedAt)})
		if err != nil || !matched {
			return unavailable
		}
	} else if !after.CommittedAt.IsZero() || after.LastErrorCode != code {
		return unavailable
	}
	revision, ok := s.observationScope(ctx, after.Key)
	if !ok || revision != before.revision {
		return unavailable
	}
	check, checkedOwner, err := s.readDeliveryObservationRow(ctx, effect)
	if err != nil || checkedOwner != owner || !sameDeliveryObservationRow(after, check) {
		return unavailable
	}
	if state == StateCommitted {
		matched, err = s.receiptMatches(ctx, sourceRow{state: check.State, effectKey: check.EffectKey, digest: check.PayloadDigest, committedAt: timestamp(check.CommittedAt)})
		if err != nil || !matched {
			return unavailable
		}
	}
	reason := runtimeview.ParseReason(code)
	if state == StateCommitted {
		reason = runtimeview.ReasonNone
	} else if reason == runtimeview.ReasonNone {
		reason = runtimeview.ReasonUnavailable
	}
	return DeliveryObservation{Available: true, Key: after.Key, WorkspaceRevision: revision,
		ConnectorID: after.Destination.ConnectorID, ConnectorRevision: after.Destination.ConnectorRevision,
		EffectKey: after.EffectKey, RecordID: after.RecordID, PayloadDigest: after.PayloadDigest,
		State: state, Reason: reason, CommittedAt: after.CommittedAt, ReceiptMatched: matched}
}

func sameDeliveryIdentity(a, b OutboxItem) bool {
	return bytes.Equal(a.Payload, b.Payload) && a.Key == b.Key && a.EffectKey == b.EffectKey && a.RecordID == b.RecordID && a.PayloadDigest == b.PayloadDigest && a.Destination == b.Destination && a.EntityKey == b.EntityKey && a.PartitionKey == b.PartitionKey && a.BucketStart.Equal(b.BucketStart) && a.DueAt.Equal(b.DueAt)
}

func sameDeliveryObservationRow(a, b OutboxItem) bool {
	return sameDeliveryIdentity(a, b) && a.State == b.State && a.ClaimEpoch == b.ClaimEpoch && a.RetryCount == b.RetryCount && a.NextRetryAt.Equal(b.NextRetryAt) && a.LastErrorCode == b.LastErrorCode && a.CommittedAt.Equal(b.CommittedAt)
}

func validDeliveryIdentity(item OutboxItem) bool {
	if !item.Key.Valid() || item.EffectKey == "" || item.RecordID == "" || item.PayloadDigest == "" || item.Destination.ConnectorID == "" || item.Destination.ConnectorRevision == "" {
		return false
	}
	// Stable columns plus a matching receipt do not prove that the frozen
	// payload still has its original digest and encoded identities.
	_, _, _, valid := openFrozenRow(SourceScope{
		WorkspaceID: item.Key.WorkspaceID, GroupID: item.Key.GroupID, GroupRevision: item.Key.GroupRevision,
		ConnectorID: item.Destination.ConnectorID, ConnectorRevision: item.Destination.ConnectorRevision,
	}, sourceRow{
		workspaceID: item.Key.WorkspaceID, groupID: item.Key.GroupID, groupRevision: item.Key.GroupRevision,
		connectorID: item.Destination.ConnectorID, connectorRevision: item.Destination.ConnectorRevision,
		effectKey: item.EffectKey, recordID: item.RecordID, entityKey: item.EntityKey,
		payload: string(item.Payload), digest: item.PayloadDigest,
	})
	return valid
}

// Claim owner is diagnostic proof only, never part of the public carrier.
func (s *Store) readDeliveryObservationRow(ctx context.Context, effect string) (OutboxItem, string, error) {
	item, err := s.GetOutbox(ctx, effect)
	if err != nil {
		return OutboxItem{}, "", err
	}
	var owner string
	err = s.db.QueryRowContext(ctx, `SELECT claim_owner FROM wg_delivery_outbox WHERE effect_key=?`, effect).Scan(&owner)
	return item, owner, err
}
