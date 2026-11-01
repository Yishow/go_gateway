package main

import (
	"fmt"
	"testing"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"

	"github.com/stretchr/testify/require"
)

func TestProductionScheduledAcquisitionScopedDisabledAndSkippedGeometry(t *testing.T) {
	for _, operation := range []string{"disable-enable", "skip-only-create-delete"} {
		t.Run(operation, func(t *testing.T) {
			memory := map[int]uint16{}
			for address := range 21 {
				memory[address] = uint16(address + 100)
			}
			peer, config := newAcquisitionWirePeer(t, schema.ProtocolMC3E, memory)
			config = withAcquisitionDeclaration(config, `{"version":1,"profile":"safe-v1","coverage":{"revision":"declared-v1","readable":[{"area":"D","start":0,"end":21}]}}`)
			f := newScheduledProductionFixture(t, schema.ProtocolMC3E, config)
			var target *schema.SourceRule
			var a *schema.Point
			for i, address := range []string{"D0", "D10", "D20"} {
				rule, links := f.createSource(t, address, 1, schema.DataTypeInt16, fmt.Sprint("source-", i))
				applyAcquisitionTags(t, f.services, rule)
				if i == 1 {
					target = rule
				}
				if i == 0 {
					var err error
					a, err = f.services.point.GetByID(t.Context(), links[0].PointID)
					require.NoError(t, err)
				}
			}
			require.NoError(t, f.services.runtime.Start(t.Context()))
			f.samples(t, 3)
			leaseA := f.services.scheduler.CapturePointLease(a.ID, a.RevisionID)
			count := 3
			if operation == "disable-enable" {
				outcome, err := f.services.sourceRule.DisableWithRuntimeReconcile(t.Context(), target.ID)
				require.NoError(t, err)
				require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, outcome.Status)
				count = 2
			} else {
				var err error
				var outcome sourcerule.RuntimeReconcileOutcome
				target, outcome, err = f.services.sourceRule.CreateWithRuntimeReconcile(t.Context(), sourcerule.CreateRuleRequest{DeviceID: "batch-device", StartAddress: "D5", Count: 1, DataType: schema.DataTypeInt16, Enabled: true, SkippedAddresses: []string{"D5"}})
				require.NoError(t, err)
				require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, outcome.Status)
			}
			require.Equal(t, leaseA, f.services.scheduler.CapturePointLease(a.ID, a.RevisionID))
			baseline := len(peer.log())
			f.samples(t, count)
			requests := peer.log()[baseline:]
			require.Len(t, requests, 2, "installed disabled/skipped exclusion must split the declared readable region")
			for _, request := range requests {
				for _, address := range request.Addresses {
					if operation == "disable-enable" {
						require.NotEqual(t, 10, address)
					} else {
						require.NotEqual(t, 5, address)
					}
				}
			}
			if operation == "disable-enable" {
				outcome, err := f.services.sourceRule.EnableWithRuntimeReconcile(t.Context(), target.ID)
				require.NoError(t, err)
				require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, outcome.Status)
			} else {
				outcome, err := f.services.sourceRule.DeleteWithRuntimeReconcile(t.Context(), target.ID)
				require.NoError(t, err)
				require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, outcome.Status)
			}
			require.Equal(t, leaseA, f.services.scheduler.CapturePointLease(a.ID, a.RevisionID))
			baseline = len(peer.log())
			f.samples(t, 3)
			require.NoError(t, f.services.scheduler.Stop())
			require.Len(t, peer.log()[baseline:], 1, "fresh complete witness must release only the removed exclusion, never reinstall A")
		})
	}
}
