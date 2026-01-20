// Package tag 提供 SQL Repository 的單元測試。
package tag

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
 * createTestTag 建立測試用的標籤資料
 * @param t 測試實例
 * @returns *schema.Tag 測試標籤物件
 */
func createTestTag(t *testing.T) *schema.Tag {
	now := time.Now().UTC()
	return &schema.Tag{
		ID:          "test-tag-001",
		Key:         "test.tag.001",
		DisplayName: "測試標籤 001",
		DataType:    schema.DataTypeFloat64,
		Unit:        "°C",
		Description: "這是一個測試標籤",
		Status:      schema.TagStatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

/**
 * TestSQLRepository_Create 測試建立標籤功能
 */
func TestSQLRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	tag := createTestTag(t)

	// 執行 Create
	err := repo.Create(ctx, tag)
	assert.NoError(t, err)

	// 驗證標籤已建立
	retrieved, err := repo.GetByID(ctx, tag.ID)
	assert.NoError(t, err)
	assert.Equal(t, tag.ID, retrieved.ID)
	assert.Equal(t, tag.Key, retrieved.Key)
	assert.Equal(t, tag.DisplayName, retrieved.DisplayName)
	assert.Equal(t, tag.DataType, retrieved.DataType)
}

/**
 * TestSQLRepository_Create_DuplicateKey 測試建立重複 Key 的標籤應失敗
 */
func TestSQLRepository_Create_DuplicateKey(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	tag := createTestTag(t)

	// 第一次建立應成功
	err := repo.Create(ctx, tag)
	assert.NoError(t, err)

	// 第二次建立相同 Key 應失敗（透過 key_lower 唯一約束）
	tag2 := createTestTag(t)
	tag2.ID = "test-tag-002"
	tag2.Key = "TEST.TAG.001" // 大寫版本，應該因為 key_lower 衝突而失敗

	err = repo.Create(ctx, tag2)
	assert.Error(t, err)
}

/**
 * TestSQLRepository_Update 測試更新標籤功能
 */
func TestSQLRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	tag := createTestTag(t)
	err := repo.Create(ctx, tag)
	require.NoError(t, err)

	// 更新標籤資料
	tag.DisplayName = "更新後的標籤名稱"
	tag.DataType = schema.DataTypeInt32
	tag.Unit = "rpm"
	tag.Description = "更新後的描述"
	tag.Status = schema.TagStatusRetired

	err = repo.Update(ctx, tag)
	assert.NoError(t, err)

	// 驗證更新結果
	retrieved, err := repo.GetByID(ctx, tag.ID)
	assert.NoError(t, err)
	assert.Equal(t, "更新後的標籤名稱", retrieved.DisplayName)
	assert.Equal(t, schema.DataTypeInt32, retrieved.DataType)
	assert.Equal(t, "rpm", retrieved.Unit)
	assert.Equal(t, "更新後的描述", retrieved.Description)
	assert.Equal(t, schema.TagStatusRetired, retrieved.Status)
}

/**
 * TestSQLRepository_Update_NotFound 測試更新不存在的標籤應失敗
 */
func TestSQLRepository_Update_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	tag := createTestTag(t)
	tag.ID = "non-existent-id"

	err := repo.Update(ctx, tag)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "標籤不存在")
}

/**
 * TestSQLRepository_Delete 測試刪除標籤功能
 */
func TestSQLRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	tag := createTestTag(t)
	err := repo.Create(ctx, tag)
	require.NoError(t, err)

	// 刪除標籤
	err = repo.Delete(ctx, tag.ID)
	assert.NoError(t, err)

	// 驗證標籤已刪除
	_, err = repo.GetByID(ctx, tag.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "標籤不存在")
}

/**
 * TestSQLRepository_Delete_NotFound 測試刪除不存在的標籤應失敗
 */
func TestSQLRepository_Delete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	err := repo.Delete(ctx, "non-existent-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "標籤不存在")
}

/**
 * TestSQLRepository_GetByID 測試根據 ID 取得標籤
 */
func TestSQLRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	tag := createTestTag(t)
	err := repo.Create(ctx, tag)
	require.NoError(t, err)

	// 取得標籤
	retrieved, err := repo.GetByID(ctx, tag.ID)
	assert.NoError(t, err)
	assert.Equal(t, tag.ID, retrieved.ID)
	assert.Equal(t, tag.Key, retrieved.Key)
	assert.Equal(t, tag.DisplayName, retrieved.DisplayName)
}

/**
 * TestSQLRepository_GetByID_NotFound 測試取得不存在的標籤
 */
func TestSQLRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	_, err := repo.GetByID(ctx, "non-existent-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "標籤不存在")
}

/**
 * TestSQLRepository_GetByKey 測試根據 Key 取得標籤
 */
func TestSQLRepository_GetByKey(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	tag := createTestTag(t)
	err := repo.Create(ctx, tag)
	require.NoError(t, err)

	// 使用原 Key 取得
	retrieved, err := repo.GetByKey(ctx, tag.Key)
	assert.NoError(t, err)
	assert.Equal(t, tag.ID, retrieved.ID)
	assert.Equal(t, tag.Key, retrieved.Key)

	// 使用大寫 Key 取得（應透過 key_lower 找到）
	retrieved, err = repo.GetByKey(ctx, "TEST.TAG.001")
	assert.NoError(t, err)
	assert.Equal(t, tag.ID, retrieved.ID)
}

/**
 * TestSQLRepository_GetByKey_NotFound 測試取得不存在的 Key
 */
func TestSQLRepository_GetByKey_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	_, err := repo.GetByKey(ctx, "non.existent.key")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "標籤不存在")
}

/**
 * TestSQLRepository_List 測試列出標籤
 */
func TestSQLRepository_List(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	// 建立多個標籤
	tags := []*schema.Tag{
		createTestTag(t),
		{
			ID:          "test-tag-002",
			Key:         "test.tag.002",
			DisplayName: "測試標籤 002",
			DataType:    schema.DataTypeInt32,
			Unit:        "rpm",
			Status:      schema.TagStatusActive,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
		{
			ID:          "test-tag-003",
			Key:         "test.tag.003",
			DisplayName: "測試標籤 003",
			DataType:    schema.DataTypeBool,
			Unit:        "",
			Status:      schema.TagStatusDraft,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
	}

	for _, tag := range tags {
		err := repo.Create(ctx, tag)
		require.NoError(t, err)
	}

	// 列出所有標籤
	result, err := repo.List(ctx, ListFilter{})
	assert.NoError(t, err)
	assert.Len(t, result, 3)

	// 驗證排序（按 key ASC）
	assert.Equal(t, "test.tag.001", result[0].Key)
	assert.Equal(t, "test.tag.002", result[1].Key)
	assert.Equal(t, "test.tag.003", result[2].Key)
}

/**
 * TestSQLRepository_List_WithFilter 測試使用過濾條件列出標籤
 */
func TestSQLRepository_List_WithFilter(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	// 建立不同狀態和類型的標籤
	tags := []*schema.Tag{
		{
			ID:          "test-float64-active",
			Key:         "test.float64.active",
			DisplayName: "Float64 Active",
			DataType:    schema.DataTypeFloat64,
			Status:      schema.TagStatusActive,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
		{
			ID:          "test-float64-draft",
			Key:         "test.float64.draft",
			DisplayName: "Float64 Draft",
			DataType:    schema.DataTypeFloat64,
			Status:      schema.TagStatusDraft,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
		{
			ID:          "test-int32-active",
			Key:         "test.int32.active",
			DisplayName: "Int32 Active",
			DataType:    schema.DataTypeInt32,
			Status:      schema.TagStatusActive,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
	}

	for _, tag := range tags {
		err := repo.Create(ctx, tag)
		require.NoError(t, err)
	}

	// 測試狀態過濾
	status := schema.TagStatusActive
	result, err := repo.List(ctx, ListFilter{Status: &status})
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	// 測試資料類型過濾
	dataType := schema.DataTypeFloat64
	result, err = repo.List(ctx, ListFilter{DataType: &dataType})
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	// 測試組合過濾
	result, err = repo.List(ctx, ListFilter{
		Status:   &status,
		DataType: &dataType,
	})
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "test.float64.active", result[0].Key)
}

/**
 * TestSQLRepository_List_WithPagination 測試分頁功能
 */
func TestSQLRepository_List_WithPagination(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	// 建立 5 個標籤
	for i := 1; i <= 5; i++ {
		tag := createTestTag(t)
		tag.ID = string(rune('0' + i))
		tag.Key = string(rune('0' + i))
		tag.DisplayName = string(rune('0' + i))
		err := repo.Create(ctx, tag)
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
 * TestSQLRepository_Count 測試計算標籤數量
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

	// 建立 3 個標籤
	for i := 0; i < 3; i++ {
		tag := createTestTag(t)
		tag.ID = string(rune('0' + i))
		tag.Key = string(rune('0' + i))
		err := repo.Create(ctx, tag)
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
	tag := &schema.Tag{
		ID:          "test-null-fields",
		Key:         "test.null.fields",
		DisplayName: "測試 NULL 欄位",
		DataType:    schema.DataTypeInt16,
		Status:      schema.TagStatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
		// Unit, Description 為空字串
	}

	err := repo.Create(ctx, tag)
	require.NoError(t, err)

	// 取回並驗證 NULL 欄位正確處理
	retrieved, err := repo.GetByID(ctx, tag.ID)
	assert.NoError(t, err)
	assert.Equal(t, "", retrieved.Unit)
	assert.Equal(t, "", retrieved.Description)
}

/**
 * TestSQLRepository_KeyLowerAutoFill 測試 Key 自動轉小寫功能
 */
func TestSQLRepository_KeyLowerAutoFill(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	now := time.Now().UTC()
	tag := &schema.Tag{
		ID:          "test-mixed-case",
		Key:         "Test.Mixed.Case.Key",
		DisplayName: "混合大小寫的 Key",
		DataType:    schema.DataTypeInt16,
		Status:      schema.TagStatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err := repo.Create(ctx, tag)
	require.NoError(t, err)

	// 驗證可以使用小寫 Key 取得
	retrieved, err := repo.GetByKey(ctx, "test.mixed.case.key")
	assert.NoError(t, err)
	assert.Equal(t, tag.ID, retrieved.ID)
	assert.Equal(t, "Test.Mixed.Case.Key", retrieved.Key) // 原 Key 保持不變
}

/**
 * TestSQLRepository_DataTypeValidation 測試資料類型驗證
 */
func TestSQLRepository_DataTypeValidation(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

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
		tag := createTestTag(t)
		tag.ID = string(rune('0' + i))
		tag.Key = string(rune('0' + i))
		tag.DisplayName = string(rune('0' + i))
		tag.DataType = dataType
		err := repo.Create(ctx, tag)
		assert.NoError(t, err, "DataType %s should be valid", dataType)
	}
}

/**
 * TestSQLRepository_StatusValidation 測試狀態驗證
 */
func TestSQLRepository_StatusValidation(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	// 測試所有支援的狀態
	statuses := []schema.TagStatus{
		schema.TagStatusDraft,
		schema.TagStatusActive,
		schema.TagStatusRetired,
	}

	for i, status := range statuses {
		tag := createTestTag(t)
		tag.ID = string(rune('0' + i))
		tag.Key = string(rune('0' + i))
		tag.DisplayName = string(rune('0' + i))
		tag.Status = status
		err := repo.Create(ctx, tag)
		assert.NoError(t, err, "Status %s should be valid", status)
	}
}

/**
 * TestSQLRepository_WithUnit 測試單位欄位
 */
func TestSQLRepository_WithUnit(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	units := []string{"°C", "°F", "bar", "Pa", "rpm", "m/s", "%", ""}

	for i, unit := range units {
		tag := createTestTag(t)
		tag.ID = string(rune('0' + i))
		tag.Key = string(rune('0' + i))
		tag.DisplayName = string(rune('0' + i))
		tag.Unit = unit
		err := repo.Create(ctx, tag)
		require.NoError(t, err)

		// 驗證單位正確儲存
		retrieved, err := repo.GetByID(ctx, tag.ID)
		assert.NoError(t, err)
		assert.Equal(t, unit, retrieved.Unit)
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

	tag := createTestTag(t)
	tag.Description = "這是一個非常長的標籤描述，用於測試描述欄位的儲存和讀取功能是否正常運作。"

	err := repo.Create(ctx, tag)
	require.NoError(t, err)

	// 驗證描述正確儲存
	retrieved, err := repo.GetByID(ctx, tag.ID)
	assert.NoError(t, err)
	assert.Equal(t, tag.Description, retrieved.Description)
}

/**
 * TestSQLRepository_KeyCaseInsensitivity 測試 Key 的不區分大小寫特性
 */
func TestSQLRepository_KeyCaseInsensitivity(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	tag := createTestTag(t)
	err := repo.Create(ctx, tag)
	require.NoError(t, err)

	// 測試各種大小寫組合都能找到標籤
	testKeys := []string{
		"test.tag.001",
		"TEST.TAG.001",
		"Test.Tag.001",
		"tEsT.tAg.001",
	}

	for _, testKey := range testKeys {
		retrieved, err := repo.GetByKey(ctx, testKey)
		assert.NoError(t, err, "Key %s should work", testKey)
		assert.Equal(t, tag.ID, retrieved.ID)
	}
}
