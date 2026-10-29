package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"
)

type recoveryErrorReader struct {
	result *recordingplan.TestWriteRecoveryList
	err    error
}

func (r recoveryErrorReader) GetSchemaOperation(context.Context, string, string) (*recordingplan.SchemaOperation, error) {
	panic("listing must not fall back to an operation ID lookup")
}
func (r recoveryErrorReader) ListUnresolvedTestWrites(context.Context, string) (*recordingplan.TestWriteRecoveryList, error) {
	return r.result, r.err
}

type ordinaryOnlyOperationReader struct{}

func (ordinaryOnlyOperationReader) GetSchemaOperation(context.Context, string, string) (*recordingplan.SchemaOperation, error) {
	panic("unsupported listing must fail closed")
}

func TestTestWriteRecoveryHandlerSanitizesFailedObservations(t *testing.T) {
	for _, tc := range []struct {
		name   string
		reader schemaOperationReader
	}{
		{"raw repository failure", recoveryErrorReader{err: errors.New("postgres://user:password@secret-host/private")}},
		{"partial result", recoveryErrorReader{result: &recordingplan.TestWriteRecoveryList{Complete: false}}},
		{"foreign response", recoveryErrorReader{result: &recordingplan.TestWriteRecoveryList{WorkspaceID: "foreign", Complete: true}}},
		{"nil response", recoveryErrorReader{}}, {"unsupported repository", ordinaryOnlyOperationReader{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ws := workspace.NewService(workspace.NewMemoryRepository())
			h := NewStudioV2WorkspaceDatabaseOperationsHandler(ws, tc.reader)
			r := gin.New()
			r.GET("/recovery", h.ListTestWrites)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/recovery", nil))
			require.Equal(t, http.StatusServiceUnavailable, w.Code)
			require.Contains(t, w.Body.String(), `"code":"RECORDING_PLAN_UNAVAILABLE"`)
			require.Contains(t, w.Body.String(), `"retryable":true`)
			for _, secret := range []string{"secret-host", "password", "postgres://", "foreign", `"data"`, `"operations"`} {
				require.NotContains(t, w.Body.String(), secret)
			}
		})
	}
}

type recoveryFailingWorkspaceRepository struct{ workspace.Repository }

func (recoveryFailingWorkspaceRepository) Get(context.Context) (*workspace.Record, error) {
	return nil, errors.New("password=private at secret-host")
}

func TestTestWriteRecoveryHandlerWorkspaceFailureIsTyped(t *testing.T) {
	h := NewStudioV2WorkspaceDatabaseOperationsHandler(workspace.NewService(recoveryFailingWorkspaceRepository{}), recoveryErrorReader{})
	r := gin.New()
	r.GET("/recovery", h.ListTestWrites)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/recovery", nil))
	require.Equal(t, http.StatusServiceUnavailable, w.Code)
	require.Contains(t, w.Body.String(), `"code":"RECORDING_PLAN_UNAVAILABLE"`)
	require.Contains(t, w.Body.String(), `"retryable":true`)
	require.NotContains(t, w.Body.String(), "private")
	require.NotContains(t, w.Body.String(), "secret-host")
}
