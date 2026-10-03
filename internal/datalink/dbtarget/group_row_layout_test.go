package dbtarget

import (
	"database/sql"
	"encoding/json"
	"math"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"

	"github.com/stretchr/testify/require"
)

var layoutBucket = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

func baseColumns() []ColumnInfo {
	return []ColumnInfo{
		{Name: "temperature", DataType: "TEXT", Nullable: false},
		{Name: "running", DataType: "INTEGER", Nullable: false},
		{Name: "counter", DataType: "INTEGER", Nullable: true},
		{Name: "bucket_start", DataType: "TEXT", Nullable: false},
		{Name: "record_id", DataType: "TEXT", Nullable: false},
		{Name: "provenance", DataType: "TEXT", Nullable: true},
		{Name: "entity", DataType: "TEXT", Nullable: false},
	}
}

func baseSpec() GroupRowSpec {
	return GroupRowSpec{
		Dialect: SQLDialectSQLite,
		Columns: baseColumns(),
		Members: []GroupRowMember{
			{MemberKey: "temp", Column: "temperature", Type: measurement.ExactText, Required: true},
			{MemberKey: "run", Column: "running", Type: measurement.ExactBool, Required: true},
			{MemberKey: "cnt", Column: "counter", Type: measurement.ExactUint64, Required: false},
		},
	}
}

func issueCodes(issues []LayoutIssue) []string {
	codes := make([]string, 0, len(issues))
	for _, issue := range issues {
		codes = append(codes, issue.Code)
	}
	return codes
}

func rowOutcome(members ...snapshot.MemberResult) snapshot.Outcome {
	return snapshot.Outcome{
		Kind: snapshot.OutcomeRow, RecordID: "record-1", EffectKey: "effect-1",
		BucketStart: layoutBucket, BucketEnd: layoutBucket.Add(10 * time.Second), Members: members,
	}
}

func okMember(key string, value measurement.ExactValue) snapshot.MemberResult {
	return snapshot.MemberResult{
		MemberKey: key, Status: snapshot.MemberOK,
		Sample: &snapshot.Sample{
			SampleID: "sample-" + key, MemberKey: key, ObservedAt: layoutBucket.Add(5 * time.Second),
			Quality: schema.QualityGood, Value: value,
		},
	}
}

func TestExactMixedValueRoundTripLayoutAcceptsInspectedColumns(t *testing.T) {
	layout, issues := NewGroupRowLayout(baseSpec())
	require.Empty(t, issues)
	require.NotNil(t, layout)
	require.False(t, layout.Capabilities().SupportsPartial(), "no provenance column, so partial is not supported")
}

func TestExactMixedValueRoundTripLayoutBlocksUnsupportedSQLTypesAndIdentity(t *testing.T) {
	cases := map[string]struct {
		mutate func(*GroupRowSpec)
		code   string
		column string
	}{
		"missing column": {func(s *GroupRowSpec) { s.Members[0].Column = "absent" }, "column-missing", "absent"},
		"decimal on NUMERIC": {func(s *GroupRowSpec) {
			s.Columns[0].DataType = "NUMERIC"
			s.Members[0].Type = measurement.ExactDecimal
		}, "unsupported-sql-type", "temperature"},
		"uint64 on REAL":   {func(s *GroupRowSpec) { s.Columns[2].DataType = "REAL" }, "unsupported-sql-type", "counter"},
		"bool on TEXT":     {func(s *GroupRowSpec) { s.Columns[1].DataType = "TEXT" }, "unsupported-sql-type", "running"},
		"duplicate target": {func(s *GroupRowSpec) { s.Members[1].Column = "temperature" }, "duplicate-column", "temperature"},
		"unknown dialect":  {func(s *GroupRowSpec) { s.Dialect = "oracle" }, "unsupported-dialect", ""},
		"record key not text": {func(s *GroupRowSpec) {
			s.RecordKeyColumn = "counter"
		}, "identity-column-unsupported", "counter"},
		"bucket column wrong type": {func(s *GroupRowSpec) {
			s.BucketStartColumn = "running"
		}, "identity-column-unsupported", "running"},
		"entity key unstored":    {func(s *GroupRowSpec) { s.EntityKeyed = true }, "identity-column-missing", ""},
		"identity column absent": {func(s *GroupRowSpec) { s.EntityKeyed, s.EntityKeyColumn = true, "ghost" }, "column-missing", "ghost"},
		"member without key":     {func(s *GroupRowSpec) { s.Members[0].MemberKey = "" }, "member-incomplete", ""},
	}
	for name, c := range cases {
		spec := baseSpec()
		spec.Columns = append([]ColumnInfo(nil), spec.Columns...)
		spec.Members = append([]GroupRowMember(nil), spec.Members...)
		c.mutate(&spec)
		layout, issues := NewGroupRowLayout(spec)
		require.Nil(t, layout, name)
		require.Contains(t, issueCodes(issues), c.code, name)
		for _, issue := range issues {
			if issue.Code == c.code {
				require.Equal(t, c.column, issue.Column, name)
			}
		}
	}
}

func TestExplicitPartialPolicyLayoutCapabilitiesComeFromRealColumns(t *testing.T) {
	spec := baseSpec()
	spec.Partial = true
	spec.ProvenanceColumn = "provenance"
	layout, issues := NewGroupRowLayout(spec)
	require.Empty(t, issues)
	require.True(t, layout.Capabilities().SupportsPartial())

	notNullable := baseSpec()
	notNullable.Partial = true
	notNullable.ProvenanceColumn = "provenance"
	notNullable.Columns = append([]ColumnInfo(nil), notNullable.Columns...)
	notNullable.Columns[2].Nullable = false // optional counter column cannot hold NULL
	layout, issues = NewGroupRowLayout(notNullable)
	require.Nil(t, layout)
	require.Contains(t, issueCodes(issues), "column-not-nullable")

	noProvenance := baseSpec()
	noProvenance.Partial = true
	layout, issues = NewGroupRowLayout(noProvenance)
	require.Nil(t, layout)
	require.Contains(t, issueCodes(issues), "partial-requires-provenance")

	badProvenance := baseSpec()
	badProvenance.Partial = true
	badProvenance.ProvenanceColumn = "running"
	_, issues = NewGroupRowLayout(badProvenance)
	require.Contains(t, issueCodes(issues), "identity-column-unsupported")
}

func TestExactMixedValueRoundTripEncodeRowIntoDisposableSQLite(t *testing.T) {
	spec := baseSpec()
	spec.ProvenanceColumn = "provenance"
	spec.RecordKeyColumn = "record_id"
	spec.BucketStartColumn = "bucket_start"
	spec.EntityKeyed, spec.EntityKeyColumn = true, "entity"
	layout, issues := NewGroupRowLayout(spec)
	require.Empty(t, issues)

	outcome := rowOutcome(
		okMember("temp", measurement.NewText("batch-001")),
		okMember("run", measurement.NewBool(true)),
		okMember("cnt", measurement.NewUint64(9007199254740993)),
	)
	outcome.EntityKey = "plant-1"
	encoded, err := layout.EncodeRow(outcome)
	require.NoError(t, err)
	require.Equal(t, "record-1", encoded.RecordID)
	require.Equal(t, "effect-1", encoded.EffectKey)
	require.False(t, encoded.Partial)

	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "rows.db"))
	require.NoError(t, err)
	defer db.Close()
	_, err = db.ExecContext(t.Context(), `CREATE TABLE readings (
		temperature TEXT NOT NULL, running INTEGER NOT NULL, counter INTEGER, bucket_start TEXT NOT NULL,
		record_id TEXT NOT NULL, provenance TEXT, entity TEXT NOT NULL)`)
	require.NoError(t, err)
	columns := make([]string, 0, len(encoded.Cells))
	marks := make([]string, 0, len(encoded.Cells))
	args := make([]any, 0, len(encoded.Cells))
	for _, cell := range encoded.Cells {
		columns = append(columns, cell.Column)
		marks = append(marks, "?")
		args = append(args, cell.Value)
	}
	_, err = db.ExecContext(t.Context(), "INSERT INTO readings ("+join(columns)+") VALUES ("+join(marks)+")", args...)
	require.NoError(t, err)

	var temperature, bucket, recordID, provenance, entity string
	var running, counter int64
	require.NoError(t, db.QueryRowContext(t.Context(),
		`SELECT temperature, running, counter, bucket_start, record_id, provenance, entity FROM readings`).
		Scan(&temperature, &running, &counter, &bucket, &recordID, &provenance, &entity))
	require.Equal(t, "batch-001", temperature)
	require.Equal(t, int64(1), running)
	require.Equal(t, int64(9007199254740993), counter)
	require.Equal(t, "2026-01-01T00:00:00Z", bucket)
	require.Equal(t, "record-1", recordID)
	require.Equal(t, "plant-1", entity)

	var entries []map[string]any
	require.NoError(t, json.Unmarshal([]byte(provenance), &entries))
	require.Len(t, entries, 3)
	require.Equal(t, "temp", entries[0]["member"])
	require.Equal(t, "sample-temp", entries[0]["sample_id"])
	require.Equal(t, "ok", entries[0]["status"])
}

func join(parts []string) string {
	out := ""
	for i, part := range parts {
		if i > 0 {
			out += ","
		}
		out += part
	}
	return out
}

func TestExplicitPartialPolicyEncodesNullWithReasonNeverZero(t *testing.T) {
	spec := baseSpec()
	spec.Partial = true
	spec.ProvenanceColumn = "provenance"
	layout, issues := NewGroupRowLayout(spec)
	require.Empty(t, issues)

	outcome := rowOutcome(
		okMember("temp", measurement.NewText("batch-001")),
		okMember("run", measurement.NewBool(false)),
		snapshot.MemberResult{
			MemberKey: "cnt", Status: snapshot.MemberBad, Reason: "read-failed",
			Sample: &snapshot.Sample{SampleID: "bad-1", MemberKey: "cnt", ObservedAt: layoutBucket.Add(9 * time.Second), Quality: schema.QualityBad, Value: measurement.NewUint64(77)},
		},
	)
	outcome.Partial = true
	encoded, err := layout.EncodeRow(outcome)
	require.NoError(t, err)
	require.True(t, encoded.Partial)

	cells := map[string]any{}
	for _, cell := range encoded.Cells {
		cells[cell.Column] = cell.Value
	}
	require.Nil(t, cells["counter"], "a bad reading is NULL, never zero or its stale value")
	require.Contains(t, cells, "counter")
	require.Equal(t, int64(0), cells["running"], "a real false stays false")
	var entries []map[string]any
	require.NoError(t, json.Unmarshal([]byte(cells["provenance"].(string)), &entries))
	require.Equal(t, "bad", entries[2]["status"])
	require.Equal(t, "read-failed", entries[2]["reason"])
}

func TestExactMixedValueRoundTripEncodeRowBlocksTypeMismatchAndNonRows(t *testing.T) {
	layout, issues := NewGroupRowLayout(baseSpec())
	require.Empty(t, issues)

	mismatch := rowOutcome(
		okMember("temp", measurement.NewInt64(5)),
		okMember("run", measurement.NewBool(true)),
		okMember("cnt", measurement.NewUint64(1)),
	)
	_, err := layout.EncodeRow(mismatch)
	var rowErr *GroupRowError
	require.ErrorAs(t, err, &rowErr)
	require.Equal(t, "type-mismatch", rowErr.Code)
	require.Equal(t, "temperature", rowErr.Column)
	require.NotContains(t, err.Error(), "5")

	overflow := rowOutcome(
		okMember("temp", measurement.NewText("x")),
		okMember("run", measurement.NewBool(true)),
		okMember("cnt", measurement.NewUint64(math.MaxUint64)),
	)
	_, err = layout.EncodeRow(overflow)
	require.ErrorAs(t, err, &rowErr)
	require.Equal(t, "sql-value-blocked", rowErr.Code)
	require.Equal(t, "counter", rowErr.Column)

	for _, kind := range []snapshot.OutcomeKind{snapshot.OutcomeSkipped, snapshot.OutcomeNoData} {
		_, err = layout.EncodeRow(snapshot.Outcome{Kind: kind})
		require.ErrorAs(t, err, &rowErr, string(kind))
		require.Equal(t, "not-a-row", rowErr.Code)
	}

	unexpectedNull := rowOutcome(
		okMember("temp", measurement.NewText("x")),
		okMember("run", measurement.NewBool(true)),
		snapshot.MemberResult{MemberKey: "cnt", Status: snapshot.MemberMissing, Reason: "no-samples"},
	)
	_, err = layout.EncodeRow(unexpectedNull)
	require.ErrorAs(t, err, &rowErr)
	require.Equal(t, "unexpected-null", rowErr.Code, "a non-partial layout never writes NULL for a member")

	missing := rowOutcome(okMember("temp", measurement.NewText("x")), okMember("run", measurement.NewBool(true)))
	_, err = layout.EncodeRow(missing)
	require.ErrorAs(t, err, &rowErr)
	require.Equal(t, "member-missing", rowErr.Code)
}

func TestExactMixedValueRoundTripPostgresTimestampBucketColumn(t *testing.T) {
	spec := GroupRowSpec{
		Dialect: SQLDialectPostgres,
		Columns: []ColumnInfo{
			{Name: "reading", DataType: "numeric(28,18)"},
			{Name: "ts", DataType: "timestamp with time zone"},
		},
		Members:           []GroupRowMember{{MemberKey: "r", Column: "reading", Type: measurement.ExactDecimal, Required: true}},
		BucketStartColumn: "ts",
	}
	layout, issues := NewGroupRowLayout(spec)
	require.Empty(t, issues)
	decimal, err := measurement.NewDecimal("1234567890.123456789012345678")
	require.NoError(t, err)
	encoded, err := layout.EncodeRow(rowOutcome(okMember("r", decimal)))
	require.NoError(t, err)
	cells := map[string]any{}
	for _, cell := range encoded.Cells {
		cells[cell.Column] = cell.Value
	}
	require.Equal(t, "1234567890.123456789012345678", cells["reading"])
	require.Equal(t, layoutBucket, cells["ts"], "timestamp columns receive UTC time.Time")
}

func TestGroupRowLayoutScopesMembersAndClaimsByEntity(t *testing.T) {
	spec := GroupRowSpec{
		Dialect: SQLDialectSQLite,
		Columns: []ColumnInfo{
			{Name: "shared_reading", DataType: "TEXT", Nullable: false},
			{Name: "entity", DataType: "TEXT", Nullable: false},
		},
		Members: []GroupRowMember{
			{MemberKey: "line-a-reading", EntityKey: "line-a", Column: "shared_reading", Type: measurement.ExactText, Required: true},
			{MemberKey: "line-b-reading", EntityKey: "line-b", Column: "shared_reading", Type: measurement.ExactText, Required: true},
		},
		EntityKeyed:     true,
		EntityKeyColumn: "entity",
	}
	layout, issues := NewGroupRowLayout(spec)
	require.Empty(t, issues)

	for _, test := range []struct {
		entity string
		key    string
		value  string
	}{
		{entity: "line-a", key: "line-a-reading", value: "a"},
		{entity: "line-b", key: "line-b-reading", value: "b"},
	} {
		outcome := rowOutcome(okMember(test.key, measurement.NewText(test.value)))
		outcome.EntityKey = test.entity
		encoded, err := layout.EncodeRow(outcome)
		require.NoError(t, err, test.entity)
		require.Equal(t, test.value, encoded.Cells[0].Value, test.entity)
	}

	conflict := spec
	conflict.Members = slices.Clone(spec.Members)
	conflict.Members[1].EntityKey = conflict.Members[0].EntityKey
	layout, issues = NewGroupRowLayout(conflict)
	require.Nil(t, layout)
	require.Contains(t, issueCodes(issues), "duplicate-column")
}
