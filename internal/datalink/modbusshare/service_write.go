package modbusshare

import (
	"context"
	"fmt"
)

// WriteTagValue writes tag value to mapped Modbus registers.
func (s *Service) WriteTagValue(ctx context.Context, tagID string, value interface{}) error {
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

	words, err := encodeToWords(t.DataType, value)
	if err != nil {
		return err
	}

	for i, word := range words {
		offset := int(mapping.Register+uint16(i)) * 2
		if err := s.bank.WriteWord(offset, word); err != nil {
			return fmt.Errorf("write register %d failed: %w", int(mapping.Register)+i, err)
		}
	}

	return nil
}
