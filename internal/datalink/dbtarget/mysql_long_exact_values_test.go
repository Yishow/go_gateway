package dbtarget

import (
	"encoding/json"
	"math"
	"strings"
	"testing"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/measurement"

	"github.com/stretchr/testify/require"
)

const mysqlDialectForTest SQLDialect = "mysql"

func TestMySQLLongExactValues(t *testing.T) {
	t.Run("uint64 maximum survives fixed-scale DECIMAL readback", func(t *testing.T) {
		value := measurement.NewUint64(math.MaxUint64)
		wire, err := json.Marshal(value)
		require.NoError(t, err)
		var fromJSON measurement.ExactValue
		require.NoError(t, json.Unmarshal(wire, &fromJSON))

		encoded, err := EncodeExactValue(mysqlDialectForTest, "DECIMAL(65,30)", fromJSON)
		require.NoError(t, err)
		require.Equal(t, "18446744073709551615", encoded)

		got, err := DecodeExactValue(
			mysqlDialectForTest,
			"DECIMAL(65,30)",
			measurement.ExactUint64,
			[]byte("18446744073709551615.000000000000000000000000000000"),
		)
		require.NoError(t, err)
		require.True(t, value.Equal(got))
	})

	t.Run("large integers retain a unit difference", func(t *testing.T) {
		from := measurement.NewInt64(9007199254740992)
		to := measurement.NewInt64(9007199254740993)
		read := func(value measurement.ExactValue) measurement.ExactValue {
			bound, err := EncodeExactValue(mysqlDialectForTest, "BIGINT", value)
			require.NoError(t, err)
			got, err := DecodeExactValue(mysqlDialectForTest, "BIGINT", measurement.ExactInt64, bound)
			require.NoError(t, err)
			return got
		}
		gotFrom, gotTo := read(from), read(to)
		delta, err := measurement.Subtract(gotTo, gotFrom)
		require.NoError(t, err)
		require.True(t, decimalValue(t, "1").Equal(delta))
	})

	t.Run("decimal values retain a four-tenths difference", func(t *testing.T) {
		from := decimalValue(t, "12000.0")
		to := decimalValue(t, "12000.4")
		read := func(value measurement.ExactValue, raw string) measurement.ExactValue {
			bound, err := EncodeExactValue(mysqlDialectForTest, "DECIMAL(65,30)", value)
			require.NoError(t, err)
			require.Equal(t, string(value.Value().(measurement.Decimal)), bound)
			got, err := DecodeExactValue(mysqlDialectForTest, "DECIMAL(65,30)", measurement.ExactDecimal, []byte(raw))
			require.NoError(t, err)
			return got
		}
		gotFrom := read(from, "12000.000000000000000000000000000000")
		gotTo := read(to, "12000.400000000000000000000000000000")
		delta, err := measurement.Subtract(gotTo, gotFrom)
		require.NoError(t, err)
		require.True(t, decimalValue(t, "0.4").Equal(delta))
	})

	t.Run("DECIMAL bounds reject lossy values without rounding", func(t *testing.T) {
		maximum := strings.Repeat("9", 35) + "." + strings.Repeat("9", 30)
		bound, err := EncodeExactValue(mysqlDialectForTest, "DECIMAL(65,30)", decimalValue(t, maximum))
		require.NoError(t, err)
		require.Equal(t, maximum, bound)

		trailingZero := maximum + "0"
		bound, err = EncodeExactValue(mysqlDialectForTest, "DECIMAL(65,30)", decimalValue(t, trailingZero))
		require.NoError(t, err, "discarding a trailing fractional zero is exact")
		require.Equal(t, trailingZero, bound)

		for _, digits := range []string{
			strings.Repeat("9", 36),
			"0." + strings.Repeat("1", 31),
		} {
			_, err := EncodeExactValue(mysqlDialectForTest, "DECIMAL(65,30)", decimalValue(t, digits))
			require.ErrorIs(t, err, ErrExactSQLBlocked)
			require.ErrorIs(t, err, ErrExactValueUnrepresentable)
			require.Contains(t, err.Error(), "exact_value_unrepresentable")
			require.NotContains(t, err.Error(), digits, "safe errors must not echo the rejected value")
		}
	})

	t.Run("fixed hash IDs use bounded CHAR(64)", func(t *testing.T) {
		id := strings.Repeat("a", 64)
		bound, err := EncodeExactValue(mysqlDialectForTest, "CHAR(64)", measurement.NewText(id))
		require.NoError(t, err)
		require.Equal(t, id, bound)
		got, err := DecodeExactValue(mysqlDialectForTest, "CHAR(64)", measurement.ExactText, []byte(id))
		require.NoError(t, err)
		require.Equal(t, id, got.Value())

		for _, declared := range []string{"CHAR", "CHAR(63)", "CHAR(65)"} {
			_, err := EncodeExactValue(mysqlDialectForTest, declared, measurement.NewText(id))
			require.ErrorIs(t, err, ErrExactSQLBlocked, declared)
			require.NotErrorIs(t, err, ErrExactValueUnrepresentable, declared)
		}

		_, err = EncodeExactValue(mysqlDialectForTest, "CHAR(64)", measurement.NewText(id+"a"))
		require.ErrorIs(t, err, ErrExactSQLBlocked)
		require.ErrorIs(t, err, ErrExactValueUnrepresentable)
	})

	t.Run("Unicode text and case remain unchanged", func(t *testing.T) {
		for _, declared := range []string{"VARCHAR(32)", "LONGTEXT"} {
			for _, value := range []string{"設備-Alpha-😀", "設備-alpha-😀"} {
				bound, err := EncodeExactValue(mysqlDialectForTest, declared, measurement.NewText(value))
				require.NoError(t, err)
				got, err := DecodeExactValue(mysqlDialectForTest, declared, measurement.ExactText, []byte(bound.(string)))
				require.NoError(t, err)
				require.Equal(t, value, got.Value())
			}
		}
	})

	t.Run("fractional and out-of-range integer readbacks are refused", func(t *testing.T) {
		for _, raw := range []string{
			"18446744073709551615.000000000000000000000000000001",
			"18446744073709551616.000000000000000000000000000000",
		} {
			_, err := DecodeExactValue(mysqlDialectForTest, "DECIMAL(65,30)", measurement.ExactUint64, []byte(raw))
			require.Error(t, err)
		}
	})

	t.Run("unbounded or lossy declarations fail closed", func(t *testing.T) {
		for _, declared := range []string{"DECIMAL", "DECIMAL(66,30)", "DECIMAL(65,31)", "FLOAT"} {
			_, err := EncodeExactValue(mysqlDialectForTest, declared, decimalValue(t, "1.25"))
			require.ErrorIs(t, err, ErrExactSQLBlocked, declared)
			require.NotErrorIs(t, err, ErrExactValueUnrepresentable, declared)
		}
		_, err := EncodeExactValue(mysqlDialectForTest, "VARCHAR(1)", measurement.NewText("測試"))
		require.ErrorIs(t, err, ErrExactSQLBlocked)
		require.ErrorIs(t, err, ErrExactValueUnrepresentable)
	})
}

func TestMySQLLongColumnDeclarations(t *testing.T) {
	datetime := func(declared string) bool {
		return longColumnUsable(SQLDialectMySQL, ColumnInfo{DataType: declared}, longTime)
	}
	for _, declared := range []string{"DATETIME", "DATETIME(3)", "TIMESTAMP(6)"} {
		require.False(t, datetime(declared), declared)
	}
	require.True(t, datetime("DATETIME(6)"))

	provenance := func(declared string) bool {
		return longColumnUsable(SQLDialectMySQL, ColumnInfo{DataType: declared}, longProvenanceRole)
	}
	require.True(t, provenance("JSON"))
	require.True(t, provenance("LONGTEXT"), "MariaDB's JSON alias is LONGTEXT")
	require.False(t, provenance("VARCHAR(191)"))
}

func TestMySQLLongScope(t *testing.T) {
	spec := mysqlLongIdentitySpec()
	spec.WorkspaceID = strings.Repeat("測", 191)
	spec.GroupID = strings.Repeat("g", 191)
	spec.GroupRevision = strings.Repeat("r", 191)
	spec.ProfileKind = strings.Repeat("p", 191)
	spec.Members[0].DeviceID = strings.Repeat("d", 191)
	spec.Members[0].PointID = strings.Repeat("p", 191)
	spec.Members[0].TagID = strings.Repeat("t", 191)
	spec.Members[0].CircuitID = strings.Repeat("c", 191)
	spec.Members[0].MeasurementID = strings.Repeat("m", 191)
	spec.Members[0].Counter.SeriesEpoch = strings.Repeat("e", 191)
	spec.Members[0].Counter.Unit = strings.Repeat("u", 191)
	require.Empty(t, validateMySQLLongIdentity(spec), "VARCHAR length counts Unicode code points")

	spec.WorkspaceID += "w"
	spec.GroupID += "g"
	spec.GroupRevision += "r"
	spec.ProfileKind += "p"
	spec.Members[0].DeviceID += "d"
	spec.Members[0].PointID += "p"
	spec.Members[0].TagID += "t"
	spec.Members[0].CircuitID += "c"
	spec.Members[0].MeasurementID += "m"
	spec.Members[0].Counter.SeriesEpoch += "e"
	spec.Members[0].Counter.Unit += "u"
	issues := validateMySQLLongIdentity(spec)
	for _, column := range []string{
		longColWorkspaceID, longColGroupID, longColGroupRev, longColProfileKind,
		longColDeviceID, longColPointID, longColTagID, longColCircuitID, longColMeasurement,
		longColSeriesEpoch, longColUnit,
	} {
		memberKey := ""
		if column != longColWorkspaceID && column != longColGroupID && column != longColGroupRev && column != longColProfileKind {
			memberKey = "member-1"
		}
		require.Contains(t, issues, LayoutIssue{Code: issueIdentityUnsupported, Column: column, MemberKey: memberKey}, column)
	}

	unbounded := mysqlLongIdentitySpec()
	unbounded.Columns[0].DataType = "VARCHAR"
	require.Contains(t, validateMySQLLongIdentity(unbounded), LayoutIssue{Code: issueIdentityUnsupported, Column: longColWorkspaceID})

	spec.Dialect = SQLDialectPostgres
	require.Empty(t, validateMySQLLongIdentity(spec), "the MySQL-only gate does not alter PostgreSQL")
}

func mysqlLongIdentitySpec() LongRowSpec {
	columns := make([]ColumnInfo, 0, 12)
	for _, name := range []string{
		longColWorkspaceID, longColGroupID, longColGroupRev, longColProfileKind,
		longColDeviceID, longColPointID, longColTagID, longColCircuitID, longColMeasurement,
		longColSeriesEpoch, longColUnit,
	} {
		columns = append(columns, ColumnInfo{Name: name, DataType: "VARCHAR(191)"})
	}
	return LongRowSpec{
		Dialect: SQLDialectMySQL, Columns: columns,
		WorkspaceID: "workspace", GroupID: "group", GroupRevision: "revision",
		ProfileKind: "instantaneous", ProfileUnit: "kWh",
		Members: []LongRowMember{{
			MemberKey: "member-1", DeviceID: "device", PointID: "point", TagID: "tag",
			CircuitID: "circuit", MeasurementID: "measurement", Type: measurement.ExactDecimal,
			Counter: &aggregation.CounterBinding{SeriesEpoch: "epoch", Unit: "kWh"},
		}},
	}
}

func TestExactValueNonMySQLErrorContracts(t *testing.T) {
	_, err := EncodeExactValue(SQLDialectSQLite, "INTEGER", measurement.NewUint64(math.MaxUint64))
	require.EqualError(t, err, "exact sql type blocked: uint64 value exceeds the signed 64-bit column range")
	require.NotErrorIs(t, err, ErrExactValueUnrepresentable)

	_, err = EncodeExactValue(SQLDialectPostgres, "NUMERIC(18,0)", measurement.NewUint64(math.MaxUint64))
	require.EqualError(t, err, "exact sql type blocked: decimal exceeds the column precision")
	require.NotErrorIs(t, err, ErrExactValueUnrepresentable)

	_, err = EncodeExactValue(SQLDialectPostgres, "VARCHAR(3)", measurement.NewText("four"))
	require.EqualError(t, err, "exact sql type blocked: text exceeds the declared column length")
	require.NotErrorIs(t, err, ErrExactValueUnrepresentable)
}

type longSampleValueValidator interface {
	ValidateSampleValue(measurement.ExactValue) error
}

func TestMySQLLongExactValueIntakeGate(t *testing.T) {
	layout := &LongRowLayout{
		spec: LongRowSpec{Dialect: mysqlDialectForTest},
		columns: map[string]ColumnInfo{
			longColValueDecimal: {Name: longColValueDecimal, DataType: "DECIMAL(65,30)"},
		},
	}
	validator, ok := any(layout).(longSampleValueValidator)
	require.True(t, ok, "the intake hook must be available to the shared boundary")

	require.NoError(t, validator.ValidateSampleValue(measurement.NewUint64(math.MaxUint64)))
	value := "0." + strings.Repeat("1", 31)
	err := validator.ValidateSampleValue(decimalValue(t, value))
	require.ErrorIs(t, err, ErrExactSQLBlocked)
	require.ErrorIs(t, err, ErrExactValueUnrepresentable)
	require.NotContains(t, err.Error(), value)
	require.ErrorIs(t, validator.ValidateSampleValue(measurement.ExactValue{}), ErrExactSQLBlocked)

	withoutSlot := &LongRowLayout{spec: LongRowSpec{Dialect: mysqlDialectForTest}}
	withoutSlotValidator, ok := any(withoutSlot).(longSampleValueValidator)
	require.True(t, ok)
	require.ErrorIs(t, withoutSlotValidator.ValidateSampleValue(measurement.NewUint64(1)), ErrExactSQLBlocked)

	unboundedSlot := &LongRowLayout{
		spec: LongRowSpec{Dialect: mysqlDialectForTest},
		columns: map[string]ColumnInfo{
			longColValueDecimal: {Name: longColValueDecimal, DataType: "DECIMAL"},
		},
	}
	unboundedValidator, ok := any(unboundedSlot).(longSampleValueValidator)
	require.True(t, ok)
	require.ErrorIs(t, unboundedValidator.ValidateSampleValue(measurement.NewUint64(1)), ErrExactSQLBlocked)

	var unavailable *LongRowLayout
	unavailableValidator, ok := any(unavailable).(longSampleValueValidator)
	require.True(t, ok)
	require.ErrorIs(t, unavailableValidator.ValidateSampleValue(measurement.NewUint64(1)), ErrExactSQLBlocked)

	nonMySQL := &LongRowLayout{spec: LongRowSpec{Dialect: SQLDialectSQLite}}
	validator, ok = any(nonMySQL).(longSampleValueValidator)
	require.True(t, ok)
	require.NoError(t, validator.ValidateSampleValue(measurement.ExactValue{}))
}
