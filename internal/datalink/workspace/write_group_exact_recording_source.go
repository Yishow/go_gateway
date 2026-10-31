package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
)

// ExactRecordingSource is the shared immutable scaled-source facts contract.
type ExactRecordingSource = counterproof.ExactRecordingSource

// ValidateExactRecordingSource validates the member's frozen exact recording
// source proof against the confirmed counter contract and tag data type.
func ValidateExactRecordingSource(member WriteGroupMember, tagType schema.DataType) error {
	proof := member.ExactRecordingSource
	if proof == nil {
		return fmt.Errorf("%w: exact recording source proof is required", ErrWriteGroupValidation)
	}
	if strings.TrimSpace(proof.MeasurementID) == "" {
		return fmt.Errorf("%w: measurement_id is required in exact recording source", ErrWriteGroupValidation)
	}
	if member.MeasurementID == nil || strings.TrimSpace(*member.MeasurementID) == "" || *member.MeasurementID != proof.MeasurementID {
		return fmt.Errorf("%w: measurement_id mismatch between member and exact recording source", ErrWriteGroupValidation)
	}
	if proof.OriginalDataType != tagType {
		return fmt.Errorf("%w: original data type %q does not match tag type %q", ErrWriteGroupValidation, proof.OriginalDataType, tagType)
	}
	exactKind, ok := measurement.ExactTypeForTag(tagType)
	if !ok || (exactKind != measurement.ExactInt64 && exactKind != measurement.ExactUint64) {
		return fmt.Errorf("%w: exact recording source requires an integer tag type, got %q", ErrWriteGroupValidation, tagType)
	}
	switch member.ValueRole {
	case WriteGroupValueRoleGauge, WriteGroupValueRoleState:
		return fmt.Errorf("%w: value role %q cannot carry exact recording source", ErrWriteGroupValidation, member.ValueRole)
	case WriteGroupValueRoleCounterRaw, "":
	default:
		return fmt.Errorf("%w: unsupported value role %q for exact recording source", ErrWriteGroupValidation, member.ValueRole)
	}
	if strings.TrimSpace(proof.Scale) == "" {
		return fmt.Errorf("%w: scale cannot be empty in exact recording source", ErrWriteGroupValidation)
	}
	if _, err := measurement.ScaleExactInteger(measurement.NewInt64(1), proof.Scale); err != nil {
		return fmt.Errorf("%w: invalid exact recording scale %q: %w", ErrWriteGroupValidation, proof.Scale, err)
	}
	def := proof.Definition()
	if confirmation := measurement.ConfirmCounter(def); confirmation != measurement.CounterConfirmed {
		return fmt.Errorf("%w: counter semantics not confirmed (%s)", ErrWriteGroupValidation, confirmation)
	}
	if member.SourceValueType != proof.OriginalDataType {
		return fmt.Errorf("%w: member source_value_type %q does not match proof %q", ErrWriteGroupValidation, member.SourceValueType, proof.OriginalDataType)
	}
	if member.ResultValueType != measurement.ExactDecimal {
		return fmt.Errorf("%w: member result_value_type %q must be %s", ErrWriteGroupValidation, member.ResultValueType, measurement.ExactDecimal)
	}
	return nil
}

// LongMemberStoredTypeForMember returns measurement.ExactDecimal only when the
// member carries a valid exact recording source proof. An absent proof keeps
// the legacy type contract; an invalid present proof cannot fall back to it.
func LongMemberStoredTypeForMember(reduction WriteGroupSampleReduction, member WriteGroupMember, tagType schema.DataType) (measurement.ExactType, bool) {
	if member.ExactRecordingSource != nil {
		if ValidateExactRecordingSource(member, tagType) != nil {
			return "", false
		}
		return measurement.ExactDecimal, true
	}
	return LongMemberStoredType(reduction, member.ValueRole, tagType)
}

func cloneExactRecordingSource(s *ExactRecordingSource) *ExactRecordingSource {
	if s == nil {
		return nil
	}
	cloned := *s
	cloned.CounterPolicy = cloneExactCounterPolicy(s.CounterPolicy)
	return &cloned
}

func cloneExactCounterPolicy(policy *measurement.CounterPolicy) *measurement.CounterPolicy {
	if policy == nil {
		return nil
	}
	cloned := *policy
	if policy.ResetThreshold != nil {
		value := *policy.ResetThreshold
		cloned.ResetThreshold = &value
	}
	if policy.MaxRollOver != nil {
		value := *policy.MaxRollOver
		cloned.MaxRollOver = &value
	}
	return &cloned
}

func sameExactRecordingSource(left, right *ExactRecordingSource) bool {
	return reflect.DeepEqual(left, right)
}

func storedExactRecordingSource(value sql.NullString) (*ExactRecordingSource, error) {
	if !value.Valid {
		return nil, nil
	}
	if strings.TrimSpace(value.String) == "" {
		return nil, fmt.Errorf("%w: stored exact recording source must not be empty", ErrWriteGroupValidation)
	}
	var source ExactRecordingSource
	if err := json.Unmarshal([]byte(value.String), &source); err != nil {
		return nil, fmt.Errorf("%w: decode stored exact recording source: %w", ErrWriteGroupValidation, err)
	}
	id := source.MeasurementID
	if err := ValidateExactRecordingSource(WriteGroupMember{
		MeasurementID: &id, ExactRecordingSource: &source,
		SourceValueType: source.OriginalDataType, ResultValueType: measurement.ExactDecimal,
	}, source.OriginalDataType); err != nil {
		return nil, err
	}
	return &source, nil
}

func (r *SQLWriteGroupRepository) deriveAndValidateExactRecordingSourceInTx(
	ctx context.Context,
	runner writeGroupSQLRunner,
	group *WriteGroup,
	member *WriteGroupMember,
	index int,
	tagType schema.DataType,
	definitions map[string]*measurement.MeasurementDefinition,
) error {
	if member.MeasurementID == nil || strings.TrimSpace(*member.MeasurementID) == "" {
		if member.ExactRecordingSource != nil {
			return fmt.Errorf("%w: member %d carries exact recording source without measurement_id", ErrWriteGroupValidation, index)
		}
		return nil
	}
	tx, ok := runner.(*sql.Tx)
	if !ok {
		return fmt.Errorf("derive exact recording source: validation runner is not a transaction")
	}
	def, err := measurement.NewSQLRepository(r.db).GetByIDInTx(ctx, tx, *member.MeasurementID)
	if err != nil {
		return fmt.Errorf("read write-group measurement: %w", err)
	}
	if def == nil {
		return fmt.Errorf("read write-group measurement: the counter binding is unavailable")
	}
	definitions[def.ID] = def
	hasScale := strings.TrimSpace(def.ExactRecordingScale) != ""
	if !hasScale {
		if member.ExactRecordingSource != nil {
			return fmt.Errorf("%w: member %d carries exact recording source but measurement has no scale", ErrWriteGroupValidation, index)
		}
		return nil
	}

	if !group.isLongLayout() || group.Destination.StorageStrategy != WriteGroupStorageStrategyManaged {
		return fmt.Errorf("%w: member %d exact recording scale requires managed_long_v1 plain layout", ErrWriteGroupValidation, index)
	}
	cumulative := group.RecordingProfile != nil && group.RecordingProfile.Kind == WriteGroupProfileCumulative && member.ValueRole == ""
	counterRaw := group.RecordingProfile != nil && group.RecordingProfile.Kind == WriteGroupProfileInstantaneous && member.ValueRole == WriteGroupValueRoleCounterRaw
	// Mixed scale applies only to the delta member's existing confirmed counter
	// binding. It does not make latest/mean consume counter scaling or a baseline.
	mixedDelta := isMixedWriteGroup(group) && member.Calculation == WriteGroupMemberCalculationDelta &&
		(member.ValueRole == "" || member.ValueRole == WriteGroupValueRoleCounterRaw) &&
		(group.StorageMode == "" || group.StorageMode == WriteGroupStorageModePlain)
	if !cumulative && !counterRaw && !mixedDelta {
		return fmt.Errorf("%w: member %d exact recording scale requires cumulative, instantaneous counter_raw, or managed-long mixed delta", ErrWriteGroupValidation, index)
	}

	exactKind, ok := measurement.ExactTypeForTag(tagType)
	if !ok || (exactKind != measurement.ExactInt64 && exactKind != measurement.ExactUint64) {
		return fmt.Errorf("%w: member %d exact recording scale requires integer tag, got %q", ErrWriteGroupValidation, index, tagType)
	}

	if confirmation := measurement.ConfirmCounter(*def); confirmation != measurement.CounterConfirmed {
		return fmt.Errorf("%w: member %d counter semantics not confirmed (%s)", ErrWriteGroupValidation, index, confirmation)
	}
	if strings.ToLower(strings.TrimSpace(def.NumericEncoding)) != string(tagType) {
		return fmt.Errorf("%w: member %d measurement encoding does not match the integer tag", ErrWriteGroupValidation, index)
	}
	if member.SourceValueType != "" && member.SourceValueType != tagType ||
		member.ResultValueType != "" && member.ResultValueType != measurement.ExactDecimal {
		return fmt.Errorf("%w: member %d exact recording metadata conflicts with the saved definition", ErrWriteGroupValidation, index)
	}
	authoritativeProof := &ExactRecordingSource{
		MeasurementID:         def.ID,
		DefinitionRevision:    def.DefinitionRevision,
		SourceBindingRevision: def.SourceBindingRevision,
		SeriesEpoch:           def.SeriesEpoch,
		OriginalDataType:      tagType,
		Scale:                 def.ExactRecordingScale,
		Unit:                  strings.TrimSpace(def.Unit),
		SemanticKind:          def.SemanticKind,
		CounterPolicy:         cloneExactCounterPolicy(def.CounterPolicy),
	}

	if member.ExactRecordingSource != nil {
		if !sameExactRecordingSource(member.ExactRecordingSource, authoritativeProof) {
			return fmt.Errorf("%w: member %d exact recording source proof is stale or does not match authoritative definition", ErrWriteGroupValidation, index)
		}
	} else {
		member.ExactRecordingSource = authoritativeProof
	}

	member.SourceValueType = tagType
	member.ResultValueType = measurement.ExactDecimal
	return ValidateExactRecordingSource(*member, tagType)
}
