package runtime

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"go-gateway/internal/datalink/groupdelivery"

	"github.com/stretchr/testify/require"
)

func longLedger(db *sql.DB) GroupLedger {
	store := groupdelivery.NewStore(db)
	return &frozenFixtureLedger{store: store, Ledger: store.Ledger(groupdelivery.GroupKey{
		WorkspaceID: "workspace-A", GroupID: "group-L", GroupRevision: "rev-1",
	})}
}

// outboxRows decodes every closed long record from the durable outbox by tag.
func outboxRows(t *testing.T, db *sql.DB) map[string]map[string]any {
	t.Helper()
	rows, err := db.QueryContext(t.Context(), `SELECT payload FROM wg_delivery_outbox`)
	require.NoError(t, err)
	defer rows.Close()
	decoded := map[string]map[string]any{}
	for rows.Next() {
		var payload string
		require.NoError(t, rows.Scan(&payload))
		row, err := groupdelivery.DecodeRowPayload([]byte(payload))
		require.NoError(t, err)
		cells := map[string]any{}
		for _, cell := range row.Cells {
			cells[cell.Column] = cell.Value
		}
		decoded[cells["tag_id"].(string)] = cells
	}
	require.NoError(t, rows.Err())
	return decoded
}

// TestProfileReceivedAtSurvivesRestart proves the receive time is journaled with
// the accepted sample: a restart replays the original time, it neither drops it
// nor substitutes the restart or tick time, and a window without a sample keeps
// NULL times.
func TestProfileReceivedAtSurvivesRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "gateway.db")
	db := openDurableDB(t, path)
	f := newLongFixture(t)
	f.config.Ledger = longLedger(db)
	boundary := f.build(t)
	original := f.envelope(0, "a-1", at(6*time.Second), 21.5)
	original.ReceivedAt = at(6*time.Second + 123456*time.Microsecond)
	f.accept(t, boundary, original)
	f.accept(t, boundary, f.envelope(4, "count-1", at(9*time.Second), uint64(18446744073709551615)))
	require.Equal(t, 2, countJournal(t, db))

	// Crash after the ACKs; the process restarts well after the original receive time.
	require.NoError(t, db.Close())
	reopened := openDurableDB(t, path)
	again := newLongFixture(t)
	again.config.Ledger = longLedger(reopened)
	again.clock = at(10 * time.Minute) // the restart happens long after the original receive time
	recovered := again.build(t)
	require.NoError(t, recovered.Tick(t.Context(), at(16*time.Second)))

	rows := outboxRows(t, reopened)
	require.Len(t, rows, 5, "one durable record per member for the window")
	a := rows["tag-A"]
	require.Equal(t, at(6*time.Second).Format(longTimeLayout), a["observed_at"])
	require.Equal(t, at(6*time.Second+123456*time.Microsecond).Format(longTimeLayout), a["received_at"],
		"the replayed receive time is the original one, not the restart time")
	require.Equal(t, 21.5, a["value_float"])
	count := rows["tag-count"]
	require.Equal(t, maxUint64Text, count["value_decimal"], "uint64 survives journal, closure and outbox exactly")
	for _, tag := range []string{"tag-B", "tag-run", "tag-batch"} {
		require.Equal(t, "missing", rows[tag]["status"], tag)
		require.Nil(t, rows[tag]["observed_at"], tag)
		require.Nil(t, rows[tag]["received_at"], tag)
	}
	require.Equal(t, 5, countTable(t, reopened, "wg_delivery_outbox"))
}

// TestProfileDurableClosureIsOneEffectPerRecord checks that the closure keeps
// each member's record and effect identity distinct and that a repeated Tick
// neither reissues them nor changes the closed window.
func TestProfileDurableClosureIsOneEffectPerRecord(t *testing.T) {
	db := openDurableDB(t, filepath.Join(t.TempDir(), "gateway.db"))
	f := newLongFixture(t)
	f.config.Ledger = longLedger(db)
	boundary := f.build(t)
	f.accept(t, boundary, f.envelope(0, "a-1", at(2*time.Second), 1.0), f.envelope(1, "b-1", at(3*time.Second), int64(5)))

	require.NoError(t, boundary.Tick(t.Context(), at(15*time.Second)))
	require.Equal(t, 5, countTable(t, db, "wg_delivery_outbox"))
	var distinct int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(DISTINCT effect_key) FROM wg_delivery_outbox`).Scan(&distinct))
	require.Equal(t, 5, distinct)

	require.NoError(t, boundary.Tick(t.Context(), at(16*time.Second)))
	require.NoError(t, boundary.Tick(t.Context(), at(17*time.Second)))
	require.Equal(t, 5, countTable(t, db, "wg_delivery_outbox"), "repeated ticks never issue new identities for a closed window")
}
