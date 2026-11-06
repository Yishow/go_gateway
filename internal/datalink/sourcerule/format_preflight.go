package sourcerule

import (
	"context"
	"errors"
	"fmt"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
)

// No write has occurred when this error is returned; outer compensation must
// not replay a stale snapshot over an unrelated concurrent point edit.
type formatPreflightError struct{ error }

func (e *formatPreflightError) Unwrap() error { return e.error }

func (s *Service) validateRulePointFormats(ctx context.Context, current, next *schema.SourceRule, protocol schema.ProtocolType, links []*schema.SourceRuleLink, layout rulePointPlan, matched map[string]*schema.SourceRuleLink) error {
	if current.DataType == next.DataType && !isMixedRule(current) && !layout.mixed() {
		return nil
	}
	targetByPoint := make(map[string]string, len(matched))
	for address, link := range matched {
		targetByPoint[link.PointID] = address
	}
	for _, link := range links {
		address := link.Address
		if layout.mixed() {
			var retained bool
			address, retained = targetByPoint[link.PointID]
			if !retained {
				continue // removed items are deleted, not retyped
			}
		}
		p, err := s.pointSvc.GetByID(ctx, link.PointID)
		if errors.Is(err, point.ErrPointNotFound) {
			continue
		}
		if err != nil {
			return &formatPreflightError{err}
		}
		if p.FormatMode != schema.PointFormatOverride {
			continue
		}
		candidate := *p
		candidate.Address = address
		candidate.DataType = layout.dataType(address, next.DataType)
		if _, err = decoding.ResolvePoint(candidate, decoding.FormatContext{Protocol: protocol, RuleOwned: true, RuleFormat: layout.dataFormat(address, next.DataFormat)}); err != nil {
			return &formatPreflightError{fmt.Errorf("%w: reset the affected point format override before changing its type", ErrValidation)}
		}
	}
	return nil
}
