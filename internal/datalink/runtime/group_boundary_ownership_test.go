package runtime

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/snapshot"
)

type ownershipRefusingLedger struct {
	GroupLedger
	refusal error
	appends int
}

func (l *ownershipRefusingLedger) AppendSample(ctx context.Context, start time.Time, sample snapshot.Sample) error {
	l.appends++
	if l.refusal != nil {
		return l.refusal
	}
	return l.GroupLedger.AppendSample(ctx, start, sample)
}

func TestGroupBoundaryOwnershipDuplicateRechecksDurableAuthority(t *testing.T) {
	f := newBoundaryFixture(t)
	db := openDurableDB(t, filepath.Join(t.TempDir(), "ownership.db"))
	f.withLedger(db)
	ledger := &ownershipRefusingLedger{GroupLedger: f.config.Ledger}
	f.config.Ledger = ledger
	b := f.build(t)
	env := f.envelope(0, "duplicate", 4, 21.5)
	require.NoError(t, b.AcceptOwnedSample(t.Context(), env))
	require.Equal(t, 1, ledger.appends)
	ledger.refusal = groupdelivery.ErrOwnershipConflict
	for _, accept := range []func(context.Context, measurement.SampleEnvelope) error{b.AcceptOwnedSample, b.AcceptSample} {
		err := accept(t.Context(), env)
		require.ErrorIs(t, err, groupdelivery.ErrOwnershipConflict)
		var owner *groupdelivery.OwnershipError
		require.ErrorAs(t, err, &owner)
		var refused *GroupSampleError
		require.ErrorAs(t, err, &refused)
		require.Equal(t, "recording-ownership-conflict", refused.Reason)
	}
	require.Equal(t, 3, ledger.appends)
	require.Equal(t, 1, countJournal(t, db))
	ledger.refusal = nil
	require.NoError(t, b.AcceptOwnedSample(t.Context(), env))
	require.Equal(t, 1, countJournal(t, db))
}

func TestGroupBoundaryOwnershipNeverAcknowledgesIgnoredSample(t *testing.T) {
	for _, mode := range []string{"no ledger", "before cached start", "at cached end", "unknown member"} {
		t.Run(mode, func(t *testing.T) {
			f := newBoundaryFixture(t)
			if mode != "no ledger" {
				f.withLedger(openDurableDB(t, filepath.Join(t.TempDir(), "ownership.db")))
			}
			env := f.envelope(0, "ignored", 4, 21.5)
			if mode == "before cached start" {
				f.config.FirstBucket = boundaryAt(10)
			}
			if mode == "unknown member" {
				env.TagID = "unowned"
			}
			b := f.build(t)
			if mode == "at cached end" {
				b.SetUntil(boundaryAt(4))
			}
			err := b.AcceptOwnedSample(t.Context(), env)
			want := groupdelivery.ErrOwnershipConflict
			if mode == "no ledger" {
				want = groupdelivery.ErrOwnershipUnavailable
			}
			require.ErrorIs(t, err, want)
		})
	}
}
