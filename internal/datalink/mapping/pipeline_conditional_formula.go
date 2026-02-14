package mapping

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// executeConditional 條件轉換
func executeConditional(input interface{}, params map[string]interface{}) (interface{}, error) {
	if params == nil {
		return input, nil
	}

	v, _ := toFloat64Value(input)

	operator, threshold, err := parseConditionalParams(params)
	if err != nil {
		return input, err
	}
	trueValue := params["true_value"]
	falseValue := params["false_value"]

	result := evaluateCondition(v, operator, threshold)

	if result {
		return trueValue, nil
	}
	return falseValue, nil
}

// evaluateCondition 評估條件
func evaluateCondition(value float64, operator string, threshold float64) bool {
	switch operator {
	case "==", "eq":
		return value == threshold
	case "!=", "ne":
		return value != threshold
	case ">", "gt":
		return value > threshold
	case ">=", "gte":
		return value >= threshold
	case "<", "lt":
		return value < threshold
	case "<=", "lte":
		return value <= threshold
	default:
		return false
	}
}

func parseConditionalParams(params map[string]interface{}) (string, float64, error) {
	if params == nil {
		return "", 0, fmt.Errorf("缺少參數")
	}

	if condition, ok := params["condition"].(string); ok && strings.TrimSpace(condition) != "" {
		return parseConditionString(condition)
	}

	operator := ""
	if rawOperator, exists := params["operator"]; exists {
		if op, ok := rawOperator.(string); ok {
			operator = op
		}
	}
	operator = normalizeOperator(strings.TrimSpace(operator))
	if operator == "" {
		return "", 0, fmt.Errorf("缺少 operator 或 condition")
	}

	threshold, ok := toFloat64Value(params["threshold"])
	if !ok {
		return "", 0, fmt.Errorf("無效的 threshold")
	}

	return operator, threshold, nil
}

func normalizeOperator(operator string) string {
	switch operator {
	case "eq":
		return "=="
	case "ne":
		return "!="
	case "gt":
		return ">"
	case "gte":
		return ">="
	case "lt":
		return "<"
	case "lte":
		return "<="
	default:
		return operator
	}
}

func parseConditionString(condition string) (string, float64, error) {
	cond := strings.TrimSpace(condition)
	if strings.HasPrefix(cond, "value") {
		cond = strings.TrimSpace(strings.TrimPrefix(cond, "value"))
	}
	if strings.HasPrefix(cond, "x") {
		cond = strings.TrimSpace(strings.TrimPrefix(cond, "x"))
	}

	operators := []string{"==", "!=", ">=", "<=", ">", "<"}
	for _, op := range operators {
		idx := strings.Index(cond, op)
		if idx < 0 {
			continue
		}
		left := strings.TrimSpace(cond[:idx])
		right := strings.TrimSpace(cond[idx+len(op):])
		if left != "" && left != "value" && left != "x" {
			return "", 0, fmt.Errorf("無法解析 condition: %s", condition)
		}
		threshold, ok := toFloat64Value(right)
		if !ok {
			return "", 0, fmt.Errorf("無法解析門檻值: %s", right)
		}
		return op, threshold, nil
	}

	return "", 0, fmt.Errorf("無法解析 condition: %s", condition)
}

// executeFormula 公式轉換
func executeFormula(input interface{}, params map[string]interface{}) (interface{}, error) {
	if params == nil {
		return input, nil
	}

	expression := ""
	if rawExpression, exists := params["expression"]; exists {
		if expr, ok := rawExpression.(string); ok {
			expression = expr
		}
	}
	if expression == "" {
		return input, nil
	}

	v, ok := toFloat64Value(input)
	if !ok {
		return input, fmt.Errorf("無法轉換為數值: %v", input)
	}

	expr := strings.TrimSpace(expression)
	expr = strings.ReplaceAll(expr, "value", strconv.FormatFloat(v, 'f', -1, 64))
	expr = strings.ReplaceAll(expr, "x", strconv.FormatFloat(v, 'f', -1, 64))

	result, err := evaluateSimpleExpression(expr, v)
	if err != nil {
		return input, err
	}

	return result, nil
}

// evaluateSimpleExpression 評估簡易表達式
func evaluateSimpleExpression(expr string, value float64) (float64, error) {
	expr = strings.TrimSpace(expr)

	if strings.HasPrefix(expr, "abs(") {
		return math.Abs(value), nil
	}
	if strings.HasPrefix(expr, "sqrt(") {
		return math.Sqrt(value), nil
	}
	if strings.HasPrefix(expr, "round(") {
		return math.Round(value), nil
	}
	if strings.HasPrefix(expr, "floor(") {
		return math.Floor(value), nil
	}
	if strings.HasPrefix(expr, "ceil(") {
		return math.Ceil(value), nil
	}

	if v, err := strconv.ParseFloat(expr, 64); err == nil {
		return v, nil
	}

	return value, nil
}
