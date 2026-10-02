package dbtarget

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"math"
	"strconv"
	"unicode/utf8"

	"go-gateway/internal/datalink/measurement"
)

// SQLDialect selects the type rules used by the exact value SQL codec.
type SQLDialect string

const (
	SQLDialectSQLite   SQLDialect = "sqlite"
	SQLDialectPostgres SQLDialect = "postgres"
)

// ErrExactSQLBlocked marks a value/type combination that cannot store the value exactly.
var ErrExactSQLBlocked = errors.New("exact sql type blocked")

// EncodeExactValue returns a parameterized driver value for a declared SQL
// column type, or ErrExactSQLBlocked when storing it could lose information.
// It is a pure type/range gate: the declaration is a caller claim, not
// verified column metadata.
func EncodeExactValue(dialect SQLDialect, declared string, value measurement.ExactValue) (driver.Value, error) {
	if value.Type() == "" {
		return nil, fmt.Errorf("%w: value is unset", measurement.ErrExactValueInvalid)
	}
	strategy, err := exactStrategyFor(dialect, declared, value.Type())
	if err != nil {
		return nil, err
	}
	switch v := value.Value().(type) {
	case bool:
		if dialect == SQLDialectPostgres {
			return v, nil
		}
		if v {
			return int64(1), nil
		}
		return int64(0), nil
	case string:
		if strategy.bounded && utf8.RuneCountInString(v) > strategy.precision {
			return nil, blockedf("text exceeds the declared column length")
		}
		return v, nil
	case int64:
		return v, nil
	case uint64:
		return encodeUint64(strategy, v)
	case measurement.Decimal:
		return encodeDecimal(strategy, string(v))
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return nil, fmt.Errorf("%w: float64 must be finite", measurement.ErrExactValueInvalid)
		}
		return v, nil
	}
	return nil, fmt.Errorf("%w: unsupported value", measurement.ErrExactValueInvalid)
}

func encodeUint64(strategy exactStrategy, v uint64) (driver.Value, error) {
	switch strategy.kind {
	case storageInteger:
		if v > math.MaxInt64 {
			return nil, blockedf("uint64 value exceeds the signed 64-bit column range")
		}
		return int64(v), nil
	case storageNumeric:
		digits := strconv.FormatUint(v, 10)
		return encodeDecimal(strategy, digits)
	}
	return strconv.FormatUint(v, 10), nil
}

func encodeDecimal(strategy exactStrategy, digits string) (driver.Value, error) {
	if strategy.kind == storageNumeric && strategy.bounded {
		integerDigits, fractionDigits, err := measurement.DecimalShape(digits)
		if err != nil {
			return nil, err
		}
		if fractionDigits > strategy.scale {
			return nil, blockedf("decimal would be rounded by the column scale")
		}
		if integerDigits > strategy.precision-strategy.scale {
			return nil, blockedf("decimal exceeds the column precision")
		}
	}
	return digits, nil
}

// DecodeExactValue converts a database/sql readback value to the expected
// exact type, rejecting shapes that could hide a lossy conversion.
func DecodeExactValue(dialect SQLDialect, declared string, expected measurement.ExactType, raw any) (measurement.ExactValue, error) {
	if _, err := exactStrategyFor(dialect, declared, expected); err != nil {
		return measurement.ExactValue{}, err
	}
	mismatch := func() (measurement.ExactValue, error) {
		return measurement.ExactValue{}, fmt.Errorf("%w: readback is not a %s", measurement.ErrExactValueInvalid, expected)
	}
	switch expected {
	case measurement.ExactBool:
		switch v := raw.(type) {
		case bool:
			return measurement.NewBool(v), nil
		case int64:
			if v == 0 || v == 1 {
				return measurement.NewBool(v == 1), nil
			}
		}
	case measurement.ExactText:
		if s, ok := readbackString(raw); ok {
			return measurement.NewText(s), nil
		}
	case measurement.ExactInt64:
		if v, ok := raw.(int64); ok {
			return measurement.NewInt64(v), nil
		}
	case measurement.ExactUint64:
		if v, ok := raw.(int64); ok && v >= 0 {
			return measurement.NewUint64(uint64(v)), nil
		}
		if s, ok := readbackString(raw); ok {
			if n, err := strconv.ParseUint(s, 10, 64); err == nil {
				return measurement.NewUint64(n), nil
			}
		}
	case measurement.ExactDecimal:
		if s, ok := readbackString(raw); ok {
			return measurement.NewDecimal(s)
		}
	case measurement.ExactFloat64:
		if v, ok := raw.(float64); ok {
			return measurement.NewFloat64(v)
		}
	}
	return mismatch()
}

func readbackString(raw any) (string, bool) {
	switch v := raw.(type) {
	case string:
		return v, true
	case []byte:
		return string(v), true
	}
	return "", false
}
