// Package point 提供點位管理功能的 SQL Repository 實現。
package point

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// SQL Repository 實現
// =============================================================================

// SQLRepository SQL 點位儲存庫
type SQLRepository struct {
	db *sql.DB
}

// NewSQLRepository 建立新的 SQL 儲存庫
func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

// Create 建立點位
func (r *SQLRepository) Create(ctx context.Context, point *schema.Point) error {
	if point.ID == "" {
		id, err := common.NewUUID()
		if err != nil {
			return fmt.Errorf("建立點位 ID 失敗: %w", err)
		}
		point.ID = id
	}

	query := `
		INSERT INTO points (id, device_id, name, description, address, function, data_type, mode, polling_group_id, enabled, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		point.ID,
		point.DeviceID,
		point.Name,
		point.Description,
		point.Address,
		point.Function,
		point.DataType,
		point.Mode,
		point.PollingGroupID,
		point.Enabled,
		point.CreatedAt,
		point.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("建立點位失敗: %w", err)
	}

	return nil
}

// Update 更新點位
func (r *SQLRepository) Update(ctx context.Context, point *schema.Point) error {
	query := `
		UPDATE points 
		SET name = ?, description = ?, address = ?, function = ?, data_type = ?, mode = ?, 
		    polling_group_id = ?, last_read_at = ?, last_value = ?, last_error = ?, enabled = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		point.Name,
		point.Description,
		point.Address,
		point.Function,
		point.DataType,
		point.Mode,
		point.PollingGroupID,
		point.LastReadAt,
		point.LastValue,
		point.LastError,
		point.Enabled,
		time.Now(),
		point.ID,
	)

	if err != nil {
		return fmt.Errorf("更新點位失敗: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("點位不存在: %s", point.ID)
	}

	return nil
}

// Delete 刪除點位
func (r *SQLRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM points WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("刪除點位失敗: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("點位不存在: %s", id)
	}

	return nil
}

// GetByID 根據 ID 取得點位
func (r *SQLRepository) GetByID(ctx context.Context, id string) (*schema.Point, error) {
	query := `
		SELECT id, device_id, name, description, address, function, data_type, mode, 
		       polling_group_id, last_read_at, last_value, last_error, enabled, created_at, updated_at
		FROM points WHERE id = ?
	`

	row := r.db.QueryRowContext(ctx, query, id)
	return r.scanPoint(row)
}

// ListByDevice 根據設備 ID 取得點位列表
func (r *SQLRepository) ListByDevice(ctx context.Context, deviceID string) ([]*schema.Point, error) {
	query := `
		SELECT id, device_id, name, description, address, function, data_type, mode,
		       polling_group_id, last_read_at, last_value, last_error, enabled, created_at, updated_at
		FROM points WHERE device_id = ?
		ORDER BY name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, deviceID)
	if err != nil {
		return nil, fmt.Errorf("查詢點位失敗: %w", err)
	}
	defer rows.Close()

	points := make([]*schema.Point, 0)
	for rows.Next() {
		point, err := r.scanPointFromRows(rows)
		if err != nil {
			return nil, err
		}
		points = append(points, point)
	}

	return points, nil
}

// ListByPollingGroup 根據輪詢群組 ID 取得點位列表
func (r *SQLRepository) ListByPollingGroup(ctx context.Context, groupID string) ([]*schema.Point, error) {
	query := `
		SELECT id, device_id, name, description, address, function, data_type, mode,
		       polling_group_id, last_read_at, last_value, last_error, enabled, created_at, updated_at
		FROM points WHERE polling_group_id = ?
		ORDER BY name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, fmt.Errorf("查詢點位失敗: %w", err)
	}
	defer rows.Close()

	points := make([]*schema.Point, 0)
	for rows.Next() {
		point, err := r.scanPointFromRows(rows)
		if err != nil {
			return nil, err
		}
		points = append(points, point)
	}

	return points, nil
}

// List 列出點位
func (r *SQLRepository) List(ctx context.Context, filter ListFilter) ([]*schema.Point, error) {
	query := `
		SELECT id, device_id, name, description, address, function, data_type, mode,
		       polling_group_id, last_read_at, last_value, last_error, enabled, created_at, updated_at
		FROM points
		WHERE 1=1
	`
	args := []interface{}{}

	if filter.DeviceID != nil {
		query += ` AND device_id = ?`
		args = append(args, *filter.DeviceID)
	}

	if filter.PollingGroupID != nil {
		query += ` AND polling_group_id = ?`
		args = append(args, *filter.PollingGroupID)
	}

	query += ` ORDER BY name ASC`

	if filter.Limit > 0 {
		query += fmt.Sprintf(` LIMIT %d`, filter.Limit)
	}
	if filter.Offset > 0 {
		query += fmt.Sprintf(` OFFSET %d`, filter.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("查詢點位失敗: %w", err)
	}
	defer rows.Close()

	points := make([]*schema.Point, 0)
	for rows.Next() {
		point, err := r.scanPointFromRows(rows)
		if err != nil {
			return nil, err
		}
		points = append(points, point)
	}

	return points, nil
}

// scanPoint 從單一行掃描點位
func (r *SQLRepository) scanPoint(row *sql.Row) (*schema.Point, error) {
	var point schema.Point
	var description, function, lastValue, lastError sql.NullString
	var pollingGroupID sql.NullString
	var lastReadAt sql.NullString
	var createdAt, updatedAt string

	err := row.Scan(
		&point.ID,
		&point.DeviceID,
		&point.Name,
		&description,
		&point.Address,
		&function,
		&point.DataType,
		&point.Mode,
		&pollingGroupID,
		&lastReadAt,
		&lastValue,
		&lastError,
		&point.Enabled,
		&createdAt,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("點位不存在")
	}
	if err != nil {
		return nil, fmt.Errorf("掃描點位失敗: %w", err)
	}

	if description.Valid {
		point.Description = description.String
	}
	if function.Valid {
		point.Function = function.String
	}
	if pollingGroupID.Valid {
		point.PollingGroupID = &pollingGroupID.String
	}
	if lastReadAt.Valid {
		parsed, err := common.ParseTimeString(lastReadAt.String)
		if err != nil {
			return nil, fmt.Errorf("解析最後讀取時間失敗: %w", err)
		}
		point.LastReadAt = &parsed
	}
	if lastValue.Valid {
		point.LastValue = &lastValue.String
	} else {
		point.LastValue = nil
	}
	if lastError.Valid {
		point.LastError = lastError.String
	} else {
		point.LastError = ""
	}

	parsedCreatedAt, err := common.ParseTimeString(createdAt)
	if err != nil {
		return nil, fmt.Errorf("解析建立時間失敗: %w", err)
	}
	parsedUpdatedAt, err := common.ParseTimeString(updatedAt)
	if err != nil {
		return nil, fmt.Errorf("解析更新時間失敗: %w", err)
	}
	point.CreatedAt = parsedCreatedAt
	point.UpdatedAt = parsedUpdatedAt

	return &point, nil
}

// scanPointFromRows 從多行結果掃描點位
func (r *SQLRepository) scanPointFromRows(rows *sql.Rows) (*schema.Point, error) {
	var point schema.Point
	var description, function, lastValue, lastError sql.NullString
	var pollingGroupID sql.NullString
	var lastReadAt sql.NullString
	var createdAt, updatedAt string

	err := rows.Scan(
		&point.ID,
		&point.DeviceID,
		&point.Name,
		&description,
		&point.Address,
		&function,
		&point.DataType,
		&point.Mode,
		&pollingGroupID,
		&lastReadAt,
		&lastValue,
		&lastError,
		&point.Enabled,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("掃描點位失敗: %w", err)
	}

	if description.Valid {
		point.Description = description.String
	}
	if function.Valid {
		point.Function = function.String
	}
	if pollingGroupID.Valid {
		point.PollingGroupID = &pollingGroupID.String
	}
	if lastReadAt.Valid {
		parsed, err := common.ParseTimeString(lastReadAt.String)
		if err != nil {
			return nil, fmt.Errorf("解析最後讀取時間失敗: %w", err)
		}
		point.LastReadAt = &parsed
	}
	if lastValue.Valid {
		point.LastValue = &lastValue.String
	} else {
		point.LastValue = nil
	}
	if lastError.Valid {
		point.LastError = lastError.String
	} else {
		point.LastError = ""
	}

	parsedCreatedAt, err := common.ParseTimeString(createdAt)
	if err != nil {
		return nil, fmt.Errorf("解析建立時間失敗: %w", err)
	}
	parsedUpdatedAt, err := common.ParseTimeString(updatedAt)
	if err != nil {
		return nil, fmt.Errorf("解析更新時間失敗: %w", err)
	}
	point.CreatedAt = parsedCreatedAt
	point.UpdatedAt = parsedUpdatedAt

	return &point, nil
}
// BatchUpdateReadResult 批次更新點位讀取結果
func (r *SQLRepository) BatchUpdateReadResult(ctx context.Context, results []ReadResultUpdate) error {
	// 由於 SQLite 不支援複雜的批次更新語法，這裡使用交易 + 逐條更新
	// 對於高頻數據，這可能不是最高效的，但在 SQLite 場景下是可行的
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("開啟交易失敗: %w", err)
	}
	defer tx.Rollback() // 如果沒有 Commit，則自動 Rollback

	query := `
		UPDATE points 
		SET last_value = ?, last_read_at = ?, last_error = ?, updated_at = ?
		WHERE id = ?
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("準備更新語句失敗: %w", err)
	}
	defer stmt.Close()

	now := time.Now()

	for _, result := range results {
		var lastValue *string
		var lastError string

		if result.Error != "" {
			lastError = result.Error
		} else {
			strVal := fmt.Sprintf("%v", result.Value)
			lastValue = &strVal
		}

		_, err := stmt.ExecContext(ctx, lastValue, now, lastError, now, result.PointID)
		if err != nil {
			return fmt.Errorf("更新點位 %s 失敗: %w", result.PointID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交交易失敗: %w", err)
	}

	return nil
}

// UpdateReadResult 更新點位讀取結果
func (r *SQLRepository) UpdateReadResult(ctx context.Context, pointID string, value interface{}, errMsg string) error {
	query := `
		UPDATE points 
		SET last_value = ?, last_read_at = ?, last_error = ?, updated_at = ?
		WHERE id = ?
	`

	var lastValue *string
	if errMsg == "" {
		strVal := fmt.Sprintf("%v", value)
		lastValue = &strVal
	}

	result, err := r.db.ExecContext(ctx, query, lastValue, time.Now(), errMsg, time.Now(), pointID)
	if err != nil {
		return fmt.Errorf("更新點位讀取結果失敗: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("點位不存在: %s", pointID)
	}

	return nil
}
