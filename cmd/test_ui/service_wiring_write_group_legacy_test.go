package main

import (
	"database/sql"
	"encoding/json"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestProductionServicesLegacyWriteGuardBeforeRouterCreation(t *testing.T) {
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "configuration.db"))
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	for _, statement := range []string{
		`INSERT INTO devices (id,name,protocol,status,connection_config,readiness_status)
		 VALUES ('device-A','PLC','modbus_tcp','draft','{}','{}')`,
		`INSERT INTO points (id,device_id,name,address,function,data_type,mode,enabled)
		 VALUES ('point-A','device-A','Temperature','40001','FC03','float32','read',1)`,
		`INSERT INTO tags (id,key,key_lower,display_name,data_type,status,labels)
		 VALUES ('tag-A','temperature','temperature','Temperature','float32','active','{}')`,
		`INSERT INTO mappings (id,point_id,tag_id,transform_pipeline,status,enabled)
		 VALUES ('mapping-A','point-A','tag-A','[]','active',1)`,
	} {
		_, err := db.ExecContext(t.Context(), statement)
		require.NoError(t, err)
	}
	config, err := json.Marshal(map[string]string{"dsn": filepath.Join(t.TempDir(), "unopened-target.db")})
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), `INSERT INTO database_connectors
		(id,name,kind,connection_config,identity_revision,status,enabled)
		VALUES ('connector-A','Target','sqlite',?,'connector-1','ready',1)`, string(config))
	require.NoError(t, err)
	services := wireGatewayServices(db, connector.NewConnectionManager(connector.ConnectionManagerConfig{}))
	record, err := services.workspace.AttachDevice(t.Context(), "device-A")
	require.NoError(t, err)
	group, err := services.writeGroups.Create(t.Context(), workspace.WriteGroupMutation{
		WorkspaceID: record.ID, ExpectedWorkspaceRevision: record.DatabaseSetupRevision, ExpectedConnectorRevision: "connector-1",
		Group: &workspace.WriteGroup{
			WorkspaceID: record.ID, Name: "Temperature",
			Members:     []workspace.WriteGroupMember{{DeviceID: "device-A", PointID: "point-A", TagID: "tag-A", TargetColumn: "temperature", Required: true}},
			Destination: workspace.WriteGroupDestination{ConnectorID: "connector-A", ConnectorRevision: "connector-1", TableSchema: "main", TableName: "raw_values", StorageStrategy: workspace.WriteGroupStorageStrategyCustom},
			RowPolicy:   workspace.WriteGroupRowPolicy{IntervalSeconds: 15},
		},
	})
	require.NoError(t, err)
	_, err = services.dbMapping.Create(t.Context(), dbtarget.CreateTargetMappingRequest{
		TagID: "tag-A", ConnectorID: "connector-A", TableSchema: "main", TableName: "other_table", ColumnName: "other_column", WriteMode: "insert",
	})
	require.ErrorIs(t, err, workspace.ErrWriteGroupLegacyWriteConflict)
	current, err := services.writeGroups.Get(t.Context(), group.Group.ID)
	require.NoError(t, err)
	require.Equal(t, group, current)
	var targets, versions int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT count(*) FROM database_target_mappings`).Scan(&targets))
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT count(*) FROM write_group_versions`).Scan(&versions))
	require.Zero(t, targets)
	require.Zero(t, versions)
}
