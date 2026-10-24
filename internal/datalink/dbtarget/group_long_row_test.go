package dbtarget

import (
	"testing"
	"time"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"

	"github.com/stretchr/testify/require"
)

var longWindow = time.Date(2026, 10, 24, 0, 0, 15, 0, time.UTC)

// longColumnsFor returns declared-type column metadata of the fixed layout.
func longColumnsFor(t *testing.T, dialect string) []ColumnInfo {
	t.Helper()
	layout, err := recordingplan.ManagedLongGroupLayout(dialect, indexTestTable, indexTestOwner)
	require.NoError(t, err)
	columns := make([]ColumnInfo, len(layout.Columns))
	for i, column := range layout.Columns {
		columns[i] = ColumnInfo{Name: column.Name, DataType: column.SQLType, Nullable: column.Nullable, PrimaryKey: column.PrimaryKey}
	}
	return columns
}

type longMemberCase struct {
	member LongRowMember
	value  measurement.ExactValue
}

func longCases(t *testing.T) []longMemberCase {
	t.Helper()
	decimal, err := measurement.NewDecimal("12000.4")
	require.NoError(t, err)
	float, err := measurement.NewFloat64(0.1)
	require.NoError(t, err)
	return []longMemberCase{
		{LongRowMember{MemberKey: "m-int", DeviceID: "dev", PointID: "p-int", TagID: "t-int", CircuitID: "circuit-a", Type: measurement.ExactInt64}, measurement.NewInt64(9007199254740993)},
		{LongRowMember{MemberKey: "m-uint", DeviceID: "dev", PointID: "p-uint", TagID: "t-uint", MeasurementID: "measurement-1", Type: measurement.ExactUint64}, measurement.NewUint64(18446744073709551615)},
		{LongRowMember{MemberKey: "m-dec", DeviceID: "dev", PointID: "p-dec", TagID: "t-dec", Type: measurement.ExactDecimal}, decimal},
		{LongRowMember{MemberKey: "m-float", DeviceID: "dev", PointID: "p-float", TagID: "t-float", Type: measurement.ExactFloat64}, float},
		{LongRowMember{MemberKey: "m-bool", DeviceID: "dev", PointID: "p-bool", TagID: "t-bool", Type: measurement.ExactBool}, measurement.NewBool(true)},
		{LongRowMember{MemberKey: "m-text", DeviceID: "dev", PointID: "p-text", TagID: "t-text", Type: measurement.ExactText}, measurement.NewText("run-42")},
	}
}

func longSpec(dialect SQLDialect, columns []ColumnInfo, cases []longMemberCase) LongRowSpec {
	members := make([]LongRowMember, len(cases))
	for i, c := range cases {
		members[i] = c.member
	}
	return LongRowSpec{
		Dialect: dialect, Columns: columns, WorkspaceID: "ws", GroupID: "g1", GroupRevision: "r1",
		ProfileKind: "instantaneous", ProfileInterval: 15, ProfileUnit: "second", Members: members,
	}
}

func longOutcome(member LongRowMember, status snapshot.MemberStatus, reason string, sample *snapshot.Sample) snapshot.Outcome {
	return snapshot.Outcome{
		Kind: snapshot.OutcomeRow, EntityKey: member.MemberKey, RecordID: "rec-" + member.MemberKey, EffectKey: "eff-" + member.MemberKey,
		BucketStart: longWindow, BucketEnd: longWindow.Add(15 * time.Second),
		Members: []snapshot.MemberResult{{MemberKey: member.MemberKey, Status: status, Reason: reason, Sample: sample}},
	}
}

func longSample(id string, value measurement.ExactValue) *snapshot.Sample {
	return &snapshot.Sample{
		SampleID: id, ObservedAt: longWindow.Add(14*time.Second + 123456789*time.Nanosecond),
		ReceivedAt: longWindow.Add(14*time.Second + 500*time.Millisecond), Quality: schema.QualityGood, Value: value,
		SourceRevision: "src-1", MappingRevision: "map-1",
	}
}

func cellMap(row EncodedRow) map[string]any {
	cells := map[string]any{}
	for _, cell := range row.Cells {
		cells[cell.Column] = cell.Value
	}
	return cells
}

func TestLongRowLayoutRejectsIncompatibleTables(t *testing.T) {
	cases := longCases(t)
	columns := longColumnsFor(t, "sqlite")
	_, issues := NewLongRowLayout(longSpec(SQLDialectSQLite, columns, cases))
	require.Empty(t, issues)

	mutate := func(edit func([]ColumnInfo) []ColumnInfo) []LayoutIssue {
		edited := edit(append([]ColumnInfo(nil), columns...))
		_, issues := NewLongRowLayout(longSpec(SQLDialectSQLite, edited, cases))
		return issues
	}
	codes := func(issues []LayoutIssue) []string {
		out := make([]string, len(issues))
		for i, issue := range issues {
			out[i] = issue.Code + ":" + issue.Column
		}
		return out
	}
	require.Contains(t, codes(mutate(func(c []ColumnInfo) []ColumnInfo { return c[1:] })), "column-missing:bucket_start")
	require.Contains(t, codes(mutate(func(c []ColumnInfo) []ColumnInfo {
		for i := range c {
			if c[i].Name == "value_int" {
				c[i].DataType = "TEXT"
			}
		}
		return c
	})), "unsupported-sql-type:value_int", "an integer slot must be able to hold int64 exactly")
	require.Contains(t, codes(mutate(func(c []ColumnInfo) []ColumnInfo {
		for i := range c {
			if c[i].Name == "value_float" {
				c[i].DataType = "INTEGER"
			}
		}
		return c
	})), "unsupported-sql-type:value_float", "a float slot must not be an integer column")
	require.Contains(t, codes(mutate(func(c []ColumnInfo) []ColumnInfo {
		for i := range c {
			if c[i].Name == "value_decimal" {
				c[i].DataType = "REAL"
			}
		}
		return c
	})), "unsupported-sql-type:value_decimal", "exact decimals and uint64 must never land in a float column")
	require.Contains(t, codes(mutate(func(c []ColumnInfo) []ColumnInfo {
		for i := range c {
			if c[i].Name == "value_text" {
				c[i].Nullable = false
			}
		}
		return c
	})), "column-not-nullable:value_text", "a value slot that rejects NULL cannot store a missing reading")

	incomplete := longSpec(SQLDialectSQLite, columns, cases)
	incomplete.GroupRevision = ""
	_, issues = NewLongRowLayout(incomplete)
	require.Contains(t, codes(issues), "identity-column-missing:")
	incomplete = longSpec(SQLDialectSQLite, columns, cases)
	incomplete.Members[0].Type = "complex"
	_, issues = NewLongRowLayout(incomplete)
	require.Contains(t, codes(issues), "member-incomplete:")
	_, issues = NewLongRowLayout(longSpec("oracle", columns, cases))
	require.Contains(t, codes(issues), "unsupported-dialect:")
}

func TestLongRowEncodingUsesOneSlotAndNormalizedTimes(t *testing.T) {
	cases := longCases(t)
	slots := map[string]string{"m-int": "value_int", "m-uint": "value_decimal", "m-dec": "value_decimal", "m-float": "value_float", "m-bool": "value_bool", "m-text": "value_text"}
	for _, dialect := range []SQLDialect{SQLDialectSQLite, SQLDialectPostgres} {
		t.Run(string(dialect), func(t *testing.T) {
			layout, issues := NewLongRowLayout(longSpec(dialect, longColumnsFor(t, string(dialect)), cases))
			require.Empty(t, issues)
			for _, c := range cases {
				sample := longSample("s-"+c.member.MemberKey, c.value)
				row, err := layout.EncodeRow(longOutcome(c.member, snapshot.MemberOK, "", sample))
				require.NoError(t, err, c.member.MemberKey)
				cells := cellMap(row)
				require.Len(t, row.Cells, len(longColumnRoles)+1, "all fixed columns and the owner marker are addressed explicitly")
				require.Equal(t, "gateway", cells[indexTestOwner])
				var filled []string
				for _, slot := range []string{"value_int", "value_decimal", "value_float", "value_bool", "value_text"} {
					if cells[slot] != nil {
						filled = append(filled, slot)
					}
				}
				require.Equal(t, []string{slots[c.member.MemberKey]}, filled, "exactly the matching typed slot is filled")
				require.Equal(t, string(c.member.Type), cells["value_type"])
				require.Equal(t, "ok", cells["status"])
				require.Equal(t, "instantaneous", cells["profile_kind"])
				require.Nil(t, cells["delta_value"])
				require.Nil(t, cells["delta_from"])
				require.Nil(t, cells["delta_to"])
				require.Equal(t, "rec-"+c.member.MemberKey, row.RecordID)
				require.Equal(t, "rec-"+c.member.MemberKey, cells["record_id"])
				if dialect == SQLDialectPostgres {
					require.Equal(t, longWindow, cells["bucket_start"])
					require.Equal(t, longWindow.Add(14*time.Second+123456*time.Microsecond), cells["observed_at"], "nanoseconds truncate to microseconds in UTC")
				} else {
					require.Equal(t, "2026-10-24T00:00:15.000000Z", cells["bucket_start"])
					require.Equal(t, "2026-10-24T00:00:29.123456Z", cells["observed_at"])
					require.Equal(t, "2026-10-24T00:00:29.500000Z", cells["received_at"])
				}
			}
			uint64Cells := cellMap(mustEncode(t, layout, cases[1].member, snapshot.MemberOK, "", longSample("u", cases[1].value)))
			require.Equal(t, "18446744073709551615", uint64Cells["value_decimal"], "uint64 is carried as exact digits")
			require.Equal(t, "measurement-1", uint64Cells["measurement_id"])
			require.Nil(t, uint64Cells["circuit_id"])
		})
	}

	t.Run("a non-UTC instant is normalized, the zero instant stays NULL", func(t *testing.T) {
		layout, issues := NewLongRowLayout(longSpec(SQLDialectSQLite, longColumnsFor(t, "sqlite"), cases))
		require.Empty(t, issues)
		sample := longSample("tz", cases[0].value)
		sample.ObservedAt = time.Date(2026, 10, 24, 8, 0, 29, 0, time.FixedZone("UTC+8", 8*3600))
		sample.ReceivedAt = time.Time{}
		cells := cellMap(mustEncode(t, layout, cases[0].member, snapshot.MemberOK, "", sample))
		require.Equal(t, "2026-10-24T00:00:29.000000Z", cells["observed_at"])
		require.Nil(t, cells["received_at"], "an unknown receive time is NULL, never invented")
	})

	t.Run("an unusable reading is NULL with its status and reason, and a missing one has no times", func(t *testing.T) {
		layout, issues := NewLongRowLayout(longSpec(SQLDialectSQLite, longColumnsFor(t, "sqlite"), cases))
		require.Empty(t, issues)
		failed := longSample("bad", measurement.ExactValue{})
		failed.Quality = schema.QualityBad
		for _, tc := range []struct {
			status snapshot.MemberStatus
			reason string
			sample *snapshot.Sample
		}{
			{snapshot.MemberBad, "read-timeout", failed},
			{snapshot.MemberStale, "max-age-exceeded", longSample("old", cases[0].value)},
			{snapshot.MemberInvalid, "type-mismatch", failed},
			{snapshot.MemberMissing, "no-samples", nil},
		} {
			cells := cellMap(mustEncode(t, layout, cases[0].member, tc.status, tc.reason, tc.sample))
			require.Equal(t, string(tc.status), cells["status"])
			for _, slot := range []string{"value_int", "value_decimal", "value_float", "value_bool", "value_text"} {
				require.Nil(t, cells[slot], "%s %s", tc.status, slot)
			}
			require.Contains(t, cells["provenance"], tc.reason)
			require.Equal(t, "int64", cells["value_type"], "the declared source type survives")
			if tc.sample == nil {
				require.Nil(t, cells["observed_at"])
				require.Nil(t, cells["quality"])
			}
		}
	})

	t.Run("structural mistakes are refused with safe codes", func(t *testing.T) {
		layout, issues := NewLongRowLayout(longSpec(SQLDialectSQLite, longColumnsFor(t, "sqlite"), cases))
		require.Empty(t, issues)
		skipped := longOutcome(cases[0].member, snapshot.MemberOK, "", longSample("s", cases[0].value))
		skipped.Kind = snapshot.OutcomeSkipped
		_, err := layout.EncodeRow(skipped)
		require.Equal(t, rowErrNotARow, rowCode(t, err))
		unknown := longOutcome(LongRowMember{MemberKey: "stranger"}, snapshot.MemberOK, "", longSample("s", cases[0].value))
		_, err = layout.EncodeRow(unknown)
		require.Equal(t, rowErrMemberMissing, rowCode(t, err))
		wrongType := longOutcome(cases[0].member, snapshot.MemberOK, "", longSample("s", measurement.NewText("not a number")))
		_, err = layout.EncodeRow(wrongType)
		require.Equal(t, rowErrTypeMismatch, rowCode(t, err))
		anonymous := longOutcome(cases[0].member, snapshot.MemberOK, "", longSample("s", cases[0].value))
		anonymous.EffectKey = ""
		_, err = layout.EncodeRow(anonymous)
		require.Equal(t, rowErrNotARow, rowCode(t, err), "a record without an effect identity is never sent")
	})
}

func mustEncode(t *testing.T, layout *LongRowLayout, member LongRowMember, status snapshot.MemberStatus, reason string, sample *snapshot.Sample) EncodedRow {
	t.Helper()
	row, err := layout.EncodeRow(longOutcome(member, status, reason, sample))
	require.NoError(t, err)
	return row
}

func rowCode(t *testing.T, err error) string {
	t.Helper()
	var rowErr *GroupRowError
	require.ErrorAs(t, err, &rowErr)
	return rowErr.Code
}
