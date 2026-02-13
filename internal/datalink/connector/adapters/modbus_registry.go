package adapters

import (
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// 註冊連接器
// =============================================================================

func init() {
	connector.Register(schema.ProtocolModbusTCP, NewModbusTCPConnector)
	connector.Register(schema.ProtocolModbusRTU, NewModbusRTUConnector)
	connector.Register(schema.ProtocolModbusUDP, NewModbusUDPConnector)
}
