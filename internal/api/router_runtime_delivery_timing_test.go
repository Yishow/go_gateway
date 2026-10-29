package api

import (
	"math"
	"net/http"
	"testing"
	"time"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestRuntimeDeliveryTimingOldestPendingAgeUsesDurableCreatedAtAcrossAllRevisions(t *testing.T) {
	r := newRuntimeDeliveryRig(t)
	r.historicalRows(t, 1)
	createdAt := r.clock
	var storedCreated, bucket string
	require.NoError(t, r.f.db.QueryRow(`SELECT created_at,bucket_start FROM wg_delivery_outbox`).Scan(&storedCreated, &bucket))
	parsedCreated, err := time.Parse(time.RFC3339Nano, storedCreated)
	require.NoError(t, err)
	require.Equal(t, createdAt, parsedCreated)
	require.NotEqual(t, storedCreated, bucket)
	r.clock = r.clock.Add(125*time.Second + 500*time.Millisecond)
	raw, err := r.pipe.DeliveryForRevision(t.Context(), r.f.group["id"].(string), r.f.group["applied_revision"].(string))
	require.NoError(t, err)
	require.Equal(t, 125.5, raw.OldestPendingSeconds, "production Store age is measured from durable created_at, not the older bucket")
	g := runtimeTimingReadOnly(t, r)
	require.NotNil(t, g.Delivery)
	require.NotNil(t, g.Delivery.OldestPendingSeconds, "scoped GET must retain the actual fractional age across all group revisions")
	require.Equal(t, raw.OldestPendingSeconds, *g.Delivery.OldestPendingSeconds)
	require.Equal(t, "group-all-revisions", g.Delivery.OldestPendingScope)
	require.Zero(t, g.Delivery.Current.Stages.Queued, "historical age is not restricted to current pending rows")
	require.EqualValues(t, 1, g.Delivery.Backlog[0].Pending)
}

func TestRuntimeDeliveryTimingHistoricalAgeNotCurrentRevisionOrDraftDestination(t *testing.T) {
	r := newRuntimeDeliveryRig(t)
	r.historicalRows(t, 1)
	oldCreated := r.clock
	r.clock = r.clock.Add(time.Second)
	r.update(t, "meter_B", 30)
	b := r.apply(t).Group.AppliedRevision
	frozenB, err := r.f.groups.ResolveAppliedAt(t.Context(), r.f.group["id"].(string), r.clock.Add(time.Minute))
	require.NoError(t, err)
	require.Equal(t, b, frozenB.AppliedRevision)
	r.clock = frozenB.EffectiveAt.Add(time.Second)
	require.NoError(t, r.pipe.Reconcile(t.Context()))
	r.accept(t, "current-B")
	r.clock = frozenB.EffectiveAt.Add(30 * time.Second)
	r.pipe.TickAll(t.Context())
	currentCreated := r.clock
	r.update(t, "meter_C", 45)
	r.clock = r.clock.Add(125*time.Second + 500*time.Millisecond)
	r.useReader(r.pipe)
	raw, err := r.pipe.DeliveryForRevision(t.Context(), frozenB.GroupID, b)
	require.NoError(t, err)
	require.Positive(t, raw.RevisionStages.Stages.Queued)
	require.Equal(t, r.clock.Sub(oldCreated).Seconds(), raw.OldestPendingSeconds)
	require.Greater(t, raw.OldestPendingSeconds, r.clock.Sub(currentCreated).Seconds())
	g := runtimeTimingReadOnly(t, r)
	require.Equal(t, "meter_C", g.Saved.Destination.Table)
	require.Equal(t, b, *g.EffectiveRevision)
	require.Equal(t, "meter_B", g.Effective.Destination.Table)
	require.Equal(t, raw.OldestPendingSeconds, *g.Delivery.OldestPendingSeconds)
	require.Equal(t, b, g.Delivery.Current.Binding.EntityRevision)
	require.Len(t, g.Delivery.Backlog, 1)
	require.Equal(t, "rev7", g.Delivery.Backlog[0].Binding.EntityRevision)
	require.Equal(t, "C1", g.Delivery.Backlog[0].Destination.ConnectorID)
	require.Equal(t, "7", g.Delivery.Backlog[0].Destination.ConnectorRevision)
	require.Equal(t, "old_table", g.Delivery.Backlog[0].Destination.Table)
}

func TestRuntimeDeliveryTimingDisabledAndMissingEffectiveProofRetainHistoricalAge(t *testing.T) {
	for _, missingProof := range []bool{false, true} {
		t.Run(map[bool]string{false: "disabled", true: "missing-proof"}[missingProof], func(t *testing.T) {
			r := newRuntimeDeliveryRig(t)
			r.historicalRows(t, 3)
			r.update(t, "new_draft_table", 30)
			saved, err := r.f.groups.Get(t.Context(), r.f.group["id"].(string))
			require.NoError(t, err)
			_, err = r.f.groups.Disable(t.Context(), saved.Group.ID, workspace.WriteGroupMutation{WorkspaceID: saved.Group.WorkspaceID, ExpectedWorkspaceRevision: saved.WorkspaceRevision, ExpectedGroupRevision: saved.Group.Revision, ExpectedConnectorRevision: saved.Group.Destination.ConnectorRevision})
			require.NoError(t, err)
			require.NoError(t, r.pipe.Reconcile(t.Context()))
			if missingProof {
				_, err = r.f.db.Exec(`DELETE FROM wg_runtime_versions`)
				require.NoError(t, err)
			}
			r.clock = r.clock.Add(125*time.Second + 500*time.Millisecond)
			g := runtimeTimingReadOnly(t, r)
			require.False(t, g.Saved.Enabled)
			require.Equal(t, "new_draft_table", g.Saved.Destination.Table)
			require.NotNil(t, g.Delivery.OldestPendingSeconds)
			require.Equal(t, 125.5, *g.Delivery.OldestPendingSeconds)
			require.Equal(t, "group-all-revisions", g.Delivery.OldestPendingScope)
			require.Len(t, g.Delivery.Backlog, 1)
			require.Equal(t, "C1", g.Delivery.Backlog[0].Destination.ConnectorID)
			require.Equal(t, "old_table", g.Delivery.Backlog[0].Destination.Table)
			require.EqualValues(t, 3, g.Delivery.Backlog[0].Pending)
			if missingProof {
				require.Nil(t, g.Effective)
				require.Nil(t, g.Delivery.Current, "retained age does not establish current applied proof")
			} else if g.Delivery.Current.Intake != nil {
				require.NotEqual(t, "active", g.Delivery.Current.Intake.State)
			}
		})
	}
}

func TestRuntimeDeliveryTimingZeroAgePresenceAmbiguityStaysUnavailable(t *testing.T) {
	for _, mode := range []string{"no-row", "new-row-exactly-zero", "only-committed"} {
		t.Run(mode, func(t *testing.T) {
			r := newRuntimeDeliveryRig(t)
			if mode != "no-row" {
				r.accept(t, "zero-age")
				key := r.tick(t)
				if mode == "only-committed" {
					r.send(t, key)
					r.clock = r.clock.Add(125*time.Second + 500*time.Millisecond)
				}
			}
			raw, err := r.pipe.DeliveryForRevision(t.Context(), r.f.group["id"].(string), r.f.group["applied_revision"].(string))
			require.NoError(t, err)
			require.Zero(t, raw.OldestPendingSeconds)
			g := runtimeTimingReadOnly(t, r)
			require.Nil(t, g.Delivery.OldestPendingSeconds, "non-null zero would claim presence the producer cannot distinguish")
			switch mode {
			case "no-row":
				require.Zero(t, g.Delivery.Current.Stages.Queued)
				require.Zero(t, g.Delivery.Current.Stages.SQLCommitted)
			case "new-row-exactly-zero":
				require.EqualValues(t, 1, g.Delivery.Current.Stages.Queued)
			case "only-committed":
				require.Zero(t, g.Delivery.Current.Stages.Queued)
				require.EqualValues(t, 1, g.Delivery.Current.Stages.SQLCommitted)
				require.NotNil(t, g.Delivery.Current.LastCommit)
			}
		})
	}
}

func TestRuntimeDeliveryTimingInvalidAgeFailsClosedWithoutClamp(t *testing.T) {
	for _, tc := range []struct {
		name string
		age  float64
	}{{"negative", -0.5}, {"nan", math.NaN()}, {"positive-infinity", math.Inf(1)}, {"negative-infinity", math.Inf(-1)}, {"over-safe-bound", 9007199254740992}} {
		t.Run(tc.name, func(t *testing.T) {
			r := newRuntimeDeliveryRig(t)
			r.historicalRows(t, 1)
			r.useReader(&runtimeDeliveryFault{Pipeline: r.pipe, change: func(v *grouppipeline.DeliveryView) { v.OldestPendingSeconds = tc.age }})
			before := lifecycleAPIState(t, r.f.db)
			response := performJSONRequest(t, r.router, http.MethodGet, runtimeRuleContextPath+"?"+r.query(t).Encode(), nil)
			require.Equal(t, http.StatusServiceUnavailable, response.Code, response.Body.String())
			require.Contains(t, response.Body.String(), "unavailable")
			require.NotContains(t, response.Body.String(), `"rule_evidence"`)
			require.NotContains(t, response.Body.String(), r.f.file)
			require.Equal(t, before, lifecycleAPIState(t, r.f.db))
		})
	}
}

func TestRuntimeDeliveryTimingNoSchedulingEvidenceDoesNotBorrowPerEffectRetry(t *testing.T) {
	r := newRuntimeDeliveryRig(t)
	r.accept(t, "retry-timing")
	key := r.tick(t)
	claim, err := r.store.BeginDelivery(t.Context(), key, "timing-fixture", time.Minute)
	require.NoError(t, err)
	next := r.clock.Add(47*time.Second + 250*time.Millisecond)
	require.NoError(t, r.store.MarkRetry(t.Context(), claim, groupdelivery.StateRetrying, "temporary-unavailable", next))
	row, err := r.store.GetOutbox(t.Context(), key)
	require.NoError(t, err)
	require.Equal(t, next, row.NextRetryAt, "the fixture really has per-effect retry evidence, not a new group authority")
	r.clock = r.clock.Add(500 * time.Millisecond)
	g := runtimeTimingReadOnly(t, r)
	require.EqualValues(t, 1, g.Delivery.Current.Stages.Retrying)
	require.Equal(t, 0.5, *g.Delivery.OldestPendingSeconds)
	require.Nil(t, g.Delivery.Current.NextDueAt)
	require.Nil(t, g.Delivery.Current.NextRetryAt, "a selected outbox effect cannot supply a scoped group deadline")
	require.Nil(t, g.Delivery.Current.AttemptDurationMS)
}

func TestRuntimeDeliveryTimingSourceWindowAndCommittedTimesRetainExactExistingProof(t *testing.T) {
	for _, wide := range []bool{false, true} {
		t.Run(map[bool]string{false: "long", true: "wide"}[wide], func(t *testing.T) {
			r := newRuntimeSourceDeliveryRig(t, wide)
			r.offer(t, "mean", 0, false)
			r.offer(t, "mean", 0, true)
			r.closeWindow(t, 0)
			r.clock = r.clock.Add(3*time.Second + 250*time.Millisecond)
			r.deliver(t)
			committedAt := r.clock
			prior := runtimeDeliveredSource(t, r.evidence(t), "tag-mean")
			r.offer(t, "mean", 1, false)
			r.closeWindow(t, 1)
			r.clock = r.clock.Add(125*time.Second + 500*time.Millisecond)
			before := r.durableState(t)
			g := runtimeTimingReadOnly(t, r.runtimeDeliveryRig)
			require.Equal(t, before, r.durableState(t))
			source := runtimeDeliveredSource(t, g, "tag-mean")
			require.Equal(t, r.snapshot.EffectiveAt.Add(r.window).Format(time.RFC3339Nano), *source.Latest.WindowStart)
			require.Equal(t, r.snapshot.EffectiveAt.Add(2*r.window).Format(time.RFC3339Nano), *source.Latest.WindowEnd)
			require.Nil(t, source.Latest.Receipt)
			require.Equal(t, prior.LastCommitted, source.LastCommitted)
			require.Equal(t, committedAt.Format(time.RFC3339Nano), source.LastCommitted.Receipt.CommittedAt)
			require.Equal(t, committedAt.Format(time.RFC3339Nano), g.Delivery.Current.LastCommit.CommittedAt)
			require.Equal(t, 125.5, *g.Delivery.OldestPendingSeconds)
			require.Nil(t, g.Delivery.TestWrite)
			require.Nil(t, g.Delivery.Current.NextDueAt)
			require.Nil(t, g.Delivery.Current.NextRetryAt)
			require.Nil(t, g.Delivery.Current.AttemptDurationMS)
		})
	}
}

// Setup, Tick, Apply, recovery and Sender work happen before this witness.
// Repeated scoped GETs must preserve SQL, journal, receipts and frozen proof.
func runtimeTimingReadOnly(t *testing.T, r *runtimeDeliveryRig) runtimeview.GroupView {
	t.Helper()
	before := lifecycleAPIState(t, r.f.db)
	var changes int
	require.NoError(t, r.f.db.QueryRow(`SELECT total_changes()`).Scan(&changes))
	g := r.evidence(t)
	require.Equal(t, g.Delivery, r.evidence(t).Delivery, "fixed durable clock preserves the observation; server as-of is independently read")
	var after int
	require.NoError(t, r.f.db.QueryRow(`SELECT total_changes()`).Scan(&after))
	require.Equal(t, changes, after)
	require.Equal(t, before, lifecycleAPIState(t, r.f.db))
	return g
}
