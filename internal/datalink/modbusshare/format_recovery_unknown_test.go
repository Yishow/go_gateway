package modbusshare

import (
	"context"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sqltx"
	"go-gateway/internal/datalink/tag"
	"testing"
)

type recoveryUnknownStore struct{ *mockWorkspaceRevisionStore }

func (s recoveryUnknownStore) CheckFormatRecoveryWitness(context.Context, string, string) error {
	return nil
}

func TestFormatRecoveryRealPlanUnknownCannotRemainRecoverablePending(t *testing.T) {
	r, s, tags, store := setupReconcilerFixture(t)
	r.revStore = recoveryUnknownStore{store}
	settings := s.Settings()
	state := HydrationState{State: HydrationStatePending, WorkspaceID: "ws-1", WorkspaceRevision: "rev-1", SettingsRevision: settings.SettingsRevision}
	s.SetHydrationState(state)
	fence, err := s.CaptureFormatRecoveryFence("ws-1", "rev-1", settings.SettingsRevision)
	require.NoError(t, err)
	tagRecord, err := tags.Create(t.Context(), tag.CreateTagRequest{Key: "recovery-unknown", DisplayName: "value", DataType: schema.DataTypeInt16})
	require.NoError(t, err)
	store.failCommitAt = 1
	store.updateErr = sqltx.ErrCommitOutcomeUnknown
	_, err = r.ReconcileFormatRecovery(t.Context(), ReconcileRequest{WorkspaceID: "ws-1", ExpectedWorkspaceRevision: "rev-1", ExpectedSettingsRevision: settings.SettingsRevision, DesiredMappings: []DesiredMapping{{WorkspaceID: "ws-1", SourceRuleID: "rule-1", SourceRuleRevision: "rule-rev", TagID: tagRecord.ID, DataType: schema.DataTypeInt16, ShareStartRegister: 40001, SpanRegisters: 1, StrideRegisters: 1}}}, fence, "clean")
	require.ErrorIs(t, err, sqltx.ErrCommitOutcomeUnknown)
	_, dirty, e := store.GetRevision(t.Context(), "ws-1")
	require.NoError(t, e)
	require.True(t, dirty, "unknown real commit must revoke known-pending recovery")
	require.Error(t, s.CheckFormatRecoveryFence(t.Context(), fence))
	require.Equal(t, 1, store.commitCalls, "never inverse an unknown real commit")
}
