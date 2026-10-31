package optimizer

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
)

func TestFatekNativeFootprintBoundsAndForbidden(t *testing.T) {
	for _, tc := range []struct {
		address, area     string
		dtype             schema.DataType
		width, start, end int
		forbidden         []NativeRange
		code              string
	}{
		{"DR5", "DR", schema.DataTypeFloat32, 1, 5, 7, nil, ""},
		{"DR5", "DR", schema.DataTypeFloat32, 1, 5, 6, nil, "address-out-of-bounds"},
		{"WY8", "WY", schema.DataTypeUint32, 2, 8, 40, nil, ""},
		{"WY8", "WY", schema.DataTypeUint32, 2, 8, 39, nil, "address-out-of-bounds"},
		{"WY8", "WY", schema.DataTypeUint32, 2, 8, 40, []NativeRange{{Area: "WY", Start: 31, End: 32}}, "forbidden-address"},
		{"DWM8", "DWM", schema.DataTypeUint64, 2, 8, 72, nil, ""},
		{"DWM8", "DWM", schema.DataTypeUint64, 2, 8, 71, nil, "address-out-of-bounds"},
		{"DR65532", "DR", schema.DataTypeUint64, 2, 0, 65536, nil, ""},
		{"DR65533", "DR", schema.DataTypeUint64, 2, 0, 65536, nil, "invalid-address"},
	} {
		t.Run(tc.address+tc.code, func(t *testing.T) {
			in := randomInput("fatek")
			s := &in.Sources[0]
			s.Points = []EffectivePoint{{ID: "p", Address: tc.address, DataType: tc.dtype, NativeWidth: tc.width}}
			s.Capability.DeviceAddressBounds = []NativeAddressBound{{Area: tc.area, StartInclusive: tc.start, EndExclusive: tc.end}}
			s.Coverage = ReadableCoverage{Revision: "coverage", Forbidden: tc.forbidden}
			plan := BuildAcquisitionPlan(in)
			if tc.code != "" {
				require.Len(t, plan.Errors(), 1)
				require.Equal(t, tc.code, plan.Errors()[0].Code)
				require.Empty(t, plan.Requests())
			} else {
				require.Empty(t, plan.Errors())
				require.Len(t, plan.Requests(), 1)
			}
		})
	}
}
func TestFatekNativeUnknownCapabilityKeepsPhaseAndReadableGap(t *testing.T) {
	in := randomInput("fatek")
	s := &in.Sources[0]
	s.Capability.RandomSupport = "unknown"
	s.Points = []EffectivePoint{{ID: "a", Address: "DR0", DataType: schema.DataTypeUint32, NativeWidth: 1}, {ID: "b", Address: "DR1", DataType: schema.DataTypeUint32, NativeWidth: 1}}
	plan := BuildAcquisitionPlan(in)
	require.Empty(t, plan.Errors())
	require.Len(t, plan.Requests(), 2)
	for _, r := range plan.Requests() {
		require.Equal(t, 2, r.AddressStep)
		require.Equal(t, 1, r.Length)
	}
	s.Points[1].Address = "DR4"
	plan = BuildAcquisitionPlan(in)
	require.Len(t, plan.Requests(), 2)
	s.Coverage = ReadableCoverage{Revision: "coverage", Readable: []NativeRange{{Area: "DR", Start: 2, End: 4}}}
	plan = BuildAcquisitionPlan(in)
	require.Empty(t, plan.Errors())
	require.Len(t, plan.Requests(), 1)
	require.Equal(t, 3, plan.Requests()[0].Length)
	require.Equal(t, 2, plan.Requests()[0].Points[1].Offset)
	s.Coverage.Forbidden = []NativeRange{{Area: "DR", Start: 3, End: 4}}
	plan = BuildAcquisitionPlan(in)
	require.Empty(t, plan.Errors())
	require.Len(t, plan.Requests(), 2)
}
func TestFatekNativeStepSnapshotIsDetached(t *testing.T) {
	in := randomInput("fatek")
	in.Sources[0].Points = []EffectivePoint{{ID: "p", Address: "DR0", DataType: schema.DataTypeUint64, NativeWidth: 2}}
	plan := BuildAcquisitionPlan(in)
	request, ok := plan.Request(0)
	require.True(t, ok)
	require.Equal(t, 2, request.AddressStep)
	request.AddressStep = 99
	again, ok := plan.Request(0)
	require.True(t, ok)
	require.Equal(t, 2, again.AddressStep)
	requests := plan.Requests()
	requests[0].AddressStep = 0
	require.Equal(t, 2, plan.Requests()[0].AddressStep)
}
