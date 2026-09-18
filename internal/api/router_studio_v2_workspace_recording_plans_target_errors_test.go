package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"go-gateway/internal/api/handlers"
	"go-gateway/internal/datalink/dbtarget"
)

func TestStudioV2WorkspaceRecordingPlans_ResolverSentinelsUseSafeResponses(t *testing.T) {
	tests := []struct {
		name        string
		resolverErr error
		status      int
		code        string
		retryable   bool
	}{
		{name: "unknown connector", resolverErr: fmt.Errorf("foreign connector detail: %w", dbtarget.ErrConnectorNotFound), status: http.StatusNotFound, code: "RECORDING_PLAN_NOT_FOUND"},
		{name: "stale revision", resolverErr: fmt.Errorf("revision detail: %w", dbtarget.ErrConnectorRevisionConflict), status: http.StatusConflict, code: "RECORDING_CONNECTOR_REVISION_CONFLICT"},
		{name: "resolver failure", resolverErr: errors.New("private resolver diagnostic"), status: http.StatusServiceUnavailable, code: "RECORDING_PLAN_UNAVAILABLE", retryable: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := &recordingTargetResolverStub{err: tt.resolverErr}
			fixture := newRecordingTargetResolverRouter(t, resolver)
			w := serveRecordingTargetResolverRequest(fixture.router, fixture.previewBody(t, nil))
			if w.Code != tt.status {
				t.Fatalf("resolver error status = %d, want %d: %s", w.Code, tt.status, w.Body.String())
			}
			var response struct {
				Success bool                           `json:"success"`
				Error   handlers.TypedAPIErrorEnvelope `json:"error"`
			}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("decode resolver error: %v; body=%s", err, w.Body.String())
			}
			if response.Success || response.Error.Code != tt.code || response.Error.Retryable != tt.retryable || response.Error.RequestID == "" {
				t.Fatalf("unsafe resolver error response: %+v; body=%s", response, w.Body.String())
			}
			if strings.Contains(w.Body.String(), "diagnostic") || strings.Contains(w.Body.String(), "foreign connector detail") || strings.Contains(w.Body.String(), "revision detail") {
				t.Fatalf("resolver diagnostic leaked into response: %s", w.Body.String())
			}
		})
	}
}
