package modbusshare

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReconciler_NilRevisionStoreFailsClosedBeforeMutation(t *testing.T) {
	ctx := context.Background()
	svc := NewService(nil, 4096)
	svc.SetHydrationState(HydrationState{State: HydrationStateReady, WorkspaceID: "ws-nil-store", Readiness: true})
	reconciler := NewReconciler(svc, nil).WithOwnershipValidator(func(context.Context, DesiredMapping) error {
		return nil
	})

	outcome, err := reconciler.Reconcile(ctx, ReconcileRequest{
		WorkspaceID: "ws-nil-store",
		DesiredMappings: []DesiredMapping{{
			WorkspaceID: "ws-nil-store", TagID: "tag-1", DataType: "int16",
			ZeroBasedRegister: 0, ShareStartRegister: 40001, SpanRegisters: 1, StrideRegisters: 1,
		}},
	})
	require.Error(t, err)
	var shareErr *Error
	require.ErrorAs(t, err, &shareErr)
	require.Equal(t, ErrCodeProjectionRequired, shareErr.Code)
	require.Equal(t, HydrationStateFailed, outcome.Outcome)
	require.Empty(t, svc.ListMappings())
}

func TestReconciler_EmptyExpectedRevisionsRejectBeforeMutation(t *testing.T) {
	ctx := context.Background()
	svc := NewService(nil, 4096)
	svc.SetHydrationState(HydrationState{State: HydrationStateReady, WorkspaceID: "ws-empty", Readiness: true})
	reconciler := NewReconciler(svc, newMockWorkspaceRevisionStore()).WithOwnershipValidator(func(context.Context, DesiredMapping) error { return nil })
	before := svc.MemorySnapshot()
	_, err := reconciler.Reconcile(ctx, ReconcileRequest{WorkspaceID: "ws-empty", DesiredMappings: []DesiredMapping{{WorkspaceID: "ws-empty", TagID: "tag", DataType: "int16", ZeroBasedRegister: 0, SpanRegisters: 1, StrideRegisters: 1}}})
	var shareErr *Error
	require.ErrorAs(t, err, &shareErr)
	require.Equal(t, ErrCodeRevisionRequired, shareErr.Code)
	require.Equal(t, before, svc.MemorySnapshot())
	require.Empty(t, svc.ListMappings())
}
