package modbusshare

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReconciler_IdempotentDoubleSubmit(t *testing.T) {
	ctx := context.Background()
	reconciler, svc, tagSvc, _ := setupReconcilerFixture(t)

	t1, _ := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "t1", DisplayName: "Tag 1", DataType: schema.DataTypeInt16})

	req := ReconcileRequest{
		WorkspaceID:               "ws-1",
		ExpectedWorkspaceRevision: "rev-1",
		ExpectedSettingsRevision:  "set-1",
		DesiredMappings: []DesiredMapping{
			{
				WorkspaceID:        "ws-1",
				TagID:              t1.ID,
				DataType:           schema.DataTypeInt16,
				ZeroBasedRegister:  0,
				ShareStartRegister: 40001,
				SpanRegisters:      1,
				StrideRegisters:    1,
			},
		},
	}

	out1, err := reconciler.Reconcile(ctx, req)
	require.NoError(t, err)
	require.Equal(t, "applied", out1.Outcome)

	// Resubmitting the same desired set with the new revision is idempotent (aligned)
	req2 := ReconcileRequest{
		WorkspaceID:               "ws-1",
		ExpectedWorkspaceRevision: out1.NewWorkspaceRevision,
		ExpectedSettingsRevision:  "set-1",
		DesiredMappings:           req.DesiredMappings,
	}

	out2, err := reconciler.Reconcile(ctx, req2)
	require.NoError(t, err)
	assert.Equal(t, "aligned", out2.Outcome)
	assert.Equal(t, out1.NewWorkspaceRevision, out2.NewWorkspaceRevision)
	assert.True(t, svc.HasMapping(t1.ID))
}

func TestReconciler_IdempotentSameExpectedRevisionResubmit(t *testing.T) {
	ctx := context.Background()
	reconciler, svc, tagSvc, _ := setupReconcilerFixture(t)
	t1, err := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "same-rev", DisplayName: "Same revision", DataType: schema.DataTypeInt16})
	require.NoError(t, err)
	req := ReconcileRequest{WorkspaceID: "ws-1", ExpectedWorkspaceRevision: "rev-1", ExpectedSettingsRevision: "set-1", DesiredMappings: []DesiredMapping{{WorkspaceID: "ws-1", TagID: t1.ID, DataType: schema.DataTypeInt16, ZeroBasedRegister: 0, ShareStartRegister: 40001, SpanRegisters: 1, StrideRegisters: 1}}}
	first, err := reconciler.Reconcile(ctx, req)
	require.NoError(t, err)
	second, err := reconciler.Reconcile(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, first.Outcome, second.Outcome)
	assert.Equal(t, first.NewWorkspaceRevision, second.NewWorkspaceRevision)
	assert.True(t, svc.HasMapping(t1.ID))
}

func TestReconciler_ExactSuccessfulRequestReplaysAfterReadinessTokenRotation(t *testing.T) {
	ctx := context.Background()
	reconciler, _, tagSvc, _ := setupReconcilerFixture(t)
	t1, err := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "stale-token", DisplayName: "Stale token", DataType: schema.DataTypeInt16})
	require.NoError(t, err)
	reconciler.svc.SetHydrationState(HydrationState{State: HydrationStateReady, WorkspaceID: "ws-1", WorkspaceRevision: "rev-1", SettingsRevision: "set-1", Readiness: true, ReadinessToken: "token-1"})
	req := ReconcileRequest{WorkspaceID: "ws-1", ExpectedWorkspaceRevision: "rev-1", ExpectedSettingsRevision: "set-1", ReadinessToken: "token-1", DesiredMappings: []DesiredMapping{{WorkspaceID: "ws-1", TagID: t1.ID, DataType: schema.DataTypeInt16, ZeroBasedRegister: 0, ShareStartRegister: 40001, SpanRegisters: 1, StrideRegisters: 1}}}
	first, err := reconciler.Reconcile(ctx, req)
	require.NoError(t, err)
	assert.NotEmpty(t, first.NewReadinessToken)
	second, err := reconciler.Reconcile(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, first, second)
}

func TestReconciler_SameOldTokenDifferentDesiredSetConflicts(t *testing.T) {
	ctx := context.Background()
	reconciler, _, tagSvc, _ := setupReconcilerFixture(t)
	t1, err := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "old-token-a", DataType: schema.DataTypeInt16})
	require.NoError(t, err)
	t2, err := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "old-token-b", DataType: schema.DataTypeInt16})
	require.NoError(t, err)
	reconciler.svc.SetHydrationState(HydrationState{State: HydrationStateReady, WorkspaceID: "ws-1", WorkspaceRevision: "rev-1", SettingsRevision: "set-1", Readiness: true, ReadinessToken: "token-1"})
	base := DesiredMapping{WorkspaceID: "ws-1", SourceRuleID: "rule-1", SourceRuleRevision: "rev-1", TagID: t1.ID, DataType: schema.DataTypeInt16, ZeroBasedRegister: 0, ShareStartRegister: 40001, SpanRegisters: 1, StrideRegisters: 1}
	req := ReconcileRequest{WorkspaceID: "ws-1", ExpectedWorkspaceRevision: "rev-1", ExpectedSettingsRevision: "set-1", ReadinessToken: "token-1", DesiredMappings: []DesiredMapping{base}}
	first, err := reconciler.Reconcile(ctx, req)
	require.NoError(t, err)
	base.TagID = t2.ID
	base.ZeroBasedRegister = 1
	base.ShareStartRegister = 40002
	req.DesiredMappings = []DesiredMapping{base}
	req.ExpectedWorkspaceRevision = first.NewWorkspaceRevision
	_, err = reconciler.Reconcile(ctx, req)
	var shareErr *Error
	require.ErrorAs(t, err, &shareErr)
	assert.Equal(t, ErrCodeHydrationRequired, shareErr.Code)
}

func TestReconciler_ExactReplayRejectsExternalWorkspaceRevisionDrift(t *testing.T) {
	ctx := context.Background()
	reconciler, _, tagSvc, revStore := setupReconcilerFixture(t)
	t1, err := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "external-drift", DataType: schema.DataTypeInt16})
	require.NoError(t, err)
	reconciler.svc.SetHydrationState(HydrationState{State: HydrationStateReady, WorkspaceID: "ws-1", WorkspaceRevision: "rev-1", SettingsRevision: "set-1", Readiness: true, ReadinessToken: "token-1"})
	req := ReconcileRequest{WorkspaceID: "ws-1", ExpectedWorkspaceRevision: "rev-1", ExpectedSettingsRevision: "set-1", ReadinessToken: "token-1", DesiredMappings: []DesiredMapping{{WorkspaceID: "ws-1", TagID: t1.ID, DataType: schema.DataTypeInt16, ZeroBasedRegister: 0, ShareStartRegister: 40001, SpanRegisters: 1, StrideRegisters: 1}}}
	first, err := reconciler.Reconcile(ctx, req)
	require.NoError(t, err)
	revStore.revisions["ws-1"] = "external-revision"
	_, err = reconciler.Reconcile(ctx, req)
	var shareErr *Error
	require.ErrorAs(t, err, &shareErr)
	assert.Equal(t, ErrCodeRevisionConflict, shareErr.Code)
	assert.NotEqual(t, first.NewWorkspaceRevision, "external-revision")
}

func TestReconciler_ExactReplayDoesNotCrossProcessBoundary(t *testing.T) {
	ctx := context.Background()
	reconciler, svc, tagSvc, revStore := setupReconcilerFixture(t)
	t1, err := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "process-restart", DataType: schema.DataTypeInt16})
	require.NoError(t, err)
	reconciler.svc.SetHydrationState(HydrationState{State: HydrationStateReady, WorkspaceID: "ws-1", WorkspaceRevision: "rev-1", SettingsRevision: "set-1", Readiness: true, ReadinessToken: "token-1"})
	req := ReconcileRequest{WorkspaceID: "ws-1", ExpectedWorkspaceRevision: "rev-1", ExpectedSettingsRevision: "set-1", ReadinessToken: "token-1", DesiredMappings: []DesiredMapping{{WorkspaceID: "ws-1", TagID: t1.ID, DataType: schema.DataTypeInt16, ZeroBasedRegister: 0, ShareStartRegister: 40001, SpanRegisters: 1, StrideRegisters: 1}}}
	_, err = reconciler.Reconcile(ctx, req)
	require.NoError(t, err)
	restarted := NewReconciler(svc, revStore).WithOwnershipValidator(func(context.Context, DesiredMapping) error {
		return nil
	})
	restarted.svc.SetHydrationState(HydrationState{State: HydrationStateReady, WorkspaceID: "ws-1", WorkspaceRevision: revStore.revisions["ws-1"], SettingsRevision: "set-1", Readiness: true, ReadinessToken: "token-2"})
	_, err = restarted.Reconcile(ctx, req)
	var shareErr *Error
	require.ErrorAs(t, err, &shareErr)
	assert.Equal(t, ErrCodeRevisionConflict, shareErr.Code)
}
