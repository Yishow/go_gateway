package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-gateway/internal/datalink/audit"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestStudioV2WorkspaceDatabaseHandler_RecordsDeliveryImpactingAuditEntries(t *testing.T) {
	fixture := newWorkspaceDatabaseFixture(t)

	saveWorkspaceDatabaseConfig(t, fixture)
	saveValidTarget(t, fixture, fixture.pointIDs[0], "line_a")

	record, err := fixture.workspaceSvc.GetOrCreate(context.Background())
	require.NoError(t, err)
	entries, err := fixture.auditSvc.List(context.Background(), audit.ListFilter{WorkspaceID: record.ID, Limit: 10})
	require.NoError(t, err)
	require.Len(t, entries, 2)

	require.Equal(t, audit.EventTypeDatabaseTargetSaved, entries[0].EventType)
	require.Equal(t, audit.ResultSuccess, entries[0].Result)
	require.Contains(t, entries[0].Scope, fixture.pointIDs[0])
	require.Contains(t, entries[0].Details, "line_a")

	require.Equal(t, audit.EventTypeDatabaseConfigSaved, entries[1].EventType)
	require.Equal(t, audit.ResultSuccess, entries[1].Result)
	require.Contains(t, entries[1].Scope, "database_connector")
}

func TestStudioV2WorkspaceAuditHandler_ListReturnsRecentWorkspaceHistory(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctx := context.Background()
	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository())
	auditSvc := audit.NewService(audit.NewMemoryRepository())
	record, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)

	require.NoError(t, auditSvc.Record(ctx, audit.RecordEvent{
		WorkspaceID: record.ID,
		EventType:   audit.EventTypeWorkspaceActivation,
		Result:      audit.ResultPartialSuccess,
		Scope:       "devices:dev-A,dev-B",
		ReferenceID: "activation-1",
		Details:     map[string]string{"message": "dev-B failed"},
		OccurredAt:  time.Date(2026, 5, 29, 10, 10, 0, 0, time.UTC),
	}))
	require.NoError(t, auditSvc.Record(ctx, audit.RecordEvent{
		WorkspaceID: record.ID,
		EventType:   audit.EventTypeDatabaseConfigSaved,
		Result:      audit.ResultSuccess,
		Scope:       "database_connector:sqlite",
		ReferenceID: "connector-1",
		Details:     map[string]string{"runtime_apply_status": "applied"},
		OccurredAt:  time.Date(2026, 5, 29, 10, 12, 0, 0, time.UTC),
	}))

	handler := NewStudioV2WorkspaceAuditHandler(workspaceSvc, auditSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/datalink/studio-v2/workspace/audit-history?limit=5", nil)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.List(c)

	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
	var body struct {
		Success bool `json:"success"`
		Data    struct {
			Entries []audit.Entry `json:"entries"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body))
	require.True(t, body.Success)
	require.Len(t, body.Data.Entries, 2)
	require.Equal(t, audit.EventTypeDatabaseConfigSaved, body.Data.Entries[0].EventType)
	require.Equal(t, audit.EventTypeWorkspaceActivation, body.Data.Entries[1].EventType)
	require.Equal(t, record.ID, body.Data.Entries[0].WorkspaceID)
	require.Contains(t, body.Data.Entries[1].Details, "dev-B failed")
}
