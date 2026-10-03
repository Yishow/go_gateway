//go:build f_write_group_fixture

package main

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFixtureCleanupReportsCapacityRestoreFailure(t *testing.T) {
	fixtureOwnedDatabasePath(t, fixtureMarkerValue)
	env := newOutageEnv(t)
	fixture, ok := env.services.fixture.(*groupFixture)
	if !ok {
		t.Fatalf("fixture type = %T, want tagged controller", env.services.fixture)
	}
	fixture.port = 0
	cleanup, err := fixture.configure(t.Context(), env.db, &env.services)
	if err != nil {
		t.Fatalf("configure fixture: %v", err)
	}
	if _, err := fixture.capacity.configure(t.Context(), true); err != nil {
		t.Fatalf("enable fixture capacity: %v", err)
	}
	if err := env.db.Close(); err != nil {
		t.Fatalf("close fixture database: %v", err)
	}
	results := reflect.ValueOf(cleanup).Call(nil)
	if len(results) != 1 {
		t.Fatalf("cleanup returned %d values, want one error", len(results))
	}
	cleanupErr, ok := results[0].Interface().(error)
	if !ok || cleanupErr == nil {
		t.Fatalf("cleanup result = %#v, want non-nil error", results[0].Interface())
	}
}

func TestFixtureCleanupReleasesHeldClosureBeforeCapacitySQL(t *testing.T) {
	fixtureOwnedDatabasePath(t, fixtureMarkerValue)
	env := newOutageEnv(t)
	fixture, ok := env.services.fixture.(*groupFixture)
	if !ok {
		t.Fatalf("fixture type = %T, want tagged controller", env.services.fixture)
	}
	fixture.port = 0
	cleanup, err := fixture.configure(t.Context(), env.db, &env.services)
	if err != nil {
		t.Fatalf("configure fixture: %v", err)
	}
	activeAt := time.Now().UTC()
	fixture.clock.Set(activeAt)
	env.services.writeGroups.WithClock(fixture.clock.Now)
	group := env.createAndApply(t, "A")
	fixture.pipe = env.services.groupPipe
	resolved, err := env.services.writeGroups.ResolveAppliedAt(t.Context(), group.ID, activeAt.Add(time.Minute))
	if err != nil {
		t.Fatalf("resolve applied group: %v", err)
	}
	base := resolved.EffectiveAt.Truncate(10 * time.Second).Add(10 * time.Second)
	fixture.clock.Set(base.Add(time.Second))
	if err := fixture.pipe.Reconcile(t.Context()); err != nil {
		t.Fatalf("reconcile fixture group: %v", err)
	}
	if _, err := fixture.capacity.configure(t.Context(), true); err != nil {
		t.Fatalf("enable fixture capacity: %v", err)
	}
	if _, err := fixture.faults.configure(t.Context(), fixtureFaultRequest{
		Kind: fixtureFaultClosureHold, GroupID: group.ID, Enabled: true,
	}); err != nil {
		t.Fatalf("enable closure hold: %v", err)
	}
	observed := base.Add(5 * time.Second)
	fixture.clock.Set(observed.Add(time.Second))
	if err := fixture.pipe.AcceptSample(t.Context(), env.envelope(group, 0, "cleanup-held-temperature", observed, 21.5)); err != nil {
		t.Fatalf("accept temperature sample: %v", err)
	}
	if err := fixture.pipe.AcceptSample(t.Context(), env.envelope(group, 1, "cleanup-held-pressure", observed, int64(1013))); err != nil {
		t.Fatalf("accept pressure sample: %v", err)
	}
	fixture.clock.Set(base.Add(11 * time.Second))
	tickDone := make(chan struct{})
	go func() {
		fixture.pipe.TickAll(t.Context())
		close(tickDone)
	}()
	require.Eventually(t, func() bool {
		for _, state := range fixture.faults.snapshot() {
			if state.GroupID == group.ID {
				return state.Holding
			}
		}
		return false
	}, time.Second, time.Millisecond)

	cleanupDone := make(chan struct{})
	cleanupResult := make(chan error, 1)
	go func() {
		cleanupResult <- cleanup()
		close(cleanupDone)
	}()
	select {
	case <-cleanupDone:
	case <-time.After(2 * time.Second):
		releaseFixtureClosureGate(group.ID)
		select {
		case <-cleanupDone:
		case <-time.After(time.Second):
		}
		t.Fatal("fixture cleanup blocked behind a held closure transaction")
	}
	if err := <-cleanupResult; err != nil {
		t.Fatalf("fixture cleanup returned error: %v", err)
	}
	select {
	case <-tickDone:
	case <-time.After(time.Second):
		t.Fatal("held closure tick did not finish during fixture cleanup")
	}
}
