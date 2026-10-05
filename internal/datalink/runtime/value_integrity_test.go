package runtime

import (
	"encoding/json"
	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/schema"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRuntimeValueWirePreservesUnsafeIntegers(t *testing.T) {
	for _, tc := range []struct {
		input any
		want  any
	}{{uint64(9007199254740993), "9007199254740993"}, {int64(math.MinInt64), "-9223372036854775808"}, {uint64(math.MaxUint64), "18446744073709551615"}, {int64(215), float64(215)}, {float64(131.5), float64(131.5)}} {
		wire, err := json.Marshal(ValueEvent{RawValue: tc.input, TransformedValue: tc.input})
		require.NoError(t, err)
		var received map[string]any
		require.NoError(t, json.Unmarshal(wire, &received))
		require.Equal(t, tc.want, received["raw_value"])
		require.Equal(t, tc.want, received["transformed_value"])
	}
}

func TestRuntimeInvalidCastDoesNotPublishGoodValueOrWrite(t *testing.T) {
	writer, target := &mockWriter{}, &mockTargetWriter{}
	svc := &Service{config: Config{UpdatePointState: false}, writer: writer, target: target, mappingIndex: map[string][]mappingBinding{"point-1": {{TagID: "tag-1", TagDataType: schema.DataTypeUint64, TransformPipeline: `[{"type":"cast","params":{"target_type":"uint64"}}]`}}}}
	events, unsubscribe := svc.SubscribeValueEvents("device-1", []string{"point-1"})
	defer unsubscribe()
	svc.handleCollectedValue(t.Context(), collector.CollectedValue{PointID: "point-1", DeviceID: "device-1", Value: "bad", Timestamp: time.Now(), Quality: schema.QualityGood})
	require.Empty(t, writer.records)
	require.Empty(t, target.calls)
	event := <-events
	require.Equal(t, schema.QualityBad, event.Quality)
	require.Nil(t, event.TransformedValue)
	require.Equal(t, "bad", event.RawValue)
}
