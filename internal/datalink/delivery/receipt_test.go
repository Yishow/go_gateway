package delivery

import (
	"testing"
	"time"
)

func TestReceiptLedger_RecordAndVerify(t *testing.T) {
	ledger := NewMemoryReceiptLedger()

	now := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	receipt := &Receipt{
		DestinationID:       "dest-1",
		RecordID:            "rec-100",
		CalculationRevision: 1,
		Table:               "gw_record_samples",
		DeliveredAt:         now,
	}

	exists, err := ledger.HasReceipt("dest-1", "rec-100", 1)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if exists {
		t.Errorf("receipt should not exist initially")
	}

	err = ledger.SaveReceipt(receipt)
	if err != nil {
		t.Fatalf("save receipt failed: %v", err)
	}

	exists, err = ledger.HasReceipt("dest-1", "rec-100", 1)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if !exists {
		t.Errorf("receipt should exist after saving")
	}
}
