package workspace

import "context"

// ReadRuntimeProjection rebuilds a detached persisted projection without bootstrap.
// It neither installs nor repairs sources and is safe for read-only previews.
func (s *Service) ReadRuntimeProjection(ctx context.Context) (*RuntimeProjection, error) {
	if err := s.ensureRuntimeProjectionReaders(); err != nil {
		return nil, err
	}
	record, err := s.ReadSnapshot(ctx)
	if err != nil {
		return nil, err
	}
	return s.runtimeProjectionFromRecord(ctx, record)
}
