package dbtarget

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// UpdateWriteOutcome 僅更新資料庫連接器的最後寫入狀態、時間與錯誤摘要，不覆蓋其他連線設定。
func (r *SQLConnectorRepository) UpdateWriteOutcome(ctx context.Context, id string, at time.Time, status, errorSummary string) error {
	query := `UPDATE database_connectors
		SET last_write_at = ?, last_write_status = ?, last_write_error = ?, updated_at = ?
		WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, at.UTC(), strings.TrimSpace(status), strings.TrimSpace(errorSummary), time.Now(), id)
	if err != nil {
		return fmt.Errorf("更新資料庫連接器寫入結果失敗: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("取得資料庫連接器寫入結果更新筆數失敗: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("資料庫連接器不存在")
	}
	return nil
}

// UpdateFlushOutcome 僅更新資料庫連接器的最後 flush 狀態、時間與錯誤摘要。
func (r *SQLConnectorRepository) UpdateFlushOutcome(ctx context.Context, id string, at time.Time, status, errorSummary string) error {
	query := `UPDATE database_connectors
		SET last_flush_at = ?, last_flush_status = ?, last_flush_error = ?, updated_at = ?
		WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, at.UTC(), strings.TrimSpace(status), strings.TrimSpace(errorSummary), time.Now(), id)
	if err != nil {
		return fmt.Errorf("更新資料庫連接器 flush 結果失敗: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("取得資料庫連接器 flush 結果更新筆數失敗: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("資料庫連接器不存在")
	}
	return nil
}

// UpdateSchemaEnsureOutcome 僅更新資料庫連接器的最後 schema ensure 狀態、時間與錯誤摘要。
func (r *SQLConnectorRepository) UpdateSchemaEnsureOutcome(ctx context.Context, id string, at time.Time, status, errorSummary string) error {
	query := `UPDATE database_connectors
		SET last_schema_ensure_at = ?, last_schema_ensure_status = ?, last_schema_ensure_error = ?, updated_at = ?
		WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, at.UTC(), strings.TrimSpace(status), strings.TrimSpace(errorSummary), time.Now(), id)
	if err != nil {
		return fmt.Errorf("更新資料庫連接器 schema ensure 結果失敗: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("取得資料庫連接器 schema ensure 結果更新筆數失敗: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("資料庫連接器不存在")
	}
	return nil
}
