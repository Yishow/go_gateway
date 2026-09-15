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

	var response map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))

	data := response["data"].([]interface{})
	expectedTypes := []string{"modbus_tcp", "modbus_rtu", "modbus_udp"}
	expectedEnum := []interface{}{"ABCD", "BADC", "CDAB", "DCBA"}

	for _, protocolType := range expectedTypes {
		t.Run(protocolType, func(t *testing.T) {
			t.Parallel()

			var protocol map[string]interface{}
			for _, item := range data {
				candidate := item.(map[string]interface{})
				if candidate["type"] == protocolType {
					protocol = candidate
					break
				}
			}
			require.NotNil(t, protocol, "protocol %s not found", protocolType)

			configSchema := protocol["config_schema"].(string)
			var schema map[string]interface{}
			require.NoError(t, json.Unmarshal([]byte(configSchema), &schema))

			properties := schema["properties"].(map[string]interface{})
			dataFormatValue, ok := properties["data_format"]
			require.True(t, ok, "%s schema should expose data_format", protocolType)

			dataFormat := dataFormatValue.(map[string]interface{})
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

	var response map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))

	data := response["data"].([]interface{})
	var mqtt map[string]interface{}
	for _, item := range data {
		candidate := item.(map[string]interface{})
		if candidate["type"] == "mqtt" {
			mqtt = candidate
			break
		}
	}
	require.NotNil(t, mqtt, "protocol mqtt not found")

	configSchema := mqtt["config_schema"].(string)
	var schema map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(configSchema), &schema))

	properties := schema["properties"].(map[string]interface{})
	assert.NotContains(t, properties, "data_format")
}
