package runtime

import (
	"database/sql"
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// Construct the same real SQLite authority as the existing replay fixture,
// with the two already-supported target layouts. No accepted proof is edited.
func newReductionHandoffFixture(t *testing.T, wideTarget bool) *runtimeHandoffFixture {
	t.Helper()
	source, target := wideCounterFixture(t), newCounterFixture(t)
	if wideTarget {
		target = wideCounterFixture(t)
		target.config.Group.ID = "group-next-wide"
	}
	source.members = source.members[:1]
	source.config.Group.Members = source.members
	target.members = append([]workspace.WriteGroupMember(nil), source.members...)
	if !wideTarget {
		target.members[0].SourceValueType, target.members[0].ResultValueType = "", ""
	}
	target.config.Group.Members = target.members
	target.config.TagTypes, target.config.Counters = source.config.TagTypes, source.config.Counters
	source.clock, target.clock = at(11*time.Second), at(11*time.Second)
	path := filepath.Join(t.TempDir(), "reduction-handoff.db")
	store := groupdelivery.NewStore(openDurableDB(t, path)).WithClock(func() time.Time { return source.clock })
	boundary := at(12500 * time.Millisecond)
	f := &runtimeHandoffFixture{source: source, target: target, store: store, boundary: boundary, databasePath: path}
	f.sourceKey = freezeHandoffRuntime(t, store, source, boundaryStart)
	f.targetKey = freezeHandoffRuntime(t, store, target, boundary)
	target.config.StartsAt = boundary
	f.old, f.emptyTarget = source.build(t), target.build(t)
	for _, input := range []struct {
		id    string
		at    time.Time
		value uint64
	}{{"after", at(13 * time.Second), 200}, {"exact", boundary, 150}, {"before", boundary.Add(-time.Millisecond), 100}} {
		require.NoError(t, f.old.AcceptSample(t.Context(), source.envelope(0, input.id, input.at, input.value)))
	}
	tx, err := store.DB().BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	runtimeHandoffVersion(t, tx, source, boundaryStart)
	runtimeHandoffVersion(t, tx, target, boundary)
	m := source.members[0]
	detail, err := json.Marshal(map[string]any{
		"ownership_fence_revision": "runtime-reduction",
		"ended":                    []workspace.CutoverEnd{{GroupID: f.sourceKey.GroupID, Revision: f.sourceKey.GroupRevision}},
		"started":                  []workspace.CutoverStart{{GroupID: f.targetKey.GroupID, Revision: f.targetKey.GroupRevision}},
		"moved":                    []workspace.CutoverMemberRef{{DeviceID: m.DeviceID, PointID: m.PointID, TagID: m.TagID}},
	})
	require.NoError(t, err)
	_, err = tx.Exec(`INSERT INTO write_group_cutovers(id,workspace_id,kind,source_group_id,boundary,review_digest,detail,intent) VALUES('runtime-reduction',?,'forward',?,?,'reduction-review',?,'{}')`, f.sourceKey.WorkspaceID, f.sourceKey.GroupID, boundary.Format(time.RFC3339Nano), string(detail))
	require.NoError(t, err)
	_, err = tx.Exec(`INSERT INTO write_group_revision_ends(workspace_id,group_id,group_revision,ends_at,cutover_id) VALUES(?,?,?,?,'runtime-reduction')`, f.sourceKey.WorkspaceID, f.sourceKey.GroupID, f.sourceKey.GroupRevision, boundary.Format(time.RFC3339Nano))
	require.NoError(t, err)
	for direction, owner := range map[string]*longFixture{"ended": source, "started": target} {
		g := owner.config.Group
		_, err = tx.Exec(`INSERT INTO write_group_ownership_fences(fence_revision,workspace_id,direction,group_id,group_revision,device_id,point_id,tag_id,profile_kind,connector_id,connector_revision,database_name,table_schema,table_name) VALUES('runtime-reduction',?,?,?,?,?,?,?,?,?,?,?,?,?)`, g.WorkspaceID, direction, g.ID, g.AppliedRevision, m.DeviceID, m.PointID, m.TagID, string(g.RecordingProfile.Kind), g.Destination.ConnectorID, g.Destination.ConnectorRevision, g.Destination.Database, g.Destination.TableSchema, g.Destination.TableName)
		require.NoError(t, err)
	}
	require.NoError(t, tx.Commit())
	f.old.SetEndsAt(boundary)
	source.config.EndsAt, source.config.Until = boundary, at(15*time.Second)
	return f
}

type reductionFrozenRow struct{ payload, digest, scope, connector, table string }

func reductionOutbox(t *testing.T, db *sql.DB, group string) reductionFrozenRow {
	t.Helper()
	var row reductionFrozenRow
	require.NoError(t, db.QueryRow(`SELECT payload,payload_digest,destination_scope,connector_id,table_name FROM wg_delivery_outbox WHERE group_id=? ORDER BY bucket_start LIMIT 1`, group).Scan(&row.payload, &row.digest, &row.scope, &row.connector, &row.table))
	return row
}

func reductionOriginal(t *testing.T, db *sql.DB, group, id string) [3]string {
	t.Helper()
	var result [3]string
	require.NoError(t, db.QueryRow(`SELECT payload,payload_digest,created_at FROM wg_delivery_samples WHERE group_id=? AND json_extract(payload,'$.sample_id')=?`, group, id).Scan(&result[0], &result[1], &result[2]))
	return result
}

func reductionReadback(t *testing.T, f *longFixture, frozen reductionFrozenRow, expectedValue, side string, boundary time.Time) {
	t.Helper()
	row, err := groupdelivery.DecodeRowPayload([]byte(frozen.payload))
	require.NoError(t, err)
	_, err = dbtarget.InsertGroupRow(t.Context(), f.db, dbtarget.GroupInsertRequest{
		Kind: schema.DatabaseConnectorKindSQLite, SchemaName: "main", TableName: frozen.table, Row: row,
		PayloadDigest: frozen.digest, Strategy: dbtarget.GroupEffectReceipt, CommittedAt: at(700 * time.Second).Format(time.RFC3339Nano),
	})
	require.NoError(t, err)
	var value, provenance string
	if f.config.Group.LayoutVersion == workspace.WriteGroupLayoutManagedLongV1 {
		var status string
		var delta, from, to sql.NullString
		require.NoError(t, f.db.QueryRow(`SELECT value_decimal,status,delta_value,delta_from,delta_to,provenance FROM `+longFixtureTable).Scan(&value, &status, &delta, &from, &to, &provenance))
		require.Equal(t, "baseline", status)
		require.False(t, delta.Valid || from.Valid || to.Valid, "baseline delta and both bounds must be SQL NULL")
	} else {
		require.True(t, row.Partial, "wide rows carry their partial flag; long rows preserve baseline status and window evidence")
		require.NoError(t, f.db.QueryRow(`SELECT v_gauge,provenance FROM `+wideFixtureTable).Scan(&value, &provenance))
		var entries []struct {
			Counter struct {
				Status string  `json:"status"`
				Delta  *string `json:"delta"`
			} `json:"counter"`
		}
		require.NoError(t, json.Unmarshal([]byte(provenance), &entries))
		require.Len(t, entries, 1)
		require.Equal(t, "baseline", entries[0].Counter.Status)
		require.Nil(t, entries[0].Counter.Delta, "wide baseline carries no fabricated zero delta")
	}
	require.Equal(t, expectedValue, value)
	require.Contains(t, provenance, `"truncated":"`+side+`"`)
	require.Contains(t, provenance, boundary.Format(longTimeLayout), "partial evidence keeps the exact actual boundary")
}

func TestGroupBoundaryCutoverReductionVerifiedCumulativeMatrix(t *testing.T) {
	for _, wideTarget := range []bool{false, true} {
		name := "managed-long"
		if wideTarget {
			name = "managed-wide"
		}
		t.Run(name, func(t *testing.T) {
			f := newReductionHandoffFixture(t, wideTarget)
			original, err := f.store.Restore(t.Context(), f.sourceKey)
			require.NoError(t, err)
			require.Len(t, original.Samples, 3)
			bytes := map[string][3]string{}
			for _, sample := range original.Samples {
				bytes[sample.SampleID] = reductionOriginal(t, f.store.DB(), f.sourceKey.GroupID, sample.SampleID)
			}
			batch, err := f.store.PrepareCutoverHandoffs(t.Context(), f.sourceKey, 0, 128)
			require.NoError(t, err)
			require.Equal(t, 2, batch.Accepted, "exactly-at and after transfer; before remains old-owned")
			restored, err := f.store.Restore(t.Context(), f.targetKey)
			require.NoError(t, err)
			require.Len(t, restored.Samples, 2)
			for _, sample := range restored.Samples {
				require.Equal(t, bytes[sample.SampleID], reductionOriginal(t, f.store.DB(), f.targetKey.GroupID, sample.SampleID), "payload/digest/original accepted_at are unchanged")
				require.Equal(t, f.sourceKey.GroupID, sample.CounterProof.Member.GroupID)
				require.Equal(t, snapshot.AcceptanceOrigin{WorkspaceID: f.sourceKey.WorkspaceID, GroupID: f.sourceKey.GroupID, GroupRevision: f.sourceKey.GroupRevision}, sample.AcceptanceOrigin)
			}
			require.Error(t, f.old.Tick(t.Context(), at(700*time.Second)))
			require.Error(t, f.emptyTarget.Tick(t.Context(), at(700*time.Second)))
			require.NoError(t, f.store.DB().Close())
			f.store = groupdelivery.NewStore(openDurableDB(t, f.databasePath))
			f.source.config.Ledger = &frozenFixtureLedger{Ledger: f.store.Ledger(f.sourceKey), store: f.store}
			f.target.config.Ledger = &frozenFixtureLedger{Ledger: f.store.Ledger(f.targetKey), store: f.store}
			// Only frozen descriptors and verified Restore may supply counter meaning.
			f.source.config.Counters, f.target.config.Counters = nil, nil
			old, target := f.source.build(t), f.target.build(t)
			require.NoError(t, old.Tick(t.Context(), at(15*time.Second)))
			targetEnd := at(time.Duration(f.target.config.Group.RowPolicy.IntervalSeconds) * time.Second)
			require.NoError(t, target.Tick(t.Context(), targetEnd))
			for _, tc := range []struct {
				key  groupdelivery.GroupKey
				ids  []string
				last string
				end  time.Time
			}{{f.sourceKey, []string{"before"}, "before", f.boundary}, {f.targetKey, []string{"exact", "after"}, "after", targetEnd}} {
				rows, err := f.store.DB().Query(`SELECT json_extract(payload,'$.sample_id') FROM wg_delivery_samples WHERE group_id=? AND consumed=1`, tc.key.GroupID)
				require.NoError(t, err)
				var ids []string
				for rows.Next() {
					var id string
					require.NoError(t, rows.Scan(&id))
					ids = append(ids, id)
				}
				require.NoError(t, rows.Err())
				require.NoError(t, rows.Close())
				require.ElementsMatch(t, tc.ids, ids, "consuming-owner sets are disjoint and their union is the accepted set")
				state, err := f.store.Restore(t.Context(), tc.key)
				require.NoError(t, err)
				require.Len(t, state.Baselines, 1)
				checkpoint, err := aggregation.DecodeCounterCheckpoint(state.Baselines[0].Payload)
				require.NoError(t, err)
				require.Equal(t, tc.last, checkpoint.Reading.SampleID)
				require.Equal(t, tc.end, checkpoint.WindowEnd)
			}
			oldFrozen, newFrozen := reductionOutbox(t, f.store.DB(), f.sourceKey.GroupID), reductionOutbox(t, f.store.DB(), f.targetKey.GroupID)
			reductionReadback(t, f.source, oldFrozen, "100", "end", f.boundary)
			reductionReadback(t, f.target, newFrozen, "200", "start", f.boundary)
			// Restart/retry cannot reopen finalized outputs or change frozen bytes/destination.
			require.NoError(t, f.source.build(t).Tick(t.Context(), at(15*time.Second)))
			require.NoError(t, f.target.build(t).Tick(t.Context(), targetEnd))
			require.Equal(t, oldFrozen, reductionOutbox(t, f.store.DB(), f.sourceKey.GroupID))
			require.Equal(t, newFrozen, reductionOutbox(t, f.store.DB(), f.targetKey.GroupID))
			for id, originalBytes := range bytes {
				require.Equal(t, originalBytes, reductionOriginal(t, f.store.DB(), f.sourceKey.GroupID, id))
			}
		})
	}
}
