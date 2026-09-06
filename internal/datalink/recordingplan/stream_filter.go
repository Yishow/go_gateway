package recordingplan

import (
	"fmt"
	"math"
	"sync"
	"time"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
)

type streamState struct {
	lastValue        interface{}
	lastQuality      schema.QualityFlag
	lastObservedAt   time.Time
	lastEmittedAt    time.Time
	lastTriggerVal   interface{}
	hasEmitted       bool
}

// StreamFilter 實作即時樣本流的明細過濾、死區、心跳與批次觸發去重。
type StreamFilter struct {
	mu     sync.Mutex
	states map[string]*streamState
}

// NewStreamFilter 建立新的 StreamFilter。
func NewStreamFilter() *StreamFilter {
	return &StreamFilter{
		states: make(map[string]*streamState),
	}
}

// ShouldEmitSample 依 Stream 模式與策略評估是否應將此樣本輸出至明細流。
func (f *StreamFilter) ShouldEmitSample(stream PlanStream, sample measurement.SampleEnvelope) (bool, string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	state, exists := f.states[stream.StreamID]
	if !exists {
		state = &streamState{}
		f.states[stream.StreamID] = state
	}

	switch stream.Mode {
	case StreamModeLatestOnly:
		if state.hasEmitted && !sample.ObservedAt.After(state.lastObservedAt) {
			return false, "sample_not_newer_than_latest"
		}
		state.lastValue = sample.Value
		state.lastQuality = sample.Quality
		state.lastObservedAt = sample.ObservedAt
		state.lastEmittedAt = sample.ObservedAt
		state.hasEmitted = true
		return true, "latest_sample"

	case StreamModeRawHistory:
		switch stream.RawPolicy {
		case RawPolicyEverySample, "":
			state.lastValue = sample.Value
			state.lastQuality = sample.Quality
			state.lastObservedAt = sample.ObservedAt
			state.lastEmittedAt = sample.ObservedAt
			state.hasEmitted = true
			return true, "every_sample"

		case RawPolicyOnChange:
			if !state.hasEmitted {
				state.lastValue = sample.Value
				state.lastQuality = sample.Quality
				state.lastObservedAt = sample.ObservedAt
				state.lastEmittedAt = sample.ObservedAt
				state.hasEmitted = true
				return true, "initial_sample"
			}

			// 品質變更絕對不被死區隱藏
			if sample.Quality != state.lastQuality {
				state.lastValue = sample.Value
				state.lastQuality = sample.Quality
				state.lastObservedAt = sample.ObservedAt
				state.lastEmittedAt = sample.ObservedAt
				return true, "quality_change"
			}

			// 檢查最大心跳 (Heartbeat)
			heartbeatSec := 60
			if stream.MaxHeartbeatSeconds != nil && *stream.MaxHeartbeatSeconds > 0 {
				heartbeatSec = *stream.MaxHeartbeatSeconds
			}
			if sample.ObservedAt.Sub(state.lastEmittedAt) >= time.Duration(heartbeatSec)*time.Second {
				state.lastValue = sample.Value
				state.lastQuality = sample.Quality
				state.lastObservedAt = sample.ObservedAt
				state.lastEmittedAt = sample.ObservedAt
				return true, "heartbeat"
			}

			// 檢查數值死區
			deadband := 0.0
			if stream.OnChangeDeadband != nil && *stream.OnChangeDeadband >= 0 {
				deadband = *stream.OnChangeDeadband
			}

			currNum, currIsNum := toNumeric(sample.Value)
			lastNum, lastIsNum := toNumeric(state.lastValue)

			if currIsNum && lastIsNum {
				if math.Abs(currNum-lastNum) > deadband {
					state.lastValue = sample.Value
					state.lastQuality = sample.Quality
					state.lastObservedAt = sample.ObservedAt
					state.lastEmittedAt = sample.ObservedAt
					return true, "deadband_exceeded"
				}
				return false, "within_deadband"
			}

			// 非數值相等性比較
			if fmt.Sprintf("%v", sample.Value) != fmt.Sprintf("%v", state.lastValue) {
				state.lastValue = sample.Value
				state.lastQuality = sample.Quality
				state.lastObservedAt = sample.ObservedAt
				state.lastEmittedAt = sample.ObservedAt
				return true, "value_change"
			}

			return false, "value_unchanged"
		}
	}

	return true, "default_emit"
}

// EvaluateBatchTrigger 評估批次快照的觸發邊緣（0->1 上升緣觸發），避免 held-high 重複產生。
func (f *StreamFilter) EvaluateBatchTrigger(stream PlanStream, triggerSample measurement.SampleEnvelope) bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	state, exists := f.states[stream.StreamID]
	if !exists {
		state = &streamState{}
		f.states[stream.StreamID] = state
	}

	currActive := isTruthy(triggerSample.Value)
	prevActive := isTruthy(state.lastTriggerVal)
	state.lastTriggerVal = triggerSample.Value

	// 僅在上升緣 (0 -> 1) 觸發
	if currActive && !prevActive {
		return true
	}
	return false
}

func toNumeric(val interface{}) (float64, bool) {
	switch v := val.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int16:
		return float64(v), true
	case uint16:
		return float64(v), true
	case int32:
		return float64(v), true
	case uint32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint64:
		return float64(v), true
	default:
		return 0, false
	}
}

func isTruthy(val interface{}) bool {
	if val == nil {
		return false
	}
	switch v := val.(type) {
	case bool:
		return v
	case int, int16, int32, int64, uint16, uint32, uint64:
		n, _ := toNumeric(v)
		return n != 0
	case float32, float64:
		n, _ := toNumeric(v)
		return n != 0
	case string:
		return v == "1" || v == "true" || v == "TRUE"
	default:
		return false
	}
}
