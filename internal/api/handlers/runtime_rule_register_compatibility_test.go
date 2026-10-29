package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-gateway/internal/datalink/point"
	datalinkruntime "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestRuntimeSourcePagePreservesSupportedRegister(t *testing.T) {
	for _, address := range []string{"D0000000", "W000000001A", "00040001", "0x000000FFFF", "D1000000", "D16777215", "M1000000"} {
		t.Run(address, func(t *testing.T) {
			f, query := runtimeSourcePageFixture(t, 1)
			f.h.sourceRuleReader = runtimeSetupRuleReaderStub{rules: []*schema.SourceRule{{
				ID: "rule-padded", DeviceID: "device-1", RevisionID: "R1",
				StartAddress: address, Count: 1, DataType: schema.DataTypeInt16, Enabled: true,
			}}}

			page := runtimeSourcePageGet(t, f.h, query)
			require.Len(t, page.Rules, 1)
			require.Equal(t, address, page.Rules[0]["start_address"])
		})
	}
}

func TestRuntimeSupportedRegisterPointAndStream(t *testing.T) {
	for _, address := range []string{"D0000000", "W000000001", "00040001", "D1000000", "D16777215"} {
		t.Run(address, func(t *testing.T) {
			points := point.NewService(point.NewMemoryRepository(), nil)
			saved, err := points.Create(t.Context(), point.CreatePointRequest{
				DeviceID: "device", Name: "Point", Address: address, DataType: schema.DataTypeInt16,
			})
			require.NoError(t, err)
			response := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(response)
			ctx.Request = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/points/"+saved.ID, http.NoBody)
			ctx.Request.RemoteAddr = "203.0.113.9:1234"
			ctx.Params = gin.Params{{Key: "id", Value: saved.ID}}
			NewPointHandler(points).Get(ctx)
			require.Equal(t, http.StatusOK, response.Code)
			require.Contains(t, response.Body.String(), `"address":"`+address+`"`)
			require.NotContains(t, response.Body.String(), `"metadata_state":"unsupported"`)

			event := datalinkruntime.ValueEvent{DeviceID: "device", PointID: saved.ID, Address: address,
				RawValue: 1, TransformedValue: 1, Quality: schema.QualityGood}
			body := runtimePaddedAddressStream(t, event)
			require.Contains(t, body, `"address":"`+address+`"`)
			require.Contains(t, body, `"quality":"good"`)
			require.NotContains(t, body, `"value_state":"unsupported"`)
		})
	}
}

func TestRuntimeRegisterProjectionRetainsSafetyBounds(t *testing.T) {
	for _, address := range []string{"W00080000000", "SB000FFFFFFFF", "D16777216", "HOST1000000", "ſ1000000", "D0000000:502", "D" + strings.Repeat("0", 128)} {
		t.Run(address, func(t *testing.T) {
			event := datalinkruntime.ValueEvent{DeviceID: "device", PointID: "point", Address: address,
				RawValue: 1, TransformedValue: 1, Quality: schema.QualityGood}
			body := runtimePaddedAddressStream(t, event)
			require.NotContains(t, body, address)
			require.Contains(t, body, `"address":""`)
			require.Contains(t, body, `"quality":"invalid"`)
			require.Contains(t, body, `"value_state":"unsupported"`)
		})
	}
}

func runtimePaddedAddressStream(t *testing.T, event datalinkruntime.ValueEvent) string {
	t.Helper()
	source := &stubRuntimeStreamSource{valueCh: make(chan datalinkruntime.ValueEvent, 1), statusCh: make(chan datalinkruntime.DeviceStatusEvent)}
	source.valueCh <- event
	close(source.valueCh)
	response := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(response)
	ctx.Request = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/runtime/stream?device_id=device", http.NoBody)
	ctx.Request.RemoteAddr = "203.0.113.9:1234"
	NewRuntimeStreamHandler(source).Stream(ctx)
	require.Equal(t, http.StatusOK, response.Code)
	return response.Body.String()
}
