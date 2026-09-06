package delivery

import (
	"testing"
	"time"
)

func TestJournal_AppendAndScan(t *testing.T) {
	journal := NewMemoryJournal()

	now := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	item1 := &JournalEntry{
		Sequence:   1,
		RecordID:   "rec-1",
		ObservedAt: now,
		Payload:    []byte(`{"val": 10}`),
	}
	item2 := &JournalEntry{
		Sequence:   2,
		RecordID:   "rec-2",
		ObservedAt: now.Add(1 * time.Second),
		Payload:    []byte(`{"val": 20}`),
	}

	err := journal.Append(item1)
	if err != nil {
		t.Fatalf("failed to append item1: %v", err)
	}
	err = journal.Append(item2)
	if err != nil {
		t.Fatalf("failed to append item2: %v", err)
	}

	entries, err := journal.ScanFrom(1, 10)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].RecordID != "rec-1" || entries[1].RecordID != "rec-2" {
		t.Errorf("unexpected scan content")
	}
}

func TestJournal_TruncateBefore(t *testing.T) {
	journal := NewMemoryJournal()
	now := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)

	_ = journal.Append(&JournalEntry{Sequence: 1, RecordID: "rec-1", ObservedAt: now})
	_ = journal.Append(&JournalEntry{Sequence: 2, RecordID: "rec-2", ObservedAt: now.Add(1 * time.Second)})
	_ = journal.Append(&JournalEntry{Sequence: 3, RecordID: "rec-3", ObservedAt: now.Add(2 * time.Second)})

	err := journal.TruncateBefore(3)
	if err != nil {
		t.Fatalf("truncate failed: %v", err)
	}

	entries, err := journal.ScanFrom(1, 10)
	if err != nil {
		t.Fatalf("scan after truncate failed: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry after truncate, got %d", len(entries))
	}
	if entries[0].Sequence != 3 {
		t.Errorf("expected sequence 3, got %d", entries[0].Sequence)
	}
}
