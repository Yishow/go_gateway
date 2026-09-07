package delivery

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open sqlite in-memory: %v", err)
	}

	ddl := `
	CREATE TABLE gw_delivery_journal (
		sequence INTEGER PRIMARY KEY AUTOINCREMENT,
		record_id TEXT NOT NULL,
		observed_at DATETIME NOT NULL,
		payload BLOB NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE gw_delivery_outbox (
		id TEXT PRIMARY KEY,
		destination_id TEXT NOT NULL,
		destination_revision TEXT NOT NULL DEFAULT '',
		plan_revision TEXT NOT NULL DEFAULT '',
		record_id TEXT NOT NULL,
		calculation_revision INTEGER NOT NULL DEFAULT 1,
		table_name TEXT NOT NULL,
		payload BLOB NOT NULL,
		status TEXT NOT NULL DEFAULT 'pending',
		retry_count INTEGER NOT NULL DEFAULT 0,
		next_retry_at DATETIME NOT NULL,
		last_error TEXT,
		observed_at DATETIME NOT NULL,
		delivered_at DATETIME,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE UNIQUE INDEX idx_gw_delivery_outbox_dest_rec 
		ON gw_delivery_outbox(destination_id, record_id);

	CREATE TABLE gw_delivery_receipts (
		destination_id TEXT NOT NULL,
		record_id TEXT NOT NULL,
		calculation_revision INTEGER NOT NULL DEFAULT 1,
		table_name TEXT NOT NULL,
		delivered_at DATETIME NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (destination_id, record_id, calculation_revision)
	);
	`
	if _, err := db.ExecContext(context.Background(), ddl); err != nil {
		t.Fatalf("failed to create tables: %v", err)
	}
	return db
}

func TestSQLStorage_JournalPersistentScanAndTruncate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	journal := NewSQLJournal(db)
	now := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)

	err := journal.Append(&JournalEntry{
		RecordID:   "rec-1",
		ObservedAt: now,
		Payload:    []byte(`{"test":1}`),
	})
	if err != nil {
		t.Fatalf("append failed: %v", err)
	}

	entries, err := journal.ScanFrom(1, 10)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	if len(entries) != 1 || entries[0].RecordID != "rec-1" {
		t.Errorf("unexpected journal scan entries: %+v", entries)
	}

	err = journal.TruncateBefore(entries[0].Sequence + 1)
	if err != nil {
		t.Fatalf("truncate failed: %v", err)
	}

	entriesAfter, _ := journal.ScanFrom(1, 10)
	if len(entriesAfter) != 0 {
		t.Errorf("expected 0 entries after truncate, got %d", len(entriesAfter))
	}
}

func TestSQLStorage_OutboxAndReceiptIntegration(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	outbox := NewSQLOutbox(db)
	receiptLedger := NewSQLReceiptLedger(db)
	now := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)

	item := &OutboxItem{
		ID:                  "out-1",
		DestinationID:       "dest-pg",
		RecordID:            "rec-meter-1",
		CalculationRevision: 2,
		Table:               "gw_record_samples",
		Payload:             []byte(`{"val":120.0}`),
		ObservedAt:          now,
	}

	// 1. 寫入 Outbox
	err := outbox.Enqueue(item)
	if err != nil {
		t.Fatalf("enqueue failed: %v", err)
	}

	// 2. 舊 Revision 覆寫應失敗
	itemStale := &OutboxItem{
		ID:                  "out-2",
		DestinationID:       "dest-pg",
		RecordID:            "rec-meter-1",
		CalculationRevision: 1, // 舊 revision
		Table:               "gw_record_samples",
		Payload:             []byte(`{"val":100.0}`),
		ObservedAt:          now,
	}
	err = outbox.Enqueue(itemStale)
	if err == nil {
		t.Errorf("expected error when enqueuing stale revision")
	}

	// 3. 取得 Pending
	pending, err := outbox.FetchPending("dest-pg", 10)
	if err != nil {
		t.Fatalf("fetch pending failed: %v", err)
	}
	if len(pending) != 1 || pending[0].CalculationRevision != 2 {
		t.Errorf("unexpected pending items: %+v", pending)
	}

	// 4. Receipt 記錄與去重檢查
	hasReceipt, err := receiptLedger.HasReceipt("dest-pg", "rec-meter-1", 2)
	if err != nil || hasReceipt {
		t.Errorf("receipt should not exist before delivery")
	}

	err = receiptLedger.SaveReceipt(&Receipt{
		DestinationID:       "dest-pg",
		RecordID:            "rec-meter-1",
		CalculationRevision: 2,
		Table:               "gw_record_samples",
		DeliveredAt:         now.Add(1 * time.Second),
	})
	if err != nil {
		t.Fatalf("save receipt failed: %v", err)
	}

	hasReceipt, err = receiptLedger.HasReceipt("dest-pg", "rec-meter-1", 2)
	if err != nil || !hasReceipt {
		t.Errorf("receipt should exist after saving")
	}

	// 5. 標記 Delivered
	err = outbox.MarkDelivered(item.ID, now.Add(1*time.Second))
	if err != nil {
		t.Fatalf("mark delivered failed: %v", err)
	}

	metrics, err := outbox.GetMetrics("dest-pg")
	if err != nil {
		t.Fatalf("get metrics failed: %v", err)
	}
	if metrics.PendingCount != 0 || metrics.DeliveredCount != 1 {
		t.Errorf("unexpected metrics after delivery: %+v", metrics)
	}
}
