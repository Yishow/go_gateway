package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"
)

func TestSQLServerSuccessorLocalReviewRejectsCorruptFullProof(t *testing.T) {
	sqlServerSuccessorFixtures(t) // Missing/unsafe real fixtures must mark the entrypoint NOT_RUN.
	for _, mode := range []workspace.WriteGroupStorageMode{"plain", "native_monthly"} {
		for _, fault := range []string{"owner", "database", "source_digest", "capability", "partition", "mode_exchange", "operation_digest", "token_digest"} {
			t.Run(string(mode)+"/"+fault, func(t *testing.T) {
				r := newSQLServerSuccessorRig(t, mode, false, true)
				var tokenID, raw string
				require.NoError(t, r.f.db.QueryRow(`SELECT token,tables FROM managed_schema_preview_tokens WHERE plan_id=?`, r.prepared.Group.ID).Scan(&tokenID, &raw))
				var body map[string]any
				require.NoError(t, json.Unmarshal([]byte(raw), &body))
				layout := body["group_layout"].(map[string]any)
				options := layout["sqlserver"].(map[string]any)
				switch fault {
				case "owner":
					layout["owner_column"] = "_gw_owner_000000000000000000000000"
				case "database":
					options["database"] = r.fixtures[0].Manifest.Database
				case "source_digest":
					body["source_digest"] = "foreign-source"
				case "capability":
					options["capability_revision"] = ""
				case "partition":
					if mode == "native_monthly" {
						options["native_monthly"].(map[string]any)["function_name"] = "pf_foreign"
					} else {
						options["compatibility_level"] = 150
					}
				case "mode_exchange":
					if mode == "native_monthly" {
						delete(options, "native_monthly")
					} else {
						// Use a valid canonical native payload from the ordinary builder,
						// but never claim it is this plain scope's confirmed authority.
						table := layout["table_name"].(string)
						fn, scheme, err := recordingplan.SQLServerNativeMonthlyObjectNames(table)
						require.NoError(t, err)
						native, err := recordingplan.NewSQLServerNativeMonthlyPartition(r.boundary, fn, scheme)
						require.NoError(t, err)
						options["native_monthly"] = native
					}
				case "operation_digest":
					_, err := r.f.db.Exec(`UPDATE managed_schema_operations SET verified_digest='foreign' WHERE token=?`, tokenID)
					require.NoError(t, err)
				case "token_digest":
					_, err := r.f.db.Exec(`UPDATE managed_schema_preview_tokens SET digest='foreign' WHERE token=?`, tokenID)
					require.NoError(t, err)
				}
				updated, err := json.Marshal(body)
				require.NoError(t, err)
				_, err = r.f.db.Exec(`UPDATE managed_schema_preview_tokens SET tables=? WHERE token=?`, string(updated), tokenID)
				require.NoError(t, err)
				spy := &sqlServerSuccessorProbeSpy{ManagedTableInspector: r.f.groups.ManagedTableInspector(r.targets), offline: true}
				r.f.groups.WithManagedTableInspector(spy)
				before, external := sqlServerSuccessorLocalState(t, r.f.db), r.targetState(t)
				response := performJSONRequest(t, r.f.router, http.MethodPost, writeGroupsPath+"/"+r.groupID+"/successor", r.activationAPI().activationRequest(t))
				if response.Code == http.StatusOK {
					require.NotEmpty(t, decodeJSONBody(t, response)["data"].(map[string]any)["blockers"])
				} else {
					require.GreaterOrEqual(t, response.Code, 400)
				}
				require.Zero(t, spy.calls, "corrupt stored proof cannot fall back to external inspection")
				require.Equal(t, before, sqlServerSuccessorLocalState(t, r.f.db))
				require.Equal(t, external, r.targetState(t))
			})
		}
	}
}

func TestSQLServerSuccessorSchemaProofReplayOffline(t *testing.T) {
	sqlServerSuccessorFixtures(t) // Missing/unsafe real fixtures must mark the entrypoint NOT_RUN.
	for _, mode := range []workspace.WriteGroupStorageMode{"plain", "native_monthly"} {
		t.Run(string(mode), func(t *testing.T) {
			r := newSQLServerSuccessorRig(t, mode, false, false)
			r.preparationReview = r.review(t)
			r.prepared = r.reserve(t, r.preparationReview)
			token := successorAPIPreview(t, r.successorAPIFixture, r.preparationReview, r.prepared)
			request := successorAPIConfirmRequest(t, r.successorAPIFixture, r.preparationReview, r.prepared, token)
			f := r.f
			f.target, f.namespace = r.targetDB, r.targetNamespace
			registerSQLServerSuccessorPreparationAudit(t, f, token)
			response := performJSONRequest(t, r.f.router, http.MethodPost, writeGroupsPath+"/"+r.groupID+"/successor", request)
			require.Equal(t, http.StatusOK, response.Code, response.Body.String())
			var outcome struct {
				Data recordingplan.SchemaOperation `json:"data"`
			}
			require.NoError(t, json.Unmarshal(response.Body.Bytes(), &outcome))
			require.Equal(t, recordingplan.SchemaOperationSucceeded, outcome.Data.Status)
			registerSQLServerSuccessorSchemaCleanup(t, f, token)
			spy := &sqlServerSuccessorProbeSpy{ManagedTableInspector: r.f.groups.ManagedTableInspector(r.targets), offline: true}
			r.f.groups.WithManagedTableInspector(spy)
			before, external := sqlServerSuccessorLocalState(t, r.f.db), r.targetState(t)
			replay := performJSONRequest(t, r.f.router, http.MethodPost, writeGroupsPath+"/"+r.groupID+"/successor", request)
			require.Equal(t, http.StatusOK, replay.Code, replay.Body.String())
			require.Equal(t, outcome.Data.OperationID, decodeJSONBody(t, replay)["data"].(map[string]any)["operation_id"])
			require.Zero(t, spy.calls, "exact schema confirmation replay cannot re-probe or run DDL")
			require.Equal(t, before, sqlServerSuccessorLocalState(t, r.f.db))
			require.Equal(t, external, r.targetState(t))
		})
	}
}
