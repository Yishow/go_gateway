package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go-gateway/internal/datalink/storage"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	datalinkruntime "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
)

func TestRuntimeRuleStatusSafetyLegacyNestedDiagnostics(t *testing.T) {
	hostile := "postgres://operator:password-fixture@private-host:5432/private-db SELECT secret FROM readings /private/sqlite.db driver exception"
	response := runtimeStatusResponse{
		SnapshotState:       datalinkruntime.RuntimeUnavailableTruthState(hostile),
		Collectors:          []runtimeCollectorResponse{{DeviceID: "device-1", DeviceName: "Meter A", LastError: &hostile, AvailabilityReason: &hostile, ProjectionMessage: hostile}},
		DatabaseDelivery:    []datalinkruntime.DatabaseDeliveryDiagnostic{{PointID: "point-1", TagID: "tag-1", Status: "failed", Error: hostile, LastFailureReason: hostile}},
		ModbusShareDelivery: []datalinkruntime.ModbusShareDeliveryDiagnostic{{PointID: "point-1", TagID: "tag-1", Status: "failed", Error: hostile}},
		Diagnostics:         []datalinkruntime.RuntimeFlowDiagnostic{{PointID: "point-1", FailureReason: hostile, Stages: []datalinkruntime.RuntimeFlowDiagnosticStage{{Stage: "database_delivery", Status: "failed", Reason: hostile}}}},
	}
	wire, err := json.Marshal(response)
	require.NoError(t, err)
	require.NotContains(t, string(wire), "password-fixture")
	require.NotContains(t, string(wire), "private-host")
	require.NotContains(t, string(wire), "SELECT secret")
	require.NotContains(t, string(wire), "/private/")
	require.Contains(t, string(wire), `"device_name":"Meter A"`)
	require.Contains(t, string(wire), `"status":"failed"`)
	require.Equal(t, hostile, response.DatabaseDelivery[0].Error, "projection must not mutate runtime evidence")
	require.Equal(t, hostile, response.Diagnostics[0].Stages[0].Reason)
}

func TestRuntimeRuleStreamSafetyLegacyValuesAndStatus(t *testing.T) {
	hostile := "postgres://operator:password-fixture@private-host:5432/private-db"
	for _, tc := range []struct {
		name    string
		event   string
		payload any
	}{
		{"status", "status", datalinkruntime.DeviceStatusEvent{DeviceID: "device-1", LastError: &hostile, AvailabilityReason: &hostile, ProjectionMessage: hostile}},
		{"nested-value", "value", datalinkruntime.ValueEvent{DeviceID: "device-1", PointID: "point-1", Address: "40001", RawValue: map[string]any{"safe": []any{map[string]any{"password": hostile}}}, TransformedValue: hostile, Quality: schema.QualityGood}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			require.NoError(t, writeRuntimeSSE(response, response, tc.event, tc.payload))
			require.NotContains(t, response.Body.String(), "password-fixture")
			require.NotContains(t, response.Body.String(), "private-host")
			if tc.event == "value" {
				require.Contains(t, response.Body.String(), `"address":"40001"`)
				require.Contains(t, response.Body.String(), `"value_state":"unsupported"`)
			}
		})
	}
}

func TestRuntimeRuleStreamSafetyPreservesSensorScalars(t *testing.T) {
	for _, tc := range []struct {
		value any
		wire  string
	}{
		{37.5, `37.5`}, {int64(9223372036854775807), `"9223372036854775807"`}, {uint64(18446744073709551615), `"18446744073709551615"`}, {true, `true`}, {"running", `"running"`}, {"SELECT mode", `"SELECT mode"`}, {"Driver status", `"Driver status"`}, {"OK/NG", `"OK/NG"`}, {"23.5 °C", `"23.5 °C"`}, {"2026/10/28", `"2026/10/28"`}, {nil, `null`},
	} {
		response := httptest.NewRecorder()
		event := datalinkruntime.ValueEvent{DeviceID: "device-1", PointID: "point-1", Address: "40001", RawValue: tc.value, TransformedValue: tc.value, Timestamp: time.Date(2026, 10, 28, 0, 0, 0, 0, time.UTC)}
		require.NoError(t, writeRuntimeSSE(response, response, "value", event))
		require.True(t, strings.Contains(response.Body.String(), `"raw_value":`+tc.wire), response.Body.String())
	}
}

func TestRuntimeRuleStreamSafetyHostileMetadata(t *testing.T) {
	for _, hostile := range []string{"private.example.com", "192.0.2.1", "SELECT secret FROM readings", "/private/file.db", "password=fixture"} {
		response := httptest.NewRecorder()
		event := datalinkruntime.ValueEvent{DeviceID: "device-1", PointID: "point-1", Address: hostile, RawValue: hostile}
		require.NoError(t, writeRuntimeSSE(response, response, "value", event))
		require.NotContains(t, response.Body.String(), hostile)
	}
}

func TestRuntimeRuleStatusSafetyHostileEnums(t *testing.T) {
	hostile := "postgres://operator:password-fixture@private-host:5432/private-db"
	response := runtimeStatusResponse{
		DatabaseDelivery:    []datalinkruntime.DatabaseDeliveryDiagnostic{{Status: datalinkruntime.DatabaseDeliveryStatus(hostile), FailedStage: datalinkruntime.DatabaseDeliveryStage(hostile), Stages: []datalinkruntime.DatabaseDeliveryStage{datalinkruntime.DatabaseDeliveryStage(hostile)}}},
		ModbusShareDelivery: []datalinkruntime.ModbusShareDeliveryDiagnostic{{Status: datalinkruntime.ModbusShareDeliveryStatus(hostile)}},
		Diagnostics:         []datalinkruntime.RuntimeFlowDiagnostic{{FailureStage: datalinkruntime.RuntimeFlowStage(hostile), LatestSuccessfulStage: datalinkruntime.RuntimeFlowStage(hostile), Stages: []datalinkruntime.RuntimeFlowDiagnosticStage{{Stage: datalinkruntime.RuntimeFlowStage(hostile), Status: datalinkruntime.RuntimeFlowStageStatus(hostile)}}}},
	}
	wire, err := json.Marshal(response)
	require.NoError(t, err)
	require.NotContains(t, string(wire), "password-fixture")
}

func TestRuntimeRuleStreamSafetyRemoteRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	source := &stubRuntimeStreamSource{valueCh: make(chan datalinkruntime.ValueEvent, 1), statusCh: make(chan datalinkruntime.DeviceStatusEvent)}
	source.valueCh <- datalinkruntime.ValueEvent{DeviceID: "device-1", PointID: "point-1", Address: "40001", RawValue: map[string]any{"private_path": "/private/password-fixture.db"}}
	close(source.valueCh)
	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)
	c.Request = httptest.NewRequest(http.MethodGet, "/runtime/stream?device_id=postgres%3A%2F%2Fpassword-fixture", http.NoBody)
	c.Request.RemoteAddr = "203.0.113.10:4567"
	NewRuntimeStreamHandler(source).Stream(c)
	require.Equal(t, http.StatusOK, response.Code)
	require.NotContains(t, response.Body.String(), "password-fixture")
	require.Contains(t, response.Body.String(), `"address":"40001"`)
	require.Contains(t, response.Body.String(), `"value_state":"unsupported"`)
}

func TestRuntimeRuleStatusSafetyRemoteRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	now := time.Now()
	service, err := datalinkruntime.NewService(datalinkruntime.Config{Writer: storage.NewMemoryStorage(4), Snapshot: datalinkruntime.Snapshot{
		Devices: []*schema.Device{{ID: "device-1", Name: "Meter A", Status: schema.DeviceStatusActive}},
		Points:  []*schema.Point{{ID: "point-1", DeviceID: "device-1", LastReadAt: &now, LastError: "driver password-fixture /private/file.db"}},
	}})
	require.NoError(t, err)
	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)
	c.Request = httptest.NewRequest(http.MethodGet, "/runtime/status", http.NoBody)
	c.Request.RemoteAddr = "203.0.113.10:4567"
	NewRuntimeHandler(nil, nil, nil, nil, service, nil).Status(c)
	require.Equal(t, http.StatusOK, response.Code)
	require.NotContains(t, response.Body.String(), "password-fixture")
	require.Contains(t, response.Body.String(), `"device_name":"Meter A"`)
}

func TestRuntimeRuleStreamSafetyUnavailableRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)
	c.Request = httptest.NewRequest(http.MethodGet, "/runtime/stream?device_id=postgres%3A%2F%2Fpassword-fixture", http.NoBody)
	c.Request.RemoteAddr = "203.0.113.10:4567"
	NewRuntimeStreamHandler(&stubRuntimeStreamSource{}).Stream(c)
	require.Equal(t, http.StatusServiceUnavailable, response.Code)
	require.NotContains(t, response.Body.String(), "password-fixture")
	require.Contains(t, response.Body.String(), `"unavailable":true`)
}

func TestRuntimeRuleStatusSafetyUnknownStatesAndSafeScope(t *testing.T) {
	for _, input := range []string{"", "future-state", "postgres://password-fixture@private-host"} {
		response := runtimeStatusResponse{SnapshotState: datalinkruntime.RuntimeTruthState{State: input}, Collectors: []runtimeCollectorResponse{{Protocol: input, Status: input, AvailabilityStatus: input, BreakerState: input, ProjectionAlignment: input, Running: true}}, Diagnostics: []datalinkruntime.RuntimeFlowDiagnostic{{Scope: "device:dev-A"}, {Scope: "point:point-1"}, {Scope: "tag:tag-1"}, {Scope: "runtime"}, {Scope: "device:postgres://password-fixture"}}}
		raw, err := json.Marshal(response)
		require.NoError(t, err)
		var wire struct {
			SnapshotState datalinkruntime.RuntimeTruthState       `json:"snapshot_state"`
			Collectors    []runtimeCollectorResponse              `json:"collectors"`
			Diagnostics   []datalinkruntime.RuntimeFlowDiagnostic `json:"diagnostics"`
		}
		require.NoError(t, json.Unmarshal(raw, &wire))
		require.Equal(t, "unavailable", wire.SnapshotState.State)
		require.True(t, wire.SnapshotState.Unavailable)
		require.Equal(t, "unavailable", wire.Collectors[0].Status)
		require.Equal(t, "unavailable", wire.Collectors[0].Protocol)
		require.Equal(t, "unavailable", wire.Collectors[0].AvailabilityStatus)
		require.Equal(t, "unavailable", wire.Collectors[0].BreakerState)
		require.Equal(t, "unavailable", wire.Collectors[0].ProjectionAlignment)
		require.False(t, wire.Collectors[0].Running)
		require.Equal(t, "device:dev-A", wire.Diagnostics[0].Scope)
		require.Equal(t, "point:point-1", wire.Diagnostics[1].Scope)
		require.Equal(t, "tag:tag-1", wire.Diagnostics[2].Scope)
		require.Equal(t, "runtime", wire.Diagnostics[3].Scope)
		require.Empty(t, wire.Diagnostics[4].Scope)
		require.NotContains(t, string(raw), "password-fixture")
	}
}

func TestRuntimeRuleStatusSafetyPreservesKnownStates(t *testing.T) {
	for _, alignment := range []string{"aligned", "deferred", "restart-required", "restart_required", "stale", "degraded"} {
		response := runtimeStatusResponse{SnapshotState: datalinkruntime.RuntimeReadyTruthState(), Collectors: []runtimeCollectorResponse{{Protocol: "fatek_fbs", Status: "running", AvailabilityStatus: "available", BreakerState: "half-open", ProjectionAlignment: alignment, Running: true}}}
		raw, err := json.Marshal(response)
		require.NoError(t, err)
		require.Contains(t, string(raw), `"state":"ready"`)
		require.Contains(t, string(raw), `"protocol":"fatek_fbs"`)
		require.Contains(t, string(raw), `"status":"running"`)
		require.Contains(t, string(raw), `"availability_status":"available"`)
		require.Contains(t, string(raw), `"breaker_state":"half-open"`)
		require.Contains(t, string(raw), `"projection_alignment":"`+alignment+`"`)
	}
}

func TestRuntimeRuleStatusSafetyScopeCompatibility(t *testing.T) {
	raw, err := json.Marshal(runtimeStatusResponse{Diagnostics: []datalinkruntime.RuntimeFlowDiagnostic{{Scope: "device:dev-A"}, {Scope: "point:point-1"}, {Scope: "tag:tag-1"}, {Scope: "runtime"}}})
	require.NoError(t, err)
	for _, scope := range []string{"device:dev-A", "point:point-1", "tag:tag-1", "runtime"} {
		require.Contains(t, string(raw), `"scope":"`+scope+`"`)
	}
}
