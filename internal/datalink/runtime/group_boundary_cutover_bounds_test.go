package runtime

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// This is a constructor/replay test seam, not a replacement SQL authority.
// Durable cutover transactions and closure guards are covered by the real
// SQLite handoff fixtures; this ledger supplies an already accepted set only.
type boundsReplayLedger struct{ restored groupdelivery.Restored }

func (l *boundsReplayLedger) Restore(context.Context) (groupdelivery.Restored, error) {
	return l.restored, nil
}
func (l *boundsReplayLedger) AppendSample(context.Context, time.Time, snapshot.Sample) error {
	return errors.New("unexpected test append")
}
func (l *boundsReplayLedger) CommitClosure(context.Context, groupdelivery.Closure) error {
	return errors.New("unexpected test closure")
}

func TestGroupBoundaryActualBoundsReplayKeepsEarlierLatest(t *testing.T) {
	for _, shape := range []string{"legacy wide", "managed long", "managed wide"} {
		t.Run(shape, func(t *testing.T) {
			var cfg GroupBoundaryConfig
			switch shape {
			case "legacy wide":
				cfg = newBoundaryFixture(t).config
			case "managed long":
				cfg = newLongFixture(t).config
			case "managed wide":
				cfg = wideProfileFixture(t).config
				// Initial supported latest fixture, before any accepted journal.
				cfg.Group.RecordingProfile.SampleReduction = workspace.WriteGroupSampleReductionLatest
				cfg.Group.Members[0].SourceValueType = schema.DataTypeUint64
				cfg.Group.Members[0].ResultValueType = measurement.ExactUint64
				cfg.TagTypes[cfg.Group.Members[0].TagID] = schema.DataTypeUint64
			}
			cfg.Group.Members = cfg.Group.Members[:1]
			m := cfg.Group.Members[0]
			cfg.FirstBucket = boundaryStart
			interval := time.Duration(cfg.Group.RowPolicy.IntervalSeconds) * time.Second
			end := boundaryStart.Add(12500 * time.Millisecond)
			if interval == 10*time.Second {
				cfg.FirstBucket = boundaryAt(10)
			}
			cfg.EndsAt, cfg.Until = end, snapshot.BucketStart(end, interval).Add(interval)
			cfg.Clock = func() time.Time { return boundaryAt(200) }
			key := `["` + m.DeviceID + `","` + m.PointID + `","` + m.TagID + `"]`
			before := snapshot.Sample{SampleID: "before", MemberKey: key, ObservedAt: end.Add(-time.Millisecond), ReceivedAt: boundaryAt(11), Quality: schema.QualityGood, Value: boundsFloat(t, 100), SourceRevision: m.SourceRevision, MappingRevision: m.MappingRevision}
			after := before
			after.SampleID, after.ObservedAt, after.Value = "after", boundaryAt(13), boundsFloat(t, 200)
			if shape == "managed wide" {
				before.Value, after.Value = measurement.NewUint64(100), measurement.NewUint64(200)
			}
			ledger := &boundsReplayLedger{restored: groupdelivery.Restored{Samples: []snapshot.Sample{after, before}}}
			cfg.Ledger = ledger
			b, err := NewGroupBoundary(cfg)
			require.NoError(t, err)
			plan := b.assembler.Plan(boundaryAt(200))
			require.Len(t, plan.Outcomes, 1)
			require.Equal(t, before, *plan.Outcomes[0].Members[0].Sample)
			require.Equal(t, snapshot.OfferSelected, b.assembler.Check(after, after.ObservedAt).Outcome, "outside original was not replayed into the eligible memory buffer")
			require.Equal(t, []snapshot.Sample{after, before}, ledger.restored.Samples, "filter changes memory only, never the original accepted set")
			require.Equal(t, ledger.restored.HandoffView, b.handoffView)
		})
	}
}

func TestGroupBoundaryActualBoundsValidatesProofBeforeFiltering(t *testing.T) {
	f := newBoundaryFixture(t)
	f.config.Group.Members = f.config.Group.Members[:1]
	m := f.members[0]
	sample := snapshot.Sample{SampleID: "outside-corrupt", MemberKey: `["` + m.DeviceID + `","` + m.PointID + `","` + m.TagID + `"]`, ObservedAt: boundaryAt(9), Quality: schema.QualityGood, Value: boundsFloat(t, 99), SourceRevision: m.SourceRevision, MappingRevision: m.MappingRevision}
	sample.CounterProof.Version = 1 // A non-counter member cannot carry a counter proof.
	f.config.Ledger = &boundsReplayLedger{restored: groupdelivery.Restored{Samples: []snapshot.Sample{sample}}}
	f.config.EndsAt, f.config.Until = boundaryAt(5), boundaryAt(10)
	b, err := NewGroupBoundary(f.config)
	require.Nil(t, b)
	require.ErrorIs(t, err, counterproof.ErrCounterBindingProofConflict, "out-of-bounds filtering must not hide corrupt accepted evidence")
}

func boundsFloat(t *testing.T, value float64) measurement.ExactValue {
	t.Helper()
	exact, err := measurement.NewFloat64(value)
	require.NoError(t, err)
	return exact
}

func TestGroupBoundaryActualBoundsReplayPreservesOriginalOnTimeAcceptance(t *testing.T) {
	f := newBoundaryFixture(t)
	f.config.Group.Members = f.config.Group.Members[:1]
	f.config.Group.RowPolicy.IntervalSeconds = 5
	f.config.Group.RowPolicy.AllowedLatenessSeconds = 1
	f.config.FirstBucket, f.config.StartsAt, f.config.EndsAt, f.config.Until = boundaryAt(10), boundaryAt(10), boundaryAt(15), boundaryAt(15)
	f.config.Clock = func() time.Time { return boundaryAt(20) }
	m := f.members[0]
	sample := snapshot.Sample{SampleID: "on-time", MemberKey: memberKey(m.DeviceID, m.PointID, m.TagID), ObservedAt: boundaryAt(14), ReceivedAt: boundaryAt(14), Quality: schema.QualityGood, Value: boundsFloat(t, 7), SourceRevision: m.SourceRevision, MappingRevision: m.MappingRevision}
	f.config.Ledger = &boundsReplayLedger{restored: groupdelivery.Restored{Samples: []snapshot.Sample{sample}}}
	b, err := NewGroupBoundary(f.config)
	require.NoError(t, err)
	plan := b.assembler.Plan(boundaryAt(20))
	require.Len(t, plan.Outcomes, 1)
	require.Equal(t, sample, *plan.Outcomes[0].Members[0].Sample, "recovery at20 cannot turn original acceptance14/deadline16 into newlateinput")
	require.Zero(t, b.assembler.Late())
}

func TestGroupBoundaryActualBoundsDoesNotHideMalformedOriginalInput(t *testing.T) {
	for _, damage := range []string{"sample identity", "observation time", "source revision", "mapping revision"} {
		t.Run(damage, func(t *testing.T) {
			f := newBoundaryFixture(t)
			f.config.Group.Members = f.config.Group.Members[:1]
			m := f.members[0]
			sample := snapshot.Sample{SampleID: "outside", MemberKey: memberKey(m.DeviceID, m.PointID, m.TagID), ObservedAt: boundaryAt(9), Quality: schema.QualityGood, Value: boundsFloat(t, 99), SourceRevision: m.SourceRevision, MappingRevision: m.MappingRevision}
			switch damage {
			case "sample identity":
				sample.SampleID = ""
			case "observation time":
				sample.ObservedAt = time.Time{}
			case "source revision":
				sample.SourceRevision = "different-source"
			case "mapping revision":
				sample.MappingRevision = "different-mapping"
			}
			f.config.Ledger = &boundsReplayLedger{restored: groupdelivery.Restored{Samples: []snapshot.Sample{sample}}}
			f.config.StartsAt, f.config.EndsAt, f.config.Until = boundaryAt(2), boundaryAt(5), boundaryAt(10)
			b, err := NewGroupBoundary(f.config)
			require.Nil(t, b, "filtering cannot hide a preexisting replay validation failure")
			var blocked *GroupBoundaryError
			require.ErrorAs(t, err, &blocked)
			require.Equal(t, boundaryRestoreFailed, blocked.Code)
		})
	}
}
