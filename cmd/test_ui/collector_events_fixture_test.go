package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"
	"go-gateway/internal/diagnostics"
	"go-gateway/internal/virtual/memory"
	virtualmodbus "go-gateway/internal/virtual/server/modbus"

	"github.com/stretchr/testify/require"
)

type collectorEventDevice struct {
	server *virtualmodbus.Server
	port   int
	rule   *schema.SourceRule
	point  *schema.Point
}

type collectorEventFixture struct {
	services gatewayServices
	db       *sql.DB
	broker   *diagnostics.Broker
	manager  *connector.ConnectionManager
	record   *workspace.Record
	devices  []collectorEventDevice
}

func newCollectorEventFixture(t *testing.T, count int, observedRevision, periodic bool) collectorEventFixture {
	t.Helper()
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "collector-events.db"))
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	t.Cleanup(func() { require.NoError(t, manager.CloseAll()) })
	broker, err := diagnostics.New(diagnostics.Options{})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, broker.Close(context.Background())) })
	services, err := wireGatewayServicesContext(t.Context(), db, manager, broker)
	require.NoError(t, err)
	f := collectorEventFixture{services: services, db: db, broker: broker, manager: manager}
	for i := range count {
		bank := memory.NewMemoryBank(64)
		require.NoError(t, bank.WriteWord(0, 12345))
		server := virtualmodbus.NewServer(bank)
		require.NoError(t, server.Start(0))
		t.Cleanup(func() { require.NoError(t, server.Stop()) })
		config, err := json.Marshal(schema.ConnectionConfigModbusTCP{Host: "127.0.0.1", Port: server.Port(), SlaveID: 1, Timeout: 1})
		require.NoError(t, err)
		id := fmt.Sprintf("collector-device-%d", i)
		_, err = db.ExecContext(t.Context(), `INSERT INTO devices(id,name,description,protocol,status,connection_config,readiness_status) VALUES(?, 'password_fixture', '', 'modbus_tcp', 'active', ?, '{}')`, id, string(config))
		require.NoError(t, err)
		probe, err := services.device.TestConnectionWithResult(t.Context(), id)
		require.NoError(t, err)
		require.True(t, probe.CanActivate, "fixture requires successful native connect and read probe")
		record, err := services.workspace.AttachDevice(t.Context(), id)
		require.NoError(t, err)
		f.record = record
		rule, err := services.sourceRule.Create(t.Context(), sourcerule.CreateRuleRequest{ID: fmt.Sprintf("collector-rule-%d", i), DeviceID: id, StartAddress: "40001", Count: 1, DataType: schema.DataTypeInt16, NamingPrefix: "Source", Enabled: true})
		require.NoError(t, err)
		require.True(t, rule.Enabled, "fixture must exercise actual enabled acquisition")
		links, err := services.sourceRule.ListLinks(t.Context(), rule.ID)
		require.NoError(t, err)
		require.Len(t, links, 1)
		point, err := services.point.GetByID(t.Context(), links[0].PointID)
		require.NoError(t, err)
		require.True(t, point.Enabled)
		f.devices = append(f.devices, collectorEventDevice{server: server, port: server.Port(), rule: rule, point: point})
	}
	if observedRevision {
		f.record.DatabaseSetupRevision = "observed-workspace-revision"
	} else {
		f.record.DatabaseSetupRevision = ""
	}
	require.NoError(t, workspace.NewSQLRepository(db).Save(t.Context(), f.record))
	groups, err := services.pollingGroup.List(t.Context())
	require.NoError(t, err)
	for _, group := range groups {
		interval := 100
		_, err := services.pollingGroup.Update(t.Context(), group.ID, pollinggroup.UpdateRequest{IntervalMs: &interval, Enabled: &periodic})
		require.NoError(t, err)
	}
	t.Cleanup(func() { require.NoError(t, services.runtime.Stop(context.Background())) })
	return f
}

func (f collectorEventFixture) start(t *testing.T) {
	t.Helper()
	require.NoError(t, f.services.runtime.Start(t.Context()))
}
func (f collectorEventFixture) stopDevice(t *testing.T, i int) {
	t.Helper()
	require.NoError(t, f.manager.Close(f.devices[i].rule.DeviceID))
	require.NoError(t, f.devices[i].server.Stop())
}
func (f collectorEventFixture) recoverDevice(t *testing.T, i int) {
	t.Helper()
	require.NoError(t, f.devices[i].server.Start(f.devices[i].port))
	require.NoError(t, f.manager.Close(f.devices[i].rule.DeviceID))
}
func (f collectorEventFixture) poll(t *testing.T, i int) collector.CollectedValue {
	t.Helper()
	results := f.services.scheduler.PollNowContext(t.Context(), []string{f.devices[i].point.ID})
	require.Len(t, results, 1)
	return results[0]
}
func collectorEventRecords(t *testing.T, b *diagnostics.Broker) []diagnostics.Event {
	t.Helper()
	require.NoError(t, b.Flush(t.Context()))
	snapshot, err := b.Snapshot(diagnostics.Query{})
	require.NoError(t, err)
	return snapshot.Records
}
