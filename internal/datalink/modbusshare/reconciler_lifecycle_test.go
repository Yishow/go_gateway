package modbusshare

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReconciler_StagedTransactionIsRejectedAfterGlobalDisable(t *testing.T) {
	ctx := context.Background()
	reconciler, svc, tagSvc, revStore := setupReconcilerFixture(t)
	tagRecord, err := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "lifecycle", DataType: schema.DataTypeInt16})
	require.NoError(t, err)
	svc.mu.Lock()
	svc.settings.Enabled = true
	svc.settings.SettingsRevision = "set-1"
	svc.authoritativeSettings = true
	svc.mu.Unlock()

	stageReady := make(chan struct{})
	disableDone := make(chan struct{})
	revStore.beforeGetDesired = func() {
		close(stageReady)
		<-disableDone
	}
	go func() {
		<-stageReady
		disabled := svc.Settings()
		disabled.Enabled = false
		if disableErr := svc.ApplySettingsCAS(ctx, disabled, "set-1"); disableErr != nil {
			t.Errorf("disable settings: %v", disableErr)
		}
		close(disableDone)
	}()

	outcome, err := reconciler.Reconcile(ctx, ReconcileRequest{
		WorkspaceID: "ws-1", ExpectedWorkspaceRevision: "rev-1", ExpectedSettingsRevision: "set-1",
		DesiredMappings: []DesiredMapping{{WorkspaceID: "ws-1", TagID: tagRecord.ID,
			DataType: schema.DataTypeInt16, ZeroBasedRegister: 0, ShareStartRegister: 40001,
			SpanRegisters: 1, StrideRegisters: 1}},
	})
	var shareErr *Error
	require.ErrorAs(t, err, &shareErr)
	assert.Equal(t, ErrCodeDisabled, shareErr.Code)
	assert.Equal(t, shareStateDisabled, outcome.Outcome)
	revStore.mu.Lock()
	commitCalls := revStore.commitCalls
	revStore.mu.Unlock()
	assert.Zero(t, commitCalls)
	assert.False(t, svc.HasMapping(tagRecord.ID))
}
