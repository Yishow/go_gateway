// Package point 提供 SQL Repository 的單元測試。
package point

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

	_, err = db.Exec(`ALTER TABLE points ADD COLUMN data_format TEXT`)
	require.NoError(t, err)

	return db
}

/**
 * setupTestData 建立測試所需的基礎資料（設備、輪詢群組）
 * @param t 測試實例
 * @param db 資料庫連線
 * @returns deviceID, pollingGroupID 設備 ID 與輪詢群組 ID
 */
func setupTestData(t *testing.T, db *sql.DB) (string, string) {
	ctx := context.Background()
	now := time.Now().UTC()

	// 建立設備
	deviceID := "test-device-001"
	_, err := db.ExecContext(ctx, `
		INSERT INTO devices (id, name, protocol, status, connection_config, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, deviceID, "測試設備", schema.ProtocolModbusTCP, schema.DeviceStatusActive, "{}", now, now)
	require.NoError(t, err)

	// 建立輪詢群組
	pollingGroupID := "test-group-001"
	_, err = db.ExecContext(ctx, `
		INSERT INTO polling_groups (id, name, description, interval_ms, priority, enabled, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, pollingGroupID, "測試群組", "這是測試群組", 1000, 100, 1, now, now)
	require.NoError(t, err)

	return deviceID, pollingGroupID
}

/**
 * createTestPoint 建立測試用的點位資料
 * @param t 測試實例
 * @param deviceID 設備 ID
 * @returns *schema.Point 測試點位物件
 */
func createTestPoint(t *testing.T, deviceID string) *schema.Point {
	now := time.Now().UTC()
	return &schema.Point{
		ID:             "test-point-001",
		DeviceID:       deviceID,
		Name:           "測試點位",
		Description:    "這是一個測試點位",
		Address:        "40001",
		Function:       "read_holding_registers",
		DataType:       schema.DataTypeInt16,
		Mode:           schema.PointModeReadOnly,
		PollingGroupID: nil,
		Enabled:        true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

/**
 * TestSQLRepository_Create 測試建立點位功能
 */
func TestSQLRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	deviceID, _ := setupTestData(t, db)
	point := createTestPoint(t, deviceID)

	// 執行 Create
	err := repo.Create(ctx, point)
	assert.NoError(t, err)

	// 驗證點位已建立
	retrieved, err := repo.GetByID(ctx, point.ID)
	assert.NoError(t, err)
	assert.Equal(t, point.ID, retrieved.ID)
	assert.Equal(t, point.Name, retrieved.Name)
	assert.Equal(t, point.Address, retrieved.Address)
	assert.Equal(t, point.DataType, retrieved.DataType)
}

/**
 * TestSQLRepository_Create_Duplicate 測試建立重複點位應失敗
 */
func TestSQLRepository_Create_Duplicate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	deviceID, _ := setupTestData(t, db)
	point := createTestPoint(t, deviceID)

	// 第一次建立應成功
	err := repo.Create(ctx, point)
	assert.NoError(t, err)

	// 第二次建立相同 ID 應失敗
	err = repo.Create(ctx, point)
	assert.Error(t, err)
}

/**
 * TestSQLRepository_Update 測試更新點位功能
 */
func TestSQLRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	deviceID, _ := setupTestData(t, db)
	point := createTestPoint(t, deviceID)
	err := repo.Create(ctx, point)
	require.NoError(t, err)

	// 更新點位資料
	point.Name = "更新後的點位名稱"
	point.Description = "更新後的描述"
	point.Address = "40002"
	point.DataType = schema.DataTypeFloat64
	point.Enabled = false

	err = repo.Update(ctx, point)
	assert.NoError(t, err)

	// 驗證更新結果
	retrieved, err := repo.GetByID(ctx, point.ID)
	assert.NoError(t, err)
	assert.Equal(t, "更新後的點位名稱", retrieved.Name)
	assert.Equal(t, "更新後的描述", retrieved.Description)
	assert.Equal(t, "40002", retrieved.Address)
	assert.Equal(t, schema.DataTypeFloat64, retrieved.DataType)
	assert.False(t, retrieved.Enabled)
}

/**
 * TestSQLRepository_Update_NotFound 測試更新不存在的點位應失敗
 */
func TestSQLRepository_Update_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	deviceID, _ := setupTestData(t, db)
	point := createTestPoint(t, deviceID)
	point.ID = "non-existent-id"

	err := repo.Update(ctx, point)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrPointNotFound))
}

/**
 * TestSQLRepository_Delete 測試刪除點位功能
 */
func TestSQLRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	deviceID, _ := setupTestData(t, db)
	point := createTestPoint(t, deviceID)
	err := repo.Create(ctx, point)
	require.NoError(t, err)

	// 刪除點位
	err = repo.Delete(ctx, point.ID)
	assert.NoError(t, err)

	// 驗證點位已刪除
	_, err = repo.GetByID(ctx, point.ID)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrPointNotFound))
}

/**
 * TestSQLRepository_Delete_NotFound 測試刪除不存在的點位應失敗
 */
func TestSQLRepository_Delete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	err := repo.Delete(ctx, "non-existent-id")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrPointNotFound))
}

/**
 * TestSQLRepository_GetByID 測試根據 ID 取得點位
 */
func TestSQLRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	deviceID, _ := setupTestData(t, db)
	point := createTestPoint(t, deviceID)
	err := repo.Create(ctx, point)
	require.NoError(t, err)

	// 取得點位
	retrieved, err := repo.GetByID(ctx, point.ID)
	assert.NoError(t, err)
	assert.Equal(t, point.ID, retrieved.ID)
	assert.Equal(t, point.Name, retrieved.Name)
	assert.Equal(t, point.DeviceID, retrieved.DeviceID)
}

/**
 * TestSQLRepository_GetByID_NotFound 測試取得不存在的點位
 */
func TestSQLRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	_, err := repo.GetByID(ctx, "non-existent-id")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrPointNotFound))
}

/**
 * TestSQLRepository_ListByDevice 測試根據設備 ID 取得點位列表
 */
func TestSQLRepository_ListByDevice(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	deviceID, _ := setupTestData(t, db)

	// 建立多個點位
	for i := 1; i <= 3; i++ {
		point := createTestPoint(t, deviceID)
		point.ID = string(rune('0' + i))
		point.Name = string(rune('0' + i))
		point.Address = fmt.Sprintf("4000%d", i)
		err := repo.Create(ctx, point)
		require.NoError(t, err)
	}

	// 根據設備 ID 取得點位
	points, err := repo.ListByDevice(ctx, deviceID)
	assert.NoError(t, err)
	assert.Len(t, points, 3)

	// 驗證排序（按 name ASC）
	assert.Equal(t, "1", points[0].Name)
	assert.Equal(t, "2", points[1].Name)
	assert.Equal(t, "3", points[2].Name)
}

/**
 * TestSQLRepository_List 測試列出點位
 */
func TestSQLRepository_List(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	deviceID, _ := setupTestData(t, db)

	// 建立多個點位
	points := []*schema.Point{
		createTestPoint(t, deviceID),
		{
			ID:        "test-point-002",
			DeviceID:  deviceID,
			Name:      "測試點位 2",
			Address:   "40002",
			DataType:  schema.DataTypeFloat32,
			Mode:      schema.PointModeReadOnly,
			Enabled:   true,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		{
			ID:        "test-point-003",
			DeviceID:  deviceID,
			Name:      "測試點位 3",
			Address:   "40003",
			DataType:  schema.DataTypeBool,
			Mode:      schema.PointModeReadWrite,
			Enabled:   false,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
	}

	for _, point := range points {
		err := repo.Create(ctx, point)
		require.NoError(t, err)
	}

	// 列出所有點位
	result, err := repo.List(ctx, ListFilter{})
	assert.NoError(t, err)
	assert.Len(t, result, 3)

	// 驗證排序（按 name ASC）
	assert.Equal(t, "測試點位", result[0].Name)
	assert.Equal(t, "測試點位 2", result[1].Name)
	assert.Equal(t, "測試點位 3", result[2].Name)
}

/**
 * TestSQLRepository_List_WithFilter 測試使用過濾條件列出點位
 */
func TestSQLRepository_List_WithFilter(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	deviceID, pollingGroupID := setupTestData(t, db)

	// 建立不同設備和群組的點位
	points := []*schema.Point{
		{
			ID:             "test-point-device1",
			DeviceID:       deviceID,
			Name:           "Device 1 Point",
			Address:        "40001",
			DataType:       schema.DataTypeInt16,
			Mode:           schema.PointModeReadOnly,
			PollingGroupID: &pollingGroupID,
			Enabled:        true,
			CreatedAt:      time.Now().UTC(),
			UpdatedAt:      time.Now().UTC(),
		},
		{
			ID:             "test-point-device2",
			DeviceID:       "other-device-id",
			Name:           "Device 2 Point",
			Address:        "40002",
			DataType:       schema.DataTypeFloat64,
			Mode:           schema.PointModeReadOnly,
			PollingGroupID: &pollingGroupID,
			Enabled:        true,
			CreatedAt:      time.Now().UTC(),
			UpdatedAt:      time.Now().UTC(),
		},
		{
			ID:             "test-point-disabled",
			DeviceID:       deviceID,
			Name:           "Disabled Point",
			Address:        "40003",
			DataType:       schema.DataTypeBool,
			Mode:           schema.PointModeReadWrite,
			PollingGroupID: &pollingGroupID,
			Enabled:        false,
			CreatedAt:      time.Now().UTC(),
			UpdatedAt:      time.Now().UTC(),
		},
	}

	// 建立第二個設備
	_, err := db.ExecContext(ctx, `
		INSERT INTO devices (id, name, protocol, status, connection_config, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, "other-device-id", "其他設備", schema.ProtocolModbusTCP, schema.DeviceStatusActive, "{}", time.Now().UTC(), time.Now().UTC())
	require.NoError(t, err)

	for _, point := range points {
		err := repo.Create(ctx, point)
		require.NoError(t, err)
	}

	// 測試設備過濾
	result, err := repo.List(ctx, ListFilter{DeviceID: &deviceID})
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	// 測試輪詢群組過濾
	result, err = repo.List(ctx, ListFilter{PollingGroupID: &pollingGroupID})
	assert.NoError(t, err)
	assert.Len(t, result, 3)

	// 測試組合過濾
	result, err = repo.List(ctx, ListFilter{
		DeviceID:       &deviceID,
		PollingGroupID: &pollingGroupID,
	})
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

/**
 * TestSQLRepository_List_WithPagination 測試分頁功能
 */
func TestSQLRepository_List_WithPagination(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	deviceID, _ := setupTestData(t, db)

	// 建立 5 個點位
	for i := 1; i <= 5; i++ {
		point := createTestPoint(t, deviceID)
		point.ID = string(rune('0' + i))
		point.Name = string(rune('0' + i))
		point.Address = fmt.Sprintf("4000%d", i)
		err := repo.Create(ctx, point)
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
 * TestSQLRepository_NullFields 測試處理 NULL 欄位
 */
func TestSQLRepository_NullFields(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	deviceID, _ := setupTestData(t, db)

	now := time.Now().UTC()
	point := &schema.Point{
		ID:        "test-null-fields",
		DeviceID:  deviceID,
		Name:      "測試 NULL 欄位",
		Address:   "40001",
		DataType:  schema.DataTypeInt16,
		Mode:      schema.PointModeReadOnly,
		Enabled:   true,
		CreatedAt: now,
		UpdatedAt: now,
		// Description, Function, PollingGroupID, LastReadAt, LastValue, LastError 為 nil 或空字串
	}

	err := repo.Create(ctx, point)
	require.NoError(t, err)

	// 取回並驗證 NULL 欄位正確處理
	retrieved, err := repo.GetByID(ctx, point.ID)
	assert.NoError(t, err)
	assert.Equal(t, "", retrieved.Description)
	assert.Equal(t, "", retrieved.Function)
	assert.Nil(t, retrieved.PollingGroupID)
	assert.Nil(t, retrieved.LastReadAt)
	assert.Nil(t, retrieved.LastValue)
	assert.Equal(t, "", retrieved.LastError)
}

/**
 * TestSQLRepository_WithPollingGroup 測試點位與輪詢群組的關聯
 */
func TestSQLRepository_WithPollingGroup(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	deviceID, pollingGroupID := setupTestData(t, db)

	point := createTestPoint(t, deviceID)
	point.PollingGroupID = &pollingGroupID

	err := repo.Create(ctx, point)
	require.NoError(t, err)

	// 驗證輪詢群組 ID 正確關聯
	retrieved, err := repo.GetByID(ctx, point.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved.PollingGroupID)
	assert.Equal(t, pollingGroupID, *retrieved.PollingGroupID)
}

/**
 * TestSQLRepository_UpdateLastRead 測試更新最後讀取相關欄位
 */
func TestSQLRepository_UpdateLastRead(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	deviceID, _ := setupTestData(t, db)
	point := createTestPoint(t, deviceID)
	err := repo.Create(ctx, point)
	require.NoError(t, err)

	// 更新最後讀取結果（透過 Update）
	now := time.Now().UTC()
	lastValue := "123.45"
	lastError := ""

	point.LastReadAt = &now
	point.LastValue = &lastValue
	point.LastError = lastError

	err = repo.Update(ctx, point)
	assert.NoError(t, err)

	// 驗證最後讀取結果已更新
	retrieved, err := repo.GetByID(ctx, point.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved.LastReadAt)
	assert.NotNil(t, retrieved.LastValue)
	assert.Equal(t, lastValue, *retrieved.LastValue)
	assert.Equal(t, lastError, retrieved.LastError)
}

/**
 * TestSQLRepository_UniquenessConstraint 測試設備內 address+function 的唯一性約束
 */
func TestSQLRepository_UniquenessConstraint(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	deviceID, _ := setupTestData(t, db)

	// 建立第一個點位
	point1 := createTestPoint(t, deviceID)
	point1.Function = "03"
	err := repo.Create(ctx, point1)
	require.NoError(t, err)

	// 相同設備和位址，但不同 function 應成功
	point2 := createTestPoint(t, deviceID)
	point2.ID = "test-point-002"
	point2.Name = "不同功能點位"
	point2.Function = "04"
	err = repo.Create(ctx, point2)
	assert.NoError(t, err)

	// 相同設備、位址與 function 應失敗
	point3 := createTestPoint(t, deviceID)
	point3.ID = "test-point-003"
	point3.Name = "重複位址功能點位"
	point3.Function = "03"
	err = repo.Create(ctx, point3)
	assert.Error(t, err)

	// 不同設備但相同位址與 function 應成功
	otherDeviceID := "other-device"
	_, err = db.ExecContext(ctx, `
		INSERT INTO devices (id, name, protocol, status, connection_config, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, otherDeviceID, "其他設備", schema.ProtocolModbusTCP, schema.DeviceStatusActive, "{}", time.Now().UTC(), time.Now().UTC())
	require.NoError(t, err)

	point4 := createTestPoint(t, otherDeviceID)
	point4.ID = "test-point-004"
	point4.Function = "03"
	err = repo.Create(ctx, point4)
	assert.NoError(t, err)
}

/**
 * TestSQLRepository_DataTypeValidation 測試資料類型驗證
 */
func TestSQLRepository_DataTypeValidation(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	deviceID, _ := setupTestData(t, db)

	// 測試所有支援的資料類型
	dataTypes := []schema.DataType{
		schema.DataTypeBool,
		schema.DataTypeInt16,
		schema.DataTypeUint16,
		schema.DataTypeInt32,
		schema.DataTypeUint32,
		schema.DataTypeInt64,
		schema.DataTypeUint64,
		schema.DataTypeFloat32,
		schema.DataTypeFloat64,
		schema.DataTypeString,
	}

	for i, dataType := range dataTypes {
		point := createTestPoint(t, deviceID)
		point.ID = string(rune('0' + i))
		point.Address = fmt.Sprintf("5000%d", i)
		point.DataType = dataType
		err := repo.Create(ctx, point)
		assert.NoError(t, err, "DataType %s should be valid", dataType)
	}
}

/**
 * TestSQLRepository_ModeValidation 測試模式驗證
 */
func TestSQLRepository_ModeValidation(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	deviceID, _ := setupTestData(t, db)

	// 測試所有支援的模式
	modes := []schema.PointMode{
		schema.PointModeReadOnly,
		schema.PointModeReadWrite,
	}

	for i, mode := range modes {
		point := createTestPoint(t, deviceID)
		point.ID = string(rune('0' + i))
		point.Address = fmt.Sprintf("6000%d", i)
		point.Mode = mode
		err := repo.Create(ctx, point)
		assert.NoError(t, err, "Mode %s should be valid", mode)
	}
}

/**
 * TestSQLRepository_EnabledStatus 測試啟用狀態
 */
func TestSQLRepository_EnabledStatus(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	deviceID, _ := setupTestData(t, db)

	// 建立啟用和停用的點位
	enabledPoint := createTestPoint(t, deviceID)
	enabledPoint.ID = "enabled-point"
	enabledPoint.Address = "70001"
	enabledPoint.Enabled = true
	err := repo.Create(ctx, enabledPoint)
	require.NoError(t, err)

	disabledPoint := createTestPoint(t, deviceID)
	disabledPoint.ID = "disabled-point"
	disabledPoint.Address = "70002"
	disabledPoint.Enabled = false
	err = repo.Create(ctx, disabledPoint)
	require.NoError(t, err)

	// 驗證啟用狀態
	retrieved, err := repo.GetByID(ctx, enabledPoint.ID)
	assert.NoError(t, err)
	assert.True(t, retrieved.Enabled)

	retrieved, err = repo.GetByID(ctx, disabledPoint.ID)
	assert.NoError(t, err)
	assert.False(t, retrieved.Enabled)
}

func TestSQLRepository_UpdateReadResult(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	deviceID, _ := setupTestData(t, db)
	point := createTestPoint(t, deviceID)
	require.NoError(t, repo.Create(ctx, point))

	require.NoError(t, repo.UpdateReadResult(ctx, point.ID, 123.45, ""))

	retrieved, err := repo.GetByID(ctx, point.ID)
	require.NoError(t, err)
	require.NotNil(t, retrieved.LastValue)
	assert.Equal(t, "123.45", *retrieved.LastValue)
	assert.NotNil(t, retrieved.LastReadAt)
	assert.Equal(t, "", retrieved.LastError)
}

func TestSQLRepository_BatchUpdateReadResult(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	deviceID, _ := setupTestData(t, db)
	point1 := createTestPoint(t, deviceID)
	point1.ID = "batch-point-1"
	point1.Address = "80001"
	point2 := createTestPoint(t, deviceID)
	point2.ID = "batch-point-2"
	point2.Address = "80002"

	require.NoError(t, repo.Create(ctx, point1))
	require.NoError(t, repo.Create(ctx, point2))

	err := repo.BatchUpdateReadResult(ctx, []ReadResultUpdate{
		{PointID: point1.ID, Value: 11},
		{PointID: point2.ID, Error: "read failed"},
	})
	require.NoError(t, err)

	updated1, err := repo.GetByID(ctx, point1.ID)
	require.NoError(t, err)
	require.NotNil(t, updated1.LastValue)
	assert.Equal(t, "11", *updated1.LastValue)
	assert.Equal(t, "", updated1.LastError)

	updated2, err := repo.GetByID(ctx, point2.ID)
	require.NoError(t, err)
	assert.Nil(t, updated2.LastValue)
	assert.Equal(t, "read failed", updated2.LastError)
}
