package mapping

import (
	"fmt"

	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// 轉換管線驗證
// =============================================================================

// ValidateTransformPipeline 驗證轉換管線
func ValidateTransformPipeline(steps []schema.TransformStep) error {
	validTypes := map[schema.TransformType]bool{
		schema.TransformDecode:      true,
		schema.TransformCast:        true,
		schema.TransformScale:       true,
		schema.TransformLookup:      true,
		schema.TransformConditional: true,
		schema.TransformFormula:     true,
	}

	for i, step := range steps {
		if step.Order < 0 {
			return fmt.Errorf("步驟 %d: order 不能為負數", i+1)
		}
		if !validTypes[step.Type] {
			return fmt.Errorf("步驟 %d: 不支援的轉換類型 '%s'", i+1, step.Type)
		}

		switch step.Type {
		case schema.TransformScale:
			if step.Params == nil {
				return fmt.Errorf("步驟 %d: scale 需要參數", i+1)
			}
		case schema.TransformCast:
			if step.Params == nil {
				return fmt.Errorf("步驟 %d: cast 需要參數", i+1)
			}
		case schema.TransformFormula:
			if step.Params == nil {
				return fmt.Errorf("步驟 %d: formula 需要參數", i+1)
			}
		case schema.TransformConditional:
			if step.Params == nil {
				return fmt.Errorf("步驟 %d: conditional 需要參數", i+1)
			}
			if _, exists := step.Params["true_value"]; !exists {
				return fmt.Errorf("步驟 %d: conditional 缺少 true_value", i+1)
			}
			if _, exists := step.Params["false_value"]; !exists {
				return fmt.Errorf("步驟 %d: conditional 缺少 false_value", i+1)
			}
			if _, _, err := parseConditionalParams(step.Params); err != nil {
				return fmt.Errorf("步驟 %d: conditional 參數無效: %w", i+1, err)
			}
		}
	}

	return nil
}
