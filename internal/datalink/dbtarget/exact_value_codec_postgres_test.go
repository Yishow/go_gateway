package dbtarget

import (
	"math"
	"strconv"
	"testing"
	"time"

	"go-gateway/internal/datalink/measurement"

	"github.com/stretchr/testify/require"
)

// Live PostgreSQL fixture for the exact codec. It needs POSTGRES_DSN and uses
// an isolated schema that is dropped afterwards.
func TestExactMixedValueRoundTripLivePostgres(t *testing.T) {
	dsn, _ := postgresKeywordDSN(t)
	admin := postgresAdminDB(t, dsn)
	schemaName := postgresScopedSchema(t, admin)
	table := pgIdentifier(t, schemaName) + `.exact_row`
	_, err := admin.ExecContext(t.Context(), `CREATE TABLE `+table+` (
		running BOOLEAN, batch TEXT, minimum BIGINT, big NUMERIC(20,0), reading NUMERIC(28,18),
		narrow NUMERIC(16,4), bucket TIMESTAMPTZ)`)
	require.NoError(t, err)

	decimal, err := measurement.NewDecimal("1234567890.123456789012345678")
	require.NoError(t, err)
	type column struct {
		name, declared string
		value          measurement.ExactValue
	}
	columns := []column{
		{"running", "BOOLEAN", measurement.NewBool(true)},
		{"batch", "TEXT", measurement.NewText("batch-001")},
		{"minimum", "BIGINT", measurement.NewInt64(math.MinInt64)},
		{"big", "NUMERIC(20,0)", measurement.NewUint64(math.MaxUint64)},
		{"reading", "NUMERIC(28,18)", decimal},
	}
	args := make([]any, 0, len(columns)+1)
	names := ""
	marks := ""
	for i, c := range columns {
		bound, err := EncodeExactValue(SQLDialectPostgres, c.declared, c.value)
		require.NoError(t, err, c.name)
		args = append(args, bound)
		if i > 0 {
			names += ","
			marks += ","
		}
		names += c.name
		marks += "$" + strconv.Itoa(i+1)
	}
	bucket := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	args = append(args, bucket)
	_, err = admin.ExecContext(t.Context(), `INSERT INTO `+table+` (`+names+`,bucket) VALUES (`+marks+`,$6)`, args...)
	require.NoError(t, err)

	readback := make([]any, len(columns))
	dest := make([]any, len(columns))
	for i := range dest {
		dest[i] = &readback[i]
	}
	require.NoError(t, admin.QueryRowContext(t.Context(),
		`SELECT running,batch,minimum,big::text,reading::text FROM `+table).Scan(dest...))
	for i, c := range columns {
		got, err := DecodeExactValue(SQLDialectPostgres, c.declared, c.value.Type(), readback[i])
		require.NoError(t, err, c.name)
		require.True(t, c.value.Equal(got), "%s: want %v got %v", c.name, c.value.Value(), got.Value())
	}
	var storedBucket time.Time
	require.NoError(t, admin.QueryRowContext(t.Context(), `SELECT bucket FROM `+table).Scan(&storedBucket))
	require.True(t, bucket.Equal(storedBucket))
}

// The gate matters: PostgreSQL itself silently rounds a NUMERIC(16,4) value,
// so the codec must refuse before the driver ever sees it.
func TestExactPostgresNumericScaleRoundsSilentlyWithoutTheGate(t *testing.T) {
	dsn, _ := postgresKeywordDSN(t)
	admin := postgresAdminDB(t, dsn)
	schemaName := postgresScopedSchema(t, admin)
	table := pgIdentifier(t, schemaName) + `.narrow_row`
	_, err := admin.ExecContext(t.Context(), `CREATE TABLE `+table+` (narrow NUMERIC(16,4))`)
	require.NoError(t, err)

	lossy := "1234567890.123456789012345678"
	_, err = admin.ExecContext(t.Context(), `INSERT INTO `+table+` (narrow) VALUES ($1)`, lossy)
	require.NoError(t, err, "the database accepts the value and rounds it")
	var stored string
	require.NoError(t, admin.QueryRowContext(t.Context(), `SELECT narrow::text FROM `+table).Scan(&stored))
	require.Equal(t, "1234567890.1235", stored)

	decimal, err := measurement.NewDecimal(lossy)
	require.NoError(t, err)
	_, err = EncodeExactValue(SQLDialectPostgres, "NUMERIC(16,4)", decimal)
	require.ErrorIs(t, err, ErrExactSQLBlocked)

	// A real PostgreSQL integer column also refuses to hold MaxUint64.
	_, err = admin.ExecContext(t.Context(), `CREATE TABLE `+pgIdentifier(t, schemaName)+`.int_row (big BIGINT)`)
	require.NoError(t, err)
	_, err = admin.ExecContext(t.Context(), `INSERT INTO `+pgIdentifier(t, schemaName)+`.int_row (big) VALUES ($1)`, "18446744073709551615")
	require.Error(t, err)
	_, err = EncodeExactValue(SQLDialectPostgres, "BIGINT", measurement.NewUint64(math.MaxUint64))
	require.ErrorIs(t, err, ErrExactSQLBlocked)
}
