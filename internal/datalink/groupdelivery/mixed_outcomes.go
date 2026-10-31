package groupdelivery

import (
	"go-gateway/internal/datalink/snapshot"
	"time"
)

func restoredThrough(r Restored) time.Time {
	if r.Mixed != nil {
		return r.Mixed.Progress.ClosedThrough
	}
	return r.NextClose
}
func validateMixedOutcomes(key GroupKey, c Closure, intervals map[string]time.Duration) error {
	seen := map[string]bool{}
	for _, b := range c.Buckets {
		o := b.Outcome
		if o.Kind != snapshot.OutcomeRow && o.Kind != snapshot.OutcomeSkipped && o.Kind != snapshot.OutcomeNoData {
			return mixedInvalid("unknown outcome kind")
		}
		if o.RecordID == "" || seen[o.RecordID] || o.DueAt.IsZero() || !o.DueAt.After(c.Mixed.Before.Progress.ClosedThrough) || o.DueAt.After(closureThrough(c)) {
			return mixedInvalid("outcome dispatch")
		}
		seen[o.RecordID] = true
		if (o.Kind == snapshot.OutcomeRow) != (b.Row != nil) {
			return mixedInvalid("row disposition")
		}
		if b.Row != nil && (b.Row.RecordID != o.RecordID || b.Row.EffectKey != o.EffectKey || !b.Row.DueAt.Equal(o.DueAt) || !b.Row.BucketStart.Equal(o.BucketStart)) {
			return mixedInvalid("encoded row scope")
		}
		var expected string
		var err error
		if o.BucketStart.IsZero() {
			expected, err = snapshot.MixedDispatchRecordID(key.WorkspaceID, key.GroupID, key.GroupRevision, o.DueAt)
		} else {
			if len(o.Members) != 1 {
				return mixedInvalid("long member count")
			}
			expected, err = snapshot.ProfileRecordID(key.WorkspaceID, key.GroupID, key.GroupRevision, o.Members[0].MemberKey, "mixed", o.BucketStart)
		}
		if err != nil || expected != o.RecordID {
			return mixedInvalid("record identity")
		}
		members := map[string]bool{}
		for _, m := range o.Members {
			interval, ok := intervals[m.MemberKey]
			if !ok || members[m.MemberKey] {
				return mixedInvalid("outcome member")
			}
			members[m.MemberKey] = true
			if m.Status == snapshot.MemberNotDue {
				if m.Window != nil && (!m.Window.ActualStart.IsZero() || !m.Window.ActualEnd.IsZero()) {
					return mixedInvalid("not due window")
				}
				continue
			}
			w := m.Window
			if w == nil || w.Interval != interval || !w.AlignedEnd.Equal(o.DueAt) {
				return mixedInvalid("member bounds")
			}
			start, end, err := snapshot.CheckedWindowBounds(w.AlignedStart, interval)
			if err != nil || !start.Equal(w.AlignedStart) || !end.Equal(w.AlignedEnd) {
				return mixedInvalid("member alignment")
			}
			if c.Mixed.After.StartsAt.After(start) {
				start = c.Mixed.After.StartsAt
			}
			if !c.Mixed.After.EndsAt.IsZero() && c.Mixed.After.EndsAt.Before(end) {
				end = c.Mixed.After.EndsAt
			}
			if !start.Before(end) || !w.ActualStart.Equal(start) || !w.ActualEnd.Equal(end) {
				return mixedInvalid("member actual bounds")
			}
			if m.Sample != nil && (m.Sample.ObservedAt.Before(start) || !m.Sample.ObservedAt.Before(end)) {
				return mixedInvalid("sample outside member bounds")
			}
		}
		if o.BucketStart.IsZero() && len(members) != len(intervals) {
			return mixedInvalid("wide member set")
		}
	}
	return nil
}
