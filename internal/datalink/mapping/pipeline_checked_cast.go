package mapping

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"go-gateway/internal/datalink/schema"
)

// numericValue preserves integer digits and the exact value of finite floats.
func numericValue(input any) (*big.Rat, error) {
	if input == nil {
		return nil, fmt.Errorf("invalid numeric value")
	}
	value := reflect.ValueOf(input)
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return new(big.Rat).SetInt64(value.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return new(big.Rat).SetInt(new(big.Int).SetUint64(value.Uint())), nil
	case reflect.Float32, reflect.Float64:
		f := value.Float()
		if math.IsNaN(f) || math.IsInf(f, 0) {
			return nil, fmt.Errorf("non-finite numeric value")
		}
		return new(big.Rat).SetFloat64(f), nil
	case reflect.Bool:
		if value.Bool() {
			return new(big.Rat).SetInt64(1), nil
		}
		return new(big.Rat), nil
	}
	var text string
	switch v := input.(type) {
	case string:
		text = v
	case json.Number:
		text = v.String()
	default:
		return nil, fmt.Errorf("invalid numeric value")
	}
	// ParseFloat validates numeric syntax; Rat retains the original decimal digits.
	f, err := strconv.ParseFloat(text, 64)
	if err != nil || math.IsNaN(f) || math.IsInf(f, 0) || strings.ContainsAny(text, "/xXpP") {
		return nil, fmt.Errorf("invalid numeric text")
	}
	r, ok := new(big.Rat).SetString(text)
	if !ok {
		return nil, fmt.Errorf("invalid numeric text")
	}
	return r, nil
}

func checkedFloat64(input any) (float64, error) {
	r, err := numericValue(input)
	if err != nil {
		return 0, err
	}
	f, exact := r.Float64()
	if math.IsNaN(f) || math.IsInf(f, 0) || (f == 0 && r.Sign() != 0) || (r.IsInt() && !exact) {
		return 0, fmt.Errorf("numeric conversion loses precision")
	}
	return f, nil
}

func checkedCast(input any, target schema.DataType) (any, error) {
	if target == schema.DataTypeString {
		return fmt.Sprint(input), nil
	}
	if target == schema.DataTypeBool {
		if v, ok := input.(string); ok {
			switch v {
			case "true", "1", "on", "yes":
				return true, nil
			case "false", "0", "off", "no":
				return false, nil
			}
			return nil, fmt.Errorf("invalid boolean text")
		}
		r, err := numericValue(input)
		if err != nil {
			return nil, err
		}
		return r.Sign() != 0, nil
	}
	r, err := numericValue(input)
	if err != nil {
		return nil, err
	}
	switch target {
	case schema.DataTypeFloat32, schema.DataTypeFloat64:
		f, err := checkedFloat64(input)
		if err != nil {
			return nil, err
		}
		if target == schema.DataTypeFloat64 {
			return f, nil
		}
		narrowed := float32(f)
		if math.IsInf(float64(narrowed), 0) || (narrowed == 0 && f != 0) || (r.IsInt() && new(big.Rat).SetFloat64(float64(narrowed)).Cmp(r) != 0) {
			return nil, fmt.Errorf("numeric conversion loses precision")
		}
		return narrowed, nil
	case schema.DataTypeInt16, schema.DataTypeInt32, schema.DataTypeInt64:
		if !r.IsInt() || !r.Num().IsInt64() {
			return nil, fmt.Errorf("integer conversion out of range or fractional")
		}
		n := r.Num().Int64()
		switch target {
		case schema.DataTypeInt16:
			if n < math.MinInt16 || n > math.MaxInt16 {
				return nil, fmt.Errorf("int16 out of range")
			}
			return int16(n), nil
		case schema.DataTypeInt32:
			if n < math.MinInt32 || n > math.MaxInt32 {
				return nil, fmt.Errorf("int32 out of range")
			}
			return int32(n), nil
		default:
			return n, nil
		}
	case schema.DataTypeUint16, schema.DataTypeUint32, schema.DataTypeUint64:
		if !r.IsInt() || !r.Num().IsUint64() {
			return nil, fmt.Errorf("unsigned conversion out of range or fractional")
		}
		n := r.Num().Uint64()
		switch target {
		case schema.DataTypeUint16:
			if n > math.MaxUint16 {
				return nil, fmt.Errorf("uint16 out of range")
			}
			return uint16(n), nil
		case schema.DataTypeUint32:
			if n > math.MaxUint32 {
				return nil, fmt.Errorf("uint32 out of range")
			}
			return uint32(n), nil
		default:
			return n, nil
		}
	default:
		return nil, fmt.Errorf("invalid cast target")
	}
}
