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

func (s *Service) validateRulePointFormats(ctx context.Context, current, next *schema.SourceRule, protocol schema.ProtocolType, links []*schema.SourceRuleLink) error {
	if current.DataType == next.DataType {
		return nil
	}
	for _, link := range links {
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
		candidate.DataType = next.DataType
		if _, err = decoding.ResolvePoint(candidate, decoding.FormatContext{Protocol: protocol, RuleOwned: true, RuleFormat: next.DataFormat}); err != nil {
			return &formatPreflightError{fmt.Errorf("%w: reset the affected point format override before changing its type", ErrValidation)}
		}
	}
	return nil
}
