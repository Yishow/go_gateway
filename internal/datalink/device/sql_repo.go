// Package device 提供設備管理功能的 SQL Repository 實現。
package device

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

// SQLRepository SQL 設備儲存庫
type SQLRepository struct {
	db *sql.DB
}

// NewSQLRepository 建立新的 SQL 儲存庫
func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

// Create 建立設備
func (r *SQLRepository) Create(ctx context.Context, device *schema.Device) error {
	query := `
		INSERT INTO devices (id, name, description, protocol, status, connection_config, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		device.ID,
		device.Name,
		device.Description,
		device.Protocol,
		device.Status,
		device.ConnectionConfig,
		device.CreatedAt,
		device.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("建立設備失敗: %w", err)
	}

	return nil
}

// Update 更新設備
func (r *SQLRepository) Update(ctx context.Context, device *schema.Device) error {
	query := `
		UPDATE devices 
		SET name = ?, description = ?, protocol = ?, status = ?, 
		    connection_config = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		device.Name,
		device.Description,
		device.Protocol,
		device.Status,
		device.ConnectionConfig,
		time.Now(),
		device.ID,
	)

	if err != nil {
		return fmt.Errorf("更新設備失敗: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("設備不存在: %s", device.ID)
	}

	return nil
}

// Delete 刪除設備
func (r *SQLRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM devices WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("刪除設備失敗: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("設備不存在: %s", id)
	}

	return nil
}

// GetByID 根據 ID 取得設備
func (r *SQLRepository) GetByID(ctx context.Context, id string) (*schema.Device, error) {
	query := `
		SELECT id, name, description, protocol, status, connection_config, 
		       last_test_at, last_test_success, last_test_error, created_at, updated_at
		FROM devices WHERE id = ?
	`

	row := r.db.QueryRowContext(ctx, query, id)
	return r.scanDevice(row)
}

// List 列出設備
func (r *SQLRepository) List(ctx context.Context, filter ListFilter) ([]*schema.Device, error) {
	query := `
		SELECT id, name, description, protocol, status, connection_config,
		       last_test_at, last_test_success, last_test_error, created_at, updated_at
		FROM devices
		WHERE 1=1
	`
	args := []interface{}{}

	if filter.Protocol != nil {
		query += ` AND protocol = ?`
		args = append(args, *filter.Protocol)
	}

	if filter.Status != nil {
		query += ` AND status = ?`
		args = append(args, *filter.Status)
	}

	query += ` ORDER BY created_at DESC`

	if filter.Limit > 0 {
		query += fmt.Sprintf(` LIMIT %d`, filter.Limit)
	}
	if filter.Offset > 0 {
		query += fmt.Sprintf(` OFFSET %d`, filter.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("查詢設備失敗: %w", err)
	}
	defer rows.Close()

	devices := make([]*schema.Device, 0)
	for rows.Next() {
		device, err := r.scanDeviceFromRows(rows)
		if err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}

	return devices, nil
}

// Count 計算設備數量
func (r *SQLRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM devices`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("計算設備數量失敗: %w", err)
	}

	return count, nil
}

// scanDevice 從單一行掃描設備
func (r *SQLRepository) scanDevice(row *sql.Row) (*schema.Device, error) {
	var device schema.Device
	var description, connectionConfig, lastTestError sql.NullString
	var lastTestAt sql.NullTime
	var lastTestSuccess sql.NullBool

	err := row.Scan(
		&device.ID,
		&device.Name,
		&description,
		&device.Protocol,
		&device.Status,
		&connectionConfig,
		&lastTestAt,
		&lastTestSuccess,
		&lastTestError,
		&device.CreatedAt,
		&device.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("設備不存在")
	}
	if err != nil {
		return nil, fmt.Errorf("掃描設備失敗: %w", err)
	}

	if description.Valid {
		device.Description = description.String
	}
	if connectionConfig.Valid {
		device.ConnectionConfig = connectionConfig.String
	}
	if lastTestAt.Valid {
		device.LastTestAt = &lastTestAt.Time
	}
	if lastTestSuccess.Valid {
		device.LastTestSuccess = &lastTestSuccess.Bool
	}
	if lastTestError.Valid {
		device.LastTestError = lastTestError.String
	}

	return &device, nil
}

// scanDeviceFromRows 從多行結果掃描設備
func (r *SQLRepository) scanDeviceFromRows(rows *sql.Rows) (*schema.Device, error) {
	var device schema.Device
	var description, connectionConfig, lastTestError sql.NullString
	var lastTestAt sql.NullTime
	var lastTestSuccess sql.NullBool

	err := rows.Scan(
		&device.ID,
		&device.Name,
		&description,
		&device.Protocol,
		&device.Status,
		&connectionConfig,
		&lastTestAt,
		&lastTestSuccess,
		&lastTestError,
		&device.CreatedAt,
		&device.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("掃描設備失敗: %w", err)
	}

	if description.Valid {
		device.Description = description.String
	}
	if connectionConfig.Valid {
		device.ConnectionConfig = connectionConfig.String
	}
	if lastTestAt.Valid {
		device.LastTestAt = &lastTestAt.Time
	}
	if lastTestSuccess.Valid {
		device.LastTestSuccess = &lastTestSuccess.Bool
	}
	if lastTestError.Valid {
		device.LastTestError = lastTestError.String
	}

	return &device, nil
}
