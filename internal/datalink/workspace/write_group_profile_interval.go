package workspace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Interval bounds of the two profile kinds, in each kind's own public unit.
const (
	instantaneousIntervalMin     = 1
	instantaneousIntervalMax     = 300
	instantaneousIntervalDefault = 15
	cumulativeIntervalMin        = 1
	cumulativeIntervalMax        = 600
	cumulativeIntervalDefault    = 10
)

// ErrWriteGroupProfileIntervalInvalid is the stable profile_interval_invalid
// failure for an out-of-range, fractional or mis-unit profile interval.
var ErrWriteGroupProfileIntervalInvalid = fmt.Errorf("profile_interval_invalid: %w", ErrWriteGroupValidation)

var jsonIntegerLiteral = regexp.MustCompile(`^-?\d+$`)

// IntervalSeconds returns the execution interval in seconds, or zero when the
// unit is not one this contract knows.
func (p WriteGroupRecordingProfile) IntervalSeconds() int {
	switch p.Unit {
	case WriteGroupIntervalUnitSecond:
		return p.Interval
	case WriteGroupIntervalUnitMinute:
		return p.Interval * 60
	default:
		return 0
	}
}

// DefaultWriteGroupRecordingProfile returns the default of a new profile
// draft. It is a constructor default only: a request that omits or zeroes an
// interval is rejected, never completed from these values. A new instantaneous
// recording defaults to a mean over the 15-second interval; a cumulative
// profile counts and never averages.
func DefaultWriteGroupRecordingProfile(kind WriteGroupProfileKind) (WriteGroupRecordingProfile, error) {
	switch kind {
	case WriteGroupProfileInstantaneous:
		return WriteGroupRecordingProfile{
			Kind: kind, Interval: instantaneousIntervalDefault, Unit: WriteGroupIntervalUnitSecond,
			SampleReduction: WriteGroupSampleReductionMean,
		}, nil
	case WriteGroupProfileCumulative:
		return WriteGroupRecordingProfile{
			Kind: kind, Interval: cumulativeIntervalDefault, Unit: WriteGroupIntervalUnitMinute,
			SampleReduction: WriteGroupSampleReductionLatest,
		}, nil
	default:
		return WriteGroupRecordingProfile{}, fmt.Errorf("%w: unsupported profile kind %q", ErrWriteGroupValidation, kind)
	}
}

// validate checks the kind, its fixed unit, the integer interval bounds and
// the sample reduction: an unknown value is never read as latest, and a mean
// cumulative profile is rejected before anything can be saved.
func (p WriteGroupRecordingProfile) validate() error {
	var unit WriteGroupIntervalUnit
	var minimum, maximum int
	switch p.Kind {
	case WriteGroupProfileInstantaneous:
		unit, minimum, maximum = WriteGroupIntervalUnitSecond, instantaneousIntervalMin, instantaneousIntervalMax
	case WriteGroupProfileCumulative:
		unit, minimum, maximum = WriteGroupIntervalUnitMinute, cumulativeIntervalMin, cumulativeIntervalMax
	default:
		return fmt.Errorf("%w: unsupported profile kind %q", ErrWriteGroupValidation, p.Kind)
	}
	switch p.SampleReduction {
	case "", WriteGroupSampleReductionLatest, WriteGroupSampleReductionMean:
	default:
		return fmt.Errorf("%w: unsupported sample reduction %q", ErrWriteGroupValidation, p.SampleReduction)
	}
	if p.SampleReduction == WriteGroupSampleReductionMean && p.Kind != WriteGroupProfileInstantaneous {
		return fmt.Errorf("%w: mean reduction requires an instantaneous profile, not %q", ErrWriteGroupValidation, p.Kind)
	}
	if p.Unit != unit || p.Interval < minimum || p.Interval > maximum {
		return fmt.Errorf("%w: a %s profile needs an integer interval of %d to %d %ss",
			ErrWriteGroupProfileIntervalInvalid, p.Kind, minimum, maximum, unit)
	}
	return nil
}

// UnmarshalJSON accepts only a JSON integer literal for the interval and only
// a known sample reduction. A fraction, exponent, string, null or missing
// interval is rejected instead of being rounded, truncated or replaced by a
// default, and an absent reduction stays latest while an explicit empty,
// null or non-string one is rejected instead of degrading silently.
func (p *WriteGroupRecordingProfile) UnmarshalJSON(data []byte) error {
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		return nil
	}
	var raw struct {
		Kind            WriteGroupProfileKind  `json:"kind"`
		ResultShape     json.RawMessage        `json:"result_shape"`
		Interval        json.RawMessage        `json:"interval"`
		Unit            WriteGroupIntervalUnit `json:"unit"`
		SampleReduction json.RawMessage        `json:"sample_reduction"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	shape, err := decodeMemberEnum(raw.ResultShape, func(value string) error {
		if WriteGroupResultShape(value) != WriteGroupResultShapeMemberCalculationV1 {
			return fmt.Errorf("unsupported result shape %q", value)
		}
		return nil
	})
	if err != nil {
		return err
	}
	profile := WriteGroupRecordingProfile{Kind: raw.Kind}
	if shape != nil {
		profile.ResultShape = WriteGroupResultShape(*shape)
	}
	if err := validateCalculationProfileShape(&profile); err != nil {
		return err
	}
	if raw.Kind == WriteGroupProfileMixed {
		if raw.SampleReduction != nil || raw.Interval != nil || raw.Unit != "" {
			return fmt.Errorf("%w: mixed profile uses member calculations and intervals, not group reduction or interval", ErrWriteGroupValidation)
		}
		*p = profile
		return nil
	}
	text := string(bytes.TrimSpace(raw.Interval))
	if !jsonIntegerLiteral.MatchString(text) {
		return fmt.Errorf("%w: the interval must be a JSON integer", ErrWriteGroupProfileIntervalInvalid)
	}
	interval, err := strconv.Atoi(text)
	if err != nil {
		return fmt.Errorf("%w: the interval is out of range", ErrWriteGroupProfileIntervalInvalid)
	}
	reduction := WriteGroupSampleReduction("")
	if raw.SampleReduction != nil {
		wire := string(bytes.TrimSpace(raw.SampleReduction))
		if wire == "null" || !strings.HasPrefix(wire, `"`) {
			return fmt.Errorf("%w: the sample reduction must be the JSON strings %q or %q",
				ErrWriteGroupValidation, WriteGroupSampleReductionLatest, WriteGroupSampleReductionMean)
		}
		var parsed WriteGroupSampleReduction
		if err := json.Unmarshal(raw.SampleReduction, &parsed); err != nil {
			return fmt.Errorf("%w: the sample reduction must be the JSON strings %q or %q",
				ErrWriteGroupValidation, WriteGroupSampleReductionLatest, WriteGroupSampleReductionMean)
		}
		switch parsed {
		case WriteGroupSampleReductionLatest, WriteGroupSampleReductionMean:
			reduction = parsed
		default:
			return fmt.Errorf("%w: unsupported sample reduction %q", ErrWriteGroupValidation, parsed)
		}
	}
	*p = WriteGroupRecordingProfile{Kind: raw.Kind, Interval: interval, Unit: raw.Unit, SampleReduction: reduction}
	return nil
}
