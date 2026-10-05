package main

import (
	"context"
	"database/sql"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/workspace"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func migratedLegacyGroup(t *testing.T, env *outageEnv) *workspace.WriteGroup {
	t.Helper()
	target, err := sql.Open("sqlite", env.targets["connector-A"])
	require.NoError(t, err)
	_, err = target.ExecContext(t.Context(), `CREATE TABLE legacy_values (temperature REAL)`)
	require.NoError(t, err)
	require.NoError(t, target.Close())
	_, err = env.db.ExecContext(t.Context(), `INSERT INTO database_target_mappings (id,tag_id,connector_id,table_schema,table_name,column_name,write_mode,write_interval_seconds,enabled) VALUES ('legacy-A','tag-t-A','connector-A','main','legacy_values','temperature','insert',10,1)`)
	require.NoError(t, err)
	record, err := env.services.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	_, err = env.services.workspace.UpdateDatabaseSetup(t.Context(), record.DatabaseSetupRevision, func(_ context.Context, _ *sql.Tx, r *workspace.Record) error {
		r.DatabaseConnectorID = "connector-A"
		return nil
	})
	require.NoError(t, err)
	preview, err := env.services.writeGroups.PreviewSingleMappingMigration(t.Context(), env.workspID, []string{"legacy-A"})
	require.NoError(t, err)
	reviewed, err := env.services.writeGroups.ReviewSingleMappingMigration(t.Context(), workspace.WriteGroupMigrationReviewRequest{WorkspaceID: env.workspID, ExpectedWorkspaceRevision: preview.WorkspaceRevision, ExpectedConnectorRevision: preview.ConnectorRevision, ReviewDigest: preview.ReviewDigest, SourceIDs: []string{"legacy-A"}, ConfirmSnapshotConversion: true})
	require.NoError(t, err)
	require.Len(t, reviewed.Groups, 1)
	return reviewed.Groups[0]
}

func legacyWriterCount(t *testing.T, env *outageEnv, connectorID string) int {
	t.Helper()
	// Same production SuppressMapping seam, backed by the production services.
	writer := dbtarget.NewWriterWithConfig(dbtarget.NewSQLConnectorRepository(env.db), dbtarget.NewSQLTargetMappingRepository(env.db), dbtarget.WriterConfig{TagReader: env.services.tag, SuppressMapping: env.services.groupPipe.Owns})
	defer func() { require.NoError(t, writer.Close(t.Context())) }()
	err := writer.WriteTagValue(t.Context(), "tag-t-A", float64(215), time.Now())
	if err != nil {
		require.ErrorIs(t, err, dbtarget.ErrOutputOwnedByWriteGroup)
	}
	target, err := sql.Open("sqlite", env.targets[connectorID])
	require.NoError(t, err)
	defer target.Close()
	var count int
	require.NoError(t, target.QueryRowContext(t.Context(), `SELECT count(*) FROM legacy_values`).Scan(&count))
	return count
}

func TestProductionLegacyTakeoverSurvivesLifecycle(t *testing.T) {
	for _, action := range []string{"disable", "delete", "destination"} {
		t.Run(action, func(t *testing.T) {
			env, clock := newOutageEnv(t), &testClock{}
			base := time.Now().UTC().Truncate(10 * time.Second).Add(-time.Minute)
			clock.set(base.Add(-time.Second))
			env.services.writeGroups.WithClock(clock.now)
			group := migratedLegacyGroup(t, env)
			require.Equal(t, 1, legacyWriterCount(t, env, "connector-A"), "draft migration does not take over")
			mutation := lifecycleMutation(t, env, group.ID)
			stale := mutation
			stale.ExpectedGroupRevision = "stale"
			_, err := env.services.writeGroups.Apply(t.Context(), group.ID, stale)
			require.Error(t, err)
			require.Equal(t, 2, legacyWriterCount(t, env, "connector-A"), "failed Apply does not take over")
			applied, err := env.services.writeGroups.Apply(t.Context(), group.ID, mutation)
			require.NoError(t, err)
			group = applied.Group
			require.Equal(t, 2, legacyWriterCount(t, env, "connector-A"), "successful Apply owns before runtime reconcile")
			clock.set(base.Add(time.Second))
			require.NoError(t, env.services.groupPipe.Reconcile(t.Context()))
			switch action {
			case "disable":
				_, err = env.services.writeGroups.Disable(t.Context(), group.ID, lifecycleMutation(t, env, group.ID))
			case "delete":
				_, err = env.services.writeGroups.Delete(t.Context(), group.ID, lifecycleMutation(t, env, group.ID))
			case "destination":
				target, e := sql.Open("sqlite", env.targets["connector-B"])
				require.NoError(t, e)
				_, e = target.ExecContext(t.Context(), `CREATE TABLE legacy_values(temperature REAL)`)
				require.NoError(t, e)
				require.NoError(t, target.Close())
				_, e = env.db.ExecContext(t.Context(), `INSERT INTO database_target_mappings (id,tag_id,connector_id,table_schema,table_name,column_name,write_mode,enabled) VALUES ('legacy-B','tag-t-A','connector-B','main','legacy_values','temperature','insert',1)`)
				require.NoError(t, e)
				updated := lifecycleMutation(t, env, group.ID)
				updated.ExpectedConnectorRevision = "connector-1"
				updated.Group = group
				updated.Group.Destination.ConnectorID = "connector-B"
				_, err = env.services.writeGroups.Update(t.Context(), group.ID, updated)
				require.NoError(t, err)
				_, err = env.services.writeGroups.Apply(t.Context(), group.ID, lifecycleMutation(t, env, group.ID))
			}
			require.NoError(t, err)
			clock.set(base.Add(11 * time.Second))
			require.NoError(t, env.services.groupPipe.Reconcile(t.Context()))
			env.services.groupPipe.TickAll(t.Context())
			require.Equal(t, 2, legacyWriterCount(t, env, "connector-A"))
			if action == "destination" {
				require.Zero(t, legacyWriterCount(t, env, "connector-B"), "current destination also taken over")
			}
			reopenLifecycleConfiguration(t, env, clock)
			require.Equal(t, 2, legacyWriterCount(t, env, "connector-A"), "restart retains original destination ownership")
			if action == "destination" {
				require.Zero(t, legacyWriterCount(t, env, "connector-B"))
			}
		})
	}
}
