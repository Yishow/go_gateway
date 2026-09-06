package aggregation

import (
	"fmt"
	"sort"
	"time"
)

// StateTrackingConfig 狀態追蹤配置。
type StateTrackingConfig struct {
	MaxHoldSeconds      float64 `json:"max_hold_seconds"`
	CalculationRevision int64   `json:"calculation_revision"`
}

// TrackStateDurations 統計給定視窗內的狀態轉換次數、各狀態持續時間與未知 (unknown) 時長。
func TrackStateDurations(
	measurementID string,
	windowStart time.Time,
	windowEnd time.Time,
	samples []TelemetrySample,
	config StateTrackingConfig,
) StateTrackingResult {
	totalSeconds := windowEnd.Sub(windowStart).Seconds()
	if totalSeconds <= 0 {
		totalSeconds = 1.0
	}

	maxHold := config.MaxHoldSeconds
	if maxHold <= 0 {
		maxHold = 60.0
	}

	result := StateTrackingResult{
		MeasurementID:       measurementID,
		WindowStart:         windowStart,
		WindowEnd:           windowEnd,
		CalculationRevision: config.CalculationRevision,
	}

	// 提取字串狀態
	type stateSample struct {
		state string
		at    time.Time
		good  bool
	}

	var stateSamples []stateSample
	for _, s := range samples {
		if s.ObservedAt.Before(windowStart) || !s.ObservedAt.Before(windowEnd) {
			continue
		}

		st := ""
		if s.ValueString != nil {
			st = *s.ValueString
		} else if s.ValueBool != nil {
			if *s.ValueBool {
				st = "ON"
			} else {
				st = "OFF"
			}
		} else if s.Value != nil {
			st = fmt.Sprintf("%v", s.Value)
		}

		stateSamples = append(stateSamples, stateSample{
			state: st,
			at:    s.ObservedAt,
			good:  s.Quality == QualityGood,
		})
	}

	if len(stateSamples) == 0 {
		result.UnknownDurationSec = totalSeconds
		return result
	}

	sort.SliceStable(stateSamples, func(i, j int) bool {
		return stateSamples[i].at.Before(stateSamples[j].at)
	})

	result.InitialState = stateSamples[0].state
	result.FinalState = stateSamples[len(stateSamples)-1].state

	durationsMap := make(map[string]float64)
	var transitions int
	var validTrackedDuration float64

	for i := 0; i < len(stateSamples); i++ {
		curr := stateSamples[i]
		if i > 0 && curr.state != stateSamples[i-1].state {
			transitions++
		}

		if !curr.good || curr.state == "" {
			continue
		}

		segStart := curr.at
		if segStart.Before(windowStart) {
			segStart = windowStart
		}

		holdEnd := curr.at.Add(time.Duration(maxHold * float64(time.Second)))
		segEnd := holdEnd

		if i+1 < len(stateSamples) {
			nextTime := stateSamples[i+1].at
			if nextTime.Before(segEnd) {
				segEnd = nextTime
			}
		}

		if segEnd.After(windowEnd) {
			segEnd = windowEnd
		}

		if segEnd.After(segStart) {
			dur := segEnd.Sub(segStart).Seconds()
			durationsMap[curr.state] += dur
			validTrackedDuration += dur
		}
	}

	result.TransitionsCount = transitions
	result.UnknownDurationSec = totalSeconds - validTrackedDuration
	if result.UnknownDurationSec < 0 {
		result.UnknownDurationSec = 0
	}

	for st, dur := range durationsMap {
		result.StateDurations = append(result.StateDurations, StateDuration{
			State:           st,
			DurationSeconds: dur,
			Percentage:      (dur / totalSeconds) * 100.0,
		})
	}

	return result
}
