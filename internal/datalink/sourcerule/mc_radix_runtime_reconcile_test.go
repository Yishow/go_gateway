package sourcerule

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServiceUpdateWithRuntimeReconcile_MCRadixConflictDoesNotRestore(t *testing.T) {
	fixture := newMCFixture(t, "W0", 17)
	before := snapshotMCFixture(t, fixture)
	spy := &mcGuardRepositorySpy{MemoryRepository: fixture.repo}
	fixture.svc.repo = spy
	prefix := "RUNTIME-RENAMED"

	_, _, err := fixture.svc.UpdateWithRuntimeReconcile(t.Context(), fixture.rule.ID, UpdateRuleRequest{NamingPrefix: &prefix})
	require.Error(t, err)
	require.ErrorIs(t, err, ErrValidation)
	require.ErrorIs(t, err, errMCAddressPlanConflict)
	require.Contains(t, err.Error(), "mc_address_plan_conflict")
	require.Empty(t, spy.calls)
	assertMCFixtureSnapshotEqual(t, before, snapshotMCFixture(t, fixture))
}
