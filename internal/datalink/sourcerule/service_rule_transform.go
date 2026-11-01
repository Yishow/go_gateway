package sourcerule

import "go-gateway/internal/datalink/schema"

// buildRuleTransformPipeline 根據來源規則建構 Mapping 的轉換管線。
// 順序：先 cast（若目標型別與 Point 型別不同），再 scale（若有設定）。
func (s *Service) buildRuleTransformPipeline(rule *schema.SourceRule, pointRecord *schema.Point) []schema.TransformStep {
	steps := make([]schema.TransformStep, 0)
	order := 0

	targetDataType := desiredRuleTargetDataType(rule, pointRecord)

	// 若目標型別與 Point 型別不同，加入 cast 步驟
	if targetDataType != pointRecord.DataType && targetDataType.IsValid() {
		steps = append(steps, schema.TransformStep{
			Type:  schema.TransformCast,
			Order: order,
			Params: map[string]any{
				"target_type": string(targetDataType),
			},
		})
		order++
	}

	// 若有 scale 設定，加入 scale 步驟
	if rule.ScaleMultiplier != nil || rule.ScaleOffset != nil {
		scale := 1.0
		if rule.ScaleMultiplier != nil {
			scale = *rule.ScaleMultiplier
		}
		offset := 0.0
		if rule.ScaleOffset != nil {
			offset = *rule.ScaleOffset
		}

		steps = append(steps, schema.TransformStep{
			Type:  schema.TransformScale,
			Order: order,
			Params: map[string]any{
				"scale":  scale,
				"offset": offset,
			},
		})
	}

	return steps
}
