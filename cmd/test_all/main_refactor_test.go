package main

import (
	"testing"
	"time"
)

func TestTestResultsAddResult(t *testing.T) {
	results := &TestResults{}

	results.addResult(TestResult{Protocol: "Modbus", Mode: "TCP", Success: true, Duration: time.Millisecond})
	results.addResult(TestResult{Protocol: "Fatek", Mode: "RTU", Success: false, Duration: 2 * time.Millisecond})

	if results.Total != 2 {
		t.Fatalf("Total 預期 2，實際 %d", results.Total)
	}
	if results.Passed != 1 {
		t.Fatalf("Passed 預期 1，實際 %d", results.Passed)
	}
	if results.Failed != 1 {
		t.Fatalf("Failed 預期 1，實際 %d", results.Failed)
	}
	if len(results.Results) != 2 {
		t.Fatalf("Results 長度預期 2，實際 %d", len(results.Results))
	}
}
