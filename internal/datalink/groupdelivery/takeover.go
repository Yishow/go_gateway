package groupdelivery

import (
	"context"
	"fmt"
)

// HasAppliedOwner reads the immutable, transactionally committed Apply history.
// Ownership survives active-worker retirement, deletion and destination edits.
// Drafts and rolled-back Apply transactions have no version and claim nothing.
func (s *Store) HasAppliedOwner(ctx context.Context, connectorID, tagID string) (bool, error) {
	var owned bool
	err := s.db.QueryRowContext(ctx, `SELECT EXISTS (
  SELECT 1 FROM write_group_versions v, json_each(v.payload, '$.members') m
  WHERE json_extract(v.payload, '$.destination.connector_id') = ?
   AND json_extract(m.value, '$.tag_id') = ?
   AND json_extract(v.payload, '$.applied_revision') = v.group_revision
   AND v.group_revision <> ''
 )`, connectorID, tagID).Scan(&owned)
	if err != nil {
		return false, fmt.Errorf("read applied writer ownership: %w", err)
	}
	return owned, nil
}
