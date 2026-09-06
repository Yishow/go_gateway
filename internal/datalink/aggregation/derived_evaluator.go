package aggregation

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// DerivedFormula 衍生計算公式模型。
type DerivedFormula struct {
	FormulaID           string `json:"formula_id"`
	OutputMeasurementID string `json:"output_measurement_id"`
	Operation           string `json:"operation"` // "add", "sub", "mul", "div"
	OperandLeft         string `json:"operand_left"`
	OperandRight        string `json:"operand_right"`
	CalculationRevision int64  `json:"calculation_revision"`
}

// MeteringScopeMember 計量範圍成員。
type MeteringScopeMember struct {
	MeasurementID       string `json:"measurement_id"`
	ParentMeasurementID string `json:"parent_measurement_id,omitempty"`
}

// ValidateMeteringScope 驗證計量範圍中是否包含重疊的母表與子表（防止重複計量）。
func ValidateMeteringScope(members []MeteringScopeMember) error {
	memberMap := make(map[string]bool)
	for _, m := range members {
		memberMap[m.MeasurementID] = true
	}

	for _, m := range members {
		if m.ParentMeasurementID != "" && memberMap[m.ParentMeasurementID] {
			return fmt.Errorf("overlapping metering scope: %s is a sub-meter of %s", m.MeasurementID, m.ParentMeasurementID)
		}
	}
	return nil
}

// ValidateFormulaDependencies 檢查公式之間的循環依賴 (DAG Cycle Detection)。
func ValidateFormulaDependencies(formulas []DerivedFormula) error {
	adj := make(map[string][]string)
	for _, f := range formulas {
		out := f.OutputMeasurementID
		if f.OperandLeft != "" && !isNumber(f.OperandLeft) {
			adj[out] = append(adj[out], f.OperandLeft)
		}
		if f.OperandRight != "" && !isNumber(f.OperandRight) {
			adj[out] = append(adj[out], f.OperandRight)
		}
	}

	visited := make(map[string]int) // 0: unvisited, 1: visiting, 2: visited

	var dfs func(node string) error
	dfs = func(node string) error {
		visited[node] = 1
		for _, neighbor := range adj[node] {
			if visited[neighbor] == 1 {
				return fmt.Errorf("cyclic dependency detected involving %s and %s", node, neighbor)
			}
			if visited[neighbor] == 0 {
				if err := dfs(neighbor); err != nil {
					return err
				}
			}
		}
		visited[node] = 2
		return nil
	}

	for _, f := range formulas {
		if visited[f.OutputMeasurementID] == 0 {
			if err := dfs(f.OutputMeasurementID); err != nil {
				return err
			}
		}
	}
	return nil
}

func isNumber(s string) bool {
	_, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return err == nil
}

func resolveOperand(op string, values map[string]float64) (float64, error) {
	if isNumber(op) {
		return strconv.ParseFloat(strings.TrimSpace(op), 64)
	}
	v, ok := values[op]
	if !ok {
		return 0, fmt.Errorf("missing operand value for %s", op)
	}
	return v, nil
}

// EvaluateDerivedFormula 依白名單運算符安全執行衍生公式。
func EvaluateDerivedFormula(
	formula DerivedFormula,
	observedAt time.Time,
	values map[string]float64,
) DerivedCalculationResult {
	result := DerivedCalculationResult{
		FormulaID:           formula.FormulaID,
		OutputMeasurementID: formula.OutputMeasurementID,
		ObservedAt:          observedAt,
		CalculationRevision: formula.CalculationRevision,
	}

	leftVal, err := resolveOperand(formula.OperandLeft, values)
	if err != nil {
		result.ErrorMessage = err.Error()
		return result
	}

	rightVal, err := resolveOperand(formula.OperandRight, values)
	if err != nil {
		result.ErrorMessage = err.Error()
		return result
	}

	var res float64
	switch strings.ToLower(formula.Operation) {
	case "add", "+":
		res = leftVal + rightVal
	case "sub", "-":
		res = leftVal - rightVal
	case "mul", "*":
		res = leftVal * rightVal
	case "div", "/":
		if rightVal == 0 {
			result.ErrorMessage = "division by zero"
			return result
		}
		res = leftVal / rightVal
	default:
		result.ErrorMessage = fmt.Sprintf("unsupported operation: %s", formula.Operation)
		return result
	}

	result.ValueNumeric = &res
	result.IsValid = true
	return result
}
