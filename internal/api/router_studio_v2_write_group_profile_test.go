package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const ensureCircuitProfilePath = writeGroupsPath + "/circuit-profiles/ensure"

// seedProfileSource adds one more mapped point and tag to device-A.
func seedProfileSource(t *testing.T, f writeGroupRouterFixture, name string) {
	t.Helper()
	ctx := t.Context()
	_, err := f.db.ExecContext(ctx, `INSERT INTO points (id,device_id,name,address,function,data_type,mode,enabled)
		VALUES (?,'device-A',?,?,'FC03','float32','read',1)`, "point-"+name, name, "4"+fmt.Sprint(len(name))+"0"+name)
	require.NoError(t, err)
	_, err = f.db.ExecContext(ctx, `INSERT INTO tags (id,key,key_lower,display_name,data_type,status,labels)
		VALUES (?,?,?,?,'float32','active','{}')`, "tag-"+name, name, name, name)
	require.NoError(t, err)
	_, err = f.db.ExecContext(ctx, `INSERT INTO mappings (id,point_id,tag_id,transform_pipeline,status,enabled)
		VALUES (?,?,?,'[]','active',1)`, "mapping-"+name, "point-"+name, "tag-"+name)
	require.NoError(t, err)
}

// profileRequest builds a managed circuit-profile request. interval and unit
// are sent exactly as given; nil omits the key. The request carries the
// current workspace revision, as a client that has just reloaded would.
func profileRequest(t *testing.T, f writeGroupRouterFixture, kind string, interval, unit any, circuitID, source string) map[string]any {
	profile := map[string]any{"kind": kind}
	if interval != nil {
		profile["interval"] = interval
	}
	if unit != nil {
		profile["unit"] = unit
	}
	current, err := f.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	pointID, tagID := "point-A", "tag-A"
	if source != "" {
		pointID, tagID = "point-"+source, "tag-"+source
	}
	return map[string]any{
		"workspace_id": f.record.ID, "expected_workspace_revision": current.DatabaseSetupRevision,
		"expected_connector_revision": "connector-1",
		"group": map[string]any{
			"workspace_id": f.record.ID, "name": circuitID + " " + kind,
			"members": []map[string]any{{
				"device_id": "device-A", "point_id": pointID, "tag_id": tagID, "required": true,
				"circuit_id": circuitID, "circuit_name": "Circuit " + circuitID,
			}},
			"destination":       map[string]any{"connector_id": "connector-A", "connector_revision": "connector-1", "storage_strategy": "managed"},
			"row_policy":        map[string]any{"interval_seconds": 999, "allowed_lateness_seconds": 0},
			"recording_profile": profile,
		},
	}
}

func requireProfile(t *testing.T, group map[string]any, kind string, interval float64, unit string, seconds float64) {
	t.Helper()
	profile, ok := group["recording_profile"].(map[string]any)
	require.True(t, ok, "recording_profile missing from %v", group)
	require.Equal(t, kind, profile["kind"])
	require.Equal(t, interval, profile["interval"])
	require.Equal(t, unit, profile["unit"], "the public unit must survive")
	require.Equal(t, seconds, group["row_policy"].(map[string]any)["interval_seconds"])
}

func TestStudioV2WriteGroupsEnsureCircuitProfileRoundTrip(t *testing.T) {
	f := newWriteGroupRouterFixture(t)
	seedProfileSource(t, f, "energy")

	power := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath,
		profileRequest(t, f, "instantaneous", 15, "second", "main", ""))
	require.Equal(t, http.StatusOK, power.Code, power.Body.String())
	_, powerGroup := groupSaveData(t, decodeJSONBody(t, power))
	requireProfile(t, powerGroup, "instantaneous", 15, "second", 15)
	require.Equal(t, "main", powerGroup["members"].([]any)[0].(map[string]any)["circuit_id"])
	require.NotContains(t, powerGroup, "layout_version", "an unset layout stays absent")

	energy := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath,
		profileRequest(t, f, "cumulative", 10, "minute", "main", "energy"))
	require.Equal(t, http.StatusOK, energy.Code, energy.Body.String())
	_, energyGroup := groupSaveData(t, decodeJSONBody(t, energy))
	requireProfile(t, energyGroup, "cumulative", 10, "minute", 600)
	require.NotEqual(t, powerGroup["id"], energyGroup["id"])

	t.Run("a resent create returns the same group", func(t *testing.T) {
		lost := profileRequest(t, f, "instantaneous", 15, "second", "main", "")
		lost["expected_workspace_revision"] = f.record.DatabaseSetupRevision // the response never arrived
		resent := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, lost)
		require.Equal(t, http.StatusOK, resent.Code, resent.Body.String())
		_, again := groupSaveData(t, decodeJSONBody(t, resent))
		require.Equal(t, powerGroup["id"], again["id"])
	})

	t.Run("a resent create with another interval is a conflict", func(t *testing.T) {
		changed := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath,
			profileRequest(t, f, "instantaneous", 30, "second", "main", ""))
		require.Equal(t, http.StatusConflict, changed.Code, changed.Body.String())
		assertWriteGroupError(t, decodeJSONBody(t, changed), "WRITE_GROUP_CIRCUIT_PROFILE_CONFLICT")
	})

	t.Run("get and list return the saved profile", func(t *testing.T) {
		got := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+powerGroup["id"].(string), nil)
		require.Equal(t, http.StatusOK, got.Code, got.Body.String())
		_, gotGroup := groupSaveData(t, decodeJSONBody(t, got))
		require.Equal(t, powerGroup, gotGroup)

		listed := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath, nil)
		require.Equal(t, http.StatusOK, listed.Code, listed.Body.String())
		groups := decodeJSONBody(t, listed)["data"].(map[string]any)["groups"].([]any)
		require.Len(t, groups, 2)
		requireProfile(t, groups[0].(map[string]any), "instantaneous", 15, "second", 15)
		requireProfile(t, groups[1].(map[string]any), "cumulative", 10, "minute", 600)
	})

	t.Run("update keeps the unit, re-derives seconds and refuses out-of-range values", func(t *testing.T) {
		current := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+energyGroup["id"].(string), nil)
		data, group := groupSaveData(t, decodeJSONBody(t, current))
		group["recording_profile"].(map[string]any)["interval"] = 600
		updated := performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+energyGroup["id"].(string), groupEditRequest(f, data, group))
		require.Equal(t, http.StatusOK, updated.Code, updated.Body.String())
		updatedData, updatedGroup := groupSaveData(t, decodeJSONBody(t, updated))
		requireProfile(t, updatedGroup, "cumulative", 600, "minute", 36000)

		updatedGroup["recording_profile"].(map[string]any)["interval"] = 601
		rejected := performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+energyGroup["id"].(string), groupEditRequest(f, updatedData, updatedGroup))
		require.Equal(t, http.StatusUnprocessableEntity, rejected.Code, rejected.Body.String())
		assertWriteGroupError(t, decodeJSONBody(t, rejected), "profile_interval_invalid")

		after := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+energyGroup["id"].(string), nil)
		_, afterGroup := groupSaveData(t, decodeJSONBody(t, after))
		requireProfile(t, afterGroup, "cumulative", 600, "minute", 36000)
	})
}

func TestStudioV2WriteGroupsProfileLayoutRoundTrip(t *testing.T) {
	f := newWriteGroupRouterFixture(t)
	seedProfileSource(t, f, "third")

	explicit := profileRequest(t, f, "instantaneous", 15, "second", "wide", "third")
	explicit["group"].(map[string]any)["layout_version"] = "managed_wide_v1"
	created := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, explicit)
	require.Equal(t, http.StatusOK, created.Code, created.Body.String())
	_, group := groupSaveData(t, decodeJSONBody(t, created))
	require.Equal(t, "managed_wide_v1", group["layout_version"])
	got := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+group["id"].(string), nil)
	_, gotGroup := groupSaveData(t, decodeJSONBody(t, got))
	require.Equal(t, "managed_wide_v1", gotGroup["layout_version"])

	unknown := profileRequest(t, f, "instantaneous", 15, "second", "long", "")
	unknown["group"].(map[string]any)["layout_version"] = "managed_future_v9"
	rejected := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, unknown)
	require.Equal(t, http.StatusUnprocessableEntity, rejected.Code, rejected.Body.String())
	assertWriteGroupError(t, decodeJSONBody(t, rejected), "WRITE_GROUP_INVALID")
}

func TestStudioV2WriteGroupsProfileIntervalBoundaries(t *testing.T) {
	type boundary struct {
		name     string
		kind     string
		interval any
		unit     any
		seconds  float64 // zero means the request must be rejected
	}
	cases := []boundary{
		{"instantaneous 1 second", "instantaneous", 1, "second", 1},
		{"instantaneous 300 seconds", "instantaneous", 300, "second", 300},
		{"cumulative 1 minute", "cumulative", 1, "minute", 60},
		{"cumulative 600 minutes", "cumulative", 600, "minute", 36000},
		{"instantaneous zero", "instantaneous", 0, "second", 0},
		{"instantaneous 301", "instantaneous", 301, "second", 0},
		{"instantaneous negative", "instantaneous", -5, "second", 0},
		{"cumulative zero", "cumulative", 0, "minute", 0},
		{"cumulative 601", "cumulative", 601, "minute", 0},
		{"instantaneous fraction", "instantaneous", 15.5, "second", 0},
		{"cumulative fraction", "cumulative", 2.5, "minute", 0},
		{"integral float literal", "instantaneous", json.RawMessage("15.0"), "second", 0},
		{"exponent literal", "instantaneous", json.RawMessage("1e1"), "second", 0},
		{"quoted number", "instantaneous", "15", "second", 0},
		{"null interval", "instantaneous", json.RawMessage("null"), "second", 0},
		{"missing interval", "cumulative", nil, "minute", 0},
		{"missing unit", "instantaneous", 15, nil, 0},
		{"seconds for a cumulative profile", "cumulative", 10, "second", 0},
		{"minutes for an instantaneous profile", "instantaneous", 15, "minute", 0},
		{"unknown unit", "instantaneous", 15, "hour", 0},
	}
	routes := map[string]string{"ensure": ensureCircuitProfilePath, "create": writeGroupsPath}
	for routeName, path := range routes {
		for _, tc := range cases {
			t.Run(routeName+" "+tc.name, func(t *testing.T) {
				f := newWriteGroupRouterFixture(t)
				request := profileRequest(t, f, tc.kind, tc.interval, tc.unit, "main", "")
				response := performJSONRequest(t, f.router, http.MethodPost, path, request)
				if tc.seconds > 0 {
					require.Less(t, response.Code, 300, response.Body.String())
					_, group := groupSaveData(t, decodeJSONBody(t, response))
					profile := group["recording_profile"].(map[string]any)
					require.Equal(t, tc.unit, profile["unit"])
					require.Equal(t, tc.seconds, group["row_policy"].(map[string]any)["interval_seconds"])
					return
				}
				require.Equal(t, http.StatusUnprocessableEntity, response.Code, response.Body.String())
				assertWriteGroupError(t, decodeJSONBody(t, response), "profile_interval_invalid")
				body := response.Body.String()
				for _, leak := range []string{"cannot unmarshal", "json:", "Go struct", "workspace."} {
					require.NotContains(t, body, leak)
				}
				listed := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath, nil)
				require.Empty(t, decodeJSONBody(t, listed)["data"].(map[string]any)["groups"])
			})
		}
	}
}

func TestStudioV2WriteGroupsProfileConflictAndLegacyCompatibility(t *testing.T) {
	f := newWriteGroupRouterFixture(t)
	seedProfileSource(t, f, "second")

	t.Run("a legacy group keeps its 60 second interval and gains no profile", func(t *testing.T) {
		request := f.createRequest()
		request["group"].(map[string]any)["row_policy"] = map[string]any{"interval_seconds": 60, "allowed_lateness_seconds": 0}
		created := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, request)
		require.Equal(t, http.StatusCreated, created.Code, created.Body.String())
		data, group := groupSaveData(t, decodeJSONBody(t, created))
		require.Equal(t, float64(60), group["row_policy"].(map[string]any)["interval_seconds"])
		for _, key := range []string{"recording_profile", "layout_version"} {
			require.NotContains(t, group, key)
		}
		group["name"] = "legacy renamed"
		updated := performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+group["id"].(string), groupEditRequest(f, data, group))
		require.Equal(t, http.StatusOK, updated.Code, updated.Body.String())
		_, updatedGroup := groupSaveData(t, decodeJSONBody(t, updated))
		require.Equal(t, float64(60), updatedGroup["row_policy"].(map[string]any)["interval_seconds"])
		require.NotContains(t, updatedGroup, "recording_profile")
		require.NotContains(t, strings.ToLower(updated.Body.String()), "recording_profile")
	})

	t.Run("a profile cannot be attached to a legacy group by update", func(t *testing.T) {
		listed := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath, nil)
		listData := decodeJSONBody(t, listed)["data"].(map[string]any)
		group := listData["groups"].([]any)[0].(map[string]any)
		group["recording_profile"] = map[string]any{"kind": "instantaneous", "interval": 15, "unit": "second"}
		request := map[string]any{
			"workspace_id": f.record.ID, "expected_workspace_revision": listData["workspace_revision"],
			"expected_group_revision": group["revision"], "expected_connector_revision": "connector-1", "group": group,
		}
		response := performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+group["id"].(string), request)
		require.Equal(t, http.StatusUnprocessableEntity, response.Code, response.Body.String())
		assertWriteGroupError(t, decodeJSONBody(t, response), "WRITE_GROUP_INVALID")
	})

	t.Run("another source for the same circuit profile is a conflict", func(t *testing.T) {
		first := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath,
			profileRequest(t, f, "instantaneous", 15, "second", "main", ""))
		require.Equal(t, http.StatusOK, first.Code, first.Body.String())
		other := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath,
			profileRequest(t, f, "instantaneous", 15, "second", "main", "second"))
		require.Equal(t, http.StatusConflict, other.Code, other.Body.String())
		assertWriteGroupError(t, decodeJSONBody(t, other), "WRITE_GROUP_CIRCUIT_PROFILE_CONFLICT")
	})
}
