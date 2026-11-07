package api

import (
	"fmt"
	"testing"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// sqlServerCircuitGateReason evaluates whether the SQL Server circuit test fixture
// should skip due to capability gating.
func sqlServerCircuitGateReason(mode workspace.WriteGroupStorageMode, capability *dbtarget.SQLServerRecordingCapability) string {
	if capability == nil || !capability.SupportsManagedLong {
		return "NOT_RUN: SQL Server managed-long operation gate is closed"
	}
	if mode == workspace.WriteGroupStorageModeNativeMonthly && !capability.SupportsNativePartition {
		return "NOT_RUN: SQL Server native partition operation gate is closed"
	}
	return ""
}

type circuitGateReporter interface {
	Helper()
	Skip(...any)
}

func checkSQLServerCircuitGate(t circuitGateReporter, mode workspace.WriteGroupStorageMode, capability *dbtarget.SQLServerRecordingCapability) {
	t.Helper()
	if reason := sqlServerCircuitGateReason(mode, capability); reason != "" {
		t.Skip(reason)
	}
}

func TestSQLServerCircuitNativeGate(t *testing.T) {
	tests := []struct {
		name       string
		mode       workspace.WriteGroupStorageMode
		capability *dbtarget.SQLServerRecordingCapability
		wantReason string
		wantSkip   bool
	}{
		{
			name: "plain mode with managed long supported is open",
			mode: workspace.WriteGroupStorageModePlain,
			capability: &dbtarget.SQLServerRecordingCapability{
				SupportsManagedLong:     true,
				SupportsNativePartition: false,
			},
			wantReason: "",
			wantSkip:   false,
		},
		{
			name: "plain mode without managed long is closed",
			mode: workspace.WriteGroupStorageModePlain,
			capability: &dbtarget.SQLServerRecordingCapability{
				SupportsManagedLong:     false,
				SupportsNativePartition: false,
			},
			wantReason: "NOT_RUN: SQL Server managed-long operation gate is closed",
			wantSkip:   true,
		},
		{
			name: "native monthly mode without native partition support is closed",
			mode: workspace.WriteGroupStorageModeNativeMonthly,
			capability: &dbtarget.SQLServerRecordingCapability{
				SupportsManagedLong:     true,
				SupportsNativePartition: false,
			},
			wantReason: "NOT_RUN: SQL Server native partition operation gate is closed",
			wantSkip:   true,
		},
		{
			name: "native monthly mode with both supported is open",
			mode: workspace.WriteGroupStorageModeNativeMonthly,
			capability: &dbtarget.SQLServerRecordingCapability{
				SupportsManagedLong:     true,
				SupportsNativePartition: true,
			},
			wantReason: "",
			wantSkip:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotReason := sqlServerCircuitGateReason(tt.mode, tt.capability)
			require.Equal(t, tt.wantReason, gotReason)

			reporter := &circuitGateRecorder{}
			checkSQLServerCircuitGate(reporter, tt.mode, tt.capability)
			require.Equal(t, tt.wantSkip, reporter.skipped)
			require.Equal(t, tt.wantReason, reporter.reason)
		})
	}
}

// Record the gate decision without deliberately skipping a Go subtest.
type circuitGateRecorder struct {
	skipped bool
	reason  string
}

func (*circuitGateRecorder) Helper() {}
func (r *circuitGateRecorder) Skip(args ...any) {
	r.skipped = true
	r.reason = fmt.Sprint(args...)
}
