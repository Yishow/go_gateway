package runtime

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/sqltx"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

func mixedDurableFixture(t *testing.T, f *longFixture, db *sql.DB) *groupdelivery.Store {
	t.Helper()
	store := groupdelivery.NewStore(db)
	g := f.config.Group
	key := groupdelivery.GroupKey{WorkspaceID: g.WorkspaceID, GroupID: g.ID, GroupRevision: g.AppliedRevision}
	payload, err := json.Marshal(workspace.WriteGroupAppliedSnapshot{WorkspaceID: g.WorkspaceID, GroupID: g.ID, AppliedRevision: g.AppliedRevision, EffectiveAt: f.config.StartsAt, Group: g, CounterBindings: f.config.CounterBindings, RuntimeLayout: &workspace.WriteGroupRuntimeLayout{TagTypes: f.config.TagTypes, Dialect: f.config.Dialect, Columns: f.config.Columns, ReceiptTableReady: f.config.ReceiptTableReady}})
	require.NoError(t, err)
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), key, payload))
	f.config.Ledger = store.Ledger(key)
	return store
}

func TestMixedDurableProductionBuilder(t *testing.T) {
	f := mixedBoundaryFixture(t)
	db := openDurableDB(t, filepath.Join(t.TempDir(), "gateway.db"))
	mixedDurableFixture(t, f, db)
	_, err := NewGroupBoundaryContext(t.Context(), f.config)
	require.NoError(t, err)
}

func TestMixedDurableDeltaSkipRestart(t *testing.T) {
	for _, scaled := range []bool{false, true} {
		t.Run(fmt.Sprint(scaled), func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "gateway.db")
			db := openDurableDB(t, path)
			f := mixedCounterBoundaryFixture(t, scaled)
			mixedDurableFixture(t, f, db)
			b, err := NewGroupBoundaryContext(t.Context(), f.config)
			require.NoError(t, err)
			f.accept(t, b, f.envelope(0, "base", boundaryAt(14), uint64(9007199254740993)))
			require.NoError(t, b.Tick(t.Context(), boundaryAt(15)))
			require.Equal(t, 0, countTable(t, db, "wg_delivery_outbox"))
			require.Equal(t, 1, countTable(t, db, "wg_delivery_mixed_outcomes"))
			require.Equal(t, 1, countTable(t, db, "wg_delivery_counter_baselines"))
			require.Equal(t, 0, countTable(t, db, "wg_delivery_checkpoints"))
			require.NoError(t, db.Close())
			reopened := openDurableDB(t, path)
			again := mixedCounterBoundaryFixture(t, scaled)
			mixedDurableFixture(t, again, reopened)
			restored, err := NewGroupBoundaryContext(t.Context(), again.config)
			require.NoError(t, err)
			again.accept(t, restored, again.envelope(0, "next", boundaryAt(29), uint64(9007199254741001)))
			require.NoError(t, restored.Tick(t.Context(), boundaryAt(30)))
			require.Equal(t, 1, countTable(t, reopened, "wg_delivery_outbox"))
			var payload string
			require.NoError(t, reopened.QueryRow(`SELECT outcome FROM wg_delivery_mixed_outcomes WHERE kind='row'`).Scan(&payload))
			var outcome snapshot.Outcome
			require.NoError(t, json.Unmarshal([]byte(payload), &outcome))
			digits, ok := outcome.Members[0].Counter.Delta.Value.DecimalDigits()
			require.True(t, ok)
			if scaled {
				require.Equal(t, "0.8", digits)
			} else {
				require.Equal(t, "8", digits)
			}
			require.NoError(t, restored.Tick(t.Context(), boundaryAt(30)))
			require.Equal(t, 1, countTable(t, reopened, "wg_delivery_outbox"))
		})
	}
}

func TestMixedDurableSlowMemberSurvivesFastCloseRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "gateway.db")
	db := openDurableDB(t, path)
	f := mixedBoundaryFixture(t)
	mixedDurableFixture(t, f, db)
	b, err := NewGroupBoundaryContext(t.Context(), f.config)
	require.NoError(t, err)
	f.accept(t, b, f.envelope(0, "mean", boundaryAt(10), int64(22)), f.envelope(1, "state", boundaryAt(10), "ON"))
	require.NoError(t, b.Tick(t.Context(), boundaryAt(15)))
	var pending int
	require.NoError(t, db.QueryRow(`SELECT count(*) FROM wg_delivery_samples WHERE consumed=0`).Scan(&pending))
	require.Equal(t, 1, pending)
	var legacy, due string
	require.NoError(t, db.QueryRow(`SELECT bucket_start,due_at FROM wg_delivery_outbox`).Scan(&legacy, &due))
	require.Empty(t, legacy)
	parsed, err := time.Parse(time.RFC3339Nano, due)
	require.NoError(t, err)
	require.True(t, boundaryAt(15).Equal(parsed))
	require.NoError(t, db.Close())
	reopened := openDurableDB(t, path)
	again := mixedBoundaryFixture(t)
	mixedDurableFixture(t, again, reopened)
	restored, err := NewGroupBoundaryContext(t.Context(), again.config)
	require.NoError(t, err)
	again.accept(t, restored, again.envelope(0, "last-mean", boundaryAt(59), int64(24)))
	require.NoError(t, restored.Tick(t.Context(), boundaryAt(60)))
	var payload string
	require.NoError(t, reopened.QueryRow(`SELECT outcome FROM wg_delivery_mixed_outcomes WHERE due_at=?`, boundaryAt(60).Format("2006-01-02T15:04:05.000000000Z")).Scan(&payload))
	var outcome snapshot.Outcome
	require.NoError(t, json.Unmarshal([]byte(payload), &outcome))
	require.Equal(t, "state", outcome.Members[1].Sample.SampleID)
	require.Equal(t, boundaryAt(0), outcome.Members[1].Window.AlignedStart)
	require.Equal(t, boundaryAt(45), outcome.Members[0].Window.AlignedStart)
}

type mixedUnknownLedger struct {
	*groupdelivery.Ledger
	appendUnknown  bool
	closureUnknown bool
	calls          int
}

func (l *mixedUnknownLedger) AppendSample(ctx context.Context, start time.Time, s snapshot.Sample) error {
	l.calls++
	if err := l.Ledger.AppendSample(ctx, start, s); err != nil {
		return err
	}
	if l.appendUnknown {
		return sqltx.ErrCommitOutcomeUnknown
	}
	return nil
}
func (l *mixedUnknownLedger) CommitClosure(ctx context.Context, c groupdelivery.Closure) error {
	l.calls++
	if err := l.Ledger.CommitClosure(ctx, c); err != nil {
		return err
	}
	if l.closureUnknown {
		return sqltx.ErrCommitOutcomeUnknown
	}
	return nil
}
func TestMixedDurableUnknownCommitRequiresRestore(t *testing.T) {
	for _, phase := range []string{"append", "closure"} {
		t.Run(phase, func(t *testing.T) {
			db := openDurableDB(t, filepath.Join(t.TempDir(), "gateway.db"))
			f := mixedBoundaryFixture(t)
			mixedDurableFixture(t, f, db)
			ledger := &mixedUnknownLedger{Ledger: f.config.Ledger.(*groupdelivery.Ledger), appendUnknown: phase == "append", closureUnknown: phase == "closure"}
			f.config.Ledger = ledger
			b, err := NewGroupBoundaryContext(t.Context(), f.config)
			require.NoError(t, err)
			err = b.AcceptSample(t.Context(), f.envelope(0, "accepted", boundaryAt(10), int64(22)))
			if phase == "append" {
				require.ErrorIs(t, err, sqltx.ErrCommitOutcomeUnknown)
			} else {
				require.NoError(t, err)
				require.ErrorIs(t, b.Tick(t.Context(), boundaryAt(15)), sqltx.ErrCommitOutcomeUnknown)
			}
			calls := ledger.calls
			require.ErrorIs(t, b.Tick(t.Context(), boundaryAt(30)), sqltx.ErrCommitOutcomeUnknown)
			require.ErrorIs(t, b.AcceptSample(t.Context(), f.envelope(0, "later", boundaryAt(20), int64(24))), sqltx.ErrCommitOutcomeUnknown)
			require.Equal(t, calls, ledger.calls, "unknown state must not retry/recompute accepted effects")
			f.config.Ledger = ledger.Ledger
			recovered, err := NewGroupBoundaryContext(t.Context(), f.config)
			require.NoError(t, err)
			require.NoError(t, recovered.Tick(t.Context(), boundaryAt(15)))
			require.Equal(t, 1, countTable(t, db, "wg_delivery_outbox"))
		})
	}
}

func TestMixedDurableRejectsChangedFrozenDescriptor(t *testing.T) {
	for _, field := range []string{"required", "cadence", "method", "record", "tagtype", "columns", "dialect", "receipt", "starts"} {
		t.Run(field, func(t *testing.T) {
			db := openDurableDB(t, filepath.Join(t.TempDir(), "gateway.db"))
			f := mixedBoundaryFixture(t)
			mixedDurableFixture(t, f, db)
			switch field {
			case "required":
				f.config.Group.Members[0].Required = !f.config.Group.Members[0].Required
			case "cadence":
				f.config.Group.Members[0].Cadence.Interval = 30
			case "method":
				f.config.Group.Members[0].Calculation = workspace.WriteGroupMemberCalculationLatest
			case "record":
				f.config.Group.RecordID = "ad8a4612-505b-4c13-a306-d76772f2887d"
			case "tagtype":
				f.config.TagTypes[f.members[0].TagID] = schema.DataTypeUint64
			case "columns":
				f.config.Columns = append(f.config.Columns, dbtarget.ColumnInfo{Name: "added", Nullable: true})
			case "dialect":
				f.config.Dialect = dbtarget.SQLDialectPostgres
			case "receipt":
				f.config.ReceiptTableReady = !f.config.ReceiptTableReady
			case "starts":
				f.config.StartsAt = boundaryAt(1)
			}
			_, err := NewGroupBoundaryContext(t.Context(), f.config)
			require.Error(t, err)
			require.Equal(t, 0, countTable(t, db, "wg_delivery_samples"))
			require.Equal(t, 0, countTable(t, db, "wg_delivery_mixed_progress"))
		})
	}
}

func TestMixedDurableCutoffAfterProgressRestoresExactEnd(t *testing.T) {
	db := openDurableDB(t, filepath.Join(t.TempDir(), "gateway.db"))
	f := mixedBoundaryFixture(t)
	mixedDurableFixture(t, f, db)
	b, err := NewGroupBoundaryContext(t.Context(), f.config)
	require.NoError(t, err)
	f.accept(t, b, f.envelope(0, "first", boundaryAt(10), int64(20)), f.envelope(1, "before", boundaryAt(10), "ON"), f.envelope(1, "future", boundaryAt(29), "OFF"))
	require.NoError(t, b.Tick(t.Context(), boundaryAt(15)))
	require.NoError(t, b.SetUntil(boundaryAt(23)))
	require.NoError(t, b.Tick(t.Context(), boundaryAt(60)))
	restored, err := NewGroupBoundaryContext(t.Context(), f.config)
	require.NoError(t, err)
	require.True(t, restored.Retired())
	require.True(t, restored.until.Equal(boundaryAt(23)))
	var payload string
	require.NoError(t, db.QueryRow(`SELECT outcome FROM wg_delivery_mixed_outcomes WHERE due_at=?`, boundaryAt(60).Format("2006-01-02T15:04:05.000000000Z")).Scan(&payload))
	var outcome snapshot.Outcome
	require.NoError(t, json.Unmarshal([]byte(payload), &outcome))
	require.Equal(t, "before", outcome.Members[1].Sample.SampleID)
	require.True(t, outcome.Members[1].Window.ActualEnd.Equal(boundaryAt(23)))
}

type mixedChangedRestoreLedger struct {
	*groupdelivery.Ledger
	digest string
}

func (l mixedChangedRestoreLedger) Restore(ctx context.Context) (groupdelivery.Restored, error) {
	restored, err := l.Ledger.Restore(ctx)
	restored.DescriptorDigest = l.digest
	return restored, err
}
func TestMixedDurableRejectsUnprovenDescriptorDigest(t *testing.T) {
	db := openDurableDB(t, filepath.Join(t.TempDir(), "gateway.db"))
	f := mixedBoundaryFixture(t)
	mixedDurableFixture(t, f, db)
	original := f.config.Ledger.(*groupdelivery.Ledger)
	for _, digest := range []string{"", "foreign"} {
		f.config.Ledger = mixedChangedRestoreLedger{Ledger: original, digest: digest}
		_, err := NewGroupBoundaryContext(t.Context(), f.config)
		require.Error(t, err)
	}
}

func TestMixedDurableBaselineFaultDoesNotPartiallyAdvance(t *testing.T) {
	db := openDurableDB(t, filepath.Join(t.TempDir(), "gateway.db"))
	f := mixedCounterBoundaryFixture(t, true)
	mixedDurableFixture(t, f, db)
	b, err := NewGroupBoundaryContext(t.Context(), f.config)
	require.NoError(t, err)
	f.accept(t, b, f.envelope(0, "base", boundaryAt(14), uint64(9007199254740993)))
	before, err := b.assembler.MixedCheckpoint()
	require.NoError(t, err)
	_, err = db.Exec(`CREATE TRIGGER fail_mixed_baseline BEFORE INSERT ON wg_delivery_counter_baselines BEGIN SELECT RAISE(ABORT,'baseline fault'); END`)
	require.NoError(t, err)
	require.ErrorContains(t, b.Tick(t.Context(), boundaryAt(15)), "baseline fault")
	for _, table := range []string{"wg_delivery_mixed_progress", "wg_delivery_mixed_outcomes", "wg_delivery_counter_baselines", "wg_delivery_outbox"} {
		require.Equal(t, 0, countTable(t, db, table), table)
	}
	var pending int
	require.NoError(t, db.QueryRow(`SELECT count(*) FROM wg_delivery_samples WHERE consumed=0`).Scan(&pending))
	require.Equal(t, 1, pending)
	after, err := b.assembler.MixedCheckpoint()
	require.NoError(t, err)
	require.Equal(t, before, after)
	require.Empty(t, b.assembler.CounterCheckpoints())
	_, err = db.Exec(`DROP TRIGGER fail_mixed_baseline`)
	require.NoError(t, err)
	require.NoError(t, b.Tick(t.Context(), boundaryAt(15)))
	var kind string
	require.NoError(t, db.QueryRow(`SELECT kind FROM wg_delivery_mixed_outcomes`).Scan(&kind))
	require.Equal(t, string(snapshot.OutcomeSkipped), kind)
	recovered, err := NewGroupBoundaryContext(t.Context(), f.config)
	require.NoError(t, err)
	f.accept(t, recovered, f.envelope(0, "next", boundaryAt(29), uint64(9007199254741001)))
	require.NoError(t, recovered.Tick(t.Context(), boundaryAt(30)))
	require.Equal(t, 1, countTable(t, db, "wg_delivery_outbox"))
}
