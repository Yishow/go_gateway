package groupdelivery

import (
	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/snapshot"
	"time"
)

func validateMixedSchedule(d *mixedDescriptor, c Closure, intervals map[string]time.Duration) error {
	mode := snapshot.ModeMixedWide
	if d.Group.LayoutVersion == mixedLongLayout {
		mode = snapshot.ModeMixedLong
	}
	expected, err := snapshot.MixedCheckpointSchedule(c.Mixed.Before, c.Mixed.After, intervals, mode)
	if err != nil {
		return err
	}
	if len(expected) != len(c.Buckets) {
		return mixedInvalid("incomplete closure schedule")
	}
	ids := map[string]snapshot.Outcome{}
	for _, o := range expected {
		ids[o.RecordID] = o
	}
	for _, b := range c.Buckets {
		expectedOutcome, ok := ids[b.Outcome.RecordID]
		if !ok {
			return mixedInvalid("foreign closure schedule")
		}
		if err := validateMixedMemberSchedule(expectedOutcome, b.Outcome); err != nil {
			return err
		}
		delete(ids, b.Outcome.RecordID)
	}
	methods := map[string]string{}
	delta := map[string]bool{}
	for _, m := range d.Group.Members {
		member := (counterproof.CounterBindingMemberIdentity{DeviceID: m.DeviceID, PointID: m.PointID, TagID: m.TagID}).SampleMemberKey()
		methods[member] = m.Calculation
		if m.Calculation == mixedDeltaMethod {
			key := (counterproof.CounterBindingMemberIdentity{DeviceID: m.DeviceID, PointID: m.PointID, TagID: m.TagID}).SampleMemberKey()
			delta[key] = true
		}
	}
	for _, b := range c.Buckets {
		for _, m := range b.Outcome.Members {
			if m.Window == nil || string(m.Window.Calculation) != methods[m.MemberKey] {
				return mixedInvalid("changed calculation")
			}
		}
	}
	if len(c.CounterProofs) != len(delta) {
		return mixedInvalid("counter proof set")
	}
	for _, b := range c.Baselines {
		cp, err := aggregation.DecodeCounterCheckpoint(b.Payload)
		if err != nil {
			return err
		}
		bound, ok := c.Mixed.After.Members[b.MemberKey]
		if !ok || !delta[b.MemberKey] || cp.WindowEnd.After(bound.Start) {
			return mixedInvalid("baseline exceeds member progress")
		}
	}
	return nil
}

// The existing planner is the sole cadence authority. Its empty-input results
// prove scheduling only: never compare missing/value status or recompute values.
func validateMixedMemberSchedule(expected, actual snapshot.Outcome) error {
	if !expected.DueAt.Equal(actual.DueAt) || !expected.BucketStart.Equal(actual.BucketStart) || len(expected.Members) != len(actual.Members) {
		return mixedInvalid("changed outcome schedule")
	}
	members := map[string]snapshot.MemberResult{}
	for _, member := range expected.Members {
		members[member.MemberKey] = member
	}
	for _, member := range actual.Members {
		wanted, ok := members[member.MemberKey]
		if !ok || (wanted.Status == snapshot.MemberNotDue) != (member.Status == snapshot.MemberNotDue) {
			return mixedInvalid("changed member due mask")
		}
		delete(members, member.MemberKey)
		w, got := wanted.Window, member.Window
		if w == nil || got == nil || w.Interval != got.Interval ||
			!w.AlignedStart.Equal(got.AlignedStart) || !w.AlignedEnd.Equal(got.AlignedEnd) ||
			!w.ActualStart.Equal(got.ActualStart) || !w.ActualEnd.Equal(got.ActualEnd) ||
			!w.NextDue.Equal(got.NextDue) || w.Truncated != got.Truncated {
			return mixedInvalid("changed scheduled member bounds")
		}
	}
	return nil
}
