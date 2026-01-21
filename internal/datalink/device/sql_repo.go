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
			last_test_at, last_test_success, last_test_error, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		dev.ID,
		dev.Name,
		dev.Description,
		dev.Protocol,
		dev.Status,
		dev.ConnectionConfig,
		dev.LastTestAt,
		dev.LastTestSuccess,
		dev.LastTestError,
		dev.CreatedAt,
		dev.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create device: %w", err)
	}
	return nil
}

func (r *SQLRepository) GetByID(ctx context.Context, id string) (*schema.Device, error) {
	query := `
		SELECT id, name, description, protocol, status, connection_config,
			   last_test_at, last_test_success, last_test_error, created_at, updated_at
		FROM devices WHERE id = ?
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var dev schema.Device
	var lastTestAt sql.NullString
	var lastTestSuccess sql.NullBool
	var lastTestError sql.NullString
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

	return &dev, nil
}

func (r *SQLRepository) List(ctx context.Context, filter ListFilter) ([]*schema.Device, error) {
	query := `
		SELECT id, name, description, protocol, status, connection_config,
			   last_test_at, last_test_success, last_test_error, created_at, updated_at
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

		devices = append(devices, &dev)
	}

	return devices, nil
}

func (r *SQLRepository) Update(ctx context.Context, dev *schema.Device) error {
	query := `
		UPDATE devices 
		SET name=?, description=?, protocol=?, status=?, connection_config=?, updated_at=?
		WHERE id=?
	`
	
	_, err := r.db.ExecContext(ctx, query,
		dev.Name,
		dev.Description,
		dev.Protocol,
		dev.Status,
		dev.ConnectionConfig,
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
