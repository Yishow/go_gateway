package schema

import "testing"

func TestRegisterCountForDataType(t *testing.T) {
	tests := map[DataType]int{
		DataTypeBool:    1,
		DataTypeInt16:   1,
		DataTypeUint16:  1,
		DataTypeInt32:   2,
		DataTypeUint32:  2,
		DataTypeFloat32: 2,
		DataTypeInt64:   4,
		DataTypeUint64:  4,
		DataTypeFloat64: 4,
		DataTypeString:  1,
	}

	for dt, expected := range tests {
		if got := RegisterCountForDataType(dt); got != expected {
			t.Fatalf("data type %s expected %d got %d", dt, expected, got)
		}
	}
}

func TestRegisterCountForUnknownDataType(t *testing.T) {
	if got := RegisterCountForDataType(DataType("custom")); got != 1 {
		t.Fatalf("unknown data type expected 1 got %d", got)
	}
}
