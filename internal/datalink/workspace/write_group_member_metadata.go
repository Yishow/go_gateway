package workspace

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
)

// ErrWriteGroupMemberMetadataInvalid is the stable member_metadata_invalid
// failure for a member role or type that is unknown, explicitly empty or not
// a JSON string. It wraps ErrWriteGroupValidation so every existing mapping
// keeps treating it as a validation refusal.
var ErrWriteGroupMemberMetadataInvalid = fmt.Errorf("member_metadata_invalid: %w", ErrWriteGroupValidation)

// WriteGroupValueRole names how one instantaneous member is read. It is an
// explicit operator selection: nothing infers it from names or units.
type WriteGroupValueRole string

const (
	// WriteGroupValueRoleGauge uses the selected profile reduction for a gauge.
	WriteGroupValueRoleGauge WriteGroupValueRole = "gauge"
	// WriteGroupValueRoleCounterRaw keeps the raw counter value. It requires
	// the existing confirmed counter binding and is never averaged.
	WriteGroupValueRoleCounterRaw WriteGroupValueRole = "counter_raw"
	// WriteGroupValueRoleState keeps latest eligible values and never enters
	// an average, whatever the tag type is.
	WriteGroupValueRoleState WriteGroupValueRole = "state"
)

// memberMetadataUnsupported reports a member that carries new recording
// metadata. No runtime proof serves these shapes yet, so such a group stays
// savable but never starts or restores a writer.
func memberMetadataUnsupported(member *WriteGroupMember) bool {
	return member != nil && (member.Cadence != nil || member.Calculation != "" || member.ValueRole != "" || member.SourceValueType != "" || member.ResultValueType != "" || member.ExactRecordingSource != nil)
}

// UnmarshalJSON decodes a member and validates its calculation and recording
// metadata fields: an unknown value, or an explicit empty, null or non-string
// one, is rejected instead of silently degrading. Every other field keeps the
// plain decoding; schema.DataType itself is not changed globally.
func (m *WriteGroupMember) UnmarshalJSON(data []byte) error {
	type memberAlias WriteGroupMember
	var raw struct {
		Cadence         json.RawMessage `json:"cadence"`
		Calculation     json.RawMessage `json:"calculation"`
		ValueRole       json.RawMessage `json:"value_role"`
		SourceValueType json.RawMessage `json:"source_value_type"`
		ResultValueType json.RawMessage `json:"result_value_type"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if bytes.Equal(bytes.TrimSpace(raw.Cadence), []byte("null")) {
		return fmt.Errorf("%w: member cadence cannot be null", ErrWriteGroupValidation)
	}
	_, err := decodeMemberEnum(raw.Calculation, func(value string) error {
		return WriteGroupMemberCalculation(value).validate()
	})
	if err != nil {
		return err
	}
	valueRole, err := decodeMemberEnum(raw.ValueRole, func(value string) error {
		switch WriteGroupValueRole(value) {
		case WriteGroupValueRoleGauge, WriteGroupValueRoleCounterRaw, WriteGroupValueRoleState:
			return nil
		}
		return fmt.Errorf("unsupported value role %q", value)
	})
	if err != nil {
		return err
	}
	sourceType, err := decodeMemberEnum(raw.SourceValueType, func(value string) error {
		if !schema.DataType(value).IsValid() {
			return fmt.Errorf("unsupported source value type %q", value)
		}
		return nil
	})
	if err != nil {
		return err
	}
	resultType, err := decodeMemberEnum(raw.ResultValueType, func(value string) error {
		switch measurement.ExactType(value) {
		case measurement.ExactBool, measurement.ExactText, measurement.ExactInt64,
			measurement.ExactUint64, measurement.ExactDecimal, measurement.ExactFloat64:
			return nil
		}
		return fmt.Errorf("unsupported result value type %q", value)
	})
	if err != nil {
		return err
	}
	var member memberAlias
	if err := json.Unmarshal(data, &member); err != nil {
		return err
	}
	if valueRole != nil {
		member.ValueRole = WriteGroupValueRole(*valueRole)
	}
	if sourceType != nil {
		member.SourceValueType = schema.DataType(*sourceType)
	}
	if resultType != nil {
		member.ResultValueType = measurement.ExactType(*resultType)
	}
	*m = WriteGroupMember(member)
	return nil
}

// decodeMemberEnum validates one optional JSON string field and returns the
// decoded value, or nil when the key is absent. An explicit null, a non-string
// or an unknown value is a typed rejection instead of a silent default.
func decodeMemberEnum(raw json.RawMessage, known func(string) error) (*string, error) {
	if raw == nil {
		return nil, nil
	}
	wire := string(bytes.TrimSpace(raw))
	if wire == "null" || !strings.HasPrefix(wire, `"`) {
		return nil, fmt.Errorf("%w: recording metadata must be a JSON string", ErrWriteGroupMemberMetadataInvalid)
	}
	var value string
	if err := json.Unmarshal(raw, &value); err != nil {
		return nil, fmt.Errorf("%w: recording metadata must be a JSON string", ErrWriteGroupMemberMetadataInvalid)
	}
	if err := known(value); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrWriteGroupMemberMetadataInvalid, err)
	}
	return &value, nil
}

// validateMemberMetadataShape rejects the member metadata shapes that have no
// canonical meaning: roles exist only for instantaneous recording profiles.
// It runs before the database is touched and never infers a role.
func validateMemberMetadataShape(group *WriteGroup) error {
	if isMixedWriteGroup(group) {
		if err := ValidateMemberCalculationIntent(group); err != nil {
			return err
		}
		return validateMixedMemberCadences(group)
	}
	if group.RecordID != "" {
		return fmt.Errorf("%w: record_id requires mixed profile", ErrWriteGroupValidation)
	}
	if err := validateMixedMemberCadences(group); err != nil {
		return err
	}
	if err := validateCalculationProfileShape(group.RecordingProfile); err != nil {
		return err
	}
	for i, member := range group.Members {
		if member.Calculation != "" {
			return fmt.Errorf("%w: member %d calculation requires unsupported mixed recording runtime", ErrWriteGroupValidation, i)
		}
	}

	for i := range group.Members {
		member := &group.Members[i]
		if member.ValueRole == "" && member.SourceValueType == "" && member.ResultValueType == "" {
			continue
		}
		// New recording metadata never attaches to a group without a recording
		// profile, whatever field carries it.
		if group.RecordingProfile == nil {
			return fmt.Errorf("%w: member %d carries recording metadata without a recording profile", ErrWriteGroupValidation, i)
		}
		if member.ValueRole == "" {
			// A raw source or result type on a profile group is verified
			// against the real tag later, per the profile's own kind.
			continue
		}
		switch member.ValueRole {
		case WriteGroupValueRoleGauge, WriteGroupValueRoleCounterRaw, WriteGroupValueRoleState:
		default:
			return fmt.Errorf("%w: unsupported value role %q", ErrWriteGroupValidation, member.ValueRole)
		}
		if group.RecordingProfile.Kind != WriteGroupProfileInstantaneous {
			return fmt.Errorf("%w: member %d value roles record instantaneous profiles only", ErrWriteGroupValidation, i)
		}
	}
	return nil
}

// deriveMemberResultType maps a verified source type to the result type the
// runtime must write: latest keeps the tag's exact type, a mean of a numeric
// gauge yields decimal for exact integers and float64 for floats, and state or
// counter_raw members are never averaged.
func deriveMemberResultType(reduction WriteGroupSampleReduction, role WriteGroupValueRole, source schema.DataType) (measurement.ExactType, error) {
	exact, ok := measurement.ExactTypeForTag(source)
	if !ok {
		return "", fmt.Errorf("%w: the tag type %q cannot be verified for recording metadata", ErrWriteGroupValidation, source)
	}
	if reduction == WriteGroupSampleReductionMean && role == WriteGroupValueRoleGauge {
		switch exact {
		case measurement.ExactInt64, measurement.ExactUint64:
			return measurement.ExactDecimal, nil
		case measurement.ExactFloat64:
			return measurement.ExactFloat64, nil
		default:
			return "", fmt.Errorf("%w: a %s gauge cannot enter a mean", ErrWriteGroupValidation, exact)
		}
	}
	return exact, nil
}

// storedMemberEnum decodes a nullable stored metadata column: NULL is the
// legacy absence and reads empty, while a non-NULL empty or unknown value is
// corrupt data that fails loudly instead of degrading to a default.
func storedMemberEnum(value sql.NullString, known func(string) error) (string, error) {
	if !value.Valid {
		return "", nil
	}
	if value.String == "" {
		return "", fmt.Errorf("%w: a stored recording metadata value must not be empty", ErrWriteGroupValidation)
	}
	if err := known(value.String); err != nil {
		return "", fmt.Errorf("%w: %w", ErrWriteGroupValidation, err)
	}
	return value.String, nil
}

// sameStoredMetadata keeps a resend of the same recording intent idempotent:
// the server-derived types may be omitted by the caller and resolve to the
// stored facts, while an explicit value must match them. A role compares by
// its effective value, so an omitted gauge replays a saved gauge.
func sameStoredMetadata(candidate, stored string) bool {
	return candidate == "" || candidate == stored
}
