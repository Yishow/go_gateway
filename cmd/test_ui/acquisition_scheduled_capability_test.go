package main

import (
	"fmt"
	"strings"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func declaredAcquisitionProfile(protocol schema.ProtocolType) string {
	if protocol == schema.ProtocolMC3E {
		return `{"version":1,"profile":"mitsubishi-fx5-cpu-3e-binary-v1"}`
	}
	return `{"version":1,"profile":"fatek-fbs-ascii-v1"}`
}
func createScheduledSparseSource(t *testing.T, f scheduledProductionFixture, protocol schema.ProtocolType) {
	t.Helper()
	for i, offset := range []int{0, 10, 20} {
		address := fmt.Sprintf("D%d", offset)
		if protocol == schema.ProtocolModbusTCP {
			address = fmt.Sprint(40001 + offset)
		}
		rule, _ := f.createSource(t, address, 1, schema.DataTypeInt16, fmt.Sprintf("sparse-%d", i))
		applyAcquisitionTags(t, f.services, rule)
	}
}
func TestProductionScheduledAcquisitionDeclaredSparseVersusUnknown(t *testing.T) {
	for _, protocol := range []schema.ProtocolType{schema.ProtocolMC3E, schema.ProtocolFatekFBs} {
		for _, declared := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/declared=%t", protocol, declared), func(t *testing.T) {
				peer, config := newAcquisitionWirePeer(t, protocol, map[int]uint16{0: 11, 10: 22, 20: 33})
				if declared {
					config = withAcquisitionDeclaration(config, declaredAcquisitionProfile(protocol))
				}
				f := newScheduledProductionFixture(t, protocol, config)
				createScheduledSparseSource(t, f, protocol)
				baseline := len(peer.log())
				require.NoError(t, f.services.runtime.Start(t.Context()))
				samples := f.samples(t, 3)
				require.NoError(t, f.services.scheduler.Stop())
				requests := peer.log()[baseline:]
				want := 3
				if declared {
					want = 1
				}
				require.Len(t, requests, want)
				if declared {
					require.Equal(t, "random", requests[0].Kind)
					require.Equal(t, []int{0, 10, 20}, requests[0].Addresses)
				} else {
					for _, request := range requests {
						require.Equal(t, "contiguous", request.Kind)
						require.Len(t, request.Addresses, 1)
					}
				}
				values, ids, requestIDs := map[any]bool{}, map[string]bool{}, map[string]bool{}
				for _, sample := range samples {
					require.Equal(t, schema.QualityGood, sample.Quality)
					require.Equal(t, samples[0].AcquisitionID, sample.AcquisitionID)
					require.False(t, ids[sample.SampleID])
					ids[sample.SampleID] = true
					values[sample.Value] = true
				}
				require.Equal(t, map[any]bool{int16(11): true, int16(22): true, int16(33): true}, values)
				for range 3 {
					observation := <-f.capture.observations
					require.Equal(t, samples[0].AcquisitionID, observation.AcquisitionID)
					require.NotEmpty(t, observation.RequestID)
					requestIDs[observation.RequestID] = true
				}
				require.Len(t, requestIDs, want)
				if declared {
					for _, sample := range samples {
						require.Equal(t, samples[0].ReceivedAt, sample.ReceivedAt)
					}
				}
			})
		}
	}
}

func TestProductionScheduledAcquisitionExplicitUnsupportedFallsBackOnceAndCaches(t *testing.T) {
	for _, protocol := range []schema.ProtocolType{schema.ProtocolMC3E, schema.ProtocolFatekFBs} {
		t.Run(string(protocol), func(t *testing.T) {
			peer, config := newAcquisitionWirePeer(t, protocol, map[int]uint16{0: 11, 10: 22, 20: 33})
			f := newScheduledProductionFixture(t, protocol, withAcquisitionDeclaration(config, declaredAcquisitionProfile(protocol)))
			createScheduledSparseSource(t, f, protocol)
			peer.setAction(func(request acquisitionWireRequest) acquisitionWireAction {
				if request.Kind == "random" {
					return acquisitionWireAction{reject: '9'}
				}
				return acquisitionWireAction{}
			})
			baseline := len(peer.log())
			require.NoError(t, f.services.runtime.Start(t.Context()))
			first := f.samples(t, 3)
			firstRequests := peer.log()[baseline:]
			require.Len(t, firstRequests, 4)
			require.Equal(t, "random", firstRequests[0].Kind)
			for _, request := range firstRequests[1:] {
				require.Equal(t, "contiguous", request.Kind)
				require.Len(t, request.Addresses, 1)
			}
			for _, sample := range first {
				require.Equal(t, schema.QualityGood, sample.Quality)
				require.Equal(t, first[0].AcquisitionID, sample.AcquisitionID)
			}
			for range 3 {
				observation := <-f.capture.observations
				require.Contains(t, observation.RequestID, "/fallback/")
			}
			second := f.samples(t, 3)
			require.NoError(t, f.services.scheduler.Stop())
			require.NotEqual(t, first[0].AcquisitionID, second[0].AcquisitionID)
			all := peer.log()[baseline:]
			require.Len(t, all, 7, "next same-connection scan must retain declared-command downgrade")
			for _, request := range all[4:] {
				require.Equal(t, "contiguous", request.Kind)
			}
			for range 3 {
				observation := <-f.capture.observations
				require.True(t, strings.Contains(observation.RequestID, "/safe/"))
			}
			for _, sample := range second {
				require.Equal(t, schema.QualityGood, sample.Quality)
				require.Equal(t, second[0].AcquisitionID, sample.AcquisitionID)
			}
		})
	}
}

func TestProductionScheduledAcquisitionFatekFourDoesNotAuthorizeFallback(t *testing.T) {
	peer, config := newAcquisitionWirePeer(t, schema.ProtocolFatekFBs, map[int]uint16{0: 11, 10: 22, 20: 33})
	// A one-second original request timeout prevents its default two-second retry
	// wait; the test does not alter retry policy or mint a new deadline.
	config = strings.Replace(config, `"timeout":5`, `"timeout":1`, 1)
	f := newScheduledProductionFixture(t, schema.ProtocolFatekFBs, withAcquisitionDeclaration(config, declaredAcquisitionProfile(schema.ProtocolFatekFBs)))
	createScheduledSparseSource(t, f, schema.ProtocolFatekFBs)
	peer.setAction(func(request acquisitionWireRequest) acquisitionWireAction {
		if request.Kind == "random" {
			return acquisitionWireAction{reject: '4'}
		}
		return acquisitionWireAction{}
	})
	baseline := len(peer.log())
	require.NoError(t, f.services.runtime.Start(t.Context()))
	samples := f.samples(t, 3)
	require.NoError(t, f.services.scheduler.Stop())
	requests := peer.log()[baseline:]
	require.NotEmpty(t, requests)
	for _, request := range requests {
		require.Equal(t, "random", request.Kind, "ambiguous status4 must never fan out safe requests")
	}
	for _, sample := range samples {
		require.Equal(t, schema.QualityBad, sample.Quality)
		require.Nil(t, sample.Value)
	}
}
