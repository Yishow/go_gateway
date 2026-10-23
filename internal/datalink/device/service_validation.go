package device

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// 輔助函數
// =============================================================================

// isValidProtocol 檢查協議類型是否有效
func isValidProtocol(protocol schema.ProtocolType) bool {
	return connector.IsRegistered(protocol)
}

// validateConnectionConfig 驗證連線配置
func validateConnectionConfig(protocol schema.ProtocolType, config map[string]any) error {
	// 取得協議資訊
	info, ok := connector.GetProtocolInfo(protocol)
	if !ok {
		return fmt.Errorf("未知的協議類型: %s", protocol)
	}

	// 解析配置 Schema
	var configSchema map[string]any
	if err := json.Unmarshal(info.ConfigSchema, &configSchema); err != nil {
		return fmt.Errorf("invalid protocol configuration schema: %w", err)
	}

	// 取得 properties
	properties, ok := configSchema["properties"].(map[string]any)
	if !ok {
		properties = make(map[string]any)
	}

	// 檢查必填欄位：缺少與空值一次收集，不在第一個錯誤就返回
	var (
		missingFields []string
		emptyFields   []string
		valueErr      error // 第一個數值範圍或 enum 錯誤，僅在必填欄位皆齊全時回報
	)
	required, ok := configSchema["required"].([]any)
	if !ok {
		required = nil // 沒有 required 或型別不符時，視為沒有必填欄位，不進入迴圈
	}
	for _, field := range required {
		fieldName, ok := field.(string)
		if !ok {
			continue
		}

		if isConditionallyOptional(protocol, fieldName, config) {
			continue
		}

		// 檢查欄位是否存在
		value, exists := config[fieldName]
		if !exists {
			missingFields = append(missingFields, fieldName)
			continue
		}

		// 檢查值是否為 nil（JSON null）或空字串
		if isEmptyValue(value) {
			emptyFields = append(emptyFields, fieldName)
			continue
		}

		// 驗證數字範圍和 enum，只保留第一個錯誤以維持原有的回報順序
		if prop, ok := properties[fieldName].(map[string]any); ok && valueErr == nil {
			valueErr = validateFieldValue(fieldName, value, prop)
		}
	}

	// 必填問題一次回報：缺少在前、空值在後，修正後才回報數值範圍／enum。
	var requiredErrs []error
	if len(missingFields) > 0 {
		requiredErrs = append(requiredErrs, fmt.Errorf("缺少必填欄位: %s", strings.Join(missingFields, ", ")))
	}
	if len(emptyFields) == 1 {
		requiredErrs = append(requiredErrs, fmt.Errorf("必填欄位 %s 不能為空", emptyFields[0]))
	}
	if len(emptyFields) > 1 {
		requiredErrs = append(requiredErrs, fmt.Errorf("必填欄位不能為空: %s", strings.Join(emptyFields, ", ")))
	}
	if len(requiredErrs) > 0 {
		return errors.Join(requiredErrs...)
	}

	return valueErr
}

// isConditionallyOptional 判斷必填欄位是否因其他欄位的值而不需檢查
// （如 fatek_fbs 的 host/port 在 mode=serial 時非必填）
func isConditionallyOptional(protocol schema.ProtocolType, fieldName string, config map[string]any) bool {
	if protocol != schema.ProtocolFatekFBs {
		return false
	}

	mode, ok := config["mode"].(string)
	if !ok {
		return false
	}

	switch fieldName {
	case "host", "port":
		return mode == "serial" // serial 模式不需要 host 和 port
	case "serial_port", "baud_rate":
		return mode != "serial" // 非 serial 模式不需要 serial_port 和 baud_rate
	}
	return false
}

// isEmptyValue 判斷必填欄位的值是否為空（nil 或空字串）
func isEmptyValue(value any) bool {
	if value == nil {
		return true
	}
	strValue, ok := value.(string)
	return ok && strValue == ""
}

// validateFieldValue 依欄位 Schema 驗證數字範圍與 enum
func validateFieldValue(fieldName string, value any, prop map[string]any) error {
	// 檢查數字範圍
	if numValue, ok := value.(float64); ok {
		if minimum, ok := prop["minimum"].(float64); ok {
			if numValue < minimum {
				return fmt.Errorf("欄位 %s 的值 %v 小於最小值 %v", fieldName, numValue, minimum)
			}
		}
		if maximum, ok := prop["maximum"].(float64); ok {
			if numValue > maximum {
				return fmt.Errorf("欄位 %s 的值 %v 大於最大值 %v", fieldName, numValue, maximum)
			}
		}
		// 檢查 enum（對於數字類型）
		if enum, ok := prop["enum"].([]any); ok {
			valid := slices.ContainsFunc(enum, func(e any) bool {
				eNum, ok := e.(float64)
				return ok && eNum == numValue
			})
			if !valid {
				return fmt.Errorf("欄位 %s 的值 %v 不在允許的選項中", fieldName, numValue)
			}
		}
	}

	// 檢查字串 enum
	if strValue, ok := value.(string); ok {
		if enum, ok := prop["enum"].([]any); ok {
			valid := slices.ContainsFunc(enum, func(e any) bool {
				eStr, ok := e.(string)
				return ok && eStr == strValue
			})
			if !valid {
				return fmt.Errorf("欄位 %s 的值 %s 不在允許的選項中", fieldName, strValue)
			}
		}
	}

	return nil
}
