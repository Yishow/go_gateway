package dbtarget

import (
	"database/sql"
	"encoding/json"
	"math"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink/measurement"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func decimalValue(t *testing.T, digits string) measurement.ExactValue {
	t.Helper()
	v, err := measurement.NewDecimal(digits)
	require.NoError(t, err)
	return v
}

func TestExactMixedValueRoundTripSQLite(t *testing.T) {
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "exact.db"))
	require.NoError(t, err)
	defer db.Close()
	_, err = db.ExecContext(t.Context(), `CREATE TABLE exact_row (
		running INTEGER, batch TEXT, minimum INTEGER, counter INTEGER, big TEXT, reading TEXT)`)
	require.NoError(t, err)

	type column struct {
		name, declared string
		value          measurement.ExactValue
	}
	columns := []column{
		{"running", "INTEGER", measurement.NewBool(true)},
		{"batch", "TEXT", measurement.NewText("batch-001")},
		{"minimum", "INTEGER", measurement.NewInt64(math.MinInt64)},
		{"counter", "INTEGER", measurement.NewUint64(9007199254740993)},
		{"big", "TEXT", measurement.NewUint64(math.MaxUint64)},
		{"reading", "TEXT", decimalValue(t, "1234567890.123456789012345678")},
	}
	// Values cross JSON before reaching the codec, as they do between runtime and sender.
	args := make([]any, 0, len(columns))
	for _, c := range columns {
		raw, err := json.Marshal(c.value)
		require.NoError(t, err)
		var decoded measurement.ExactValue
		require.NoError(t, json.Unmarshal(raw, &decoded))
		bound, err := EncodeExactValue(SQLDialectSQLite, c.declared, decoded)
		require.NoError(t, err, c.name)
		args = append(args, bound)
	}
	_, err = db.ExecContext(t.Context(), `INSERT INTO exact_row VALUES (?,?,?,?,?,?)`, args...)
	require.NoError(t, err)

	readback := make([]any, len(columns))
	dest := make([]any, len(columns))
	for i := range dest {
		dest[i] = &readback[i]
	}
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT running,batch,minimum,counter,big,reading FROM exact_row`).Scan(dest...))
	for i, c := range columns {
		got, err := DecodeExactValue(SQLDialectSQLite, c.declared, c.value.Type(), readback[i])
		require.NoError(t, err, c.name)
		require.Equal(t, c.value.Type(), got.Type(), c.name)
		require.True(t, c.value.Equal(got), "%s: want %v got %v", c.name, c.value.Value(), got.Value())
	}
}

func TestExactSQLiteRejectsLossyStrategies(t *testing.T) {
	maxUint := measurement.NewUint64(math.MaxUint64)
	_, err := EncodeExactValue(SQLDialectSQLite, "INTEGER", maxUint)
	require.ErrorIs(t, err, ErrExactSQLBlocked)
	for _, declared := range []string{"NUMERIC", "REAL", "INTEGER", "FLOAT", "DECIMAL(20,4)", "BLOB"} {
		_, err := EncodeExactValue(SQLDialectSQLite, declared, decimalValue(t, "1.25"))
		require.ErrorIs(t, err, ErrExactSQLBlocked, declared)
	}
	_, err = EncodeExactValue(SQLDialectSQLite, "REAL", measurement.NewUint64(1))
	require.ErrorIs(t, err, ErrExactSQLBlocked)
	_, err = EncodeExactValue(SQLDialectSQLite, "TEXT", measurement.NewInt64(1))
	require.ErrorIs(t, err, ErrExactSQLBlocked)
	_, err = EncodeExactValue(SQLDialectSQLite, "TEXT", measurement.NewBool(true))
	require.ErrorIs(t, err, ErrExactSQLBlocked)
	_, err = EncodeExactValue(SQLDialectSQLite, "", measurement.NewText("x"))
	require.ErrorIs(t, err, ErrExactSQLBlocked)
}

func TestExactSQLiteBindShapes(t *testing.T) {
	bound, err := EncodeExactValue(SQLDialectSQLite, "INTEGER", measurement.NewBool(true))
	require.NoError(t, err)
	require.Equal(t, int64(1), bound)
	bound, err = EncodeExactValue(SQLDialectSQLite, "INTEGER", measurement.NewUint64(math.MaxInt64))
	require.NoError(t, err)
	require.Equal(t, int64(math.MaxInt64), bound)
	bound, err = EncodeExactValue(SQLDialectSQLite, "TEXT", measurement.NewUint64(math.MaxUint64))
	require.NoError(t, err)
	require.Equal(t, "18446744073709551615", bound)
}

func TestExactPostgresCodecDeclarations(t *testing.T) {
	bound, err := EncodeExactValue(SQLDialectPostgres, "BIGINT", measurement.NewInt64(math.MinInt64))
	require.NoError(t, err)
	require.Equal(t, int64(math.MinInt64), bound)
	bound, err = EncodeExactValue(SQLDialectPostgres, "BOOLEAN", measurement.NewBool(false))
	require.NoError(t, err)
	require.Equal(t, false, bound)
	bound, err = EncodeExactValue(SQLDialectPostgres, "NUMERIC(20,0)", measurement.NewUint64(math.MaxUint64))
	require.NoError(t, err)
	require.Equal(t, "18446744073709551615", bound)
	bound, err = EncodeExactValue(SQLDialectPostgres, "NUMERIC(28,18)", decimalValue(t, "1234567890.123456789012345678"))
	require.NoError(t, err)
	require.Equal(t, "1234567890.123456789012345678", bound)
	bound, err = EncodeExactValue(SQLDialectPostgres, "numeric", decimalValue(t, "0.1"))
	require.NoError(t, err)
	require.Equal(t, "0.1", bound)
	bound, err = EncodeExactValue(SQLDialectPostgres, "NUMERIC(16,4)", decimalValue(t, "1.2500"))
	require.NoError(t, err, "trailing zeros beyond scale lose nothing")
	require.Equal(t, "1.2500", bound)
}

func TestExactPostgresRejectsRoundingAndOverflow(t *testing.T) {
	_, err := EncodeExactValue(SQLDialectPostgres, "NUMERIC(16,4)", decimalValue(t, "1234567890.123456789012345678"))
	require.ErrorIs(t, err, ErrExactSQLBlocked)
	_, err = EncodeExactValue(SQLDialectPostgres, "NUMERIC(18,0)", measurement.NewUint64(math.MaxUint64))
	require.ErrorIs(t, err, ErrExactSQLBlocked)
	_, err = EncodeExactValue(SQLDialectPostgres, "NUMERIC(5,4)", decimalValue(t, "12.5"))
	require.ErrorIs(t, err, ErrExactSQLBlocked)
	_, err = EncodeExactValue(SQLDialectPostgres, "BIGINT", measurement.NewUint64(math.MaxUint64))
	require.ErrorIs(t, err, ErrExactSQLBlocked)
	for _, declared := range []string{"INTEGER", "SMALLINT", "REAL", "MONEY", "JSONB", "CITEXT2"} {
		_, err := EncodeExactValue(SQLDialectPostgres, declared, measurement.NewInt64(1))
		require.ErrorIs(t, err, ErrExactSQLBlocked, declared)
	}
	_, err = EncodeExactValue(SQLDialectPostgres, "VARCHAR(3)", measurement.NewText("abcd"))
	require.ErrorIs(t, err, ErrExactSQLBlocked)
	_, err = EncodeExactValue("oracle", "NUMBER", measurement.NewInt64(1))
	require.ErrorIs(t, err, ErrExactSQLBlocked)
}

func TestExactCodecRejectsUnsetAndNonFiniteValues(t *testing.T) {
	_, err := EncodeExactValue(SQLDialectSQLite, "TEXT", measurement.ExactValue{})
	require.ErrorIs(t, err, measurement.ErrExactValueInvalid)
	_, err = measurement.NewFloat64(math.NaN())
	require.Error(t, err)
}

func TestExactDecodeRejectsWrongReadbackShapes(t *testing.T) {
	_, err := DecodeExactValue(SQLDialectSQLite, "TEXT", measurement.ExactUint64, "12x")
	require.Error(t, err)
	_, err = DecodeExactValue(SQLDialectSQLite, "INTEGER", measurement.ExactBool, int64(2))
	require.Error(t, err)
	_, err = DecodeExactValue(SQLDialectSQLite, "INTEGER", measurement.ExactInt64, float64(1))
	require.Error(t, err, "REAL readback must not satisfy an integer type")
	_, err = DecodeExactValue(SQLDialectSQLite, "TEXT", measurement.ExactText, nil)
	require.Error(t, err)
	got, err := DecodeExactValue(SQLDialectPostgres, "NUMERIC(28,18)", measurement.ExactDecimal, []byte("1.50"))
	require.NoError(t, err)
	require.True(t, decimalValue(t, "1.5").Equal(got))
}
