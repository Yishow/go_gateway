package sourcerule

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
	"testing"
)

func TestSourceProjectionProofNeverSkipsForeignEpoch(t *testing.T) {
	svc := &Service{}
	ctx, release := svc.AcquireRuleMutation(t.Context(), "R")
	defer release()
	j := svc.pointJournal(ctx)
	j.deviceID = "D"
	j.projectionEpoch = 10
	j.projectionEpochKnown = true
	j.ownedEpochs = map[uint64]bool{12: true}
	j.points["absent"] = &ownedPointWrite{before: &schema.Point{ID: "absent", DeviceID: "D", RevisionID: "r"}, removedEpoch: 12}
	proof := svc.sourceProjectionProof(ctx)
	require.Equal(t, uint64(10), proof.Epoch, "foreign epoch 11 cannot be covered by later own receipt 12")
	j.ownedEpochs[11] = true
	require.Equal(t, uint64(12), svc.sourceProjectionProof(ctx).Epoch)
}

func TestSourceProjectionForeignCallbackDoesNotBecomeOwnedEpoch(t *testing.T) {
	svc := &Service{}
	ctx, release := svc.AcquireRuleMutation(t.Context(), "R")
	defer release()
	j := svc.pointJournal(ctx)
	j.deviceID = "D"
	j.points["owned"] = &ownedPointWrite{current: &schema.Point{ID: "owned", DeviceID: "D", RevisionID: "r"}}
	svc.observeRuntimeInstallation(ctx, &schema.Point{ID: "foreign", DeviceID: "D", RevisionID: "r"}, 11)
	svc.observeRuntimeInstallation(ctx, &schema.Point{ID: "owned", DeviceID: "D", RevisionID: "newer"}, 12)
	require.Empty(t, j.ownedEpochs)
	require.Zero(t, j.points["owned"].currentLease)
	svc.observeRuntimeInstallation(ctx, &schema.Point{ID: "owned", DeviceID: "D", RevisionID: "r"}, 13)
	require.Equal(t, map[uint64]bool{13: true}, j.ownedEpochs)
}

func TestCapturedProjectionRollbackRequiresRuntimeAuthority(t *testing.T) {
	svc := &Service{}
	require.Nil(t, svc.snapshotEmptyRuntimeProjection("D"))
	require.NoError(t, svc.restoreCapturedRuntimeProjection(t.Context(), nil))
	require.Error(t, svc.restoreCapturedRuntimeProjection(t.Context(), struct{}{}))
}
