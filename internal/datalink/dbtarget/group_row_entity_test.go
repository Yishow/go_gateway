package dbtarget

import (
	"testing"

	"go-gateway/internal/datalink/measurement"

	"github.com/stretchr/testify/require"
)

func TestEntityEncodingUsesVerifiedColumnIdentifiers(t *testing.T) {
	for _, dialect := range []SQLDialect{SQLDialectSQLite, SQLDialectPostgres} {
		t.Run(string(dialect), func(t *testing.T) {
			layout, issues := NewGroupRowLayout(GroupRowSpec{
				Dialect:     dialect,
				Columns:     []ColumnInfo{{Name: "temperature", DataType: "DOUBLE PRECISION"}, {Name: "entity", DataType: "TEXT"}},
				Members:     []GroupRowMember{{MemberKey: "temp", EntityKey: "A", Column: "Temperature", Type: measurement.ExactFloat64, Required: true}},
				EntityKeyed: true, EntityKeyColumn: "ENTITY",
			})
			require.Empty(t, issues, "readiness accepts aliases of the inspected identifiers")
			value, err := measurement.NewFloat64(21.5)
			require.NoError(t, err)
			outcome := rowOutcome(okMember("temp", value))
			outcome.EntityKey = "A"
			row, err := layout.EncodeRow(outcome)
			require.NoError(t, err)
			require.Equal(t, []EncodedCell{{Column: "temperature", Value: 21.5}, {Column: "entity", Value: "A"}}, row.Cells,
				"SQL quoting must use verified actual names, including identity columns")
		})
	}
}

func TestEntityEncodingPreservesExactQuotedPostgresColumn(t *testing.T) {
	layout, issues := NewGroupRowLayout(GroupRowSpec{
		Dialect: SQLDialectPostgres,
		Columns: []ColumnInfo{{Name: "Temperature", DataType: "DOUBLE PRECISION"}, {Name: "temperature", DataType: "DOUBLE PRECISION"}, {Name: "entity", DataType: "TEXT"}},
		Members: []GroupRowMember{
			{MemberKey: "a", EntityKey: "A", Column: "Temperature", Type: measurement.ExactFloat64, Required: true},
			{MemberKey: "b", EntityKey: "B", Column: "temperature", Type: measurement.ExactFloat64, Required: true},
		},
		EntityKeyed: true, EntityKeyColumn: "entity",
	})
	require.Empty(t, issues)
	for _, row := range []struct{ entity, member, column string }{{"A", "a", "Temperature"}, {"B", "b", "temperature"}} {
		value, err := measurement.NewFloat64(21.5)
		require.NoError(t, err)
		outcome := rowOutcome(okMember(row.member, value))
		outcome.EntityKey = row.entity
		encoded, err := layout.EncodeRow(outcome)
		require.NoError(t, err)
		require.Equal(t, row.column, encoded.Cells[0].Column, "an exact inspected identifier must not be redirected to a case variant")
	}
}
