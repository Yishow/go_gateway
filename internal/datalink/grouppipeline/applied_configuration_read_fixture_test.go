package grouppipeline

import (
	"context"
	"database/sql"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"
)

type appliedConfigurationFixture struct {
	db     *sql.DB
	groups *workspace.WriteGroupService
	pipe   *Pipeline
	now    time.Time
	id     string
}

// Create, schema inspection, Apply and immutable-version reads are production
// code backed by two disposable SQLite files; no effective provider is mocked.
func newAppliedConfigurationFixture(t *testing.T) *appliedConfigurationFixture {
	t.Helper()
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "config.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	db.SetMaxOpenConns(1)
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	ws := workspace.NewService(workspace.NewSQLRepository(db))
	record, err := ws.GetOrCreate(t.Context())
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO devices(id,name,protocol,status,connection_config,readiness_status) VALUES('device-read','Read fixture','modbus_tcp','draft','{}','{}')`)
	require.NoError(t, err)
	_, err = ws.AttachDevice(t.Context(), "device-read")
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO points(id,device_id,name,address,function,data_type,mode,enabled) VALUES('point-read','device-read','temperature','40001','FC03','float32','read',1)`)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO tags(id,key,key_lower,display_name,data_type,status,labels) VALUES('tag-read','temperature','temperature','Temperature','float32','active','{}')`)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO mappings(id,point_id,tag_id,transform_pipeline,status,enabled) VALUES('mapping-read','point-read','tag-read','[]','active',1)`)
	require.NoError(t, err)
	targetPath := filepath.Join(t.TempDir(), "destination.db")
	raw, err := json.Marshal(map[string]string{"dsn": targetPath, "database": targetPath, "schema": "main"})
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO database_connectors(id,name,kind,connection_config,identity_revision,status,enabled) VALUES('connector-read','Read destination','sqlite',?,'connector-revision','ready',1)`, string(raw))
	require.NoError(t, err)
	f := &appliedConfigurationFixture{db: db, now: at(7)}
	f.groups = workspace.NewWriteGroupService(ws, workspace.NewSQLWriteGroupRepository(db)).WithClock(func() time.Time { return f.now }).WithBacklogOwnershipGuard(BacklogGuard{})
	connectors := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(db))
	f.groups.WithManagedTableInspector(dbtarget.NewManagedTableInspector(connectors, db))
	current, err := ws.GetOrCreate(t.Context())
	require.NoError(t, err)
	created, err := f.groups.Create(t.Context(), workspace.WriteGroupMutation{WorkspaceID: record.ID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision, ExpectedConnectorRevision: "connector-revision", Group: &workspace.WriteGroup{
		WorkspaceID: record.ID, Name: "applied A", Members: []workspace.WriteGroupMember{{DeviceID: "device-read", PointID: "point-read", TagID: "tag-read", Required: true}},
		Destination: workspace.WriteGroupDestination{ConnectorID: "connector-read", ConnectorRevision: "connector-revision", Database: targetPath, TableSchema: "main", StorageStrategy: workspace.WriteGroupStorageStrategyManaged},
		RowPolicy:   workspace.WriteGroupRowPolicy{IntervalSeconds: 15},
	}})
	require.NoError(t, err)
	f.id = created.Group.ID
	scope, err := f.groups.ManagedSchemaScope(t.Context(), f.id, f.mutation(t, nil))
	require.NoError(t, err)
	ledger := recordingplan.NewService(recordingplan.NewSQLRepository(db))
	token, err := ledger.PrepareSchemaPreview(t.Context(), scope, func(context.Context, string) (recordingplan.TargetTableInspection, error) {
		return recordingplan.TargetTableInspection{Status: "missing"}, nil
	})
	require.NoError(t, err)
	target, err := sql.Open("sqlite", targetPath)
	require.NoError(t, err)
	for _, statement := range token.Statements {
		_, err = target.Exec(statement)
		require.NoError(t, err)
	}
	require.NoError(t, target.Close())
	_, err = f.groups.ConfirmManagedSchema(t.Context(), f.id, f.mutation(t, nil), "schema-read", strings.Repeat("a", 64))
	require.NoError(t, err)
	f.pipe = New(Dependencies{Groups: f.groups, Store: groupdelivery.NewStore(db)}, Config{Now: func() time.Time { return f.now }})
	return f
}

func (f *appliedConfigurationFixture) mutation(t *testing.T, group *workspace.WriteGroup) workspace.WriteGroupMutation {
	t.Helper()
	saved, err := f.groups.Get(t.Context(), f.id)
	require.NoError(t, err)
	return workspace.WriteGroupMutation{WorkspaceID: saved.Group.WorkspaceID, ExpectedWorkspaceRevision: saved.WorkspaceRevision, ExpectedGroupRevision: saved.Group.Revision, ExpectedConnectorRevision: saved.Group.Destination.ConnectorRevision, Group: group}
}
func (f *appliedConfigurationFixture) apply(t *testing.T) *workspace.WriteGroupSaveResult {
	t.Helper()
	current, err := f.groups.Get(t.Context(), f.id)
	require.NoError(t, err)
	if current.Group.Destination.SchemaRevision == "" {
		_, err = f.groups.ConfirmManagedSchema(t.Context(), f.id, f.mutation(t, nil), "schema-"+current.Group.Revision, strings.Repeat("a", 64))
		require.NoError(t, err)
	}
	ready, err := f.groups.Readiness(t.Context(), f.id)
	require.NoError(t, err)
	require.True(t, ready.Ready, "%+v", ready.Issues)
	saved, err := f.groups.Apply(t.Context(), f.id, f.mutation(t, nil))
	require.NoError(t, err)
	return saved
}
func (f *appliedConfigurationFixture) update(t *testing.T, interval int, name string) *workspace.WriteGroupSaveResult {
	t.Helper()
	saved, err := f.groups.Get(t.Context(), f.id)
	require.NoError(t, err)
	saved.Group.RowPolicy.IntervalSeconds, saved.Group.Name = interval, name
	updated, err := f.groups.Update(t.Context(), f.id, f.mutation(t, saved.Group))
	require.NoError(t, err)
	return updated
}
func appliedConfigurationState(t *testing.T, db *sql.DB) map[string]string {
	t.Helper()
	state := map[string]string{}
	for _, table := range []string{"system_settings", "write_groups", "write_group_members", "write_group_versions", "wg_runtime_versions", "write_group_revision_ends", "wg_delivery_samples", "wg_delivery_outbox", "wg_delivery_receipts", "wg_delivery_checkpoints"} {
		rows, err := db.Query(`SELECT * FROM ` + table + ` ORDER BY rowid`)
		require.NoError(t, err)
		cols, err := rows.Columns()
		require.NoError(t, err)
		all := [][]any{}
		for rows.Next() {
			values := make([]any, len(cols))
			pointers := make([]any, len(cols))
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
		state[table] = string(raw)
	}
	return state
}
