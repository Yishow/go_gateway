package grouppipeline

import (
	"database/sql"
	"encoding/json"
	"testing"

	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestBacklogGuardRejectsUnresolvedImmutableOwnership(t *testing.T) {
	for _, kind := range []string{"outbox", "sample"} {
		for _, damage := range []string{"missing", "empty", "malformed", "workspace", "group", "revision", "applied", "destination", "blank_destination", "blank_connector", "blank_table"} {
			t.Run(kind+"/"+damage, func(t *testing.T) {
				db, group := backlogOwnershipFixture(t, kind)
				version := *group
				switch damage {
				case "missing":
					_, err := db.Exec(`UPDATE ` + backlogTable(kind) + ` SET group_revision='R0'`)
					require.NoError(t, err)
				case "workspace":
					version.WorkspaceID = "foreign"
				case "group":
					version.ID = "foreign"
				case "revision":
					version.Revision = "foreign"
				case "applied":
					version.AppliedRevision = "foreign"
				case "destination":
					version.Destination.ConnectorRevision = ""
				case "blank_destination":
					version.Destination.ConnectorRevision = " "
				case "blank_connector":
					version.Destination.ConnectorID = " "
				case "blank_table":
					version.Destination.TableName = " "
				}
				raw, err := json.Marshal(version)
				require.NoError(t, err)
				if damage == "empty" {
					raw = []byte(`{}`)
				} else if damage == "malformed" {
					raw = []byte(`{`)
				}
				_, err = db.Exec(`UPDATE write_group_versions SET payload=?`, string(raw))
				require.NoError(t, err)
				before := backlogOwnershipState(t, db)
				require.ErrorIs(t, checkBacklogOwnership(t, db, group), ErrBacklogOwnershipUnproven)
				require.Equal(t, before, backlogOwnershipState(t, db), "failed proof is read-only")
			})
		}
	}
}

func TestBacklogGuardRejectsFrozenDestinationMismatch(t *testing.T) {
	for _, column := range []string{"workspace_id", "connector_id", "connector_revision", "database_name", "table_schema", "table_name", "destination_scope", "dedupe_capability"} {
		t.Run(column, func(t *testing.T) {
			db, group := backlogOwnershipFixture(t, "outbox")
			_, err := db.Exec(`UPDATE wg_delivery_outbox SET ` + column + `='foreign'`)
			require.NoError(t, err)
			before := backlogOwnershipState(t, db)
			require.ErrorIs(t, checkBacklogOwnership(t, db, group), ErrBacklogOwnershipUnproven)
			require.Equal(t, before, backlogOwnershipState(t, db))
		})
	}
}

func TestBacklogGuardRejectsSampleWorkspaceMismatch(t *testing.T) {
	db, group := backlogOwnershipFixture(t, "sample")
	_, err := db.Exec(`UPDATE wg_delivery_samples SET workspace_id='foreign'`)
	require.NoError(t, err)
	require.ErrorIs(t, checkBacklogOwnership(t, db, group), ErrBacklogOwnershipUnproven)
}

func TestBacklogGuardAndImmutableReaderRequireScopedTransaction(t *testing.T) {
	db, group := backlogOwnershipFixture(t, "outbox")
	tx, err := db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	for _, invalid := range []*workspace.WriteGroup{nil, {ID: group.ID}, {WorkspaceID: group.WorkspaceID}} {
		require.ErrorIs(t, (BacklogGuard{}).CheckWriteGroupBacklog(t.Context(), tx, invalid), ErrBacklogOwnershipUnproven)
	}
	require.ErrorIs(t, (BacklogGuard{}).CheckWriteGroupBacklog(t.Context(), nil, group), ErrBacklogOwnershipUnproven)
	repo := workspace.NewSQLWriteGroupRepository(nil)
	for _, scope := range [][3]string{{"", "G1", "R1"}, {group.WorkspaceID, "", "R1"}, {group.WorkspaceID, "G1", ""}} {
		_, err := repo.VersionByRevisionInTx(t.Context(), tx, scope[0], scope[1], scope[2])
		require.Error(t, err)
	}
	_, err = repo.VersionByRevisionInTx(t.Context(), nil, group.WorkspaceID, group.ID, group.Revision)
	require.ErrorIs(t, err, workspace.ErrWriteGroupValidation)
	var absent *workspace.SQLWriteGroupRepository
	_, err = absent.VersionByRevisionInTx(t.Context(), tx, group.WorkspaceID, group.ID, group.Revision)
	require.ErrorIs(t, err, workspace.ErrWriteGroupValidation)
}

func TestBacklogGuardPreservesHistoricalDestinationWithoutLiveConnector(t *testing.T) {
	for _, kind := range []string{"outbox", "sample"} {
		t.Run(kind, func(t *testing.T) {
			db, group := backlogOwnershipFixture(t, kind)
			// R1/D1 remains the owner even though the saved draft now names R2/D2.
			// This legacy revision deliberately has no runtime descriptor or connector.
			group.Revision, group.AppliedRevision = "R2", "R2"
			group.Destination.ConnectorRevision, group.Destination.TableName = "C2", "D2"
			_, err := db.Exec(`UPDATE write_groups SET revision='R2', applied_revision='R2', destination_connector_revision='C2', destination_table_name='D2'`)
			require.NoError(t, err)
			before := backlogOwnershipState(t, db)
			require.NoError(t, checkBacklogOwnership(t, db, group))
			require.Equal(t, before, backlogOwnershipState(t, db))
		})
	}
}

func TestBacklogGuardChecksPendingStatesButRetainsFinalizedHistory(t *testing.T) {
	for _, state := range []string{"pending", "retrying", "blocked", "unknown", "sql_committed", "operator_skipped"} {
		t.Run(state, func(t *testing.T) {
			db, group := backlogOwnershipFixture(t, "outbox")
			_, err := db.Exec(`UPDATE wg_delivery_outbox SET group_revision='R0', state=?`, state)
			require.NoError(t, err)
			before := backlogOwnershipState(t, db)
			err = checkBacklogOwnership(t, db, group)
			if state == "sql_committed" || state == "operator_skipped" {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, ErrBacklogOwnershipUnproven)
			}
			require.Equal(t, before, backlogOwnershipState(t, db))
		})
	}
}

func backlogOwnershipFixture(t *testing.T, kind string) (*sql.DB, *workspace.WriteGroup) {
	t.Helper()
	f := newFixture(t)
	db := f.store.DB()
	db.SetMaxOpenConns(1)
	group := groupAt("R1")
	group.ID, group.Destination.ConnectorRevision, group.Destination.TableName = "G1", "C1", "D1"
	tx, err := db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	persistOwnershipVersion(t, tx, snapshotOf(group, at(0)))
	require.NoError(t, tx.Commit())
	if kind == "sample" {
		_, err = db.Exec(`INSERT INTO wg_delivery_samples(workspace_id,group_id,group_revision,sample_id,member_key,observed_at,bucket_start,payload,payload_digest,created_at)
		 VALUES(?,'G1','R1','Q1','member','2026-01-01T00:00:01Z','2026-01-01T00:00:00Z','accepted sample bytes','sample-digest','2026-01-01T00:00:01Z')`, group.WorkspaceID)
	} else {
		scope, scopeErr := snapshot.DestinationScope(group.Destination.ConnectorID, group.Destination.ConnectorRevision, group.Destination.Database, group.Destination.TableSchema, group.Destination.TableName)
		require.NoError(t, scopeErr)
		_, err = db.Exec(`INSERT INTO wg_delivery_outbox(effect_key,record_id,workspace_id,group_id,group_revision,bucket_start,partition_key,destination_scope,connector_id,connector_revision,database_name,table_schema,table_name,payload,payload_digest,state,next_retry_at,created_at,updated_at)
		 VALUES('Q1','record-1',?,'G1','R1','2026-01-01T00:00:00Z','partition',?,?,'C1',?,?,'D1','accepted row bytes','row-digest','pending','2026-01-01T00:00:00Z','created','updated')`, group.WorkspaceID, scope, group.Destination.ConnectorID, group.Destination.Database, group.Destination.TableSchema)
	}
	require.NoError(t, err)
	return db, group
}

func checkBacklogOwnership(t *testing.T, db *sql.DB, group *workspace.WriteGroup) error {
	t.Helper()
	tx, err := db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	return (BacklogGuard{}).CheckWriteGroupBacklog(t.Context(), tx, group)
}

func backlogTable(kind string) string {
	if kind == "sample" {
		return "wg_delivery_samples"
	}
	return "wg_delivery_outbox"
}

func backlogOwnershipState(t *testing.T, db *sql.DB) map[string]string {
	t.Helper()
	result := map[string]string{}
	for _, table := range []string{"system_settings", "write_groups", "write_group_members", "write_group_versions", "wg_runtime_versions", "wg_delivery_samples", "wg_delivery_outbox", "wg_delivery_receipts"} {
		rows, err := db.Query(`SELECT * FROM ` + table + ` ORDER BY rowid`)
		require.NoError(t, err)
		columns, err := rows.Columns()
		require.NoError(t, err)
		all := [][]any{}
		for rows.Next() {
			values, pointers := make([]any, len(columns)), make([]any, len(columns))
			for i := range values {
				pointers[i] = &values[i]
			}
			require.NoError(t, rows.Scan(pointers...))
			all = append(all, values)
		}
		require.NoError(t, rows.Err())
		require.NoError(t, rows.Close())
		raw, err := json.Marshal(all)
		require.NoError(t, err)
		result[table] = string(raw)
	}
	return result
}
