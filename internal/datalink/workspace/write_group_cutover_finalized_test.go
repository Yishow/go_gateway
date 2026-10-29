package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
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

// The source has its real immutable counter proof; the closure uses the real
// journal/outbox/baseline/checkpoint transaction, not independently seeded rows.
func finalizedClosure(t *testing.T, f counterApplyFixture, g *WriteGroup, end time.Time) (*groupdelivery.Store, groupdelivery.GroupKey, groupdelivery.Closure) {
	t.Helper()
	var raw string
	require.NoError(t, f.db.QueryRow(`SELECT payload FROM wg_runtime_versions WHERE group_id=? AND group_revision=?`, g.ID, g.AppliedRevision).Scan(&raw))
	var saved WriteGroupAppliedSnapshot
	require.NoError(t, json.Unmarshal([]byte(raw), &saved))
	require.Len(t, saved.CounterBindings, 1)
	var sample snapshot.Sample
	var binding aggregation.CounterBinding
	for _, proof := range saved.CounterBindings {
		ref, err := proof.Reference()
		require.NoError(t, err)
		sample = snapshot.Sample{SampleID: "s-1", MemberKey: proof.Member.SampleMemberKey(), ObservedAt: end.Add(-7 * time.Second), Quality: schema.QualityGood, Value: measurement.NewInt64(7), CounterProof: ref, SourceRevision: proof.SourceRevision, MappingRevision: proof.MappingRevision}
		binding = aggregation.CounterBinding{WorkspaceID: g.WorkspaceID, GroupID: g.ID, ProfileRevision: g.AppliedRevision, MemberKey: sample.MemberKey, MeasurementID: proof.MeasurementID, DefinitionRevision: proof.DefinitionRevision, SourceBindingRevision: proof.SourceBindingRevision, SeriesEpoch: proof.SeriesEpoch, Kind: aggregation.CounterKindMonotonic, Unit: proof.Unit}
	}
	start := end.Add(-15 * time.Second)
	store := groupdelivery.NewStore(f.db)
	key := groupdelivery.GroupKey{WorkspaceID: g.WorkspaceID, GroupID: g.ID, GroupRevision: g.AppliedRevision}
	require.NoError(t, store.AppendSample(t.Context(), key, start, sample))
	baseline, err := aggregation.EncodeCounterCheckpoint(aggregation.CounterCheckpoint{Binding: binding, WindowEnd: end, Reading: aggregation.CounterReading{SampleID: sample.SampleID, ObservedAt: sample.ObservedAt, Value: sample.Value}})
	require.NoError(t, err)
	outcome := snapshot.Outcome{Kind: snapshot.OutcomeRow, RecordID: "record-1", EffectKey: "e-1", BucketStart: start, BucketEnd: end, Members: []snapshot.MemberResult{{MemberKey: sample.MemberKey, Status: snapshot.MemberOK, Sample: &sample}}}
	row := dbtarget.EncodedRow{RecordID: outcome.RecordID, EffectKey: outcome.EffectKey, BucketStart: start, Cells: []dbtarget.EncodedCell{{Column: "value", Value: int64(7)}, {Column: "observed_at", Value: sample.ObservedAt}}}
	return store, key, groupdelivery.Closure{Destination: groupdelivery.FrozenDestination{Scope: "frozen-source", ConnectorID: g.Destination.ConnectorID, ConnectorRevision: g.Destination.ConnectorRevision, Database: g.Destination.Database, TableSchema: g.Destination.TableSchema, TableName: g.Destination.TableName}, NextClose: end, Buckets: []groupdelivery.ClosedBucket{{Outcome: outcome, Row: &row}}, Baselines: []groupdelivery.CounterBaseline{{MemberKey: sample.MemberKey, Payload: baseline}}}
}

func finalizedState(t *testing.T, db *sql.DB) map[string][]string {
	t.Helper()
	state := map[string][]string{}
	for _, table := range []string{"system_settings", "write_groups", "write_group_versions", "wg_runtime_versions", "write_group_cutovers", "write_group_revision_ends", "write_group_ownership_fences", "wg_delivery_samples", "wg_delivery_buckets", "wg_delivery_outbox", "wg_delivery_checkpoints", "wg_delivery_counter_baselines"} {
		rows, err := db.Query(`SELECT * FROM ` + table + ` ORDER BY rowid`)
		require.NoError(t, err)
		columns, err := rows.Columns()
		require.NoError(t, err)
		for rows.Next() {
			values := make([]any, len(columns))
			ptrs := make([]any, len(columns))
			for i := range values {
				ptrs[i] = &values[i]
			}
			require.NoError(t, rows.Scan(ptrs...))
			raw, err := json.Marshal(values)
			require.NoError(t, err)
			state[table] = append(state[table], string(raw))
		}
		require.NoError(t, rows.Err())
		require.NoError(t, rows.Close())
	}
	return state
}

func TestCutoverFinalizedOverlap(t *testing.T) {
	for _, mode := range []string{"future observation", "post reclaim", "boundary equality", "nonoverlap", "one nanosecond overlap", "unfinalized future"} {
		t.Run(mode, func(t *testing.T) {
			f, _, req := counterCutoverFixture(t)
			anchor := req.Boundary.Add(-22*time.Second - 500*time.Millisecond)
			end := anchor.Add(30 * time.Second)
			if mode == "boundary equality" {
				end = req.Boundary
			}
			if mode == "nonoverlap" {
				end = req.Boundary.Add(-time.Nanosecond)
			}
			if mode == "one nanosecond overlap" {
				end = req.Boundary.Add(time.Nanosecond)
			}
			store, key, closure := finalizedClosure(t, f, f.group, end)
			if mode != "unfinalized future" {
				require.NoError(t, store.CommitClosure(t.Context(), key, closure))
			}
			if mode == "post reclaim" {
				_, err := f.db.Exec(`UPDATE wg_delivery_outbox SET state='sql_committed', committed_at='2000-01-01T00:00:00.000000000Z'`)
				require.NoError(t, err)
				removed, err := store.Reclaim(t.Context(), 0)
				require.NoError(t, err)
				require.EqualValues(t, 3, removed)
			}
			before := finalizedState(t, f.db)
			result, err := f.service.ApplyCutover(t.Context(), req)
			if mode == "boundary equality" || mode == "nonoverlap" || mode == "unfinalized future" {
				require.NoError(t, err)
				require.NotNil(t, result)
				return
			}
			require.ErrorContains(t, err, "recording-cutover-finalized-overlap")
			require.Nil(t, result)
			require.Equal(t, before, finalizedState(t, f.db), "refusal cannot mutate history, frozen payload/digest/destination, consumed state, proof, baseline or revision")
		})
	}
}

// Pause the final configuration transaction after the readiness inspector has
// completed. The closure wins while Apply is in flight, after its review and
// preflight, before the SQL transaction that must recheck closure authority.
type finalizedBarrierRepository struct {
	*SQLRepository
	armed            bool
	entered, release chan struct{}
}

func (r *finalizedBarrierRepository) BeginDatabaseSetup(ctx context.Context) (DatabaseSetupTx, error) {
	if r.armed {
		r.armed = false
		close(r.entered)
		select {
		case <-r.release:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return r.SQLRepository.BeginDatabaseSetup(ctx)
}
func TestCutoverFinalizedClosureWinsAfterPreflight(t *testing.T) {
	f, _, req := counterCutoverFixture(t)
	store, key, closure := finalizedClosure(t, f, f.group, req.Boundary.Add(7500*time.Millisecond))
	repo := &finalizedBarrierRepository{SQLRepository: NewSQLRepository(f.db), entered: make(chan struct{}), release: make(chan struct{})}
	f.service.workspaceSvc.repo = repo
	f.inspector.mutate = func() { repo.armed = true }
	done := make(chan error, 1)
	go func() { _, err := f.service.ApplyCutover(t.Context(), req); done <- err }()
	<-repo.entered
	require.NoError(t, store.CommitClosure(t.Context(), key, closure))
	before := finalizedState(t, f.db)
	close(repo.release)
	require.ErrorContains(t, <-done, "recording-cutover-finalized-overlap")
	require.Equal(t, before, finalizedState(t, f.db))
}

func TestCutoverFinalizedRollbackChecksTargetRevision(t *testing.T) {
	for _, layout := range []WriteGroupLayoutVersion{WriteGroupLayoutManagedLongV1, WriteGroupLayoutManagedWideProfileV1} {
		t.Run(string(layout), func(t *testing.T) {
			f, target, req := counterCutoverFixtureWithTarget(t, layout)
			forward, err := f.service.ApplyCutover(t.Context(), req)
			require.NoError(t, err)
			rollback := reviewCounterRollback(t, f, forward)
			target.AppliedRevision = target.Revision
			// Finalize a real aligned target window beyond the proposed rollback.
			// A forward cutover already exists, so closure must carry its restored view.
			store, key, closure := finalizedClosure(t, f, target, rollback.Boundary.Add(15*time.Second))
			restored, err := store.Restore(t.Context(), key)
			require.NoError(t, err)
			closure.OwnershipView, closure.HandoffView = restored.OwnershipView, restored.HandoffView
			require.NoError(t, store.CommitClosure(t.Context(), key, closure))
			before := finalizedState(t, f.db)
			_, err = f.service.ApplyCutover(t.Context(), rollback)
			require.ErrorContains(t, err, "recording-cutover-finalized-overlap")
			require.Equal(t, before, finalizedState(t, f.db))
		})
	}
}

func TestCutoverFinalizedMissingCheckpointCannotEraseHistory(t *testing.T) {
	for _, evidence := range []string{"samples", "buckets", "outbox", "counter_baselines"} {
		t.Run(evidence, func(t *testing.T) {
			f, _, req := counterCutoverFixture(t)
			store, key, closure := finalizedClosure(t, f, f.group, req.Boundary.Add(7500*time.Millisecond))
			require.NoError(t, store.CommitClosure(t.Context(), key, closure))
			// A missing cursor is not a normal initial state if ANY retained closure
			// evidence remains, even if retention has removed all other evidence.
			for _, table := range []string{"checkpoints", "samples", "buckets", "outbox", "counter_baselines"} {
				if table != evidence {
					_, err := f.db.Exec(`DELETE FROM wg_delivery_` + table)
					require.NoError(t, err)
				}
			}
			before := finalizedState(t, f.db)
			_, err := f.service.ApplyCutover(t.Context(), req)
			require.ErrorIs(t, err, groupdelivery.ErrOwnershipUnavailable)
			require.Equal(t, before, finalizedState(t, f.db))
		})
	}
}

func TestCutoverFinalizedUnreadableCheckpointFailsSafely(t *testing.T) {
	for _, mode := range []string{"malformed", "missing table"} {
		t.Run(mode, func(t *testing.T) {
			f, _, req := counterCutoverFixture(t)
			store, key, closure := finalizedClosure(t, f, f.group, req.Boundary.Add(7500*time.Millisecond))
			require.NoError(t, store.CommitClosure(t.Context(), key, closure))
			statement := `UPDATE wg_delivery_checkpoints SET next_close='private-driver-dsn-and-value'`
			if mode == "missing table" {
				statement = `ALTER TABLE wg_delivery_checkpoints RENAME TO unavailable_checkpoints`
			}
			_, err := f.db.Exec(statement)
			require.NoError(t, err)
			_, err = f.service.ApplyCutover(t.Context(), req)
			require.ErrorIs(t, err, groupdelivery.ErrOwnershipUnavailable)
			require.NotContains(t, err.Error(), "private-driver")
			require.NotContains(t, err.Error(), "SQLITE")
			var count int
			require.NoError(t, f.db.QueryRow(`SELECT COUNT(*) FROM write_group_cutovers`).Scan(&count))
			require.Zero(t, count)
		})
	}
}
