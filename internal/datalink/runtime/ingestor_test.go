package runtime

import (
	"context"
	"testing"
	"time"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/storage"

	"github.com/stretchr/testify/require"
)

type mockWriter struct {
	records []storage.TimeSeriesRecord
}

func (m *mockWriter) Write(ctx context.Context, record storage.TimeSeriesRecord) error {
	m.records = append(m.records, record)
	return nil
}

func (m *mockWriter) WriteBatch(ctx context.Context, records []storage.TimeSeriesRecord) error {
	m.records = append(m.records, records...)
	return nil
}

func (m *mockWriter) Flush(ctx context.Context) error { return nil }
func (m *mockWriter) Close() error                    { return nil }

type mockTargetWriter struct {
	calls []targetWriteCall
}

type targetWriteCall struct {
	tagID      string
	value      any
	observedAt time.Time
}

func (m *mockTargetWriter) WriteTagValue(ctx context.Context, tagID string, value any, observedAt time.Time) error {
	m.calls = append(m.calls, targetWriteCall{
		tagID:      tagID,
		value:      value,
		observedAt: observedAt,
	})
	return nil
}

func TestHandleCollectedValue_RunPipelineAndWrite(t *testing.T) {
	mw := &mockWriter{}
	s := &Service{
		config: Config{UpdatePointState: false},
		writer: mw,
		mappingIndex: map[string][]mappingBinding{
			"p1": {
				{
					TagID:             "t1",
					TagDataType:       schema.DataTypeUint16,
					TransformPipeline: `[{"type":"scale","order":1,"params":{"scale":2}}]`,
				},
			},
		},
	}

	s.handleCollectedValue(context.Background(), collector.CollectedValue{
		PointID:   "p1",
		Value:     uint16(21),
		Timestamp: time.Now(),
		Quality:   schema.QualityGood,
		RawBytes:  []byte{0x00, 0x15},
	})

	require.Len(t, mw.records, 1)
	require.Equal(t, "t1", mw.records[0].TagID)
	require.NotNil(t, mw.records[0].ValueNum)
	require.Equal(t, 42.0, *mw.records[0].ValueNum)
}

func TestHandleCollectedValue_NoMapping_NoWrite(t *testing.T) {
	mw := &mockWriter{}
	s := &Service{
		config:       Config{UpdatePointState: false},
		writer:       mw,
		mappingIndex: map[string][]mappingBinding{},
	}

	s.handleCollectedValue(context.Background(), collector.CollectedValue{
		PointID:   "missing",
		Value:     uint16(1),
		Timestamp: time.Now(),
		Quality:   schema.QualityGood,
	})

	require.Len(t, mw.records, 0)
}

func TestHandleCollectedValue_WritesToTargetWriter(t *testing.T) {
	mw := &mockWriter{}
	targetWriter := &mockTargetWriter{}
	ts := time.Date(2026, 3, 16, 12, 0, 0, 0, time.UTC)
	s := &Service{
		config: Config{UpdatePointState: false},
		writer: mw,
		target: targetWriter,
		mappingIndex: map[string][]mappingBinding{
			"p1": {
				{
					TagID:             "t1",
					TagDataType:       schema.DataTypeFloat64,
					TransformPipeline: `[{"type":"scale","order":1,"params":{"scale":0.5}}]`,
				},
			},
		},
	}

	s.handleCollectedValue(context.Background(), collector.CollectedValue{
		PointID:   "p1",
		Value:     20.0,
		Timestamp: ts,
		Quality:   schema.QualityGood,
	})

	require.Len(t, targetWriter.calls, 1)
	require.Equal(t, "t1", targetWriter.calls[0].tagID)
	require.Equal(t, 10.0, targetWriter.calls[0].value)
	require.Equal(t, ts, targetWriter.calls[0].observedAt)
}
