package grouppipeline

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/workspace"
	"testing"
)

func TestMixedPipelineShutdownRestartKeepsOpenCoverage(t *testing.T) {
	f, pipe, snap := mixedPipelineFixture(t)
	f.clock = at(10)
	require.NoError(t, pipe.Start(t.Context()))
	require.NoError(t, pipe.AcceptSample(t.Context(), envelope(snap.Group, "before-shutdown", at(10), 2.0)))
	require.NoError(t, pipe.StopContext(t.Context()))
	// Process shutdown is not a canonical revision end or a Disable operation.
	restarted := f.pipeline()
	require.NoError(t, restarted.Reconcile(t.Context()))
	require.True(t, restarted.WantsSample("device-1", "point-1", "tag-t"))
	f.clock = at(15)
	restarted.TickAll(t.Context())
	restored, err := f.store.Restore(t.Context(), groupKey(snap))
	require.NoError(t, err)
	require.NotNil(t, restored.Mixed)
	require.True(t, restored.Mixed.EndsAt.IsZero())
	f.clock = at(20)
	require.NoError(t, restarted.AcceptSample(t.Context(), envelope(snap.Group, "after-restart", at(20), 4.0)))
	f.clock = at(30)
	restarted.TickAll(t.Context())
	restored, err = f.store.Restore(t.Context(), groupKey(snap))
	require.NoError(t, err)
	require.Equal(t, at(30), restored.Mixed.Progress.ClosedThrough)
	require.True(t, restored.Mixed.EndsAt.IsZero())
}

func TestMixedPipelineDisabledRestartDrainsOriginalJournal(t *testing.T) {
	f, pipe, snap := mixedPipelineFixture(t)
	f.clock = at(10)
	require.NoError(t, pipe.AcceptSample(t.Context(), envelope(snap.Group, "accepted-before-disable", at(10), 2.0)))
	require.NoError(t, pipe.StopContext(t.Context()))
	f.groups.groups[0] = liveGroup(snap.Group, workspace.WriteGroupStatusDisabled)
	restarted := f.pipeline()
	require.NoError(t, restarted.Reconcile(t.Context()))
	require.False(t, restarted.WantsSample("device-1", "point-1", "tag-t"))
	require.Len(t, restarted.Status(), 1, "disabled restart still owns its accepted journal")
	f.clock = at(15)
	restarted.TickAll(t.Context())
	restored, err := f.store.Restore(t.Context(), groupKey(snap))
	require.NoError(t, err)
	require.NotNil(t, restored.Mixed)
	require.True(t, restored.Mixed.EndsAt.IsZero())
	require.Empty(t, restored.Samples)
	require.Empty(t, restarted.Status())
}
