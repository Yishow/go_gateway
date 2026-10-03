package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go-gateway/internal/datalink/dbtarget"

	"github.com/stretchr/testify/require"
)

type revisionManagedReadinessInspector struct {
	delegate *dbtarget.ManagedTableInspector
	calls    int
}

func (i *revisionManagedReadinessInspector) InspectTableAtRevision(ctx context.Context, connectorID, expectedRevision, schemaName, tableName string) (*dbtarget.TableInspection, error) {
	i.calls++
	return i.delegate.InspectTableAtRevision(ctx, connectorID, expectedRevision, schemaName, tableName)
}

func newManagedReadinessDatabase(t *testing.T) (*sql.DB, writeGroupFixture, string) {
	t.Helper()
	db := openWorkspaceTestDB(t, filepath.Join(t.TempDir(), "workspace.db"))
	fixture := seedWriteGroupFixture(t.Context(), t, db)
	var sequence int
	var databaseName, internalPath string
	require.NoError(t, db.QueryRowContext(t.Context(), `PRAGMA database_list`).Scan(&sequence, &databaseName, &internalPath))
	require.NotEmpty(t, internalPath)
	return db, fixture, internalPath
}

func setSQLiteConnectorConfig(t *testing.T, db *sql.DB, connectorID, dsn string) {
	t.Helper()
	config, err := json.Marshal(map[string]string{
		"dsn": dsn, "database": dsn, "schema": "main",
	})
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), `UPDATE database_connectors SET connection_config = ? WHERE id = ?`, string(config), connectorID)
	require.NoError(t, err)
}

func saveManagedReadinessProof(t *testing.T, db *sql.DB, groupID string) {
	t.Helper()
	_, err := db.ExecContext(t.Context(), `UPDATE write_groups SET destination_schema_revision = ?, destination_schema_digest = ? WHERE id = ?`,
		"schema-proof-1", strings.Repeat("a", 64), groupID)
	require.NoError(t, err)
}

func TestManagedReadinessUsesRevisionGuardBeforeGenericInspector(t *testing.T) {
	ctx := t.Context()
	db, fixture, internalPath := newManagedReadinessDatabase(t)
	defer db.Close()
	setSQLiteConnectorConfig(t, db, fixture.connectorID, internalPath)

	workspaceSvc := NewService(NewSQLRepository(db))
	service := NewWriteGroupService(workspaceSvc, NewSQLWriteGroupRepository(db))
	current, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	group := newBasicWriteGroup(fixture)
	group.Destination.Database = internalPath
	group.Destination.TableSchema = "main"
	group.Destination.TableName = "raw_values"
	group.Destination.StorageStrategy = WriteGroupStorageStrategyManaged
	group.RowPolicy.IntervalSeconds = 15
	created, err := service.Create(ctx, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedConnectorRevision: "connector-revision-1", Group: group,
	})
	require.NoError(t, err)
	saveManagedReadinessProof(t, db, created.Group.ID)

	generic := &writeGroupReadinessInspector{inspection: &dbtarget.TableInspection{
		Status: dbtarget.TableInspectionExists, Schema: "main", Table: "raw_values",
	}}
	connectorService := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(db))
	managed := &revisionManagedReadinessInspector{
		delegate: dbtarget.NewManagedTableInspector(connectorService, db),
	}
	service.WithTableInspector(generic).WithManagedTableInspector(managed)

	readiness, err := service.Readiness(ctx, created.Group.ID)
	require.NoError(t, err)
	require.True(t, readiness.ConfigReady)
	require.False(t, readiness.SchemaReady)
	require.False(t, readiness.Ready)
	require.Contains(t, readinessIssueCodes(readiness.Issues), "schema-inspection-failed")
	require.Zero(t, generic.calls)
	require.Equal(t, 1, managed.calls)
}

func TestManagedReadinessAcceptsOpaqueSQLiteURIForMissingTarget(t *testing.T) {
	ctx := t.Context()
	db, fixture, _ := newManagedReadinessDatabase(t)
	defer db.Close()
	target := filepath.Join(t.TempDir(), "missing target.db")
	opaqueURI := (&url.URL{Scheme: "file", Opaque: url.PathEscape(target)}).String()
	setSQLiteConnectorConfig(t, db, fixture.connectorID, opaqueURI)

	workspaceSvc := NewService(NewSQLRepository(db))
	service := NewWriteGroupService(workspaceSvc, NewSQLWriteGroupRepository(db))
	current, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	group := newBasicWriteGroup(fixture)
	group.Destination.Database = opaqueURI
	group.Destination.TableSchema = "main"
	group.Destination.TableName = "raw_values"
	group.Destination.StorageStrategy = WriteGroupStorageStrategyManaged
	group.RowPolicy.IntervalSeconds = 15
	created, err := service.Create(ctx, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedConnectorRevision: "connector-revision-1", Group: group,
	})
	require.NoError(t, err)
	saveManagedReadinessProof(t, db, created.Group.ID)

	generic := &writeGroupReadinessInspector{inspection: &dbtarget.TableInspection{
		Status: dbtarget.TableInspectionExists, Schema: "main", Table: "raw_values",
	}}
	connectorService := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(db))
	managed := &revisionManagedReadinessInspector{
		delegate: dbtarget.NewManagedTableInspector(connectorService, db),
	}
	service.WithTableInspector(generic).WithManagedTableInspector(managed)

	readiness, err := service.Readiness(ctx, created.Group.ID)
	require.NoError(t, err)
	require.True(t, readiness.ConfigReady)
	require.False(t, readiness.SchemaReady)
	require.False(t, readiness.Ready)
	require.Contains(t, readinessIssueCodes(readiness.Issues), "schema-target-missing")
	require.NotContains(t, readinessIssueCodes(readiness.Issues), "schema-scope-unverified")
	require.Zero(t, generic.calls)
	require.Equal(t, 1, managed.calls)
	_, err = os.Stat(target)
	require.ErrorIs(t, err, os.ErrNotExist)
}
