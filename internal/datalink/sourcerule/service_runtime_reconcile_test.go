package sourcerule

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestService_DisableWithRuntimeReconcileReturnsAlignedRuleOutcome(t *testing.T) {
	fixture := newRuntimeReconcileRuleFixture(t, true)

	outcome, err := fixture.svc.DisableWithRuntimeReconcile(fixture.ctx, fixture.rule.ID)
	require.NoError(t, err)

	requireRuleReconcileOutcome(t, outcome, fixture.runtimeSync, RuntimeReconcileOperationDisable, fixture.rule.ID, fixture.deviceID)
}

func TestService_EnableWithRuntimeReconcileReturnsAlignedRuleOutcome(t *testing.T) {
	fixture := newRuntimeReconcileRuleFixture(t, false)

	outcome, err := fixture.svc.EnableWithRuntimeReconcile(fixture.ctx, fixture.rule.ID)
	require.NoError(t, err)

	requireRuleReconcileOutcome(t, outcome, fixture.runtimeSync, RuntimeReconcileOperationEnable, fixture.rule.ID, fixture.deviceID)
}

func TestService_UpdateWithRuntimeReconcileReturnsAlignedRuleOutcome(t *testing.T) {
	fixture := newRuntimeReconcileRuleFixture(t, true)
	prefix := "UPDATED"

	updatedRule, outcome, err := fixture.svc.UpdateWithRuntimeReconcile(fixture.ctx, fixture.rule.ID, UpdateRuleRequest{
		NamingPrefix: &prefix,
	})
	require.NoError(t, err)
	require.Equal(t, prefix, updatedRule.NamingPrefix)
	requireRuleReconcileOutcome(t, outcome, fixture.runtimeSync, RuntimeReconcileOperationUpdate, fixture.rule.ID, fixture.deviceID)
}

func TestService_DeleteWithRuntimeReconcileReturnsAlignedRuleOutcome(t *testing.T) {
	fixture := newRuntimeReconcileRuleFixture(t, true)

	outcome, err := fixture.svc.DeleteWithRuntimeReconcile(fixture.ctx, fixture.rule.ID)
	require.NoError(t, err)

	requireRuleReconcileOutcome(t, outcome, fixture.runtimeSync, RuntimeReconcileOperationDelete, fixture.rule.ID, fixture.deviceID)
}

type runtimeReconcileRuleFixture struct {
	ctx         context.Context
	svc         *Service
	rule        *schema.SourceRule
	runtimeSync *reconcileRuntimeSync
	deviceID    string
}

func newRuntimeReconcileRuleFixture(t *testing.T, enabled bool) runtimeReconcileRuleFixture {
	t.Helper()
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	runtimeSync := &reconcileRuntimeSync{}
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, pointSvc, runtimeSync)

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-1")
	require.NoError(t, err)
	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-A",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "RULE",
		Enabled:      enabled,
	})
	require.NoError(t, err)
	runtimeSync.reconcileRequests = nil

	return runtimeReconcileRuleFixture{
		ctx:         ctx,
		svc:         svc,
		rule:        rule,
		runtimeSync: runtimeSync,
		deviceID:    dev.ID,
	}
}

func requireRuleReconcileOutcome(
	t *testing.T,
	outcome RuntimeReconcileOutcome,
	runtimeSync *reconcileRuntimeSync,
	operation RuntimeReconcileOperation,
	ruleID string,
	deviceID string,
) {
	t.Helper()
	require.Equal(t, RuntimeReconcileStatusAligned, outcome.Status)
	require.Equal(t, ruleID, outcome.Scope.RuleID)
	require.Equal(t, deviceID, outcome.Scope.DeviceID)
	require.Len(t, runtimeSync.reconcileRequests, 1)
	require.Equal(t, operation, runtimeSync.reconcileRequests[0].Operation)
}

type reconcileRuntimeSync struct {
	stubRuntimeSync
	reconcileOutcome  RuntimeReconcileOutcome
	reconcileRequests []RuntimeReconcileRequest
}

func (s *reconcileRuntimeSync) ReconcileSourceRule(_ context.Context, req RuntimeReconcileRequest) RuntimeReconcileOutcome {
	s.reconcileRequests = append(s.reconcileRequests, req)
	if s.reconcileOutcome.Status == "" {
		s.reconcileOutcome = RuntimeReconcileOutcome{Status: RuntimeReconcileStatusAligned}
	}
	s.reconcileOutcome.Scope = req.Scope
	return s.reconcileOutcome
}
