package api

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Observation helpers of the cutover rig: they read the destination, the local
// database and the public routes, never the code under test.

func (r *cutoverRig) localCount(query string, args ...any) int {
	r.t.Helper()
	var n int
	require.NoError(r.t, r.f.db.QueryRowContext(r.t.Context(), query, args...).Scan(&n))
	return n
}

func (r *cutoverRig) tableExists(table string) bool {
	r.t.Helper()
	return (&recoveryRig{t: r.t, f: r.f, kind: r.kind, dest: r.dest}).tableExists(table)
}

func (r *cutoverRig) receipts() map[string]string {
	r.t.Helper()
	return (&recoveryRig{t: r.t, f: r.f, kind: r.kind, dest: r.dest}).receipts()
}

func (r *cutoverRig) outboxDigests() map[string]string {
	r.t.Helper()
	return (&recoveryRig{t: r.t, f: r.f}).outboxDigests()
}

func (r *cutoverRig) delivery(groupID string) map[string]any {
	r.t.Helper()
	return (&recoveryRig{t: r.t, f: r.f}).delivery(groupID)
}

// versionPayload is the stored, immutable content of one applied revision.
func (r *cutoverRig) versionPayload(groupID, revision string) string {
	r.t.Helper()
	var payload string
	require.NoError(r.t, r.f.db.QueryRowContext(r.t.Context(),
		`SELECT payload FROM write_group_versions WHERE group_id = ? AND group_revision = ?`, groupID, revision).Scan(&payload))
	return payload
}

// rollbackBody is the request that returns a cutover's sources at boundary.
func (r *cutoverRig) rollbackBody(cutoverID string, boundary time.Time, extra map[string]any) map[string]any {
	r.t.Helper()
	body := r.cutoverBody(boundary, extra)
	delete(body, "targets")
	body["kind"], body["cutover_id"] = "rollback", cutoverID
	return body
}

func (r *cutoverRig) rollbackReview(cutoverID string, boundary time.Time) map[string]any {
	r.t.Helper()
	code, body := r.post(writeGroupsPath+"/"+r.wide["id"].(string)+"/cutover-review", r.rollbackBody(cutoverID, boundary, nil))
	require.Equal(r.t, http.StatusOK, code, "%v", body)
	return body["data"].(map[string]any)
}

func toStrings(values any) []string {
	items, _ := values.([]any)
	out := make([]string, 0, len(items))
	for _, item := range items {
		if text, ok := item.(string); ok {
			out = append(out, text)
		}
	}
	return out
}

// wideRow is one row of the wide destination table as the destination stores it.
type wideRow struct {
	Start      time.Time
	First      *float64 // the first source's column; nil is SQL NULL
	Second     *float64 // the second source's column, when the group has one
	Provenance string
}

// wideRowsInOrder reads the wide table in window order.
func (r *cutoverRig) wideRowsInOrder() []wideRow {
	r.t.Helper()
	second := "NULL"
	if r.second.TargetColumn != "" {
		second = `"` + r.second.TargetColumn + `"`
	}
	rows, err := r.dest.QueryContext(r.t.Context(), `SELECT CAST(bucket_start AS TEXT), "`+r.member.TargetColumn+`", `+second+`, COALESCE(CAST(provenance AS TEXT), '')
		FROM `+r.quoted(r.wideTable())+` ORDER BY bucket_start, record_id`)
	require.NoError(r.t, err)
	defer rows.Close()
	var out []wideRow
	for rows.Next() {
		var row wideRow
		var start string
		var first, secondValue *float64
		require.NoError(r.t, rows.Scan(&start, &first, &secondValue, &row.Provenance))
		row.Start, row.First, row.Second = storedInstant(r.t, start), first, secondValue
		out = append(out, row)
	}
	require.NoError(r.t, rows.Err())
	return out
}

// windowOfWide returns the truncation evidence one wide row's provenance carries.
func windowOfWide(t *testing.T, provenance string) map[string]any {
	t.Helper()
	var entries []struct {
		Window map[string]any `json:"window"`
	}
	require.NoError(t, json.Unmarshal([]byte(provenance), &entries))
	for _, entry := range entries {
		if entry.Window != nil {
			return entry.Window
		}
	}
	return nil
}

// cutoverCounts is what a cutover writes locally, to prove a replay adds nothing.
func (r *cutoverRig) cutoverCounts() [4]int {
	r.t.Helper()
	return [4]int{
		r.localCount(`SELECT COUNT(*) FROM write_group_cutovers`),
		r.localCount(`SELECT COUNT(*) FROM write_group_revision_ends`),
		r.localCount(`SELECT COUNT(*) FROM write_group_versions`),
		r.localCount(`SELECT COUNT(*) FROM write_group_storage_scopes`),
	}
}
