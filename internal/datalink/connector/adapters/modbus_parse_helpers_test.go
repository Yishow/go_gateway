package adapters

import (
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/assert"
)

func TestConvertModbusValue_Float32_DataFormat(t *testing.T) {
	// 暫存器順序為 Modbus 傳回順序：高位字在前、低位字在後（與舊版 ABCD 解碼一致時為 1.0）
	regs := []uint16{0x3f80, 0x0000}

	gotABCD := convertModbusValue(regs, schema.DataTypeFloat32, "ABCD")
	assert.Equal(t, float32(1), gotABCD.(float32))

	gotDefault := convertModbusValue(regs, schema.DataTypeFloat32, "")
	assert.Equal(t, float32(1), gotDefault.(float32))

	gotCDAB := convertModbusValue(regs, schema.DataTypeFloat32, "CDAB")
	assert.NotEqual(t, float32(1), gotCDAB.(float32))
}
