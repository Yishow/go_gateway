package sourcerule

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type candidateScopeWorkspaceStub struct {
	id      string
	devices []string
}

func (s candidateScopeWorkspaceStub) ShareOwnershipSnapshot(context.Context) (workspaceID string, deviceIDs []string, err error) {
	return s.id, s.devices, nil
}

type candidateScopeHydrationStub struct {
	state modbusshare.HydrationState
	err   error
}

func (s candidateScopeHydrationStub) CheckHydration(context.Context) (modbusshare.HydrationState, error) {
	return s.state, s.err
}

func TestCandidateScopeRejectsHydrationBeforeRuleLookup(t *testing.T) {
	svc := &Service{}
	svc.SetCandidateScope(candidateScopeWorkspaceStub{id: "ws-a", devices: []string{"device-a"}}, candidateScopeHydrationStub{
		state: modbusshare.HydrationState{State: modbusshare.HydrationStatePending},
	})

	err := svc.ValidateCandidateScope(context.Background(), "missing-rule", CandidateScopeRequest{
		WorkspaceID:               "ws-a",
		ExpectedWorkspaceRevision: "rev-a",
		RevisionID:                "rule-rev-a",
	})

	var typed *modbusshare.Error
	require.ErrorAs(t, err, &typed)
	require.Equal(t, modbusshare.ErrCodeHydrationRequired, typed.Code)
}

func TestCandidateScopeRejectsForeignRuleAndStaleRevision(t *testing.T) {
	repo := NewMemoryRepository()
	require.NoError(t, repo.Create(context.Background(), &schema.SourceRule{ID: "rule-a", DeviceID: "device-a", RevisionID: "rule-rev-a"}))
	svc := &Service{repo: repo}
	svc.SetCandidateScope(candidateScopeWorkspaceStub{id: "ws-a", devices: []string{"device-a"}}, candidateScopeHydrationStub{
		state: modbusshare.HydrationState{State: modbusshare.HydrationStateReady, Readiness: true, WorkspaceRevision: "rev-a"},
	})

	err := svc.ValidateCandidateScope(context.Background(), "rule-a", CandidateScopeRequest{
		WorkspaceID:               "ws-a",
		ExpectedWorkspaceRevision: "rev-a",
		RevisionID:                "stale-rev",
	})
	var typed *modbusshare.Error
	require.ErrorAs(t, err, &typed)
	require.Equal(t, modbusshare.ErrCodeRevisionConflict, typed.Code)

	err = svc.ValidateCandidateScope(context.Background(), "rule-a", CandidateScopeRequest{
		WorkspaceID:               "ws-b",
		ExpectedWorkspaceRevision: "rev-a",
		RevisionID:                "rule-rev-a",
	})
	require.ErrorAs(t, err, &typed)
	require.Equal(t, modbusshare.ErrCodeWorkspaceScope, typed.Code)
}

func TestConfiguredLocalModbusRegisterRejectsShortStride(t *testing.T) {
	rule := &schema.SourceRule{
		ShareStartRegister: intPtrForCandidateScopeTest(40001),
		ShareStride:        intPtrForCandidateScopeTest(1),
	}

	_, err := configuredLocalModbusRegister(rule, "40001", 0, 2)
	var typed *modbusshare.Error
	require.ErrorAs(t, err, &typed)
	require.Equal(t, modbusshare.ErrCodeInvalidGeometry, typed.Code)
}

func intPtrForCandidateScopeTest(value int) *int { return &value }
