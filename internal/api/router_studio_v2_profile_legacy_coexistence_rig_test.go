package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
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

	"github.com/stretchr/testify/require"
)

// coexistRig holds, on one real destination, a legacy 60 second managed wide group, a custom table group and two
// profile groups of different circuits that share one long table. All four are applied and run by the production
// pipeline with deliveries held back, so their accepted rows stay pending for their frozen destinations.
type coexistRig struct {
	t      *testing.T
	f      managedGroupFixture
	kind   string
	dest   *sql.DB
	groups map[string]string // name -> group ID
	tags   map[string]string // name -> tag ID
}

func (r *coexistRig) quoted(table string) string {
	if r.kind == "sqlite" {
		return `"` + table + `"`
	}
	return `"` + r.f.namespace + `"."` + table + `"`
}

func (r *coexistRig) post(path string, body any) (code int, decoded map[string]any) {
	r.t.Helper()
	response := performJSONRequest(r.t, r.f.router, http.MethodPost, path, body)
	return response.Code, decodeJSONBody(r.t, response)
}

// apply applies a saved group through the public route with the revisions it was saved with.
func (r *coexistRig) apply(id string) (code int, decoded map[string]any) {
	r.t.Helper()
	return r.post(writeGroupsPath+"/"+id+"/apply", scopeBody(r.t, r.f, id))
}

// saveProfile saves one instantaneous profile group of its own circuit for one source.
func (r *coexistRig) saveProfile(circuit, source string) string {
	r.t.Helper()
	request := profileRequest(r.t, r.f.writeGroupRouterFixture, "instantaneous", 15, "second", circuit, source)
	group := request["group"].(map[string]any)
	group["layout_version"] = "managed_long_v1"
	group["destination"].(map[string]any)["table_schema"] = r.f.namespace
	code, body := r.post(ensureCircuitProfilePath, request)
	require.Equal(r.t, http.StatusOK, code, body)
	_, saved := groupSaveData(r.t, body)
	return saved["id"].(string)
}

// fresh stamps a request with the workspace revision a client that has just reloaded would send.
func (r *coexistRig) fresh(request map[string]any) map[string]any {
	r.t.Helper()
	current, err := r.f.workspace.GetOrCreate(r.t.Context())
	require.NoError(r.t, err)
	request["expected_workspace_revision"] = current.DatabaseSetupRevision
	return request
}

func newCoexistRig(t *testing.T, kind string) *coexistRig {
	t.Helper()
	f := newManagedDestinationFixture(t, kind)
	for _, name := range []string{"custom", "energy", "voltage", "extra"} {
		seedProfileSource(t, f.writeGroupRouterFixture, name)
	}
	r := &coexistRig{t: t, f: f, kind: kind, groups: map[string]string{}, tags: map[string]string{
		"wide": "tag-A", "custom": "tag-custom", "energy": "tag-energy", "voltage": "tag-voltage",
	}}
	r.dest = (&recoveryRig{t: t, f: f, kind: kind}).openDestination()

	// The legacy wide group: managed layout, one row per 60 seconds.
	request := f.createRequest()
	wide := request["group"].(map[string]any)
	wide["destination"].(map[string]any)["storage_strategy"] = "managed"
	wide["destination"].(map[string]any)["table_schema"] = f.namespace
	wide["row_policy"] = map[string]any{"interval_seconds": 60}
	code, body := r.post(writeGroupsPath, request)
	require.Equal(t, http.StatusCreated, code, body)
	_, saved := groupSaveData(t, body)
	r.groups["wide"] = saved["id"].(string)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, r.groups["wide"], previewLong(t, f, r.groups["wide"])).Status)

	// The custom group: a table the operator owns, columns they chose.
	floating := map[bool]string{true: "REAL", false: "DOUBLE PRECISION"}[kind == "sqlite"]
	_, err := r.dest.ExecContext(t.Context(), `CREATE TABLE `+r.quoted("coexist_custom")+` (c_col `+floating+` NOT NULL)`)
	require.NoError(t, err)
	request = r.fresh(f.createRequest())
	custom := request["group"].(map[string]any)
	custom["name"] = "Plant table"
	custom["members"] = []map[string]any{{"device_id": "device-A", "point_id": "point-custom", "tag_id": "tag-custom", "target_column": "c_col", "required": true}}
	custom["destination"] = map[string]any{"connector_id": "connector-A", "connector_revision": "connector-1", "table_schema": f.namespace, "table_name": "coexist_custom", "storage_strategy": "custom"}
	custom["write_policy"] = map[string]any{"dedupe_capability": "none"}
	custom["row_policy"] = map[string]any{"interval_seconds": 60}
	code, body = r.post(writeGroupsPath, request)
	require.Equal(t, http.StatusCreated, code, body)
	_, saved = groupSaveData(t, body)
	r.groups["custom"] = saved["id"].(string)

	// Two profile groups on two circuits share one long table; the first confirmation creates it.
	r.groups["energy"], r.groups["voltage"] = r.saveProfile("c-energy", "energy"), r.saveProfile("c-voltage", "voltage")
	for _, name := range []string{"energy", "voltage"} {
		require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, r.groups[name], previewLong(t, f, r.groups[name])).Status)
	}
	for _, name := range []string{"wide", "custom", "energy", "voltage"} {
		code, body := r.apply(r.groups[name])
		require.Equal(t, http.StatusOK, code, "apply %s: %v", name, body)
	}
	return r
}

// holdPendingRows runs the production pipeline with deliveries held back and closes one window per source.
func (r *coexistRig) holdPendingRows() {
	r.t.Helper()
	targets := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(r.f.db))
	var clock atomic.Int64
	pipe := grouppipeline.New(grouppipeline.Dependencies{
		Groups: r.f.groups, Tags: tag.NewService(tag.NewSQLRepository(r.f.db)), Destinations: targets,
		Inspector: dbtarget.NewReadOnlyTableInspector(targets), Store: groupdelivery.NewStore(r.f.db),
	}, grouppipeline.Config{
		NodeID: "node", Owner: "node/test", Now: func() time.Time { return time.Unix(0, clock.Load()).UTC() },
		TickInterval: 20 * time.Millisecond, WorkerInterval: time.Hour, ReconcileInterval: time.Hour, MaxFutureSkew: time.Minute,
	})
	latest := time.Time{}
	effective := map[string]time.Time{}
	for name, id := range r.groups {
		snapshot, err := r.f.groups.ResolveAppliedAt(r.t.Context(), id, time.Now().Add(time.Hour))
		require.NoError(r.t, err)
		effective[name] = snapshot.EffectiveAt
		if snapshot.EffectiveAt.After(latest) {
			latest = snapshot.EffectiveAt
		}
	}
	clock.Store(latest.Add(5 * time.Second).UnixNano())
	require.NoError(r.t, pipe.Start(r.t.Context()))
	r.t.Cleanup(func() { _ = pipe.Stop(time.Second) })
	for name, id := range r.groups {
		snapshot, err := r.f.groups.ResolveAppliedAt(r.t.Context(), id, time.Now().Add(time.Hour))
		require.NoError(r.t, err)
		member := snapshot.Group.Members[0]
		observed := effective[name].Add(2 * time.Second)
		require.NoError(r.t, pipe.AcceptSample(r.t.Context(), measurement.SampleEnvelope{
			SampleID: "coexist-" + name, WorkspaceID: r.f.record.ID, DeviceID: member.DeviceID, PointID: member.PointID, TagID: member.TagID,
			SourceRevision: member.SourceRevision, MappingRevision: member.MappingRevision, ObservedAt: observed, ReceivedAt: observed.Add(30 * time.Millisecond),
			TimeOrigin: "gateway", Value: 1.5, Quality: schema.QualityGood,
		}))
	}
	clock.Store(latest.Add(5 * time.Minute).UnixNano())
	require.Eventually(r.t, func() bool { return r.pendingRows() >= len(r.groups) }, 15*time.Second, 50*time.Millisecond, "every group closed a window into its pending outbox")
	time.Sleep(200 * time.Millisecond) // the closed windows are all in; nothing closes after the clock stops
}

func (r *coexistRig) pendingRows() int {
	var n int
	require.NoError(r.t, r.f.db.QueryRowContext(r.t.Context(), `SELECT COUNT(*) FROM wg_delivery_outbox WHERE state = 'pending'`).Scan(&n))
	return n
}

// ground is everything about a group that another group's lifecycle must never change: the saved group, its applied
// snapshots, its durable source ownership and the frozen destination of its pending rows.
func (r *coexistRig) ground(name string) string {
	r.t.Helper()
	id := r.groups[name]
	saved := getGroup(r.t, r.f, id)
	keep := map[string]any{}
	for _, key := range []string{"id", "revision", "applied_revision", "status", "layout_version", "storage_mode", "recording_profile", "row_policy", "write_policy", "members", "name"} {
		keep[key] = saved[key]
	}
	destination := cloneMap(destinationOf(saved))
	delete(destination, "schema_revision")
	delete(destination, "schema_digest")
	keep["destination"] = destination
	var versions, pending []string
	for _, row := range []struct {
		query string
		into  *[]string
	}{
		{`SELECT group_revision || ':' || payload FROM write_group_versions WHERE group_id = ? ORDER BY group_revision`, &versions},
		{`SELECT effect_key || ':' || connector_id || ':' || connector_revision || ':' || table_name || ':' || state || ':' || payload_digest FROM wg_delivery_outbox WHERE group_id = ? ORDER BY effect_key`, &pending},
	} {
		rows, err := r.f.db.QueryContext(r.t.Context(), row.query, id)
		require.NoError(r.t, err)
		for rows.Next() {
			var item string
			require.NoError(r.t, rows.Scan(&item))
			*row.into = append(*row.into, item)
		}
		require.NoError(r.t, rows.Err())
		require.NoError(r.t, rows.Close())
	}
	owned, err := groupdelivery.NewStore(r.f.db).HasAppliedOwner(r.t.Context(), "connector-A", r.tags[name])
	require.NoError(r.t, err)
	keep["versions"], keep["pending"], keep["owner"] = versions, pending, owned
	encoded, err := json.Marshal(keep)
	require.NoError(r.t, err)
	return string(encoded)
}

func (r *coexistRig) tableCount() int {
	var n int
	query, args := `SELECT count(*) FROM sqlite_master WHERE type='table'`, []any{}
	if r.kind == "postgres" {
		query, args = `SELECT count(*) FROM information_schema.tables WHERE table_schema=$1`, []any{r.f.namespace}
	}
	require.NoError(r.t, r.dest.QueryRowContext(r.t.Context(), query, args...).Scan(&n))
	return n
}
