// Package mapping 提供 SQL Repository 的單元測試。
package mapping

import (
	"context"
	"database/sql"
	"fmt"
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
 * setupTestData 建立測試所需的基礎資料（設備、點位、標籤）
 * @param t 測試實例
 * @param db 資料庫連線
 * @returns pointID, tagID 點位 ID 與標籤 ID
 */
func setupTestData(t *testing.T, db *sql.DB) (pointIDResult, tagIDResult string) {
	ctx := context.Background()
	now := time.Now().UTC()

	// 建立設備
	deviceID := "test-device-001"
	_, err := db.ExecContext(ctx, `
		INSERT INTO devices (id, name, protocol, status, connection_config, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, deviceID, "測試設備", schema.ProtocolModbusTCP, schema.DeviceStatusActive, "{}", now, now)
	require.NoError(t, err)

	// 建立點位
	pointID := "test-point-001"
	_, err = db.ExecContext(ctx, `
		INSERT INTO points (id, device_id, name, address, data_type, mode, enabled, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, pointID, deviceID, "測試點位", "40001", schema.DataTypeInt16, schema.PointModeReadOnly, 1, now, now)
	require.NoError(t, err)

	// 建立標籤
	tagID := "test-tag-001"
	_, err = db.ExecContext(ctx, `
		INSERT INTO tags (id, key, key_lower, display_name, data_type, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, tagID, "test.tag.001", "test.tag.001", "測試標籤 001", schema.DataTypeFloat64, schema.TagStatusActive, now, now)
	require.NoError(t, err)

	return pointID, tagID
}

/**
 * createTestMapping 建立測試用的映射資料
 * @param t 測試實例
 * @param pointID 點位 ID
 * @param tagID 標籤 ID
 * @returns *schema.Mapping 測試映射物件
 */
func createTestMapping(t *testing.T, pointID, tagID string) *schema.Mapping {
	now := time.Now().UTC()
	id, err := common.NewUUID()
	require.NoError(t, err)
	pipeline := `[
		{"type": "scale", "params": {"multiplier": 10, "offset": 0}},
		{"type": "round", "params": {"decimals": 2}}
	]`
	return &schema.Mapping{
		ID:                id,
		PointID:           pointID,
		TagID:             tagID,
		TransformPipeline: pipeline,
		Enabled:           true,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

/**
 * TestSQLRepository_Create 測試建立映射功能
 */
func TestSQLRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	pointID, tagID := setupTestData(t, db)
	mapping := createTestMapping(t, pointID, tagID)

	// 執行 Create
	err := repo.Create(ctx, mapping)
	assert.NoError(t, err)
	assert.NotEmpty(t, mapping.ID)

	// 驗證映射已建立
	retrieved, err := repo.GetByID(ctx, mapping.ID)
	assert.NoError(t, err)
	assert.Equal(t, mapping.PointID, retrieved.PointID)
	assert.Equal(t, mapping.TagID, retrieved.TagID)
	assert.Equal(t, mapping.TransformPipeline, retrieved.TransformPipeline)
}

/**
 * TestSQLRepository_Create_Duplicate 測試建立重複映射應失敗
 */
func TestSQLRepository_Create_Duplicate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	pointID, tagID := setupTestData(t, db)
	mapping := createTestMapping(t, pointID, tagID)

	// 第一次建立應成功
	err := repo.Create(ctx, mapping)
	require.NoError(t, err)

	// 第二次建立相同 point_id 和 tag_id 應失敗（透過 UNIQUE 約束）
	mapping2 := createTestMapping(t, pointID, tagID)
	mapping2.ID = "" // 清空 ID 以便生成新 ID
	err = repo.Create(ctx, mapping2)
	assert.Error(t, err)
}

/**
 * TestSQLRepository_Update 測試更新映射功能
 */
func TestSQLRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	pointID, tagID := setupTestData(t, db)
	mapping := createTestMapping(t, pointID, tagID)
	err := repo.Create(ctx, mapping)
	require.NoError(t, err)

	// 更新映射資料
	newPipeline := `[
		{"type": "scale", "params": {"multiplier": 100, "offset": 5}},
		{"type": "clamp", "params": {"min": 0, "max": 100}}
	]`
	mapping.TransformPipeline = newPipeline
	mapping.Enabled = false

	err = repo.Update(ctx, mapping)
	assert.NoError(t, err)

	// 驗證更新結果
	retrieved, err := repo.GetByID(ctx, mapping.ID)
	assert.NoError(t, err)
	assert.Equal(t, newPipeline, retrieved.TransformPipeline)
	assert.False(t, retrieved.Enabled)
}

/**
 * TestSQLRepository_Update_NotFound 測試更新不存在的映射應失敗
 */
func TestSQLRepository_Update_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	pointID, tagID := setupTestData(t, db)
	mapping := createTestMapping(t, pointID, tagID)
	mapping.ID = "non-existent-id"

	err := repo.Update(ctx, mapping)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "映射不存在")
}

/**
 * TestSQLRepository_Delete 測試刪除映射功能
 */
func TestSQLRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	pointID, tagID := setupTestData(t, db)
	mapping := createTestMapping(t, pointID, tagID)
	err := repo.Create(ctx, mapping)
	require.NoError(t, err)

	// 刪除映射
	err = repo.Delete(ctx, mapping.ID)
	assert.NoError(t, err)

	// 驗證映射已刪除
	_, err = repo.GetByID(ctx, mapping.ID)
	assert.Error(t, err)
}

/**
 * TestSQLRepository_Delete_NotFound 測試刪除不存在的映射應失敗
 */
func TestSQLRepository_Delete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	err := repo.Delete(ctx, "non-existent-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "映射不存在")
}

/**
 * TestSQLRepository_GetByID 測試根據 ID 取得映射
 */
func TestSQLRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	pointID, tagID := setupTestData(t, db)
	mapping := createTestMapping(t, pointID, tagID)
	err := repo.Create(ctx, mapping)
	require.NoError(t, err)

	// 取得映射
	retrieved, err := repo.GetByID(ctx, mapping.ID)
	assert.NoError(t, err)
	assert.Equal(t, mapping.ID, retrieved.ID)
	assert.Equal(t, mapping.PointID, retrieved.PointID)
	assert.Equal(t, mapping.TagID, retrieved.TagID)
}

/**
 * TestSQLRepository_GetByID_NotFound 測試取得不存在的映射
 */
func TestSQLRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	_, err := repo.GetByID(ctx, "non-existent-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "映射不存在")
}

/**
 * TestSQLRepository_GetByPointID 測試根據點位 ID 取得映射列表
 */
func TestSQLRepository_GetByPointID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	pointID, _ := setupTestData(t, db)

	// 建立多個標籤和映射
	tagIDs := []string{"tag-001", "tag-002", "tag-003"}
	for _, tagID := range tagIDs {
		// 建立標籤
		_, err := db.ExecContext(ctx, `
			INSERT INTO tags (id, key, key_lower, display_name, data_type, status, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, tagID, tagID, tagID, "測試標籤", schema.DataTypeFloat64, schema.TagStatusActive, time.Now().UTC(), time.Now().UTC())
		require.NoError(t, err)

		// 建立映射
		mapping := createTestMapping(t, pointID, tagID)
		err = repo.Create(ctx, mapping)
		require.NoError(t, err)
	}

	// 根據點位 ID 取得映射
	mappings, err := repo.GetByPointID(ctx, pointID)
	assert.NoError(t, err)
	assert.Len(t, mappings, 3)

	// 驗證所有映射都指向同一點位
	for _, mapping := range mappings {
		assert.Equal(t, pointID, mapping.PointID)
	}
}

/**
 * TestSQLRepository_GetByTagID 測試根據標籤 ID 取得映射列表
 */
func TestSQLRepository_GetByTagID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	_, tagID := setupTestData(t, db)

	// 建立多個點位和映射
	pointIDs := []string{"point-001", "point-002", "point-003"}
	for _, pointID := range pointIDs {
		// 建立點位（需要設備）
		deviceID := "device-" + pointID
		_, err := db.ExecContext(ctx, `
			INSERT INTO devices (id, name, protocol, status, connection_config, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, deviceID, deviceID, schema.ProtocolModbusTCP, schema.DeviceStatusActive, "{}", time.Now().UTC(), time.Now().UTC())
		require.NoError(t, err)

		_, err = db.ExecContext(ctx, `
			INSERT INTO points (id, device_id, name, address, data_type, mode, enabled, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, pointID, deviceID, pointID, "40001", schema.DataTypeInt16, schema.PointModeReadOnly, 1, time.Now().UTC(), time.Now().UTC())
		require.NoError(t, err)

		// 建立映射
		mapping := createTestMapping(t, pointID, tagID)
		err = repo.Create(ctx, mapping)
		require.NoError(t, err)
	}

	// 根據標籤 ID 取得映射
	mappings, err := repo.GetByTagID(ctx, tagID)
	assert.NoError(t, err)
	assert.Len(t, mappings, 3)

	// 驗證所有映射都指向同一標籤
	for _, mapping := range mappings {
		assert.Equal(t, tagID, mapping.TagID)
	}
}

/**
 * TestSQLRepository_List 測試列出映射
 */
func TestSQLRepository_List(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	pointID, tagID := setupTestData(t, db)

	// 建立多個映射
	mapping := createTestMapping(t, pointID, tagID)
	err := repo.Create(ctx, mapping)
	require.NoError(t, err)
	for i := 1; i <= 3; i++ {
		// 建立額外的點位和標籤
		newPointID := "point-00" + string(rune('0'+i))
		newTagID := "tag-00" + string(rune('0'+i))

		// 建立設備
		_, err := db.ExecContext(ctx, `
			INSERT INTO devices (id, name, protocol, status, connection_config, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, "device-"+newPointID, "Device", schema.ProtocolModbusTCP, schema.DeviceStatusActive, "{}", time.Now().UTC(), time.Now().UTC())
		require.NoError(t, err)

		// 建立點位
		address := fmt.Sprintf("4000%d", i)
		_, err = db.ExecContext(ctx, `
			INSERT INTO points (id, device_id, name, address, data_type, mode, enabled, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, newPointID, "device-"+newPointID, "Point", address, schema.DataTypeInt16, schema.PointModeReadOnly, 1, time.Now().UTC(), time.Now().UTC())
		require.NoError(t, err)

		// 建立標籤
		_, err = db.ExecContext(ctx, `
			INSERT INTO tags (id, key, key_lower, display_name, data_type, status, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, newTagID, newTagID, newTagID, "Tag", schema.DataTypeFloat64, schema.TagStatusActive, time.Now().UTC(), time.Now().UTC())
		require.NoError(t, err)

		// 建立映射
		newMapping := createTestMapping(t, newPointID, newTagID)
		err = repo.Create(ctx, newMapping)
		require.NoError(t, err)
	}

	// 列出所有映射
	result, err := repo.List(ctx, ListFilter{})
	assert.NoError(t, err)
	assert.Len(t, result, 4)
}

/**
 * TestSQLRepository_ReadsNullLifecycleColumns 測試讀取可為 NULL 的 lifecycle 欄位。
 */
func TestSQLRepository_ReadsNullLifecycleColumns(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	pointID, tagID := setupTestData(t, db)
	now := time.Now().UTC()

	_, err := db.ExecContext(ctx, `
		INSERT INTO mappings (
			id, point_id, tag_id, transform_pipeline, status, rule_candidate_id,
			proposed_signature, last_applied_signature, blocking_reason, enabled, created_at, updated_at
		)
		VALUES (?, ?, ?, ?, ?, NULL, NULL, NULL, NULL, ?, ?, ?)
	`, "mapping-null-lifecycle", pointID, tagID, `[]`, schema.MappingStatusActive, 1, now, now)
	require.NoError(t, err)

	result, err := repo.List(ctx, ListFilter{})
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Empty(t, result[0].RuleCandidateID)
	assert.Empty(t, result[0].ProposedSignature)
	assert.Empty(t, result[0].LastAppliedSignature)
	assert.Empty(t, result[0].BlockingReason)

	retrieved, err := repo.GetByID(ctx, "mapping-null-lifecycle")
	require.NoError(t, err)
	assert.Empty(t, retrieved.RuleCandidateID)
	assert.Empty(t, retrieved.ProposedSignature)
	assert.Empty(t, retrieved.LastAppliedSignature)
	assert.Empty(t, retrieved.BlockingReason)
}

/**
 * TestSQLRepository_List_WithFilter 測試使用過濾條件列出映射
 */
func TestSQLRepository_List_WithFilter(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	pointID, tagID := setupTestData(t, db)

	// 建立多個映射
	mappings := []struct {
		PointID string
		TagID   string
		Enabled bool
	}{
		{pointID, tagID, true},
		{pointID, "other-tag-001", false},
		{"other-point-001", tagID, true},
	}

	// 建立額外資料
	for i, m := range mappings {
		pid := m.PointID
		tid := m.TagID

		if pid != pointID {
			// 建立設備
			_, err := db.ExecContext(ctx, `
				INSERT INTO devices (id, name, protocol, status, connection_config, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, ?, ?)
			`, "device-"+pid, "Device", schema.ProtocolModbusTCP, schema.DeviceStatusActive, "{}", time.Now().UTC(), time.Now().UTC())
			require.NoError(t, err)

			// 建立點位
			address := fmt.Sprintf("4100%d", i)
			_, err = db.ExecContext(ctx, `
				INSERT INTO points (id, device_id, name, address, data_type, mode, enabled, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
			`, pid, "device-"+pid, "Point", address, schema.DataTypeInt16, schema.PointModeReadOnly, 1, time.Now().UTC(), time.Now().UTC())
			require.NoError(t, err)
		}

		if tid != tagID {
			// 建立標籤
			_, err := db.ExecContext(ctx, `
				INSERT INTO tags (id, key, key_lower, display_name, data_type, status, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?)
			`, tid, tid, tid, "Tag", schema.DataTypeFloat64, schema.TagStatusActive, time.Now().UTC(), time.Now().UTC())
			require.NoError(t, err)
		}

		// 建立映射
		newMapping := createTestMapping(t, pid, tid)
		newMapping.Enabled = m.Enabled
		err := repo.Create(ctx, newMapping)
		require.NoError(t, err)
	}

	// 測試點位過濾
	result, err := repo.List(ctx, ListFilter{PointID: &pointID})
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	// 測試標籤過濾
	result, err = repo.List(ctx, ListFilter{TagID: &tagID})
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	// 測試啟用過濾
	enabled := true
	result, err = repo.List(ctx, ListFilter{Enabled: &enabled})
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	// 測試停用過濾
	enabled = false
	result, err = repo.List(ctx, ListFilter{Enabled: &enabled})
	assert.NoError(t, err)
	assert.Len(t, result, 1)

	// 測試組合過濾
	enabled = true
	result, err = repo.List(ctx, ListFilter{
		PointID: &pointID,
		TagID:   &tagID,
		Enabled: &enabled,
	})
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

/**
 * TestSQLRepository_List_WithPagination 測試分頁功能
 */
func TestSQLRepository_List_WithPagination(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	setupTestData(t, db)

	// 建立 5 個映射
	for i := 0; i < 5; i++ {
		// 建立額外的點位和標籤
		newPointID := "point-00" + string(rune('0'+i))
		newTagID := "tag-00" + string(rune('0'+i))

		// 建立設備
		_, err := db.ExecContext(ctx, `
			INSERT INTO devices (id, name, protocol, status, connection_config, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, "device-"+newPointID, "Device", schema.ProtocolModbusTCP, schema.DeviceStatusActive, "{}", time.Now().UTC(), time.Now().UTC())
		require.NoError(t, err)

		// 建立點位
		address := fmt.Sprintf("4300%d", i)
		_, err = db.ExecContext(ctx, `
			INSERT INTO points (id, device_id, name, address, data_type, mode, enabled, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, newPointID, "device-"+newPointID, "Point", address, schema.DataTypeInt16, schema.PointModeReadOnly, 1, time.Now().UTC(), time.Now().UTC())
		require.NoError(t, err)

		// 建立標籤
		_, err = db.ExecContext(ctx, `
			INSERT INTO tags (id, key, key_lower, display_name, data_type, status, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, newTagID, newTagID, newTagID, "Tag", schema.DataTypeFloat64, schema.TagStatusActive, time.Now().UTC(), time.Now().UTC())
		require.NoError(t, err)

		// 建立映射
		mapping := createTestMapping(t, newPointID, newTagID)
		err = repo.Create(ctx, mapping)
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
 * TestSQLRepository_EnabledStatus 測試啟用狀態
 */
func TestSQLRepository_EnabledStatus(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	pointID, tagID := setupTestData(t, db)

	// 建立啟用和停用的映射
	enabledMapping := createTestMapping(t, pointID, tagID)
	enabledMapping.Enabled = true
	err := repo.Create(ctx, enabledMapping)
	require.NoError(t, err)

	// 建立其他標籤和映射
	otherTagID := "other-tag-001"
	_, err = db.ExecContext(ctx, `
		INSERT INTO tags (id, key, key_lower, display_name, data_type, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, otherTagID, otherTagID, otherTagID, "Other Tag", schema.DataTypeFloat64, schema.TagStatusActive, time.Now().UTC(), time.Now().UTC())
	require.NoError(t, err)

	disabledMapping := createTestMapping(t, pointID, otherTagID)
	disabledMapping.Enabled = false
	err = repo.Create(ctx, disabledMapping)
	require.NoError(t, err)

	// 驗證啟用狀態
	enabled := true
	result, err := repo.List(ctx, ListFilter{Enabled: &enabled})
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.True(t, result[0].Enabled)

	// 驗證停用狀態
	enabled = false
	result, err = repo.List(ctx, ListFilter{Enabled: &enabled})
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.False(t, result[0].Enabled)
}

/**
 * TestSQLRepository_TransformPipeline 測試轉換管線的 JSON 儲存
 */
func TestSQLRepository_TransformPipeline(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	pointID, tagID := setupTestData(t, db)

	// 測試複雜的轉換管線
	complexPipeline := `[
		{"type": "scale", "params": {"multiplier": 10, "offset": 0}},
		{"type": "offset", "params": {"offset": 5}},
		{"type": "clamp", "params": {"min": 0, "max": 100}},
		{"type": "round", "params": {"decimals": 2}}
	]`

	mapping := createTestMapping(t, pointID, tagID)
	mapping.TransformPipeline = complexPipeline

	err := repo.Create(ctx, mapping)
	require.NoError(t, err)

	// 驗證轉換管線正確儲存
	retrieved, err := repo.GetByID(ctx, mapping.ID)
	assert.NoError(t, err)
	assert.Equal(t, complexPipeline, retrieved.TransformPipeline)
}

/**
 * TestSQLRepository_UniquenessConstraint 測試點位與標籤的唯一性約束
 */
func TestSQLRepository_UniquenessConstraint(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	pointID, tagID := setupTestData(t, db)

	// 建立第一個映射
	mapping1 := createTestMapping(t, pointID, tagID)
	err := repo.Create(ctx, mapping1)
	require.NoError(t, err)

	// 建立相同點位和標籤的第二個映射應失敗
	mapping2 := createTestMapping(t, pointID, tagID)
	err = repo.Create(ctx, mapping2)
	assert.Error(t, err)

	// 不同點位但相同標籤應成功
	otherPointID := "other-point"
	_, err = db.ExecContext(ctx, `
		INSERT INTO devices (id, name, protocol, status, connection_config, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, "other-device", "Other Device", schema.ProtocolModbusTCP, schema.DeviceStatusActive, "{}", time.Now().UTC(), time.Now().UTC())
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		INSERT INTO points (id, device_id, name, address, data_type, mode, enabled, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, otherPointID, "other-device", "Other Point", "40001", schema.DataTypeInt16, schema.PointModeReadOnly, 1, time.Now().UTC(), time.Now().UTC())
	require.NoError(t, err)

	mapping3 := createTestMapping(t, otherPointID, tagID)
	err = repo.Create(ctx, mapping3)
	assert.NoError(t, err)

	// 相同點位但不同標籤應成功
	otherTagID := "other-tag"
	_, err = db.ExecContext(ctx, `
		INSERT INTO tags (id, key, key_lower, display_name, data_type, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, otherTagID, otherTagID, otherTagID, "Other Tag", schema.DataTypeFloat64, schema.TagStatusActive, time.Now().UTC(), time.Now().UTC())
	require.NoError(t, err)

	mapping4 := createTestMapping(t, pointID, otherTagID)
	err = repo.Create(ctx, mapping4)
	assert.NoError(t, err)
}
