package workspace

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestManagedSchemaScopeRejectsSQLiteNonMainSchema(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	fixture := seedWriteGroupFixture(ctx, t, db)
	target := filepath.Join(t.TempDir(), "managed-target.db")
	config, err := json.Marshal(map[string]string{"dsn": target, "schema": "other"})
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `UPDATE database_connectors SET connection_config = ? WHERE id = ?`, string(config), fixture.connectorID)
	require.NoError(t, err)

	workspaceSvc := NewService(NewSQLRepository(db))
	service := NewWriteGroupService(workspaceSvc, NewSQLWriteGroupRepository(db))
	current, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	group := newBasicWriteGroup(fixture)
	group.Destination.Database = target
	group.Destination.TableSchema = "other"
	group.Destination.StorageStrategy = WriteGroupStorageStrategyManaged
	created, err := service.Create(ctx, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedConnectorRevision: "connector-revision-1", Group: group,
	})
	require.NoError(t, err)

	_, err = service.ManagedSchemaScope(ctx, created.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: created.WorkspaceRevision,
		ExpectedGroupRevision: created.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
	})
	require.ErrorIs(t, err, ErrWriteGroupValidation)
	require.NoFileExists(t, target)
}

func TestManagedReadinessRejectsMissingOwnerMarker(t *testing.T) {
	group := &WriteGroup{
		ID: "group-1", WorkspaceID: "workspace-1",
		Destination: WriteGroupDestination{
			TableSchema: "main", TableName: "raw_values", StorageStrategy: WriteGroupStorageStrategyManaged,
		},
		Members: []WriteGroupMember{{DeviceID: "device-1", PointID: "point-1", TagID: "tag-1", Required: true}},
	}
	require.NoError(t, prepareManagedWriteGroup(group))
	layout, err := managedSchemaLayout(group, map[string]schema.DataType{"tag-1": schema.DataTypeFloat32}, writeGroupConnectorKindSQLite)
	require.NoError(t, err)

	ready, issues := evaluateWriteGroupInspection(group, map[string]schema.DataType{"tag-1": schema.DataTypeFloat32},
		managedTableInspection(layout, layout.OwnerColumn), writeGroupConnectorKindSQLite)
	require.False(t, ready)
	require.Contains(t, readinessIssueCodes(issues), "managed-schema-unverified")
}

func managedTableInspection(layout *recordingplan.GroupSchemaLayout, omitted string) *dbtarget.TableInspection {
	inspection := &dbtarget.TableInspection{
		Status: dbtarget.TableInspectionExists, Schema: "main", Table: layout.TableName,
	}
	for _, column := range layout.Columns {
		if column.Name == omitted {
			continue
		}
		inspection.Columns = append(inspection.Columns, dbtarget.ColumnInfo{
			Name: column.Name, DataType: column.SQLType, Nullable: column.Nullable, PrimaryKey: column.PrimaryKey,
		})
	}
	return inspection
}
