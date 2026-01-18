package mapping

import (
	"encoding/json"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/assert"
)

func TestExecutePipeline_Scale(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		steps    []schema.TransformStep
		expected interface{}
		wantErr  bool
	}{
		{
			name:  "Scale Simple",
			input: float64(10),
			steps: []schema.TransformStep{
				{
					Type: schema.TransformScale,
					Params: map[string]interface{}{
						"multiplier": 2.0,
						"offset":     5.0,
					},
				},
			},
			expected: float64(25),
			wantErr:  false,
		},
		{
			name:  "Scale with Divisor",
			input: float64(100),
			steps: []schema.TransformStep{
				{
					Type: schema.TransformScale,
					Params: map[string]interface{}{
						"divisor": 10.0,
					},
				},
			},
			expected: float64(10),
			wantErr:  false,
		},
		{
			name:  "Scale Min Cap",
			input: float64(5),
			steps: []schema.TransformStep{
				{
					Type: schema.TransformScale,
					Params: map[string]interface{}{
						"min": 10.0,
					},
				},
			},
			expected: float64(10), // Should be capped at min
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pipelineJSON, _ := json.Marshal(tt.steps)
			ctx, err := ExecutePipeline(tt.input, string(pipelineJSON))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, ctx.CurrentValue)
			}
		})
	}
}

func TestExecutePipeline_Conditional(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		steps    []schema.TransformStep
		expected interface{}
	}{
		{
			name:  "Conditional GT",
			input: float64(50),
			steps: []schema.TransformStep{
				{
					Type: schema.TransformConditional,
					Params: map[string]interface{}{
						"operator":    "gt",
						"threshold":   40.0,
						"true_value":  "High",
						"false_value": "Low",
					},
				},
			},
			expected: "High",
		},
		{
			name:  "Conditional LTE",
			input: float64(50),
			steps: []schema.TransformStep{
				{
					Type: schema.TransformConditional,
					Params: map[string]interface{}{
						"operator":    "lte",
						"threshold":   40.0,
						"true_value":  "High",
						"false_value": "Low",
					},
				},
			},
			expected: "Low",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pipelineJSON, _ := json.Marshal(tt.steps)
			ctx, err := ExecutePipeline(tt.input, string(pipelineJSON))
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, ctx.CurrentValue)
		})
	}
}

func TestExecutePipeline_Formula(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		steps    []schema.TransformStep
		expected interface{}
	}{
		{
			name:  "Formula Abs",
			input: float64(-10),
			steps: []schema.TransformStep{
				{
					Type: schema.TransformFormula,
					Params: map[string]interface{}{
						"expression": "abs(x)",
					},
				},
			},
			expected: float64(10),
		},
		{
			name:  "Formula Sqrt",
			input: float64(16),
			steps: []schema.TransformStep{
				{
					Type: schema.TransformFormula,
					Params: map[string]interface{}{
						"expression": "sqrt(x)",
					},
				},
			},
			expected: float64(4),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pipelineJSON, _ := json.Marshal(tt.steps)
			ctx, err := ExecutePipeline(tt.input, string(pipelineJSON))
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, ctx.CurrentValue)
		})
	}
}

func TestExecutePipeline_Chained(t *testing.T) {
	// Scale -> Cast -> Conditional
	input := float64(10)
	steps := []schema.TransformStep{
		{
			Type:  schema.TransformScale,
			Order: 1,
			Params: map[string]interface{}{
				"multiplier": 10.0, // = 100
			},
		},
		{
			Type:  schema.TransformConditional,
			Order: 2,
			Params: map[string]interface{}{
				"operator":    "gte",
				"threshold":   100.0,
				"true_value":  "Max",
				"false_value": "Normal",
			},
		},
	}

	pipelineJSON, _ := json.Marshal(steps)
	ctx, err := ExecutePipeline(input, string(pipelineJSON))
	assert.NoError(t, err)
	assert.Equal(t, "Max", ctx.CurrentValue)
	assert.Len(t, ctx.StepResults, 2)
	assert.Equal(t, "scale", ctx.StepResults[0].StepType)
	assert.Equal(t, "conditional", ctx.StepResults[1].StepType)
}
