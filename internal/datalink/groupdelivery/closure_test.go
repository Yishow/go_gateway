package groupdelivery

import (
	"database/sql"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"

	"github.com/stretchr/testify/require"
)

var testDestination = FrozenDestination{
	Scope: "scope-X", ConnectorID: "connector-1", ConnectorRevision: "crev-1",
	Database: "db", TableSchema: "public", TableName: "readings", DedupeCapability: "none",
}

func rowBucket(entity string, start time.Time, counter uint64) ClosedBucket {
	outcome := snapshot.Outcome{
		Kind: snapshot.OutcomeRow, EntityKey: entity, RecordID: "record-" + entity + start.Format("150405"),
		EffectKey: "effect-" + entity + start.Format("150405"), BucketStart: start, BucketEnd: start.Add(10 * time.Second),
		Members: []snapshot.MemberResult{{
			MemberKey: "temp", Status: snapshot.MemberOK,
			Sample: &snapshot.Sample{SampleID: "s-" + entity, MemberKey: "temp", ObservedAt: start.Add(5 * time.Second), Quality: schema.QualityGood, Value: measurement.NewUint64(counter)},
		}},
	}
	row := dbtarget.EncodedRow{
		RecordID: outcome.RecordID, EffectKey: outcome.EffectKey, EntityKey: entity, BucketStart: start,
		Cells: []dbtarget.EncodedCell{
			{Column: "record_id", Value: outcome.RecordID},
			{Column: "counter", Value: int64(counter)},
			{Column: "label", Value: "batch-001"},
			{Column: "ratio", Value: 0.1},
			{Column: "running", Value: true},
			{Column: "optional", Value: nil},
			{Column: "bucket_start", Value: start},
		},
	}
	return ClosedBucket{Outcome: outcome, Row: &row}
}

func silentBucket(entity string, start time.Time) ClosedBucket {
	return ClosedBucket{Outcome: snapshot.Outcome{
		Kind: snapshot.OutcomeNoData, EntityKey: entity, RecordID: "silent-" + start.Format("150405"),
		BucketStart: start, BucketEnd: start.Add(10 * time.Second), Reason: snapshot.ReasonNoSamples,
		Members: []snapshot.MemberResult{{MemberKey: "temp", Status: snapshot.MemberMissing, Reason: snapshot.ReasonNoSamples}},
	}}
}

func skippedOutcome(start time.Time) snapshot.Outcome {
	return snapshot.Outcome{
		Kind: snapshot.OutcomeSkipped, RecordID: "skip-" + start.Format("150405"), BucketStart: start, BucketEnd: start.Add(10 * time.Second),
		Reason:  snapshot.ReasonIncompleteRequired,
		Members: []snapshot.MemberResult{{MemberKey: "temp", Status: snapshot.MemberBad, Reason: "read-failed"}},
	}
}

func openCount(t *testing.T, db *sql.DB, table, where string) int {
	t.Helper()
	var n int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM `+table+` WHERE `+where).Scan(&n))
	return n
}

func seedSamples(t *testing.T, store *Store) {
	t.Helper()
	require.NoError(t, store.AppendSample(t.Context(), testKey, t0, sampleAt("s-a", "temp", 4, 1)))
	require.NoError(t, store.AppendSample(t.Context(), testKey, t0, sampleAt("s-b", "temp", 8, 2)))
	require.NoError(t, store.AppendSample(t.Context(), testKey, at(10), sampleAt("s-next", "temp", 14, 3)))
}

func TestAtomicRowOutboxCheckpointCommitsEverythingTogether(t *testing.T) {
	store, db, _ := newTestStore(t)
	seedSamples(t, store)
	closure := Closure{Destination: testDestination, NextClose: at(10), Buckets: []ClosedBucket{rowBucket("", t0, 9007199254740993)}}

	require.NoError(t, store.CommitClosure(t.Context(), testKey, closure))

	require.Equal(t, 1, openCount(t, db, "wg_delivery_outbox", "state = 'pending'"))
	require.Equal(t, 1, openCount(t, db, "wg_delivery_buckets", "kind = 'row'"))
	require.Equal(t, 2, openCount(t, db, "wg_delivery_samples", "consumed = 1"), "only the closed bucket's samples are consumed")
	require.Equal(t, 1, openCount(t, db, "wg_delivery_samples", "consumed = 0"))

	var scope, connector, revision, table, capability, effect, digest, partition string
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT destination_scope, connector_id, connector_revision, table_name,
		dedupe_capability, effect_key, payload_digest, partition_key FROM wg_delivery_outbox`).
		Scan(&scope, &connector, &revision, &table, &capability, &effect, &digest, &partition))
	require.Equal(t, "scope-X", scope)
	require.Equal(t, "connector-1", connector)
	require.Equal(t, "crev-1", revision, "the destination identity is frozen with the row")
	require.Equal(t, "readings", table)
	require.Equal(t, "none", capability)
	require.Equal(t, closure.Buckets[0].Outcome.EffectKey, effect)
	require.NotEmpty(t, digest)
	require.NotEmpty(t, partition)

	restored, err := store.Restore(t.Context(), testKey)
	require.NoError(t, err)
	require.Equal(t, at(10), restored.NextClose)
	require.Len(t, restored.Samples, 1)
	require.Equal(t, "s-next", restored.Samples[0].SampleID, "a closed bucket's samples are never reused")
}

func TestAtomicRowOutboxCheckpointPayloadKeepsExactTypedCells(t *testing.T) {
	store, db, _ := newTestStore(t)
	require.NoError(t, store.CommitClosure(t.Context(), testKey,
		Closure{Destination: testDestination, NextClose: at(10), Buckets: []ClosedBucket{rowBucket("plant-1", t0, 9007199254740993)}}))

	var payload string
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT payload FROM wg_delivery_outbox`).Scan(&payload))
	row, err := DecodeRowPayload([]byte(payload))
	require.NoError(t, err)
	require.Equal(t, "plant-1", row.EntityKey)
	require.Equal(t, t0, row.BucketStart)
	cells := map[string]any{}
	for _, cell := range row.Cells {
		cells[cell.Column] = cell.Value
	}
	require.Equal(t, int64(9007199254740993), cells["counter"], "an integer beyond 2^53 survives the payload")
	require.Equal(t, "batch-001", cells["label"])
	require.Equal(t, 0.1, cells["ratio"])
	require.Equal(t, true, cells["running"])
	require.Contains(t, cells, "optional")
	require.Nil(t, cells["optional"], "NULL stays NULL")
	require.Equal(t, t0, cells["bucket_start"])
}

func TestAtomicRowOutboxCheckpointSilentAndSkippedBucketsArePersisted(t *testing.T) {
	store, db, _ := newTestStore(t)
	skipped := ClosedBucket{Outcome: snapshot.Outcome{
		Kind: snapshot.OutcomeSkipped, RecordID: "skip-1", BucketStart: at(10), BucketEnd: at(20), Reason: snapshot.ReasonIncompleteRequired,
		Members: []snapshot.MemberResult{{MemberKey: "temp", Status: snapshot.MemberBad, Reason: "read-failed"}},
	}}
	require.NoError(t, store.CommitClosure(t.Context(), testKey,
		Closure{Destination: testDestination, NextClose: at(30), Buckets: []ClosedBucket{silentBucket("", t0), skipped}}))

	require.Equal(t, 0, countRows(t, db, "wg_delivery_outbox"), "no SQL row exists for silent or skipped buckets")
	require.Equal(t, 1, openCount(t, db, "wg_delivery_buckets", "kind = 'no_data' AND reason = 'no-samples'"))
	require.Equal(t, 1, openCount(t, db, "wg_delivery_buckets", "kind = 'skipped' AND reason = 'incomplete-required-member'"))
	var members string
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT members FROM wg_delivery_buckets WHERE kind = 'skipped'`).Scan(&members))
	require.Contains(t, members, "read-failed", "per-member causes are kept for the status view")
	restored, err := store.Restore(t.Context(), testKey)
	require.NoError(t, err)
	require.Equal(t, at(30), restored.NextClose)
}

func TestAtomicRowOutboxCheckpointFailureAtEveryBoundaryRollsBackEverything(t *testing.T) {
	boundaries := map[string]string{
		"outbox insert":     `CREATE TRIGGER inject BEFORE INSERT ON wg_delivery_outbox BEGIN SELECT RAISE(ABORT, 'injected commit failure'); END`,
		"bucket insert":     `CREATE TRIGGER inject BEFORE INSERT ON wg_delivery_buckets BEGIN SELECT RAISE(ABORT, 'injected commit failure'); END`,
		"samples consumed":  `CREATE TRIGGER inject BEFORE UPDATE ON wg_delivery_samples BEGIN SELECT RAISE(ABORT, 'injected commit failure'); END`,
		"checkpoint insert": `CREATE TRIGGER inject BEFORE INSERT ON wg_delivery_checkpoints BEGIN SELECT RAISE(ABORT, 'injected commit failure'); END`,
	}
	for name, trigger := range boundaries {
		store, db, _ := newTestStore(t)
		seedSamples(t, store)
		closure := Closure{Destination: testDestination, NextClose: at(10), Buckets: []ClosedBucket{rowBucket("", t0, 1), silentBucket("plant-2", t0)}}
		_, err := db.ExecContext(t.Context(), trigger)
		require.NoError(t, err, name)

		err = store.CommitClosure(t.Context(), testKey, closure)
		require.ErrorContains(t, err, "injected commit failure", name)
		require.Zero(t, countRows(t, db, "wg_delivery_outbox"), name)
		require.Zero(t, countRows(t, db, "wg_delivery_buckets"), name)
		require.Zero(t, countRows(t, db, "wg_delivery_checkpoints"), name)
		require.Equal(t, 3, openCount(t, db, "wg_delivery_samples", "consumed = 0"), name)

		// After the fault clears the same closure replays cleanly.
		_, err = db.ExecContext(t.Context(), `DROP TRIGGER inject`)
		require.NoError(t, err, name)
		require.NoError(t, store.CommitClosure(t.Context(), testKey, closure), name)
		require.Equal(t, 1, countRows(t, db, "wg_delivery_outbox"), name)
		require.Equal(t, 2, countRows(t, db, "wg_delivery_buckets"), name)
		require.Equal(t, 2, openCount(t, db, "wg_delivery_samples", "consumed = 1"), name)
	}
}

func TestAtomicRowOutboxCheckpointRepeatedCommitIsIdempotentAndConflictsAreRefused(t *testing.T) {
	store, db, _ := newTestStore(t)
	closure := Closure{Destination: testDestination, NextClose: at(10), Buckets: []ClosedBucket{rowBucket("", t0, 7)}}
	require.NoError(t, store.CommitClosure(t.Context(), testKey, closure))
	require.NoError(t, store.CommitClosure(t.Context(), testKey, closure), "the commit response may be lost; repeating is safe")
	require.Equal(t, 1, countRows(t, db, "wg_delivery_outbox"))
	require.Equal(t, 1, countRows(t, db, "wg_delivery_buckets"))

	different := rowBucket("", t0, 8) // same record and effect key, different values
	different.Outcome.RecordID, different.Outcome.EffectKey = closure.Buckets[0].Outcome.RecordID, closure.Buckets[0].Outcome.EffectKey
	different.Row.RecordID, different.Row.EffectKey = different.Outcome.RecordID, different.Outcome.EffectKey
	err := store.CommitClosure(t.Context(), testKey, Closure{Destination: testDestination, NextClose: at(10), Buckets: []ClosedBucket{different}})
	require.ErrorIs(t, err, ErrEffectConflict)
	require.Equal(t, 1, countRows(t, db, "wg_delivery_outbox"))
}

func TestAtomicRowOutboxCheckpointNeverMovesBackwardsAndValidatesInput(t *testing.T) {
	store, db, _ := newTestStore(t)
	require.NoError(t, store.CommitClosure(t.Context(), testKey,
		Closure{Destination: testDestination, NextClose: at(30), Buckets: []ClosedBucket{silentBucket("", at(20))}}))
	require.NoError(t, store.CommitClosure(t.Context(), testKey,
		Closure{Destination: testDestination, NextClose: at(10), Buckets: []ClosedBucket{silentBucket("", t0)}}))
	var next string
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT next_close FROM wg_delivery_checkpoints`).Scan(&next))
	require.Equal(t, timestamp(at(30)), next)

	require.ErrorIs(t, store.CommitClosure(t.Context(), GroupKey{}, Closure{Destination: testDestination, NextClose: at(10)}), ErrInvalidGroupKey)
	require.ErrorIs(t, store.CommitClosure(t.Context(), testKey, Closure{Destination: testDestination}), ErrInvalidClosure)
	missingRow := ClosedBucket{Outcome: rowBucket("", t0, 1).Outcome} // a row outcome without encoded values
	require.ErrorIs(t, store.CommitClosure(t.Context(), testKey,
		Closure{Destination: testDestination, NextClose: at(10), Buckets: []ClosedBucket{missingRow}}), ErrInvalidClosure)
	noDestination := Closure{NextClose: at(10), Buckets: []ClosedBucket{rowBucket("", t0, 1)}}
	require.ErrorIs(t, store.CommitClosure(t.Context(), testKey, noDestination), ErrInvalidClosure)
}
