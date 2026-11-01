package runtime

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestAcquisitionScopedProjectionNeverAdoptsLateForeignScope(t *testing.T) {
	f := newMappingObservationFixture(t, true)
	f.projection.Rules = []*schema.SourceRule{{ID: "R", DeviceID: "D", RevisionID: "rule-original", Enabled: true}}
	f.projection.RuleLinks = []*schema.SourceRuleLink{{RuleID: "R", PointID: "P", Address: "40001"}}
	f.install(t)
	f.service.running.Store(true)
	installed, ok := f.service.scheduler.CaptureOrdinaryPoint("P")
	require.True(t, ok)
	require.NotZero(t, installed.Lease)
	scope := sourcerule.RuntimeReconcileScope{DeviceID: "D", RuleID: "R"}
	snapshot := f.service.SnapshotSourceRuleProjection("D", []string{"P"})
	ctx, finish := f.service.BeginSourceRuleAcquisitionMutation(t.Context(), scope)
	defer finish()
	foreign := f.service.deviceProjectionState("D")
	foreign.RuntimeVersion = "foreign-runtime-version"
	foreign.WorkspaceVersion = "foreign-workspace-version"
	foreign.Alignment = workspace.RuntimeProjectionAlignmentStale
	foreign.Code = "foreign-source-installation"
	foreign.Message = "foreign installation must survive"
	changed := false
	f.reader.beforeRead = func(int) {
		if changed {
			return
		}
		changed = true
		// This is only a scope publication: mapping pointers and source leases
		// deliberately stay unchanged, so their guards cannot stand in for it.
		f.service.projectionMu.Lock()
		f.service.projections["D"] = foreign
		f.service.projectionMu.Unlock()
	}
	outcome := f.service.ReconcileSourceRule(ctx, sourcerule.RuntimeReconcileRequest{
		Operation: sourcerule.RuntimeReconcileOperationUpdate, Scope: scope,
		AffectedPointIDs: []string{"P"}, ProjectionSnapshot: snapshot,
		ProjectionProof: &sourcerule.RuntimeProjectionProof{Epoch: installed.Epoch, Points: map[string]sourcerule.RuntimePointProof{
			"P": {Revision: installed.Revision, Lease: installed.Lease},
		}},
	})
	require.True(t, changed, "the real mapping observation provider must reenter publication")
	require.Equal(t, foreign, f.service.deviceProjectionState("D"), "a late getter must not adopt and overwrite the foreign scope")
	require.Equal(t, sourcerule.RuntimeReconcileStatusStale, outcome.Status)
}

// The embedded service retains the real SQL graph acceptance and projection
// readers. Only its existing optional observation read can synchronously reenter.
type acquisitionInputReentryWorkspace struct {
	*workspace.Service
	beforeRead func()
}

func (r *acquisitionInputReentryWorkspace) ReadSnapshot(ctx context.Context) (*workspace.Record, error) {
	if r.beforeRead != nil {
		r.beforeRead()
	}
	return r.Service.ReadSnapshot(ctx)
}

func newAcquisitionInputSQLFixture(t *testing.T) (*Service, *acquisitionInputReentryWorkspace, *workspace.RuntimeProjection) {
	t.Helper()
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "source-input.db"))
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	for _, statement := range []string{
		`INSERT INTO devices(id,name,description,protocol,status,connection_config) VALUES('d','PLC','','mc_3e','active','{"host":"127.0.0.1","port":5007,"acquisition":{"version":1,"profile":"safe-v1","coverage":{"revision":"readable-original","readable":[{"area":"D","start":0,"end":21}]}}}')`,
		`INSERT INTO points(id,device_id,name,address,data_type,revision_id) VALUES('a','d','a','D0','uint16','pa'),('b','d','b','D2','uint16','pb')`,
		`INSERT INTO source_rules(id,device_id,start_address,count,data_type,naming_prefix,revision_id) VALUES('ra','d','D0',1,'uint16','a','ra-original'),('rb','d','D2',1,'uint16','b','rb-original')`,
		`INSERT INTO source_rule_links(id,rule_id,address,point_id) VALUES('la','ra','D0','a'),('lb','rb','D2','b')`,
		`INSERT INTO system_settings(key,value) VALUES('studio_v2_workspace','{"id":"w","ordered_device_ids":["d"],"database_setup_revision":"workspace-original"}')`,
	} {
		_, err := db.ExecContext(t.Context(), statement)
		require.NoError(t, err)
	}
	ws := workspace.NewService(workspace.NewSQLRepository(db)).WithRuntimeProjectionServices(
		device.NewSQLRepository(db), sourcerule.NewService(sourcerule.NewSQLRepository(db), nil, nil, nil),
		point.NewSQLRepository(db), mapping.NewSQLRepository(db), tag.NewSQLRepository(db), nil, nil,
		point.NewMemoryPollingGroupRepository(),
	)
	projection, err := ws.RuntimeProjection(t.Context())
	require.NoError(t, err)
	graph, complete := projection.AcquisitionSourceGraph("d")
	require.True(t, complete, graph.Reason())
	require.Len(t, graph.Points(), 2, "the graph must come from the complete canonical SQL source set")
	for _, p := range projection.Points {
		require.False(t, modernPoint(p), "legacy fixture must not borrow a modern ReadProfile guard")
	}
	reader := &acquisitionInputReentryWorkspace{Service: ws}
	cfg := collector.DefaultSchedulerConfig()
	cfg.AcquisitionObserver = func(collector.AcquisitionObservation) {}
	live := &Service{scheduler: collector.NewScheduler(cfg, nil), workspace: reader, mappingIndex: map[string][]mappingBinding{}, pointMetaIndex: map[string]pointMeta{}}
	return live, reader, projection
}

func TestAcquisitionFullProjectionFreezesOriginalSQLSourcePayload(t *testing.T) {
	for _, mutation := range []string{"points-slice", "legacy-address"} {
		t.Run(mutation, func(t *testing.T) {
			live, reader, projection := newAcquisitionInputSQLFixture(t)
			original := make(map[string]schema.Point, len(projection.Points))
			for _, p := range projection.Points {
				original[p.ID] = *p
			}
			changed := false
			reader.beforeRead = func() {
				if changed {
					return
				}
				changed = true
				// The SQL witness remains untouched. Only the caller-owned payload
				// changes after its original graph was captured by the full apply.
				if mutation == "points-slice" {
					projection.Points = projection.Points[:1]
				} else {
					projection.Points[0].Address = "D20"
				}
			}
			err := live.ApplyWorkspaceProjection(t.Context(), projection)
			require.True(t, changed, "the actual optional observation provider must run")
			_, stillOriginal := projection.AcquisitionSourceGraph("d")
			require.False(t, stillOriginal, "the mutated public payload no longer matches the opaque original SQL graph")
			state := live.deviceProjectionState("d")
			if err != nil {
				require.False(t, state.acquisitionGraph.Complete(), "rejected input must not publish a complete source grant")
				return
			}
			// Success is legal only if the full apply detached its original input
			// before provider reentry, not if original SQL is joined to new slots.
			for id, p := range original {
				require.NotZero(t, live.scheduler.CapturePointLease(id, p.RevisionID), "successful full publication must retain every frozen SQL point")
				require.Equal(t, pointMeta{DeviceID: p.DeviceID, Address: p.Address}, live.lookupPointMeta(id), "complete SQL graph must never authorize caller-mutated source geometry")
			}
		})
	}
}
