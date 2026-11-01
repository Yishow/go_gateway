package main

import (
	"context"
	"errors"
	"testing"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

type acquisitionRollbackAcceptanceFailure struct {
	*workspace.Service
	calls  int
	refuse bool
}

func (r *acquisitionRollbackAcceptanceFailure) AcceptAcquisitionSourceGraph(ctx context.Context, graph workspace.AcquisitionSourceGraph, stage func() error) error {
	r.calls++
	if r.refuse && r.calls == 2 {
		return errors.New("fixture rollback SQL reservation unavailable")
	}
	return r.Service.AcceptAcquisitionSourceGraph(ctx, graph, stage)
}

func TestProductionScheduledAcquisitionFailedZeroPointRollbackRevokesForwardGapGrant(t *testing.T) {
	testProductionZeroPointRollback(t, true)
}

func TestProductionScheduledAcquisitionSuccessfulZeroPointRollbackKeepsRestoredGeometry(t *testing.T) {
	testProductionZeroPointRollback(t, false)
}

func testProductionZeroPointRollback(t *testing.T, refuse bool) {
	t.Helper()
	memory := map[int]uint16{}
	for address := range 21 {
		memory[address] = uint16(address + 100)
	}
	peer, config := newAcquisitionWirePeer(t, schema.ProtocolMC3E, memory)
	config = withAcquisitionDeclaration(config, scheduledReadableCoverage(schema.ProtocolMC3E))
	f := newScheduledProductionFixture(t, schema.ProtocolMC3E, config)
	createScheduledSparseSource(t, f, schema.ProtocolMC3E)
	skip, err := f.services.sourceRule.Create(t.Context(), sourcerule.CreateRuleRequest{DeviceID: "batch-device", StartAddress: "D5", Count: 1, DataType: schema.DataTypeInt16, Enabled: true, SkippedAddresses: []string{"D5"}})
	require.NoError(t, err)
	reader := &acquisitionRollbackAcceptanceFailure{Service: f.services.workspace, refuse: refuse}
	f.services.runtime.SetWorkspaceProjectionReader(reader)
	require.NoError(t, f.services.runtime.Start(t.Context()))
	f.samples(t, 3)
	f.services.sourceRule.SetShareRuntimeReconciler(sourcerule.ShareRuntimeReconcilerFunc(func(context.Context, sourcerule.RuntimeReconcileRequest) (sourcerule.RuntimeReconcileOutcome, error) {
		return sourcerule.RuntimeReconcileOutcome{Status: sourcerule.RuntimeReconcileStatusFailed}, errors.New("fixture share refused")
	}))
	_, err = f.services.sourceRule.DeleteWithRuntimeReconcile(t.Context(), skip.ID)
	require.Error(t, err)
	require.Equal(t, 2, reader.calls, "forward and compensated SQL acceptance each attempted once")
	_, err = f.services.sourceRule.GetByID(t.Context(), skip.ID)
	require.NoError(t, err, "durable skipped-only rule was compensated")
	baseline := len(peer.log())
	f.samples(t, 3)
	require.NoError(t, f.services.scheduler.Stop())
	requests := peer.log()[baseline:]
	if refuse {
		require.Len(t, requests, 3, "failed compensation cannot release the pending forward gap authorization")
	} else {
		require.Len(t, requests, 2, "final error cleanup must preserve the successfully restored old exclusion geometry")
	}
	for _, request := range requests {
		if refuse {
			require.Len(t, request.Addresses, 1)
		}
		for _, address := range request.Addresses {
			require.NotEqual(t, 5, address)
		}
	}
}
