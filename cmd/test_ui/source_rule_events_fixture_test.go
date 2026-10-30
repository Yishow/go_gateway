package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"
	"go-gateway/internal/diagnostics"
	"go-gateway/internal/virtual/memory"
	virtualmodbus "go-gateway/internal/virtual/server/modbus"
)

type sourceEventFixture struct {
	services gatewayServices
	db       *sql.DB
	broker   *diagnostics.Broker
	rule     *schema.SourceRule
	record   *workspace.Record
}

func newSourceEventFixture(t *testing.T, running bool) sourceEventFixture {
	t.Helper()
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "configuration.db"))
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	server := virtualmodbus.NewServer(memory.NewMemoryBank(64))
	require.NoError(t, server.Start(0))
	t.Cleanup(func() { require.NoError(t, server.Stop()) })
	config, err := json.Marshal(schema.ConnectionConfigModbusTCP{Host: "127.0.0.1", Port: server.Port(), SlaveID: 1, Timeout: 1})
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO devices(id,name,description,protocol,status,connection_config,readiness_status) VALUES('source-device','Simulator','','modbus_tcp','active',?,'{}')`, string(config))
	require.NoError(t, err)
	broker, err := diagnostics.New(diagnostics.Options{})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, broker.Close(context.Background())) })
	services, err := wireGatewayServicesContext(t.Context(), db, connector.NewConnectionManager(connector.DefaultConnectionManagerConfig()), broker)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, services.runtime.Stop(context.Background())) })
	record, err := services.workspace.AttachDevice(t.Context(), "source-device")
	require.NoError(t, err)
	record.DatabaseSetupRevision = "observed-workspace-revision"
	require.NoError(t, workspace.NewSQLRepository(db).Save(t.Context(), record))
	rule, err := services.sourceRule.Create(t.Context(), sourcerule.CreateRuleRequest{ID: "source-rule", DeviceID: "source-device", StartAddress: "40001", Count: 1, DataType: schema.DataTypeInt16, NamingPrefix: "Source", Enabled: true})
	require.NoError(t, err)
	if running {
		require.NoError(t, services.runtime.Start(t.Context()))
	}
	return sourceEventFixture{services, db, broker, rule, record}
}

func sourceEventRecords(t *testing.T, b *diagnostics.Broker) []diagnostics.Event {
	t.Helper()
	require.NoError(t, b.Flush(t.Context()))
	snapshot, err := b.Snapshot(diagnostics.Query{})
	require.NoError(t, err)
	return snapshot.Records
}
