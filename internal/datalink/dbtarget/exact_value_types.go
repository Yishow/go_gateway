package dbtarget

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"go-gateway/internal/datalink/measurement"
)

// exactStrategy is how a declared SQL column stores one exact value type.
type exactStrategy struct {
	kind      exactStorage
	precision int  // numeric precision; 0 with !bounded means unconstrained
	scale     int  // numeric scale
	bounded   bool // numeric(p,s) or varchar(n) declared with a limit
}

type exactStorage int

const (
	storageBoolean exactStorage = iota + 1
	storageInteger
	storageText
	storageNumeric
	storageReal
)

const (
	affinityInteger = "INTEGER"
	affinityText    = "TEXT"
	affinityReal    = "REAL"
	affinityNumeric = "NUMERIC"
	affinityBlob    = "BLOB"
	typeBoolean     = "BOOLEAN"
)

var declaredTypePattern = regexp.MustCompile(`^([A-Z][A-Z0-9 ]*?)\s*(?:\(\s*(\d+)\s*(?:,\s*(\d+)\s*)?\))?$`)

func blockedf(format string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrExactSQLBlocked, fmt.Sprintf(format, args...))
}

type declaredType struct {
	base       string
	p, s       int
	hasP, hasS bool
}

func parseDeclaredType(declared string) (declaredType, error) {
	match := declaredTypePattern.FindStringSubmatch(strings.ToUpper(strings.TrimSpace(declared)))
	if match == nil {
		return declaredType{}, blockedf("unrecognized column type declaration")
	}
	out := declaredType{base: strings.Join(strings.Fields(match[1]), " ")}
	if match[2] != "" {
		p, err := strconv.Atoi(match[2])
		if err != nil {
			return declaredType{}, blockedf("column type precision out of range")
		}
		out.p, out.hasP = p, true
	}
	if match[3] != "" {
		s, err := strconv.Atoi(match[3])
		if err != nil {
			return declaredType{}, blockedf("column type scale out of range")
		}
		out.s, out.hasS = s, true
	}
	return out, nil
}

func exactStrategyFor(dialect SQLDialect, declared string, typ measurement.ExactType) (exactStrategy, error) {
	parsed, err := parseDeclaredType(declared)
	if err != nil {
		return exactStrategy{}, err
	}
	switch dialect {
	case SQLDialectSQLite:
		return sqliteStrategy(parsed, typ)
	case SQLDialectPostgres:
		return postgresStrategy(parsed, typ)
	}
	return exactStrategy{}, blockedf("unsupported SQL dialect")
}

// hasDoubleMarker matches the double-precision marker of the SQLite affinity rules.
func hasDoubleMarker(base string) bool {
	return strings.Contains(base, "DOUB") //nolint:misspell // substring named by the SQLite affinity rules
}

// sqliteAffinity follows the SQLite type affinity rules for declared types.
func sqliteAffinity(base string) string {
	switch {
	case strings.Contains(base, "INT"):
		return affinityInteger
	case strings.Contains(base, "CHAR"), strings.Contains(base, "CLOB"), strings.Contains(base, affinityText):
		return affinityText
	case strings.Contains(base, affinityBlob):
		return affinityBlob
	case strings.Contains(base, affinityReal), strings.Contains(base, "FLOA"), hasDoubleMarker(base):
		return affinityReal
	}
	return affinityNumeric
}

func sqliteStrategy(d declaredType, typ measurement.ExactType) (exactStrategy, error) {
	affinity := sqliteAffinity(d.base)
	switch typ {
	case measurement.ExactBool:
		if affinity == affinityInteger || d.base == typeBoolean || d.base == "BOOL" {
			return exactStrategy{kind: storageBoolean}, nil
		}
	case measurement.ExactText:
		if affinity == affinityText {
			return exactStrategy{kind: storageText}, nil
		}
	case measurement.ExactInt64:
		if affinity == affinityInteger {
			return exactStrategy{kind: storageInteger}, nil
		}
	case measurement.ExactUint64:
		switch affinity {
		case affinityInteger:
			return exactStrategy{kind: storageInteger}, nil
		case affinityText:
			return exactStrategy{kind: storageText}, nil
		}
	case measurement.ExactDecimal:
		if affinity == affinityText {
			return exactStrategy{kind: storageText}, nil
		}
	case measurement.ExactFloat64:
		if affinity == affinityReal {
			return exactStrategy{kind: storageReal}, nil
		}
	default:
		return exactStrategy{}, blockedf("unknown exact type")
	}
	return exactStrategy{}, blockedf("SQLite %s column cannot store %s exactly", d.base, typ)
}

func postgresStrategy(d declaredType, typ measurement.ExactType) (exactStrategy, error) {
	numeric := d.base == "NUMERIC" || d.base == "DECIMAL"
	bigint := d.base == "BIGINT" || d.base == "INT8"
	text := d.base == affinityText || ((d.base == "VARCHAR" || d.base == "CHARACTER VARYING") && !d.hasS)
	switch typ {
	case measurement.ExactBool:
		if d.base == typeBoolean || d.base == "BOOL" {
			return exactStrategy{kind: storageBoolean}, nil
		}
	case measurement.ExactText:
		if text {
			return exactStrategy{kind: storageText, precision: d.p, bounded: d.hasP}, nil
		}
	case measurement.ExactInt64:
		if bigint {
			return exactStrategy{kind: storageInteger}, nil
		}
	case measurement.ExactUint64:
		switch {
		case bigint:
			return exactStrategy{kind: storageInteger}, nil
		case numeric:
			return numericStrategy(d)
		case d.base == affinityText:
			return exactStrategy{kind: storageText}, nil
		}
	case measurement.ExactDecimal:
		switch {
		case numeric:
			return numericStrategy(d)
		case d.base == affinityText:
			return exactStrategy{kind: storageText}, nil
		}
	case measurement.ExactFloat64:
		if d.base == "DOUBLE PRECISION" || d.base == "FLOAT8" {
			return exactStrategy{kind: storageReal}, nil
		}
	default:
		return exactStrategy{}, blockedf("unknown exact type")
	}
	return exactStrategy{}, blockedf("PostgreSQL %s column cannot store %s exactly", d.base, typ)
}

func numericStrategy(d declaredType) (exactStrategy, error) {
	if !d.hasP {
		return exactStrategy{kind: storageNumeric}, nil
	}
	if d.p < 1 || d.s > d.p {
		return exactStrategy{}, blockedf("invalid numeric precision/scale")
	}
	return exactStrategy{kind: storageNumeric, precision: d.p, scale: d.s, bounded: true}, nil
}
