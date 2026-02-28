package mapping

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/schema"
)

func (s *Service) validateMappingPreviewGate(ctx context.Context, tagID, pipelineJSON string, previewRawValue interface{}) error {
	normalizedPipeline := normalizePipelineJSON(pipelineJSON)

	var steps []schema.TransformStep
	if err := json.Unmarshal([]byte(normalizedPipeline), &steps); err != nil {
		return fmt.Errorf("pipeline 解析失敗: %w", err)
	}
	if err := ValidateTransformPipeline(steps); err != nil {
		return fmt.Errorf("pipeline 驗證失敗: %w", err)
	}

	sample := previewRawValue
	if sample == nil {
		sample = defaultPreviewSampleValue()
	}
	result, err := ExecutePipeline(sample, normalizedPipeline)
	if err != nil {
		return fmt.Errorf("sample 執行失敗: %w", err)
	}

	targetType, err := s.resolveTagDataType(ctx, tagID)
	if err != nil {
		return fmt.Errorf("取得 tag 資料型別失敗: %w", err)
	}

	if err := ensureCastableToDataType(result.CurrentValue, targetType); err != nil {
		return fmt.Errorf("輸出無法轉換為 tag data_type=%s: %w", targetType, err)
	}
	return nil
}

func normalizePipelineJSON(pipelineJSON string) string {
	if strings.TrimSpace(pipelineJSON) == "" {
		return "[]"
	}
	return pipelineJSON
}

func defaultPreviewSampleValue() interface{} {
	return float64(1)
}

func (s *Service) resolveTagDataType(ctx context.Context, tagID string) (schema.DataType, error) {
	s.mu.RLock()
	resolver := s.tagResolver
	s.mu.RUnlock()
	if resolver == nil {
		return "", ErrTagResolverNotConfigured
	}

	tag, err := resolver(ctx, tagID)
	if err != nil {
		return "", err
	}
	if tag == nil {
		return "", fmt.Errorf("tag resolver 回傳空值: tag_id=%s", tagID)
	}
	if tag.DataType == "" {
		return "", fmt.Errorf("tag data_type 為空: tag_id=%s", tagID)
	}
	return tag.DataType, nil
}

func ensureCastableToDataType(value interface{}, targetType schema.DataType) error {
	switch targetType {
	case schema.DataTypeString:
		return nil
	case schema.DataTypeBool:
		return ensureBoolCastable(value)
	case schema.DataTypeInt16:
		return ensureIntegerRange(value, -32768, 32767, "int16")
	case schema.DataTypeUint16:
		return ensureUnsignedRange(value, 65535, "uint16")
	case schema.DataTypeInt32:
		return ensureIntegerRange(value, -2147483648, 2147483647, "int32")
	case schema.DataTypeUint32:
		return ensureUnsignedRange(value, 4294967295, "uint32")
	case schema.DataTypeInt64:
		return ensureNumeric(value, "int64")
	case schema.DataTypeUint64:
		return ensureNonNegativeNumeric(value, "uint64")
	case schema.DataTypeFloat32, schema.DataTypeFloat64:
		if _, ok := toFloat64Value(value); !ok {
			return fmt.Errorf("值 %v (%T) 不是可轉數值", value, value)
		}
		return nil
	default:
		return fmt.Errorf("不支援的目標型別: %s", targetType)
	}
}

func ensureBoolCastable(value interface{}) error {
	switch v := value.(type) {
	case bool:
		return nil
	case string:
		switch strings.ToLower(strings.TrimSpace(v)) {
		case "true", "1", "on", "yes", "false", "0", "off", "no":
			return nil
		default:
			return fmt.Errorf("字串 %q 不是合法布林值", v)
		}
	default:
		f, ok := toFloat64Value(value)
		if !ok {
			return fmt.Errorf("值 %v (%T) 不是可轉布林類型", value, value)
		}
		if f != 0 && f != 1 {
			return fmt.Errorf("數值 %v 不是 0/1", value)
		}
		return nil
	}
}

func ensureIntegerRange(value interface{}, min, max float64, typeName string) error {
	f, ok := toFloat64Value(value)
	if !ok {
		return fmt.Errorf("值 %v (%T) 不是可轉數值", value, value)
	}
	if f < min || f > max {
		return fmt.Errorf("數值 %v 超出 %s 範圍", value, typeName)
	}
	return nil
}

func ensureUnsignedRange(value interface{}, max float64, typeName string) error {
	f, ok := toFloat64Value(value)
	if !ok {
		return fmt.Errorf("值 %v (%T) 不是可轉數值", value, value)
	}
	if f < 0 || f > max {
		return fmt.Errorf("數值 %v 超出 %s 範圍", value, typeName)
	}
	return nil
}

func ensureNumeric(value interface{}, typeName string) error {
	if _, ok := toFloat64Value(value); !ok {
		return fmt.Errorf("值 %v (%T) 不是可轉 %s 的數值", value, value, typeName)
	}
	return nil
}

func ensureNonNegativeNumeric(value interface{}, typeName string) error {
	f, ok := toFloat64Value(value)
	if !ok {
		return fmt.Errorf("值 %v (%T) 不是可轉 %s 的數值", value, value, typeName)
	}
	if f < 0 {
		return fmt.Errorf("數值 %v 不可為負數", value)
	}
	return nil
}
