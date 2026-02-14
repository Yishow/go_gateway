package device

import (
	"encoding/json"
	"fmt"

	"go-gateway/internal/datalink/schema"
)

// generateUUID 產生 UUID
// =============================================================================
// 連線配置輔助
// =============================================================================

// ParseConnectionConfig 解析連線配置為具體類型
func ParseConnectionConfig(device *schema.Device) (interface{}, error) {
	var config interface{}

	switch device.Protocol {
	case schema.ProtocolModbusTCP:
		var c schema.ConnectionConfigModbusTCP
		if err := json.Unmarshal([]byte(device.ConnectionConfig), &c); err != nil {
			return nil, err
		}
		config = c
	case schema.ProtocolModbusRTU:
		var c schema.ConnectionConfigModbusRTU
		if err := json.Unmarshal([]byte(device.ConnectionConfig), &c); err != nil {
			return nil, err
		}
		config = c
	case schema.ProtocolModbusUDP:
		var c schema.ConnectionConfigModbusUDP
		if err := json.Unmarshal([]byte(device.ConnectionConfig), &c); err != nil {
			return nil, err
		}
		config = c
	case schema.ProtocolFatekFBs:
		var c schema.ConnectionConfigFatekFBs
		if err := json.Unmarshal([]byte(device.ConnectionConfig), &c); err != nil {
			return nil, err
		}
		config = c
	case schema.ProtocolMC3E:
		var c schema.ConnectionConfigMC3E
		if err := json.Unmarshal([]byte(device.ConnectionConfig), &c); err != nil {
			return nil, err
		}
		config = c
	case schema.ProtocolMQTT:
		var c schema.ConnectionConfigMQTT
		if err := json.Unmarshal([]byte(device.ConnectionConfig), &c); err != nil {
			return nil, err
		}
		config = c
	default:
		return nil, fmt.Errorf("未知的協議類型: %s", device.Protocol)
	}

	return config, nil
}
