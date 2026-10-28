package groupdelivery

import "context"

// PendingSampleCount counts accepted, unfinalized inputs of one complete scope.
// It is not a count of SQL outbox rows or a bounded diagnostic sample.
func (s *Store) PendingSampleCount(ctx context.Context, key GroupKey) (int, error) {
	if !key.Valid() {
		return 0, ErrInvalidGroupKey
	}
	if err := validateHandoffScope(ctx, s.db, key.GroupID, key.GroupRevision); err != nil {
		return 0, err
	}
	var count int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM wg_delivery_samples WHERE workspace_id=? AND group_id=? AND group_revision=? AND consumed=0 AND NOT EXISTS(SELECT 1 FROM wg_delivery_sample_handoffs h WHERE h.source_seq=wg_delivery_samples.seq)`, key.WorkspaceID, key.GroupID, key.GroupRevision).Scan(&count)
	return count, err
}
