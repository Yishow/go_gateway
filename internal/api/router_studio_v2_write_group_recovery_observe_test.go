package api

// Observation helpers of the recovery rig: they read the destination's catalog and
// the local delivery tables, never the code under test.

import (
	"net/http"
	"time"

	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// quoted names a table of the destination in the SQL dialect of the rig.
func (r *recoveryRig) quoted(table string) string {
	if r.kind == "sqlite" {
		return `"` + table + `"`
	}
	return `"` + r.f.namespace + `"."` + table + `"`
}

// tableExists asks the destination catalog, not the code under test.
func (r *recoveryRig) tableExists(table string) bool {
	r.t.Helper()
	var n int
	if r.kind == "sqlite" {
		require.NoError(r.t, r.dest.QueryRowContext(r.t.Context(), `SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = ?`, table).Scan(&n))
	} else {
		require.NoError(r.t, r.dest.QueryRowContext(r.t.Context(), `SELECT COUNT(*) FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace WHERE n.nspname = $1 AND c.relname = $2 AND c.relkind IN ('r','p')`, r.f.namespace, table).Scan(&n))
	}
	return n > 0
}

func (r *recoveryRig) count(table string) int {
	r.t.Helper()
	var n int
	require.NoError(r.t, r.dest.QueryRowContext(r.t.Context(), `SELECT COUNT(*) FROM `+r.quoted(table)).Scan(&n))
	return n
}

func (r *recoveryRig) drop(table string) {
	r.t.Helper()
	_, err := r.dest.ExecContext(r.t.Context(), `DROP TABLE `+r.quoted(table))
	require.NoError(r.t, err)
}

// receipts returns effect key → payload digest|committed_at of every receipt.
func (r *recoveryRig) receipts() map[string]string {
	r.t.Helper()
	rows, err := r.dest.QueryContext(r.t.Context(), `SELECT effect_key, payload_digest, committed_at FROM `+r.quoted("gw_effect_receipts")+` ORDER BY effect_key`)
	require.NoError(r.t, err)
	defer rows.Close()
	out := map[string]string{}
	for rows.Next() {
		var key, digest, committed string
		require.NoError(r.t, rows.Scan(&key, &digest, &committed))
		out[key] = digest + "|" + committed
	}
	require.NoError(r.t, rows.Err())
	return out
}

// outbox returns state → count over the local delivery outbox.
func (r *recoveryRig) outbox() map[string]int {
	r.t.Helper()
	rows, err := r.f.db.QueryContext(r.t.Context(), `SELECT state, COUNT(*) FROM wg_delivery_outbox GROUP BY state`)
	require.NoError(r.t, err)
	defer rows.Close()
	out := map[string]int{}
	for rows.Next() {
		var state string
		var n int
		require.NoError(r.t, rows.Scan(&state, &n))
		out[state] = n
	}
	require.NoError(r.t, rows.Err())
	return out
}

// outboxCodes returns the last error code of every undelivered outbox row.
func (r *recoveryRig) outboxCodes() []string {
	r.t.Helper()
	rows, err := r.f.db.QueryContext(r.t.Context(), `SELECT DISTINCT last_error_code FROM wg_delivery_outbox WHERE state <> 'sql_committed' AND last_error_code <> '' ORDER BY 1`)
	require.NoError(r.t, err)
	defer rows.Close()
	var codes []string
	for rows.Next() {
		var code string
		require.NoError(r.t, rows.Scan(&code))
		codes = append(codes, code)
	}
	require.NoError(r.t, rows.Err())
	return codes
}

// committed waits until n outbox rows are sql_committed, and on failure reports
// where the outbox stood then, not when the wait began.
func (r *recoveryRig) committed(n int, message string) {
	r.t.Helper()
	r.waitFor(func() bool { return r.outbox()["sql_committed"] >= n }, message)
}

// delivery fetches the public delivery view of a group over the HTTP route.
func (r *recoveryRig) delivery(groupID string) map[string]any {
	r.t.Helper()
	res := performJSONRequest(r.t, r.f.router, http.MethodGet, writeGroupsPath+"/"+groupID+"/delivery", nil)
	require.Equal(r.t, http.StatusOK, res.Code, res.Body.String())
	return decodeJSONBody(r.t, res)["data"].(map[string]any)
}

// lineage reads the local storage scope row of the rig's shared long table.
func (r *recoveryRig) lineage() *recordingplan.StorageLineage {
	r.t.Helper()
	lineage, err := recordingplan.ReadStorageLineage(r.t.Context(), r.f.db, workspace.StorageScopeOf(r.groups[0].snapshot.Group))
	require.NoError(r.t, err)
	return lineage
}

// waitFor polls until cond holds; its failure message carries the outbox state
// at the moment the wait gave up.
func (r *recoveryRig) waitFor(cond func() bool, message string) {
	r.t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(25 * time.Millisecond)
	}
	require.Fail(r.t, "condition never satisfied", "%s: outbox %v codes %v", message, r.outbox(), r.outboxCodes())
}
