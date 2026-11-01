package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type mixedScopedProductionFixture struct {
	scheduledProductionFixture
	group         *workspace.WriteGroup
	request       workspace.RecordingStartRequest
	excludedPoint string
	targetFile    string
}

// Real production services, native loopback Modbus, disposable configuration and
// destination SQLite. No fake activator, samples or clock are used.
func newMixedScopedProductionFixture(t *testing.T, prepare bool) mixedScopedProductionFixture {
	t.Helper()
	plc, config := wordWireServer(t, schema.ProtocolModbusTCP)
	plc.set(0, []uint16{22, 12, 1000, 99})
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "mixed-start-config.db"))
	require.NoError(t, err)
	datalink.ApplySQLitePoolDefaults(db)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	t.Cleanup(func() { require.NoError(t, manager.CloseAll()) })
	capture := &acquisitionWiringCapture{groupFixtureHooks: newGroupFixture(), samples: make(chan measurement.SampleEnvelope, 256), observations: make(chan collector.AcquisitionObservation, 256)}
	s, err := wireGatewayServicesContext(t.Context(), db, manager, nil, capture)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, s.runtime.Stop(context.Background())) })
	for _, id := range []string{"mixed-selected", "mixed-selected-second", "mixed-unselected"} {
		_, err = db.ExecContext(t.Context(), `INSERT INTO devices(id,name,description,protocol,status,connection_config,readiness_status) VALUES(?,?,'','modbus_tcp','draft',?,'{}')`, id, id, config)
		require.NoError(t, err)
		_, err = s.workspace.AttachDevice(t.Context(), id)
		require.NoError(t, err)
		result, err := s.device.TestConnectionWithResult(t.Context(), id)
		require.NoError(t, err)
		require.True(t, result.CanActivate)
	}
	router := wordProductionRouter(s, db)
	revisions := modbusshare.NewSQLWorkspaceRevisionStore(db)
	hydrateStartupShare(t.Context(), s, revisions, modbusshare.NewReconciler(s.modbusShare, revisions))
	f := mixedScopedProductionFixture{scheduledProductionFixture: scheduledProductionFixture{s, db, manager, router, capture}, targetFile: filepath.Join(t.TempDir(), "mixed-target.db")}
	var links []*schema.SourceRuleLink
	var secondLinks []*schema.SourceRuleLink
	for _, id := range []string{"mixed-selected", "mixed-selected-second", "mixed-unselected"} {
		rule, err := s.sourceRule.Create(t.Context(), sourcerule.CreateRuleRequest{DeviceID: id, StartAddress: "40001", Count: 4, DataType: schema.DataTypeUint16, NamingPrefix: id, Enabled: true})
		require.NoError(t, err)
		applyAcquisitionTags(t, s, rule)
		switch id {
		case "mixed-selected":
			links, err = s.sourceRule.ListLinks(t.Context(), rule.ID)
			require.NoError(t, err)
		case "mixed-selected-second":
			secondLinks, err = s.sourceRule.ListLinks(t.Context(), rule.ID)
			require.NoError(t, err)
		}
	}
	require.Len(t, links, 4)
	require.Len(t, secondLinks, 4)
	links[2] = secondLinks[2]
	record, err := s.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	members := make([]workspace.WriteGroupMember, 0, 3)
	for i, link := range links {
		require.NotNil(t, link.TagID)
		if i == 3 {
			f.excludedPoint = link.PointID
			continue
		}
		method := workspace.WriteGroupMemberCalculation("mean")
		cadence := int64(15)
		member := workspace.WriteGroupMember{DeviceID: "mixed-selected", PointID: link.PointID, TagID: *link.TagID, ValueRole: "gauge", Calculation: method, Cadence: &workspace.WriteGroupMemberCadence{Interval: cadence, Unit: "second", Alignment: "utc_epoch"}}
		if i == 2 {
			id := "mixed-start-counter"
			member.DeviceID = "mixed-selected-second"
			member.MeasurementID = &id
			member.ValueRole = ""
			member.Calculation = "delta"
			member.Cadence.Interval = 60
			_, err = db.ExecContext(t.Context(), `INSERT INTO measurement_definitions(id,workspace_id,device_id,point_id,tag_id,equipment_id,definition_revision,source_binding_revision,series_epoch,name,quantity,unit,semantic_kind,numeric_encoding) VALUES(?,?,'mixed-selected-second',?,?,'mixed-selected-second','def-1','bind-1','epoch-1','Energy','energy','kWh','counter','uint16')`, id, record.ID, link.PointID, *link.TagID)
			require.NoError(t, err)
		}
		members = append(members, member)
	}
	cfg, err := json.Marshal(map[string]string{"dsn": f.targetFile})
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), `INSERT INTO database_connectors(id,name,kind,connection_config,identity_revision,status,enabled) VALUES('mixed-target','Target','sqlite',?,'target-v1','ready',1)`, string(cfg))
	require.NoError(t, err)
	created, err := s.writeGroups.Create(t.Context(), workspace.WriteGroupMutation{WorkspaceID: record.ID, ExpectedWorkspaceRevision: record.DatabaseSetupRevision, ExpectedConnectorRevision: "target-v1", Group: &workspace.WriteGroup{
		WorkspaceID: record.ID, RecordID: uuid.NewString(), Name: "Independent mixed", LayoutVersion: "managed_long_v1", StorageMode: "plain", RecordingProfile: &workspace.WriteGroupRecordingProfile{Kind: "mixed", ResultShape: "member-calculation-v1"}, Members: members, RowPolicy: workspace.WriteGroupRowPolicy{IncompletePolicy: "partial"}, Destination: workspace.WriteGroupDestination{ConnectorID: "mixed-target", ConnectorRevision: "target-v1", TableSchema: "main", StorageStrategy: "managed"},
	}})
	require.NoError(t, err)
	if prepare {
		body := map[string]any{"workspace_id": record.ID, "expected_workspace_revision": created.WorkspaceRevision, "expected_group_revision": created.Group.Revision, "expected_connector_revision": "target-v1"}
		path := testWriteBase + "/write-groups/" + created.Group.ID
		status, response := httpJSON(t, router, http.MethodPost, path+"/schema-preview", body)
		require.Equal(t, http.StatusOK, status, "%+v", response)
		token := response["data"].(map[string]any)
		body["token"], body["operation_id"] = token["token"], token["operation_id"]
		status, response = httpJSON(t, router, http.MethodPost, path+"/schema-apply", body)
		require.Equal(t, http.StatusOK, status, "%+v", response)
		require.Equal(t, "succeeded", response["data"].(map[string]any)["status"])
	}
	current, err := s.writeGroups.Get(t.Context(), created.Group.ID)
	require.NoError(t, err)
	f.group = current.Group
	record, err = s.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	f.request = workspace.RecordingStartRequest{RequestID: "independent-mixed-start", WorkspaceID: record.ID, ExpectedWorkspaceRevision: record.DatabaseSetupRevision, DeviceIDs: []string{"mixed-selected-second", "mixed-selected"}, Groups: []workspace.RecordingStartGroupIntent{{GroupID: f.group.ID, ExpectedGroupRevision: f.group.Revision, ExpectedConnectorRevision: "target-v1"}}}
	require.NoError(t, s.runtime.Start(t.Context()))
	require.Empty(t, s.scheduler.CollectionStatus(), "draft devices must not collect before scoped start")
	hydration, err := s.modbusShare.CheckHydration(t.Context())
	require.NoError(t, err)
	require.True(t, hydration.Readiness, "%+v", hydration)
	f.request.WorkspaceRevision = hydration.WorkspaceRevision
	f.request.SettingsRevision = hydration.SettingsRevision
	f.request.ReadinessToken = hydration.ReadinessToken
	return f
}
