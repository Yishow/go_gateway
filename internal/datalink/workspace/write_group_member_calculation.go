package workspace

import (
	"encoding/json"
	"fmt"
)

// WriteGroupMemberCalculation is the explicit calculation of one member in a
// mixed recording. It is not a source attribute or a legacy group reduction.
type WriteGroupMemberCalculation string

const (
	WriteGroupMemberCalculationLatest WriteGroupMemberCalculation = "latest"
	WriteGroupMemberCalculationMean   WriteGroupMemberCalculation = "mean"
	WriteGroupMemberCalculationDelta  WriteGroupMemberCalculation = "delta"
)

func (method WriteGroupMemberCalculation) validate() error {
	switch method {
	case WriteGroupMemberCalculationLatest, WriteGroupMemberCalculationMean, WriteGroupMemberCalculationDelta:
		return nil
	default:
		return fmt.Errorf("%w: unsupported member calculation %q", ErrWriteGroupValidation, method)
	}
}

// WriteGroupResultShape identifies persisted result semantics, not a runtime
// capability. Legacy raw-plus-provenance output has no mixed shape.
type WriteGroupResultShape string

// WriteGroupResultShapeMemberCalculationV1 stores each member's chosen result.
const WriteGroupResultShapeMemberCalculationV1 WriteGroupResultShape = "member-calculation-v1"

func validateCalculationProfileShape(p *WriteGroupRecordingProfile) error {
	if p == nil {
		return nil
	}
	if p.Kind != WriteGroupProfileMixed {
		if p.ResultShape != "" {
			return fmt.Errorf("%w: result shape requires mixed profile", ErrWriteGroupValidation)
		}
		return nil
	}
	if p.ResultShape != WriteGroupResultShapeMemberCalculationV1 {
		return fmt.Errorf("%w: mixed profile requires member-calculation-v1 result shape", ErrWriteGroupValidation)
	}
	if p.SampleReduction != "" || p.Interval != 0 || p.Unit != "" {
		return fmt.Errorf("%w: mixed calculations cannot use group sample reduction or interval", ErrWriteGroupValidation)
	}
	return nil
}

// InitializeMemberCalculation initializes a new draft member only; it must not
// be used to repair an omitted method in submitted mixed intent.
func InitializeMemberCalculation(member WriteGroupMember, groupDefault WriteGroupMemberCalculation) (WriteGroupMember, error) {
	if member.Calculation != "" {
		return member, member.Calculation.validate()
	}
	if groupDefault == "" {
		groupDefault = WriteGroupMemberCalculationLatest
	}
	if err := groupDefault.validate(); err != nil {
		return member, err
	}
	member.Calculation = groupDefault
	return member, nil
}

// ValidateMemberCalculationIntent checks only the typed intent boundary.
// Success is not applicability, persistence, readiness, or runtime support.
func ValidateMemberCalculationIntent(group *WriteGroup) error {
	if group == nil {
		return fmt.Errorf("%w: recording group is required", ErrWriteGroupValidation)
	}
	if err := validateCalculationProfileShape(group.RecordingProfile); err != nil {
		return err
	}
	mixed := group.RecordingProfile != nil && group.RecordingProfile.Kind == WriteGroupProfileMixed
	if mixed && len(group.Members) == 0 {
		return fmt.Errorf("%w: mixed recording requires members", ErrWriteGroupValidation)
	}
	for i, member := range group.Members {
		if !mixed {
			if member.Calculation != "" {
				return fmt.Errorf("%w: member %d calculation requires mixed profile", ErrWriteGroupValidation, i)
			}
			continue
		}
		if err := member.Calculation.validate(); err != nil {
			return fmt.Errorf("member %d: %w", i, err)
		}
	}
	return nil
}

// MarshalJSON preserves legacy profile representation and emits mixed intent
// without legacy group cadence fields, so its own wire form round trips.
func (p WriteGroupRecordingProfile) MarshalJSON() ([]byte, error) {
	type profileAlias WriteGroupRecordingProfile
	if p.Kind != WriteGroupProfileMixed {
		return json.Marshal(profileAlias(p))
	}
	if err := validateCalculationProfileShape(&p); err != nil {
		return nil, err
	}
	return json.Marshal(struct {
		Kind        WriteGroupProfileKind `json:"kind"`
		ResultShape WriteGroupResultShape `json:"result_shape"`
	}{p.Kind, p.ResultShape})
}
