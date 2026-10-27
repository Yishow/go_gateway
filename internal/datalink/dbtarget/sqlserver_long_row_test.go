package dbtarget

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"

	"github.com/stretchr/testify/require"
)

// These are declared SQL Server columns, not evidence of an inspected server.
func sqlServerLongColumnsForTest() []ColumnInfo {
	columns := make([]ColumnInfo, 0, len(longColumnRoles)+1)
	for _, column := range longColumnRoles {
		declared := "NVARCHAR(191)"
		switch column.role {
		case longTime:
			declared = "DATETIME2(6)"
		case longSlotInt:
			declared = "BIGINT"
		case longSlotDecimal, longSlotDelta:
			declared = "DECIMAL(38,18)"
		case longSlotFloat:
			declared = "FLOAT(53)"
		case longSlotBool:
			declared = "BIT"
		case longSlotText, longProvenanceRole:
			declared = "NVARCHAR(MAX)"
		}
		if column.name == longColRecordID {
			declared = "CHAR(64)"
		}
		columns = append(columns, ColumnInfo{Name: column.name, DataType: declared, Nullable: !column.required,
			PrimaryKey: column.name == longColBucketStart || column.name == longColRecordID})
	}
	return append(columns, ColumnInfo{Name: indexTestOwner, DataType: "NVARCHAR(16)"})
}

func sqlServerLongSpecForTest(t *testing.T) LongRowSpec {
	t.Helper()
	return longSpec(SQLDialectSQLServer, sqlServerLongColumnsForTest(), longCases(t))
}

func sqlServerOutcomeForTest(member LongRowMember, value measurement.ExactValue) snapshot.Outcome {
	outcome := longOutcome(member, snapshot.MemberOK, "", longSample("sqlserver-sample", value))
	outcome.RecordID = strings.Repeat("a", 64)
	outcome.EffectKey = strings.Repeat("b", 64)
	return outcome
}

func TestSQLServerLongExactValues(t *testing.T) {
	spec := sqlServerLongSpecForTest(t)
	layout, issues := NewLongRowLayout(spec)
	require.Empty(t, issues)
	for _, tc := range longCases(t) {
		t.Run(tc.member.MemberKey, func(t *testing.T) {
			outcome := sqlServerOutcomeForTest(tc.member, tc.value)
			outcome.Members[0].Sample.ObservedAt = time.Date(2026, 10, 24, 8, 0, 29, 123456789, time.FixedZone("UTC+8", 8*3600))
			row, err := layout.EncodeRow(outcome)
			require.NoError(t, err)
			cells := cellMap(row)
			require.Equal(t, time.Date(2026, 10, 24, 0, 0, 29, 123456000, time.UTC), cells[longColObservedAt])
			require.Nil(t, cells[longColDeltaValue], "a first/instantaneous raw record has no invented delta")
			require.True(t, json.Valid([]byte(cells[longColProvenance].(string))))
			switch tc.value.Type() {
			case measurement.ExactUint64:
				require.Equal(t, "18446744073709551615", cells[longColValueDecimal])
			case measurement.ExactInt64:
				require.Equal(t, int64(9007199254740993), cells[longColValueInt])
			case measurement.ExactDecimal:
				require.Equal(t, "12000.4", cells[longColValueDecimal])
			case measurement.ExactBool:
				require.Equal(t, true, cells[longColValueBool])
			}
		})
	}
	t.Run("Unicode and case exact source bytes", func(t *testing.T) {
		upper, lower := spec.Members[5], spec.Members[5]
		upper.TagID, lower.TagID = "電路A", "電路a"
		for _, member := range []LongRowMember{upper, lower} {
			local := spec
			local.Members = []LongRowMember{member}
			encoder, issues := NewLongRowLayout(local)
			require.Empty(t, issues)
			row, err := encoder.EncodeRow(sqlServerOutcomeForTest(member, measurement.NewText("運轉中🟢")))
			require.NoError(t, err)
			require.Equal(t, member.TagID, cellMap(row)[longColTagID])
			require.Equal(t, "運轉中🟢", cellMap(row)[longColValueText])
		}
	})
	t.Run("unknown values and times stay NULL", func(t *testing.T) {
		outcome := sqlServerOutcomeForTest(spec.Members[0], measurement.NewInt64(1))
		outcome.Members[0] = snapshot.MemberResult{MemberKey: spec.Members[0].MemberKey, Status: snapshot.MemberMissing}
		row, err := layout.EncodeRow(outcome)
		require.NoError(t, err)
		for _, name := range []string{longColValueInt, longColValueDecimal, longColValueBool, longColDeltaValue, longColObservedAt, longColReceivedAt} {
			require.Nil(t, cellMap(row)[name], name)
		}
	})
}

func TestSQLServerDecimalRejected(t *testing.T) {
	layout := &LongRowLayout{spec: LongRowSpec{Dialect: SQLDialectSQLServer}, columns: map[string]ColumnInfo{
		longColValueDecimal: {Name: longColValueDecimal, DataType: "DECIMAL(38,18)"},
	}}
	require.NoError(t, layout.ValidateSampleValue(measurement.NewUint64(math.MaxUint64)))
	for _, digits := range []string{"100000000000000000000", "0.0000000000000000001"} {
		value, err := measurement.NewDecimal(digits)
		require.NoError(t, err)
		err = layout.ValidateSampleValue(value)
		require.ErrorIs(t, err, ErrExactValueUnrepresentable)
		require.NotContains(t, err.Error(), digits)
	}
	// Extra fractional zeros do not lose value, matching the shared exact codec.
	value, err := measurement.NewDecimal("1.0000000000000000000")
	require.NoError(t, err)
	require.NoError(t, layout.ValidateSampleValue(value))
}

func TestSQLServerIdentifierBinding(t *testing.T) {
	kind := schema.DatabaseConnectorKindSQLServer
	require.Equal(t, []string{"@p1", "@p2", "@p3"}, buildPlaceholders(kind, 3))
	require.Equal(t, "[電路]]名稱]", quoteIdentifier(kind, "電路]名稱"))
	require.Equal(t, "[dbo].[table]]one]", qualifiedTableName(kind, "dbo", "table]one"))
}

func TestSQLServerLongCounterUsesSharedExactResults(t *testing.T) {
	for _, tc := range []struct {
		name          string
		typ           measurement.ExactType
		first, second *measurement.ExactValue
		delta         string
	}{
		{"decimal", measurement.ExactDecimal, exactDecimal(t, "12000.0"), exactDecimal(t, "12000.4"), "0.4"},
		{"large integer", measurement.ExactUint64, exactUint(9007199254740992), exactUint(9007199254740993), "1"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			member := counterMemberFor(tc.typ)
			layout, issues := NewLongRowLayout(counterSpecFor(SQLDialectSQLServer, sqlServerLongColumnsForTest(), member))
			require.Empty(t, issues)
			outcomes := counterOutcomes(t, member, tc.first, tc.second)
			for i, outcome := range outcomes {
				outcome.RecordID = fmt.Sprintf("%064x", i+1)
				row, err := layout.EncodeRow(outcome)
				require.NoError(t, err)
				cells := cellMap(row)
				if i == 0 {
					require.Nil(t, cells[longColDeltaValue])
					require.Equal(t, "baseline", cells[longColStatus])
				} else {
					require.Equal(t, tc.delta, cells[longColDeltaValue])
					require.IsType(t, time.Time{}, cells[longColDeltaFrom])
				}
			}
		})
	}
}

func TestSQLServerLongIdentityAndTimeDeclarations(t *testing.T) {
	for _, tag := range []string{strings.Repeat("🟢", 96), "電路 ", string([]byte{0xff})} {
		spec := sqlServerLongSpecForTest(t)
		spec.Members[0].TagID = tag
		_, issues := NewLongRowLayout(spec)
		require.Contains(t, issues, LayoutIssue{Code: issueIdentityUnsupported, Column: longColTagID, MemberKey: spec.Members[0].MemberKey})
	}
	spec := sqlServerLongSpecForTest(t)
	spec.Members[0].TagID = strings.Repeat("🟢", 95) + "a"
	_, issues := NewLongRowLayout(spec)
	require.Empty(t, issues, "NVARCHAR(191) measures UTF-16 code units")
	for _, declared := range []string{"TIMESTAMP", "ROWVERSION", "DATETIME", "DATETIME2", "DATETIME2(7)", "DATETIME2(3)"} {
		require.False(t, longColumnUsable(SQLDialectSQLServer, ColumnInfo{DataType: declared}, longTime), declared)
	}
	require.True(t, longColumnUsable(SQLDialectSQLServer, ColumnInfo{DataType: "DATETIME2(6)"}, longTime))
}

func TestSQLServerLongRefusesUnnormalizedWindowIdentity(t *testing.T) {
	spec := sqlServerLongSpecForTest(t)
	layout, issues := NewLongRowLayout(spec)
	require.Empty(t, issues)
	outcome := sqlServerOutcomeForTest(spec.Members[0], measurement.NewInt64(1))
	outcome.BucketStart = outcome.BucketStart.Add(time.Nanosecond)
	_, err := layout.EncodeRow(outcome)
	var blocked *GroupRowError
	require.ErrorAs(t, err, &blocked, "a caller must normalize the window before deriving its identity")
	require.Equal(t, longColBucketStart, blocked.Column)
}
