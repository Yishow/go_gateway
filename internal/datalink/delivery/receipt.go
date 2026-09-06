package delivery

import (
	"fmt"
	"sync"
)

// ReceiptLedger 定義回執記錄簿的操作介面。
type ReceiptLedger interface {
	SaveReceipt(receipt *Receipt) error
	HasReceipt(destinationID string, recordID string, revision int64) (bool, error)
}

// MemoryReceiptLedger 記憶體測試用 ReceiptLedger 實作。
type MemoryReceiptLedger struct {
	mu       sync.RWMutex
	receipts map[string]*Receipt // key: destID + ":" + recordID + ":" + rev
}

// NewMemoryReceiptLedger 建立記憶體回執記錄簿。
func NewMemoryReceiptLedger() *MemoryReceiptLedger {
	return &MemoryReceiptLedger{
		receipts: make(map[string]*Receipt),
	}
}

func (m *MemoryReceiptLedger) SaveReceipt(receipt *Receipt) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if receipt == nil {
		return fmt.Errorf("receipt cannot be nil")
	}

	key := fmt.Sprintf("%s:%s:%d", receipt.DestinationID, receipt.RecordID, receipt.CalculationRevision)
	m.receipts[key] = receipt
	return nil
}

func (m *MemoryReceiptLedger) HasReceipt(destinationID string, recordID string, revision int64) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	key := fmt.Sprintf("%s:%s:%d", destinationID, recordID, revision)
	_, exists := m.receipts[key]
	return exists, nil
}
