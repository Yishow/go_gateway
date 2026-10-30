package grouppipeline

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func ownerCapture(p *Pipeline, out *[]OwnerObservation, reads *int) {
	p.deps.OwnerObservation = OwnerObservationConfig{
		RecordLimit: 1024,
		ScopeRevision: func(_ context.Context, key groupdelivery.GroupKey) (string, bool) {
			*reads++
			return "observed-workspace-revision", key.WorkspaceID == "workspace-A"
		},
		Observe: func(o OwnerObservation) { *out = append(*out, o) },
	}
}

func ownerCounterFixture(t *testing.T) *longPipeline {
	t.Helper()
	f := newLongPipeline(t)
	require.NoError(t, f.reading(0, "accepted-input", at(1), 100))
	return f
}

func requireOwnerUnavailable(t *testing.T, out []OwnerObservation) {
	t.Helper()
	require.NotEmpty(t, out)
	for _, o := range out {
		require.Equal(t, OwnerObservation{}, o)
	}
}

func ownerRefusalFixture(t *testing.T) (*ownershipPipelineFixture, *managed, measurement.SampleEnvelope) {
	t.Helper()
	f := newOwnershipPipelineFixture(t, shapeMeanWide)
	f.clock = at(21)
	env := envelope(f.source.Group, "accepted-before-cutover", at(23), uint64(7))
	require.NoError(t, f.pipe.AcceptSample(t.Context(), env))
	f.commit()
	f.clock = at(24)
	f.groups.groups = append(f.groups.groups, f.target.Group)
	f.groups.snapshots[f.target.GroupID] = []*workspace.WriteGroupAppliedSnapshot{f.target}
	target, reason := f.pipe.build(t.Context(), f.target)
	require.Empty(t, reason)
	f.pipe.installManaged(target)
	original, err := f.store.OwnershipView(t.Context(), groupKey(f.target))
	require.NoError(t, err)
	require.Equal(t, original, target.boundary.ReplayState().OwnershipView)
	return f, target, env
}
