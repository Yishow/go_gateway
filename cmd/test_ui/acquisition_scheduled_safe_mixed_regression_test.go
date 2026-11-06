package main

import (
	"testing"
	"time"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"

	"github.com/stretchr/testify/require"
)

func TestProductionScheduledUnsupportedMixedRuleIsNotRead(t *testing.T) {
	for _, tc := range []struct {
		name, declaration, regular, first, second string
		protocol                                  schema.ProtocolType
		legacy                                    bool
	}{
		{"safe", `{"version":1,"profile":"safe-v1"}`, "D0", "D39", "D100", schema.ProtocolMC3E, false},
		{"undeclared", "", "D0", "D39", "D100", schema.ProtocolMC3E, false},
		{"legacy-safe", `{"version":1,"profile":"safe-v1"}`, "D0", "D39", "D100", schema.ProtocolMC3E, true},
		{"modbus", "", "40001", "40040", "40101", schema.ProtocolModbusTCP, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			peer, config := newAcquisitionWirePeer(t, tc.protocol, map[int]uint16{0: 1, 39: 40, 100: 101})
			if tc.declaration != "" {
				config = withAcquisitionDeclaration(config, tc.declaration)
			}
			f := newScheduledProductionFixture(t, tc.protocol, config)
			regular, _ := f.createSource(t, tc.regular, 1, schema.DataTypeInt16, "regular")
			applyAcquisitionTags(t, f.services, regular)
			items := []measurement.MixedItem{
				{ItemID: "a", Address: tc.first, DataType: schema.DataTypeInt16, SemanticKind: measurement.SemanticKindGauge},
				{ItemID: "b", Address: tc.second, DataType: schema.DataTypeInt16, SemanticKind: measurement.SemanticKindGauge},
			}
			mixed, err := f.services.sourceRule.Create(t.Context(), sourcerule.CreateRuleRequest{DeviceID: "batch-device", NamingPrefix: "mixed", Enabled: true, LayoutMode: "mixed", MixedItems: items})
			require.NoError(t, err)
			applyAcquisitionTags(t, f.services, mixed)
			if tc.legacy {
				links, err := f.services.sourceRule.ListLinks(t.Context(), mixed.ID)
				require.NoError(t, err)
				for _, link := range links {
					_, err := f.db.ExecContext(t.Context(), `UPDATE points SET format_mode='legacy', format_override='', format_version=0 WHERE id=?`, link.PointID)
					require.NoError(t, err)
				}
			}
			_, err = f.services.point.SetDeviceAcquisitionInterval(t.Context(), "batch-device", 100)
			require.NoError(t, err)
			baseline := len(peer.log())
			require.NoError(t, f.services.runtime.Start(t.Context()))
			deadline := time.NewTimer(5 * time.Second)
			defer deadline.Stop()
			regularSamples := 0
			for {
				select {
				case sample := <-f.capture.samples:
					if sample.Quality == schema.QualityGood && sample.Value == int16(1) {
						regularSamples++
						if regularSamples < 4 {
							continue
						}
						require.NoError(t, f.services.scheduler.Stop())
						for _, request := range peer.log()[baseline:] {
							require.NotContains(t, request.Addresses, 39, "mixed point must not reach the device")
							require.NotContains(t, request.Addresses, 100, "mixed point must not reach the device")
						}
						return
					}
				case <-deadline.C:
					t.Fatal("supported homogeneous point was not acquired")
				}
			}
		})
	}
}
