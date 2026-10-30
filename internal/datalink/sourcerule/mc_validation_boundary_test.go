package sourcerule

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestServiceUpdateWithRuntimeReconcile_MCInvalidInputDoesNotRestore(t *testing.T) {
	fixture := newMCFixture(t, "W0", 1)
	before := snapshotMCFixture(t, fixture)
	spy := &mcGuardRepositorySpy{MemoryRepository: fixture.repo}
	fixture.svc.repo = spy
	start := "X8"

	_, _, err := fixture.svc.UpdateWithRuntimeReconcile(t.Context(), fixture.rule.ID, UpdateRuleRequest{StartAddress: &start})
	require.ErrorIs(t, err, ErrValidation)
	require.Contains(t, err.Error(), "start_address=X8 protocol=mc_3e")
	require.Empty(t, spy.calls)
	assertMCFixtureSnapshotEqual(t, before, snapshotMCFixture(t, fixture))
}

type mcValidationAfterWriteRepository struct {
	*mcGuardRepositorySpy
	failed bool
}

func (r *mcValidationAfterWriteRepository) Update(ctx context.Context, rule *schema.SourceRule) error {
	if err := r.mcGuardRepositorySpy.Update(ctx, rule); err != nil {
		return err
	}
	if !r.failed {
		r.failed = true
		return validationError("injected validation failure after repository mutation")
	}
	return nil
}

func TestServiceUpdateWithRuntimeReconcile_MCValidationAfterWriteStillRestores(t *testing.T) {
	fixture := newMCFixture(t, "W0", 1)
	before := snapshotMCFixture(t, fixture)
	spy := &mcValidationAfterWriteRepository{mcGuardRepositorySpy: &mcGuardRepositorySpy{MemoryRepository: fixture.repo}}
	fixture.svc.repo = spy
	prefix := "CHANGED"

	_, _, err := fixture.svc.UpdateWithRuntimeReconcile(t.Context(), fixture.rule.ID, UpdateRuleRequest{NamingPrefix: &prefix})
	require.ErrorIs(t, err, ErrValidation)
	require.Equal(t, []string{"Update", "DeleteCandidateSnapshots", "Update", "ReplaceLinks"}, spy.calls)
	assertMCFixtureSnapshotEqual(t, before, snapshotMCFixture(t, fixture))
}
