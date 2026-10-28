package api

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteGroupLifecycleAppliedDraftAndUnknownOwner(t *testing.T) {
	f := lifecycleAPIApplied(t)
	id := f.group["id"].(string)
	applied := f.group["applied_revision"].(string)
	require.NotEmpty(t, applied)
	require.EqualValues(t, 15, f.group["row_policy"].(map[string]any)["interval_seconds"])
	var immutable string
	require.NoError(t, f.db.QueryRow(`SELECT payload FROM wg_runtime_versions WHERE group_id=? AND group_revision=?`, id, applied).Scan(&immutable))
	changed := cloneMap(f.group)
	changed["row_policy"].(map[string]any)["interval_seconds"] = 30
	response := performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+id, groupEditRequest(f.writeGroupRouterFixture, f.data, changed))
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	savedData, saved := groupSaveData(t, decodeJSONBody(t, response))
	data := lifecycleAPIRead(t, f.writeGroupRouterFixture, id)
	require.Equal(t, saved["revision"], data["group_revision"])
	require.Equal(t, applied, data["applied_revision"])
	require.NotEqual(t, data["group_revision"], data["applied_revision"])
	require.Equal(t, savedData["workspace_revision"], data["workspace_revision"])
	var after string
	require.NoError(t, f.db.QueryRow(`SELECT payload FROM wg_runtime_versions WHERE group_id=? AND group_revision=?`, id, applied).Scan(&after))
	require.Equal(t, immutable, after)
	for _, mode := range []string{"nil", "failure", "proven"} {
		t.Run(mode, func(t *testing.T) {
			guard := &lifecycleAPIGuard{}
			if mode == "nil" {
				f.groups.WithBacklogOwnershipGuard(nil)
			} else {
				if mode == "failure" {
					guard.err = errors.New("private-sql-error secret-dsn")
				}
				f.groups.WithBacklogOwnershipGuard(guard)
			}
			before := lifecycleAPIState(t, f.db)
			result := lifecycleAPIRead(t, f.writeGroupRouterFixture, id)
			actions := lifecycleAPIActions(t, result)
			require.Equal(t, true, actions["stop"]["allowed"])
			require.Equal(t, mode == "proven", actions["remove"]["allowed"])
			if mode != "proven" {
				require.Equal(t, "ownership_unavailable", actions["remove"]["reason_code"])
				require.Equal(t, "stop", actions["remove"]["next_action"])
			}
			if mode != "nil" {
				require.Equal(t, 1, guard.calls)
			}
			require.Equal(t, before, lifecycleAPIState(t, f.db))
			read := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+id+"/lifecycle", nil)
			require.NotContains(t, read.Body.String(), "private-sql")
			require.NotContains(t, read.Body.String(), "secret-dsn")
		})
	}
	// A successful action read is not a reusable mutation grant: old CAS fails.
	stale := performJSONRequest(t, f.router, http.MethodDelete, writeGroupsPath+"/"+id, groupEditRequest(f.writeGroupRouterFixture, f.data, f.group))
	require.Equal(t, http.StatusConflict, stale.Code, stale.Body.String())
	current := lifecycleAPIRead(t, f.writeGroupRouterFixture, id)
	require.Equal(t, saved["revision"], current["group_revision"])
}
