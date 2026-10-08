package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProtocolHandler_List_ConfigSchema_ModbusVariantsExposeDataFormat(t *testing.T) {
	t.Parallel()

	r := setupProtocolRouter()
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "/datalink/protocols", http.NoBody)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))

	data := response["data"].([]any)
	expectedTypes := []string{"modbus_tcp", "modbus_rtu", "modbus_udp"}
	expectedEnum := []any{"ABCD", "BADC", "CDAB", "DCBA"}

	for _, protocolType := range expectedTypes {
		t.Run(protocolType, func(t *testing.T) {
			t.Parallel()

			var protocol map[string]any
			for _, item := range data {
				candidate := item.(map[string]any)
				if candidate["type"] == protocolType {
					protocol = candidate
					break
				}
			}
			require.NotNil(t, protocol, "protocol %s not found", protocolType)

			configSchema := protocol["config_schema"].(string)
			var schema map[string]any
			require.NoError(t, json.Unmarshal([]byte(configSchema), &schema))

			properties := schema["properties"].(map[string]any)
			dataFormatValue, ok := properties["data_format"]
			require.True(t, ok, "%s schema should expose data_format", protocolType)

			dataFormat := dataFormatValue.(map[string]any)
			assert.Equal(t, "string", dataFormat["type"])
			assert.Equal(t, expectedEnum, dataFormat["enum"])
			assert.Equal(t, "ABCD", dataFormat["default"])
		})
	}
}

func TestProtocolHandler_List_ConfigSchema_MQTTOmitsDataFormat(t *testing.T) {
	t.Parallel()

	r := setupProtocolRouter()
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "/datalink/protocols", http.NoBody)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))

	data := response["data"].([]any)
	var mqtt map[string]any
	for _, item := range data {
		candidate := item.(map[string]any)
		if candidate["type"] == "mqtt" {
			mqtt = candidate
			break
		}
	}
	require.NotNil(t, mqtt, "protocol mqtt not found")

	configSchema := mqtt["config_schema"].(string)
	var schema map[string]any
	require.NoError(t, json.Unmarshal([]byte(configSchema), &schema))

	properties := schema["properties"].(map[string]any)
	assert.NotContains(t, properties, "data_format")
}
