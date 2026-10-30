package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/workspace"
	"go-gateway/internal/diagnostics"

	"github.com/stretchr/testify/require"
)

func TestProductionIntakeEventsCurrentAppliedIsNotSQLCommit(t *testing.T) {
	f := newIntakeEventFixture(t)
	f.reconcileCurrent(t)
	require.Equal(t, groupdelivery.StatePending, f.state(t).State)
	events := deliveryEventsOfCode(t, f.broker, "domain.group.applied")
	require.Len(t, events, 1, "actual eligible installed revision must reach production intake diagnostics")
	fields := events[0].Fields
	require.Equal(t, f.record.ID, fields["workspace_id"])
	require.Equal(t, f.record.DatabaseSetupRevision, fields["workspace_revision"])
	require.Equal(t, f.item.Key.GroupID, fields["group_id"])
	require.Equal(t, f.item.Key.GroupRevision, fields["entity_revision"])
	require.Equal(t, f.item.Destination.ConnectorID, fields["connector_id"])
	require.Equal(t, f.item.Destination.ConnectorRevision, fields["connector_revision"])
	require.Equal(t, "collecting", fields["stage"])
	require.Equal(t, "active", fields["intake_state"])
	require.Equal(t, "intake-status", fields["evidence_kind"])
	require.NotEmpty(t, fields["reference_id"])
	require.NotContains(t, fields, "effect_key")
	require.NotContains(t, fields, "committed_at")
	require.Empty(t, deliveryEventsOfCode(t, f.broker, "domain.delivery.first-confirmed"))
	for range 3 {
		require.NoError(t, f.env.services.groupPipe.Reconcile(t.Context()))
	}
	require.Len(t, deliveryEventsOfCode(t, f.broker, "domain.group.applied"), 1, "unchanged state does not emit per-Reconcile success spam")
}

func TestProductionIntakeEventsDisableKeepsHistoricalBacklog(t *testing.T) {
	f := newIntakeEventFixture(t)
	f.reconcileCurrent(t)
	disabled := f.disable(t)
	require.NoError(t, f.env.services.groupPipe.Reconcile(t.Context()))
	events := deliveryEventsOfCode(t, f.broker, "domain.group.paused")
	require.Len(t, events, 1)
	fields := events[0].Fields
	require.Equal(t, f.item.Key.GroupRevision, fields["entity_revision"], "paused evidence retains original applied revision")
	require.Equal(t, disabled.WorkspaceRevision, fields["workspace_revision"], "current observed workspace scope is not historical output revision")
	require.Equal(t, "not_running", fields["intake_state"])
	require.Equal(t, "blocked", fields["stage"])
	require.Equal(t, groupdelivery.StatePending, f.state(t).State)
	require.Equal(t, f.item.EffectKey, f.state(t).EffectKey)
	require.Empty(t, deliveryEventsOfCode(t, f.broker, "domain.delivery.first-confirmed"))
}

func TestProductionIntakeEventsFutureGroupIsNotCurrentApplied(t *testing.T) {
	f := newIntakeEventFixture(t)
	future := time.Now().UTC().Add(time.Minute)
	f.env.services.writeGroups.WithClock(func() time.Time { return future })
	group := f.env.createAndApply(t, "B")
	f.reconcileCurrent(t)
	for _, e := range deliveryEventsOfCode(t, f.broker, "domain.group.applied") {
		require.NotEqual(t, group.ID, e.Fields["group_id"])
	}
	require.Equal(t, groupdelivery.StatePending, f.state(t).State)
}

func TestProductionQuotaEventsPreserveActualGlobalAndGroupRefusal(t *testing.T) {
	for _, scope := range []groupdelivery.QuotaScope{groupdelivery.QuotaScopeGlobal, groupdelivery.QuotaScopeGroup} {
		t.Run(string(scope), func(t *testing.T) {
			f := newIntakeEventFixture(t)
			clock := intakeEventClock()
			config := groupdelivery.QuotaConfig{GlobalMaxBytes: 1}
			if scope == groupdelivery.QuotaScopeGroup {
				config = groupdelivery.QuotaConfig{GroupMaxBytes: 1}
			}
			pipe, _ := newObservedQuotaPipeline(t, f, f.broker, config, clock)
			quota := quotaEventRefusal(t, f, pipe, clock)
			require.Equal(t, scope, quota.Scope)
			event := waitIntakeEvent(t, f.broker, "domain.quota.refused")
			require.Equal(t, string(scope), event.Fields["quota_scope"])
			require.Equal(t, quota.UsedBytes, event.Fields["quota_used_bytes"])
			require.Equal(t, quota.MaxBytes, event.Fields["quota_max_bytes"])
			require.Equal(t, true, event.Fields["intake_refused"])
			require.Equal(t, "hard_limit", event.Fields["quota_state"])
			require.Equal(t, f.item.Key.GroupRevision, event.Fields["entity_revision"])
			require.Equal(t, f.item.Destination.ConnectorID, event.Fields["connector_id"])
			require.NotContains(t, event.Fields, "effect_key")
			require.Equal(t, groupdelivery.StatePending, f.state(t).State)
			// Ignored foreign samples return nil natively, but cannot prove recovery.
			ignored := f.env.envelope(f.group, 0, "ignored", clock.now(), 1.0)
			ignored.TagID = "not-owned"
			require.NoError(t, pipe.AcceptSample(t.Context(), ignored))
			require.Empty(t, deliveryEventsOfCode(t, f.broker, "domain.quota.recovered"))
		})
	}
}

func TestProductionQuotaEventsRecoveryUsesActualExactScopeAfterCommit(t *testing.T) {
	f := newIntakeEventFixture(t)
	clock := intakeEventClock()
	var callback atomic.Bool
	pipe, store := newObservedQuotaPipeline(t, f, f.broker, groupdelivery.QuotaConfig{GlobalMaxBytes: 1}, clock, func(cfg *grouppipeline.IntakeQuotaObservationConfig) {
		take := cfg.TakeQuotaRecoveryHint
		cfg.TakeQuotaRecoveryHint = func(ctx context.Context) (grouppipeline.QuotaRecoveryHint, bool) {
			hint, ok := take(ctx)
			if !ok {
				return hint, false
			}
			observe := hint.Observe
			hint.Observe = func(o grouppipeline.IntakeQuotaObservation) {
				// Inject only callback faults; real native proof and actual capacity remain
				// unchanged. Invalid identities must not consume the retained exact token.
				for _, change := range []func(*grouppipeline.IntakeQuotaObservation){
					func(v *grouppipeline.IntakeQuotaObservation) { v.Key.GroupID = "foreign" },
					func(v *grouppipeline.IntakeQuotaObservation) { v.WorkspaceRevision = "later" },
					func(v *grouppipeline.IntakeQuotaObservation) { v.ConnectorRevision = "later" },
					func(v *grouppipeline.IntakeQuotaObservation) { v.ObservationID = "new-reference" },
					func(v *grouppipeline.IntakeQuotaObservation) { v.Quota.Scope = groupdelivery.QuotaScopeGroup },
				} {
					bad := o
					change(&bad)
					observe(bad)
				}
				callback.Store(true)
				observe(o)
			}
			return hint, true
		}
	})
	quotaEventRefusal(t, f, pipe, clock)
	fault := waitIntakeEvent(t, f.broker, "domain.quota.refused")
	require.NoError(t, pipe.Start(t.Context()))
	f.waitCommitted(t)
	current, err := store.QuotaStatusForScope(t.Context(), f.item.Key, groupdelivery.QuotaScopeGlobal)
	require.NoError(t, err)
	require.Zero(t, current.UsedBytes)
	// Real broker heartbeat, no new clock/timer authority: only existing Reconcile
	// obtains one due retained-fault hint and measures this exact scope.
	require.Eventually(t, func() bool {
		require.NoError(t, pipe.Reconcile(t.Context()))
		return len(deliveryEventsOfCode(t, f.broker, "domain.quota.recovered")) == 1
	}, diagnostics.RuntimeLogHeartbeat+5*time.Second, 100*time.Millisecond)
	require.True(t, callback.Load())
	event := waitIntakeEvent(t, f.broker, "domain.quota.recovered")
	require.Equal(t, fault.Fields["reference_id"], event.Fields["reference_id"])
	require.Equal(t, "global", event.Fields["quota_scope"])
	require.Equal(t, int64(0), event.Fields["quota_used_bytes"])
	require.Equal(t, false, event.Fields["intake_refused"])
	require.Equal(t, "accepted", event.Fields["stage"])
	require.NotContains(t, event.Fields, "intake_state")
	require.NotContains(t, event.Fields, "effect_key")
	require.Contains(t, event.Message, "capacity")
}

func TestProductionIntakeEventsMissingObservedRevisionIsUnavailable(t *testing.T) {
	f := newIntakeEventFixture(t)
	f.record.DatabaseSetupRevision = ""
	require.NoError(t, workspace.NewSQLRepository(f.env.db).Save(t.Context(), f.record))
	f.reconcileCurrent(t)
	require.Empty(t, deliveryEventsOfCode(t, f.broker, "domain.group.applied"))
	event := waitIntakeEvent(t, f.broker, "domain.unavailable")
	require.Empty(t, event.Fields)
	require.Equal(t, groupdelivery.StatePending, f.state(t).State)
}

func TestIntakeQuotaEventObserverNilAndLevelConfiguration(t *testing.T) {
	zero := newIntakeQuotaEventObserver(nil, nil)
	require.Nil(t, zero.Observe)
	require.Nil(t, zero.ScopeRevision)
	require.Nil(t, zero.TakeQuotaRecoveryHint)
	require.Nil(t, zero.OnDrop)
	for _, level := range []string{"error", "warn", "info"} {
		t.Run(level, func(t *testing.T) {
			broker, err := diagnostics.New(diagnostics.Options{Level: level})
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, broker.Close(context.Background())) })
			cfg := newIntakeQuotaEventObserver(broker, func(context.Context, groupdelivery.GroupKey) (string, bool) {
				t.Fatal("no hint must not read scope")
				return "", false
			})
			if level == "error" {
				require.Nil(t, cfg.Observe)
				require.Nil(t, cfg.ScopeRevision)
				require.Nil(t, cfg.TakeQuotaRecoveryHint)
				require.Nil(t, cfg.OnDrop)
				require.Zero(t, cfg.RecordLimit)
				return
			}
			require.Equal(t, diagnostics.AdmissionRecords, cfg.RecordLimit)
			require.Equal(t, level == "info", cfg.IntakeEnabled)
			require.True(t, cfg.QuotaEnabled)
			_, ok := cfg.TakeQuotaRecoveryHint(t.Context())
			require.False(t, ok)
			records := collectorEventRecords(t, broker)
			require.Empty(t, records)
		})
	}
}

func TestIntakeQuotaEventObserverRejectsUnsafeUnknownOrContradictoryProof(t *testing.T) {
	cases := map[string]func(*grouppipeline.IntakeQuotaObservation){
		"unavailable":    func(o *grouppipeline.IntakeQuotaObservation) { o.Available = false },
		"missing-scope":  func(o *grouppipeline.IntakeQuotaObservation) { o.WorkspaceRevision = "" },
		"unsafe":         func(o *grouppipeline.IntakeQuotaObservation) { o.Key.GroupID = "password=secret" },
		"unknown-kind":   func(o *grouppipeline.IntakeQuotaObservation) { o.Kind = "future" },
		"state":          func(o *grouppipeline.IntakeQuotaObservation) { o.IntakeState = "not_running" },
		"no-reference":   func(o *grouppipeline.IntakeQuotaObservation) { o.ObservationID = "" },
		"mixed-evidence": func(o *grouppipeline.IntakeQuotaObservation) { o.Quota = groupdelivery.QuotaStatus{Configured: true} },
		"unsolicited-recovery": func(o *grouppipeline.IntakeQuotaObservation) {
			o.Kind = grouppipeline.IntakeQuotaRecovered
			o.IntakeState = ""
			o.Quota = groupdelivery.QuotaStatus{Configured: true, Scope: groupdelivery.QuotaScopeGlobal, State: groupdelivery.QuotaOK, MaxBytes: 100}
		},
		"disk-full-max-unknown": func(o *grouppipeline.IntakeQuotaObservation) {
			o.Kind = grouppipeline.IntakeQuotaRefused
			o.IntakeState = ""
			o.Quota = groupdelivery.QuotaStatus{Configured: true, Scope: groupdelivery.QuotaScopeGlobal, State: groupdelivery.QuotaHard, IntakeRefused: true}
		},
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			broker, err := diagnostics.New(diagnostics.Options{})
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, broker.Close(context.Background())) })
			o := grouppipeline.IntakeQuotaObservation{Available: true, Key: groupdelivery.GroupKey{WorkspaceID: "W", GroupID: "G", GroupRevision: "R1"}, WorkspaceRevision: "W1", ConnectorID: "C", ConnectorRevision: "C1", ObservationID: "O", Kind: grouppipeline.IntakeGroupApplied, IntakeState: "active"}
			mutate(&o)
			newIntakeQuotaEventObserver(broker, nil).Observe(o)
			records := collectorEventRecords(t, broker)
			require.Len(t, records, 1)
			require.Equal(t, "domain.unavailable", records[0].Code)
			require.Empty(t, records[0].Fields)
		})
	}
}

func TestProductionIntakeQuotaEventsSinksNeverExposeLabelsOrNotice(t *testing.T) {
	f := newIntakeEventFixture(t, true)
	_, err := f.env.db.ExecContext(t.Context(), "UPDATE write_groups SET name=? WHERE id=?", "password_fixture nested_secret private_host", f.group.ID)
	require.NoError(t, err)
	sub, err := f.broker.Subscribe(diagnostics.Query{}, "")
	require.NoError(t, err)
	defer sub.Close()
	clock := intakeEventClock()
	pipe, _ := newObservedQuotaPipeline(t, f, f.broker, groupdelivery.QuotaConfig{GlobalMaxBytes: 1}, clock, func(cfg *grouppipeline.IntakeQuotaObservationConfig) {
		observe := cfg.Observe
		cfg.Observe = func(o grouppipeline.IntakeQuotaObservation) {
			if o.Kind == grouppipeline.IntakeQuotaRefused {
				o.Quota.LossRiskNotice = "Server=private_host;Password=password_fixture nested_secret"
			}
			observe(o)
		}
	})
	quotaEventRefusal(t, f, pipe, clock)
	event := waitIntakeEvent(t, f.broker, "domain.group.applied")
	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()
	live, err := sub.Next(ctx)
	require.NoError(t, err)
	require.Equal(t, event, live)
	sub.Ack(live.Cursor())
	records := collectorEventRecords(t, f.broker)
	raw, err := json.Marshal(records)
	require.NoError(t, err)
	fileBytes, err := os.ReadFile(f.logPath)
	require.NoError(t, err)
	require.Contains(t, string(fileBytes), `"code":"domain.group.applied"`)
	require.Contains(t, string(fileBytes), `"code":"domain.quota.refused"`)
	for _, secret := range []string{"password_fixture", "nested_secret", "private_host", "connection_config", "loss_risk_notice"} {
		require.NotContains(t, string(raw), secret)
		require.NotContains(t, string(fileBytes), secret)
	}
}

func TestProductionIntakeQuotaFullOrClosedCaptureKeepsNativeOutcome(t *testing.T) {
	for _, mode := range []string{"full", "closed"} {
		t.Run(mode, func(t *testing.T) {
			f := newIntakeEventFixture(t)
			clock := intakeEventClock()
			pipe, _ := newObservedQuotaPipeline(t, f, f.broker, groupdelivery.QuotaConfig{GlobalMaxBytes: 1}, clock)
			var dropped uint64
			if mode == "closed" {
				require.NoError(t, f.broker.Close(t.Context()))
			} else {
				scope := runtimeview.Scope{WorkspaceID: f.record.ID, WorkspaceRevision: f.record.DatabaseSetupRevision, DeviceID: "device-1"}
				for i := range diagnostics.AdmissionRecords {
					id := fmt.Sprintf("retained-intake-%d", i)
					ref := diagnostics.DomainEventRef{ContractVersion: 1, Code: "domain.source.rejected", Binding: runtimeview.Binding{WorkspaceID: scope.WorkspaceID, WorkspaceRevision: scope.WorkspaceRevision, EntityRevision: "R1"}, Locator: runtimeview.EntityLocator{Step: 2, WorkspaceID: scope.WorkspaceID, DeviceID: &scope.DeviceID, RuleID: &id}, Stage: "blocked", Reason: "unavailable", Action: "refresh", Evidence: diagnostics.DomainEvidenceRef{Kind: "source-runtime-transition", ReferenceID: id}}
					raw, err := json.Marshal(ref)
					require.NoError(t, err)
					f.broker.EmitDomain(t.Context(), raw, scope)
				}
				require.NoError(t, f.broker.Flush(t.Context()))
				view, err := f.broker.Snapshot(diagnostics.Query{})
				require.NoError(t, err)
				dropped = view.CaptureDropped
			}
			quota := quotaEventRefusal(t, f, pipe, clock)
			require.Equal(t, groupdelivery.ReasonQuotaHardLimit, quota.Reason)
			require.Equal(t, groupdelivery.StatePending, f.state(t).State)
			f.disable(t)
			require.NoError(t, pipe.Reconcile(t.Context()))
			require.Equal(t, groupdelivery.StatePending, f.state(t).State)
			if mode == "full" {
				require.NoError(t, f.broker.Flush(t.Context()))
				view, err := f.broker.Snapshot(diagnostics.Query{})
				require.NoError(t, err)
				require.Greater(t, view.CaptureDropped, dropped)
			}
		})
	}
}

func TestProductionIntakeEventsFutureBoundaryKeepsRevisionsSeparate(t *testing.T) {
	f := newIntakeEventFixture(t)
	original, err := f.env.services.writeGroups.ResolveAppliedAt(t.Context(), f.group.ID, time.Now().UTC())
	require.NoError(t, err)
	clock := &testClock{}
	clock.set(original.EffectiveAt.Add(time.Second))
	pipe, _ := newObservedQuotaPipeline(t, f, f.broker, groupdelivery.QuotaConfig{GlobalMaxBytes: 1 << 20}, clock)
	current, err := f.env.services.writeGroups.Get(t.Context(), f.group.ID)
	require.NoError(t, err)
	current.Group.Name = "second revision"
	current.Group.RowPolicy.IntervalSeconds = 20
	f.env.services.writeGroups.WithClock(clock.now)
	saved, err := f.env.services.writeGroups.Update(t.Context(), f.group.ID, workspace.WriteGroupMutation{WorkspaceID: f.record.ID, ExpectedWorkspaceRevision: f.record.DatabaseSetupRevision, ExpectedGroupRevision: current.Group.Revision, ExpectedConnectorRevision: current.Group.Destination.ConnectorRevision, Group: current.Group})
	require.NoError(t, err)
	applied, err := f.env.services.writeGroups.Apply(t.Context(), f.group.ID, workspace.WriteGroupMutation{WorkspaceID: f.record.ID, ExpectedWorkspaceRevision: saved.WorkspaceRevision, ExpectedGroupRevision: saved.Group.Revision, ExpectedConnectorRevision: saved.Group.Destination.ConnectorRevision})
	require.NoError(t, err)
	require.NotEqual(t, original.AppliedRevision, applied.Group.AppliedRevision)
	require.NoError(t, pipe.Reconcile(t.Context()))
	futureInstalled := false
	for _, status := range pipe.Status() {
		if status.GroupID == f.group.ID && status.Revision == applied.Group.AppliedRevision {
			futureInstalled = true
		}
	}
	require.True(t, futureInstalled, "native lookahead may install a future boundary without making it intake now")
	for _, event := range deliveryEventsOfCode(t, f.broker, "domain.group.applied") {
		require.NotEqual(t, applied.Group.AppliedRevision, event.Fields["entity_revision"])
	}
	retiring := waitIntakeEvent(t, f.broker, "domain.group.retiring")
	require.Equal(t, original.AppliedRevision, retiring.Fields["entity_revision"])
	next, err := f.env.services.writeGroups.ResolveAppliedAt(t.Context(), f.group.ID, original.EffectiveAt.Add(time.Hour))
	require.NoError(t, err)
	require.Equal(t, applied.Group.AppliedRevision, next.AppliedRevision)
	clock.set(next.EffectiveAt.Add(time.Second))
	require.NoError(t, pipe.Reconcile(t.Context()))
	found := false
	for _, event := range deliveryEventsOfCode(t, f.broker, "domain.group.applied") {
		if event.Fields["entity_revision"] == applied.Group.AppliedRevision {
			found = true
			require.Equal(t, applied.WorkspaceRevision, event.Fields["workspace_revision"])
		}
	}
	require.True(t, found, "future revision becomes applied only after its effective boundary")
	require.Equal(t, groupdelivery.StatePending, f.state(t).State)
}

func TestProductionIntakeQuotaObserverPanicDoesNotChangeNativeOutcomes(t *testing.T) {
	f := newIntakeEventFixture(t)
	clock := intakeEventClock()
	pipe, _ := newObservedQuotaPipeline(t, f, f.broker, groupdelivery.QuotaConfig{GlobalMaxBytes: 1}, clock, func(cfg *grouppipeline.IntakeQuotaObservationConfig) {
		cfg.Observe = func(grouppipeline.IntakeQuotaObservation) { panic("private diagnostic callback") }
	})
	require.NotEmpty(t, pipe.Status())
	quota := quotaEventRefusal(t, f, pipe, clock)
	require.Equal(t, groupdelivery.ReasonQuotaHardLimit, quota.Reason)
	require.Equal(t, groupdelivery.StatePending, f.state(t).State)
}
