package sourcerule

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestMixedReviewDecisionMemoryRollback(t *testing.T) {
	dr := device.NewMemoryRepository()
	success := true
	require.NoError(t, dr.Create(t.Context(), &schema.Device{ID: "d", Name: "MC", Protocol: schema.ProtocolMC3E, Status: schema.DeviceStatusActive, LastTestSuccess: &success, ConnectionConfig: `{"host":"127.0.0.1","port":5000,"acquisition":{"version":1,"profile":"mitsubishi-q-iqr-3e-binary-v1"}}`, CreatedAt: time.Now(), UpdatedAt: time.Now()}))
	repo := NewMemoryRepository()
	svc := NewService(repo, device.NewService(dr, nil), point.NewService(point.NewMemoryRepository(), nil), nil)
	rule, err := svc.Create(t.Context(), CreateRuleRequest{ID: "move", DeviceID: "d", NamingPrefix: "M", Enabled: true, LayoutMode: "mixed", MixedItems: []measurement.MixedItem{mixedItem("a", "D100", schema.DataTypeUint16)}})
	require.NoError(t, err)
	before := firstTagCandidateFromMemoryRepo(t, repo, rule.ID, rule.RevisionID)
	_, err = svc.UpsertTagReviewDecision(t.Context(), rule.ID, UpsertTagReviewDecisionRequest{CandidateID: before.ID, Action: schema.SourceRuleTagReviewDecisionActionRename, TagKey: "factory.old"})
	require.NoError(t, err)
	svc.SetShareRuntimeReconciler(ShareRuntimeReconcilerFunc(func(context.Context, RuntimeReconcileRequest) (RuntimeReconcileOutcome, error) {
		return RuntimeReconcileOutcome{}, errors.New("injected reconcile failure")
	}))
	items := []measurement.MixedItem{mixedItem("a", "D200", schema.DataTypeUint16)}
	_, _, err = svc.UpdateWithRuntimeReconcile(t.Context(), rule.ID, UpdateRuleRequest{MixedItems: &items})
	require.Error(t, err)
	restored, err := svc.GetByID(t.Context(), rule.ID)
	require.NoError(t, err)
	require.Equal(t, rule.RevisionID, restored.RevisionID)
	decision, err := svc.GetTagReviewDecision(t.Context(), rule.ID, before.ID)
	require.NoError(t, err)
	require.NotNil(t, decision)
	require.False(t, decision.Stale)
	require.Equal(t, "factory.old", decision.TagKey)
}
