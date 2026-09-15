package adapters

import (
	"encoding/json"
	"fmt"
	"strconv"

	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// MQTT 主題映射輔助
// =============================================================================

// TopicMapping MQTT 主題到點位的映射配置
type TopicMapping struct {
	// Topic MQTT 主題 (支援萬用字元)
	Topic string `json:"topic"`

	// PayloadPath JSON 路徑 (用於從 payload 擷取值，如 "data.temperature")
	PayloadPath string `json:"payload_path,omitempty"`

	// DataType 資料型別
	DataType schema.DataType `json:"data_type"`
}

// ExtractValueFromPayload 從 MQTT payload 擷取值
func ExtractValueFromPayload(payload []byte, path string, dataType schema.DataType) (interface{}, error) {
	// 如果沒有路徑，直接解析整個 payload
	if path == "" {
		var value interface{}
		if err := json.Unmarshal(payload, &value); err != nil {
			// 非 JSON，嘗試直接解析為字串
			return string(payload), nil //nolint:nilerr // Plain-text MQTT payloads are supported when no JSON path is requested.
		}
		return convertMQTTValue(value, dataType), nil
	}

	// 解析 JSON 並按路徑擷取
	var obj map[string]interface{}
	if err := json.Unmarshal(payload, &obj); err != nil {
		return nil, fmt.Errorf("payload 不是 JSON 物件: %w", err)
	}

	// 簡單的點號路徑解析
	parts := splitPath(path)
	current := interface{}(obj)

	for _, part := range parts {
		switch v := current.(type) {
		case map[string]interface{}:
			var ok bool
			current, ok = v[part]
			if !ok {
				return nil, fmt.Errorf("路徑 '%s' 中找不到欄位 '%s'", path, part)
			}
		default:
			return nil, fmt.Errorf("路徑 '%s' 無效，'%s' 不是物件", path, part)
		}
	}

	return convertMQTTValue(current, dataType), nil
}

// splitPath 分割 JSON 路徑
func splitPath(path string) []string {
	// 簡單實作：按點號分割
	result := make([]string, 0)
	current := ""
	for _, ch := range path {
		if ch == '.' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

// convertMQTTValue 將 MQTT 值轉換為指定型別
func convertMQTTValue(value interface{}, dataType schema.DataType) interface{} {
	switch dataType {
	case schema.DataTypeBool:
		switch v := value.(type) {
		case bool:
			return v
		case float64:
			return v != 0
		case string:
			return v == "true" || v == "1" || v == "on"
		}
	case schema.DataTypeInt16, schema.DataTypeInt32, schema.DataTypeInt64:
		switch v := value.(type) {
		case float64:
			return int64(v)
		case string:
			if n, err := json.Number(v).Int64(); err == nil {
				return n
			}
		}
	case schema.DataTypeUint16, schema.DataTypeUint32, schema.DataTypeUint64:
		switch v := value.(type) {
		case float64:
			return uint64(v)
		case string:
			if n, err := strconv.ParseUint(v, 10, 64); err == nil {
				return n
			}
		}
	case schema.DataTypeFloat32, schema.DataTypeFloat64:
		switch v := value.(type) {
		case float64:
			return v
		case string:
			if n, err := json.Number(v).Float64(); err == nil {
				return n
			}
		}
	case schema.DataTypeString:
		switch v := value.(type) {
		case string:
			return v
		default:
			data, err := json.Marshal(v)
			if err != nil {
				return fmt.Sprintf("%v", v)
			}
			return string(data)
		}
	}

	return value
}
