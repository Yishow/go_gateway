package dbtarget

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func mysqlGroupTestInspection(table string, primaryKey []string) *MySQLManagedInspection {
	columns := []ColumnInfo{
		{Name: "bucket_start", DataType: "datetime(6)", PrimaryKey: true},
		{Name: "record_id", DataType: "char(64)", PrimaryKey: true},
	}
	defaults := map[string]*string{"bucket_start": nil, "record_id": nil}
	if table == EffectReceiptTable {
		columns = []ColumnInfo{{Name: "effect_key", DataType: "varchar(191)", PrimaryKey: true}}
		defaults = map[string]*string{"effect_key": nil}
	}
	return &MySQLManagedInspection{
		Inspection: &TableInspection{
			Status: TableInspectionExists, Schema: "gw", Table: table, Columns: columns,
			Indexes: []IndexInfo{}, IndexesInspected: true,
		},
		Evidence: &MySQLManagedTableEvidence{
			ServerFamily: "mysql", ServerVersion: "8.4.13", CapabilityRevision: "cap-1", Database: "gw",
			Engine: "InnoDB", Charset: "utf8mb4", Collation: "utf8mb4_0900_bin", TableType: "BASE TABLE",
			PrimaryKey: primaryKey, ColumnCollations: map[string]string{}, ColumnDefaults: defaults, JSONColumns: map[string]bool{},
		},
		Indexes: []IndexInfo{},
	}
}

func mysqlGroupTestRequest() GroupInsertRequest {
	return GroupInsertRequest{
		Kind: schema.DatabaseConnectorKindMySQL, SchemaName: "gw", TableName: "gw_long_table",
		Strategy: GroupEffectReceipt, PayloadDigest: "digest-1",
		CommittedAt: "2026-01-02T03:04:05.123456Z",
		Row:         EncodedRow{EffectKey: "effect-1", Cells: []EncodedCell{{Column: "record_id", Value: "row-1"}}},
		MySQL: &MySQLManagedRowProof{
			ConnectorID: "conn-mysql", ConnectorRevision: "rev-1",
			Group:   mysqlGroupTestInspection("gw_long_table", []string{"bucket_start", "record_id"}),
			Receipt: mysqlGroupTestInspection(EffectReceiptTable, []string{"effect_key"}),
		},
	}
}

func TestMySQLGroupInsertRequiresProof(t *testing.T) {
	req := mysqlGroupTestRequest()
	req.MySQL = nil
	err := validateGroupInsert(&sql.DB{}, req)
	require.ErrorIs(t, err, ErrGroupInsertUnsupported)

	req = mysqlGroupTestRequest()
	req.MySQL.ConnectorID = ""
	require.ErrorIs(t, validateGroupInsert(&sql.DB{}, req), ErrGroupInsertUnsupported)

	req = mysqlGroupTestRequest()
	req.SchemaName = "other"
	require.ErrorIs(t, validateGroupInsert(&sql.DB{}, req), ErrGroupInsertUnsupported)

	req = mysqlGroupTestRequest()
	req.MySQL.Group.Evidence.Engine = "MyISAM"
	require.ErrorIs(t, validateGroupInsert(&sql.DB{}, req), ErrGroupInsertUnsupported)
}

func TestMySQLGroupInsertAcceptsCompleteProof(t *testing.T) {
	require.NoError(t, validateGroupInsert(&sql.DB{}, mysqlGroupTestRequest()))
}

func TestMySQLDeliveryDatabaseRequiresVerifiedTLS(t *testing.T) {
	database, err := mysqlDeliveryDatabase(ConnectionConfig{"dsn": "user:pass@tcp(127.0.0.1:3306)/gw?tls=true"})
	require.NoError(t, err)
	require.Equal(t, "gw", database)

	_, err = mysqlDeliveryDatabase(ConnectionConfig{"dsn": "user:pass@tcp(127.0.0.1:3306)/gw"})
	require.Error(t, err)
	require.NotContains(t, err.Error(), "pass")
}

func TestMySQLCommittedAtStrictUTC(t *testing.T) {
	parsed, err := parseMySQLCommittedAt("2026-01-02T03:04:05.123456Z")
	require.NoError(t, err)
	require.Equal(t, time.Date(2026, 1, 2, 3, 4, 5, 123456000, time.UTC), parsed)

	parsed, err = parseMySQLCommittedAt("2026-01-02T11:04:05.123456+08:00")
	require.NoError(t, err)
	require.Equal(t, time.Date(2026, 1, 2, 3, 4, 5, 123456000, time.UTC), parsed)

	for _, invalid := range []string{
		"2026-01-02T03:04:05Z",
		"2026-01-02T03:04:05.1234Z",
		"2026-01-02 03:04:05.123456",
		"not-a-time",
	} {
		_, err := parseMySQLCommittedAt(invalid)
		require.Error(t, err, invalid)
	}
}

func TestMySQLGroupInsertStatementsArePlain(t *testing.T) {
	req := mysqlGroupTestRequest()
	receipts := "`gw_effect_receipts`"
	group := "`gw_long_table`"
	statements := []string{
		receiptInsertStatement(schema.DatabaseConnectorKindMySQL, receipts),
		groupInsertStatement(req, group, []string{"`record_id`"}),
	}
	for _, statement := range statements {
		upper := strings.ToUpper(statement)
		require.NotContains(t, upper, "ON DUPLICATE")
		require.NotContains(t, upper, "ON CONFLICT")
		require.NotContains(t, upper, "INSERT IGNORE")
		require.NotContains(t, upper, "REPLACE")
	}
	require.Contains(t, mysqlReceiptLookupForUpdateStatement(receipts), "FOR UPDATE")
}

func TestMySQLManagedRowProofMatchesPhysicalProof(t *testing.T) {
	proof := mysqlGroupTestRequest().MySQL
	facts := mysqlManagedTableFacts{
		status: TableInspectionExists, tableType: "BASE TABLE", engine: "InnoDB",
		charset: "utf8mb4", collation: "utf8mb4_0900_bin",
		columns: []mysqlManagedColumnFacts{
			{name: "bucket_start", dataType: "datetime", columnType: "datetime(6)", primaryKey: true},
			{name: "record_id", dataType: "char", columnType: "char(64)", primaryKey: true},
		},
		primaryKey:  []string{"bucket_start", "record_id"},
		jsonColumns: map[string]bool{},
	}
	mismatch := facts
	mismatch.engine = "MyISAM"
	require.ErrorIs(t, matchMySQLManagedInspection(proof.Group, mismatch), ErrGroupInsertUnsupported)
	require.ErrorIs(t, matchMySQLManagedInspection(proof.Group, mismatch), ErrMySQLManagedPhysicalMismatch)

	missing := facts
	missing.status = TableInspectionMissing
	err := matchMySQLManagedInspection(proof.Group, missing)
	require.ErrorIs(t, err, ErrGroupInsertUnsupported)
	require.False(t, errors.Is(err, ErrMySQLManagedPhysicalMismatch), "missing is not physical drift")

	require.NoError(t, matchMySQLManagedInspection(proof.Group, facts))
}

func TestMySQLManagedRowProofMatchesFrozenJSON(t *testing.T) {
	proof := mysqlGroupTestRequest().MySQL
	encoded, err := json.Marshal(proof)
	require.NoError(t, err)
	var frozen MySQLManagedRowProof
	require.NoError(t, json.Unmarshal(encoded, &frozen))
	require.NoError(t, validateMySQLManagedRowProof(GroupInsertRequest{
		Kind: schema.DatabaseConnectorKindMySQL, SchemaName: "gw", TableName: "gw_long_table",
		Strategy: GroupEffectReceipt, PayloadDigest: "digest-1", CommittedAt: "2026-01-02T03:04:05.123456Z",
		Row:   EncodedRow{EffectKey: "effect-1", Cells: []EncodedCell{{Column: "record_id", Value: "row-1"}}},
		MySQL: &frozen,
	}))

	facts := mysqlManagedTableFacts{
		status: TableInspectionExists, tableType: "BASE TABLE", engine: "InnoDB",
		charset: "utf8mb4", collation: "utf8mb4_0900_bin",
		columns: []mysqlManagedColumnFacts{
			{name: "bucket_start", dataType: "datetime", columnType: "datetime(6)", primaryKey: true},
			{name: "record_id", dataType: "char", columnType: "char(64)", primaryKey: true},
		},
		primaryKey:  []string{"bucket_start", "record_id"},
		jsonColumns: map[string]bool{},
	}
	require.NoError(t, matchMySQLManagedInspection(frozen.Group, facts))
}
