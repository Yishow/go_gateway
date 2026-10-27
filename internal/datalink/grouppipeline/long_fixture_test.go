package grouppipeline

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
	"go-gateway/internal/testutil/lostcommit"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

const (
	longWindow    = 10 * time.Minute
	longKwhFirst  = "tag-kwh-1"
	longKwhSecond = "tag-kwh-2"
)

// liveDestinations resolves the saved connector to a real destination database.
// While offline is set it behaves like an unreachable endpoint.
type liveDestinations struct {
	connector *schema.DatabaseConnector
	target    *sql.DB
	kind      schema.DatabaseConnectorKind
	offline   atomic.Bool
}

func (d *liveDestinations) GetByID(context.Context, string) (*schema.DatabaseConnector, error) {
	return d.connector, nil
}

func (d *liveDestinations) OpenDestination(context.Context, string, string) (*dbtarget.OpenedDestination, error) {
	if d.offline.Load() {
		return nil, errors.New("dial tcp: connection refused")
	}
	return &dbtarget.OpenedDestination{DB: d.target, Kind: d.kind, Close: func() error { return nil }}, nil
}

type confirmedCounters struct {
	sources map[string]runtime.CounterSource
}

func (c confirmedCounters) CounterSources(context.Context, *workspace.WriteGroup) (map[string]runtime.CounterSource, error) {
	return c.sources, nil
}

// longPipeline is the production pipeline over an applied cumulative group on
// the shared long table, with a real local database (whose commits can be made
// to fail or to lose their response) and a real destination database.
type longPipeline struct {
	t          *testing.T
	group      *workspace.WriteGroup
	snapshot   *workspace.WriteGroupAppliedSnapshot
	groups     *fakeGroups
	dest       *liveDestinations
	counters   confirmedCounters
	localPath  string
	local      *sql.DB
	localFault *lostcommit.Faults
	targetFlt  *lostcommit.Faults
	store      *groupdelivery.Store
	pipe       *Pipeline
	table      string
	// postgres is set when the destination is a real PostgreSQL schema.
	postgres *longPostgres
	// columns are the destination columns the group's runtime layout froze.
	columns []dbtarget.ColumnInfo
	clock   time.Time
	mu      sync.Mutex
	errs    []error
}

// longPostgres is the disposable PostgreSQL schema a PostgreSQL pipeline writes to.
type longPostgres struct {
	dsn, schema string
	admin       *sql.DB
	// timescale makes the long table an owned hypertable (the group's mode is
	// timescale), created by the layout's own DDL.
	timescale bool
}

// newLongPipeline builds an applied cumulative long group over SQLite.
func newLongPipeline(t *testing.T) *longPipeline { return newLongPipelineOn(t, nil) }

// newLongPipelinePostgres builds the same group over a disposable PostgreSQL
// schema and fails the test when no database is configured: a skipped
// PostgreSQL case is not evidence.
func newLongPipelinePostgres(t *testing.T) *longPipeline {
	t.Helper()
	dsn := strings.TrimSpace(os.Getenv("POSTGRES_DSN"))
	if dsn == "" {
		t.Skip("POSTGRES_DSN is not set; the PostgreSQL case is NOT RUN")
	}
	admin, err := sql.Open("pgx", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = admin.Close() })
	require.NoError(t, admin.PingContext(t.Context()))
	name := fmt.Sprintf("gw_e2e_%d", time.Now().UnixNano())
	_, err = admin.ExecContext(t.Context(), `CREATE SCHEMA "`+name+`"`)
	require.NoError(t, err)
	t.Cleanup(func() { _, _ = admin.ExecContext(context.Background(), `DROP SCHEMA "`+name+`" CASCADE`) })
	return newLongPipelineOn(t, &longPostgres{dsn: dsn, schema: name, admin: admin})
}

// newLongPipelineTimescale builds the same group in the explicit timescale mode
// over a disposable PostgreSQL schema whose server offers TimescaleDB. The
// extension is installed by the fixture's database image, never by the code
// under test. Without a server that offers it the case is NOT RUN.
func newLongPipelineTimescale(t *testing.T) *longPipeline {
	t.Helper()
	f := newLongPipelinePostgresMode(t, true)
	return f
}

func newLongPipelinePostgresMode(t *testing.T, timescale bool) *longPipeline {
	t.Helper()
	dsn := strings.TrimSpace(os.Getenv("POSTGRES_DSN"))
	if dsn == "" {
		t.Skip("POSTGRES_DSN is not set; the PostgreSQL case is NOT RUN")
	}
	admin, err := sql.Open("pgx", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = admin.Close() })
	require.NoError(t, admin.PingContext(t.Context()))
	if timescale {
		var installed int
		require.NoError(t, admin.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM pg_extension WHERE extname = 'timescaledb'`).Scan(&installed))
		if installed == 0 {
			t.Skip("the destination database has no TimescaleDB extension installed; the timescale case is NOT RUN")
		}
	}
	name := fmt.Sprintf("gw_e2e_%d", time.Now().UnixNano())
	_, err = admin.ExecContext(t.Context(), `CREATE SCHEMA "`+name+`"`)
	require.NoError(t, err)
	t.Cleanup(func() { _, _ = admin.ExecContext(context.Background(), `DROP SCHEMA "`+name+`" CASCADE`) })
	return newLongPipelineOn(t, &longPostgres{dsn: dsn, schema: name, admin: admin, timescale: timescale})
}

func newLongPipelineOn(t *testing.T, pg *longPostgres, scaled ...bool) *longPipeline {
	t.Helper()
	f := &longPipeline{t: t, localPath: filepath.Join(t.TempDir(), "gateway.db"), clock: at(1), postgres: pg}
	f.localFault = &lostcommit.Faults{Lose: &atomic.Bool{}}
	f.targetFlt = &lostcommit.Faults{Lose: &atomic.Bool{}}

	f.group = &workspace.WriteGroup{
		ID: "group-C", WorkspaceID: "workspace-A", Revision: "rev-1", AppliedRevision: "rev-1", Status: workspace.WriteGroupStatusReady,
		LayoutVersion: workspace.WriteGroupLayoutManagedLongV1, StorageMode: workspace.WriteGroupStorageModePlain,
		RecordingProfile: &workspace.WriteGroupRecordingProfile{Kind: workspace.WriteGroupProfileCumulative, Interval: 10, Unit: workspace.WriteGroupIntervalUnitMinute},
		Members: []workspace.WriteGroupMember{
			{DeviceID: "device-1", PointID: "point-1", TagID: longKwhFirst, TargetColumn: "kwh_1", CircuitID: "circuit-1", MeasurementID: ptr("meas-1"), SourceRevision: "src-1", MappingRevision: "map-1"},
			{DeviceID: "device-1", PointID: "point-2", TagID: longKwhSecond, TargetColumn: "kwh_2", CircuitID: "circuit-1", MeasurementID: ptr("meas-2"), SourceRevision: "src-2", MappingRevision: "map-2"},
		},
		Destination: workspace.WriteGroupDestination{ConnectorID: "connector-1", ConnectorRevision: "crev-1", TableSchema: f.schemaName(), StorageStrategy: workspace.WriteGroupStorageStrategyManaged},
		RowPolicy:   workspace.WriteGroupRowPolicy{IntervalSeconds: int(longWindow / time.Second), RecordKeyColumn: "record_id"},
		WritePolicy: workspace.WriteGroupWritePolicy{Mode: "append", DedupeCapability: groupdelivery.DedupeReceipt},
	}
	if pg != nil && pg.timescale {
		f.group.StorageMode = workspace.WriteGroupStorageModeTimescale
	}
	table, owner, ok := workspace.LongStorageIdentity(f.group)
	require.True(t, ok)
	f.table = table
	f.group.Destination.TableName = table
	f.counters = confirmedCounters{sources: map[string]runtime.CounterSource{}}
	for _, member := range f.group.Members {
		f.counters.sources[member.TagID] = runtime.CounterSource{
			MeasurementID: *member.MeasurementID, DefinitionRevision: "def-1", SourceBindingRevision: "bind-1",
			SeriesEpoch: "epoch-1", Unit: "kWh", Kind: aggregation.CounterKindMonotonic,
		}
	}

	dialect, kind := "sqlite", schema.DatabaseConnectorKindSQLite
	if pg != nil {
		dialect, kind = "postgres", schema.DatabaseConnectorKindPostgres
	}
	layout, err := recordingplan.ManagedLongGroupLayout(dialect, table, owner)
	require.NoError(t, err)
	if pg != nil && pg.timescale {
		layout, err = recordingplan.ManagedLongTimescaleLayout(table, owner)
		require.NoError(t, err)
		layout.Partition.ExtensionSchema = "public" // where the fixture's image installs it
	}
	statements, err := recordingplan.GenerateGroupSchemaDDL(dialect, layout)
	require.NoError(t, err)
	var columns []dbtarget.ColumnInfo
	var target *sql.DB
	if pg == nil {
		targetPath := filepath.Join(t.TempDir(), "target.db")
		plain, err := sql.Open("sqlite", targetPath+"?_pragma=busy_timeout(5000)")
		require.NoError(t, err)
		for _, statement := range statements {
			_, err := plain.ExecContext(t.Context(), statement)
			require.NoError(t, err, statement)
		}
		columns = sqliteColumnInfo(t, plain, table)
		require.NoError(t, plain.Close())
		target = lostcommit.OpenFaults(targetPath+"?_pragma=busy_timeout(5000)", f.targetFlt)
	} else {
		tx, err := pg.admin.BeginTx(t.Context(), nil)
		require.NoError(t, err)
		_, err = tx.ExecContext(t.Context(), `SET LOCAL search_path TO "`+pg.schema+`"`)
		require.NoError(t, err)
		for _, statement := range statements {
			_, err = tx.ExecContext(t.Context(), statement)
			require.NoError(t, err, statement)
		}
		require.NoError(t, tx.Commit())
		for _, column := range layout.Columns {
			columns = append(columns, dbtarget.ColumnInfo{Name: column.Name, DataType: column.SQLType, Nullable: column.Nullable, PrimaryKey: column.PrimaryKey, Unique: column.PrimaryKey})
		}
		target = lostcommit.WrapFaults(stdlib.GetDefaultDriver(), pg.dsn, f.targetFlt)
	}
	t.Cleanup(func() { _ = target.Close() })
	f.columns = columns
	f.dest = &liveDestinations{
		connector: &schema.DatabaseConnector{ID: "connector-1", Kind: kind, IdentityRevision: "crev-1", Enabled: true},
		target:    target, kind: kind,
	}

	f.snapshot = &workspace.WriteGroupAppliedSnapshot{
		WorkspaceID: f.group.WorkspaceID, GroupID: f.group.ID, GroupRevision: f.group.Revision, AppliedRevision: f.group.AppliedRevision,
		EffectiveAt: at(0), Group: f.group,
	}
	digest := sha256.Sum256([]byte(table))
	frozen := *f.snapshot
	frozen.RuntimeLayout = &workspace.WriteGroupRuntimeLayout{
		Dialect: dbtarget.SQLDialect(dialect), Columns: columns,
		TagTypes:          map[string]schema.DataType{longKwhFirst: schema.DataTypeUint64, longKwhSecond: schema.DataTypeUint64},
		ReceiptTableReady: true, SchemaDigest: hex.EncodeToString(digest[:]),
	}
	if pg != nil && pg.timescale {
		frozen.RuntimeLayout.Hypertable = &workspace.WriteGroupHypertableEvidence{TimeColumn: "bucket_start", ChunkIntervalMicros: 7 * 24 * 3600 * 1000000}
	}
	f.groups = &fakeGroups{
		groups:    []*workspace.WriteGroup{liveGroup(f.group, workspace.WriteGroupStatusReady)},
		snapshots: map[string][]*workspace.WriteGroupAppliedSnapshot{f.group.ID: {f.snapshot}},
	}
	if len(scaled) > 0 && scaled[0] {
		scaleFrozenFixture(&frozen)
	}
	freezeFixtureCounters(t, &frozen, f.counters.sources)
	f.openLocal(&frozen)
	return f
}

func ptr[T any](value T) *T { return &value }

// schemaName is the destination schema: "main" for SQLite, the disposable schema for PostgreSQL.
func (f *longPipeline) schemaName() string {
	if f.postgres != nil {
		return f.postgres.schema
	}
	return "main"
}

func sqliteColumnInfo(t *testing.T, db *sql.DB, table string) []dbtarget.ColumnInfo {
	t.Helper()
	rows, err := db.QueryContext(t.Context(), fmt.Sprintf(`PRAGMA table_info('%s')`, table))
	require.NoError(t, err)
	defer rows.Close()
	var columns []dbtarget.ColumnInfo
	for rows.Next() {
		var cid, notNull, pk int
		var name, dataType string
		var defaultValue sql.NullString
		require.NoError(t, rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &pk))
		columns = append(columns, dbtarget.ColumnInfo{Name: name, DataType: dataType, Nullable: notNull == 0, PrimaryKey: pk > 0, Unique: pk > 0})
	}
	require.NoError(t, rows.Err())
	return columns
}

// openLocal opens the gateway's local database, migrating it once, and starts a
// pipeline over it. A non-nil frozen snapshot is saved the way Apply does.
func (f *longPipeline) openLocal(frozen *workspace.WriteGroupAppliedSnapshot) {
	f.t.Helper()
	plain, err := sql.Open("sqlite", f.localPath+"?_pragma=busy_timeout(5000)")
	require.NoError(f.t, err)
	require.NoError(f.t, datalink.NewMigrator().Migrate(plain))
	require.NoError(f.t, plain.Close())
	f.local = lostcommit.OpenFaults(f.localPath+"?_pragma=busy_timeout(5000)", f.localFault)
	f.t.Cleanup(func() { _ = f.local.Close() })
	f.store = groupdelivery.NewStore(f.local)
	if frozen != nil {
		payload, err := json.Marshal(frozen)
		require.NoError(f.t, err)
		require.NoError(f.t, f.store.SaveRuntimeVersion(f.t.Context(), groupdelivery.GroupKey{
			WorkspaceID: frozen.WorkspaceID, GroupID: frozen.GroupID, GroupRevision: frozen.AppliedRevision,
		}, payload))
	}
	f.startPipeline()
}

func (f *longPipeline) startPipeline() {
	f.t.Helper()
	immediate := func(int) (string, time.Time) { return groupdelivery.StateRetrying, time.Now().UTC().Add(-time.Second) }
	f.pipe = New(Dependencies{
		Groups: f.groups, Tags: fakeTags{missing: map[string]bool{}}, Destinations: f.dest, Inspector: &fakeInspector{},
		Store: f.store, Counters: f.counters,
	}, Config{
		NodeID: "node-1", Owner: "node-1/test", Now: func() time.Time { return f.clock }, ReconcileInterval: 5 * time.Second,
		Sender: groupdelivery.SenderConfig{Backoff: immediate},
		OnError: func(err error) {
			f.mu.Lock()
			defer f.mu.Unlock()
			f.errs = append(f.errs, err)
		},
	})
	require.NoError(f.t, f.pipe.Reconcile(f.t.Context()))
	require.Equal(f.t, StateActive, f.statusOf().State, "the applied long group is built from its frozen proof: %s", f.statusOf().Reason)
}

func (f *longPipeline) statusOf() GroupStatus {
	for _, status := range f.pipe.Status() {
		if status.GroupID == f.group.ID {
			return status
		}
	}
	return GroupStatus{}
}

// restart simulates a process restart: the local database and the pipeline are
// reopened from the files; nothing in memory survives.
func (f *longPipeline) restart() {
	f.t.Helper()
	require.NoError(f.t, f.local.Close())
	f.errs = nil
	f.openLocal(nil)
}

func (f *longPipeline) reportedErrors() []error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]error(nil), f.errs...)
}

// reading delivers one counter reading for a member through the production intake.
func (f *longPipeline) reading(member int, id string, observed time.Time, value uint64) error {
	f.t.Helper()
	m := f.group.Members[member]
	f.clock = observed.Add(time.Second)
	return f.pipe.AcceptSample(f.t.Context(), measurement.SampleEnvelope{
		SampleID: id, WorkspaceID: f.group.WorkspaceID, DeviceID: m.DeviceID, PointID: m.PointID, TagID: m.TagID,
		SourceRevision: m.SourceRevision, MappingRevision: m.MappingRevision, ObservedAt: observed, ReceivedAt: observed.Add(40 * time.Millisecond),
		TimeOrigin: "gateway", Value: value, Quality: schema.QualityGood,
	})
}

// closeUntil advances the gateway clock to when and closes every due window.
func (f *longPipeline) closeUntil(when time.Time) {
	f.t.Helper()
	f.clock = when
	f.pipe.TickAll(f.t.Context())
}

// deliver runs delivery cycles until the outbox has nothing left to try.
func (f *longPipeline) deliver() { f.deliverCycles(4) }

// deliverCycles runs a fixed number of production delivery cycles: dispatcher,
// sender, the saved destination and the local settlement.
func (f *longPipeline) deliverCycles(cycles int) {
	f.t.Helper()
	sender := groupdelivery.NewSender(f.store, destinationResolver{destinations: f.dest}, f.pipe.senderConfig())
	dispatcher := groupdelivery.NewDispatcher(f.store, sender, f.pipe.config.Delivery)
	for range cycles {
		_, err := dispatcher.RunOnce(f.t.Context())
		require.NoError(f.t, err)
	}
}

func (f *longPipeline) count(db *sql.DB, query string) int {
	f.t.Helper()
	var n int
	require.NoError(f.t, db.QueryRowContext(f.t.Context(), query).Scan(&n))
	return n
}

// plainLocal and plainTarget are fault-free views for assertions.
func (f *longPipeline) plain(path string) *sql.DB {
	f.t.Helper()
	db, err := sql.Open("sqlite", path+"?_pragma=busy_timeout(5000)")
	require.NoError(f.t, err)
	f.t.Cleanup(func() { _ = db.Close() })
	return db
}

func (f *longPipeline) plainLocal() *sql.DB { return f.plain(f.localPath) }

func (f *longPipeline) plainTarget() *sql.DB {
	f.t.Helper()
	var path string
	require.NoError(f.t, f.dest.target.QueryRowContext(f.t.Context(), `SELECT file FROM pragma_database_list WHERE name = 'main'`).Scan(&path))
	return f.plain(path)
}

// targetRow is one stored record of the destination's long table.
type targetRow struct {
	status, raw, delta, deltaFrom, deltaTo, provenance sql.NullString
}

func (f *longPipeline) targetRows() map[string]targetRow {
	f.t.Helper()
	query := fmt.Sprintf(`SELECT tag_id, bucket_start, status, CAST(value_decimal AS TEXT), CAST(delta_value AS TEXT), delta_from, delta_to, provenance FROM %s`, f.table)
	db := f.plainTarget
	if f.postgres != nil {
		utc := func(column string) string {
			return `to_char(` + column + ` AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS.US"Z"')`
		}
		query = fmt.Sprintf(`SELECT tag_id, %s, status, CAST(value_decimal AS TEXT), CAST(delta_value AS TEXT), %s, %s, CAST(provenance AS TEXT) FROM %q.%q`,
			utc("bucket_start"), utc("delta_from"), utc("delta_to"), f.postgres.schema, f.table)
		db = func() *sql.DB { return f.postgres.admin }
	}
	rows, err := db().QueryContext(f.t.Context(), query)
	require.NoError(f.t, err)
	defer rows.Close()
	result := map[string]targetRow{}
	for rows.Next() {
		var tag, bucket string
		var row targetRow
		require.NoError(f.t, rows.Scan(&tag, &bucket, &row.status, &row.raw, &row.delta, &row.deltaFrom, &row.deltaTo, &row.provenance))
		result[tag+"@"+bucket] = row
	}
	require.NoError(f.t, rows.Err())
	return result
}

const longTime = "2006-01-02T15:04:05.000000Z"

func windowKey(tag string, window int) string {
	return tag + "@" + at(window*int(longWindow/time.Second)).Format(longTime)
}

// baseline reads the durable baseline of a member: its window end and the
// serialized checkpoint, or ("", "") when none is committed.
func (f *longPipeline) baseline(member int) (windowEnd, payload string) {
	f.t.Helper()
	m := f.group.Members[member]
	key := fmt.Sprintf("[%q,%q,%q]", m.DeviceID, m.PointID, m.TagID)
	err := f.plainLocal().QueryRowContext(f.t.Context(), `SELECT window_end, payload FROM wg_delivery_counter_baselines WHERE group_id = ? AND member_key = ?`,
		f.group.ID, key).Scan(&windowEnd, &payload)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ""
	}
	require.NoError(f.t, err)
	return windowEnd, payload
}

func (f *longPipeline) outboxDigests() map[string]string {
	f.t.Helper()
	rows, err := f.plainLocal().QueryContext(f.t.Context(), `SELECT effect_key, payload_digest FROM wg_delivery_outbox`)
	require.NoError(f.t, err)
	defer rows.Close()
	digests := map[string]string{}
	for rows.Next() {
		var key, digest string
		require.NoError(f.t, rows.Scan(&key, &digest))
		digests[key] = digest
	}
	require.NoError(f.t, rows.Err())
	return digests
}
