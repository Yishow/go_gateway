package point

import (
	"context"
	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
	"time"
)

// UpdateMaterialized applies a rule plan only to the exact revision it read.
// Defaults are not operator intent; callers must omit all format edit fields.
func (s *Service) UpdateMaterialized(ctx context.Context, planned *schema.Point, req UpdatePointRequest) (*schema.Point, error) {
	if planned == nil || req.DataFormat != nil || req.FormatMode != nil || req.ConfiguredFormat != nil {
		return nil, ErrFormatInvalid
	}
	repo, ok := s.repo.(pointRevisionUpdater)
	if !ok {
		return nil, ErrFormatUnavailable
	}
	next := clonePoint(planned)
	if err := applyPointUpdate(next, req); err != nil {
		return nil, err
	}
	revision, err := common.NewUUID()
	if err != nil {
		return nil, ErrFormatUnavailable
	}
	next.RevisionID = revision
	next.UpdatedAt = time.Now()
	if err = repo.UpdateRevision(ctx, next, planned.RevisionID); err != nil {
		return nil, err
	}
	return next, nil
}
