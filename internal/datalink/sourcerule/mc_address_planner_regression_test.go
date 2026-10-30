package sourcerule

import (
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/lib/hsllogic"

	"github.com/stretchr/testify/require"
)

func newMCPlannerService(t *testing.T) (service *Service, rules *MemoryRepository, points *point.Service, deviceRecord *schema.Device) {
	t.Helper()
	ctx := t.Context()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, nil)
	pointSvc := point.NewService(pointRepo, nil)
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, pointSvc, nil)

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-mc-planner")
	require.NoError(t, err)
	dev.Protocol = schema.ProtocolMC3E
	require.NoError(t, deviceRepo.Update(ctx, dev))
	return svc, repo, pointSvc, dev
}

func TestOffsetAddress_MCUsesCorrectedRadixAndRuntimeOffsets(t *testing.T) {
	tests := []struct {
		name       string
		start      string
		delta      int
		expected   string
		runtimeOff int
	}{
		{name: "W offset ten", start: "W0", delta: 10, expected: "WA", runtimeOff: 10},
		{name: "W offset sixteen", start: "W0", delta: 16, expected: "W10", runtimeOff: 16},
		{name: "W carry", start: "WA", delta: 1, expected: "WB", runtimeOff: 11},
		{name: "B carry", start: "BFF", delta: 1, expected: "B100", runtimeOff: 256},
		{name: "X offset sixteen", start: "X0", delta: 16, expected: "X20", runtimeOff: 16},
		{name: "X carry", start: "X7", delta: 1, expected: "X10", runtimeOff: 8},
		{name: "Y offset sixteen", start: "Y0", delta: 16, expected: "Y20", runtimeOff: 16},
		{name: "D decimal", start: "D0", delta: 16, expected: "D16", runtimeOff: 16},
		{name: "M decimal", start: "M0", delta: 16, expected: "M16", runtimeOff: 16},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := offsetAddress(tt.start, tt.delta, schema.ProtocolMC3E)
			require.NoError(t, err)
			require.Equal(t, tt.expected, got)
			parsed, err := hsllogic.ParseAddress(hsllogic.ProtocolMitsubishi, got)
			require.NoError(t, err)
			require.Equal(t, tt.runtimeOff, parsed.Offset)
		})
	}
}

func TestOffsetAddress_NonMCProtocolRadixBoundariesRemainStable(t *testing.T) {
	tests := []struct {
		protocol schema.ProtocolType
		start    string
		delta    int
		expected string
	}{
		{protocol: schema.ProtocolFatekFBs, start: "R0", delta: 16, expected: "R16"},
		{protocol: schema.ProtocolModbusTCP, start: "40001", delta: 2, expected: "40003"},
	}
	for _, tt := range tests {
		got, err := offsetAddress(tt.start, tt.delta, tt.protocol)
		require.NoError(t, err)
		require.Equal(t, tt.expected, got)
	}
}

func TestBuildPlannedPointAddresses_MCW0Count17UsesCompleteCorrectedRange(t *testing.T) {
	addresses, err := buildPlannedPointAddresses("W0", 17, schema.DataTypeInt16, schema.ProtocolMC3E)
	require.NoError(t, err)
	require.Equal(t, []string{"W0", "W1", "W2", "W3", "W4", "W5", "W6", "W7", "W8", "W9", "WA", "WB", "WC", "WD", "WE", "WF", "W10"}, addresses)
}

func TestServiceCreate_MCW0Count17PersistsCompleteCorrectedRange(t *testing.T) {
	svc, repo, pointSvc, dev := newMCPlannerService(t)
	rule, err := svc.Create(t.Context(), CreateRuleRequest{
		ID: "rule-mc-planner", DeviceID: dev.ID, StartAddress: "W0", Count: 17,
		DataType: schema.DataTypeInt16, NamingPrefix: "MC", Enabled: true,
	})
	require.NoError(t, err)

	links, err := repo.ListLinks(t.Context(), rule.ID)
	require.NoError(t, err)
	addresses := make([]string, 0, len(links))
	for _, link := range links {
		addresses = append(addresses, link.Address)
	}
	require.ElementsMatch(t, []string{"W0", "W1", "W2", "W3", "W4", "W5", "W6", "W7", "W8", "W9", "WA", "WB", "WC", "WD", "WE", "WF", "W10"}, addresses)

	points, err := pointSvc.List(t.Context(), point.ListFilter{DeviceID: &dev.ID})
	require.NoError(t, err)
	require.Len(t, points, 17)
}

func TestServiceCreate_MCInvalidCorrectedAddressesFailBeforePersistence(t *testing.T) {
	for _, start := range []string{"X8", "X9", "YA", "WG", "Z0"} {
		t.Run(start, func(t *testing.T) {
			svc, repo, pointSvc, dev := newMCPlannerService(t)
			_, err := svc.Create(t.Context(), CreateRuleRequest{
				ID: "rule-mc-invalid-" + start, DeviceID: dev.ID, StartAddress: start, Count: 1,
				DataType: schema.DataTypeInt16, NamingPrefix: "MC", Enabled: true,
			})
			require.ErrorIs(t, err, ErrValidation)
			require.ErrorContains(t, err, "rule_id=rule-mc-invalid-"+start)
			require.ErrorContains(t, err, "start_address="+start)

			rules, listErr := repo.List(t.Context(), ListFilter{DeviceID: &dev.ID})
			require.NoError(t, listErr)
			require.Empty(t, rules)
			points, pointErr := pointSvc.List(t.Context(), point.ListFilter{DeviceID: &dev.ID})
			require.NoError(t, pointErr)
			require.Empty(t, points)
		})
	}
}

func TestServiceUpdate_MCInvalidCorrectedAddressesFailBeforeMutation(t *testing.T) {
	for _, start := range []string{"X8", "X9", "YA", "WG", "Z0"} {
		t.Run(start, func(t *testing.T) {
			svc, repo, pointSvc, dev := newMCPlannerService(t)
			rule, err := svc.Create(t.Context(), CreateRuleRequest{
				ID: "rule-mc-update-" + start, DeviceID: dev.ID, StartAddress: "D0", Count: 1,
				DataType: schema.DataTypeInt16, NamingPrefix: "MC", Enabled: true,
			})
			require.NoError(t, err)
			beforeLinks, err := repo.ListLinks(t.Context(), rule.ID)
			require.NoError(t, err)

			_, err = svc.Update(t.Context(), rule.ID, UpdateRuleRequest{StartAddress: &start})
			require.ErrorIs(t, err, ErrValidation)
			require.ErrorContains(t, err, "rule_id=rule-mc-update-"+start)
			require.ErrorContains(t, err, "start_address="+start)
			stored, getErr := repo.GetByID(t.Context(), rule.ID)
			require.NoError(t, getErr)
			require.Equal(t, "D0", stored.StartAddress)
			afterLinks, listErr := repo.ListLinks(t.Context(), rule.ID)
			require.NoError(t, listErr)
			require.Equal(t, beforeLinks, afterLinks)
			points, pointErr := pointSvc.List(t.Context(), point.ListFilter{DeviceID: &dev.ID})
			require.NoError(t, pointErr)
			require.Len(t, points, 1)
		})
	}
}
