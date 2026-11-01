package groupdelivery

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"testing"

	"go-gateway/internal/datalink/counterproof"

	"github.com/stretchr/testify/require"
)

func TestHandoffReceiptConstructionRejectsUnverifiedPayload(t *testing.T) {
	payload, digest, err := encodeSample(sampleAt("sample", "member", 3, 7))
	require.NoError(t, err)
	row := handoffSample{Key: testKey, Seq: 1, SampleID: "sample", Member: "member", Payload: string(payload), Digest: digest}
	receipt, err := newHandoffReceipt(row, row, otherKey, "fence")
	require.NoError(t, err)
	require.NotEmpty(t, receipt.ID)
	for _, bad := range []handoffSample{
		{Key: testKey, Payload: `{"private_secret":`, Digest: digest},
		{Key: testKey, Payload: string(payload), Digest: "wrong"},
	} {
		got, err := newHandoffReceipt(bad, row, otherKey, "fence")
		require.ErrorIs(t, err, counterproof.ErrCounterBindingProofConflict)
		require.Equal(t, handoffReceipt{}, got)
		require.NotContains(t, err.Error(), "private_secret")
	}
}

func TestHandoffIdentityPreservesOriginalEncoding(t *testing.T) {
	for _, targets := range [][]GroupKey{nil, {}, {testKey, otherKey}} {
		r := handoffReceipt{Fence: "fence", Source: testKey, Target: otherKey, Origin: testKey, SourceSeq: 1, OriginSeq: 2, SampleID: "sample", Member: "member", Targets: targets}
		// This is the original identity wire ordering; neither receipt status nor
		// the new error return may change persisted transfer IDs.
		original, err := json.Marshal([]any{1, r.Fence, r.Source, r.SourceSeq, r.Target, r.Origin, r.OriginSeq, r.SampleID, r.Member, r.ObservedAt, r.Digest, r.RawDigest, r.AcceptedAt, r.Targets})
		require.NoError(t, err)
		sum := sha256.Sum256(original)
		got, err := handoffID(r)
		require.NoError(t, err)
		require.Equal(t, hex.EncodeToString(sum[:]), got)
	}
}
