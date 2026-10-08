package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-gateway/internal/datalink/modbusshare"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestModbusShareHandler_ReconcileEmptyMappings(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, test := range []struct {
		name, field string
		wantStatus  int
	}{
		{name: "omitted", wantStatus: http.StatusOK},
		{name: "null", field: `,"desired_mappings":null`, wantStatus: http.StatusOK},
		{name: "empty", field: `,"desired_mappings":[]`, wantStatus: http.StatusOK},
		{name: "unexpected mapping", field: `,"desired_mappings":[{"tag_id":"foreign"}]`, wantStatus: http.StatusUnprocessableEntity},
	} {
		t.Run(test.name, func(t *testing.T) {
			hydration := modbusshare.HydrationState{
				State: modbusshare.HydrationStateReady, Readiness: true,
				WorkspaceID: "ws", WorkspaceRevision: "workspace-rev", SettingsRevision: "settings-rev", ReadinessToken: "token",
			}
			gate := &mockHydrationGate{hydrationState: hydration, settings: modbusshare.Settings{Enabled: true}}
			handler, svc, _ := setupGatedShareHandler(t, gate)
			svc.SetHydrationState(hydration)
			handler.WithReconciler(modbusshare.NewReconciler(svc, emptyReconcileStore{}))
			handler.WithDesiredMappingBuilder(func(context.Context, string, modbusshare.Settings) ([]modbusshare.DesiredMapping, error) {
				return []modbusshare.DesiredMapping{}, nil
			})
			router := gin.New()
			router.POST("/reconcile", handler.Reconcile)
			body := `{"workspace_id":"ws","expected_workspace_revision":"workspace-rev","expected_settings_revision":"settings-rev","readiness_token":"token"` + test.field + `}`
			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/reconcile", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			router.ServeHTTP(response, req)
			require.Equal(t, test.wantStatus, response.Code, response.Body.String())
			if test.wantStatus == http.StatusOK {
				require.Contains(t, response.Body.String(), `"outcome":"aligned"`)
			} else {
				require.Contains(t, response.Body.String(), modbusshare.ErrCodeProjectionRequired)
			}
		})
	}
}

// Empty reconciliation must be an aligned no-op and never mutate the store.
type emptyReconcileStore struct{}

func (emptyReconcileStore) GetRevision(context.Context, string) (revision string, dirty bool, err error) {
	return "workspace-rev", false, nil
}

func (emptyReconcileStore) GetDesiredMappings(context.Context, string) ([]modbusshare.DesiredMapping, error) {
	return nil, nil
}

func (emptyReconcileStore) UpdateRevision(context.Context, string, string, string) error {
	return errors.New("unexpected revision mutation")
}

func (emptyReconcileStore) MarkDirty(context.Context, string) error {
	return errors.New("unexpected dirty mutation")
}

func (emptyReconcileStore) CommitDesiredMappings(context.Context, string, string, string, []modbusshare.DesiredMapping) error {
	return errors.New("unexpected desired mappings mutation")
}
