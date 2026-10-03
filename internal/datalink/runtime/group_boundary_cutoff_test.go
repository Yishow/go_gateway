package runtime

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGroupBoundarySetUntilAndAcceptSampleSynchronizeCutoff(t *testing.T) {
	f := newBoundaryFixture(t)
	boundary := f.build(t)
	boundary.SetUntil(boundaryAt(20))

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Go(func() { boundary.SetUntil(boundaryAt(10)) })
		wg.Go(func() {
			// This sample is after either cutoff, so it must never open a bucket
			// while SetUntil updates the boundary concurrently.
			_ = boundary.AcceptSample(t.Context(), f.envelope(i%len(f.members), "cutoff", 21, 1.0))
		})
	}
	wg.Wait()

	require.NoError(t, boundary.AcceptSample(t.Context(), f.envelope(0, "after-cutoff", 21, 1.0)))
	boundary.mu.Lock()
	defer boundary.mu.Unlock()
	require.Zero(t, boundary.assembler.OpenBuckets())
}
