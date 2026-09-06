package aggregation

import (
	"math"
	"testing"
	"time"
)

func TestDerivedEvaluator_D01_DivisionByZeroAndValidation(t *testing.T) {
	// D01: 衍生公式除零/單位不合/循環 -> 阻擋或 invalid，不吐假數字
	formula := DerivedFormula{
		FormulaID:           "form-eff",
		OutputMeasurementID: "meas-eff",
		Operation:           "div",
		OperandLeft:         "meas-output-kw",
		OperandRight:        "meas-input-kw",
	}

	values := map[string]float64{
		"meas-output-kw": 50.0,
		"meas-input-kw":  0.0, // 除零
	}

	result := EvaluateDerivedFormula(formula, time.Now(), values)
	if result.IsValid {
		t.Errorf("expected invalid result on division by zero")
	}
	if result.ErrorMessage == "" {
		t.Errorf("expected error message for division by zero")
	}
}

func TestDerivedEvaluator_D01_CycleDetection(t *testing.T) {
	// A = B + C, B = A * 2 (循環)
	formulas := []DerivedFormula{
		{
			FormulaID:           "f1",
			OutputMeasurementID: "A",
			Operation:           "add",
			OperandLeft:         "B",
			OperandRight:        "C",
		},
		{
			FormulaID:           "f2",
			OutputMeasurementID: "B",
			Operation:           "mul",
			OperandLeft:         "A",
			OperandRight:        "2",
		},
	}

	err := ValidateFormulaDependencies(formulas)
	if err == nil {
		t.Errorf("expected error on cyclic dependency")
	}
}

func TestDerivedEvaluator_D02_ParentSubmeterOverlapCheck(t *testing.T) {
	// D02: 總表+其分表欲自動加總 -> 警告／阻擋重複計量
	totalMembers := []MeteringScopeMember{
		{MeasurementID: "m-main", ParentMeasurementID: ""},
		{MeasurementID: "m-sub-1", ParentMeasurementID: "m-main"}, // 子表
	}

	err := ValidateMeteringScope(totalMembers)
	if err == nil {
		t.Errorf("expected error on overlapping parent/child metering scope")
	}
}

func TestDerivedEvaluator_ValidArithmetic(t *testing.T) {
	formula := DerivedFormula{
		FormulaID:           "f-diff",
		OutputMeasurementID: "meas-delta-p",
		Operation:           "sub",
		OperandLeft:         "p_in",
		OperandRight:        "p_out",
	}

	values := map[string]float64{
		"p_in":  10.5,
		"p_out": 2.5,
	}

	result := EvaluateDerivedFormula(formula, time.Now(), values)
	if !result.IsValid {
		t.Fatalf("expected valid result: %s", result.ErrorMessage)
	}
	if result.ValueNumeric == nil || math.Abs(*result.ValueNumeric-8.0) > 1e-6 {
		t.Errorf("expected 8.0, got %v", result.ValueNumeric)
	}
}
