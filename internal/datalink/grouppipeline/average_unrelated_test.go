package grouppipeline

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"

	"github.com/stretchr/testify/require"
)

// unrelatedKey is another group of the same workspace that already recorded in
// the same local database. Its member key is deliberately the first member's key
// of the group under test, so anything keyed by member alone would collide.
var unrelatedKey = groupdelivery.GroupKey{WorkspaceID: "workspace-A", GroupID: "group-OTHER", GroupRevision: "rev-1"}

const unrelatedMember = `["device-1","point-1","tag-kwh-1"]`

func unrelatedBucket(label string, start time.Time, value uint64) groupdelivery.ClosedBucket {
	outcome := snapshot.Outcome{
		Kind: snapshot.OutcomeRow, RecordID: "record-other-" + label, EffectKey: "effect-other-" + label,
		BucketStart: start, BucketEnd: start.Add(10 * time.Second),
		Members: []snapshot.MemberResult{{
			MemberKey: unrelatedMember, Status: snapshot.MemberOK,
			Sample: &snapshot.Sample{
				SampleID: "s-" + label, MemberKey: unrelatedMember, ObservedAt: start.Add(5 * time.Second),
				Quality: schema.QualityGood, Value: measurement.NewUint64(value),
			},
		}},
	}
	row := dbtarget.EncodedRow{
		RecordID: outcome.RecordID, EffectKey: outcome.EffectKey, BucketStart: start,
		Cells: []dbtarget.EncodedCell{
			{Column: "record_id", Value: outcome.RecordID}, {Column: "counter", Value: int64(value)}, {Column: "bucket_start", Value: start},
		},
	}
	return groupdelivery.ClosedBucket{Outcome: outcome, Row: &row}
}

// seedUnrelatedGroup saves what an earlier run of another group left in the same
// Store, through the Store's public interface: two frozen results (one delivered
// with its local receipt, one held), the bucket outcomes, the consumed readings,
// a reading still waiting in the journal, the checkpoint and a counter baseline.
// It returns the effect keys.
func (f *longPipeline) seedUnrelatedGroup() []string {
	f.t.Helper()
	ctx := f.t.Context()
	start := base.Add(-2 * time.Hour)
	nextClose := start.Add(20 * time.Second)
	sample := func(id string, seconds int, value uint64) snapshot.Sample {
		observed := start.Add(time.Duration(seconds) * time.Second)
		return snapshot.Sample{
			SampleID: id, MemberKey: unrelatedMember, ObservedAt: observed, ReceivedAt: observed.Add(40 * time.Millisecond),
			Quality: schema.QualityGood, Value: measurement.NewUint64(value), SourceRevision: "src-1", MappingRevision: "map-1",
		}
	}
	require.NoError(f.t, f.store.AppendSample(ctx, unrelatedKey, start, sample("s-closed", 4, 40)))
	require.NoError(f.t, f.store.AppendSample(ctx, unrelatedKey, nextClose, sample("s-open", 24, 41)))

	checkpoint, err := aggregation.EncodeCounterCheckpoint(aggregation.CounterCheckpoint{
		Binding: aggregation.CounterBinding{
			WorkspaceID: unrelatedKey.WorkspaceID, GroupID: unrelatedKey.GroupID, ProfileRevision: unrelatedKey.GroupRevision,
			MemberKey: unrelatedMember, MeasurementID: "meas-other", DefinitionRevision: "def-9", SourceBindingRevision: "bind-9",
			SeriesEpoch: "epoch-9", Kind: aggregation.CounterKindMonotonic, Unit: "kWh",
		},
		WindowEnd: nextClose,
		Reading: aggregation.CounterReading{
			SampleID: "s-closed", ObservedAt: start.Add(4 * time.Second), ReceivedAt: start.Add(4*time.Second + 40*time.Millisecond),
			Value: measurement.NewUint64(40),
		},
	})
	require.NoError(f.t, err)
	delivered, held := unrelatedBucket("delivered", start, 40), unrelatedBucket("held", start.Add(10*time.Second), 41)
	require.NoError(f.t, f.store.CommitClosure(ctx, unrelatedKey, groupdelivery.Closure{
		Destination: groupdelivery.FrozenDestination{
			Scope: "scope-other", ConnectorID: "connector-other", ConnectorRevision: "crev-9", Database: "db-other",
			TableSchema: "main", TableName: "other_readings", DedupeCapability: groupdelivery.DedupeReceipt, RecordKeyColumn: "record_id",
		},
		NextClose: nextClose, Buckets: []groupdelivery.ClosedBucket{delivered, held},
		Baselines: []groupdelivery.CounterBaseline{{MemberKey: unrelatedMember, Payload: checkpoint}},
	}))

	digest, err := groupdelivery.RowPayloadDigest(*delivered.Row)
	require.NoError(f.t, err)
	claim, err := f.store.BeginDelivery(ctx, delivered.Outcome.EffectKey, "other-owner", time.Minute)
	require.NoError(f.t, err)
	require.NoError(f.t, f.store.CompleteDelivery(ctx, claim, digest))
	require.NoError(f.t, f.store.MarkBlocked(ctx, held.Outcome.EffectKey, "operator-hold"))
	return []string{delivered.Outcome.EffectKey, held.Outcome.EffectKey}
}

func dumpRows(t *testing.T, db *sql.DB, query string, args ...any) []string {
	t.Helper()
	rows, err := db.QueryContext(t.Context(), query, args...)
	require.NoError(t, err)
	defer rows.Close()
	columns, err := rows.Columns()
	require.NoError(t, err)
	var out []string
	for rows.Next() {
		values := make([]any, len(columns))
		pointers := make([]any, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}
		require.NoError(t, rows.Scan(pointers...))
		cells := make([]string, len(columns))
		for i, value := range values {
			if raw, ok := value.([]byte); ok {
				value = string(raw)
			}
			cells[i] = fmt.Sprintf("%s=%v", columns[i], value)
		}
		out = append(out, fmt.Sprint(cells))
	}
	require.NoError(t, rows.Err())
	return out
}

// unrelatedState is every stored byte of the other group, table by table.
func (f *longPipeline) unrelatedState(effects []string) map[string][]string {
	f.t.Helper()
	local := f.plainLocal()
	state := map[string][]string{}
	for _, table := range []string{
		"wg_delivery_outbox", "wg_delivery_buckets", "wg_delivery_samples", "wg_delivery_checkpoints", "wg_delivery_counter_baselines",
	} {
		state[table] = dumpRows(f.t, local, `SELECT * FROM `+table+` WHERE group_id = ? ORDER BY 1, 2, 3`, unrelatedKey.GroupID)
	}
	state["wg_delivery_receipts"] = dumpRows(f.t, local, `SELECT * FROM wg_delivery_receipts WHERE effect_key IN (?, ?) ORDER BY 1`, effects[0], effects[1])
	return state
}

// TestAverageClosureKeepsUnrelatedGroups pins the other half of the closure
// failure scenario: a closure transaction of the average that fails after it has
// already written its rows, then retries and survives a restart, leaves every
// stored byte of an unrelated group in the same Store exactly as it was —
// frozen results and their identity, delivery state and receipt, bucket outcomes,
// consumed and still-waiting readings, checkpoint and counter baseline.
func TestAverageClosureKeepsUnrelatedGroups(t *testing.T) {
	for _, shape := range averageShapes {
		t.Run(string(shape), func(t *testing.T) {
			f := newAveragePipeline(t, shape)
			f.feed(shape, 0)
			f.closeUntil(shape.windowEnd(0))
			f.deliver()
			f.feed(shape, 1)

			effects := f.seedUnrelatedGroup()
			frozen := f.unrelatedState(effects)
			counts := map[string]int{
				"wg_delivery_outbox": 2, "wg_delivery_buckets": 2, "wg_delivery_samples": 2,
				"wg_delivery_checkpoints": 1, "wg_delivery_counter_baselines": 1, "wg_delivery_receipts": 1,
			}
			for table, want := range counts {
				require.Len(t, frozen[table], want, "the unrelated group has its own %s rows", table)
			}
			require.Contains(t, frozen["wg_delivery_outbox"][0]+frozen["wg_delivery_outbox"][1], "state=sql_committed")
			require.Contains(t, frozen["wg_delivery_outbox"][0]+frozen["wg_delivery_outbox"][1], "state=blocked")

			hits := f.failClosureAt("UPDATE wg_delivery_samples SET consumed")
			f.closeUntil(shape.windowEnd(1))
			require.Positive(t, hits.Load(), "the injected fault fired inside the average's closure")
			require.NotEmpty(t, f.reportedErrors())
			require.Equal(t, frozen, f.unrelatedState(effects), "the failed closure changed nothing of the other group")

			f.localFault.Reject.Store(nil)
			f.closeUntil(shape.windowEnd(1))
			f.deliver()
			require.Equal(t, shape.want(0, 1), f.averageRows(shape), "the average closed once, exactly")
			require.Equal(t, frozen, f.unrelatedState(effects), "the retry and its delivery left the other group alone")

			f.restart()
			f.closeUntil(shape.windowEnd(1))
			f.deliver()
			require.Equal(t, shape.want(0, 1), f.averageRows(shape))
			require.Equal(t, frozen, f.unrelatedState(effects), "a restart left the other group's bytes unchanged")
		})
	}
}
