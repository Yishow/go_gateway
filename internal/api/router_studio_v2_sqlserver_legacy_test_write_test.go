package api

import (
	"context"
	"net/http"
	"testing"
	"time"

	"go-gateway/internal/api/handlers"
	"go-gateway/internal/datalink/grouptestwrite"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

const (
	sqlServerLegacyGroupID = "group-sql-legacy"
	sqlServerLegacyPlanID  = "plan-sql-legacy"
)

type sqlServerLegacyResolver struct {
	groupID string
	err     error
	calls   int
}

func (r *sqlServerLegacyResolver) ResolveRecordingPlanGroup(context.Context, string) (string, error) {
	r.calls++
	return r.groupID, r.err
}

type sqlServerLegacyWriter struct {
	preview     *grouptestwrite.Preview
	outcome     *grouptestwrite.Outcome
	err         error
	previews    int
	confirms    int
	gotGroup    string
	gotWS       grouptestwrite.Workspace
	gotConfirms []grouptestwrite.Confirmation
}

func (w *sqlServerLegacyWriter) Preview(_ context.Context, ws grouptestwrite.Workspace, groupID string) (*grouptestwrite.Preview, error) {
	w.previews++
	w.gotWS, w.gotGroup = ws, groupID
	return w.preview, w.err
}

func (w *sqlServerLegacyWriter) Confirm(_ context.Context, ws grouptestwrite.Workspace, req grouptestwrite.Confirmation) (*grouptestwrite.Outcome, error) {
	w.confirms++
	w.gotWS, w.gotConfirms = ws, append(w.gotConfirms, req)
	if w.err != nil {
		return w.outcome, w.err
	}
	if w.preview == nil || req.Token != w.preview.Token || req.OperationID != w.preview.OperationID {
		return nil, recordingplan.ErrSchemaOperationMismatch
	}
	if req.GroupID != w.preview.GroupID {
		return w.outcome, grouptestwrite.ErrGroupMismatch
	}
	return w.outcome, nil
}

func newSQLServerLegacyTestWriteRouter(t *testing.T, writer *sqlServerLegacyWriter, resolver *sqlServerLegacyResolver) (*gin.Engine, *workspace.Service) {
	t.Helper()
	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository())
	_, err := workspaceSvc.GetOrCreate(t.Context())
	require.NoError(t, err)
	handler := handlers.NewStudioV2WorkspaceWriteGroupTestWriteHandler(workspaceSvc, writer, resolver)
	router := gin.New()
	router.POST("/api/v1/datalink/studio-v2/workspace/recording-plans/test-write-preview", handler.LegacyPreview)
	router.POST("/api/v1/datalink/studio-v2/workspace/recording-plans/test-write", handler.LegacyConfirm)
	return router, workspaceSvc
}

func sqlServerLegacyPreview() *grouptestwrite.Preview {
	return &grouptestwrite.Preview{
		Token: "sql-token-1", OperationID: "op-sql-1", Action: recordingplan.TestWriteAction,
		GroupID: sqlServerLegacyGroupID, GroupRevision: "group-rev-1",
		ExpiresAt: time.Now().UTC().Add(time.Hour),
		Target: grouptestwrite.PreviewTarget{
			ConnectorID: "conn-sql", Dialect: "sqlserver", Database: "recordings", Schema: "dbo", Table: "gw_long_sql",
		},
		OwnerColumn: "_gw_owner_sql", OwnerValue: "gw-test-op-sql", Dedupe: "receipt",
		RecordIDs: []string{"record-sql-1"}, Cleanup: "remove only the exact test record IDs and receipts",
	}
}

func sqlServerLegacyOperation() *grouptestwrite.Outcome {
	now := time.Now().UTC()
	return &grouptestwrite.Outcome{
		Claim: recordingplan.ClaimCompleted,
		Operation: &recordingplan.SchemaOperation{
			OperationID: "op-sql-1", Action: recordingplan.TestWriteAction, Status: recordingplan.SchemaOperationSucceeded,
			WriteOutcome: grouptestwrite.WriteVerified, CleanupStatus: grouptestwrite.CleanupCleaned,
			CreatedAt: now, UpdatedAt: now, CompletedAt: &now,
		},
	}
}

func TestSQLServerLegacyTestWriteBridgeRequiresPreviewTokenAndKeepsScope(t *testing.T) {
	preview := sqlServerLegacyPreview()
	writer := &sqlServerLegacyWriter{preview: preview, outcome: sqlServerLegacyOperation()}
	resolver := &sqlServerLegacyResolver{groupID: sqlServerLegacyGroupID}
	router, workspaceSvc := newSQLServerLegacyTestWriteRouter(t, writer, resolver)

	previewResponse := performJSONRequest(t, router, http.MethodPost,
		"/api/v1/datalink/studio-v2/workspace/recording-plans/test-write-preview", map[string]any{"plan_id": sqlServerLegacyPlanID})
	require.Equal(t, http.StatusOK, previewResponse.Code, previewResponse.Body.String())
	previewBody := decodeJSONBody(t, previewResponse)
	previewData := previewBody["data"].(map[string]any)
	require.Equal(t, preview.Token, previewData["token"])
	require.Equal(t, preview.OperationID, previewData["operation_id"])
	require.Equal(t, string(recordingplan.TestWriteAction), previewData["action"])
	require.Equal(t, "sqlserver", previewData["target"].(map[string]any)["dialect"])
	require.Equal(t, sqlServerLegacyGroupID, writer.gotGroup)
	require.Equal(t, sqlServerLegacyGroupID, previewData["group_id"])

	confirmResponse := performJSONRequest(t, router, http.MethodPost,
		"/api/v1/datalink/studio-v2/workspace/recording-plans/test-write",
		map[string]any{"plan_id": sqlServerLegacyPlanID, "token": preview.Token, "operation_id": preview.OperationID})
	require.Equal(t, http.StatusOK, confirmResponse.Code, confirmResponse.Body.String())
	require.Equal(t, 1, writer.confirms)
	require.Equal(t, sqlServerLegacyGroupID, writer.gotConfirms[0].GroupID)
	require.Equal(t, preview.Token, writer.gotConfirms[0].Token)
	require.Equal(t, preview.OperationID, writer.gotConfirms[0].OperationID)
	require.Equal(t, "op-sql-1", decodeJSONBody(t, confirmResponse)["data"].(map[string]any)["operation_id"])
	require.Equal(t, writer.gotWS.ID, workspaceIDFromStore(t, workspaceSvc), "workspace identity is server-resolved")
}

func TestSQLServerLegacyTestWriteNakedAndStaleRequestsStayClosed(t *testing.T) {
	preview := sqlServerLegacyPreview()
	writer := &sqlServerLegacyWriter{preview: preview, outcome: sqlServerLegacyOperation()}
	resolver := &sqlServerLegacyResolver{groupID: sqlServerLegacyGroupID}
	router, _ := newSQLServerLegacyTestWriteRouter(t, writer, resolver)

	naked := performJSONRequest(t, router, http.MethodPost,
		"/api/v1/datalink/studio-v2/workspace/recording-plans/test-write", map[string]any{"plan_id": sqlServerLegacyPlanID})
	require.Equal(t, http.StatusBadRequest, naked.Code, naked.Body.String())
	require.Zero(t, writer.confirms)
	require.Zero(t, resolver.calls)

	writer.err = recordingplan.ErrPreviewTokenStale
	stale := performJSONRequest(t, router, http.MethodPost,
		"/api/v1/datalink/studio-v2/workspace/recording-plans/test-write",
		map[string]any{"plan_id": sqlServerLegacyPlanID, "token": "sql-token-stale", "operation_id": preview.OperationID})
	require.Equal(t, http.StatusConflict, stale.Code, stale.Body.String())
	require.Equal(t, "WRITE_GROUP_TEST_WRITE_PREVIEW_STALE", decodeJSONBody(t, stale)["error"].(map[string]any)["code"])

	writer.err = recordingplan.ErrSchemaOperationMismatch
	mismatch := performJSONRequest(t, router, http.MethodPost,
		"/api/v1/datalink/studio-v2/workspace/recording-plans/test-write",
		map[string]any{"plan_id": sqlServerLegacyPlanID, "token": preview.Token, "operation_id": "op-other"})
	require.Equal(t, http.StatusConflict, mismatch.Code, mismatch.Body.String())
	require.Equal(t, "WRITE_GROUP_TEST_WRITE_OPERATION_MISMATCH", decodeJSONBody(t, mismatch)["error"].(map[string]any)["code"])

	writer.err = recordingplan.ErrPreviewTokenKind
	wrongKind := performJSONRequest(t, router, http.MethodPost,
		"/api/v1/datalink/studio-v2/workspace/recording-plans/test-write",
		map[string]any{"plan_id": sqlServerLegacyPlanID, "token": preview.Token, "operation_id": preview.OperationID})
	require.Equal(t, http.StatusUnprocessableEntity, wrongKind.Code, wrongKind.Body.String())
	require.Equal(t, "WRITE_GROUP_TEST_WRITE_TOKEN_KIND", decodeJSONBody(t, wrongKind)["error"].(map[string]any)["code"])
	require.False(t, recordingplan.ManagedSchemaExecutionVerified("sqlserver"), "generic SQL DDL must remain closed")
}

func workspaceIDFromStore(t *testing.T, svc *workspace.Service) string {
	t.Helper()
	record, err := svc.GetOrCreate(t.Context())
	require.NoError(t, err)
	return record.ID
}
