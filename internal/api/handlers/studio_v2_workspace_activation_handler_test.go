package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

type stubWorkspaceActivator struct {
	response *workspace.ActivationResponse
	err      error
	called   bool
}

func (s *stubWorkspaceActivator) ActivateEligible(_ context.Context) (*workspace.ActivationResponse, error) {
	s.called = true
	return s.response, s.err
}

type stubWorkspaceSchemaEnsurer struct {
	err    error
	called bool
}

func (s *stubWorkspaceSchemaEnsurer) EnsureWorkspaceSchema(_ context.Context) error {
	s.called = true
	return s.err
}

func TestStudioV2WorkspaceActivationHandler_ReturnsPerDeviceResultsOnPartialSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewStudioV2WorkspaceActivationHandler(&stubWorkspaceActivator{
		response: &workspace.ActivationResponse{
			WorkspaceID: "workspace-1",
			Results: []workspace.ActivationResult{
				{DeviceID: "dev-A", Status: workspace.ActivationResultStatusSuccess, Message: "activated"},
				{DeviceID: "dev-B", Status: workspace.ActivationResultStatusFailed, Message: "activation timeout"},
			},
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/datalink/studio-v2/workspace/activate", nil)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.Activate(c)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", resp.Code, resp.Body.String())
	}
	if body := resp.Body.String(); body == "" || !containsAll(body, "workspace-1", "dev-A", "success", "dev-B", "failed", "activation timeout") {
		t.Fatalf("expected per-device partial success payload, got %s", body)
	}
}

func TestStudioV2WorkspaceActivationHandler_HidesGenericSchemaEnsureErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)

	activator := &stubWorkspaceActivator{}
	ensurer := &stubWorkspaceSchemaEnsurer{err: errors.New("dial tcp 127.0.0.1:5432: connect: connection refused")}
	handler := NewStudioV2WorkspaceActivationHandler(activator, ensurer)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/datalink/studio-v2/workspace/activate", nil)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.Activate(c)

	if resp.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d body=%s", resp.Code, resp.Body.String())
	}
	if !ensurer.called {
		t.Fatal("expected schema ensurer to be called")
	}
	if activator.called {
		t.Fatal("expected activator to be skipped when schema ensure fails")
	}
	if body := resp.Body.String(); !containsAll(body, "success", "false", "Studio V2 database operation failed") {
		t.Fatalf("expected generic schema error payload, got %s", body)
	} else if strings.Contains(body, "connection refused") {
		t.Fatalf("expected generic schema error payload without raw infra details, got %s", body)
	}
}

func TestStudioV2WorkspaceActivationHandler_Returns422ForValidationSchemaEnsureError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	activator := &stubWorkspaceActivator{}
	ensurer := &stubWorkspaceSchemaEnsurer{err: fmt.Errorf("%w: schema preflight failed", dbtarget.ErrValidation)}
	handler := NewStudioV2WorkspaceActivationHandler(activator, ensurer)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/datalink/studio-v2/workspace/activate", nil)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.Activate(c)

	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status 422, got %d body=%s", resp.Code, resp.Body.String())
	}
	if !ensurer.called {
		t.Fatal("expected schema ensurer to be called")
	}
	if activator.called {
		t.Fatal("expected activator to be skipped when schema ensure fails")
	}
	if body := resp.Body.String(); !containsAll(body, "success", "false", "validation", "schema preflight failed") {
		t.Fatalf("expected validation schema error payload, got %s", body)
	}
}

func TestStudioV2WorkspaceActivationHandler_ReturnsBlockingReadinessIssues(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewStudioV2WorkspaceActivationHandler(&stubWorkspaceActivator{
		err: &workspace.ReadinessBlockedError{
			Operation: "activation",
			Summary: &workspace.ReadinessSummary{
				BlockingCount: 1,
				Issues: []workspace.ReadinessIssue{
					{
						Code:     "device-probe-required",
						Severity: workspace.ReadinessSeverityBlocking,
						Step:     workspace.ReadinessStep1,
						Scope:    "dev-A",
						Message:  "probe diagnostics failed",
					},
				},
			},
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/datalink/studio-v2/workspace/activate", nil)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.Activate(c)

	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status 422, got %d body=%s", resp.Code, resp.Body.String())
	}
	if body := resp.Body.String(); !containsAll(body, "success", "false", "readiness_blocked", "device-probe-required", "dev-A") {
		t.Fatalf("expected readiness blocker payload, got %s", body)
	}
}

func containsAll(body string, parts ...string) bool {
	for _, part := range parts {
		if !strings.Contains(body, part) {
			return false
		}
	}
	return true
}
