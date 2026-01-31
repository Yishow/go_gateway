package optimizer

import (
	"testing"

	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// 任務 2.1: [RED] Block 資料解析回填測試
// =============================================================================

func TestDispatcher_ParseBlockData(t *testing.T) {
	// 建立一個 Block(D100, Len=3) 的場景
	block := Block{
		MemoryArea:   "D",
		StartAddress: 100,
		Length:       3,
		Points: []PointInfo{
			{ID: "p1", Address: "D100", DataType: schema.DataTypeInt16},
			{ID: "p2", Address: "D101", DataType: schema.DataTypeInt16},
			{ID: "p3", Address: "D102", DataType: schema.DataTypeInt16},
		},
		PointOffsets: map[string]int{
			"p1": 0,
			"p2": 1,
			"p3": 2,
		},
	}

	// Raw Bytes: [00 01 00 02 00 03] (Big Endian)
	rawBytes := []byte{0x00, 0x01, 0x00, 0x02, 0x00, 0x03}

	dispatcher := NewDispatcher(DispatcherConfig{ByteOrder: BigEndian})
	values, err := dispatcher.Dispatch(block, rawBytes)

	if err != nil {
		t.Fatalf("Dispatch 失敗: %v", err)
	}

	if len(values) != 3 {
		t.Errorf("預期 3 個值，實際得到 %d 個", len(values))
		return
	}

	// 驗證值
	expected := map[string]interface{}{
		"p1": int16(1),
		"p2": int16(2),
		"p3": int16(3),
	}

	for id, exp := range expected {
		val, exists := values[id]
		if !exists {
			t.Errorf("缺少點位 %s 的值", id)
			continue
		}
		if val != exp {
			t.Errorf("點位 %s: 預期 %v，實際 %v", id, exp, val)
		}
	}
}

func TestDispatcher_LittleEndian(t *testing.T) {
	block := Block{
		MemoryArea:   "D",
		StartAddress: 100,
		Length:       1,
		Points: []PointInfo{
			{ID: "p1", Address: "D100", DataType: schema.DataTypeInt16},
		},
		PointOffsets: map[string]int{"p1": 0},
	}

	// Little Endian: [01 00] = 1
	rawBytes := []byte{0x01, 0x00}

	dispatcher := NewDispatcher(DispatcherConfig{ByteOrder: LittleEndian})
	values, err := dispatcher.Dispatch(block, rawBytes)

	if err != nil {
		t.Fatalf("Dispatch 失敗: %v", err)
	}

	if val, ok := values["p1"].(int16); !ok || val != 1 {
		t.Errorf("預期 1，實際 %v", values["p1"])
	}
}

func TestDispatcher_Int32(t *testing.T) {
	block := Block{
		MemoryArea:   "D",
		StartAddress: 100,
		Length:       2,
		Points: []PointInfo{
			{ID: "p1", Address: "D100", DataType: schema.DataTypeInt32},
		},
		PointOffsets: map[string]int{"p1": 0},
	}

	// Big Endian: [00 00 01 00] = 256
	rawBytes := []byte{0x00, 0x00, 0x01, 0x00}

	dispatcher := NewDispatcher(DispatcherConfig{ByteOrder: BigEndian})
	values, err := dispatcher.Dispatch(block, rawBytes)

	if err != nil {
		t.Fatalf("Dispatch 失敗: %v", err)
	}

	if val, ok := values["p1"].(int32); !ok || val != 256 {
		t.Errorf("預期 256，實際 %v", values["p1"])
	}
}

func TestDispatcher_Float32(t *testing.T) {
	block := Block{
		MemoryArea:   "D",
		StartAddress: 100,
		Length:       2,
		Points: []PointInfo{
			{ID: "p1", Address: "D100", DataType: schema.DataTypeFloat32},
		},
		PointOffsets: map[string]int{"p1": 0},
	}

	// Big Endian float32: 3.14159 ≈ 0x40490FDB
	rawBytes := []byte{0x40, 0x49, 0x0F, 0xDB}

	dispatcher := NewDispatcher(DispatcherConfig{ByteOrder: BigEndian})
	values, err := dispatcher.Dispatch(block, rawBytes)

	if err != nil {
		t.Fatalf("Dispatch 失敗: %v", err)
	}

	val, ok := values["p1"].(float32)
	if !ok {
		t.Errorf("預期 float32，實際類型 %T", values["p1"])
		return
	}

	// 允許小誤差
	if val < 3.14 || val > 3.15 {
		t.Errorf("預期約 3.14159，實際 %v", val)
	}
}

func TestDispatcher_MixedTypes(t *testing.T) {
	// 複合場景：D100(int16) + D101-D102(int32)
	block := Block{
		MemoryArea:   "D",
		StartAddress: 100,
		Length:       3,
		Points: []PointInfo{
			{ID: "p1", Address: "D100", DataType: schema.DataTypeInt16},
			{ID: "p2", Address: "D101", DataType: schema.DataTypeInt32},
		},
		PointOffsets: map[string]int{
			"p1": 0,
			"p2": 1,
		},
	}

	// [00 0A] [00 00 00 64] = 10, 100
	rawBytes := []byte{0x00, 0x0A, 0x00, 0x00, 0x00, 0x64}

	dispatcher := NewDispatcher(DispatcherConfig{ByteOrder: BigEndian})
	values, err := dispatcher.Dispatch(block, rawBytes)

	if err != nil {
		t.Fatalf("Dispatch 失敗: %v", err)
	}

	if val, ok := values["p1"].(int16); !ok || val != 10 {
		t.Errorf("p1: 預期 10，實際 %v", values["p1"])
	}

	if val, ok := values["p2"].(int32); !ok || val != 100 {
		t.Errorf("p2: 預期 100，實際 %v", values["p2"])
	}
}

func TestDispatcher_InsufficientData(t *testing.T) {
	block := Block{
		MemoryArea:   "D",
		StartAddress: 100,
		Length:       3,
		Points: []PointInfo{
			{ID: "p1", Address: "D100", DataType: schema.DataTypeInt16},
		},
		PointOffsets: map[string]int{"p1": 0},
	}

	// 資料不足
	rawBytes := []byte{0x00}

	dispatcher := NewDispatcher(DispatcherConfig{ByteOrder: BigEndian})
	_, err := dispatcher.Dispatch(block, rawBytes)

	if err == nil {
		t.Error("預期資料不足錯誤，但未發生")
	}
}
