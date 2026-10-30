package runtime

import (
	"context"
	"fmt"
	"go-gateway/internal/datalink/schema"
)

type sourceProjectionLeaseKey struct{}
type sourceProjectionLeaseGuard struct {
	deviceID string
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
	s.pointMetaMu.Lock()
	defer s.pointMetaMu.Unlock()
	guard := sourceProjectionGuard(ctx)
	if guard == nil {
		s.setPointMetaLocked(p)
		return s.scheduler.AddPointWithObservationLease(p), nil
	}
	if p.DeviceID != guard.deviceID {
		return 0, fmt.Errorf("source projection point scope changed")
	}
	lease, ok := s.scheduler.InstallPointAtLease(p, guard.expected[p.ID])
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
		if !s.scheduler.RemovePointAtLease(id, guard.expected[id]) {
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
