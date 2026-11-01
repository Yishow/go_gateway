package point

import (
	"context"

	"go-gateway/internal/datalink/schema"
)

// RestoreRevision compensates only the configuration revision owned by a
// materializer. It neither resurrects deleted identities nor replaces telemetry.
// Repositories without atomic configuration CAS must fail closed.
func (s *Service) RestoreRevision(ctx context.Context, snapshot *schema.Point, appliedRevision string) error {
	if snapshot == nil || appliedRevision == "" {
		return ErrFormatUnavailable
	}
	repo, ok := s.repo.(pointRevisionUpdater)
	if !ok {
		return ErrFormatUnavailable
	}
	finish := s.beginAcquisitionMutation(ctx, snapshot.DeviceID)
	defer finish()
	return repo.UpdateRevision(ctx, clonePoint(snapshot), appliedRevision)
}
