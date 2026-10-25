package dbtarget

import (
	"time"

	"go-gateway/internal/datalink/snapshot"
)

// windowStartOf is where the row's window begins: its own start, or the instant
// a cutover put there when the window was cut at its start.
func windowStartOf(outcome snapshot.Outcome) time.Time {
	start, _ := outcome.Window()
	return start.UTC()
}

// windowEvidence is the truncation evidence of a window a cutover cut short, or
// nil for any other window, so ordinary rows keep exactly the JSON they had.
func windowEvidence(outcome snapshot.Outcome) *longWindowProvenance {
	if outcome.Truncated == "" {
		return nil
	}
	start, end := outcome.Window()
	return &longWindowProvenance{
		Truncated: outcome.Truncated, Start: start.UTC().Format(longTimeText), End: end.UTC().Format(longTimeText),
		NominalStart: outcome.BucketStart.UTC().Format(longTimeText), NominalEnd: outcome.BucketEnd.UTC().Format(longTimeText),
	}
}

type provenanceEntry struct {
	Member     string `json:"member"`
	Status     string `json:"status"`
	Reason     string `json:"reason,omitempty"`
	SampleID   string `json:"sample_id,omitempty"`
	ObservedAt string `json:"observed_at,omitempty"`
	Quality    string `json:"quality,omitempty"`
	// Role, SourceType and ResultType are the frozen recording metadata of a
	// latest state or counter_raw member of a wide profile; Mean is the evidence of
	// a mean member. All stay absent for a legacy member, so its entry keeps
	// exactly the JSON it had.
	Role       string              `json:"role,omitempty"`
	SourceType string              `json:"source_type,omitempty"`
	ResultType string              `json:"result_type,omitempty"`
	Mean       *longMeanProvenance `json:"mean,omitempty"`
	// Counter is the cumulative evidence of a cumulative member: status, typed
	// exact difference, actual bounds and the policy it rests on.
	Counter *longCounterProvenance `json:"counter,omitempty"`
	// Window is set only on the entries of a row whose window a cutover cut short.
	Window *longWindowProvenance `json:"window,omitempty"`
}

// provenanceEntries is the member-level evidence of one row. A member that
// carries frozen recording metadata adds its role, types or mean evidence; every
// other member is described exactly as before.
func (l *GroupRowLayout) provenanceEntries(outcome snapshot.Outcome) []provenanceEntry {
	entries := provenanceEntries(outcome)
	for i := range entries {
		member, ok := l.members[entries[i].Member]
		if !ok {
			continue
		}
		if member.Counter != nil {
			entries[i].Counter = wideCounterProvenance(member.Counter, outcome.Members[i].Counter)
		}
		if member.Reduction == longReductionMean {
			entries[i].Mean = meanProvenance(member.ValueRole, member.SourceValueType, member.Type, outcome.Members[i])
		} else if member.SourceValueType != "" || member.ValueRole != "" {
			entries[i].Role, entries[i].SourceType, entries[i].ResultType = member.ValueRole, string(member.SourceValueType), string(member.Type)
		}
	}
	return entries
}

func provenanceEntries(outcome snapshot.Outcome) []provenanceEntry {
	entries := make([]provenanceEntry, 0, len(outcome.Members))
	window := windowEvidence(outcome)
	for _, result := range outcome.Members {
		entry := provenanceEntry{Member: result.MemberKey, Status: string(result.Status), Reason: result.Reason, Window: window}
		if result.Status == snapshot.MemberOK {
			entry.Reason = ""
		}
		if result.Sample != nil {
			entry.SampleID = result.Sample.SampleID
			entry.ObservedAt = result.Sample.ObservedAt.UTC().Format(time.RFC3339Nano)
			entry.Quality = string(result.Sample.Quality)
		}
		entries = append(entries, entry)
	}
	return entries
}
