package edge

import (
	"context"
	"math"

	"go-gateway/lib/algorithms"
)

// =============================================================================
// 內建處理器
// =============================================================================

// ScaleProcessor 縮放處理器
type ScaleProcessor struct {
	factor float64
	offset float64
}

func NewScaleProcessor(factor, offset float64) *ScaleProcessor {
	return &ScaleProcessor{factor: factor, offset: offset}
}

func (p *ScaleProcessor) Name() string { return "scale" }

func (p *ScaleProcessor) Process(ctx context.Context, input []float64) ([]float64, error) {
	result := make([]float64, len(input))
	for i, v := range input {
		result[i] = v*p.factor + p.offset
	}
	return result, nil
}

// FilterProcessor 濾波處理器 (移動平均)
type FilterProcessor struct {
	windowSize int
	buffer     []float64
}

func NewFilterProcessor(windowSize int) *FilterProcessor {
	return &FilterProcessor{
		windowSize: windowSize,
		buffer:     make([]float64, 0, windowSize),
	}
}

func (p *FilterProcessor) Name() string { return "filter" }

func (p *FilterProcessor) Process(ctx context.Context, input []float64) ([]float64, error) {
	result := make([]float64, len(input))
	
	for i, v := range input {
		// 加入緩衝區
		p.buffer = append(p.buffer, v)
		if len(p.buffer) > p.windowSize {
			p.buffer = p.buffer[1:]
		}
		
		// 計算移動平均
		sum := 0.0
		for _, bv := range p.buffer {
			sum += bv
		}
		result[i] = sum / float64(len(p.buffer))
	}
	
	return result, nil
}

// ThresholdProcessor 閾值處理器
type ThresholdProcessor struct {
	lowThreshold  float64
	highThreshold float64
	lowValue      float64
	highValue     float64
}

func NewThresholdProcessor(low, high, lowVal, highVal float64) *ThresholdProcessor {
	return &ThresholdProcessor{
		lowThreshold:  low,
		highThreshold: high,
		lowValue:      lowVal,
		highValue:     highVal,
	}
}

func (p *ThresholdProcessor) Name() string { return "threshold" }

func (p *ThresholdProcessor) Process(ctx context.Context, input []float64) ([]float64, error) {
	result := make([]float64, len(input))
	for i, v := range input {
		if v < p.lowThreshold {
			result[i] = p.lowValue
		} else if v > p.highThreshold {
			result[i] = p.highValue
		} else {
			result[i] = v
		}
	}
	return result, nil
}

// DeadbandProcessor 死區處理器
type DeadbandProcessor struct {
	threshold float64
	lastValue float64
	hasLast   bool
}

func NewDeadbandProcessor(threshold float64) *DeadbandProcessor {
	return &DeadbandProcessor{threshold: threshold}
}

func (p *DeadbandProcessor) Name() string { return "deadband" }

func (p *DeadbandProcessor) Process(ctx context.Context, input []float64) ([]float64, error) {
	result := make([]float64, len(input))
	for i, v := range input {
		if !p.hasLast || math.Abs(v-p.lastValue) >= p.threshold {
			result[i] = v
			p.lastValue = v
			p.hasLast = true
		} else {
			result[i] = p.lastValue
		}
	}
	return result, nil
}

// FFTProcessor FFT 處理器
type FFTProcessor struct {
	outputType string // "magnitude", "phase", "power"
}

func NewFFTProcessor(outputType string) *FFTProcessor {
	return &FFTProcessor{outputType: outputType}
}

func (p *FFTProcessor) Name() string { return "fft" }

func (p *FFTProcessor) Process(ctx context.Context, input []float64) ([]float64, error) {
	// 轉換為複數
	complexInput := make([]complex128, len(input))
	for i, v := range input {
		complexInput[i] = complex(v, 0)
	}

	// 執行 FFT
	spectrum := algorithms.FFT(complexInput)

	// 根據輸出類型返回結果
	switch p.outputType {
	case "phase":
		return algorithms.Phase(spectrum), nil
	case "power":
		return algorithms.PowerSpectrum(spectrum), nil
	default: // magnitude
		return algorithms.Magnitude(spectrum), nil
	}
}

// PIDProcessor PID 控制處理器
type PIDProcessor struct {
	pid *algorithms.PID
}

func NewPIDProcessor(kp, ki, kd, target float64) *PIDProcessor {
	pid := algorithms.NewPID(kp, ki, kd)
	pid.SetTarget(target)
	return &PIDProcessor{pid: pid}
}

func (p *PIDProcessor) Name() string { return "pid" }

func (p *PIDProcessor) Process(ctx context.Context, input []float64) ([]float64, error) {
	result := make([]float64, len(input))
	dt := 0.1 // 預設時間間隔
	for i, v := range input {
		result[i] = p.pid.Update(v, dt)
	}
	return result, nil
}

// =============================================================================
// 處理器註冊
// =============================================================================

// RegisterBuiltinProcessors 註冊內建處理器
func RegisterBuiltinProcessors(registry *ProcessorRegistry) {
	registry.Register("scale", func(params map[string]interface{}) (Processor, error) {
		factor := getFloat(params, "factor", 1.0)
		offset := getFloat(params, "offset", 0.0)
		return NewScaleProcessor(factor, offset), nil
	})

	registry.Register("filter", func(params map[string]interface{}) (Processor, error) {
		windowSize := getInt(params, "window_size", 5)
		return NewFilterProcessor(windowSize), nil
	})

	registry.Register("threshold", func(params map[string]interface{}) (Processor, error) {
		low := getFloat(params, "low", 0)
		high := getFloat(params, "high", 100)
		lowVal := getFloat(params, "low_value", 0)
		highVal := getFloat(params, "high_value", 100)
		return NewThresholdProcessor(low, high, lowVal, highVal), nil
	})

	registry.Register("deadband", func(params map[string]interface{}) (Processor, error) {
		threshold := getFloat(params, "threshold", 1.0)
		return NewDeadbandProcessor(threshold), nil
	})

	registry.Register("fft", func(params map[string]interface{}) (Processor, error) {
		outputType := getString(params, "output", "magnitude")
		return NewFFTProcessor(outputType), nil
	})

	registry.Register("pid", func(params map[string]interface{}) (Processor, error) {
		kp := getFloat(params, "kp", 1.0)
		ki := getFloat(params, "ki", 0.0)
		kd := getFloat(params, "kd", 0.0)
		target := getFloat(params, "target", 0.0)
		return NewPIDProcessor(kp, ki, kd, target), nil
	})
}

// 輔助函數
func getFloat(params map[string]interface{}, key string, defaultVal float64) float64 {
	if v, ok := params[key]; ok {
		switch val := v.(type) {
		case float64:
			return val
		case int:
			return float64(val)
		}
	}
	return defaultVal
}

func getInt(params map[string]interface{}, key string, defaultVal int) int {
	if v, ok := params[key]; ok {
		switch val := v.(type) {
		case int:
			return val
		case float64:
			return int(val)
		}
	}
	return defaultVal
}

func getString(params map[string]interface{}, key string, defaultVal string) string {
	if v, ok := params[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return defaultVal
}
