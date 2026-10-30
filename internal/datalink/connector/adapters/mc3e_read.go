package adapters

import (
	"context"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/lib/hsllogic"
)

func (c *MC3EConnector) executeRead(req connector.ReadRequest) (connector.ReadResult, error) {
	return c.executeReadContext(context.Background(), req)
}

func (c *MC3EConnector) executeReadContext(ctx context.Context, req connector.ReadRequest) (connector.ReadResult, error) {
	if err := validateMCReadProfile(req); err != nil {
		return connector.ReadResult{Quality: schema.QualityBad, Error: err.Error()}, err
	}
	client := c.client.WithContext(ctx)
	result := connector.ReadResult{
		Timestamp: time.Now(),
		Quality:   schema.QualityGood,
	}

	// 使用 hsllogic 解析地址 (格式: "D100", "M0", "X0", "Y0", "R100", "W100", etc.)
	parsedAddr, err := hsllogic.ParseAddress(hsllogic.ProtocolMitsubishi, req.Address)
	if err != nil {
		return result, err
	}

	count := 1
	if req.Count > 0 {
		count = req.Count
	}

	// 根據設備類型選擇讀取方式
	if parsedAddr.IsBitDevice {
		// 位元設備
		values, err := client.BatchReadBit(parsedAddr.DeviceType, parsedAddr.Offset, count)
		if err != nil {
			return result, err
		}
		if len(values) > 0 {
			if count == 1 {
				result.Value = values[0]
			} else {
				result.Value = values
			}
		}
	} else {
		// 字組設備
		regPerValue := hsllogic.RegisterCountForDataType(hsllogic.DataType(req.DataType))
		modern := modernMCScalar(req)
		if modern && (req.Count < 0 || regPerValue < 1 || count > 125/regPerValue) {
			return result, errMCReadCount
		}
		wordCount := count * regPerValue

		var values []int
		if modern {
			values, err = client.BatchReadWordExact(parsedAddr.DeviceType, parsedAddr.Offset, wordCount)
		} else {
			values, err = client.BatchReadWord(parsedAddr.DeviceType, parsedAddr.Offset, wordCount)
		}
		if err != nil {
			return result, err
		}
		result.RawBytes = intSliceToBytes(values)

		registers := make([]uint16, len(values))
		for i, v := range values {
			registers[i] = uint16(v) // #nosec G115
		}

		if modern {
			result.Value, err = decodeMCReadWords(registers, req, count)
			if err != nil {
				return result, err
			}
			return result, nil
		}
		if count == 1 {
			result.Value = c.dataConverter.RegistersToValue(registers, hsllogic.DataType(req.DataType))
		} else {
			result.Value = c.dataConverter.RegistersToValues(registers, hsllogic.DataType(req.DataType), count)
		}
	}

	return result, nil
}
