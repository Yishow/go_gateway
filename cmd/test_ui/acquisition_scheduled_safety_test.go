package main

import (
	"context"
	"sync"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"

	"github.com/stretchr/testify/require"
)

func TestProductionScheduledAcquisitionOldInstallationCannotRefreshReplacement(t *testing.T) {
	for _, mutation := range []string{"same-revision-reinstall", "source-revision"} {
		t.Run(mutation, func(t *testing.T) {
			peer, config := newAcquisitionWirePeer(t, schema.ProtocolModbusTCP, map[int]uint16{0: 11, 1: 22, 2: 33, 3: 44, 4: 55, 5: 66})
			f := newScheduledProductionFixture(t, schema.ProtocolModbusTCP, config)
			rule, links := f.createSource(t, "40001", 3, schema.DataTypeUint32, "stale")
			applyAcquisitionTags(t, f.services, rule)
			entered, release := make(chan struct{}), make(chan struct{})
			closeRelease := sync.OnceFunc(func() { close(release) })
			t.Cleanup(closeRelease)
			var first sync.Once
			peer.setAction(func(acquisitionWireRequest) acquisitionWireAction {
				blocked := false
				first.Do(func() { blocked = true; close(entered) })
				if blocked {
					return acquisitionWireAction{wait: release}
				}
				return acquisitionWireAction{}
			})
			baseline := len(peer.log())
			require.NoError(t, f.services.runtime.Start(t.Context()))
			select {
			case <-entered:
			case <-time.After(3 * time.Second):
				t.Fatal("real scheduled adapter did not send")
			}
			before, err := f.services.point.GetByID(t.Context(), links[0].PointID)
			require.NoError(t, err)
			require.Nil(t, before.LastReadAt)
			oldLease, ok := f.services.scheduler.CaptureOrdinaryPoint(before.ID)
			require.True(t, ok)
			if mutation == "same-revision-reinstall" {
				projection, err := f.services.workspace.RuntimeProjection(t.Context())
				require.NoError(t, err)
				require.NoError(t, f.services.runtime.ApplyWorkspaceProjection(t.Context(), projection))
			} else {
				format := "CDAB"
				_, _, err := f.services.sourceRule.UpdateWithRuntimeReconcile(t.Context(), rule.ID, sourcerule.UpdateRuleRequest{DataFormat: &format})
				require.NoError(t, err)
			}
			newLease, ok := f.services.scheduler.CaptureOrdinaryPoint(before.ID)
			require.True(t, ok)
			require.NotEqual(t, oldLease.Lease, newLease.Lease)
			replacement, err := f.services.point.GetByID(t.Context(), before.ID)
			require.NoError(t, err)
			if mutation == "same-revision-reinstall" {
				require.Equal(t, before.RevisionID, replacement.RevisionID)
			} else {
				require.NotEqual(t, before.RevisionID, replacement.RevisionID)
			}
			closeRelease()
			require.Eventually(t, func() bool { return f.services.runtime.RejectedObservationCount() == 3 }, time.Second, 5*time.Millisecond, "old raw results must reach the original lease guard, not be rebound")
			require.NoError(t, f.services.scheduler.Stop())
			require.Len(t, peer.log()[baseline:], 1, "only the already-admitted old batch may send")
			for _, link := range links {
				p, err := f.services.point.GetByID(t.Context(), link.PointID)
				require.NoError(t, err)
				require.Nil(t, p.LastReadAt)
				require.Nil(t, p.LastValue)
			}
			require.Empty(t, f.capture.samples, "stale samples must not reach SampleSink with replacement bindings")
			require.Zero(t, f.services.runtime.Snapshot().WriteSuccess)
		})
	}
}

func TestProductionScheduledAcquisitionDisableDuringBackoffDeniesEveryLaterSend(t *testing.T) {
	for _, protocol := range []schema.ProtocolType{schema.ProtocolMC3E, schema.ProtocolFatekFBs, schema.ProtocolModbusTCP} {
		t.Run(string(protocol), func(t *testing.T) {
			peer, config := newAcquisitionWirePeer(t, protocol, map[int]uint16{0: 11, 1: 22, 2: 33})
			f := newScheduledProductionFixture(t, protocol, config)
			address := "D0"
			if protocol == schema.ProtocolModbusTCP {
				address = "40001"
			}
			rule, links := f.createSource(t, address, 3, schema.DataTypeInt16, "disable")
			applyAcquisitionTags(t, f.services, rule)
			sent := make(chan struct{}, 1)
			peer.setAction(func(acquisitionWireRequest) acquisitionWireAction {
				select {
				case sent <- struct{}{}:
				default:
				}
				return acquisitionWireAction{malformed: true}
			})
			baseline := len(peer.log())
			require.NoError(t, f.services.runtime.Start(t.Context()))
			select {
			case <-sent:
			case <-time.After(3 * time.Second):
				t.Fatal("no initial scheduled send")
			}
			disabled := false
			_, _, err := f.services.sourceRule.UpdateWithRuntimeReconcile(t.Context(), rule.ID, sourcerule.UpdateRuleRequest{Enabled: &disabled})
			require.NoError(t, err)
			// Configured request timeout is five seconds and retry delay one second.
			// This spans the original adapter/scheduler retry opportunity, not a new one.
			require.Eventually(t, func() bool { return f.services.runtime.RejectedObservationCount() == 3 }, 4*time.Second, 10*time.Millisecond)
			require.NoError(t, f.services.scheduler.Stop())
			require.Len(t, peer.log()[baseline:], 1, "disabled source must reject retry/fallback at every actual-send admission")
			for _, link := range links {
				p, err := f.services.point.GetByID(t.Context(), link.PointID)
				require.NoError(t, err)
				require.False(t, p.Enabled)
				require.Nil(t, p.LastReadAt)
			}
			require.Empty(t, f.capture.samples)
		})
	}
}

func TestProductionScheduledAcquisitionStopCancelsBlockedWireWithoutLateSample(t *testing.T) {
	for _, protocol := range []schema.ProtocolType{schema.ProtocolMC3E, schema.ProtocolFatekFBs, schema.ProtocolModbusTCP} {
		t.Run(string(protocol), func(t *testing.T) {
			peer, config := newAcquisitionWirePeer(t, protocol, map[int]uint16{0: 11, 1: 22, 2: 33})
			f := newScheduledProductionFixture(t, protocol, config)
			address := "D0"
			if protocol == schema.ProtocolModbusTCP {
				address = "40001"
			}
			rule, links := f.createSource(t, address, 3, schema.DataTypeInt16, "stop")
			applyAcquisitionTags(t, f.services, rule)
			entered, release := make(chan struct{}), make(chan struct{})
			closeRelease := sync.OnceFunc(func() { close(release) })
			t.Cleanup(closeRelease)
			var first sync.Once
			peer.setAction(func(acquisitionWireRequest) acquisitionWireAction {
				first.Do(func() { close(entered) })
				return acquisitionWireAction{wait: release}
			})
			baseline := len(peer.log())
			require.NoError(t, f.services.runtime.Start(t.Context()))
			select {
			case <-entered:
			case <-time.After(3 * time.Second):
				t.Fatal("no scheduled actual send")
			}
			ctx, cancel := context.WithTimeout(t.Context(), time.Second)
			defer cancel()
			require.NoError(t, f.services.runtime.Stop(ctx), "native stop must cancel and join original blocked adapter")
			closeRelease()
			require.Len(t, peer.log()[baseline:], 1)
			require.Empty(t, f.capture.samples)
			require.Empty(t, f.capture.observations)
			for _, link := range links {
				p, err := f.services.point.GetByID(t.Context(), link.PointID)
				require.NoError(t, err)
				require.Nil(t, p.LastReadAt)
			}
		})
	}
}
