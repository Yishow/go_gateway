package groupdelivery

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLifecyclePendingKeepsEveryRevisionAndFrozenDestination(t *testing.T) {
	store, db, _ := newTestStore(t)
	next := testKey
	next.GroupRevision = "rev-2"
	destination := testDestination
	destination.TableName = "new_table"
	destination.Scope = "new-scope"
	enqueueFor(t, store, testKey, testDestination, "old", t0, 1)
	enqueueFor(t, store, next, destination, "new", t0, 1)
	tx, err := db.BeginTx(t.Context(), &sql.TxOptions{ReadOnly: true})
	require.NoError(t, err)
	defer tx.Rollback()
	revisions, err := LifecyclePendingRevisions(t.Context(), tx, testKey.WorkspaceID, testKey.GroupID)
	require.NoError(t, err)
	require.Equal(t, []string{"rev-1", "rev-2"}, revisions)
	for _, item := range []struct {
		key         GroupKey
		destination FrozenDestination
	}{{testKey, testDestination}, {next, destination}} {
		got, err := ReadLifecyclePending(t.Context(), tx, item.key, item.destination)
		require.NoError(t, err)
		require.Equal(t, 1, got.Outbox.Total)
	}
	_, err = ReadLifecyclePending(t.Context(), tx, testKey, destination)
	require.ErrorIs(t, err, ErrOwnershipUnavailable, "current target must never relabel historical rows")
}

func TestLifecyclePendingOutboxStatesAreSubsetsNotAdditionalRows(t *testing.T) {
	store, db, _ := newTestStore(t)
	states := []string{StatePending, StateSending, StateRetrying, StateBlocked, StateUnknown, StateQuarantined, StateCommitted, StateSkipped}
	for i, state := range states {
		entity := fmt.Sprintf("e%d", i)
		enqueueFor(t, store, testKey, testDestination, entity, at(i*10), 1)
		_, err := db.Exec(`UPDATE wg_delivery_outbox SET state=? WHERE entity_key=?`, state, entity)
		require.NoError(t, err)
	}
	got, err := lifecycleRead(t, db, testKey, testDestination)
	require.NoError(t, err)
	require.Equal(t, LifecycleOutboxCounts{Total: 8, Outstanding: 6, Queued: 2, Retrying: 1, Blocked: 1, Unknown: 1, Quarantined: 1, SQLCommitted: 1, Skipped: 1}, got.Outbox)
}

func TestLifecyclePendingUsesCallerSnapshot(t *testing.T) {
	store, db, _ := newTestStore(t)
	_, err := db.Exec(`PRAGMA journal_mode=WAL`)
	require.NoError(t, err)
	enqueueFor(t, store, testKey, testDestination, "e", t0, 1)
	db.SetMaxOpenConns(2)
	tx, err := db.BeginTx(t.Context(), &sql.TxOptions{ReadOnly: true})
	require.NoError(t, err)
	defer tx.Rollback()
	first, err := ReadLifecyclePending(t.Context(), tx, testKey, testDestination)
	require.NoError(t, err)
	_, err = db.Exec(`UPDATE wg_delivery_outbox SET state='unknown'`)
	require.NoError(t, err)
	second, err := ReadLifecyclePending(t.Context(), tx, testKey, testDestination)
	require.NoError(t, err)
	require.Equal(t, first, second)
	require.NoError(t, tx.Rollback())
	current, err := lifecycleRead(t, db, testKey, testDestination)
	require.NoError(t, err)
	require.Equal(t, 1, current.Outbox.Unknown)
	require.Zero(t, current.Outbox.Queued)
}

func TestLifecyclePendingRejectsMissingTransactionAndStorage(t *testing.T) {
	_, err := ReadLifecyclePending(t.Context(), nil, testKey, testDestination)
	require.ErrorIs(t, err, ErrOwnershipUnavailable)
	_, err = LifecyclePendingRevisions(t.Context(), nil, testKey.WorkspaceID, testKey.GroupID)
	require.ErrorIs(t, err, ErrOwnershipUnavailable)
	_, db, _ := newTestStore(t)
	tx, err := db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	require.NoError(t, tx.Rollback())
	_, err = ReadLifecyclePending(t.Context(), tx, testKey, testDestination)
	require.ErrorIs(t, err, ErrOwnershipUnavailable)
	require.ErrorIs(t, err, sql.ErrTxDone)
	require.Equal(t, ErrOwnershipUnavailable.Error(), err.Error())
}

func TestLifecyclePendingProjectionExcludesPrivateFields(t *testing.T) {
	store, db, _ := newTestStore(t)
	destination := testDestination
	destination.Database = "/private/password-secret.db"
	enqueueFor(t, store, testKey, destination, "e", t0, 1)
	_, err := db.Exec(`UPDATE wg_delivery_outbox SET last_error_code='private-raw-exception'`)
	require.NoError(t, err)
	got, err := lifecycleRead(t, db, testKey, destination)
	require.NoError(t, err)
	raw, err := json.Marshal(got)
	require.NoError(t, err)
	for _, secret := range []string{"password-secret", "private-raw-exception", "database", "payload", "connector"} {
		require.NotContains(t, string(raw), secret)
	}
}

func TestLifecyclePendingHandoffJSONUsesAllowlistedSnakeCaseKeys(t *testing.T) {
	raw, err := json.Marshal(LifecycleHandoff{Origin: testKey, Source: testKey, Target: otherKey, State: "pending", Samples: 1})
	require.NoError(t, err)
	var decoded map[string]any
	require.NoError(t, json.Unmarshal(raw, &decoded))
	require.Equal(t, map[string]any{"workspace_id": "workspace-A", "group_id": "group-G", "group_revision": "rev-1"}, decoded["origin"])
	require.NotContains(t, string(raw), "WorkspaceID")
}
