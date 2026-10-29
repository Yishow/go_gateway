package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
	"time"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestRuntimeRuleEvidenceSavedCRequestedBEffectiveAThenB(t *testing.T) {
	r := newRuntimeRuleEvidenceRig(t)
	a := r.f.group["applied_revision"].(string)
	r.clock = time.Date(2026, 10, 2, 12, 0, 17, 0, time.UTC)
	r.update(t, "meter_new", 30)
	b := r.apply(t).Group.AppliedRevision
	c := r.update(t, "meter_draft", 45).Group.Revision
	require.NotEqual(t, a, b)
	require.NotEqual(t, b, c)
	router := r.deterministicRouter(r.f.groups)
	for _, tc := range []struct {
		second          int
		revision, table string
		interval        int64
	}{
		{29, a, "meter_old", 15}, {30, b, "meter_new", 30},
	} {
		r.clock = time.Date(2026, 10, 2, 12, 0, tc.second, 0, time.UTC)
		before := lifecycleAPIState(t, r.f.db)
		response := performJSONRequest(t, router, http.MethodGet, runtimeRuleContextPath+"?"+r.query(t).Encode(), nil)
		require.Equal(t, http.StatusOK, response.Code, response.Body.String())
		var body struct {
			Data struct {
				Evidence runtimeview.RuleEvidenceView `json:"rule_evidence"`
			} `json:"data"`
		}
		require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
		v := body.Data.Evidence
		require.Len(t, v.Groups, 1)
		g := v.Groups[0]
		require.Equal(t, c, *g.SavedRevision)
		require.Equal(t, b, *g.AppliedRevision)
		require.Equal(t, tc.revision, *g.EffectiveRevision)
		require.Equal(t, tc.table, g.Effective.Destination.Table)
		require.Equal(t, tc.interval, *g.Effective.IntervalSeconds)
		require.Equal(t, "meter_draft", g.Saved.Destination.Table)
		require.EqualValues(t, 45, *g.Saved.IntervalSeconds)
		require.Equal(t, r.clock.Format(time.RFC3339Nano), *v.ObservedAt)
		require.Equal(t, *v.ObservedAt, *g.AsOf)
		if tc.revision == b {
			require.NotNil(t, g.AppliedEffectiveAt)
			require.Equal(t, "2026-10-02T12:00:30Z", *g.AppliedEffectiveAt)
		} else {
			require.Nil(t, g.AppliedEffectiveAt, "requested B activation is not proven by the resolved A snapshot")
		}
		require.Nil(t, g.Delivery)
		require.Equal(t, before, lifecycleAPIState(t, r.f.db))
	}
}

func TestRuntimeRuleEvidenceMissingProofNeverUsesSavedAsEffective(t *testing.T) {
	for _, mutation := range []string{`DELETE FROM wg_runtime_versions`, `UPDATE wg_runtime_versions SET payload='{"private":"/private/fixture-secret.sqlite password=fixture-password"}'`} {
		t.Run(mutation[:6], func(t *testing.T) {
			r := newRuntimeRuleEvidenceRig(t)
			_, err := r.f.db.Exec(mutation)
			require.NoError(t, err)
			before := lifecycleAPIState(t, r.f.db)
			response := performJSONRequest(t, r.router, http.MethodGet, runtimeRuleContextPath+"?"+r.query(t).Encode(), nil)
			require.Equal(t, http.StatusOK, response.Code, response.Body.String())
			evidence := decodeJSONBody(t, response)["data"].(map[string]any)["rule_evidence"].(map[string]any)
			group := evidence["groups"].([]any)[0].(map[string]any)
			require.Equal(t, "unavailable", group["state"])
			require.Equal(t, "unavailable", group["proof"])
			require.Nil(t, group["effective"])
			require.NotNil(t, group["saved"])
			require.NotContains(t, response.Body.String(), "fixture-secret")
			require.NotContains(t, response.Body.String(), "fixture-password")
			require.Equal(t, before, lifecycleAPIState(t, r.f.db))
		})
	}
}

func TestRuntimeRuleEvidenceActualSaveDuringReadIsStale(t *testing.T) {
	r := newRuntimeRuleEvidenceRig(t)
	r.clock = time.Date(2026, 10, 2, 12, 0, 17, 0, time.UTC)
	changing := &runtimeRuleChangingGroups{groups: r.f.groups, afterFirst: func() { r.update(t, "meter_changed", 30) }}
	response := performJSONRequest(t, r.deterministicRouter(changing), http.MethodGet, runtimeRuleContextPath+"?"+r.query(t).Encode(), nil)
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
	require.Contains(t, response.Body.String(), "stale-reference")
	require.NotContains(t, response.Body.String(), `"rule_evidence"`)
}

func TestRuntimeRuleEvidenceStoppedOrDeletedFactsDoNotGrantIntake(t *testing.T) {
	for _, removed := range []bool{false, true} {
		t.Run(map[bool]string{false: "disabled", true: "deleted"}[removed], func(t *testing.T) {
			r := newRuntimeRuleEvidenceRig(t)
			r.clock = time.Date(2026, 10, 2, 12, 0, 17, 0, time.UTC)
			current, err := r.f.groups.Get(t.Context(), r.f.group["id"].(string))
			require.NoError(t, err)
			mutation := workspace.WriteGroupMutation{WorkspaceID: current.Group.WorkspaceID, ExpectedWorkspaceRevision: current.WorkspaceRevision, ExpectedGroupRevision: current.Group.Revision, ExpectedConnectorRevision: current.Group.Destination.ConnectorRevision}
			if removed {
				r.f.groups.WithBacklogOwnershipGuard(&lifecycleAPIGuard{})
				_, err = r.f.groups.Delete(t.Context(), current.Group.ID, mutation)
			} else {
				_, err = r.f.groups.Disable(t.Context(), current.Group.ID, mutation)
			}
			require.NoError(t, err)
			before := lifecycleAPIState(t, r.f.db)
			response := performJSONRequest(t, r.deterministicRouter(r.f.groups), http.MethodGet, runtimeRuleContextPath+"?"+r.query(t).Encode(), nil)
			require.Equal(t, http.StatusOK, response.Code, response.Body.String())
			group := decodeJSONBody(t, response)["data"].(map[string]any)["rule_evidence"].(map[string]any)["groups"].([]any)[0].(map[string]any)
			require.NotNil(t, group["effective"], "retained immutable configuration may explain history")
			require.Nil(t, group["delivery"], "snapshot presence must never imply intake permission")
			require.Equal(t, false, group["saved"].(map[string]any)["enabled"])
			require.Equal(t, before, lifecycleAPIState(t, r.f.db))
		})
	}
}

func TestRuntimeRuleEvidenceRetainsActualCutoverEndWithoutIntake(t *testing.T) {
	cut, _ := finalizedAPIRig(t, cutoverOptions{})
	boundary := cut.at(22500 * time.Millisecond)
	cut.cutover(boundary)
	_, err := cut.f.db.Exec(`UPDATE devices SET description='' WHERE id='device-A'`)
	require.NoError(t, err)
	id := cut.wide["id"].(string)
	r := &runtimeRuleEvidenceRig{f: cut.f, pipe: cut.pipe, clock: boundary.Add(time.Second)}
	r.f.data = getGroupData(t, r.f, id)
	r.f.group = r.f.data["group"].(map[string]any)
	before := lifecycleAPIState(t, cut.f.db)
	response := performJSONRequest(t, r.deterministicRouter(r.f.groups), http.MethodGet, runtimeRuleContextPath+"?"+r.query(t).Encode(), nil)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	var body struct {
		Data struct {
			Evidence runtimeview.RuleEvidenceView `json:"rule_evidence"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
	var old *runtimeview.GroupView
	for i := range body.Data.Evidence.Groups {
		if body.Data.Evidence.Groups[i].GroupID == id {
			old = &body.Data.Evidence.Groups[i]
		}
	}
	require.NotNil(t, old)
	require.NotNil(t, old.Effective)
	require.NotNil(t, old.Effective.EndsAt)
	ended, err := time.Parse(time.RFC3339Nano, *old.Effective.EndsAt)
	require.NoError(t, err)
	require.True(t, boundary.Equal(ended))
	require.Nil(t, old.Delivery, "an ended retained snapshot cannot grant intake")
	require.Equal(t, before, lifecycleAPIState(t, cut.f.db))
}

func TestRuntimeRuleEvidenceFrozenDestinationDoesNotNeedLiveConnector(t *testing.T) {
	r := newRuntimeRuleEvidenceRig(t)
	_, err := r.f.db.Exec(`DELETE FROM database_connectors WHERE id='connector-A'`)
	require.NoError(t, err)
	before := lifecycleAPIState(t, r.f.db)
	response := performJSONRequest(t, r.router, http.MethodGet, runtimeRuleContextPath+"?"+r.query(t).Encode(), nil)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	require.Contains(t, response.Body.String(), `"state":"effective"`)
	require.Contains(t, response.Body.String(), `"table":"meter_old"`)
	require.NotContains(t, response.Body.String(), r.f.file)
	for _, key := range []string{`"database"`, `"dsn"`, `"connection_config"`, `"password"`, `"runtime_layout"`} {
		require.NotContains(t, response.Body.String(), key)
	}
	require.Equal(t, before, lifecycleAPIState(t, r.f.db))
}

func TestRuntimeRuleEvidenceProductionRouteUsesAppliedProof(t *testing.T) {
	r := newRuntimeRuleEvidenceRig(t)
	before := lifecycleAPIState(t, r.f.db)
	response := performJSONRequest(t, r.router, http.MethodGet, runtimeRuleContextPath+"?"+r.query(t).Encode(), nil)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	data := decodeJSONBody(t, response)["data"].(map[string]any)
	require.Nil(t, data["setup"], "scoped evidence must not invoke bootstrap-capable setup readers")
	require.Empty(t, data["devices"])
	require.Contains(t, data, "rule_evidence", "the real route must wire the verified production provider")
	evidence := data["rule_evidence"].(map[string]any)
	require.Equal(t, r.f.record.ID, evidence["workspace_id"])
	require.Equal(t, "device-A", evidence["device_id"])
	groups := evidence["groups"].([]any)
	require.Len(t, groups, 1)
	group := groups[0].(map[string]any)
	require.Equal(t, "effective", group["state"])
	require.Equal(t, r.f.group["applied_revision"], group["effective_revision"])
	delivery, ok := group["delivery"].(map[string]any)
	require.True(t, ok)
	current, ok := delivery["current"].(map[string]any)
	require.True(t, ok)
	binding := current["binding"].(map[string]any)
	require.Equal(t, r.f.record.ID, binding["workspace_id"])
	require.Equal(t, evidence["workspace_revision"], binding["workspace_revision"])
	require.Equal(t, group["effective_revision"], binding["entity_revision"])
	require.Equal(t, "group", current["totals_scope"])
	require.Equal(t, map[string]any{
		"collecting": float64(0), "queued": float64(0), "retrying": float64(0), "blocked": float64(0),
		"quarantined": float64(0), "unknown": float64(0), "sql_committed": float64(0), "skipped": float64(0),
	}, current["stages"])
	require.Nil(t, current["last_commit"], "observed empty delivery is not a commit")
	require.Nil(t, current["intake"], "configuration proof is not revision-bound intake permission")
	require.Empty(t, delivery["backlog"])
	require.Equal(t, before, lifecycleAPIState(t, r.f.db))
}

func TestRuntimeRuleEvidenceConcurrentWorkspaceDeletionDoesNotRecreate(t *testing.T) {
	r := newRuntimeRuleEvidenceRig(t)
	r.clock = time.Date(2026, 10, 2, 12, 0, 17, 0, time.UTC)
	query := r.query(t)
	deleting := &runtimeRuleDeletingWorkspaceGroups{groups: r.f.groups, delete: func() {
		_, err := r.f.db.Exec(`DELETE FROM system_settings WHERE key='studio_v2_workspace'`)
		require.NoError(t, err)
	}}
	response := performJSONRequest(t, r.deterministicRouter(deleting), http.MethodGet, runtimeRuleContextPath+"?"+query.Encode(), nil)
	var count int
	require.NoError(t, r.f.db.QueryRow(`SELECT COUNT(*) FROM system_settings WHERE key='studio_v2_workspace'`).Scan(&count))
	require.Zero(t, count, "final consistency read must never replace the deleted workspace")
	require.Equal(t, http.StatusServiceUnavailable, response.Code, response.Body.String())
	require.NotContains(t, response.Body.String(), `"rule_evidence"`)
}

func TestRuntimeRuleEvidenceMissingProviderIsUnavailableButLegacyCompatible(t *testing.T) {
	r := newRuntimeRuleEvidenceRig(t)
	router := NewRouter(&DatalinkServices{Workspace: r.f.workspace, WriteGroups: r.f.groups,
		Device: device.NewService(device.NewSQLRepository(r.f.db), nil)})
	before := lifecycleAPIState(t, r.f.db)
	response := performJSONRequest(t, router, http.MethodGet, runtimeRuleContextPath+"?"+r.query(t).Encode(), nil)
	require.Equal(t, http.StatusServiceUnavailable, response.Code, response.Body.String())
	require.Contains(t, response.Body.String(), "unavailable")
	require.NotContains(t, response.Body.String(), `"rule_evidence"`)
	legacy := performJSONRequest(t, router, http.MethodGet, runtimeRuleContextPath, nil)
	require.Equal(t, http.StatusOK, legacy.Code, legacy.Body.String())
	require.NotContains(t, legacy.Body.String(), `"rule_evidence"`)
	require.Equal(t, before, lifecycleAPIState(t, r.f.db))
}

func TestRuntimeRuleEvidenceScopedQueryFailsClosed(t *testing.T) {
	r := newRuntimeRuleEvidenceRig(t)
	for _, tc := range []struct {
		name   string
		edit   func(url.Values)
		code   int
		reason string
	}{
		{"partial", func(q url.Values) { q.Del("workspace_revision") }, http.StatusBadRequest, "unavailable"},
		{"empty", func(q url.Values) { q.Set("device_id", "") }, http.StatusBadRequest, "unavailable"},
		{"duplicate", func(q url.Values) { q.Add("device_id", "device-A") }, http.StatusBadRequest, "unavailable"},
		{"unsafe", func(q url.Values) { q.Set("device_id", "private.example.internal") }, http.StatusBadRequest, "unavailable"},
		{"foreign_workspace", func(q url.Values) { q.Set("workspace_id", "foreign-workspace") }, http.StatusForbidden, "foreign-scope"},
		{"foreign_device", func(q url.Values) { q.Set("device_id", "foreign-device") }, http.StatusForbidden, "foreign-scope"},
		{"stale", func(q url.Values) { q.Set("workspace_revision", "old-revision") }, http.StatusConflict, "stale-reference"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			q := r.query(t)
			tc.edit(q)
			response := performJSONRequest(t, r.router, http.MethodGet, runtimeRuleContextPath+"?"+q.Encode(), nil)
			require.Equal(t, tc.code, response.Code, response.Body.String())
			require.Contains(t, response.Body.String(), tc.reason)
			require.NotContains(t, response.Body.String(), `"rule_evidence"`)
		})
	}
	legacy := performJSONRequest(t, r.router, http.MethodGet, runtimeRuleContextPath, nil)
	require.Equal(t, http.StatusOK, legacy.Code, legacy.Body.String())
	require.NotContains(t, legacy.Body.String(), `"rule_evidence"`)
}
