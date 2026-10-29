package grouppipeline

import (
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"go-gateway/internal/datalink/counterproof"

	"github.com/stretchr/testify/require"
)

func TestCounterQuarantineRefusesIntakeAlreadyWaitingForProof(t *testing.T) {
	f := newLongPipeline(t)
	f.clock = at(3)
	damageCounterProof(t, f, "counter-replay-proof-conflict")

	entered, release := make(chan struct{}), make(chan struct{})
	var once, unblock sync.Once
	t.Cleanup(func() { unblock.Do(func() { close(release) }) })
	reject := func(query string) error {
		if strings.Contains(query, "SELECT payload FROM wg_runtime_versions") {
			once.Do(func() { close(entered) })
			select {
			case <-release:
			case <-t.Context().Done():
				return t.Context().Err()
			}
		}
		return nil
	}
	f.localFault.Reject.Store(&reject)
	t.Cleanup(func() { f.localFault.Reject.Store(nil) })

	first, waiting := make(chan error, 1), make(chan error, 1)
	go func() { first <- f.pipe.AcceptSample(t.Context(), envelope(f.group, "first", at(1), uint64(100))) }()
	<-entered
	go func() { waiting <- f.pipe.AcceptSample(t.Context(), envelope(f.group, "waiting", at(2), uint64(105))) }()
	// The second intake must already hold a reference to the group before the
	// first one detects the damaged proof and quarantines it.
	require.Eventually(t, func() bool {
		stacks := make([]byte, 1<<20)
		n := runtime.Stack(stacks, true)
		for stack := range strings.SplitSeq(string(stacks[:n]), "\n\n") {
			if strings.Contains(stack, "grouppipeline.(*Pipeline).withCounterProofGate") && strings.Contains(strings.ToLower(stack), "semacquire") {
				return true
			}
		}
		return false
	}, 5*time.Second, time.Millisecond)
	unblock.Do(func() { close(release) })
	require.ErrorIs(t, <-first, counterproof.ErrCounterBindingProofConflict)
	require.ErrorIs(t, <-waiting, counterproof.ErrCounterBindingProofConflict)
	f.localFault.Reject.Store(nil)
	entries, err := f.pipe.CounterReplayQuarantines(t.Context(), f.group.ID)
	require.NoError(t, err)
	require.Len(t, entries, 1)
	require.Zero(t, entries[0].PendingSamples)
}
