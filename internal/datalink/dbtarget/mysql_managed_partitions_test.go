package dbtarget

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

func nativePartitionRows() []mysqlPartitionRow {
	return []mysqlPartitionRow{
		{
			name: sql.NullString{String: "p2026_11", Valid: true}, method: sql.NullString{String: "RANGE COLUMNS", Valid: true},
			expression:  sql.NullString{String: "`bucket_start`", Valid: true},
			description: sql.NullString{String: "'2026-11-01 00:00:00'", Valid: true},
			ordinal:     sql.NullInt64{Int64: 1, Valid: true},
			engine:      sql.NullString{String: "InnoDB", Valid: true},
		},
		{
			name: sql.NullString{String: "pmax", Valid: true}, method: sql.NullString{String: "RANGE COLUMNS", Valid: true},
			expression:  sql.NullString{String: "`bucket_start`", Valid: true},
			description: sql.NullString{String: "MAXVALUE", Valid: true},
			ordinal:     sql.NullInt64{Int64: 2, Valid: true},
			engine:      sql.NullString{String: "InnoDB", Valid: true},
		},
	}
}

func TestMySQLManagedPartitionEvidence(t *testing.T) {
	plain, err := mysqlManagedPartitionEvidence([]mysqlPartitionRow{{}}, "InnoDB")
	require.NoError(t, err)
	require.Nil(t, plain, "an unpartitioned table must stay plain nil")

	native, err := mysqlManagedPartitionEvidence(nativePartitionRows(), "InnoDB")
	require.NoError(t, err)
	require.NotNil(t, native)
	require.Equal(t, "RANGE COLUMNS", native.Method)
	require.Equal(t, "bucket_start", native.Expression)
	require.Equal(t, []MySQLPartitionBoundary{
		{Name: "p2026_11", LessThan: "2026-11-01T00:00:00Z", Ordinal: 1},
		{Name: "pmax", LessThan: "MAXVALUE", Ordinal: 2},
	}, native.Partitions)

	for name, testCase := range map[string]struct {
		rows   []mysqlPartitionRow
		engine string
	}{
		"empty catalog":          {nil, "InnoDB"},
		"multiple plain rows":    {[]mysqlPartitionRow{{}, {}}, "InnoDB"},
		"non innodb":             {nativePartitionRows(), "MyISAM"},
		"non innodb partition":   {withPartitionMutation(nativePartitionRows(), func(row *mysqlPartitionRow) { row.engine.String = "MyISAM" }), "InnoDB"},
		"foreign expression":     {withPartitionMutation(nativePartitionRows(), func(row *mysqlPartitionRow) { row.expression.String = "`other_column`" }), "InnoDB"},
		"hash method":            {withPartitionMutation(nativePartitionRows(), func(row *mysqlPartitionRow) { row.method.String = "HASH" }), "InnoDB"},
		"subpartition":           {withPartitionMutation(nativePartitionRows(), func(row *mysqlPartitionRow) { row.subpartitionName = sql.NullString{String: "sp1", Valid: true} }), "InnoDB"},
		"duplicate name":         {withPartitionMutation(nativePartitionRows(), func(row *mysqlPartitionRow) { row.name = sql.NullString{String: "p2026_11", Valid: true} }), "InnoDB"},
		"duplicate ordinal":      {withPartitionMutation(nativePartitionRows(), func(row *mysqlPartitionRow) { row.ordinal.Int64 = 1 }), "InnoDB"},
		"missing maxvalue":       {nativePartitionRows()[:1], "InnoDB"},
		"missing partition name": {withPartitionMutation(nativePartitionRows(), func(row *mysqlPartitionRow) { row.name = sql.NullString{} }), "InnoDB"},
		"foreign boundary":       {withPartitionMutation(nativePartitionRows(), func(row *mysqlPartitionRow) { row.description.String = "'not-a-date'" }), "InnoDB"},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := mysqlManagedPartitionEvidence(testCase.rows, testCase.engine)
			require.Error(t, err)
		})
	}
}

func withPartitionMutation(rows []mysqlPartitionRow, mutate func(*mysqlPartitionRow)) []mysqlPartitionRow {
	mutated := append([]mysqlPartitionRow(nil), rows...)
	mutate(&mutated[len(mutated)-1])
	return mutated
}

type failingMySQLPartitionQuerier struct{ err error }

func (f failingMySQLPartitionQuerier) QueryContext(context.Context, string, ...any) (*sql.Rows, error) {
	return nil, f.err
}

func (failingMySQLPartitionQuerier) QueryRowContext(context.Context, string, ...any) *sql.Row {
	return nil
}

func TestMySQLManagedPartitionQueryFailureIsNotPlain(t *testing.T) {
	_, err := readMySQLManagedPartitions(t.Context(), failingMySQLPartitionQuerier{err: &mysqldriver.MySQLError{Number: 1142}}, "gw", "t", "InnoDB")
	require.Error(t, err)
}

func TestMySQLManagedPartitionEvidenceEquality(t *testing.T) {
	native, err := mysqlManagedPartitionEvidence(nativePartitionRows(), "InnoDB")
	require.NoError(t, err)
	require.NotNil(t, native)
	require.True(t, sameMySQLPartitionEvidence(nil, nil))
	require.False(t, sameMySQLPartitionEvidence(nil, native))
	require.False(t, sameMySQLPartitionEvidence(native, nil))
	require.True(t, sameMySQLPartitionEvidence(native, cloneMySQLPartitionEvidence(native)))

	mutated := cloneMySQLPartitionEvidence(native)
	mutated.Partitions[1].LessThan = "MAXVALUE "
	require.False(t, sameMySQLPartitionEvidence(native, mutated))
}

func TestMySQLManagedPlainPartitionJSONIsUnchanged(t *testing.T) {
	plain, err := json.Marshal(&MySQLManagedTableEvidence{ServerFamily: "mysql"})
	require.NoError(t, err)
	require.False(t, strings.Contains(string(plain), "partitions"))

	native, err := mysqlManagedPartitionEvidence(nativePartitionRows(), "InnoDB")
	require.NoError(t, err)
	require.NotNil(t, native)
	encoded, err := json.Marshal(&MySQLManagedTableEvidence{ServerFamily: "mysql", Partitions: native})
	require.NoError(t, err)
	require.Contains(t, string(encoded), `"partitions":`)
}

func TestMySQLManagedPartitionJSONContract(t *testing.T) {
	encoded, err := json.Marshal(MySQLPartitionEvidence{
		Method: "RANGE COLUMNS", Expression: "bucket_start",
		Partitions: []MySQLPartitionBoundary{{Name: "p1", LessThan: "2026-11-01T00:00:00Z", Ordinal: 1}},
	})
	require.NoError(t, err)
	for _, key := range []string{`"method"`, `"expression"`, `"partitions"`, `"name"`, `"less_than"`, `"ordinal"`} {
		require.Contains(t, string(encoded), key)
	}
}

func TestMySQLManagedPartitionHookPropagation(t *testing.T) {
	native, err := mysqlManagedPartitionEvidence(nativePartitionRows(), "InnoDB")
	require.NoError(t, err)
	require.NotNil(t, native)
	facts := mysqlManagedTableFacts{status: TableInspectionExists, engine: "InnoDB"}
	merged, err := mergeMySQLManagedPartitions(facts, native, nil)
	require.NoError(t, err)
	require.NotNil(t, merged.partitions)

	forbidden, err := mergeMySQLManagedPartitions(facts, nil, &mysqldriver.MySQLError{Number: 1142})
	require.NoError(t, err)
	require.Equal(t, TableInspectionForbidden, forbidden.status)
	require.Nil(t, forbidden.partitions)

	_, err = mergeMySQLManagedPartitions(facts, nil, errors.New("catalog query failed"))
	require.Error(t, err)
}

func TestMySQLManagedPlainProjectionKeepsPartitionsNil(t *testing.T) {
	facts := mysqlManagedTableFacts{
		status: TableInspectionExists, tableType: "BASE TABLE", engine: "InnoDB",
		charset: "utf8mb4", collation: "utf8mb4_0900_bin", jsonColumns: map[string]bool{},
	}
	result, err := buildMySQLManagedInspection("gw", "plain_table", mysqlManagedServerFacts{
		family: "mysql", version: "8.4.5", capabilityRevision: "cap", database: "gw",
		metadataStatus: mysqlRecordingMetadataVerified, transactionalEngine: true,
	}, facts)
	require.NoError(t, err)
	require.NotNil(t, result.Evidence)
	require.Nil(t, result.Evidence.Partitions)
	encoded, err := json.Marshal(result.Evidence)
	require.NoError(t, err)
	require.False(t, strings.Contains(string(encoded), `"partitions"`))
}
