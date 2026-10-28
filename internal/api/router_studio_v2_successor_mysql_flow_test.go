package api

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"
	"net/http"
	"testing"
	"time"
)

func TestMySQLSuccessorPublicFlow(t *testing.T) {
	for _, family := range []string{"mysql", "mariadb"} {
		t.Run(family, func(t *testing.T) {
			for _, mode := range []workspace.WriteGroupStorageMode{"plain", "native_monthly"} {
				t.Run(string(mode), func(t *testing.T) {
					r := newMySQLSuccessorRig(t, family, family, mode, false, false)
					before := successorAPIAllState(t, r.f.db)
					review := r.review(t)
					require.Equal(t, before, successorAPIAllState(t, r.f.db), "local review is read-only")
					r.preparationReview = review
					r.prepared = r.reserve(t, review)
					token := successorAPIPreview(t, r.successorAPIFixture, review, r.prepared)
					response := performJSONRequest(t, r.f.router, http.MethodPost, writeGroupsPath+"/"+r.groupID+"/successor", successorAPIConfirmRequest(t, r.successorAPIFixture, review, r.prepared, token))
					require.Equal(t, http.StatusOK, response.Code, response.Body.String())
					var outcome struct {
						Data recordingplan.SchemaOperation `json:"data"`
					}
					require.NoError(t, json.Unmarshal(response.Body.Bytes(), &outcome))
					require.Equal(t, recordingplan.SchemaOperationSucceeded, outcome.Data.Status)
					f := r.activationAPI()
					spy := &mysqlSuccessorProbeSpy{ManagedTableInspector: r.f.groups.ManagedTableInspector(r.targets), offline: true}
					r.f.groups.WithManagedTableInspector(spy)
					activationLocal := successorAPIAllState(t, r.f.db)
					localReview := f.postActivation(t, f.activationRequest(t))
					require.Empty(t, localReview["blockers"])
					require.Zero(t, spy.calls, "activation review must use confirmed durable token without network")
					require.Equal(t, activationLocal, successorAPIAllState(t, r.f.db))
					spy.offline = false
					local := successorAPIAllState(t, r.f.db)
					external := r.targetState(t)
					confirmation := f.activationConfirmation(t)
					require.Equal(t, local, successorAPIAllState(t, r.f.db), "review/readiness read-only")
					event := f.postActivation(t, confirmation)
					require.NotEmpty(t, event["id"])
					next, err := r.f.groups.ResolveAppliedAt(t.Context(), r.prepared.Group.ID, r.boundary)
					require.NoError(t, err)
					require.Equal(t, "connector-B", next.Group.Destination.ConnectorID)
					require.Equal(t, r.targetNamespace, next.Group.Destination.Database)
					require.Equal(t, mode, next.Group.StorageMode)
					runtimeBytes, err := r.store.RuntimeVersion(t.Context(), groupdelivery.GroupKey{WorkspaceID: next.WorkspaceID, GroupID: next.GroupID, GroupRevision: next.GroupRevision})
					require.NoError(t, err)
					var runtimeSnapshot workspace.WriteGroupAppliedSnapshot
					require.NoError(t, json.Unmarshal(runtimeBytes, &runtimeSnapshot))
					require.NoError(t, workspace.ValidateMySQLRuntimeLayout(runtimeSnapshot.Group, runtimeSnapshot.RuntimeLayout))
					require.Equal(t, family, runtimeSnapshot.RuntimeLayout.MySQL.Group.Evidence.ServerFamily)
					_, err = r.f.groups.ResolveAppliedAt(t.Context(), r.prepared.Group.ID, r.boundary.Add(-time.Nanosecond))
					require.ErrorIs(t, err, workspace.ErrWriteGroupLifecycleBlocked)
					after := successorAPIAllState(t, r.f.db)
					require.Equal(t, local["write_group_storage_scopes"], after["write_group_storage_scopes"])
					require.Equal(t, event["id"], f.postActivation(t, confirmation)["id"])
					require.Equal(t, after, successorAPIAllState(t, r.f.db))
					raw, err := json.Marshal(confirmation)
					require.NoError(t, err)
					var changed map[string]any
					require.NoError(t, json.Unmarshal(raw, &changed))
					changed["envelope"].(map[string]any)["request"].(map[string]any)["boundary"] = r.boundary.Add(time.Second)
					rejected := performJSONRequest(t, r.f.router, http.MethodPost, writeGroupsPath+"/"+r.groupID+"/successor", changed)
					require.GreaterOrEqual(t, rejected.Code, 400)
					require.Less(t, rejected.Code, 500)
					require.Equal(t, after, successorAPIAllState(t, r.f.db))
					require.Equal(t, external, r.targetState(t))
				})
			}
		})
	}
}
