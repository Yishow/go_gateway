package dbtarget

import (
	"context"
	"testing"

	mysqldriver "github.com/go-sql-driver/mysql"
	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildExternalDBConfig_MySQLDefaultsEnablePreferredTLSAndCleartextPasswords(t *testing.T) {
	t.Parallel()

	dbConfig, err := buildExternalDBConfig(
		schema.DatabaseConnectorKindMySQL,
		ConnectionConfig{
			"host":     "mysql.internal",
			"port":     "3307",
			"user":     "gateway",
			"password": "secret",
			"database": "gateway_metrics",
		},
	)
	require.NoError(t, err)

	parsed, err := mysqldriver.ParseDSN(dbConfig.DSN)
	require.NoError(t, err)

	assert.Equal(t, datalinkbase.DBTypeMySQL, dbConfig.Type)
	assert.Equal(t, "gateway", parsed.User)
	assert.Equal(t, "secret", parsed.Passwd)
	assert.Equal(t, "tcp", parsed.Net)
	assert.Equal(t, "mysql.internal:3307", parsed.Addr)
	assert.Equal(t, "gateway_metrics", parsed.DBName)
	assert.True(t, parsed.AllowCleartextPasswords)
	assert.Equal(t, "preferred", parsed.TLSConfig)
}

func TestProbeConnector_MySQLClassifiesConnectionFailuresAsReachableConnectorErrors(t *testing.T) {
	t.Parallel()

	status, checkedAt, message := probeConnector(
		context.Background(),
		schema.DatabaseConnectorKindMySQL,
		ConnectionConfig{
			"host":     "127.0.0.1",
			"port":     "1",
			"user":     "root",
			"password": "secret",
			"database": "gateway_metrics",
			"timeout":  "100ms",
		},
	)

	require.NotNil(t, checkedAt)
	assert.Equal(t, schema.DatabaseConnectorStatusUnreachable, status)
	assert.NotContains(t, message, "不支援的資料庫類型")
}

func TestConnectorErrorMessage_MySQLUnknownPluginReturnsActionableGuidance(t *testing.T) {
	t.Parallel()

	message := connectorErrorMessage(
		schema.DatabaseConnectorKindMySQL,
		mysqldriver.ErrUnknownPlugin,
	)

	assert.Contains(t, message, "未支援的 MySQL/MariaDB 驗證插件")
	assert.Contains(t, message, "pam_use_cleartext_plugin")
	assert.Contains(t, message, "mysql_native_password")
	assert.Contains(t, message, "mysql_clear_password")
}
