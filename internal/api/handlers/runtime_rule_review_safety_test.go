package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/device"
	datalinkruntime "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

func TestRuntimeRuleReviewSafetyDNSAndConnectionSyntax(t *testing.T) {
	for _, value := range []string{"db.private.internal", "uid=fixture-user;pwd=fixture-password", "file:private-fixture.sqlite", "file:/private/fixture.db"} {
		t.Run(value, func(t *testing.T) {
			response := httptest.NewRecorder()
			require.NoError(t, writeRuntimeSSE(response, response, "value", datalinkruntime.ValueEvent{DeviceID: value, PointID: value, Address: "40001", RawValue: value, TransformedValue: value}))
			require.NotContains(t, response.Body.String(), value)
			require.Contains(t, response.Body.String(), `"value_state":"unsupported"`)
		})
	}
}

func TestRuntimeRuleReviewSafetyPreservesMCAndUnicode(t *testing.T) {
	for _, address := range []string{"BFF", "SB1A", "W1A", "SWFF", "D0100", "40001"} {
		response := httptest.NewRecorder()
		require.NoError(t, writeRuntimeSSE(response, response, "value", datalinkruntime.ValueEvent{DeviceID: "device-1", PointID: "point-1", Address: address, RawValue: 23.5, TransformedValue: 23.5, Quality: schema.QualityGood}))
		require.Contains(t, response.Body.String(), `"address":"`+address+`"`)
		require.NotContains(t, response.Body.String(), `"value_state":"unsupported"`)
	}
	for _, name := range []string{"泵浦 A (°C)", "Cafe\u0301 A"} {
		actual := mapRuntimeWorkspaceContextDevice(&schema.Device{ID: "device-1", Name: name, Protocol: schema.ProtocolMC3E}, datalinkruntime.DeviceRuntimeStatus{})
		require.Equal(t, name, actual.Name)
	}
}

func TestRuntimeRuleReviewSafetyContextUnknownIsNotRunningOrDefault(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, status := range []string{"future-availability", "db.private.internal"} {
		ws := workspace.NewService(workspace.NewMemoryRepository())
		_, err := ws.AttachDevice(t.Context(), "device-1")
		require.NoError(t, err)
		repo := device.NewMemoryRepository()
		readiness, err := json.Marshal(map[string]string{"availability_status": status})
		require.NoError(t, err)
		require.NoError(t, repo.Create(t.Context(), &schema.Device{ID: "device-1", Name: "Meter", Protocol: schema.ProtocolMC3E, Status: schema.DeviceStatusActive, ReadinessStatus: string(readiness)}))
		response := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(response)
		c.Request = httptest.NewRequest(http.MethodGet, "/runtime-context", http.NoBody)
		NewRuntimeHandler(device.NewService(repo, nil), nil, nil, nil, nil, ws).WorkspaceContext(c)
		require.Equal(t, http.StatusOK, response.Code, response.Body.String())
		require.Contains(t, response.Body.String(), `"availability_status":"unavailable"`)
		require.Contains(t, response.Body.String(), `"running":false`)
		require.Contains(t, response.Body.String(), `"default_device_id":null`)
		require.NotContains(t, response.Body.String(), status)
	}
	for _, status := range []string{"", "future-availability"} {
		actual := mapRuntimeWorkspaceContextDevice(&schema.Device{ID: "device-1", Protocol: schema.ProtocolMC3E}, datalinkruntime.DeviceRuntimeStatus{DeviceID: "device-1", AvailabilityStatus: status, Running: true, ProjectionAlignment: "future-alignment"})
		require.Equal(t, "unavailable", actual.AvailabilityStatus)
		require.Equal(t, "unavailable", actual.ProjectionAlignment)
		require.False(t, actual.Running)
	}
}

func TestRuntimeRuleReviewSafetyRemoteStreamRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, available := range []bool{false, true} {
		for _, hostile := range []string{"db.private.internal", "file:private-fixture.sqlite", "uid=fixture;pwd=fixture"} {
			source := &stubRuntimeStreamSource{}
			if available {
				source.valueCh = make(chan datalinkruntime.ValueEvent)
				source.statusCh = make(chan datalinkruntime.DeviceStatusEvent)
				close(source.valueCh)
			}
			response := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(response)
			c.Request = httptest.NewRequest(http.MethodGet, "/runtime/stream?device_id=device-1", http.NoBody)
			c.Request.RemoteAddr = "203.0.113.10:4567"
			c.Request.Header.Set("X-Request-ID", hostile)
			NewRuntimeStreamHandler(source).Stream(c)
			require.NotContains(t, response.Body.String(), hostile)
			require.Contains(t, response.Body.String(), `"request_id":"req-`)
		}
	}
}

func TestRuntimeRuleReviewSafetyUnicodeLabels(t *testing.T) {
	for _, name := range []string{"泵浦 A (°C)", "Cafe\u0301 A"} {
		actual := mapRuntimeWorkspaceContextDevice(&schema.Device{ID: "device-1", Name: name, Protocol: schema.ProtocolMC3E}, datalinkruntime.DeviceRuntimeStatus{})
		require.Equal(t, name, actual.Name)
	}
}

func TestRuntimeRuleReviewSafetyRelativeFileURI(t *testing.T) {
	response := httptest.NewRecorder()
	require.NoError(t, writeRuntimeSSE(response, response, "value", datalinkruntime.ValueEvent{DeviceID: "device-1", PointID: "point-1", Address: "40001", RawValue: "file:fixture", TransformedValue: 23.5}))
	require.NotContains(t, response.Body.String(), "file:fixture")
	require.Contains(t, response.Body.String(), `"value_state":"unsupported"`)
}

func TestRuntimeRuleReviewSafetyMissingValueIdentity(t *testing.T) {
	for _, event := range []datalinkruntime.ValueEvent{
		{DeviceID: "db.private.internal", PointID: "point-1", Address: "40001", RawValue: 23.5, Quality: schema.QualityGood},
		{DeviceID: "device-1", PointID: "", Address: "40001", RawValue: 23.5, Quality: schema.QualityGood},
		{DeviceID: "device-1", PointID: "point-1", Address: "", RawValue: 23.5, Quality: schema.QualityGood},
	} {
		response := httptest.NewRecorder()
		require.NoError(t, writeRuntimeSSE(response, response, "value", event))
		require.Contains(t, response.Body.String(), `"value_state":"unsupported"`)
		require.Contains(t, response.Body.String(), `"quality":"invalid"`)
	}
}

func TestRuntimeRuleReviewSafetyMCAddressBounds(t *testing.T) {
	for _, address := range []string{"W80000000", "WFFFFFFFF", "SBFFFFFFFF"} {
		response := httptest.NewRecorder()
		require.NoError(t, writeRuntimeSSE(response, response, "value", datalinkruntime.ValueEvent{DeviceID: "device-1", PointID: "point-1", Address: address, RawValue: 23.5, Quality: schema.QualityGood}))
		require.Contains(t, response.Body.String(), `"address":""`)
		require.Contains(t, response.Body.String(), `"value_state":"unsupported"`)
	}
	require.Equal(t, "W7FFFFFFF", runtimeSafeAddress("W7FFFFFFF"))
}

func TestRuntimeRuleReviewSafetyWindowsForwardSlashPath(t *testing.T) {
	for _, value := range []string{"C:/private/fixture.sqlite", "D:/plant/private"} {
		response := httptest.NewRecorder()
		require.NoError(t, writeRuntimeSSE(response, response, "value", datalinkruntime.ValueEvent{DeviceID: "device-1", PointID: "point-1", Address: "40001", RawValue: value, TransformedValue: 23.5}))
		require.NotContains(t, response.Body.String(), value)
		require.Contains(t, response.Body.String(), `"value_state":"unsupported"`)
		_, supported := runtimeSafeSensorValue(value)
		require.False(t, supported)
	}
	for _, value := range []string{"23.5 °C", "C temperature", "C: temperature"} {
		actual, supported := runtimeSafeSensorValue(value)
		require.True(t, supported)
		require.Equal(t, value, actual)
	}
}

func TestRuntimeRuleReviewSafetyUnsupportedCollectorCannotRun(t *testing.T) {
	for _, input := range []runtimeCollectorResponse{
		{DeviceID: "device-1", Protocol: "future-protocol", Status: "running", AvailabilityStatus: "available", Running: true},
		{DeviceID: "db.private.internal", Protocol: "mc_3e", Status: "running", AvailabilityStatus: "available", Running: true},
		{DeviceID: "", Protocol: "mc_3e", Status: "running", AvailabilityStatus: "available", Running: true},
	} {
		raw, err := json.Marshal(runtimeStatusResponse{Collectors: []runtimeCollectorResponse{input}})
		require.NoError(t, err)
		require.Contains(t, string(raw), `"running":false`)
		var result struct {
			Collectors []runtimeCollectorResponse `json:"collectors"`
		}
		require.NoError(t, json.Unmarshal(raw, &result))
		require.False(t, result.Collectors[0].Running)
	}
}
