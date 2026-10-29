package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/api/handlers"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/workspace"
)

type runtimeDeliveryRig struct {
	*runtimeRuleEvidenceRig
	store *groupdelivery.Store
}

// Configuration Apply and durable acceptance are production services; only the
// sample envelope and clock are simulated. No background workers or real PLC.
func newRuntimeDeliveryRig(t *testing.T) *runtimeDeliveryRig {
	t.Helper()
	r := &runtimeDeliveryRig{runtimeRuleEvidenceRig: newRuntimeRuleEvidenceRig(t)}
	r.clock = time.Date(2026, 10, 2, 12, 0, 16, 0, time.UTC)
	r.store = groupdelivery.NewStore(r.f.db).WithClock(func() time.Time { return r.clock })
	r.pipe = grouppipeline.New(grouppipeline.Dependencies{Groups: r.f.groups, Store: r.store,
		Destinations: dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(r.f.db))},
		grouppipeline.Config{Now: func() time.Time { return r.clock }})
	require.NoError(t, r.pipe.Reconcile(t.Context()))
	r.router = NewRouter(&DatalinkServices{Workspace: r.f.workspace, WriteGroups: r.f.groups, WriteGroupDelivery: r.pipe,
		Device: device.NewService(device.NewSQLRepository(r.f.db), nil)})
	return r
}

func (r *runtimeDeliveryRig) accept(t *testing.T, id string) {
	t.Helper()
	snap, err := r.f.groups.ResolveAppliedAt(t.Context(), r.f.group["id"].(string), r.clock)
	require.NoError(t, err)
	for _, member := range snap.Group.Members {
		require.NoError(t, r.pipe.AcceptSample(t.Context(), measurement.SampleEnvelope{SampleID: id + "-" + member.DeviceID,
			WorkspaceID: r.f.record.ID, DeviceID: member.DeviceID, PointID: member.PointID, TagID: member.TagID,
			SourceRevision: member.SourceRevision, MappingRevision: member.MappingRevision, ObservedAt: r.clock,
			ReceivedAt: r.clock, Value: float64(12), Quality: schema.QualityGood}))
	}
}

func (r *runtimeDeliveryRig) evidence(t *testing.T) runtimeview.GroupView {
	t.Helper()
	response := performJSONRequest(t, r.router, http.MethodGet, runtimeRuleContextPath+"?"+r.query(t).Encode(), nil)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	for _, unsafe := range []string{r.f.file, "private-reading-secret", "/private/fixture-only.sqlite", `"payload"`, `"connection_config"`} {
		require.NotContains(t, response.Body.String(), unsafe)
	}
	var body struct {
		Data struct {
			Evidence runtimeview.RuleEvidenceView `json:"rule_evidence"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
	require.Len(t, body.Data.Evidence.Groups, 1)
	return body.Data.Evidence.Groups[0]
}

func (r *runtimeDeliveryRig) useReader(reader handlers.RuntimeRevisionDeliveryReader) {
	h := handlers.NewRuntimeHandler(device.NewService(device.NewSQLRepository(r.f.db), nil), nil, nil, nil, nil, r.f.workspace).
		WithRuleEvidenceServices(r.f.groups, r.pipe).WithRuleEvidenceClock(func() time.Time { return r.clock }).WithRuleDeliveryServices(reader)
	r.router = gin.New()
	r.router.GET(runtimeRuleContextPath, h.WorkspaceContext)
}

func (r *runtimeDeliveryRig) tick(t *testing.T) string {
	t.Helper()
	r.clock = r.clock.Truncate(15 * time.Second).Add(15 * time.Second)
	r.pipe.TickAll(t.Context())
	var key string
	err := r.f.db.QueryRow(`SELECT effect_key FROM wg_delivery_outbox WHERE group_revision=? ORDER BY bucket_start DESC LIMIT 1`, r.f.group["applied_revision"]).Scan(&key)
	if err != nil {
		t.Logf("tick clock=%s applied=%s statuses=%+v", r.clock, r.f.group["applied_revision"], r.pipe.Status())
		status, statusErr := r.store.GroupStatus(t.Context(), r.f.group["id"].(string))
		t.Logf("durable status=%+v err=%v", status, statusErr)
	}
	require.NoError(t, err)
	return key
}

func (r *runtimeDeliveryRig) Resolve(ctx context.Context, item groupdelivery.OutboxItem) (groupdelivery.Target, error) {
	connectors := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(r.f.db))
	opened, err := connectors.OpenDestination(ctx, item.Destination.ConnectorID, item.Destination.ConnectorRevision)
	if err != nil {
		return groupdelivery.Target{}, err
	}
	return groupdelivery.Target{DB: opened.DB, Kind: opened.Kind, Close: opened.Close}, nil
}

func (r *runtimeDeliveryRig) send(t *testing.T, key string) {
	t.Helper()
	sender := groupdelivery.NewSender(r.store, r, groupdelivery.SenderConfig{Owner: "runtime-delivery-api", Now: func() time.Time { return r.clock }})
	result, err := sender.Deliver(t.Context(), key)
	require.NoError(t, err)
	require.Equal(t, groupdelivery.StateCommitted, result.State)
}

// Seed only retained durable output, not effective configuration proof. The
// same production Store commit/transition APIs persist all stage/backlog rows.
func (r *runtimeDeliveryRig) historicalRows(t *testing.T, n int) []string {
	t.Helper()
	key := groupdelivery.GroupKey{WorkspaceID: r.f.record.ID, GroupID: r.f.group["id"].(string), GroupRevision: "rev7"}
	dest := groupdelivery.FrozenDestination{Scope: "historical-scope", ConnectorID: "C1", ConnectorRevision: "7", Database: "/private/fixture-only.sqlite", TableSchema: "main", TableName: "old_table", DedupeCapability: "none"}
	keys := make([]string, 0, n)
	for i := range n {
		start := r.clock.Add(time.Duration(i-100) * 15 * time.Second)
		effect := fmt.Sprintf("historical-effect-%d", i)
		record := fmt.Sprintf("historical-record-%d", i)
		outcome := snapshot.Outcome{Kind: snapshot.OutcomeRow, RecordID: record, EffectKey: effect, BucketStart: start, BucketEnd: start.Add(15 * time.Second)}
		row := dbtarget.EncodedRow{RecordID: record, EffectKey: effect, BucketStart: start, Cells: []dbtarget.EncodedCell{{Column: "value", Value: float64(12)}}}
		require.NoError(t, r.store.CommitClosure(t.Context(), key, groupdelivery.Closure{Destination: dest, NextClose: start.Add(15 * time.Second), Buckets: []groupdelivery.ClosedBucket{{Outcome: outcome, Row: &row}}}))
		keys = append(keys, effect)
	}
	return keys
}

type runtimeDeliveryFault struct {
	*grouppipeline.Pipeline
	after    func()
	change   func(*grouppipeline.DeliveryView)
	err      error
	nilView  bool
	statuses func(int, []grouppipeline.GroupStatus) []grouppipeline.GroupStatus
	calls    int
}

func (f *runtimeDeliveryFault) DeliveryForRevision(ctx context.Context, id, revision string) (*grouppipeline.DeliveryView, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.nilView {
		return nil, nil
	}
	v, err := f.Pipeline.DeliveryForRevision(ctx, id, revision)
	if err == nil && f.change != nil {
		f.change(v)
	}
	if f.after != nil {
		f.after()
	}
	return v, err
}
func (f *runtimeDeliveryFault) Status() []grouppipeline.GroupStatus {
	f.calls++
	s := f.Pipeline.Status()
	if f.statuses != nil {
		return f.statuses(f.calls, s)
	}
	return s
}

func newRuntimeCounterDeliveryRig(t *testing.T, reason string) *runtimeDeliveryRig {
	t.Helper()
	f := newManagedDestinationFixture(t, "sqlite", "uint64")
	r := &runtimeDeliveryRig{runtimeRuleEvidenceRig: &runtimeRuleEvidenceRig{f: f, clock: time.Date(2026, 10, 1, 0, 0, 1, 0, time.UTC)}}
	f.groups.WithClock(func() time.Time { return r.clock })
	definition := map[string]any{"id": "meas-energy", "workspace_id": f.record.ID, "device_id": "device-A", "point_id": "point-A", "tag_id": "tag-A", "equipment_id": "main", "name": "Energy", "quantity": "energy", "unit": "kWh", "semantic_kind": "counter", "numeric_encoding": "uint64", "definition_revision": "def-1", "source_binding_revision": "bind-1", "series_epoch": "epoch-1"}
	response := performJSONRequest(t, f.router, http.MethodPost, measurementsPath, definition)
	require.Equal(t, http.StatusCreated, response.Code, response.Body.String())
	saved := ensureProfileGroup(t, f, "cumulative", 10, "minute", "meas-energy")
	id := saved["id"].(string)
	confirmLong(t, f, id, previewLong(t, f, id))
	response = performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/apply", scopeBody(t, f, id))
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	r.f.data = getGroupData(t, f, id)
	r.f.group = r.f.data["group"].(map[string]any)
	r.store = groupdelivery.NewStore(f.db).WithClock(func() time.Time { return r.clock })
	key := groupdelivery.GroupKey{WorkspaceID: f.record.ID, GroupID: id, GroupRevision: appliedRevisionOf(t, f, id)}
	original, raw := counterReplayAPISnapshot(t, r.store, key)
	r.clock = original.EffectiveAt.Add(time.Second)
	deps := grouppipeline.Dependencies{Groups: f.groups, Store: r.store, Destinations: dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(f.db)), Counters: grouppipeline.NewMeasurementCounterSources(measurement.NewSQLRepository(f.db))}
	r.pipe = grouppipeline.New(deps, grouppipeline.Config{Now: func() time.Time { return r.clock }})
	require.NoError(t, r.pipe.Reconcile(t.Context()))
	m := original.Group.Members[0]
	for i := range 3 {
		at := r.clock.Add(time.Duration(i) * time.Second)
		require.NoError(t, r.pipe.AcceptSample(t.Context(), measurement.SampleEnvelope{SampleID: fmt.Sprintf("replay-%d", i), WorkspaceID: f.record.ID, DeviceID: m.DeviceID, PointID: m.PointID, TagID: m.TagID, SourceRevision: m.SourceRevision, MappingRevision: m.MappingRevision, ObservedAt: at, ReceivedAt: at, Value: uint64(100 + i), Quality: schema.QualityGood}))
	}
	var descriptor map[string]any
	require.NoError(t, json.Unmarshal(raw, &descriptor))
	if reason == "counter-replay-proof-missing" {
		delete(descriptor, "counter_bindings")
	} else {
		for _, entry := range descriptor["counter_bindings"].(map[string]any) {
			proof := entry.(map[string]any)
			if reason == "counter-replay-proof-conflict" {
				proof["unit"] = "private-reading-secret"
			} else {
				proof["version"] = 77
			}
		}
	}
	damaged, err := json.Marshal(descriptor)
	require.NoError(t, err)
	_, err = f.db.Exec(`UPDATE wg_runtime_versions SET payload=? WHERE group_id=? AND group_revision=?`, string(damaged), id, key.GroupRevision)
	require.NoError(t, err)
	// Recovery is through public Reconcile; no private quarantine map injection.
	r.pipe = grouppipeline.New(deps, grouppipeline.Config{Now: func() time.Time { return r.clock }})
	require.NoError(t, r.pipe.Reconcile(t.Context()))
	_, err = f.db.Exec(`UPDATE devices SET description='' WHERE id='device-A'`)
	require.NoError(t, err)
	r.router = NewRouter(&DatalinkServices{Workspace: f.workspace, WriteGroups: f.groups, WriteGroupDelivery: r.pipe, Device: device.NewService(device.NewSQLRepository(f.db), nil)})
	return r
}

func (r *runtimeDeliveryRig) sharedGroup(t *testing.T) {
	t.Helper()
	_, err := r.f.db.Exec(`INSERT INTO devices(id,name,description,protocol,status,connection_config,readiness_status) VALUES('device-B','Second','','modbus_tcp','draft','{}','{}')`)
	require.NoError(t, err)
	_, err = r.f.workspace.AttachDevice(t.Context(), "device-B")
	require.NoError(t, err)
	_, err = r.f.db.Exec(`INSERT INTO points(id,device_id,name,address,function,data_type,mode,enabled) VALUES('point-B','device-B','Second','40001','FC03','float32','read',1)`)
	require.NoError(t, err)
	_, err = r.f.db.Exec(`INSERT INTO tags(id,key,key_lower,display_name,data_type,status,labels) VALUES('tag-B','second','second','Second','float32','active','{}')`)
	require.NoError(t, err)
	_, err = r.f.db.Exec(`INSERT INTO mappings(id,point_id,tag_id,transform_pipeline,status,enabled) VALUES('mapping-B','point-B','tag-B','[]','active',1)`)
	require.NoError(t, err)
	current, err := r.f.groups.Get(t.Context(), r.f.group["id"].(string))
	require.NoError(t, err)
	current.Group.Destination.TableName = "meter_shared"
	current.Group.Members = append(current.Group.Members, workspace.WriteGroupMember{DeviceID: "device-B", PointID: "point-B", TagID: "tag-B", TargetColumn: "second", Required: true})
	_, err = r.f.groups.Update(t.Context(), current.Group.ID, workspace.WriteGroupMutation{WorkspaceID: current.Group.WorkspaceID, ExpectedWorkspaceRevision: current.WorkspaceRevision, ExpectedGroupRevision: current.Group.Revision, ExpectedConnectorRevision: current.Group.Destination.ConnectorRevision, Group: current.Group})
	require.NoError(t, err)
	r.f.data = getGroupData(t, r.f, current.Group.ID)
	r.f.group = r.f.data["group"].(map[string]any)
	r.apply(t)
	r.clock = r.clock.Truncate(15 * time.Second).Add(16 * time.Second)
	require.NoError(t, r.pipe.Reconcile(t.Context()))
	require.Contains(t, r.pipe.Status(), grouppipeline.GroupStatus{GroupID: current.Group.ID, Revision: r.f.group["applied_revision"].(string), State: grouppipeline.StateActive})
}
