package grouppipeline

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/groupdelivery"

	"github.com/stretchr/testify/require"
)

func intakeCapture(p *Pipeline, out *[]IntakeQuotaObservation, reads *int) {
	p.deps.IntakeQuotaObservation = IntakeQuotaObservationConfig{
		RecordLimit: 1024, IntakeEnabled: true, QuotaEnabled: true,
		ScopeRevision: func(_ context.Context, key groupdelivery.GroupKey) (string, bool) {
			*reads++
			return "workspace-revision", key.WorkspaceID == "workspace-A"
		},
		Observe: func(o IntakeQuotaObservation) { *out = append(*out, o) },
	}
}
func intakeKind(out []IntakeQuotaObservation, kind IntakeQuotaObservationKind) []IntakeQuotaObservation {
	var found []IntakeQuotaObservation
	for _, o := range out {
		if o.Available && o.Kind == kind {
			found = append(found, o)
		}
	}
	return found
}
func intakeQuotaFixture(t *testing.T, scope groupdelivery.QuotaScope) (*fixture, *Pipeline, *[]IntakeQuotaObservation, *int) {
	t.Helper()
	f := newFixture(t)
	cfg := groupdelivery.QuotaConfig{GroupMaxBytes: 1}
	if scope == groupdelivery.QuotaScopeGlobal {
		cfg = groupdelivery.QuotaConfig{GlobalMaxBytes: 1}
	}
	var err error
	f.store, err = f.store.WithQuota(cfg)
	require.NoError(t, err)
	p := f.pipeline()
	require.NoError(t, p.Reconcile(t.Context()))
	group := f.groups.groups[0]
	require.NoError(t, p.AcceptSample(t.Context(), envelope(group, "first", at(1), 1.0)))
	var out []IntakeQuotaObservation
	reads := 0
	intakeCapture(p, &out, &reads)
	return f, p, &out, &reads
}
