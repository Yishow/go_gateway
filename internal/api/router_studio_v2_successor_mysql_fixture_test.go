package api

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"encoding/json"
	"net"
	"os"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"
)

// Only disposable loopback A/B databases are accepted. Each database must be
// empty before ownership is acquired; cleanup checks all owned objects are gone.
func newMySQLSuccessorDestination(t *testing.T, family, targetFamily string) (managedGroupFixture, *sql.DB, string) {
	t.Helper()
	pool := x509.NewCertPool()
	found := false
	for _, key := range []string{"MYSQL_TEST_CA_FILE", "MARIADB_TEST_CA_FILE"} {
		if path := os.Getenv(key); path != "" {
			raw, err := os.ReadFile(path)
			require.NoError(t, err)
			require.True(t, pool.AppendCertsFromPEM(raw))
			found = true
		}
	}
	if !found {
		t.Skip("NOT RUN: owned verified TLS fixture unavailable")
	}
	require.NoError(t, mysql.RegisterTLSConfig("gateway-local-verified", &tls.Config{RootCAs: pool, ServerName: "127.0.0.1", MinVersion: tls.VersionTLS12}))
	t.Cleanup(func() { mysql.DeregisterTLSConfig("gateway-local-verified") })
	f := managedGroupFixture{writeGroupRouterFixture: newWriteGroupRouterFixture(t)}
	var targetB *sql.DB
	var namespaceB string
	for i, engine := range []string{family, targetFamily} {
		key, database, port := "MYSQL_DSN", "gateway_mysql_test", "55440"
		if engine == "mariadb" {
			key, database, port = "MARIADB_DSN", "gateway_mariadb_test", "55441"
		}
		if i == 1 {
			key += "_B"
			database += "_b"
		}
		rawDSN := os.Getenv(key)
		if rawDSN == "" {
			t.Skip("NOT RUN: owned loopback MySQL-family fixture unavailable")
		}
		cfg, err := mysql.ParseDSN(rawDSN)
		require.NoError(t, err)
		require.Equal(t, "tcp", cfg.Net)
		require.Equal(t, "gateway-local-verified", cfg.TLSConfig)
		require.NotNil(t, cfg.TLS)
		require.False(t, cfg.TLS.InsecureSkipVerify)
		require.False(t, cfg.AllowFallbackToPlaintext)
		require.Equal(t, net.JoinHostPort("127.0.0.1", port), cfg.Addr)
		require.Equal(t, database, cfg.DBName)
		db, err := sql.Open("mysql", rawDSN)
		require.NoError(t, err)
		t.Cleanup(func() { require.NoError(t, db.Close()) })
		require.NoError(t, db.PingContext(t.Context()))
		var version string
		require.NoError(t, db.QueryRowContext(t.Context(), "SELECT VERSION()").Scan(&version))
		expectedVersion := "8.4.13"
		if engine == "mariadb" {
			expectedVersion = "11.4.13-MariaDB-ubu2404"
		}
		require.Equal(t, expectedVersion, version, "actual binary must match the existing verified release, never a spoofed family proof")
		var count int
		require.NoError(t, db.QueryRow(`SELECT count(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA=DATABASE()`).Scan(&count))
		require.Zero(t, count, "fixture database must start empty")
		t.Cleanup(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			rows, err := db.QueryContext(ctx, `SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_SCHEMA=DATABASE()`)
			require.NoError(t, err)
			var names []string
			for rows.Next() {
				var name string
				require.NoError(t, rows.Scan(&name))
				names = append(names, name)
			}
			require.NoError(t, rows.Err())
			require.NoError(t, rows.Close())
			for _, name := range names {
				require.Regexp(t, `^(gw_long_[0-9a-f]{16}|gw_effect_receipts)$`, name)
				_, err = db.ExecContext(ctx, "DROP TABLE `"+name+"`")
				require.NoError(t, err)
			}
			require.NoError(t, db.QueryRowContext(ctx, `SELECT count(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA=DATABASE()`).Scan(&count))
			require.Zero(t, count, "owned object cleanup")
		})
		config, err := json.Marshal(map[string]string{"dsn": rawDSN, "schema": database})
		require.NoError(t, err)
		if i == 0 {
			_, err = f.db.Exec(`UPDATE database_connectors SET kind='mysql',connection_config=? WHERE id='connector-A'`, string(config))
			f.target, f.namespace = db, database
		} else {
			_, err = f.db.Exec(`INSERT INTO database_connectors(id,name,kind,connection_config,identity_revision,status,enabled) VALUES('connector-B','Other','mysql',?,'connector-2','ready',1)`, string(config))
			targetB, namespaceB = db, database
		}
		require.NoError(t, err)
	}
	var err error
	f.record, err = f.workspace.BindDatabaseConnector(t.Context(), "connector-A")
	require.NoError(t, err)
	_, err = f.db.Exec(`UPDATE tags SET data_type='uint64' WHERE id='tag-A'`)
	require.NoError(t, err)
	targets := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(f.db))
	f.groups = workspace.NewWriteGroupService(f.workspace, workspace.NewSQLWriteGroupRepository(f.db)).WithTableInspector(dbtarget.NewReadOnlyTableInspector(targets)).WithBacklogOwnershipGuard(grouppipeline.BacklogGuard{})
	f.groups.WithManagedTableInspector(f.groups.ManagedTableInspector(targets))
	f.router = NewRouter(&DatalinkServices{Workspace: f.workspace, WriteGroups: f.groups, DBTarget: targets, RecordingPlan: recordingplan.NewService(recordingplan.NewSQLRepository(f.db))})
	return f, targetB, namespaceB
}
