// Package device 提供 SQL Repository 的單元測試。
package device

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

/**
 * setupTestDB 建立測試用的記憶體資料庫並執行 migration
 * @param t 測試實例
 * @returns *sql.DB 資料庫連線
 */
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)

	// 讀取並執行 migration 腳本
	migrationPath := "../schema/migrations/001_initial_schema_sqlite.sql"
	migrationContent, err := os.ReadFile(migrationPath)
	require.NoError(t, err)

	_, err = db.Exec(string(migrationContent))
	require.NoError(t, err)

	return db
}

/**
 * createTestDevice 建立測試用的設備資料
 * @param t 測試實例
 * @returns *schema.Device 測試設備物件
 */
func createTestDevice(t *testing.T) *schema.Device {
	now := time.Now().UTC()
	return &schema.Device{
		ID:               "test-device-001",
		Name:             "測試設備",
		Description:      "這是一個測試設備",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusActive,
		ConnectionConfig: `{"host":"127.0.0.1","port":502,"slave_id":1}`,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

/**
 * TestSQLRepository_Create 測試建立設備功能
 */
func TestSQLRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	device := createTestDevice(t)

	// 執行 Create
	err := repo.Create(ctx, device)
	assert.NoError(t, err)

	// 驗證設備已建立
	retrieved, err := repo.GetByID(ctx, device.ID)
	assert.NoError(t, err)
	assert.Equal(t, device.ID, retrieved.ID)
	assert.Equal(t, device.Name, retrieved.Name)
	assert.Equal(t, device.Protocol, retrieved.Protocol)
	assert.Equal(t, device.Status, retrieved.Status)
}

/**
 * TestSQLRepository_Create_Duplicate 測試建立重複設備應失敗
 */
func TestSQLRepository_Create_Duplicate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	device := createTestDevice(t)

	// 第一次建立應成功
	err := repo.Create(ctx, device)
	assert.NoError(t, err)

	// 第二次建立相同 ID 應失敗
	err = repo.Create(ctx, device)
	assert.Error(t, err)
}

/**
 * TestSQLRepository_Update 測試更新設備功能
 */
func TestSQLRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	device := createTestDevice(t)
	err := repo.Create(ctx, device)
	require.NoError(t, err)

	// 更新設備資料
	device.Name = "更新後的設備名稱"
	device.Description = "更新後的描述"
	device.Status = schema.DeviceStatusDisabled
	device.ConnectionConfig = `{"host":"192.168.1.100","port":502,"slave_id":2}`

	err = repo.Update(ctx, device)
	assert.NoError(t, err)

	// 驗證更新結果
	retrieved, err := repo.GetByID(ctx, device.ID)
	assert.NoError(t, err)
	assert.Equal(t, "更新後的設備名稱", retrieved.Name)
	assert.Equal(t, "更新後的描述", retrieved.Description)
	assert.Equal(t, schema.DeviceStatusDisabled, retrieved.Status)
	assert.Contains(t, retrieved.ConnectionConfig, "192.168.1.100")
}

/**
 * TestSQLRepository_Update_NotFound 測試更新不存在的設備應失敗
 */
func TestSQLRepository_Update_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	device := createTestDevice(t)
	device.ID = "non-existent-id"

	err := repo.Update(ctx, device)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "設備不存在")
}

/**
 * TestSQLRepository_Delete 測試刪除設備功能
 */
func TestSQLRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	device := createTestDevice(t)
	err := repo.Create(ctx, device)
	require.NoError(t, err)

	// 刪除設備
	err = repo.Delete(ctx, device.ID)
	assert.NoError(t, err)

	// 驗證設備已刪除
	_, err = repo.GetByID(ctx, device.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "設備不存在")
}

/**
 * TestSQLRepository_Delete_NotFound 測試刪除不存在的設備應失敗
 */
func TestSQLRepository_Delete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	err := repo.Delete(ctx, "non-existent-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "設備不存在")
}

/**
 * TestSQLRepository_GetByID 測試根據 ID 取得設備
 */
func TestSQLRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	device := createTestDevice(t)
	err := repo.Create(ctx, device)
	require.NoError(t, err)

	// 取得設備
	retrieved, err := repo.GetByID(ctx, device.ID)
	assert.NoError(t, err)
	assert.Equal(t, device.ID, retrieved.ID)
	assert.Equal(t, device.Name, retrieved.Name)
}

/**
 * TestSQLRepository_GetByID_NotFound 測試取得不存在的設備
 */
func TestSQLRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	_, err := repo.GetByID(ctx, "non-existent-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "設備不存在")
}

/**
 * TestSQLRepository_List 測試列出設備
 */
func TestSQLRepository_List(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	// 建立多個設備
	devices := []*schema.Device{
		createTestDevice(t),
		{
			ID:               "test-device-002",
			Name:             "測試設備 2",
			Protocol:         schema.ProtocolModbusTCP,
			Status:           schema.DeviceStatusDraft,
			ConnectionConfig: "{}",
			CreatedAt:        time.Now().UTC().Add(time.Second),
			UpdatedAt:        time.Now().UTC().Add(time.Second),
		},
		{
			ID:               "test-device-003",
			Name:             "測試設備 3",
			Protocol:         schema.ProtocolFatekFBs,
			Status:           schema.DeviceStatusActive,
			ConnectionConfig: "{}",
			CreatedAt:        time.Now().UTC().Add(2 * time.Second),
			UpdatedAt:        time.Now().UTC().Add(2 * time.Second),
		},
	}

	for _, device := range devices {
		err := repo.Create(ctx, device)
		require.NoError(t, err)
	}

	// 列出所有設備
	result, err := repo.List(ctx, ListFilter{})
	assert.NoError(t, err)
	assert.Len(t, result, 3)

	// 驗證排序（按 created_at DESC）
	assert.Equal(t, "test-device-003", result[0].ID)
	assert.Equal(t, "test-device-002", result[1].ID)
	assert.Equal(t, "test-device-001", result[2].ID)
}

/**
 * TestSQLRepository_List_WithFilter 測試使用過濾條件列出設備
 */
func TestSQLRepository_List_WithFilter(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	// 建立不同協議和狀態的設備
	devices := []*schema.Device{
		{
			ID:               "test-modbus-active",
			Name:             "Modbus Active",
			Protocol:         schema.ProtocolModbusTCP,
			Status:           schema.DeviceStatusActive,
			ConnectionConfig: "{}",
			CreatedAt:        time.Now().UTC(),
			UpdatedAt:        time.Now().UTC(),
		},
		{
			ID:               "test-modbus-draft",
			Name:             "Modbus Draft",
			Protocol:         schema.ProtocolModbusTCP,
			Status:           schema.DeviceStatusDraft,
			ConnectionConfig: "{}",
			CreatedAt:        time.Now().UTC(),
			UpdatedAt:        time.Now().UTC(),
		},
		{
			ID:               "test-fatek-active",
			Name:             "Fatek Active",
			Protocol:         schema.ProtocolFatekFBs,
			Status:           schema.DeviceStatusActive,
			ConnectionConfig: "{}",
			CreatedAt:        time.Now().UTC(),
			UpdatedAt:        time.Now().UTC(),
		},
	}

	for _, device := range devices {
		err := repo.Create(ctx, device)
		require.NoError(t, err)
	}

	// 測試協議過濾
	protocol := schema.ProtocolModbusTCP
	result, err := repo.List(ctx, ListFilter{Protocol: &protocol})
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	// 測試狀態過濾
	status := schema.DeviceStatusActive
	result, err = repo.List(ctx, ListFilter{Status: &status})
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	// 測試組合過濾
	result, err = repo.List(ctx, ListFilter{
		Protocol: &protocol,
		Status:   &status,
	})
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "test-modbus-active", result[0].ID)
}

/**
 * TestSQLRepository_List_WithPagination 測試分頁功能
 */
func TestSQLRepository_List_WithPagination(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	// 建立 5 個設備
	for i := 1; i <= 5; i++ {
		device := createTestDevice(t)
		device.ID = string(rune('0' + i))
		device.Name = string(rune('0' + i))
		err := repo.Create(ctx, device)
		require.NoError(t, err)
	}

	// 測試 limit
	result, err := repo.List(ctx, ListFilter{Limit: 3})
	assert.NoError(t, err)
	assert.Len(t, result, 3)

	// 測試 offset
	result, err = repo.List(ctx, ListFilter{
		Limit:  2,
		Offset: 2,
	})
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

/**
 * TestSQLRepository_Count 測試計算設備數量
 */
func TestSQLRepository_Count(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	// 初始計數應為 0
	count, err := repo.Count(ctx)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)

	// 建立 3 個設備
	for i := 0; i < 3; i++ {
		device := createTestDevice(t)
		device.ID = string(rune('0' + i))
		err := repo.Create(ctx, device)
		require.NoError(t, err)
	}

	// 計數應為 3
	count, err = repo.Count(ctx)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), count)
}

/**
 * TestSQLRepository_NullFields 測試處理 NULL 欄位
 */
func TestSQLRepository_NullFields(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	now := time.Now().UTC()
	device := &schema.Device{
		ID:               "test-null-fields",
		Name:             "測試 NULL 欄位",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusActive,
		ConnectionConfig: "{}",
		// Description 為空字串
		// LastTestAt, LastTestSuccess, LastTestError 為 nil
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := repo.Create(ctx, device)
	require.NoError(t, err)

	// 取回並驗證 NULL 欄位正確處理
	retrieved, err := repo.GetByID(ctx, device.ID)
	assert.NoError(t, err)
	assert.Equal(t, "", retrieved.Description)
	assert.Nil(t, retrieved.LastTestAt)
	assert.Nil(t, retrieved.LastTestSuccess)
	assert.Equal(t, "", retrieved.LastTestError)
}

/**
 * TestSQLRepository_UpdateLastTest 測試更新測試相關欄位
 */
func TestSQLRepository_UpdateLastTest(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	device := createTestDevice(t)
	err := repo.Create(ctx, device)
	require.NoError(t, err)

	// 更新測試結果（透過 Update）
	now := time.Now().UTC()
	success := true
	testError := ""

	device.LastTestAt = &now
	device.LastTestSuccess = &success
	device.LastTestError = testError

	err = repo.Update(ctx, device)
	assert.NoError(t, err)

	// 驗證測試結果已更新
	retrieved, err := repo.GetByID(ctx, device.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved.LastTestAt)
	assert.NotNil(t, retrieved.LastTestSuccess)
	assert.True(t, *retrieved.LastTestSuccess)
	assert.Equal(t, testError, retrieved.LastTestError)
}

/**
 * TestSQLRepository_JSONConfig 測試處理 JSON 格式的 ConnectionConfig
 */
func TestSQLRepository_JSONConfig(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	config := `{"host":"192.168.1.1","port":502,"slave_id":1,"timeout":10}`
	device := createTestDevice(t)
	device.ConnectionConfig = config

	err := repo.Create(ctx, device)
	require.NoError(t, err)

	// 驗證 JSON 配置正確儲存和讀取
	retrieved, err := repo.GetByID(ctx, device.ID)
	assert.NoError(t, err)
	assert.Equal(t, config, retrieved.ConnectionConfig)
	assert.Contains(t, retrieved.ConnectionConfig, "192.168.1.1")
}
