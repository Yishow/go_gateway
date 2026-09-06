package device

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-gateway/internal/datalink/schema"
)

// SQLRepository implements the Repository interface using SQL
type SQLRepository struct {
	db *sql.DB
}

// NewSQLRepository creates a new SQL repository
func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) Create(ctx context.Context, dev *schema.Device) error {
	query := `
		INSERT INTO devices (
			id, name, description, protocol, status, connection_config,
			last_test_at, last_test_success, last_test_error, readiness_status, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	// 處理可選欄位：將 nil 指針轉換為適當的 SQL 值
	var lastTestAt interface{}
	if dev.LastTestAt != nil {
		// SQLite 使用 TEXT 存儲時間，格式為 "2006-01-02 15:04:05"
		lastTestAt = dev.LastTestAt.Format("2006-01-02 15:04:05")
	} else {
		lastTestAt = nil
	}

	var lastTestSuccess interface{}
	if dev.LastTestSuccess != nil {
		// SQLite 使用 INTEGER 存儲布林值 (0 或 1)
		if *dev.LastTestSuccess {
			lastTestSuccess = 1
		} else {
			lastTestSuccess = 0
		}
	} else {
		lastTestSuccess = nil
	}

	var lastTestError interface{}
	if dev.LastTestError != "" {
		lastTestError = dev.LastTestError
	} else {
		lastTestError = nil
	}

	// 格式化時間為 SQLite 格式
	createdAtStr := dev.CreatedAt.Format("2006-01-02 15:04:05")
	updatedAtStr := dev.UpdatedAt.Format("2006-01-02 15:04:05")

	_, err := r.db.ExecContext(ctx, query,
		dev.ID,
		dev.Name,
		dev.Description,
		dev.Protocol,
		dev.Status,
		dev.ConnectionConfig,
		lastTestAt,
		lastTestSuccess,
		lastTestError,
		dev.ReadinessStatus,
		createdAtStr,
		updatedAtStr,
	)
	if err != nil {
		return fmt.Errorf("failed to create device: %w", err)
	}
	return nil
}

func (r *SQLRepository) GetByID(ctx context.Context, id string) (*schema.Device, error) {
	query := `
		SELECT id, name, description, protocol, status, connection_config,
			   last_test_at, last_test_success, last_test_error, readiness_status, created_at, updated_at
		FROM devices WHERE id = ?
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var dev schema.Device
	var lastTestAt sql.NullString
	var lastTestSuccess sql.NullBool
	var lastTestError sql.NullString
	var readinessStatus sql.NullString
	var createdAtStr, updatedAtStr string

	err := row.Scan(
		&dev.ID,
		&dev.Name,
		&dev.Description,
		&dev.Protocol,
		&dev.Status,
		&dev.ConnectionConfig,
		&lastTestAt,
		&lastTestSuccess,
		&lastTestError,
		&readinessStatus,
		&createdAtStr,
		&updatedAtStr,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("device not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	dev.CreatedAt = parseTime(createdAtStr)
	dev.UpdatedAt = parseTime(updatedAtStr)

	if lastTestAt.Valid {
		t := parseTime(lastTestAt.String)
		dev.LastTestAt = &t
	}
	if lastTestSuccess.Valid {
		dev.LastTestSuccess = &lastTestSuccess.Bool
	}
	if lastTestError.Valid {
		dev.LastTestError = lastTestError.String
	}
	if readinessStatus.Valid {
		dev.ReadinessStatus = readinessStatus.String
	}

	return &dev, nil
}

func (r *SQLRepository) List(ctx context.Context, filter ListFilter) ([]*schema.Device, error) {
	query := `
		SELECT id, name, description, protocol, status, connection_config,
			   last_test_at, last_test_success, last_test_error, readiness_status, created_at, updated_at
		FROM devices WHERE 1=1
	`
	var args []interface{}

	if filter.Protocol != nil {
		query += " AND protocol = ?"
		args = append(args, *filter.Protocol)
	}
	if filter.Status != nil {
		query += " AND status = ?"
		args = append(args, *filter.Status)
	}

	// Always order by name or created_at for consistency
	query += " ORDER BY name ASC"

	if filter.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, filter.Limit)
	}
	if filter.Offset > 0 {
		query += " OFFSET ?"
		args = append(args, filter.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list devices: %w", err)
	}
	defer rows.Close()

	var devices []*schema.Device
	for rows.Next() {
		var dev schema.Device
		var lastTestAt sql.NullString
		var lastTestSuccess sql.NullBool
		var lastTestError sql.NullString
		var readinessStatus sql.NullString
		var createdAtStr, updatedAtStr string

		err := rows.Scan(
			&dev.ID,
			&dev.Name,
			&dev.Description,
			&dev.Protocol,
			&dev.Status,
			&dev.ConnectionConfig,
			&lastTestAt,
			&lastTestSuccess,
			&lastTestError,
			&readinessStatus,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan device: %w", err)
		}

		dev.CreatedAt = parseTime(createdAtStr)
		dev.UpdatedAt = parseTime(updatedAtStr)

		if lastTestAt.Valid {
			t := parseTime(lastTestAt.String)
			dev.LastTestAt = &t
		}
		if lastTestSuccess.Valid {
			dev.LastTestSuccess = &lastTestSuccess.Bool
		}
		if lastTestError.Valid {
			dev.LastTestError = lastTestError.String
		}
		if readinessStatus.Valid {
			dev.ReadinessStatus = readinessStatus.String
		}

		devices = append(devices, &dev)
	}

	return devices, nil
}

func (r *SQLRepository) Update(ctx context.Context, dev *schema.Device) error {
	query := `
		UPDATE devices 
		SET name=?, description=?, protocol=?, status=?, connection_config=?, readiness_status=?, updated_at=?
		WHERE id=?
	`

	_, err := r.db.ExecContext(ctx, query,
		dev.Name,
		dev.Description,
		dev.Protocol,
		dev.Status,
		dev.ConnectionConfig,
		dev.ReadinessStatus,
		time.Now(), // Update updated_at
		dev.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update device: %w", err)
	}
	return nil
}

func (r *SQLRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM devices WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete device: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("device not found (delete failed)")
	}

	return nil
}

func (r *SQLRepository) UpdateTestResult(ctx context.Context, id string, success bool, errMsg string) error {
	query := `
		UPDATE devices
		SET last_test_at=?, last_test_success=?, last_test_error=?, updated_at=?
		WHERE id=?
	`
	_, err := r.db.ExecContext(ctx, query,
		time.Now(),
		success,
		errMsg,
		time.Now(),
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to update test result: %w", err)
	}
	return nil
}

// ClearTestResult 清除已失效的連線測試結果
func (r *SQLRepository) ClearTestResult(ctx context.Context, id string) error {
	query := `
		UPDATE devices
		SET last_test_at=NULL, last_test_success=NULL, last_test_error='', updated_at=?
		WHERE id=?
	`
	if _, err := r.db.ExecContext(ctx, query, time.Now(), id); err != nil {
		return fmt.Errorf("failed to clear test result: %w", err)
	}
	return nil
}

func (r *SQLRepository) UpdateStatus(ctx context.Context, id string, status schema.DeviceStatus) error {
	query := `
		UPDATE devices
		SET status=?, updated_at=?
		WHERE id=?
	`
	_, err := r.db.ExecContext(ctx, query,
		status,
		time.Now(),
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}
	return nil
}

// parseTime tries to parse sqlite time strings
func parseTime(s string) time.Time {
	// Try standard SQL format first
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err == nil {
		return t // Assuming UTC or local based on context. Hsl likely assumes local.
	}
	// Try RFC3339
	t, err = time.Parse(time.RFC3339, s)
	if err == nil {
		return t
	}
	// If empty string or unparseable, return zero time
	return time.Time{}
}
