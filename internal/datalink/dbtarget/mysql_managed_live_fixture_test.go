package dbtarget

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

func mysqlManagedLiveInspector(t *testing.T, cfg mysqlLongFixtureConfig) (*ManagedTableInspector, *MySQLRecordingCapability) {
	t.Helper()
	pem, err := os.ReadFile(cfg.CAPath)
	require.NoError(t, err)
	pool := x509.NewCertPool()
	require.True(t, pool.AppendCertsFromPEM(pem))
	key := "managed_live_" + randomHex(t, 6)
	require.NoError(t, mysqldriver.RegisterTLSConfig(key, &tls.Config{RootCAs: pool, ServerName: cfg.Host, MinVersion: tls.VersionTLS12}))
	t.Cleanup(func() { mysqldriver.DeregisterTLSConfig(key) })
	dsn := mysqldriver.NewConfig()
	dsn.User, dsn.Passwd, dsn.DBName, dsn.Net = cfg.User, cfg.Password, cfg.Database, "tcp"
	dsn.Addr, dsn.TLSConfig = net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)), key
	dsn.ParseTime, dsn.Loc = true, time.UTC
	dsn.Timeout, dsn.ReadTimeout, dsn.WriteTimeout = 5*time.Second, 10*time.Second, 10*time.Second
	raw, err := json.Marshal(map[string]string{"dsn": dsn.FormatDSN()})
	require.NoError(t, err)
	connector := &schema.DatabaseConnector{ID: "live-" + randomHex(t, 6), IdentityRevision: "rev-1", Kind: schema.DatabaseConnectorKindMySQL, Enabled: true, ConnectionConfig: string(raw)}
	service := NewConnectorService(&mysqlManagedTestRepo{connector: connector})
	inspector := NewManagedTableInspector(service, nil)
	capability, err := inspector.ProbeMySQLRecordingAtRevision(t.Context(), connector.ID, connector.IdentityRevision)
	require.NoError(t, err)
	require.True(t, capability.Connected)
	require.Equal(t, "verified", capability.MetadataStatus)
	require.True(t, capability.TransactionalEngine)
	return inspector, capability
}

func mysqlManagedLiveState(ctx context.Context, inspector *ManagedTableInspector, capability *MySQLRecordingCapability, table string) (recordingplan.TargetTableInspection, error) {
	result, err := inspector.InspectMySQLManagedTableAtRevision(ctx, capability.ConnectorID, capability.ConnectorRevision, capability.Database, capability.CapabilityRevision, table)
	if err != nil {
		return recordingplan.TargetTableInspection{}, err
	}
	state := recordingplan.TargetTableInspection{Status: string(result.Inspection.Status), ColumnTypes: map[string]string{}, ColumnNullable: map[string]bool{}, ColumnPrimaryKey: map[string]bool{}, IndexesInspected: result.Inspection.IndexesInspected}
	for _, column := range result.Inspection.Columns {
		state.Columns = append(state.Columns, column.Name)
		state.ColumnTypes[column.Name], state.ColumnNullable[column.Name], state.ColumnPrimaryKey[column.Name] = column.DataType, column.Nullable, column.PrimaryKey
	}
	for _, index := range result.Indexes {
		state.Indexes = append(state.Indexes, recordingplan.TargetIndexInspection{Name: index.Name, Columns: index.Columns, Unique: index.Unique})
	}
	e := result.Evidence
	state.MySQL = &recordingplan.TargetMySQLTableEvidence{
		ServerFamily: e.ServerFamily, ServerVersion: e.ServerVersion, CapabilityRevision: e.CapabilityRevision, Database: e.Database,
		Engine: e.Engine, Charset: e.Charset, Collation: e.Collation, TableType: e.TableType, PrimaryKey: e.PrimaryKey,
		ColumnCollations: e.ColumnCollations, ColumnDefaults: e.ColumnDefaults, JSONColumns: e.JSONColumns,
	}
	return state, nil
}

func mysqlManagedLocalLedger(t *testing.T, path string, initialize bool) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", path+"?_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)")
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	if initialize {
		for _, name := range []string{"019_recording_plans_sqlite.up.sql", "022_schema_preview_scope_sqlite.up.sql", "023_schema_operations_sqlite.up.sql", "026_operation_test_write_sqlite.up.sql"} {
			ddl, err := os.ReadFile(filepath.Join("..", "schema", "migrations", name))
			require.NoError(t, err)
			_, err = db.ExecContext(t.Context(), string(ddl))
			require.NoError(t, err)
		}
	}
	return db
}
