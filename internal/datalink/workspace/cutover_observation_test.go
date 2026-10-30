package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCutoverObservationUsesOriginalSourceAndDurableCutoverID(t *testing.T) {
	h := newCutoverObservationFixture(t, true, true)
	before := h.req.ExpectedWorkspaceRevision
	result, err := h.apply(t)
	require.NoError(t, err)
	require.NotNil(t, result)
	current, err := h.f.service.workspaceSvc.ReadSnapshot(t.Context())
	require.NoError(t, err)
	require.NotEqual(t, before, current.DatabaseSetupRevision)
	require.Len(t, h.observations, 1, "native commit lacked confirmed original-scope observation")
	require.Equal(t, observedCutover{Available: true, WorkspaceID: h.f.group.WorkspaceID, WorkspaceRevision: current.DatabaseSetupRevision, GroupID: h.f.group.ID, GroupRevision: h.f.group.AppliedRevision, ConnectorID: h.f.group.Destination.ConnectorID, ConnectorRevision: h.f.group.Destination.ConnectorRevision, OperationID: result.ID, Outcome: "applied"}, h.observations[0])
	require.Equal(t, 2, h.scopeCalls)
}
func TestCutoverObservationReplayAndRollbackDoNotDuplicateCommit(t *testing.T) {
	h := newCutoverObservationFixture(t, true, true)
	forward, err := h.apply(t)
	require.NoError(t, err)
	replay, err := h.apply(t)
	require.NoError(t, err)
	require.Equal(t, forward.ID, replay.ID)
	require.Len(t, h.observations, 1, "recorded replay emitted again")
	rollback := reviewCounterRollback(t, h.f, forward)
	h.req = rollback
	restored, err := h.apply(t)
	require.NoError(t, err)
	require.Len(t, h.observations, 2)
	require.Equal(t, h.f.group.AppliedRevision, h.observations[1].GroupRevision)
	require.NotEqual(t, restored.Started[0].Revision, h.observations[1].GroupRevision)
	require.Equal(t, restored.ID, h.observations[1].OperationID)
	replay, err = h.apply(t)
	require.NoError(t, err)
	require.Equal(t, restored.ID, replay.ID)
	require.Len(t, h.observations, 2)
}
func TestCutoverObservationFinalizedOverlapWithoutLedgerIsUnavailable(t *testing.T) {
	h := newCutoverObservationFixture(t, true, true)
	store, key, closure := finalizedClosure(t, h.f, h.f.group, h.req.Boundary.Add(time.Second))
	require.NoError(t, store.CommitClosure(t.Context(), key, closure))
	before := finalizedState(t, h.f.db)
	result, err := h.apply(t)
	require.ErrorContains(t, err, "recording-cutover-finalized-overlap")
	require.Nil(t, result)
	require.Equal(t, before, finalizedState(t, h.f.db))
	h.requireUnavailable(t)
	require.Zero(t, h.scopeCalls)
}
func TestCutoverObservationMissingChangedOrForeignProofIsUnavailable(t *testing.T) {
	for _, mode := range []string{"missing-frozen", "corrupt-frozen", "changed-workspace", "foreign-workspace", "missing-ledger", "provider-panic"} {
		t.Run(mode, func(t *testing.T) {
			h := newCutoverObservationFixture(t, true, true)
			h.beforeScope = func(call int) {
				if call != 1 {
					return
				}
				var err error
				switch mode {
				case "missing-frozen":
					_, err = h.f.db.Exec(`DELETE FROM wg_runtime_versions WHERE group_id=? AND group_revision=?`, h.f.group.ID, h.f.group.AppliedRevision)
				case "corrupt-frozen":
					_, err = h.f.db.Exec(`UPDATE wg_runtime_versions SET payload='{}' WHERE group_id=? AND group_revision=?`, h.f.group.ID, h.f.group.AppliedRevision)
				case "changed-workspace":
					h.beforeScope = func(call int) {
						if call == 2 {
							record, e := h.f.service.workspaceSvc.ReadSnapshot(t.Context())
							require.NoError(t, e)
							_, e = h.f.service.workspaceSvc.UpdateDatabaseSetup(t.Context(), record.DatabaseSetupRevision, func(context.Context, *sql.Tx, *Record) error { return nil })
							require.NoError(t, e)
						}
					}
				case "foreign-workspace":
					h.f.group.WorkspaceID = "foreign"
				case "missing-ledger":
					_, err = h.f.db.Exec(`DELETE FROM write_group_cutovers`)
				case "provider-panic":
					panic("password_fixture")
				}
				require.NoError(t, err)
			}
			result, err := h.apply(t)
			require.NoError(t, err)
			require.NotNil(t, result)
			h.requireUnavailable(t)
		})
	}
}
func TestCutoverObservationCallbacksAfterWorkspaceUnlockAndNoExtraWrites(t *testing.T) {
	h := newCutoverObservationFixture(t, true, true)
	var settledChanges int64
	h.beforeScope = func(call int) {
		if call == 1 {
			require.NoError(t, h.f.db.QueryRow(`SELECT total_changes()`).Scan(&settledChanges))
		}
	}
	h.onObserve = func(v observedCutover) {
		var observedChanges int64
		require.NoError(t, h.f.db.QueryRow(`SELECT total_changes()`).Scan(&observedChanges))
		require.Equal(t, settledChanges, observedChanges, "observation performed SQL writes")
		require.True(t, h.f.service.workspaceSvc.mu.TryLock(), "callback under workspace lock")
		h.f.service.workspaceSvc.mu.Unlock()
		var count int
		require.NoError(t, h.f.db.QueryRow(`SELECT COUNT(*) FROM write_group_cutovers WHERE id=?`, v.OperationID).Scan(&count))
		require.Equal(t, 1, count)
	}
	result, err := h.apply(t)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, h.observations, 1)
}
func TestCutoverObservationFilterAndPanicDoNotChangeNativeResult(t *testing.T) {
	for _, mode := range []string{"filtered", "panic"} {
		t.Run(mode, func(t *testing.T) {
			h := newCutoverObservationFixture(t, mode != "filtered", mode != "filtered")
			if mode == "panic" {
				h.onObserve = func(observedCutover) { panic("password_fixture") }
			}
			result, err := h.apply(t)
			require.NoError(t, err)
			require.NotNil(t, result)
			if mode == "filtered" {
				require.Empty(t, h.observations)
				require.Zero(t, h.scopeCalls)
			} else {
				require.Len(t, h.observations, 1)
			}
		})
	}
}

func TestCutoverObservationUnacknowledgedSettlementNeverClaimsApplied(t *testing.T) {
	for _, after := range []bool{false, true} {
		t.Run(fmt.Sprint(after), func(t *testing.T) {
			h := newCutoverObservationFixture(t, true, true)
			repo := &ownershipFenceCommitRepository{SQLRepository: NewSQLRepository(h.f.db), digest: h.req.ReviewDigest, fail: true, after: after}
			h.f.service.workspaceSvc.repo = repo
			result, err := h.apply(t)
			require.Error(t, err)
			require.Nil(t, result)
			h.requireUnavailable(t)
			require.Zero(t, h.scopeCalls)
			var count int
			require.NoError(t, h.f.db.QueryRow(`SELECT COUNT(*) FROM write_group_cutovers`).Scan(&count))
			if after {
				require.Equal(t, 1, count)
				replayed, e := h.apply(t)
				require.NoError(t, e)
				require.NotNil(t, replayed)
				require.Len(t, h.observations, 1, "recorded replay cannot upgrade uncertainty into a new commit event")
			} else {
				require.Zero(t, count)
			}
		})
	}
}

func TestCutoverObservationKeepsOriginalVersionAfterLaterDraft(t *testing.T) {
	h := newCutoverObservationFixture(t, true, true)
	var latest *WriteGroupSaveResult
	h.beforeScope = func(call int) {
		if call != 1 {
			return
		}
		saved, err := h.f.service.Get(t.Context(), h.f.group.ID)
		require.NoError(t, err)
		draft := cloneWriteGroup(saved.Group)
		draft.Name = "password_fixture-later-draft"
		latest, err = h.f.service.Update(t.Context(), draft.ID, WriteGroupMutation{WorkspaceID: draft.WorkspaceID, ExpectedWorkspaceRevision: saved.WorkspaceRevision, ExpectedGroupRevision: draft.Revision, ExpectedConnectorRevision: draft.Destination.ConnectorRevision, Group: draft})
		require.NoError(t, err)
	}
	result, err := h.apply(t)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, h.observations, 1)
	require.NotNil(t, latest)
	require.True(t, h.observations[0].Available)
	require.Equal(t, h.f.group.AppliedRevision, h.observations[0].GroupRevision)
	require.NotEqual(t, latest.Group.Revision, h.observations[0].GroupRevision)
	require.Equal(t, latest.WorkspaceRevision, h.observations[0].WorkspaceRevision)
	raw, err := json.Marshal(h.observations)
	require.NoError(t, err)
	require.NotContains(t, string(raw), "password_fixture")
}

func TestCutoverObservationCaptureLimitAndWarnFilteringPrecedeProofReads(t *testing.T) {
	for _, mode := range []string{"limit", "filtered-limit", "warn-only-success", "warn-filtered-failure"} {
		t.Run(mode, func(t *testing.T) {
			h := newCutoverObservationFixture(t, true, true)
			var dropped uint64
			h.f.service.cutoverObservation.OnDrop = func(n uint64) { dropped += n }
			if mode == "limit" || mode == "filtered-limit" {
				h.f.service.cutoverObservation.RecordLimit = 0
			}
			if mode == "filtered-limit" || mode == "warn-only-success" {
				h.f.service.cutoverObservation.InfoEnabled = false
			}
			if mode == "warn-filtered-failure" {
				h.f.service.cutoverObservation.WarnEnabled = false
				h.req.ReviewDigest = "stale-review"
			}
			result, err := h.apply(t)
			if mode == "warn-filtered-failure" {
				require.Error(t, err)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
			}
			require.Empty(t, h.observations)
			require.Zero(t, h.scopeCalls)
			if mode == "limit" {
				require.EqualValues(t, 1, dropped)
			} else {
				require.Zero(t, dropped, "filtered observations are not loss")
			}
		})
	}
}
func TestCutoverObservationUnsafeScopeAndWrongFrozenIdentityFailClosed(t *testing.T) {
	for _, mode := range []string{"unsafe-scope", "wrong-revision", "wrong-connector", "foreign-descriptor"} {
		t.Run(mode, func(t *testing.T) {
			h := newCutoverObservationFixture(t, true, true)
			if mode == "unsafe-scope" {
				h.f.service.cutoverObservation.ScopeRevision = func(context.Context, string, string, string, string, string) (string, bool) {
					return "postgres://private-host/password_fixture", true
				}
			} else {
				h.beforeScope = func(call int) {
					if call != 1 {
						return
					}
					var raw string
					require.NoError(t, h.f.db.QueryRow(`SELECT payload FROM wg_runtime_versions WHERE group_id=? AND group_revision=?`, h.f.group.ID, h.f.group.AppliedRevision).Scan(&raw))
					if mode == "foreign-descriptor" {
						_, err := h.f.db.Exec(`INSERT INTO wg_runtime_versions(workspace_id,group_id,group_revision,payload)VALUES('foreign',?,?,?)`, h.f.group.ID, h.f.group.AppliedRevision, raw)
						require.NoError(t, err)
						return
					}
					var frozen WriteGroupAppliedSnapshot
					require.NoError(t, json.Unmarshal([]byte(raw), &frozen))
					if mode == "wrong-revision" {
						frozen.AppliedRevision = "other"
					} else {
						frozen.Group.Destination.ConnectorRevision = "other"
					}
					changed, err := json.Marshal(frozen)
					require.NoError(t, err)
					_, err = h.f.db.Exec(`UPDATE wg_runtime_versions SET payload=? WHERE group_id=? AND group_revision=?`, string(changed), h.f.group.ID, h.f.group.AppliedRevision)
					require.NoError(t, err)
				}
			}
			result, err := h.apply(t)
			require.NoError(t, err)
			require.NotNil(t, result)
			h.requireUnavailable(t)
		})
	}
}

func TestCutoverObservationFinalProofRaceCannotClaimApplied(t *testing.T) {
	for _, mode := range []string{"frozen", "frozen-layout", "ledger", "original-version", "prior-ledger"} {
		t.Run(mode, func(t *testing.T) {
			h := newCutoverObservationFixture(t, true, true)
			priorID := ""
			if mode == "prior-ledger" {
				forward, err := h.apply(t)
				require.NoError(t, err)
				priorID = forward.ID
				h.req = reviewCounterRollback(t, h.f, forward)
				h.observations = nil
				h.scopeCalls = 0
			}
			h.beforeScope = func(call int) {
				if call != 2 {
					return
				}
				before, err := h.f.service.workspaceSvc.ReadSnapshot(t.Context())
				require.NoError(t, err)
				switch mode {
				case "frozen":
					_, err = h.f.db.Exec(`UPDATE wg_runtime_versions SET payload='{}' WHERE group_id=? AND group_revision=?`, h.f.group.ID, h.f.group.AppliedRevision)
				case "frozen-layout":
					var raw string
					require.NoError(t, h.f.db.QueryRow(`SELECT payload FROM wg_runtime_versions WHERE group_id=? AND group_revision=?`, h.f.group.ID, h.f.group.AppliedRevision).Scan(&raw))
					var frozen WriteGroupAppliedSnapshot
					require.NoError(t, json.Unmarshal([]byte(raw), &frozen))
					require.NotEmpty(t, frozen.RuntimeLayout.Columns)
					frozen.RuntimeLayout.Columns[0].Nullable = !frozen.RuntimeLayout.Columns[0].Nullable
					changed, e := json.Marshal(frozen)
					require.NoError(t, e)
					_, err = h.f.db.Exec(`UPDATE wg_runtime_versions SET payload=? WHERE group_id=? AND group_revision=?`, string(changed), h.f.group.ID, h.f.group.AppliedRevision)
				case "ledger":
					_, err = h.f.db.Exec(`UPDATE write_group_cutovers SET review_digest=review_digest||'-changed'`)
				case "original-version":
					_, err = h.f.db.Exec(`UPDATE write_group_versions SET payload='{}' WHERE group_id=? AND group_revision=?`, h.f.group.ID, h.f.group.AppliedRevision)
				case "prior-ledger":
					_, err = h.f.db.Exec(`UPDATE write_group_cutovers SET review_digest=review_digest||'-changed' WHERE id=?`, priorID)
				}
				require.NoError(t, err)
				after, err := h.f.service.workspaceSvc.ReadSnapshot(t.Context())
				require.NoError(t, err)
				require.Equal(t, before.DatabaseSetupRevision, after.DatabaseSetupRevision, "test specifically mutates immutable proof without W revision movement")
			}
			result, err := h.apply(t)
			require.NoError(t, err)
			require.NotNil(t, result)
			h.requireUnavailable(t)
		})
	}
}
