package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"go-gateway/internal/datalink/recordingplan"

	"github.com/stretchr/testify/require"
)

// columns returns the column names of a PostgreSQL table in the rig's schema.
func (r *recoveryRig) columns(table string) []string {
	r.t.Helper()
	rows, err := r.dest.QueryContext(r.t.Context(), `SELECT column_name FROM information_schema.columns
		WHERE table_schema = $1 AND table_name = $2 ORDER BY ordinal_position`, r.f.namespace, table)
	require.NoError(r.t, err)
	defer rows.Close()
	var names []string
	for rows.Next() {
		var name string
		require.NoError(r.t, rows.Scan(&name))
		names = append(names, name)
	}
	require.NoError(r.t, rows.Err())
	return names
}

// repairOperations counts the restore operations that ran DDL.
func (r *recoveryRig) repairOperations() int {
	r.t.Helper()
	var n int
	require.NoError(r.t, r.f.db.QueryRowContext(r.t.Context(),
		`SELECT COUNT(*) FROM managed_schema_operations WHERE action = ? AND executed_statements > 0`,
		recordingplan.StorageRepairAction).Scan(&n))
	return n
}

// limitedRole replaces the connector's login with a role that may read and
// insert into the destination schema but not create anything in it. The role is
// the fixture's own and is dropped with it.
func (r *recoveryRig) limitedRole() {
	r.t.Helper()
	secret := make([]byte, 16)
	_, err := rand.Read(secret)
	require.NoError(r.t, err)
	role, password := fmt.Sprintf("gw_norestore_%d", time.Now().UnixNano()), hex.EncodeToString(secret)
	schema := `"` + r.f.namespace + `"`
	for _, statement := range []string{
		`CREATE ROLE "` + role + `" LOGIN PASSWORD '` + password + `'`,
		`GRANT USAGE ON SCHEMA ` + schema + ` TO "` + role + `"`,
		`GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA ` + schema + ` TO "` + role + `"`,
	} {
		_, err := r.dest.ExecContext(r.t.Context(), statement)
		require.NoError(r.t, err)
	}
	r.t.Cleanup(func() {
		// The test context is already canceled when cleanups run.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		for _, statement := range []string{`DROP OWNED BY "` + role + `"`, `DROP ROLE IF EXISTS "` + role + `"`} {
			_, err := r.dest.ExecContext(ctx, statement)
			require.NoError(r.t, err, "the fixture's own role is removed")
		}
	})

	var raw string
	require.NoError(r.t, r.f.db.QueryRowContext(r.t.Context(), `SELECT connection_config FROM database_connectors WHERE id='connector-A'`).Scan(&raw))
	config := map[string]string{}
	require.NoError(r.t, json.Unmarshal([]byte(raw), &config))
	config["user"], config["password"] = role, password
	encoded, err := json.Marshal(config)
	require.NoError(r.t, err)
	_, err = r.f.db.ExecContext(r.t.Context(), `UPDATE database_connectors SET connection_config = ? WHERE id='connector-A'`, string(encoded))
	require.NoError(r.t, err)
}

// extensionCount reports whether the TimescaleDB extension is installed (1) or not (0).
func (r *recoveryRig) extensionCount() int {
	r.t.Helper()
	var n int
	require.NoError(r.t, r.dest.QueryRowContext(r.t.Context(), `SELECT COUNT(*) FROM pg_extension WHERE extname = 'timescaledb'`).Scan(&n))
	return n
}

// removeTimescaleExtension drops the extension of the disposable fixture
// database and puts it back when the test ends. Only a fixture may do this: the
// gateway itself never installs or removes an extension.
func (r *recoveryRig) removeTimescaleExtension() {
	r.t.Helper()
	_, err := r.dest.ExecContext(r.t.Context(), `DROP EXTENSION timescaledb`)
	require.NoError(r.t, err)
	r.t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err := r.dest.ExecContext(ctx, `CREATE EXTENSION IF NOT EXISTS timescaledb`)
		require.NoError(r.t, err, "the fixture's extension is put back")
	})
}

// TestLongStorageRecoveryDeniedPostgres is the PostgreSQL side of what recovery
// must refuse: a foreign table holding the name and a role that may not create
// tables, or a hypertable layout whose extension is gone. Delivery stays blocked with a stable code, nothing is created, altered
// or deleted, and receipts stay as they were.
func TestLongStorageRecoveryDeniedPostgres(t *testing.T) {
	t.Run("a foreign table took the name", func(t *testing.T) {
		r := newRecoveryRig(t, "postgres")
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
		require.Equal(t, []string{"note"}, r.columns(r.table), "the foreign table is not altered")
		require.Equal(t, receipts, r.receipts())
		require.Zero(t, r.repairOperations(), "no DDL ran")
		require.Equal(t, 1, r.lineage().Generation)
		require.False(t, r.lineage().HistoryMayBeLost)
	})

	t.Run("the role may not create tables", func(t *testing.T) {
		r := newRecoveryRig(t, "postgres")
		r.start()
		r.record(10.5)
		r.committed(1, "first window")
		receipts := r.receipts()

		r.limitedRole()
		r.drop(r.table)
		r.record(11.5)

		r.waitFor(func() bool { return r.outbox()["blocked"] == 1 }, "a missing privilege blocks delivery")
		require.Equal(t, []string{recordingplan.StorageCodeRepairBlocked}, r.outboxCodes())
		require.False(t, r.tableExists(r.table), "no table is created without the privilege")
		require.Zero(t, r.repairOperations(), "no DDL ran")
		require.Equal(t, receipts, r.receipts())
		require.Equal(t, 1, r.lineage().Generation)

		storage := r.delivery(r.groups[0].id)["storage"].(map[string]any)
		events, _ := storage["events"].([]any)
		var reasons []string
		for _, event := range events {
			if reason, ok := event.(map[string]any)["reason"].(string); ok {
				reasons = append(reasons, reason)
			}
		}
		require.Contains(t, reasons, recordingplan.StorageReasonPermissionDenied, "the actionable cause is shown")
	})

	t.Run("the TimescaleDB extension is gone", func(t *testing.T) {
		r := newRecoveryRig(t, "timescale")
		r.start()
		r.record(10.5)
		r.committed(1, "first window")
		receipts := r.receipts()

		r.drop(r.table)
		r.removeTimescaleExtension()
		r.record(11.5)

		r.waitFor(func() bool { return r.outbox()["blocked"] == 1 }, "a hypertable cannot be rebuilt without its extension")
		require.Equal(t, []string{recordingplan.StorageCodeRepairBlocked}, r.outboxCodes())
		require.False(t, r.tableExists(r.table), "no plain table stands in for the hypertable")
		require.Zero(t, r.repairOperations(), "no DDL ran, and the extension is not installed on the gateway's behalf")
		require.Zero(t, r.extensionCount(), "the extension stays absent")
		require.Equal(t, receipts, r.receipts())
		require.Equal(t, 1, r.lineage().Generation)
		require.False(t, r.lineage().HistoryMayBeLost)

		storage := r.delivery(r.groups[0].id)["storage"].(map[string]any)
		var reasons []string
		for _, event := range storage["events"].([]any) {
			if reason, ok := event.(map[string]any)["reason"].(string); ok {
				reasons = append(reasons, reason)
			}
		}
		require.Contains(t, reasons, recordingplan.StorageReasonTimescaleUnavailable, "the actionable cause is shown")
	})
}
