package point

import (
	"context"
	"go-gateway/internal/datalink/schema"
)

type pointRevisionDeleter interface {
	DeleteRevision(context.Context, string, string, string) error
}

func (s *Service) DeleteMaterializedObserved(ctx context.Context, planned *schema.Point, observer MaterializedWriteObserver) error {
	if planned == nil {
		return ErrFormatInvalid
	}
	repo, ok := s.repo.(pointRevisionDeleter)
	if !ok {
		return ErrFormatUnavailable
	}
	complete := observeMaterialized(observer, planned, nil)
	err := repo.DeleteRevision(ctx, planned.ID, planned.DeviceID, planned.RevisionID)
	complete(err)
	return err
}
func (s *Service) DeleteRevision(ctx context.Context, planned *schema.Point) error {
	return s.DeleteMaterializedObserved(ctx, planned, nil)
}

// RestoreAbsent never overwrites an identity another operation has restored.
func (s *Service) RestoreAbsent(ctx context.Context, p *schema.Point) error {
	if p == nil {
		return ErrFormatInvalid
	}
	return s.repo.Create(ctx, clonePoint(p))
}
func (r *SQLRepository) DeleteRevision(ctx context.Context, id, deviceID, expected string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM points WHERE id=? AND device_id=? AND revision_id=?`, id, deviceID, expected)
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
func (r *MemoryRepository) DeleteRevision(_ context.Context, id, deviceID, expected string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.points[id]
	if !ok || p.DeviceID != deviceID || p.RevisionID != expected {
		return ErrFormatConflict
	}
	delete(r.points, id)
	return nil
}
