package dbtarget

import (
	"regexp"
	"strconv"
	"strings"
	"unicode/utf16"

	"go-gateway/internal/datalink/measurement"
)

const (
	sqlServerMaxNumericPrecision = 38
	sqlServerNumericScale        = 18
	sqlServerMaxNvarcharLength   = 4000
	sqlServerBigintSQLType       = "BIGINT"
	sqlServerDecimalSQLType      = "DECIMAL"
	sqlServerHashSQLType         = "CHAR(64)"
	sqlServerDateTime2SQLType    = "DATETIME2(6)"
	sqlServerNvarcharMaxSQLType  = "NVARCHAR(MAX)"
	sqlServerBitSQLType          = "BIT"
)

var sqlServerNvarcharMaxPattern = regexp.MustCompile(`(?i)^NVARCHAR\s*\(\s*MAX\s*\)$`)

func sqlServerStrategy(declared declaredType, typ measurement.ExactType) (exactStrategy, error) {
	rangeErrorUnrepresentable := true
	switch typ {
	case measurement.ExactBool:
		if declared.base == sqlServerBitSQLType && !declared.hasP {
			return exactStrategy{kind: storageBoolean}, nil
		}
	case measurement.ExactText:
		if declared.base == "CHAR" && declared.hasP && declared.p == 64 && !declared.hasS {
			return exactStrategy{
				kind: storageText, precision: 64, bounded: true,
				rangeErrorUnrepresentable: rangeErrorUnrepresentable, sqlServerHashID: true,
			}, nil
		}
		if declared.base != "NVARCHAR" || declared.hasS {
			break
		}
		if declared.hasP && declared.p >= 1 && declared.p <= sqlServerMaxNvarcharLength {
			return exactStrategy{
				kind: storageText, precision: declared.p, bounded: true,
				rangeErrorUnrepresentable: rangeErrorUnrepresentable,
			}, nil
		}
	case measurement.ExactInt64:
		if declared.base == sqlServerBigintSQLType && !declared.hasP {
			return exactStrategy{kind: storageInteger}, nil
		}
	case measurement.ExactUint64, measurement.ExactDecimal:
		if (declared.base == sqlServerDecimalSQLType || declared.base == affinityNumeric) && declared.hasP &&
			declared.p == sqlServerMaxNumericPrecision && declared.hasS && declared.s == sqlServerNumericScale {
			return exactStrategy{
				kind: storageNumeric, precision: sqlServerMaxNumericPrecision, scale: sqlServerNumericScale,
				bounded: true, rangeErrorUnrepresentable: rangeErrorUnrepresentable,
			}, nil
		}
	case measurement.ExactFloat64:
		if declared.base == "FLOAT" && declared.hasP && declared.p == 53 && !declared.hasS {
			return exactStrategy{kind: storageReal}, nil
		}
	default:
		return exactStrategy{}, blockedf("unknown exact type")
	}
	return exactStrategy{}, blockedf("SQL Server %s column cannot store %s exactly", declared.base, typ)
}

func sqlServerTextLength(value string) int {
	return len(utf16.Encode([]rune(value)))
}

func sqlServerStableHashID(value string) bool {
	if len(value) != 64 {
		return false
	}
	for i := range len(value) {
		digit := value[i]
		if (digit < '0' || digit > '9') && (digit < 'a' || digit > 'f') {
			return false
		}
	}
	return true
}

// sqlServerUint64Readback accepts fixed-scale decimal text only when every
// fractional digit is zero. It never parses through float64.
func sqlServerUint64Readback(value string) (uint64, bool) {
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
