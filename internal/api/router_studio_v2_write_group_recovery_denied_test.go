package api

import (
	"os"
	"strings"
	"testing"

	"go-gateway/internal/datalink/recordingplan"

	"github.com/stretchr/testify/require"
)

// TestLongStorageRecoveryDenied lists what must NOT be repaired: the delivery
// stays blocked with a stable code, nothing is created, altered or deleted, and
// existing receipts are untouched.
func TestLongStorageRecoveryDenied(t *testing.T) {
	t.Run("a foreign table took the name", func(t *testing.T) {
		r := newRecoveryRig(t, "sqlite")
		r.start()
		r.record(10.5)
		r.committed(1, "first window")
		receipts := r.receipts()

		r.drop(r.table)
		_, err := r.dest.ExecContext(t.Context(), `CREATE TABLE `+r.quoted(r.table)+` (note TEXT)`)
		require.NoError(t, err)
		_, err = r.dest.ExecContext(t.Context(), `INSERT INTO `+r.quoted(r.table)+` VALUES ('someone else''s data')`)
		require.NoError(t, err)
		r.record(11.5)

		r.waitFor(func() bool { return r.outbox()["blocked"] == 1 }, "the delivery is blocked")
		require.Equal(t, []string{recordingplan.StorageCodeOwnerMismatch}, r.outboxCodes())
		require.Equal(t, 1, r.count(r.table), "the foreign table keeps its row")
		require.Equal(t, receipts, r.receipts())
		var columns []string
		rows, err := r.dest.QueryContext(t.Context(), `SELECT name FROM pragma_table_info(?)`, r.table)
		require.NoError(t, err)
		defer rows.Close()
		for rows.Next() {
			var name string
			require.NoError(t, rows.Scan(&name))
			columns = append(columns, name)
		}
		require.NoError(t, rows.Err())
		require.Equal(t, []string{"note"}, columns, "the foreign table is not altered")
		require.Equal(t, 1, r.lineage().Generation)
		require.False(t, r.lineage().HistoryMayBeLost)
	})

	t.Run("a same-name table lacks the confirmed owner marker", func(t *testing.T) {
		r := newRecoveryRig(t, "sqlite")
		r.start()
		r.record(10.5)
		r.committed(1, "first window")

		r.drop(r.table)
		_, err := r.dest.ExecContext(t.Context(), `CREATE TABLE `+r.quoted(r.table)+` (bucket_start TEXT NOT NULL, record_id TEXT NOT NULL, PRIMARY KEY (bucket_start, record_id))`)
		require.NoError(t, err)
		r.record(11.5)

		r.waitFor(func() bool { return r.outbox()["blocked"] == 1 }, "the incompatible table blocks delivery")
		require.Equal(t, []string{recordingplan.StorageCodeOwnerMismatch}, r.outboxCodes())
		require.Equal(t, 1, r.lineage().Generation)
	})

	t.Run("recreate_missing was turned off before confirmation", func(t *testing.T) {
		off := false
		r := newRecoveryRigWithPolicy(t, "sqlite", &off)
		r.start()
		r.record(10.5)
		r.committed(1, "first window")

		r.drop(r.table)
		r.record(11.5)

		r.waitFor(func() bool { return r.outbox()["blocked"] == 1 }, "disabled recovery blocks delivery without DDL")
		require.Equal(t, []string{recordingplan.StorageCodeRepairBlocked}, r.outboxCodes())
		require.False(t, r.tableExists(r.table), "the dropped table stays absent when recreate_missing is false")
		storage := r.delivery(r.groups[0].id)["storage"].(map[string]any)
		require.Equal(t, false, storage["recreate_missing"])
		require.Equal(t, false, storage["history_may_be_lost"])
	})

	t.Run("the receipts are gone too", func(t *testing.T) {
		r := newRecoveryRig(t, "sqlite")
		r.start()
		r.record(10.5)
		r.committed(1, "first window")

		r.drop(r.table)
		r.drop("gw_effect_receipts")
		r.record(11.5)

		r.waitFor(func() bool { return r.outbox()["blocked"] == 1 }, "the delivery is blocked")
		require.Equal(t, []string{recordingplan.StorageCodeRepairBlocked}, r.outboxCodes())
		require.False(t, r.tableExists(r.table), "no table is created while the receipt history is gone")
		require.False(t, r.tableExists("gw_effect_receipts"), "an empty receipt table is never created to look like no effect was committed")
	})

	t.Run("local confirmation evidence is missing", func(t *testing.T) {
		r := newRecoveryRig(t, "sqlite")
		r.start()
		r.record(10.5)
		r.committed(1, "first window")

		r.drop(r.table)
		_, err := r.f.db.ExecContext(t.Context(), `UPDATE write_group_storage_scopes SET confirmed_operation_id = '', confirmed_layout = ''`)
		require.NoError(t, err)
		r.record(11.5)

		r.waitFor(func() bool { return r.outbox()["blocked"] == 1 }, "unconfirmed scope blocks delivery")
		require.Equal(t, []string{recordingplan.StorageCodeOwnerMismatch}, r.outboxCodes())
		require.False(t, r.tableExists(r.table))
	})

	t.Run("the connector identity changed", func(t *testing.T) {
		r := newRecoveryRig(t, "sqlite")
		r.start()
		r.record(10.5)
		r.committed(1, "first window")

		r.drop(r.table)
		_, err := r.f.db.ExecContext(t.Context(), `UPDATE database_connectors SET identity_revision='connector-2' WHERE id='connector-A'`)
		require.NoError(t, err)
		r.record(11.5)

		r.waitFor(func() bool { return r.outbox()["blocked"] == 1 }, "the delivery is blocked")
		require.Equal(t, []string{recordingplan.StorageCodeRepairBlocked}, r.outboxCodes())
		require.False(t, r.tableExists(r.table))
	})

	t.Run("the SQLite file is gone", func(t *testing.T) {
		r := newRecoveryRig(t, "sqlite")
		r.start()
		r.record(10.5)
		r.committed(1, "first window")

		require.NoError(t, os.Remove(r.f.file))
		r.record(11.5)

		r.waitFor(func() bool { return r.outbox()["retrying"] == 1 }, "the delivery waits")
		require.Equal(t, []string{"target-unavailable"}, r.outboxCodes(), "an outage is reported as such and waits")
		_, err := os.Stat(r.f.file)
		require.True(t, os.IsNotExist(err), "no replacement database file is created: %v", err)
		require.False(t, strings.Contains(strings.Join(r.outboxCodes(), ","), "repair"))
	})
}
