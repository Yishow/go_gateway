package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/point"
	datalinkruntime "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRuntimeRuleNetworkSafetyRejectsIPv6AndPairedDSN(t *testing.T) {
	for _, value := range []string{"[2001:db8::1]:5432", "[::1]:502", "[fe80::1%eth0]:502", "Server=localhost;Database=metrics;Integrated Security=true", "Data Source=localhost;Initial Catalog=metrics", "host=localhost dbname=metrics", "host=localhost port=5432", "Server=localhost;User Id=operator"} {
		t.Run(value, func(t *testing.T) {
			require.False(t, runtimeSafeSensorText(value))
		})
	}
}
func TestRuntimeRuleNetworkSafetyPreservesOrdinaryScalarsAndAddresses(t *testing.T) {
	for _, value := range []string{"Server mode", "Server=ready", "Database=metrics", "m/s", "OK/NG", "40001/16", "12:34:56", "2026-10-29T01:02:03.456Z", "ratio 1:2", "[normal]", "[12:30:00]", "SELECT mode", "23.5 °C"} {
		require.True(t, runtimeSafeSensorText(value), value)
	}
	for _, address := range []string{"40001", "BFF", "SB1A", "W1A", "D100", "R100.2"} {
		require.Equal(t, address, runtimeSafeAddress(address))
	}
}
func TestRuntimeRuleNetworkSafetyRemotePointAndSSE(t *testing.T) {
	for _, value := range []string{"[2001:db8::1]:5432", "[::1]:502", "[fe80::1%eth0]:502", "Server=localhost;Database=metrics;Integrated Security=true"} {
		t.Run(value, func(t *testing.T) {
			repo := point.NewMemoryRepository()
			encoded, _ := json.Marshal(value)
			raw := string(encoded)
			require.NoError(t, repo.Create(t.Context(), &schema.Point{ID: "point", DeviceID: "device", Address: "40001", DataType: schema.DataTypeString, LastValue: &raw}))
			handler := NewPointHandler(point.NewService(repo, nil))
			for _, get := range []bool{false, true} {
				response := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(response)
				c.Request = httptest.NewRequest(http.MethodGet, "/points?device_id=device", nil)
				c.Request.RemoteAddr = "203.0.113.9:1234"
				c.Params = gin.Params{{Key: "id", Value: "point"}}
				if get {
					handler.Get(c)
				} else {
					handler.List(c)
				}
				require.Equal(t, 200, response.Code)
				require.NotContains(t, response.Body.String(), value)
				require.Contains(t, response.Body.String(), `"value_state":"unsupported"`)
				require.Contains(t, response.Body.String(), `"address":"40001"`)
			}
			saved, err := repo.GetByID(t.Context(), "point")
			require.NoError(t, err)
			require.Equal(t, raw, *saved.LastValue)
			source := &stubRuntimeStreamSource{valueCh: make(chan datalinkruntime.ValueEvent, 1), statusCh: make(chan datalinkruntime.DeviceStatusEvent)}
			source.valueCh <- datalinkruntime.ValueEvent{DeviceID: "device", PointID: "point", Address: "40001", RawValue: value, TransformedValue: value, Quality: schema.QualityGood}
			close(source.valueCh)
			response := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(response)
			c.Request = httptest.NewRequest(http.MethodGet, "/runtime/stream?device_id=device", nil)
			c.Request.RemoteAddr = "203.0.113.9:1234"
			NewRuntimeStreamHandler(source).Stream(c)
			require.Equal(t, 200, response.Code)
			require.NotContains(t, response.Body.String(), value)
			require.Contains(t, response.Body.String(), `"value_state":"unsupported"`)
		})
	}
}
