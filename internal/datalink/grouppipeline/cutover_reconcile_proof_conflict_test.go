package grouppipeline

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

// These are original descriptors fixed before the first accepted sample, not
// an attempt to repair or reinterpret already accepted proof in recovery.
func originalConflictFixture(t *testing.T, shape averageShape, mode string) (*ownershipPipelineFixture, *changedCounterReader, *workspace.WriteGroupAppliedSnapshot) {
	t.Helper()
	f, reader := reductionHandoffFixture(t, shape, "cumulative")
	for _, owner := range []*workspace.WriteGroupAppliedSnapshot{f.source, f.target} {
		for key, original := range owner.CounterBindings {
			facts := original.CounterBindingFacts
			facts.SeriesEpoch = "epoch-1"
			if owner == f.target {
				switch mode {
				case "epoch-and-unit":
					facts.SeriesEpoch, facts.Unit = "epoch-2", "MWh"
				case "epoch":
					facts.SeriesEpoch = "epoch-2"
				case "unit":
					facts.Unit = "MWh"
				case "definition-proof":
					facts.DefinitionRevision = "definition-next"
				}
			}
			proof, err := workspace.NewCounterBindingProof(facts)
			require.NoError(t, err)
			owner.CounterBindings[key] = proof
		}
		raw, err := json.Marshal(owner)
		require.NoError(t, err)
		_, err = f.local.Exec(`UPDATE wg_runtime_versions SET payload=? WHERE group_id=?`, string(raw), owner.GroupID)
		require.NoError(t, err)
	}
	// Independent existing legacy scope: it has no counter proof to reseal and
	// its real SQLite output must progress while only the cutover pair is blocked.
	healthyGroup := &workspace.WriteGroup{
		ID: "healthy-proof-scope", WorkspaceID: f.source.WorkspaceID, Revision: "healthy-rev", AppliedRevision: "healthy-rev", Status: workspace.WriteGroupStatusReady,
		Members:     []workspace.WriteGroupMember{{DeviceID: "healthy-device", PointID: "healthy-point", TagID: "healthy-tag", TargetColumn: "reading", Required: true, SourceRevision: "healthy-source", MappingRevision: "healthy-map"}},
		Destination: workspace.WriteGroupDestination{ConnectorID: "connector-1", ConnectorRevision: "crev-1", TableSchema: "main", TableName: "healthy_proof_readings", StorageStrategy: workspace.WriteGroupStorageStrategyCustom},
		RowPolicy:   workspace.WriteGroupRowPolicy{IntervalSeconds: 10}, WritePolicy: workspace.WriteGroupWritePolicy{Mode: "append", DedupeCapability: groupdelivery.DedupeNone},
	}
	_, err := f.plainTarget().Exec(`CREATE TABLE healthy_proof_readings(reading TEXT)`)
	require.NoError(t, err)
	healthy := snapshotOf(healthyGroup, at(0))
	healthy.RuntimeLayout = &workspace.WriteGroupRuntimeLayout{Dialect: dbtarget.SQLDialectSQLite, Columns: []dbtarget.ColumnInfo{{Name: "reading", DataType: "TEXT", Nullable: true}}, TagTypes: map[string]schema.DataType{"healthy-tag": schema.DataTypeUint64}}
	require.NoError(t, f.pipe.saveRuntimeSnapshot(t.Context(), healthy))
	f.groups.groups = append(f.groups.groups, healthyGroup)
	f.groups.snapshots[healthy.GroupID] = []*workspace.WriteGroupAppliedSnapshot{healthy}
	f.pipe = New(f.pipe.deps, f.pipe.config)
	require.NoError(t, f.pipe.Reconcile(t.Context()))
	require.NotNil(t, f.pipe.boundaries[boundaryKey(healthy.GroupID, healthy.AppliedRevision)], "%+v", f.pipe.Status())
	return f, reader, healthy
}

// Exclude the legitimate handoff audit and independent healthy group; every
// original or derived fact belonging to the blocked pair must remain unchanged.
func conflictPairState(t *testing.T, f *ownershipPipelineFixture) map[string][]string {
	t.Helper()
	state := map[string][]string{}
	for _, table := range []string{"wg_delivery_samples", "wg_delivery_checkpoints", "wg_delivery_counter_baselines", "wg_delivery_outbox", "wg_runtime_versions", "wg_counter_proof_associations"} {
		rows, err := f.local.Query(`SELECT * FROM `+table+` WHERE group_id IN (?,?) ORDER BY rowid`, f.source.GroupID, f.target.GroupID)
		require.NoError(t, err)
		columns, err := rows.Columns()
		require.NoError(t, err)
		for rows.Next() {
			values, pointers := make([]any, len(columns)), make([]any, len(columns))
			for i := range values {
				pointers[i] = &values[i]
			}
			require.NoError(t, rows.Scan(pointers...))
			raw, err := json.Marshal(values)
			require.NoError(t, err)
			state[table] = append(state[table], string(raw))
		}
		require.NoError(t, rows.Err())
		require.NoError(t, rows.Close())
	}
	return state
}

func TestCutoverReconcileProofConflictRemainsDurableAndScoped(t *testing.T) {
	for _, shape := range []averageShape{shapeMeanLong, shapeMeanWide} {
		for _, mode := range []string{"epoch-and-unit", "epoch", "unit", "definition-proof", "missing-authority", "authority-read-failure", "receipt-write-failure"} {
			t.Run(string(shape)+"/"+mode, func(t *testing.T) {
				f, reader, healthy := originalConflictFixture(t, shape, mode)
				f.clock = at(-2)
				require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(f.source.Group, "baseline", at(-1), uint64(5))))
				source := f.pipe.boundaries[boundaryKey(f.source.GroupID, f.source.AppliedRevision)]
				require.NoError(t, source.boundary.Tick(t.Context(), at(0)))
				f.deliver() // Settle the already finalized baseline before freezing blocked-pair state.
				f.clock = at(7)
				require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(f.source.Group, "before", at(9), uint64(20))))
				input := envelope(f.source.Group, "s-1", at(11), uint64(100))
				require.NoError(t, f.pipe.AcceptSample(t.Context(), input))
				require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(healthy.Group, "healthy", at(11), uint64(7))))
				accepted, err := f.store.Restore(t.Context(), groupKey(f.source))
				require.NoError(t, err)
				require.Len(t, accepted.Samples, 2)
				require.Equal(t, "s-1", accepted.Samples[1].SampleID)
				require.Len(t, accepted.Baselines, 1)
				for _, proof := range f.source.CounterBindings {
					require.Equal(t, "epoch-1", proof.SeriesEpoch)
					require.Equal(t, "kWh", proof.Unit)
				}
				f.commit()
				f.source.EndsAt = f.boundary
				f.groups.groups = append(f.groups.groups, f.target.Group)
				f.groups.snapshots[f.target.GroupID] = []*workspace.WriteGroupAppliedSnapshot{f.target}
				if mode == "missing-authority" {
					_, err = f.local.Exec(`DELETE FROM wg_runtime_versions WHERE group_id=?`, f.target.GroupID)
					require.NoError(t, err)
				}
				before, originalRows := conflictPairState(t, f), handoffOriginalRows(t, f)
				wantErr, wantState := error(groupdelivery.ErrHandoffProofConflict), "quarantined"
				if strings.Contains(mode, "authority") || mode == "receipt-write-failure" {
					wantErr, wantState = groupdelivery.ErrOwnershipUnavailable, "pending"
				}
				faultCalls := 0
				reject := func(query string) error {
					if mode == "authority-read-failure" && strings.Contains(query, "SELECT payload FROM wg_runtime_versions WHERE workspace_id=?") ||
						mode == "receipt-write-failure" && strings.Contains(query, "INSERT INTO wg_delivery_sample_handoffs") {
						faultCalls++
						return errors.New("private-proof-db-path-and-SQL")
					}
					return nil
				}
				f.localFault.Reject.Store(&reject)
				t.Cleanup(func() { f.localFault.Reject.Store(nil) })
				f.clock = at(12)
				var receiptID string
				for attempt := 0; attempt < 2; attempt++ {
					err = f.pipe.Reconcile(t.Context())
					require.ErrorIs(t, err, wantErr)
					for _, secret := range []string{"private-", "epoch-1", "epoch-2", "kWh", "MWh", "100"} {
						require.NotContains(t, err.Error(), secret)
					}
					require.Contains(t, err.Error(), wantErr.Error(), "typed refusal retains a safe public reason")
					var receipts int
					require.NoError(t, f.local.QueryRow(`SELECT COUNT(*) FROM wg_delivery_sample_handoffs`).Scan(&receipts))
					if mode == "receipt-write-failure" {
						require.Zero(t, receipts, "uncommitted transfer cannot masquerade as durable pending")
					} else {
						require.Equal(t, 1, receipts)
						var id, state, raw string
						var targetSequence sql.NullInt64
						require.NoError(t, f.local.QueryRow(`SELECT transfer_id,state,target_seq,record FROM wg_delivery_sample_handoffs`).Scan(&id, &state, &targetSequence, &raw))
						require.Equal(t, wantState, state)
						require.False(t, targetSequence.Valid, "blocked receipt is not target acceptance")
						var receipt struct {
							SampleID, Reason string
							SampleProof      counterproof.CounterBindingReference
						}
						require.NoError(t, json.Unmarshal([]byte(raw), &receipt))
						require.Equal(t, "s-1", receipt.SampleID)
						require.Equal(t, wantErr.Error(), receipt.Reason)
						require.Equal(t, accepted.Samples[1].CounterProof, receipt.SampleProof)
						if attempt == 0 {
							receiptID = id
						}
						require.Equal(t, receiptID, id, "restart/retry retains stable transfer identity")
					}
					if attempt == 0 {
						// This retry is still before the new cumulative closure deadline;
						// late-after-close cannot disguise an unpersisted transfer ACK.
						resendReason := error(groupdelivery.ErrOwnershipUnavailable)
						if mode == "receipt-write-failure" {
							resendReason = groupdelivery.ErrOwnershipConflict
						}
						require.ErrorIs(t, f.pipe.AcceptSample(t.Context(), input), resendReason)
						f.clock = at(70)
					}
					if current := f.pipe.boundaries[boundaryKey(f.source.GroupID, f.source.AppliedRevision)]; current != nil {
						require.Error(t, current.boundary.Tick(t.Context(), f.clock))
					}
					if current := f.pipe.boundaries[boundaryKey(f.target.GroupID, f.target.AppliedRevision)]; current != nil {
						require.Error(t, current.boundary.Tick(t.Context(), f.clock))
					}
					require.Equal(t, originalRows, handoffOriginalRows(t, f))
					require.Equal(t, before, conflictPairState(t, f), "no accepted bytes, baseline, checkpoint, proof or outbox may advance")
					require.Zero(t, reader.calls, "current definitions cannot resolve immutable handoff meaning")
					f.pipe = New(f.pipe.deps, f.pipe.config)
				}
				if strings.HasSuffix(mode, "failure") {
					require.Positive(t, faultCalls, "the intended infrastructure boundary was reached")
				}
				// Already accepted healthy input still closes under the injected fault.
				// A global runtime-proof read failure cannot promise new native intake.
				require.ErrorIs(t, f.pipe.Reconcile(t.Context()), wantErr)
				require.NotNil(t, f.pipe.boundaries[boundaryKey(healthy.GroupID, healthy.AppliedRevision)], "%+v", f.pipe.Status())
				require.NoError(t, f.pipe.boundaries[boundaryKey(healthy.GroupID, healthy.AppliedRevision)].boundary.Tick(t.Context(), at(80)))
				f.deliver()
				var value string
				require.NoError(t, f.plainTarget().QueryRow(`SELECT reading FROM healthy_proof_readings`).Scan(&value))
				require.Equal(t, "7", value)
				var targetRows int
				require.NoError(t, f.plainTarget().QueryRow(`SELECT COUNT(*) FROM `+f.target.Group.Destination.TableName+` WHERE group_id=?`, f.target.GroupID).Scan(&targetRows))
				require.Zero(t, targetRows, "blocked is neither SQL committed nor accepted")
				require.Equal(t, before, conflictPairState(t, f))
				require.Zero(t, reader.calls)
			})
		}
	}
}
