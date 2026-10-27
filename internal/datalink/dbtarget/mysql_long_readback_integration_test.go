package dbtarget

// TestMySQLLongExactValuesReadback exercises the managed_long_v1 MySQL/MariaDB
// plain layout end to end on the disposable gw-circuit fixtures: validated
// layout, preview DDL, bound inserts, exact readback. This is fixture evidence
// for the typed plain slice only; journal, outbox, ledger, partition, repair,
// the profile engine and durable delivery stay outside this slice.

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"math"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
)

// Fixed expectations. None of them are derived from the functions under test.
// record_id values are fixed 64-hex literals matching the managed_long_v1
// record identity shape; Unicode/case-distinct identities live in tag_id.
const (
	wantCaseID        = "IdA-座標-Γ1"
	wantCaseVariantID = "ida-座標-γ1"
	wantTrailingTag   = "tag-pad  "
	wantUnicodeText   = "溫度_Γραμμή🌡"
	wantCaseVarText   = "溫度_γραμμή🌡"
	wantTrailingText  = "ok  "
	wantIntBase       = int64(9007199254740992)
	wantIntNext       = int64(9007199254740993)
	wantDecBase       = "12000.0"
	wantDecNext       = "12000.4"
	wantDecDiff       = "0.4"
	wantDecBaseCan    = "12000"
	wantDecNextCan    = "12000.4"
	wantDecMaxInt     = "99999999999999999999999999999999999" // 35 integer digits
	wantDecMaxFrac    = "0.123456789012345678901234567890"    // 30 fraction digits
	wantDecMaxFracCan = "0.12345678901234567890123456789"
	wantMaxUintDigits = "18446744073709551615"
	wantProvenance    = `{"seq":1,"src":"gw_ltest"}`

	wantRecordIDA = "1111111111111111111111111111111111111111111111111111111111111111"
	wantRecordIDB = "2222222222222222222222222222222222222222222222222222222222222222"
	wantRecordIDC = "3333333333333333333333333333333333333333333333333333333333333333"
	wantRecordIDD = "4444444444444444444444444444444444444444444444444444444444444444"
	wantRecordIDE = "5555555555555555555555555555555555555555555555555555555555555555"
	wantRecordIDF = "6666666666666666666666666666666666666666666666666666666666666666"
	wantRecordIDG = "7777777777777777777777777777777777777777777777777777777777777777"
	wantRecordIDH = "8888888888888888888888888888888888888888888888888888888888888888"
	wantRecordIDI = "9999999999999999999999999999999999999999999999999999999999999999"
	wantRecordIDJ = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	wantRecordIDK = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	wantRecordIDL = "cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"
)

func mysqlFamilyCollationFor(family string) string {
	if family == mariaDBFixtureFamily {
		return "utf8mb4_nopad_bin"
	}
	return "utf8mb4_0900_bin"
}

func TestMySQLLongExactValuesReadback(t *testing.T) {
	t.Run("mysql", func(t *testing.T) { runMySQLLongExactValuesReadback(t, mysqlFixtureFamily) })
	t.Run("mariadb", func(t *testing.T) { runMySQLLongExactValuesReadback(t, mariaDBFixtureFamily) })
}

// longExactRow is one managed long record: exactly one raw value slot plus the
// typed record metadata. An all-slot-empty row is the typed-NULL record.
type longExactRow struct {
	id        string // fixed 64-hex record identity
	bucket    time.Time
	valueType string
	observed  time.Time // zero means NULL
	intValue  *int64
	uintValue *uint64
	decValue  string // empty means NULL
	boolValue *bool
	textValue string // empty means NULL
	tag       string
}

func int64Ptr(v int64) *int64    { return &v }
func uint64Ptr(v uint64) *uint64 { return &v }
func boolPtr(v bool) *bool       { return &v }

func (r longExactRow) rawSlots() int {
	count := 0
	for _, set := range []bool{r.intValue != nil, r.uintValue != nil, r.decValue != "", r.boolValue != nil, r.textValue != ""} {
		if set {
			count++
		}
	}
	return count
}

// longExactRows pins the exact-value matrix as legal typed rows: the int64
// pair differing by 1 past the float64 integer range, the decimal pair
// differing by 0.4, both DECIMAL(65,30) boundary shapes, the uint64 maximum,
// both bools, exact Unicode and trailing-space texts, and one typed-NULL
// record. Case-exact and no-pad identity is carried by tag_id.
func longExactRows() []longExactRow {
	return []longExactRow{
		{id: wantRecordIDA, valueType: "int64", observed: time.Date(2026, 11, 6, 8, 30, 15, 123456000, time.FixedZone("UTC+8", 8*3600)),
			intValue: int64Ptr(wantIntBase), decValue: "", tag: wantCaseID},
		{id: wantRecordIDB, valueType: "int64", observed: time.Date(2026, 11, 6, 0, 31, 45, 999999000, time.UTC),
			intValue: int64Ptr(wantIntNext), tag: wantCaseVariantID},
		{id: wantRecordIDC, valueType: "decimal", decValue: wantDecBase, tag: "tag-座標-3"},
		{id: wantRecordIDD, valueType: "decimal", decValue: wantDecNext, tag: "tag-座標-4"},
		{id: wantRecordIDE, valueType: "uint64", uintValue: uint64Ptr(math.MaxUint64), tag: "tag-uint64"},
		{id: wantRecordIDF, valueType: "decimal", decValue: wantDecMaxInt, tag: "tag-35int"},
		{id: wantRecordIDG, valueType: "decimal", decValue: wantDecMaxFrac, tag: "tag-30frac"},
		{id: wantRecordIDH, valueType: "bool", boolValue: boolPtr(true), tag: wantTrailingTag},
		{id: wantRecordIDI, valueType: "bool", boolValue: boolPtr(false), tag: "tag-false"},
		{id: wantRecordIDJ, valueType: "text", textValue: wantUnicodeText, tag: "tag-text"},
		{id: wantRecordIDK, valueType: "text", textValue: wantTrailingText, tag: "tag-trailtext"},
		{id: wantRecordIDL, valueType: "int64", tag: "tag-null"}, // typed NULL record
	}
}

func runMySQLLongExactValuesReadback(t *testing.T, family string) {
	cfg := loadMySQLLongFixtureConfig(t, family)
	db := openMySQLLongFixtureDB(t, cfg)
	version := requireMySQLFixtureFamily(t, db, family)

	options := recordingplan.MySQLGroupSchemaOptions{
		ServerFamily:       family,
		ServerVersion:      version,
		CapabilityRevision: family + "-unverified-1", // fixture marker, never a production capability release
		Engine:             "InnoDB",
		Collation:          mysqlFamilyCollationFor(family),
	}
	table := uniqueLongTestTable(t)
	owner := uniqueLongTestOwnerColumn(t)
	layout, err := recordingplan.ManagedMySQLLongGroupLayout(table, owner, options)
	require.NoError(t, err)
	stmts, err := recordingplan.GenerateGroupSchemaDDL("mysql", layout)
	require.NoError(t, err)

	// Cleanup is registered before any DDL and owns only the tables this test
	// successfully created. A pre-existing receipt table is never dropped or
	// altered; ownership is tracked per executed statement, never by a
	// before/after table scan.
	ownsReceipt := !fixtureTableExists(t, db, longTestReceiptTable)
	var owned []string
	t.Cleanup(func() { dropFixtureTablesOwned(t, db, owned) })
	execManagedDDL(t, db, stmts, table, ownsReceipt, &owned)

	rows := longExactRows()
	base := time.Date(2026, 11, 6, 0, 30, 15, 123456000, time.UTC)
	for i := range rows {
		rows[i].bucket = base.Add(time.Duration(i) * time.Hour)
		require.Len(t, rows[i].id, 64)
		require.Equal(t, "", strings.Trim(rows[i].id, "0123456789abcdef"), "record id must be 64-hex")
		require.LessOrEqual(t, rows[i].rawSlots(), 1, "managed long record has one raw value slot")
	}
	// The observed_at source is a fixed non-UTC wall clock; it must map to the
	// known UTC microsecond literal.
	require.True(t, rows[0].observed.In(time.UTC).Equal(rows[0].bucket),
		"non-UTC source must convert to the known UTC microsecond literal")

	types := layoutColumnTypes(t, layout)
	insertLongExactRows(t, db, table, types, rows)
	readbackLongExactRows(t, db, table, owner, types, rows)

	// Case-exact and no-pad equality is proven through tag_id equality reads:
	// each distinct-case tag matches exactly its own row, and a trailing-space
	// tag matches neither its trimmed form nor any other row.
	requireTagCount(t, db, table, wantCaseID, 1)
	requireTagCount(t, db, table, wantCaseVariantID, 1)
	requireTagCount(t, db, table, "tag-pad", 0)
	requireTagCount(t, db, table, wantTrailingTag, 1)

	var total int
	require.NoError(t, db.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM `"+table+"`").Scan(&total))
	require.Equal(t, len(rows), total)
}

func insertLongExactRows(t *testing.T, db *sql.DB, table string, types map[string]string, rows []longExactRow) {
	t.Helper()
	query := "INSERT INTO `" + table + "` (bucket_start, bucket_end, record_id, workspace_id, group_id, " +
		"group_revision, device_id, point_id, tag_id, profile_kind, value_type, status, observed_at, " +
		"value_int, value_decimal, value_bool, value_text, provenance) " +
		"VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	for _, row := range rows {
		var valueInt, valueDec, valueBool, valueText, observed any
		if row.intValue != nil {
			valueInt = encodeExact(t, types, "value_int", measurement.NewInt64(*row.intValue))
		}
		if row.uintValue != nil {
			valueDec = encodeExact(t, types, "value_decimal", measurement.NewUint64(*row.uintValue))
		}
		if row.decValue != "" {
			valueDec = encodeExact(t, types, "value_decimal", requireDecimal(t, row.decValue))
		}
		if row.boolValue != nil {
			valueBool = encodeExact(t, types, "value_bool", measurement.NewBool(*row.boolValue))
		}
		if row.textValue != "" {
			valueText = encodeExact(t, types, "value_text", measurement.NewText(row.textValue))
		}
		if !row.observed.IsZero() {
			observed = row.observed
		}
		status := "ok"
		if row.rawSlots() == 0 {
			status = "missing"
		}
		_, err := db.ExecContext(t.Context(), query, row.bucket, row.bucket.Add(time.Hour), row.id,
			"gw_ltest_ws", "gw_ltest_group", "rev-1", "dev-1", "pt-1", row.tag,
			"instantaneous", row.valueType, status, observed, valueInt, valueDec, valueBool, valueText, wantProvenance)
		require.NoError(t, err, "insert %s", row.id)
	}
}

func readbackLongExactRows(t *testing.T, db *sql.DB, table, owner string, types map[string]string, rows []longExactRow) {
	t.Helper()
	for _, row := range rows {
		var observed sql.Null[time.Time]
		var valueInt sql.Null[int64]
		var valueDec sql.Null[string]
		var valueBool sql.Null[int64]
		var valueText sql.Null[string]
		var provenance []byte
		var ownerValue, tagValue string
		err := db.QueryRowContext(t.Context(),
			"SELECT observed_at, value_int, value_decimal, value_bool, value_text, provenance, tag_id, `"+owner+"` "+
				"FROM `"+table+"` WHERE bucket_start = ? AND record_id = ?", row.bucket, row.id).
			Scan(&observed, &valueInt, &valueDec, &valueBool, &valueText, &provenance, &tagValue, &ownerValue)
		require.NoError(t, err, "readback %s", row.id)
		require.Equal(t, longTestOwnerDefault, ownerValue, "owner marker for %s", row.id)
		require.Equal(t, row.tag, tagValue, "tag_id readback for %s", row.id)

		if row.observed.IsZero() {
			require.False(t, observed.Valid, "observed_at must stay NULL for %s", row.id)
		} else {
			require.True(t, observed.Valid, "observed_at lost for %s", row.id)
			require.Equal(t, time.UTC, observed.V.Location(), "readback time must be UTC for %s", row.id)
			require.True(t, observed.V.Equal(row.observed), "observed_at microsecond drift for %s", row.id)
		}

		switch {
		case row.intValue != nil:
			require.Equal(t, measurement.NewInt64(*row.intValue),
				decodeExact(t, types, "value_int", measurement.ExactInt64, valueInt.V), "int64 %s", row.id)
		case row.uintValue != nil:
			gotUint := decodeExact(t, types, "value_decimal", measurement.ExactUint64, valueDec.V)
			require.True(t, gotUint.Equal(measurement.NewUint64(*row.uintValue)), "uint64 max %s", row.id)
			gotDigits, ok := gotUint.CanonicalDigits()
			require.True(t, ok)
			require.Equal(t, wantMaxUintDigits, gotDigits, "uint64 max canonical digits")
		default:
			require.False(t, valueInt.Valid, "value_int must stay NULL for %s", row.id)
		}

		switch {
		case row.uintValue != nil:
			// Already asserted in the uint64 branch above.
		case row.decValue != "":
			got := decodeExact(t, types, "value_decimal", measurement.ExactDecimal, valueDec.V)
			require.True(t, got.Equal(requireDecimal(t, row.decValue)),
				"decimal %s: %v != %s", row.id, got.Value(), row.decValue)
			gotDigits, ok := got.CanonicalDigits()
			require.True(t, ok)
			switch row.decValue {
			case wantDecBase:
				require.Equal(t, wantDecBaseCan, gotDigits)
			case wantDecNext:
				require.Equal(t, wantDecNextCan, gotDigits)
			case wantDecMaxInt:
				require.Equal(t, wantDecMaxInt, gotDigits)
			case wantDecMaxFrac:
				require.Equal(t, wantDecMaxFracCan, gotDigits)
			}
		default:
			require.False(t, valueDec.Valid, "value_decimal must stay NULL for %s", row.id)
		}

		if row.boolValue != nil {
			require.Equal(t, measurement.NewBool(*row.boolValue),
				decodeExact(t, types, "value_bool", measurement.ExactBool, valueBool.V), "bool %s", row.id)
		} else {
			require.False(t, valueBool.Valid, "value_bool must stay NULL for %s", row.id)
		}

		if row.textValue != "" {
			require.Equal(t, measurement.NewText(row.textValue),
				decodeExact(t, types, "value_text", measurement.ExactText, valueText.V), "text %s", row.id)
		} else {
			require.False(t, valueText.Valid, "value_text must stay NULL for %s", row.id)
		}

		var wantProvenanceAny, gotProvenanceAny any
		require.NoError(t, json.Unmarshal([]byte(wantProvenance), &wantProvenanceAny))
		require.NoError(t, json.Unmarshal(provenance, &gotProvenanceAny), "provenance %s", row.id)
		require.Equal(t, wantProvenanceAny, gotProvenanceAny, "provenance semantic equality for %s", row.id)
	}

	baseDec := decodeExact(t, types, "value_decimal", measurement.ExactDecimal,
		decimalReadback(t, db, table, wantRecordIDC, rows[2].bucket))
	nextDec := decodeExact(t, types, "value_decimal", measurement.ExactDecimal,
		decimalReadback(t, db, table, wantRecordIDD, rows[3].bucket))
	baseRat, ok := baseDec.Rat()
	require.True(t, ok)
	nextRat, ok := nextDec.Rat()
	require.True(t, ok)
	wantRat, ok := requireDecimal(t, wantDecDiff).Rat()
	require.True(t, ok)
	require.Equal(t, 0, new(big.Rat).Sub(nextRat, baseRat).Cmp(wantRat),
		"decimal difference must stay exactly 0.4")
	var baseInt, nextInt int64
	require.NoError(t, db.QueryRowContext(t.Context(), "SELECT value_int FROM `"+table+"` WHERE record_id = ?", wantRecordIDA).Scan(&baseInt))
	require.NoError(t, db.QueryRowContext(t.Context(), "SELECT value_int FROM `"+table+"` WHERE record_id = ?", wantRecordIDB).Scan(&nextInt))
	require.Equal(t, int64(1), nextInt-baseInt, "SQL readback integer pair differs by exactly 1")
}

func requireTagCount(t *testing.T, db *sql.DB, table, tag string, want int) {
	t.Helper()
	var count int
	require.NoError(t, db.QueryRowContext(t.Context(),
		"SELECT COUNT(*) FROM `"+table+"` WHERE tag_id = ?", tag).Scan(&count))
	require.Equal(t, want, count, "tag_id equality under the family binary collation for %q", tag)
}

func decimalReadback(t *testing.T, db *sql.DB, table, id string, bucket time.Time) []byte {
	t.Helper()
	var raw []byte
	require.NoError(t, db.QueryRowContext(t.Context(),
		"SELECT value_decimal FROM `"+table+"` WHERE bucket_start = ? AND record_id = ?", bucket, id).Scan(&raw))
	return raw
}

func layoutColumnTypes(t *testing.T, layout recordingplan.GroupSchemaLayout) map[string]string {
	t.Helper()
	types := make(map[string]string, len(layout.Columns))
	for _, column := range layout.Columns {
		types[column.Name] = column.SQLType
	}
	return types
}

func declaredTypeOf(t *testing.T, types map[string]string, column string) string {
	t.Helper()
	declared, ok := types[column]
	require.True(t, ok, "layout column %q missing", column)
	return declared
}

func requireDecimal(t *testing.T, digits string) measurement.ExactValue {
	t.Helper()
	value, err := measurement.NewDecimal(digits)
	require.NoError(t, err)
	return value
}

func encodeExact(t *testing.T, types map[string]string, column string, value measurement.ExactValue) driver.Value {
	t.Helper()
	encoded, err := EncodeExactValue(SQLDialectMySQL, declaredTypeOf(t, types, column), value)
	require.NoError(t, err)
	return encoded
}

func decodeExact(t *testing.T, types map[string]string, column string, expected measurement.ExactType, raw any) measurement.ExactValue {
	t.Helper()
	decoded, err := DecodeExactValue(SQLDialectMySQL, declaredTypeOf(t, types, column), expected, raw)
	require.NoError(t, err)
	return decoded
}

// execManagedDDL runs the preview statements and records, per executed
// statement, only the tables this test created. A mid-batch failure still
// leaves every already-created object tracked for cleanup.
func execManagedDDL(t *testing.T, db *sql.DB, statements []string, table string, ownsReceipt bool, owned *[]string) {
	t.Helper()
	for i, statement := range statements {
		isReceipt := strings.HasPrefix(statement, "CREATE TABLE `"+longTestReceiptTable+"`")
		if isReceipt && !ownsReceipt {
			continue // This test neither creates nor changes a pre-existing receipt.
		}
		if _, err := db.ExecContext(t.Context(), statement); err != nil {
			t.Fatalf("managed DDL statement %d failed: %v", i, err)
		}
		if strings.Contains(statement, "CREATE TABLE `"+table+"`") {
			*owned = append(*owned, table)
		}
		if isReceipt {
			*owned = append(*owned, longTestReceiptTable)
		}
	}
}
