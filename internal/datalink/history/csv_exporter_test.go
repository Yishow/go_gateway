package history

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestCSVExporter_ExportPoints(t *testing.T) {
	baseTime := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	points := []HistoryPoint{
		{
			MeasurementID: "meas-kw",
			ObservedAt:    baseTime,
			ValueNumeric:  floatPtr(120.5),
			Quality:       "good",
			CoverageRatio: 1.0,
		},
		{
			MeasurementID: "meas-kwh",
			ObservedAt:    baseTime.Add(10 * time.Minute),
			UsageDelta:    floatPtr(0.7),
			Quality:       "good",
			CoverageRatio: 1.0,
		},
	}

	var buf bytes.Buffer
	err := ExportCSV(&buf, points)
	if err != nil {
		t.Fatalf("export CSV failed: %v", err)
	}

	csvStr := buf.String()
	if !strings.Contains(csvStr, "observed_at,measurement_id,value_numeric") {
		t.Errorf("expected CSV header")
	}
	if !strings.Contains(csvStr, "120.5") || !strings.Contains(csvStr, "0.7") {
		t.Errorf("expected numeric values in CSV")
	}
}

func TestCSVExporter_FormulaInjectionSanitization(t *testing.T) {
	// Formula injection defense: 对以 =, +, -, @ 开头的恶意字串进行前缀单引号跳脱
	maliciousStr := "=SUM(A1:A10)"
	points := []HistoryPoint{
		{
			MeasurementID: "meas-text",
			ObservedAt:    time.Now().UTC(),
			ValueString:   &maliciousStr,
			Quality:       "good",
		},
	}

	var buf bytes.Buffer
	err := ExportCSV(&buf, points)
	if err != nil {
		t.Fatalf("export CSV failed: %v", err)
	}

	csvStr := buf.String()
	// 期望跳脫成 '=SUM(A1:A10)
	if !strings.Contains(csvStr, "'=SUM(A1:A10)") {
		t.Errorf("expected formula injection sanitization, got: %s", csvStr)
	}
}
