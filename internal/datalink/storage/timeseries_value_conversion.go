package storage

import (
	"encoding/json"
	"fmt"
	"time"

	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// 值轉換輔助
// =============================================================================

// ValueToRecord 將收集值轉換為時序記錄
func ValueToRecord(tagID string, value interface{}, rawValue interface{}, timestamp time.Time, quality schema.QualityFlag, dataType schema.DataType) TimeSeriesRecord {
	record := TimeSeriesRecord{
		TagID:     tagID,
		Timestamp: timestamp,
		Quality:   quality,
	}

	// 設定原始值
	if rawValue != nil {
		record.RawValue = rawValue
	}

	// 根據資料型別設定值
	switch dataType {
	case schema.DataTypeBool:
		if v, ok := value.(bool); ok {
			record.ValueBool = &v
		}
	case schema.DataTypeString:
		if v, ok := value.(string); ok {
			record.ValueText = &v
		} else {
			s := fmt.Sprintf("%v", value)
			record.ValueText = &s
		}
	default:
		// 數值型別
		if v, ok := toFloat64(value); ok {
			record.ValueNum = &v
		}
	}

	return record
}

// toFloat64 將值轉換為 float64
func toFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int16:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case uint:
		return float64(val), true
	case uint16:
		return float64(val), true
	case uint32:
		return float64(val), true
	case uint64:
		return float64(val), true
	case bool:
		if val {
			return 1, true
		}
		return 0, true
	case json.Number:
		if f, err := val.Float64(); err == nil {
			return f, true
		}
	}
	return 0, false
}
