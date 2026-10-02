package groupdelivery_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

// These tests kill a real child process with SIGKILL. The child runs the real
// durable store, boundary and sender; only the process boundary is the test's.

const (
	roleEnv       = "GW_CRASH_ROLE"
	localDBEnv    = "GW_CRASH_LOCAL_DB"
	pgDSNEnv      = "GW_CRASH_PG_DSN"
	readyEnv      = "GW_CRASH_READY"
	schemaEnv     = "GW_CRASH_SCHEMA"
	capabilityEnv = "GW_CRASH_CAPABILITY"
	effectEnv     = "GW_CRASH_EFFECT"
)

var crashBase = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

func crashGroup(pgSchema, capability string) *workspace.WriteGroup {
	return &workspace.WriteGroup{
		ID: "group-crash", WorkspaceID: "workspace-crash", Revision: "rev-1", AppliedRevision: "rev-1", Status: workspace.WriteGroupStatusReady,
		Members: []workspace.WriteGroupMember{{
			DeviceID: "device-1", PointID: "point-1", TagID: "tag-t", TargetColumn: "temperature", Required: true,
			SourceRevision: "src-1", MappingRevision: "map-1",
		}},
		Destination: workspace.WriteGroupDestination{ConnectorID: "connector-pg", ConnectorRevision: "crev-1", TableSchema: pgSchema, TableName: "readings"},
		RowPolicy:   workspace.WriteGroupRowPolicy{IntervalSeconds: 10},
		WritePolicy: workspace.WriteGroupWritePolicy{DedupeCapability: capability},
	}
}

func crashBoundary(t testing.TB, store *groupdelivery.Store, pgSchema, capability string, now time.Time) *runtime.GroupBoundary {
	t.Helper()
	group := crashGroup(pgSchema, capability)
	boundary, err := runtime.NewGroupBoundary(runtime.GroupBoundaryConfig{
		Group: group, TagTypes: map[string]schema.DataType{"tag-t": schema.DataTypeFloat64},
		Dialect: dbtarget.SQLDialectPostgres, Columns: []dbtarget.ColumnInfo{{Name: "temperature", DataType: "DOUBLE PRECISION"}},
		FirstBucket: crashBase, MaxFutureSkew: 5 * time.Second, ReceiptTableReady: capability == groupdelivery.DedupeReceipt,
		Ledger: store.Ledger(groupdelivery.GroupKey{WorkspaceID: group.WorkspaceID, GroupID: group.ID, GroupRevision: "rev-1"}),
		Clock:  func() time.Time { return now },
	})
	require.NoError(t, err)
	return boundary
}

func crashEnvelope(id string, observed time.Time, value float64) measurement.SampleEnvelope {
	return measurement.SampleEnvelope{
		SampleID: id, WorkspaceID: "workspace-crash", DeviceID: "device-1", PointID: "point-1", TagID: "tag-t",
		SourceRevision: "src-1", MappingRevision: "map-1", ObservedAt: observed, ReceivedAt: observed,
		TimeOrigin: "gateway", Value: value, Quality: schema.QualityGood,
	}
}

type pgResolver struct{ dsn string }

func (r pgResolver) Resolve(ctx context.Context, _ groupdelivery.OutboxItem) (groupdelivery.Target, error) {
	db, err := sql.Open("pgx", r.dsn)
	if err != nil {
		return groupdelivery.Target{}, err
	}
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return groupdelivery.Target{}, err
	}
	return groupdelivery.Target{DB: db, Kind: schema.DatabaseConnectorKindPostgres, Close: db.Close}, nil
}

func openLocal(t testing.TB, path string) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", path+"?_pragma=busy_timeout(10000)")
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	return db
}

// TestCrashChildProcess is the child: it does nothing unless the parent set a role.
func TestCrashChildProcess(t *testing.T) {
	role := os.Getenv(roleEnv)
	if role == "" {
		t.Skip("helper process; only runs when spawned by the crash tests")
	}
	ctx := context.Background()
	local := openLocal(t, os.Getenv(localDBEnv))
	store := groupdelivery.NewStore(local)
	pgSchema, capability := os.Getenv(schemaEnv), os.Getenv(capabilityEnv)
	switch role {
	case "ack":
		now := crashBase.Add(8 * time.Second)
		boundary := crashBoundary(t, store, pgSchema, capability, now)
		require.NoError(t, boundary.AcceptSample(ctx, crashEnvelope("crash-sample-1", crashBase.Add(5*time.Second), 21.5)))
	case "deliver":
		sender := groupdelivery.NewSender(store, pgResolver{dsn: os.Getenv(pgDSNEnv)}, groupdelivery.SenderConfig{
			Owner: "node-1/first", DeliveryTimeout: time.Minute,
		})
		go func() { _, _ = sender.Deliver(ctx, os.Getenv(effectEnv)) }()
	default:
		t.Fatalf("unknown role %q", role)
	}
	require.NoError(t, os.WriteFile(os.Getenv(readyEnv), []byte("ready"), 0o600))
	time.Sleep(time.Hour) // wait to be killed (select{} would be a Go deadlock exit, not a kill)
}

func spawnChild(t *testing.T, role string, env map[string]string) *exec.Cmd {
	t.Helper()
	cmd := exec.CommandContext(t.Context(), os.Args[0], "-test.run=^TestCrashChildProcess$", "-test.count=1")
	cmd.Env = append(os.Environ(), roleEnv+"="+role)
	for key, value := range env {
		cmd.Env = append(cmd.Env, key+"="+value)
	}
	require.NoError(t, cmd.Start())
	t.Cleanup(func() { _ = cmd.Process.Kill(); _, _ = cmd.Process.Wait() })
	return cmd
}

func killHard(t *testing.T, cmd *exec.Cmd) {
	t.Helper()
	require.NoError(t, cmd.Process.Signal(syscall.SIGKILL))
	err := cmd.Wait()
	var exitErr *exec.ExitError
	require.ErrorAs(t, err, &exitErr, "the child must die by signal")
	status, ok := exitErr.Sys().(syscall.WaitStatus)
	require.True(t, ok && status.Signaled() && status.Signal() == syscall.SIGKILL, "killed by SIGKILL, not a graceful exit")
}

func waitFor(t *testing.T, what string, condition func() bool) {
	t.Helper()
	require.Eventually(t, condition, 30*time.Second, 25*time.Millisecond, what)
}

func postgresDSN(t *testing.T) string {
	t.Helper()
	dsn := strings.TrimSpace(os.Getenv("POSTGRES_DSN"))
	if dsn == "" {
		t.Skip("Skipping PostgreSQL crash test: POSTGRES_DSN not set")
	}
	return dsn
}

func newPGSchema(t *testing.T, admin *sql.DB, withSlowInsert bool) string {
	t.Helper()
	name := fmt.Sprintf("gw_crash_%d", time.Now().UnixNano())
	run := func(statement string) {
		t.Helper()
		_, err := admin.ExecContext(t.Context(), statement)
		require.NoError(t, err)
	}
	run(`CREATE SCHEMA "` + name + `"`)
	t.Cleanup(func() { _, _ = admin.ExecContext(context.Background(), `DROP SCHEMA "`+name+`" CASCADE`) })
	run(`CREATE TABLE "` + name + `".readings (temperature DOUBLE PRECISION NOT NULL)`)
	require.NoError(t, dbtarget.CreateEffectReceiptTable(t.Context(), admin, schema.DatabaseConnectorKindPostgres, name))
	if withSlowInsert {
		run(`CREATE FUNCTION "` + name + `".slow() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN PERFORM pg_sleep(60); RETURN NEW; END $$`)
		run(`CREATE TRIGGER slow BEFORE INSERT ON "` + name + `".readings FOR EACH ROW EXECUTE FUNCTION "` + name + `".slow()`)
	}
	return name
}

func countPG(t *testing.T, admin *sql.DB, pgSchema, table string) int {
	t.Helper()
	var n int
	require.NoError(t, admin.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM "`+pgSchema+`".`+table).Scan(&n))
	return n
}

func TestCrashKillAfterAckLosesNothingAndDeliversOnce(t *testing.T) {
	dsn := postgresDSN(t)
	admin, err := sql.Open("pgx", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = admin.Close() })
	pgSchema := newPGSchema(t, admin, false)

	dir := t.TempDir()
	localPath, ready := filepath.Join(dir, "gateway.db"), filepath.Join(dir, "ready")
	openLocal(t, localPath) // migrate before the child opens it
	child := spawnChild(t, "ack", map[string]string{
		localDBEnv: localPath, readyEnv: ready, schemaEnv: pgSchema, capabilityEnv: groupdelivery.DedupeReceipt,
	})
	waitFor(t, "child ACKed its sample", func() bool { _, err := os.Stat(ready); return err == nil })
	killHard(t, child) // dies after the ACK, before any bucket closed

	local := openLocal(t, localPath)
	var journaled, outbox int
	require.NoError(t, local.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM wg_delivery_samples WHERE consumed = 0`).Scan(&journaled))
	require.NoError(t, local.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM wg_delivery_outbox`).Scan(&outbox))
	require.Equal(t, 1, journaled, "the ACKed sample survived the SIGKILL in the journal")
	require.Zero(t, outbox)

	// A new process recovers from durable state alone.
	store := groupdelivery.NewStore(local)
	recovered := crashBoundary(t, store, pgSchema, groupdelivery.DedupeReceipt, crashBase.Add(11*time.Second))
	require.NoError(t, recovered.Tick(t.Context(), crashBase.Add(11*time.Second)))
	sender := groupdelivery.NewSender(store, pgResolver{dsn: dsn}, groupdelivery.SenderConfig{Owner: "node-1/second"})
	dispatcher := groupdelivery.NewDispatcher(store, sender, groupdelivery.DispatcherConfig{})
	for i := 0; i < 3; i++ { // repeated cycles must not duplicate
		_, err := dispatcher.RunOnce(t.Context())
		require.NoError(t, err)
	}
	require.Equal(t, 1, countPG(t, admin, pgSchema, "readings"), "the row reaches PostgreSQL exactly once")
	var value float64
	require.NoError(t, admin.QueryRowContext(t.Context(), `SELECT temperature FROM "`+pgSchema+`".readings`).Scan(&value))
	require.Equal(t, 21.5, value)
	require.Equal(t, 1, countPG(t, admin, pgSchema, dbtarget.EffectReceiptTable))
}

func TestCrashKillDuringPostgresInsertIsRecoveredByCapability(t *testing.T) {
	dsn := postgresDSN(t)
	admin, err := sql.Open("pgx", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = admin.Close() })

	for _, capability := range []string{groupdelivery.DedupeReceipt, groupdelivery.DedupeNone} {
		t.Run(capability, func(t *testing.T) {
			pgSchema := newPGSchema(t, admin, true) // every INSERT sleeps inside its transaction
			dir := t.TempDir()
			localPath, ready := filepath.Join(dir, "gateway.db"), filepath.Join(dir, "ready")
			local := openLocal(t, localPath)
			store := groupdelivery.NewStore(local)

			// Close one bucket into the outbox in this process; the child will try to deliver it.
			boundary := crashBoundary(t, store, pgSchema, capability, crashBase.Add(8*time.Second))
			require.NoError(t, boundary.AcceptSample(t.Context(), crashEnvelope("crash-sample-1", crashBase.Add(5*time.Second), 21.5)))
			require.NoError(t, boundary.Tick(t.Context(), crashBase.Add(11*time.Second)))
			var effectKey string
			require.NoError(t, local.QueryRowContext(t.Context(), `SELECT effect_key FROM wg_delivery_outbox`).Scan(&effectKey))

			child := spawnChild(t, "deliver", map[string]string{
				localDBEnv: localPath, readyEnv: ready, pgDSNEnv: dsn, schemaEnv: pgSchema, capabilityEnv: capability, effectEnv: effectKey,
			})
			waitFor(t, "child is inside the PostgreSQL INSERT transaction", func() bool {
				var active int
				require.NoError(t, admin.QueryRowContext(t.Context(),
					`SELECT COUNT(*) FROM pg_stat_activity WHERE state = 'active' AND query LIKE $1`, `INSERT INTO "`+pgSchema+`"%`).Scan(&active))
				var state string
				require.NoError(t, local.QueryRowContext(t.Context(), `SELECT state FROM wg_delivery_outbox WHERE effect_key = ?`, effectKey).Scan(&state))
				return active == 1 && state == groupdelivery.StateSending
			})
			killHard(t, child) // killed mid-transaction: nothing was committed

			// The server rolls back the orphaned transaction.
			_, err := admin.ExecContext(t.Context(),
				`SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE query LIKE $1 AND pid <> pg_backend_pid()`, `INSERT INTO "`+pgSchema+`"%`)
			require.NoError(t, err)
			_, err = admin.ExecContext(t.Context(), `DROP TRIGGER slow ON "`+pgSchema+`".readings`)
			require.NoError(t, err)
			require.Zero(t, countPG(t, admin, pgSchema, "readings"), "the killed transaction left no row")

			// A new incarnation of the same node recovers the interrupted claim.
			recovered, err := store.RecoverStaleClaims(t.Context(), "node-1", "node-1/second", time.Now().UTC())
			require.NoError(t, err)
			require.Equal(t, 1, recovered)
			sender := groupdelivery.NewSender(store, pgResolver{dsn: dsn}, groupdelivery.SenderConfig{Owner: "node-1/second"})
			dispatcher := groupdelivery.NewDispatcher(store, sender, groupdelivery.DispatcherConfig{})
			_, err = dispatcher.RunOnce(t.Context())
			require.NoError(t, err)

			var state string
			require.NoError(t, local.QueryRowContext(t.Context(), `SELECT state FROM wg_delivery_outbox WHERE effect_key = ?`, effectKey).Scan(&state))
			if capability == groupdelivery.DedupeReceipt {
				time.Sleep(100 * time.Millisecond)
				_, err = dispatcher.RunOnce(t.Context())
				require.NoError(t, err)
				require.NoError(t, local.QueryRowContext(t.Context(), `SELECT state FROM wg_delivery_outbox WHERE effect_key = ?`, effectKey).Scan(&state))
				require.Equal(t, groupdelivery.StateCommitted, state, "a destination receipt makes recovery after a crash safe")
				require.Equal(t, 1, countPG(t, admin, pgSchema, "readings"), "exactly one row")
				require.Equal(t, 1, countPG(t, admin, pgSchema, dbtarget.EffectReceiptTable))
				return
			}
			// Without dedupe the earlier attempt's outcome could not be known, so
			// the row waits for reconciliation instead of being blindly resent.
			require.Equal(t, groupdelivery.StateUnknown, state)
			require.Zero(t, countPG(t, admin, pgSchema, "readings"), "no blind retry was made")
		})
	}
}
