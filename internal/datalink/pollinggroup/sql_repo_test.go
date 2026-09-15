// Package pollinggroup 提供 SQL Repository 的單元測試。
package pollinggroup

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"go-gateway/internal/datalink/common"
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
	ctx := t.Context()
	// 讀取並執行 migration 腳本
	migrationPath := "../schema/migrations/001_initial_schema_sqlite.sql"
	migrationContent, err := os.ReadFile(migrationPath)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, string(migrationContent))
	require.NoError(t, err)

	return db
}

/**
 * createTestPollingGroup 建立測試用的輪詢群組資料
 * @param t 測試實例
 * @returns *schema.PollingGroup 測試輪詢群組物件
 */
func createTestPollingGroup(t *testing.T) *schema.PollingGroup {
	now := time.Now().UTC()
	id, err := common.NewUUID()
	require.NoError(t, err)
	return &schema.PollingGroup{
		ID:          id,
		Name:        "測試群組",
		Description: "這是一個測試輪詢群組",
		IntervalMs:  1000,
		Priority:    100,
		Enabled:     true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

/**
 * TestSQLRepository_Create 測試建立輪詢群組功能
 */
func TestSQLRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	group := createTestPollingGroup(t)

	// 執行 Create
	err := repo.Create(ctx, group)
	assert.NoError(t, err)
	assert.NotEmpty(t, group.ID)

	// 驗證輪詢群組已建立
	retrieved, err := repo.GetByID(ctx, group.ID)
	assert.NoError(t, err)
	assert.Equal(t, group.ID, retrieved.ID)
	assert.Equal(t, group.Name, retrieved.Name)
	assert.Equal(t, group.IntervalMs, retrieved.IntervalMs)
}

/**
 * TestSQLRepository_Create_DuplicateName 測試建立重複名稱的輪詢群組應失敗
 */
func TestSQLRepository_Create_DuplicateName(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	group := createTestPollingGroup(t)

	// 第一次建立應成功
	err := repo.Create(ctx, group)
	assert.NoError(t, err)

	// 第二次建立相同名稱應失敗（透過 name UNIQUE 約束）
	group2 := createTestPollingGroup(t)
	group2.ID = "" // 清空 ID 以便生成新 ID
	err = repo.Create(ctx, group2)
	assert.Error(t, err)
}

/**
 * TestSQLRepository_Update 測試更新輪詢群組功能
 */
func TestSQLRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	group := createTestPollingGroup(t)
	err := repo.Create(ctx, group)
	require.NoError(t, err)

	// 更新輪詢群組資料
	group.Name = "更新後的群組名稱"
	group.Description = "更新後的描述"
	group.IntervalMs = 2000
	group.Priority = 200
	group.Enabled = false

	err = repo.Update(ctx, group)
	assert.NoError(t, err)

	// 驗證更新結果
	retrieved, err := repo.GetByID(ctx, group.ID)
	assert.NoError(t, err)
	assert.Equal(t, "更新後的群組名稱", retrieved.Name)
	assert.Equal(t, "更新後的描述", retrieved.Description)
	assert.Equal(t, 2000, retrieved.IntervalMs)
	assert.Equal(t, 200, retrieved.Priority)
	assert.False(t, retrieved.Enabled)
}

/**
 * TestSQLRepository_Update_NotFound 測試更新不存在的輪詢群組應失敗
 */
func TestSQLRepository_Update_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	group := createTestPollingGroup(t)
	group.ID = "non-existent-id"

	err := repo.Update(ctx, group)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "輪詢群組不存在")
}

/**
 * TestSQLRepository_Delete 測試刪除輪詢群組功能
 */
func TestSQLRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	group := createTestPollingGroup(t)
	err := repo.Create(ctx, group)
	require.NoError(t, err)

	// 刪除輪詢群組
	err = repo.Delete(ctx, group.ID)
	assert.NoError(t, err)

	// 驗證輪詢群組已刪除
	_, err = repo.GetByID(ctx, group.ID)
	assert.Error(t, err)
}

/**
 * TestSQLRepository_Delete_NotFound 測試刪除不存在的輪詢群組應失敗
 */
func TestSQLRepository_Delete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	err := repo.Delete(ctx, "non-existent-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "輪詢群組不存在")
}

/**
 * TestSQLRepository_GetByID 測試根據 ID 取得輪詢群組
 */
func TestSQLRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	group := createTestPollingGroup(t)
	err := repo.Create(ctx, group)
	require.NoError(t, err)

	// 取得輪詢群組
	retrieved, err := repo.GetByID(ctx, group.ID)
	assert.NoError(t, err)
	assert.Equal(t, group.ID, retrieved.ID)
	assert.Equal(t, group.Name, retrieved.Name)
	assert.Equal(t, group.IntervalMs, retrieved.IntervalMs)
}

/**
 * TestSQLRepository_GetByID_NotFound 測試取得不存在的輪詢群組
 */
func TestSQLRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	_, err := repo.GetByID(ctx, "non-existent-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "輪詢群組不存在")
}

/**
 * TestSQLRepository_List 測試列出輪詢群組
 */
func TestSQLRepository_List(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	// 建立多個輪詢群組
	groups := []*schema.PollingGroup{
		createTestPollingGroup(t),
		{
			Name:        "測試群組 2",
			Description: "這是第二個測試群組",
			IntervalMs:  500,
			Priority:    200,
			Enabled:     true,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
		{
			Name:        "測試群組 3",
			Description: "這是第三個測試群組",
			IntervalMs:  2000,
			Priority:    50,
			Enabled:     false,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
	}

	for _, group := range groups {
		err := repo.Create(ctx, group)
		require.NoError(t, err)
	}

	// 列出所有輪詢群組
	result, err := repo.List(ctx)
	assert.NoError(t, err)
	assert.Len(t, result, 3)

	// 驗證排序（按 priority DESC, name ASC）
	assert.Equal(t, "測試群組 2", result[0].Name)
	assert.Equal(t, "測試群組", result[1].Name)
	assert.Equal(t, "測試群組 3", result[2].Name)
}

/**
 * TestSQLRepository_Count 測試計算輪詢群組數量
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

	// 建立 3 個輪詢群組
	for i := 0; i < 3; i++ {
		group := createTestPollingGroup(t)
		group.Name = string(rune('0' + i))
		err := repo.Create(ctx, group)
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
	group := &schema.PollingGroup{
		Name:       "測試 NULL 欄位",
		IntervalMs: 1000,
		Priority:   100,
		Enabled:    true,
		CreatedAt:  now,
		UpdatedAt:  now,
		// Description 為空字串
	}

	err := repo.Create(ctx, group)
	require.NoError(t, err)

	// 取回並驗證 NULL 欄位正確處理
	retrieved, err := repo.GetByID(ctx, group.ID)
	assert.NoError(t, err)
	assert.Equal(t, "", retrieved.Description)
}

/**
 * TestSQLRepository_Clear 測試清空所有輪詢群組
 */
func TestSQLRepository_Clear(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	// 建立 3 個輪詢群組
	for i := 0; i < 3; i++ {
		group := createTestPollingGroup(t)
		group.Name = string(rune('0' + i))
		err := repo.Create(ctx, group)
		require.NoError(t, err)
	}

	// 驗證存在
	count, err := repo.Count(ctx)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), count)

	// 清空所有輪詢群組
	err = repo.Clear(ctx)
	assert.NoError(t, err)

	// 驗證已清空
	count, err = repo.Count(ctx)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

/**
 * TestSQLRepository_IntervalValidation 測試間隔時間驗證
 */
func TestSQLRepository_IntervalValidation(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	// 測試有效間隔時間
	validIntervals := []int32{100, 500, 1000, 5000, 10000}

	for i, interval := range validIntervals {
		group := createTestPollingGroup(t)
		group.Name = string(rune('0' + i))
		group.IntervalMs = int(interval)
		err := repo.Create(ctx, group)
		assert.NoError(t, err, "Interval %d should be valid", interval)
	}

	// 測試無效間隔時間（小於 100ms）
	group := createTestPollingGroup(t)
	group.Name = "invalid-interval"
	group.IntervalMs = 50
	err := repo.Create(ctx, group)
	assert.Error(t, err)
}

/**
 * TestSQLRepository_EnabledStatus 測試啟用狀態
 */
func TestSQLRepository_EnabledStatus(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	// 建立啟用和停用的輪詢群組
	enabledGroup := createTestPollingGroup(t)
	enabledGroup.Name = "啟用群組"
	enabledGroup.Enabled = true
	err := repo.Create(ctx, enabledGroup)
	require.NoError(t, err)

	disabledGroup := createTestPollingGroup(t)
	disabledGroup.Name = "停用群組"
	disabledGroup.Enabled = false
	err = repo.Create(ctx, disabledGroup)
	require.NoError(t, err)

	// 驗證啟用狀態
	retrieved, err := repo.GetByID(ctx, enabledGroup.ID)
	assert.NoError(t, err)
	assert.True(t, retrieved.Enabled)

	retrieved, err = repo.GetByID(ctx, disabledGroup.ID)
	assert.NoError(t, err)
	assert.False(t, retrieved.Enabled)
}

/**
 * TestSQLRepository_Priority 測試優先級欄位
 */
func TestSQLRepository_Priority(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	// 建立不同優先級的群組
	priorities := []int32{10, 50, 100, 200, 500}

	for i, priority := range priorities {
		group := createTestPollingGroup(t)
		group.Name = string(rune('0' + i))
		group.Priority = int(priority)
		err := repo.Create(ctx, group)
		require.NoError(t, err, "Priority %d should be valid", priority)
	}

	// 驗證排序（按 priority DESC）
	result, err := repo.List(ctx)
	assert.NoError(t, err)
	assert.Len(t, result, 5)

	// 檢查優先級降序排列
	for i := 0; i < len(result)-1; i++ {
		assert.GreaterOrEqual(t, result[i].Priority, result[i+1].Priority,
			"Group at index %d should have higher or equal priority than group at index %d", i, i+1)
	}
}

/**
 * TestSQLRepository_WithDescription 測試描述欄位
 */
func TestSQLRepository_WithDescription(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	description := "這是一個非常長的輪詢群組描述，用於測試描述欄位的儲存和讀取功能是否正常運作。"
	group := createTestPollingGroup(t)
	group.Description = description

	err := repo.Create(ctx, group)
	require.NoError(t, err)

	// 驗證描述正確儲存
	retrieved, err := repo.GetByID(ctx, group.ID)
	assert.NoError(t, err)
	assert.Equal(t, description, retrieved.Description)
}

/**
 * TestSQLRepository_NameUniqueness 測試名稱唯一性約束
 */
func TestSQLRepository_NameUniqueness(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	// 建立第一個群組
	group1 := createTestPollingGroup(t)
	group1.Name = "唯一名稱"
	err := repo.Create(ctx, group1)
	require.NoError(t, err)

	// 建立相同名稱的第二個群組應失敗
	group2 := createTestPollingGroup(t)
	group2.Name = "唯一名稱" // 相同名稱
	err = repo.Create(ctx, group2)
	assert.Error(t, err)

	// 不同名稱應成功
	group3 := createTestPollingGroup(t)
	group3.Name = "另一個名稱"
	err = repo.Create(ctx, group3)
	assert.NoError(t, err)
}
