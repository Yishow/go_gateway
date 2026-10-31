package snapshot

import (
	"fmt"
	"reflect"
	"time"
)

// MixedCheckpointSchedule derives only expected closure identities and bounds
// with the existing bounded dispatch planner. No accepted sample or counter is
// replayed here; actual results remain those computed by the live assembler.
func MixedCheckpointSchedule(before, after MixedCheckpoint, intervals map[string]time.Duration, mode Mode) ([]Outcome, error) {
	config := Config{WorkspaceID: before.Progress.WorkspaceID, GroupID: before.Progress.GroupID, GroupRevision: before.Progress.GroupRevision, DestinationScope: "checkpoint-projection", StartsAt: before.StartsAt, EndsAt: before.EndsAt, Mode: mode, Storage: StorageCapabilities{NullableValues: true, MemberQuality: true}}
	for key, interval := range intervals {
		config.Members = append(config.Members, Member{Key: key, Interval: interval, Calculation: CalculationLatest, SourceValueType: "int64"})
	}
	a, err := NewAssembler(config)
	if err != nil {
		return nil, err
	}
	if err = a.RestoreMixedCheckpoint(before); err != nil {
		return nil, err
	}
	plan := a.Plan(after.Progress.ClosedThrough)
	if plan.Err != nil {
		return nil, plan.Err
	}
	cp, err := a.MixedCheckpointForPlan(plan)
	if err != nil {
		return nil, err
	}
	if !reflect.DeepEqual(cp, after) {
		return nil, fmt.Errorf("%w: incomplete mixed schedule", ErrInvalidConfig)
	}
	return plan.Outcomes, nil
}
