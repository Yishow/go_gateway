package modbusshare

import (
	"context"
	"errors"
	"testing"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReconciler_AtomicCommitSuccess(t *testing.T) {
	ctx := context.Background()
	reconciler, svc, tagSvc, _ := setupReconcilerFixture(t)

	t1, err := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "temp", DisplayName: "Temp", DataType: schema.DataTypeInt32})
	require.NoError(t, err)

	req := ReconcileRequest{
		WorkspaceID:               "ws-1",
		ExpectedWorkspaceRevision: "rev-1",
		ExpectedSettingsRevision:  "set-1",
		DesiredMappings: []DesiredMapping{
			{
				WorkspaceID:        "ws-1",
				SourceRuleID:       "rule-1",
				SourceRuleRevision: "rrev-1",
				TagID:              t1.ID,
				DataType:           schema.DataTypeInt32,
				ShareStartRegister: 40001,
				ZeroBasedRegister:  0,
				SpanRegisters:      2,
				StrideRegisters:    2,
			},
		},
	}

	outcome, err := reconciler.Reconcile(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, "applied", outcome.Outcome)
	assert.Equal(t, 1, outcome.AppliedCount)
	assert.NotEmpty(t, outcome.NewWorkspaceRevision)
	assert.NotEqual(t, "rev-1", outcome.NewWorkspaceRevision)
	assert.True(t, svc.HasMapping(t1.ID))
}

func TestReconciler_StaleRevisionCASConflict(t *testing.T) {
	ctx := context.Background()
	reconciler, _, tagSvc, _ := setupReconcilerFixture(t)

	t1, _ := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "t1", DisplayName: "Tag 1", DataType: schema.DataTypeInt16})

	req := ReconcileRequest{
		WorkspaceID:               "ws-1",
		ExpectedWorkspaceRevision: "stale-rev-999", // mismatch
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

	_, err := reconciler.Reconcile(ctx, req)
	require.Error(t, err)
	var shareErr *Error
	require.ErrorAs(t, err, &shareErr)
	assert.Equal(t, ErrCodeRevisionConflict, shareErr.Code)
}

func TestReconciler_RollbackUncertaintyReportsDirtyUnknown(t *testing.T) {
	ctx := context.Background()
	reconciler, svc, tagSvc, revStore := setupReconcilerFixture(t)
	t1, err := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "dirty", DisplayName: "Dirty", DataType: schema.DataTypeInt16})
	require.NoError(t, err)
	first, err := reconciler.Reconcile(ctx, ReconcileRequest{WorkspaceID: "ws-1", ExpectedWorkspaceRevision: "rev-1", ExpectedSettingsRevision: "set-1", DesiredMappings: []DesiredMapping{{WorkspaceID: "ws-1", SourceRuleID: "rule-1", SourceRuleRevision: "rev-1", TagID: t1.ID, DataType: schema.DataTypeInt16, ZeroBasedRegister: 0, ShareStartRegister: 40001, SpanRegisters: 1, StrideRegisters: 1}}})
	require.NoError(t, err)
	svc.SetProjectionSwapHook(func(map[string]TagMirrorMapping) error { return errors.New("runtime projection failed") })
	revStore.failCommitAt = 3 // forward commit succeeds; durable rollback fails
	revStore.updateErr = errors.New("durable commit failed")
	svc.SetRestoreMemoryHook(func([]byte) error { return errors.New("memory rollback failed") })
	out, err := reconciler.Reconcile(ctx, ReconcileRequest{WorkspaceID: "ws-1", ExpectedWorkspaceRevision: first.NewWorkspaceRevision, ExpectedSettingsRevision: "set-1", DesiredMappings: []DesiredMapping{}})
	require.Error(t, err)
	assert.Equal(t, "dirty_unknown", out.Outcome)
	var shareErr *Error
	require.ErrorAs(t, err, &shareErr)
	assert.Equal(t, ErrCodeDirtyUnknown, shareErr.Code)
}

func TestReconciler_RollbackUncertaintyFailsClosedWhenMarkDirtyFails(t *testing.T) {
	ctx := context.Background()
	reconciler, svc, tagSvc, revStore := setupReconcilerFixture(t)
	tagRecord, err := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "dirty-mark-failure", DataType: schema.DataTypeInt16})
	require.NoError(t, err)
	first, err := reconciler.Reconcile(ctx, ReconcileRequest{WorkspaceID: "ws-1", ExpectedWorkspaceRevision: "rev-1", ExpectedSettingsRevision: "set-1", DesiredMappings: []DesiredMapping{{WorkspaceID: "ws-1", SourceRuleID: "rule-1", SourceRuleRevision: "r1", TagID: tagRecord.ID, DataType: schema.DataTypeInt16, ZeroBasedRegister: 0, ShareStartRegister: 40001, SpanRegisters: 1, StrideRegisters: 1}}})
	require.NoError(t, err)
	svc.SetProjectionSwapHook(func(map[string]TagMirrorMapping) error { return errors.New("projection failed") })
	svc.SetRestoreMemoryHook(func([]byte) error { return errors.New("memory restore failed") })
	revStore.failCommitAt = 3
	revStore.updateErr = errors.New("durable rollback failed")
	revStore.markDirtyErr = errors.New("mark dirty failed")
	_, err = reconciler.Reconcile(ctx, ReconcileRequest{WorkspaceID: "ws-1", ExpectedWorkspaceRevision: first.NewWorkspaceRevision, ExpectedSettingsRevision: "set-1", DesiredMappings: []DesiredMapping{}})
	var shareErr *Error
	require.ErrorAs(t, err, &shareErr)
	require.Equal(t, ErrCodeDirtyUnknown, shareErr.Code)
	hydration, hydrationErr := svc.CheckHydration(ctx)
	require.NoError(t, hydrationErr)
	require.Equal(t, HydrationStateFailed, hydration.State)
	require.False(t, hydration.Readiness)
}
