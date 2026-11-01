package optimizer

import (
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestAcquisitionWireUnitsPhysicalCells(t *testing.T) {
	for _, tc := range []struct {
		name        string
		points      []EffectivePoint
		bits, units int
	}{
		{"DR", []EffectivePoint{{ID: "a", Address: "DR0", DataType: schema.DataTypeUint32, NativeWidth: 1}}, 16, 2},
		{"DR overlap phases", []EffectivePoint{{ID: "a", Address: "DR0", DataType: schema.DataTypeUint32, NativeWidth: 1}, {ID: "b", Address: "DR1", DataType: schema.DataTypeUint32, NativeWidth: 1}}, 16, 3},
		{"DR wide", []EffectivePoint{{ID: "a", Address: "DR0", DataType: schema.DataTypeUint64, NativeWidth: 2}}, 16, 4},
		{"WM", []EffectivePoint{{ID: "a", Address: "WM8", DataType: schema.DataTypeUint16, NativeWidth: 1}}, 1, 16},
		{"WM overlap", []EffectivePoint{{ID: "a", Address: "WM8", DataType: schema.DataTypeUint16, NativeWidth: 1}, {ID: "b", Address: "WM16", DataType: schema.DataTypeUint16, NativeWidth: 1}}, 1, 24},
		{"DWM", []EffectivePoint{{ID: "a", Address: "DWM8", DataType: schema.DataTypeUint64, NativeWidth: 2}}, 1, 64},
	} {
		t.Run(tc.name, func(t *testing.T) {
			in := randomInput(acquisitionProtocolFatek)
			in.Sources[0].Points = tc.points
			plan := BuildAcquisitionPlan(in)
			require.Empty(t, plan.Errors())
			units := plan.WireUnits()
			require.Len(t, units, 1)
			require.Equal(t, tc.bits, units[0].AddressUnitBits)
			require.Equal(t, tc.units, units[0].UniqueAddressUnits)
			units[0].UniqueAddressUnits = 999
			require.Equal(t, tc.units, plan.WireUnits()[0].UniqueAddressUnits)
		})
	}
}

func TestAcquisitionWireUnitsExcludeGapsAndRejectedScalars(t *testing.T) {
	in := randomInput(acquisitionProtocolFatek)
	s := &in.Sources[0]
	s.Capability.RandomSupport = "unknown"
	s.Points = []EffectivePoint{{ID: "a", Address: "DR0", DataType: schema.DataTypeUint32, NativeWidth: 1}, {ID: "b", Address: "DR4", DataType: schema.DataTypeUint32, NativeWidth: 1}, {ID: "bad", Address: "DR8", DataType: schema.DataTypeUint32, NativeWidth: 7}}
	s.Coverage = ReadableCoverage{Revision: "coverage", Readable: []NativeRange{{Area: "DR", Start: 2, End: 4}}}
	p := BuildAcquisitionPlan(in)
	require.Len(t, p.Requests(), 1)
	require.Len(t, p.Errors(), 1)
	require.Equal(t, 6, p.Requests()[0].Length*p.Requests()[0].AddressStep)
	require.Len(t, p.WireUnits(), 1)
	require.Equal(t, 4, p.WireUnits()[0].UniqueAddressUnits)
}

func TestAcquisitionWireUnitsCompatibilityNamespaces(t *testing.T) {
	base := randomInput(acquisitionProtocolFatek).Sources[0]
	base.Points = []EffectivePoint{{ID: "a", Address: "DR0", DataType: schema.DataTypeUint32, NativeWidth: 1}}
	for _, tc := range []struct {
		name   string
		change func(*AcquisitionSource)
	}{
		{"group", func(s *AcquisitionSource) { s.Identity.PollingGroup = "other" }},
		{"cadence", func(s *AcquisitionSource) { s.Identity.Cadence = time.Second * 7 }},
		{"function", func(s *AcquisitionSource) {
			s.Points = append([]EffectivePoint(nil), s.Points...)
			s.Points[0].Function = "other"
		}},
		{"station", func(s *AcquisitionSource) { s.Identity.Station = "other" }},
		{"connection", func(s *AcquisitionSource) { s.Identity.ConnectionRevision = "other" }},
		{"coverage declaration", func(s *AcquisitionSource) { s.Coverage.Readable = []NativeRange{{Area: "DR", Start: 5, End: 7}} }},
		{"capability declaration", func(s *AcquisitionSource) { s.Capability.DeviceLimit = 3 }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			other := base
			tc.change(&other)
			p := BuildAcquisitionPlan(AcquisitionInput{Sources: []AcquisitionSource{base, other}})
			require.Empty(t, p.Errors())
			require.Len(t, p.WireUnits(), 2)
			for _, u := range p.WireUnits() {
				require.Equal(t, 2, u.UniqueAddressUnits)
			}
		})
	}
	base.Points = append(base.Points, EffectivePoint{ID: "b", Address: "D0", DataType: schema.DataTypeUint16, NativeWidth: 1})
	p := BuildAcquisitionPlan(AcquisitionInput{Sources: []AcquisitionSource{base}})
	require.Empty(t, p.Errors())
	require.Len(t, p.WireUnits(), 2)
	require.NotEqual(t, p.WireUnits()[0].Group.Area, p.WireUnits()[1].Group.Area)
}

func TestAcquisitionWireUnitsDeterministicAndFingerprintUnchanged(t *testing.T) {
	in := randomInput(acquisitionProtocolFatek, 0, 1, 9)
	p := BuildAcquisitionPlan(in)
	before := p.Fingerprint()
	units := p.WireUnits()
	require.Equal(t, before, p.Fingerprint())
	for i, j := 0, len(in.Sources[0].Points)-1; i < j; i, j = i+1, j-1 {
		in.Sources[0].Points[i], in.Sources[0].Points[j] = in.Sources[0].Points[j], in.Sources[0].Points[i]
	}
	other := BuildAcquisitionPlan(in)
	require.Equal(t, before, other.Fingerprint())
	require.Equal(t, units, other.WireUnits())
}

func TestAcquisitionWireUnitsNativeBitAndWordAreas(t *testing.T) {
	for _, tc := range []struct {
		protocol, address, function string
		dtype                       schema.DataType
		width, bits, units          int
	}{
		{acquisitionProtocolModbus, "40001", "", schema.DataTypeUint32, 2, 16, 2},
		{acquisitionProtocolModbus, "0", "01", schema.DataTypeBool, 1, 1, 1},
		{"mc", "D0", "", schema.DataTypeBool, 1, 16, 1},
		{"mc", "M0", "", schema.DataTypeBool, 1, 1, 1},
	} {
		t.Run(tc.protocol+tc.address, func(t *testing.T) {
			in := randomInput(tc.protocol)
			in.Sources[0].Points = []EffectivePoint{{ID: "a", Address: tc.address, Function: tc.function, DataType: tc.dtype, NativeWidth: tc.width}}
			p := BuildAcquisitionPlan(in)
			require.Empty(t, p.Errors())
			require.Len(t, p.WireUnits(), 1)
			require.Equal(t, tc.bits, p.WireUnits()[0].AddressUnitBits)
			require.Equal(t, tc.units, p.WireUnits()[0].UniqueAddressUnits)
		})
	}
	require.Empty(t, BuildAcquisitionPlan(AcquisitionInput{}).WireUnits())
}
