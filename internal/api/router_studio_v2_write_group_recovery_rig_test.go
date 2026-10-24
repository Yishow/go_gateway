package api

import (
	"database/sql"
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
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

const (
	recoveryWindow           = 60 * time.Second
	recoveryCumulativeWindow = 10 * time.Minute
	recoveryWorkerEvery      = 20 * time.Millisecond
)

// recoveryGroup is one applied long group of the rig. The groups of one rig
// share a destination and therefore one shared long table.
type recoveryGroup struct {
	id, circuit, device, point, tag string
	snapshot                        *workspace.WriteGroupAppliedSnapshot
}

// recoveryRig is the production path of a long group on a disposable
// destination: saved through the public routes, confirmed and applied, then run
// by the production pipeline. Nothing about the confirmed state is built by the
// test; it only drops tables and reads the destination back.
type recoveryRig struct {
	t      *testing.T
	f      managedGroupFixture
	kind   string
	table  string
	groups []recoveryGroup
	dest   *sql.DB
	store  *groupdelivery.Store
	clock  atomic.Int64
	pipe   *grouppipeline.Pipeline
	// base is the start of window zero, shared by every group.
	base time.Time
	next int
	// window is the closing interval of the rig's groups; cumulative rigs run
	// one counter group on a saved measurement definition.
	window     time.Duration
	cumulative bool
	options    []func(*grouppipeline.Dependencies)
}

// newRecoveryRig prepares the groups (the first on the seeded tag, the others
// on extra sources) up to Apply, without starting the pipeline. kind is
// "sqlite", "postgres" or "timescale".
func newRecoveryRig(t *testing.T, kind string, extraSources ...string) *recoveryRig {
	t.Helper()
	return newRecoveryRigWithPolicy(t, kind, nil, extraSources...)
}

func newRecoveryRigWithPolicy(t *testing.T, kind string, recreateMissing *bool, extraSources ...string) *recoveryRig {
	t.Helper()
	return buildRecoveryRig(t, kind, recreateMissing, false, extraSources...)
}

// newCumulativeRecoveryRig prepares one cumulative group counting a uint64
// energy counter from a saved measurement definition, windows of ten minutes.
func newCumulativeRecoveryRig(t *testing.T, kind string) *recoveryRig {
	t.Helper()
	return buildRecoveryRig(t, kind, nil, true)
}

func buildRecoveryRig(t *testing.T, kind string, recreateMissing *bool, cumulative bool, extraSources ...string) *recoveryRig {
	t.Helper()
	destination := kind
	if kind == "timescale" {
		destination = "postgres"
	}
	var tagTypes []string
	if cumulative {
		tagTypes = []string{"uint64"}
	}
	f := newManagedDestinationFixture(t, destination, tagTypes...)
	if kind == "timescale" && !timescaleInstalled(t, f) {
		t.Skip("this destination database has no TimescaleDB extension; the timescale case is NOT RUN here")
	}
	r := &recoveryRig{t: t, f: f, kind: kind, window: recoveryWindow, cumulative: cumulative}
	if cumulative {
		r.window = recoveryCumulativeWindow
		r.options = []func(*grouppipeline.Dependencies){func(deps *grouppipeline.Dependencies) {
			deps.Counters = grouppipeline.NewMeasurementCounterSources(measurement.NewSQLRepository(f.db))
		}}
		_, err := f.db.ExecContext(t.Context(), `
			INSERT INTO measurement_definitions (
				id, workspace_id, device_id, point_id, tag_id, equipment_id, definition_revision, source_binding_revision,
				series_epoch, name, quantity, unit, semantic_kind
			) VALUES ('meas-energy', ?, 'device-A', 'point-A', 'tag-A', 'device-A', 'def-5', 'bind-5', 'epoch-3', 'Energy', 'energy', 'kWh', 'counter')`,
			f.record.ID)
		require.NoError(t, err)
	}
	r.dest = r.openDestination()
	sources := append([]string{""}, extraSources...)
	for i, source := range sources {
		if source != "" {
			seedProfileSource(t, f.writeGroupRouterFixture, source)
		}
		r.groups = append(r.groups, r.saveGroup(i, source, recreateMissing))
	}
	// Every group is confirmed first and then applied back to back: a group takes
	// effect at the next boundary of its window, so applies spread over a slow
	// setup could straddle one and leave an earlier group an empty first window.
	for i := range r.groups {
		r.confirm(&r.groups[i])
	}
	for i := range r.groups {
		r.apply(&r.groups[i])
		require.True(t, r.groups[i].snapshot.EffectiveAt.Equal(r.groups[0].snapshot.EffectiveAt),
			"the groups take effect together: %s and %s; run the test again", r.groups[0].snapshot.EffectiveAt, r.groups[i].snapshot.EffectiveAt)
	}
	r.table = destinationOf(getGroup(t, f, r.groups[0].id))["table_name"].(string)
	return r
}

// openDestination returns an administrative handle on the destination database
// for catalog reads and for dropping the tables the tests remove.
func (r *recoveryRig) openDestination() *sql.DB {
	r.t.Helper()
	if r.kind != "sqlite" {
		return r.f.target
	}
	db, err := sql.Open("sqlite", r.f.file+"?_pragma=busy_timeout(5000)")
	require.NoError(r.t, err)
	r.t.Cleanup(func() { _ = db.Close() })
	return db
}

func (r *recoveryRig) saveGroup(index int, source string, recreateMissing *bool) recoveryGroup {
	r.t.Helper()
	circuit := "recovery-" + string(rune('a'+index))
	request := profileRequest(r.t, r.f.writeGroupRouterFixture, "instantaneous", int(recoveryWindow/time.Second), "second", circuit, source)
	group := request["group"].(map[string]any)
	if r.cumulative {
		request = profileRequest(r.t, r.f.writeGroupRouterFixture, "cumulative", 10, "minute", circuit, source)
		group = request["group"].(map[string]any)
		group["members"].([]map[string]any)[0]["measurement_id"] = "meas-energy"
	}
	group["layout_version"] = "managed_long_v1"
	group["storage_mode"] = "plain"
	if r.kind == "timescale" {
		group["storage_mode"] = "timescale"
	}
	if recreateMissing != nil {
		// The request carries no policy of its own: only the explicit choice is sent.
		group["write_policy"] = map[string]any{"recreate_missing": *recreateMissing}
	}
	group["destination"].(map[string]any)["table_schema"] = r.f.namespace
	response := performJSONRequest(r.t, r.f.router, http.MethodPost, ensureCircuitProfilePath, request)
	require.Equal(r.t, http.StatusOK, response.Code, response.Body.String())
	_, saved := groupSaveData(r.t, decodeJSONBody(r.t, response))
	device, point, tagID := "device-A", "point-A", "tag-A"
	if source != "" {
		point, tagID = "point-"+source, "tag-"+source
	}
	return recoveryGroup{id: saved["id"].(string), circuit: circuit, device: device, point: point, tag: tagID}
}

func (r *recoveryRig) confirm(g *recoveryGroup) {
	r.t.Helper()
	operation := confirmLong(r.t, r.f, g.id, previewLong(r.t, r.f, g.id))
	require.Equal(r.t, recordingplan.SchemaOperationSucceeded, operation.Status)
	ready, err := r.f.groups.Readiness(r.t.Context(), g.id)
	require.NoError(r.t, err)
	require.True(r.t, ready.Ready, "%+v", ready.Issues)
}

func (r *recoveryRig) apply(g *recoveryGroup) {
	r.t.Helper()
	applied := performJSONRequest(r.t, r.f.router, http.MethodPost, writeGroupsPath+"/"+g.id+"/apply", scopeBody(r.t, r.f, g.id))
	require.Equal(r.t, http.StatusOK, applied.Code, applied.Body.String())
	// A cumulative group takes effect at the next window boundary, ten minutes at most.
	snapshot, err := r.f.groups.ResolveAppliedAt(r.t.Context(), g.id, time.Now().Add(time.Hour))
	require.NoError(r.t, err)
	g.snapshot = snapshot
}

// start runs the production pipeline over the applied groups. The clock is
// driven by the test, so windows close exactly when it says they do.
func (r *recoveryRig) start() {
	r.t.Helper()
	for _, g := range r.groups {
		if g.snapshot.EffectiveAt.After(r.base) {
			r.base = g.snapshot.EffectiveAt
		}
	}
	if aligned := r.base.Truncate(r.window); aligned.Before(r.base) {
		r.base = aligned.Add(r.window)
	}
	r.clock.Store(r.base.Add(time.Second).UnixNano())
	r.launch(recoveryWorkerEvery)
}

// restart stops the running pipeline and starts a new one over the same local
// database, as a process restart does. workerEvery is the pause between delivery
// cycles: a long one keeps closed windows pending after the first cycle.
func (r *recoveryRig) restart(workerEvery time.Duration) {
	r.t.Helper()
	require.NoError(r.t, r.pipe.Stop(time.Second))
	r.launch(workerEvery)
}

func (r *recoveryRig) launch(workerEvery time.Duration) {
	r.t.Helper()
	r.store = groupdelivery.NewStore(r.f.db)
	targets := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(r.f.db))
	deps := grouppipeline.Dependencies{
		Groups: r.f.groups, Tags: tag.NewService(tag.NewSQLRepository(r.f.db)), Destinations: targets,
		Inspector: dbtarget.NewReadOnlyTableInspector(targets), Store: r.store,
	}
	for _, option := range r.options {
		option(&deps)
	}
	r.pipe = grouppipeline.New(deps, grouppipeline.Config{
		NodeID: "node", Owner: "node/test", Now: func() time.Time { return time.Unix(0, r.clock.Load()).UTC() },
		TickInterval: 20 * time.Millisecond, WorkerInterval: workerEvery, ReconcileInterval: time.Hour, MaxFutureSkew: time.Minute,
	})
	r.f.router = NewRouter(&DatalinkServices{
		Workspace: r.f.workspace, WriteGroups: r.f.groups, WriteGroupDelivery: r.pipe,
		RecordingPlan: recordingplan.NewService(recordingplan.NewSQLRepository(r.f.db)),
	})
	require.NoError(r.t, r.pipe.Start(r.t.Context()))
	pipe := r.pipe
	r.t.Cleanup(func() { _ = pipe.Stop(time.Second) })
	for _, g := range r.groups {
		require.Equal(r.t, grouppipeline.StateActive, pipeStatus(r.pipe, g.id).State, "%+v", r.pipe.Status())
	}
}

// record accepts one reading per group inside the next window and moves the
// clock past its end, so the window closes.
func (r *recoveryRig) record(value any) {
	r.t.Helper()
	window := r.next
	r.next++
	start := r.base.Add(time.Duration(window) * r.window)
	for _, g := range r.groups {
		member := g.snapshot.Group.Members[0]
		observed := start.Add(2 * time.Second)
		r.clock.Store(observed.Add(time.Second).UnixNano())
		require.NoError(r.t, r.pipe.AcceptSample(r.t.Context(), measurement.SampleEnvelope{
			SampleID: g.id + "-" + start.Format(time.RFC3339), WorkspaceID: r.f.record.ID, DeviceID: g.device, PointID: g.point, TagID: g.tag,
			SourceRevision: member.SourceRevision, MappingRevision: member.MappingRevision, ObservedAt: observed,
			ReceivedAt: observed.Add(30 * time.Millisecond), TimeOrigin: "gateway", Value: value, Quality: schema.QualityGood,
		}))
	}
	r.clock.Store(start.Add(r.window + time.Second).UnixNano())
}
