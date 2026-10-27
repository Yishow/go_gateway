package groupdelivery

import "context"

// PendingSampleCount counts accepted, unfinalized inputs of one complete scope.
// It is not a count of SQL outbox rows or a bounded diagnostic sample.
func (s *Store) PendingSampleCount(ctx context.Context, key GroupKey) (int, error) {
	if !key.Valid() {
		return 0, ErrInvalidGroupKey
	}
	var count int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM wg_delivery_samples WHERE workspace_id=? AND group_id=? AND group_revision=? AND consumed=0`, key.WorkspaceID, key.GroupID, key.GroupRevision).Scan(&count)
	return count, err
}
