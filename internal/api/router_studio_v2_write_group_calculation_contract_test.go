package api

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"go-gateway/docs/swagger"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/workspace"
	"net/http"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestWriteGroupCanonicalCalculationContractKeepsLegacyEntry(t *testing.T) {
	f := newWriteGroupRouterFixture(t)
	created := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, f.createRequest())
	require.Equal(t, http.StatusCreated, created.Code, created.Body.String())
	data, group := groupSaveData(t, decodeJSONBody(t, created))
	capability := data["calculation_capabilities"].(map[string]any)
	require.Equal(t, true, capability["draft_supported"])
	require.Equal(t, true, capability["activation_supported"])
	require.Equal(t, "record_readiness_required", capability["reason"])
	require.Equal(t, false, capability["editor_supported"])
	require.NotContains(t, group, "recording_profile")
	require.NotContains(t, group["members"].([]any)[0].(map[string]any), "calculation")
	get := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+group["id"].(string), nil)
	require.Equal(t, http.StatusOK, get.Code)
	got, legacy := groupSaveData(t, decodeJSONBody(t, get))
	require.Equal(t, group, legacy)
	require.Equal(t, capability, got["calculation_capabilities"])
	list := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath, nil)
	require.Equal(t, http.StatusOK, list.Code)
	require.Equal(t, capability, decodeJSONBody(t, list)["data"].(map[string]any)["calculation_capabilities"])
	assertGeneratedCalculationSchema(t, capability)
}

func TestWriteGroupMixedAPIReadsMemberIntentWithoutChangingPoint(t *testing.T) {
	f := newWriteGroupRouterFixture(t)
	req := f.createRequest()
	g := req["group"].(map[string]any)
	g["record_id"] = "6b2669b5-02df-455a-ab6f-0614438257e8"
	g["layout_version"] = "managed_wide_profile_v1"
	g["recording_profile"] = map[string]any{"kind": "mixed", "result_shape": "member-calculation-v1"}
	g["row_policy"] = map[string]any{"incomplete_policy": "partial"}
	g["destination"].(map[string]any)["storage_strategy"] = "managed"
	member := g["members"].([]map[string]any)[0]
	member["calculation"] = "mean"
	member["cadence"] = map[string]any{"interval": 15, "unit": "second", "alignment": "utc_epoch"}
	member["value_role"] = "gauge"
	saved := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, req)
	require.Equal(t, http.StatusCreated, saved.Code, saved.Body.String())
	_, group := groupSaveData(t, decodeJSONBody(t, saved))
	got := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+group["id"].(string), nil)
	require.Equal(t, http.StatusOK, got.Code)
	data, readback := groupSaveData(t, decodeJSONBody(t, got))
	require.Equal(t, group, readback)
	require.Equal(t, g["record_id"], readback["record_id"])
	require.NotContains(t, readback["row_policy"], "interval_seconds")
	m := readback["members"].([]any)[0].(map[string]any)
	require.Equal(t, "mean", m["calculation"])
	require.Equal(t, float64(15), m["cadence"].(map[string]any)["interval"])
	require.NotEmpty(t, m["source_revision"])
	require.NotEmpty(t, m["mapping_revision"])
	require.Equal(t, "float32", m["source_value_type"])
	require.Equal(t, "float64", m["result_value_type"])
	require.Equal(t, true, data["calculation_capabilities"].(map[string]any)["activation_supported"])
	require.Equal(t, false, data["calculation_capabilities"].(map[string]any)["editor_supported"])
	require.Equal(t, "record_readiness_required", data["calculation_capabilities"].(map[string]any)["reason"])
	var sourceType string
	require.NoError(t, f.db.QueryRow(`SELECT data_type FROM points WHERE id='point-A'`).Scan(&sourceType))
	require.Equal(t, "float32", sourceType)
}

func assertGeneratedCalculationSchema(t *testing.T, actual map[string]any) {
	t.Helper()
	var spec map[string]any
	require.NoError(t, json.Unmarshal([]byte(swagger.SwaggerInfo.ReadDoc()), &spec))
	defs := spec["definitions"].(map[string]any)
	statusField, ok := reflect.TypeOf(runtimeview.CalculationResult{}).FieldByName("Status")
	require.True(t, ok)
	statusDoc := defs["runtimeview.CalculationResult"].(map[string]any)["properties"].(map[string]any)["status"].(map[string]any)
	var documentedStatus []string
	for _, value := range statusDoc["enum"].([]any) {
		documentedStatus = append(documentedStatus, value.(string))
	}
	require.Equal(t, strings.Split(statusField.Tag.Get("enum"), ","), documentedStatus, "generated member statuses match runtime validator; partial is coverage/state only")

	for _, name := range []string{"workspace.WriteGroupSaveResult", "workspace.WriteGroupListResult"} {
		props := defs[name].(map[string]any)["properties"].(map[string]any)
		field, ok := props["calculation_capabilities"].(map[string]any)
		require.True(t, ok, "generated canonical response must expose actual wire capability")
		require.Equal(t, true, field["readOnly"])
		require.Equal(t, "#/definitions/workspace.WriteGroupCalculationCapabilities", field["allOf"].([]any)[0].(map[string]any)["$ref"])
	}
	for name, value := range map[string]map[string]any{
		"workspace.WriteGroupCalculationCapabilities": actual,
		"workspace.WriteGroupCadenceCapability":       actual["cadence"].(map[string]any),
	} {
		props := defs[name].(map[string]any)["properties"].(map[string]any)
		var documented, wire []string
		for key := range props {
			documented = append(documented, key)
		}
		for key := range value {
			wire = append(wire, key)
		}
		sort.Strings(documented)
		sort.Strings(wire)
		require.Equal(t, wire, documented, name)
	}
	assertGeneratedCalculationRestriction(t, defs, actual["unsupported_combinations"].([]any))
}

func assertGeneratedCalculationRestriction(t *testing.T, defs map[string]any, wire []any) {
	t.Helper()
	doc := defs["workspace.WriteGroupCalculationRestriction"].(map[string]any)
	props := doc["properties"].(map[string]any)
	typ := reflect.TypeOf(workspace.WriteGroupCalculationRestriction{})
	required, all := []string{}, []string{}
	for i := 0; i < typ.NumField(); i++ {
		tag := strings.Split(typ.Field(i).Tag.Get("json"), ",")
		all = append(all, tag[0])
		if len(tag) == 1 {
			required = append(required, tag[0])
		}
	}
	var documented []string
	for key := range props {
		documented = append(documented, key)
	}
	require.ElementsMatch(t, all, documented, "schema declares every required and optional Go wire field, and no unknown property")
	require.ElementsMatch(t, []string{"method", "phase", "source_transform", "reason"}, required)
	if requiredSchema, present := doc["required"]; present {
		require.NotContains(t, requiredSchema, "layout", "layout is an optional narrower restriction, not a required default")
	}
	layout := props["layout"].(map[string]any)
	require.Equal(t, "#/definitions/workspace.WriteGroupLayoutVersion", layout["allOf"].([]any)[0].(map[string]any)["$ref"])
	require.Equal(t, "string", defs["workspace.WriteGroupLayoutVersion"].(map[string]any)["type"])
	require.Equal(t, []any{"managed_wide_profile_v1", "managed_long_v1"}, layout["enum"])
	foundLayout := false
	foundOmitted := false
	for _, raw := range wire {
		restriction := raw.(map[string]any)
		for _, key := range required {
			require.Contains(t, restriction, key)
		}
		for key := range restriction {
			require.Contains(t, props, key, "actual HTTP must not emit undocumented restriction keys")
		}
		if value, present := restriction["layout"]; present {
			foundLayout = true
			require.IsType(t, "", value)
			require.Equal(t, "managed_wide_profile_v1", value)
			require.Equal(t, "delta", restriction["method"])
			require.Contains(t, layout["enum"], value)
		} else {
			foundOmitted = true
		}
	}
	require.True(t, foundLayout, "actual HTTP includes the supported-long/refused-wide delta distinction")
	require.True(t, foundOmitted, "all-layout restrictions omit layout rather than inventing a default")
}
