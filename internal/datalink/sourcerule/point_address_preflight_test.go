package sourcerule

import (
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestPointAddressPreflightPreservesOwnOtherDeviceAndSkippedAddresses(t *testing.T) {
	ctx := t.Context()
	devices := device.NewMemoryRepository()
	_, err := seedActiveDevice(ctx, devices, "device-a")
	require.NoError(t, err)
	_, err = seedActiveDevice(ctx, devices, "device-b")
	require.NoError(t, err)
	svc := NewService(NewMemoryRepository(), device.NewService(devices, nil), point.NewService(point.NewMemoryRepository(), nil), nil)
	create := func(id, deviceID string, skipped []string) *schema.SourceRule {
		start := "40001"
		if deviceID == "device-b" {
			start = " 40001 "
		}
		rule, _, err := svc.CreateWithRuntimeReconcile(ctx, CreateRuleRequest{
			ID: id, DeviceID: deviceID, StartAddress: start, Count: 1,
			DataType: schema.DataTypeInt16, NamingPrefix: id, Enabled: true, SkippedAddresses: skipped,
		})
		require.NoError(t, err)
		require.Equal(t, "40001", rule.StartAddress)
		return rule
	}
	a := create("rule-a", "device-a", nil)
	create("rule-b", "device-b", nil)
	create("rule-skipped", "device-a", []string{"40001"})
	before, err := svc.ListLinks(ctx, a.ID)
	require.NoError(t, err)
	address := "40001"
	_, _, err = svc.UpdateWithRuntimeReconcile(ctx, a.ID, UpdateRuleRequest{StartAddress: &address})
	require.NoError(t, err)
	after, err := svc.ListLinks(ctx, a.ID)
	require.NoError(t, err)
	require.Len(t, after, 1)
	require.Equal(t, before[0].PointID, after[0].PointID)
}

func TestPointAddressPreflightRejectsForeignReplacementBeforeOrphanCleanup(t *testing.T) {
	ctx := t.Context()
	devices := device.NewMemoryRepository()
	_, err := seedActiveDevice(ctx, devices, "device-a")
	require.NoError(t, err)
	points := point.NewMemoryRepository()
	pointSvc := point.NewService(points, nil)
	rules := NewMemoryRepository()
	svc := NewService(rules, device.NewService(devices, nil), pointSvc, nil)
	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID: "rule-a", DeviceID: "device-a", StartAddress: "40001", Count: 2,
		DataType: schema.DataTypeInt16, NamingPrefix: "ORIGINAL", Enabled: true,
	})
	require.NoError(t, err)
	beforeLinks, err := svc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	require.NoError(t, points.Delete(ctx, beforeLinks[0].PointID))
	foreign, err := pointSvc.Create(ctx, point.CreatePointRequest{
		DeviceID: "device-a", Name: "Foreign replacement", Address: beforeLinks[0].Address, DataType: schema.DataTypeInt16,
	})
	require.NoError(t, err)
	spy := &mcGuardRepositorySpy{MemoryRepository: rules}
	svc.repo = spy
	prefix := "DO_NOT_SAVE"
	_, _, err = svc.UpdateWithRuntimeReconcile(ctx, rule.ID, UpdateRuleRequest{NamingPrefix: &prefix})
	require.ErrorIs(t, err, point.ErrPointNotFound, "the outer snapshot must still reject an already-broken source")
	require.Empty(t, spy.calls)
	// Direct Update historically repairs orphan links; reject a conflicting
	// address plan before that repair is allowed to delete any resources.
	_, err = svc.Update(ctx, rule.ID, UpdateRuleRequest{NamingPrefix: &prefix})
	require.ErrorIs(t, err, ErrValidation)
	var preflight *pointAddressPreflightError
	require.ErrorAs(t, err, &preflight)
	require.Empty(t, spy.calls, "must not clean or restore any source resource")
	afterLinks, err := svc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	require.Equal(t, beforeLinks, afterLinks)
	afterRule, err := svc.GetByID(ctx, rule.ID)
	require.NoError(t, err)
	require.Equal(t, rule, afterRule)
	afterForeign, err := pointSvc.GetByID(ctx, foreign.ID)
	require.NoError(t, err)
	require.Equal(t, foreign, afterForeign)
}
