package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRouter_WriteGroupInactiveDraftStillRequiresSavedSourceOwnership(t *testing.T) {
	for _, mutation := range []string{
		`DELETE FROM points WHERE id=?`,
		`DELETE FROM tags WHERE id=(SELECT tag_id FROM mappings WHERE point_id=?)`,
		`UPDATE points SET device_id='foreign-device' WHERE id=?`,
		`DELETE FROM source_rule_links WHERE point_id=?`,
		`UPDATE source_rule_links SET mapping_id='mapping-A' WHERE point_id=?`,
		`UPDATE source_rule_links SET tag_id='tag-A' WHERE point_id=?`,
		`UPDATE write_group_members SET mapping_revision='' WHERE point_id=?`,
	} {
		t.Run(mutation, func(t *testing.T) {
			f := newInactiveDraftFixture(t)
			f.stopSource(t)
			_, err := f.db.ExecContext(t.Context(), mutation, f.link.PointID)
			require.NoError(t, err)
			request := f.editRequest(t)
			before := inactiveDraftState(t, f.db)
			response := performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+f.group["id"].(string), request)
			require.Contains(t, []int{http.StatusNotFound, http.StatusConflict}, response.Code, response.Body.String())
			require.Equal(t, before, inactiveDraftState(t, f.db))
		})
	}
}

func TestNewRouter_WriteGroupInactiveDraftRejectsReparentedRule(t *testing.T) {
	f := newInactiveDraftFixture(t)
	f.stopSource(t)
	_, err := f.db.ExecContext(t.Context(), `INSERT INTO source_rules(id,device_id,start_address,count,data_type,naming_prefix,revision_id,enabled)
		VALUES('replacement-rule','device-A','41201',1,'int16','Other','replacement-revision',0)`)
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), `UPDATE source_rule_links SET rule_id='replacement-rule' WHERE point_id=?`, f.link.PointID)
	require.NoError(t, err)
	before := inactiveDraftState(t, f.db)
	response := performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+f.group["id"].(string), f.editRequest(t))
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
	require.Equal(t, before, inactiveDraftState(t, f.db))
}

func TestNewRouter_WriteGroupInactiveDraftStorageRejectsAdditionalRuleOwner(t *testing.T) {
	f := newInactiveDraftFixture(t)
	f.stopSource(t)
	_, err := f.db.ExecContext(t.Context(), `INSERT INTO source_rules(id,device_id,start_address,count,data_type,naming_prefix,revision_id,enabled)
		VALUES('additional-rule','device-A','41201',1,'int16','Other','additional-revision',0)`)
	require.NoError(t, err)
	before := inactiveDraftState(t, f.db)
	_, err = f.db.ExecContext(t.Context(), `INSERT INTO source_rule_links(id,rule_id,address,point_id,tag_id,mapping_id)
		VALUES('additional-link','additional-rule','41201',?,?,?)`, f.link.PointID, *f.link.TagID, *f.link.MappingID)
	require.ErrorContains(t, err, "UNIQUE constraint failed: source_rule_links.point_id")
	require.Equal(t, before, inactiveDraftState(t, f.db))
	var owners int
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM source_rule_links WHERE point_id=?`, f.link.PointID).Scan(&owners))
	require.Equal(t, 1, owners)
}

func TestNewRouter_WriteGroupInactiveDraftDoesNotTrustExplicitStaleProof(t *testing.T) {
	for _, proof := range []string{"source_revision", "mapping_revision"} {
		t.Run(proof, func(t *testing.T) {
			f := newInactiveDraftFixture(t)
			previous := f.group["members"].([]any)[0].(map[string]any)[proof]
			f.stopSource(t)
			request := f.editRequest(t)
			f.group["members"].([]any)[0].(map[string]any)[proof] = previous
			before := inactiveDraftState(t, f.db)
			response := performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+f.group["id"].(string), request)
			require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
			assertWriteGroupError(t, decodeJSONBody(t, response), "revision_mismatch")
			require.Equal(t, before, inactiveDraftState(t, f.db))
		})
	}
}

func TestNewRouter_WriteGroupInactiveMixedDraftRetainsExplicitProof(t *testing.T) {
	f := newInactiveDraftFixture(t, func(group map[string]any) {
		group["record_id"] = "6b2669b5-02df-455a-ab6f-0614438257e8"
		group["layout_version"] = "managed_wide_profile_v1"
		group["recording_profile"] = map[string]any{"kind": "mixed", "result_shape": "member-calculation-v1"}
		group["row_policy"] = map[string]any{"incomplete_policy": "partial"}
		group["destination"].(map[string]any)["storage_strategy"] = "managed"
		member := group["members"].([]map[string]any)[0]
		member["calculation"], member["value_role"] = "mean", "gauge"
		member["cadence"] = map[string]any{"interval": 15, "unit": "second", "alignment": "utc_epoch"}
	})
	f.stopSource(t)
	current, err := f.workspace.ReadSnapshot(t.Context())
	require.NoError(t, err)
	request := groupEditRequest(f.writeGroupRouterFixture, f.data, f.group)
	request["expected_workspace_revision"] = current.DatabaseSetupRevision
	f.group["name"] = "Mixed stale proof must not be refreshed"
	before := inactiveDraftState(t, f.db)
	updated := performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+f.group["id"].(string), request)
	require.Equal(t, http.StatusConflict, updated.Code, updated.Body.String())
	assertWriteGroupError(t, decodeJSONBody(t, updated), "revision_mismatch")
	require.Equal(t, before, inactiveDraftState(t, f.db))
}

func TestNewRouter_WriteGroupInactiveSourceCannotBecomeNewMemberOrGroup(t *testing.T) {
	f := newInactiveDraftFixture(t)
	f.stopSource(t)
	record, err := f.workspace.ReadSnapshot(t.Context())
	require.NoError(t, err)
	f.record = record
	request := f.createRequest()
	created := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, request)
	require.Equal(t, http.StatusCreated, created.Code, created.Body.String())
	data, group := groupSaveData(t, decodeJSONBody(t, created))
	inactiveMembers := f.editRequest(t)["group"].(map[string]any)["members"]
	group["members"] = inactiveMembers
	before := inactiveDraftState(t, f.db)
	updated := performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+group["id"].(string), groupEditRequest(f.writeGroupRouterFixture, data, group))
	require.Equal(t, http.StatusNotFound, updated.Code, updated.Body.String())
	request = f.createRequest()
	request["expected_workspace_revision"] = data["workspace_revision"]
	request["group"].(map[string]any)["members"] = inactiveMembers
	created = performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, request)
	require.Equal(t, http.StatusNotFound, created.Code, created.Body.String())
	require.Equal(t, before, inactiveDraftState(t, f.db))
}

func TestNewRouter_WriteGroupInactiveDraftRenamePreservesAppliedWriter(t *testing.T) {
	f := lifecycleAPIApplied(t)
	id := f.group["id"].(string)
	member := f.group["members"].([]any)[0].(map[string]any)
	_, err := f.db.ExecContext(t.Context(), `UPDATE mappings SET enabled=0 WHERE point_id=? AND tag_id=?`, member["point_id"], member["tag_id"])
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), `UPDATE points SET enabled=0 WHERE id=?`, member["point_id"])
	require.NoError(t, err)
	before := inactiveDraftState(t, f.db)
	previousApplied := f.group["applied_revision"]
	f.group["name"] = "Inactive source, unchanged writer"
	delete(member, "source_revision")
	delete(member, "mapping_revision")
	updated := performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+id, groupEditRequest(f.writeGroupRouterFixture, f.data, f.group))
	require.Equal(t, http.StatusOK, updated.Code, updated.Body.String())
	_, group := groupSaveData(t, decodeJSONBody(t, updated))
	require.Equal(t, "draft", group["status"])
	require.Equal(t, previousApplied, group["applied_revision"])
	after := inactiveDraftState(t, f.db)
	for _, table := range []string{"write_group_versions", "wg_runtime_versions", "wg_delivery_samples", "wg_delivery_outbox", "wg_delivery_receipts"} {
		require.Equal(t, before[table], after[table], table)
	}
}
