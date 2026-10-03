package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

type entityTarget struct {
	env  *outageEnv
	db   *sql.DB
	name string
}

func newEntityTarget(t *testing.T, kind string) *entityTarget {
	t.Helper()
	target := &entityTarget{}
	setup := func(t *testing.T, env *outageEnv, name string) (string, string, string) {
		t.Helper()
		var db *sql.DB
		var namespace, config string
		if kind == "sqlite" {
			file := filepath.Join(t.TempDir(), "entity-"+name+".db")
			var err error
			db, err = sql.Open("sqlite", file+"?_pragma=busy_timeout(5000)")
			require.NoError(t, err)
			namespace = "main"
			payload, err := json.Marshal(map[string]string{"dsn": file})
			require.NoError(t, err)
			config = string(payload)
		} else {
			dsn, fields := postgresEnvFields(t)
			require.Equal(t, "127.0.0.1", fields["host"])
			require.Equal(t, "55432", fields["port"])
			require.Equal(t, "gwtest", fields["dbname"])
			var err error
			db, err = sql.Open("pgx", dsn)
			require.NoError(t, err)
			namespace = fmt.Sprintf("gw_entity_%s_%d", name, time.Now().UnixNano())
			_, err = db.ExecContext(t.Context(), `CREATE SCHEMA "`+namespace+`"`)
			require.NoError(t, err)
			payload, err := json.Marshal(map[string]string{"host": fields["host"], "port": fields["port"],
				"user": fields["user"], "password": fields["password"], "database": fields["dbname"], "sslmode": "disable"})
			require.NoError(t, err)
			config = string(payload)
		}
		t.Cleanup(func() {
			defer func() { _ = db.Close() }()
			if kind == "postgres" {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				_, err := db.ExecContext(ctx, `DROP SCHEMA "`+namespace+`" CASCADE`)
				require.NoError(t, err)
			}
		})
		table := `"` + namespace + `".readings`
		_, err := db.ExecContext(t.Context(), `CREATE TABLE `+table+` (
			temperature DOUBLE PRECISION, pressure BIGINT, other_temperature DOUBLE PRECISION,
			other_pressure BIGINT, entity TEXT NOT NULL, prov TEXT)`)
		require.NoError(t, err)
		dbKind := schema.DatabaseConnectorKindSQLite
		if kind == "postgres" {
			dbKind = schema.DatabaseConnectorKindPostgres
		}
		require.NoError(t, dbtarget.CreateEffectReceiptTable(t.Context(), db, dbKind, namespace))
		if name == "A" {
			target.db, target.name = db, table
		}
		return kind, config, namespace
	}
	target.env = newOutageEnvWith(t, setup)
	target.env.dedupe = "receipt"
	return target
}

func (e *entityTarget) create(t *testing.T, mutate func(*workspace.WriteGroup)) *workspace.WriteGroup {
	t.Helper()
	env := e.env
	group := &workspace.WriteGroup{WorkspaceID: env.workspID, Name: "One group, two entities",
		Destination: workspace.WriteGroupDestination{ConnectorID: "connector-A", ConnectorRevision: "connector-1",
			TableSchema: env.schemas["A"], TableName: "readings", StorageStrategy: workspace.WriteGroupStorageStrategyCustom},
		RowPolicy:   workspace.WriteGroupRowPolicy{IntervalSeconds: 10, EntityKeyColumn: "entity", ProvenanceColumn: "prov"},
		WritePolicy: workspace.WriteGroupWritePolicy{DedupeCapability: "receipt"}}
	for _, entity := range []string{"A", "B"} {
		for _, member := range []struct{ prefix, column string }{{"t", "temperature"}, {"p", "pressure"}} {
			group.Members = append(group.Members, workspace.WriteGroupMember{DeviceID: "device-1", EntityKey: entity,
				PointID: "point-" + member.prefix + "-" + entity, TagID: "tag-" + member.prefix + "-" + entity,
				TargetColumn: member.column, Required: true})
		}
	}
	if mutate != nil {
		mutate(group)
	}
	current, err := env.services.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	created, err := env.services.writeGroups.Create(t.Context(), workspace.WriteGroupMutation{
		WorkspaceID: env.workspID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedConnectorRevision: "connector-1", Group: group})
	require.NoError(t, err)
	return created.Group
}

func TestProductionEntityLayoutReadinessAndApplyAgree(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		for _, collision := range []bool{false, true} {
			name := "distinct-rows"
			if collision {
				name = "same-row-case-variant"
			}
			t.Run(kind+"/"+name, func(t *testing.T) {
				target := newEntityTarget(t, kind)
				group := target.create(t, func(g *workspace.WriteGroup) {
					if collision {
						g.Members = []workspace.WriteGroupMember{g.Members[0], g.Members[2]}
						g.Members[1].EntityKey = "A"
						g.Members[1].TargetColumn = "Temperature"
					}
				})
				readiness, err := target.env.services.writeGroups.Readiness(t.Context(), group.ID)
				require.NoError(t, err)
				require.Equal(t, !collision, readiness.Ready, "%+v", readiness.Issues)
				_, err = target.env.services.writeGroups.Apply(t.Context(), group.ID, lifecycleMutation(t, target.env, group.ID))
				if collision {
					require.NotEmpty(t, readiness.Issues)
					require.Contains(t, readiness.Issues[0].Message, "same entity row", "readiness must explain how to repair the collision, not call it an SQL type mismatch")
					require.Error(t, err)
					live, getErr := target.env.services.writeGroups.Get(t.Context(), group.ID)
					require.NoError(t, getErr)
					require.Empty(t, live.Group.AppliedRevision)
				} else {
					require.NoError(t, err)
					require.Len(t, group.Members, 4)
				}
			})
		}
	}
}
