package api

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/runtimeview"
)

func assertMixedScaledReadback(t *testing.T, r *mixedAcceptanceRig) {
	t.Helper()
	router := NewRouter(&DatalinkServices{Workspace: r.f.workspace, WriteGroups: r.f.groups, WriteGroupDelivery: r.pipe})
	response := performJSONRequest(t, router, http.MethodGet, writeGroupsPath+"/"+r.frozen.GroupID+"/delivery", nil)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	data := decodeJSONBody(t, response)["data"].(map[string]any)
	source := data["source_evidence"].(map[string]any)["sources"].([]any)[0].(map[string]any)
	latest := source["latest"].(map[string]any)
	require.Equal(t, true, latest["available"], latest)
	require.Equal(t, "sql_committed", latest["delivery"])
	calc := latest["calculation"].(map[string]any)
	require.Equal(t, "uint64", calc["source_type"])
	require.Equal(t, "decimal", calc["result_type"])
	require.Equal(t, true, calc["result_present"])
	require.Equal(t, true, calc["raw_present"])
	require.Equal(t, true, latest["delta_known"])
	if path := os.Getenv("MEMBER42_SOURCE_EVIDENCE"); path != "" {
		require.NoError(t, os.WriteFile(path, response.Body.Bytes(), 0600))
	}
	_, err := r.f.db.Exec(`UPDATE devices SET description='' WHERE description IS NULL`)
	require.NoError(t, err)
	resolved, resolveErr := r.pipe.ResolveAppliedConfigurationAt(t.Context(), r.frozen.GroupID, r.clock)
	require.NoError(t, resolveErr, "the same production frozen proof must be readable")
	require.NotNil(t, resolved)
	rr := &runtimeRuleEvidenceRig{f: r.f, pipe: r.pipe, clock: r.clock}
	dr := &runtimeDeliveryRig{runtimeRuleEvidenceRig: rr}
	dr.useReader(r.pipe)
	q := rr.query(t)
	response = performJSONRequest(t, dr.router, http.MethodGet, runtimeRuleContextPath+"?"+q.Encode(), nil)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	var body struct {
		Data struct {
			Evidence json.RawMessage `json:"rule_evidence"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
	if path := os.Getenv("MEMBER42_RUNTIME_EVIDENCE"); path != "" {
		require.NoError(t, os.WriteFile(path, body.Data.Evidence, 0600))
	}
	evidence, err := runtimeview.ParseRuleEvidence(body.Data.Evidence, runtimeview.Scope{WorkspaceID: r.f.record.ID, WorkspaceRevision: q.Get("workspace_revision"), DeviceID: "device-A"})
	require.NoError(t, err, string(body.Data.Evidence))
	require.Len(t, evidence.Groups, 1)
	require.NotNil(t, evidence.Groups[0].Delivery, string(body.Data.Evidence))
	current := evidence.Groups[0].Delivery.Current
	require.NotNil(t, current)
	require.Len(t, current.Sources, 1)
	selected := current.Sources[0].Latest
	require.NotNil(t, selected)
	require.NotNil(t, selected.Calculation)
	require.Equal(t, "uint64", selected.Calculation.SourceType)
	require.Equal(t, "decimal", selected.Calculation.ResultType)
	require.True(t, selected.Calculation.ResultPresent)
	if path := os.Getenv("MEMBER42_RUNTIME_EVIDENCE"); path != "" {
		require.NoError(t, os.WriteFile(path, body.Data.Evidence, 0600))
	}
}

// Corrupt retained proof must not become readable merely because mixed has a
// dedicated projection branch. Restore exact bytes before the next case.
func assertMixedReadbackRejectsForeignProof(t *testing.T, r *mixedAcceptanceRig) {
	t.Helper()
	var original string
	query := `SELECT payload FROM wg_runtime_versions WHERE workspace_id=? AND group_id=? AND group_revision=?`
	args := []any{r.frozen.WorkspaceID, r.frozen.GroupID, r.frozen.AppliedRevision}
	require.NoError(t, r.f.db.QueryRow(query, args...).Scan(&original))
	for _, field := range []string{"foreign group", "unsupported dialect", "missing digest", "unknown shape"} {
		var payload map[string]any
		require.NoError(t, json.Unmarshal([]byte(original), &payload))
		group := payload["group"].(map[string]any)
		layout := payload["runtime_layout"].(map[string]any)
		switch field {
		case "foreign group":
			group["id"] = "foreign-group"
		case "unsupported dialect":
			layout["dialect"] = "mysql"
		case "missing digest":
			layout["schema_digest"] = ""
		case "unknown shape":
			group["recording_profile"].(map[string]any)["result_shape"] = "future-shape"
		}
		raw, err := json.Marshal(payload)
		require.NoError(t, err)
		_, err = r.f.db.Exec(`UPDATE wg_runtime_versions SET payload=? WHERE workspace_id=? AND group_id=? AND group_revision=?`, string(raw), args[0], args[1], args[2])
		require.NoError(t, err)
		got, readErr := r.pipe.ResolveAppliedConfigurationAt(t.Context(), r.frozen.GroupID, r.clock)
		_, restoreErr := r.f.db.Exec(`UPDATE wg_runtime_versions SET payload=? WHERE workspace_id=? AND group_id=? AND group_revision=?`, original, args[0], args[1], args[2])
		require.NoError(t, restoreErr)
		require.Error(t, readErr, field)
		require.Nil(t, got, field)
	}
	got, err := r.pipe.ResolveAppliedConfigurationAt(t.Context(), r.frozen.GroupID, r.clock)
	require.NoError(t, err)
	require.NotNil(t, got)
}
