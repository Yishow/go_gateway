package dbtarget

import (
	"context"
	"errors"
	"strings"
	"testing"

	"go-gateway/internal/datalink/schema"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

const mysqlManagedTestDSN = "user:pass@tcp(127.0.0.1:3306)/gw?tls=true"

type mysqlManagedTestRepo struct {
	ConnectorRepository
	connector *schema.DatabaseConnector
}

func (r *mysqlManagedTestRepo) GetByID(context.Context, string) (*schema.DatabaseConnector, error) {
	return r.connector, nil
}

type mysqlManagedTestSession struct {
	server    mysqlManagedServerFacts
	table     mysqlManagedTableFacts
	serverErr error
	tableErr  error
	execErr   error
	execs     []string
}

func (s *mysqlManagedTestSession) ServerFacts(context.Context, string) (mysqlManagedServerFacts, error) {
	return s.server, s.serverErr
}

func (s *mysqlManagedTestSession) TableFacts(context.Context, string, string) (mysqlManagedTableFacts, error) {
	return s.table, s.tableErr
}

func (s *mysqlManagedTestSession) Exec(_ context.Context, statement string) error {
	s.execs = append(s.execs, statement)
	return s.execErr
}

func (s *mysqlManagedTestSession) Close() error { return nil }

func mysqlManagedTestConnector() *schema.DatabaseConnector {
	return &schema.DatabaseConnector{
		ID: "conn-mysql", Kind: schema.DatabaseConnectorKindMySQL, Enabled: true,
		IdentityRevision: "rev-1",
		ConnectionConfig: `{"dsn":"` + mysqlManagedTestDSN + `"}`,
	}
}

func newMySQLManagedTestInspector(t *testing.T, session *mysqlManagedTestSession, connector *schema.DatabaseConnector) *ManagedTableInspector {
	t.Helper()
	if connector == nil {
		connector = mysqlManagedTestConnector()
	}
	svc := NewConnectorService(&mysqlManagedTestRepo{connector: connector})
	original := openMySQLManagedSessionFunc
	openMySQLManagedSessionFunc = func(*schema.DatabaseConnector, ConnectionConfig, string) (mysqlManagedSession, error) {
		return session, nil
	}
	t.Cleanup(func() { openMySQLManagedSessionFunc = original })
	return NewManagedTableInspector(svc, nil)
}

func mysqlManagedBaselineFacts() mysqlManagedServerFacts {
	return mysqlManagedServerFacts{
		family: "mysql", version: "8.4.13", capabilityRevision: "cap-1", database: "gw",
		metadataStatus: "verified", transactionalEngine: true,
	}
}

func mysqlManagedInspectionCall(t *testing.T, inspector *ManagedTableInspector) (*MySQLManagedInspection, error) {
	t.Helper()
	return inspector.InspectMySQLManagedTableAtRevision(t.Context(), "conn-mysql", "rev-1", "gw", "cap-1", "gw_long_table")
}

func TestMySQLManagedInspectionMissingKeepsBindingEvidence(t *testing.T) {
	session := &mysqlManagedTestSession{
		server: mysqlManagedBaselineFacts(),
		table:  mysqlManagedTableFacts{status: TableInspectionMissing},
	}
	inspection, err := mysqlManagedInspectionCall(t, newMySQLManagedTestInspector(t, session, nil))
	require.NoError(t, err)
	require.NotNil(t, inspection)
	require.NotNil(t, inspection.Inspection)
	require.Equal(t, TableInspectionMissing, inspection.Inspection.Status)
	require.NotNil(t, inspection.Evidence, "missing still needs verified server binding")
	require.Equal(t, "mysql", inspection.Evidence.ServerFamily)
	require.Equal(t, "8.4.13", inspection.Evidence.ServerVersion)
	require.Equal(t, "cap-1", inspection.Evidence.CapabilityRevision)
	require.Equal(t, "gw", inspection.Evidence.Database)
	require.Empty(t, inspection.Evidence.Engine)
}

func TestMySQLManagedInspectionForbiddenIsNotMissing(t *testing.T) {
	session := &mysqlManagedTestSession{
		server: mysqlManagedBaselineFacts(),
		table:  mysqlManagedTableFacts{status: TableInspectionForbidden, reason: "permission_denied"},
	}
	inspection, err := mysqlManagedInspectionCall(t, newMySQLManagedTestInspector(t, session, nil))
	require.ErrorIs(t, err, ErrMySQLManagedTableForbidden)
	require.NotNil(t, inspection)
	require.Equal(t, TableInspectionForbidden, inspection.Inspection.Status)
	require.NotNil(t, inspection.Evidence)
	require.Equal(t, "cap-1", inspection.Evidence.CapabilityRevision)
}

func TestMySQLManagedInspectionUnconfirmedIsNotMissing(t *testing.T) {
	session := &mysqlManagedTestSession{
		server: mysqlManagedBaselineFacts(),
		table:  mysqlManagedTableFacts{status: TableInspectionFailed, reason: "existence_unconfirmed"},
	}
	inspection, err := mysqlManagedInspectionCall(t, newMySQLManagedTestInspector(t, session, nil))
	require.ErrorIs(t, err, ErrMySQLManagedTableUnconfirmed)
	require.NotNil(t, inspection)
	require.Equal(t, TableInspectionFailed, inspection.Inspection.Status)
	require.NotNil(t, inspection.Evidence)
}

func TestMySQLManagedInspectionPreciseColumnsAndPrimaryKey(t *testing.T) {
	defaultValue := "42"
	session := &mysqlManagedTestSession{
		server: mysqlManagedBaselineFacts(),
		table: mysqlManagedTableFacts{
			status: TableInspectionExists, tableType: "BASE TABLE", engine: "InnoDB",
			charset: "utf8mb4", collation: "utf8mb4_0900_bin",
			columns: []mysqlManagedColumnFacts{
				{name: "bucket_start", dataType: "datetime", columnType: "datetime(6)", primaryKey: true},
				{name: "record_id", dataType: "char", columnType: "char(64)", primaryKey: true, collation: "utf8mb4_0900_bin"},
				{name: "value_decimal", dataType: "decimal", columnType: "decimal(65,30)", nullable: true, hasDefault: true, defaultValue: &defaultValue},
				{name: "provenance", dataType: "json", columnType: "json", jsonValid: true},
			},
			primaryKey:  []string{"bucket_start", "record_id"},
			jsonColumns: map[string]bool{"provenance": true},
		},
	}
	inspection, err := mysqlManagedInspectionCall(t, newMySQLManagedTestInspector(t, session, nil))
	require.NoError(t, err)
	require.Equal(t, "datetime(6)", inspection.Inspection.Columns[0].DataType)
	require.Equal(t, "char(64)", inspection.Inspection.Columns[1].DataType)
	require.Equal(t, "decimal(65,30)", inspection.Inspection.Columns[2].DataType)
	require.True(t, inspection.Inspection.Columns[0].PrimaryKey)
	require.True(t, inspection.Inspection.Columns[1].PrimaryKey)
	require.False(t, inspection.Inspection.Columns[2].PrimaryKey)
	require.Equal(t, []string{"bucket_start", "record_id"}, inspection.Evidence.PrimaryKey)
	require.Equal(t, "utf8mb4_0900_bin", inspection.Evidence.Collation)
	require.Equal(t, "utf8mb4_0900_bin", inspection.Evidence.ColumnCollations["record_id"])
	require.Equal(t, "42", *inspection.Evidence.ColumnDefaults["value_decimal"])
	require.True(t, inspection.Evidence.JSONColumns["provenance"])
}

func TestMySQLManagedInspectionRejectsUnusableIndexes(t *testing.T) {
	cases := map[string]mysqlManagedIndexFacts{
		"prefix":     {name: "idx", columns: []string{"workspace_id"}, indexType: "BTREE", prefix: true, visible: true},
		"functional": {name: "idx", columns: []string{"LOWER(workspace_id)"}, indexType: "BTREE", visible: true},
		"invisible":  {name: "idx", columns: []string{"workspace_id"}, indexType: "BTREE", visible: false},
		"ignored":    {name: "idx", columns: []string{"workspace_id"}, indexType: "BTREE", visible: true, ignored: true},
		"non_btree":  {name: "idx", columns: []string{"workspace_id"}, indexType: "HASH", visible: true},
	}
	for name, index := range cases {
		t.Run(name, func(t *testing.T) {
			session := &mysqlManagedTestSession{
				server: mysqlManagedBaselineFacts(),
				table: mysqlManagedTableFacts{
					status: TableInspectionExists, tableType: "BASE TABLE", engine: "InnoDB",
					columns: []mysqlManagedColumnFacts{{name: "workspace_id", dataType: "varchar", columnType: "varchar(191)"}},
					indexes: []mysqlManagedIndexFacts{index},
				},
			}
			inspection, err := mysqlManagedInspectionCall(t, newMySQLManagedTestInspector(t, session, nil))
			require.ErrorIs(t, err, ErrMySQLManagedIndexUnusable)
			require.NotNil(t, inspection)
			require.Equal(t, TableInspectionExists, inspection.Inspection.Status)
		})
	}
}

func TestMySQLManagedInspectionAllowsOrderedUsableIndex(t *testing.T) {
	session := &mysqlManagedTestSession{
		server: mysqlManagedBaselineFacts(),
		table: mysqlManagedTableFacts{
			status: TableInspectionExists, tableType: "BASE TABLE", engine: "InnoDB",
			columns: []mysqlManagedColumnFacts{
				{name: "workspace_id", dataType: "varchar", columnType: "varchar(191)"},
				{name: "tag_id", dataType: "varchar", columnType: "varchar(191)"},
			},
			indexes: []mysqlManagedIndexFacts{{
				name: "idx_ws_tag", columns: []string{"workspace_id", "tag_id"}, indexType: "BTREE", visible: true,
			}},
		},
	}
	inspection, err := mysqlManagedInspectionCall(t, newMySQLManagedTestInspector(t, session, nil))
	require.NoError(t, err)
	require.Len(t, inspection.Indexes, 1)
	require.Equal(t, []string{"workspace_id", "tag_id"}, inspection.Indexes[0].Columns)
}

func TestMySQLManagedInspectionCapabilityAndServerBinding(t *testing.T) {
	cases := map[string]struct {
		mutate func(*mysqlManagedServerFacts)
	}{
		"capability changed": {func(f *mysqlManagedServerFacts) { f.capabilityRevision = "cap-2" }},
		"empty capability":   {func(f *mysqlManagedServerFacts) { f.capabilityRevision = "" }},
		"database mismatch":  {func(f *mysqlManagedServerFacts) { f.database = "other" }},
		"family mismatch":    {func(f *mysqlManagedServerFacts) { f.family = "postgres" }},
		"version unverified": {func(f *mysqlManagedServerFacts) { f.version = "8.0.36" }},
		"metadata forbidden": {func(f *mysqlManagedServerFacts) { f.metadataStatus = "forbidden" }},
	}
	for name, testCase := range cases {
		t.Run(name, func(t *testing.T) {
			server := mysqlManagedBaselineFacts()
			testCase.mutate(&server)
			session := &mysqlManagedTestSession{server: server, table: mysqlManagedTableFacts{status: TableInspectionMissing}}
			_, err := mysqlManagedInspectionCall(t, newMySQLManagedTestInspector(t, session, nil))
			require.Error(t, err)
		})
	}
}

func TestMySQLManagedInspectionRequiresVerifiedTransport(t *testing.T) {
	connector := mysqlManagedTestConnector()
	connector.ConnectionConfig = `{"dsn":"user:pass@tcp(127.0.0.1:3306)/gw"}`
	session := &mysqlManagedTestSession{server: mysqlManagedBaselineFacts(), table: mysqlManagedTableFacts{status: TableInspectionMissing}}
	_, err := mysqlManagedInspectionCall(t, newMySQLManagedTestInspector(t, session, connector))
	require.Error(t, err)
	require.Empty(t, session.execs)
}

func TestMySQLManagedStatementOnlyAllowsCreateTableOrIndex(t *testing.T) {
	allowed := []string{
		"CREATE TABLE `gw_long_table` (`bucket_start` DATETIME(6) NOT NULL);",
		"CREATE INDEX `idx_ws` ON `gw_long_table` (`workspace_id`);",
	}
	for _, statement := range allowed {
		t.Run("allows "+statement[:12], func(t *testing.T) {
			session := &mysqlManagedTestSession{server: mysqlManagedBaselineFacts()}
			inspector := newMySQLManagedTestInspector(t, session, nil)
			require.NoError(t, inspector.ExecuteMySQLManagedStatementAtRevision(t.Context(), "conn-mysql", "rev-1", "gw", "cap-1", statement))
			require.Equal(t, []string{statement}, session.execs)
		})
	}

	refused := []string{
		"ALTER TABLE `gw_long_table` ADD COLUMN x INT;",
		"DROP TABLE `gw_long_table`;",
		"TRUNCATE TABLE `gw_long_table`;",
		"RENAME TABLE `a` TO `b`;",
		"CREATE DATABASE other;",
		"CREATE TABLE IF NOT EXISTS `gw_long_table` (`x` INT);",
		"CREATE TABLE `a` (`x` INT); CREATE INDEX `i` ON `a` (`x`);",
	}
	for _, statement := range refused {
		t.Run("refuses "+statement[:len(statement)-1], func(t *testing.T) {
			session := &mysqlManagedTestSession{server: mysqlManagedBaselineFacts()}
			inspector := newMySQLManagedTestInspector(t, session, nil)
			err := inspector.ExecuteMySQLManagedStatementAtRevision(t.Context(), "conn-mysql", "rev-1", "gw", "cap-1", statement)
			require.ErrorIs(t, err, ErrMySQLManagedStatementRefused)
			require.Empty(t, session.execs)
		})
	}
}

func TestMySQLManagedStatementErrorIsNotRolledBack(t *testing.T) {
	session := &mysqlManagedTestSession{server: mysqlManagedBaselineFacts(), execErr: errors.New("ddl outcome unknown")}
	inspector := newMySQLManagedTestInspector(t, session, nil)
	err := inspector.ExecuteMySQLManagedStatementAtRevision(t.Context(), "conn-mysql", "rev-1", "gw", "cap-1",
		"CREATE TABLE `gw_long_table` (`x` INT);")
	require.Error(t, err)
	require.False(t, strings.Contains(strings.ToLower(err.Error()), "rolled back"))
	require.Len(t, session.execs, 1)
}

func TestMySQLManagedAbsenceClassification(t *testing.T) {
	require.Equal(t, TableInspectionMissing, mysqlManagedAbsenceFacts(&mysqldriver.MySQLError{Number: 1146}).status)
	require.Equal(t, TableInspectionForbidden, mysqlManagedAbsenceFacts(&mysqldriver.MySQLError{Number: 1142}).status)
	require.Equal(t, TableInspectionForbidden, mysqlManagedAbsenceFacts(&mysqldriver.MySQLError{Number: 1044}).status)
	require.Equal(t, TableInspectionFailed, mysqlManagedAbsenceFacts(errors.New("driver unavailable")).status)
	require.Equal(t, inspectionReasonExistenceUnconfirmed, mysqlManagedAbsenceFacts(errors.New("driver unavailable")).reason)
}

func TestMySQLRecordingDatabaseScope(t *testing.T) {
	database, err := MySQLRecordingDatabaseScope(`{"dsn":"user:pass@tcp(127.0.0.1:3306)/gw?tls=true"}`)
	require.NoError(t, err)
	require.Equal(t, "gw", database)

	database, err = MySQLRecordingDatabaseScope(`{"host":"127.0.0.1","user":"user","password":"secret","database":"gw2"}`)
	require.NoError(t, err)
	require.Equal(t, "gw2", database)

	for _, raw := range []string{
		`{"dsn":"user:pass@tcp(127.0.0.1:3306)/?tls=true"}`,
		`{"host":"127.0.0.1","user":"user","password":"secret"}`,
	} {
		_, err := MySQLRecordingDatabaseScope(raw)
		require.Error(t, err)
		require.NotContains(t, err.Error(), "pass")
		require.NotContains(t, err.Error(), "secret")
	}

	_, err = MySQLRecordingDatabaseScope(`not-json`)
	require.Error(t, err)
	require.NotContains(t, err.Error(), "not-json")
}

func TestManagedTableInspectorProbeMySQLRequiresService(t *testing.T) {
	var inspector *ManagedTableInspector
	_, err := inspector.ProbeMySQLRecordingAtRevision(t.Context(), "conn-mysql", "rev-1")
	require.Error(t, err)
}

func TestMySQLManagedMariaDBDefaultDecoding(t *testing.T) {
	stringDefault := "'gateway'"
	decoded, err := decodeMySQLManagedDefault("mariadb", &stringDefault)
	require.NoError(t, err)
	require.Equal(t, "gateway", *decoded)

	quoted := "'a''b'"
	decoded, err = decodeMySQLManagedDefault("mariadb", &quoted)
	require.NoError(t, err)
	require.Equal(t, "a'b", *decoded)

	introduced := "_utf8mb4'gateway'"
	decoded, err = decodeMySQLManagedDefault("mariadb", &introduced)
	require.NoError(t, err)
	require.Equal(t, "gateway", *decoded)

	numeric := "12000.4"
	decoded, err = decodeMySQLManagedDefault("mariadb", &numeric)
	require.NoError(t, err)
	require.Equal(t, "12000.4", *decoded)

	for _, unsupported := range []string{"current_timestamp()", "(uuid())", "b'1'", "'unterminated"} {
		raw := unsupported
		_, err := decodeMySQLManagedDefault("mariadb", &raw)
		require.ErrorIs(t, err, ErrMySQLManagedTableUnconfirmed, unsupported)
	}

	backslash := `'a\'b'`
	_, err = decodeMySQLManagedDefault("mariadb", &backslash)
	require.ErrorIs(t, err, ErrMySQLManagedTableUnconfirmed)

	mysqlRaw := "gateway"
	decoded, err = decodeMySQLManagedDefault("mysql", &mysqlRaw)
	require.NoError(t, err)
	require.Equal(t, "gateway", *decoded)

	decoded, err = decodeMySQLManagedDefault("mariadb", nil)
	require.NoError(t, err)
	require.Nil(t, decoded)
}
