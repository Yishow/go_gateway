package workspace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"go-gateway/internal/datalink/snapshot"
)

// WriteGroupAlignmentUTCEpoch fixes member window phase independently of activation.
const WriteGroupAlignmentUTCEpoch = "utc_epoch"

// WriteGroupMemberCadence is member-owned intent. Nil preserves legacy absence;
// explicit mixed drafts must supply a complete cadence rather than a default.
type WriteGroupMemberCadence struct {
	Interval  int64                  `json:"interval"`
	Unit      WriteGroupIntervalUnit `json:"unit"`
	Alignment string                 `json:"alignment"`
}

// Duration validates exact integer intent against the backend representation
// bound. This is not a public advertisement of runtime capability.
func (c WriteGroupMemberCadence) Duration() (time.Duration, error) {
	if c.Alignment != WriteGroupAlignmentUTCEpoch {
		return 0, fmt.Errorf("%w: unsupported member alignment %q", ErrWriteGroupValidation, c.Alignment)
	}
	duration, err := snapshot.CheckedCadenceDuration(c.Interval, string(c.Unit))
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrWriteGroupValidation, err)
	}
	return duration, nil
}

func (c *WriteGroupMemberCadence) UnmarshalJSON(data []byte) error {
	var raw struct {
		Interval  json.RawMessage        `json:"interval"`
		Unit      WriteGroupIntervalUnit `json:"unit"`
		Alignment string                 `json:"alignment"`
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&raw); err != nil {
		return err
	}
	text := string(bytes.TrimSpace(raw.Interval))
	if !jsonIntegerLiteral.MatchString(text) {
		return fmt.Errorf("%w: member interval must be a JSON integer", ErrWriteGroupValidation)
	}
	value, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return fmt.Errorf("%w: member interval exceeds int64", ErrWriteGroupValidation)
	}
	result := WriteGroupMemberCadence{Interval: value, Unit: raw.Unit, Alignment: raw.Alignment}
	if _, err := result.Duration(); err != nil {
		return err
	}
	*c = result
	return nil
}

func validateMixedMemberCadences(group *WriteGroup) error {
	mixed := group.RecordingProfile != nil && group.RecordingProfile.Kind == WriteGroupProfileMixed
	for i, member := range group.Members {
		if !mixed {
			if member.Cadence != nil {
				return fmt.Errorf("%w: member %d cadence requires mixed profile", ErrWriteGroupValidation, i)
			}
			continue
		}
		if member.Cadence == nil {
			return fmt.Errorf("%w: member %d cadence is required", ErrWriteGroupValidation, i)
		}
		if _, err := member.Cadence.Duration(); err != nil {
			return fmt.Errorf("member %d: %w", i, err)
		}
	}
	return nil
}

func memberCadenceStorage(cadence *WriteGroupMemberCadence) (any, error) {
	if cadence == nil {
		return nil, nil
	}
	if _, err := cadence.Duration(); err != nil {
		return nil, err
	}
	raw, err := json.Marshal(cadence)
	if err != nil {
		return nil, err
	}
	return string(raw), nil
}

func sameMemberCadence(a, b *WriteGroupMemberCadence) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return *a == *b
}

// WindowAt computes metadata bounds without activating a runtime. The epoch
// phase is identical before/after restart and does not depend on activation.
func (c WriteGroupMemberCadence) WindowAt(at time.Time) (alignedStart, alignedEnd time.Time, boundsErr error) {
	duration, err := c.Duration()
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return snapshot.CheckedWindowBounds(at, duration)
}
