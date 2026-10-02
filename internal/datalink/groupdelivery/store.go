package groupdelivery

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"go-gateway/internal/datalink/snapshot"
)

// Store is the local SQLite-backed durable state for write groups.
type Store struct {
	db    *sql.DB
	now   func() time.Time
	quota *QuotaConfig
	usage *usageCache
}

// NewStore wraps a migrated local database handle.
func NewStore(db *sql.DB) *Store {
	return &Store{db: db, now: func() time.Time { return time.Now().UTC() }, usage: newUsageCache()}
}

// timestampLayout is fixed-width UTC so text columns sort chronologically.
const timestampLayout = "2006-01-02T15:04:05.000000000Z"

// DB exposes the underlying handle for read-only diagnostics and tests.
func (s *Store) DB() *sql.DB { return s.db }

func timestamp(t time.Time) string { return t.UTC().Format(timestampLayout) }

// AppendSample durably records an accepted sample in its own transaction.
// Returning nil means the sample is committed and may be ACKed. Resending the
// same sample is a no-op; the same ID with different content is ErrSampleConflict.
func (s *Store) AppendSample(ctx context.Context, key GroupKey, bucketStart time.Time, sample snapshot.Sample) error {
	if !key.Valid() {
		return ErrInvalidGroupKey
	}
	if sample.SampleID == "" || sample.MemberKey == "" || sample.ObservedAt.IsZero() || bucketStart.IsZero() {
		return ErrInvalidSample
	}
	payload, digest, err := encodeSample(sample)
	if err != nil {
		return err
	}
	err = s.inTx(ctx, "journal", func(tx *sql.Tx) error {
		var existing string
		err := tx.QueryRowContext(ctx,
			`SELECT payload_digest FROM wg_delivery_samples WHERE group_id = ? AND group_revision = ? AND sample_id = ?`,
			key.GroupID, key.GroupRevision, sample.SampleID).Scan(&existing)
		switch {
		case err == nil && existing == digest:
			return nil
		case err == nil:
			return ErrSampleConflict
		case !errors.Is(err, sql.ErrNoRows):
			return fmt.Errorf("check journal sample: %w", err)
		}
		if err := s.checkQuota(ctx, tx, key); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO wg_delivery_samples (
				workspace_id, group_id, group_revision, sample_id, member_key,
				observed_at, bucket_start, payload, payload_digest, consumed, created_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 0, ?)`,
			key.WorkspaceID, key.GroupID, key.GroupRevision, sample.SampleID, sample.MemberKey,
			timestamp(sample.ObservedAt), timestamp(bucketStart), string(payload), digest, timestamp(s.now()),
		); err != nil {
			return fmt.Errorf("append journal sample: %w", err)
		}
		return nil
	})
	if err == nil {
		s.usage.add(key.GroupID, int64(len(payload)))
	}
	return s.capacityError(ctx, err)
}

// capacityError turns an out-of-space failure into a capacity refusal; other
// errors pass through unchanged.
func (s *Store) capacityError(ctx context.Context, err error) error {
	if err == nil || !isDiskFull(err) {
		return err
	}
	used, _ := measureUsage(ctx, s.db, "") //nolint:errcheck // best effort: the refusal matters more than the figure
	return diskFullError(used)
}

// Restore returns the closure checkpoint and every unconsumed sample.
func (s *Store) Restore(ctx context.Context, key GroupKey) (Restored, error) {
	if !key.Valid() {
		return Restored{}, ErrInvalidGroupKey
	}
	var restored Restored
	var next string
	err := s.db.QueryRowContext(ctx,
		`SELECT next_close FROM wg_delivery_checkpoints WHERE group_id = ? AND group_revision = ?`,
		key.GroupID, key.GroupRevision).Scan(&next)
	switch {
	case err == nil:
		parsed, perr := time.Parse(time.RFC3339Nano, next)
		if perr != nil {
			return Restored{}, fmt.Errorf("parse checkpoint: %w", perr)
		}
		restored.NextClose = parsed.UTC()
	case !errors.Is(err, sql.ErrNoRows):
		return Restored{}, fmt.Errorf("read checkpoint: %w", err)
	}

	rows, err := s.db.QueryContext(ctx, `
		SELECT payload FROM wg_delivery_samples
		WHERE group_id = ? AND group_revision = ? AND consumed = 0
		ORDER BY seq ASC`, key.GroupID, key.GroupRevision)
	if err != nil {
		return Restored{}, fmt.Errorf("read journal samples: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var payload string
		if err := rows.Scan(&payload); err != nil {
			return Restored{}, fmt.Errorf("scan journal sample: %w", err)
		}
		sample, err := decodeSample([]byte(payload))
		if err != nil {
			return Restored{}, err
		}
		restored.Samples = append(restored.Samples, sample)
	}
	return restored, rows.Err()
}
