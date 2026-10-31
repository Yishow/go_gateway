package main

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
)

func TestProductionWordFormatLegacyAndRegister16RemainCompatible(t *testing.T) {
	for _, protocol := range []schema.ProtocolType{schema.ProtocolModbusTCP, schema.ProtocolMC3E, schema.ProtocolFatekFBs} {
		t.Run(string(protocol), func(t *testing.T) {
			f := newWordProductionFixture(t, protocol)
			_, legacy, _ := f.createRule(t, protocol, 0, schema.DataTypeUint32, "ABCD")
			// Seed an actual pre-intent durable point, rather than treating a new override
			// as legacy. The production source materializer/restart must preserve this.
			legacy.FormatMode = schema.PointFormatLegacy
			legacy.FormatVersion = 0
			legacy.FormatOverride = ""
			legacy.DataFormat = "CDAB"
			require.NoError(t, f.services.point.Restore(t.Context(), legacy))
			var want any = uint32(0x11223344)
			switch protocol {
			case schema.ProtocolModbusTCP:
				f.plc.set(0, []uint16{0x3344, 0x1122})
			case schema.ProtocolMC3E, schema.ProtocolFatekFBs:
				f.plc.set(0, []uint16{0x1122, 0x3344}) // MC connection ABCD / FATEK native D-word concatenation, ignore point CDAB
			}
			_, single, _ := f.createRule(t, protocol, 1, schema.DataTypeUint16, "")
			f.plc.set(8, []uint16{0x1234})
			require.NoError(t, f.services.runtime.Start(t.Context()))
			for _, tc := range []struct {
				p    *schema.Point
				want any
			}{{legacy, want}, {single, uint16(0x1234)}} {
				result := f.services.scheduler.PollNowContext(t.Context(), []string{tc.p.ID})
				require.Len(t, result, 1)
				require.Equal(t, schema.QualityGood, result[0].Quality)
				require.Equal(t, tc.want, result[0].Value)
				raw, err := json.Marshal(tc.want)
				require.NoError(t, err)
				require.Eventually(t, func() bool {
					p, e := f.services.point.GetByID(t.Context(), tc.p.ID)
					return e == nil && p.LastValue != nil && *p.LastValue == string(raw)
				}, 5*time.Second, 20*time.Millisecond)
			}
		})
	}
}

func TestProductionWordFormatSourceDRRemainsExplicitlyUnsupported(t *testing.T) {
	f := newWordProductionFixture(t, schema.ProtocolFatekFBs)
	// Adapter/direct DR support does not silently widen the source authoring
	// address grammar. A refused rule cannot leave hidden points behind.
	_, err := f.services.sourceRule.Create(t.Context(), sourcerule.CreateRuleRequest{DeviceID: "word-device", StartAddress: "DR192", Count: 1, DataType: schema.DataTypeUint64, NamingPrefix: "unsupported", Enabled: true, DataFormat: "DCBA"})
	require.ErrorIs(t, err, sourcerule.ErrValidation)
	var rules, points int
	require.NoError(t, f.db.QueryRow(`SELECT count(*) FROM source_rules`).Scan(&rules))
	require.NoError(t, f.db.QueryRow(`SELECT count(*) FROM points`).Scan(&points))
	require.Zero(t, rules)
	require.Zero(t, points)
}
