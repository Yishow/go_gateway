package dbtarget

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	mysqlRecordingMetadataQuery = "SELECT VERSION(), @@version_comment, DATABASE()"
	mysqlRecordingSchemaQuery   = "SELECT SCHEMA_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?"
	mysqlRecordingEngineQuery   = "SELECT SUPPORT, TRANSACTIONS FROM information_schema.ENGINES WHERE ENGINE = 'InnoDB'"
)

type mysqlRecordingMetadataScript struct {
	version            string
	versionComment     string
	currentDatabase    string
	visibleDatabase    string
	engineSupport      string
	engineTransactions string
	currentUser        string
	currentRole        string
	currentRoleNull    bool
	globalPrivileges   []mysqlRecordingPrivilegeFixture
	schemaPrivileges   []mysqlRecordingPrivilegeFixture
	showGrants         []string
	errorQuery         string
	queryError         error
	rollbackError      error
	rollbackCalls      int
	readOnly           bool
	queries            []string
}

type mysqlRecordingPrivilegeFixture struct {
	schemaName    string
	grantee       string
	privilegeType string
}

type mysqlRecordingScriptConnector struct {
	script *mysqlRecordingMetadataScript
}

func (c mysqlRecordingScriptConnector) Connect(context.Context) (driver.Conn, error) {
	return &mysqlRecordingScriptConn{script: c.script}, nil
}

func (mysqlRecordingScriptConnector) Driver() driver.Driver {
	return mysqlRecordingScriptDriver{}
}

type mysqlRecordingScriptDriver struct{}

func (mysqlRecordingScriptDriver) Open(string) (driver.Conn, error) {
	return nil, errors.New("scripted connector must be used")
}

type mysqlRecordingScriptConn struct {
	script *mysqlRecordingMetadataScript
}

func (c *mysqlRecordingScriptConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("prepared statements are not scripted")
}

func (*mysqlRecordingScriptConn) Close() error { return nil }

func (c *mysqlRecordingScriptConn) Begin() (driver.Tx, error) {
	return c.BeginTx(context.Background(), driver.TxOptions{})
}

func (c *mysqlRecordingScriptConn) BeginTx(_ context.Context, options driver.TxOptions) (driver.Tx, error) {
	c.script.readOnly = options.ReadOnly
	return mysqlRecordingScriptTx{script: c.script}, nil
}

func (c *mysqlRecordingScriptConn) QueryContext(_ context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	normalized := strings.Join(strings.Fields(query), " ")
	c.script.queries = append(c.script.queries, normalized)
	if normalized == c.script.errorQuery && c.script.queryError != nil {
		return nil, c.script.queryError
	}

	switch normalized {
	case mysqlRecordingMetadataQuery:
		return mysqlRecordingRows([]string{"version", "version_comment", "database"}, []driver.Value{
			c.script.version, c.script.versionComment, c.script.currentDatabase,
		}), nil
	case mysqlRecordingSchemaQuery:
		if len(args) != 1 || args[0].Value != c.script.visibleDatabase || c.script.visibleDatabase == "" {
			return mysqlRecordingRows([]string{"schema_name"}), nil
		}
		return mysqlRecordingRows([]string{"schema_name"}, []driver.Value{c.script.visibleDatabase}), nil
	case mysqlRecordingEngineQuery:
		if c.script.engineSupport == "" {
			return mysqlRecordingRows([]string{"support", "transactions"}), nil
		}
		return mysqlRecordingRows([]string{"support", "transactions"}, []driver.Value{
			c.script.engineSupport, c.script.engineTransactions,
		}), nil
	case mysqlRecordingCurrentUserQuery:
		return mysqlRecordingRows([]string{"current_user"}, []driver.Value{c.script.currentUser}), nil
	case mysqlRecordingCurrentRoleQuery:
		if c.script.currentRoleNull {
			return mysqlRecordingRows([]string{"current_role"}, []driver.Value{nil}), nil
		}
		return mysqlRecordingRows([]string{"current_role"}, []driver.Value{c.script.currentRole}), nil
	case mysqlRecordingGlobalPrivilegesQuery:
		return mysqlRecordingPrivilegeRows(c.script.globalPrivileges), nil
	case mysqlRecordingSchemaPrivilegesQuery:
		return mysqlRecordingSchemaPrivilegeRows(c.script.schemaPrivileges, c.script.currentDatabase), nil
	case mysqlRecordingShowGrantsQuery:
		values := make([][]driver.Value, 0, len(c.script.showGrants))
		for _, grant := range c.script.showGrants {
			values = append(values, []driver.Value{grant})
		}
		return mysqlRecordingRows([]string{"grants"}, values...), nil
	default:
		return nil, errors.New("unexpected metadata or privilege query")
	}
}

func (c *mysqlRecordingScriptConn) ExecContext(context.Context, string, []driver.NamedValue) (driver.Result, error) {
	return nil, errors.New("metadata probe must not execute statements")
}

type mysqlRecordingScriptTx struct {
	script *mysqlRecordingMetadataScript
}

func (mysqlRecordingScriptTx) Commit() error { return nil }

func (tx mysqlRecordingScriptTx) Rollback() error {
	tx.script.rollbackCalls++
	return tx.script.rollbackError
}

type mysqlRecordingScriptRows struct {
	columns []string
	values  [][]driver.Value
	next    int
}

func mysqlRecordingRows(columns []string, values ...[]driver.Value) *mysqlRecordingScriptRows {
	return &mysqlRecordingScriptRows{columns: columns, values: values}
}

func mysqlRecordingPrivilegeRows(privileges []mysqlRecordingPrivilegeFixture) *mysqlRecordingScriptRows {
	values := make([][]driver.Value, 0, len(privileges))
	for _, privilege := range privileges {
		values = append(values, []driver.Value{privilege.grantee, privilege.privilegeType})
	}
	return mysqlRecordingRows([]string{"grantee", "privilege_type"}, values...)
}

func mysqlRecordingSchemaPrivilegeRows(privileges []mysqlRecordingPrivilegeFixture, defaultSchema string) *mysqlRecordingScriptRows {
	values := make([][]driver.Value, 0, len(privileges))
	for _, privilege := range privileges {
		schemaName := privilege.schemaName
		if schemaName == "" {
			schemaName = defaultSchema
		}
		values = append(values, []driver.Value{schemaName, privilege.grantee, privilege.privilegeType})
	}
	return mysqlRecordingRows([]string{"table_schema", "grantee", "privilege_type"}, values...)
}

func (r *mysqlRecordingScriptRows) Columns() []string { return r.columns }
func (*mysqlRecordingScriptRows) Close() error        { return nil }

func (r *mysqlRecordingScriptRows) Next(destination []driver.Value) error {
	if r.next >= len(r.values) {
		return io.EOF
	}
	copy(destination, r.values[r.next])
	r.next++
	return nil
}

func openMySQLRecordingMetadataScript(t *testing.T, script *mysqlRecordingMetadataScript) *sql.DB {
	t.Helper()
	db := sql.OpenDB(mysqlRecordingScriptConnector{script: script})
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	return db
}

func mysqlRecordingMetadataFixture(version, comment string) *mysqlRecordingMetadataScript {
	return &mysqlRecordingMetadataScript{
		version:            version,
		versionComment:     comment,
		currentDatabase:    "metrics",
		visibleDatabase:    "metrics",
		engineSupport:      "DEFAULT",
		engineTransactions: "YES",
		currentUser:        "fixture@%",
		currentRole:        "NONE",
		schemaPrivileges:   mysqlRecordingAllPrivilegeFixtures("'fixture'@'%'"),
	}
}

func mysqlRecordingAllPrivilegeFixtures(grantee string) []mysqlRecordingPrivilegeFixture {
	return []mysqlRecordingPrivilegeFixture{
		{grantee: grantee, privilegeType: "CREATE"},
		{grantee: grantee, privilegeType: "ALTER"},
		{grantee: grantee, privilegeType: "INDEX"},
		{grantee: grantee, privilegeType: "SELECT"},
		{grantee: grantee, privilegeType: "INSERT"},
		{grantee: grantee, privilegeType: "DELETE"},
	}
}
