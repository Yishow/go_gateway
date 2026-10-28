package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func lifecycleErrorResponse(t *testing.T, err error) (int, map[string]any) {
	t.Helper()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/write-groups/G1", nil)
	c.Request.Header.Set("X-Request-ID", "req-lifecycle")
	renderWriteGroupError(c, err)
	var response map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	require.Equal(t, false, response["success"])
	require.NotContains(t, w.Body.String(), "private-")
	return w.Code, response["error"].(map[string]any)
}

func TestWriteGroupLifecycleErrorDetailsKeepStatusAndCode(t *testing.T) {
	for _, tc := range []struct {
		name, reason, next, code string
		cause                    error
		status                   int
		retryable                bool
	}{
		{"stale", "connector_identity_stale", "stop", "revision_mismatch", workspace.ErrSetupRevisionConflict, 409, false},
		{"missing", "connector_unavailable", "stop", "WRITE_GROUP_NOT_FOUND", workspace.ErrWriteGroupConnectorNotFound, 404, false},
		{"read_failed", "connector_unavailable", "review_connector", "WRITE_GROUP_UNAVAILABLE", errors.New("private-database"), 500, true},
		{"unproven", "ownership_unavailable", "stop", "WRITE_GROUP_LIFECYCLE_BLOCKED", workspace.ErrWriteGroupLifecycleBlocked, 409, false},
		{"stopped", "ownership_unavailable", "inspect_ownership", "WRITE_GROUP_LIFECYCLE_BLOCKED", workspace.ErrWriteGroupLifecycleBlocked, 409, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := errors.Join(tc.cause, &workspace.WriteGroupLifecycleActionError{Action: "remove", ReasonCode: tc.reason, NextAction: tc.next}, errors.New("private-proof private-dsn"))
			status, envelope := lifecycleErrorResponse(t, err)
			require.Equal(t, tc.status, status)
			require.Equal(t, tc.code, envelope["code"])
			require.Equal(t, tc.reason, envelope["reason_code"])
			require.Equal(t, tc.next, envelope["next_action"])
			legacy := "review_group"
			if tc.next == "stop" {
				legacy = "disable_group"
			}
			require.Equal(t, legacy, envelope["action"])
			require.Equal(t, tc.retryable, envelope["retryable"])
			require.Equal(t, "req-lifecycle", envelope["request_id"])
		})
	}
}

func TestWriteGroupLifecycleErrorUnknownDetailsUseSafeLegacyFallback(t *testing.T) {
	for _, detail := range []*workspace.WriteGroupLifecycleActionError{
		{},
		{Action: "private-action", ReasonCode: "ownership_unavailable", NextAction: "stop"},
		{Action: "stop", ReasonCode: "ownership_unavailable", NextAction: "stop"},
		{Action: "remove", ReasonCode: "private-reason", NextAction: "stop"},
		{Action: "remove", ReasonCode: "ownership_unavailable", NextAction: "private-next"},
		{Action: "remove", ReasonCode: "ownership_unavailable", NextAction: "review_connector"},
		{Action: "remove", ReasonCode: "connector_identity_stale", NextAction: "inspect_ownership"},
	} {
		err := errors.Join(workspace.ErrWriteGroupLifecycleBlocked, detail, errors.New("private-error"))
		status, envelope := lifecycleErrorResponse(t, err)
		require.Equal(t, http.StatusConflict, status)
		require.Equal(t, "WRITE_GROUP_LIFECYCLE_BLOCKED", envelope["code"])
		require.Equal(t, "disable_group", envelope["action"])
		require.Equal(t, "write group cannot be deleted while accepted data ownership is unresolved", envelope["message"])
		require.NotContains(t, envelope, "reason_code")
		require.NotContains(t, envelope, "next_action")
	}
}
