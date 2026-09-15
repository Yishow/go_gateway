package delivery

import (
	"testing"
	"time"
)

func TestSQLJournalRejectsCorruptObservedTime(t *testing.T) {
	db := setupTestDB(t)
	t.Cleanup(func() { _ = db.Close() })
	journal := NewSQLJournal(db)
	if err := journal.Append(&JournalEntry{RecordID: "corrupt-time", ObservedAt: time.Now(), Payload: []byte("sample")}); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(t.Context(), "UPDATE gw_delivery_journal SET observed_at = ?", "invalid-time"); err != nil {
		t.Fatal(err)
	}
	if got, err := journal.ScanFrom(0, 10); err == nil || got != nil {
		t.Fatalf("corrupt journal time must fail: got=%v err=%v", got, err)
	}
}

func TestSQLOutboxRejectsCorruptTimes(t *testing.T) {
	for _, query := range []string{
		"UPDATE gw_delivery_outbox SET observed_at = ?",
		"UPDATE gw_delivery_outbox SET next_retry_at = ?",
	} {
		t.Run(query, func(t *testing.T) {
			db := setupTestDB(t)
			t.Cleanup(func() { _ = db.Close() })
			outbox := NewSQLOutbox(db)
			item := &OutboxItem{ID: "time-item", DestinationID: "time-dest", RecordID: "time-record", Payload: []byte("sample"), ObservedAt: time.Now()}
			if err := outbox.Enqueue(item); err != nil {
				t.Fatal(err)
			}
			if _, err := db.ExecContext(t.Context(), query, "invalid-time"); err != nil {
				t.Fatal(err)
			}
			if got, err := outbox.FetchPending(item.DestinationID, 10); err == nil || got != nil {
				t.Fatalf("corrupt outbox time must fail: got=%v err=%v", got, err)
			}
		})
	}
}
