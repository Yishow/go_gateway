package point

import (
	"context"
	"go-gateway/internal/datalink/schema"
)

// MaterializedWriteObserver is operation-local evidence, never an API intent.
// It runs before storage and receives the exact attempted identities/revisions.
// Its completion callback also receives uncertain storage errors.
type MaterializedWriteObserver func(before, after *schema.Point) func(error)

func observeMaterialized(observer MaterializedWriteObserver, before, after *schema.Point) func(error) {
	if observer == nil {
		return func(error) {}
	}
	return observer(clonePoint(before), clonePoint(after))
}

func (s *Service) CreateMaterializedObserved(ctx context.Context, req CreatePointRequest, observer MaterializedWriteObserver) (*schema.Point, error) {
	return s.create(ctx, req, observer)
}
