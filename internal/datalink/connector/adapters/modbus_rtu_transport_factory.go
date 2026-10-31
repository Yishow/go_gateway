package adapters

import (
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/modbus"
	"time"
)

// NewModbusRTUConnectorWithTransportFactory supplies only the raw byte transport.
// Connect still owns configuration defaults, the production client and lifecycle.
// The default registry constructor continues to use the physical RTU transport.
func NewModbusRTUConnectorWithTransportFactory(factory func(schema.ConnectionConfigModbusRTU) modbus.Transport) connector.Protocol {
	if factory == nil {
		factory = func(schema.ConnectionConfigModbusRTU) modbus.Transport { return nil }
	}
	return &ModbusRTUConnector{transportFactory: factory}
}

func (c *ModbusRTUConnector) newTransport() modbus.Transport {
	if c.transportFactory != nil {
		return c.transportFactory(c.config)
	}
	return modbus.NewRTUTransport(c.config.SerialPort, c.config.BaudRate, c.config.DataBits, c.config.StopBits, c.config.Parity, time.Duration(c.config.Timeout)*time.Second)
}
