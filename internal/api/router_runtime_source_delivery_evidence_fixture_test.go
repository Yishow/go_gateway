package api

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-gateway/internal/api/handlers"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/tag"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

type runtimeSourceDeliveryRig struct {
	*runtimeDeliveryRig
	snapshot *workspace.WriteGroupAppliedSnapshot
	window   time.Duration
}

// Every applied/frozen fact comes from production create, preview, confirm and
// Apply. Public Reconcile builds the boundary; no background worker is started.
type runtimeSourceDeliveryOptions struct{ wide, cumulative, multi bool }

func newRuntimeSourceDeliveryRig(t *testing.T, wide bool) *runtimeSourceDeliveryRig {
	return newRuntimeSourceDeliveryOptions(t, runtimeSourceDeliveryOptions{wide: wide})
}
func newRuntimeSourceDeliveryOptions(t *testing.T, options runtimeSourceDeliveryOptions) *runtimeSourceDeliveryRig {
	t.Helper()
	f := newManagedDestinationFixture(t, "sqlite")
	base := &runtimeRuleEvidenceRig{f: f, clock: time.Date(2026, 10, 2, 12, 0, 7, 0, time.UTC)}
	r := &runtimeSourceDeliveryRig{runtimeDeliveryRig: &runtimeDeliveryRig{runtimeRuleEvidenceRig: base}, window: 15 * time.Second}
	f.groups.WithClock(func() time.Time { return r.clock })
	sources := []meanSource{{name: "mean", tagType: "int16", role: "gauge"}, {name: "other", tagType: "int16", role: "gauge"}, {name: "quiet", tagType: "int16", role: "gauge"}}
	if options.cumulative {
		r.window = time.Minute
		for i := range sources {
			sources[i].tagType = "uint64"
			sources[i].role = ""
			sources[i].measurementID = "meas-" + sources[i].name
		}
	}
	seedMeanSources(t, f, sources)
	for _, source := range sources {
		if source.measurementID != "" {
			insertMeasurement(t, f, source.measurementID, source.name, "epoch-1", "counter")
		}
	}
	if options.multi {
		_, err := f.db.Exec(`INSERT INTO devices(id,name,description,protocol,status,connection_config,readiness_status) VALUES('device-B','Second','','modbus_tcp','draft','{}','{}')`)
		require.NoError(t, err)
		_, err = f.workspace.AttachDevice(t.Context(), "device-B")
		require.NoError(t, err)
		_, err = f.db.Exec(`UPDATE points SET device_id='device-B' WHERE id='point-other'`)
		require.NoError(t, err)
	}
	request := meanLongRequest(t, f, sources)
	if options.wide {
		request = wideProfileRequest(t, f, sources)
	}
	group := request["group"].(map[string]any)
	if options.cumulative {
		group["recording_profile"] = map[string]any{"kind": "cumulative", "interval": 1, "unit": "minute"}
		for _, member := range group["members"].([]map[string]any) {
			delete(member, "value_role")
		}
	}
	if options.multi {
		group["members"].([]map[string]any)[1]["device_id"] = "device-B"
	}
	created := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, request)
	require.Equal(t, http.StatusOK, created.Code, created.Body.String())
	r.f.data, r.f.group = groupSaveData(t, decodeJSONBody(t, created))
	id := r.f.group["id"].(string)
	confirmLong(t, f, id, previewLong(t, f, id))
	applied := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/apply", scopeBody(t, f, id))
	require.Equal(t, http.StatusOK, applied.Code, applied.Body.String())
	r.f.data = getGroupData(t, f, id)
	r.f.group = r.f.data["group"].(map[string]any)
	var err error
	r.snapshot, err = f.groups.ResolveAppliedAt(t.Context(), id, r.clock.Add(time.Minute))
	require.NoError(t, err)
	r.clock = r.snapshot.EffectiveAt.Add(13 * time.Second)
	r.store = groupdelivery.NewStore(f.db).WithClock(func() time.Time { return r.clock })
	r.pipe = grouppipeline.New(grouppipeline.Dependencies{Groups: f.groups, Store: r.store, Destinations: dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(f.db)), Counters: grouppipeline.NewMeasurementCounterSources(measurement.NewSQLRepository(f.db))}, grouppipeline.Config{Now: func() time.Time { return r.clock }, OnError: func(err error) { t.Errorf("pipeline: %v", err) }})
	require.NoError(t, r.pipe.Reconcile(t.Context()))
	require.Contains(t, r.pipe.Status(), grouppipeline.GroupStatus{GroupID: id, Revision: r.snapshot.AppliedRevision, State: grouppipeline.StateActive})
	_, err = f.db.Exec(`UPDATE devices SET description='' WHERE id='device-A'`)
	require.NoError(t, err)
	r.router = NewRouter(&DatalinkServices{Workspace: f.workspace, WriteGroups: f.groups, WriteGroupDelivery: r.pipe, Device: device.NewService(device.NewSQLRepository(f.db), nil)})
	return r
}

func (r *runtimeSourceDeliveryRig) offer(t *testing.T, source string, window int, bad bool) {
	t.Helper()
	var member workspace.WriteGroupMember
	for _, m := range r.snapshot.Group.Members {
		if m.TagID == "tag-"+source {
			member = m
		}
	}
	require.NotEmpty(t, member.TagID)
	at := r.snapshot.EffectiveAt.Add(time.Duration(window)*r.window + time.Second)
	quality := schema.QualityGood
	if bad {
		at = at.Add(time.Second)
		quality = schema.QualityBad
	}
	if r.clock.Before(at) {
		r.clock = at
	}
	var value any = int16(12)
	if r.snapshot.Group.RecordingProfile.Kind == workspace.WriteGroupProfileCumulative {
		value = uint64(100 + window)
	}
	require.NoError(t, r.pipe.AcceptSample(t.Context(), measurement.SampleEnvelope{SampleID: fmt.Sprintf("%s-%d-%t", source, window, bad), WorkspaceID: r.f.record.ID, DeviceID: member.DeviceID, PointID: member.PointID, TagID: member.TagID, SourceRevision: member.SourceRevision, MappingRevision: member.MappingRevision, ObservedAt: at, ReceivedAt: at, Value: value, Quality: quality, QualityReason: "password=private-driver-secret /private/raw-value"}))
}
func (r *runtimeSourceDeliveryRig) closeWindow(t *testing.T, window int) {
	t.Helper()
	r.clock = r.snapshot.EffectiveAt.Add(time.Duration(window+1) * r.window)
	r.pipe.TickAll(t.Context())
}
func (r *runtimeSourceDeliveryRig) raw(t *testing.T) *grouppipeline.DeliveryView {
	t.Helper()
	v, err := r.pipe.DeliveryForRevision(t.Context(), r.snapshot.GroupID, r.snapshot.AppliedRevision)
	require.NoError(t, err)
	return v
}

func (r *runtimeSourceDeliveryRig) deliver(t *testing.T) {
	t.Helper()
	rows, err := r.f.db.Query(`SELECT effect_key FROM wg_delivery_outbox WHERE state='pending' ORDER BY bucket_start,effect_key`)
	require.NoError(t, err)
	var keys []string
	for rows.Next() {
		var key string
		require.NoError(t, rows.Scan(&key))
		keys = append(keys, key)
	}
	require.NoError(t, rows.Err())
	require.NoError(t, rows.Close())
	require.NotEmpty(t, keys)
	for _, key := range keys {
		r.send(t, key)
	}
}
func runtimeDeliveredSource(t *testing.T, g runtimeview.GroupView, tagID string) runtimeview.SourceEvidence {
	t.Helper()
	require.NotNil(t, g.Delivery)
	require.NotNil(t, g.Delivery.Current)
	require.Equal(t, "available", g.Delivery.Current.SourceState)
	for _, s := range g.Delivery.Current.Sources {
		if s.TagID == tagID {
			return s
		}
	}
	t.Fatalf("source %s missing", tagID)
	return runtimeview.SourceEvidence{}
}
func (r *runtimeSourceDeliveryRig) get(t *testing.T) *httptest.ResponseRecorder {
	t.Helper()
	return performJSONRequest(t, r.router, http.MethodGet, runtimeRuleContextPath+"?"+r.query(t).Encode(), nil)
}

// The GET witness covers configuration, accepted journal, descriptors, outbox
// and receipts byte-for-byte. Setup/tick/sender writes occur before this read.
func (r *runtimeSourceDeliveryRig) durableState(t *testing.T) map[string]string {
	t.Helper()
	state := map[string]string{}
	for _, table := range []string{"system_settings", "write_groups", "write_group_versions", "wg_delivery_samples", "wg_delivery_outbox", "wg_delivery_receipts", "wg_runtime_versions"} {
		rows, err := r.f.db.Query(`SELECT * FROM ` + table + ` ORDER BY rowid`)
		require.NoError(t, err)
		columns, err := rows.Columns()
		require.NoError(t, err)
		var content [][]any
		for rows.Next() {
			values := make([]any, len(columns))
			dest := make([]any, len(columns))
			for i := range values {
				dest[i] = &values[i]
			}
			require.NoError(t, rows.Scan(dest...))
			content = append(content, values)
		}
		require.NoError(t, rows.Err())
		require.NoError(t, rows.Close())
		raw, err := json.Marshal(content)
		require.NoError(t, err)
		state[table] = fmt.Sprintf("%x", sha256.Sum256(raw))
	}
	return state
}

func (r *runtimeSourceDeliveryRig) useAsOf(asOf time.Time) {
	h := handlers.NewRuntimeHandler(device.NewService(device.NewSQLRepository(r.f.db), nil), nil, nil, nil, nil, r.f.workspace).
		WithRuleEvidenceServices(r.f.groups, r.pipe).WithRuleEvidenceClock(func() time.Time { return asOf }).WithRuleDeliveryServices(r.pipe)
	r.router = gin.New()
	r.router.GET(runtimeRuleContextPath, h.WorkspaceContext)
}
func newRuntimeSourceCustomRig(t *testing.T) *runtimeDeliveryRig {
	t.Helper()
	f := newManagedDestinationFixture(t, "sqlite")
	r := &runtimeDeliveryRig{runtimeRuleEvidenceRig: &runtimeRuleEvidenceRig{f: f, clock: time.Date(2026, 10, 2, 12, 0, 7, 0, time.UTC)}}
	f.groups.WithClock(func() time.Time { return r.clock })
	target, err := sql.Open("sqlite", f.file)
	require.NoError(t, err)
	_, err = target.Exec(`CREATE TABLE custom_values(temperature REAL NOT NULL)`)
	require.NoError(t, err)
	require.NoError(t, target.Close())
	request := f.createRequest()
	g := request["group"].(map[string]any)
	g["destination"].(map[string]any)["table_name"] = "custom_values"
	g["destination"].(map[string]any)["table_schema"] = "main"
	g["write_policy"] = map[string]any{"dedupe_capability": "none"}
	created := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, request)
	require.Equal(t, 201, created.Code, created.Body.String())
	r.f.data, r.f.group = groupSaveData(t, decodeJSONBody(t, created))
	id := r.f.group["id"].(string)
	applied := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/apply", scopeBody(t, f, id))
	require.Equal(t, 200, applied.Code, applied.Body.String())
	r.f.data = getGroupData(t, f, id)
	r.f.group = r.f.data["group"].(map[string]any)
	r.clock = r.clock.Add(9 * time.Second)
	targets := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(f.db))
	r.store = groupdelivery.NewStore(f.db).WithClock(func() time.Time { return r.clock })
	r.pipe = grouppipeline.New(grouppipeline.Dependencies{Groups: f.groups, Store: r.store, Destinations: targets, Inspector: dbtarget.NewReadOnlyTableInspector(targets), Tags: tag.NewService(tag.NewSQLRepository(f.db))}, grouppipeline.Config{Now: func() time.Time { return r.clock }})
	require.NoError(t, r.pipe.Reconcile(t.Context()))
	require.Contains(t, r.pipe.Status(), grouppipeline.GroupStatus{GroupID: id, Revision: r.f.group["applied_revision"].(string), State: grouppipeline.StateActive})
	_, err = f.db.Exec(`UPDATE devices SET description='' WHERE id='device-A'`)
	require.NoError(t, err)
	r.router = NewRouter(&DatalinkServices{Workspace: f.workspace, WriteGroups: f.groups, WriteGroupDelivery: r.pipe, Device: device.NewService(device.NewSQLRepository(f.db), nil)})
	return r
}
