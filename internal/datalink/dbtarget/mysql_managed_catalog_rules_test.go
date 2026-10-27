package dbtarget

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMySQLManagedJSONCheckMatcher(t *testing.T) {
	valid := []string{
		"json_valid(`provenance`)",
		"JSON_VALID(provenance)",
		"  ( json_valid( `provenance` ) )  ",
		"((JSON_VALID(provenance)))",
	}
	for _, clause := range valid {
		t.Run("valid "+clause, func(t *testing.T) {
			require.True(t, mysqlManagedCheckIsJSONValid(clause, "provenance"), clause)
		})
	}
	invalid := []string{
		"json_valid(other_provenance)",
		"json_valid(`other_provenance`)",
		"json_valid('provenance')",
		"json_valid(provenance) OR TRUE",
		"NOT json_valid(provenance)",
		"json_valid(provenance) = 1",
		"json_valid(provenance, 1)",
		"json_valid(provenance_x)",
		"json_valid(`provenance`) OR TRUE",
	}
	for _, clause := range invalid {
		t.Run("invalid "+clause, func(t *testing.T) {
			require.False(t, mysqlManagedCheckIsJSONValid(clause, "provenance"), clause)
		})
	}
}

func TestMySQLManagedPrimaryKeyRules(t *testing.T) {
	valid := []mysqlManagedPrimaryKeyRow{
		{column: sql.NullString{String: "bucket_start", Valid: true}, indexType: mysqlManagedIndexTypeBTREE},
		{column: sql.NullString{String: "record_id", Valid: true}, indexType: mysqlManagedIndexTypeBTREE},
	}
	columns, err := mysqlManagedPrimaryKeyColumns(valid)
	require.NoError(t, err)
	require.Equal(t, []string{"bucket_start", "record_id"}, columns)

	cases := map[string][]mysqlManagedPrimaryKeyRow{
		"prefix": {
			{column: sql.NullString{String: "bucket_start", Valid: true}, indexType: mysqlManagedIndexTypeBTREE},
			{column: sql.NullString{String: "record_id", Valid: true}, subPart: sql.NullInt64{Int64: 10, Valid: true}, indexType: mysqlManagedIndexTypeBTREE},
		},
		"functional": {
			{column: sql.NullString{String: "bucket_start", Valid: true}, indexType: mysqlManagedIndexTypeBTREE},
			{column: sql.NullString{}, indexType: mysqlManagedIndexTypeBTREE},
		},
		"non_btree": {
			{column: sql.NullString{String: "bucket_start", Valid: true}, indexType: "HASH"},
		},
		"non_unique": {
			{column: sql.NullString{String: "bucket_start", Valid: true}, indexType: mysqlManagedIndexTypeBTREE, nonUnique: 1},
		},
	}
	for name, rows := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := mysqlManagedPrimaryKeyColumns(rows)
			require.ErrorIs(t, err, ErrMySQLManagedTableUnconfirmed)
		})
	}
}

func TestMySQLManagedIndexVisibilityRules(t *testing.T) {
	yes := sql.NullString{String: "YES", Valid: true}
	no := sql.NullString{String: "NO", Valid: true}
	unknown := sql.NullString{}

	indexes := func() map[string]*mysqlManagedIndexFacts {
		return map[string]*mysqlManagedIndexFacts{"idx": {name: "idx", visible: true}}
	}

	require.NoError(t, applyMySQLManagedIndexVisibility(mysqlManagedFamilyMySQL, map[string]*mysqlManagedIndexFacts{},
		nil, nil), "no expected index needs no visibility columns")

	require.ErrorIs(t, applyMySQLManagedIndexVisibility(mysqlManagedFamilyMySQL, indexes(),
		[]string{"Key_name"}, []mysqlManagedVisibilityRow{{name: "idx", visible: yes}}),
		ErrMySQLManagedTableUnconfirmed, "MySQL requires Visible")

	require.ErrorIs(t, applyMySQLManagedIndexVisibility(mysqlManagedFamilyMySQL, indexes(),
		[]string{"Key_name", "Visible"}, []mysqlManagedVisibilityRow{{name: "idx", visible: unknown}}),
		ErrMySQLManagedTableUnconfirmed, "MySQL NULL visibility is unknown")

	require.ErrorIs(t, applyMySQLManagedIndexVisibility(mysqlManagedFamilyMySQL, indexes(),
		[]string{"Key_name", "Visible"}, []mysqlManagedVisibilityRow{{name: "idx", visible: no}}),
		ErrMySQLManagedTableUnconfirmed, "MySQL invisible index is unusable")

	require.NoError(t, applyMySQLManagedIndexVisibility(mysqlManagedFamilyMySQL, indexes(),
		[]string{"Key_name", "Visible"}, []mysqlManagedVisibilityRow{{name: "idx", visible: yes}}))

	require.ErrorIs(t, applyMySQLManagedIndexVisibility(mysqlManagedFamilyMariaDB, indexes(),
		[]string{"Key_name"}, []mysqlManagedVisibilityRow{{name: "idx", ignored: no}}),
		ErrMySQLManagedTableUnconfirmed, "MariaDB requires Ignored")

	require.ErrorIs(t, applyMySQLManagedIndexVisibility(mysqlManagedFamilyMariaDB, indexes(),
		[]string{"Key_name", "Ignored"}, []mysqlManagedVisibilityRow{{name: "idx", ignored: unknown}}),
		ErrMySQLManagedTableUnconfirmed, "MariaDB NULL ignored is unknown")

	require.ErrorIs(t, applyMySQLManagedIndexVisibility(mysqlManagedFamilyMariaDB, indexes(),
		[]string{"Key_name", "Ignored"}, []mysqlManagedVisibilityRow{{name: "idx", ignored: yes}}),
		ErrMySQLManagedTableUnconfirmed, "MariaDB ignored index is unusable")

	require.NoError(t, applyMySQLManagedIndexVisibility(mysqlManagedFamilyMariaDB, indexes(),
		[]string{"Key_name", "Ignored"}, []mysqlManagedVisibilityRow{{name: "idx", ignored: no}}))

	require.ErrorIs(t, applyMySQLManagedIndexVisibility(mysqlManagedFamilyMySQL, indexes(),
		[]string{"Key_name", "Visible"}, nil),
		ErrMySQLManagedTableUnconfirmed, "missing expected index row is unknown")
}
