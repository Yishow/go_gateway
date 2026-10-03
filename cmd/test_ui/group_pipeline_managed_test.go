package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
	"go-gateway/internal/testutil/lostcommit"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

type managedPipelineTarget struct {
	env                  *outageEnv
	db                   *sql.DB
	file, dsn, namespace string
	kind                 schema.DatabaseConnectorKind
}

func newManagedPipelineTarget(t *testing.T, kind string) *managedPipelineTarget {
	t.Helper()
	f := &managedPipelineTarget{kind: schema.DatabaseConnectorKindSQLite}
	setup := func(t *testing.T, env *outageEnv, name string) (string, string, string) {
		t.Helper()
		config := map[string]string{}
		namespace := "main"
		if kind == "sqlite" {
			file := filepath.Join(t.TempDir(), "managed-"+name+".db")
			config["dsn"] = file
			env.targets["connector-"+name] = file
			if name == "A" {
				f.file = file
				f.namespace = namespace
			}
		} else {
			dsn, fields := postgresEnvFields(t)
			require.Equal(t, "127.0.0.1", fields["host"])
			require.Equal(t, "55432", fields["port"])
			require.Equal(t, "gwtest", fields["dbname"])
			db, err := sql.Open("pgx", dsn)
			require.NoError(t, err)
			namespace = fmt.Sprintf("gw_managed_pipe_%s_%d", name, time.Now().UnixNano())
			_, err = db.ExecContext(t.Context(), `CREATE SCHEMA "`+namespace+`"`)
			require.NoError(t, err)
			t.Cleanup(func() {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				_, err := db.ExecContext(ctx, `DROP SCHEMA "`+namespace+`" CASCADE`)
				require.NoError(t, err)
				require.NoError(t, db.Close())
			})
			config = map[string]string{"host": fields["host"], "port": fields["port"], "user": fields["user"], "password": fields["password"], "database": fields["dbname"], "sslmode": "disable", "schema": namespace}
			if name == "A" {
				f.db = db
				f.dsn = dsn
				f.namespace = namespace
				f.kind = schema.DatabaseConnectorKindPostgres
			}
		}
		raw, err := json.Marshal(config)
		require.NoError(t, err)
		return kind, string(raw), namespace
	}
	f.env = newOutageEnvWith(t, setup)
	return f
}

func (f *managedPipelineTarget) prepare(t *testing.T, members []workspace.WriteGroupMember, policy workspace.WriteGroupRowPolicy) *workspace.WriteGroup {
	t.Helper()
	e := f.env
	record, err := e.services.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	created, err := e.services.writeGroups.Create(t.Context(), workspace.WriteGroupMutation{WorkspaceID: e.workspID, ExpectedWorkspaceRevision: record.DatabaseSetupRevision, ExpectedConnectorRevision: "connector-1", Group: &workspace.WriteGroup{
		WorkspaceID: e.workspID, Name: "Managed recording", Members: members,
		Destination: workspace.WriteGroupDestination{ConnectorID: "connector-A", ConnectorRevision: "connector-1", TableSchema: f.namespace, StorageStrategy: workspace.WriteGroupStorageStrategyManaged}, RowPolicy: policy}})
	require.NoError(t, err)
	group := created.Group
	request := map[string]any{"workspace_id": e.workspID, "expected_workspace_revision": created.WorkspaceRevision, "expected_group_revision": group.Revision, "expected_connector_revision": "connector-1"}
	path := testWriteBase + "/write-groups/" + group.ID
	status, body := httpJSON(t, e.router(), http.MethodPost, path+"/schema-preview", request)
	require.Equal(t, http.StatusOK, status, "%+v", body)
	if f.kind == schema.DatabaseConnectorKindSQLite {
		require.NoFileExists(t, f.file)
	}
	token := body["data"].(map[string]any)
	request["token"], request["operation_id"] = token["token"], token["operation_id"]
	status, body = httpJSON(t, e.router(), http.MethodPost, path+"/schema-apply", request)
	require.Equal(t, http.StatusOK, status, "%+v", body)
	require.Equal(t, "succeeded", body["data"].(map[string]any)["status"])
	if f.kind == schema.DatabaseConnectorKindSQLite {
		f.db, err = sql.Open("sqlite", f.file)
		require.NoError(t, err)
		t.Cleanup(func() { require.NoError(t, f.db.Close()) })
	}
	ready, err := e.services.writeGroups.Readiness(t.Context(), group.ID)
	require.NoError(t, err)
	require.True(t, ready.Ready, "%+v", ready.Issues)
	applied, err := e.services.writeGroups.Apply(t.Context(), group.ID, lifecycleMutation(t, e, group.ID))
	require.NoError(t, err)
	return applied.Group
}

type managedAckResolver struct {
	service *dbtarget.ConnectorService
	db      *sql.DB
	kind    schema.DatabaseConnectorKind
}

func (r managedAckResolver) Resolve(ctx context.Context, item groupdelivery.OutboxItem) (groupdelivery.Target, error) {
	opened, err := r.service.OpenDestination(ctx, item.Destination.ConnectorID, item.Destination.ConnectorRevision)
	if err != nil {
		return groupdelivery.Target{}, err
	}
	if err := opened.Close(); err != nil {
		return groupdelivery.Target{}, err
	}
	return groupdelivery.Target{DB: r.db, Kind: r.kind}, nil
}

func TestProductionManagedRowsKeepTimeIdentityAndExactTypesAfterRestart(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		t.Run(kind, func(t *testing.T) {
			f := newManagedPipelineTarget(t, kind)
			e := f.env
			types := []schema.DataType{schema.DataTypeBool, schema.DataTypeInt16, schema.DataTypeInt32, schema.DataTypeInt64, schema.DataTypeUint16, schema.DataTypeUint32, schema.DataTypeUint64, schema.DataTypeFloat32, schema.DataTypeFloat64, schema.DataTypeString}
			values := []any{true, int16(-32768), int32(-2147483648), int64(math.MinInt64), uint16(65535), uint32(math.MaxUint32), uint64(9007199254740993), float32(1.25), float64(2.5), "同名保留字 select"}
			members := make([]workspace.WriteGroupMember, 0, len(types))
			for i, typ := range types {
				id := "managed-" + strconv.Itoa(i)
				_, err := e.db.ExecContext(t.Context(), `INSERT INTO points(id,device_id,name,address,function,data_type,mode,enabled) VALUES(?,'device-1','Same label',?,'FC03',?,'read',1)`, id, strconv.Itoa(40101+i*4), string(typ))
				require.NoError(t, err)
				_, err = e.db.ExecContext(t.Context(), `INSERT INTO tags(id,key,key_lower,display_name,data_type,status,labels) VALUES(?,?,?,'同名保留字',?,'active','{}')`, id, id, id, string(typ))
				require.NoError(t, err)
				_, err = e.db.ExecContext(t.Context(), `INSERT INTO mappings(id,point_id,tag_id,transform_pipeline,status,enabled) VALUES(?,?,?,'[]','active',1)`, id, id, id)
				require.NoError(t, err)
				members = append(members, workspace.WriteGroupMember{DeviceID: "device-1", PointID: id, TagID: id, Required: true})
			}
			base := time.Now().UTC().Truncate(10 * time.Second).Add(-time.Minute)
			clock := &testClock{}
			clock.set(base.Add(-time.Second))
			e.services.writeGroups.WithClock(clock.now)
			group := f.prepare(t, members, workspace.WriteGroupRowPolicy{IntervalSeconds: 10})
			seen := map[string]bool{}
			for _, member := range group.Members {
				require.False(t, seen[member.TargetColumn])
				seen[member.TargetColumn] = true
			}
			pipe := e.pipeline(clock, "node-1/managed-before-restart")
			clock.set(base.Add(time.Second))
			require.NoError(t, pipe.Reconcile(t.Context()))
			for bucket, big := range []uint64{9007199254740993, 9223372036854775808, math.MaxUint64} {
				at := base.Add(time.Duration(bucket)*10*time.Second + 2*time.Second + 123*time.Nanosecond)
				clock.set(at)
				values[6] = big
				for i, value := range values {
					require.NoError(t, pipe.AcceptSample(t.Context(), e.envelope(group, i, fmt.Sprintf("managed-%d-%d", bucket, i), at, value)))
				}
			}
			var frozenBefore string
			require.NoError(t, e.db.QueryRowContext(t.Context(), `SELECT payload FROM wg_runtime_versions WHERE group_id=?`, group.ID).Scan(&frozenBefore))
			_, err := e.db.ExecContext(t.Context(), `UPDATE tags SET display_name='Renamed after acquisition' WHERE id LIKE 'managed-%'`)
			require.NoError(t, err)
			if kind == "sqlite" {
				require.NoError(t, f.db.Close())
				require.NoError(t, os.Rename(f.file, f.file+".offline"))
			}
			reopenLifecycleConfiguration(t, e, clock)
			recovered := e.pipeline(clock, "node-1/managed-after-restart")
			require.NoError(t, recovered.Reconcile(t.Context()))
			clock.set(base.Add(31 * time.Second))
			recovered.TickAll(t.Context())
			require.Equal(t, 3, outboxCount(t, e.db, "connector-A", "1=1"))
			var frozenAfter string
			require.NoError(t, e.db.QueryRowContext(t.Context(), `SELECT payload FROM wg_runtime_versions WHERE group_id=?`, group.ID).Scan(&frozenAfter))
			require.Equal(t, frozenBefore, frozenAfter)
			if kind == "sqlite" {
				require.NoError(t, os.Rename(f.file+".offline", f.file))
				var err error
				f.db, err = sql.Open("sqlite", f.file)
				require.NoError(t, err)
			}
			lose := &atomic.Bool{}
			lose.Store(true)
			var faultDB *sql.DB
			if kind == "sqlite" {
				faultDB = lostcommit.Open(f.file, lose)
			} else {
				faultDB = lostcommit.Wrap(stdlib.GetDefaultDriver(), f.dsn, lose)
			}
			t.Cleanup(func() { require.NoError(t, faultDB.Close()) })
			store := groupdelivery.NewStore(e.db)
			sender := groupdelivery.NewSender(store, managedAckResolver{service: e.services.dbTarget, db: faultDB, kind: f.kind}, groupdelivery.SenderConfig{
				Owner: "node-1/managed-lost-ack", Backoff: func(int) (string, time.Time) {
					return groupdelivery.StateRetrying, time.Now().UTC().Add(-time.Second)
				},
			})
			rows, err := e.db.QueryContext(t.Context(), `SELECT effect_key FROM wg_delivery_outbox ORDER BY bucket_start`)
			require.NoError(t, err)
			var effects []string
			for rows.Next() {
				var effect string
				require.NoError(t, rows.Scan(&effect))
				effects = append(effects, effect)
			}
			require.NoError(t, rows.Err())
			require.NoError(t, rows.Close())
			for _, effect := range effects {
				result, err := sender.Deliver(t.Context(), effect)
				require.NoError(t, err)
				require.Equal(t, groupdelivery.StateRetrying, result.State)
			}
			var localReceipts int
			require.NoError(t, e.db.QueryRowContext(t.Context(), `SELECT count(*) FROM wg_delivery_receipts`).Scan(&localReceipts))
			require.Zero(t, localReceipts, "a lost commit reply is not yet local proof")
			lose.Store(false)
			for _, effect := range effects {
				result, err := sender.Deliver(t.Context(), effect)
				require.NoError(t, err)
				require.Equal(t, groupdelivery.StateCommitted, result.State)
			}
			for _, effect := range effects {
				result, err := sender.Deliver(t.Context(), effect)
				require.NoError(t, err)
				require.True(t, result.Already, "confirmed replay must not send another SQL row")
			}
			assertManagedSQL(t, f, group, base)
		})
	}
}
