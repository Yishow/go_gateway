package sourcerule

import (
	"context"
	"errors"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
)

func (s *Service) createRulePoint(ctx context.Context, req point.CreatePointRequest) (*schema.Point, error) {
	return s.pointSvc.CreateMaterializedObserved(ctx, req, s.observePointWrite(ctx))
}
func (s *Service) updateRulePoint(ctx context.Context, before *schema.Point, req point.UpdatePointRequest) (*schema.Point, error) {
	return s.pointSvc.UpdateMaterializedObserved(ctx, before, req, s.observePointWrite(ctx))
}
func (s *Service) deleteRulePoint(ctx context.Context, id string) error {
	var p *schema.Point
	var err error
	if j := s.pointJournal(ctx); j != nil && j.points[id] != nil {
		p = clonePointSnapshot(j.points[id].current)
		if p == nil {
			return nil
		}
	} else {
		p, err = s.pointSvc.GetByID(ctx, id)
	}
	if errors.Is(err, point.ErrPointNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	return s.pointSvc.DeleteMaterializedObserved(ctx, p, s.observePointWrite(ctx))
}
