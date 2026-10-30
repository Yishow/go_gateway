package workspace

import (
	"fmt"
	"strings"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
)

func applyMixedRecordingMetadata(member *WriteGroupMember, index int, tagType schema.DataType) error {
	switch member.ValueRole {
	case "", WriteGroupValueRoleGauge, WriteGroupValueRoleCounterRaw, WriteGroupValueRoleState:
	default:
		return fmt.Errorf("%w: unsupported mixed value role", ErrWriteGroupValidation)
	}
	if err := member.Calculation.validate(); err != nil {
		return err
	}
	exact, ok := measurement.ExactTypeForTag(tagType)
	if !ok {
		return fmt.Errorf("%w: member %d source type is unsupported", ErrWriteGroupValidation, index)
	}
	if member.SourceValueType != "" && member.SourceValueType != tagType {
		return fmt.Errorf("%w: member %d source type is stale", ErrWriteGroupValidation, index)
	}
	if member.ValueRole == WriteGroupValueRoleCounterRaw {
		if exact == measurement.ExactBool || exact == measurement.ExactText || member.MeasurementID == nil || strings.TrimSpace(*member.MeasurementID) == "" {
			return fmt.Errorf("%w: counter_raw requires numeric source and confirmed binding", ErrWriteGroupValidation)
		}
	}
	result := exact
	switch member.Calculation {
	case WriteGroupMemberCalculationMean:
		if member.ValueRole != "" && member.ValueRole != WriteGroupValueRoleGauge {
			return fmt.Errorf("%w: mean requires a gauge", ErrWriteGroupValidation)
		}
		var err error
		result, err = deriveMemberResultType(WriteGroupSampleReductionMean, WriteGroupValueRoleGauge, tagType)
		if err != nil {
			return err
		}
	case WriteGroupMemberCalculationDelta:
		if exact == measurement.ExactBool || exact == measurement.ExactText {
			return fmt.Errorf("%w: delta requires numeric source", ErrWriteGroupValidation)
		}
		if member.ValueRole != "" && member.ValueRole != WriteGroupValueRoleCounterRaw {
			return fmt.Errorf("%w: delta requires counter semantics", ErrWriteGroupValidation)
		}
		if member.MeasurementID == nil || strings.TrimSpace(*member.MeasurementID) == "" {
			return fmt.Errorf("%w: delta requires confirmed counter binding", ErrWriteGroupValidation)
		}
		if exact == measurement.ExactInt64 || exact == measurement.ExactUint64 {
			result = measurement.ExactDecimal
		}
	}
	if member.ExactRecordingSource != nil {
		result = measurement.ExactDecimal
	}
	if member.ResultValueType != "" && member.ResultValueType != result {
		return fmt.Errorf("%w: member %d result type is stale", ErrWriteGroupValidation, index)
	}
	member.SourceValueType, member.ResultValueType = tagType, result
	return nil
}
