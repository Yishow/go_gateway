package runtime

import (
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/snapshot"

	"github.com/stretchr/testify/require"
)

func TestEntityStructuralFailureKeepsAcknowledgedJournalAndCheckpoint(t *testing.T) {
	path := filepath.Join(t.TempDir(), "configuration.db")
	db := openDurableDB(t, path)
	f := newBoundaryFixture(t)
	f.withLedger(db)
	boundary := f.build(t)
	f.feedAll(t, boundary)
	require.Equal(t, 5, countJournal(t, db))
	// Represent a recovered layout that no longer agrees with the accepted
	// snapshot. The encoder expects a member absent from that snapshot row.
	inconsistent, issues := dbtarget.NewGroupRowLayout(dbtarget.GroupRowSpec{
		Dialect: dbtarget.SQLDialectSQLite,
		Columns: []dbtarget.ColumnInfo{{Name: "unresolved_member", DataType: "REAL"}},
		Members: []dbtarget.GroupRowMember{{MemberKey: "unresolved", Column: "unresolved_member", Type: measurement.ExactFloat64, Required: true}},
	})
	require.Empty(t, issues)
	boundary.layout = inconsistent
	err := boundary.Tick(t.Context(), boundaryAt(10))
	var blocked *GroupBoundaryError
	require.ErrorAs(t, err, &blocked, "structural mismatch must not become an ordinary quality skip")
	require.Equal(t, "layout-blocked", blocked.Code)
	var rowErr *dbtarget.GroupRowError
	require.ErrorAs(t, err, &rowErr)
	require.Equal(t, "member-missing", rowErr.Code)
	require.NotContains(t, err.Error(), "unresolved_member", "only safe diagnostic codes reach operators")
	var open int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM wg_delivery_samples WHERE consumed = 0`).Scan(&open))
	require.Equal(t, 5, open)
	require.Zero(t, countTable(t, db, "wg_delivery_checkpoints"))
	require.Zero(t, countTable(t, db, "wg_delivery_buckets"))
	require.Zero(t, countTable(t, db, "wg_delivery_outbox"))

	// Repair via the original verified layout on an actual reopened store.
	require.NoError(t, db.Close())
	reopened := openDurableDB(t, path)
	again := newBoundaryFixture(t)
	again.withLedger(reopened)
	recovered := again.build(t)
	require.NoError(t, recovered.Tick(t.Context(), boundaryAt(10)))
	require.Equal(t, 1, countTable(t, reopened, "wg_delivery_outbox"))
	require.Equal(t, 1, countTable(t, reopened, "wg_delivery_checkpoints"))
	require.Equal(t, 21.5, outboxCells(t, reopened)["temperature"])
}

func TestEntityRowsKeepMissingPartialAndSilentOutcomesLocal(t *testing.T) {
	for _, partial := range []bool{false, true} {
		t.Run(map[bool]string{false: "skip-row", true: "partial"}[partial], func(t *testing.T) {
			db := openDurableDB(t, filepath.Join(t.TempDir(), "configuration.db"))
			f := newBoundaryFixture(t)
			group := copyGroup(f.config.Group)
			group.Members = group.Members[:4]
			for i := range group.Members {
				group.Members[i].EntityKey = "A"
				if i > 1 {
					group.Members[i].EntityKey = "B"
				}
			}
			group.RowPolicy.EntityKeyColumn = "entity"
			if partial {
				group.RowPolicy.IncompletePolicy = "partial"
				group.RowPolicy.ProvenanceColumn = "provenance"
				group.Members[1].Required = false
				f.config.Columns[1].Nullable = true
			}
			f.config.Columns = append(f.config.Columns, dbtarget.ColumnInfo{Name: "entity", DataType: "TEXT"})
			f.config.Group, f.members = group, group.Members
			f.withLedger(db)
			boundary := f.build(t)
			require.NoError(t, boundary.AcceptSample(t.Context(), f.envelope(0, "a-good", 4, 21.5)))
			require.NoError(t, boundary.AcceptSample(t.Context(), f.envelope(2, "b-running", 4, true)))
			require.NoError(t, boundary.AcceptSample(t.Context(), f.envelope(3, "b-batch", 4, "batch-B")))
			require.NoError(t, boundary.Tick(t.Context(), boundaryAt(10)))
			var aKind, bKind string
			require.NoError(t, db.QueryRowContext(t.Context(), `SELECT kind FROM wg_delivery_buckets WHERE entity_key = 'A'`).Scan(&aKind))
			require.NoError(t, db.QueryRowContext(t.Context(), `SELECT kind FROM wg_delivery_buckets WHERE entity_key = 'B'`).Scan(&bKind))
			require.Equal(t, string(snapshot.OutcomeRow), bKind)
			wantA := snapshot.OutcomeSkipped
			if partial {
				wantA = snapshot.OutcomeRow
			}
			require.Equal(t, string(wantA), aKind)
			f.clock = boundaryAt(21)
			require.NoError(t, boundary.Tick(t.Context(), boundaryAt(20)))
			var silent int
			require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM wg_delivery_buckets WHERE kind = 'no_data'`).Scan(&silent))
			require.Equal(t, 2, silent, "one silent bucket for each entity, without NULL/default rows")
		})
	}
}
