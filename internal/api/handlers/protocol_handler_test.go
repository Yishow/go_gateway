// Package handlers 提供 Protocol API Handler 的單元測試。
package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/**
 * setupProtocolRouter 建立測試用的 Protocol Router
 * @param t 測試實例
 * @returns *gin.Engine 測試路由器
 */
func setupProtocolRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h := NewProtocolHandler()
	r.GET("/datalink/protocols", h.List)
	return r
}

/**
 * TestProtocolHandler_List 測試列出所有協議
 */
func TestProtocolHandler_List(t *testing.T) {
	r := setupProtocolRouter()

	req, _ := http.NewRequest("GET", "/datalink/protocols", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].([]interface{})
	assert.NotEmpty(t, data)
}

/**
 * TestProtocolHandler_List_Count 測試協議數量
 */
func TestProtocolHandler_List_Count(t *testing.T) {
	r := setupProtocolRouter()

	req, _ := http.NewRequest("GET", "/datalink/protocols", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].([]interface{})

	// 應該有 6 個協議
	assert.Equal(t, 6, len(data))
}

/**
 * TestProtocolHandler_List_Structure 測試回應結構
 */
func TestProtocolHandler_List_Structure(t *testing.T) {
	r := setupProtocolRouter()

	req, _ := http.NewRequest("GET", "/datalink/protocols", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// 驗證頂層結構
	assert.True(t, response["success"].(bool))
	assert.NotNil(t, response["data"])

	// 驗證協議陣列
	data := response["data"].([]interface{})
	assert.Greater(t, len(data), 0)

	// 驗證第一個協議的結構
	firstProtocol := data[0].(map[string]interface{})
	assert.Contains(t, firstProtocol, "type")
	assert.Contains(t, firstProtocol, "name")
	assert.Contains(t, firstProtocol, "description")
	assert.Contains(t, firstProtocol, "config_schema")
}

/**
 * TestProtocolHandler_List_ModbusTCP 測試 Modbus TCP 協議
 */
func TestProtocolHandler_List_ModbusTCP(t *testing.T) {
	r := setupProtocolRouter()

	req, _ := http.NewRequest("GET", "/datalink/protocols", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].([]interface{})

	// 尋找 Modbus TCP 協議
	var modbusTCP map[string]interface{}
	for _, item := range data {
		protocol := item.(map[string]interface{})
		if protocol["type"] == "modbus_tcp" {
			modbusTCP = protocol
			break
		}
	}

	require.NotNil(t, modbusTCP, "Modbus TCP protocol not found")

	// 驗證欄位
	assert.Equal(t, "modbus_tcp", modbusTCP["type"])
	assert.Equal(t, "Modbus TCP", modbusTCP["name"])
	assert.Equal(t, "Modbus TCP/IP 協議，用於乙太網路連接的 PLC 和設備", modbusTCP["description"])
	assert.NotNil(t, modbusTCP["config_schema"])

	// 驗證 Config Schema 是有效的 JSON
	configSchema, ok := modbusTCP["config_schema"].(string)
	assert.True(t, ok)

	var schema map[string]interface{}
	err := json.Unmarshal([]byte(configSchema), &schema)
	assert.NoError(t, err)
	assert.Equal(t, "object", schema["type"])
	assert.NotNil(t, schema["properties"])
	assert.NotNil(t, schema["required"])
}

/**
 * TestProtocolHandler_List_ModbusRTU 測試 Modbus RTU 協議
 */
func TestProtocolHandler_List_ModbusRTU(t *testing.T) {
	r := setupProtocolRouter()

	req, _ := http.NewRequest("GET", "/datalink/protocols", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].([]interface{})

	// 尋找 Modbus RTU 協議
	var modbusRTU map[string]interface{}
	for _, item := range data {
		protocol := item.(map[string]interface{})
		if protocol["type"] == "modbus_rtu" {
			modbusRTU = protocol
			break
		}
	}

	require.NotNil(t, modbusRTU, "Modbus RTU protocol not found")

	assert.Equal(t, "modbus_rtu", modbusRTU["type"])
	assert.Equal(t, "Modbus RTU", modbusRTU["name"])
	assert.Equal(t, "Modbus RTU 協議，用於 RS-232/RS-485 串列通訊", modbusRTU["description"])
}

/**
 * TestProtocolHandler_List_ModbusUDP 測試 Modbus UDP 協議
 */
func TestProtocolHandler_List_ModbusUDP(t *testing.T) {
	r := setupProtocolRouter()

	req, _ := http.NewRequest("GET", "/datalink/protocols", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].([]interface{})

	// 尋找 Modbus UDP 協議
	var modbusUDP map[string]interface{}
	for _, item := range data {
		protocol := item.(map[string]interface{})
		if protocol["type"] == "modbus_udp" {
			modbusUDP = protocol
			break
		}
	}

	require.NotNil(t, modbusUDP, "Modbus UDP protocol not found")

	assert.Equal(t, "modbus_udp", modbusUDP["type"])
	assert.Equal(t, "Modbus UDP", modbusUDP["name"])
	assert.Equal(t, "Modbus UDP 協議，用於 UDP 通訊", modbusUDP["description"])
}

/**
 * TestProtocolHandler_List_FatekFBs 測試 FATEK FBs 協議
 */
func TestProtocolHandler_List_FatekFBs(t *testing.T) {
	r := setupProtocolRouter()

	req, _ := http.NewRequest("GET", "/datalink/protocols", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].([]interface{})

	// 尋找 FATEK FBs 協議
	var fatekFBs map[string]interface{}
	for _, item := range data {
		protocol := item.(map[string]interface{})
		if protocol["type"] == "fatek_fbs" {
			fatekFBs = protocol
			break
		}
	}

	require.NotNil(t, fatekFBs, "FATEK FBs protocol not found")

	assert.Equal(t, "fatek_fbs", fatekFBs["type"])
	assert.Equal(t, "FATEK FBs", fatekFBs["name"])
	assert.Equal(t, "FATEK FBs 系列 PLC ASCII 協議", fatekFBs["description"])
}

/**
 * TestProtocolHandler_List_MC3E 測試 Mitsubishi MC 3E 協議
 */
func TestProtocolHandler_List_MC3E(t *testing.T) {
	r := setupProtocolRouter()

	req, _ := http.NewRequest("GET", "/datalink/protocols", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].([]interface{})

	// 尋找 MC 3E 協議
	var mc3e map[string]interface{}
	for _, item := range data {
		protocol := item.(map[string]interface{})
		if protocol["type"] == "mc_3e" {
			mc3e = protocol
			break
		}
	}

	require.NotNil(t, mc3e, "MC 3E protocol not found")

	assert.Equal(t, "mc_3e", mc3e["type"])
	assert.Equal(t, "Mitsubishi MC 3E", mc3e["name"])
	assert.Equal(t, "三菱 MC Protocol 3E Frame (Binary)", mc3e["description"])
}

/**
 * TestProtocolHandler_List_MQTT 測試 MQTT 協議
 */
func TestProtocolHandler_List_MQTT(t *testing.T) {
	r := setupProtocolRouter()

	req, _ := http.NewRequest("GET", "/datalink/protocols", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].([]interface{})

	// 尋找 MQTT 協議
	var mqtt map[string]interface{}
	for _, item := range data {
		protocol := item.(map[string]interface{})
		if protocol["type"] == "mqtt" {
			mqtt = protocol
			break
		}
	}

	require.NotNil(t, mqtt, "MQTT protocol not found")

	assert.Equal(t, "mqtt", mqtt["type"])
	assert.Equal(t, "MQTT", mqtt["name"])
	assert.Equal(t, "MQTT 訂閱接收，用於接收 IoT 設備資料", mqtt["description"])
}

/**
 * TestProtocolHandler_List_ConfigSchema_ModbusTCP 測試 Modbus TCP Config Schema
 */
func TestProtocolHandler_List_ConfigSchema_ModbusTCP(t *testing.T) {
	r := setupProtocolRouter()

	req, _ := http.NewRequest("GET", "/datalink/protocols", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].([]interface{})

	// 尋找 Modbus TCP 協議
	var modbusTCP map[string]interface{}
	for _, item := range data {
		protocol := item.(map[string]interface{})
		if protocol["type"] == "modbus_tcp" {
			modbusTCP = protocol
			break
		}
	}
	require.NotNil(t, modbusTCP)

	// 解析 Config Schema
	configSchema := modbusTCP["config_schema"].(string)
	var schema map[string]interface{}
	json.Unmarshal([]byte(configSchema), &schema)

	// 驗證結構
	assert.Equal(t, "object", schema["type"])
	assert.NotNil(t, schema["properties"])
	assert.NotNil(t, schema["required"])

	// 驗證 required 欄位
	required := schema["required"].([]interface{})
	assert.Contains(t, required, "host")
	assert.Contains(t, required, "slave_id")

	// 驗證 properties
	properties := schema["properties"].(map[string]interface{})
	assert.Contains(t, properties, "host")
	assert.Contains(t, properties, "port")
	assert.Contains(t, properties, "slave_id")
	assert.Contains(t, properties, "timeout")

	// 驗證 host 屬性
	host := properties["host"].(map[string]interface{})
	assert.Equal(t, "string", host["type"])
	assert.Equal(t, "主機位址", host["title"])

	// 驗證 port 屬性
	port := properties["port"].(map[string]interface{})
	assert.Equal(t, "integer", port["type"])
	assert.Equal(t, float64(502), port["default"])
	assert.Equal(t, "埠號", port["title"])

	// 驗證 slave_id 屬性
	slaveID := properties["slave_id"].(map[string]interface{})
	assert.Equal(t, "integer", slaveID["type"])
	assert.Equal(t, float64(1), slaveID["minimum"])
	assert.Equal(t, float64(247), slaveID["maximum"])
	assert.Equal(t, "從站 ID", slaveID["title"])
}

/**
 * TestProtocolHandler_List_ConfigSchema_FatekFBs 測試 FATEK FBs Config Schema
 */
func TestProtocolHandler_List_ConfigSchema_FatekFBs(t *testing.T) {
	r := setupProtocolRouter()

	req, _ := http.NewRequest("GET", "/datalink/protocols", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].([]interface{})

	// 尋找 FATEK FBs 協議
	var fatekFBs map[string]interface{}
	for _, item := range data {
		protocol := item.(map[string]interface{})
		if protocol["type"] == "fatek_fbs" {
			fatekFBs = protocol
			break
		}
	}
	require.NotNil(t, fatekFBs)

	// 解析 Config Schema
	configSchema := fatekFBs["config_schema"].(string)
	var schema map[string]interface{}
	json.Unmarshal([]byte(configSchema), &schema)

	// 驗證 properties
	properties := schema["properties"].(map[string]interface{})
	assert.Contains(t, properties, "mode")
	assert.Contains(t, properties, "data_bits")
	assert.Contains(t, properties, "stop_bits")
	assert.Contains(t, properties, "parity")
	assert.Contains(t, properties, "station_no")

	// 驗證 mode 屬性
	mode := properties["mode"].(map[string]interface{})
	assert.Equal(t, "string", mode["type"])
	modeEnum := mode["enum"].([]interface{})
	assert.Contains(t, modeEnum, "tcp")
	assert.Contains(t, modeEnum, "serial")
}

/**
 * TestProtocolHandler_List_ConfigSchema_MC3E 測試 MC 3E Config Schema
 */
func TestProtocolHandler_List_ConfigSchema_MC3E(t *testing.T) {
	r := setupProtocolRouter()

	req, _ := http.NewRequest("GET", "/datalink/protocols", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].([]interface{})

	var mc3e map[string]interface{}
	for _, item := range data {
		protocol := item.(map[string]interface{})
		if protocol["type"] == "mc_3e" {
			mc3e = protocol
			break
		}
	}
	require.NotNil(t, mc3e)

	configSchema := mc3e["config_schema"].(string)
	var schema map[string]interface{}
	json.Unmarshal([]byte(configSchema), &schema)

	properties := schema["properties"].(map[string]interface{})
	dataFormatValue, ok := properties["data_format"]
	require.True(t, ok, "MC 3E schema should expose data_format")

	dataFormat := dataFormatValue.(map[string]interface{})
	assert.Equal(t, "string", dataFormat["type"])
	dataFormatEnum := dataFormat["enum"].([]interface{})
	assert.Contains(t, dataFormatEnum, "ABCD")
	assert.Contains(t, dataFormatEnum, "BADC")
	assert.Contains(t, dataFormatEnum, "CDAB")
	assert.Contains(t, dataFormatEnum, "DCBA")
}

/**
 * TestProtocolHandler_List_AllProtocolsPresent 測試所有協議都存在
 */
func TestProtocolHandler_List_AllProtocolsPresent(t *testing.T) {
	r := setupProtocolRouter()

	req, _ := http.NewRequest("GET", "/datalink/protocols", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].([]interface{})

	expectedProtocols := []string{
		"modbus_tcp",
		"modbus_rtu",
		"modbus_udp",
		"fatek_fbs",
		"mc_3e",
		"mqtt",
	}

	// 建立實際協議類型列表
	actualProtocols := make(map[string]bool)
	for _, item := range data {
		protocol := item.(map[string]interface{})
		protocolType := protocol["type"].(string)
		actualProtocols[protocolType] = true
	}

	// 驗證所有預期協議都存在
	for _, expectedProtocol := range expectedProtocols {
		assert.True(t, actualProtocols[expectedProtocol], "Protocol %s not found", expectedProtocol)
	}
}

/**
 * TestProtocolHandler_List_JSONValid 測試回應是有效的 JSON
 */
func TestProtocolHandler_List_JSONValid(t *testing.T) {
	r := setupProtocolRouter()

	req, _ := http.NewRequest("GET", "/datalink/protocols", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response["data"])
}
