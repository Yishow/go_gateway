package sourcerule

import (
	"testing"

	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestMCAddressBounds_AcceptsMaximumClientOffset(t *testing.T) {
	tests := []struct {
		name     string
		start    string
		expected string
	}{
		{name: "hex word", start: "WFFFFFF", expected: "WFFFFFF"},
		{name: "octal contact", start: "X77777777", expected: "X77777777"},
		{name: "decimal register", start: "D16777215", expected: "D16777215"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildPlannedPointAddresses(tt.start, 1, schema.DataTypeInt16, schema.ProtocolMC3E)
			require.NoError(t, err)
			require.Equal(t, []string{tt.expected}, got)
		})
	}
}

func TestMCAddressBounds_RejectsInvalidAddressOrExpansion(t *testing.T) {
	tests := []struct {
		name  string
		start string
		count int
	}{
		{name: "hex address over client maximum", start: "W1000000", count: 1},
		{name: "hex address far over client maximum", start: "WFFFFFFFF", count: 1},
		{name: "large hex address cannot wrap after offset", start: "W7FFFFFFFFFFFFFFF", count: 1},
		{name: "decimal max int64 input", start: "D9223372036854775807", count: 1},
		{name: "expansion crosses client maximum", start: "WFFFFFE", count: 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildPlannedPointAddresses(tt.start, tt.count, schema.DataTypeInt16, schema.ProtocolMC3E)
			require.ErrorIs(t, err, ErrValidation)
			require.Empty(t, got)
		})
	}
}

func TestMCAddressBounds_RejectsOffsetAcrossMaximum(t *testing.T) {
	got, err := offsetAddress("WFFFFFE", 1, schema.ProtocolMC3E)
	require.NoError(t, err)
	require.Equal(t, "WFFFFFF", got)

	got, err = offsetAddress("WFFFFFF", 1, schema.ProtocolMC3E)
	require.ErrorIs(t, err, ErrValidation)
	require.Empty(t, got)
}

func TestMCAddressBounds_RejectsCreateBeforePersistence(t *testing.T) {
	for _, start := range []string{"W1000000", "WFFFFFFFF", "W7FFFFFFFFFFFFFFF", "D9223372036854775807", "WFFFFFE"} {
		t.Run(start, func(t *testing.T) {
			svc, repo, pointSvc, dev := newMCPlannerService(t)
			spy := &mcGuardRepositorySpy{MemoryRepository: repo}
			svc.repo = spy
			count := 1
			if start == "WFFFFFE" {
				count = 3
			}

			_, err := svc.Create(t.Context(), CreateRuleRequest{
				ID:           "rule-mc-bounds-" + start,
				DeviceID:     dev.ID,
				StartAddress: start,
				Count:        count,
				DataType:     schema.DataTypeInt16,
				NamingPrefix: "MC",
				Enabled:      true,
			})
			require.ErrorIs(t, err, ErrValidation)
			require.Empty(t, spy.calls)

			rules, listErr := repo.List(t.Context(), ListFilter{DeviceID: &dev.ID})
			require.NoError(t, listErr)
			require.Empty(t, rules)
			points, pointErr := pointSvc.List(t.Context(), point.ListFilter{DeviceID: &dev.ID})
			require.NoError(t, pointErr)
			require.Empty(t, points)
		})
	}
}
