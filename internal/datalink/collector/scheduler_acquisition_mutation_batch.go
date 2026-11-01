package collector

import "context"

// TryBeginAcquisitionGeometryMutations quarantines an exact pre-enumerated set
// without waiting for scheduler.mu. It is safe to call while holding a SQL
// writer reservation. Either every device receives a token or none does.
// The caller must release its SQL reservation before finalizing these tokens.
func (s *Scheduler) TryBeginAcquisitionGeometryMutations(ctx context.Context, deviceIDs []string) ([]AcquisitionGeometryMutation, bool) {
	if s == nil || ctx.Err() != nil || !s.mu.TryLock() {
		return nil, false
	}
	defer s.mu.Unlock()
	if ctx.Err() != nil {
		return nil, false
	}
	for i, id := range deviceIDs {
		if id == "" {
			return nil, false
		}
		for _, previous := range deviceIDs[:i] {
			if previous == id {
				return nil, false
			}
		}
	}
	tokens := make([]AcquisitionGeometryMutation, 0, len(deviceIDs))
	for _, id := range deviceIDs {
		tokens = append(tokens, s.beginAcquisitionGeometryMutationLocked(id))
	}
	return tokens, true
}
