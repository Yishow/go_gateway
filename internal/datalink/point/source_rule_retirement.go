package point

import (
	"context"

	"go-gateway/internal/datalink/schema"
)

type unmappedPointRevisionDeleter interface {
	deleteUnmappedRevision(context.Context, string, string, string) error
}

// DeleteUnmappedMaterializedObserved guards SQL deletion against cascading a
// new mapping. The independent memory repository has no cascading deletes.
func (s *Service) DeleteUnmappedMaterializedObserved(ctx context.Context, planned *schema.Point, observer MaterializedWriteObserver) error {
	repo, ok := s.repo.(unmappedPointRevisionDeleter)
	if !ok || planned == nil {
		return ErrFormatUnavailable
	}
	finish := s.beginAcquisitionMutation(ctx, planned.DeviceID)
	defer finish()
	complete := observeMaterialized(observer, planned, nil)
	err := repo.deleteUnmappedRevision(ctx, planned.ID, planned.DeviceID, planned.RevisionID)
	complete(err)
	return err
}

func (r *SQLRepository) deleteUnmappedRevision(ctx context.Context, id, deviceID, expected string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM points WHERE id=? AND device_id=? AND revision_id=? AND NOT EXISTS (SELECT 1 FROM mappings WHERE point_id=points.id)`, id, deviceID, expected)
	if err != nil {
		return ErrFormatUnavailable
	}
	n, err := result.RowsAffected()
	if err != nil {
		return ErrFormatUnavailable
	}
	if n != 1 {
		return ErrFormatConflict
	}
	return nil
}

func (r *MemoryRepository) deleteUnmappedRevision(ctx context.Context, id, deviceID, expected string) error {
	// Memory point deletion has no cascading mapping side effects.
	return r.DeleteRevision(ctx, id, deviceID, expected)
}
