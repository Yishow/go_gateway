package point

import (
	"context"
	"time"

	common "go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
)

// MaterializedChange updates one exact point revision as part of a rule edit.
type MaterializedChange struct {
	Before  *schema.Point
	Request UpdatePointRequest
}

// RevisionRestore is one original point and the revision owned by the caller.
type RevisionRestore struct {
	Point            *schema.Point
	ExpectedRevision string
}

type materializedBatchRepository interface {
	updateMaterializedRevisions(context.Context, []RevisionRestore) error
}

// UpdateMaterializedBatchObserved commits a group of rule-owned point changes
// atomically. This permits address exchanges without exposing temporary values.
func (s *Service) UpdateMaterializedBatchObserved(ctx context.Context, changes []MaterializedChange, observer MaterializedWriteObserver) ([]*schema.Point, error) {
	repo, ok := s.repo.(materializedBatchRepository)
	if !ok || len(changes) == 0 {
		return nil, ErrFormatUnavailable
	}
	updates := make([]RevisionRestore, 0, len(changes))
	points := make([]*schema.Point, 0, len(changes))
	deviceID := ""
	for _, change := range changes {
		before, req := change.Before, change.Request
		if before == nil || before.RevisionID == "" || req.DataFormat != nil || req.FormatMode != nil || req.ConfiguredFormat != nil {
			return nil, ErrFormatInvalid
		}
		if deviceID == "" {
			deviceID = before.DeviceID
		} else if deviceID != before.DeviceID {
			return nil, ErrFormatInvalid
		}
		next := clonePoint(before)
		if err := applyPointUpdate(next, req); err != nil {
			return nil, err
		}
		revision, err := common.NewUUID()
		if err != nil {
			return nil, ErrFormatUnavailable
		}
		next.RevisionID, next.UpdatedAt = revision, time.Now()
		updates = append(updates, RevisionRestore{Point: next, ExpectedRevision: before.RevisionID})
		points = append(points, next)
	}
	finish := s.beginAcquisitionMutation(ctx, deviceID)
	defer finish()
	callbacks := make([]func(error), len(changes))
	for i, change := range changes {
		callbacks[i] = observeMaterialized(observer, change.Before, points[i])
	}
	err := repo.updateMaterializedRevisions(ctx, updates)
	for _, complete := range callbacks {
		complete(err)
	}
	if err != nil {
		return nil, err
	}
	return points, nil
}

// RestoreRevisions restores a set of owned point revisions atomically.
func (s *Service) RestoreRevisions(ctx context.Context, restores []RevisionRestore) error {
	repo, ok := s.repo.(materializedBatchRepository)
	if !ok || len(restores) == 0 || restores[0].Point == nil {
		return ErrFormatUnavailable
	}
	finish := s.beginAcquisitionMutation(ctx, restores[0].Point.DeviceID)
	defer finish()
	return repo.updateMaterializedRevisions(ctx, restores)
}
