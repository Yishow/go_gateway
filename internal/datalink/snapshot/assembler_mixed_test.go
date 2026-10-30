package snapshot

import (
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"testing"
	"time"
)

func mixedTime(seconds int) time.Time {
	return time.Date(2026, 10, 30, 12, 0, 0, 0, time.UTC).Add(time.Duration(seconds) * time.Second)
}
func mixedConfig(mode Mode) Config {
	return Config{
		WorkspaceID: "w", GroupID: "g", GroupRevision: "r", DestinationScope: "db", Mode: mode,
		StartsAt: mixedTime(0), Storage: StorageCapabilities{NullableValues: true, MemberQuality: true},
		Members: []Member{
			{Key: "a", Required: true, Calculation: CalculationMean, Interval: 15 * time.Second, SourceValueType: schema.DataTypeFloat64},
			{Key: "b", Required: true, Calculation: CalculationLatest, Interval: 60 * time.Second, SourceValueType: schema.DataTypeFloat64},
		},
	}
}
func mixedSample(key, id string, sec int, value float64) Sample {
	v, _ := measurement.NewFloat64(value)
	return Sample{MemberKey: key, SampleID: id, ObservedAt: mixedTime(sec), Quality: schema.QualityGood, Value: v}
}
func mixedOffer(t *testing.T, a *Assembler, s Sample) {
	t.Helper()
	if r := a.Offer(s, s.ObservedAt); r.Outcome != OfferSelected && r.Outcome != OfferNotSelected {
		t.Fatalf("offer %s: %+v", s.SampleID, r)
	}
}

func TestMixedIndependentWindowRetention(t *testing.T) {
	a, err := NewAssembler(mixedConfig(ModeMixedWide))
	if err != nil {
		t.Fatal(err)
	}
	mixedOffer(t, a, mixedSample("a", "a1", 10, 20))
	mixedOffer(t, a, mixedSample("b", "b1", 10, 100))
	rows := a.Tick(mixedTime(15))
	if len(rows) != 1 {
		t.Fatalf("rows=%+v", rows)
	}
	if rows[0].Kind != OutcomeRow || rows[0].Members[1].Status != MemberNotDue {
		t.Fatalf("required not-due suppressed row: %+v", rows[0])
	}
	if a.OpenBuckets() != 1 {
		t.Fatal("15-second close deleted the 60-second selector")
	}
	mixedOffer(t, a, mixedSample("b", "b2", 20, 101))
	mixedOffer(t, a, mixedSample("a", "a2", 50, 22))
	rows = a.Tick(mixedTime(60))
	last := rows[len(rows)-1]
	if last.Members[1].Sample == nil || last.Members[1].Sample.SampleID != "b2" {
		t.Fatalf("slow member lost: %+v", last)
	}
	if !last.DueAt.Equal(mixedTime(60)) || !last.BucketStart.IsZero() {
		t.Fatal("mixed wide claims a common window")
	}
	if !last.Members[0].Window.AlignedStart.Equal(mixedTime(45)) || !last.Members[1].Window.AlignedStart.Equal(mixedTime(0)) {
		t.Fatalf("wrong individual bounds: %+v", last.Members)
	}
}
