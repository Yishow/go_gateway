package handlers

import (
	"encoding/json"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/schema"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWorkspaceScaleMatchesDeclaredTarget(t *testing.T) {
	for _, tc := range []struct {
		name          string
		target        schema.DataType
		scale, offset float64
		want          any
		invalid       bool
	}{
		{"fraction rejected", schema.DataTypeInt16, 0.5, 10, nil, true},
		{"integer result", schema.DataTypeInt16, 2, 0, int16(486), false},
		{"compatible floating result", schema.DataTypeFloat64, 0.5, 10, float64(131.5), false},
		{"overflow rejected", schema.DataTypeInt16, 1000, 0, nil, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			steps := buildWorkspaceMappingPipeline(schema.DataTypeInt16, tc.target, tc.scale, tc.offset)
			raw, err := json.Marshal(steps)
			require.NoError(t, err)
			result, err := mapping.ExecutePipeline(int16(243), string(raw))
			if tc.invalid {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.want, result.CurrentValue)
			}
		})
	}
}

func TestStoredPipelineKeepsPublishedOrder(t *testing.T) {
	// The already-published cast-before-scale sequence must not be reordered.
	result, err := mapping.ExecutePipeline(float64(243), `[{"type":"cast","order":0,"params":{"target_type":"int16"}},{"type":"scale","order":1,"params":{"scale":0.5,"offset":10}}]`)
	require.NoError(t, err)
	require.Equal(t, float64(131.5), result.CurrentValue)
	_, err = mapping.ExecutePipeline(uint64(9007199254740993), `[{"type":"cast","order":0,"params":{"target_type":"float64"}},{"type":"scale","order":1,"params":{"scale":0.5,"offset":10}}]`)
	require.Error(t, err)
	_, err = mapping.ExecutePipeline(math.MaxFloat64, `[{"type":"scale","params":{"scale":2,"offset":0}}]`)
	require.Error(t, err)
	_, err = mapping.ExecutePipeline(float64(1), `[{"type":"scale","params":{"scale":"bad","offset":0}}]`)
	require.Error(t, err)
}
