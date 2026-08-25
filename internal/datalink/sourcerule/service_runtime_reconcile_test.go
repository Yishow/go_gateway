package sourcerule

import (
	"context"
	"errors"
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/modbusshare"
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

func TestService_GenericMutationsHandOffShareAuthority(t *testing.T) {
	fixture := newRuntimeReconcileRuleFixture(t, true)
	recorder := &shareReconcileRecorder{}
	fixture.svc.SetShareRuntimeReconciler(recorder)

	_, err := fixture.svc.DisableWithRuntimeReconcile(fixture.ctx, fixture.rule.ID)
	require.NoError(t, err)
	_, err = fixture.svc.EnableWithRuntimeReconcile(fixture.ctx, fixture.rule.ID)
	require.NoError(t, err)
	prefix := "SHARE-UPDATED"
	_, _, err = fixture.svc.UpdateWithRuntimeReconcile(fixture.ctx, fixture.rule.ID, UpdateRuleRequest{NamingPrefix: &prefix})
	require.NoError(t, err)
	created, _, err := fixture.svc.CreateWithRuntimeReconcile(fixture.ctx, CreateRuleRequest{
		ID: "rule-share-created", DeviceID: fixture.deviceID, StartAddress: "40002", Count: 1,
		DataType: schema.DataTypeInt16, NamingPrefix: "SHARE", Enabled: true,
	})
	require.NoError(t, err)
	_, err = fixture.svc.DeleteWithRuntimeReconcile(fixture.ctx, created.ID)
	require.NoError(t, err)

	require.ElementsMatch(t, []RuntimeReconcileOperation{
		RuntimeReconcileOperationDisable,
		RuntimeReconcileOperationEnable,
		RuntimeReconcileOperationUpdate,
		RuntimeReconcileOperationCreate,
		RuntimeReconcileOperationDelete,
	}, recorder.operations)
}

func TestService_GenericRuntimeFailureIsNotMaskedByAlignedShare(t *testing.T) {
	fixture := newRuntimeReconcileRuleFixture(t, true)
	fixture.runtimeSync.reconcileOutcomes = []RuntimeReconcileStatus{RuntimeReconcileStatusStale, RuntimeReconcileStatusAligned}
	share := &shareReconcileRecorder{}
	fixture.svc.SetShareRuntimeReconciler(share)

	_, err := fixture.svc.DisableWithRuntimeReconcile(fixture.ctx, fixture.rule.ID)
	var reconcileErr *modbusshare.Error
	require.ErrorAs(t, err, &reconcileErr)
	require.Equal(t, modbusshare.ErrCodeReconcileFailed, reconcileErr.Code)
	require.Empty(t, share.operations)
}

func TestService_UpdateWithRuntimeReconcileRestoresRuleOnShareFailure(t *testing.T) {
	fixture := newRuntimeReconcileRuleFixture(t, true)
	fixture.svc.SetShareRuntimeReconciler(failingShareReconciler{})
	prefix := "SHARE-FAIL"

	_, _, err := fixture.svc.UpdateWithRuntimeReconcile(fixture.ctx, fixture.rule.ID, UpdateRuleRequest{NamingPrefix: &prefix})
	requireShareReconcileFailure(t, err)
	restored, getErr := fixture.svc.GetByID(fixture.ctx, fixture.rule.ID)
	require.NoError(t, getErr)
	require.Equal(t, fixture.rule.NamingPrefix, restored.NamingPrefix)
}

func TestService_EnableWithRuntimeReconcileRestoresRuleOnShareFailure(t *testing.T) {
	fixture := newRuntimeReconcileRuleFixture(t, false)
	fixture.svc.SetShareRuntimeReconciler(failingShareReconciler{})

	_, err := fixture.svc.EnableWithRuntimeReconcile(fixture.ctx, fixture.rule.ID)
	requireShareReconcileFailure(t, err)
	restored, getErr := fixture.svc.GetByID(fixture.ctx, fixture.rule.ID)
	require.NoError(t, getErr)
	require.False(t, restored.Enabled)
}

func TestService_DeleteWithRuntimeReconcileRestoresRuleOnShareFailure(t *testing.T) {
	fixture := newRuntimeReconcileRuleFixture(t, true)
	fixture.svc.SetShareRuntimeReconciler(failingShareReconciler{})

	_, err := fixture.svc.DeleteWithRuntimeReconcile(fixture.ctx, fixture.rule.ID)
	requireShareReconcileFailure(t, err)
	_, getErr := fixture.svc.GetByID(fixture.ctx, fixture.rule.ID)
	require.NoError(t, getErr)
}

func TestService_CreateWithRuntimeReconcileRemovesRuleOnShareFailure(t *testing.T) {
	fixture := newRuntimeReconcileRuleFixture(t, true)
	fixture.svc.SetShareRuntimeReconciler(failingShareReconciler{})

	created, _, err := fixture.svc.CreateWithRuntimeReconcile(fixture.ctx, CreateRuleRequest{
		ID: "rule-share-create-fail", DeviceID: fixture.deviceID, StartAddress: "40002", Count: 1,
		DataType: schema.DataTypeInt16, NamingPrefix: "FAIL", Enabled: true,
	})
	requireShareReconcileFailure(t, err)
	require.NotNil(t, created)
	_, getErr := fixture.svc.GetByID(fixture.ctx, created.ID)
	require.Error(t, getErr)
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
	reconcileOutcomes []RuntimeReconcileStatus
	reconcileRequests []RuntimeReconcileRequest
}

type shareReconcileRecorder struct {
	operations []RuntimeReconcileOperation
}

var errShareReconcile = errors.New("share reconcile failed")

func requireShareReconcileFailure(t *testing.T, err error) {
	t.Helper()
	var shareErr *modbusshare.Error
	require.ErrorAs(t, err, &shareErr)
	require.Equal(t, modbusshare.ErrCodeReconcileFailed, shareErr.Code)
}

type failingShareReconciler struct{}

func (failingShareReconciler) ReconcileSourceRuleShare(context.Context, RuntimeReconcileRequest) (RuntimeReconcileOutcome, error) {
	return RuntimeReconcileOutcome{}, errShareReconcile
}

func (r *shareReconcileRecorder) ReconcileSourceRuleShare(_ context.Context, req RuntimeReconcileRequest) (RuntimeReconcileOutcome, error) {
	r.operations = append(r.operations, req.Operation)
	return RuntimeReconcileOutcome{Status: RuntimeReconcileStatusAligned, Scope: req.Scope}, nil
}

func (s *reconcileRuntimeSync) ReconcileSourceRule(_ context.Context, req RuntimeReconcileRequest) RuntimeReconcileOutcome {
	s.reconcileRequests = append(s.reconcileRequests, req)
	if len(s.reconcileOutcomes) > 0 {
		status := s.reconcileOutcomes[0]
		s.reconcileOutcomes = s.reconcileOutcomes[1:]
		return RuntimeReconcileOutcome{Status: status, Scope: req.Scope}
	}
	if s.reconcileOutcome.Status == "" {
		s.reconcileOutcome = RuntimeReconcileOutcome{Status: RuntimeReconcileStatusAligned}
	}
	s.reconcileOutcome.Scope = req.Scope
	return s.reconcileOutcome
}
