package workspace

import "context"

// ReadSnapshot reads an existing workspace without initializing or repairing it.
// The returned record is detached, using the same copy policy as saved reads.
func (s *Service) ReadSnapshot(ctx context.Context) (*Record, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if s == nil || s.repo == nil {
		return nil, ErrDatabaseSetupUnavailable
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	record, err := s.repo.Get(ctx)
	if err != nil {
		return nil, err
	}
	return cloneRecord(record), nil
}
