package delivery

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// SQLOutbox 實作基於關聯資料庫的 Outbox 佇列。
type SQLOutbox struct {
	db *sql.DB
}

// NewSQLOutbox 建立 SQLOutbox 實例。
func NewSQLOutbox(db *sql.DB) *SQLOutbox {
	return &SQLOutbox{db: db}
}

func (o *SQLOutbox) Enqueue(item *OutboxItem) error {
	if item == nil {
		return fmt.Errorf("outbox item cannot be nil")
	}

	ctx := context.Background()

	// 檢查是否已有相同 destination_id 與 record_id
	var existingRev int64
	queryCheck := `SELECT calculation_revision FROM gw_delivery_outbox WHERE destination_id = ? AND record_id = ?`
	err := o.db.QueryRowContext(ctx, queryCheck, item.DestinationID, item.RecordID).Scan(&existingRev)
	if err == nil {
		if item.CalculationRevision < existingRev {
			return fmt.Errorf("stale calculation revision %d, current is %d", item.CalculationRevision, existingRev)
		}
		// 更新現有記錄
		queryUpdate := `
			UPDATE gw_delivery_outbox
			SET calculation_revision = ?, payload = ?, status = 'pending', updated_at = ?
			WHERE destination_id = ? AND record_id = ?
		`
		_, err = o.db.ExecContext(ctx, queryUpdate,
			item.CalculationRevision,
			item.Payload,
			time.Now().UTC().Format(time.RFC3339Nano),
			item.DestinationID,
			item.RecordID,
		)
		return err
	}

	now := time.Now().UTC()
	if item.Status == "" {
		item.Status = StatusPending
	}
	if item.NextRetryAt.IsZero() {
		item.NextRetryAt = now
	}

	queryInsert := `
		INSERT INTO gw_delivery_outbox (
			id, destination_id, destination_revision, plan_revision, record_id,
			calculation_revision, table_name, payload, status, retry_count,
			next_retry_at, observed_at, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = o.db.ExecContext(ctx, queryInsert,
		item.ID,
		item.DestinationID,
		item.DestinationRevision,
		item.PlanRevision,
		item.RecordID,
		item.CalculationRevision,
		item.Table,
		item.Payload,
		string(item.Status),
		item.RetryCount,
		item.NextRetryAt.Format(time.RFC3339Nano),
		item.ObservedAt.Format(time.RFC3339Nano),
		now.Format(time.RFC3339Nano),
		now.Format(time.RFC3339Nano),
	)
	return err
}

func (o *SQLOutbox) FetchPending(destinationID string, limit int) ([]*OutboxItem, error) {
	if limit <= 0 {
		limit = 100
	}

	ctx := context.Background()
	nowStr := time.Now().UTC().Format(time.RFC3339Nano)

	query := `
		SELECT id, destination_id, destination_revision, plan_revision, record_id,
		       calculation_revision, table_name, payload, status, retry_count,
		       next_retry_at, observed_at
		FROM gw_delivery_outbox
		WHERE destination_id = ?
		  AND (status = 'pending' OR (status = 'retrying' AND next_retry_at <= ?))
		ORDER BY observed_at ASC
		LIMIT ?
	`
	rows, err := o.db.QueryContext(ctx, query, destinationID, nowStr, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pending outbox: %w", err)
	}
	defer rows.Close()

	var items []*OutboxItem
	for rows.Next() {
		var (
			it         OutboxItem
			statusStr  string
			retryAtStr string
			obsStr     string
		)
		err := rows.Scan(
			&it.ID, &it.DestinationID, &it.DestinationRevision, &it.PlanRevision, &it.RecordID,
			&it.CalculationRevision, &it.Table, &it.Payload, &statusStr, &it.RetryCount,
			&retryAtStr, &obsStr,
		)
		if err != nil {
			return nil, err
		}
		it.Status = OutboxStatus(statusStr)
		it.NextRetryAt, _ = time.Parse(time.RFC3339Nano, retryAtStr)
		it.ObservedAt, _ = time.Parse(time.RFC3339Nano, obsStr)
		items = append(items, &it)
	}
	return items, nil
}

func (o *SQLOutbox) MarkDelivered(itemID string, deliveredAt time.Time) error {
	ctx := context.Background()
	query := `
		UPDATE gw_delivery_outbox
		SET status = 'delivered', delivered_at = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := o.db.ExecContext(ctx, query,
		deliveredAt.Format(time.RFC3339Nano),
		time.Now().UTC().Format(time.RFC3339Nano),
		itemID,
	)
	return err
}

func (o *SQLOutbox) MarkFailed(itemID string, errStr string, maxRetries int) error {
	ctx := context.Background()
	var retryCount int
	queryGet := `SELECT retry_count FROM gw_delivery_outbox WHERE id = ?`
	if err := o.db.QueryRowContext(ctx, queryGet, itemID).Scan(&retryCount); err != nil {
		return err
	}

	retryCount++
	newStatus, nextRetry := CalculateBackoff(retryCount, maxRetries)

	queryUpdate := `
		UPDATE gw_delivery_outbox
		SET status = ?, retry_count = ?, next_retry_at = ?, last_error = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := o.db.ExecContext(ctx, queryUpdate,
		string(newStatus),
		retryCount,
		nextRetry.Format(time.RFC3339Nano),
		errStr,
		time.Now().UTC().Format(time.RFC3339Nano),
		itemID,
	)
	return err
}

func (o *SQLOutbox) GetMetrics(destinationID string) (DestinationMetrics, error) {
	ctx := context.Background()
	var metrics DestinationMetrics
	metrics.DestinationID = destinationID
	metrics.Status = StatusPending

	query := `
		SELECT
			COALESCE(SUM(CASE WHEN status IN ('pending', 'retrying') THEN 1 ELSE 0 END), 0) AS pending_count,
			COALESCE(SUM(CASE WHEN status IN ('pending', 'retrying') THEN LENGTH(payload) ELSE 0 END), 0) AS pending_bytes,
			COALESCE(SUM(CASE WHEN status = 'delivered' THEN 1 ELSE 0 END), 0) AS delivered_count,
			COALESCE(SUM(CASE WHEN status IN ('blocked', 'quarantined') THEN 1 ELSE 0 END), 0) AS failed_count,
			MIN(CASE WHEN status IN ('pending', 'retrying') THEN observed_at END) AS oldest_observed_at
		FROM gw_delivery_outbox
		WHERE destination_id = ?
	`
	var (
		oldestObs sql.NullString
	)
	err := o.db.QueryRowContext(ctx, query, destinationID).Scan(
		&metrics.PendingCount,
		&metrics.PendingBytes,
		&metrics.DeliveredCount,
		&metrics.FailedCount,
		&oldestObs,
	)
	if err != nil {
		return metrics, err
	}

	if oldestObs.Valid && oldestObs.String != "" {
		t, err := time.Parse(time.RFC3339Nano, oldestObs.String)
		if err == nil {
			metrics.OldestPendingAgeSec = time.Now().UTC().Sub(t).Seconds()
		}
	}

	return metrics, nil
}

// SQLReceiptLedger 實作基於關聯資料庫的回執記錄簿。
type SQLReceiptLedger struct {
	db *sql.DB
}

// NewSQLReceiptLedger 建立 SQLReceiptLedger 實例。
func NewSQLReceiptLedger(db *sql.DB) *SQLReceiptLedger {
	return &SQLReceiptLedger{db: db}
}

func (r *SQLReceiptLedger) SaveReceipt(receipt *Receipt) error {
	if receipt == nil {
		return fmt.Errorf("receipt cannot be nil")
	}

	query := `
		INSERT INTO gw_delivery_receipts (destination_id, record_id, calculation_revision, table_name, delivered_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(context.Background(), query,
		receipt.DestinationID,
		receipt.RecordID,
		receipt.CalculationRevision,
		receipt.Table,
		receipt.DeliveredAt.Format(time.RFC3339Nano),
		time.Now().UTC().Format(time.RFC3339Nano),
	)
	return err
}

func (r *SQLReceiptLedger) HasReceipt(destinationID string, recordID string, revision int64) (bool, error) {
	query := `SELECT 1 FROM gw_delivery_receipts WHERE destination_id = ? AND record_id = ? AND calculation_revision = ?`
	var dummy int
	err := r.db.QueryRowContext(context.Background(), query, destinationID, recordID, revision).Scan(&dummy)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
