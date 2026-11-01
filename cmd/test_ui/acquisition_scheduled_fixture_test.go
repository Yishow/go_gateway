package main

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	appconfig "go-gateway/internal/config"
	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"

	"github.com/stretchr/testify/require"
)

type scheduledProductionFixture struct {
	services gatewayServices
	db       *sql.DB
	manager  *connector.ConnectionManager
	router   http.Handler
	capture  *acquisitionWiringCapture
}

func newScheduledProductionFixture(t *testing.T, protocol schema.ProtocolType, config string) scheduledProductionFixture {
	t.Helper()
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "acquisition.db"))
	require.NoError(t, err)
	datalink.ApplySQLitePoolDefaults(db)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	_, err = db.ExecContext(t.Context(), `INSERT INTO devices(id,name,description,protocol,status,connection_config,readiness_status) VALUES('batch-device','PLC','',?,'draft',?,'{}')`, protocol, config)
	require.NoError(t, err)
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	t.Cleanup(func() { require.NoError(t, manager.CloseAll()) })
	capture := &acquisitionWiringCapture{groupFixtureHooks: newGroupFixture(), samples: make(chan measurement.SampleEnvelope, 256), observations: make(chan collector.AcquisitionObservation, 256)}
	services, err := wireGatewayServicesContext(t.Context(), db, manager, nil, capture)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, services.runtime.Stop(context.Background())) })
	require.NoError(t, services.device.ProbeAndActivate(t.Context(), "batch-device"))
	_, err = services.workspace.AttachDevice(t.Context(), "batch-device")
	require.NoError(t, err)
	router := wordProductionRouter(services, db)
	revisions := modbusshare.NewSQLWorkspaceRevisionStore(db)
	hydrateStartupShare(t.Context(), services, revisions, modbusshare.NewReconciler(services.modbusShare, revisions))
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodPost, appconfig.Get().API.BasePath+"/datalink/devices/batch-device/activate", http.NoBody))
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	return scheduledProductionFixture{services, db, manager, router, capture}
}

func (f scheduledProductionFixture) createSource(t *testing.T, address string, count int, dtype schema.DataType, name string) (*schema.SourceRule, []*schema.SourceRuleLink) {
	t.Helper()
	rule, err := f.services.sourceRule.Create(t.Context(), sourcerule.CreateRuleRequest{DeviceID: "batch-device", StartAddress: address, Count: count, DataType: dtype, NamingPrefix: name, Enabled: true})
	require.NoError(t, err)
	links, err := f.services.sourceRule.ListLinks(t.Context(), rule.ID)
	require.NoError(t, err)
	require.Len(t, links, count)
	return rule, links
}

func (f scheduledProductionFixture) samples(t *testing.T, count int) []measurement.SampleEnvelope {
	t.Helper()
	result := make([]measurement.SampleEnvelope, 0, count)
	deadline := time.NewTimer(7 * time.Second)
	defer deadline.Stop()
	for len(result) < count {
		select {
		case sample := <-f.capture.samples:
			result = append(result, sample)
		case <-deadline.C:
			t.Fatalf("scheduled SampleSink timeout: received %d/%d; runtime=%+v rejected=%d", len(result), count, f.services.runtime.Snapshot(), f.services.runtime.RejectedObservationCount())
		}
	}
	return result
}
