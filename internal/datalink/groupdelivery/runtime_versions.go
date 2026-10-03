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
	rows, err := s.db.QueryContext(ctx, `SELECT workspace_id, group_id, group_revision,
		MIN(bucket_start), MAX(bucket_start) FROM wg_delivery_samples WHERE consumed = 0
		GROUP BY workspace_id, group_id, group_revision`)
	if err != nil {
		return nil, fmt.Errorf("read open journals: %w", err)
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
	return journals, rows.Err()
}
