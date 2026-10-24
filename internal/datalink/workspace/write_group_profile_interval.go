package workspace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
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
// interval is rejected, never completed from these values.
func DefaultWriteGroupRecordingProfile(kind WriteGroupProfileKind) (WriteGroupRecordingProfile, error) {
	switch kind {
	case WriteGroupProfileInstantaneous:
		return WriteGroupRecordingProfile{Kind: kind, Interval: instantaneousIntervalDefault, Unit: WriteGroupIntervalUnitSecond}, nil
	case WriteGroupProfileCumulative:
		return WriteGroupRecordingProfile{Kind: kind, Interval: cumulativeIntervalDefault, Unit: WriteGroupIntervalUnitMinute}, nil
	default:
		return WriteGroupRecordingProfile{}, fmt.Errorf("%w: unsupported profile kind %q", ErrWriteGroupValidation, kind)
	}
}

// validate checks the kind, its fixed unit and the integer interval bounds.
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
	if p.Unit != unit || p.Interval < minimum || p.Interval > maximum {
		return fmt.Errorf("%w: a %s profile needs an integer interval of %d to %d %ss",
			ErrWriteGroupProfileIntervalInvalid, p.Kind, minimum, maximum, unit)
	}
	return nil
}

// UnmarshalJSON accepts only a JSON integer literal for the interval. A
// fraction, exponent, string, null or missing interval is rejected instead of
// being rounded, truncated or replaced by a default.
func (p *WriteGroupRecordingProfile) UnmarshalJSON(data []byte) error {
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		return nil
	}
	var raw struct {
		Kind     WriteGroupProfileKind  `json:"kind"`
		Interval json.RawMessage        `json:"interval"`
		Unit     WriteGroupIntervalUnit `json:"unit"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	text := string(bytes.TrimSpace(raw.Interval))
	if !jsonIntegerLiteral.MatchString(text) {
		return fmt.Errorf("%w: the interval must be a JSON integer", ErrWriteGroupProfileIntervalInvalid)
	}
	interval, err := strconv.Atoi(text)
	if err != nil {
		return fmt.Errorf("%w: the interval is out of range", ErrWriteGroupProfileIntervalInvalid)
	}
	*p = WriteGroupRecordingProfile{Kind: raw.Kind, Interval: interval, Unit: raw.Unit}
	return nil
}
