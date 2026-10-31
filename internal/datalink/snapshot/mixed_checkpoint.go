package snapshot

import (
	"fmt"
	"time"
)

// MixedCheckpoint is a versioned projection of committed mixed progress. It
// contains no samples or counter baselines; callers persist those atomically.
// GroupRevision must identify the immutable member descriptors being restored.
type MixedCheckpoint struct {
	Version  int                       `json:"version"`
	Progress DispatchWatermark         `json:"progress"`
	StartsAt time.Time                 `json:"starts_at"`
	EndsAt   time.Time                 `json:"ends_at"`
	Members  map[string]MemberProgress `json:"members"`
}

// MixedCheckpoint returns a detached projection of the current mixed cursor.
func (a *Assembler) MixedCheckpoint() (MixedCheckpoint, error) {
	progress, ok := a.MixedProgress()
	if !ok {
		return MixedCheckpoint{}, fmt.Errorf("%w: not mixed", ErrInvalidConfig)
	}
	members, err := a.MemberProgress()
	if err != nil {
		return MixedCheckpoint{}, err
	}
	return MixedCheckpoint{1, progress, a.config.StartsAt, a.config.EndsAt, members}, nil
}

// MixedCheckpointForPlan projects the cursor that Commit would install. It
// rejects stale plans and does not advance the assembler or counter baselines.
func (a *Assembler) MixedCheckpointForPlan(plan Plan) (MixedCheckpoint, error) {
	if plan.Err != nil {
		return MixedCheckpoint{}, plan.Err
	}
	if !a.config.Mode.mixed() || plan.mixedProgress == nil || !plan.mixedBase.Equal(a.mixedThrough) || plan.mixedProgress.ClosedThrough.Before(a.mixedThrough) {
		return MixedCheckpoint{}, fmt.Errorf("%w: foreign or stale mixed plan", ErrInvalidConfig)
	}
	projected := *a
	projected.mixedThrough = plan.mixedProgress.ClosedThrough
	projected.mixedNext = plan.mixedNext
	cp, err := projected.MixedCheckpoint()
	if err != nil {
		return MixedCheckpoint{}, err
	}
	cp.Progress = *plan.mixedProgress
	if _, err = a.validateMixedCheckpoint(cp); err != nil {
		return MixedCheckpoint{}, err
	}
	return cp, nil
}

func (a *Assembler) validateMixedCheckpoint(cp MixedCheckpoint) (*Assembler, error) {
	invalid := func(reason string) (*Assembler, error) {
		return nil, fmt.Errorf("%w: mixed checkpoint %s", ErrInvalidConfig, reason)
	}
	if !a.config.Mode.mixed() || cp.Version != 1 {
		return invalid("mode or version")
	}
	if !cp.StartsAt.Equal(a.config.StartsAt) {
		return invalid("activation mismatch")
	}
	if !a.config.EndsAt.IsZero() && (cp.EndsAt.IsZero() || cp.EndsAt.After(a.config.EndsAt)) {
		return invalid("extends configured end")
	}
	config := a.config
	config.EndsAt = cp.EndsAt
	candidate, err := NewAssembler(config)
	if err != nil {
		return nil, err
	}
	if err := candidate.RestoreMixedProgress(cp.Progress); err != nil {
		return nil, err
	}
	expected, err := candidate.MemberProgress()
	if err != nil {
		return nil, err
	}
	if len(cp.Members) != len(expected) {
		return invalid("member set mismatch")
	}
	for key, want := range expected {
		got, ok := cp.Members[key]
		if !ok || !got.Start.Equal(want.Start) || !got.End.Equal(want.End) {
			return invalid("member phase or bounds mismatch")
		}
	}
	return candidate, nil
}

// RestoreMixedCheckpoint restores only into an untouched mixed assembler. All
// scope, cutoff and member bounds are checked before any live state changes.
// A tighter persisted end survives restart; counter baselines are not touched.
func (a *Assembler) RestoreMixedCheckpoint(cp MixedCheckpoint) error {
	if len(a.buckets) != 0 || !a.mixedThrough.IsZero() {
		return fmt.Errorf("%w: mixed checkpoint restore requires untouched assembler", ErrInvalidConfig)
	}
	candidate, err := a.validateMixedCheckpoint(cp)
	if err != nil {
		return err
	}
	a.config.EndsAt = candidate.config.EndsAt
	a.mixedThrough = candidate.mixedThrough
	a.mixedNext = candidate.mixedNext
	return nil
}

// ValidateMixedCheckpointCadences checks a durable scheduling projection using
// the same cursor implementation as the live assembler. Calculation/type facts
// are not consumed here: no sample is offered and no reduction is performed.
func ValidateMixedCheckpointCadences(cp MixedCheckpoint, intervals map[string]time.Duration) error {
	config := Config{WorkspaceID: cp.Progress.WorkspaceID, GroupID: cp.Progress.GroupID, GroupRevision: cp.Progress.GroupRevision, DestinationScope: "checkpoint-projection", Storage: StorageCapabilities{NullableValues: true, MemberQuality: true}, StartsAt: cp.StartsAt, EndsAt: cp.EndsAt, Mode: ModeMixedWide}
	for key, interval := range intervals {
		config.Members = append(config.Members, Member{Key: key, Interval: interval, Calculation: CalculationLatest, SourceValueType: "int64"})
	}
	a, err := NewAssembler(config)
	if err != nil {
		return err
	}
	_, err = a.validateMixedCheckpoint(cp)
	return err
}
