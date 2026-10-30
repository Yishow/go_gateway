package recordingplan

import (
	"context"
	"math"
	"strings"
	"sync"
	"time"
)

// StorageSettlement is diagnostic evidence of one newly settled repair ledger
// row. It contains no execution token, owner, raw detail, SQL or destination.
// An unavailable settlement is the zero value, never permission to retry DDL.
type StorageSettlement struct {
	Available bool
	// OwnedExecution is copied from the native owned finish, never inferred from context.
	OwnedExecution           bool
	OperationID, WorkspaceID string
	Status                   SchemaOperationStatus
	VerifiedDigest           string
	UpdatedAt, CompletedAt   time.Time
}

// StorageObservationConfig is copied into an operation-local bounded capture.
// Observe and OnDrop must be nonblocking. Their panics cannot alter native work.
// RecordLimit comes from the existing diagnostic admission cap; zero admits none.
type StorageObservationConfig struct {
	RecordLimit              int
	OnDrop                   func(uint64)
	InfoEnabled, WarnEnabled bool
	Observe                  func(StorageSettlement)
}

type storageObservationKey struct{}
type storageObservationBatch struct {
	mu      sync.Mutex
	config  StorageObservationConfig
	records []StorageSettlement
	dropped uint64
	closed  bool
}

// CaptureStorageSettlements defers diagnostic callbacks until finish is called.
// The caller must register finish before acquiring outer native locks so those
// locks are released first. Nested captures share the outer batch. Finish is
// idempotent and invokes callbacks without holding the capture mutex.
func CaptureStorageSettlements(ctx context.Context, cfg StorageObservationConfig) (context.Context, func()) {
	if ctx == nil || cfg.Observe == nil || !cfg.InfoEnabled && !cfg.WarnEnabled {
		return ctx, func() {}
	}
	if _, ok := ctx.Value(storageObservationKey{}).(*storageObservationBatch); ok {
		return ctx, func() {}
	}
	batch := &storageObservationBatch{config: cfg}
	return context.WithValue(ctx, storageObservationKey{}, batch), batch.finish
}

func (b *storageObservationBatch) finish() {
	b.mu.Lock()
	if b.closed {
		b.mu.Unlock()
		return
	}
	b.closed = true
	records, dropped := b.records, b.dropped
	b.records = nil
	b.mu.Unlock()
	for _, record := range records {
		func() {
			defer func() {
				if recover() != nil && dropped < math.MaxUint64 {
					dropped++
				}
			}()
			b.config.Observe(record)
		}()
	}
	if dropped != 0 && b.config.OnDrop != nil {
		func() { defer func() { _ = recover() }(); b.config.OnDrop(dropped) }()
	}
}

// recordStorageSettlement reads no database and retains only closed scalar
// metadata. It is deliberately not an observer invocation inside native work.
func recordStorageSettlement(ctx context.Context, operationID string, owned bool, expected SchemaOperationResult, settledAt time.Time, op *SchemaOperation, err error) {
	defer func() { _ = recover() }()
	if ctx == nil {
		return
	}
	b, ok := ctx.Value(storageObservationKey{}).(*storageObservationBatch)
	if !ok {
		return
	}
	record := storageSettlement(operationID, owned, expected, settledAt, op, err)
	info := record.Available && record.OwnedExecution && record.Status == SchemaOperationSucceeded
	if info && !b.config.InfoEnabled || !info && !b.config.WarnEnabled {
		return // Filtering is not admission loss and consumes no descriptor slot.
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return // The operation-local lifetime has ended; never call out here.
	}
	if b.config.RecordLimit <= 0 || len(b.records) >= b.config.RecordLimit {
		if b.dropped < math.MaxUint64 {
			b.dropped++
		}
		return
	}
	record.OperationID = strings.Clone(record.OperationID)
	record.WorkspaceID = strings.Clone(record.WorkspaceID)
	record.VerifiedDigest = strings.Clone(record.VerifiedDigest)
	b.records = append(b.records, record)
}

func storageSettlement(operationID string, owned bool, expected SchemaOperationResult, settledAt time.Time, op *SchemaOperation, err error) StorageSettlement {
	if err != nil || op == nil || op.OperationID != operationID || op.Action != StorageRepairAction ||
		op.Status != expected.Status || !isTerminalSchemaOperationStatus(op.Status) ||
		op.VerifiedDigest != expected.VerifiedDigest || op.CompletedAt == nil || settledAt.IsZero() ||
		!op.UpdatedAt.Equal(settledAt) || !op.CompletedAt.Equal(settledAt) ||
		!storageObservationID(op.OperationID) || !storageObservationID(op.WorkspaceID) ||
		(op.VerifiedDigest != "" && !storageObservationID(op.VerifiedDigest)) {
		return StorageSettlement{}
	}
	return StorageSettlement{Available: true, OwnedExecution: owned, OperationID: op.OperationID, WorkspaceID: op.WorkspaceID,
		Status: op.Status, VerifiedDigest: op.VerifiedDigest, UpdatedAt: op.UpdatedAt, CompletedAt: *op.CompletedAt}
}

func storageObservationID(value string) bool {
	if len(value) == 0 || len(value) > 128 {
		return false
	}
	for i, c := range []byte(value) {
		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '_' || c == '-' || i > 0 && c == '.' {
			continue
		}
		return false
	}
	return true
}
