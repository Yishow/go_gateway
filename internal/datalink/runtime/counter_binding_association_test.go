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

	"github.com/stretchr/testify/require"
)

func TestCounterBindingAssociationFailureBlocksReplayBeforeDuplicateACK(t *testing.T) {
	f, store, key, descriptor := counterAcceptanceFixture(t)
	legacy := snapshot.Sample{SampleID: "legacy-1", MemberKey: `["device-A","point-int","tag-int"]`, ObservedAt: at(time.Second), Value: measurement.NewInt64(100), Quality: schema.QualityGood, SourceRevision: "src-int", MappingRevision: "map-int"}
	require.NoError(t, store.AppendSample(t.Context(), key, at(0), legacy))
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), key, descriptor))
	before, err := store.Restore(t.Context(), key)
	require.NoError(t, err)
	var payload, digest string
	require.NoError(t, store.DB().QueryRowContext(t.Context(), `SELECT payload,payload_digest FROM wg_delivery_samples`).Scan(&payload, &digest))
	_, err = store.DB().ExecContext(t.Context(), `CREATE TRIGGER fail_association BEFORE INSERT ON wg_counter_proof_associations BEGIN SELECT RAISE(ABORT,'injected association failure'); END`)
	require.NoError(t, err)
	boundary, err := NewGroupBoundary(f.config)
	require.Nil(t, boundary, "no memory or duplicate ACK path is exposed after failed publication")
	var blocked *GroupBoundaryError
	require.ErrorAs(t, err, &blocked)
	require.Equal(t, boundaryRestoreFailed, blocked.Code)
	require.Zero(t, countTable(t, store.DB(), "wg_counter_proof_associations"))
	after, err := store.Restore(t.Context(), key)
	require.NoError(t, err)
	require.Equal(t, before, after)
	_, err = store.DB().ExecContext(t.Context(), `DROP TRIGGER fail_association`)
	require.NoError(t, err)
	boundary, err = NewGroupBoundary(f.config)
	require.NoError(t, err)
	require.Equal(t, 1, countTable(t, store.DB(), "wg_counter_proof_associations"))
	require.NoError(t, boundary.AcceptSample(t.Context(), f.envelope(0, "legacy-1", at(time.Second), int64(100))))
	require.Equal(t, 1, countJournal(t, store.DB()))
	var original, originalDigest string
	require.NoError(t, store.DB().QueryRowContext(t.Context(), `SELECT payload,payload_digest FROM wg_delivery_samples`).Scan(&original, &originalDigest))
	require.Equal(t, payload, original)
	require.Equal(t, digest, originalDigest)
}

func TestCounterBindingLegacyReplayRequiresAssociationCapableLedger(t *testing.T) {
	f, store, key, descriptor := counterAcceptanceFixture(t)
	legacy := snapshot.Sample{SampleID: "legacy-1", MemberKey: `["device-A","point-int","tag-int"]`, ObservedAt: at(time.Second), Value: measurement.NewInt64(100), Quality: schema.QualityGood, SourceRevision: "src-int", MappingRevision: "map-int"}
	require.NoError(t, store.AppendSample(t.Context(), key, at(0), legacy))
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), key, descriptor))
	// This old capability-preserving test wrapper intentionally does not expose
	// the new association method; it must not silently upgrade a legacy sample.
	f.config.Ledger = duplicatingLedger{inner: f.config.Ledger}
	boundary, err := NewGroupBoundary(f.config)
	require.Nil(t, boundary)
	require.ErrorIs(t, err, counterproof.ErrCounterBindingProofMissing)
	require.Zero(t, countTable(t, store.DB(), "wg_counter_proof_associations"))
	require.Equal(t, 1, countJournal(t, store.DB()))
}

type unavailableAssociationLedger struct {
	*groupdelivery.Ledger
	calls int
}

func (l *unavailableAssociationLedger) PrepareCounterReplay(context.Context, []counterproof.CounterBindingReference) error {
	l.calls++
	return errors.New("legacy migration unavailable")
}

func TestCounterBindingVersionedReplayDoesNotRunLegacyMigration(t *testing.T) {
	f, store, key, descriptor := counterAcceptanceFixture(t)
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), key, descriptor))
	ledger := &unavailableAssociationLedger{Ledger: store.Ledger(key)}
	f.config.Ledger = ledger
	boundary, err := NewGroupBoundary(f.config)
	require.NoError(t, err)
	require.NoError(t, boundary.AcceptSample(t.Context(), f.envelope(0, "versioned-1", at(time.Second), int64(100))))
	recovered, err := NewGroupBoundary(f.config)
	require.NoError(t, err)
	require.NotNil(t, recovered)
	require.Zero(t, ledger.calls)
	require.Zero(t, countTable(t, store.DB(), "wg_counter_proof_associations"))
}
