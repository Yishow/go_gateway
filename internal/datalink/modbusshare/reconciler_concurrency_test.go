package modbusshare

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/require"
)

func TestReconciler_FailedWorkspaceRollbackDoesNotOverwriteConcurrentWorkspaceCommit(t *testing.T) {
	ctx := context.Background()
	reconciler, svc, tagSvc, _ := setupReconcilerFixture(t)
	tagA, err := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "workspace-a", DataType: schema.DataTypeInt16})
	require.NoError(t, err)
	tagB, err := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "workspace-b", DataType: schema.DataTypeInt16})
	require.NoError(t, err)

	aEntered := make(chan struct{})
	releaseA := make(chan struct{})
	bEntered := make(chan struct{})
	var aTagID, bTagID = tagA.ID, tagB.ID
	svc.SetProjectionSwapHook(func(staged map[string]TagMirrorMapping) error {
		if _, ok := staged[aTagID]; ok {
			close(aEntered)
			<-releaseA
			return errors.New("workspace A projection failure")
		}
		if _, ok := staged[bTagID]; ok {
			close(bEntered)
		}
		return nil
	})

	aDone := make(chan error, 1)
	go func() {
		_, reconcileErr := reconciler.Reconcile(ctx, ReconcileRequest{
			WorkspaceID: "ws-a", ExpectedWorkspaceRevision: "rev-1", ExpectedSettingsRevision: "set-1",
			DesiredMappings: []DesiredMapping{{WorkspaceID: "ws-a", TagID: aTagID,
				DataType: schema.DataTypeInt16, ZeroBasedRegister: 0,
				ShareStartRegister: 40001, SpanRegisters: 1, StrideRegisters: 1}},
		})
		aDone <- reconcileErr
	}()
	<-aEntered

	bDone := make(chan error, 1)
	go func() {
		_, reconcileErr := reconciler.Reconcile(ctx, ReconcileRequest{
			WorkspaceID: "ws-b", ExpectedWorkspaceRevision: "rev-1", ExpectedSettingsRevision: "set-1",
			DesiredMappings: []DesiredMapping{{WorkspaceID: "ws-b", TagID: bTagID,
				DataType: schema.DataTypeInt16, ZeroBasedRegister: 2,
				ShareStartRegister: 40003, SpanRegisters: 1, StrideRegisters: 1}},
		})
		if reconcileErr == nil {
			reconcileErr = svc.WriteProjectedTagValue(ctx, bTagID, int16(222))
		}
		bDone <- reconcileErr
	}()

	// A correct global projection transaction lock makes B wait here. The
	// existing per-workspace-only implementation reaches bEntered first and
	// completes B before A restores its stale global snapshot.
	bCompleted := false
	select {
	case <-bEntered:
		select {
		case err := <-bDone:
			require.NoError(t, err)
			bCompleted = true
		case <-time.After(time.Second):
			t.Fatal("workspace B did not finish before workspace A rollback")
		}
	case <-time.After(100 * time.Millisecond):
	}
	close(releaseA)

	require.Error(t, <-aDone)
	if !bCompleted {
		require.NoError(t, <-bDone)
	}
	require.False(t, svc.HasMapping(aTagID))
	require.True(t, svc.HasMapping(bTagID))
	words, err := svc.ReadHoldingWords(2, 1)
	require.NoError(t, err)
	require.Equal(t, []uint16{222}, words)
}
