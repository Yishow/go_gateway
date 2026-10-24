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
	// Window is set only on the entries of a row whose window a cutover cut short.
	Window *longWindowProvenance `json:"window,omitempty"`
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
