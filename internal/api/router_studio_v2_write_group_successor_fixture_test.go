package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

type successorAPIFixture struct {
	*wideRig
	request workspace.WriteGroupSuccessorReviewRequest
}

func newSuccessorAPIFixture(t *testing.T) successorAPIFixture {
	t.Helper()
	rig := newWideRig(t, "sqlite", []meanSource{{name: "pressure", tagType: "float32", role: "gauge"}})
	require.NoError(t, rig.pipe.Stop(time.Second))
	rig.f.groups.WithBacklogOwnershipGuard(grouppipeline.BacklogGuard{})
	current, err := rig.f.groups.Get(t.Context(), rig.groupID)
	require.NoError(t, err)
	raw, err := json.Marshal(current.Group)
	require.NoError(t, err)
	var target workspace.WriteGroup
	require.NoError(t, json.Unmarshal(raw, &target))
	target.Name = "Prepared successor"
	target.Destination.TableName = ""
	var lineage string
	var revision int64
	require.NoError(t, rig.f.db.QueryRow(`SELECT lineage_id,revision FROM write_group_lineages WHERE current_group_id=?`, rig.groupID).Scan(&lineage, &revision))
	return successorAPIFixture{wideRig: rig, request: workspace.WriteGroupSuccessorReviewRequest{
		WorkspaceID: rig.f.record.ID, ExpectedWorkspaceRevision: current.WorkspaceRevision,
		LineageID: lineage, ExpectedLineageRevision: revision, ExpectedPredecessorRevision: current.Group.Revision,
		ExpectedPredecessorConnectorRevision: "connector-1", ExpectedTargetConnectorRevision: "connector-1",
		OperationID: "successor-O2", Boundary: rig.snapshot.EffectiveAt.Add(time.Minute), Group: &target,
	}}
}

func (f successorAPIFixture) review(t *testing.T) workspace.WriteGroupSuccessorReview {
	t.Helper()
	response := performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+f.groupID+"/successor-review", f.request)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	var body struct {
		Data workspace.WriteGroupSuccessorReview `json:"data"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
	require.Empty(t, body.Data.Blockers, "%+v", body.Data.Blockers)
	return body.Data
}

func (f successorAPIFixture) reserve(t *testing.T, review workspace.WriteGroupSuccessorReview) workspace.WriteGroupSuccessorPreparation {
	t.Helper()
	response := performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+f.groupID+"/successor", map[string]any{
		"action": "reserve", "envelope": review.Envelope, "review_digest": review.ReviewDigest,
	})
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	var body struct {
		Data workspace.WriteGroupSuccessorPreparation `json:"data"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
	require.NotNil(t, body.Data.Group)
	return body.Data
}

func successorAPIStage(t *testing.T, f successorAPIFixture, review workspace.WriteGroupSuccessorReview, prepared workspace.WriteGroupSuccessorPreparation) map[string]any {
	t.Helper()
	current, err := f.f.groups.Get(t.Context(), prepared.Group.ID)
	require.NoError(t, err)
	predecessor, err := f.f.groups.Get(t.Context(), f.groupID)
	require.NoError(t, err)
	var lineageRevision int64
	require.NoError(t, f.f.db.QueryRow(`SELECT revision FROM write_group_lineages WHERE lineage_id=?`, prepared.LineageID).Scan(&lineageRevision))
	return map[string]any{"envelope": review.Envelope, "review_digest": review.ReviewDigest,
		"current": workspace.WriteGroupSuccessorCurrentCAS{
			WorkspaceRevision: current.WorkspaceRevision, LineageRevision: lineageRevision,
			PredecessorRevision: predecessor.Group.Revision, PredecessorConnectorRevision: predecessor.Group.Destination.ConnectorRevision,
			GroupID: current.Group.ID, GroupRevision: current.Group.Revision, TargetConnectorRevision: current.Group.Destination.ConnectorRevision,
		}}
}

// Capture every local SQL table, including preview tokens, ledger and lineage.
// Compare sorted row bytes so query order is not confused with mutation.
func successorAPIAllState(t *testing.T, db *sql.DB) map[string][]string {
	t.Helper()
	rows, err := db.Query(`SELECT name FROM sqlite_schema WHERE type='table' ORDER BY name`)
	require.NoError(t, err)
	var names []string
	for rows.Next() {
		var name string
		require.NoError(t, rows.Scan(&name))
		names = append(names, name)
	}
	require.NoError(t, rows.Err())
	require.NoError(t, rows.Close())
	state := map[string][]string{}
	for _, name := range names {
		rows, err := db.Query(`SELECT * FROM "` + strings.ReplaceAll(name, `"`, `""`) + `"`)
		require.NoError(t, err)
		columns, err := rows.Columns()
		require.NoError(t, err)
		state[name] = []string{}
		for rows.Next() {
			values := make([]any, len(columns))
			dest := make([]any, len(columns))
			for i := range values {
				dest[i] = &values[i]
			}
			require.NoError(t, rows.Scan(dest...))
			raw, err := json.Marshal(values)
			require.NoError(t, err)
			state[name] = append(state[name], string(raw))
		}
		require.NoError(t, rows.Err())
		require.NoError(t, rows.Close())
		sort.Strings(state[name])
	}
	return state
}

// A review must never ask either destination inspection interface to open SQL.
type successorAPIForbiddenInspector struct{ calls int }

func (i *successorAPIForbiddenInspector) InspectTable(context.Context, string, string, string) (*dbtarget.TableInspection, error) {
	i.calls++
	return nil, nil
}
func (i *successorAPIForbiddenInspector) InspectTableAtRevision(context.Context, string, string, string, string) (*dbtarget.TableInspection, error) {
	i.calls++
	return nil, nil
}
