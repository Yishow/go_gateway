package storage

import (
	"context"
	"sync"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBatchWriterFlushOnBatchSize(t *testing.T) {
	writer := &mockTimeSeriesWriter{}
	bw := NewBatchWriter(writer, BatchWriterConfig{
		BatchSize:     2,
		FlushInterval: time.Hour,
		WriteTimeout:  time.Second,
	})
	defer func() {
		require.NoError(t, bw.Close())
	}()

	err := bw.Write(context.Background(), TimeSeriesRecord{TagID: "t1", Timestamp: time.Now()})
	require.NoError(t, err)
	assert.Equal(t, 0, writer.BatchCalls())

	err = bw.Write(context.Background(), TimeSeriesRecord{TagID: "t1", Timestamp: time.Now()})
	require.NoError(t, err)
	assert.Equal(t, 1, writer.BatchCalls())
	assert.Len(t, writer.LastBatch(), 2)
}

func TestMemoryStorageGetLatestBatch(t *testing.T) {
	store := NewMemoryStorage(10)
	now := time.Now()

	require.NoError(t, store.Write(context.Background(), TimeSeriesRecord{TagID: "a", Timestamp: now.Add(-2 * time.Minute)}))
	require.NoError(t, store.Write(context.Background(), TimeSeriesRecord{TagID: "b", Timestamp: now.Add(-1 * time.Minute)}))
	require.NoError(t, store.Write(context.Background(), TimeSeriesRecord{TagID: "a", Timestamp: now}))

	latest, err := store.GetLatestBatch(context.Background(), []string{"a", "b"})
	require.NoError(t, err)
	require.Len(t, latest, 2)
	assert.Equal(t, now, latest["a"].Timestamp)
	assert.Equal(t, now.Add(-1*time.Minute), latest["b"].Timestamp)
}

func TestValueToRecordConvertsBoolAndString(t *testing.T) {
	now := time.Now()

	boolRecord := ValueToRecord("t-bool", true, nil, now, schema.QualityGood, schema.DataTypeBool)
	require.NotNil(t, boolRecord.ValueBool)
	assert.True(t, *boolRecord.ValueBool)

	stringRecord := ValueToRecord("t-str", 123, nil, now, schema.QualityGood, schema.DataTypeString)
	require.NotNil(t, stringRecord.ValueText)
	assert.Equal(t, "123", *stringRecord.ValueText)
}

type mockTimeSeriesWriter struct {
	mu      sync.Mutex
	batches [][]TimeSeriesRecord
}

func (m *mockTimeSeriesWriter) Write(ctx context.Context, record TimeSeriesRecord) error {
	return nil
}

func (m *mockTimeSeriesWriter) WriteBatch(ctx context.Context, records []TimeSeriesRecord) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	copied := make([]TimeSeriesRecord, len(records))
	copy(copied, records)
	m.batches = append(m.batches, copied)
	return nil
}

func (m *mockTimeSeriesWriter) Flush(ctx context.Context) error {
	return nil
}

func (m *mockTimeSeriesWriter) Close() error {
	return nil
}

func (m *mockTimeSeriesWriter) BatchCalls() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.batches)
}

func (m *mockTimeSeriesWriter) LastBatch() []TimeSeriesRecord {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.batches) == 0 {
		return nil
	}
	return m.batches[len(m.batches)-1]
}
