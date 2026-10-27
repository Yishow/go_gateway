package grouppipeline

import (
	"errors"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/counterproof"

	"github.com/stretchr/testify/require"
)

func TestCounterReplayRestoreInfrastructureFailuresRemainErrorsAndRecoverOtherScopes(t *testing.T) {
	for _, stage := range []string{"checkpoint", "journal", "association"} {
		t.Run(stage, func(t *testing.T) {
			f := newLongPipeline(t)
			if stage == "association" {
				f = legacyQuarantineFixture(t)
			} else {
				require.NoError(t, f.reading(0, "pending", at(1), 100))
			}
			// A healthy pending scope is recovered after the first scope fails.
			healthy := healthyCounterRevision(t, f, "healthy-group", "rev-healthy")
			ready, reason := f.pipe.build(t.Context(), healthy)
			require.Empty(t, reason)
			require.NoError(t, ready.boundary.AcceptSample(t.Context(), envelope(healthy.Group, "healthy-1", at(2), uint64(7))))
			before := counterDurableState(t, f.store)
			f.groups.groups = nil
			f.pipe.boundaries = map[string]*managed{}
			injected := errors.New("private restore driver failure")
			calls := 0
			reject := func(query string) error {
				matches := stage == "checkpoint" && strings.Contains(query, "SELECT next_close FROM wg_delivery_checkpoints") ||
					stage == "journal" && strings.Contains(query, "SELECT payload, payload_digest, workspace_id FROM wg_delivery_samples") ||
					stage == "association" && strings.Contains(query, "INSERT INTO wg_counter_proof_associations")
				if matches && calls == 0 {
					calls++
					return injected
				}
				return nil
			}
			f.localFault.Reject.Store(&reject)
			err := f.pipe.Reconcile(t.Context())
			f.localFault.Reject.Store(nil)
			require.Equal(t, 1, calls, "the intended durable boundary must actually fail")
			require.ErrorIs(t, err, injected)
			require.NotContains(t, err.Error(), "private restore driver failure", "the boundary retains its safe error text")
			entries, err := f.pipe.CounterReplayQuarantines(t.Context(), f.group.ID)
			require.NoError(t, err)
			require.Empty(t, entries, "an unclassified outage is not damaged original proof")
			require.Equal(t, before, counterDurableState(t, f.store))
			recovered := f.pipe.boundaries[boundaryKey(healthy.GroupID, healthy.AppliedRevision)]
			require.NotNil(t, recovered, "recovery continues to the other accepted scope after the infrastructure failure")
			require.True(t, recovered.drainOnly, "a historical journal without a live group remains recovery-only")
		})
	}
}

func TestCounterReplayQuarantineStopsPreviouslyCapturedUnstartedTick(t *testing.T) {
	f := newLongPipeline(t)
	require.NoError(t, f.reading(0, "pending", at(1), 100))
	captured, resume, done := make(chan struct{}), make(chan struct{}), make(chan struct{})
	f.pipe.config.Now = func() time.Time { close(captured); <-resume; return at(600) }
	go func() { f.pipe.TickAll(t.Context()); close(done) }()
	<-captured
	damageCounterProof(t, f, "counter-replay-proof-conflict")
	before := counterDurableState(t, f.store)
	err := f.reading(0, "refused", at(2), 105)
	require.ErrorIs(t, err, counterproof.ErrCounterBindingProofConflict)
	require.False(t, f.pipe.WantsSample("device-1", "point-1", longKwhFirst))
	close(resume)
	<-done
	require.Equal(t, before, counterDurableState(t, f.store), "a copied boundary cannot begin closure after its proof is quarantined")
}

func TestCounterReplayQuarantineGateDoesNotSerializeUnrelatedIntake(t *testing.T) {
	f := newLongPipeline(t)
	healthy := healthyCounterRevision(t, f, "healthy-group", "rev-healthy")
	good, reason := f.pipe.build(t.Context(), healthy)
	require.Empty(t, reason)
	f.pipe.installManaged(good)
	old := f.pipe.boundaries[boundaryKey(f.group.ID, f.group.AppliedRevision)]
	entered, release, finished := make(chan struct{}), make(chan struct{}), make(chan struct{})
	go func() {
		_ = f.pipe.withCounterProofGate(old, func() error { close(entered); <-release; return nil })
		close(finished)
	}()
	<-entered
	quarantined := make(chan struct{})
	go func() {
		f.pipe.setCounterQuarantine(groupKey(f.snapshot), "counter-replay-proof-missing")
		close(quarantined)
	}()
	accepted := make(chan error, 1)
	go func() {
		accepted <- f.pipe.AcceptSample(t.Context(), envelope(healthy.Group, "healthy-1", at(2), uint64(7)))
	}()
	received := false
	select {
	case err := <-accepted:
		received = true
		require.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Error("an unrelated source must not wait behind the affected revision's operation gate")
	}
	close(release)
	<-finished
	<-quarantined
	if !received {
		require.NoError(t, <-accepted)
	}
	// A previously captured operation cannot execute after the gate is closed.
	called := false
	require.NoError(t, f.pipe.withCounterProofGate(old, func() error { called = true; return nil }))
	require.False(t, called)
}
