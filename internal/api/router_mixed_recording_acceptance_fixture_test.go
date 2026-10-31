package api

import (
	"context"
	"database/sql"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"
)

type mixedAcceptanceRig struct {
	f      managedGroupFixture
	clock  time.Time
	frozen *workspace.WriteGroupAppliedSnapshot
	store  *groupdelivery.Store
	pipe   *grouppipeline.Pipeline
}

// Only acquisition observations and clock are simulated. Configuration, DDL,
// Apply, snapshot recovery, accepted samples, closure, sender and SQL are real.
func applyMixedAcceptance(t *testing.T, f managedGroupFixture, request map[string]any, at time.Time) *mixedAcceptanceRig {
	t.Helper()
	r := &mixedAcceptanceRig{f: f, clock: at}
	f.groups.WithClock(func() time.Time { return r.clock })
	created := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, request)
	require.Equal(t, http.StatusCreated, created.Code, created.Body.String())
	r.f.data, r.f.group = groupSaveData(t, decodeJSONBody(t, created))
	id := r.f.group["id"].(string)
	confirmLong(t, r.f, id, previewLong(t, r.f, id))
	applied := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/apply", scopeBody(t, r.f, id))
	require.Equal(t, http.StatusOK, applied.Code, applied.Body.String())
	var err error
	r.frozen, err = f.groups.ResolveAppliedAt(t.Context(), id, at)
	require.NoError(t, err)
	r.store = newMixedAcceptanceStore(t, f.db, func() time.Time { return r.clock })
	r.reconcile(t)
	r.f.target, err = sql.Open("sqlite", f.file)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, r.f.target.Close()) })
	return r
}
func (r *mixedAcceptanceRig) reconcile(t *testing.T) {
	t.Helper()
	targets := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(r.f.db))
	r.pipe = grouppipeline.New(grouppipeline.Dependencies{Groups: r.f.groups, Store: r.store, Destinations: targets, Tags: tag.NewService(tag.NewSQLRepository(r.f.db)), Inspector: dbtarget.NewReadOnlyTableInspector(targets), Counters: grouppipeline.NewMeasurementCounterSources(measurement.NewSQLRepository(r.f.db))}, grouppipeline.Config{Now: func() time.Time { return r.clock }, OnError: func(err error) { t.Errorf("production mixed pipeline: %v", err) }})
	require.NoError(t, r.pipe.Reconcile(t.Context()))
	require.Contains(t, r.pipe.Status(), grouppipeline.GroupStatus{GroupID: r.frozen.GroupID, Revision: r.frozen.AppliedRevision, State: grouppipeline.StateActive})
}
func (r *mixedAcceptanceRig) offer(t *testing.T, index int, id string, at time.Time, value any) {
	t.Helper()
	m := r.frozen.Group.Members[index]
	if r.clock.Before(at) {
		r.clock = at
	}
	require.NoError(t, r.pipe.AcceptSample(t.Context(), measurement.SampleEnvelope{SampleID: id, WorkspaceID: r.f.record.ID, DeviceID: m.DeviceID, PointID: m.PointID, TagID: m.TagID, SourceRevision: m.SourceRevision, MappingRevision: m.MappingRevision, ObservedAt: at, ReceivedAt: r.clock, Value: value, Quality: schema.QualityGood}))
}
func (r *mixedAcceptanceRig) tick(t *testing.T, at time.Time) {
	t.Helper()
	r.clock = at
	r.pipe.TickAll(t.Context())
}
func (r *mixedAcceptanceRig) Resolve(ctx context.Context, item groupdelivery.OutboxItem) (groupdelivery.Target, error) {
	opened, err := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(r.f.db)).OpenDestination(ctx, item.Destination.ConnectorID, item.Destination.ConnectorRevision)
	if err != nil {
		return groupdelivery.Target{}, err
	}
	return groupdelivery.Target{DB: opened.DB, Kind: opened.Kind, Close: opened.Close}, nil
}
func (r *mixedAcceptanceRig) deliver(t *testing.T) []string {
	t.Helper()
	rows, err := r.f.db.Query(`SELECT effect_key FROM wg_delivery_outbox WHERE group_id=? AND state='pending' ORDER BY due_at,effect_key`, r.frozen.GroupID)
	require.NoError(t, err)
	var keys []string
	for rows.Next() {
		var key string
		require.NoError(t, rows.Scan(&key))
		keys = append(keys, key)
	}
	require.NoError(t, rows.Err())
	require.NoError(t, rows.Close())
	sender := groupdelivery.NewSender(r.store, r, groupdelivery.SenderConfig{Owner: "mixed-api", Now: func() time.Time { return r.clock }})
	for _, key := range keys {
		result, err := sender.Deliver(t.Context(), key)
		require.NoError(t, err)
		require.Equal(t, groupdelivery.StateCommitted, result.State)
		var local, target int
		require.NoError(t, r.f.db.QueryRow(`SELECT count(*) FROM wg_delivery_receipts WHERE effect_key=?`, key).Scan(&local))
		require.NoError(t, r.f.target.QueryRow(`SELECT count(*) FROM "`+dbtarget.EffectReceiptTable+`" WHERE effect_key=?`, key).Scan(&target))
		require.Equal(t, 1, local)
		require.Equal(t, 1, target)
	}
	return keys
}

// Match production canonical-intake and ordinary default capacity settings;
// these tiny fixtures do not claim quota-pressure or acquisition-loss testing.
func newMixedAcceptanceStore(t *testing.T, db *sql.DB, now func() time.Time) *groupdelivery.Store {
	t.Helper()
	store, err := groupdelivery.NewStore(db).WithCanonicalIntake().WithClock(now).WithQuota(groupdelivery.QuotaConfig{GlobalMaxBytes: 500 * 1024 * 1024})
	require.NoError(t, err)
	return store
}
