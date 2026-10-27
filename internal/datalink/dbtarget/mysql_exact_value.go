package dbtarget

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"go-gateway/internal/datalink/measurement"
)

const (
	mysqlTypeBool            = "BOOL"
	mysqlTypeBigint          = "BIGINT"
	mysqlTypeVarchar         = "VARCHAR"
	mysqlTypeDoublePrecision = "DOUBLE PRECISION"
	mysqlTypeDatetime        = "DATETIME"
)

// ErrExactValueUnrepresentable marks an exact value that exceeds the declared
// SQL column's lossless range.
var ErrExactValueUnrepresentable = errors.New("exact_value_unrepresentable")

func exactRangeBlocked(strategy exactStrategy, reason string) error {
	if strategy.mysql || strategy.rangeErrorUnrepresentable {
		return fmt.Errorf("%w: %s", errors.Join(ErrExactSQLBlocked, ErrExactValueUnrepresentable), reason)
	}
	return blockedf("%s", reason)
}

func mysqlStrategy(d declaredType, typ measurement.ExactType) (exactStrategy, error) {
	switch typ {
	case measurement.ExactBool:
		if d.base == mysqlTypeBool || d.base == typeBoolean || d.base == "TINYINT" {
			return exactStrategy{kind: storageBoolean, mysql: true}, nil
		}
	case measurement.ExactText:
		switch d.base {
		case mysqlTypeVarchar:
			if d.hasP {
				return exactStrategy{kind: storageText, precision: d.p, bounded: true, mysql: true}, nil
			}
		case "CHAR":
			if d.hasP && d.p == 64 {
				return exactStrategy{kind: storageText, precision: 64, bounded: true, mysql: true}, nil
			}
		case "TINYTEXT", affinityText, "MEDIUMTEXT", "LONGTEXT":
			return exactStrategy{kind: storageText, mysql: true}, nil
		}
	case measurement.ExactInt64:
		if d.base == mysqlTypeBigint {
			return exactStrategy{kind: storageInteger, mysql: true}, nil
		}
	case measurement.ExactUint64, measurement.ExactDecimal:
		if d.base == "DECIMAL" || d.base == "NUMERIC" {
			return mysqlNumericStrategy(d)
		}
	case measurement.ExactFloat64:
		if d.base == "DOUBLE" || d.base == mysqlTypeDoublePrecision {
			return exactStrategy{kind: storageReal, mysql: true}, nil
		}
	default:
		return exactStrategy{}, blockedf("unknown exact type")
	}
	return exactStrategy{}, blockedf("MySQL %s column cannot store %s exactly", d.base, typ)
}

func mysqlNumericStrategy(d declaredType) (exactStrategy, error) {
	if !d.hasP || d.p < 1 || d.p > 65 || d.s > 30 || d.s > d.p {
		return exactStrategy{}, blockedf("MySQL DECIMAL precision/scale is not explicitly supported")
	}
	return exactStrategy{kind: storageNumeric, precision: d.p, scale: d.s, bounded: true, mysql: true}, nil
}

// mysqlUint64Readback accepts MySQL DECIMAL's fixed-scale integer rendering
// only when every fractional digit is zero. It never parses through float64.
func mysqlUint64Readback(value string) (uint64, bool) {
	integer, fraction, hasFraction := strings.Cut(value, ".")
	if integer == "" || (hasFraction && fraction == "") {
		return 0, false
	}
	if hasFraction {
		for _, digit := range fraction {
			if digit != '0' {
				return 0, false
			}
		}
	}
	parsed, err := strconv.ParseUint(integer, 10, 64)
	return parsed, err == nil
}
