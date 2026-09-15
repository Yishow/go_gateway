package dbtarget

import (
	"context"
	"testing"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"

	mysqldriver "github.com/go-sql-driver/mysql"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildExternalDBConfig_MySQLDefaultsDisableCleartextPasswords(t *testing.T) {
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
	assert.False(t, parsed.AllowCleartextPasswords)
	// 預設維持機會性加密：有 TLS 就用，沒有才退回明文，且密碼永遠不是明文。
	assert.Equal(t, "preferred", parsed.TLSConfig)
	assert.True(t, parsed.AllowFallbackToPlaintext)
}

// cleartext 認證會把密碼明文送給伺服器，因此啟用時必須加密且不得退回明文。
func TestBuildExternalDBConfig_MySQLResolvesCleartextAndTLS(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name              string
		extraConfig       ConnectionConfig
		wantCleartext     bool
		wantTLSConfig     string
		wantAllowFallback bool
	}{
		{
			name:              "defaults use opportunistic tls",
			extraConfig:       ConnectionConfig{},
			wantCleartext:     false,
			wantTLSConfig:     "preferred",
			wantAllowFallback: true,
		},
		{
			name:              "use_tls encrypts without plaintext fallback",
			extraConfig:       ConnectionConfig{"use_tls": true},
			wantCleartext:     false,
			wantTLSConfig:     "skip-verify",
			wantAllowFallback: false,
		},
		{
			name:              "cleartext requires an encrypted transport",
			extraConfig:       ConnectionConfig{"allow_cleartext_passwords": true},
			wantCleartext:     true,
			wantTLSConfig:     "skip-verify",
			wantAllowFallback: false,
		},
		{
			name:              "explicit tls value wins",
			extraConfig:       ConnectionConfig{"allow_cleartext_passwords": true, "tls": "true"},
			wantCleartext:     true,
			wantTLSConfig:     "true",
			wantAllowFallback: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			connectionConfig := ConnectionConfig{
				"host":     "mysql.internal",
				"port":     "3306",
				"user":     "gateway",
				"password": "secret",
				"database": "gateway_metrics",
			}
			for key, value := range testCase.extraConfig {
				connectionConfig[key] = value
			}

			dbConfig, err := buildExternalDBConfig(schema.DatabaseConnectorKindMySQL, connectionConfig)
			require.NoError(t, err)
			parsed, err := mysqldriver.ParseDSN(dbConfig.DSN)
			require.NoError(t, err)

			assert.Equal(t, testCase.wantCleartext, parsed.AllowCleartextPasswords)
			assert.Equal(t, testCase.wantTLSConfig, parsed.TLSConfig)
			assert.Equal(t, testCase.wantAllowFallback, parsed.AllowFallbackToPlaintext)

			// 探測用的管理連線必須沿用同一組解析。
			adminConfig := newMySQLDriverConfig(connectionConfig, "")
			assert.Equal(t, parsed.AllowCleartextPasswords, adminConfig.AllowCleartextPasswords)
			assert.Equal(t, parsed.TLSConfig, adminConfig.TLSConfig)
			assert.Equal(t, parsed.AllowFallbackToPlaintext, adminConfig.AllowFallbackToPlaintext)
		})
	}
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

func TestIsMySQLMissingDatabaseError_DetectsError1049(t *testing.T) {
	t.Parallel()

	mysqlErr := &mysqldriver.MySQLError{
		Number:  1049,
		Message: "Unknown database 'gateway_metrics'",
	}
	assert.True(t, isMySQLMissingDatabaseError(mysqlErr))

	strErr := mysqldriver.MySQLError{Number: 1049, Message: "Unknown database 'test'"}
	assert.True(t, isMySQLMissingDatabaseError(&strErr))

	otherErr := &mysqldriver.MySQLError{
		Number:  1045,
		Message: "Access denied for user 'root'@'localhost'",
	}
	assert.False(t, isMySQLMissingDatabaseError(otherErr))
}

func TestProbeConnector_MySQLAutoCreatesDatabaseOnMissingError(t *testing.T) {
	// 覆寫 openExternalDBManagerFunc 與 ensureMySQLDatabaseIfMissingFunc
	originalOpen := openExternalDBManagerFunc
	originalEnsure := ensureMySQLDatabaseIfMissingFunc
	defer func() {
		openExternalDBManagerFunc = originalOpen
		ensureMySQLDatabaseIfMissingFunc = originalEnsure
	}()

	attempts := 0
	ensureCalled := false

	openExternalDBManagerFunc = func(kind schema.DatabaseConnectorKind, config ConnectionConfig) (*datalinkbase.DBManager, error) {
		attempts++
		if attempts == 1 {
			return nil, &mysqldriver.MySQLError{
				Number:  1049,
				Message: "Unknown database 'gateway_metrics'",
			}
		}
		// 第二次成功（模擬建庫後連接成功）
		mgr := datalinkbase.NewDBManager(datalinkbase.DBConfig{
			Type: datalinkbase.DBTypeSQLite,
			DSN:  ":memory:",
		})
		require.NoError(t, mgr.Connect())
		return mgr, nil
	}

	ensureMySQLDatabaseIfMissingFunc = func(ctx context.Context, config ConnectionConfig, connectErr error) error {
		ensureCalled = true
		assert.True(t, isMySQLMissingDatabaseError(connectErr))
		return nil
	}

	status, checkedAt, errMsg := probeConnector(
		context.Background(),
		schema.DatabaseConnectorKindMySQL,
		ConnectionConfig{
			"host":     "127.0.0.1",
			"port":     "3306",
			"user":     "root",
			"password": "secret",
			"database": "gateway_metrics",
		},
	)

	assert.True(t, ensureCalled)
	assert.Equal(t, 2, attempts)
	assert.Equal(t, schema.DatabaseConnectorStatusReady, status)
	assert.Empty(t, errMsg)
	assert.NotNil(t, checkedAt)
}

func TestOpenExternalDBManagerWithMySQLDatabaseEnsureRetriesOnce(t *testing.T) {
	originalOpen := openExternalDBManagerFunc
	originalEnsure := ensureMySQLDatabaseIfMissingFunc
	t.Cleanup(func() {
		openExternalDBManagerFunc = originalOpen
		ensureMySQLDatabaseIfMissingFunc = originalEnsure
	})

	openAttempts := 0
	ensureCalls := 0
	openExternalDBManagerFunc = func(kind schema.DatabaseConnectorKind, config ConnectionConfig) (*datalinkbase.DBManager, error) {
		openAttempts++
		if openAttempts == 1 {
			return nil, &mysqldriver.MySQLError{Number: 1049, Message: "Unknown database 'gateway_metrics'"}
		}
		manager := datalinkbase.NewDBManager(datalinkbase.DBConfig{Type: datalinkbase.DBTypeSQLite, DSN: ":memory:"})
		require.NoError(t, manager.Connect())
		return manager, nil
	}
	ensureMySQLDatabaseIfMissingFunc = func(ctx context.Context, config ConnectionConfig, connectErr error) error {
		ensureCalls++
		return nil
	}

	manager, err := openExternalDBManagerWithMySQLDatabaseEnsure(
		context.Background(),
		schema.DatabaseConnectorKindMySQL,
		ConnectionConfig{"user": "root", "database": "gateway_metrics"},
	)
	require.NoError(t, err)
	require.NotNil(t, manager)
	t.Cleanup(func() { require.NoError(t, manager.Close()) })
	assert.Equal(t, 2, openAttempts)
	assert.Equal(t, 1, ensureCalls)
}

// 非「資料庫不存在」的失敗不得觸發自動建庫，也不得重複嘗試連線。
func TestProbeConnector_MySQLDoesNotAttemptDatabaseCreationForOtherErrors(t *testing.T) {
	originalOpen := openExternalDBManagerFunc
	originalEnsure := ensureMySQLDatabaseIfMissingFunc
	t.Cleanup(func() {
		openExternalDBManagerFunc = originalOpen
		ensureMySQLDatabaseIfMissingFunc = originalEnsure
	})

	openAttempts := 0
	ensureCalls := 0
	openExternalDBManagerFunc = func(kind schema.DatabaseConnectorKind, config ConnectionConfig) (*datalinkbase.DBManager, error) {
		openAttempts++
		return nil, &mysqldriver.MySQLError{Number: 1045, Message: "Access denied for user 'gateway'"}
	}
	ensureMySQLDatabaseIfMissingFunc = func(ctx context.Context, config ConnectionConfig, connectErr error) error {
		ensureCalls++
		return connectErr
	}

	status, checkedAt, errMsg := probeConnector(
		context.Background(),
		schema.DatabaseConnectorKindMySQL,
		ConnectionConfig{
			"host":     "127.0.0.1",
			"port":     "3306",
			"user":     "gateway",
			"password": "secret",
			"database": "gateway_metrics",
		},
	)

	assert.Equal(t, 0, ensureCalls, "非 1049 的失敗不得觸發自動建庫")
	assert.Equal(t, 1, openAttempts, "非 1049 的失敗不得重試連線")
	assert.NotEqual(t, schema.DatabaseConnectorStatusReady, status)
	assert.NotEmpty(t, errMsg)
	assert.NotNil(t, checkedAt)
}

// 資料庫不存在時，一次探測只嘗試建庫一次並重試連線一次。
func TestProbeConnector_MySQLAttemptsDatabaseCreationOncePerProbe(t *testing.T) {
	originalOpen := openExternalDBManagerFunc
	originalEnsure := ensureMySQLDatabaseIfMissingFunc
	t.Cleanup(func() {
		openExternalDBManagerFunc = originalOpen
		ensureMySQLDatabaseIfMissingFunc = originalEnsure
	})

	openAttempts := 0
	ensureCalls := 0
	openExternalDBManagerFunc = func(kind schema.DatabaseConnectorKind, config ConnectionConfig) (*datalinkbase.DBManager, error) {
		openAttempts++
		if openAttempts == 1 {
			return nil, &mysqldriver.MySQLError{Number: 1049, Message: "Unknown database 'gateway_metrics'"}
		}
		manager := datalinkbase.NewDBManager(datalinkbase.DBConfig{Type: datalinkbase.DBTypeSQLite, DSN: ":memory:"})
		require.NoError(t, manager.Connect())
		return manager, nil
	}
	ensureMySQLDatabaseIfMissingFunc = func(ctx context.Context, config ConnectionConfig, connectErr error) error {
		ensureCalls++
		return nil
	}

	status, _, errMsg := probeConnector(
		context.Background(),
		schema.DatabaseConnectorKindMySQL,
		ConnectionConfig{
			"host":     "127.0.0.1",
			"port":     "3306",
			"user":     "root",
			"password": "secret",
			"database": "gateway_metrics",
		},
	)

	assert.Equal(t, 1, ensureCalls)
	assert.Equal(t, 2, openAttempts)
	assert.Equal(t, schema.DatabaseConnectorStatusReady, status)
	assert.Empty(t, errMsg)
}
