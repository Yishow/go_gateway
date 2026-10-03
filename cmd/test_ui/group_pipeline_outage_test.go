package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"
	"go-gateway/internal/protocol/modbus"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

// outageEnv is a real gateway configuration database with two SQLite
// destinations and one applied write group per destination.
type outageEnv struct {
	db       *sql.DB
	services gatewayServices
	targets  map[string]string // connector ID -> destination file
	groups   map[string]*workspace.WriteGroup
	workspID string
	// schemas maps a destination name to the table schema groups use for it;
	// dedupe is the dedupe capability the groups declare.
	schemas map[string]string
	dedupe  string
}

// destinationSetup creates one destination and returns its connector kind, its
// connection config JSON and the table schema groups should use.
type destinationSetup func(t *testing.T, env *outageEnv, name string) (kind, configJSON, tableSchema string)

func newOutageEnv(t *testing.T) *outageEnv {
	t.Helper()
	return newOutageEnvWith(t, sqliteDestinations)
}

func sqliteDestinations(t *testing.T, env *outageEnv, name string) (kind, configJSON, tableSchema string) {
	t.Helper()
	file := filepath.Join(t.TempDir(), "destination-"+name+".db")
	destination, err := sql.Open("sqlite", file)
	require.NoError(t, err)
	_, err = destination.ExecContext(t.Context(), `CREATE TABLE readings (temperature REAL NOT NULL, pressure INTEGER NOT NULL)`)
	require.NoError(t, err)
	require.NoError(t, destination.Close())
	config, err := json.Marshal(map[string]string{"dsn": file})
	require.NoError(t, err)
	env.targets["connector-"+name] = file
	return "sqlite", string(config), "main"
}

func newOutageEnvWith(t *testing.T, setup destinationSetup) *outageEnv {
	t.Helper()
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "configuration.db")+"?_pragma=busy_timeout(5000)")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	datalink.ApplySQLitePoolDefaults(db) // the same single-connection pool production uses
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	exec := func(statement string, args ...any) {
		t.Helper()
		_, err := db.ExecContext(t.Context(), statement, args...)
		require.NoError(t, err)
	}
	exec(`INSERT INTO devices (id,name,protocol,status,connection_config,readiness_status) VALUES ('device-1','PLC','modbus_tcp','draft','{}','{}')`)
	env := &outageEnv{db: db, targets: map[string]string{}, groups: map[string]*workspace.WriteGroup{}, schemas: map[string]string{}}
	for _, name := range []string{"A", "B"} {
		exec(`INSERT INTO points (id,device_id,name,address,function,data_type,mode,enabled) VALUES (?, 'device-1', ?, ?, 'FC03', 'float64', 'read', 1)`, "point-t-"+name, "Temp "+name, "4000"+name)
		exec(`INSERT INTO points (id,device_id,name,address,function,data_type,mode,enabled) VALUES (?, 'device-1', ?, ?, 'FC03', 'int64', 'read', 1)`, "point-p-"+name, "Pressure "+name, "5000"+name)
		exec(`INSERT INTO tags (id,key,key_lower,display_name,data_type,status,labels) VALUES (?, ?, ?, 'Temperature', 'float64', 'active', '{}')`, "tag-t-"+name, "temperature."+name, "temperature."+name)
		exec(`INSERT INTO tags (id,key,key_lower,display_name,data_type,status,labels) VALUES (?, ?, ?, 'Pressure', 'int64', 'active', '{}')`, "tag-p-"+name, "pressure."+name, "pressure."+name)
		exec(`INSERT INTO mappings (id,point_id,tag_id,transform_pipeline,status,enabled) VALUES (?, ?, ?, '[]', 'active', 1)`, "mapping-t-"+name, "point-t-"+name, "tag-t-"+name)
		exec(`INSERT INTO mappings (id,point_id,tag_id,transform_pipeline,status,enabled) VALUES (?, ?, ?, '[]', 'active', 1)`, "mapping-p-"+name, "point-p-"+name, "tag-p-"+name)

		kind, configJSON, tableSchema := setup(t, env, name)
		env.schemas[name] = tableSchema
		exec(`INSERT INTO database_connectors (id,name,kind,connection_config,identity_revision,status,enabled) VALUES (?, ?, ?, ?, 'connector-1', 'ready', 1)`,
			"connector-"+name, "Target "+name, kind, configJSON)
	}
	env.services = wireGatewayServices(db, connector.NewConnectionManager(connector.ConnectionManagerConfig{}))
	record, err := env.services.workspace.AttachDevice(t.Context(), "device-1")
	require.NoError(t, err)
	env.workspID = record.ID
	return env
}

func (e *outageEnv) createAndApply(t *testing.T, name string) *workspace.WriteGroup {
	t.Helper()
	ctx := t.Context()
	record, err := e.services.workspace.GetOrCreate(ctx)
	require.NoError(t, err)
	created, err := e.services.writeGroups.Create(ctx, workspace.WriteGroupMutation{
		WorkspaceID: record.ID, ExpectedWorkspaceRevision: record.DatabaseSetupRevision, ExpectedConnectorRevision: "connector-1",
		Group: &workspace.WriteGroup{
			WorkspaceID: record.ID, Name: "Group " + name,
			Members: []workspace.WriteGroupMember{
				{DeviceID: "device-1", PointID: "point-t-" + name, TagID: "tag-t-" + name, TargetColumn: "temperature", Required: true},
				{DeviceID: "device-1", PointID: "point-p-" + name, TagID: "tag-p-" + name, TargetColumn: "pressure", Required: true},
			},
			Destination: workspace.WriteGroupDestination{ConnectorID: "connector-" + name, ConnectorRevision: "connector-1", TableSchema: e.schemas[name], TableName: "readings", StorageStrategy: workspace.WriteGroupStorageStrategyCustom},
			RowPolicy:   workspace.WriteGroupRowPolicy{IntervalSeconds: 10},
			WritePolicy: workspace.WriteGroupWritePolicy{DedupeCapability: e.dedupe},
		},
	})
	require.NoError(t, err)
	record, err = e.services.workspace.GetOrCreate(ctx)
	require.NoError(t, err)
	applied, err := e.services.writeGroups.Apply(ctx, created.Group.ID, workspace.WriteGroupMutation{
		WorkspaceID: record.ID, ExpectedWorkspaceRevision: record.DatabaseSetupRevision,
		ExpectedGroupRevision: created.Group.Revision, ExpectedConnectorRevision: "connector-1",
	})
	require.NoError(t, err)
	e.groups[name] = applied.Group
	return applied.Group
}

type testClock struct{ nanos atomic.Int64 }

func (c *testClock) now() time.Time  { return time.Unix(0, c.nanos.Load()).UTC() }
func (c *testClock) set(t time.Time) { c.nanos.Store(t.UnixNano()) }

func (e *outageEnv) pipeline(clock *testClock, owner string) *grouppipeline.Pipeline {
	immediate := func(int) (string, time.Time) { return groupdelivery.StateRetrying, clock.now().Add(-time.Hour) }
	return grouppipeline.New(grouppipeline.Dependencies{
		Groups: e.services.writeGroups, Tags: e.services.tag, Destinations: e.services.dbTarget,
		Inspector: dbtarget.NewReadOnlyTableInspector(e.services.dbTarget), Store: groupdelivery.NewStore(e.db),
	}, grouppipeline.Config{
		NodeID: "node-1", Owner: owner, Now: clock.now,
		TickInterval: time.Hour, ReconcileInterval: time.Hour, // the test drives ticks and reconciles itself
		WorkerInterval: 10 * time.Millisecond,
		Sender:         groupdelivery.SenderConfig{Backoff: immediate, DeliveryTimeout: 2 * time.Second},
		OnError:        func(error) {},
	})
}

func (e *outageEnv) envelope(group *workspace.WriteGroup, memberIndex int, id string, at time.Time, value any) measurement.SampleEnvelope {
	member := group.Members[memberIndex]
	return measurement.SampleEnvelope{
		SampleID: id, WorkspaceID: group.WorkspaceID, DeviceID: member.DeviceID, PointID: member.PointID, TagID: member.TagID,
		SourceRevision: member.SourceRevision, MappingRevision: member.MappingRevision,
		ObservedAt: at, ReceivedAt: at, TimeOrigin: "gateway", Value: value, Quality: schema.QualityGood,
	}
}

// feed gives every group one complete sample set observed inside bucket, then
// moves the clock past the bucket so it closes.
func (e *outageEnv) feed(t *testing.T, pipe *grouppipeline.Pipeline, clock *testClock, groups []*workspace.WriteGroup, bucket time.Time, temperature float64, pressure int64) {
	t.Helper()
	ctx := t.Context()
	observed := bucket.Add(5 * time.Second)
	clock.set(observed.Add(time.Second))
	for _, g := range groups {
		require.NoError(t, pipe.AcceptSample(ctx, e.envelope(g, 0, g.ID+"-t-"+bucket.Format("150405"), observed, temperature)))
		require.NoError(t, pipe.AcceptSample(ctx, e.envelope(g, 1, g.ID+"-p-"+bucket.Format("150405"), observed, pressure)))
	}
	clock.set(bucket.Add(11 * time.Second))
	pipe.TickAll(ctx)
}

func destinationRows(t *testing.T, file string) [][2]float64 {
	t.Helper()
	// Observe the real commit after a short lock clears; retain SQL errors if
	// the bounded wait expires instead of failing during normal contention.
	db, err := sql.Open("sqlite", file+"?_pragma=busy_timeout(15000)")
	require.NoError(t, err)
	defer db.Close()
	rows, err := db.QueryContext(t.Context(), `SELECT temperature, pressure FROM readings ORDER BY rowid`)
	require.NoError(t, err)
	defer rows.Close()
	var out [][2]float64
	for rows.Next() {
		var temperature, pressure float64
		require.NoError(t, rows.Scan(&temperature, &pressure))
		out = append(out, [2]float64{temperature, pressure})
	}
	require.NoError(t, rows.Err())
	return out
}

func TestProductionGroupOutageRecoveryOneWriterRestartAndIsolation(t *testing.T) {
	env := newOutageEnv(t)
	groupA := env.createAndApply(t, "A")
	groupB := env.createAndApply(t, "B")
	ctx := t.Context()

	clock := &testClock{}
	// Group effective boundaries are the next 10s boundary after Apply.
	snapshotAt := time.Now().UTC()
	for _, g := range []*workspace.WriteGroup{groupA, groupB} {
		resolved, err := env.services.writeGroups.ResolveAppliedAt(ctx, g.ID, snapshotAt.Add(time.Minute))
		require.NoError(t, err)
		if resolved.EffectiveAt.After(snapshotAt) {
			snapshotAt = resolved.EffectiveAt
		}
	}
	base := snapshotAt.Truncate(10 * time.Second).Add(10 * time.Second)
	clock.set(base.Add(time.Second))

	first := env.pipeline(clock, "node-1/first")
	require.NoError(t, first.Start(ctx))
	for _, status := range first.Status() {
		require.Equal(t, grouppipeline.StateActive, status.State, "group %s blocked: %s", status.GroupID, status.Reason)
	}
	require.True(t, first.Owns("connector-A", "tag-t-A"), "an applied group owns its tags on its destination: the legacy writer must stand down")
	require.False(t, first.Owns("connector-B", "tag-t-A"), "ownership is per output, not per tag")

	feed := func(pipe *grouppipeline.Pipeline, bucket time.Time, temperature float64, pressure int64) {
		t.Helper()
		env.feed(t, pipe, clock, []*workspace.WriteGroup{groupA, groupB}, bucket, temperature, pressure)
	}

	// Bucket 1: both destinations healthy.
	feed(first, base, 21.5, 100)
	require.Eventually(t, func() bool {
		return len(destinationRows(t, env.targets["connector-A"])) == 1 && len(destinationRows(t, env.targets["connector-B"])) == 1
	}, 10*time.Second, 20*time.Millisecond)
	require.Equal(t, [][2]float64{{21.5, 100}}, destinationRows(t, env.targets["connector-A"]))

	// Destination A goes offline (file removed); B and Share stay healthy.
	offline := env.targets["connector-A"] + ".offline"
	require.NoError(t, os.Rename(env.targets["connector-A"], offline))
	shareValue := publishToShare(t, 4321)
	feed(first, base.Add(10*time.Second), 22.5, 200)
	require.Eventually(t, func() bool { return len(destinationRows(t, env.targets["connector-B"])) == 2 }, 10*time.Second, 20*time.Millisecond)
	require.Equal(t, uint16(4321), shareValue(), "Modbus Share keeps publishing while a database destination is down")
	feed(first, base.Add(20*time.Second), 23.5, 300)
	require.Eventually(t, func() bool { return len(destinationRows(t, env.targets["connector-B"])) == 3 }, 10*time.Second, 20*time.Millisecond)

	var aBacklog int
	require.NoError(t, env.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM wg_delivery_outbox WHERE connector_id = 'connector-A' AND state != 'sql_committed'`).Scan(&aBacklog))
	require.Equal(t, 2, aBacklog, "A's two outage buckets are accepted and waiting, not lost and not reported as written")
	var aCommitted int
	require.NoError(t, env.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM wg_delivery_outbox WHERE connector_id = 'connector-A' AND state = 'sql_committed'`).Scan(&aCommitted))
	require.Equal(t, 1, aCommitted, "SQL committed appears only for the row the destination really holds")

	// The gateway restarts during the outage: the old incarnation stops and a
	// new one recovers from durable state alone. The offline file must not be
	// recreated as an empty database by the restart.
	require.NoError(t, first.Stop(2*time.Second))
	second := env.pipeline(clock, "node-1/second")
	require.NoError(t, second.Start(ctx))
	time.Sleep(300 * time.Millisecond)
	_, err := os.Stat(env.targets["connector-A"])
	require.ErrorIs(t, err, os.ErrNotExist, "an offline destination is never silently recreated")
	require.Len(t, destinationRows(t, env.targets["connector-B"]), 3, "B is unaffected by A's outage and by the restart")
	require.EqualValues(t, 2, outboxCount(t, env.db, "connector-A", "state != 'sql_committed'"))

	// Destination A comes back: its backlog is delivered once, in order, with
	// exact values, and nothing else is written twice.
	require.NoError(t, os.Rename(offline, env.targets["connector-A"]))
	require.Eventually(t, func() bool { return len(destinationRows(t, env.targets["connector-A"])) == 3 }, 15*time.Second, 25*time.Millisecond)
	require.Equal(t, [][2]float64{{21.5, 100}, {22.5, 200}, {23.5, 300}}, destinationRows(t, env.targets["connector-A"]))
	require.Equal(t, [][2]float64{{21.5, 100}, {22.5, 200}, {23.5, 300}}, destinationRows(t, env.targets["connector-B"]))
	require.Eventually(t, func() bool { return outboxCount(t, env.db, "connector-A", "state = 'sql_committed'") == 3 }, 5*time.Second, 25*time.Millisecond)
	require.Zero(t, outboxCount(t, env.db, "connector-B", "state != 'sql_committed'"))
	require.NoError(t, second.Stop(2*time.Second))
	require.Len(t, destinationRows(t, env.targets["connector-A"]), 3, "recovery never duplicates a row")
	require.NotNil(t, env.services.groupPipe, "production wiring owns a group pipeline")
}

func outboxCount(t *testing.T, db *sql.DB, connectorID, where string) int {
	t.Helper()
	var n int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM wg_delivery_outbox WHERE connector_id = ? AND `+where, connectorID).Scan(&n))
	return n
}

// publishToShare publishes a value through the production Share target writer
// and returns a reader of the Modbus register it landed in.
func publishToShare(t *testing.T, value int16) func() uint16 {
	t.Helper()
	ctx := context.Background()
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	created, err := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "outage.share", DataType: schema.DataTypeInt16})
	require.NoError(t, err)
	share := modbusshare.NewService(tagSvc, 4096)
	share.SetHydrationState(modbusshare.HydrationState{State: modbusshare.HydrationStateReady, Readiness: true})
	_, err = share.UpsertMapping(ctx, created.ID, 25)
	require.NoError(t, err)
	port := freeTCPPort(t)
	require.NoError(t, share.ApplySettings(ctx, modbusshare.Settings{Enabled: true, BindAddress: "127.0.0.1", Port: port, SlaveID: 1, CapacityRegisters: 2048}))
	t.Cleanup(func() { _ = share.CloseRuntime() })
	require.NoError(t, newProductionTargetWriter(nil, share).WriteTagValue(ctx, created.ID, value, time.Now()))
	client := modbus.CreateTCPClient("127.0.0.1", port, 1, time.Second)
	require.NoError(t, client.Connect())
	t.Cleanup(func() { _ = client.Close() })
	return func() uint16 {
		words, err := client.ReadHoldingRegisters(25, 1)
		require.NoError(t, err)
		require.Len(t, words, 1)
		return words[0]
	}
}

func TestRevisionBoundBacklogEndpointEditAndDisableNeverRetargetOldRows(t *testing.T) {
	env := newOutageEnv(t)
	group := env.createAndApply(t, "A")
	ctx := t.Context()
	resolved, err := env.services.writeGroups.ResolveAppliedAt(ctx, group.ID, time.Now().UTC().Add(time.Minute))
	require.NoError(t, err)
	base := resolved.EffectiveAt.Truncate(10 * time.Second).Add(10 * time.Second)
	clock := &testClock{}
	clock.set(base.Add(time.Second))

	pipe := env.pipeline(clock, "node-1/first")
	require.NoError(t, pipe.Start(ctx))
	defer func() { _ = pipe.Stop(2 * time.Second) }()

	// Destination A is offline when the first bucket closes.
	offline := env.targets["connector-A"] + ".offline"
	require.NoError(t, os.Rename(env.targets["connector-A"], offline))
	observed := base.Add(5 * time.Second)
	clock.set(observed.Add(time.Second))
	require.NoError(t, pipe.AcceptSample(ctx, env.envelope(group, 0, "t1", observed, 21.5)))
	require.NoError(t, pipe.AcceptSample(ctx, env.envelope(group, 1, "p1", observed, int64(100))))
	clock.set(base.Add(11 * time.Second))
	pipe.TickAll(ctx)
	require.Eventually(t, func() bool { return outboxCount(t, env.db, "connector-A", "state = 'retrying'") == 1 }, 10*time.Second, 25*time.Millisecond)

	// An operator points the connector at a different endpoint (new database file).
	elsewhere := filepath.Join(t.TempDir(), "elsewhere.db")
	other, err := sql.Open("sqlite", elsewhere)
	require.NoError(t, err)
	_, err = other.ExecContext(ctx, `CREATE TABLE readings (temperature REAL NOT NULL, pressure INTEGER NOT NULL)`)
	require.NoError(t, err)
	require.NoError(t, other.Close())
	newConfig := dbtarget.ConnectionConfig{"dsn": elsewhere}
	before, err := env.services.dbTarget.GetByID(ctx, "connector-A")
	require.NoError(t, err)
	updated, err := env.services.dbTarget.Update(ctx, "connector-A", dbtarget.UpdateConnectorRequest{
		ConnectionConfig: &newConfig, ExpectedIdentityRevision: &before.IdentityRevision,
	})
	require.NoError(t, err)
	require.NotEqual(t, before.IdentityRevision, updated.IdentityRevision, "editing the endpoint is a new destination identity")

	// The accepted row stays with the endpoint it was accepted for: it is blocked,
	// never delivered to the new endpoint, and the status says so truthfully.
	require.Eventually(t, func() bool { return outboxCount(t, env.db, "connector-A", "state = 'blocked'") == 1 }, 10*time.Second, 25*time.Millisecond)
	require.Empty(t, destinationRows(t, elsewhere), "old backlog is never sent to the edited endpoint")
	view, err := pipe.Delivery(ctx, group.ID)
	require.NoError(t, err)
	require.Equal(t, 1, view.Stages.Blocked)
	require.Zero(t, view.Stages.SQLCommitted)
	require.Nil(t, view.LastSQLCommittedAt, "blocked and buffered data is not written data")
	require.Len(t, view.Backlog, 1)
	require.Equal(t, before.IdentityRevision, view.Backlog[0].ConnectorRevision, "the backlog keeps the destination revision it was accepted for")
	require.Equal(t, []string{"target-blocked"}, view.Backlog[0].ErrorCodes)
	for _, entry := range view.Backlog {
		require.NotContains(t, entry.TableName+entry.ConnectorID, elsewhere)
	}

	// Disabling the group stops new intake but leaves the old backlog untouched.
	record, err := env.services.workspace.GetOrCreate(ctx)
	require.NoError(t, err)
	live, err := env.services.writeGroups.Get(ctx, group.ID)
	require.NoError(t, err)
	_, err = env.services.writeGroups.Disable(ctx, group.ID, workspace.WriteGroupMutation{
		WorkspaceID: record.ID, ExpectedWorkspaceRevision: record.DatabaseSetupRevision,
		ExpectedGroupRevision: live.Group.Revision, ExpectedConnectorRevision: before.IdentityRevision,
	})
	require.NoError(t, err, "an operator can stop intake even though the endpoint was edited")
	require.NoError(t, pipe.Reconcile(ctx))
	clock.set(base.Add(40 * time.Second))
	pipe.TickAll(ctx)
	require.Empty(t, destinationRows(t, elsewhere))
	view, err = pipe.Delivery(ctx, group.ID)
	require.NoError(t, err)
	require.Equal(t, 1, view.Stages.Blocked, "disabling never changes where accepted rows go")
	require.Equal(t, before.IdentityRevision, view.Backlog[0].ConnectorRevision)
	require.Equal(t, "not_running", view.Intake.State, "a disabled group is no longer accepting")
}

func TestProductionGroupOutageRecoveryDeleteWithBacklogIsCompletedByWorkers(t *testing.T) {
	env := newOutageEnv(t)
	group := env.createAndApply(t, "A")
	ctx := t.Context()
	resolved, err := env.services.writeGroups.ResolveAppliedAt(ctx, group.ID, time.Now().UTC().Add(time.Minute))
	require.NoError(t, err)
	base := resolved.EffectiveAt.Truncate(10 * time.Second).Add(10 * time.Second)
	clock := &testClock{}
	clock.set(base.Add(time.Second))
	pipe := env.pipeline(clock, "node-1/first")
	require.NoError(t, pipe.Start(ctx))
	defer func() { _ = pipe.Stop(2 * time.Second) }()

	// Destination offline: two accepted buckets wait in the durable outbox.
	offline := env.targets["connector-A"] + ".offline"
	require.NoError(t, os.Rename(env.targets["connector-A"], offline))
	env.feed(t, pipe, clock, []*workspace.WriteGroup{group}, base, 21.5, 100)
	env.feed(t, pipe, clock, []*workspace.WriteGroup{group}, base.Add(10*time.Second), 22.5, 200)
	require.Eventually(t, func() bool { return outboxCount(t, env.db, "connector-A", "state != 'sql_committed'") == 2 }, 10*time.Second, 25*time.Millisecond)

	// The operator deletes the applied group while it still has accepted data.
	record, err := env.services.workspace.GetOrCreate(ctx)
	require.NoError(t, err)
	live, err := env.services.writeGroups.Get(ctx, group.ID)
	require.NoError(t, err)
	deleted, err := env.services.writeGroups.Delete(ctx, group.ID, workspace.WriteGroupMutation{
		WorkspaceID: record.ID, ExpectedWorkspaceRevision: record.DatabaseSetupRevision,
		ExpectedGroupRevision: live.Group.Revision, ExpectedConnectorRevision: live.Group.Destination.ConnectorRevision,
	})
	require.NoError(t, err, "a group whose backlog still names its immutable revision and frozen destination can be deleted")
	require.Equal(t, workspace.WriteGroupStatusDeleted, deleted.Group.Status)
	require.NoError(t, pipe.Reconcile(ctx))

	// New intake stops; the tombstone keeps the backlog queryable and its owner proven.
	require.Equal(t, 2, outboxCount(t, env.db, "connector-A", "state != 'sql_committed'"))
	view, err := pipe.Delivery(ctx, group.ID)
	require.NoError(t, err)
	require.Len(t, view.Backlog, 1)
	require.Equal(t, 2, view.Backlog[0].Pending)

	// The destination returns: workers deliver the original accepted rows from
	// their frozen revisions, once and in order, with nothing orphaned or replayed.
	require.NoError(t, os.Rename(offline, env.targets["connector-A"]))
	require.Eventually(t, func() bool { return len(destinationRows(t, env.targets["connector-A"])) == 2 }, 15*time.Second, 25*time.Millisecond)
	require.Equal(t, [][2]float64{{21.5, 100}, {22.5, 200}}, destinationRows(t, env.targets["connector-A"]))
	require.Eventually(t, func() bool { return outboxCount(t, env.db, "connector-A", "state = 'sql_committed'") == 2 }, 5*time.Second, 25*time.Millisecond)
	view, err = pipe.Delivery(ctx, group.ID)
	require.NoError(t, err)
	require.Equal(t, 2, view.Stages.SQLCommitted)
	require.Empty(t, view.Backlog)
}
