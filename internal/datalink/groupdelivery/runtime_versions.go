package groupdelivery

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// SaveRuntimeVersion preserves verified recovery metadata before intake starts.
// A revision's first descriptor is immutable.
func (s *Store) SaveRuntimeVersion(ctx context.Context, key GroupKey, payload []byte) error {
	if !key.Valid() || len(payload) == 0 {
		return ErrInvalidGroupKey
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO wg_runtime_versions
		(workspace_id, group_id, group_revision, payload) VALUES (?, ?, ?, ?)
		ON CONFLICT (workspace_id, group_id, group_revision) DO NOTHING`,
		key.WorkspaceID, key.GroupID, key.GroupRevision, string(payload))
	if err != nil {
		return fmt.Errorf("save runtime version: %w", err)
	}
	return nil
}

// RuntimeVersion returns the frozen descriptor, or nil for a legacy revision.
func (s *Store) RuntimeVersion(ctx context.Context, key GroupKey) ([]byte, error) {
	var payload string
	err := s.db.QueryRowContext(ctx, `SELECT payload FROM wg_runtime_versions
		WHERE workspace_id = ? AND group_id = ? AND group_revision = ?`,
		key.WorkspaceID, key.GroupID, key.GroupRevision).Scan(&payload)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read runtime version: %w", err)
	}
	return []byte(payload), nil
}

// OpenJournal identifies revisions with accepted samples awaiting closure.
type OpenJournal struct {
	Key                     GroupKey
	FirstBucket, LastBucket time.Time
}

// OpenJournals discovers historical revisions independently of the live groups.
func (s *Store) OpenJournals(ctx context.Context) ([]OpenJournal, error) {
	// Receipt validation performs nested ownership reads. Keep discovery and
	// exclusions in one snapshot and on one connection, even with a pool of one.
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, &handoffStorageError{cause: err}
	}
	defer tx.Rollback()
	rows, err := tx.QueryContext(ctx, `SELECT workspace_id, group_id, group_revision,
		MIN(bucket_start), MAX(bucket_start) FROM wg_delivery_samples WHERE consumed = 0
		GROUP BY workspace_id, group_id, group_revision`)
	if err != nil {
		return nil, &handoffStorageError{cause: err}
	}
	defer rows.Close()
	var journals []OpenJournal
	for rows.Next() {
		var journal OpenJournal
		var first, last string
		if err := rows.Scan(&journal.Key.WorkspaceID, &journal.Key.GroupID, &journal.Key.GroupRevision, &first, &last); err != nil {
			return nil, fmt.Errorf("scan open journal: %w", err)
		}
		journal.FirstBucket, err = time.Parse(time.RFC3339Nano, first)
		if err != nil {
			return nil, fmt.Errorf("parse first journal bucket: %w", err)
		}
		journal.LastBucket, err = time.Parse(time.RFC3339Nano, last)
		if err != nil {
			return nil, fmt.Errorf("parse last journal bucket: %w", err)
		}
		journals = append(journals, journal)
	}
	if err := rows.Err(); err != nil {
		return nil, &handoffStorageError{cause: err}
	}
	if err := rows.Close(); err != nil {
		return nil, &handoffStorageError{cause: err}
	}
	var active []OpenJournal
	for _, journal := range journals {
		observed := &handoffReadObserver{ownershipQuery: tx}
		_, _, validationErr := handoffView(ctx, observed, journal.Key)
		if observed.failure != nil {
			return nil, &handoffStorageError{cause: observed.failure}
		}
		var storage *handoffStorageError
		if errors.As(validationErr, &storage) {
			return nil, validationErr
		}
		if validationErr != nil {
			// Discovery does not authorize exclusion. Keep a damaged scope so
			// its decisive Restore can quarantine it without hiding healthy ones.
			active = append(active, journal)
			continue
		}
		var first, last sql.NullString
		err := tx.QueryRowContext(ctx, `SELECT MIN(bucket_start),MAX(bucket_start) FROM wg_delivery_samples WHERE workspace_id=? AND group_id=? AND group_revision=? AND consumed=0 AND NOT EXISTS(SELECT 1 FROM wg_delivery_sample_handoffs h WHERE h.source_seq=wg_delivery_samples.seq)`, journal.Key.WorkspaceID, journal.Key.GroupID, journal.Key.GroupRevision).Scan(&first, &last)
		if err != nil {
			return nil, &handoffStorageError{cause: err}
		}
		if !first.Valid {
			continue
		}
		journal.FirstBucket, err = time.Parse(time.RFC3339Nano, first.String)
		if err != nil {
			return nil, ErrOwnershipUnavailable
		}
		journal.LastBucket, err = time.Parse(time.RFC3339Nano, last.String)
		if err != nil {
			return nil, ErrOwnershipUnavailable
		}
		active = append(active, journal)
	}
	if err := tx.Commit(); err != nil {
		return nil, &handoffStorageError{cause: err}
	}
	return active, nil
}
