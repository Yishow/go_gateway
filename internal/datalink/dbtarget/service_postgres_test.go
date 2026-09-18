package dbtarget

import (
	"context"
	"path/filepath"
	"testing"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSavedConnectorIdentityPostgresPasswordBytes(t *testing.T) {
	for _, password := range []string{"  fixture password  ", "   ", "fixture'\\password\t"} {
		config, err := buildExternalDBConfig(schema.DatabaseConnectorKindPostgres, ConnectionConfig{
			"host": "127.0.0.1", "user": "fixture_writer", "database": "fixture_metrics",
			"password": password, "sslmode": "disable",
		})
		require.NoError(t, err)
		parsed, err := pgconn.ParseConfig(config.DSN)
		require.NoError(t, err)
		require.Equal(t, password, parsed.Password)
		_, err = pq.NewConnector(config.DSN)
		require.NoError(t, err)
	}
}

func TestProbeConnector_PostgresCreatesMissingDatabaseAndRetries(t *testing.T) {
	originalOpen := openExternalDBManagerFunc
	originalEnsure := ensurePostgresDatabaseIfMissingFunc
	t.Cleanup(func() {
		openExternalDBManagerFunc = originalOpen
		ensurePostgresDatabaseIfMissingFunc = originalEnsure
	})

	tempDSN := filepath.Join(t.TempDir(), "probe-retry.db")
	openAttempts := 0
	ensureCalls := 0
	openExternalDBManagerFunc = func(kind schema.DatabaseConnectorKind, config ConnectionConfig) (*datalinkbase.DBManager, error) {
		openAttempts++
		if openAttempts == 1 {
			return nil, &pq.Error{Code: "3D000", Message: `database "gateway_metrics" does not exist`}
		}

		manager := datalinkbase.NewDBManager(datalinkbase.DBConfig{
			Type:         datalinkbase.DBTypeSQLite,
			DSN:          tempDSN,
			MaxOpenConns: 1,
			MaxIdleConns: 1,
		})
		if err := manager.Connect(); err != nil {
			return nil, err
		}
		return manager, nil
	}
	ensurePostgresDatabaseIfMissingFunc = func(ctx context.Context, config ConnectionConfig, connectErr error) error {
		ensureCalls++
		return nil
	}

	status, checkedAt, message := probeConnector(
		context.Background(),
		schema.DatabaseConnectorKindPostgres,
		ConnectionConfig{
			"host":     "127.0.0.1",
			"port":     "5432",
			"user":     "postgres",
			"password": "secret",
			"database": "gateway_metrics",
			"sslmode":  "disable",
		},
	)

	require.NotNil(t, checkedAt)
	assert.Equal(t, schema.DatabaseConnectorStatusReady, status)
	assert.Empty(t, message)
	assert.Equal(t, 2, openAttempts)
	assert.Equal(t, 1, ensureCalls)
}
