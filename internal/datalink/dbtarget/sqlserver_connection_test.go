package dbtarget

import (
	"crypto/tls"
	"maps"
	"net/url"
	"path/filepath"
	"testing"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"

	"github.com/microsoft/go-mssqldb/msdsn"
	"github.com/stretchr/testify/require"
)

func sqlServerConnectorConfigForTest() ConnectionConfig {
	return ConnectionConfig{
		"host": "sql.fixture.invalid", "port": "1433", "database": "fixture_db", "schema": "dbo",
		"user": "fixture_user", "password": " test-only +;&@ 密碼 ",
	}
}

func TestSQLServerConnectorScope(t *testing.T) {
	t.Parallel()
	t.Run("structured credentials retain exact password and explicit database", func(t *testing.T) {
		config := sqlServerConnectorConfigForTest()
		before := maps.Clone(config)
		got, err := buildExternalDBConfig(schema.DatabaseConnectorKindSQLServer, config)
		require.NoError(t, err)
		require.Equal(t, datalinkbase.DBType("sqlserver"), got.Type)
		dsn, err := url.Parse(got.DSN)
		require.NoError(t, err)
		require.Equal(t, "sqlserver", dsn.Scheme)
		require.Equal(t, "sql.fixture.invalid:1433", dsn.Host)
		require.Equal(t, "fixture_user", dsn.User.Username())
		password, present := dsn.User.Password()
		require.True(t, present)
		require.Equal(t, config["password"], password)
		require.Equal(t, "fixture_db", dsn.Query().Get("database"))
		require.Equal(t, "true", dsn.Query().Get("encrypt"))
		require.Equal(t, "false", dsn.Query().Get("TrustServerCertificate"))
		require.Equal(t, "1.2", dsn.Query().Get("tlsmin"))
		require.Equal(t, before, config, "building a connection must not rewrite saved settings")
	})
	t.Run("schema identity changes invalidate old confirmation", func(t *testing.T) {
		before := sqlServerConnectorConfigForTest()
		after := maps.Clone(before)
		after["schema"] = "other_schema"
		require.True(t, connectorIdentityChanged(schema.DatabaseConnectorKindSQLServer, before,
			schema.DatabaseConnectorKindSQLServer, after))
		require.False(t, connectorIdentityChanged(schema.DatabaseConnectorKindPostgres, before,
			schema.DatabaseConnectorKindPostgres, after), "the existing PostgreSQL identity contract stays unchanged")
	})
	for _, key := range []string{"host", "database", "schema", "user", "password"} {
		t.Run("missing "+key+" is rejected", func(t *testing.T) {
			config := sqlServerConnectorConfigForTest()
			delete(config, key)
			_, err := buildExternalDBConfig(schema.DatabaseConnectorKindSQLServer, config)
			require.Error(t, err)
			require.NotContains(t, err.Error(), "test-only")
		})
	}
}

func TestSQLServerTLSRejected(t *testing.T) {
	t.Parallel()
	for name, patch := range map[string]ConnectionConfig{
		"encryption disabled":                       {"encrypt": "disable"},
		"login only encryption":                     {"encrypt": "false"},
		"opportunistic encryption":                  {"encrypt": "optional"},
		"certificate verification bypass":           {"trustservercertificate": true},
		"alternate verification bypass":             {"trust_server_certificate": "true"},
		"legacy insecure TLS alias":                 {"tls": "skip-verify"},
		"legacy TLS disabled":                       {"use_tls": false},
		"server certificate pin skips hostname":     {"servercertificate": "fixture-server.pem"},
		"obsolete TLS minimum":                      {"tlsmin": "1.0"},
		"raw DSN cannot bypass structured identity": {"dsn": "sqlserver://other.invalid?encrypt=disable"},
		"invalid port":                              {"port": "1433;encrypt=disable"},
		"integrated authentication user":            {"user": `DOMAIN\fixture_user`},
	} {
		t.Run(name, func(t *testing.T) {
			config := sqlServerConnectorConfigForTest()
			maps.Copy(config, patch)
			_, err := buildExternalDBConfig(schema.DatabaseConnectorKindSQLServer, config)
			require.Error(t, err)
			require.NotContains(t, err.Error(), "test-only")
		})
	}
}

func TestSQLServerVerifiedTLSDriverConfiguration(t *testing.T) {
	t.Parallel()
	config := sqlServerConnectorConfigForTest()
	config["encrypt"] = "strict"
	config["hostnameincertificate"] = "certificate.fixture.invalid"
	settings, err := sqlServerConnectionSettingsFor(config)
	require.NoError(t, err)
	parsed, err := msdsn.Parse(settings.DSN)
	require.NoError(t, err)
	require.Equal(t, msdsn.Encryption(msdsn.EncryptionStrict), parsed.Encryption)
	require.False(t, parsed.TrustServerCertificate)
	require.NotNil(t, parsed.TLSConfig)
	require.False(t, parsed.TLSConfig.InsecureSkipVerify)
	require.Equal(t, "certificate.fixture.invalid", parsed.TLSConfig.ServerName)
	require.Equal(t, uint16(tls.VersionTLS12), parsed.TLSConfig.MinVersion)
}

func TestSQLServerInvalidCertificateConfigurationRejected(t *testing.T) {
	t.Parallel()
	config := sqlServerConnectorConfigForTest()
	config["certificate"] = filepath.Join(t.TempDir(), "missing-public-ca.pem")
	_, err := sqlServerConnectionSettingsFor(config)
	require.Error(t, err)
	require.NotContains(t, err.Error(), "missing-public-ca.pem")
	require.NotContains(t, err.Error(), "test-only")
}
