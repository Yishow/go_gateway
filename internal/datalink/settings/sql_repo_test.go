// Package settings 提供 SQL Repository 的單元測試。
package settings

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

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
 * TestSQLRepository_Get 測試取得設定值
 */
func TestSQLRepository_Get(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	// 建立測試設定
	valueJSON := `{"enabled":true,"count":10}`
	now := time.Now().UTC()

	_, err := db.ExecContext(ctx, `
		INSERT INTO system_settings (key, value, updated_at)
		VALUES (?, ?, ?)
	`, "test.setting", valueJSON, now)
	require.NoError(t, err)

	// 取得設定
	item, err := repo.Get(ctx, "test.setting")
	assert.NoError(t, err)
	assert.Equal(t, "test.setting", item.Key)
	assert.NotNil(t, item.Value)
}

/**
 * TestSQLRepository_Get_NotFound 測試取得不存在的設定
 */
func TestSQLRepository_Get_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	_, err := repo.Get(ctx, "non.existent.setting")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "設定不存在")
}

/**
 * TestSQLRepository_Set 測試設定值
 */
func TestSQLRepository_Set(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	// 設定新值
	value := map[string]interface{}{
		"enabled": true,
		"count":   10,
	}

	err := repo.Set(ctx, "new.setting", value)
	assert.NoError(t, err)

	// 驗證值已設定
	item, err := repo.Get(ctx, "new.setting")
	assert.NoError(t, err)
	assert.Equal(t, "new.setting", item.Key)
	assert.NotNil(t, item.Value)
}

/**
 * TestSQLRepository_Set_Update 測試更新現有設定
 */
func TestSQLRepository_Set_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	// 建立初始設定
	initialValue := map[string]interface{}{
		"enabled": true,
		"count":   10,
	}
	err := repo.Set(ctx, "test.setting", initialValue)
	require.NoError(t, err)

	// 更新設定
	updatedValue := map[string]interface{}{
		"enabled": false,
		"count":   20,
	}
	err = repo.Set(ctx, "test.setting", updatedValue)
	assert.NoError(t, err)

	// 驗證值已更新
	item, err := repo.Get(ctx, "test.setting")
	assert.NoError(t, err)
	assert.Equal(t, "test.setting", item.Key)
	assert.NotNil(t, item.Value)
}

/**
 * TestSQLRepository_Set_SimpleType 測試設定簡單類型值
 */
func TestSQLRepository_Set_SimpleType(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	// 測試各種簡單類型
	testCases := []struct {
		key   string
		value interface{}
	}{
		{"string.value", "hello world"},
		{"int.value", 42},
		{"float.value", 3.14},
		{"bool.value", true},
	}

	for _, tc := range testCases {
		err := repo.Set(ctx, tc.key, tc.value)
		assert.NoError(t, err, "Failed to set %s", tc.key)

		item, err := repo.Get(ctx, tc.key)
		assert.NoError(t, err, "Failed to get %s", tc.key)
		assert.Equal(t, tc.key, item.Key)
	}
}

/**
 * TestSQLRepository_Set_ComplexType 測試設定複雜類型值
 */
func TestSQLRepository_Set_ComplexType(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	// 測試複雜類型（陣列、物件）
	value := map[string]interface{}{
		"array": []int{1, 2, 3, 4, 5},
		"nested": map[string]interface{}{
			"key1": "value1",
			"key2": 123,
		},
	}

	err := repo.Set(ctx, "complex.setting", value)
	assert.NoError(t, err)

	item, err := repo.Get(ctx, "complex.setting")
	assert.NoError(t, err)
	assert.Equal(t, "complex.setting", item.Key)
}

/**
 * TestSQLRepository_List 測試列出所有設定
 */
func TestSQLRepository_List(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	// 清除預設設定
	_, err := db.ExecContext(ctx, `DELETE FROM system_settings`)
	require.NoError(t, err)

	// 建立多個設定
	settings := []struct {
		key   string
		value interface{}
	}{
		{"setting.1", "value1"},
		{"setting.2", 42},
		{"setting.3", true},
		{"setting.a", "valuea"},
		{"setting.b", 3.14},
	}

	for _, s := range settings {
		err := repo.Set(ctx, s.key, s.value)
		require.NoError(t, err)
	}

	// 列出所有設定
	items, err := repo.List(ctx)
	assert.NoError(t, err)
	assert.Len(t, items, len(settings))

	// 驗證排序（按 key ASC）
	assert.Equal(t, "setting.1", items[0].Key)
	assert.Equal(t, "setting.2", items[1].Key)
	assert.Equal(t, "setting.3", items[2].Key)
	assert.Equal(t, "setting.a", items[3].Key)
	assert.Equal(t, "setting.b", items[4].Key)
}

/**
 * TestSQLRepository_List_Empty 測試列出空設定
 */
func TestSQLRepository_List_Empty(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	// 清除預設設定
	_, err := db.ExecContext(ctx, `DELETE FROM system_settings`)
	require.NoError(t, err)

	// 列出應為空
	items, err := repo.List(ctx)
	assert.NoError(t, err)
	assert.Len(t, items, 0)
}

/**
 * TestSQLRepository_Delete 測試刪除設定
 */
func TestSQLRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	// 建立設定
	err := repo.Set(ctx, "test.setting", "value")
	require.NoError(t, err)

	// 驗證設定存在
	_, err = repo.Get(ctx, "test.setting")
	assert.NoError(t, err)

	// 刪除設定
	err = repo.Delete(ctx, "test.setting")
	assert.NoError(t, err)

	// 驗證設定已刪除
	_, err = repo.Get(ctx, "test.setting")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "設定不存在")
}

/**
 * TestSQLRepository_Delete_NotFound 測試刪除不存在的設定
 */
func TestSQLRepository_Delete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	err := repo.Delete(ctx, "non.existent.setting")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "設定不存在")
}

/**
 * TestSQLRepository_InitDefaults 測試初始化預設設定
 */
func TestSQLRepository_InitDefaults(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	// 清除所有設定
	_, err := db.ExecContext(ctx, `DELETE FROM system_settings`)
	require.NoError(t, err)

	// 初始化預設設定
	err = repo.InitDefaults(ctx)
	assert.NoError(t, err)

	// 驗證預設設定已建立
	items, err := repo.List(ctx)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(items), 5) // 應至少有 5 個預設設定

	// 驗證特定預設設定存在
	expectedKeys := []string{
		KeyWritePrecision,
		KeyPartitionInterval,
		KeyBatchSize,
		KeyDefaultRetryCount,
		KeyDefaultRetryDelay,
	}

	for _, key := range expectedKeys {
		_, err := repo.Get(ctx, key)
		assert.NoError(t, err, "Default setting %s should exist", key)
	}
}

/**
 * TestSQLRepository_InitDefaults_AlreadyExists 測試已存在時不覆蓋
 */
func TestSQLRepository_InitDefaults_AlreadyExists(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	// 先初始化預設設定
	err := repo.InitDefaults(ctx)
	require.NoError(t, err)

	// 修改某個預設設定的值
	err = repo.Set(ctx, KeyBatchSize, 9999)
	require.NoError(t, err)

	// 再次初始化（不應覆蓋現有值）
	err = repo.InitDefaults(ctx)
	assert.NoError(t, err)

	// 驗證值未被覆蓋
	item, err := repo.Get(ctx, KeyBatchSize)
	assert.NoError(t, err)
	assert.Equal(t, float64(9999), item.Value)
}

/**
 * TestSQLRepository_WithDescription 測試帶描述的設定
 */
func TestSQLRepository_WithDescription(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	// 建立帶描述的設定
	description := "這是一個測試設定，用於測試描述欄位的儲存和讀取功能。"
	now := time.Now().UTC()

	_, err := db.ExecContext(ctx, `
		INSERT INTO system_settings (key, value, description, updated_at)
		VALUES (?, ?, ?, ?)
	`, "test.with.description", `{"value":123}`, description, now)
	require.NoError(t, err)

	// 驗證描述正確儲存
	item, err := repo.Get(ctx, "test.with.description")
	assert.NoError(t, err)
	assert.Equal(t, description, item.Description)
}

/**
 * TestSQLRepository_JSONParsing 測試 JSON 解析
 */
func TestSQLRepository_JSONParsing(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	// 測試有效的 JSON 值
	validJSON := `{"key1":"value1","key2":123,"key3":true}`
	err := repo.Set(ctx, "json.valid", validJSON)
	require.NoError(t, err)

	item, err := repo.Get(ctx, "json.valid")
	assert.NoError(t, err)
	assert.NotNil(t, item.Value)

	// 測試無效的 JSON 值（應使用原始字串）
	invalidJSON := `not a valid json {`
	err = repo.Set(ctx, "json.invalid", invalidJSON)
	require.NoError(t, err)

	item, err = repo.Get(ctx, "json.invalid")
	assert.NoError(t, err)
	assert.Equal(t, invalidJSON, item.Value)
}

/**
 * TestSQLRepository_Concurrency 測試並發存取
 */
func TestSQLRepository_Concurrency(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	db.SetMaxOpenConns(1)

	repo := NewSQLRepository(db)
	ctx := context.Background()

	// 並發設定多個值
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(index int) {
			key := "concurrent.setting." + string(rune('0'+index))
			err := repo.Set(ctx, key, index)
			assert.NoError(t, err)
			done <- true
		}(i)
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 驗證所有值都設定成功
	for i := 0; i < 10; i++ {
		key := "concurrent.setting." + string(rune('0'+i))
		_, err := repo.Get(ctx, key)
		assert.NoError(t, err)
	}
}

/**
 * TestSQLRepository_UpdatedAt 測試更新時間
 */
func TestSQLRepository_UpdatedAt(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	// 建立設定
	err := repo.Set(ctx, "time.test", "initial")
	require.NoError(t, err)

	// 取得初始更新時間
	item1, err := repo.Get(ctx, "time.test")
	require.NoError(t, err)
	initialTime := item1.UpdatedAt

	// 等待一秒
	time.Sleep(time.Second)

	// 更新設定
	err = repo.Set(ctx, "time.test", "updated")
	require.NoError(t, err)

	// 取得更新後的時間
	item2, err := repo.Get(ctx, "time.test")
	require.NoError(t, err)
	updatedTime := item2.UpdatedAt

	// 驗證時間已更新
	assert.True(t, updatedTime.After(initialTime), "UpdatedAt should be updated")
}

/**
 * TestSQLRepository_SpecialCharacters 測試特殊字元的 Key
 */
func TestSQLRepository_SpecialCharacters(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	// 測試各種特殊字元
	keys := []string{
		"setting.with.dots",
		"setting-with-dashes",
		"setting_with_underscores",
		"setting.with.multiple.dots.and.dashes",
		"Setting.With.CamelCase",
	}

	for _, key := range keys {
		err := repo.Set(ctx, key, key)
		assert.NoError(t, err, "Failed to set key: %s", key)

		item, err := repo.Get(ctx, key)
		assert.NoError(t, err, "Failed to get key: %s", key)
		assert.Equal(t, key, item.Key)
	}
}

/**
 * TestSQLRepository_NumericKeys 測試數字類型的 Key
 */
func TestSQLRepository_NumericKeys(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	// 測試數字 Key
	err := repo.Set(ctx, "12345", "value for numeric key")
	assert.NoError(t, err)

	item, err := repo.Get(ctx, "12345")
	assert.NoError(t, err)
	assert.Equal(t, "12345", item.Key)
}

/**
 * TestSQLRepository_EmptyStringKey 測試空字串 Key
 */
func TestSQLRepository_EmptyStringKey(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	// 設定空字串 Key（應允許）
	err := repo.Set(ctx, "", "value for empty key")
	assert.NoError(t, err)

	item, err := repo.Get(ctx, "")
	assert.NoError(t, err)
	assert.Equal(t, "", item.Key)
}

/**
 * TestSQLRepository_VeryLongValue 測試非常長的值
 */
func TestSQLRepository_VeryLongValue(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	// 建立非常長的字串
	longString := ""
	for i := 0; i < 1000; i++ {
		longString += "這是一個非常長的字串。"
	}

	err := repo.Set(ctx, "long.value", longString)
	assert.NoError(t, err)

	item, err := repo.Get(ctx, "long.value")
	assert.NoError(t, err)
	assert.Equal(t, longString, item.Value)
}
