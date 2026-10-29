package api

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/workspace"
)

func successorRuntimeOriginals(t *testing.T, r *successorRuntimeRig) map[string]string {
	t.Helper()
	rows, err := r.f.db.Query(`SELECT sample_id,payload,payload_digest,created_at FROM wg_delivery_samples WHERE group_id=? ORDER BY seq`, r.groupID)
	require.NoError(t, err)
	defer rows.Close()
	result := map[string]string{}
	for rows.Next() {
		var id, payload, digest, created string
		require.NoError(t, rows.Scan(&id, &payload, &digest, &created))
		result[id] = payload + "|" + digest + "|" + created
	}
	require.NoError(t, rows.Err())
	return result
}

func successorRuntimeHandoffs(t *testing.T, r *successorRuntimeRig, originals map[string]string, consumed bool) {
	t.Helper()
	require.Equal(t, originals, successorRuntimeOriginals(t, r))
	rows, err := r.f.db.Query(`SELECT s.sample_id,s.payload,s.payload_digest,s.created_at,s.consumed,t.origin_sample_id,t.payload,t.payload_digest,t.created_at,t.consumed,h.transfer_id,t.handoff_id,t.origin_group_id
 FROM wg_delivery_sample_handoffs h JOIN wg_delivery_samples s ON s.seq=h.source_seq JOIN wg_delivery_samples t ON t.seq=h.target_seq
 WHERE h.source_group_id=? AND h.target_group_id=? ORDER BY s.sample_id`, r.groupID, r.prepared.Group.ID)
	require.NoError(t, err)
	defer rows.Close()
	ids := []string{}
	for rows.Next() {
		var id, payload, digest, created, originID, targetPayload, targetDigest, targetCreated, transfer, handoff, originGroup string
		var oldConsumed, newConsumed bool
		require.NoError(t, rows.Scan(&id, &payload, &digest, &created, &oldConsumed, &originID, &targetPayload, &targetDigest, &targetCreated, &newConsumed, &transfer, &handoff, &originGroup))
		require.Equal(t, id, originID)
		require.Equal(t, originals[id], payload+"|"+digest+"|"+created)
		require.Equal(t, originals[id], targetPayload+"|"+targetDigest+"|"+targetCreated)
		require.Equal(t, r.groupID, originGroup)
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer, handoff)
		require.False(t, oldConsumed)
		require.Equal(t, consumed, newConsumed)
		ids = append(ids, id)
	}
	require.NoError(t, rows.Err())
	require.Equal(t, []string{"after", "equal"}, ids, "only observations at/after B transfer, preserving their first acceptance")
}

func successorRuntimeValues(t *testing.T, db *sql.DB, g *workspace.WriteGroup) []string {
	t.Helper()
	column := `COALESCE(CAST(value_decimal AS TEXT),CAST(value_float AS TEXT),CAST(value_int AS TEXT))`
	if g.EffectiveLayoutVersion() == workspace.WriteGroupLayoutManagedWideProfileV1 {
		column = `CAST("` + g.Members[0].TargetColumn + `" AS TEXT)`
	}
	rows, err := db.Query(`SELECT ` + column + ` FROM "` + g.Destination.TableName + `" WHERE ` + column + ` IS NOT NULL ORDER BY bucket_start`)
	require.NoError(t, err)
	defer rows.Close()
	values := []string{}
	for rows.Next() {
		var value string
		require.NoError(t, rows.Scan(&value))
		values = append(values, plainDecimal(value))
	}
	require.NoError(t, rows.Err())
	return values
}

func successorRuntimeRows(t *testing.T, db *sql.DB, table string) map[string]string {
	t.Helper()
	return successorRuntimeEvidenceRows(t, db, table, "bucket_start")
}

func successorRuntimeEvidenceRows(t *testing.T, db *sql.DB, table, identity string) map[string]string {
	t.Helper()
	rows, err := db.Query(`SELECT * FROM "` + table + `" ORDER BY "` + identity + `"`)
	require.NoError(t, err)
	defer rows.Close()
	columns, err := rows.Columns()
	require.NoError(t, err)
	result := map[string]string{}
	for rows.Next() {
		values := make([]any, len(columns))
		dest := make([]any, len(values))
		for i := range values {
			dest[i] = &values[i]
		}
		require.NoError(t, rows.Scan(dest...))
		raw, err := json.Marshal(values)
		require.NoError(t, err)
		// The stored bucket_start is immutable identity, not the mutable row count.
		for i, name := range columns {
			if name == identity {
				key, err := json.Marshal(values[i])
				require.NoError(t, err)
				result[string(key)] = string(raw)
			}
		}
	}
	require.NoError(t, rows.Err())
	return result
}
func successorRuntimePreservedRows(t *testing.T, before, after map[string]string) {
	t.Helper()
	for key, row := range before {
		require.Equal(t, row, after[key], "old finalized SQL row must remain byte-identical")
	}
}
