package api

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

type fakeDeliveryReader struct {
	view *grouppipeline.DeliveryView
	err  error
}

func (f fakeDeliveryReader) Delivery(context.Context, string) (*grouppipeline.DeliveryView, error) {
	return f.view, f.err
}

func (f writeGroupRouterFixture) routerWithDelivery(reader fakeDeliveryReader) http.Handler {
	return NewRouter(&DatalinkServices{
		Workspace: f.workspace, WriteGroups: workspace.NewWriteGroupService(f.workspace, workspace.NewSQLWriteGroupRepository(f.db)),
		WriteGroupDelivery: reader,
	})
}

func TestRevisionBoundBacklogDeliveryRouteReportsTruthfulStages(t *testing.T) {
	f := newWriteGroupRouterFixture(t)
	created := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, f.createRequest())
	require.Equal(t, http.StatusCreated, created.Code, created.Body.String())
	_, group := groupSaveData(t, decodeJSONBody(t, created))
	id := group["id"].(string)

	view := &grouppipeline.DeliveryView{
		GroupID: id,
		Intake:  grouppipeline.IntakeView{State: "active"},
		Stages:  grouppipeline.StagesView{Collecting: 2, Queued: 3, Blocked: 1, SQLCommitted: 5},
		Backlog: []grouppipeline.BacklogView{{
			GroupRevision: "rev-1", ConnectorID: "connector-A", ConnectorRevision: "connector-1", TableName: "raw_values", Pending: 4,
			ErrorCodes: []string{"target-blocked"},
		}},
		Quota: groupdelivery.QuotaStatus{Configured: true, State: groupdelivery.QuotaOK, Scope: groupdelivery.QuotaScopeGlobal, UsedBytes: 10, MaxBytes: 100}.View(),
	}
	router := f.routerWithDelivery(fakeDeliveryReader{view: view})

	response := performJSONRequest(t, router, http.MethodGet, writeGroupsPath+"/"+id+"/delivery", nil)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	body := decodeJSONBody(t, response)
	require.Equal(t, true, body["success"])
	data := body["data"].(map[string]any)
	stages := data["stages"].(map[string]any)
	require.EqualValues(t, 5, stages["sql_committed"])
	require.EqualValues(t, 3, stages["queued"])
	require.EqualValues(t, 2, stages["collecting"])
	require.Nil(t, data["last_sql_committed_at"], "nothing confirmed means no committed time, whatever is queued")
	backlog := data["backlog"].([]any)[0].(map[string]any)
	require.Equal(t, "connector-1", backlog["connector_revision"], "backlog keeps the destination revision it was accepted for")
	require.Equal(t, []any{"target-blocked"}, backlog["error_codes"])
	require.Equal(t, "ok", data["quota"].(map[string]any)["state"])
	require.NotContains(t, response.Body.String(), "dsn")
	require.NotContains(t, response.Body.String(), "password")
}

func TestRevisionBoundBacklogDeliveryRouteUnknownGroupAndUnavailableReader(t *testing.T) {
	f := newWriteGroupRouterFixture(t)
	router := f.routerWithDelivery(fakeDeliveryReader{view: &grouppipeline.DeliveryView{}})
	missing := performJSONRequest(t, router, http.MethodGet, writeGroupsPath+"/does-not-exist/delivery", nil)
	require.Equal(t, http.StatusNotFound, missing.Code, missing.Body.String())
	assertWriteGroupError(t, decodeJSONBody(t, missing), "WRITE_GROUP_NOT_FOUND")

	created := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, f.createRequest())
	require.Equal(t, http.StatusCreated, created.Code)
	_, group := groupSaveData(t, decodeJSONBody(t, created))
	id := group["id"].(string)

	unwired := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+id+"/delivery", nil)
	require.Equal(t, http.StatusServiceUnavailable, unwired.Code, "no delivery source is an outage, never an empty success")
	assertWriteGroupError(t, decodeJSONBody(t, unwired), "WRITE_GROUP_UNAVAILABLE")

	failing := f.routerWithDelivery(fakeDeliveryReader{err: errors.New("disk I/O error at /private/path")})
	failed := performJSONRequest(t, failing, http.MethodGet, writeGroupsPath+"/"+id+"/delivery", nil)
	require.Equal(t, http.StatusServiceUnavailable, failed.Code)
	require.NotContains(t, failed.Body.String(), "/private/path")
}
