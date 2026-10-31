package runtime

import (
	"context"
	"fmt"
	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
)

type sourceProjectionLeaseKey struct{}
type sourceProjectionLeaseGuard struct {
	deviceID string
	proof    *sourcerule.RuntimeProjectionProof
	removed  func(string, uint64)
	expected map[string]uint64
}

func (s *Service) beginSourceProjectionLeaseGuard(ctx context.Context, deviceID string) context.Context {
	if s == nil || s.scheduler == nil {
		return ctx
	}
	return context.WithValue(ctx, sourceProjectionLeaseKey{}, &sourceProjectionLeaseGuard{deviceID: deviceID, expected: s.scheduler.CaptureDevicePointLeases(deviceID)})
}
func sourceProjectionGuard(ctx context.Context) *sourceProjectionLeaseGuard {
	g, _ := ctx.Value(sourceProjectionLeaseKey{}).(*sourceProjectionLeaseGuard)
	return g
}
func (s *Service) installProjectionPoint(ctx context.Context, p *schema.Point) (uint64, error) {
	profile, err := s.pointReadProfile(ctx, p)
	if err != nil {
		return 0, err
	}
	s.pointMetaMu.Lock()
	defer s.pointMetaMu.Unlock()
	guard := sourceProjectionGuard(ctx)
	if guard == nil {
		lease := s.scheduler.AddPointWithObservationLease(p, profile)
		if modernPoint(p) && lease == 0 {
			return 0, connector.ErrReadProfile
		}
		s.setPointMetaLocked(p)
		return lease, nil
	}
	if p.DeviceID != guard.deviceID {
		return 0, fmt.Errorf("source projection point scope changed")
	}
	if !sourceProjectionPointOwned(guard, p) {
		return 0, fmt.Errorf("source projection point ownership changed")
	}
	var lease uint64
	var ok bool
	if guard.proof != nil {
		lease, ok = s.scheduler.ApplyPointMutationReceipt(p, collector.PointMutationLease{PointID: p.ID, Lease: guard.expected[p.ID], Epoch: guard.proof.Epoch}, profile)
		if ok {
			guard.proof.Epoch = lease
		}
	} else {
		lease, ok = s.scheduler.InstallPointAtLease(p, guard.expected[p.ID], profile)
	}
	if !ok {
		return 0, fmt.Errorf("source projection installation changed")
	}
	s.setPointMetaLocked(p)
	guard.expected[p.ID] = lease
	return lease, nil
}
func (s *Service) removeProjectionPoint(ctx context.Context, id string) error {
	return s.stopProjectionPoint(ctx, id, nil)
}

// A disabled point keeps its metadata; a deleted point removes it. Both updates
// serialize with point installation, so a rejected stale removal changes neither.
func (s *Service) stopProjectionPoint(ctx context.Context, id string, keep *schema.Point) error {
	s.pointMetaMu.Lock()
	defer s.pointMetaMu.Unlock()
	guard := sourceProjectionGuard(ctx)
	if guard == nil {
		s.scheduler.RemovePoint(id)
	} else {
		if guard.proof != nil {
			proof, owned := guard.proof.Points[id]
			if !owned || (keep == nil && !proof.Deleted) || (keep != nil && !sourceProjectionPointOwned(guard, keep)) {
				return fmt.Errorf("source projection removal ownership changed")
			}
			if _, ok := s.scheduler.ApplyPointMutationReceipt(&schema.Point{ID: id}, collector.PointMutationLease{PointID: id, Lease: guard.expected[id], Epoch: guard.proof.Epoch}); !ok {
				return fmt.Errorf("source projection installation changed")
			}
			guard.proof.Epoch++
			if guard.removed != nil {
				guard.removed(id, guard.proof.Epoch)
			}
		} else if !s.scheduler.RemovePointAtLease(id, guard.expected[id]) {
			return fmt.Errorf("source projection installation changed")
		}
		guard.expected[id] = 0
	}
	if keep != nil {
		s.setPointMetaLocked(keep)
	} else {
		delete(s.pointMetaIndex, id)
	}
	return nil
}
