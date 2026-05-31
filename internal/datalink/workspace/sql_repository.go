package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const storageKey = "studio_v2_workspace"

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) Get(ctx context.Context) (*Record, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT value
		FROM system_settings
		WHERE key = ?
	`, storageKey)

	var raw string
	if err := row.Scan(&raw); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query studio v2 workspace: %w", err)
	}

	var record Record
	if err := json.Unmarshal([]byte(raw), &record); err != nil {
		return nil, fmt.Errorf("decode studio v2 workspace: %w", err)
	}
	if record.OrderedDeviceIDs == nil {
		record.OrderedDeviceIDs = []string{}
	}

	return cloneRecord(&record), nil
}

func (r *SQLRepository) Save(ctx context.Context, record *Record) error {
	payload, err := json.Marshal(cloneRecord(record))
	if err != nil {
		return fmt.Errorf("encode studio v2 workspace: %w", err)
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO system_settings (key, value, updated_at)
		VALUES (?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET
			value = excluded.value,
			updated_at = excluded.updated_at
	`, storageKey, string(payload), time.Now().UTC())
	if err != nil {
		return fmt.Errorf("persist studio v2 workspace: %w", err)
	}

	return nil
}
