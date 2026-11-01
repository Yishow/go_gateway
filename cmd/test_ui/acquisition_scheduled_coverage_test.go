package main

import (
	"fmt"
	"testing"

	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestProductionScheduledAcquisitionUnverifiedDeclaredCoverageNeverReadsGap(t *testing.T) {
	for _, protocol := range []schema.ProtocolType{schema.ProtocolMC3E, schema.ProtocolFatekFBs, schema.ProtocolModbusTCP} {
		t.Run(string(protocol), func(t *testing.T) {
			// The peer has no fabricated memory for the gaps. A broad read is not
			// permitted merely because the active source installation omits them.
			peer, config := newAcquisitionWirePeer(t, protocol, map[int]uint16{0: 11, 10: 22, 20: 33})
			config = withAcquisitionDeclaration(config, scheduledReadableCoverage(protocol))
			f := newScheduledProductionFixture(t, protocol, config)
			createScheduledSparseSource(t, f, protocol)
			// An unlinked disabled source configuration is absent from the legacy
			// RuntimeProjection. Canonical SQL must expose that incomplete proof.
			address := "D5"
			if protocol == schema.ProtocolModbusTCP {
				address = "40006"
			}
			standalone, err := f.services.point.Create(t.Context(), point.CreatePointRequest{DeviceID: "batch-device", Address: address, DataType: schema.DataTypeInt16})
			require.NoError(t, err)
			disabled := false
			_, err = f.services.point.Update(t.Context(), standalone.ID, point.UpdatePointRequest{Enabled: &disabled})
			require.NoError(t, err)
			baseline := len(peer.log())
			require.NoError(t, f.services.runtime.Start(t.Context()))
			samples := f.samples(t, 3)
			require.NoError(t, f.services.scheduler.Stop())
			requests := peer.log()[baseline:]
			require.Len(t, requests, 3)
			for _, request := range requests {
				require.Equal(t, "contiguous", request.Kind)
				require.Len(t, request.Addresses, 1)
			}
			for _, sample := range samples {
				require.Equal(t, schema.QualityGood, sample.Quality)
				require.Empty(t, sample.QualityReason)
				require.Equal(t, samples[0].AcquisitionID, sample.AcquisitionID)
			}
			for range 3 {
				observation := <-f.capture.observations
				require.Equal(t, "source_coverage_unverified", observation.PlanReason)
				require.Empty(t, observation.QualityReason)
			}
		})
	}
}

func TestProductionScheduledAcquisitionVerifiedDeclaredCoverageReadsOneRequest(t *testing.T) {
	for _, protocol := range []schema.ProtocolType{schema.ProtocolMC3E, schema.ProtocolFatekFBs, schema.ProtocolModbusTCP} {
		t.Run(string(protocol), func(t *testing.T) {
			memory := map[int]uint16{}
			for address := range 21 {
				memory[address] = uint16(address + 100)
			}
			peer, config := newAcquisitionWirePeer(t, protocol, memory)
			config = withAcquisitionDeclaration(config, scheduledReadableCoverage(protocol))
			f := newScheduledProductionFixture(t, protocol, config)
			createScheduledSparseSource(t, f, protocol)
			projection, err := f.services.workspace.RuntimeProjection(t.Context())
			require.NoError(t, err)
			_, complete := projection.AcquisitionSourceGraph("batch-device")
			require.True(t, complete, "the original production SQL projection must establish completeness before runtime installation")
			baseline := len(peer.log())
			require.NoError(t, f.services.runtime.Start(t.Context()))
			samples := f.samples(t, 3)
			require.NoError(t, f.services.scheduler.Stop())
			requests := peer.log()[baseline:]
			require.Len(t, requests, 1, "SQL-backed source geometry must reach the original scheduled adapter path")
			require.Equal(t, "contiguous", requests[0].Kind)
			require.Len(t, requests[0].Addresses, 21)
			seen := map[string]bool{}
			requestID := ""
			for _, sample := range samples {
				require.Equal(t, schema.QualityGood, sample.Quality)
				require.Empty(t, sample.QualityReason)
				require.False(t, seen[sample.PointID])
				seen[sample.PointID] = true
				require.Equal(t, samples[0].AcquisitionID, sample.AcquisitionID)
				require.Equal(t, samples[0].ReceivedAt, sample.ReceivedAt)
				observation := <-f.capture.observations
				require.Empty(t, observation.PlanReason)
				if requestID == "" {
					requestID = observation.RequestID
				}
				require.Equal(t, requestID, observation.RequestID)
			}
		})
	}
}

func scheduledReadableCoverage(protocol schema.ProtocolType) string {
	area := "D"
	if protocol == schema.ProtocolModbusTCP {
		area = "HR"
	}
	return fmt.Sprintf(`{"version":1,"profile":"safe-v1","coverage":{"revision":"declared-v1","readable":[{"area":%q,"start":0,"end":21}]}}`, area)
}
