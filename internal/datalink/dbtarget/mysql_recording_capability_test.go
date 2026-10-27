package dbtarget

import (
	"context"
	"crypto/tls"
	"errors"
	"maps"
	"strings"
	"testing"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

func TestMySQLFamilyCapability(t *testing.T) {
	t.Run("actual MySQL family and version stay behind unverified operation gates", func(t *testing.T) {
		script := mysqlRecordingMetadataFixture("8.4.4", "MySQL Community Server - GPL")
		metadata := inspectMySQLRecordingMetadata(t.Context(), openMySQLRecordingMetadataScript(t, script), "metrics")
		connector := mysqlRecordingTestConnector()
		capability := projectMySQLRecordingCapability(connector, metadata, true)

		require.True(t, script.readOnly, "metadata must be read inside a read-only transaction")
		require.Equal(t, "mysql", capability.ServerFamily)
		require.Equal(t, "8.4.4", capability.ServerVersion)
		require.Equal(t, "metrics", capability.Database)
		require.Equal(t, "verified", capability.MetadataStatus)
		require.True(t, capability.Connected)
		require.True(t, capability.TransactionalEngine)
		require.False(t, capability.SupportsManagedLong)
		require.False(t, capability.SupportsNativePartition)
		require.False(t, capability.SupportsGroupTestWrite)
		require.False(t, capability.SupportsMissingRepair)
		require.Contains(t, capability.Reasons, "adapter_unverified")
		require.NotEmpty(t, capability.CapabilityRevision)
		require.NotContains(t, capability.CapabilityRevision, "saved-secret")
		require.Equal(t, []string{
			mysqlRecordingMetadataQuery,
			mysqlRecordingSchemaQuery,
			mysqlRecordingEngineQuery,
			mysqlRecordingCurrentUserQuery,
			mysqlRecordingCurrentRoleQuery,
			mysqlRecordingShowGrantsQuery,
			mysqlRecordingGlobalPrivilegesQuery,
			mysqlRecordingSchemaPrivilegesQuery,
		}, script.queries)

		changedServer := metadata
		changedServer.serverVersion = "8.4.5"
		require.NotEqual(t,
			capability.CapabilityRevision,
			projectMySQLRecordingCapability(connector, changedServer, true).CapabilityRevision,
			"a different server version invalidates the projection revision",
		)
		changedScope := metadata
		changedScope.database = "other_metrics"
		require.NotEqual(t,
			capability.CapabilityRevision,
			projectMySQLRecordingCapability(connector, changedScope, true).CapabilityRevision,
			"a different database scope invalidates the projection revision",
		)
		otherConnector := *connector
		otherConnector.ID = "connector-2"
		require.NotEqual(t,
			capability.CapabilityRevision,
			projectMySQLRecordingCapability(&otherConnector, metadata, true).CapabilityRevision,
			"different connectors must not share a projection revision",
		)
	})

	t.Run("connected but forbidden metadata remains distinct from missing", func(t *testing.T) {
		script := mysqlRecordingMetadataFixture("8.4.4", "MySQL Community Server - GPL")
		script.errorQuery = mysqlRecordingSchemaQuery
		script.queryError = &mysqldriver.MySQLError{Number: 1142, Message: "SELECT command denied"}
		metadata := inspectMySQLRecordingMetadata(t.Context(), openMySQLRecordingMetadataScript(t, script), "metrics")
		capability := projectMySQLRecordingCapability(mysqlRecordingTestConnector(), metadata, true)

		require.True(t, capability.Connected)
		require.Equal(t, "forbidden", capability.MetadataStatus)
		require.Equal(t, "mysql", capability.ServerFamily, "successful server metadata survives a later permission denial")
		require.False(t, capability.SupportsManagedLong)
		require.False(t, capability.SupportsNativePartition)
		require.Contains(t, capability.Reasons, "metadata_forbidden")
		require.NotContains(t, script.queries, "CREATE DATABASE")
		require.NotContains(t, script.queries, "CREATE TABLE")
	})

	t.Run("zero visible schemas is unconfirmed rather than missing", func(t *testing.T) {
		script := mysqlRecordingMetadataFixture("8.4.4", "MySQL Community Server - GPL")
		script.visibleDatabase = ""
		metadata := inspectMySQLRecordingMetadata(t.Context(), openMySQLRecordingMetadataScript(t, script), "metrics")
		capability := projectMySQLRecordingCapability(mysqlRecordingTestConnector(), metadata, true)

		require.True(t, capability.Connected)
		require.Equal(t, "unconfirmed", capability.MetadataStatus)
		require.Contains(t, capability.Reasons, "database_visibility_unconfirmed")
		require.NotContains(t, capability.Reasons, "database_missing")
	})

	t.Run("unknown family and versions stay unverified", func(t *testing.T) {
		for _, testCase := range []struct {
			name, version, comment, family string
		}{
			{name: "unknown family", version: "8.4.4", comment: "Percona Server", family: "unknown"},
			{name: "unknown version", version: "9.0.1", comment: "MySQL Community Server", family: "mysql"},
		} {
			t.Run(testCase.name, func(t *testing.T) {
				script := mysqlRecordingMetadataFixture(testCase.version, testCase.comment)
				metadata := inspectMySQLRecordingMetadata(t.Context(), openMySQLRecordingMetadataScript(t, script), "metrics")
				capability := projectMySQLRecordingCapability(mysqlRecordingTestConnector(), metadata, true)

				require.Equal(t, testCase.family, capability.ServerFamily)
				require.Contains(t, capability.Reasons, "adapter_unverified")
				require.Contains(t, capability.Reasons, "server_version_unverified")
				require.False(t, capability.SupportsManagedLong)
				require.False(t, capability.SupportsNativePartition)
			})
		}
	})

	t.Run("non-permission metadata errors remain unconfirmed", func(t *testing.T) {
		script := mysqlRecordingMetadataFixture("8.4.4", "MySQL Community Server - GPL")
		script.errorQuery = mysqlRecordingEngineQuery
		script.queryError = errors.New("catalog query failed")
		metadata := inspectMySQLRecordingMetadata(t.Context(), openMySQLRecordingMetadataScript(t, script), "metrics")
		capability := projectMySQLRecordingCapability(mysqlRecordingTestConnector(), metadata, true)

		require.True(t, capability.Connected)
		require.Equal(t, "mysql", capability.ServerFamily)
		require.Equal(t, "unconfirmed", capability.MetadataStatus)
		require.Contains(t, capability.Reasons, "metadata_unconfirmed")
		require.NotContains(t, strings.Join(capability.Reasons, ","), "catalog query failed")
		require.False(t, capability.TransactionalEngine)
	})

	t.Run("rollback failure retains metadata as unconfirmed", func(t *testing.T) {
		script := mysqlRecordingMetadataFixture("8.4.4", "MySQL Community Server - GPL")
		script.errorQuery = mysqlRecordingEngineQuery
		script.queryError = errors.New("catalog query failed")
		script.rollbackError = errors.New("rollback failed")
		metadata := inspectMySQLRecordingMetadata(t.Context(), openMySQLRecordingMetadataScript(t, script), "metrics")
		capability := projectMySQLRecordingCapability(mysqlRecordingTestConnector(), metadata, true)

		require.Equal(t, 1, script.rollbackCalls)
		require.True(t, capability.Connected)
		require.Equal(t, "unconfirmed", capability.MetadataStatus)
		require.Contains(t, capability.Reasons, "metadata_rollback_failed")
		require.False(t, capability.SupportsManagedLong)
	})

	t.Run("missing InnoDB is reported without opening recording gates", func(t *testing.T) {
		script := mysqlRecordingMetadataFixture("8.4.4", "MySQL Community Server - GPL")
		script.engineSupport = ""
		metadata := inspectMySQLRecordingMetadata(t.Context(), openMySQLRecordingMetadataScript(t, script), "metrics")
		capability := projectMySQLRecordingCapability(mysqlRecordingTestConnector(), metadata, true)

		require.Equal(t, "verified", capability.MetadataStatus)
		require.False(t, capability.TransactionalEngine)
		require.Contains(t, capability.Reasons, "innodb_unavailable")
		require.False(t, capability.SupportsManagedLong)
	})
}

func TestMariaDBUnverifiedGate(t *testing.T) {
	script := mysqlRecordingMetadataFixture("11.4.2-MariaDB-1:11.4.2+maria~ubu2204", "MariaDB Server")
	metadata := inspectMySQLRecordingMetadata(t.Context(), openMySQLRecordingMetadataScript(t, script), "metrics")
	capability := projectMySQLRecordingCapability(mysqlRecordingTestConnector(), metadata, true)

	require.Equal(t, "mariadb", capability.ServerFamily)
	require.Equal(t, "11.4.2-MariaDB-1:11.4.2+maria~ubu2204", capability.ServerVersion)
	require.True(t, capability.Connected)
	require.False(t, capability.SupportsManagedLong, "MySQL verification must not unlock MariaDB")
	require.False(t, capability.SupportsNativePartition)
	require.False(t, capability.SupportsGroupTestWrite)
	require.False(t, capability.SupportsMissingRepair)
	require.Contains(t, capability.Reasons, "adapter_unverified")
}

func TestMySQLFamilyCapabilityTransport(t *testing.T) {
	t.Run("verified TLS permits configured cleartext authentication", func(t *testing.T) {
		config := mysqlRecordingConnectionConfig(ConnectionConfig{
			"tls": "true", "allow_cleartext_passwords": true,
		})
		parsed := parseMySQLRecordingTransport(t, config)
		require.True(t, parsed.AllowCleartextPasswords)
		require.False(t, parsed.AllowFallbackToPlaintext)
		require.True(t, mysqlRecordingTransportVerified(parsed))
		require.Equal(t, "  password with spaces \t", parsed.Passwd)
	})

	t.Run("registered TLS config is checked by the driver's parsed TLS config", func(t *testing.T) {
		const key = "mysql-recording-capability-verified-test"
		require.NoError(t, mysqldriver.RegisterTLSConfig(key, &tls.Config{ServerName: "mysql.internal"}))
		t.Cleanup(func() { mysqldriver.DeregisterTLSConfig(key) })
		parsed := parseMySQLRecordingTransport(t, mysqlRecordingConnectionConfig(ConnectionConfig{"tls": key}))
		require.NotNil(t, parsed.TLS)
		require.False(t, parsed.TLS.InsecureSkipVerify)
		require.True(t, mysqlRecordingTransportVerified(parsed))
	})

	t.Run("opportunistic, skip-verify and unregistered TLS are refused", func(t *testing.T) {
		defaultConfig := parseMySQLRecordingTransport(t, mysqlRecordingConnectionConfig(nil))
		require.False(t, mysqlRecordingTransportVerified(defaultConfig))

		skipVerify := parseMySQLRecordingTransport(t, mysqlRecordingConnectionConfig(ConnectionConfig{"use_tls": true}))
		require.False(t, mysqlRecordingTransportVerified(skipVerify))

		unregistered, err := mysqldriver.ParseDSN("user:secret@tcp(mysql.internal:3306)/metrics?tls=not-registered")
		require.Error(t, err)
		require.Nil(t, unregistered)
	})

	t.Run("a registered TLS config that skips verification is refused", func(t *testing.T) {
		const key = "mysql-recording-capability-insecure-test"
		require.NoError(t, mysqldriver.RegisterTLSConfig(key, &tls.Config{InsecureSkipVerify: true}))
		t.Cleanup(func() { mysqldriver.DeregisterTLSConfig(key) })
		parsed := parseMySQLRecordingTransport(t, mysqlRecordingConnectionConfig(ConnectionConfig{"tls": key}))
		require.NotNil(t, parsed.TLS)
		require.True(t, parsed.TLS.InsecureSkipVerify)
		require.False(t, mysqlRecordingTransportVerified(parsed))
	})
}

func TestMySQLFamilyCapabilityProbeUsesSavedIdentityAndDirectOpen(t *testing.T) {
	internalDB := openMigratedTestDB(t)
	repo := NewSQLConnectorRepository(internalDB)
	service := NewConnectorService(repo)
	password := "  password with spaces \t"
	rawConfig, err := serializeConnectionConfig(ConnectionConfig{
		"host": "mysql.internal", "user": "reader", "password": password, "database": "metrics",
		"tls": "true", "allow_cleartext_passwords": true,
	})
	require.NoError(t, err)
	connector := &schema.DatabaseConnector{
		ID: "mysql-capability-test", Name: "mysql capability test", Kind: schema.DatabaseConnectorKindMySQL,
		ConnectionConfig: rawConfig, Enabled: true,
	}
	require.NoError(t, repo.Create(t.Context(), connector), "seed the saved connector without running the legacy probe")

	originalOpen := openExternalDBManagerFunc
	originalEnsure := ensureMySQLDatabaseIfMissingFunc
	t.Cleanup(func() {
		openExternalDBManagerFunc = originalOpen
		ensureMySQLDatabaseIfMissingFunc = originalEnsure
	})
	openCalls := 0
	ensureCalls := 0
	openExternalDBManagerFunc = func(kind schema.DatabaseConnectorKind, config ConnectionConfig) (*datalinkbase.DBManager, error) {
		openCalls++
		require.Equal(t, schema.DatabaseConnectorKindMySQL, kind)
		require.Equal(t, password, config["password"], "probing must not trim or rewrite saved credentials")
		dbConfig, configErr := buildExternalDBConfig(kind, config)
		require.NoError(t, configErr)
		parsed, parseErr := mysqldriver.ParseDSN(dbConfig.DSN)
		require.NoError(t, parseErr)
		require.True(t, mysqlRecordingTransportVerified(parsed), "explicit cleartext auth is allowed only under verified TLS")
		return nil, &mysqldriver.MySQLError{Number: 1049, Message: "Unknown database 'metrics'"}
	}
	ensureMySQLDatabaseIfMissingFunc = func(context.Context, ConnectionConfig, error) error {
		ensureCalls++
		return nil
	}

	capability, err := service.ProbeMySQLRecordingAtRevision(t.Context(), connector.ID, connector.IdentityRevision)
	require.NoError(t, err)
	require.Equal(t, 1, openCalls)
	require.Zero(t, ensureCalls, "the capability probe must not create a missing database")
	require.Equal(t, connector.ID, capability.ConnectorID)
	require.Equal(t, connector.IdentityRevision, capability.ConnectorRevision)
	require.Equal(t, "metrics", capability.Database)
	require.False(t, capability.Connected)
	require.Contains(t, capability.Reasons, "connection_failed")
	require.NotContains(t, capability.CapabilityRevision, password)
	require.NotContains(t, strings.Join(capability.Reasons, ","), "Unknown database", "driver errors must not enter the projection")

	_, err = service.ProbeMySQLRecordingAtRevision(t.Context(), connector.ID, "stale-revision")
	require.ErrorIs(t, err, ErrConnectorRevisionConflict)
	require.Equal(t, 1, openCalls, "a stale saved identity must not initiate a connection")
	require.Zero(t, ensureCalls, "a stale saved identity must not initiate database creation")
}

func mysqlRecordingTestConnector() *schema.DatabaseConnector {
	return &schema.DatabaseConnector{
		ID: "connector-1", Kind: schema.DatabaseConnectorKindMySQL,
		IdentityRevision: "saved-revision-1",
	}
}

func mysqlRecordingConnectionConfig(extra ConnectionConfig) ConnectionConfig {
	config := ConnectionConfig{
		"host": "mysql.internal", "user": "reader", "password": "  password with spaces \t", "database": "metrics",
	}
	maps.Copy(config, extra)
	return config
}

func parseMySQLRecordingTransport(t *testing.T, config ConnectionConfig) *mysqldriver.Config {
	t.Helper()
	dbConfig, err := buildExternalDBConfig(schema.DatabaseConnectorKindMySQL, config)
	require.NoError(t, err)
	parsed, err := mysqldriver.ParseDSN(dbConfig.DSN)
	require.NoError(t, err)
	return parsed
}
