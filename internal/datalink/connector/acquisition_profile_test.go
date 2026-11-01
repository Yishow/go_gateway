package connector

import (
	"fmt"
	"testing"
	"time"

	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestAcquisitionConfigUnknownDoesNotInferFromProtocolOrProbe(t *testing.T) {
	for _, tc := range []struct {
		protocol               schema.ProtocolType
		family, frame, station string
	}{
		{schema.ProtocolMC3E, "mc", "3e-binary", "0"}, {schema.ProtocolFatekFBs, "fatek", "ascii", "1"},
		{schema.ProtocolModbusTCP, "modbus", "tcp", "1"}, {schema.ProtocolModbusUDP, "modbus", "tcp", "1"}, {schema.ProtocolModbusRTU, "modbus", "rtu", "1"},
	} {
		t.Run(string(tc.protocol), func(t *testing.T) {
			cfg, err := ResolveAcquisitionConfig(tc.protocol, `{"probe_success":true}`)
			require.NoError(t, err)
			require.False(t, cfg.Declared())
			require.Empty(t, cfg.ProfileID())
			require.Equal(t, "unknown", cfg.Capability(16).RandomSupport)
			require.Equal(t, 5*time.Second, cfg.ReadTimeout())
			identity := cfg.Identity("device")
			require.Equal(t, tc.family, identity.Protocol)
			require.Equal(t, tc.frame, identity.Frame)
			require.Equal(t, tc.station, identity.Station)
			require.Equal(t, ReadConnectionRevision(`{"probe_success":true}`), identity.ConnectionRevision)
			require.Empty(t, cfg.Coverage().Readable)
		})
	}
}

func TestAcquisitionConfigProfilesHaveCommandSpecificLimits(t *testing.T) {
	cfg, err := ResolveAcquisitionConfig(schema.ProtocolMC3E, `{"acquisition":{"version":1,"profile":"mitsubishi-fx5-cpu-3e-binary-v1"}}`)
	require.NoError(t, err)
	require.Equal(t, 960, cfg.Capability(16).ProtocolLimit)
	require.Equal(t, 192, cfg.Capability(16).RandomProtocolLimit)
	require.Equal(t, 255, cfg.Capability(16).ClientLimit)
	require.Equal(t, 3584, cfg.Capability(1).ProtocolLimit)
	require.Zero(t, cfg.Capability(1).RandomProtocolLimit)
	require.Equal(t, 2000, cfg.Capability(1).ClientLimit)
	require.Equal(t, "unsupported", cfg.Capability(1).RandomSupport)
	for name, points := range map[string][]optimizer.EffectivePoint{
		"193 sparse words":     acquisitionConfigPoints(193, 10, "D", schema.DataTypeUint16),
		"126 contiguous words": acquisitionConfigPoints(126, 1, "D", schema.DataTypeUint16),
		"193 contiguous bits":  acquisitionConfigPoints(193, 1, "M", schema.DataTypeBool),
	} {
		t.Run(name, func(t *testing.T) {
			bits := 16
			if points[0].DataType == schema.DataTypeBool {
				bits = 1
			}
			plan := acquisitionConfigPlan(cfg, bits, points)
			require.Empty(t, plan.Errors())
			switch {
			case bits == 1:
				require.Len(t, plan.Requests(), 1)
				require.Equal(t, 193, plan.Requests()[0].Length)
			case name == "193 sparse words":
				require.Len(t, plan.Requests(), 2)
				require.Equal(t, 192, plan.Requests()[0].Length)
			default:
				// The optimizer may choose one random command over two contiguous ones.
				require.Len(t, plan.Requests(), 1)
				require.Equal(t, "random", plan.Requests()[0].Kind)
				capability := cfg.Capability(bits)
				capability.RandomSupport = "unsupported"
				fallback := optimizer.BuildAcquisitionPlan(optimizer.AcquisitionInput{Sources: []optimizer.AcquisitionSource{{Identity: cfg.Identity("device"), Capability: capability, Points: points}}})
				require.Empty(t, fallback.Errors())
				require.Len(t, fallback.Requests(), 2)
				require.Equal(t, 125, fallback.Requests()[0].Length)
			}
		})
	}
	fatek, err := ResolveAcquisitionConfig(schema.ProtocolFatekFBs, `{"acquisition":{"version":1,"profile":"fatek-fbs-ascii-v1"}}`)
	require.NoError(t, err)
	require.Equal(t, "supported", fatek.Capability(1).RandomSupport)
	require.Equal(t, 255, fatek.Capability(1).ClientLimit)
	for _, tc := range []struct {
		name, area          string
		count, stride, bits int
		dtype               schema.DataType
		requests            int
	}{
		{"mixed native limit", "D", 65, 10, 16, schema.DataTypeUint16, 2},
		{"32-bit weighted word limit", "DR", 33, 10, 32, schema.DataTypeUint32, 2},
		{"contiguous bits not word cap", "X", 200, 1, 1, schema.DataTypeBool, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			plan := acquisitionConfigPlan(fatek, tc.bits, acquisitionConfigPoints(tc.count, tc.stride, tc.area, tc.dtype))
			require.Empty(t, plan.Errors())
			require.Len(t, plan.Requests(), tc.requests)
		})
	}
}

func TestAcquisitionConfigLimitsCannotLoosenClientAndProfile(t *testing.T) {
	cfg, err := ResolveAcquisitionConfig(schema.ProtocolMC3E, `{"timeout":2,"acquisition":{"version":1,"profile":"mitsubishi-fx5-cpu-3e-binary-v1","limits":{"native_units":999,"request_bytes":100000,"response_bytes":100000}}}`)
	require.NoError(t, err)
	require.Equal(t, 2*time.Second, cfg.ReadTimeout())
	plan := acquisitionConfigPlan(cfg, 16, acquisitionConfigPoints(193, 10, "D", schema.DataTypeUint16))
	require.Empty(t, plan.Errors())
	require.Len(t, plan.Requests(), 2)
	require.Equal(t, 192, plan.Requests()[0].Length)
	tight, err := ResolveAcquisitionConfig(schema.ProtocolMC3E, `{"acquisition":{"version":1,"profile":"mitsubishi-fx5-cpu-3e-binary-v1","limits":{"native_units":2}}}`)
	require.NoError(t, err)
	require.Len(t, acquisitionConfigPlan(tight, 16, acquisitionConfigPoints(3, 10, "D", schema.DataTypeUint16)).Requests(), 2)
	safe, err := ResolveAcquisitionConfig(schema.ProtocolMC3E, `{"acquisition":{"version":1,"profile":"safe-v1"}}`)
	require.NoError(t, err)
	require.Equal(t, "unsupported", safe.Capability(16).RandomSupport)
	require.Len(t, acquisitionConfigPlan(safe, 16, acquisitionConfigPoints(3, 10, "D", schema.DataTypeUint16)).Requests(), 3)
}

func TestAcquisitionConfigCoverageDetachedAndContentBound(t *testing.T) {
	raw := `{"acquisition":{"version":1,"profile":"safe-v1","coverage":{"revision":"coverage-1","readable":[{"area":"HR","start":0,"end":21}],"excluded":[],"forbidden":[]}}}`
	cfg, err := ResolveAcquisitionConfig(schema.ProtocolModbusTCP, raw)
	require.NoError(t, err)
	points := acquisitionConfigPoints(3, 10, "HR", schema.DataTypeUint16)
	first := acquisitionConfigPlan(cfg, 16, points)
	require.Len(t, first.Requests(), 1)
	require.Equal(t, 21, first.Requests()[0].Length)
	detached := cfg.Coverage()
	detached.Readable[0].End = 1
	require.Equal(t, 21, cfg.Coverage().Readable[0].End)
	changed, err := ResolveAcquisitionConfig(schema.ProtocolModbusTCP, `{"acquisition":{"version":1,"profile":"safe-v1","coverage":{"revision":"coverage-1","readable":[{"area":"HR","start":0,"end":20}]}}}`)
	require.NoError(t, err)
	require.NotEqual(t, cfg.Revision(), changed.Revision())
	require.NotEqual(t, first.Fingerprint(), acquisitionConfigPlan(changed, 16, points).Fingerprint())
	denied, err := ResolveAcquisitionConfig(schema.ProtocolModbusTCP, `{"acquisition":{"version":1,"profile":"safe-v1","coverage":{"revision":"coverage-1","readable":[{"area":"HR","start":0,"end":21}],"forbidden":[{"area":"HR","start":5,"end":6}]}}}`)
	require.NoError(t, err)
	require.Greater(t, len(acquisitionConfigPlan(denied, 16, points).Requests()), 1)
}

func acquisitionConfigPoints(count, stride int, area string, dtype schema.DataType) []optimizer.EffectivePoint {
	points := make([]optimizer.EffectivePoint, count)
	for i := range points {
		points[i] = optimizer.EffectivePoint{ID: fmt.Sprint(i), Address: fmt.Sprintf("%s%d", area, i*stride), DataType: dtype, NativeWidth: 1}
	}
	return points
}
func acquisitionConfigPlan(cfg AcquisitionConfig, bits int, points []optimizer.EffectivePoint) optimizer.AcquisitionPlan {
	return optimizer.BuildAcquisitionPlan(optimizer.AcquisitionInput{ConfigurationRevision: cfg.Revision(), Sources: []optimizer.AcquisitionSource{{Identity: cfg.Identity("device"), Capability: cfg.Capability(bits), Coverage: cfg.Coverage(), Points: points}}})
}
