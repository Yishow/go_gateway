package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/workspace"
	"go-gateway/internal/diagnostics"

	"github.com/stretchr/testify/require"
)

type deliveryEventFixture struct {
	env         *outageEnv
	broker      *diagnostics.Broker
	destination *sql.DB
	group       *workspace.WriteGroup
	record      *workspace.Record
	item        groupdelivery.OutboxItem
	logPath     string
}

func newDeliveryEventFixture(t *testing.T, fileLog ...bool) deliveryEventFixture {
	t.Helper()
	var destination *sql.DB
	setup := func(t *testing.T, env *outageEnv, name string) (string, string, string) {
		t.Helper()
		file := filepath.Join(t.TempDir(), "delivery-"+name+".db")
		db, err := sql.Open("sqlite", file)
		require.NoError(t, err)
		t.Cleanup(func() { require.NoError(t, db.Close()) })
		_, err = db.ExecContext(t.Context(), `CREATE TABLE readings (temperature REAL NOT NULL, pressure INTEGER NOT NULL)`)
		require.NoError(t, err)
		require.NoError(t, dbtarget.CreateEffectReceiptTable(t.Context(), db, schema.DatabaseConnectorKindSQLite, "main"))
		env.targets["connector-"+name] = file
		if name == "A" {
			destination = db
		}
		config, err := json.Marshal(map[string]string{"dsn": file + "?_pragma=busy_timeout(0)"})
		require.NoError(t, err)
		return "sqlite", string(config), "main"
	}
	env := newOutageEnvWith(t, setup)
	// Replace only the test's original nil-broker composition. The production
	// composition below owns its Store/dispatcher, not a manually configured clone.
	require.NoError(t, env.services.runtime.Stop(context.Background()))
	var err error
	options, logPath := diagnostics.Options{}, ""
	if len(fileLog) > 0 && fileLog[0] {
		var sequence int
		var name, dbPath string
		require.NoError(t, env.db.QueryRowContext(t.Context(), "PRAGMA database_list").Scan(&sequence, &name, &dbPath))
		logPath = filepath.Join(filepath.Dir(dbPath), "logs", "delivery.jsonl")
		options.File, err = diagnostics.OpenFileSink(diagnostics.FileOptions{Path: logPath, DatabasePath: dbPath, DatabaseID: strings.Repeat("a", 64)})
		require.NoError(t, err)
	}
	broker, err := diagnostics.New(options)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, broker.Close(context.Background())) })
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	t.Cleanup(func() { require.NoError(t, manager.CloseAll()) })
	env.services, err = wireGatewayServicesContext(t.Context(), env.db, manager, broker)
	require.NoError(t, err)
	env.dedupe = "receipt"
	clock := &testClock{}
	base := time.Now().UTC().Truncate(10 * time.Second).Add(-time.Minute)
	clock.set(base.Add(-time.Second))
	env.services.writeGroups.WithClock(clock.now)
	group := env.createAndApply(t, "A")
	clock.set(base.Add(time.Second))
	seed := env.pipeline(clock, "event-seed")
	require.NoError(t, seed.Reconcile(t.Context()))
	env.feed(t, seed, clock, []*workspace.WriteGroup{group}, base, 21.5, 100)
	var effect string
	require.NoError(t, env.db.QueryRowContext(t.Context(), `SELECT effect_key FROM wg_delivery_outbox WHERE group_id=?`, group.ID).Scan(&effect))
	item, err := groupdelivery.NewStore(env.db).GetOutbox(t.Context(), effect)
	require.NoError(t, err)
	require.Equal(t, groupdelivery.StatePending, item.State, "fixture must produce a real durable queued output")
	record, err := env.services.workspace.ReadSnapshot(t.Context())
	require.NoError(t, err)
	require.NotEmpty(t, record.DatabaseSetupRevision)
	t.Cleanup(func() {
		require.NoError(t, env.services.groupPipe.Stop(3*time.Second))
		require.NoError(t, env.services.runtime.Stop(context.Background()))
	})
	return deliveryEventFixture{env, broker, destination, group, record, item, logPath}
}

func (f deliveryEventFixture) start(t *testing.T) {
	t.Helper()
	require.NoError(t, f.env.services.groupPipe.Start(t.Context()))
}

func (f deliveryEventFixture) state(t *testing.T) groupdelivery.OutboxItem {
	t.Helper()
	item, err := groupdelivery.NewStore(f.env.db).GetOutbox(t.Context(), f.item.EffectKey)
	require.NoError(t, err)
	return item
}

func (f deliveryEventFixture) waitCommitted(t *testing.T) groupdelivery.OutboxItem {
	t.Helper()
	require.Eventually(t, func() bool { return f.state(t).State == groupdelivery.StateCommitted }, 8*time.Second, 20*time.Millisecond, "native durable delivery must commit independently of diagnostics")
	item := f.state(t)
	var matched int
	require.NoError(t, f.env.db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM wg_delivery_receipts WHERE effect_key=? AND payload_digest=?`, item.EffectKey, item.PayloadDigest).Scan(&matched))
	require.Equal(t, 1, matched)
	var rows int
	require.NoError(t, f.destination.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM readings`).Scan(&rows))
	require.Equal(t, 1, rows, "actual destination row, not just local success status")
	return item
}

func deliveryEventsOfCode(t *testing.T, broker *diagnostics.Broker, code string) []diagnostics.Event {
	t.Helper()
	var events []diagnostics.Event
	for _, e := range collectorEventRecords(t, broker) {
		if e.Code == code {
			events = append(events, e)
		}
	}
	return events
}

// The SQL state becomes visible just before the post-commit diagnostic callback.
func waitDeliveryEvents(t *testing.T, broker *diagnostics.Broker, code string) []diagnostics.Event {
	t.Helper()
	var events []diagnostics.Event
	require.Eventually(t, func() bool { events = deliveryEventsOfCode(t, broker, code); return len(events) > 0 }, time.Second, 10*time.Millisecond)
	return events
}

// Journal real accepted samples before corrupting only frozen proof metadata.
// These two historical scopes let one native Reconcile exercise a reduced
// test-only observation cap, without creating an unrelated writer or device.
func seedDeliveryOwnerProofFaults(t *testing.T, f deliveryEventFixture) []groupdelivery.GroupKey {
	t.Helper()
	store := groupdelivery.NewStore(f.env.db)
	raw, err := store.RuntimeVersion(t.Context(), f.item.Key)
	require.NoError(t, err)
	keys := make([]groupdelivery.GroupKey, 0, 2)
	for i := range 2 {
		var frozen workspace.WriteGroupAppliedSnapshot
		require.NoError(t, json.Unmarshal(raw, &frozen))
		revision := fmt.Sprintf("historical-%d", i)
		frozen.AppliedRevision = revision
		frozen.Group.Revision = revision
		frozen.Group.AppliedRevision = revision
		key := groupdelivery.GroupKey{WorkspaceID: frozen.WorkspaceID, GroupID: frozen.GroupID, GroupRevision: revision}
		clean, err := json.Marshal(frozen)
		require.NoError(t, err)
		require.NoError(t, store.SaveRuntimeVersion(t.Context(), key, clean))
		member := frozen.Group.Members[0]
		sample := snapshot.Sample{SampleID: fmt.Sprintf("accepted-%d", i), MemberKey: member.DeviceID + "/" + member.PointID + "/" + member.TagID, ObservedAt: frozen.EffectiveAt.Add(time.Second), Quality: schema.QualityGood, Value: measurement.NewInt64(1), SourceRevision: member.SourceRevision, MappingRevision: member.MappingRevision}
		require.NoError(t, store.AppendSample(t.Context(), key, frozen.EffectiveAt, sample))
		var fields map[string]json.RawMessage
		require.NoError(t, json.Unmarshal(clean, &fields))
		fields["counter_bindings"] = json.RawMessage(`{"corrupted-proof":{"version":77}}`)
		damaged, err := json.Marshal(fields)
		require.NoError(t, err)
		_, err = f.env.db.ExecContext(t.Context(), "UPDATE wg_runtime_versions SET payload=? WHERE workspace_id=? AND group_id=? AND group_revision=?", string(damaged), key.WorkspaceID, key.GroupID, key.GroupRevision)
		require.NoError(t, err)
		keys = append(keys, key)
	}
	return keys
}
