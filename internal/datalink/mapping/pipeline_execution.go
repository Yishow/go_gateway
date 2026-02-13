package mapping

import (
	"encoding/json"
	"fmt"

	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// 轉換管線執行
// =============================================================================

// TransformContext 轉換上下文
type TransformContext struct {
	RawValue     interface{}
	CurrentValue interface{}
	StepResults  []StepResult
	Error        error
}

// StepResult 步驟結果
type StepResult struct {
	StepIndex int         `json:"step_index"`
	StepType  string      `json:"step_type"`
	Input     interface{} `json:"input"`
	Output    interface{} `json:"output"`
	Error     string      `json:"error,omitempty"`
}

// ExecutePipeline 執行轉換管線
func ExecutePipeline(raw interface{}, pipelineJSON string) (*TransformContext, error) {
	ctx := &TransformContext{
		RawValue:     raw,
		CurrentValue: raw,
		StepResults:  make([]StepResult, 0),
	}

	var steps []schema.TransformStep
	if err := json.Unmarshal([]byte(pipelineJSON), &steps); err != nil {
		ctx.Error = fmt.Errorf("解析轉換管線失敗: %w", err)
		return ctx, ctx.Error
	}

	orderedSteps := normalizeTransformSteps(steps)
	for i, step := range orderedSteps {
		result := StepResult{
			StepIndex: i,
			StepType:  string(step.Type),
			Input:     ctx.CurrentValue,
		}

		output, err := executeStep(ctx.CurrentValue, step)
		if err != nil {
			result.Error = err.Error()
			ctx.StepResults = append(ctx.StepResults, result)
			ctx.Error = fmt.Errorf("步驟 %d (%s) 失敗: %w", i+1, step.Type, err)
			return ctx, ctx.Error
		}

		result.Output = output
		ctx.StepResults = append(ctx.StepResults, result)
		ctx.CurrentValue = output
	}

	return ctx, nil
}

// executeStep 執行單一轉換步驟
func executeStep(input interface{}, step schema.TransformStep) (interface{}, error) {
	switch step.Type {
	case schema.TransformDecode:
		return executeDecode(input, step.Params)
	case schema.TransformCast:
		return executeCast(input, step.Params)
	case schema.TransformScale:
		return executeScale(input, step.Params)
	case schema.TransformLookup:
		return executeLookup(input, step.Params)
	case schema.TransformConditional:
		return executeConditional(input, step.Params)
	case schema.TransformFormula:
		return executeFormula(input, step.Params)
	default:
		return nil, fmt.Errorf("不支援的轉換類型: %s", step.Type)
	}
}

// =============================================================================
// 轉換步驟實作 (使用 map[string]interface{} 參數)
// =============================================================================

// executeDecode 解碼轉換
func executeDecode(input interface{}, params map[string]interface{}) (interface{}, error) {
	if params == nil {
		return input, nil
	}

	format, _ := params["format"].(string)
	switch format {
	case "swap16":
		if v, ok := toUint16Value(input); ok {
			return ((v & 0xFF) << 8) | ((v >> 8) & 0xFF), nil
		}
	case "swap32":
		if v, ok := toUint32Value(input); ok {
			return ((v&0xFF)<<24 | ((v>>8)&0xFF)<<16 | ((v>>16)&0xFF)<<8 | (v >> 24)), nil
		}
	}

	return input, nil
}

// executeCast 型別轉換
func executeCast(input interface{}, params map[string]interface{}) (interface{}, error) {
	if params == nil {
		return input, nil
	}

	toType, _ := params["to_type"].(string)
	if toType == "" {
		toType, _ = params["target_type"].(string)
	}
	targetType := schema.DataType(toType)

	switch targetType {
	case schema.DataTypeBool:
		return toBoolValue(input), nil
	case schema.DataTypeInt16:
		v, _ := toFloat64Value(input)
		return int16(v), nil
	case schema.DataTypeUint16:
		v, _ := toFloat64Value(input)
		return uint16(v), nil
	case schema.DataTypeInt32:
		v, _ := toFloat64Value(input)
		return int32(v), nil
	case schema.DataTypeUint32:
		v, _ := toFloat64Value(input)
		return uint32(v), nil
	case schema.DataTypeInt64:
		v, _ := toFloat64Value(input)
		return int64(v), nil
	case schema.DataTypeUint64:
		v, _ := toFloat64Value(input)
		return uint64(v), nil
	case schema.DataTypeFloat32:
		v, _ := toFloat64Value(input)
		return float32(v), nil
	case schema.DataTypeFloat64:
		v, _ := toFloat64Value(input)
		return v, nil
	case schema.DataTypeString:
		return fmt.Sprintf("%v", input), nil
	default:
		return input, nil
	}
}

// executeScale 縮放轉換
func executeScale(input interface{}, params map[string]interface{}) (interface{}, error) {
	if params == nil {
		return input, nil
	}

	v, ok := toFloat64Value(input)
	if !ok {
		return input, fmt.Errorf("無法轉換為數值: %v", input)
	}

	multiplier := 1.0
	if m, ok := toFloat64Value(params["multiplier"]); ok {
		multiplier = m
	} else if m, ok := toFloat64Value(params["scale"]); ok {
		multiplier = m
	}

	divisor := 1.0
	if d, ok := toFloat64Value(params["divisor"]); ok && d != 0 {
		divisor = d
	}

	offset := 0.0
	if o, ok := toFloat64Value(params["offset"]); ok {
		offset = o
	}

	result := (v*multiplier)/divisor + offset

	if minVal, ok := toFloat64Value(params["min"]); ok && result < minVal {
		result = minVal
	}
	if maxVal, ok := toFloat64Value(params["max"]); ok && result > maxVal {
		result = maxVal
	}

	return result, nil
}

// executeLookup 查表轉換
func executeLookup(input interface{}, params map[string]interface{}) (interface{}, error) {
	if params == nil {
		return input, nil
	}

	table, _ := params["table"].(map[string]interface{})
	if table == nil {
		return input, nil
	}

	key := fmt.Sprintf("%v", input)

	if value, exists := table[key]; exists {
		return value, nil
	}

	if defaultVal, exists := params["default"]; exists {
		return defaultVal, nil
	}

	return input, nil
}
