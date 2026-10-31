package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func mixedScaledSaveFixture(t *testing.T, method string) (fixture managedGroupFixture, payload map[string]any) {
	t.Helper()
	f := newManagedDestinationFixture(t, "sqlite", "uint64")
	_, err := f.db.ExecContext(t.Context(), `UPDATE points SET data_type='uint64' WHERE id='point-A'`)
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), `INSERT INTO measurement_definitions(id,workspace_id,device_id,point_id,tag_id,equipment_id,definition_revision,source_binding_revision,series_epoch,name,quantity,unit,semantic_kind,numeric_encoding,exact_recording_scale) VALUES('mixed-meter',?,'device-A','point-A','tag-A','device-A','def-1','bind-1','epoch-1','Energy','energy','kWh','counter','uint64','0.1')`, f.record.ID)
	require.NoError(t, err)
	request := f.createRequest()
	g := request["group"].(map[string]any)
	g["record_id"] = "64a57435-9b81-45d7-8d80-86000de6eaa1"
	g["layout_version"] = "managed_long_v1"
	g["storage_mode"] = "plain"
	g["recording_profile"] = map[string]any{"kind": "mixed", "result_shape": "member-calculation-v1"}
	g["row_policy"] = map[string]any{"incomplete_policy": "partial"}
	dest := g["destination"].(map[string]any)
	dest["storage_strategy"] = "managed"
	dest["table_schema"] = "main"
	dest["table_name"] = ""
	m := g["members"].([]map[string]any)[0]
	m["calculation"] = method
	m["measurement_id"] = "mixed-meter"
	m["cadence"] = map[string]any{"interval": 60, "unit": "second", "alignment": "utc_epoch"}
	delete(m, "target_column")
	return f, request
}

func TestMixedAcceptanceScaledDeltaPublicSave(t *testing.T) {
	f, request := mixedScaledSaveFixture(t, "delta")
	response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, request)

	require.Equal(t, http.StatusCreated, response.Code, response.Body.String())
	capabilities := decodeJSONBody(t, response)["data"].(map[string]any)["calculation_capabilities"].(map[string]any)
	require.Equal(t, true, capabilities["activation_supported"], "supported activation still requires this record's fresh schema/readiness proof")
	require.Equal(t, false, capabilities["editor_supported"])
	require.Equal(t, "record_readiness_required", capabilities["reason"])
	var wideOnly bool
	for _, raw := range capabilities["unsupported_combinations"].([]any) {
		restriction := raw.(map[string]any)
		if restriction["method"] == "delta" && restriction["phase"] == "draft" {
			require.Equal(t, "managed_wide_profile_v1", restriction["layout"])
			wideOnly = true
		}
	}
	require.True(t, wideOnly)
	_, g := groupSaveData(t, decodeJSONBody(t, response))
	m := g["members"].([]any)[0].(map[string]any)
	require.Equal(t, "uint64", m["source_value_type"])
	require.Equal(t, "decimal", m["result_value_type"])
	proof := m["exact_recording_source"].(map[string]any)
	require.Equal(t, "0.1", proof["scale"])
	require.Equal(t, "def-1", proof["definition_revision"])
	require.NotContains(t, g["row_policy"], "interval_seconds")
	get := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+g["id"].(string), nil)
	require.Equal(t, http.StatusOK, get.Code)
	_, read := groupSaveData(t, decodeJSONBody(t, get))
	require.Equal(t, g, read)
	var pointType, tagType string
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT data_type FROM points WHERE id='point-A'`).Scan(&pointType))
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT data_type FROM tags WHERE id='tag-A'`).Scan(&tagType))
	require.Equal(t, "uint64", pointType)
	require.Equal(t, "uint64", tagType)

}

func TestMixedAcceptanceScaledLatestAndMeanRemainRejected(t *testing.T) {
	for _, method := range []string{"latest", "mean"} {
		t.Run(method, func(t *testing.T) {
			f, request := mixedScaledSaveFixture(t, method)
			response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, request)
			require.Equal(t, http.StatusUnprocessableEntity, response.Code, response.Body.String())
			var n int
			require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT count(*) FROM write_groups`).Scan(&n))
			require.Zero(t, n)
		})
	}
}

func TestMixedAcceptanceScaledDraftRefusesForeignOrUnsupportedSource(t *testing.T) {
	for _, mode := range []string{"wide", "foreign measurement", "unconfirmed"} {
		t.Run(mode, func(t *testing.T) {
			f, request := mixedScaledSaveFixture(t, "delta")
			g := request["group"].(map[string]any)
			switch mode {
			case "wide":
				g["layout_version"] = "managed_wide_profile_v1"
				delete(g, "storage_mode")
				g["destination"].(map[string]any)["table_name"] = "fresh_scaled"
				g["members"].([]map[string]any)[0]["target_column"] = "v"
			case "foreign measurement":
				_, err := f.db.ExecContext(t.Context(), `UPDATE measurement_definitions SET point_id='foreign-point' WHERE id='mixed-meter'`)
				require.NoError(t, err)
			case "unconfirmed":
				_, err := f.db.ExecContext(t.Context(), `UPDATE measurement_definitions SET semantic_kind='gauge' WHERE id='mixed-meter'`)
				require.NoError(t, err)
			}
			response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, request)
			expected := http.StatusUnprocessableEntity
			if mode == "foreign measurement" {
				expected = http.StatusNotFound
			}
			require.Equal(t, expected, response.Code, response.Body.String())
			var n int
			require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT count(*) FROM write_groups`).Scan(&n))
			require.Zero(t, n)
		})
	}
}

func TestMixedAcceptanceScaledDraftCannotReuseOldDefinitionProofWithFreshCAS(t *testing.T) {
	f, request := mixedScaledSaveFixture(t, "delta")
	created := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, request)
	require.Equal(t, http.StatusCreated, created.Code, created.Body.String())
	data, g := groupSaveData(t, decodeJSONBody(t, created))
	id := g["id"].(string)
	_, err := f.db.ExecContext(t.Context(), `UPDATE measurement_definitions SET definition_revision='def-2',exact_recording_scale='0.01' WHERE id='mixed-meter'`)
	require.NoError(t, err)
	stale := groupEditRequest(f.writeGroupRouterFixture, data, g)
	response := performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+id, stale)
	require.Equal(t, http.StatusUnprocessableEntity, response.Code, response.Body.String())
	read := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+id, nil)
	require.Equal(t, http.StatusOK, read.Code)
	current, unchanged := groupSaveData(t, decodeJSONBody(t, read))
	require.Equal(t, g["revision"], unchanged["revision"])
	member := unchanged["members"].([]any)[0].(map[string]any)
	delete(member, "exact_recording_source")
	response = performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+id, groupEditRequest(f.writeGroupRouterFixture, current, unchanged))
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	_, fresh := groupSaveData(t, decodeJSONBody(t, response))
	proof := fresh["members"].([]any)[0].(map[string]any)["exact_recording_source"].(map[string]any)
	require.Equal(t, "def-2", proof["definition_revision"])
	require.Equal(t, "0.01", proof["scale"])
}
