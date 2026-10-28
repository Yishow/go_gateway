package grouppipeline

import (
	"encoding/json"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Canonical lifecycle commits must fence new durable acceptance even while the
// reconciler still holds the old active boundary. Existing accepted inputs are
// not canceled, and a different recording remains independent.
func TestLifecycleCommitRejectsCachedBoundaryIntakeBeforeReconcile(t *testing.T) {
	for _, action := range []string{"stop", "remove"} {
		t.Run(action, func(t *testing.T) {
			f := newFixture(t)
			f.store = f.store.WithCanonicalIntake() // same explicit canonical policy as production wiring
			db := f.store.DB()
			ws := workspace.NewService(workspace.NewSQLRepository(db))
			record, err := ws.GetOrCreate(t.Context())
			require.NoError(t, err)
			group := groupAt("R1")
			group.ID = "G1"
			group.WorkspaceID = record.ID
			group.Members = group.Members[:1]
			healthy := groupAt("healthy-R1")
			healthy.ID = "healthy-G1"
			healthy.WorkspaceID = record.ID
			healthy.Members = healthy.Members[:1]
			healthy.Members[0].DeviceID = "healthy-device"
			healthy.Members[0].PointID = "healthy-point"
			healthy.Members[0].TagID = "healthy-tag"
			for _, g := range []*workspace.WriteGroup{group, healthy} {
				persistLifecycleIntakeGroup(t, f, g)
			}
			_, err = db.Exec(`INSERT INTO database_connectors(id,name,kind,identity_revision,connection_config,status) VALUES('connector-1','fixture','sqlite','crev-1','{}','ready')`)
			require.NoError(t, err)
			service := workspace.NewWriteGroupService(ws, workspace.NewSQLWriteGroupRepository(db)).WithBacklogOwnershipGuard(BacklogGuard{}).WithClock(func() time.Time { return f.clock })
			pipe := New(Dependencies{Groups: service, Tags: f.tags, Destinations: f.dest, Inspector: f.inspect, Store: f.store}, Config{Now: func() time.Time { return f.clock }})
			require.NoError(t, pipe.Reconcile(t.Context()))
			require.Len(t, pipe.boundaries, 2, "both cached boundaries were built from canonical immutable revisions: %+v", pipe.Status())
			f.clock = at(2)
			require.NoError(t, pipe.AcceptSample(t.Context(), envelope(group, "already-closed", at(1), 10.5)))
			f.clock = at(12)
			pipe.TickAll(t.Context())
			var queued int
			require.NoError(t, db.QueryRow(`SELECT COUNT(*) FROM wg_delivery_outbox WHERE group_id='G1' AND state='pending'`).Scan(&queued))
			require.Equal(t, 1, queued)
			require.NoError(t, pipe.AcceptSample(t.Context(), envelope(group, "already-accepted", at(11), 11.5)))
			f.clock = at(18)
			require.NoError(t, pipe.AcceptSample(t.Context(), envelope(group, "preaccepted-future", at(21), 21.5)), "three seconds ahead is within the existing five-second MaxFutureSkew")
			_, err = db.Exec(`INSERT INTO wg_delivery_receipts(effect_key,payload_digest,committed_at) VALUES('prior-receipt','prior-digest','2026-01-01T00:00:00Z')`)
			require.NoError(t, err)
			current, err := service.Get(t.Context(), group.ID)
			require.NoError(t, err)
			mutation := workspace.WriteGroupMutation{WorkspaceID: record.ID, ExpectedWorkspaceRevision: current.WorkspaceRevision, ExpectedGroupRevision: current.Group.Revision, ExpectedConnectorRevision: current.Group.Destination.ConnectorRevision}
			if action == "stop" {
				_, err = service.Disable(t.Context(), group.ID, mutation)
			} else {
				_, err = service.Delete(t.Context(), group.ID, mutation)
			}
			require.NoError(t, err)
			eligibility, err := service.IntakeEligibility(t.Context(), group.ID)
			require.NoError(t, err)
			require.False(t, eligibility.Accepting)
			cached := pipe.boundaries[boundaryKey(group.ID, group.AppliedRevision)]
			require.NotNil(t, cached)
			require.False(t, cached.drainOnly, "no reconciliation has refreshed the old active cache")
			before := backlogOwnershipState(t, db)
			require.NoError(t, pipe.AcceptSample(t.Context(), envelope(group, "already-accepted", at(11), 11.5)), "an exact retry recognizes durable acceptance from before stop; it is not new intake")
			require.Equal(t, before, backlogOwnershipState(t, db))
			f.clock = at(19)
			err = pipe.AcceptSample(t.Context(), envelope(group, "after-lifecycle-commit", at(18), 12.5))
			assert.Error(t, err, "a new sample must not be ACKed after the canonical stop/remove committed")
			var refusal *groupdelivery.OwnershipError
			assert.ErrorAs(t, err, &refusal, "reuse the safe typed not-accepted ownership contract")
			after := backlogOwnershipState(t, db)
			assert.Equal(t, before, after, "refused sample must not change accepted samples, outbox, versions or receipts")
			var appended int
			require.NoError(t, db.QueryRow(`SELECT COUNT(*) FROM wg_delivery_samples WHERE sample_id='after-lifecycle-commit'`).Scan(&appended))
			assert.Zero(t, appended, "the cached boundary must not commit the post-stop sample")
			require.NoError(t, pipe.AcceptSample(t.Context(), envelope(healthy, "healthy-after-stop", at(18), 7.5)), "another active group remains usable without a global intake stop")
			require.NoError(t, db.QueryRow(`SELECT COUNT(*) FROM wg_delivery_samples WHERE group_id='healthy-G1' AND sample_id='healthy-after-stop'`).Scan(&appended))
			require.Equal(t, 1, appended)
			for _, table := range []string{"wg_delivery_outbox", "wg_delivery_receipts", "write_group_versions", "wg_runtime_versions"} {
				require.Equal(t, before[table], backlogOwnershipState(t, db)[table], "independent intake also leaves frozen history unchanged")
			}
			t.Run("preaccepted-future-drain", func(t *testing.T) {
				var payload, digest string
				require.NoError(t, db.QueryRow(`SELECT payload,payload_digest FROM wg_delivery_samples WHERE group_id='G1' AND sample_id='preaccepted-future' AND consumed=0`).Scan(&payload, &digest))
				// The existing stop boundary closes at 20, before the accepted
				// observation at 21. Historical-journal recovery must then drain it.
				require.NoError(t, pipe.Reconcile(t.Context()))
				require.Equal(t, at(20), cached.boundary.ReplayState().Until)
				f.clock = at(31)
				pipe.TickAll(t.Context())
				require.NoError(t, pipe.Reconcile(t.Context()))
				pipe.TickAll(t.Context())
				var afterPayload, afterDigest string
				var consumed, effects int
				require.NoError(t, db.QueryRow(`SELECT payload,payload_digest,consumed FROM wg_delivery_samples WHERE group_id='G1' AND sample_id='preaccepted-future'`).Scan(&afterPayload, &afterDigest, &consumed))
				require.Equal(t, payload, afterPayload)
				require.Equal(t, digest, afterDigest)
				require.Equal(t, 1, consumed, "preaccepted future observation must finish closure, not remain stranded")
				require.NoError(t, db.QueryRow(`SELECT COUNT(*) FROM wg_delivery_outbox WHERE group_id='G1' AND group_revision='R1' AND bucket_start=?`, "2026-01-01T00:00:20.000000000Z").Scan(&effects))
				if effects != 1 {
					var outcomes, queued string
					require.NoError(t, db.QueryRow(`SELECT COALESCE(json_group_array(json_object('bucket',bucket_start,'kind',kind,'reason',reason,'members',members)),'[]') FROM wg_delivery_buckets WHERE group_id='G1'`).Scan(&outcomes))
					require.NoError(t, db.QueryRow(`SELECT COALESCE(group_concat(bucket_start),'') FROM wg_delivery_outbox WHERE group_id='G1'`).Scan(&queued))
					t.Logf("future closure outcomes=%s queued=%s", outcomes, queued)
				}
				require.Equal(t, 1, effects, "the original revision keeps one frozen effect for the accepted future bucket")
				var connector, connectorRevision, database, schema, table string
				require.NoError(t, db.QueryRow(`SELECT connector_id,connector_revision,database_name,table_schema,table_name FROM wg_delivery_outbox WHERE group_id='G1' AND group_revision='R1' AND bucket_start='2026-01-01T00:00:20.000000000Z'`).Scan(&connector, &connectorRevision, &database, &schema, &table))
				require.Equal(t, []string{"connector-1", "crev-1", "", "main", "readings"}, []string{connector, connectorRevision, database, schema, table}, "recovery keeps the original frozen destination")
				require.Equal(t, before["wg_delivery_receipts"], backlogOwnershipState(t, db)["wg_delivery_receipts"])
			})
		})
	}
}

func persistLifecycleIntakeGroup(t *testing.T, f *fixture, group *workspace.WriteGroup) {
	t.Helper()
	tx, err := f.store.DB().BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	persistOwnershipVersion(t, tx, snapshotOf(group, at(0)))
	rowPolicy, err := json.Marshal(group.RowPolicy)
	require.NoError(t, err)
	writePolicy, err := json.Marshal(group.WritePolicy)
	require.NoError(t, err)
	_, err = tx.Exec(`UPDATE write_groups SET row_policy=?,write_policy=? WHERE id=?`, string(rowPolicy), string(writePolicy), group.ID)
	require.NoError(t, err)
	for i, m := range group.Members {
		_, err = tx.Exec(`INSERT INTO write_group_members(group_id,member_index,device_id,point_id,tag_id,source_revision,mapping_revision,target_column,required) VALUES(?,?,?,?,?,?,?,?,?)`, group.ID, i, m.DeviceID, m.PointID, m.TagID, m.SourceRevision, m.MappingRevision, m.TargetColumn, m.Required)
		require.NoError(t, err)
	}
	require.NoError(t, tx.Commit())
	// Freeze verified runtime storage/type facts as an ordinary applied revision;
	// no live source metadata reconstruction is needed for this lifecycle test.
	snap := snapshotOf(group, at(0))
	snap.RuntimeLayout = &workspace.WriteGroupRuntimeLayout{Dialect: dbtarget.SQLDialectSQLite, Columns: f.inspect.inspection.Columns, TagTypes: map[string]schema.DataType{group.Members[0].TagID: schema.DataTypeFloat64}}
	payload, err := json.Marshal(snap)
	require.NoError(t, err)
	require.NoError(t, f.store.SaveRuntimeVersion(t.Context(), groupKey(snap), payload))
}
