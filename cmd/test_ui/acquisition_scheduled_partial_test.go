package main

import (
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestProductionScheduledAcquisitionFailedRequestDoesNotCorruptOtherPoints(t *testing.T) {
	for _, protocol := range []schema.ProtocolType{schema.ProtocolMC3E, schema.ProtocolFatekFBs, schema.ProtocolModbusTCP} {
		t.Run(string(protocol), func(t *testing.T) {
			peer, config := newAcquisitionWirePeer(t, protocol, map[int]uint16{0: 11, 10: 22, 20: 33})
			config = strings.Replace(config, `"timeout":5`, `"timeout":1`, 1)
			f := newScheduledProductionFixture(t, protocol, config)
			createScheduledSparseSource(t, f, protocol)
			// The native request deadline expires before the next real ticker due;
			// no retry count, synthetic clock or alternate scheduler is injected.
			groups, err := f.services.pollingGroup.List(t.Context())
			require.NoError(t, err)
			interval := 1500
			for _, group := range groups {
				_, err = f.services.pollingGroup.Update(t.Context(), group.ID, pollinggroup.UpdateRequest{IntervalMs: &interval})
				require.NoError(t, err)
			}
			peer.setAction(func(request acquisitionWireRequest) acquisitionWireAction {
				return acquisitionWireAction{malformed: request.Addresses[0] == 0}
			})
			baseline := len(peer.log())
			require.NoError(t, f.services.runtime.Start(t.Context()))
			samples := f.samples(t, 3)
			require.NoError(t, f.services.scheduler.Stop())
			requests := peer.log()[baseline:]
			require.Len(t, requests, 3, "a malformed reply cannot fan out a fallback or grant attempts after the original deadline")
			good, bad := 0, 0
			points, completions := map[string]bool{}, map[time.Time]bool{}
			for _, sample := range samples {
				require.False(t, points[sample.PointID])
				points[sample.PointID] = true
				completions[sample.ReceivedAt] = true
				require.Equal(t, samples[0].AcquisitionID, sample.AcquisitionID)
				if sample.Quality == schema.QualityBad {
					bad++
					require.Nil(t, sample.Value)
				} else {
					good++
					require.Equal(t, schema.QualityGood, sample.Quality)
					require.Contains(t, []any{int16(22), int16(33)}, sample.Value)
				}
			}
			require.Equal(t, 2, good)
			require.Equal(t, 1, bad)
			require.Len(t, completions, 3, "separate requests retain separate gateway completion times, not an invented atomic source snapshot")
			requestIDs := map[string]bool{}
			for range 3 {
				observation := <-f.capture.observations
				require.False(t, requestIDs[observation.RequestID])
				requestIDs[observation.RequestID] = true
			}
			stats, exists := f.services.scheduler.GetDeviceBreakerStats("batch-device")
			require.True(t, exists)
			require.Equal(t, 3, stats.TrackerStats.TotalRequests)
			require.Equal(t, 2, stats.TrackerStats.SuccessCount)
			require.Equal(t, 1, stats.TrackerStats.FailureCount)
		})
	}
}
