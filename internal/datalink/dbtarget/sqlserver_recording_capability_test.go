package dbtarget

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"maps"
	"strings"
	"testing"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type sqlServerRecordingMetadataScript struct {
	serverVersion        string
	serverEdition        string
	engineEdition        int64
	database             string
	compatibilityLevel   int64
	expectedSchema       string
	visibleSchema        string
	viewDefinition       bool
	schemaViewDefinition bool
	createTable          bool
	alterSchema          bool
	selectSchema         bool
	insertSchema         bool
	deleteSchema         bool
	alterAnyDataspace    bool
	errorQuery           string
	queryError           error
	queryHasDeadline     bool
	beginCalls           int
	readOnly             bool
	commitCalls          int
	rollbackCalls        int
	execCalls            int
	queries              []string
}

type sqlServerRecordingScriptConnector struct {
	script *sqlServerRecordingMetadataScript
}

func (c sqlServerRecordingScriptConnector) Connect(context.Context) (driver.Conn, error) {
	return &sqlServerRecordingScriptConn{script: c.script}, nil
}

func (sqlServerRecordingScriptConnector) Driver() driver.Driver {
	return sqlServerRecordingScriptDriver{}
}

type sqlServerRecordingScriptDriver struct{}

func (sqlServerRecordingScriptDriver) Open(string) (driver.Conn, error) {
	return nil, errors.New("scripted connector must be used")
}

type sqlServerRecordingScriptConn struct {
	script *sqlServerRecordingMetadataScript
}

func (c *sqlServerRecordingScriptConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("prepared statements are not scripted")
}

func (*sqlServerRecordingScriptConn) Close() error { return nil }

func (c *sqlServerRecordingScriptConn) Begin() (driver.Tx, error) {
	return c.BeginTx(context.Background(), driver.TxOptions{})
}

func (c *sqlServerRecordingScriptConn) BeginTx(_ context.Context, options driver.TxOptions) (driver.Tx, error) {
	c.script.beginCalls++
	c.script.readOnly = options.ReadOnly
	return sqlServerRecordingScriptTx{script: c.script}, nil
}

func (c *sqlServerRecordingScriptConn) QueryContext(
	ctx context.Context,
	query string,
	args []driver.NamedValue,
) (driver.Rows, error) {
	if _, ok := ctx.Deadline(); ok {
		c.script.queryHasDeadline = true
	}
	normalized := strings.Join(strings.Fields(query), " ")
	c.script.queries = append(c.script.queries, normalized)
	if normalized == c.script.errorQuery && c.script.queryError != nil {
		return nil, c.script.queryError
	}
	switch normalized {
	case sqlServerRecordingMetadataQuery:
		return sqlServerRecordingRows([]string{"version", "edition", "engine_edition", "database", "compatibility"}, []driver.Value{
			c.script.serverVersion, c.script.serverEdition, c.script.engineEdition, c.script.database, c.script.compatibilityLevel,
		}), nil
	case sqlServerRecordingVisibilityPermissionQuery:
		if len(args) != 1 || args[0].Value != c.script.expectedSchema {
			return nil, errors.New("unexpected visibility scope parameter")
		}
		return sqlServerRecordingRows([]string{"database_view_definition", "schema_view_definition"}, []driver.Value{
			c.script.viewDefinition, c.script.schemaViewDefinition,
		}), nil
	case sqlServerRecordingSchemaQuery:
		if len(args) != 1 || args[0].Value != c.script.expectedSchema {
			return nil, errors.New("unexpected schema parameter")
		}
		if c.script.visibleSchema == "" {
			return sqlServerRecordingRows([]string{"name"}), nil
		}
		return sqlServerRecordingRows([]string{"name"}, []driver.Value{c.script.visibleSchema}), nil
	case sqlServerRecordingPermissionQuery:
		if len(args) != 1 || args[0].Value != c.script.expectedSchema {
			return nil, errors.New("unexpected permission scope parameter")
		}
		return sqlServerRecordingRows([]string{
			"create_table", "alter_schema", "select_schema", "insert_schema", "delete_schema", "alter_any_dataspace",
		}, []driver.Value{
			c.script.createTable, c.script.alterSchema, c.script.selectSchema, c.script.insertSchema,
			c.script.deleteSchema, c.script.alterAnyDataspace,
		}), nil
	default:
		return nil, errors.New("unexpected SQL Server metadata query")
	}
}

func (c *sqlServerRecordingScriptConn) ExecContext(context.Context, string, []driver.NamedValue) (driver.Result, error) {
	c.script.execCalls++
	return nil, errors.New("metadata probe must not execute statements")
}

type sqlServerRecordingScriptTx struct {
	script *sqlServerRecordingMetadataScript
}

func (tx sqlServerRecordingScriptTx) Commit() error {
	tx.script.commitCalls++
	return nil
}

func (tx sqlServerRecordingScriptTx) Rollback() error {
	tx.script.rollbackCalls++
	return nil
}

type sqlServerRecordingScriptRows struct {
	columns []string
	values  [][]driver.Value
	next    int
}

func sqlServerRecordingRows(columns []string, values ...[]driver.Value) *sqlServerRecordingScriptRows {
	return &sqlServerRecordingScriptRows{columns: columns, values: values}
}

func (r *sqlServerRecordingScriptRows) Columns() []string { return r.columns }
func (*sqlServerRecordingScriptRows) Close() error        { return nil }

func (r *sqlServerRecordingScriptRows) Next(destination []driver.Value) error {
	if r.next >= len(r.values) {
		return io.EOF
	}
	copy(destination, r.values[r.next])
	r.next++
	return nil
}

func openSQLServerRecordingMetadataScript(t *testing.T, script *sqlServerRecordingMetadataScript) *sql.DB {
	t.Helper()
	db := sql.OpenDB(sqlServerRecordingScriptConnector{script: script})
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	return db
}

func sqlServerRecordingMetadataFixture() *sqlServerRecordingMetadataScript {
	return &sqlServerRecordingMetadataScript{
		serverVersion:        "16.0.1000.6",
		serverEdition:        "Developer Edition (64-bit)",
		engineEdition:        3,
		database:             "fixture_db",
		compatibilityLevel:   160,
		expectedSchema:       "dbo",
		visibleSchema:        "dbo",
		viewDefinition:       true,
		schemaViewDefinition: true,
		createTable:          true,
		alterSchema:          true,
		selectSchema:         true,
		insertSchema:         true,
		deleteSchema:         true,
		alterAnyDataspace:    true,
	}
}

func TestSQLServerCapabilityGate(t *testing.T) {
	script := sqlServerRecordingMetadataFixture()
	metadata := inspectSQLServerRecordingMetadata(t.Context(), openSQLServerRecordingMetadataScript(t, script), "fixture_db", "dbo")
	connector := &schema.DatabaseConnector{
		ID: "sqlserver-connector", Kind: schema.DatabaseConnectorKindSQLServer, IdentityRevision: "revision-1",
	}
	capability := projectSQLServerRecordingCapability(connector, metadata, true)

	require.True(t, script.queryHasDeadline, "metadata reads must have a bounded query context")
	require.Equal(t, []string{
		sqlServerRecordingMetadataQuery,
		sqlServerRecordingVisibilityPermissionQuery,
		sqlServerRecordingSchemaQuery,
		sqlServerRecordingPermissionQuery,
	}, script.queries)
	require.Zero(t, script.beginCalls, "the probe must not assume a read-only transaction is enforced")
	require.Zero(t, script.execCalls, "the probe must not execute DDL or DML")
	require.Equal(t, "verified", capability.MetadataStatus)
	require.True(t, capability.Connected)
	require.Equal(t, "fixture_db", capability.Database)
	require.Equal(t, "dbo", capability.Schema)
	require.Equal(t, "16.0.1000.6", capability.ServerVersion)
	require.Equal(t, "Developer Edition (64-bit)", capability.ServerEdition)
	require.True(t, capability.Permissions.Verified)
	require.True(t, capability.Permissions.DatabaseViewDefinition)
	require.True(t, capability.Permissions.SchemaViewDefinition)
	require.True(t, capability.SupportsManagedLong)
	require.True(t, capability.SupportsNativePartition)
	require.True(t, capability.SupportsGroupTestWrite)
	require.True(t, capability.SupportsMissingRepair)
	require.Empty(t, capability.Reasons)
	require.NotEmpty(t, capability.CapabilityRevision)
	require.NotContains(t, capability.CapabilityRevision, "%!", "revision formatting must match all metadata argument types")
	require.NotContains(t, capability.CapabilityRevision, "test-only")

	changedVersion := metadata
	changedVersion.serverVersion = "16.0.1001.1"
	require.NotEqual(t, capability.CapabilityRevision,
		projectSQLServerRecordingCapability(connector, changedVersion, true).CapabilityRevision)
	changedConnector := *connector
	changedConnector.ID = "another-connector"
	require.NotEqual(t, capability.CapabilityRevision,
		projectSQLServerRecordingCapability(&changedConnector, metadata, true).CapabilityRevision)
}

func TestSQLServerCapabilityMetadataFailsClosed(t *testing.T) {
	connector := &schema.DatabaseConnector{ID: "sqlserver-connector", Kind: schema.DatabaseConnectorKindSQLServer}
	tests := []struct {
		name                string
		mutate              func(*sqlServerRecordingMetadataScript)
		wantStatus          string
		wantReason          string
		wantNotReason       string
		wantManagedLong     bool
		wantNativePartition bool
		wantGroupTestWrite  bool
		wantMissingRepair   bool
	}{
		{
			name: "connected but metadata permission is denied",
			mutate: func(script *sqlServerRecordingMetadataScript) {
				script.errorQuery = sqlServerRecordingVisibilityPermissionQuery
				script.queryError = errors.New("The VIEW DEFINITION permission was denied")
			},
			wantStatus: sqlServerRecordingMetadataForbidden,
			wantReason: "metadata_forbidden",
		},
		{
			name: "missing schema is distinct when database metadata visibility is verified",
			mutate: func(script *sqlServerRecordingMetadataScript) {
				script.visibleSchema = ""
				script.schemaViewDefinition = false
			},
			wantStatus: sqlServerRecordingMetadataMissing,
			wantReason: "schema_missing",
		},
		{
			name: "selected schema visibility does not require database-wide visibility",
			mutate: func(script *sqlServerRecordingMetadataScript) {
				script.viewDefinition = false
				script.schemaViewDefinition = true
			},
			wantStatus:          sqlServerRecordingMetadataVerified,
			wantReason:          "native_partition_permissions_denied",
			wantManagedLong:     true,
			wantNativePartition: false,
			wantGroupTestWrite:  true,
			wantMissingRepair:   true,
		},
		{
			name: "invisible schema is not reported missing",
			mutate: func(script *sqlServerRecordingMetadataScript) {
				script.visibleSchema = ""
				script.viewDefinition = false
				script.schemaViewDefinition = false
			},
			wantStatus:    sqlServerRecordingMetadataUnconfirmed,
			wantReason:    "schema_visibility_unconfirmed",
			wantNotReason: "schema_missing",
		},
		{
			name: "visible schema name without either view definition permission is forbidden",
			mutate: func(script *sqlServerRecordingMetadataScript) {
				script.viewDefinition = false
				script.schemaViewDefinition = false
				script.createTable = false
				script.alterSchema = false
				script.selectSchema = false
				script.insertSchema = false
				script.deleteSchema = false
				script.alterAnyDataspace = false
			},
			wantStatus: sqlServerRecordingMetadataForbidden,
			wantReason: "metadata_forbidden",
		},
		{
			name: "metadata failure is not reported missing",
			mutate: func(script *sqlServerRecordingMetadataScript) {
				script.errorQuery = sqlServerRecordingMetadataQuery
				script.queryError = errors.New("catalog query failed")
			},
			wantStatus:    sqlServerRecordingMetadataUnconfirmed,
			wantReason:    "metadata_unconfirmed",
			wantNotReason: "schema_missing",
		},
		{
			name:       "database mismatch",
			mutate:     func(script *sqlServerRecordingMetadataScript) { script.database = "other_db" },
			wantStatus: sqlServerRecordingMetadataUnconfirmed,
			wantReason: "database_scope_mismatch",
		},
		{
			name:       "schema mismatch",
			mutate:     func(script *sqlServerRecordingMetadataScript) { script.visibleSchema = "other_schema" },
			wantStatus: sqlServerRecordingMetadataUnconfirmed,
			wantReason: "schema_scope_mismatch",
		},
		{
			name: "operation permission denial stays visible and gated",
			mutate: func(script *sqlServerRecordingMetadataScript) {
				script.insertSchema = false
				script.alterAnyDataspace = false
			},
			wantStatus: sqlServerRecordingMetadataVerified,
			wantReason: "operation_permissions_denied",
		},
		{
			name: "partition permission denial does not imply plain operation denial",
			mutate: func(script *sqlServerRecordingMetadataScript) {
				script.alterAnyDataspace = false
			},
			wantStatus:          sqlServerRecordingMetadataVerified,
			wantReason:          "native_partition_permissions_denied",
			wantNotReason:       "operation_permissions_denied",
			wantManagedLong:     true,
			wantNativePartition: false,
			wantGroupTestWrite:  true,
			wantMissingRepair:   true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			script := sqlServerRecordingMetadataFixture()
			test.mutate(script)
			metadata := inspectSQLServerRecordingMetadata(t.Context(), openSQLServerRecordingMetadataScript(t, script), "fixture_db", "dbo")
			capability := projectSQLServerRecordingCapability(connector, metadata, true)
			require.True(t, capability.Connected, "metadata denial must not erase successful connection evidence")
			require.Equal(t, test.wantStatus, capability.MetadataStatus)
			if test.wantReason != "" {
				require.Contains(t, capability.Reasons, test.wantReason)
			}
			if test.wantNotReason != "" {
				require.NotContains(t, capability.Reasons, test.wantNotReason)
			}
			require.Equal(t, test.wantManagedLong, capability.SupportsManagedLong)
			require.Equal(t, test.wantNativePartition, capability.SupportsNativePartition)
			require.Equal(t, test.wantGroupTestWrite, capability.SupportsGroupTestWrite)
			require.Equal(t, test.wantMissingRepair, capability.SupportsMissingRepair)
		})
	}
}

func TestSQLServerRecordingProbeUsesSavedIdentityAndDirectOpen(t *testing.T) {
	internalDB := openMigratedTestDB(t)
	repo := NewSQLConnectorRepository(internalDB)
	service := NewConnectorService(repo)
	config := sqlServerConnectorConfigForTest()
	connector := createSQLServerRecordingSavedConnector(t, repo, "sqlserver-probe-test", schema.DatabaseConnectorKindSQLServer, config, true)

	originalOpen := openExternalDBManagerFunc
	openCalls := 0
	var openedConfig ConnectionConfig
	openExternalDBManagerFunc = func(kind schema.DatabaseConnectorKind, got ConnectionConfig) (*datalinkbase.DBManager, error) {
		openCalls++
		require.Equal(t, schema.DatabaseConnectorKindSQLServer, kind)
		openedConfig = got
		return nil, errors.New("fixture server unavailable")
	}
	t.Cleanup(func() { openExternalDBManagerFunc = originalOpen })

	capability, err := service.ProbeSQLServerRecordingAtRevision(t.Context(), connector.ID, connector.IdentityRevision)
	require.NoError(t, err)
	require.Equal(t, 1, openCalls)
	require.Equal(t, "fixture_db", capability.Database)
	require.Equal(t, config["password"], openedConfig["password"], "saved credentials must not be trimmed or rewritten")
	require.False(t, capability.Connected)
	require.Contains(t, capability.Reasons, inspectionReasonConnectionFailed)
	require.NotContains(t, capability.CapabilityRevision, config["password"])

	_, err = service.ProbeSQLServerRecordingAtRevision(t.Context(), connector.ID, "stale-revision")
	require.ErrorIs(t, err, ErrConnectorRevisionConflict)
	require.Equal(t, 1, openCalls, "a stale identity must not open an external connection")

	insecureConfig := maps.Clone(config)
	insecureConfig["trust_server_certificate"] = true
	insecureConnector := createSQLServerRecordingSavedConnector(t, repo, "sqlserver-insecure-test", schema.DatabaseConnectorKindSQLServer, insecureConfig, true)
	capability, err = service.ProbeSQLServerRecordingAtRevision(t.Context(), insecureConnector.ID, insecureConnector.IdentityRevision)
	require.NoError(t, err)
	require.Contains(t, capability.Reasons, "connection_settings_unverified")
	require.Equal(t, 1, openCalls, "TLS validation must fail before opening an external connection")
}

func TestSQLServerRecordingProbeRejectsWrongKindAndDisabledConnector(t *testing.T) {
	internalDB := openMigratedTestDB(t)
	repo := NewSQLConnectorRepository(internalDB)
	service := NewConnectorService(repo)
	config := sqlServerConnectorConfigForTest()
	wrongKind := createSQLServerRecordingSavedConnector(t, repo, "not-sqlserver", schema.DatabaseConnectorKindPostgres, config, true)
	disabled := createSQLServerRecordingSavedConnector(t, repo, "disabled-sqlserver", schema.DatabaseConnectorKindSQLServer, config, false)

	originalOpen := openExternalDBManagerFunc
	openCalls := 0
	openExternalDBManagerFunc = func(schema.DatabaseConnectorKind, ConnectionConfig) (*datalinkbase.DBManager, error) {
		openCalls++
		return nil, errors.New("unexpected open")
	}
	t.Cleanup(func() { openExternalDBManagerFunc = originalOpen })

	_, err := service.ProbeSQLServerRecordingAtRevision(t.Context(), wrongKind.ID, wrongKind.IdentityRevision)
	require.Error(t, err)
	_, err = service.ProbeSQLServerRecordingAtRevision(t.Context(), disabled.ID, disabled.IdentityRevision)
	require.ErrorIs(t, err, ErrConnectorDisabled)
	require.Zero(t, openCalls)
}

func createSQLServerRecordingSavedConnector(
	t *testing.T,
	repo *SQLConnectorRepository,
	id string,
	kind schema.DatabaseConnectorKind,
	config ConnectionConfig,
	enabled bool,
) *schema.DatabaseConnector {
	t.Helper()
	rawConfig, err := serializeConnectionConfig(config)
	require.NoError(t, err)
	connector := &schema.DatabaseConnector{
		ID: id, Name: id, Kind: kind, ConnectionConfig: rawConfig, Enabled: enabled,
	}
	require.NoError(t, repo.Create(t.Context(), connector))
	return connector
}
