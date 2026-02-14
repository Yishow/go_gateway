package hsllogic

// =============================================================================
// DataType 數據類型定義
// =============================================================================

// DataType 數據類型常數
type DataType string

const (
	// DataTypeBool 布林值 (1 位元)
	DataTypeBool DataType = "bool"
	// DataTypeInt8 有符號 8 位元整數
	DataTypeInt8 DataType = "int8"
	// DataTypeUint8 無符號 8 位元整數
	DataTypeUint8 DataType = "uint8"
	// DataTypeInt16 有符號 16 位元整數
	DataTypeInt16 DataType = "int16"
	// DataTypeUint16 無符號 16 位元整數
	DataTypeUint16 DataType = "uint16"
	// DataTypeInt32 有符號 32 位元整數
	DataTypeInt32 DataType = "int32"
	// DataTypeUint32 無符號 32 位元整數
	DataTypeUint32 DataType = "uint32"
	// DataTypeInt64 有符號 64 位元整數
	DataTypeInt64 DataType = "int64"
	// DataTypeUint64 無符號 64 位元整數
	DataTypeUint64 DataType = "uint64"
	// DataTypeFloat32 32 位元浮點數
	DataTypeFloat32 DataType = "float32"
	// DataTypeFloat64 64 位元浮點數
	DataTypeFloat64 DataType = "float64"
	// DataTypeString ASCII 字串
	DataTypeString DataType = "string"
	// DataTypeBytes 原始位元組
	DataTypeBytes DataType = "bytes"
)

// =============================================================================
// DataTypeInfo 數據類型資訊
// =============================================================================

// DataTypeInfo 數據類型的元資料
type DataTypeInfo struct {
	// Type 數據類型
	Type DataType
	// ByteSize 位元組大小
	ByteSize int
	// RegisterCount 暫存器數量 (16 位元暫存器)
	RegisterCount int
	// IsSigned 是否為有符號類型
	IsSigned bool
	// IsFloating 是否為浮點類型
	IsFloating bool
}

// dataTypeInfoMap 數據類型資訊映射
var dataTypeInfoMap = map[DataType]DataTypeInfo{
	DataTypeBool:    {DataTypeBool, 1, 1, false, false},
	DataTypeInt8:    {DataTypeInt8, 1, 1, true, false},
	DataTypeUint8:   {DataTypeUint8, 1, 1, false, false},
	DataTypeInt16:   {DataTypeInt16, 2, 1, true, false},
	DataTypeUint16:  {DataTypeUint16, 2, 1, false, false},
	DataTypeInt32:   {DataTypeInt32, 4, 2, true, false},
	DataTypeUint32:  {DataTypeUint32, 4, 2, false, false},
	DataTypeInt64:   {DataTypeInt64, 8, 4, true, false},
	DataTypeUint64:  {DataTypeUint64, 8, 4, false, false},
	DataTypeFloat32: {DataTypeFloat32, 4, 2, false, true},
	DataTypeFloat64: {DataTypeFloat64, 8, 4, false, true},
	DataTypeString:  {DataTypeString, 0, 0, false, false}, // 可變長度
	DataTypeBytes:   {DataTypeBytes, 0, 0, false, false},  // 可變長度
}

// GetDataTypeInfo 取得數據類型資訊
func GetDataTypeInfo(dt DataType) DataTypeInfo {
	if info, ok := dataTypeInfoMap[dt]; ok {
		return info
	}
	return DataTypeInfo{Type: dt}
}

// ByteSizeForDataType 取得數據類型的位元組大小
func ByteSizeForDataType(dt DataType) int {
	return GetDataTypeInfo(dt).ByteSize
}

// RegisterCountForDataType 取得數據類型需要的暫存器數量
func RegisterCountForDataType(dt DataType) int {
	info := GetDataTypeInfo(dt)
	if info.RegisterCount == 0 {
		return 1 // 預設至少 1 個暫存器
	}
	return info.RegisterCount
}
