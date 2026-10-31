package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/api/handlers"
	appconfig "go-gateway/internal/config"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
)

func TestProductionWordFormatHTTPOverrideRestartAndReset(t *testing.T) {
	for _, protocol := range []schema.ProtocolType{schema.ProtocolModbusTCP, schema.ProtocolMC3E, schema.ProtocolFatekFBs} {
		t.Run(string(protocol), func(t *testing.T) {
			f := newWordProductionFixture(t, protocol)
			require.NoError(t, f.services.modbusShare.ApplySettings(t.Context(), modbusshare.Settings{Enabled: true, BindAddress: "127.0.0.1", Port: 15020, SlaveID: 1, CapacityRegisters: 64}))
			rule, p, _ := f.createRule(t, protocol, 0, schema.DataTypeUint32, "ABCD")
			require.NoError(t, f.services.runtime.Start(t.Context()))
			router := wordProductionRouter(f.services, f.db)
			record, err := f.services.workspace.GetOrCreate(t.Context())
			require.NoError(t, err)
			// Exercise the exact main startup hydration code with new services,
			// not a test-issued Ready token or a direct runtime reinstallation.
			reconcile := func() {
				require.NoError(t, f.services.runtime.Stop(t.Context()))
				require.NoError(t, f.manager.CloseAll())
				manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
				t.Cleanup(func() { require.NoError(t, manager.CloseAll()) })
				services, e := wireGatewayServices(f.db, manager)
				require.NoError(t, e)
				t.Cleanup(func() { require.NoError(t, services.runtime.Stop(context.Background())) })
				revisions := modbusshare.NewSQLWorkspaceRevisionStore(f.db)
				reconciler := modbusshare.NewReconciler(services.modbusShare, revisions)
				hydrateStartupShare(t.Context(), services, revisions, reconciler)
				state, e := services.modbusShare.CheckHydration(t.Context())
				require.NoError(t, e)
				require.True(t, state.Readiness)
				require.NotEmpty(t, state.ReadinessToken)
				require.NoError(t, services.runtime.Start(t.Context()))
				f.services, f.manager = services, manager
				router = wordProductionRouter(services, f.db)
			}
			reconcile()
			base := appconfig.Get().API.BasePath + "/datalink/points/" + p.ID + "/format"
			read := func() *point.FormatReadResult {
				response := httptest.NewRecorder()
				router.ServeHTTP(response, httptest.NewRequest("GET", base+"?workspace_id="+record.ID+"&device_id="+p.DeviceID+"&rule_id="+rule.ID, nil))
				require.Equal(t, 200, response.Code, response.Body.String())
				var body handlers.PointFormatReadResponse
				require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
				return body.Data
			}
			save := func(mode schema.PointFormatMode, format, id string) {
				initial := read()
				require.True(t, initial.Editable, "reconciled source permits a fresh edit")
				req := point.FormatMutationRequest{FormatScope: point.FormatScope{WorkspaceID: record.ID, DeviceID: p.DeviceID, PointID: p.ID, RuleID: rule.ID, ExpectedWorkspaceRevision: initial.WorkspaceRevision, ExpectedPointRevision: &initial.PointRevision, ExpectedRuleRevision: initial.RuleRevision, RequestID: id}, FormatMode: mode, ConfiguredFormat: format}
				raw, e := json.Marshal(req)
				require.NoError(t, e)
				request := httptest.NewRequest(http.MethodPut, base, bytes.NewReader(raw))
				request.Header.Set("Content-Type", "application/json")
				response := httptest.NewRecorder()
				router.ServeHTTP(response, request)
				require.Equal(t, 200, response.Code, response.Body.String())
				require.Contains(t, response.Body.String(), `"runtime_apply_status":"pending"`)
				require.False(t, read().Editable, "durable dirty released is not live submit readiness")
				require.Empty(t, f.services.scheduler.PollNowContext(t.Context(), []string{p.ID}), "saved is invalidated, not applied")
			}
			save(schema.PointFormatOverride, "CDAB", "override-1")
			f.plc.set(0, []uint16{0x3344, 0x1122})
			reconcile()
			override := f.services.scheduler.PollNowContext(t.Context(), []string{p.ID})
			require.Len(t, override, 1)
			require.Equal(t, uint32(0x11223344), override[0].Value)
			save(schema.PointFormatInherit, "", "reset-2")
			f.plc.set(0, []uint16{0x1122, 0x3344})
			reconcile()
			reset := f.services.scheduler.PollNowContext(t.Context(), []string{p.ID})
			require.Len(t, reset, 1)
			require.Equal(t, uint32(0x11223344), reset[0].Value)
			require.NotEqual(t, override[0].ConfigFingerprint, reset[0].ConfigFingerprint)
			require.Equal(t, schema.PointFormatInherit, read().Resolution.Mode)
		})
	}
}

func TestProductionWordFormatRestartDoesNotRecoverUnknownDirty(t *testing.T) {
	f := newWordProductionFixture(t, schema.ProtocolModbusTCP)
	require.NoError(t, f.services.modbusShare.ApplySettings(t.Context(), modbusshare.Settings{Enabled: true, BindAddress: "127.0.0.1", Port: 15020, SlaveID: 1, CapacityRegisters: 64}))
	f.createRule(t, schema.ProtocolModbusTCP, 0, schema.DataTypeUint32, "ABCD")
	record, err := f.services.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	revisions := modbusshare.NewSQLWorkspaceRevisionStore(f.db)
	require.NoError(t, revisions.MarkDirty(t.Context(), record.ID))
	hydrateStartupShare(t.Context(), f.services, revisions, modbusshare.NewReconciler(f.services.modbusShare, revisions))
	state, err := f.services.modbusShare.CheckHydration(t.Context())
	require.NoError(t, err)
	require.False(t, state.Readiness)
	require.Equal(t, modbusshare.HydrationStateFailed, state.State)
	_, dirty, err := revisions.GetRevision(t.Context(), record.ID)
	require.NoError(t, err)
	require.True(t, dirty, "startup never clears unknown dirty")
}
