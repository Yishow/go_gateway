package point

import (
	"context"
	"fmt"
	"time"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
)

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

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("取得更新筆數失敗: %w", err)
	}
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

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("取得刪除筆數失敗: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("點位不存在: %s", id)
	}

	return nil
}
