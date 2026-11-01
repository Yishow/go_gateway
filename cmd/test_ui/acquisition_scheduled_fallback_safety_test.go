package main

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"

	"github.com/stretchr/testify/require"
)

func TestProductionScheduledAcquisitionUnsupportedCannotFallbackAfterSourceDisable(t *testing.T) {
	for _, protocol := range []schema.ProtocolType{schema.ProtocolMC3E, schema.ProtocolFatekFBs} {
		t.Run(string(protocol), func(t *testing.T) {
			peer, config := newAcquisitionWirePeer(t, protocol, map[int]uint16{0: 11, 10: 22, 20: 33})
			f := newScheduledProductionFixture(t, protocol, withAcquisitionDeclaration(config, declaredAcquisitionProfile(protocol)))
			var rules []*schema.SourceRule
			for i, address := range []string{"D0", "D10", "D20"} {
				rule, _ := f.createSource(t, address, 1, schema.DataTypeInt16, fmt.Sprintf("disable-fallback-%d", i))
				applyAcquisitionTags(t, f.services, rule)
				rules = append(rules, rule)
			}
			entered, release := make(chan struct{}), make(chan struct{})
			unblock := sync.OnceFunc(func() { close(release) })
			t.Cleanup(unblock)
			var first sync.Once
			peer.setAction(func(request acquisitionWireRequest) acquisitionWireAction {
				first.Do(func() { close(entered) })
				return acquisitionWireAction{wait: release, reject: '9'}
			})
			baseline := len(peer.log())
			require.NoError(t, f.services.runtime.Start(t.Context()))
			select {
			case <-entered:
			case <-time.After(3 * time.Second):
				t.Fatal("scheduled declared random request was not sent")
			}
			disabled := false
			for _, rule := range rules {
				_, _, err := f.services.sourceRule.UpdateWithRuntimeReconcile(t.Context(), rule.ID, sourcerule.UpdateRuleRequest{Enabled: &disabled})
				require.NoError(t, err)
			}
			unblock()
			require.Eventually(t, func() bool { return f.services.runtime.RejectedObservationCount() == 3 }, time.Second, 5*time.Millisecond)
			require.NoError(t, f.services.scheduler.Stop())
			requests := peer.log()[baseline:]
			require.Len(t, requests, 1, "parsed unsupported command cannot bypass fresh source eligibility for any fallback send")
			require.Equal(t, "random", requests[0].Kind)
			require.Empty(t, f.capture.samples)
			stats, exists := f.services.scheduler.GetDeviceBreakerStats("batch-device")
			require.True(t, exists)
			require.Zero(t, stats.TrackerStats.TotalRequests, "superseded fallback admission is not a new device failure")
		})
	}
}
