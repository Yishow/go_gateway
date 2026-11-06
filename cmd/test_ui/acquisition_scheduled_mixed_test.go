package main

import (
	"fmt"
	"os"
	"testing"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"

	"github.com/stretchr/testify/require"
)

// Explicitly selected simulator measurement (task 5.4): a q-iqr declared MC
// device with D0–D31 plus a mixed rule of six scattered words. The local MC 3E
// peer counts the real wire requests of one production scheduled scan.
func TestProductionScheduledAcquisitionMixedRuleWireRequests(t *testing.T) {
	if os.Getenv("GATEWAY_MIXED_SIMULATOR") != "1" {
		t.Skip("set GATEWAY_MIXED_SIMULATOR=1 to measure mixed-rule wire requests on the local MC 3E peer")
	}
	scattered := []int{39, 40, 100, 200, 1000, 2000}
	memory := map[int]uint16{}
	for address := range 32 {
		memory[address] = uint16(address + 1)
	}
	for _, address := range scattered {
		memory[address] = uint16(address + 1)
	}
	peer, config := newAcquisitionWirePeer(t, schema.ProtocolMC3E, memory)
	f := newScheduledProductionFixture(t, schema.ProtocolMC3E, withAcquisitionDeclaration(config, `{"version":1,"profile":"mitsubishi-q-iqr-3e-binary-v1"}`))
	rangeRule, _ := f.createSource(t, "D0", 32, schema.DataTypeInt16, "range")
	applyAcquisitionTags(t, f.services, rangeRule)
	items := make([]measurement.MixedItem, 0, len(scattered))
	for _, address := range scattered {
		items = append(items, measurement.MixedItem{ItemID: fmt.Sprintf("d%d", address), Address: fmt.Sprintf("D%d", address), DataType: schema.DataTypeInt16, SemanticKind: measurement.SemanticKindGauge})
	}
	mixedRule, err := f.services.sourceRule.Create(t.Context(), sourcerule.CreateRuleRequest{DeviceID: "batch-device", NamingPrefix: "mixed", Enabled: true, LayoutMode: "mixed", MixedItems: items})
	require.NoError(t, err)
	applyAcquisitionTags(t, f.services, mixedRule)
	// The random plan itself does not depend on source geometry; what 5.4
	// changes is that the mixed rule no longer leaves the device's source
	// geometry unknown (source-layout-unsupported).
	projection, err := f.services.workspace.ReadRuntimeProjection(t.Context())
	require.NoError(t, err)
	graph, ok := projection.AcquisitionSourceGraph("batch-device")
	require.True(t, ok, "source geometry must be complete: %s", graph.Reason())

	baseline := len(peer.log())
	require.NoError(t, f.services.runtime.Start(t.Context()))
	samples := f.samples(t, 32+len(scattered))
	require.NoError(t, f.services.scheduler.Stop())
	requests := peer.log()[baseline:]
	values := map[any]bool{}
	for _, sample := range samples {
		require.Equal(t, schema.QualityGood, sample.Quality)
		require.Equal(t, samples[0].AcquisitionID, sample.AcquisitionID, "all points come from one scan")
		values[sample.Value] = true
	}
	for _, address := range scattered {
		require.True(t, values[int16(address+1)], "scattered D%d sampled", address)
	}
	for _, request := range requests {
		t.Logf("MC wire request kind=%s addresses=%d first=%v", request.Kind, len(request.Addresses), request.Addresses[:min(len(request.Addresses), 3)])
	}
	require.LessOrEqual(t, len(requests), 2, "D0–D31 plus six scattered addresses must take at most two requests")
}
