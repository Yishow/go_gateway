package groupdelivery

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"go-gateway/internal/datalink/snapshot"
)

type memberRecord struct {
	Member     string `json:"member"`
	Status     string `json:"status"`
	Reason     string `json:"reason,omitempty"`
	SampleID   string `json:"sample_id,omitempty"`
	Quality    string `json:"quality,omitempty"`
	ObservedAt string `json:"observed_at,omitempty"`
}

// bucketReason is the reason stored beside a closed bucket. A written row has
// none of its own, so a row a cutover cut short carries that fact instead.
func bucketReason(outcome snapshot.Outcome) string {
	if outcome.Reason == "" && outcome.Kind == snapshot.OutcomeRow && outcome.Truncated != "" {
		return snapshot.ReasonBoundaryTruncated + ":" + outcome.Truncated
	}
	return outcome.Reason
}

func memberRecords(outcome snapshot.Outcome) string {
	records := make([]memberRecord, 0, len(outcome.Members))
	for _, member := range outcome.Members {
		record := memberRecord{Member: member.MemberKey, Status: string(member.Status), Reason: member.Reason}
		if member.Sample != nil {
			record.SampleID = member.Sample.SampleID
			record.Quality = string(member.Sample.Quality)
			record.ObservedAt = timestamp(member.Sample.ObservedAt)
		}
		records = append(records, record)
	}
	encoded, err := json.Marshal(records)
	if err != nil {
		return "[]"
	}
	return string(encoded)
}

func validateClosure(key GroupKey, closure Closure) ([]decodedBaseline, error) {
	if !key.Valid() {
		return nil, ErrInvalidGroupKey
	}
	d := closure.Destination
	if closure.NextClose.IsZero() || d.Scope == "" || d.ConnectorID == "" || d.ConnectorRevision == "" || d.TableName == "" {
		return nil, fmt.Errorf("%w: next close and a frozen destination are required", ErrInvalidClosure)
	}
	for _, bucket := range closure.Buckets {
		outcome := bucket.Outcome
		if outcome.BucketStart.IsZero() {
			return nil, fmt.Errorf("%w: bucket start is required", ErrInvalidClosure)
		}
		isRow := outcome.Kind == snapshot.OutcomeRow
		if isRow != (bucket.Row != nil) {
			return nil, fmt.Errorf("%w: only row outcomes carry encoded values", ErrInvalidClosure)
		}
		if isRow && (outcome.EffectKey == "" || bucket.Row.EffectKey != outcome.EffectKey || bucket.Row.RecordID != outcome.RecordID) {
			return nil, fmt.Errorf("%w: row identity does not match its outcome", ErrInvalidClosure)
		}
	}
	return validateBaselines(key, closure.NextClose, closure.Baselines)
}

// CommitClosure persists closed buckets and advances the checkpoint in one
// transaction: outbox rows, bucket outcomes, consumed samples and checkpoint
// all commit together or not at all. Repeating a committed closure is a no-op;
// an existing effect key with different content is ErrEffectConflict.
func (s *Store) CommitClosure(ctx context.Context, key GroupKey, closure Closure) error {
	baselines, err := validateClosure(key, closure)
	if err != nil {
		return err
	}
	err = s.inTx(ctx, "closure", func(tx *sql.Tx) error {
		if err := verifyClosureCounterProofs(ctx, tx, key, closure.CounterProofs, closure.NextClose); err != nil {
			return err
		}
		now := timestamp(s.now())
		if err := storeBaselines(ctx, tx, key, baselines, now); err != nil {
			return err
		}
		for _, bucket := range closure.Buckets {
			outcome := bucket.Outcome
			if _, err := tx.ExecContext(ctx, `
				INSERT INTO wg_delivery_buckets (group_id, group_revision, entity_key, bucket_start, kind, reason, record_id, members, created_at)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
				ON CONFLICT (group_id, group_revision, entity_key, bucket_start) DO NOTHING`,
				key.GroupID, key.GroupRevision, outcome.EntityKey, timestamp(outcome.BucketStart), string(outcome.Kind),
				bucketReason(outcome), outcome.RecordID, memberRecords(outcome), now,
			); err != nil {
				return fmt.Errorf("record bucket outcome: %w", err)
			}
			if bucket.Row == nil {
				continue
			}
			if err := s.enqueueRow(ctx, tx, key, closure.Destination, bucket, now); err != nil {
				return err
			}
		}
		if _, err := tx.ExecContext(ctx, `
			UPDATE wg_delivery_samples SET consumed = 1
			WHERE group_id = ? AND group_revision = ? AND consumed = 0 AND bucket_start < ?`,
			key.GroupID, key.GroupRevision, timestamp(closure.NextClose),
		); err != nil {
			return fmt.Errorf("consume closed samples: %w", err)
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO wg_delivery_checkpoints (group_id, group_revision, next_close, updated_at)
			VALUES (?, ?, ?, ?)
			ON CONFLICT (group_id, group_revision) DO UPDATE
			SET next_close = excluded.next_close, updated_at = excluded.updated_at
			WHERE excluded.next_close > wg_delivery_checkpoints.next_close`,
			key.GroupID, key.GroupRevision, timestamp(closure.NextClose), now,
		); err != nil {
			return fmt.Errorf("advance checkpoint: %w", err)
		}
		return nil
	})
	s.usage.invalidate()
	return s.capacityError(ctx, err)
}

func (s *Store) enqueueRow(ctx context.Context, tx *sql.Tx, key GroupKey, d FrozenDestination, bucket ClosedBucket, now string) error {
	payload, digest, err := encodeRowPayload(*bucket.Row)
	if err != nil {
		return err
	}
	outcome := bucket.Outcome
	var existing string
	err = tx.QueryRowContext(ctx, `SELECT payload_digest FROM wg_delivery_outbox WHERE effect_key = ?`, outcome.EffectKey).Scan(&existing)
	switch {
	case err == nil && existing == digest:
		return nil
	case err == nil:
		return ErrEffectConflict
	case !errors.Is(err, sql.ErrNoRows):
		return fmt.Errorf("check outbox effect: %w", err)
	}
	capability := d.DedupeCapability
	if capability == "" {
		capability = DedupeNone
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO wg_delivery_outbox (
			effect_key, record_id, workspace_id, group_id, group_revision, entity_key, bucket_start, partition_key,
			destination_scope, connector_id, connector_revision, database_name, table_schema, table_name,
			dedupe_capability, record_key_column, payload, payload_digest, state, next_retry_at, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		outcome.EffectKey, outcome.RecordID, key.WorkspaceID, key.GroupID, key.GroupRevision, outcome.EntityKey,
		timestamp(outcome.BucketStart), partitionKey(key, outcome.EntityKey),
		d.Scope, d.ConnectorID, d.ConnectorRevision, d.Database, d.TableSchema, d.TableName,
		capability, d.RecordKeyColumn, string(payload), digest, StatePending, now, now, now,
	); err != nil {
		return fmt.Errorf("enqueue outbox row: %w", err)
	}
	return nil
}

// partitionKey orders delivery per group and entity; other partitions never wait on it.
func partitionKey(key GroupKey, entityKey string) string {
	return key.GroupID + "\x1f" + entityKey
}
