package sourcerule

import (
	"encoding/json"
	"fmt"
	"testing"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

const (
	mixedQIQRConfig  = `{"host":"127.0.0.1","port":5000,"network_no":0,"pc_no":255,"acquisition":{"version":1,"profile":"mitsubishi-q-iqr-3e-binary-v1"}}`
	mixedFX5Config   = `{"host":"127.0.0.1","port":5000,"network_no":0,"pc_no":255,"acquisition":{"version":1,"profile":"mitsubishi-fx5-cpu-3e-binary-v1"}}`
	mixedSafeConfig  = `{"host":"127.0.0.1","port":5000,"network_no":0,"pc_no":255,"acquisition":{"version":1,"profile":"safe-v1"}}`
	mixedLegacyMC    = `{"host":"127.0.0.1","port":5000,"network_no":0,"pc_no":255}`
	mixedFatekConfig = `{"host":"127.0.0.1","port":500,"acquisition":{"version":1,"profile":"fatek-fbs-ascii-v1"}}`
)

type mixedSourceItem struct {
	offset   int
	dataType schema.DataType
}

// The task scenario: D0–D31 homogeneous plus six scattered mixed items.
var mixedScenarioItems = []mixedSourceItem{{39, schema.DataTypeInt16}, {40, schema.DataTypeInt16}, {100, schema.DataTypeInt32}, {200, schema.DataTypeInt16}, {1000, schema.DataTypeInt16}, {2000, schema.DataTypeInt16}}

func mixedAcquisitionRows(t *testing.T, protocol schema.ProtocolType, config, area string) AcquisitionSourceRows {
	t.Helper()
	rows := AcquisitionSourceRows{Device: schema.Device{ID: "dev", Protocol: protocol, Status: schema.DeviceStatusActive, ConnectionConfig: config}}
	rows.Rules = append(rows.Rules, &schema.SourceRule{ID: "range", DeviceID: "dev", StartAddress: area + "0", Count: 32, DataType: schema.DataTypeInt16, Enabled: true, SkippedAddresses: "[]"})
	for i := range 32 {
		address := fmt.Sprintf("%s%d", area, i)
		rows.Points = append(rows.Points, &schema.Point{ID: "range-" + address, DeviceID: "dev", Address: address, DataType: schema.DataTypeInt16, Enabled: true})
		rows.Links = append(rows.Links, &schema.SourceRuleLink{ID: "l-range-" + address, RuleID: "range", Address: address, PointID: "range-" + address})
	}
	items := make([]measurement.MixedItem, 0, len(mixedScenarioItems))
	for _, entry := range mixedScenarioItems {
		address := fmt.Sprintf("%s%d", area, entry.offset)
		items = append(items, measurement.MixedItem{ItemID: "i" + address, Address: address, DataType: entry.dataType, SemanticKind: measurement.SemanticKindGauge})
		rows.Points = append(rows.Points, &schema.Point{ID: "mixed-" + address, DeviceID: "dev", Address: address, DataType: entry.dataType, Enabled: true})
		rows.Links = append(rows.Links, &schema.SourceRuleLink{ID: "l-mixed-" + address, RuleID: "mixed", Address: address, PointID: "mixed-" + address})
	}
	encoded, err := json.Marshal(items)
	require.NoError(t, err)
	rows.Rules = append(rows.Rules, &schema.SourceRule{ID: "mixed", DeviceID: "dev", StartAddress: items[0].Address, Count: len(items), DataType: items[0].DataType, Enabled: true, SkippedAddresses: "[]", LayoutMode: "mixed", MixedItems: string(encoded)})
	return rows
}

// mixedPlanRequests plans the enabled points the way the collector snapshot
// does: one source per point, profile capability, coverage plus exclusions.
func mixedPlanRequests(t *testing.T, rows AcquisitionSourceRows, excluded []optimizer.NativeRange) []optimizer.AcquisitionRequest {
	t.Helper()
	config, err := connector.ResolveAcquisitionConfig(rows.Device.Protocol, rows.Device.ConnectionConfig)
	require.NoError(t, err)
	coverage := config.Coverage()
	coverage.Excluded = append(coverage.Excluded, excluded...)
	input := optimizer.AcquisitionInput{ConfigurationRevision: config.Revision()}
	identity := config.Identity(rows.Device.ID)
	for _, point := range rows.Points {
		if !point.Enabled {
			continue
		}
		effective := acquisitionEffectivePoint(point)
		effective.Format, effective.FormatSource = "ABCD", "connection"
		bits, width, shapeErr := optimizer.AcquisitionPointShape(identity.Protocol, effective)
		require.NoError(t, shapeErr)
		effective.NativeWidth = width
		input.Sources = append(input.Sources, optimizer.AcquisitionSource{Identity: identity, Capability: config.Capability(bits), Coverage: coverage, Points: []optimizer.EffectivePoint{effective}})
	}
	plan := optimizer.BuildAcquisitionPlan(input)
	require.Empty(t, plan.Errors())
	return plan.Requests()
}

func TestAcquisitionSourceExclusionsPlanMixedRuleOnRandomProfile(t *testing.T) {
	rows := mixedAcquisitionRows(t, schema.ProtocolMC3E, mixedQIQRConfig, "D")
	excluded, reason := AcquisitionSourceExclusions(rows)
	require.Empty(t, reason)
	require.Empty(t, excluded)
	requests := mixedPlanRequests(t, rows, excluded)
	for _, request := range requests {
		t.Logf("request kind=%s area=%s start=%d points=%d units=%d", request.Kind, request.Group.Area, request.Start, len(request.Points), len(request.Units))
	}
	require.LessOrEqual(t, len(requests), 2, "D0–D31 plus six scattered items must fit two requests: %+v", requests)
	planned := 0
	for _, request := range requests {
		planned += len(request.Points)
	}
	require.Equal(t, len(rows.Points), planned, "every point is planned")
}

func TestAcquisitionSourceExclusionsMixedProfileMatrix(t *testing.T) {
	for _, tc := range []struct {
		name     string
		protocol schema.ProtocolType
		config   string
		area     string
		accepted bool
	}{
		{"q-iqr", schema.ProtocolMC3E, mixedQIQRConfig, "D", true},
		{"fx5", schema.ProtocolMC3E, mixedFX5Config, "D", true},
		{"fatek", schema.ProtocolFatekFBs, mixedFatekConfig, "R", true},
		{"safe-v1", schema.ProtocolMC3E, mixedSafeConfig, "D", false},
		{"undeclared", schema.ProtocolMC3E, mixedLegacyMC, "D", false},
		{"unreadable-config", schema.ProtocolMC3E, `{`, "D", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			excluded, reason := AcquisitionSourceExclusions(mixedAcquisitionRows(t, tc.protocol, tc.config, tc.area))
			if tc.accepted {
				require.Empty(t, reason)
				require.Empty(t, excluded)
				return
			}
			require.Equal(t, "source-layout-unsupported", reason)
			require.Empty(t, excluded)
		})
	}
}

func TestAcquisitionSourceExclusionsMixedRuleGeometry(t *testing.T) {
	disabled := mixedAcquisitionRows(t, schema.ProtocolMC3E, mixedQIQRConfig, "D")
	disabled.Rules[1].Enabled = false
	for _, point := range disabled.Points {
		if point.ID[:6] == "mixed-" {
			point.Enabled = false
		}
	}
	excluded, reason := AcquisitionSourceExclusions(disabled)
	require.Empty(t, reason)
	require.Equal(t, []optimizer.NativeRange{
		{Area: "D", Start: 39, End: 40}, {Area: "D", Start: 40, End: 41}, {Area: "D", Start: 100, End: 102},
		{Area: "D", Start: 200, End: 201}, {Area: "D", Start: 1000, End: 1001}, {Area: "D", Start: 2000, End: 2001},
	}, excluded, "a disabled mixed rule excludes exactly its item spans")

	for _, tc := range []struct {
		name   string
		reason string
		mutate func(*AcquisitionSourceRows)
	}{
		{"missing-link", acquisitionSourceLinksIncomplete, func(s *AcquisitionSourceRows) { s.Links = s.Links[:len(s.Links)-1] }},
		{"retyped-point", acquisitionSourceLinksIncomplete, func(s *AcquisitionSourceRows) { s.Points[len(s.Points)-1].DataType = schema.DataTypeInt32 }},
		{"foreign-link-address", acquisitionSourceLinksIncomplete, func(s *AcquisitionSourceRows) {
			s.Points[len(s.Points)-1].Address, s.Links[len(s.Links)-1].Address = "D3000", "D3000"
		}},
		{"disabled-rule-active-point", acquisitionSourceLinksIncomplete, func(s *AcquisitionSourceRows) { s.Rules[1].Enabled = false }},
		{"unreadable-items", acquisitionSourceRuleGeometryUnsupported, func(s *AcquisitionSourceRows) { s.Rules[1].MixedItems = `{` }},
		{"empty-items", acquisitionSourceRuleGeometryUnsupported, func(s *AcquisitionSourceRows) { s.Rules[1].MixedItems = `[]` }},
		{"bad-item-address", acquisitionSourceRuleGeometryUnsupported, func(s *AcquisitionSourceRows) {
			s.Rules[1].MixedItems = `[{"item_id":"x","address":"QQ1","data_type":"int16","semantic_kind":"gauge"}]`
		}},
		{"items-without-mixed-layout", "source-layout-unsupported", func(s *AcquisitionSourceRows) { s.Rules[1].LayoutMode = "" }},
		{"unknown-layout", "source-layout-unsupported", func(s *AcquisitionSourceRows) { s.Rules[1].LayoutMode = "scattered" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rows := mixedAcquisitionRows(t, schema.ProtocolMC3E, mixedQIQRConfig, "D")
			tc.mutate(&rows)
			spans, reason := AcquisitionSourceExclusions(rows)
			require.Equal(t, tc.reason, reason)
			require.Empty(t, spans)
		})
	}
}
