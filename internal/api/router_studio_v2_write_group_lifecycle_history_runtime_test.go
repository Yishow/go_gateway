package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Drive the real reconciled GroupBoundary and durable closure. A manually
// filled outbox can accidentally duplicate the reader's wrong assumptions.
func TestWriteGroupLifecycleHistoryReadsActualSQLiteRuntimeOutbox(t *testing.T) {
	r, _ := finalizedAPIRig(t, cutoverOptions{})
	r.accept(r.at(23*time.Second), 7.0)
	r.advance(r.at(31 * time.Second))
	r.pipe.TickAll(t.Context())
	require.Equal(t, 1, r.outbox()["pending"])
	id := r.wide["id"].(string)
	var recordKey string
	require.NoError(t, r.f.db.QueryRow(`SELECT record_key_column FROM wg_delivery_outbox WHERE group_id=?`, id).Scan(&recordKey))
	require.Empty(t, recordKey, "SQLite runtime deliberately omits the target record-key column from its frozen delivery identity")
	before := lifecycleAPIState(t, r.f.db)
	data := lifecycleAPIRead(t, r.f.writeGroupRouterFixture, id)
	require.Contains(t, data, "history")
	generation := data["history"].(map[string]any)["generations"].([]any)[0].(map[string]any)
	revisions := generation["revisions"].([]any)
	require.Len(t, revisions, 1)
	pending := revisions[0].(map[string]any)["pending"].(map[string]any)
	require.EqualValues(t, 0, pending["samples"])
	outbox := pending["outbox"].(map[string]any)
	require.EqualValues(t, 1, outbox["total"])
	require.EqualValues(t, 1, outbox["queued"])
	require.EqualValues(t, 0, outbox["sql_committed"])
	require.Equal(t, before, lifecycleAPIState(t, r.f.db))
}
