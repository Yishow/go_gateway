package adapters

import (
	"fmt"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
)

func (c *ModbusTCPConnector) executeRead(req connector.ReadRequest) (connector.ReadResult, error) {
	result := connector.ReadResult{
		Quality: schema.QualityGood,
	}

	// 解析地址 (格式: "40001" 或 "HR0" 或 "0")
	address, function, err := parseModbusAddress(req.Address, req.Function)
	if err != nil {
		result.Quality = schema.QualityBad
		result.Error = err.Error()
		return result, err
	}

	// 計算需要讀取的暫存器數量
	count := schema.RegisterCountForDataType(req.DataType)
	if req.Count > 0 {
		count = req.Count
	}
	if count < 1 || count > 2000 {
		return result, fmt.Errorf("無效的讀取數量: %d", count)
	}
	uCount := uint16(count) // #nosec G115

	// 根據功能碼讀取
	switch function {
	case "01", "coil", "FC01":
		// 讀取線圈
		values, err := c.client.ReadCoils(address, uCount)
		if err != nil {
			return result, err
		}
		if len(values) > 0 {
			result.Value = values[0]
		}

	case "02", "discrete", "FC02":
		// 讀取離散輸入
		values, err := c.client.ReadDiscreteInputs(address, uCount)
		if err != nil {
			return result, err
		}
		if len(values) > 0 {
			result.Value = values[0]
		}

	case "03", "holding", "FC03", "":
		// 讀取保持暫存器 (預設)
		values, err := c.client.ReadHoldingRegisters(address, uCount)
		if err != nil {
			return result, err
		}
		result.RawBytes = uint16SliceToBytes(values)
		result.Value = convertModbusValue(values, req.DataType, effectiveModbusDataFormat(req.DataFormat, c.config.DataFormat))

	case "04", "input", "FC04":
		// 讀取輸入暫存器
		values, err := c.client.ReadInputRegisters(address, uCount)
		if err != nil {
			return result, err
		}
		result.RawBytes = uint16SliceToBytes(values)
		result.Value = convertModbusValue(values, req.DataType, effectiveModbusDataFormat(req.DataFormat, c.config.DataFormat))

	default:
		return result, fmt.Errorf("不支援的功能碼: %s", function)
	}

	return result, nil
}
