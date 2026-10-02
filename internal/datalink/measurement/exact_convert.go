package measurement

import (
	"fmt"
	"math"

	"go-gateway/internal/datalink/schema"
)

// ExactTypeForTag maps a persisted tag data type to the exact value type the
// typed path carries. No tag type implies decimal, so decimals are never
// guessed from names or units.
func ExactTypeForTag(dataType schema.DataType) (ExactType, bool) {
	switch dataType {
	case schema.DataTypeBool:
		return ExactBool, true
	case schema.DataTypeInt16, schema.DataTypeInt32, schema.DataTypeInt64:
		return ExactInt64, true
	case schema.DataTypeUint16, schema.DataTypeUint32, schema.DataTypeUint64:
		return ExactUint64, true
	case schema.DataTypeFloat32, schema.DataTypeFloat64:
		return ExactFloat64, true
	case schema.DataTypeString:
		return ExactText, true
	}
	return "", false
}

// ExactFromGo converts a pipeline value into the exact type declared by the
// tag. It never coerces across kinds: an int for a float tag, a float for an
// integer tag, NaN or an out-of-range integer is an error, not a changed value.
func ExactFromGo(expected ExactType, value any) (ExactValue, error) {
	mismatch := fmt.Errorf("%w: value does not match declared %s type", ErrExactValueInvalid, expected)
	switch expected {
	case ExactBool:
		if v, ok := value.(bool); ok {
			return NewBool(v), nil
		}
	case ExactText:
		if v, ok := value.(string); ok {
			return NewText(v), nil
		}
	case ExactInt64:
		if signed, unsigned, kind := integerKind(value); kind == kindSigned {
			return NewInt64(signed), nil
		} else if kind == kindUnsigned && unsigned <= math.MaxInt64 {
			return NewInt64(int64(unsigned)), nil
		} else if kind != kindNone {
			return ExactValue{}, fmt.Errorf("%w: integer out of int64 range", ErrExactValueInvalid)
		}
	case ExactUint64:
		if signed, unsigned, kind := integerKind(value); kind == kindUnsigned {
			return NewUint64(unsigned), nil
		} else if kind == kindSigned && signed >= 0 {
			return NewUint64(uint64(signed)), nil
		} else if kind != kindNone {
			return ExactValue{}, fmt.Errorf("%w: integer out of uint64 range", ErrExactValueInvalid)
		}
	case ExactFloat64:
		switch v := value.(type) {
		case float64:
			return NewFloat64(v)
		case float32:
			return NewFloat64(float64(v))
		}
	}
	return ExactValue{}, mismatch
}

type integerValueKind int

const (
	kindNone integerValueKind = iota
	kindSigned
	kindUnsigned
)

func integerKind(value any) (signed int64, unsigned uint64, kind integerValueKind) {
	switch v := value.(type) {
	case int:
		return int64(v), 0, kindSigned
	case int8:
		return int64(v), 0, kindSigned
	case int16:
		return int64(v), 0, kindSigned
	case int32:
		return int64(v), 0, kindSigned
	case int64:
		return v, 0, kindSigned
	case uint:
		return 0, uint64(v), kindUnsigned
	case uint8:
		return 0, uint64(v), kindUnsigned
	case uint16:
		return 0, uint64(v), kindUnsigned
	case uint32:
		return 0, uint64(v), kindUnsigned
	case uint64:
		return 0, v, kindUnsigned
	}
	return 0, 0, kindNone
}
