package api

import (
	"encoding/json"
	"net/http"
	"testing"

	swagger "go-gateway/docs/swagger"
	"go-gateway/internal/datalink/audit"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestRecordingDeleteWritesAuditThroughRouter(t *testing.T) {
	f := newWriteGroupRouterFixture(t)
	auditSvc := audit.NewService(audit.NewSQLRepository(f.db))
	groups := workspace.NewWriteGroupService(f.workspace, workspace.NewSQLWriteGroupRepository(f.db))
	router := NewRouter(&DatalinkServices{Workspace: f.workspace, WriteGroups: groups, Audit: auditSvc})
	created := performJSONRequest(t, router, http.MethodPost, writeGroupsPath, f.createRequest())
	require.Equal(t, http.StatusCreated, created.Code, created.Body.String())
	data, group := groupSaveData(t, decodeJSONBody(t, created))
	id := group["id"].(string)
	request := groupEditRequest(f, data, group)
	delete(request, "group")
	response := performJSONRequest(t, router, http.MethodDelete, writeGroupsPath+"/"+id, request)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	entries, err := auditSvc.List(t.Context(), audit.ListFilter{WorkspaceID: f.record.ID, Limit: 10})
	require.NoError(t, err)
	require.Len(t, entries, 1)
	require.Equal(t, audit.EventType("write_group_deleted"), entries[0].EventType)
	require.Equal(t, id, entries[0].ReferenceID)
	require.NotContains(t, entries[0].Details, "connector")
	require.NotContains(t, entries[0].Details, "password")
}

func TestRecordingDeleteSwaggerContract(t *testing.T) {
	var doc struct {
		Paths map[string]map[string]struct {
			Responses map[string]struct {
				Schema struct {
					Ref string `json:"$ref"`
				} `json:"schema"`
			} `json:"responses"`
		} `json:"paths"`
	}
	require.NoError(t, json.Unmarshal([]byte(swagger.SwaggerInfo.ReadDoc()), &doc))
	for _, tc := range []struct{ path, successSchema string }{
		{"/v1/datalink/studio-v2/workspace/write-groups/{id}", "#/definitions/handlers.WriteGroupDeleteResponse"},
		{"/v1/datalink/studio-v2/workspace/devices/{id}", "#/definitions/handlers.DeviceDeleteResponse"},
	} {
		t.Run(tc.path, func(t *testing.T) {
			op, ok := doc.Paths[tc.path]["delete"]
			require.True(t, ok)
			for _, status := range []string{"200", "404", "409", "500", "503"} {
				_, ok := op.Responses[status]
				require.True(t, ok, status)
			}
			require.Equal(t, tc.successSchema, op.Responses["200"].Schema.Ref)
			require.Equal(t, "#/definitions/handlers.APIErrorResponse", op.Responses["409"].Schema.Ref)
		})
	}
}
