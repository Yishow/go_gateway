package measurement

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// ExactType names the value types that keep exact digits across JSON and SQL.
type ExactType string

const (
	ExactBool    ExactType = "bool"
	ExactText    ExactType = "text"
	ExactInt64   ExactType = "int64"
	ExactUint64  ExactType = "uint64"
	ExactDecimal ExactType = "decimal"
	ExactFloat64 ExactType = "float64"
)

const (
	exactEncodingJSON          = "json"
	exactEncodingDecimalString = "decimal-string"
)

// ErrExactValueInvalid marks a value that cannot be represented exactly.
var ErrExactValueInvalid = errors.New("exact value invalid")

// ExactValue is a typed value whose digits survive JSON and SQL unchanged.
// The zero value is unset and never marshals as a valid value.
type ExactValue struct {
	typ ExactType
	val any
}

// Decimal is a finite base-10 string validated by NewDecimal.
type Decimal string

func NewBool(v bool) ExactValue     { return ExactValue{typ: ExactBool, val: v} }
func NewText(v string) ExactValue   { return ExactValue{typ: ExactText, val: v} }
func NewInt64(v int64) ExactValue   { return ExactValue{typ: ExactInt64, val: v} }
func NewUint64(v uint64) ExactValue { return ExactValue{typ: ExactUint64, val: v} }

// NewDecimal validates a plain decimal string (optional sign, digits, optional
// fraction). Exponents, NaN, Inf and whitespace are rejected.
func NewDecimal(digits string) (ExactValue, error) {
	if _, _, _, err := splitDecimal(digits); err != nil {
		return ExactValue{}, err
	}
	return ExactValue{typ: ExactDecimal, val: Decimal(digits)}, nil
}

// NewFloat64 accepts only finite values.
func NewFloat64(v float64) (ExactValue, error) {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return ExactValue{}, fmt.Errorf("%w: float64 must be finite", ErrExactValueInvalid)
	}
	return ExactValue{typ: ExactFloat64, val: v}, nil
}

func (v ExactValue) Type() ExactType { return v.typ }

// Value returns bool, string, int64, uint64, Decimal or float64 by Type; nil when unset.
func (v ExactValue) Value() any { return v.val }

// Equal compares type and exact value. Decimals compare numerically.
func (v ExactValue) Equal(other ExactValue) bool {
	if v.typ == "" || v.typ != other.typ {
		return false
	}
	if v.typ == ExactDecimal {
		left, lok := v.val.(Decimal)
		right, rok := other.val.(Decimal)
		if !lok || !rok {
			return false
		}
		a, aok := canonicalDecimal(string(left))
		b, bok := canonicalDecimal(string(right))
		return aok && bok && a == b
	}
	return v.val == other.val
}

// DecimalDigits returns the decimal digits for decimal-string typed values.
func (v ExactValue) DecimalDigits() (string, bool) {
	switch x := v.val.(type) {
	case Decimal:
		return string(x), true
	case int64:
		return strconv.FormatInt(x, 10), true
	case uint64:
		return strconv.FormatUint(x, 10), true
	}
	return "", false
}

type exactWire struct {
	Type     ExactType       `json:"type"`
	Encoding string          `json:"encoding"`
	Value    json.RawMessage `json:"value"`
}

func (v ExactValue) MarshalJSON() ([]byte, error) {
	wire := exactWire{Type: v.typ}
	var value any
	switch x := v.val.(type) {
	case bool:
		wire.Encoding, value = exactEncodingJSON, x
	case string:
		wire.Encoding, value = exactEncodingJSON, x
	case float64:
		if math.IsNaN(x) || math.IsInf(x, 0) {
			return nil, fmt.Errorf("%w: float64 must be finite", ErrExactValueInvalid)
		}
		wire.Encoding, value = exactEncodingJSON, x
	case int64, uint64, Decimal:
		digits, _ := v.DecimalDigits()
		wire.Encoding, value = exactEncodingDecimalString, digits
	default:
		return nil, fmt.Errorf("%w: value is unset", ErrExactValueInvalid)
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("%w: encode value", ErrExactValueInvalid)
	}
	wire.Value = raw
	return json.Marshal(wire)
}

func (v *ExactValue) UnmarshalJSON(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var wire exactWire
	if err := decoder.Decode(&wire); err != nil {
		return fmt.Errorf("%w: malformed envelope", ErrExactValueInvalid)
	}
	if _, err := decoder.Token(); err == nil {
		return fmt.Errorf("%w: trailing data", ErrExactValueInvalid)
	}
	decoded, err := decodeExactWire(wire)
	if err != nil {
		return err
	}
	*v = decoded
	return nil
}

func decodeExactWire(wire exactWire) (ExactValue, error) {
	invalid := func(reason string) (ExactValue, error) {
		return ExactValue{}, fmt.Errorf("%w: %s", ErrExactValueInvalid, reason)
	}
	wantEncoding := exactEncodingJSON
	switch wire.Type {
	case ExactInt64, ExactUint64, ExactDecimal:
		wantEncoding = exactEncodingDecimalString
	case ExactBool, ExactText, ExactFloat64:
	default:
		return invalid("unknown type")
	}
	if wire.Encoding != wantEncoding {
		return invalid("encoding does not match type")
	}
	if len(wire.Value) == 0 || bytes.Equal(bytes.TrimSpace(wire.Value), []byte("null")) {
		return invalid("value is missing")
	}
	switch wire.Type {
	case ExactBool:
		var b bool
		if err := json.Unmarshal(wire.Value, &b); err != nil {
			return invalid("bool value")
		}
		return NewBool(b), nil
	case ExactText:
		var s string
		if err := json.Unmarshal(wire.Value, &s); err != nil {
			return invalid("text value")
		}
		return NewText(s), nil
	case ExactFloat64:
		var f float64
		if err := json.Unmarshal(wire.Value, &f); err != nil {
			return invalid("float64 value")
		}
		return NewFloat64(f)
	}
	var digits string
	if err := json.Unmarshal(wire.Value, &digits); err != nil {
		return invalid("digits must be a string")
	}
	switch wire.Type {
	case ExactInt64:
		n, err := strconv.ParseInt(digits, 10, 64)
		if err != nil {
			return invalid("int64 out of range or malformed")
		}
		return NewInt64(n), nil
	case ExactUint64:
		n, err := strconv.ParseUint(digits, 10, 64)
		if err != nil {
			return invalid("uint64 out of range or malformed")
		}
		return NewUint64(n), nil
	}
	return NewDecimal(digits)
}

// splitDecimal validates plain decimal syntax and returns sign, integer and
// fractional digits without altering them.
func splitDecimal(s string) (negative bool, intPart, fracPart string, err error) {
	invalid := fmt.Errorf("%w: malformed decimal", ErrExactValueInvalid)
	body := s
	if strings.HasPrefix(body, "-") || strings.HasPrefix(body, "+") {
		negative = body[0] == '-'
		body = body[1:]
	}
	intPart, fracPart, _ = strings.Cut(body, ".")
	if intPart == "" && fracPart == "" {
		return false, "", "", invalid
	}
	if strings.Contains(fracPart, ".") {
		return false, "", "", invalid
	}
	for _, part := range []string{intPart, fracPart} {
		for _, r := range part {
			if r < '0' || r > '9' {
				return false, "", "", invalid
			}
		}
	}
	return negative, intPart, fracPart, nil
}

// canonicalDecimal strips redundant zeros so equal numbers compare equal.
func canonicalDecimal(s string) (string, bool) {
	negative, intPart, fracPart, err := splitDecimal(s)
	if err != nil {
		return "", false
	}
	intPart = strings.TrimLeft(intPart, "0")
	fracPart = strings.TrimRight(fracPart, "0")
	if intPart == "" && fracPart == "" {
		return "0", true
	}
	if intPart == "" {
		intPart = "0"
	}
	out := intPart
	if fracPart != "" {
		out += "." + fracPart
	}
	if negative {
		out = "-" + out
	}
	return out, true
}

// DecimalShape reports significant integer and fractional digit counts
// (leading integer zeros and trailing fractional zeros ignored).
func DecimalShape(digits string) (integerDigits, fractionDigits int, err error) {
	_, intPart, fracPart, err := splitDecimal(digits)
	if err != nil {
		return 0, 0, err
	}
	return len(strings.TrimLeft(intPart, "0")), len(strings.TrimRight(fracPart, "0")), nil
}
