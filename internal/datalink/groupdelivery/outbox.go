package groupdelivery

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// OutboxItem is one closed row waiting for, or finished with, delivery.
type OutboxItem struct {
	EffectKey     string
	RecordID      string
	Key           GroupKey
	EntityKey     string
	BucketStart   time.Time
	PartitionKey  string
	Destination   FrozenDestination
	Payload       []byte
	PayloadDigest string
	State         string
	ClaimEpoch    int64
	RetryCount    int
	NextRetryAt   time.Time
	LastErrorCode string
	CommittedAt   time.Time
}

// ErrNotDeliverable means the item is not in a state that may be sent.
var ErrNotDeliverable = errors.New("outbox item is not deliverable")

// ErrOutboxItemNotFound means no outbox row has the effect key.
var ErrOutboxItemNotFound = errors.New("outbox item not found")

const outboxColumns = `effect_key, record_id, workspace_id, group_id, group_revision, entity_key, bucket_start,
	partition_key, destination_scope, connector_id, connector_revision, database_name, table_schema, table_name,
	dedupe_capability, record_key_column, payload, payload_digest, state, retry_count, next_retry_at,
	last_error_code, committed_at, claim_epoch`

type rowScanner interface{ Scan(dest ...any) error }

func scanOutbox(row rowScanner) (OutboxItem, error) {
	var (
		item                         OutboxItem
		bucket, nextRetry, committed string
		payload                      string
	)
	err := row.Scan(
		&item.EffectKey, &item.RecordID, &item.Key.WorkspaceID, &item.Key.GroupID, &item.Key.GroupRevision,
		&item.EntityKey, &bucket, &item.PartitionKey, &item.Destination.Scope, &item.Destination.ConnectorID,
		&item.Destination.ConnectorRevision, &item.Destination.Database, &item.Destination.TableSchema,
		&item.Destination.TableName, &item.Destination.DedupeCapability, &item.Destination.RecordKeyColumn,
		&payload, &item.PayloadDigest, &item.State, &item.RetryCount, &nextRetry, &item.LastErrorCode, &committed, &item.ClaimEpoch,
	)
	if err != nil {
		return OutboxItem{}, err
	}
	item.Payload = []byte(payload)
	var parseErr error
	if item.BucketStart, parseErr = time.Parse(time.RFC3339Nano, bucket); parseErr != nil {
		return OutboxItem{}, fmt.Errorf("parse outbox bucket_start: %w", parseErr)
	}
	if item.NextRetryAt, parseErr = time.Parse(time.RFC3339Nano, nextRetry); parseErr != nil {
		return OutboxItem{}, fmt.Errorf("parse outbox next_retry_at: %w", parseErr)
	}
	if committed != "" {
		if item.CommittedAt, parseErr = time.Parse(time.RFC3339Nano, committed); parseErr != nil {
			return OutboxItem{}, fmt.Errorf("parse outbox committed_at: %w", parseErr)
		}
	}
	return item, nil
}

// GetOutbox reads one outbox item by effect key.
func (s *Store) GetOutbox(ctx context.Context, effectKey string) (OutboxItem, error) {
	item, err := scanOutbox(s.db.QueryRowContext(ctx,
		`SELECT `+outboxColumns+` FROM wg_delivery_outbox WHERE effect_key = ?`, effectKey))
	if errors.Is(err, sql.ErrNoRows) {
		return OutboxItem{}, ErrOutboxItemNotFound
	}
	if err != nil {
		return OutboxItem{}, fmt.Errorf("read outbox item: %w", err)
	}
	return item, nil
}

// transition moves an item between states in one statement. It reports
// ErrNotDeliverable when the item exists but is not in one of the allowed
// source states, and ErrOutboxItemNotFound when it does not exist.
func (s *Store) transition(ctx context.Context, effectKey string, from []string, set string, args ...any) error {
	params := append([]any{}, args...)
	params = append(params, timestamp(s.now()), effectKey)
	for _, state := range from {
		params = append(params, state)
	}
	result, err := s.db.ExecContext(ctx, transitionStatement(set, len(from)), params...)
	if err != nil {
		return fmt.Errorf("update outbox state: %w", err)
	}
	if affected, err := result.RowsAffected(); err == nil && affected > 0 {
		return nil
	}
	var state string
	err = s.db.QueryRowContext(ctx, `SELECT state FROM wg_delivery_outbox WHERE effect_key = ?`, effectKey).Scan(&state)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrOutboxItemNotFound
	}
	if err != nil {
		return fmt.Errorf("read outbox state: %w", err)
	}
	return ErrNotDeliverable
}

// transitionStatement builds the guarded state update. set holds only
// fixed column assignments chosen by this package; all values are parameters.
func transitionStatement(set string, fromCount int) string {
	marks := strings.TrimSuffix(strings.Repeat("?, ", fromCount), ", ")
	return "UPDATE wg_delivery_outbox SET " + set + ", updated_at = ? WHERE effect_key = ? AND state IN (" + marks + ")"
}

// Claim is the exclusive, fenced right to deliver one item. Every outcome
// write must present it; a holder whose claim was superseded is fenced out.
type Claim struct {
	EffectKey string
	Owner     string
	Epoch     int64
}

// ErrFenced means the item is no longer held by the caller's claim.
var ErrFenced = errors.New("delivery claim was superseded")

// BeginDelivery marks a pending or retrying item as being sent by owner and
// returns the claim. Only one caller can win; each claim gets a strictly newer
// fencing epoch and a lease of ttl, so a worker that stalls past its lease can
// be superseded and can never write an outcome afterwards.
func (s *Store) BeginDelivery(ctx context.Context, effectKey, owner string, ttl time.Duration) (Claim, error) {
	now := s.now()
	var epoch int64
	err := s.db.QueryRowContext(ctx, `
		UPDATE wg_delivery_outbox
		SET state = 'sending', claim_owner = ?, claim_expires_at = ?, claim_epoch = claim_epoch + 1, updated_at = ?
		WHERE effect_key = ? AND state IN ('pending', 'retrying')
		RETURNING claim_epoch`, owner, timestamp(now.Add(ttl)), timestamp(now), effectKey).Scan(&epoch)
	if errors.Is(err, sql.ErrNoRows) {
		return Claim{}, s.missingOrNotDeliverable(ctx, effectKey)
	}
	if err != nil {
		return Claim{}, fmt.Errorf("claim outbox item: %w", err)
	}
	return Claim{EffectKey: effectKey, Owner: owner, Epoch: epoch}, nil
}

func (s *Store) missingOrNotDeliverable(ctx context.Context, effectKey string) error {
	var state string
	err := s.db.QueryRowContext(ctx, `SELECT state FROM wg_delivery_outbox WHERE effect_key = ?`, effectKey).Scan(&state)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrOutboxItemNotFound
	}
	if err != nil {
		return fmt.Errorf("read outbox state: %w", err)
	}
	return ErrNotDeliverable
}

// transitionClaimed moves a claimed item out of sending, but only for the
// holder of the current claim. A superseded holder gets ErrFenced and changes
// nothing, whatever state the item has reached since.
func (s *Store) transitionClaimed(ctx context.Context, claim Claim, set string, args ...any) error {
	params := append([]any{}, args...)
	params = append(params, timestamp(s.now()), claim.EffectKey, claim.Owner, claim.Epoch)
	result, err := s.db.ExecContext(ctx, claimedStatement(set), params...)
	if err != nil {
		return fmt.Errorf("update claimed outbox state: %w", err)
	}
	if affected, err := result.RowsAffected(); err == nil && affected > 0 {
		return nil
	}
	return s.fencedOrMissing(ctx, claim)
}

func claimedStatement(set string) string {
	return "UPDATE wg_delivery_outbox SET " + set + ", updated_at = ? WHERE effect_key = ? AND state = 'sending' AND claim_owner = ? AND claim_epoch = ?"
}

func (s *Store) fencedOrMissing(ctx context.Context, claim Claim) error {
	var owner string
	var epoch int64
	err := s.db.QueryRowContext(ctx, `SELECT claim_owner, claim_epoch FROM wg_delivery_outbox WHERE effect_key = ?`, claim.EffectKey).Scan(&owner, &epoch)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrOutboxItemNotFound
	}
	if err != nil {
		return fmt.Errorf("read outbox claim: %w", err)
	}
	if owner != claim.Owner || epoch != claim.Epoch {
		return ErrFenced
	}
	return ErrNotDeliverable
}

// CompleteDelivery records destination-confirmed success: the item becomes
// sql_committed and a local receipt is written in the same transaction.
func (s *Store) CompleteDelivery(ctx context.Context, claim Claim, payloadDigest string) error {
	err := s.inTx(ctx, "delivery receipt", func(tx *sql.Tx) error {
		now := timestamp(s.now())
		result, err := tx.ExecContext(ctx, `
			UPDATE wg_delivery_outbox
			SET state = ?, committed_at = ?, last_error_code = '', updated_at = ?
			WHERE effect_key = ? AND state = 'sending' AND claim_owner = ? AND claim_epoch = ?`,
			StateCommitted, now, now, claim.EffectKey, claim.Owner, claim.Epoch)
		if err != nil {
			return fmt.Errorf("mark outbox committed: %w", err)
		}
		if affected, err := result.RowsAffected(); err != nil || affected == 0 {
			return errClaimLost
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO wg_delivery_receipts (effect_key, payload_digest, committed_at) VALUES (?, ?, ?)
			ON CONFLICT (effect_key) DO NOTHING`, claim.EffectKey, payloadDigest, now); err != nil {
			return fmt.Errorf("write local receipt: %w", err)
		}
		return nil
	})
	if errors.Is(err, errClaimLost) {
		return s.fencedOrMissing(ctx, claim)
	}
	if err == nil {
		s.usage.invalidate()
	}
	return err
}

// errClaimLost is internal: the guarded update matched nothing.
var errClaimLost = errors.New("delivery claim update matched nothing")

// MarkRetry records a failed attempt that is known not to have committed.
// state is retrying (with the next due time) or blocked once retries are
// exhausted; the data is always kept.
func (s *Store) MarkRetry(ctx context.Context, claim Claim, state, code string, next time.Time) error {
	if state != StateRetrying && state != StateBlocked {
		return fmt.Errorf("%w: retry state %q", ErrNotDeliverable, state)
	}
	return s.transitionClaimed(ctx, claim, "state = ?, retry_count = retry_count + 1, next_retry_at = ?, last_error_code = ?",
		state, timestamp(next), code)
}

// MarkUnknown records an attempt whose destination outcome cannot be known.
func (s *Store) MarkUnknown(ctx context.Context, claim Claim, code string) error {
	return s.transitionClaimed(ctx, claim, "state = 'unknown', last_error_code = ?", code)
}

// MarkBlocked stops an unclaimed item until an operator repairs its cause.
func (s *Store) MarkBlocked(ctx context.Context, effectKey, code string) error {
	return s.transition(ctx, effectKey, []string{StatePending, StateRetrying}, "state = 'blocked', last_error_code = ?", code)
}

// BlockClaimed stops a claimed item until an operator repairs its cause.
func (s *Store) BlockClaimed(ctx context.Context, claim Claim, code string) error {
	return s.transitionClaimed(ctx, claim, "state = 'blocked', last_error_code = ?", code)
}

// ResolveInterrupted settles an item left in sending by a crash or a failed
// local update. A destination that can recognize a repeat (receipt or unique
// key) makes a retry safe; otherwise the earlier attempt may have committed
// and the item becomes unknown instead of being resent.
func (s *Store) ResolveInterrupted(ctx context.Context, effectKey string) (string, error) {
	item, err := s.GetOutbox(ctx, effectKey)
	if err != nil {
		return "", err
	}
	if item.State != StateSending {
		return item.State, nil
	}
	if item.Destination.DedupeCapability == DedupeReceipt || item.Destination.DedupeCapability == DedupeUniqueKey {
		err = s.transition(ctx, effectKey, []string{StateSending}, "state = 'retrying', next_retry_at = ?, last_error_code = 'interrupted-delivery', claim_owner = '', claim_epoch = claim_epoch + 1", timestamp(s.now()))
		if err != nil && !errors.Is(err, ErrNotDeliverable) {
			return "", err
		}
		return StateRetrying, nil
	}
	err = s.transition(ctx, effectKey, []string{StateSending}, "state = 'unknown', last_error_code = 'interrupted-delivery', claim_owner = '', claim_epoch = claim_epoch + 1")
	if err != nil && !errors.Is(err, ErrNotDeliverable) {
		return "", err
	}
	return StateUnknown, nil
}

// QuarantineResolution is an explicit operator decision about a stuck row.
type QuarantineResolution string

const (
	// ResolutionRetry sends the row again after the cause was repaired.
	ResolutionRetry QuarantineResolution = "retry"
	// ResolutionSkip moves the partition past the row; the payload is kept.
	ResolutionSkip QuarantineResolution = "skip"
)

// MarkQuarantined isolates a row the destination rejects; its payload is kept
// and later rows of the same partition wait behind it.
func (s *Store) MarkQuarantined(ctx context.Context, claim Claim, code string) error {
	return s.transitionClaimed(ctx, claim, "state = 'quarantined', last_error_code = ?", code)
}

// ResolveQuarantine applies an operator decision to a quarantined or blocked
// row. Retry gives the row a fresh retry budget; skip keeps the payload on
// record but lets the partition move past it.
func (s *Store) ResolveQuarantine(ctx context.Context, effectKey string, resolution QuarantineResolution) error {
	from := []string{StateQuarantined, StateBlocked}
	switch resolution {
	case ResolutionRetry:
		return s.transition(ctx, effectKey, from,
			"state = 'pending', retry_count = 0, next_retry_at = ?, last_error_code = ''", timestamp(s.now()))
	case ResolutionSkip:
		return s.transition(ctx, effectKey, from, "state = 'operator_skipped'")
	}
	return fmt.Errorf("%w: unknown resolution %q", ErrNotDeliverable, resolution)
}

// Head identifies the row a partition may deliver next.
type Head struct {
	PartitionKey string
	EffectKey    string
}

// headQuery selects, per partition, the earliest row that is not yet finished
// (committed or operator-skipped) when that row is pending or retrying and due.
// A row in any other state (sending, unknown, blocked, quarantined, retrying
// but not due) therefore holds back everything behind it in its partition,
// while other partitions are unaffected.
const headQuery = `
	SELECT o.partition_key, o.effect_key FROM wg_delivery_outbox o
	WHERE o.state IN ('pending', 'retrying') AND o.next_retry_at <= ?%s
	  AND NOT EXISTS (
		SELECT 1 FROM wg_delivery_outbox e
		WHERE e.partition_key = o.partition_key
		  AND (e.bucket_start < o.bucket_start OR (e.bucket_start = o.bucket_start AND e.effect_key < o.effect_key))
		  AND e.state NOT IN ('sql_committed', 'operator_skipped'))
	ORDER BY o.bucket_start, o.effect_key`

// headStatement builds the head query, optionally restricted to one partition.
func headStatement(onePartition bool) string {
	if onePartition {
		return fmt.Sprintf(headQuery, " AND o.partition_key = ?") + " LIMIT 1"
	}
	return fmt.Sprintf(headQuery, "")
}

// ReadyHeads returns the deliverable head of every partition that has one.
func (s *Store) ReadyHeads(ctx context.Context, now time.Time) ([]Head, error) {
	rows, err := s.db.QueryContext(ctx, headStatement(false), timestamp(now))
	if err != nil {
		return nil, fmt.Errorf("read ready partition heads: %w", err)
	}
	defer rows.Close()
	var heads []Head
	for rows.Next() {
		var head Head
		if err := rows.Scan(&head.PartitionKey, &head.EffectKey); err != nil {
			return nil, fmt.Errorf("scan partition head: %w", err)
		}
		heads = append(heads, head)
	}
	return heads, rows.Err()
}

// PartitionHead returns the deliverable head of one partition, if any.
func (s *Store) PartitionHead(ctx context.Context, partitionKey string, now time.Time) (effectKey string, found bool, err error) {
	var head Head
	err = s.db.QueryRowContext(ctx, headStatement(true), timestamp(now), partitionKey).Scan(&head.PartitionKey, &head.EffectKey)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("read partition head: %w", err)
	}
	return head.EffectKey, true, nil
}

// RecoverStaleClaims settles rows left in sending by a worker that is gone:
// claims past their lease, and claims of earlier incarnations of this node
// (owner "<node>/<other incarnation>"), which are provably dead after a
// restart. Each recovered row gets a new fencing epoch so the old worker can no
// longer write outcomes. Without dedupe a recovered row becomes unknown, never
// a blind retry. It returns how many rows were settled.
func (s *Store) RecoverStaleClaims(ctx context.Context, nodeID, selfOwner string, now time.Time) (int, error) {
	prefix := nodeID + "/"
	rows, err := s.db.QueryContext(ctx, `
		SELECT effect_key FROM wg_delivery_outbox
		WHERE state = 'sending'
		  AND (claim_expires_at <= ? OR (? != '/' AND substr(claim_owner, 1, ?) = ? AND claim_owner != ?))`,
		timestamp(now), prefix, len(prefix), prefix, selfOwner)
	if err != nil {
		return 0, fmt.Errorf("find stale claims: %w", err)
	}
	var keys []string
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			_ = rows.Close()
			return 0, fmt.Errorf("scan stale claim: %w", err)
		}
		keys = append(keys, key)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return 0, err
	}
	if err := rows.Close(); err != nil {
		return 0, err
	}
	recovered := 0
	for _, key := range keys {
		state, err := s.ResolveInterrupted(ctx, key)
		if err != nil {
			return recovered, err
		}
		if state == StateRetrying || state == StateUnknown {
			recovered++
		}
	}
	return recovered, nil
}
