package device

import (
	"encoding/json"
	"fmt"

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
func validateConnectionConfig(protocol schema.ProtocolType, config map[string]interface{}) error {
	// 取得協議資訊
	info, ok := connector.GetProtocolInfo(protocol)
	if !ok {
		return fmt.Errorf("未知的協議類型: %s", protocol)
	}

	// 解析配置 Schema
	var configSchema map[string]interface{}
	if err := json.Unmarshal(info.ConfigSchema, &configSchema); err != nil {
		return nil // 無 Schema 驗證
	}

	// 取得 properties
	properties, ok := configSchema["properties"].(map[string]interface{})
	if !ok {
		properties = make(map[string]interface{})
	}

	// 檢查必填欄位
	required, ok := configSchema["required"].([]interface{})
	if ok {
		for _, field := range required {
			fieldName, ok := field.(string)
			if !ok {
				continue
			}

			// 處理條件必填欄位（如 fatek_fbs 的 host/port 在 mode=serial 時非必填）
			if protocol == schema.ProtocolFatekFBs {
				if fieldName == "host" || fieldName == "port" {
					if mode, ok := config["mode"].(string); ok && mode == "serial" {
						continue // 跳過 host 和 port 的必填檢查
					}
				}
				if fieldName == "serial_port" || fieldName == "baud_rate" {
					if mode, ok := config["mode"].(string); ok && mode != "serial" {
						continue // 跳過 serial_port 和 baud_rate 的必填檢查
					}
				}
			}

			// 檢查欄位是否存在
			value, exists := config[fieldName]
			if !exists {
				return fmt.Errorf("缺少必填欄位: %s", fieldName)
			}

			// 檢查值是否為 nil（JSON null）
			if value == nil {
				return fmt.Errorf("必填欄位 %s 不能為空", fieldName)
			}

			// 對於字串類型，檢查是否為空字串
			if strValue, ok := value.(string); ok && strValue == "" {
				return fmt.Errorf("必填欄位 %s 不能為空", fieldName)
			}

			// 驗證數字範圍和 enum
			if prop, ok := properties[fieldName].(map[string]interface{}); ok {
				// 檢查數字範圍
				if numValue, ok := value.(float64); ok {
					if min, ok := prop["minimum"].(float64); ok {
						if numValue < min {
							return fmt.Errorf("欄位 %s 的值 %v 小於最小值 %v", fieldName, numValue, min)
						}
					}
					if max, ok := prop["maximum"].(float64); ok {
						if numValue > max {
							return fmt.Errorf("欄位 %s 的值 %v 大於最大值 %v", fieldName, numValue, max)
						}
					}
					// 檢查 enum（對於數字類型）
					if enum, ok := prop["enum"].([]interface{}); ok {
						valid := false
						for _, e := range enum {
							if eNum, ok := e.(float64); ok && eNum == numValue {
								valid = true
								break
							}
						}
						if !valid {
							return fmt.Errorf("欄位 %s 的值 %v 不在允許的選項中", fieldName, numValue)
						}
					}
				}

				// 檢查字串 enum
				if strValue, ok := value.(string); ok {
					if enum, ok := prop["enum"].([]interface{}); ok {
						valid := false
						for _, e := range enum {
							if eStr, ok := e.(string); ok && eStr == strValue {
								valid = true
								break
							}
						}
						if !valid {
							return fmt.Errorf("欄位 %s 的值 %s 不在允許的選項中", fieldName, strValue)
						}
					}
				}
			}
		}
	}

	return nil
}
