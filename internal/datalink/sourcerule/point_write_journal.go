package sourcerule

import (
	"context"
	"errors"
	"fmt"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
)

type ownedPointWrite struct {
	leases          map[uint64]bool
	currentLease    uint64
	removedLease    uint64
	removedEpoch    uint64
	before, current *schema.Point
	attempted       []*schema.Point
	uncertain       bool
}
type pointWriteJournal struct {
	projectionEpoch      uint64
	projectionEpochKnown bool
	baselineProfiles     map[string]connector.ReadProfile
	points               map[string]*ownedPointWrite
	ruleTargets          map[[32]byte]bool
	baselineLeases       map[string]map[uint64]bool
	baselineKnown        map[string]bool
	ownedEpochs          map[uint64]bool
	deviceID             string
	runtimeUncertain     bool
	ruleUncertain        bool
	storageFailure       error
	ruleDeleted          bool
	candidateWrites      map[string]bool
}

func (s *Service) pointJournal(ctx context.Context) *pointWriteJournal {
	owner, ok := ctx.Value(ruleMutationContextKey{}).(*ruleMutationOwnership)
	if !ok || owner.service != s || !owner.active.Load() {
		return nil
	}
	if owner.points == nil {
		owner.points = &pointWriteJournal{points: map[string]*ownedPointWrite{}}
	}
	return owner.points
}
func (s *Service) observePointWrite(ctx context.Context) point.MaterializedWriteObserver {
	return func(before, after *schema.Point) func(error) {
		journal := s.pointJournal(ctx)
		if journal == nil {
			return func(error) {}
		}
		identity := before
		if identity == nil {
			identity = after
		}
		s.captureProjectionEpoch(ctx, identity.ID)
		entry := journal.points[identity.ID]
		if entry == nil {
			entry = &ownedPointWrite{before: clonePointSnapshot(before), current: clonePointSnapshot(before)}
			journal.points[identity.ID] = entry
		}
		if entry.leases == nil {
			entry.leases = map[uint64]bool{}
		}
		if before != nil {
			s.capturePointBaseline(ctx, before)
			for lease := range journal.baselineLeases[before.ID] {
				entry.leases[lease] = true
				if entry.currentLease == 0 {
					entry.currentLease = lease
				}
			}

			entry.attempted = append(entry.attempted, clonePointSnapshot(before))
		}
		if after != nil {
			entry.attempted = append(entry.attempted, clonePointSnapshot(after))
		}
		return func(err error) {
			s.noteTransactionFailure(ctx, err)
			if err == nil {
				entry.current = clonePointSnapshot(after)
				return
			}
			// CAS refusal proves no write; every other storage error is ambiguous.
			if !errors.Is(err, point.ErrFormatConflict) && !errors.Is(err, point.ErrPointNotFound) {
				entry.uncertain = true
			}
		}
	}
}
func (s *Service) notePointRestored(ctx context.Context, p *schema.Point) {
	if journal := s.pointJournal(ctx); journal != nil {
		if e := journal.points[p.ID]; e != nil {
			e.current = clonePointSnapshot(p)
		}
	}
}

// restoreOwnedPoints consumes actual write receipts, never a post-operation read
// as authority. All storage must succeed before any runtime snapshot is installed.
func (s *Service) restoreOwnedPoints(ctx context.Context, journal *pointWriteJournal) ([]*schema.Point, error) {
	if journal == nil {
		return nil, nil
	}
	for _, e := range journal.points {
		if e.uncertain {
			return nil, fmt.Errorf("point write outcome is uncertain")
		}
	}
	var restored []*schema.Point
	for _, e := range journal.points {
		if e.before == nil {
			if e.current != nil {
				if err := s.pointSvc.DeleteRevision(ctx, e.current); err != nil {
					return nil, err
				}
				e.current = nil
			}
			continue
		}
		if e.current == nil {
			if err := s.pointSvc.RestoreAbsent(ctx, e.before); err != nil {
				return nil, err
			}
		} else {
			if err := s.pointSvc.RestoreRevision(ctx, e.before, e.current.RevisionID); err != nil {
				return nil, err
			}
		}
		e.current = clonePointSnapshot(e.before)
		restored = append(restored, clonePointSnapshot(e.before))
	}
	return restored, nil
}

type pointInstallationOwner interface {
	CapturePointLease(string, string) uint64
	SuspendPointLease(string, uint64) bool
	RestorePointLeaseReceipt(*schema.Point, uint64) (uint64, bool)
	InstallPointForMutationAtLease(*schema.Point, uint64) (uint64, bool)
	RemovePointOwned(string, uint64) (uint64, bool)
	RestoreRemovedPointOwned(*schema.Point, uint64, uint64) (uint64, bool)
}

func (s *Service) suspendOwnedPoints(ctx context.Context) {
	journal := s.pointJournal(ctx)
	if journal == nil {
		return
	}
	if runtime, ok := s.runtimeSync.(pointInstallationOwner); ok {
		for id, leases := range journal.baselineLeases {
			for lease := range leases {
				runtime.SuspendPointLease(id, lease)
			}
		}
		for id, e := range journal.points {
			for lease := range e.leases {
				runtime.SuspendPointLease(id, lease)
			}
		}
	}
}
func (s *Service) suspendOwnedPoint(ctx context.Context, id string) {
	journal := s.pointJournal(ctx)
	if journal == nil {
		return
	}
	e := journal.points[id]
	if e == nil {
		return
	}
	if runtime, ok := s.runtimeSync.(pointInstallationOwner); ok {
		for lease := range e.leases {
			runtime.SuspendPointLease(id, lease)
		}
	}
}

func (s *Service) finalizePointWrites(ctx context.Context, result *error) {

	if j := s.pointJournal(ctx); j != nil {
		if j.runtimeUncertain || j.ruleUncertain || j.storageFailure != nil {
			*result = s.dirtyUnknown(ctx)
			return
		}
		for _, e := range j.points {
			if e.uncertain {
				*result = s.dirtyUnknown(ctx)
				return
			}
		}
	}
}

func (s *Service) observeRuntimeInstallation(ctx context.Context, p *schema.Point, lease uint64) {
	if p == nil || lease == 0 {
		return
	}
	j := s.pointJournal(ctx)
	if j == nil {
		return
	}
	if p.DeviceID != j.deviceID {
		return
	}
	e := j.points[p.ID]
	if e == nil || e.current == nil || e.current.RevisionID != p.RevisionID || e.current.DeviceID != p.DeviceID {
		return
	}
	if j.ownedEpochs == nil {
		j.ownedEpochs = map[uint64]bool{}
	}
	j.ownedEpochs[lease] = true
	if e.leases == nil {
		e.leases = map[uint64]bool{}
	}
	e.leases[lease] = true
	e.currentLease = lease
}

func (s *Service) capturePointBaseline(ctx context.Context, p *schema.Point) {
	if p == nil {
		return
	}
	runtime, ok := s.runtimeSync.(pointInstallationOwner)
	if !ok {
		return
	}
	j := s.pointJournal(ctx)
	if j == nil {
		return
	}
	if j.baselineKnown == nil {
		j.baselineKnown = map[string]bool{}
	}
	if j.baselineKnown[p.ID] {
		return
	}
	s.captureProjectionEpoch(ctx, p.ID)
	j.baselineKnown[p.ID] = true
	lease := runtime.CapturePointLease(p.ID, p.RevisionID)
	if lease == 0 {
		return
	}
	s.captureBaselineReadProfile(j, p, lease)
	if j.baselineLeases == nil {
		j.baselineLeases = map[string]map[uint64]bool{}
	}
	if j.baselineLeases[p.ID] == nil {
		j.baselineLeases[p.ID] = map[uint64]bool{}
	}
	j.baselineLeases[p.ID][lease] = true
}
func (s *Service) captureLinkBaselines(ctx context.Context, links []*schema.SourceRuleLink) error {
	if _, ok := s.runtimeSync.(pointInstallationOwner); !ok {
		return nil
	}
	for _, link := range links {
		p, err := s.pointSvc.GetByID(ctx, link.PointID)
		if err != nil {
			return err
		}
		s.capturePointBaseline(ctx, p)
	}
	return nil
}
