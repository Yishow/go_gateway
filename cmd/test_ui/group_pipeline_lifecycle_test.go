package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

type lifecycleTarget struct {
	env  *outageEnv
	down func()
	up   func()
	rows func() [][2]float64
}

func newLifecycleTarget(t *testing.T, kind string) *lifecycleTarget {
	t.Helper()
	if kind == "sqlite" {
		env := newOutageEnv(t)
		file := env.targets["connector-A"]
		return &lifecycleTarget{env: env,
			down: func() { require.NoError(t, os.Rename(file, file+".offline")) },
			up:   func() { require.NoError(t, os.Rename(file+".offline", file)) },
			rows: func() [][2]float64 { return destinationRows(t, file) },
		}
	}
	dsn, fields := postgresEnvFields(t)
	require.Equal(t, "127.0.0.1", fields["host"], "use the owned loopback fixture")
	require.Equal(t, "55432", fields["port"])
	require.Equal(t, "gwtest", fields["dbname"])
	admin, err := sql.Open("pgx", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = admin.Close() })
	require.NoError(t, admin.PingContext(t.Context()))
	var selectedSchema string
	var selectedProxy *tcpProxy
	setup := func(t *testing.T, env *outageEnv, name string) (string, string, string) {
		t.Helper()
		pgSchema := fmt.Sprintf("gw_lifecycle_%s_%d", name, time.Now().UnixNano())
		_, err := admin.ExecContext(t.Context(), `CREATE SCHEMA "`+pgSchema+`"`)
		require.NoError(t, err)
		t.Cleanup(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_, err := admin.ExecContext(ctx, `DROP SCHEMA "`+pgSchema+`" CASCADE`)
			require.NoError(t, err)
		})
		_, err = admin.ExecContext(t.Context(), `CREATE TABLE "`+pgSchema+`".readings (temperature DOUBLE PRECISION NOT NULL, pressure BIGINT NOT NULL)`)
		require.NoError(t, err)
		require.NoError(t, dbtarget.CreateEffectReceiptTable(t.Context(), admin, schema.DatabaseConnectorKindPostgres, pgSchema))
		proxy := newTCPProxy(t, fields["host"]+":"+fields["port"])
		host, port := proxy.host()
		config, err := json.Marshal(map[string]string{"host": host, "port": port, "user": fields["user"], "password": fields["password"], "database": fields["dbname"], "sslmode": "disable"})
		require.NoError(t, err)
		if name == "A" {
			selectedSchema, selectedProxy = pgSchema, proxy
		}
		return "postgres", string(config), pgSchema
	}
	env := newOutageEnvWith(t, setup)
	env.dedupe = "receipt"
	return &lifecycleTarget{env: env, down: selectedProxy.down,
		up: func() { selectedProxy.up(selectedProxy.addr) },
		rows: func() [][2]float64 {
			rows, err := admin.QueryContext(t.Context(), `SELECT temperature, pressure FROM "`+selectedSchema+`".readings ORDER BY temperature`)
			require.NoError(t, err)
			defer rows.Close()
			var values [][2]float64
			for rows.Next() {
				var temperature float64
				var pressure int64
				require.NoError(t, rows.Scan(&temperature, &pressure))
				values = append(values, [2]float64{temperature, float64(pressure)})
			}
			require.NoError(t, rows.Err())
			return values
		},
	}
}

var lifecycleBase = time.Date(2026, 10, 4, 0, 0, 0, 0, time.UTC)

func lifecycleMutation(t *testing.T, env *outageEnv, id string) workspace.WriteGroupMutation {
	t.Helper()
	record, err := env.services.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	live, err := env.services.writeGroups.Get(t.Context(), id)
	require.NoError(t, err)
	return workspace.WriteGroupMutation{WorkspaceID: record.ID, ExpectedWorkspaceRevision: record.DatabaseSetupRevision,
		ExpectedGroupRevision: live.Group.Revision, ExpectedConnectorRevision: live.Group.Destination.ConnectorRevision}
}

func reopenLifecycleConfiguration(t *testing.T, env *outageEnv, clock *testClock) {
	t.Helper()
	var seq int
	var name, path string
	require.NoError(t, env.db.QueryRowContext(t.Context(), `PRAGMA database_list`).Scan(&seq, &name, &path))
	require.NoError(t, env.db.Close())
	db, err := sql.Open("sqlite", path+"?_pragma=busy_timeout(5000)")
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	datalink.ApplySQLitePoolDefaults(db)
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	env.db = db
	env.services, err = wireGatewayServices(db, connector.NewConnectionManager(connector.ConnectionManagerConfig{}))
	require.NoError(t, err)
	env.services.writeGroups.WithClock(clock.now)
}

func lifecycleAccept(t *testing.T, env *outageEnv, pipe *grouppipeline.Pipeline, group *workspace.WriteGroup, clock *testClock, second int) {
	t.Helper()
	at := lifecycleBase.Add(time.Duration(second) * time.Second)
	clock.set(at)
	require.NoError(t, pipe.AcceptSample(t.Context(), env.envelope(group, 0, fmt.Sprintf("t-%d", second), at, float64(second))))
	require.NoError(t, pipe.AcceptSample(t.Context(), env.envelope(group, 1, fmt.Sprintf("p-%d", second), at, int64(second*10))))
}

func TestProductionLifecycleDraftAndOfflineColdStartIntake(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		t.Run(kind, func(t *testing.T) {
			target := newLifecycleTarget(t, kind)
			env, clock := target.env, &testClock{}
			clock.set(lifecycleBase.Add(-time.Second))
			env.services.writeGroups.WithClock(clock.now)
			group := env.createAndApply(t, "A")
			clock.set(lifecycleBase.Add(time.Second))
			first := env.pipeline(clock, "node-1/draft")
			require.NoError(t, first.Reconcile(t.Context()))
			mutation := lifecycleMutation(t, env, group.ID)
			payload, err := json.Marshal(group)
			require.NoError(t, err)
			var draft workspace.WriteGroup
			require.NoError(t, json.Unmarshal(payload, &draft))
			draft.Name, draft.Members = "Saved draft", draft.Members[:1]
			mutation.Group = &draft
			saved, err := env.services.writeGroups.Update(t.Context(), group.ID, mutation)
			require.NoError(t, err)
			require.Equal(t, group.AppliedRevision, saved.Group.AppliedRevision)
			require.NoError(t, first.Reconcile(t.Context()))
			lifecycleAccept(t, env, first, group, clock, 2)
			clock.set(lifecycleBase.Add(11 * time.Second))
			first.TickAll(t.Context())
			reopenLifecycleConfiguration(t, env, clock)
			second := env.pipeline(clock, "node-1/offline")
			target.down()
			require.NoError(t, second.Reconcile(t.Context()))
			require.True(t, second.WantsSample("device-1", "point-p-A", "tag-p-A"), "the applied member removed only in the draft still belongs to intake")
			lifecycleAccept(t, env, second, group, clock, 12)
			clock.set(lifecycleBase.Add(21 * time.Second))
			second.TickAll(t.Context())
			require.Equal(t, 2, outboxCount(t, env.db, "connector-A", "1=1"))
			if kind == "sqlite" {
				_, err := os.Stat(env.targets["connector-A"])
				require.ErrorIs(t, err, os.ErrNotExist, "cold recovery did not open or create the offline destination")
			}
			target.up()
			require.NoError(t, second.Start(t.Context()))
			t.Cleanup(func() { require.NoError(t, second.Stop(2*time.Second)) })
			require.Eventually(t, func() bool { return outboxCount(t, env.db, "connector-A", "state = 'sql_committed'") == 2 }, 10*time.Second, 20*time.Millisecond)
			require.Equal(t, [][2]float64{{2, 20}, {12, 120}}, target.rows())
		})
	}
}

// A child process uses these paths only when explicitly spawned by the crash
// matrix. Ordinary go test runs do not start a child or wait here.
const lifecycleCrashEnv = "GW_LIFECYCLE_CRASH"

func TestProductionLifecycleCrashChild(t *testing.T) {
	path := os.Getenv(lifecycleCrashEnv)
	if path == "" {
		t.Skip("helper process")
	}
	db, err := sql.Open("sqlite", path+"?_pragma=busy_timeout(5000)")
	require.NoError(t, err)
	datalink.ApplySQLitePoolDefaults(db)
	clock := &testClock{}
	clock.set(lifecycleBase.Add(2 * time.Second))
	services, err := wireGatewayServices(db, connector.NewConnectionManager(connector.ConnectionManagerConfig{}))
	require.NoError(t, err)
	env := &outageEnv{db: db, services: services}
	env.services.writeGroups.WithClock(clock.now)
	id := os.Getenv("GW_LIFECYCLE_GROUP")
	live, err := env.services.writeGroups.Get(t.Context(), id)
	require.NoError(t, err)
	pipe := env.pipeline(clock, "node-1/killed")
	require.NoError(t, pipe.Reconcile(t.Context()))
	lifecycleAccept(t, env, pipe, live.Group, clock, 2)
	clock.set(lifecycleBase.Add(3 * time.Second))
	mutation := lifecycleMutation(t, env, id)
	switch os.Getenv("GW_LIFECYCLE_ACTION") {
	case "disable":
		_, err = env.services.writeGroups.Disable(t.Context(), id, mutation)
	case "delete":
		_, err = env.services.writeGroups.Delete(t.Context(), id, mutation)
	case "replace":
		draft := live.Group
		draft.RowPolicy.IntervalSeconds = 20
		mutation.Group = draft
		_, err = env.services.writeGroups.Update(t.Context(), id, mutation)
		require.NoError(t, err)
		_, err = env.services.writeGroups.Apply(t.Context(), id, lifecycleMutation(t, env, id))
	default:
		t.Fatal("unknown crash action")
	}
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(filepath.Dir(path), "crash-ready"), []byte("ACK and lifecycle committed"), 0o600))
	<-t.Context().Done()
}
