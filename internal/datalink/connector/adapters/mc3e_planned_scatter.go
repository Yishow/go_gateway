package adapters

import (
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"
	"go-gateway/lib/hsllogic"
)

func mcPlannedResults(req connector.PlannedReadRequest, r optimizer.AcquisitionRequest, bindings []connector.PlannedPointBinding, words []int, bits []bool, completed time.Time, readErr error) []connector.PlannedPointResult {
	results := make([]connector.PlannedPointResult, len(r.Points))
	for i, p := range r.Points {
		b := connector.PlannedPointBinding{Point: p.Point}
		if i < len(bindings) {
			b = bindings[i]
		}
		result := connector.PlannedPointResult{Binding: b, Group: r.Group, PlanFingerprint: req.Plan.Fingerprint(), AcquisitionID: req.AcquisitionID, RequestID: req.RequestID, ReadResult: connector.ReadResult{Timestamp: completed, Quality: schema.QualityBad}, ObservedAt: completed, ReceivedAt: completed}
		err := readErr
		if err == nil {
			result.Value, result.RawBytes, err = decodeMCPlannedPoint(p, b.Profile, words, bits)
		}
		if err != nil {
			result.Value = nil
			result.RawBytes = nil
			result.Error = err.Error()
		} else {
			result.Quality = schema.QualityGood
		}
		results[i] = result
	}
	return results
}
func decodeMCPlannedPoint(p optimizer.PointSlice, profile connector.ReadProfile, words []int, bits []bool) (any, []byte, error) {
	if p.Offset < 0 || p.Length < 1 {
		return nil, nil, decoding.ErrShortScalar
	}
	if p.NativeBits == 1 {
		if p.Length != 1 || p.Offset >= len(bits) || p.Point.DataType != schema.DataTypeBool {
			return nil, nil, decoding.ErrShortScalar
		}
		raw := byte(0)
		if bits[p.Offset] {
			raw = 1
		}
		return bits[p.Offset], []byte{raw}, nil
	}
	width := schema.RegisterCountForDataType(p.Point.DataType)
	if p.NativeBits != 16 || p.Length != width || p.Offset > len(words)-p.Length {
		return nil, nil, decoding.ErrShortScalar
	}
	scalar := words[p.Offset : p.Offset+p.Length]
	registers := make([]uint16, len(scalar))
	for i, v := range scalar {
		registers[i] = uint16(v)
	} // #nosec G115 -- raw MC components are uint16.
	read := connector.ReadRequest{Address: p.Point.Address, Function: p.Point.Function, DataType: p.Point.DataType, Profile: profile, Count: 1}
	var value any
	var err error
	if modernMCScalar(read) {
		value, err = decodeMCReadWords(registers, read, 1)
	} else {
		// Frozen legacy policy, including unchanged register16 semantics. Never
		// read or mutate the shared connection converter while scattering points.
		converter := hsllogic.NewDataConverter(hsllogic.DataFormat(profile.Resolution().EffectiveFormat))
		value = converter.RegistersToValue(registers, hsllogic.DataType(p.Point.DataType))
	}
	return value, intSliceToBytes(scalar), err
}
