package sourcerule

import (
	"testing"

	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestQIQRDisabledMixedBitUsesHexExclusionGeometry(t *testing.T) {
	rows := AcquisitionSourceRows{
		Device: schema.Device{ID: "dev", Protocol: schema.ProtocolMC3E, ConnectionConfig: mixedQIQRConfig},
		Rules:  []*schema.SourceRule{{ID: "rule", DeviceID: "dev", StartAddress: "X10", Count: 1, DataType: schema.DataTypeBool, LayoutMode: "mixed", MixedItems: `[{"item_id":"x","address":"X10","data_type":"bool","semantic_kind":"state"}]`}},
		Points: []*schema.Point{{ID: "point", DeviceID: "dev", Address: "X10", DataType: schema.DataTypeBool}},
		Links:  []*schema.SourceRuleLink{{ID: "link", RuleID: "rule", PointID: "point", Address: "X10"}},
	}
	excluded, reason := AcquisitionSourceExclusions(rows)
	require.Empty(t, reason)
	require.Contains(t, excluded, optimizer.NativeRange{Area: "X", Start: 0x10, End: 0x11})
}
