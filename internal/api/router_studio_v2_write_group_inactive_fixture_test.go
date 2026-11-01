package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

const inactiveDraftRulePath = "/api/v1/datalink/studio-v2/workspace/source-rules"

type inactiveDraftFixture struct {
	writeGroupRouterFixture
	link        *schema.SourceRuleLink
	data, group map[string]any
}

func newInactiveDraftFixture(t *testing.T, configure ...func(map[string]any)) inactiveDraftFixture {
	t.Helper()
	f := newWriteGroupRouterFixture(t)
	_, err := f.db.ExecContext(t.Context(), `UPDATE devices SET description='',connection_config='{"host":"127.0.0.1","port":502,"slave_id":1,"timeout":2}',last_test_success=1 WHERE id='device-A'`)
	require.NoError(t, err)
	devices := device.NewService(device.NewSQLRepository(f.db), nil)
	points := point.NewService(point.NewSQLRepository(f.db), pollinggroup.NewSQLRepository(f.db))
	tags := tag.NewService(tag.NewSQLRepository(f.db))
	mappings := mapping.NewServiceWithTagResolver(mapping.NewSQLRepository(f.db), tags.GetByID)
	rules := sourcerule.NewService(sourcerule.NewSQLRepository(f.db), devices, points, nil)
	rules.SetTagMappingServices(tags, mappings)
	f.router = NewRouter(&DatalinkServices{Workspace: f.workspace, WriteGroups: workspace.NewWriteGroupService(f.workspace, workspace.NewSQLWriteGroupRepository(f.db)), Device: devices, Point: points, Tag: tags, Mapping: mappings, SourceRule: rules})
	created := performJSONRequest(t, f.router, http.MethodPost, inactiveDraftRulePath, map[string]any{
		"id": "inactive-rule", "device_id": "device-A", "start_address": "41201", "count": 1,
		"data_type": "int16", "naming_prefix": "Known", "enabled": true,
	})
	require.Equal(t, http.StatusCreated, created.Code, created.Body.String())
	candidates := performJSONRequest(t, f.router, http.MethodGet, "/api/v1/datalink/source-rules/inactive-rule/candidates", nil)
	require.Equal(t, http.StatusOK, candidates.Code, candidates.Body.String())
	view := decodeJSONBody(t, candidates)["data"].(map[string]any)
	tagCandidates := view["tags"].(map[string]any)["candidates"].([]any)
	require.Len(t, tagCandidates, 1)
	applied := performJSONRequest(t, f.router, http.MethodPost, "/api/v1/datalink/source-rules/inactive-rule/tags/apply", map[string]any{
		"revision_id": view["revision_id"], "candidate_ids": []any{tagCandidates[0].(map[string]any)["id"]},
	})
	require.Equal(t, http.StatusOK, applied.Code, applied.Body.String())
	links, err := rules.ListLinks(t.Context(), "inactive-rule")
	require.NoError(t, err)
	require.Len(t, links, 1)
	require.NotNil(t, links[0].TagID)
	require.NotNil(t, links[0].MappingID)
	var pointEnabled, mappingEnabled bool
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT p.enabled,m.enabled FROM points p JOIN mappings m ON m.point_id=p.id WHERE m.id=?`, *links[0].MappingID).Scan(&pointEnabled, &mappingEnabled))
	require.True(t, pointEnabled)
	require.True(t, mappingEnabled)
	f.record, err = f.workspace.ReadSnapshot(t.Context())
	require.NoError(t, err)
	request := f.createRequest()
	request["group"].(map[string]any)["members"] = []map[string]any{{"device_id": "device-A", "point_id": links[0].PointID, "tag_id": *links[0].TagID, "target_column": "known", "required": true}}
	for _, apply := range configure {
		apply(request["group"].(map[string]any))
	}
	created = performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, request)
	require.Equal(t, http.StatusCreated, created.Code, created.Body.String())
	data, group := groupSaveData(t, decodeJSONBody(t, created))
	return inactiveDraftFixture{writeGroupRouterFixture: f, link: links[0], data: data, group: group}
}

func (f inactiveDraftFixture) stopSource(t *testing.T) {
	t.Helper()
	response := performJSONRequest(t, f.router, http.MethodPut, inactiveDraftRulePath+"/inactive-rule", map[string]any{"device_id": "device-A", "enabled": false})
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	var pointEnabled, mappingEnabled bool
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT p.enabled,m.enabled FROM points p JOIN mappings m ON m.point_id=p.id WHERE m.id=?`, *f.link.MappingID).Scan(&pointEnabled, &mappingEnabled))
	require.False(t, pointEnabled)
	require.False(t, mappingEnabled)
}

func (f inactiveDraftFixture) editRequest(t *testing.T) map[string]any {
	t.Helper()
	record, err := f.workspace.ReadSnapshot(t.Context())
	require.NoError(t, err)
	// The ordinary editor submits descriptive intent, not stale source proof.
	// Mixed editor proof is separately covered and must never be discarded.
	for _, value := range f.group["members"].([]any) {
		member := value.(map[string]any)
		delete(member, "source_revision")
		delete(member, "mapping_revision")
	}
	request := groupEditRequest(f.writeGroupRouterFixture, f.data, f.group)
	request["expected_workspace_revision"] = record.DatabaseSetupRevision
	return request
}

func inactiveDraftState(t *testing.T, db *sql.DB) map[string][]string {
	t.Helper()
	state := lifecycleAPIState(t, db)
	rows, err := db.QueryContext(t.Context(), `SELECT * FROM write_group_members ORDER BY group_id,member_index`)
	require.NoError(t, err)
	columns, err := rows.Columns()
	require.NoError(t, err)
	for rows.Next() {
		values, pointers := make([]any, len(columns)), make([]any, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}
		require.NoError(t, rows.Scan(pointers...))
		payload, err := json.Marshal(values)
		require.NoError(t, err)
		state["write_group_members"] = append(state["write_group_members"], string(payload))
	}
	require.NoError(t, rows.Err())
	require.NoError(t, rows.Close())
	return state
}
