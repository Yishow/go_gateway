package adapters

import (
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"
	"time"
)

func modbusPlannedResults(req connector.PlannedReadRequest, r optimizer.AcquisitionRequest, bindings []connector.PlannedPointBinding, words []uint16, bits []bool, completed time.Time, readErr error) []connector.PlannedPointResult {
	results := make([]connector.PlannedPointResult, len(r.Points))
	for i, p := range r.Points {
		b := connector.PlannedPointBinding{Point: p.Point}
		if i < len(bindings) {
			b = bindings[i]
		}
		result := connector.PlannedPointResult{Binding: b, Group: r.Group, PlanFingerprint: req.Plan.Fingerprint(), AcquisitionID: req.AcquisitionID, RequestID: req.RequestID, ReadResult: connector.ReadResult{Timestamp: completed, TimeOrigin: "gateway", Quality: schema.QualityBad}, ObservedAt: completed, ReceivedAt: completed}
		err := readErr
		if err == nil {
			result.Value, result.RawBytes, err = decodeModbusPlannedPoint(p, b.Profile, words, bits)
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

func decodeModbusPlannedPoint(p optimizer.PointSlice, profile connector.ReadProfile, words []uint16, bits []bool) (any, []byte, error) {
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
	if p.NativeBits != 16 || p.Length != schema.RegisterCountForDataType(p.Point.DataType) || p.Offset > len(words)-p.Length {
		return nil, nil, decoding.ErrShortScalar
	}
	scalar := words[p.Offset : p.Offset+p.Length]
	// The profile is frozen even for legacy/native16 values. Never consult the
	// connection's current DataFormat while scattering a completed old request.
	read := connector.ReadRequest{Address: p.Point.Address, Function: p.Point.Function, DataType: p.Point.DataType, Profile: profile, DataFormat: profile.Resolution().EffectiveFormat, Count: p.Length}
	value, err := decodeModbusReadWords(scalar, read, "")
	return value, uint16SliceToBytes(scalar), err
}
