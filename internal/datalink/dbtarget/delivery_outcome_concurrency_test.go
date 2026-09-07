package dbtarget

import (
	"context"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestWriter_RecordWriteOutcome_DoesNotOverwriteConnectorConfigOrName(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mainDB := openMigratedTestDB(t)
	connectorRepo := NewSQLConnectorRepository(mainDB)
	writer := NewWriter(connectorRepo, NewSQLTargetMappingRepository(mainDB))

	// 建立初始 connector
	initConn := &schema.DatabaseConnector{
		ID:               "conn-test-safe-outcome",
		Name:             "Initial Name",
		Kind:             schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: `{"dsn":"/tmp/initial.db"}`,
		Status:           "ready",
		Enabled:          true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	require.NoError(t, connectorRepo.Create(ctx, initConn))

	// 模擬外部流程（例如 Studio / API）已將資料庫中的 connector 改名與更新設定
	updatedFromAPI := &schema.DatabaseConnector{
		ID:               initConn.ID,
		Name:             "Renamed By User",
		Kind:             schema.DatabaseConnectorKindPostgres,
		ConnectionConfig: `{"host":"10.0.0.1"}`,
		Status:           "ready",
		Enabled:          true,
		CreatedAt:        initConn.CreatedAt,
		UpdatedAt:        time.Now(),
	}
	require.NoError(t, connectorRepo.Update(ctx, updatedFromAPI))

	// 模擬 Writer 記憶體中仍保留舊 snapshot 的 connector 物件，並記錄寫入結果
	staleMemoryConnector := &schema.DatabaseConnector{
		ID:               initConn.ID,
		Name:             "Old Stale Name",
		Kind:             schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: `{"dsn":"/tmp/stale.db"}`,
	}
	now := time.Now().UTC()
	require.NoError(t, writer.recordWriteOutcome(ctx, staleMemoryConnector, now, deliveryOutcomeSuccess, ""))

	// 驗證：資料庫中的 Name 與 Kind 與 ConnectionConfig 絕不可被舊 snapshot 覆蓋
	latestInDB, err := connectorRepo.GetByID(ctx, initConn.ID)
	require.NoError(t, err)
	require.Equal(t, "Renamed By User", latestInDB.Name)
	require.Equal(t, schema.DatabaseConnectorKindPostgres, latestInDB.Kind)
	require.Equal(t, `{"host":"10.0.0.1"}`, latestInDB.ConnectionConfig)
	require.Equal(t, deliveryOutcomeSuccess, latestInDB.LastWriteStatus)
	require.NotNil(t, latestInDB.LastWriteAt)
}
