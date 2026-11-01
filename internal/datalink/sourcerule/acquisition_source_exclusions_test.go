package sourcerule

import (
	"testing"

	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func acquisitionGeometryFixture() AcquisitionSourceRows {
	return AcquisitionSourceRows{
		Device: schema.Device{ID: "d", Protocol: schema.ProtocolModbusTCP},
		Rules:  []*schema.SourceRule{{ID: "r", DeviceID: "d", StartAddress: "40001", Count: 3, DataType: schema.DataTypeUint16, Enabled: true, SkippedAddresses: `["40002"]`}},
		Points: []*schema.Point{
			{ID: "a", DeviceID: "d", Address: "40001", DataType: schema.DataTypeUint16, Enabled: true},
			{ID: "b", DeviceID: "d", Address: "40003", DataType: schema.DataTypeUint16, Enabled: true},
		},
		Links: []*schema.SourceRuleLink{{ID: "la", RuleID: "r", PointID: "a", Address: "40001"}, {ID: "lb", RuleID: "r", PointID: "b", Address: "40003"}},
	}
}

func TestAcquisitionSourceExclusionsKeepSkippedAndDisabledGeometry(t *testing.T) {
	input := acquisitionGeometryFixture()
	spans, reason := AcquisitionSourceExclusions(input)
	require.Empty(t, reason)
	require.Equal(t, []optimizer.NativeRange{{Area: "HR", Start: 1, End: 2}}, spans)
	input.Rules[0].Enabled = false
	for _, p := range input.Points {
		p.Enabled = false
	}
	spans, reason = AcquisitionSourceExclusions(input)
	require.Empty(t, reason)
	require.Equal(t, []optimizer.NativeRange{{Area: "HR", Start: 0, End: 1}, {Area: "HR", Start: 1, End: 2}, {Area: "HR", Start: 2, End: 3}}, spans)
}

func TestAcquisitionSourceExclusionsUnsupportedIsUnknown(t *testing.T) {
	for _, tc := range []struct {
		name   string
		mutate func(*AcquisitionSourceRows)
	}{
		{"missing-link", func(s *AcquisitionSourceRows) { s.Links = s.Links[:1] }},
		{"foreign-point", func(s *AcquisitionSourceRows) { s.Points[0].DeviceID = "foreign" }},
		{"duplicate-owner", func(s *AcquisitionSourceRows) { s.Links = append(s.Links, s.Links[0]) }},
		{"changed-address", func(s *AcquisitionSourceRows) { s.Points[0].Address = "40004" }},
		{"mixed-layout", func(s *AcquisitionSourceRows) { s.Rules[0].LayoutMode = "mixed" }},
		{"mixed-content", func(s *AcquisitionSourceRows) { s.Rules[0].MixedItems = `[]` }},
		{"overflow-start", func(s *AcquisitionSourceRows) { s.Rules[0].StartAddress = "499999999999999999999999999999" }},
		{"overflow-count", func(s *AcquisitionSourceRows) { s.Rules[0].Count = int(^uint(0) >> 1) }},
		{"legacy-start-clamp", func(s *AcquisitionSourceRows) {
			s.Rules[0].StartAddress, s.Rules[0].Count, s.Rules[0].SkippedAddresses = "40000", 1, ""
			s.Points, s.Links = s.Points[:1], s.Links[:1]
		}},
		{"bad-skip", func(s *AcquisitionSourceRows) { s.Rules[0].SkippedAddresses = `["4999999999999999999999999999"]` }},
		{"null-skip", func(s *AcquisitionSourceRows) { s.Rules[0].SkippedAddresses = `null` }},
		{"disabled-rule-active-point", func(s *AcquisitionSourceRows) { s.Rules[0].Enabled = false }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			input := acquisitionGeometryFixture()
			tc.mutate(&input)
			spans, reason := AcquisitionSourceExclusions(input)
			require.NotEmpty(t, reason)
			require.Empty(t, spans)
		})
	}
}

func TestAcquisitionSourceExclusionsSkipOnlyAndEmptySources(t *testing.T) {
	input := acquisitionGeometryFixture()
	input.Points, input.Links = nil, nil
	input.Rules[0].SkippedAddresses = `["40001","40002","40003"]`
	spans, reason := AcquisitionSourceExclusions(input)
	require.Empty(t, reason)
	require.Len(t, spans, 3)
	input.Rules = nil
	spans, reason = AcquisitionSourceExclusions(input)
	require.Empty(t, reason)
	require.Empty(t, spans)
}

func TestAcquisitionSourceDigestSeparatesSourceFromOutputAndReceipts(t *testing.T) {
	input := acquisitionGeometryFixture()
	semantic, err := AcquisitionSourceDigest(input, false)
	require.NoError(t, err)
	witness, err := AcquisitionSourceDigest(input, true)
	require.NoError(t, err)
	tag := "output-only"
	input.Links[0].TagID, input.Links[0].MappingID = &tag, &tag
	input.Rules[0].ShareEnabled = true
	input.Rules[0].RevisionID = "output-only-rule-revision"
	input.Rules[0].NamingPrefix = "output-only-name"
	got, err := AcquisitionSourceDigest(input, false)
	require.NoError(t, err)
	require.Equal(t, semantic, got)
	got, err = AcquisitionSourceDigest(input, true)
	require.NoError(t, err)
	require.Equal(t, witness, got)
	input.Points[0].RevisionID = "new-installation"
	got, _ = AcquisitionSourceDigest(input, false)
	require.Equal(t, semantic, got)
	got, _ = AcquisitionSourceDigest(input, true)
	require.NotEqual(t, witness, got)
	input.Rules[0].SkippedAddresses = `["40002","40005"]`
	got, _ = AcquisitionSourceDigest(input, false)
	require.NotEqual(t, semantic, got)
}
