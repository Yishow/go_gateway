package api

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/workspace"
	"net/http"
	"testing"
)

func TestMySQLSuccessorRejectsActualOtherFamilyBeforeReservation(t *testing.T) {
	for _, family := range []string{"mysql", "mariadb"} {
		t.Run(family, func(t *testing.T) {
			other := "mariadb"
			if family == other {
				other = "mysql"
			}
			for _, mode := range []workspace.WriteGroupStorageMode{"plain", "native_monthly"} {
				t.Run(string(mode), func(t *testing.T) {
					r := newMySQLSuccessorRig(t, family, other, mode, false, false)
					review := r.review(t)
					before, external := successorAPIAllState(t, r.f.db), r.targetState(t)
					response := performJSONRequest(t, r.f.router, http.MethodPost, writeGroupsPath+"/"+r.groupID+"/successor", map[string]any{"action": "reserve", "envelope": review.Envelope, "review_digest": review.ReviewDigest})
					require.GreaterOrEqual(t, response.Code, 400, response.Body.String())
					require.Less(t, response.Code, 500, response.Body.String())
					require.Equal(t, http.StatusUnprocessableEntity, response.Code)
					require.Contains(t, response.Body.String(), "WRITE_GROUP_INVALID")
					// HTTP intentionally redacts detailed structural reasons. Verify the same
					// real public service path returns the precise typed family refusal too.
					_, err := r.f.groups.ReserveSuccessor(t.Context(), r.groupID, workspace.WriteGroupSuccessorReservationRequest{Envelope: review.Envelope, ReviewDigest: review.ReviewDigest})
					var rejected *workspace.WriteGroupTransitionError
					require.ErrorAs(t, err, &rejected)
					require.Equal(t, "cross-engine", rejected.ReasonCode)
					require.Equal(t, before, successorAPIAllState(t, r.f.db), "wrong actual family cannot trap lineage with reservation")
					require.Equal(t, external, r.targetState(t), "wrong actual family cannot run DDL")
				})
			}
		})
	}
}

func TestMySQLSuccessorUnsupportedShapesHaveNoEffects(t *testing.T) {
	for _, family := range []string{"mysql", "mariadb"} {
		t.Run(family, func(t *testing.T) {
			for _, shape := range []string{"wide", "cross_mode", "custom", "partial"} {
				t.Run(shape, func(t *testing.T) {
					r := newMySQLSuccessorRig(t, family, family, "plain", false, false, meanSource{name: "voltage", tagType: "float32", role: "gauge", required: true})
					switch shape {
					case "wide":
						r.request.Group.LayoutVersion = workspace.WriteGroupLayoutManagedWideProfileV1
						r.request.Group.StorageMode = ""
					case "cross_mode":
						r.request.Group.StorageMode = "native_monthly"
					case "custom":
						r.request.Group.Destination.StorageStrategy = workspace.WriteGroupStorageStrategyCustom
					case "partial":
						require.Len(t, r.request.Group.Members, 2)
						r.request.Group.Members = r.request.Group.Members[:1]
					}
					before, external := successorAPIAllState(t, r.f.db), r.targetState(t)
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
					require.Equal(t, before, successorAPIAllState(t, r.f.db))
					require.Equal(t, external, r.targetState(t))
				})
			}
		})
	}
}
