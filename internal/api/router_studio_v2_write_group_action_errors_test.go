package api

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteGroupRemoveActionErrorsPreserveCompatibilityAndState(t *testing.T) {
	for _, tc := range []struct {
		name, reason, code string
		status             int
	}{
		{"changed", "connector_identity_stale", "revision_mismatch", http.StatusConflict},
		{"removed", "connector_unavailable", "WRITE_GROUP_NOT_FOUND", http.StatusNotFound},
		{"no_guard", "ownership_unavailable", "WRITE_GROUP_LIFECYCLE_BLOCKED", http.StatusConflict},
		{"failed_guard", "ownership_unavailable", "WRITE_GROUP_LIFECYCLE_BLOCKED", http.StatusConflict},
	} {
		for _, stopped := range []bool{false, true} {
			name := tc.name + "/running"
			if stopped {
				name = tc.name + "/stopped"
			}
			t.Run(name, func(t *testing.T) {
				f := lifecycleAPIApplied(t)
				id := f.group["id"].(string)
				f.groups.WithBacklogOwnershipGuard(&lifecycleAPIGuard{})
				if stopped {
					response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/disable", scopeBody(t, f, id))
					require.Equal(t, http.StatusOK, response.Code, response.Body.String())
				}
				switch tc.name {
				case "changed":
					_, err := f.db.Exec(`UPDATE database_connectors SET identity_revision='connector-2', connection_config='{"password":"private-password"}' WHERE id='connector-A'`)
					require.NoError(t, err)
				case "removed":
					_, err := f.db.Exec(`DELETE FROM database_connectors WHERE id='connector-A'`)
					require.NoError(t, err)
				case "no_guard":
					f.groups.WithBacklogOwnershipGuard(nil)
				case "failed_guard":
					f.groups.WithBacklogOwnershipGuard(&lifecycleAPIGuard{err: errors.New("private-proof private-dsn private-password")})
				}
				request := scopeBody(t, f, id)
				before := lifecycleAPIState(t, f.db)
				response := performJSONRequest(t, f.router, http.MethodDelete, writeGroupsPath+"/"+id, request)
				require.Equal(t, tc.status, response.Code, response.Body.String())
				body := decodeJSONBody(t, response)
				assertWriteGroupError(t, body, tc.code)
				envelope := body["error"].(map[string]any)
				require.Equal(t, tc.reason, envelope["reason_code"])
				next, legacy := "stop", "disable_group"
				if stopped {
					next, legacy = "review_connector", "review_group"
					if tc.reason == "ownership_unavailable" {
						next = "inspect_ownership"
					}
				}
				require.Equal(t, next, envelope["next_action"])
				require.Equal(t, legacy, envelope["action"])
				require.Equal(t, false, envelope["retryable"])
				require.NotContains(t, response.Body.String(), "private-")
				require.Equal(t, before, lifecycleAPIState(t, f.db))
				if !stopped {
					response = performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/disable", request)
					require.Equal(t, http.StatusOK, response.Code, response.Body.String())
				}
			})
		}
	}
}

func TestWriteGroupRemoveStaleCASDoesNotOfferStopAsABypass(t *testing.T) {
	f := lifecycleAPIApplied(t)
	id := f.group["id"].(string)
	for _, field := range []string{"expected_workspace_revision", "expected_group_revision", "expected_connector_revision"} {
		t.Run(field, func(t *testing.T) {
			f.groups.WithBacklogOwnershipGuard(&lifecycleAPIGuard{})
			request := scopeBody(t, f, id)
			request[field] = "stale-revision"
			before := lifecycleAPIState(t, f.db)
			for _, operation := range []struct{ method, suffix string }{{http.MethodDelete, ""}, {http.MethodPost, "/disable"}} {
				response := performJSONRequest(t, f.router, operation.method, writeGroupsPath+"/"+id+operation.suffix, request)
				require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
				body := decodeJSONBody(t, response)
				assertWriteGroupError(t, body, "revision_mismatch")
				envelope := body["error"].(map[string]any)
				require.Equal(t, "reload", envelope["action"])
				require.NotContains(t, envelope, "reason_code")
				require.NotContains(t, envelope, "next_action")
				require.Equal(t, before, lifecycleAPIState(t, f.db))
			}
		})
	}
}
