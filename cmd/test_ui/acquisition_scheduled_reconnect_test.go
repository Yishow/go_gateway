package main

import (
	"context"
	"fmt"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestProductionScheduledAcquisitionSuccessfulReconnectExpiresOnlyConnectionDowngrade(t *testing.T) {
	for _, protocol := range []schema.ProtocolType{schema.ProtocolMC3E, schema.ProtocolFatekFBs} {
		for _, inplace := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/inplace=%t", protocol, inplace), func(t *testing.T) {
				peer, config := newAcquisitionWirePeer(t, protocol, map[int]uint16{0: 11, 10: 22, 20: 33})
				config = withAcquisitionDeclaration(config, declaredAcquisitionProfile(protocol))
				f := newScheduledProductionFixture(t, protocol, config)
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
				require.Len(t, peer.log()[baseline:], 4)
				for range 3 {
					<-f.capture.observations
				}
				conn, err := f.manager.GetOrCreate(t.Context(), "batch-device", protocol, config)
				require.NoError(t, err)
				oldToken, err := conn.CurrentAcquisitionToken(t.Context())
				require.NoError(t, err)
				if inplace {
					reconnector, ok := conn.Protocol.(interface{ Reconnect(context.Context) error })
					require.True(t, ok)
					require.NoError(t, reconnector.Reconnect(t.Context()))
				} else {
					require.NoError(t, f.manager.Reconnect(t.Context(), "batch-device"))
				}
				conn, err = f.manager.GetOrCreate(t.Context(), "batch-device", protocol, config)
				require.NoError(t, err)
				current, err := conn.CurrentAcquisitionToken(t.Context())
				require.NoError(t, err)
				require.NotEqual(t, oldToken, current)
				if inplace {
					require.Equal(t, oldToken.ManagerGeneration, current.ManagerGeneration)
					require.Greater(t, current.TransportGeneration, oldToken.TransportGeneration)
				} else {
					require.Greater(t, current.ManagerGeneration, oldToken.ManagerGeneration)
				}
				peer.setAction(nil)
				second := f.samples(t, 3)
				require.NoError(t, f.services.scheduler.Stop())
				requests := peer.log()[baseline:]
				require.Len(t, requests, 5)
				require.Equal(t, "random", requests[4].Kind, "successful reconnect retires only the old connection's unsupported-command evidence")
				firstByPoint := map[string]string{}
				for _, sample := range first {
					firstByPoint[sample.PointID] = sample.ConfigFingerprint
				}
				for _, sample := range second {
					require.Equal(t, schema.QualityGood, sample.Quality)
					require.NotEqual(t, first[0].AcquisitionID, sample.AcquisitionID)
					require.Equal(t, firstByPoint[sample.PointID], sample.ConfigFingerprint, "connection incarnation must not change source/sample configuration identity")
				}
			})
		}
	}
}
