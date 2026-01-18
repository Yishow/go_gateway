// Package point 提供點位管理功能的 SQL Repository 實現。
package point

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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
		    polling_group_id = ?, enabled = ?, updated_at = ?
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

// GetByDeviceID 根據設備 ID 取得點位列表
func (r *SQLRepository) GetByDeviceID(ctx context.Context, deviceID string) ([]*schema.Point, error) {
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
	var lastReadAt sql.NullTime

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
		&point.CreatedAt,
		&point.UpdatedAt,
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
		point.LastReadAt = &lastReadAt.Time
	}
	if lastValue.Valid {
		point.LastValue = &lastValue.String
	}
	if lastError.Valid {
		point.LastError = lastError.String
	}

	return &point, nil
}

// scanPointFromRows 從多行結果掃描點位
func (r *SQLRepository) scanPointFromRows(rows *sql.Rows) (*schema.Point, error) {
	var point schema.Point
	var description, function, lastValue, lastError sql.NullString
	var pollingGroupID sql.NullString
	var lastReadAt sql.NullTime

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
		&point.CreatedAt,
		&point.UpdatedAt,
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
		point.LastReadAt = &lastReadAt.Time
	}
	if lastValue.Valid {
		point.LastValue = &lastValue.String
	}
	if lastError.Valid {
		point.LastError = lastError.String
	}

	return &point, nil
}
