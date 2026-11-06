package handlers

import (
	"net/http"
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func mixedHandlerRuntimeProjection(t *testing.T, f *mixedRuleHandlerFixture) *workspace.RuntimeProjection {
	t.Helper()
	service := workspace.NewService(workspace.NewSQLRepository(f.db)).WithRuntimeProjectionServices(
		device.NewSQLRepository(f.db), f.rules, f.points, mapping.NewSQLRepository(f.db), tag.NewSQLRepository(f.db), nil, nil, point.NewMemoryPollingGroupRepository(),
	)
	projection, err := service.RuntimeProjection(t.Context())
	require.NoError(t, err)
	return projection
}

func TestWorkspaceMixedRuleAcceptsFatekDoubleWordComponent(t *testing.T) {
	f := newMixedRuleHandlerFixture(t)
	_, err := f.db.ExecContext(t.Context(), `UPDATE devices SET protocol=?,connection_config=? WHERE id='mc-dev'`, schema.ProtocolFatekFBs, `{"host":"127.0.0.1","port":500,"acquisition":{"version":1,"profile":"fatek-fbs-ascii-v1"}}`)
	require.NoError(t, err)
	body := `{"id":"mixed-fatek","device_id":"mc-dev","naming_prefix":"F","enabled":true,"layout_mode":"mixed","mixed_items":[{"item_id":"dr","address":"DR100","data_type":"uint32","semantic_kind":"gauge"},{"item_id":"r","address":"R5","data_type":"int16","semantic_kind":"gauge"}]}`
	code, response := f.do(t, http.MethodPost, "/workspace/source-rules", body)
	require.Equal(t, http.StatusCreated, code, response)
	links, err := f.rules.ListLinks(t.Context(), "mixed-fatek")
	require.NoError(t, err)
	require.Len(t, links, 2)
	require.ElementsMatch(t, []string{"DR100", "R5"}, []string{links[0].Address, links[1].Address})
	projection := mixedHandlerRuntimeProjection(t, f)
	graph, complete := projection.AcquisitionSourceGraph("mc-dev")
	require.True(t, complete, graph.Reason())
}

func TestWorkspaceMixedRuleRejectsOverlappingFatekAliases(t *testing.T) {
	f := newMixedRuleHandlerFixture(t)
	_, err := f.db.ExecContext(t.Context(), `UPDATE devices SET protocol=?,connection_config=? WHERE id='mc-dev'`, schema.ProtocolFatekFBs, `{"host":"127.0.0.1","port":500,"acquisition":{"version":1,"profile":"fatek-fbs-ascii-v1"}}`)
	require.NoError(t, err)
	body := `{"id":"overlap-fatek","device_id":"mc-dev","naming_prefix":"F","enabled":true,"layout_mode":"mixed","mixed_items":[{"item_id":"double","address":"DR100","data_type":"uint32","semantic_kind":"gauge"},{"item_id":"word","address":"R101","data_type":"int16","semantic_kind":"gauge"}]}`
	code, response := f.do(t, http.MethodPost, "/workspace/source-rules", body)
	require.Equal(t, http.StatusBadRequest, code, response)
	problem := response["error"].(map[string]any)
	require.Equal(t, "mixed_items", problem["field"])
	require.ElementsMatch(t, []any{"double", "word"}, problem["item_ids"])
}

func TestWorkspaceMixedRuleAcceptsQIQRHexBitAddress(t *testing.T) {
	f := newMixedRuleHandlerFixture(t)
	body := `{"id":"mixed-q","device_id":"mc-dev","naming_prefix":"Q","enabled":true,"layout_mode":"mixed","mixed_items":[{"item_id":"x","address":"X1A","data_type":"bool","semantic_kind":"state"},{"item_id":"y","address":"Y10","data_type":"bool","semantic_kind":"state"}]}`
	code, response := f.do(t, http.MethodPost, "/workspace/source-rules", body)
	require.Equal(t, http.StatusCreated, code, response)
	links, err := f.rules.ListLinks(t.Context(), "mixed-q")
	require.NoError(t, err)
	require.Len(t, links, 2)
	projection := mixedHandlerRuntimeProjection(t, f)
	graph, complete := projection.AcquisitionSourceGraph("mc-dev")
	require.True(t, complete, graph.Reason())
}

func TestWorkspaceHomogeneousMCBitRangeUsesDeviceProfileRadix(t *testing.T) {
	for _, tc := range []struct {
		name, profile, second string
	}{
		{"q", "mitsubishi-q-iqr-3e-binary-v1", "X18"},
		{"fx5", "mitsubishi-fx5-cpu-3e-binary-v1", "X20"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newMixedRuleHandlerFixture(t)
			config := `{"host":"127.0.0.1","port":5000,"pc_no":255,"acquisition":{"version":1,"profile":"` + tc.profile + `"}}`
			_, err := f.db.ExecContext(t.Context(), `UPDATE devices SET connection_config=? WHERE id='mc-dev'`, config)
			require.NoError(t, err)
			body := `{"id":"range","device_id":"mc-dev","start_address":"X17","count":2,"data_type":"bool","naming_prefix":"X","enabled":true}`
			code, response := f.do(t, http.MethodPost, "/workspace/source-rules", body)
			require.Equal(t, http.StatusCreated, code, response)
			links, err := f.rules.ListLinks(t.Context(), "range")
			require.NoError(t, err)
			require.Len(t, links, 2)
			require.ElementsMatch(t, []string{"X17", tc.second}, []string{links[0].Address, links[1].Address})
			projection := mixedHandlerRuntimeProjection(t, f)
			graph, complete := projection.AcquisitionSourceGraph("mc-dev")
			require.True(t, complete, graph.Reason())
		})
	}
}
