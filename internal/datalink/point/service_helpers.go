package point

import (
	"fmt"
	"strings"

	"go-gateway/internal/datalink/schema"
)

// isValidDataType 檢查資料型別是否有效
func isValidDataType(dataType schema.DataType) bool {
	validTypes := map[schema.DataType]bool{
		schema.DataTypeBool:    true,
		schema.DataTypeInt16:   true,
		schema.DataTypeUint16:  true,
		schema.DataTypeInt32:   true,
		schema.DataTypeUint32:  true,
		schema.DataTypeInt64:   true,
		schema.DataTypeUint64:  true,
		schema.DataTypeFloat32: true,
		schema.DataTypeFloat64: true,
		schema.DataTypeString:  true,
	}
	return validTypes[dataType]
}

// normalizePointDataFormat 正規化點位字節序字串（大寫、去空白）。
func normalizePointDataFormat(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}

// validatePointDataFormat 驗證點位 data_format；空字串表示沿用歷史 Modbus 預設（ABCD）。
func validatePointDataFormat(value string) error {
	n := normalizePointDataFormat(value)
	if n == "" {
		return nil
	}
	switch n {
	case "ABCD", "BADC", "CDAB", "DCBA":
		return nil
	default:
		return fmt.Errorf("不支援的 data_format: %s", value)
	}
}
