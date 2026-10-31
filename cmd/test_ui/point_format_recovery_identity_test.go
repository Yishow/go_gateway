package main

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	appconfig "go-gateway/internal/config"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/workspace"
)

func TestProductionFormatRecoveryMatchesSQLMemberSourceIdentity(t *testing.T) {
	f, _ := newRecoveryFaultProductionFixture(t)
	req := recoveryFaultSave(t, f)
	raw, err := json.Marshal(map[string]any{"format_recovery": req})
	require.NoError(t, err)
	request := httptest.NewRequest("POST", appconfig.Get().API.BasePath+"/datalink/modbus-share/reconcile", bytes.NewReader(raw))
	request.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	f.router.ServeHTTP(res, request)
	require.Equal(t, 200, res.Code, res.Body.String())
	evidence, err := f.services.point.ReadFormatRecovery(t.Context(), req)
	require.NoError(t, err)
	require.Equal(t, req.RuleID, evidence.Profile.IdentityView().RuleID)
	require.Len(t, evidence.TargetInputs.Mappings, 1)
	m := evidence.TargetInputs.Mappings[0]
	tg := evidence.TargetInputs.Tags[m.TagID]
	require.NotNil(t, tg)
	_, err = f.db.Exec(`INSERT INTO database_connectors(id,name,kind,connection_config,identity_revision,status,enabled) VALUES('identity-db','SQLite','sqlite','{"database":"persisted-db","schema":"persisted-schema"}','identity-db-rev','ready',1)`)
	require.NoError(t, err)
	group := &workspace.WriteGroup{
		WorkspaceID: req.WorkspaceID, Name: "recovery identity",
		Members:     []workspace.WriteGroupMember{{DeviceID: req.DeviceID, PointID: req.PointID, TagID: tg.ID, TargetColumn: "value", Required: true}},
		Destination: workspace.WriteGroupDestination{ConnectorID: "identity-db", ConnectorRevision: "identity-db-rev", Database: "line-a", TableSchema: "main", TableName: "recovered_values", StorageStrategy: workspace.WriteGroupStorageStrategyCustom},
	}
	require.NoError(t, workspace.NewSQLWriteGroupRepository(f.db).Create(t.Context(), group))
	expected, err := measurement.SourceRevision(evidence.TargetInputs.Device, evidence.Point, *tg, evidence.Profile)
	require.NoError(t, err)
	require.Equal(t, group.Members[0].SourceRevision, expected, "SQL member source resolver and recovered producer must agree")
	fingerprint, err := measurement.AcquisitionConfigFingerprint(evidence.TargetInputs.Device, evidence.Point, evidence.Profile)
	require.NoError(t, err)
	values := f.services.scheduler.PollNowContext(t.Context(), []string{req.PointID})
	require.Len(t, values, 1)
	require.Equal(t, fingerprint, values[0].ConfigFingerprint)
}
