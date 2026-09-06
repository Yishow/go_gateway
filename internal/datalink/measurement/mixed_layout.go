package measurement

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"

	"go-gateway/internal/datalink/schema"
	"go-gateway/lib/hsllogic"
)

// LayoutMode 來源規則讀取佈局模式。
type LayoutMode string

const (
	// LayoutModeHomogeneous 均一讀取模式（預設，所有點位相同型別與長度）。
	LayoutModeHomogeneous LayoutMode = "homogeneous"
	// LayoutModeMixed 混合讀取模式（同一讀取範圍內不同型別、長度與語意）。
	LayoutModeMixed LayoutMode = "mixed"
)

// MixedItem 描述混合讀取佈局中的單一項目。
type MixedItem struct {
	ItemID          string            `json:"item_id"`
	Name            string            `json:"name"`
	Address         string            `json:"address"`
	RegisterOffset  int               `json:"register_offset"`
	BitOffset       *int              `json:"bit_offset,omitempty"`
	BitLength       *int              `json:"bit_length,omitempty"`
	DataType        schema.DataType   `json:"data_type"`
	DataFormat      string            `json:"data_format,omitempty"`
	TargetDataType  *schema.DataType  `json:"target_data_type,omitempty"`
	ScaleMultiplier *float64          `json:"scale_multiplier,omitempty"`
	ScaleOffset     *float64          `json:"scale_offset,omitempty"`
	Quantity        string            `json:"quantity,omitempty"`
	Unit            string            `json:"unit,omitempty"`
	SemanticKind    SemanticKind      `json:"semantic_kind"`
	CounterPolicy   *CounterPolicy    `json:"counter_policy,omitempty"`
	StateMap        map[string]string `json:"state_map,omitempty"`
	BitmaskLabels   map[int]string    `json:"bitmask_labels,omitempty"`
}

// DecodedItemResult 混合讀取中單一項目的解碼結果。
type DecodedItemResult struct {
	ItemID        string             `json:"item_id"`
	Value         interface{}        `json:"value"`
	RawValue      interface{}        `json:"raw_value,omitempty"`
	Quality       schema.QualityFlag `json:"quality"`
	QualityReason string             `json:"quality_reason,omitempty"`
	Transformed   bool               `json:"transformed"`
}

// WordCountForDataType 回傳給定資料型態占用的 16 位元暫存器數量。
func WordCountForDataType(dt schema.DataType) int {
	return schema.RegisterCountForDataType(dt)
}

// ValidateMixedLayout 驗證混合項目的暫存器占用、資料型態、字序與換算參數。
func ValidateMixedLayout(items []MixedItem) error {
	if len(items) == 0 {
		return fmt.Errorf("mixed layout must contain at least one item")
	}

	for i := 0; i < len(items); i++ {
		itemA := &items[i]
		if strings.TrimSpace(itemA.ItemID) == "" {
			return fmt.Errorf("item %d: item_id cannot be empty", i)
		}
		if itemA.RegisterOffset < 0 {
			return fmt.Errorf("item %s: register_offset cannot be negative", itemA.ItemID)
		}
		if !itemA.DataType.IsValid() {
			return fmt.Errorf("item %s: invalid data_type %s", itemA.ItemID, itemA.DataType)
		}
		if !itemA.SemanticKind.IsValid() {
			return fmt.Errorf("item %s: invalid semantic_kind %s", itemA.ItemID, itemA.SemanticKind)
		}
		if itemA.ScaleMultiplier != nil && (math.IsNaN(*itemA.ScaleMultiplier) || math.IsInf(*itemA.ScaleMultiplier, 0)) {
			return fmt.Errorf("item %s: scale_multiplier cannot be NaN or Inf", itemA.ItemID)
		}
		if itemA.ScaleOffset != nil && (math.IsNaN(*itemA.ScaleOffset) || math.IsInf(*itemA.ScaleOffset, 0)) {
			return fmt.Errorf("item %s: scale_offset cannot be NaN or Inf", itemA.ItemID)
		}

		wordsA := WordCountForDataType(itemA.DataType)
		startA := itemA.RegisterOffset
		endA := startA + wordsA

		// 檢查與其他項目的暫存器與位元重疊
		for j := i + 1; j < len(items); j++ {
			itemB := &items[j]
			wordsB := WordCountForDataType(itemB.DataType)
			startB := itemB.RegisterOffset
			endB := startB + wordsB

			// 檢查暫存器區間重疊
			if startA < endB && startB < endA {
				// 若為同一個暫存器內的 bit slices，允許共用 word
				if startA == startB && wordsA == 1 && wordsB == 1 && itemA.BitOffset != nil && itemB.BitOffset != nil {
					bitLenA := 1
					if itemA.BitLength != nil && *itemA.BitLength > 0 {
						bitLenA = *itemA.BitLength
					}
					bitLenB := 1
					if itemB.BitLength != nil && *itemB.BitLength > 0 {
						bitLenB = *itemB.BitLength
					}
					bStartA := *itemA.BitOffset
					bEndA := bStartA + bitLenA
					bStartB := *itemB.BitOffset
					bEndB := bStartB + bitLenB

					if bStartA < bEndB && bStartB < bEndA {
						return fmt.Errorf("bit slice conflict between item %s (bits %d..%d) and item %s (bits %d..%d) at register %d",
							itemA.ItemID, bStartA, bEndA-1, itemB.ItemID, bStartB, bEndB-1, startA)
					}
					continue
				}

				return fmt.Errorf("register overlap conflict between item %s (regs %d..%d) and item %s (regs %d..%d)",
					itemA.ItemID, startA, endA-1, itemB.ItemID, startB, endB-1)
			}
		}
	}

	return nil
}

// DecodeMixedRegisters 從原始暫存器陣列解碼出各 mixed 項目的工程值。
func DecodeMixedRegisters(registers []uint16, items []MixedItem, defaultFormat string) ([]DecodedItemResult, error) {
	results := make([]DecodedItemResult, len(items))

	for i, item := range items {
		words := WordCountForDataType(item.DataType)
		res := DecodedItemResult{
			ItemID:  item.ItemID,
			Quality: schema.QualityGood,
		}

		if item.RegisterOffset < 0 || item.RegisterOffset+words > len(registers) {
			res.Quality = schema.QualityMissing
			res.QualityReason = fmt.Sprintf("register offset %d..%d out of buffer range (len %d)",
				item.RegisterOffset, item.RegisterOffset+words, len(registers))
			results[i] = res
			continue
		}

		subRegs := registers[item.RegisterOffset : item.RegisterOffset+words]

		// Bit slice 提取
		if item.BitOffset != nil && words == 1 {
			w := subRegs[0]
			bo := *item.BitOffset
			bl := 1
			if item.BitLength != nil && *item.BitLength > 0 {
				bl = *item.BitLength
			}
			mask := uint16((1 << bl) - 1)
			rawBitVal := (w >> bo) & mask
			res.RawValue = rawBitVal

			if item.DataType == schema.DataTypeBool {
				res.Value = rawBitVal != 0
			} else {
				res.Value = rawBitVal
			}
			results[i] = res
			continue
		}

		// 暫存器轉換為位元組
		buf := make([]byte, words*2)
		for idx, r := range subRegs {
			binary.BigEndian.PutUint16(buf[idx*2:(idx+1)*2], r)
		}

		fmtStr := item.DataFormat
		if fmtStr == "" {
			fmtStr = defaultFormat
		}
		if fmtStr == "" {
			fmtStr = "ABCD"
		}

		conv := hsllogic.NewDataConverter(hsllogic.DataFormat(fmtStr))
		rawVal := conv.ReadValue(buf, 0, hsllogic.DataType(item.DataType))
		res.RawValue = rawVal

		// 應用有效換算（倍率與偏移）
		if item.ScaleMultiplier != nil || item.ScaleOffset != nil {
			numVal, ok := toFloat64(rawVal)
			if !ok {
				res.Quality = schema.QualityInvalid
				res.QualityReason = "failed to convert value to numeric for scaling"
				results[i] = res
				continue
			}

			mult := 1.0
			if item.ScaleMultiplier != nil {
				mult = *item.ScaleMultiplier
			}
			off := 0.0
			if item.ScaleOffset != nil {
				off = *item.ScaleOffset
			}

			scaledVal := numVal*mult + off
			res.Value = scaledVal
			res.Transformed = true
		} else {
			res.Value = rawVal
		}

		results[i] = res
	}

	return results, nil
}

func toFloat64(val interface{}) (float64, bool) {
	switch v := val.(type) {
	case int16:
		return float64(v), true
	case uint16:
		return float64(v), true
	case int32:
		return float64(v), true
	case uint32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint64:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return v, true
	case int:
		return float64(v), true
	default:
		return 0, false
	}
}
