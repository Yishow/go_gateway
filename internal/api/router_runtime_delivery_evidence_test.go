package api

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestRuntimeDeliveryEvidenceCollectingQueuedNeverSQLCommitted(t *testing.T) {
	r := newRuntimeDeliveryRig(t)
	r.accept(t, "accepted-1")
	authority, err := r.pipe.DeliveryForRevision(t.Context(), r.f.group["id"].(string), r.f.group["applied_revision"].(string))
	require.NoError(t, err)
	require.Equal(t, 1, authority.RevisionStages.Stages.Collecting)
	before := lifecycleAPIState(t, r.f.db)
	group := r.evidence(t)
	require.NotNil(t, group.Delivery, "scoped production GET must project the durable delivery authority, not null")
	require.NotNil(t, group.Delivery.Current)
	require.EqualValues(t, 1, group.Delivery.Current.Stages.Collecting)
	require.Zero(t, group.Delivery.Current.Stages.SQLCommitted)
	require.Nil(t, group.Delivery.Current.LastCommit)
	require.Equal(t, before, lifecycleAPIState(t, r.f.db))
	r.tick(t)
	before = lifecycleAPIState(t, r.f.db)
	current := r.evidence(t).Delivery.Current
	require.Zero(t, current.Stages.Collecting)
	require.EqualValues(t, 1, current.Stages.Queued)
	require.Zero(t, current.Stages.SQLCommitted)
	require.Nil(t, current.LastCommit)
	require.Equal(t, before, lifecycleAPIState(t, r.f.db))
}

func TestRuntimeDeliveryEvidenceUsesEffectiveRevisionNotRequestedB(t *testing.T) {
	r := newRuntimeDeliveryRig(t)
	r.accept(t, "accepted-A")
	a := r.f.group["applied_revision"].(string)
	r.clock = time.Date(2026, 10, 2, 12, 0, 17, 0, time.UTC)
	r.update(t, "meter_B", 30)
	b := r.apply(t).Group.AppliedRevision
	r.update(t, "meter_C", 45)
	r.useReader(r.pipe)
	r.clock = time.Date(2026, 10, 2, 12, 0, 29, 0, time.UTC)
	g := r.evidence(t)
	require.Equal(t, a, *g.EffectiveRevision)
	require.Equal(t, b, *g.AppliedRevision)
	require.Equal(t, a, g.Delivery.Current.Binding.EntityRevision)
	require.EqualValues(t, 1, g.Delivery.Current.Stages.Collecting)
	r.clock = r.clock.Add(time.Second)
	g = r.evidence(t)
	require.Equal(t, b, *g.EffectiveRevision)
	require.Equal(t, b, g.Delivery.Current.Binding.EntityRevision)
	require.Zero(t, g.Delivery.Current.Stages.Collecting, "A journal is not B collecting")
	require.Nil(t, g.Delivery.Current.Intake, "A intake proof is not B permission")
}

func TestRuntimeDeliveryEvidenceDisabledHistoricalDestinationNotRelabeled(t *testing.T) {
	r := newRuntimeDeliveryRig(t)
	r.historicalRows(t, 3)
	// Save a real new connector and disabled draft; retained rows own C1/rev7.
	_, err := r.f.db.Exec(`INSERT INTO database_connectors(id,name,kind,connection_config,identity_revision,status,enabled) SELECT 'C2','New draft',kind,connection_config,'8',status,enabled FROM database_connectors WHERE id='connector-A'`)
	require.NoError(t, err)
	_, err = r.f.workspace.BindDatabaseConnector(t.Context(), "C2")
	require.NoError(t, err)
	saved, err := r.f.groups.Get(t.Context(), r.f.group["id"].(string))
	require.NoError(t, err)
	saved.Group.Destination.ConnectorID = "C2"
	saved.Group.Destination.ConnectorRevision = "8"
	saved.Group.Destination.TableName = "new_table"
	updated, err := r.f.groups.Update(t.Context(), saved.Group.ID, workspace.WriteGroupMutation{WorkspaceID: saved.Group.WorkspaceID, ExpectedWorkspaceRevision: saved.WorkspaceRevision, ExpectedGroupRevision: saved.Group.Revision, ExpectedConnectorRevision: "8", Group: saved.Group})
	require.NoError(t, err)
	_, err = r.f.groups.Disable(t.Context(), updated.Group.ID, workspace.WriteGroupMutation{WorkspaceID: updated.Group.WorkspaceID, ExpectedWorkspaceRevision: updated.WorkspaceRevision, ExpectedGroupRevision: updated.Group.Revision, ExpectedConnectorRevision: "8"})
	require.NoError(t, err)
	require.NoError(t, r.pipe.Reconcile(t.Context()))
	g := r.evidence(t)
	require.False(t, g.Saved.Enabled)
	require.Equal(t, "C2", g.Saved.Destination.ConnectorID)
	require.Len(t, g.Delivery.Backlog, 1)
	b := g.Delivery.Backlog[0]
	require.Equal(t, "rev7", b.Binding.EntityRevision)
	require.Equal(t, "C1", b.Destination.ConnectorID)
	require.Equal(t, "7", b.Destination.ConnectorRevision)
	require.Equal(t, "old_table", b.Destination.Table)
	require.EqualValues(t, 3, b.Pending)
	require.Nil(t, b.Destination.Alias)
	require.Equal(t, "unavailable", b.Destination.Kind)
	if g.Delivery.Current.Intake != nil {
		require.NotEqual(t, "active", g.Delivery.Current.Intake.State)
	}
}

func TestRuntimeDeliveryEvidenceHistoricalTypedOwnershipHandoffOverlapReasons(t *testing.T) {
	for _, reason := range []string{"recording-ownership-unavailable", "recording-ownership-conflict", "recording-handoff-proof-conflict", "recording-cutover-finalized-overlap", "password=fixture-secret /private/raw.sql"} {
		t.Run(reason, func(t *testing.T) {
			r := newRuntimeDeliveryRig(t)
			keys := r.historicalRows(t, 1)
			require.NoError(t, r.store.MarkBlocked(t.Context(), keys[0], reason))
			before := lifecycleAPIState(t, r.f.db)
			response := performJSONRequest(t, r.router, http.MethodGet, runtimeRuleContextPath+"?"+r.query(t).Encode(), nil)
			require.Equal(t, http.StatusOK, response.Code, response.Body.String())
			require.NotContains(t, response.Body.String(), "fixture-secret")
			require.NotContains(t, response.Body.String(), "/private/")
			g := r.evidence(t)
			require.Len(t, g.Delivery.Backlog, 1)
			require.Equal(t, []runtimeview.Reason{runtimeview.ParseReason(reason)}, g.Delivery.Backlog[0].Reasons)
			require.Equal(t, before, lifecycleAPIState(t, r.f.db))
		})
	}
}

func TestRuntimeDeliveryEvidenceReplayInputsSeparateFromFinalizedQuarantine(t *testing.T) {
	for _, reason := range []string{"counter-replay-proof-missing", "counter-replay-proof-conflict", "counter-replay-proof-unsupported"} {
		t.Run(reason, func(t *testing.T) {
			r := newRuntimeCounterDeliveryRig(t, reason)
			authority, err := r.pipe.DeliveryForRevision(t.Context(), r.f.group["id"].(string), "")
			require.NoError(t, err)
			require.Len(t, authority.CounterReplayQuarantines, 1)
			require.Equal(t, 3, authority.Stages.Collecting)
			require.Zero(t, authority.Stages.Quarantined)
			before := lifecycleAPIState(t, r.f.db)
			journal := counterReplayAPIJournal(t, r.f.db)
			var descriptor string
			require.NoError(t, r.f.db.QueryRow(`SELECT payload FROM wg_runtime_versions WHERE group_id=?`, r.f.group["id"]).Scan(&descriptor))
			var changes int
			require.NoError(t, r.f.db.QueryRow(`SELECT total_changes()`).Scan(&changes))
			g := r.evidence(t)
			require.NotNil(t, g.Delivery)
			require.Nil(t, g.Delivery.Current, "damaged current proof cannot manufacture effective/current stages")
			require.NotNil(t, g.Delivery.CounterReplayQuarantines)
			q := *g.Delivery.CounterReplayQuarantines
			require.Len(t, q, 1)
			require.Equal(t, r.f.record.ID, q[0].WorkspaceID)
			require.Equal(t, r.f.group["id"], q[0].GroupID)
			require.Equal(t, r.f.group["applied_revision"], q[0].GroupRevision)
			require.EqualValues(t, 3, q[0].PendingSamples)
			require.Equal(t, runtimeview.ParseReason(reason), q[0].Reason)
			require.Equal(t, grouppipeline.CounterReplayRepairAction, q[0].Action)
			require.Equal(t, journal, counterReplayAPIJournal(t, r.f.db))
			require.Equal(t, before, lifecycleAPIState(t, r.f.db))
			var after string
			require.NoError(t, r.f.db.QueryRow(`SELECT payload FROM wg_runtime_versions WHERE group_id=?`, r.f.group["id"]).Scan(&after))
			require.Equal(t, descriptor, after)
			var count int
			require.NoError(t, r.f.db.QueryRow(`SELECT count(*) FROM wg_delivery_outbox WHERE state='quarantined'`).Scan(&count))
			require.Zero(t, count)
			require.NoError(t, r.f.db.QueryRow(`SELECT total_changes()`).Scan(&count))
			require.Equal(t, changes, count)
		})
	}
}

func TestRuntimeDeliveryEvidenceSharedFrozenMemberTotalsNotPerTag(t *testing.T) {
	r := newRuntimeDeliveryRig(t)
	r.sharedGroup(t)
	for i := range 10 {
		r.accept(t, fmt.Sprintf("shared-%d", i))
		r.send(t, r.tick(t))
	}
	opened, err := r.Resolve(t.Context(), groupdelivery.OutboxItem{Destination: groupdelivery.FrozenDestination{ConnectorID: "connector-A", ConnectorRevision: "connector-1"}})
	require.NoError(t, err)
	var targetRows int
	require.NoError(t, opened.DB.QueryRow(`SELECT count(*) FROM meter_shared`).Scan(&targetRows))
	require.Equal(t, 10, targetRows)
	require.NoError(t, opened.Close())
	before := lifecycleAPIState(t, r.f.db)
	g := r.evidence(t)
	require.Equal(t, before, lifecycleAPIState(t, r.f.db))
	require.NotNil(t, g.Delivery.Current)
	c := g.Delivery.Current
	require.True(t, c.Shared)
	require.Equal(t, "group", c.TotalsScope)
	require.EqualValues(t, 10, c.Stages.SQLCommitted)
	require.Equal(t, "unavailable", c.SourceState)
	require.Empty(t, c.Sources)
	require.NotNil(t, c.LastCommit)
	require.Equal(t, "runtime", c.LastCommit.Provenance)
	require.Nil(t, g.Delivery.TestWrite)
	require.Nil(t, c.NextDueAt)
	require.Nil(t, c.NextRetryAt)
	require.Nil(t, c.AttemptDurationMS)
	require.Nil(t, c.Quota)
	require.Nil(t, c.Storage)
	// A later single-device saved draft cannot relabel the frozen shared totals.
	saved, err := r.f.groups.Get(t.Context(), r.f.group["id"].(string))
	require.NoError(t, err)
	saved.Group.Members = saved.Group.Members[:1]
	_, err = r.f.groups.Update(t.Context(), saved.Group.ID, workspace.WriteGroupMutation{WorkspaceID: saved.Group.WorkspaceID, ExpectedWorkspaceRevision: saved.WorkspaceRevision, ExpectedGroupRevision: saved.Group.Revision, ExpectedConnectorRevision: saved.Group.Destination.ConnectorRevision, Group: saved.Group})
	require.NoError(t, err)
	require.True(t, r.evidence(t).Delivery.Current.Shared)
	require.NotEqual(t, before, lifecycleAPIState(t, r.f.db), "the explicit setup save is outside the GET read-only assertion")
}

func TestRuntimeDeliveryEvidenceMissingOrMismatchedCommitProofFailsClosed(t *testing.T) {
	for _, damage := range []string{"receipt_missing", "digest_mismatch", "time_mismatch"} {
		t.Run(damage, func(t *testing.T) {
			r := newRuntimeDeliveryRig(t)
			r.accept(t, "receipt-sample")
			key := r.tick(t)
			r.send(t, key)
			require.EqualValues(t, 1, r.evidence(t).Delivery.Current.Stages.SQLCommitted)
			var err error
			switch damage {
			case "receipt_missing":
				_, err = r.f.db.Exec(`DELETE FROM wg_delivery_receipts`)
			case "digest_mismatch":
				_, err = r.f.db.Exec(`UPDATE wg_delivery_receipts SET payload_digest='mismatch'`)
			case "time_mismatch":
				_, err = r.f.db.Exec(`UPDATE wg_delivery_receipts SET committed_at='2026-10-01T00:00:00Z'`)
			}
			require.NoError(t, err)
			response := performJSONRequest(t, r.router, http.MethodGet, runtimeRuleContextPath+"?"+r.query(t).Encode(), nil)
			require.Equal(t, http.StatusServiceUnavailable, response.Code, response.Body.String())
			require.Contains(t, response.Body.String(), "unavailable")
			require.NotContains(t, response.Body.String(), `"rule_evidence"`)
		})
	}
}

func TestRuntimeDeliveryEvidenceProviderIdentityErrorAndIntakeRaceFailClosed(t *testing.T) {
	for _, tc := range []struct {
		name      string
		code      int
		configure func(*runtimeDeliveryFault)
	}{
		{"failure", 503, func(f *runtimeDeliveryFault) { f.err = errors.New("password=private-secret /private/database.sqlite") }},
		{"nil", 503, func(f *runtimeDeliveryFault) { f.nilView = true }},
		{"foreign_group", 403, func(f *runtimeDeliveryFault) {
			f.change = func(v *grouppipeline.DeliveryView) { v.GroupID = "other-group" }
		}},
		{"stale_revision", 409, func(f *runtimeDeliveryFault) {
			f.change = func(v *grouppipeline.DeliveryView) { v.RevisionStages.GroupRevision = "other-revision" }
		}},
		{"intake_change", 409, func(f *runtimeDeliveryFault) {
			f.statuses = func(n int, s []grouppipeline.GroupStatus) []grouppipeline.GroupStatus {
				if n > 1 {
					s[0].State = grouppipeline.StateBlocked
				}
				return s
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := newRuntimeDeliveryRig(t)
			f := &runtimeDeliveryFault{Pipeline: r.pipe}
			tc.configure(f)
			r.useReader(f)
			before := lifecycleAPIState(t, r.f.db)
			response := performJSONRequest(t, r.router, http.MethodGet, runtimeRuleContextPath+"?"+r.query(t).Encode(), nil)
			require.Equal(t, tc.code, response.Code, response.Body.String())
			require.NotContains(t, response.Body.String(), "private-secret")
			require.NotContains(t, response.Body.String(), "/private/")
			require.NotContains(t, response.Body.String(), `"rule_evidence"`)
			require.Equal(t, before, lifecycleAPIState(t, r.f.db))
		})
	}
}

func TestRuntimeDeliveryEvidenceEmptyRevisionKeepsRetainedBacklogWithoutCurrentZeros(t *testing.T) {
	r := newRuntimeDeliveryRig(t)
	r.historicalRows(t, 3)
	_, err := r.f.db.Exec(`DELETE FROM wg_runtime_versions`)
	require.NoError(t, err)
	g := r.evidence(t)
	require.Nil(t, g.Effective)
	require.NotNil(t, g.Delivery)
	require.Nil(t, g.Delivery.Current)
	require.Len(t, g.Delivery.Backlog, 1)
	require.EqualValues(t, 3, g.Delivery.Backlog[0].Pending)
}

func TestRuntimeDeliveryEvidenceNoWritesOrWorkspaceRecreation(t *testing.T) {
	for _, mode := range []string{"unchanged", "deleted_before", "deleted_during"} {
		t.Run(mode, func(t *testing.T) {
			r := newRuntimeDeliveryRig(t)
			r.accept(t, "read-only")
			query := r.query(t)
			drop := func() {
				_, err := r.f.db.Exec(`DELETE FROM system_settings WHERE key='studio_v2_workspace'`)
				require.NoError(t, err)
			}
			if mode == "deleted_before" {
				drop()
			}
			if mode == "deleted_during" {
				r.useReader(&runtimeDeliveryFault{Pipeline: r.pipe, after: drop})
			}
			journal := counterReplayAPIJournal(t, r.f.db)
			var before int
			require.NoError(t, r.f.db.QueryRow(`SELECT total_changes()`).Scan(&before))
			response := performJSONRequest(t, r.router, http.MethodGet, runtimeRuleContextPath+"?"+query.Encode(), nil)
			var after int
			require.NoError(t, r.f.db.QueryRow(`SELECT total_changes()`).Scan(&after))
			if mode == "unchanged" {
				require.Equal(t, 200, response.Code, response.Body.String())
				require.Equal(t, before, after)
			} else {
				var count int
				require.NoError(t, r.f.db.QueryRow(`SELECT count(*) FROM system_settings WHERE key='studio_v2_workspace'`).Scan(&count))
				require.Zero(t, count, "GET must not recreate a deleted workspace")
				if mode == "deleted_during" {
					require.Equal(t, 409, response.Code, response.Body.String())
				} else {
					require.Equal(t, 503, response.Code, response.Body.String())
				}
				if mode == "deleted_during" {
					require.Equal(t, before+1, after)
				} else {
					require.Equal(t, before, after)
				}
			}
			require.Equal(t, journal, counterReplayAPIJournal(t, r.f.db))
		})
	}
}

func TestRuntimeDeliveryEvidenceMidreadSavedRevisionChangeIsStale(t *testing.T) {
	r := newRuntimeDeliveryRig(t)
	query := r.query(t)
	r.useReader(&runtimeDeliveryFault{Pipeline: r.pipe, after: func() { r.update(t, "changed_during_delivery", 30) }})
	response := performJSONRequest(t, r.router, http.MethodGet, runtimeRuleContextPath+"?"+query.Encode(), nil)
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
	require.Contains(t, response.Body.String(), "stale-reference")
	require.NotContains(t, response.Body.String(), `"rule_evidence"`)
}
