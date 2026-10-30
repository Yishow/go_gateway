package sourcerule

import (
	"context"
	"fmt"
	"go-gateway/internal/datalink/schema"
)

func (s *Service) removeOwnedPointRuntime(ctx context.Context, id string) bool {
	if s.runtimeSync == nil {
		return true
	}
	runtime, ok := s.runtimeSync.(pointInstallationOwner)
	j := s.pointJournal(ctx)
	if !ok || j == nil || j.points[id] == nil {
		return false
	}
	e := j.points[id]
	epoch, removed := runtime.RemovePointOwned(id, e.currentLease)
	if !removed {
		return false
	}
	if e.currentLease != 0 {
		e.removedLease = e.currentLease
		e.removedEpoch = epoch
	}
	e.currentLease = 0
	return true
}
func advanceOwnedEpoch(j *pointWriteJournal, epoch uint64) uint64 {
	for epoch != ^uint64(0) && j.ownedEpochs[epoch+1] {
		epoch++
	}
	return epoch
}
func (s *Service) restoreOwnedRuntime(ctx context.Context, restored []*schema.Point) error {
	if s.runtimeSync == nil {
		return nil
	}
	j := s.pointJournal(ctx)
	runtime, ok := s.runtimeSync.(pointInstallationOwner)
	if !ok || j == nil {
		return fmt.Errorf("runtime installation ownership unavailable")
	}
	restoredRemoved := map[string]bool{}
	// Restore proven removals before ordinary replacements advance the global epoch.
	for _, p := range restored {
		e := j.points[p.ID]
		if e == nil {
			return fmt.Errorf("point write ownership unavailable")
		}
		if !p.Enabled || e.removedLease == 0 {
			continue
		}
		lease, restored := s.restoreRemovedProfile(j, p, e.removedLease, advanceOwnedEpoch(j, e.removedEpoch))
		if !restored {
			return fmt.Errorf("removed installation epoch changed during compensation")
		}
		e.removedLease = 0
		s.observeRuntimeInstallation(ctx, p, lease)
		restoredRemoved[p.ID] = true
	}
	for _, p := range restored {
		if restoredRemoved[p.ID] {
			continue
		}
		e := j.points[p.ID]
		if !p.Enabled {
			if !s.removeOwnedPointRuntime(ctx, p.ID) {
				return fmt.Errorf("runtime installation changed during compensation")
			}
			continue
		}
		if e.currentLease == 0 && j.baselineKnown[p.ID] && len(j.baselineLeases[p.ID]) == 0 {
			if _, absent := runtime.RemovePointOwned(p.ID, 0); !absent {
				return fmt.Errorf("foreign installation appeared during compensation")
			}
			continue
		}
		lease, restored := s.restorePointProfile(j, p, e.currentLease)
		if !restored {
			return fmt.Errorf("runtime installation changed during compensation")
		}
		s.observeRuntimeInstallation(ctx, p, lease)
	}
	return nil
}
