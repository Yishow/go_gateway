package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"
	"go-gateway/internal/testutil/lostcommit"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

// ensureTimescaleGroup saves a 15 second instantaneous long group for the "ts"
// circuit (the plain twin of the same profile uses "main") with the explicit
// timescale mode through the public route, and returns the response whatever
// its status.
func ensureTimescaleGroup(t *testing.T, f managedGroupFixture) (status int, body map[string]any) {
	t.Helper()
	request := profileRequest(t, f.writeGroupRouterFixture, "instantaneous", 15, "second", "ts", "")
	group := request["group"].(map[string]any)
	group["layout_version"] = "managed_long_v1"
	group["storage_mode"] = "timescale"
	group["destination"].(map[string]any)["table_schema"] = f.namespace
	response := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, request)
	return response.Code, decodeJSONBody(t, response)
}

func timescaleInstalled(t *testing.T, f managedGroupFixture) bool {
	t.Helper()
	var n int
	require.NoError(t, f.target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM pg_extension WHERE extname = 'timescaledb'`).Scan(&n))
	return n > 0
}

func errorOf(t *testing.T, body map[string]any) map[string]any {
	t.Helper()
	failure, ok := body["error"].(map[string]any)
	require.True(t, ok, "%v", body)
	return failure
}

func catalogObjects(t *testing.T, f managedGroupFixture) (relations, extensions int) {
	t.Helper()
	require.NoError(t, f.target.QueryRowContext(t.Context(),
		`SELECT (SELECT COUNT(*) FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace WHERE n.nspname = $1), (SELECT COUNT(*) FROM pg_extension)`,
		f.namespace).Scan(&relations, &extensions))
	return relations, extensions
}

// TestTimescaleSelectionRejectedOnSQLite: SQLite has no hypertable. The explicit
// selection is refused with the stable code and the plain-storage repair, the
// group is not saved, and nothing is rewritten to plain behind the user's back.
func TestTimescaleSelectionRejectedOnSQLite(t *testing.T) {
	f := newManagedDestinationFixture(t, "sqlite")
	status, body := ensureTimescaleGroup(t, f)
	require.Equal(t, http.StatusUnprocessableEntity, status, "%v", body)
	failure := errorOf(t, body)
	require.Equal(t, "timescale_unavailable", failure["code"])
	require.Equal(t, "select_plain_storage", failure["action"])
	require.Contains(t, failure["message"], "dialect-unsupported")
	listed := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath, nil)
	require.NotContains(t, listed.Body.String(), `"storage_mode":"timescale"`, "no group was saved with the rejected mode")
}

// TestTimescaleUnavailableBlocksWithoutChange: on a destination without the
// extension the selected timescale mode is blocked at preview with the stable
// code, its reason and a repair action; the destination is untouched, the
// extension is not installed, and the saved selection stays timescale.
func TestTimescaleUnavailableBlocksWithoutChange(t *testing.T) {
	f := newManagedDestinationFixture(t, "postgres")
	if timescaleInstalled(t, f) {
		t.Skip("this destination database has the TimescaleDB extension installed; the extension-absent case is NOT RUN here")
	}
	status, group := ensureTimescaleGroup(t, f)
	require.Equal(t, http.StatusOK, status, "%v", group)
	_, saved := groupSaveData(t, group)
	require.Equal(t, "timescale", saved["storage_mode"])
	groupID := saved["id"].(string)

	relations, extensions := catalogObjects(t, f)
	response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+groupID+"/schema-preview", scopeBody(t, f, groupID))
	require.Equal(t, http.StatusUnprocessableEntity, response.Code, response.Body.String())
	failure := errorOf(t, decodeJSONBody(t, response))
	require.Equal(t, "timescale_unavailable", failure["code"])
	require.Equal(t, "repair_timescale_capability", failure["action"])
	require.Contains(t, failure["message"], "extension-not-installed")
	require.NotContains(t, response.Body.String(), "gwtest", "no connection detail is exposed")

	afterRelations, afterExtensions := catalogObjects(t, f)
	require.Equal(t, [2]int{relations, extensions}, [2]int{afterRelations, afterExtensions}, "nothing was created and no extension was installed")
	after := getGroup(t, f, groupID)
	require.Equal(t, "timescale", after["storage_mode"], "the selection is never silently rewritten to plain")
	ready, err := f.groups.Readiness(t.Context(), groupID)
	require.NoError(t, err)
	require.False(t, ready.Ready)
}

// TestTimescaleHypertableThroughTheProductionPath runs the whole route: an
// explicit timescale group previews its seven-day hypertable, creates nothing
// until confirmed, then confirmation creates the owned empty hypertable once;
// the catalog proves it; Apply freezes that evidence; and the production pipeline
// writes a record into it that reads back.
func TestTimescaleHypertableThroughTheProductionPath(t *testing.T) {
	f := newManagedDestinationFixture(t, "postgres")
	if !timescaleInstalled(t, f) {
		t.Skip("this destination database has no TimescaleDB extension; the capable case is NOT RUN here")
	}
	seedProfileSource(t, f.writeGroupRouterFixture, "voltage")
	_, extensionsBefore := catalogObjects(t, f)

	plainGroup := ensureLong(t, f, "instantaneous", "second", 15, "voltage")
	status, created := ensureTimescaleGroup(t, f)
	require.Equal(t, http.StatusOK, status, "%v", created)
	_, group := groupSaveData(t, created)
	groupID := group["id"].(string)
	table := destinationOf(group)["table_name"].(string)
	require.Equal(t, "timescale", group["storage_mode"])
	require.Regexp(t, longTablePattern, table)
	require.NotEqual(t, destinationOf(plainGroup)["table_name"], table, "the mode is part of the storage scope: a timescale table is never the plain table")

	token := previewLong(t, f, groupID)
	require.Len(t, token.Statements, 4, strings.Join(token.Statements, "\n"))
	require.Contains(t, token.Statements[1], "create_hypertable")
	require.Contains(t, token.Statements[1], "INTERVAL '7 days'")
	require.Contains(t, token.Statements[1], "migrate_data => false")
	require.Contains(t, token.Statements[1], "if_not_exists => false")
	for _, statement := range token.Statements {
		require.NotRegexp(t, `(?i)create\s+extension|retention|compress|continuous|add_job|partition\s+by`, statement)
	}
	require.Equal(t, "7 days", token.Tables[0].Partition.ChunkInterval, "the interval is previewed")
	require.Equal(t, "bucket_start", token.Tables[0].Partition.Column)
	require.Equal(t, "public", token.GroupLayout.Partition.ExtensionSchema, "the probed schema is frozen into the confirmed layout")
	require.Empty(t, columnsOf(t, f.target, "postgres", f.namespace, table), "a preview creates nothing")

	operation := confirmLong(t, f, groupID, token)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, operation.Status)
	require.Equal(t, len(token.Statements), operation.ExecutedStatements)
	require.NotEmpty(t, columnsOf(t, f.target, "postgres", f.namespace, table))

	var dimensions, chunks, jobs int
	var timeColumn string
	var seconds float64
	var compression bool
	require.NoError(t, f.target.QueryRowContext(t.Context(), `SELECT num_dimensions, num_chunks, compression_enabled FROM timescaledb_information.hypertables WHERE hypertable_schema = $1 AND hypertable_name = $2`, f.namespace, table).Scan(&dimensions, &chunks, &compression))
	require.NoError(t, f.target.QueryRowContext(t.Context(), `SELECT column_name, EXTRACT(EPOCH FROM time_interval)::float8 FROM timescaledb_information.dimensions WHERE hypertable_schema = $1 AND hypertable_name = $2`, f.namespace, table).Scan(&timeColumn, &seconds))
	require.NoError(t, f.target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM timescaledb_information.jobs WHERE hypertable_schema = $1 AND hypertable_name = $2`, f.namespace, table).Scan(&jobs))
	require.Equal(t, [3]any{1, "bucket_start", false}, [3]any{dimensions, timeColumn, compression})
	require.InDelta(t, 7*24*3600, seconds, 0.5)
	require.Zero(t, chunks, "created empty")
	require.Zero(t, jobs)
	var keyColumns string
	require.NoError(t, f.target.QueryRowContext(t.Context(), `
		SELECT string_agg(a.attname, ',' ORDER BY k.ord) FROM pg_index i JOIN pg_class c ON c.oid = i.indrelid JOIN pg_namespace n ON n.oid = c.relnamespace
		JOIN LATERAL unnest(i.indkey::int2[]) WITH ORDINALITY k(attnum, ord) ON true JOIN pg_attribute a ON a.attrelid = c.oid AND a.attnum = k.attnum
		WHERE n.nspname = $1 AND c.relname = $2 AND i.indisprimary`, f.namespace, table).Scan(&keyColumns))
	require.Equal(t, "bucket_start,record_id", keyColumns)

	// A second confirmation of the same operation neither runs DDL nor changes the catalog.
	relations, _ := catalogObjects(t, f)
	again := confirmLong(t, f, groupID, token)
	require.Equal(t, operation.OperationID, again.OperationID)
	afterRelations, _ := catalogObjects(t, f)
	require.Equal(t, relations, afterRelations)

	// A fresh preview of the confirmed table is a no-change preview: an owned
	// hypertable is never converted again, yet it still has to be confirmed for
	// this group revision, which records the operation and verified digest.
	next := previewLong(t, f, groupID)
	require.Empty(t, next.Statements, "an owned hypertable that already matches needs no DDL")
	require.Equal(t, recordingplan.NoChangeSchemaCompatible, next.NoChangeReason)
	noop := confirmLong(t, f, groupID, next)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, noop.Status)
	require.Zero(t, noop.ExecutedStatements)
	require.Equal(t, operation.VerifiedDigest, noop.VerifiedDigest, "the same catalog evidence gives the same digest")
	require.Equal(t, noop.OperationID, destinationOf(getGroup(t, f, groupID))["schema_revision"], "the group revision is bound to the no-change confirmation")

	ready, err := f.groups.Readiness(t.Context(), groupID)
	require.NoError(t, err)
	require.True(t, ready.Ready, "%+v", ready.Issues)
	applied := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+groupID+"/apply", scopeBody(t, f, groupID))
	require.Equal(t, http.StatusOK, applied.Code, applied.Body.String())
	snapshot, err := f.groups.ResolveAppliedAt(t.Context(), groupID, time.Now().Add(time.Minute))
	require.NoError(t, err)
	frozen, err := groupdelivery.NewStore(f.db).RuntimeVersion(t.Context(), groupdelivery.GroupKey{WorkspaceID: snapshot.WorkspaceID, GroupID: groupID, GroupRevision: snapshot.AppliedRevision})
	require.NoError(t, err)
	require.Contains(t, string(frozen), `"hypertable":{"time_column":"bucket_start","chunk_interval_micros":604800000000}`, "Apply froze the verified catalog evidence")

	writeOneSampleThroughPipeline(t, f, groupID, table)

	_, extensionsAfter := catalogObjects(t, f)
	require.Equal(t, extensionsBefore, extensionsAfter, "no extension was installed by any step")
}

// TestTimescaleOperationLostResponse loses the commit response of the
// hypertable batch at every prefix. Recovery only reads the existing operation's
// catalog evidence: a partial batch (a plain table without its hypertable, or a
// hypertable without its index and receipts) is never completed, converted or
// claimed as confirmed, and only the whole batch recovers as success.
func TestTimescaleOperationLostResponse(t *testing.T) {
	for _, committed := range []int{0, 1, 2, 4} {
		t.Run(string(rune('0'+committed)), func(t *testing.T) {
			f := newManagedDestinationFixture(t, "postgres")
			if !timescaleInstalled(t, f) {
				t.Skip("this destination database has no TimescaleDB extension; the case is NOT RUN here")
			}
			_, created := ensureTimescaleGroup(t, f)
			_, group := groupSaveData(t, created)
			groupID, table := group["id"].(string), destinationOf(group)["table_name"].(string)
			token := previewLong(t, f, groupID)
			require.Len(t, token.Statements, 4)

			ledger := recordingplan.NewService(recordingplan.NewSQLRepository(f.db))
			_, outcome, err := ledger.ClaimSchemaApply(t.Context(), &token)
			require.NoError(t, err)
			require.Equal(t, recordingplan.ClaimAcquired, outcome)

			lose := &atomic.Bool{}
			lose.Store(true)
			target := lostcommit.Wrap(stdlib.GetDefaultDriver(), os.Getenv("POSTGRES_DSN"), lose)
			t.Cleanup(func() { require.NoError(t, target.Close()) })
			tx, err := target.BeginTx(t.Context(), nil)
			require.NoError(t, err)
			_, err = tx.ExecContext(t.Context(), `SET LOCAL search_path TO "`+f.namespace+`"`)
			require.NoError(t, err)
			for _, statement := range token.Statements[:committed] {
				_, err := tx.ExecContext(t.Context(), statement)
				require.NoError(t, err, statement)
			}
			require.ErrorIs(t, tx.Commit(), lostcommit.ErrResponseLost)

			old := time.Now().UTC().Add(-recordingplan.SchemaOperationLease - time.Minute)
			_, err = f.db.ExecContext(t.Context(), `UPDATE managed_schema_operations SET updated_at=? WHERE operation_id=?`, old, token.OperationID)
			require.NoError(t, err)
			_, err = f.db.ExecContext(t.Context(), `UPDATE managed_schema_preview_tokens SET expires_at=? WHERE token=?`, old, token.Token)
			require.NoError(t, err)
			var seq int
			var dbName, file string
			require.NoError(t, f.db.QueryRowContext(t.Context(), `PRAGMA database_list`).Scan(&seq, &dbName, &file))
			require.NoError(t, f.db.Close())
			f.db, err = sql.Open("sqlite", file)
			require.NoError(t, err)
			f.db.SetMaxOpenConns(1)
			t.Cleanup(func() { require.NoError(t, f.db.Close()) })
			f.workspace = workspace.NewService(workspace.NewSQLRepository(f.db))
			targets := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(f.db))
			f.groups = workspace.NewWriteGroupService(f.workspace, workspace.NewSQLWriteGroupRepository(f.db)).WithTableInspector(dbtarget.NewReadOnlyTableInspector(targets))
			f.groups.WithManagedTableInspector(f.groups.ManagedTableInspector(targets))
			f.router = NewRouter(&DatalinkServices{Workspace: f.workspace, WriteGroups: f.groups, DBTarget: targets, RecordingPlan: recordingplan.NewService(recordingplan.NewSQLRepository(f.db))})

			want := recordingplan.SchemaOperationUnknown
			if committed == 4 {
				want = recordingplan.SchemaOperationSucceeded
			}
			for range 2 {
				request := scopeBody(t, f, groupID)
				request["token"], request["operation_id"] = token.Token, token.OperationID
				response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+groupID+"/schema-apply", request)
				require.Equal(t, http.StatusOK, response.Code, response.Body.String())
				var body struct {
					Data recordingplan.SchemaOperation `json:"data"`
				}
				require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
				require.Equal(t, token.OperationID, body.Data.OperationID)
				require.Equal(t, want, body.Data.Status)
			}

			var objects int
			require.NoError(t, f.target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
				WHERE n.nspname = $1 AND c.relname IN ($2, $3, 'gw_effect_receipts')`, f.namespace, table, "idx_"+table+"_src_time").Scan(&objects))
			wantObjects := map[int]int{0: 0, 1: 1, 2: 1, 4: 3}[committed]
			require.Equal(t, wantObjects, objects, "recovery only inspects evidence; it never runs the missing DDL")
			var hypertables int
			require.NoError(t, f.target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM timescaledb_information.hypertables WHERE hypertable_schema = $1 AND hypertable_name = $2`, f.namespace, table).Scan(&hypertables))
			require.Equal(t, map[int]int{0: 0, 1: 0, 2: 1, 4: 1}[committed], hypertables, "a plain table is never converted by recovery")
			saved := getGroup(t, f, groupID)
			if want == recordingplan.SchemaOperationSucceeded {
				require.Equal(t, token.OperationID, destinationOf(saved)["schema_revision"])
			} else {
				require.Empty(t, destinationOf(saved)["schema_revision"], "an unconfirmed result must not be saved as confirmed")
			}
		})
	}
}

// writeOneSampleThroughPipeline starts the production pipeline on the applied
// group, accepts one sample, closes its window, and reads the record back from
// the destination hypertable, which the first record gives its first chunk.
func writeOneSampleThroughPipeline(t *testing.T, f managedGroupFixture, groupID, table string) {
	t.Helper()
	snapshot, err := f.groups.ResolveAppliedAt(t.Context(), groupID, time.Now().Add(time.Minute))
	require.NoError(t, err)
	var clock atomic.Int64
	clock.Store(snapshot.EffectiveAt.Add(time.Second).UnixNano())
	targetSvc := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(f.db))
	pipe := grouppipeline.New(grouppipeline.Dependencies{
		Groups: f.groups, Tags: tag.NewService(tag.NewSQLRepository(f.db)), Destinations: targetSvc,
		Inspector: dbtarget.NewReadOnlyTableInspector(targetSvc), Store: groupdelivery.NewStore(f.db),
	}, grouppipeline.Config{
		NodeID: "node", Owner: "node/test", Now: func() time.Time { return time.Unix(0, clock.Load()).UTC() },
		TickInterval: 20 * time.Millisecond, WorkerInterval: 20 * time.Millisecond, ReconcileInterval: time.Hour, MaxFutureSkew: time.Minute,
	})
	require.NoError(t, pipe.Start(t.Context()))
	t.Cleanup(func() { _ = pipe.Stop(time.Second) })
	require.Equal(t, grouppipeline.StateActive, pipeStatus(pipe, groupID).State, "%+v", pipe.Status())

	member := snapshot.Group.Members[0]
	observed := snapshot.EffectiveAt.Add(2 * time.Second)
	require.NoError(t, pipe.AcceptSample(t.Context(), measurement.SampleEnvelope{
		SampleID: "ts-1", WorkspaceID: f.record.ID, DeviceID: "device-A", PointID: "point-A", TagID: "tag-A",
		SourceRevision: member.SourceRevision, MappingRevision: member.MappingRevision, ObservedAt: observed,
		ReceivedAt: observed.Add(30 * time.Millisecond), TimeOrigin: "gateway", Value: 21.5, Quality: schema.QualityGood,
	}))
	clock.Store(snapshot.EffectiveAt.Add(16 * time.Second).UnixNano())
	qualifiedTable := `"` + f.namespace + `"."` + table + `"`
	require.Eventually(t, func() bool {
		var n int
		return f.target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM `+qualifiedTable).Scan(&n) == nil && n == 1
	}, 15*time.Second, 50*time.Millisecond, "the closed window reaches the hypertable")
	var rowStatus, profile string
	var value float64
	require.NoError(t, f.target.QueryRowContext(t.Context(), `SELECT status, value_float, profile_kind FROM `+qualifiedTable).Scan(&rowStatus, &value, &profile))
	require.Equal(t, [3]any{"ok", 21.5, "instantaneous"}, [3]any{rowStatus, value, profile})
	var chunks int
	require.NoError(t, f.target.QueryRowContext(t.Context(), `SELECT num_chunks FROM timescaledb_information.hypertables WHERE hypertable_schema = $1 AND hypertable_name = $2`, f.namespace, table).Scan(&chunks))
	require.Equal(t, 1, chunks, "the first record created the first chunk")

}
