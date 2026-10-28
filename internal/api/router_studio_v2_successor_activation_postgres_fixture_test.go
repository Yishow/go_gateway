package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"
)

// Own two unique schemas; never touch public or any pre-existing schema.
func newPostgresSuccessorDestination(t *testing.T, timescale bool) (managedGroupFixture, string) {
	t.Helper()
	key, database := "POSTGRES_DSN", "gateway_test"
	if timescale {
		key, database = "TIMESCALE_POSTGRES_DSN", "gateway_timescale_test"
	}
	dsn := os.Getenv(key)
	if dsn == "" {
		t.Skip("owned loopback PostgreSQL fixture not configured")
	}
	fields := map[string]string{}
	for _, part := range strings.Fields(dsn) {
		k, v, ok := strings.Cut(part, "=")
		if ok {
			fields[k] = v
		}
	}
	require.Equal(t, "127.0.0.1", fields["host"])
	require.Equal(t, "55439", fields["port"])
	require.Equal(t, database, fields["dbname"])
	require.Equal(t, "gateway_test", fields["user"])
	target, err := sql.Open("pgx", dsn)
	require.NoError(t, err)
	require.NoError(t, target.PingContext(t.Context()))
	a := fmt.Sprintf("gw_successor_%d_a", time.Now().UnixNano())
	b := strings.TrimSuffix(a, "_a") + "_b"
	t.Cleanup(func() { require.NoError(t, target.Close()) })
	for _, name := range []string{a, b} {
		_, err = target.Exec(`CREATE SCHEMA "` + name + `"`)
		require.NoError(t, err)
		t.Cleanup(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			_, err := target.ExecContext(ctx, `DROP SCHEMA "`+name+`" CASCADE`)
			require.NoError(t, err)
			var n int
			require.NoError(t, target.QueryRowContext(ctx, `SELECT count(*) FROM pg_namespace WHERE nspname=$1`, name).Scan(&n))
			require.Zero(t, n, "owned schema cleanup")
		})
	}
	f := managedGroupFixture{writeGroupRouterFixture: newWriteGroupRouterFixture(t), namespace: a, target: target}
	for i, name := range []string{a, b} {
		cfg := map[string]string{"host": fields["host"], "port": fields["port"], "user": fields["user"], "database": database, "sslmode": "disable", "schema": name}
		raw, err := json.Marshal(cfg)
		require.NoError(t, err)
		if i == 0 {
			_, err = f.db.Exec(`UPDATE database_connectors SET kind='postgres',connection_config=? WHERE id='connector-A'`, string(raw))
		} else {
			_, err = f.db.Exec(`INSERT INTO database_connectors(id,name,kind,connection_config,identity_revision,status,enabled) VALUES('connector-B','Other','postgres',?,'connector-2','ready',1)`, string(raw))
		}
		require.NoError(t, err)
	}
	f.record, err = f.workspace.BindDatabaseConnector(t.Context(), "connector-A")
	require.NoError(t, err)
	_, err = f.db.Exec(`UPDATE tags SET data_type='uint64' WHERE id='tag-A'`)
	require.NoError(t, err)
	targets := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(f.db))
	f.groups = workspace.NewWriteGroupService(f.workspace, workspace.NewSQLWriteGroupRepository(f.db)).WithTableInspector(dbtarget.NewReadOnlyTableInspector(targets)).WithBacklogOwnershipGuard(grouppipeline.BacklogGuard{})
	f.groups.WithManagedTableInspector(f.groups.ManagedTableInspector(targets))
	f.router = NewRouter(&DatalinkServices{Workspace: f.workspace, WriteGroups: f.groups, DBTarget: targets, RecordingPlan: recordingplan.NewService(recordingplan.NewSQLRepository(f.db))})
	return f, b
}

func (r *postgresSuccessorRig) activationAPI() successorActivationAPIFixture {
	return successorActivationAPIFixture{successorAPIFixture: r.successorAPIFixture, reviewResult: r.preparationReview, prepared: r.prepared, boundary: r.boundary}
}

// Snapshot only the fixture-owned schemas, including marker and receipt rows.
func (r *postgresSuccessorRig) targetState(t *testing.T) map[string]string {
	t.Helper()
	rows, err := r.target.Query(`SELECT schemaname,tablename FROM pg_tables WHERE schemaname IN ($1,$2) ORDER BY schemaname,tablename`, r.f.namespace, r.targetNamespace)
	require.NoError(t, err)
	var names [][2]string
	for rows.Next() {
		var pair [2]string
		require.NoError(t, rows.Scan(&pair[0], &pair[1]))
		names = append(names, pair)
	}
	require.NoError(t, rows.Err())
	require.NoError(t, rows.Close())
	result := map[string]string{}
	for _, pair := range names {
		var body string
		require.NoError(t, r.target.QueryRow(`SELECT COALESCE(jsonb_agg(rowdata ORDER BY rowdata::text),'[]'::jsonb)::text FROM (SELECT to_jsonb(t) AS rowdata FROM "`+pair[0]+`"."`+pair[1]+`" t) q`).Scan(&body))
		result[pair[0]+"."+pair[1]] = body
	}
	return result
}
