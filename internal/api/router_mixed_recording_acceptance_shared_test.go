package api

import (
	"database/sql"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMixedAcceptanceSharedLongPreservesHistoricalRows(t *testing.T) {
	f, request := mixedNativeRequest(t, false, []int64{15, 60, 60}, []bool{true, false, false})
	base := time.Date(2026, 10, 31, 12, 0, 0, 0, time.UTC)
	old := &mixedAcceptanceRig{f: f, clock: base.Add(7 * time.Second)}
	f.groups.WithClock(func() time.Time { return old.clock })
	legacy := ensureLong(t, f, "instantaneous", "second", 15, "")
	id := legacy["id"].(string)
	confirmLong(t, f, id, previewLong(t, f, id))
	applied := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/apply", scopeBody(t, f, id))
	require.Equal(t, http.StatusOK, applied.Code, applied.Body.String())
	var err error
	old.frozen, err = f.groups.ResolveAppliedAt(t.Context(), id, base.Add(time.Minute))
	require.NoError(t, err)
	old.clock = old.frozen.EffectiveAt.Add(time.Second)
	old.store = newMixedAcceptanceStore(t, f.db, func() time.Time { return old.clock })
	old.reconcile(t)
	old.f.target, err = sql.Open("sqlite", f.file)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, old.f.target.Close()) })
	old.offer(t, 0, "historical-raw", old.clock, 21.5)
	old.tick(t, old.frozen.EffectiveAt.Add(15*time.Second))
	require.Len(t, old.deliver(t), 1)
	table := old.frozen.Group.Destination.TableName
	before := mixedHistoricalRow(t, old.f.target, table, id)

	latest, err := f.groups.Get(t.Context(), id)
	require.NoError(t, err)
	request["expected_workspace_revision"] = latest.WorkspaceRevision
	current := applyMixedAcceptance(t, f, request, base.Add(67*time.Second))
	require.Equal(t, table, current.frozen.Group.Destination.TableName, "existing shared-long authority, not an arbitrary replacement table")
	require.NotEqual(t, id, current.frozen.GroupID)
	current.offer(t, 0, "new-mean", base.Add(70*time.Second), float64(22))
	current.tick(t, base.Add(75*time.Second))
	require.Len(t, current.deliver(t), 1)
	after := mixedHistoricalRow(t, current.f.target, table, id)
	require.Equal(t, before, after, "fresh mixed schema/record preserves every legacy column and value")

	var profile string
	var result float64
	require.NoError(t, current.f.target.QueryRow(`SELECT profile_kind,value_float FROM "`+table+`" WHERE group_id=?`, current.frozen.GroupID).Scan(&profile, &result))
	require.Equal(t, "mixed", profile)
	require.Equal(t, float64(22), result)
	legacyRead, err := f.groups.Get(t.Context(), id)
	require.NoError(t, err)
	require.Equal(t, old.frozen.Group.AppliedRevision, legacyRead.Group.AppliedRevision)
	require.NotEqual(t, "disabled", string(legacyRead.Group.Status), "creating a new record must not silently stop the old record")
}

func mixedHistoricalRow(t *testing.T, db *sql.DB, table, id string) map[string]any {
	t.Helper()
	rows, err := db.Query(`SELECT * FROM "`+table+`" WHERE group_id=?`, id)
	require.NoError(t, err)
	defer rows.Close()
	columns, err := rows.Columns()
	require.NoError(t, err)
	require.True(t, rows.Next())
	values := make([]any, len(columns))
	dest := make([]any, len(columns))
	for i := range values {
		dest[i] = &values[i]
	}
	require.NoError(t, rows.Scan(dest...))
	result := map[string]any{}
	for i, c := range columns {
		result[c] = values[i]
	}
	require.False(t, rows.Next(), "one accepted legacy row, not a rewritten or duplicated history")
	require.NoError(t, rows.Err())
	return result
}
