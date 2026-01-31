package optimizer

import (
	"testing"

	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// 任務 1.1: [RED] BlockMerger 基礎行為測試
// =============================================================================

func TestMerger_Consecutive(t *testing.T) {
	// Case A: 連續地址應合併為單一 Block
	points := []PointInfo{
		{ID: "p1", Address: "D100", DataType: schema.DataTypeInt16},
		{ID: "p2", Address: "D101", DataType: schema.DataTypeInt16},
		{ID: "p3", Address: "D102", DataType: schema.DataTypeInt16},
	}

	merger := NewBlockMerger(DefaultMergerConfig())
	blocks := merger.Merge(points)

	if len(blocks) != 1 {
		t.Errorf("預期 1 個 Block，實際得到 %d 個", len(blocks))
		return
	}

	block := blocks[0]
	if block.StartAddress != 100 {
		t.Errorf("預期起始地址 100，實際為 %d", block.StartAddress)
	}
	if block.Length != 3 {
		t.Errorf("預期長度 3，實際為 %d", block.Length)
	}
	if len(block.Points) != 3 {
		t.Errorf("預期 3 個 Points，實際為 %d", len(block.Points))
	}
}

func TestMerger_Separated(t *testing.T) {
	// Case B: 分離地址應產生多個 Block
	points := []PointInfo{
		{ID: "p1", Address: "D100", DataType: schema.DataTypeInt16},
		{ID: "p2", Address: "D200", DataType: schema.DataTypeInt16},
	}

	config := DefaultMergerConfig()
	config.MaxHoleSize = 5 // 小空洞限制
	merger := NewBlockMerger(config)
	blocks := merger.Merge(points)

	if len(blocks) != 2 {
		t.Errorf("預期 2 個 Blocks，實際得到 %d 個", len(blocks))
	}
}

// =============================================================================
// 任務 1.3: [RED] PDU 限制與空洞測試
// =============================================================================

func TestMerger_HoleAllowed(t *testing.T) {
	// Case C: 空洞在限制內應合併
	points := []PointInfo{
		{ID: "p1", Address: "D100", DataType: schema.DataTypeInt16},
		{ID: "p2", Address: "D105", DataType: schema.DataTypeInt16},
	}

	config := DefaultMergerConfig()
	config.MaxHoleSize = 10 // 允許 10 個暫存器的空洞
	merger := NewBlockMerger(config)
	blocks := merger.Merge(points)

	if len(blocks) != 1 {
		t.Errorf("預期 1 個 Block (空洞 5 < 最大 10)，實際得到 %d 個", len(blocks))
		return
	}

	// Block 應包含 D100-D105，長度 6
	if blocks[0].Length != 6 {
		t.Errorf("預期長度 6，實際為 %d", blocks[0].Length)
	}
}

func TestMerger_HoleExceeded(t *testing.T) {
	// 空洞超過限制應分離
	points := []PointInfo{
		{ID: "p1", Address: "D100", DataType: schema.DataTypeInt16},
		{ID: "p2", Address: "D120", DataType: schema.DataTypeInt16},
	}

	config := DefaultMergerConfig()
	config.MaxHoleSize = 5 // 最大空洞 5
	merger := NewBlockMerger(config)
	blocks := merger.Merge(points)

	if len(blocks) != 2 {
		t.Errorf("預期 2 個 Blocks (空洞 20 > 最大 5)，實際得到 %d 個", len(blocks))
	}
}

func TestMerger_PDULimit(t *testing.T) {
	// Case D: PDU 限制應拆分為多個 Blocks
	points := make([]PointInfo, 200)
	for i := 0; i < 200; i++ {
		points[i] = PointInfo{
			ID:       "p" + string(rune(i)),
			Address:  "D" + itoa(100+i),
			DataType: schema.DataTypeInt16,
		}
	}

	config := DefaultMergerConfig()
	config.MaxBlockLength = 100 // 最大 100 個暫存器
	merger := NewBlockMerger(config)
	blocks := merger.Merge(points)

	if len(blocks) < 2 {
		t.Errorf("預期至少 2 個 Blocks (200 > MaxLen 100)，實際得到 %d 個", len(blocks))
	}

	// 每個 Block 長度不應超過 MaxBlockLength
	for i, block := range blocks {
		if block.Length > config.MaxBlockLength {
			t.Errorf("Block %d 長度 %d 超過最大限制 %d", i, block.Length, config.MaxBlockLength)
		}
	}
}

func TestMerger_MixedMemoryAreas(t *testing.T) {
	// 不同記憶體區域不應合併
	points := []PointInfo{
		{ID: "p1", Address: "D100", DataType: schema.DataTypeInt16},
		{ID: "p2", Address: "D101", DataType: schema.DataTypeInt16},
		{ID: "p3", Address: "M100", DataType: schema.DataTypeBool},
		{ID: "p4", Address: "M101", DataType: schema.DataTypeBool},
	}

	merger := NewBlockMerger(DefaultMergerConfig())
	blocks := merger.Merge(points)

	if len(blocks) != 2 {
		t.Errorf("預期 2 個 Blocks (不同記憶體區域)，實際得到 %d 個", len(blocks))
	}
}

func TestMerger_EmptyInput(t *testing.T) {
	merger := NewBlockMerger(DefaultMergerConfig())
	blocks := merger.Merge(nil)

	if len(blocks) != 0 {
		t.Errorf("空輸入應返回空 Blocks，實際得到 %d 個", len(blocks))
	}
}

func TestMerger_SinglePoint(t *testing.T) {
	points := []PointInfo{
		{ID: "p1", Address: "D100", DataType: schema.DataTypeInt16},
	}

	merger := NewBlockMerger(DefaultMergerConfig())
	blocks := merger.Merge(points)

	if len(blocks) != 1 {
		t.Errorf("預期 1 個 Block，實際得到 %d 個", len(blocks))
		return
	}

	if blocks[0].Length != 1 {
		t.Errorf("預期長度 1，實際為 %d", blocks[0].Length)
	}
}

// 輔助函數
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	return result
}
