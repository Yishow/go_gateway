package workspace

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSQLServerPartitionMaintenanceReadinessIsSQLScopedAndVisible(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	service := &WriteGroupService{repo: NewSQLWriteGroupRepository(db)}
	group := &WriteGroup{
		WorkspaceID: "workspace", StorageMode: WriteGroupStorageModeNativeMonthly,
		Destination: WriteGroupDestination{ConnectorID: "connector", Database: "recordings", TableSchema: "dbo"},
	}

	sqlResult := &WriteGroupReadiness{}
	service.populateSQLServerPartitionMaintenanceReadiness(t.Context(), writeGroupConnectorKindSQLServer, group, sqlResult)
	require.NotNil(t, sqlResult.PartitionMaintenance)
	require.Equal(t, "blocked", sqlResult.PartitionMaintenance.Status)
	require.Equal(t, partitionMaintenanceBlockedCode, sqlResult.PartitionMaintenance.Code)
	require.Equal(t, "state_unreadable", sqlResult.PartitionMaintenance.Reason)

	mysqlResult := &WriteGroupReadiness{}
	service.populateSQLServerPartitionMaintenanceReadiness(t.Context(), writeGroupConnectorKindMySQL, group, mysqlResult)
	require.Nil(t, mysqlResult.PartitionMaintenance, "MySQL readiness must not be changed by the SQL Server sidecar")
}
