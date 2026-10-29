package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/workspace"
)

func TestSQLServerSuccessorUnsupportedShapesHaveNoEffects(t *testing.T) {
	sqlServerSuccessorFixtures(t) // Missing/unsafe real fixtures must mark the entrypoint NOT_RUN.
	for _, family := range []string{"sqlserver"} {
		t.Run(family, func(t *testing.T) {
			for _, shape := range []string{"wide", "cross_mode", "custom", "partial", "meaning", "cross_engine"} {
				t.Run(shape, func(t *testing.T) {
					r := newSQLServerSuccessorRig(t, "plain", false, false, meanSource{name: "voltage", tagType: "float32", role: "gauge", required: true})
					switch shape {
					case "wide":
						r.request.Group.LayoutVersion = workspace.WriteGroupLayoutManagedWideProfileV1
						r.request.Group.StorageMode = ""
					case "cross_mode":
						r.request.Group.StorageMode = "native_monthly"
					case "custom":
						r.request.Group.Destination.StorageStrategy = workspace.WriteGroupStorageStrategyCustom
					case "meaning":
						r.request.Group.RecordingProfile.SampleReduction = workspace.WriteGroupSampleReductionLatest
					case "cross_engine":
						_, err := r.f.db.Exec(`UPDATE database_connectors SET kind='postgres' WHERE id='connector-B'`)
						require.NoError(t, err)
					case "partial":
						require.Len(t, r.request.Group.Members, 2)
						r.request.Group.Members = r.request.Group.Members[:1]
					}
					before, external := sqlServerSuccessorLocalState(t, r.f.db), r.targetState(t)
					response := performJSONRequest(t, r.f.router, http.MethodPost, writeGroupsPath+"/"+r.groupID+"/successor-review", r.request)
					if response.Code == http.StatusOK {
						var body struct {
							Data workspace.WriteGroupSuccessorReview `json:"data"`
						}
						require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
						require.NotEmpty(t, body.Data.Blockers)
					} else {
						require.GreaterOrEqual(t, response.Code, 400)
						require.Less(t, response.Code, 500)
					}
					require.Equal(t, before, sqlServerSuccessorLocalState(t, r.f.db))
					require.Equal(t, external, r.targetState(t))
				})
			}
		})
	}
}
