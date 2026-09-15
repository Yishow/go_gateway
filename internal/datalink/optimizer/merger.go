// Package optimizer 提供封包聚合優化功能
//
// 本套件實作智能封包聚合機制，在 Runtime 動態分析並合併相鄰的暫存器讀取請求，
// 以最小化網路封包交互次數並最大化資料吞吐量。
package optimizer

import (
	"regexp"
	"sort"
	"strconv"
	"strings"

	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// 類型定義
// =============================================================================

// PointInfo 點位資訊 (用於優化計算)
type PointInfo struct {
	ID       string
	Address  string
	DataType schema.DataType
	DeviceID string
}

// Block 聚合後的讀取區塊
type Block struct {
	// MemoryArea 記憶體區域 (D, M, X, Y 等)
	MemoryArea string

	// StartAddress 起始地址
	StartAddress int

	// Length 讀取長度 (暫存器數量)
	Length int

	// Points 包含的原始點位
	Points []PointInfo

	// PointOffsets 每個點位在 Block 中的偏移量 (以暫存器為單位)
	PointOffsets map[string]int
}

// MergerConfig 合併器配置
type MergerConfig struct {
	// MaxHoleSize 最大空洞大小 (允許跳過的暫存器數量)
	MaxHoleSize int

	// MaxBlockLength 單一 Block 的最大長度 (受 PDU 限制)
	MaxBlockLength int
}

// DefaultMergerConfig 預設合併器配置
func DefaultMergerConfig() MergerConfig {
	return MergerConfig{
		MaxHoleSize:    10,  // 允許 10 個暫存器的空洞
		MaxBlockLength: 125, // Modbus 最大 125 暫存器
	}
}

// =============================================================================
// BlockMerger 區塊合併器
// =============================================================================

// BlockMerger 區塊合併器
type BlockMerger struct {
	config MergerConfig
}

// NewBlockMerger 建立新的區塊合併器
func NewBlockMerger(config MergerConfig) *BlockMerger {
	return &BlockMerger{config: config}
}

// parsedPoint 解析後的點位資訊
type parsedPoint struct {
	info       PointInfo
	memoryArea string
	address    int
	regCount   int // 佔用的暫存器數量
}

// Merge 合併點位為最佳化的讀取區塊
func (m *BlockMerger) Merge(points []PointInfo) []Block {
	if len(points) == 0 {
		return nil
	}

	// 解析所有點位
	parsed := make([]parsedPoint, 0, len(points))
	for _, p := range points {
		area, addr := parseAddress(p.Address)
		if area == "" {
			continue // 無法解析的地址跳過
		}
		parsed = append(parsed, parsedPoint{
			info:       p,
			memoryArea: area,
			address:    addr,
			regCount:   registerCount(p.DataType),
		})
	}

	// 按記憶體區域分組
	groups := make(map[string][]parsedPoint)
	for _, p := range parsed {
		groups[p.memoryArea] = append(groups[p.memoryArea], p)
	}

	// 對每個區域進行合併
	var blocks []Block
	for area, pts := range groups {
		areaBlocks := m.mergeGroup(area, pts)
		blocks = append(blocks, areaBlocks...)
	}

	return blocks
}

// mergeGroup 合併同一記憶體區域的點位
func (m *BlockMerger) mergeGroup(area string, points []parsedPoint) []Block {
	if len(points) == 0 {
		return nil
	}

	// 按地址排序
	sort.Slice(points, func(i, j int) bool {
		return points[i].address < points[j].address
	})

	var blocks []Block
	var currentBlock *Block

	for _, p := range points {
		if currentBlock == nil {
			// 開始新的 Block
			currentBlock = &Block{
				MemoryArea:   area,
				StartAddress: p.address,
				Length:       p.regCount,
				Points:       []PointInfo{p.info},
				PointOffsets: map[string]int{p.info.ID: 0},
			}
			continue
		}

		// 計算與當前 Block 的距離
		currentEnd := currentBlock.StartAddress + currentBlock.Length
		hole := p.address - currentEnd

		// 計算合併後的長度
		newEnd := p.address + p.regCount
		newLength := newEnd - currentBlock.StartAddress

		// 判斷是否可以合併
		canMerge := hole <= m.config.MaxHoleSize && newLength <= m.config.MaxBlockLength

		if canMerge {
			// 合併到當前 Block
			offset := p.address - currentBlock.StartAddress
			currentBlock.Length = newLength
			currentBlock.Points = append(currentBlock.Points, p.info)
			currentBlock.PointOffsets[p.info.ID] = offset
		} else {
			// 完成當前 Block，開始新的
			blocks = append(blocks, *currentBlock)
			currentBlock = &Block{
				MemoryArea:   area,
				StartAddress: p.address,
				Length:       p.regCount,
				Points:       []PointInfo{p.info},
				PointOffsets: map[string]int{p.info.ID: 0},
			}
		}
	}

	// 加入最後一個 Block
	if currentBlock != nil {
		blocks = append(blocks, *currentBlock)
	}

	return blocks
}

// =============================================================================
// 輔助函數
// =============================================================================

// 地址解析正規表示式
var addressRegex = regexp.MustCompile(`^([A-Za-z]+)(\d+)$`)

// parseAddress 解析地址字串，返回記憶體區域和數字地址
func parseAddress(addr string) (prefix string, number int) {
	addr = strings.TrimSpace(addr)
	matches := addressRegex.FindStringSubmatch(addr)
	if len(matches) != 3 {
		return "", 0
	}

	area := strings.ToUpper(matches[1])
	num, err := strconv.Atoi(matches[2])
	if err != nil {
		return "", 0
	}

	return area, num
}

// registerCount 根據資料型別計算佔用的暫存器數量
func registerCount(dt schema.DataType) int {
	switch dt {
	case schema.DataTypeBool:
		return 1
	case schema.DataTypeInt16, schema.DataTypeUint16:
		return 1
	case schema.DataTypeInt32, schema.DataTypeUint32, schema.DataTypeFloat32:
		return 2
	case schema.DataTypeInt64, schema.DataTypeUint64, schema.DataTypeFloat64:
		return 4
	case schema.DataTypeString:
		return 8 // 預設 16 字節字串
	default:
		return 1
	}
}
