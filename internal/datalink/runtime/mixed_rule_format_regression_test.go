package runtime

import (
	"testing"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestMixedRuleItemFormatControlsInheritedReadProfile(t *testing.T) {
	device := &schema.Device{ID: "d", Protocol: schema.ProtocolMC3E, Status: schema.DeviceStatusActive, ConnectionConfig: `{"data_format":"CDAB"}`}
	point := &schema.Point{ID: "p", DeviceID: "d", Address: "D100", DataType: schema.DataTypeUint32, FormatMode: schema.PointFormatInherit, FormatVersion: 1, RevisionID: "p1"}
	rule := &schema.SourceRule{ID: "r", DeviceID: "d", RevisionID: "r1", LayoutMode: "mixed", MixedItems: `[{"item_id":"value","address":"D100","data_type":"uint32","data_format":"ABCD"}]`}
	projection := &workspace.RuntimeProjection{DeviceIDs: []string{"d"}, Devices: []*schema.Device{device}, Points: []*schema.Point{point}, Rules: []*schema.SourceRule{rule}, RuleLinks: []*schema.SourceRuleLink{{RuleID: "r", PointID: "p", Address: "D100"}}}

	profile, err := projectionReadProfile(projection, point)
	require.NoError(t, err)
	require.Equal(t, "ABCD", profile.Resolution().EffectiveFormat)

	rule.MixedItems = `[{"item_id":"value","address":"D100","data_type":"uint32","data_format":"DCBA"}]`
	rule.RevisionID = "r2"
	profile, err = projectionReadProfile(projection, point)
	require.NoError(t, err)
	require.Equal(t, "DCBA", profile.Resolution().EffectiveFormat)

	point.FormatMode = schema.PointFormatOverride
	point.FormatOverride = "BADC"
	profile, err = projectionReadProfile(projection, point)
	require.NoError(t, err)
	require.Equal(t, "BADC", profile.Resolution().EffectiveFormat, "point intent must retain precedence")

	rule.MixedItems = `[{"item_id":"value","address":"D200","data_type":"uint32","data_format":"DCBA"}]`
	_, err = projectionReadProfile(projection, point)
	require.Error(t, err, "a link without its persisted item must fail closed")
}
