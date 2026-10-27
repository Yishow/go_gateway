package measurement

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
)

// ErrExactRecordingScaleInvalid marks scale text that is not a positive finite
// plain decimal or that is attached to an unconfirmed/non-integer source.
var ErrExactRecordingScaleInvalid = errors.New("exact recording scale invalid")

func parseExactRecordingScale(text string) (ExactValue, error) {
	if text == "" || len(text) > maxExactPolicyText {
		return ExactValue{}, fmt.Errorf("%w: scale must be a positive finite decimal", ErrExactRecordingScaleInvalid)
	}
	value, err := NewDecimal(text)
	if err != nil {
		return ExactValue{}, fmt.Errorf("%w: %w", ErrExactRecordingScaleInvalid, err)
	}
	rat, ok := value.Rat()
	if !ok || rat.Sign() <= 0 {
		return ExactValue{}, fmt.Errorf("%w: scale must be a positive finite decimal", ErrExactRecordingScaleInvalid)
	}
	return value, nil
}

func sameExactRecordingScale(left, right string) bool {
	if left == "" || right == "" {
		return left == "" && right == ""
	}
	leftValue, leftErr := parseExactRecordingScale(left)
	rightValue, rightErr := parseExactRecordingScale(right)
	return leftErr == nil && rightErr == nil && leftValue.Equal(rightValue)
}

func (m *MeasurementDefinition) validateExactRecordingScale() error {
	if m.ExactRecordingScale == "" {
		return nil
	}
	if _, err := parseExactRecordingScale(m.ExactRecordingScale); err != nil {
		return err
	}
	switch m.SemanticKind {
	case SemanticKindCounter:
		if m.CounterPolicy != nil && m.CounterPolicy.AllowNegative {
			return fmt.Errorf("%w: monotonic counter cannot allow negative", ErrExactRecordingScaleInvalid)
		}
	case SemanticKindSignedCounter:
		if m.CounterPolicy == nil || !m.CounterPolicy.AllowNegative {
			return fmt.Errorf("%w: signed counter needs allow_negative", ErrExactRecordingScaleInvalid)
		}
	default:
		return fmt.Errorf("%w: scale requires counter or signed_counter", ErrExactRecordingScaleInvalid)
	}
	if !isIntegerNumericEncoding(m.NumericEncoding) {
		return fmt.Errorf("%w: scale requires an integer numeric_encoding", ErrExactRecordingScaleInvalid)
	}
	return nil
}

func isIntegerNumericEncoding(encoding string) bool {
	switch strings.ToLower(strings.TrimSpace(encoding)) {
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64", "integer":
		return true
	}
	return false
}

// ScaleExactInteger multiplies an exact integer by a confirmed decimal scale.
// It never converts through float64 and renders exactly the scale's fractional
// digits (for example 120004 at scale 0.1 is 12000.4, 120000 is 12000.0).
func ScaleExactInteger(value ExactValue, scale string) (ExactValue, error) {
	if value.Type() != ExactInt64 && value.Type() != ExactUint64 {
		return ExactValue{}, fmt.Errorf("%w: scale requires an integer value", ErrExactValueInvalid)
	}
	parsedScale, err := parseExactRecordingScale(scale)
	if err != nil {
		return ExactValue{}, err
	}
	integer, ok := scaledOf(value)
	if !ok {
		return ExactValue{}, fmt.Errorf("%w: integer value is invalid", ErrExactValueInvalid)
	}
	multiplier, ok := scaledOf(parsedScale)
	if !ok {
		return ExactValue{}, fmt.Errorf("%w: scale is invalid", ErrExactRecordingScaleInvalid)
	}
	product := new(big.Int).Mul(integer.unscaled, multiplier.unscaled)
	result, err := NewDecimal(formatScaled(product, multiplier.scale))
	if err != nil {
		return ExactValue{}, fmt.Errorf("%w: scaled result is invalid", ErrExactValueInvalid)
	}
	return result, nil
}
