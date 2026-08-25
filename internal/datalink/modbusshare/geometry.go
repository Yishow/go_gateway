package modbusshare

import (
	"fmt"
	"sort"

	"go-gateway/internal/datalink/schema"
)

const (
	// HumanHoldingRegisterBase is the canonical starting address for Modbus holding registers.
	HumanHoldingRegisterBase uint32 = 40001
	// MaxRegisterIndex is the maximum allowable 0-based register index in the
	// 64 KiB backing bank (32 Ki holding registers).
	MaxRegisterIndex uint16 = 32767
)

// HumanToZeroBased converts a human holding register address (e.g. 40001) to a 0-based index (0).
func HumanToZeroBased(human uint32) (uint16, error) {
	if human < HumanHoldingRegisterBase {
		return 0, fmt.Errorf("human register %d is below base %d", human, HumanHoldingRegisterBase)
	}
	offset := human - HumanHoldingRegisterBase
	if offset > uint32(MaxRegisterIndex) {
		return 0, fmt.Errorf("human register %d exceeds maximum allowed register address", human)
	}
	return uint16(offset), nil
}

// ZeroBasedToHuman converts a 0-based register index to a human holding register address.
func ZeroBasedToHuman(zero uint16) uint32 {
	return HumanHoldingRegisterBase + uint32(zero)
}

// DataTypeSpan returns the number of 16-bit Modbus registers occupied by the datatype.
func DataTypeSpan(dt schema.DataType) int {
	switch dt {
	case schema.DataTypeInt32, schema.DataTypeUint32, schema.DataTypeFloat32:
		return 2
	case schema.DataTypeInt64, schema.DataTypeUint64, schema.DataTypeFloat64:
		return 4
	default:
		return 1
	}
}

// DataTypeBytes returns the byte count occupied by the datatype.
func DataTypeBytes(dt schema.DataType) int {
	return DataTypeSpan(dt) * 2
}

// ValidateDesiredMappings validates datatype span, stride, capacity, and pairwise collisions for a set of desired mappings.
func ValidateDesiredMappings(mappings []DesiredMapping, capacityRegisters int) error {
	if capacityRegisters <= 0 {
		return &Error{Code: ErrCodeCapacityExceeded, Message: "capacity must be greater than zero", Retryable: false}
	}
	if capacityRegisters > int(MaxRegisterIndex)+1 {
		return &Error{Code: ErrCodeCapacityExceeded, Message: "capacity exceeds register address space", Retryable: false}
	}

	type rangeInterval struct {
		tagID string
		start int
		end   int
	}
	intervals := make([]rangeInterval, 0, len(mappings))
	seenTags := make(map[string]struct{}, len(mappings))

	for _, m := range mappings {
		if m.CapacityRegisters != 0 && m.CapacityRegisters != capacityRegisters {
			return &Error{Code: ErrCodeCapacityExceeded, Message: fmt.Sprintf("mapping tag %s capacity %d does not match configured capacity %d", m.TagID, m.CapacityRegisters, capacityRegisters), Retryable: false}
		}
		if m.TagID == "" {
			return &Error{Code: ErrCodeInvalidGeometry, Message: "mapping tag id is required", Retryable: false}
		}
		if _, exists := seenTags[m.TagID]; exists {
			return &Error{Code: ErrCodeInvalidGeometry, Message: fmt.Sprintf("mapping tag %s is duplicated", m.TagID), Retryable: false}
		}
		seenTags[m.TagID] = struct{}{}
		if !m.DataType.IsValid() || m.DataType == schema.DataTypeString {
			return &Error{Code: ErrCodeInvalidGeometry, Message: fmt.Sprintf("mapping tag %s has unsupported datatype %q", m.TagID, m.DataType), Retryable: false}
		}
		span := DataTypeSpan(m.DataType)
		if m.SpanRegisters != span {
			return &Error{
				Code:      ErrCodeInvalidGeometry,
				Message:   fmt.Sprintf("mapping tag %s span %d does not match datatype span %d", m.TagID, m.SpanRegisters, span),
				Retryable: false,
			}
		}
		if m.ShareStartRegister != 0 && m.ShareStartRegister != ZeroBasedToHuman(m.ZeroBasedRegister) {
			return &Error{
				Code:      ErrCodeInvalidGeometry,
				Message:   fmt.Sprintf("mapping tag %s human and zero-based registers disagree", m.TagID),
				Retryable: false,
			}
		}

		stride := m.StrideRegisters
		if stride < span {
			return &Error{
				Code:      ErrCodeInvalidGeometry,
				Message:   fmt.Sprintf("mapping tag %s stride %d is smaller than datatype span %d", m.TagID, stride, span),
				Retryable: false,
			}
		}

		start := int(m.ZeroBasedRegister)
		end := start + span

		if end > capacityRegisters {
			return &Error{
				Code:      ErrCodeCapacityExceeded,
				Message:   fmt.Sprintf("mapping tag %s register span [%d, %d) exceeds capacity %d", m.TagID, start, end, capacityRegisters),
				Retryable: false,
			}
		}

		intervals = append(intervals, rangeInterval{
			tagID: m.TagID,
			start: start,
			end:   end,
		})
	}

	// Sort intervals by start register
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i].start < intervals[j].start
	})

	// Pairwise overlap check
	for i := 0; i+1 < len(intervals); i++ {
		j := i + 1
		if intervals[j].start < intervals[i].end {
			return &Error{
				Code: ErrCodeRangeCollision,
				Message: fmt.Sprintf("range collision between tag %s [%d, %d) and tag %s [%d, %d)",
					intervals[i].tagID, intervals[i].start, intervals[i].end,
					intervals[j].tagID, intervals[j].start, intervals[j].end,
				),
				Retryable: false,
			}
		}
	}

	return nil
}

// ValidateDesiredMappingsAgainstExisting checks a desired workspace projection
// against already projected spans owned by other workspaces. Existing rows are
// intentionally read-only inputs; they are never replaced by reconciliation.
func ValidateDesiredMappingsAgainstExisting(mappings []DesiredMapping, existing []TagMirrorMapping, capacityRegisters int) error {
	if err := ValidateDesiredMappings(mappings, capacityRegisters); err != nil {
		return err
	}
	for _, desired := range mappings {
		desiredStart := int(desired.ZeroBasedRegister)
		desiredEnd := desiredStart + DataTypeSpan(desired.DataType)
		for _, current := range existing {
			span := current.SpanRegisters
			if span <= 0 {
				span = DataTypeSpan(current.DataType)
			}
			currentStart := int(current.ZeroBasedRegister)
			currentEnd := currentStart + span
			if desiredStart < currentEnd && currentStart < desiredEnd {
				return &Error{
					Code:      ErrCodeRangeCollision,
					Message:   fmt.Sprintf("range collision between tag %s [%d, %d) and tag %s [%d, %d)", desired.TagID, desiredStart, desiredEnd, current.TagID, currentStart, currentEnd),
					Retryable: false,
				}
			}
		}
	}
	return nil
}
