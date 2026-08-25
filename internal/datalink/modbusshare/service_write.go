package modbusshare

import (
	"context"
	"fmt"
)

// WriteProjectedTagValue writes a value from the collector into the already
// reconciled Share projection. HTTP callers must use the candidate/apply/
// reconcile seam instead of this internal delivery authority.
func (s *Service) WriteProjectedTagValue(ctx context.Context, tagID string, value interface{}) error {
	hydration, hydrationErr := s.CheckHydration(ctx)
	if hydrationErr != nil {
		return NewError(ErrCodeHydrationRequired, "workspace Share hydration could not be verified", true)
	}
	if hydration.State != HydrationStateReady || !hydration.Readiness {
		return NewError(ErrCodeHydrationRequired, "workspace Share hydration is not ready", true)
	}
	s.mu.RLock()
	mapping, ok := s.mappings[tagID]
	s.mu.RUnlock()
	if !ok {
		return fmt.Errorf("mapping not found for tag_id: %s", tagID)
	}

	t, err := s.tagSvc.GetByID(ctx, tagID)
	if err != nil {
		return fmt.Errorf("tag not found: %w", err)
	}
	if t.DataType != mapping.DataType || mapping.SpanRegisters != DataTypeSpan(mapping.DataType) || mapping.StrideRegisters < mapping.SpanRegisters {
		s.markWriteDegraded()
		return NewError(ErrCodeInvalidGeometry, "runtime mapping identity is no longer aligned with the persisted tag", false)
	}

	words, err := encodeToWords(mapping.DataType, value)
	if err != nil {
		return err
	}
	if len(words) != mapping.SpanRegisters {
		s.markWriteDegraded()
		return NewError(ErrCodeInvalidGeometry, "runtime mapping encoded span does not match its identity", false)
	}

	baseRegister := int(mapping.Register)
	if baseRegister+len(words) > s.effectiveCapacityRegisters() {
		return fmt.Errorf("register range overflow: start=%d words=%d max=%d", baseRegister, len(words), s.effectiveCapacityRegisters())
	}

	for i, word := range words {
		offset := (baseRegister + i) * 2
		if err := s.bank.WriteWord(offset, word); err != nil {
			return fmt.Errorf("write register %d failed: %w", baseRegister+i, err)
		}
	}

	return nil
}

// WriteTagValue is retained for package callers that have not yet migrated to
// the explicit projection authority. HTTP handlers never call this method.
func (s *Service) WriteTagValue(ctx context.Context, tagID string, value interface{}) error {
	return s.WriteProjectedTagValue(ctx, tagID, value)
}

func (s *Service) markWriteDegraded() {
	s.mu.Lock()
	if s.state == shareStateRunning {
		s.state = HydrationStateFailed
	}
	s.mu.Unlock()
}
