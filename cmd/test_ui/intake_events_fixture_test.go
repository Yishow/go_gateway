package main

import (
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/workspace"
	"go-gateway/internal/diagnostics"

	"github.com/stretchr/testify/require"
)

type intakeEventFixture struct{ deliveryEventFixture }

// Reuse the existing disposable SQLite/native canonical group fixture unchanged.
// It contains an actual pending output; no dispatcher starts unless a test asks.
func newIntakeEventFixture(t *testing.T, fileLog ...bool) intakeEventFixture {
	t.Helper()
	return intakeEventFixture{newDeliveryEventFixture(t, fileLog...)}
}

func (f intakeEventFixture) reconcileCurrent(t *testing.T) {
	t.Helper()
	require.NoError(t, f.env.services.groupPipe.Reconcile(t.Context()))
	found := false
	for _, status := range f.env.services.groupPipe.Status() {
		if status.GroupID == f.group.ID && status.Revision == f.group.AppliedRevision {
			require.Equal(t, grouppipeline.StateActive, status.State)
			found = true
		}
	}
	require.True(t, found, "the native pipeline must actually install the exact current revision")
}

func (f intakeEventFixture) disable(t *testing.T) *workspace.WriteGroupSaveResult {
	t.Helper()
	f.env.services.writeGroups.WithClock(func() time.Time { return time.Now().UTC() })
	current, err := f.env.services.writeGroups.Get(t.Context(), f.group.ID)
	require.NoError(t, err)
	record, err := f.env.services.workspace.ReadSnapshot(t.Context())
	require.NoError(t, err)
	disabled, err := f.env.services.writeGroups.Disable(t.Context(), f.group.ID, workspace.WriteGroupMutation{WorkspaceID: record.ID, ExpectedWorkspaceRevision: record.DatabaseSetupRevision, ExpectedGroupRevision: current.Group.Revision, ExpectedConnectorRevision: current.Group.Destination.ConnectorRevision})
	require.NoError(t, err)
	require.Equal(t, workspace.WriteGroupStatusDisabled, disabled.Group.Status)
	return disabled
}

func newObservedQuotaPipeline(t *testing.T, f intakeEventFixture, broker *diagnostics.Broker, quota groupdelivery.QuotaConfig, clock *testClock, configure ...func(*grouppipeline.IntakeQuotaObservationConfig)) (*grouppipeline.Pipeline, *groupdelivery.Store) {
	t.Helper()
	store, err := groupdelivery.NewStore(f.env.db).WithCanonicalIntake().WithQuota(quota)
	require.NoError(t, err)
	delivery, owner := newDeliveryEventObservers(broker, f.env.services.workspace)
	config := newIntakeQuotaEventObserver(broker, delivery.ScopeRevision)
	for _, apply := range configure {
		apply(&config)
	}
	pipe := grouppipeline.New(grouppipeline.Dependencies{Groups: f.env.services.writeGroups, Tags: f.env.services.tag, Destinations: f.env.services.dbTarget, Inspector: dbtarget.NewReadOnlyTableInspector(f.env.services.dbTarget), Store: store, OwnerObservation: owner, IntakeQuotaObservation: config}, grouppipeline.Config{Now: clock.now, WorkerInterval: 10 * time.Millisecond, TickInterval: time.Hour, ReconcileInterval: time.Hour})
	t.Cleanup(func() { require.NoError(t, pipe.Stop(3*time.Second)) })
	require.NoError(t, pipe.Reconcile(t.Context()))
	return pipe, store
}

func intakeEventClock() *testClock { clock := &testClock{}; clock.set(time.Now().UTC()); return clock }

func quotaEventRefusal(t *testing.T, f intakeEventFixture, pipe *grouppipeline.Pipeline, clock *testClock) *groupdelivery.QuotaError {
	t.Helper()
	sample := f.env.envelope(f.group, 0, "actual-quota-refusal", clock.now(), 123.0)
	err := pipe.AcceptSample(t.Context(), sample)
	require.ErrorIs(t, err, groupdelivery.ErrQuotaExceeded)
	var quota *groupdelivery.QuotaError
	require.ErrorAs(t, err, &quota)
	return quota
}

func waitIntakeEvent(t *testing.T, broker *diagnostics.Broker, code string) diagnostics.Event {
	t.Helper()
	events := waitDeliveryEvents(t, broker, code)
	require.Len(t, events, 1)
	return events[0]
}
