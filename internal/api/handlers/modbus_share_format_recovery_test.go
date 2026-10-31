package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/point"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestExplicitFormatRecoveryDoesNotUseActivationGateOrAcceptInjectedPlan(t *testing.T) {
	for _, tc := range []struct {
		body  string
		want  int
		calls int
	}{
		{`{"format_recovery":{"workspace_id":"ws","request_id":"original"}}`, 200, 1},
		{`{"format_recovery":{"workspace_id":"ws"},"desired_mappings":[]}`, 422, 0},
		{`{"format_recovery":{"workspace_id":"ws","restore":true}}`, 422, 0},
		{`{"format_recovery":null}`, 422, 0},
	} {
		h := NewModbusShareHandler(modbusshare.NewService(nil, 64))
		calls := 0
		h.WithFormatRecovery(func(_ context.Context, r point.FormatRecoveryRequest) (FormatRecoveryResult, error) {
			calls++
			require.Equal(t, "original", r.RequestID)
			return FormatRecoveryResult{Outcome: "recovered", RuntimeState: "not_started"}, nil
		})
		router := gin.New()
		router.POST("/reconcile", h.Reconcile)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/reconcile", strings.NewReader(tc.body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rec, req)
		require.Equal(t, tc.want, rec.Code, rec.Body.String())
		require.Equal(t, tc.calls, calls)
	}
}
func TestExplicitFormatRecoveryMissingAuthorityFailsClosed(t *testing.T) {
	h := NewModbusShareHandler(modbusshare.NewService(nil, 64))
	router := gin.New()
	router.POST("/reconcile", h.Reconcile)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/reconcile", strings.NewReader(`{"format_recovery":{}}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusServiceUnavailable, rec.Code)
}
