package dbtarget

import (
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/snapshot"

	"github.com/stretchr/testify/require"
)

func TestMySQLLongRowEncoder(t *testing.T) {
	declared, err := recordingplan.ManagedMySQLLongGroupLayout(indexTestTable, indexTestOwner, recordingplan.MySQLGroupSchemaOptions{
		ServerFamily: "mysql", ServerVersion: "8.4.13", CapabilityRevision: "encoder-declaration",
		Engine: "InnoDB", Collation: "utf8mb4_0900_bin",
	})
	require.NoError(t, err)
	columns := make([]ColumnInfo, len(declared.Columns))
	for i, column := range declared.Columns {
		columns[i] = ColumnInfo{Name: column.Name, DataType: column.SQLType, Nullable: column.Nullable, PrimaryKey: column.PrimaryKey}
	}
	cases := longCases(t)
	encoder, issues := NewLongRowLayout(longSpec(SQLDialectMySQL, columns, cases))
	require.Empty(t, issues)
	for _, testCase := range cases {
		t.Run(string(testCase.member.Type), func(t *testing.T) {
			sample := longSample("mysql-sample", testCase.value)
			sample.ObservedAt = sample.ObservedAt.In(time.FixedZone("UTC+8", 8*60*60))
			outcome := longOutcome(testCase.member, snapshot.MemberOK, "", sample)
			outcome.RecordID = strings.Repeat("a", 64)
			row, err := encoder.EncodeRow(outcome)
			require.NoError(t, err)
			cells := cellMap(row)
			wantObserved := time.Date(2026, 10, 24, 0, 0, 29, 123456000, time.UTC)
			require.Equal(t, wantObserved, cells[longColObservedAt], "bind a UTC microsecond time.Time, not an RFC3339 string")
			require.Equal(t, time.Date(2026, 10, 24, 0, 0, 15, 0, time.UTC), cells[longColBucketStart])
			require.Equal(t, time.Date(2026, 10, 24, 0, 0, 30, 0, time.UTC), cells[longColBucketEnd])
			populated := 0
			for _, column := range []string{longColValueInt, longColValueDecimal, longColValueFloat, longColValueBool, longColValueText} {
				if cells[column] != nil {
					populated++
				}
			}
			require.Equal(t, 1, populated, "one raw typed slot per record")
			require.Equal(t, "gateway", cells[indexTestOwner])
		})
	}
	missing := longOutcome(cases[0].member, snapshot.MemberMissing, "no-data", nil)
	missing.RecordID = strings.Repeat("b", 64)
	row, err := encoder.EncodeRow(missing)
	require.NoError(t, err)
	cells := cellMap(row)
	require.Equal(t, "missing", cells[longColStatus])
	for _, column := range []string{longColValueInt, longColValueDecimal, longColValueFloat, longColValueBool, longColValueText, longColObservedAt, longColReceivedAt} {
		require.Nil(t, cells[column], column)
	}
}
