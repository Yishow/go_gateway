package main

import (
	"context"
	"errors"
	"testing"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"

	"github.com/stretchr/testify/require"
)

func TestProductionScheduledAcquisitionScopedRenameReacceptsCoverageWithoutReinstallingA(t *testing.T) {
	memory := map[int]uint16{}
	for address := range 21 {
		memory[address] = uint16(address + 100)
	}
	peer, config := newAcquisitionWirePeer(t, schema.ProtocolMC3E, memory)
	config = withAcquisitionDeclaration(config, `{"version":1,"profile":"safe-v1","coverage":{"revision":"declared-v1","readable":[{"area":"D","start":0,"end":21}]}}`)
	f := newScheduledProductionFixture(t, schema.ProtocolMC3E, config)
	a, al := f.createSource(t, "D0", 1, schema.DataTypeInt16, "a")
	b, _ := f.createSource(t, "D10", 1, schema.DataTypeInt16, "b")
	c, _ := f.createSource(t, "D20", 1, schema.DataTypeInt16, "c")
	for _, rule := range []*schema.SourceRule{a, b, c} {
		applyAcquisitionTags(t, f.services, rule)
	}
	require.NoError(t, f.services.runtime.Start(t.Context()))
	f.samples(t, 3)
	point, err := f.services.point.GetByID(t.Context(), al[0].PointID)
	require.NoError(t, err)
	leaseA := f.services.scheduler.CapturePointLease(point.ID, point.RevisionID)
	name := "b-renamed"
	_, outcome, err := f.services.sourceRule.UpdateWithRuntimeReconcile(t.Context(), b.ID, sourcerule.UpdateRuleRequest{NamingPrefix: &name})
	require.NoError(t, err)
	require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, outcome.Status)
	require.Equal(t, leaseA, f.services.scheduler.CapturePointLease(point.ID, point.RevisionID), "B reconcile must not reinstall A")
	baseline := len(peer.log())
	f.samples(t, 3)
	require.NoError(t, f.services.scheduler.Stop())
	require.Len(t, peer.log()[baseline:], 1, "fresh accepted B graph must restore declared coverage without widening B source scope")
}

func TestProductionScheduledAcquisitionCompensatedScopeRestoresCoverageWithFreshLeases(t *testing.T) {
	memory := map[int]uint16{}
	for address := range 21 {
		memory[address] = uint16(address + 100)
	}
	peer, config := newAcquisitionWirePeer(t, schema.ProtocolMC3E, memory)
	config = withAcquisitionDeclaration(config, `{"version":1,"profile":"safe-v1","coverage":{"revision":"declared-v1","readable":[{"area":"D","start":0,"end":21}]}}`)
	f := newScheduledProductionFixture(t, schema.ProtocolMC3E, config)
	var target *schema.SourceRule
	var targetID string
	for _, address := range []string{"D0", "D10", "D20"} {
		rule, links := f.createSource(t, address, 1, schema.DataTypeInt16, address)
		applyAcquisitionTags(t, f.services, rule)
		if address == "D10" {
			target = rule
			targetID = links[0].PointID
		}
	}
	require.NoError(t, f.services.runtime.Start(t.Context()))
	f.samples(t, 3)
	old, err := f.services.point.GetByID(t.Context(), targetID)
	require.NoError(t, err)
	oldLease := f.services.scheduler.CapturePointLease(old.ID, old.RevisionID)
	f.services.sourceRule.SetShareRuntimeReconciler(sourcerule.ShareRuntimeReconcilerFunc(func(context.Context, sourcerule.RuntimeReconcileRequest) (sourcerule.RuntimeReconcileOutcome, error) {
		return sourcerule.RuntimeReconcileOutcome{Status: sourcerule.RuntimeReconcileStatusFailed}, errors.New("fixture share publication refused")
	}))
	name := "compensated-name"
	_, _, err = f.services.sourceRule.UpdateWithRuntimeReconcile(t.Context(), target.ID, sourcerule.UpdateRuleRequest{NamingPrefix: &name})
	require.Error(t, err)
	restored, err := f.services.point.GetByID(t.Context(), targetID)
	require.NoError(t, err)
	require.Equal(t, old.RevisionID, restored.RevisionID)
	require.NotEqual(t, oldLease, f.services.scheduler.CapturePointLease(old.ID, old.RevisionID), "rollback owns fresh successful compensation leases, not copied old leases")
	baseline := len(peer.log())
	f.samples(t, 3)
	require.NoError(t, f.services.scheduler.Stop())
	require.Len(t, peer.log()[baseline:], 1, "old SQL witness and owned compensated install must restore coverage atomically")
}
