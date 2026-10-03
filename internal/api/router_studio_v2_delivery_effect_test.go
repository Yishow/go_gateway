package api

import (
	"context"
	"net/http"
	"testing"

	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

type revisionDeliveryReader struct {
	groupID, revision string
}

func (r *revisionDeliveryReader) Delivery(context.Context, string) (*grouppipeline.DeliveryView, error) {
	panic("revision-aware delivery must be selected")
}

func (r *revisionDeliveryReader) DeliveryForRevision(_ context.Context, id, revision string) (*grouppipeline.DeliveryView, error) {
	r.groupID, r.revision = id, revision
	return &grouppipeline.DeliveryView{GroupID: id}, nil
}

func TestRecordingDeliveryUsesPersistedAppliedRevisionInsteadOfClientRevision(t *testing.T) {
	f := newWriteGroupRouterFixture(t)
	created := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, f.createRequest())
	require.Equal(t, http.StatusCreated, created.Code, created.Body.String())
	_, group := groupSaveData(t, decodeJSONBody(t, created))
	id := group["id"].(string)
	reader := &revisionDeliveryReader{}
	router := NewRouter(&DatalinkServices{Workspace: f.workspace,
		WriteGroups:        workspace.NewWriteGroupService(f.workspace, workspace.NewSQLWriteGroupRepository(f.db)),
		WriteGroupDelivery: reader,
	})
	for _, revision := range []string{"", "persisted-applied-revision"} {
		_, err := f.db.ExecContext(t.Context(), `UPDATE write_groups SET applied_revision=? WHERE id=?`, revision, id)
		require.NoError(t, err)
		response := performJSONRequest(t, router, http.MethodGet, writeGroupsPath+"/"+id+"/delivery?group_revision=foreign-revision", nil)
		require.Equal(t, http.StatusOK, response.Code, response.Body.String())
		require.Equal(t, id, reader.groupID)
		require.Equal(t, revision, reader.revision)
		require.Nil(t, decodeJSONBody(t, response)["data"].(map[string]any)["last_sql_committed_effect"])
	}
}
