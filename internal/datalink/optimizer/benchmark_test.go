package optimizer

import (
	"testing"

	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// 任務 4.1: 效能基準測試
// =============================================================================

func BenchmarkMerger_100Points(b *testing.B) {
	points := generatePoints(100)
	merger := NewBlockMerger(DefaultMergerConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		merger.Merge(points)
	}
}

func BenchmarkMerger_1000Points(b *testing.B) {
	points := generatePoints(1000)
	merger := NewBlockMerger(DefaultMergerConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		merger.Merge(points)
	}
}

func BenchmarkMerger_10000Points(b *testing.B) {
	points := generatePoints(10000)
	merger := NewBlockMerger(DefaultMergerConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		merger.Merge(points)
	}
}

func BenchmarkDispatcher_SmallBlock(b *testing.B) {
	block := Block{
		MemoryArea:   "D",
		StartAddress: 100,
		Length:       10,
		Points:       make([]PointInfo, 10),
		PointOffsets: make(map[string]int),
	}

	for i := 0; i < 10; i++ {
		id := "p" + itoa(i)
		block.Points[i] = PointInfo{ID: id, Address: "D" + itoa(100+i), DataType: schema.DataTypeInt16}
		block.PointOffsets[id] = i
	}

	rawBytes := make([]byte, 20) // 10 registers * 2 bytes
	dispatcher := NewDispatcher(DefaultDispatcherConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dispatcher.Dispatch(block, rawBytes)
	}
}

func BenchmarkDispatcher_LargeBlock(b *testing.B) {
	block := Block{
		MemoryArea:   "D",
		StartAddress: 100,
		Length:       125,
		Points:       make([]PointInfo, 125),
		PointOffsets: make(map[string]int),
	}

	for i := 0; i < 125; i++ {
		id := "p" + itoa(i)
		block.Points[i] = PointInfo{ID: id, Address: "D" + itoa(100+i), DataType: schema.DataTypeInt16}
		block.PointOffsets[id] = i
	}

	rawBytes := make([]byte, 250) // 125 registers * 2 bytes
	dispatcher := NewDispatcher(DefaultDispatcherConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dispatcher.Dispatch(block, rawBytes)
	}
}

// generatePoints 產生測試用點位
func generatePoints(count int) []PointInfo {
	points := make([]PointInfo, count)
	for i := 0; i < count; i++ {
		// 模擬一些分散的地址
		area := "D"
		if i%5 == 0 {
			area = "M"
		}
		addr := (i / 10) * 15 + (i % 10) // 每 10 個一組，中間有些空洞

		points[i] = PointInfo{
			ID:       "p" + itoa(i),
			Address:  area + itoa(addr),
			DataType: schema.DataTypeInt16,
		}
	}
	return points
}
