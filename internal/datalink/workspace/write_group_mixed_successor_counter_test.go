package workspace

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"

	"github.com/stretchr/testify/require"
)

func mixedSuccessorDeltaSource(t *testing.T, now *time.Time) (*circuitProfileHarness, *WriteGroup) {
	t.Helper()
	h := newCircuitProfileHarness(t, filepath.Join(t.TempDir(), "workspace.db"))
	target := filepath.Join(t.TempDir(), "target.db")
	db, err := sql.Open("sqlite", target)
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), `PRAGMA user_version=0`)
	require.NoError(t, err)
	require.NoError(t, db.Close())
	setSQLiteConnectorConfig(t, h.db, h.fixture.connectorID, target)
	h.service.WithClock(func() time.Time { return *now })
	id := seedMeasurement(t, h, h.primary(), "counter")
	_, err = h.db.ExecContext(t.Context(), `UPDATE measurement_definitions SET unit='kWh',series_epoch='epoch-1',numeric_encoding='float32',counter_policy='{}' WHERE id=?`, id)
	require.NoError(t, err)
	draft := mixedDraft(h, WriteGroupMemberCalculationDelta, 15, "mixed_counter")
	draft.LayoutVersion = WriteGroupLayoutManagedWideProfileV1
	draft.Destination.Database = target
	draft.Members[0].MeasurementID = &id
	created, err := h.service.Create(t.Context(), h.mutation(t.Context(), draft))
	require.NoError(t, err)
	group := mixedHistoryConfirm(t, h, created.Group)
	applied, err := h.service.Apply(t.Context(), group.ID, mixedHistoryRequest(t, h, group))
	require.NoError(t, err)
	return h, applied.Group
}

// Close a real first delta window using the existing assembler and atomic
// delivery store. A successor must never inherit this delivery checkpoint.
func mixedSuccessorBaseline(t *testing.T, h *circuitProfileHarness, source *WriteGroup, end time.Time) func(string) {
	t.Helper()
	ctx := t.Context()
	definition, err := measurement.NewSQLRepository(h.db).GetByID(ctx, *source.Members[0].MeasurementID)
	require.NoError(t, err)
	proof, err := unscaledCounterBindingProof(source, source.Members[0], schema.DataTypeFloat32, definition)
	require.NoError(t, err)
	binding := aggregation.CounterBinding{WorkspaceID: source.WorkspaceID, GroupID: source.ID, ProfileRevision: source.Revision,
		MemberKey: proof.Member.SampleMemberKey(), MeasurementID: proof.MeasurementID, DefinitionRevision: proof.DefinitionRevision,
		SourceBindingRevision: proof.SourceBindingRevision, SeriesEpoch: proof.SeriesEpoch, Kind: aggregation.CounterKindMonotonic, Unit: proof.Unit}
	value, err := measurement.NewFloat64(123.5)
	require.NoError(t, err)
	start := end
	end = start.Truncate(15 * time.Second).Add(15 * time.Second)
	memberKey := binding.MemberKey
	assembler, err := snapshot.NewAssembler(snapshot.Config{WorkspaceID: source.WorkspaceID, GroupID: source.ID, GroupRevision: source.Revision,
		DestinationScope: "successor-source", StartsAt: start, Mode: snapshot.ModeMixedWide,
		Members:  []snapshot.Member{{Key: memberKey, Required: true, Calculation: snapshot.CalculationDelta, Interval: 15 * time.Second, SourceValueType: schema.DataTypeFloat32}},
		Counters: map[string]aggregation.CounterBinding{memberKey: binding}, IncompletePolicy: snapshot.IncompletePartial,
		Storage: snapshot.StorageCapabilities{NullableValues: true, MemberQuality: true}})
	require.NoError(t, err)
	ref, err := proof.Reference()
	require.NoError(t, err)
	sample := snapshot.Sample{SampleID: "old-reading", MemberKey: memberKey, ObservedAt: end.Add(-time.Second), Quality: schema.QualityGood, Value: value, CounterProof: ref, SourceRevision: proof.SourceRevision, MappingRevision: proof.MappingRevision}
	store := groupdelivery.NewStore(h.db)
	key := groupdelivery.GroupKey{WorkspaceID: source.WorkspaceID, GroupID: source.ID, GroupRevision: source.Revision}
	require.NoError(t, store.AppendSample(ctx, key, start.Truncate(15*time.Second), sample))
	require.Equal(t, snapshot.OfferSelected, assembler.Offer(sample, sample.ObservedAt).Outcome)
	before, err := assembler.MixedCheckpoint()
	require.NoError(t, err)
	plan := assembler.Plan(end)
	require.NoError(t, plan.Err)
	after, err := assembler.MixedCheckpointForPlan(plan)
	require.NoError(t, err)
	checkpoints := plan.CounterCheckpoints()
	require.Len(t, checkpoints, 1)
	payload, err := aggregation.EncodeCounterCheckpoint(checkpoints[memberKey])
	require.NoError(t, err)
	digest := fmt.Sprintf("%x", sha256.Sum256(payload))
	restored, err := store.Restore(ctx, key)
	require.NoError(t, err)
	buckets := make([]groupdelivery.ClosedBucket, len(plan.Outcomes))
	for i, outcome := range plan.Outcomes {
		require.Equal(t, snapshot.OutcomeSkipped, outcome.Kind, "first required delta has no baseline result")
		buckets[i] = groupdelivery.ClosedBucket{Outcome: outcome}
	}
	require.NoError(t, store.CommitClosure(ctx, key, groupdelivery.Closure{
		Destination:   groupdelivery.FrozenDestination{Scope: "successor-source", ConnectorID: source.Destination.ConnectorID, ConnectorRevision: source.Destination.ConnectorRevision, Database: source.Destination.Database, TableSchema: source.Destination.TableSchema, TableName: source.Destination.TableName},
		CounterProofs: []counterproof.CounterBindingReference{ref},
		Buckets:       buckets, Baselines: []groupdelivery.CounterBaseline{{MemberKey: memberKey, Payload: payload}},
		Mixed: &groupdelivery.MixedClosureProgress{Before: before, After: after, DescriptorDigest: restored.DescriptorDigest}}))
	require.Zero(t, h.countRows("SELECT COUNT(*) FROM wg_delivery_checkpoints"), "mixed closure never fabricates a legacy cursor")
	for name, corrupt := range map[string]func(*sql.Tx){
		"missing progress": func(tx *sql.Tx) {
			_, err := tx.ExecContext(ctx, `DELETE FROM wg_delivery_mixed_progress`)
			require.NoError(t, err)
		},
		"progress replaced with legacy cursor": func(tx *sql.Tx) {
			_, err := tx.ExecContext(ctx, `DELETE FROM wg_delivery_mixed_progress`)
			require.NoError(t, err)
			_, err = tx.ExecContext(ctx, `INSERT INTO wg_delivery_checkpoints(group_id,group_revision,next_close,updated_at) VALUES(?,?,?,?)`, source.ID, source.Revision, end.Format(time.RFC3339Nano), "now")
			require.NoError(t, err)
		},
		"checkpoint hash": func(tx *sql.Tx) {
			_, err := tx.ExecContext(ctx, `UPDATE wg_delivery_mixed_progress SET checkpoint_digest='bad'`)
			require.NoError(t, err)
		},
		"descriptor hash": func(tx *sql.Tx) {
			_, err := tx.ExecContext(ctx, `UPDATE wg_delivery_mixed_progress SET descriptor_digest='bad'`)
			require.NoError(t, err)
		},
		"foreign scope": func(tx *sql.Tx) {
			wrong := after
			wrong.Progress.WorkspaceID = "foreign"
			raw, err := json.Marshal(wrong)
			require.NoError(t, err)
			_, err = tx.ExecContext(ctx, `UPDATE wg_delivery_mixed_progress SET checkpoint=?,checkpoint_digest=?`, string(raw), fmt.Sprintf("%x", sha256.Sum256(raw)))
			require.NoError(t, err)
		},
		"legacy cursor": func(tx *sql.Tx) {
			_, err := tx.ExecContext(ctx, `INSERT INTO wg_delivery_checkpoints(group_id,group_revision,next_close,updated_at) VALUES(?,?,?,?)`, source.ID, source.Revision, end.Format(time.RFC3339Nano), "now")
			require.NoError(t, err)
		},
	} {
		t.Run(name, func(t *testing.T) {
			tx, err := h.db.BeginTx(ctx, nil)
			require.NoError(t, err)
			defer tx.Rollback()
			corrupt(tx)
			require.Error(t, requireNoFinalizedOverlap(ctx, tx, &cutoverPlan{req: CutoverRequest{Boundary: end.Add(time.Minute)}, ended: []CutoverEnd{{GroupID: source.ID, Revision: source.Revision}}}))
		})
	}
	t.Run("finalized overlap", func(t *testing.T) {
		tx, err := h.db.BeginTx(ctx, nil)
		require.NoError(t, err)
		defer tx.Rollback()
		require.ErrorIs(t, requireNoFinalizedOverlap(ctx, tx, &cutoverPlan{req: CutoverRequest{Boundary: end.Add(-time.Nanosecond)}, ended: []CutoverEnd{{GroupID: source.ID, Revision: source.Revision}}}), ErrCutoverFinalizedOverlap)
	})
	return func(successorID string) {
		var got, hash string
		require.NoError(t, h.db.QueryRowContext(ctx, `SELECT payload,payload_digest FROM wg_delivery_counter_baselines WHERE group_id=?`, source.ID).Scan(&got, &hash))
		require.Equal(t, string(payload), got)
		require.Equal(t, digest, hash)
		var n int
		require.NoError(t, h.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM wg_delivery_counter_baselines WHERE group_id=?`, successorID).Scan(&n))
		require.Zero(t, n, "new reserved identity must establish its own delta baseline")
	}
}
