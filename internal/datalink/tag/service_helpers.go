package tag

import "go-gateway/internal/datalink/schema"

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
