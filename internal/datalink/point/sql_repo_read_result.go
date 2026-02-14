package point

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// BatchUpdateReadResult 批次更新點位讀取結果
func (r *SQLRepository) BatchUpdateReadResult(ctx context.Context, results []ReadResultUpdate) error {
	// 由於 SQLite 不支援複雜的批次更新語法，這裡使用交易 + 逐條更新
	// 對於高頻數據，這可能不是最高效的，但在 SQLite 場景下是可行的
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("開啟交易失敗: %w", err)
	}
	defer func() {
		if rbErr := tx.Rollback(); rbErr != nil && rbErr != sql.ErrTxDone {
			return
		}
	}()

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

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("取得更新筆數失敗: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("點位不存在: %s", pointID)
	}

	return nil
}
