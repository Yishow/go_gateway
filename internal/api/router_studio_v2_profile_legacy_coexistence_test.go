package api

import (
	"net/http"
	"slices"
	"testing"
	"time"

	"go-gateway/internal/datalink/recordingplan"

	"github.com/stretchr/testify/require"
)

func TestProfileLegacyCoexistence(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		t.Run(kind, func(t *testing.T) {
			r := newCoexistRig(t, kind)
			r.holdPendingRows()
			names := []string{"wide", "custom", "energy", "voltage"}
			baseline := map[string]string{}
			for _, name := range names {
				baseline[name] = r.ground(name)
			}
			unchangedExcept := func(t *testing.T, changed ...string) {
				t.Helper()
				for _, name := range names {
					if !slices.Contains(changed, name) {
						require.JSONEq(t, baseline[name], r.ground(name), "%s must not be touched by another group's lifecycle", name)
					}
				}
			}

			t.Run("each writer keeps its own layout, period and frozen state", func(t *testing.T) {
				wide, custom := getGroup(t, r.f, r.groups["wide"]), getGroup(t, r.f, r.groups["custom"])
				require.Contains(t, []any{nil, "", "managed_wide_v1"}, wide["layout_version"])
				require.Nil(t, wide["recording_profile"])
				require.EqualValues(t, 60, wide["row_policy"].(map[string]any)["interval_seconds"])
				require.Equal(t, "custom", destinationOf(custom)["storage_strategy"])
				require.Equal(t, "coexist_custom", destinationOf(custom)["table_name"])
				require.Nil(t, custom["recording_profile"])
				for _, name := range []string{"energy", "voltage"} {
					group := getGroup(t, r.f, r.groups[name])
					require.Equal(t, "managed_long_v1", group["layout_version"])
					requireProfile(t, group, "instantaneous", 15, "second", 15)
				}
				require.Equal(t, destinationOf(getGroup(t, r.f, r.groups["energy"]))["table_name"], destinationOf(getGroup(t, r.f, r.groups["voltage"]))["table_name"], "two profile circuits share one long table")
				for _, name := range names {
					group := getGroup(t, r.f, r.groups[name])
					require.NotEmpty(t, group["applied_revision"], name)
					destination := destinationOf(group)
					var ownPending int
					require.NoError(t, r.f.db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM wg_delivery_outbox WHERE group_id = ? AND connector_id = ? AND connector_revision = ? AND table_name = ? AND state = 'pending'`,
						r.groups[name], destination["connector_id"], destination["connector_revision"], destination["table_name"]).Scan(&ownPending))
					require.Positive(t, ownPending, "%s has a pending effect frozen to its own destination", name)
				}
				require.GreaterOrEqual(t, r.pendingRows(), len(names), "every writer has its own pending row for its own destination")
			})

			t.Run("a long profile on a source another writer owns is still refused until a reviewed cutover", func(t *testing.T) {
				for _, owner := range []string{"wide", "custom"} {
					source := map[string]string{"wide": "", "custom": "custom"}[owner]
					blocked := r.saveProfile("c-over-"+owner, source)
					require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, r.f, blocked, previewLong(t, r.f, blocked)).Status)
					code, body := r.apply(blocked)
					require.Equal(t, http.StatusConflict, code, body)
					require.Equal(t, "cutover_required", errorOf(t, body)["code"], owner)
					require.Equal(t, "", getGroup(t, r.f, blocked)["applied_revision"])
					var versions int
					require.NoError(t, r.f.db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM write_group_versions WHERE group_id = ?`, blocked).Scan(&versions))
					require.Zero(t, versions, "nothing was applied for the overlapping group")
				}
				unchangedExcept(t)
			})

			var energyOperation string
			t.Run("a new member of a shared long table needs its own confirmation and no DDL", func(t *testing.T) {
				energy := getGroup(t, r.f, r.groups["energy"])
				energyOperation = previewLong(t, r.f, r.groups["energy"]).OperationID
				before := columnsOf(t, r.dest, kind, r.f.namespace, destinationOf(energy)["table_name"].(string))
				voltageRevision := destinationOf(getGroup(t, r.f, r.groups["voltage"]))["schema_revision"]
				updated := cloneMap(energy)
				updated["members"] = append(updated["members"].([]any), map[string]any{
					"device_id": "device-A", "point_id": "point-extra", "tag_id": "tag-extra", "required": true, "circuit_id": "c-energy", "circuit_name": "Circuit c-energy",
				})
				response := performJSONRequest(t, r.f.router, http.MethodPut, writeGroupsPath+"/"+r.groups["energy"], groupEditRequest(r.f.writeGroupRouterFixture, getGroupData(t, r.f, r.groups["energy"]), updated))
				require.Equal(t, http.StatusOK, response.Code, response.Body.String())

				grown := getGroup(t, r.f, r.groups["energy"])
				require.Empty(t, destinationOf(grown)["schema_revision"], "the changed group needs a fresh confirmation")
				require.Equal(t, energy["applied_revision"], grown["applied_revision"], "its applied snapshot stays as it was")
				ready, err := r.f.groups.Readiness(t.Context(), r.groups["energy"])
				require.NoError(t, err)
				require.Contains(t, issueCodes(ready.Issues), "schema-confirmation-required")
				require.Equal(t, voltageRevision, destinationOf(getGroup(t, r.f, r.groups["voltage"]))["schema_revision"], "the other profile keeps its own confirmation")

				token := previewLong(t, r.f, r.groups["energy"])
				require.Empty(t, token.Statements, "a new point adds rows, never columns")
				require.Equal(t, recordingplan.NoChangeSchemaCompatible, token.NoChangeReason)
				require.NotEqual(t, energyOperation, token.OperationID, "a compatible preview is still a fresh operation")
				confirmed := confirmLong(t, r.f, r.groups["energy"], token)
				require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmed.Status)
				require.Zero(t, confirmed.ExecutedStatements)
				require.Equal(t, before, columnsOf(t, r.dest, kind, r.f.namespace, destinationOf(energy)["table_name"].(string)))
				unchangedExcept(t, "energy")
			})

			t.Run("stale workspace, group and connector versions are refused with no DDL while the applied writers stay valid", func(t *testing.T) {
				token := previewLong(t, r.f, r.groups["voltage"])
				tables := r.tableCount()
				for _, stale := range []struct{ name, field string }{
					{"workspace", "expected_workspace_revision"}, {"group", "expected_group_revision"}, {"connector", "expected_connector_revision"},
				} {
					for _, route := range []string{"schema-preview", "schema-apply"} {
						request := scopeBody(t, r.f, r.groups["voltage"])
						request[stale.field] = "stale-" + stale.name
						if route == "schema-apply" {
							request["token"], request["operation_id"] = token.Token, token.OperationID
						}
						code, body := r.post(writeGroupsPath+"/"+r.groups["voltage"]+"/"+route, request)
						require.Equal(t, http.StatusConflict, code, "%s %s: %v", stale.name, route, body)
					}
				}
				require.Equal(t, tables, r.tableCount(), "a refused request created nothing")
				for _, name := range names {
					snapshot, err := r.f.groups.ResolveAppliedAt(t.Context(), r.groups[name], time.Now().Add(time.Hour))
					require.NoError(t, err, name)
					require.Equal(t, getGroup(t, r.f, r.groups[name])["applied_revision"], snapshot.AppliedRevision, "the applied snapshot of %s is still in force", name)
					eligibility, err := r.f.groups.IntakeEligibility(t.Context(), r.groups[name])
					require.NoError(t, err)
					// A group applied moments ago may still be waiting for its first window; none may be ended, disabled or deleted.
					require.True(t, eligibility.Accepting || eligibility.Reason == "pending-activation", "%s was not stopped by a refused request: %+v", name, eligibility)
					require.Equal(t, getGroup(t, r.f, r.groups[name])["status"], string(eligibility.Status), name)
				}
				unchangedExcept(t, "energy")
			})
		})
	}
}
