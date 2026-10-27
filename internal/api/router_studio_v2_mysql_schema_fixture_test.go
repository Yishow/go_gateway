package api

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

// This packet is shared with the dbtarget fixture tests; credentials and DSNs
// stay outside the repository and are never included in assertion messages.
type mysqlSchemaFixturePacket struct {
	Host, Database, User, Password, CAPath string
	Port                                   int
}

func newMySQLSchemaRouterFixture(t *testing.T, family string, tagTypes ...string) managedGroupFixture {
	t.Helper()
	env := "GW_CIRCUIT_MYSQL_FIXTURE_FILE"
	if family == "mariadb" {
		env = "GW_CIRCUIT_MARIADB_FIXTURE_FILE"
	}
	path := os.Getenv(env)
	if path == "" {
		t.Skip("NOT RUN: dedicated MySQL family fixture packet unavailable")
	}
	info, err := os.Stat(path)
	require.NoError(t, err)
	require.True(t, info.Mode().IsRegular())
	require.Zero(t, info.Mode().Perm()&0o077)
	raw, err := os.ReadFile(path)
	require.NoError(t, err)
	var packet mysqlSchemaFixturePacket
	require.NoError(t, json.Unmarshal(raw, &packet))
	require.Equal(t, "127.0.0.1", packet.Host)
	require.Equal(t, "gw_circuit_test", packet.Database)
	require.Greater(t, packet.Port, 0)
	require.LessOrEqual(t, packet.Port, 65535)
	require.NotEmpty(t, packet.User)
	require.True(t, packet.Password != "", "fixture password is required")
	ca, err := os.ReadFile(packet.CAPath)
	require.NoError(t, err)
	pool := x509.NewCertPool()
	require.True(t, pool.AppendCertsFromPEM(ca))
	random := make([]byte, 8)
	_, err = rand.Read(random)
	require.NoError(t, err)
	key := "api_schema_" + hex.EncodeToString(random)
	require.NoError(t, mysql.RegisterTLSConfig(key, &tls.Config{RootCAs: pool, ServerName: packet.Host, MinVersion: tls.VersionTLS12}))
	t.Cleanup(func() { mysql.DeregisterTLSConfig(key) })
	dsn := mysql.NewConfig()
	dsn.User, dsn.Passwd, dsn.DBName, dsn.Net = packet.User, packet.Password, packet.Database, "tcp"
	dsn.Addr, dsn.TLSConfig = net.JoinHostPort(packet.Host, strconv.Itoa(packet.Port)), key
	dsn.ParseTime, dsn.Loc = true, time.UTC
	dsn.Timeout, dsn.ReadTimeout, dsn.WriteTimeout = 5*time.Second, 10*time.Second, 10*time.Second
	target, err := sql.Open("mysql", dsn.FormatDSN())
	require.NoError(t, err)
	t.Cleanup(func() { _ = target.Close() })
	require.NoError(t, target.PingContext(t.Context()))
	var version string
	require.NoError(t, target.QueryRowContext(t.Context(), "SELECT VERSION()").Scan(&version))
	if family == "mariadb" {
		require.Regexp(t, `^(5\.5\.5-)?11\.4\..*MariaDB`, version)
	} else {
		require.Regexp(t, `^8\.4\.`, version)
		require.NotContains(t, version, "MariaDB")
	}
	f := newManagedDestinationFixture(t, "sqlite", tagTypes...)
	config, err := json.Marshal(map[string]string{"dsn": dsn.FormatDSN(), "schema": packet.Database})
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), `UPDATE database_connectors SET kind='mysql', connection_config=? WHERE id='connector-A'`, string(config))
	require.NoError(t, err)
	f.record, err = f.workspace.BindDatabaseConnector(t.Context(), "connector-A")
	require.NoError(t, err)
	f.namespace, f.target = packet.Database, target
	return f
}

func cleanupMySQLSchemaTables(t *testing.T, target *sql.DB, created *[]string) {
	t.Helper()
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		for _, table := range *created {
			require.Regexp(t, `^(gw_long_[0-9a-f]{16}|gw_effect_receipts)$`, table)
			_, err := target.ExecContext(ctx, "DROP TABLE `"+table+"`")
			require.NoError(t, err)
		}
	})
}
