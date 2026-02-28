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
