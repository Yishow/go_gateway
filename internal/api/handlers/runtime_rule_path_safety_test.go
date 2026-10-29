package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/point"
	datalinkruntime "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
)

func TestRuntimeRulePathSafetyDelimitedAndUNC(t *testing.T) {
	cases := []string{
		"path=/private/db", "path=C:/private/db", "(/private/db)", `"/private/db"`, "'/private/db'", "[/private/db]", `{"path":"/private/db"}`, "`/private/db`",
		`\\server\share\private`, `path=\\server\share\private`, "./private/db", "../private/db", "~/private/db", "path=./private/db", "(../private/db)", "[~/private/db]",
		`.\private\db`, `..\private\db`, "./", "../", "~/", "path:../private/db", "reading;/private/db", "reading,/private/db", ")/private/db", "]/private/db", "}/private/db",
	}
	for _, value := range cases {
		t.Run(value, func(t *testing.T) {
			t.Run("sse", func(t *testing.T) {
				rec := httptest.NewRecorder()
				require.NoError(t, writeRuntimeSSE(rec, rec, "value", datalinkruntime.ValueEvent{DeviceID: "device", PointID: "point", Address: "40001", RawValue: value, TransformedValue: value, Quality: schema.QualityGood}))
				var payload map[string]any
				require.NoError(t, json.Unmarshal([]byte(strings.TrimSpace(strings.SplitN(rec.Body.String(), "data: ", 2)[1])), &payload))
				require.Nil(t, payload["raw_value"])
				require.Nil(t, payload["transformed_value"])
				require.Equal(t, "unsupported", payload["value_state"])
				require.Equal(t, "invalid", payload["quality"])
			})
			t.Run("point-read", func(t *testing.T) {
				raw, err := json.Marshal(value)
				require.NoError(t, err)
				stored := string(raw)
				repo := point.NewMemoryRepository()
				require.NoError(t, repo.Create(t.Context(), &schema.Point{ID: "point", DeviceID: "device", Name: value, DataType: schema.DataTypeString, Address: "40001", LastValue: &stored}))
				rec := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rec)
				c.Request = httptest.NewRequest(http.MethodGet, "/points", nil)
				c.Request.RemoteAddr = "203.0.113.7:1234"
				NewPointHandler(point.NewService(repo, nil)).List(c)
				var payload struct {
					Data []struct {
						LastValue  *string `json:"last_value"`
						Name       string  `json:"name"`
						ValueState string  `json:"value_state"`
					} `json:"data"`
				}
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &payload))
				require.Len(t, payload.Data, 1)
				require.Nil(t, payload.Data[0].LastValue)
				require.Empty(t, payload.Data[0].Name)
				require.Equal(t, "unsupported", payload.Data[0].ValueState)
				saved, err := repo.GetByID(t.Context(), "point")
				require.NoError(t, err)
				require.Equal(t, stored, *saved.LastValue)
			})
		})
	}
}

func TestRuntimeRulePathSafetyPreservesOrdinaryReadings(t *testing.T) {
	for _, value := range []string{"m/s", "OK/NG", "40001/16", "C temperature", "C: temperature", "23.5 °C", "Pump A (°C)", "(normal reading)", "[normal]", "reading=value", `\\sensor`} {
		t.Run(value, func(t *testing.T) {
			rec := httptest.NewRecorder()
			require.NoError(t, writeRuntimeSSE(rec, rec, "value", datalinkruntime.ValueEvent{DeviceID: "device", PointID: "point", Address: "40001", RawValue: value, TransformedValue: value, Quality: schema.QualityGood}))
			var payload map[string]any
			require.NoError(t, json.Unmarshal([]byte(strings.TrimSpace(strings.SplitN(rec.Body.String(), "data: ", 2)[1])), &payload))
			require.Equal(t, value, payload["raw_value"])
			require.Equal(t, "good", payload["quality"])
			require.NotContains(t, payload, "value_state")
		})
	}
}
