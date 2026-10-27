package dbtarget

import (
	"math"
	"strings"
	"testing"

	"go-gateway/internal/datalink/measurement"

	"github.com/stretchr/testify/require"
)

const sqlServerCodecTestDialect SQLDialect = SQLDialectSQLServer

func sqlServerCodecFloat(t *testing.T, value float64) measurement.ExactValue {
	t.Helper()
	result, err := measurement.NewFloat64(value)
	require.NoError(t, err)
	return result
}

func TestSQLServerExactCodecManagedLongValues(t *testing.T) {
	bit := measurement.NewBool(true)
	bound, err := EncodeExactValue(sqlServerCodecTestDialect, "BIT", bit)
	require.NoError(t, err)
	require.Equal(t, true, bound, "BIT values bind as bool rather than integer surrogates")
	gotBool, err := DecodeExactValue(sqlServerCodecTestDialect, "BIT", measurement.ExactBool, bound)
	require.NoError(t, err)
	require.True(t, bit.Equal(gotBool))

	largeFrom := measurement.NewInt64(9007199254740992)
	largeTo := measurement.NewInt64(9007199254740993)
	readInt := func(value measurement.ExactValue) measurement.ExactValue {
		bound, err := EncodeExactValue(sqlServerCodecTestDialect, "BIGINT", value)
		require.NoError(t, err)
		require.IsType(t, int64(0), bound)
		got, err := DecodeExactValue(sqlServerCodecTestDialect, "BIGINT", measurement.ExactInt64, bound)
		require.NoError(t, err)
		return got
	}
	gotFrom, gotTo := readInt(largeFrom), readInt(largeTo)
	delta, err := measurement.Subtract(gotTo, gotFrom)
	require.NoError(t, err)
	require.True(t, decimalValue(t, "1").Equal(delta), "adjacent values above 2^53 must remain one apart")

	maxUint := measurement.NewUint64(math.MaxUint64)
	bound, err = EncodeExactValue(sqlServerCodecTestDialect, "DECIMAL(38,18)", maxUint)
	require.NoError(t, err)
	require.Equal(t, "18446744073709551615", bound, "DECIMAL must be bound as exact decimal text")
	gotUint, err := DecodeExactValue(
		sqlServerCodecTestDialect,
		"DECIMAL(38,18)",
		measurement.ExactUint64,
		"18446744073709551615.000000000000000000",
	)
	require.NoError(t, err)
	require.True(t, maxUint.Equal(gotUint))

	from := decimalValue(t, "12000.0")
	to := decimalValue(t, "12000.4")
	readDecimal := func(raw string) measurement.ExactValue {
		value := decimalValue(t, raw)
		bound, err := EncodeExactValue(sqlServerCodecTestDialect, "DECIMAL(38,18)", value)
		require.NoError(t, err)
		require.Equal(t, string(value.Value().(measurement.Decimal)), bound)
		got, err := DecodeExactValue(
			sqlServerCodecTestDialect,
			"DECIMAL(38,18)",
			measurement.ExactDecimal,
			[]byte(raw),
		)
		require.NoError(t, err)
		return got
	}
	gotFromDecimal := readDecimal("12000.000000000000000000")
	gotToDecimal := readDecimal("12000.400000000000000000")
	require.True(t, from.Equal(gotFromDecimal))
	require.True(t, to.Equal(gotToDecimal))
	delta, err = measurement.Subtract(gotToDecimal, gotFromDecimal)
	require.NoError(t, err)
	require.True(t, decimalValue(t, "0.4").Equal(delta))

	float, err := measurement.NewFloat64(1.25)
	require.NoError(t, err)
	bound, err = EncodeExactValue(sqlServerCodecTestDialect, "FLOAT(53)", float)
	require.NoError(t, err)
	require.IsType(t, float64(0), bound)
	gotFloat, err := DecodeExactValue(sqlServerCodecTestDialect, "FLOAT(53)", measurement.ExactFloat64, bound)
	require.NoError(t, err)
	require.True(t, float.Equal(gotFloat))
}

func TestSQLServerExactCodecUnicodeAndCaseExactText(t *testing.T) {
	for _, value := range []string{"Device-Alpha-😀", "device-alpha-😀"} {
		text := measurement.NewText(value)
		bound, err := EncodeExactValue(sqlServerCodecTestDialect, "NVARCHAR(64)", text)
		require.NoError(t, err)
		require.Equal(t, value, bound)
		got, err := DecodeExactValue(sqlServerCodecTestDialect, "NVARCHAR(64)", measurement.ExactText, []byte(value))
		require.NoError(t, err)
		require.True(t, text.Equal(got))
	}
	upper, err := EncodeExactValue(sqlServerCodecTestDialect, "NVARCHAR(64)", measurement.NewText("Device-Alpha"))
	require.NoError(t, err)
	lower, err := EncodeExactValue(sqlServerCodecTestDialect, "NVARCHAR(64)", measurement.NewText("device-alpha"))
	require.NoError(t, err)
	require.NotEqual(t, upper, lower, "codec must not fold case-sensitive source identities")

	// SQL Server NVARCHAR(n) lengths count UTF-16 code units: the astral rune
	// occupies two units, so this value exactly fills NVARCHAR(4).
	const exactlyFourUnits = "A😀B"
	bound, err := EncodeExactValue(sqlServerCodecTestDialect, "NVARCHAR(4)", measurement.NewText(exactlyFourUnits))
	require.NoError(t, err)
	require.Equal(t, exactlyFourUnits, bound)
	got, err := DecodeExactValue(sqlServerCodecTestDialect, "NVARCHAR(4)", measurement.ExactText, []byte(exactlyFourUnits))
	require.NoError(t, err)
	require.Equal(t, exactlyFourUnits, got.Value())

	longUnicode := strings.Repeat("設備😀", 1500)
	bound, err = EncodeExactValue(sqlServerCodecTestDialect, "NVARCHAR(MAX)", measurement.NewText(longUnicode))
	require.NoError(t, err)
	require.Equal(t, longUnicode, bound)
	got, err = DecodeExactValue(sqlServerCodecTestDialect, "NVARCHAR(MAX)", measurement.ExactText, []byte(longUnicode))
	require.NoError(t, err)
	require.Equal(t, longUnicode, got.Value())

	invalidUTF8 := measurement.NewText(string([]byte{0xff}))
	_, err = EncodeExactValue(SQLDialectSQLServer, "NVARCHAR(MAX)", invalidUTF8)
	require.ErrorIs(t, err, ErrExactSQLBlocked)
	require.ErrorIs(t, err, ErrExactValueUnrepresentable)
	mysqlBound, err := EncodeExactValue(SQLDialectMySQL, "LONGTEXT", invalidUTF8)
	require.NoError(t, err, "UTF-8 rejection is SQL Server-specific and must not change MySQL semantics")
	require.Equal(t, invalidUTF8.Value(), mysqlBound)
}

func TestSQLServerExactCodecRejectsLossyDeclarations(t *testing.T) {
	for _, declared := range []string{
		"VARCHAR(64)", "VARCHAR(MAX)", "CHAR(64)", "TEXT", "NTEXT", "NVARCHAR",
		"NVARCHAR(0)", "NVARCHAR(4001)",
	} {
		_, err := EncodeExactValue(sqlServerCodecTestDialect, declared, measurement.NewText("設備😀"))
		require.ErrorIs(t, err, ErrExactSQLBlocked, declared)
	}
	for _, test := range []struct {
		declared string
		value    measurement.ExactValue
	}{
		{declared: "BOOLEAN", value: measurement.NewBool(true)},
		{declared: "INTEGER", value: measurement.NewInt64(1)},
		{declared: "FLOAT(24)", value: sqlServerCodecFloat(t, 1.25)},
		{declared: "DECIMAL(39,18)", value: decimalValue(t, "1.25")},
		{declared: "DECIMAL(38,19)", value: decimalValue(t, "1.25")},
	} {
		_, err := EncodeExactValue(sqlServerCodecTestDialect, test.declared, test.value)
		require.ErrorIs(t, err, ErrExactSQLBlocked, test.declared)
	}
	_, err := EncodeExactValue(SQLDialectPostgres, "NVARCHAR(MAX)", measurement.NewText("keep dialect rules"))
	require.ErrorIs(t, err, ErrExactSQLBlocked, "MAX support must not broaden PostgreSQL")
}

func TestSQLServerExactCodecStableHashRecordIDs(t *testing.T) {
	const recordID = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	text := measurement.NewText(recordID)
	bound, err := EncodeExactValue(SQLDialectSQLServer, "CHAR(64)", text)
	require.NoError(t, err)
	require.Equal(t, recordID, bound)
	got, err := DecodeExactValue(SQLDialectSQLServer, "CHAR(64)", measurement.ExactText, []byte(recordID))
	require.NoError(t, err)
	require.True(t, text.Equal(got))

	for _, value := range []string{
		strings.ToUpper(recordID),
		recordID[:63],
		recordID + "0",
		"g" + recordID[1:],
		"設備" + recordID[2:],
	} {
		_, err := EncodeExactValue(SQLDialectSQLServer, "CHAR(64)", measurement.NewText(value))
		require.ErrorIs(t, err, ErrExactSQLBlocked)
		require.ErrorIs(t, err, ErrExactValueUnrepresentable)
		require.NotContains(t, err.Error(), value, "rejected IDs must not leak through codec errors")
	}
	_, err = DecodeExactValue(SQLDialectSQLServer, "CHAR(64)", measurement.ExactText, []byte("not-a-record-id"))
	require.ErrorIs(t, err, measurement.ErrExactValueInvalid)
}

func TestSQLServerExactCodecRejectsUnrepresentableValues(t *testing.T) {
	maximum := strings.Repeat("9", 20) + "." + strings.Repeat("9", 18)
	bound, err := EncodeExactValue(sqlServerCodecTestDialect, "DECIMAL(38,18)", decimalValue(t, maximum))
	require.NoError(t, err)
	require.Equal(t, maximum, bound)

	for _, digits := range []string{
		strings.Repeat("9", 21),
		"0." + strings.Repeat("1", 19),
	} {
		_, err := EncodeExactValue(sqlServerCodecTestDialect, "DECIMAL(38,18)", decimalValue(t, digits))
		require.ErrorIs(t, err, ErrExactSQLBlocked)
		require.ErrorIs(t, err, ErrExactValueUnrepresentable)
		require.NotContains(t, err.Error(), digits, "range errors must not disclose rejected values")
	}

	tooLong := "A😀BC"
	_, err = EncodeExactValue(sqlServerCodecTestDialect, "NVARCHAR(4)", measurement.NewText(tooLong))
	require.ErrorIs(t, err, ErrExactSQLBlocked)
	require.ErrorIs(t, err, ErrExactValueUnrepresentable)
	require.NotContains(t, err.Error(), tooLong)

	for _, raw := range []any{
		"18446744073709551615.000000000000000001",
		"18446744073709551616.000000000000000000",
		float64(1),
	} {
		_, err := DecodeExactValue(sqlServerCodecTestDialect, "DECIMAL(38,18)", measurement.ExactUint64, raw)
		require.Error(t, err, "uint64 fixed-scale readback must reject fractions, overflow, and float64")
	}
}
