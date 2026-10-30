package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/diagnostics"
	"go-gateway/internal/testutil/lostcommit"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// The fixture uses public canonical save/schema/apply APIs and two disposable
// SQLite databases. It never seeds a succeeded operation or frozen revision.
type operationEventFixture struct {
	deliveryEventFixture
	target  *workspace.WriteGroup
	request workspace.CutoverRequest
	clock   *testClock
}

func newOperationEventFixture(t *testing.T, fileLog ...bool) operationEventFixture {
	t.Helper()
	f := newOperationDeliveryFixture(t, fileLog...)
	clock := &testClock{}
	clock.set(time.Now().UTC())
	f.env.services.writeGroups.WithClock(clock.now)
	target := operationManagedGroup(t, f, "A")
	readiness, err := f.env.services.writeGroups.Readiness(t.Context(), target.ID)
	require.NoError(t, err)
	require.True(t, readiness.Ready, "actual managed target readiness: %+v", readiness)
	current, err := f.env.services.workspace.ReadSnapshot(t.Context())
	require.NoError(t, err)
	req := workspace.CutoverRequest{Kind: workspace.WriteGroupCutoverForward,
		WorkspaceID: f.group.WorkspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		SourceGroupID: f.group.ID, ExpectedSourceRevision: f.group.Revision,
		ExpectedConnectorRevision: f.group.Destination.ConnectorRevision,
		Targets:                   []workspace.CutoverTargetRef{{GroupID: target.ID, ExpectedRevision: target.Revision}},
		Boundary:                  clock.now().Add(2 * time.Minute),
	}
	review, err := f.env.services.writeGroups.ReviewCutover(t.Context(), req)
	require.NoError(t, err)
	require.True(t, review.Applicable, "actual canonical cutover must be applicable: %v", review.Issues)
	req.Boundary, req.ReviewDigest = review.Boundary, review.Digest
	return operationEventFixture{f, target, req, clock}
}

func operationManagedGroup(t *testing.T, f deliveryEventFixture, name string) *workspace.WriteGroup {
	t.Helper()
	record, err := f.env.services.workspace.ReadSnapshot(t.Context())
	require.NoError(t, err)
	yes := true
	draft := &workspace.WriteGroup{WorkspaceID: record.ID, Name: "password_fixture postgres://private-host/config",
		LayoutVersion: workspace.WriteGroupLayoutManagedLongV1, StorageMode: workspace.WriteGroupStorageModePlain,
		RecordingProfile: &workspace.WriteGroupRecordingProfile{Kind: workspace.WriteGroupProfileInstantaneous, Interval: 10, Unit: workspace.WriteGroupIntervalUnitSecond},
		Members: []workspace.WriteGroupMember{
			{DeviceID: "device-1", PointID: "point-t-" + name, TagID: "tag-t-" + name, ValueRole: workspace.WriteGroupValueRoleGauge},
			{DeviceID: "device-1", PointID: "point-p-" + name, TagID: "tag-p-" + name, ValueRole: workspace.WriteGroupValueRoleGauge},
		},
		Destination: workspace.WriteGroupDestination{ConnectorID: "connector-" + name, ConnectorRevision: "connector-1", TableSchema: "main", StorageStrategy: workspace.WriteGroupStorageStrategyManaged},
		RowPolicy:   workspace.WriteGroupRowPolicy{IntervalSeconds: 10}, WritePolicy: workspace.WriteGroupWritePolicy{RecreateMissing: &yes},
	}
	created, err := f.env.services.writeGroups.Create(t.Context(), workspace.WriteGroupMutation{WorkspaceID: record.ID, ExpectedWorkspaceRevision: record.DatabaseSetupRevision, ExpectedConnectorRevision: "connector-1", Group: draft})
	require.NoError(t, err)
	mutation := workspace.WriteGroupMutation{WorkspaceID: record.ID, ExpectedWorkspaceRevision: created.WorkspaceRevision, ExpectedGroupRevision: created.Group.Revision, ExpectedConnectorRevision: "connector-1"}
	scope, err := f.env.services.writeGroups.ManagedSchemaScope(t.Context(), created.Group.ID, mutation)
	require.NoError(t, err)
	target, err := sql.Open("sqlite", f.env.targets["connector-"+name])
	require.NoError(t, err)
	defer target.Close()
	// The borrowed delivery fixture creates the legacy receipt layout. It has
	// never delivered here; replace only that proven-empty disposable table so
	// the native managed preview owns creation of its required marker/layout.
	var receiptRows int
	require.NoError(t, target.QueryRow(`SELECT COUNT(*) FROM gw_effect_receipts`).Scan(&receiptRows))
	require.Zero(t, receiptRows)
	_, err = target.Exec(`DROP TABLE gw_effect_receipts`)
	require.NoError(t, err)
	inspect := operationTableInspector(target)
	token, err := f.env.services.recordingPlan.PrepareSchemaPreview(t.Context(), scope, inspect)
	require.NoError(t, err)
	op, _, err := f.env.services.recordingPlan.ApplySchemaPreview(t.Context(), token, recordingplan.SchemaApplyTarget{Inspect: inspect, Execute: func(ctx context.Context, statements []string) (recordingplan.SchemaExecution, error) {
		result, err := dbtarget.RunManagedSchemaBatch(ctx, target, schema.DatabaseConnectorKindSQLite, "main", statements)
		if result == nil {
			return recordingplan.SchemaExecution{}, err
		}
		return recordingplan.SchemaExecution{Committed: result.Committed, RolledBack: result.RolledBack}, err
	}})
	require.NoError(t, err)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, op.Status)
	confirmed, err := f.env.services.writeGroups.ConfirmManagedSchemaWithStorage(t.Context(), created.Group.ID, mutation, op.OperationID, op.VerifiedDigest, &workspace.StorageEvidence{Layout: *token.GroupLayout, CreatedTable: true})
	require.NoError(t, err)
	return confirmed.Group
}

func operationTableInspector(db *sql.DB) recordingplan.TargetInspector {
	return func(ctx context.Context, table string) (recordingplan.TargetTableInspection, error) {
		info, err := dbtarget.InspectOpenManagedTable(ctx, db, schema.DatabaseConnectorKindSQLite, "main", table, table != recordingplan.ManagedEffectReceiptTable)
		if err != nil {
			return recordingplan.TargetTableInspection{}, err
		}
		state := recordingplan.TargetTableInspection{Status: string(info.Status), ColumnTypes: map[string]string{}, ColumnNullable: map[string]bool{}, ColumnPrimaryKey: map[string]bool{}, IndexesInspected: info.IndexesInspected}
		for _, column := range info.Columns {
			state.Columns = append(state.Columns, column.Name)
			state.ColumnTypes[column.Name] = column.DataType
			state.ColumnNullable[column.Name] = column.Nullable
			state.ColumnPrimaryKey[column.Name] = column.PrimaryKey
		}
		for _, index := range info.Indexes {
			state.Indexes = append(state.Indexes, recordingplan.TargetIndexInspection{Name: index.Name, Columns: index.Columns, Unique: index.Unique})
		}
		return state, nil
	}
}

type repairOperationFixture struct {
	deliveryEventFixture
	long     *workspace.WriteGroup
	clock    *testClock
	seed     *grouppipeline.Pipeline
	base     time.Time
	targetDB *sql.DB
}

func newRepairOperationFixture(t *testing.T) repairOperationFixture {
	t.Helper()
	f := newOperationDeliveryFixture(t)
	clock := &testClock{}
	base := time.Now().UTC().Truncate(10 * time.Second).Add(-40 * time.Second)
	clock.set(base.Add(-time.Second))
	f.env.services.writeGroups.WithClock(clock.now)
	group := operationManagedGroup(t, f, "B")
	current, err := f.env.services.workspace.ReadSnapshot(t.Context())
	require.NoError(t, err)
	result, err := f.env.services.writeGroups.Apply(t.Context(), group.ID, workspace.WriteGroupMutation{
		WorkspaceID: current.ID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedGroupRevision: group.Revision, ExpectedConnectorRevision: group.Destination.ConnectorRevision,
	})
	require.NoError(t, err)
	group = result.Group
	clock.set(base.Add(time.Second))
	seed := f.env.pipeline(clock, "operation-event-seed")
	require.NoError(t, seed.Reconcile(t.Context()))
	f.env.feed(t, seed, clock, []*workspace.WriteGroup{group}, base, 21.5, 100)
	target, err := sql.Open("sqlite", f.env.targets["connector-B"])
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, target.Close()) })
	return repairOperationFixture{f, group, clock, seed, base, target}
}

func (f repairOperationFixture) waitLongCommitted(t *testing.T, n int) {
	t.Helper()
	require.Eventually(t, func() bool {
		var count int
		require.NoError(t, f.env.db.QueryRow(`SELECT COUNT(*) FROM wg_delivery_outbox WHERE group_id=? AND state=?`, f.long.ID, groupdelivery.StateCommitted).Scan(&count))
		return count >= n
	}, 8*time.Second, 20*time.Millisecond, "actual managed SQLite outbox must commit")
}
func (f repairOperationFixture) pending(t *testing.T) groupdelivery.OutboxItem {
	t.Helper()
	var effect string
	require.NoError(t, f.env.db.QueryRow(`SELECT effect_key FROM wg_delivery_outbox WHERE group_id=? AND state=? ORDER BY effect_key LIMIT 1`, f.long.ID, groupdelivery.StatePending).Scan(&effect))
	item, err := groupdelivery.NewStore(f.env.db).GetOutbox(t.Context(), effect)
	require.NoError(t, err)
	return item
}
func (f repairOperationFixture) dropTable(t *testing.T) {
	t.Helper()
	_, err := f.targetDB.Exec(`DROP TABLE "` + f.long.Destination.TableName + `"`)
	require.NoError(t, err, "only the disposable, canonical fixture table is removed")
}

// Only the negative seam replaces a driver. The native restore, ledger and
// lineage still execute against the same real SQLite fixture bytes.
func (f repairOperationFixture) restoreWithFault(t *testing.T, cfg grouppipeline.StorageObservationConfig, mode string) (groupdelivery.RestoreDecision, error) {
	t.Helper()
	item := f.pending(t)
	f.dropTable(t)
	local, target := f.env.db, f.targetDB
	faults := &lostcommit.Faults{}
	reject := func(q string) error {
		if mode == "rollback" && strings.HasPrefix(strings.TrimSpace(q), "CREATE TABLE") {
			return errors.New("password_fixture private-host SQL denial")
		}
		if mode == "unacknowledged" && strings.Contains(q, "UPDATE managed_schema_operations") && strings.Contains(q, "SET status") {
			return errors.New("password_fixture local settlement failure")
		}
		return nil
	}
	faults.Reject.Store(&reject)
	if mode == "rollback" {
		target = lostcommit.OpenFaults(f.env.targets["connector-B"], faults)
		defer target.Close()
	}
	if mode == "unacknowledged" {
		var seq int
		var name, path string
		require.NoError(t, f.env.db.QueryRow("PRAGMA database_list").Scan(&seq, &name, &path))
		local = lostcommit.OpenFaults(path, faults)
		defer local.Close()
	}
	recovery := grouppipeline.NewStorageRecoveryWithObservation(groupdelivery.NewStore(local), recordingplan.NewService(recordingplan.NewSQLRepository(local)), f.env.services.dbTarget, cfg)
	return recovery.Restore(t.Context(), item, groupdelivery.Target{DB: target, Kind: schema.DatabaseConnectorKindSQLite}, errors.New("no such table: private-host"))
}

// Unlike the custom-delivery lock fixture, managed schema inspection requires
// a guarded plain file DSN (not a user-supplied _pragma URI). Resolve that at
// initial fixture construction, before any frozen connector/group identity.
func newOperationDeliveryFixture(t *testing.T, fileLog ...bool) deliveryEventFixture {
	t.Helper()
	var destination *sql.DB
	setup := func(t *testing.T, env *outageEnv, name string) (string, string, string) {
		path := filepath.Join(t.TempDir(), "operation-"+name+".db")
		db, err := sql.Open("sqlite", path)
		require.NoError(t, err)
		t.Cleanup(func() { require.NoError(t, db.Close()) })
		_, err = db.Exec(`CREATE TABLE readings (temperature REAL NOT NULL, pressure INTEGER NOT NULL)`)
		require.NoError(t, err)
		require.NoError(t, dbtarget.CreateEffectReceiptTable(t.Context(), db, schema.DatabaseConnectorKindSQLite, "main"))
		env.targets["connector-"+name] = path
		if name == "A" {
			destination = db
		}
		raw, err := json.Marshal(map[string]string{"dsn": path})
		require.NoError(t, err)
		return "sqlite", string(raw), "main"
	}
	env := newOutageEnvWith(t, setup)
	require.NoError(t, env.services.runtime.Stop(context.Background()))
	options, logPath := diagnostics.Options{}, ""
	if len(fileLog) > 0 && fileLog[0] {
		var sequence int
		var name, path string
		require.NoError(t, env.db.QueryRow("PRAGMA database_list").Scan(&sequence, &name, &path))
		logPath = filepath.Join(filepath.Dir(path), "logs", "operations.jsonl")
		var err error
		options.File, err = diagnostics.OpenFileSink(diagnostics.FileOptions{Path: logPath, DatabasePath: path, DatabaseID: strings.Repeat("a", 64)})
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
	seed := env.pipeline(clock, "operation-source-seed")
	require.NoError(t, seed.Reconcile(t.Context()))
	env.feed(t, seed, clock, []*workspace.WriteGroup{group}, base, 21.5, 100)
	var effect string
	require.NoError(t, env.db.QueryRow(`SELECT effect_key FROM wg_delivery_outbox WHERE group_id=?`, group.ID).Scan(&effect))
	item, err := groupdelivery.NewStore(env.db).GetOutbox(t.Context(), effect)
	require.NoError(t, err)
	require.Equal(t, groupdelivery.StatePending, item.State)
	record, err := env.services.workspace.ReadSnapshot(t.Context())
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, env.services.groupPipe.Stop(3*time.Second))
		require.NoError(t, env.services.runtime.Stop(context.Background()))
	})
	return deliveryEventFixture{env, broker, destination, group, record, item, logPath}
}

func (f repairOperationFixture) pendingForBucket(t *testing.T, bucket time.Time) []groupdelivery.OutboxItem {
	t.Helper()
	rows, err := f.env.db.Query(`SELECT effect_key FROM wg_delivery_outbox WHERE group_id=?`, f.long.ID)
	require.NoError(t, err)
	var effects []string
	for rows.Next() {
		var effect string
		require.NoError(t, rows.Scan(&effect))
		effects = append(effects, effect)
	}
	require.NoError(t, rows.Close())
	var items []groupdelivery.OutboxItem
	store := groupdelivery.NewStore(f.env.db)
	for _, effect := range effects {
		item, err := store.GetOutbox(t.Context(), effect)
		require.NoError(t, err)
		if item.BucketStart.Equal(bucket) {
			require.Equal(t, groupdelivery.StatePending, item.State)
			items = append(items, item)
		}
	}
	require.Len(t, items, 2, "the new exact native window owns both pending source records")
	return items
}
