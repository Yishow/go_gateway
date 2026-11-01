package collector

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAcquisitionGeometryBatchQuarantineIsNonblockingAndAllOrNone(t *testing.T) {
	s := seedAcquisitionGeometry(t)
	before := s.CaptureAcquisitionGeometry("d")
	s.mu.Lock()
	_, ok := s.TryBeginAcquisitionGeometryMutations(t.Context(), []string{"d"})
	s.mu.Unlock()
	require.False(t, ok)
	require.Same(t, before.publication, s.CaptureAcquisitionGeometry("d").publication)
	for _, ids := range [][]string{{"d", ""}, {"d", "d"}} {
		_, ok = s.TryBeginAcquisitionGeometryMutations(t.Context(), ids)
		require.False(t, ok)
		require.Same(t, before.publication, s.CaptureAcquisitionGeometry("d").publication)
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	_, ok = s.TryBeginAcquisitionGeometryMutations(ctx, []string{"d"})
	require.False(t, ok)
	tokens, ok := s.TryBeginAcquisitionGeometryMutations(t.Context(), []string{"d", "other"})
	require.True(t, ok)
	require.Len(t, tokens, 2)
	require.Nil(t, s.currentAcquisitionGeometryLocked("d"))
	for _, token := range tokens {
		s.FinishAcquisitionGeometryMutation(token)
	}
	require.Nil(t, s.currentAcquisitionGeometryLocked("d"))
	require.Equal(t, before.slots, s.CaptureAcquisitionGeometry("d").slots)
}
