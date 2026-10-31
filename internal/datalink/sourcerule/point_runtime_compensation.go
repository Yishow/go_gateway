package sourcerule

import (
	"context"
	"errors"
	"fmt"
	"go-gateway/internal/datalink/schema"
)

func (s *Service) syncPoints(ctx context.Context, points []*schema.Point) {
	if s.runtimeSync == nil {
		return
	}
	for _, pointRecord := range points {
		if runtime, ok := s.runtimeSync.(pointInstallationOwner); ok {
			if journal := s.pointJournal(ctx); journal != nil {
				if e := journal.points[pointRecord.ID]; e != nil {
					required, scopeErr := s.pointRuntimeRequired(ctx, pointRecord)
					if scopeErr != nil {
						journal.runtimeUncertain = true
						continue
					}
					if !required {
						if !s.removeOwnedPointRuntime(ctx, pointRecord.ID) {
							journal.runtimeUncertain = true
						}
					} else {
						lease, installed := runtime.InstallPointForMutationAtLease(pointRecord, e.currentLease)
						if !installed {
							journal.runtimeUncertain = true
						} else {
							s.observeRuntimeInstallation(ctx, pointRecord, lease)
						}
					}
					continue
				}
			}
		}
		s.runtimeSync.UpsertPoint(pointRecord)
	}
}

func (s *Service) rollbackUpdatedPoints(ctx context.Context, plans []pointUpdatePlan) error {
	if err := s.storageCompensationError(ctx); err != nil {
		return err
	}

	var rollbackErrs []error
	for _, plan := range plans {
		// Compensation restores the exact identity and intent, not a new operator edit.
		err := s.pointSvc.RestoreRevision(ctx, plan.point, plan.appliedRevision)
		if err != nil {
			if s.noteTransactionFailure(ctx, err) {
				return err
			}
			rollbackErrs = append(rollbackErrs, fmt.Errorf("restore point %s: %w", plan.point.ID, err))
			continue
		}
		s.notePointRestored(ctx, plan.point)
		if err := s.restoreOwnedRuntime(ctx, []*schema.Point{plan.point}); err != nil {
			rollbackErrs = append(rollbackErrs, err)
		}
	}
	return errors.Join(rollbackErrs...)
}
