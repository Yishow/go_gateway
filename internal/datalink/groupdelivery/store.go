package groupdelivery

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-gateway/internal/datalink/snapshot"
)

// Store is the local SQLite-backed durable state for write groups.
type Store struct {
	db              *sql.DB
	now             func() time.Time
	quota           *QuotaConfig
	usage           *usageCache
	canonicalIntake bool
}

// NewStore wraps a migrated local database handle for standalone ledger users.
// Production canonical write-group intake must opt into WithCanonicalIntake.
func NewStore(db *sql.DB) *Store {
	return &Store{db: db, now: func() time.Time { return time.Now().UTC() }, usage: newUsageCache()}
}

// WithClock returns a clone using a deterministic clock without mutating s.
func (s *Store) WithClock(now func() time.Time) *Store {
	clone := *s
	if now != nil {
		clone.now = now
	}
	return &clone
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
	if sample.AcceptanceOrigin != (snapshot.AcceptanceOrigin{}) {
		return ErrOwnershipConflict
	}
	if sample.SampleID == "" || sample.MemberKey == "" || sample.ObservedAt.IsZero() || bucketStart.IsZero() {
		return ErrInvalidSample
	}
	payload, digest, err := encodeSample(sample)
	if err != nil {
		return err
	}
	inserted := false
	err = s.inTx(ctx, "journal", func(tx *sql.Tx) error {
		if err := refuseOrdinaryHandoffAppend(ctx, tx, key, sample); err != nil {
			return err
		}
		if err := validateCutoverAcceptance(ctx, tx, key, sample); err != nil {
			return err
		}
		_, found, err := verifyExistingSample(ctx, tx, key, sample, digest)
		if err != nil || found {
			return err
		}
		if err := verifySampleProof(ctx, tx, key, sample, true); err != nil {
			return err
		}
		if err := s.validateLifecycleIntake(ctx, tx, key, sample.ObservedAt); err != nil {
			return err
		}
		if err := s.checkQuota(ctx, tx, key); err != nil {
			return err
		}
		address, err := allocateJournalAddress(ctx, tx, key, sample.SampleID)
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO wg_delivery_samples (
				workspace_id, group_id, group_revision, sample_id, member_key,
				observed_at, bucket_start, payload, payload_digest, consumed, created_at,origin_workspace_id,origin_group_id,origin_revision,origin_sample_id
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 0, ?, ?, ?, ?, ?)`,
			key.WorkspaceID, key.GroupID, key.GroupRevision, address, sample.MemberKey,
			timestamp(sample.ObservedAt), timestamp(bucketStart), string(payload), digest, timestamp(s.now()), key.WorkspaceID, key.GroupID, key.GroupRevision, sample.SampleID,
		); err != nil {
			return fmt.Errorf("append journal sample: %w", err)
		}
		inserted = true
		return nil
	})
	if err == nil && inserted {
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
	return s.restoreHandoffView(ctx, key)
}
