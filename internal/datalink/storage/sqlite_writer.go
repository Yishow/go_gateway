package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

// SQLiteWriter 將時序資料寫入 SQLite timeseries 表。
type SQLiteWriter struct {
	db *sql.DB
}

// NewSQLiteWriter 建立 SQLite writer。
func NewSQLiteWriter(db *sql.DB) *SQLiteWriter {
	return &SQLiteWriter{db: db}
}

// Write 寫入單筆。
func (w *SQLiteWriter) Write(ctx context.Context, record TimeSeriesRecord) error {
	return w.WriteBatch(ctx, []TimeSeriesRecord{record})
}

// WriteBatch 批次寫入。
func (w *SQLiteWriter) WriteBatch(ctx context.Context, records []TimeSeriesRecord) error {
	if len(records) == 0 {
		return nil
	}

	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx 失敗: %w", err)
	}

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO timeseries (tag_id, ts, value_num, value_text, value_bool, raw_value, quality)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("prepare insert 失敗: %w", err)
	}
	defer stmt.Close()

	for _, r := range records {
		rawVal, err := encodeRawValue(r.RawValue)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("序列化 raw_value 失敗: %w", err)
		}

		if _, err := stmt.ExecContext(ctx,
			r.TagID,
			r.Timestamp,
			r.ValueNum,
			r.ValueText,
			r.ValueBool,
			rawVal,
			r.Quality,
		); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("插入 timeseries 失敗(tag=%s): %w", r.TagID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit 失敗: %w", err)
	}
	return nil
}

// Flush 無內部緩衝，為 no-op。
func (w *SQLiteWriter) Flush(ctx context.Context) error { return nil }

// Close 無持有獨立資源，為 no-op。
func (w *SQLiteWriter) Close() error { return nil }

func encodeRawValue(v interface{}) (interface{}, error) {
	if v == nil {
		return nil, nil
	}
	switch val := v.(type) {
	case string:
		return val, nil
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		return string(b), nil
	}
}
