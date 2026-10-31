package adapters

import (
	"fmt"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"
)

type fatekPlannedRaw struct {
	random map[string]any
	words  []int
	bits   []bool
}

func fatekPlannedResults(req connector.PlannedReadRequest, r optimizer.AcquisitionRequest, bindings []connector.PlannedPointBinding, raw fatekPlannedRaw, completed time.Time, readErr error) []connector.PlannedPointResult {
	results := make([]connector.PlannedPointResult, len(r.Points))
	for i, p := range r.Points {
		b := connector.PlannedPointBinding{Point: p.Point}
		if i < len(bindings) {
			b = bindings[i]
		}
		result := connector.PlannedPointResult{Binding: b, Group: r.Group, PlanFingerprint: req.Plan.Fingerprint(), AcquisitionID: req.AcquisitionID, RequestID: req.RequestID, ReadResult: connector.ReadResult{Timestamp: completed, TimeOrigin: "gateway", Quality: schema.QualityBad}, ObservedAt: completed, ReceivedAt: completed}
		err := readErr
		if err == nil {
			result.Value, result.RawBytes, err = decodeFatekPlannedPoint(r, p, b.Profile, raw)
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
func decodeFatekPlannedPoint(r optimizer.AcquisitionRequest, p optimizer.PointSlice, profile connector.ReadProfile, raw fatekPlannedRaw) (any, []byte, error) {
	if err := validateFatekPlannedSlice(r, p); err != nil {
		return nil, nil, err
	}
	if p.Offset < 0 || p.Length < 1 {
		return nil, nil, decoding.ErrShortScalar
	}
	if p.NativeBits == 1 {
		if p.Length != 1 || p.Point.DataType != schema.DataTypeBool {
			return nil, nil, decoding.ErrShortScalar
		}
		var value bool
		if r.Kind == "random" {
			item, err := fatekPlannedItem(r, p.Offset, raw.random)
			if err != nil {
				return nil, nil, err
			}
			var ok bool
			value, ok = item.(bool)
			if !ok {
				return nil, nil, fmt.Errorf("invalid FATEK bit component")
			}
		} else {
			if p.Offset >= len(raw.bits) {
				return nil, nil, decoding.ErrShortScalar
			}
			value = raw.bits[p.Offset]
		}
		b := byte(0)
		if value {
			b = 1
		}
		return value, []byte{b}, nil
	}
	if p.NativeBits != 16 && p.NativeBits != 32 {
		return nil, nil, decoding.ErrShortScalar
	}
	var values []int
	if r.Kind == "random" {
		if p.Offset > len(r.Units)-p.Length {
			return nil, nil, decoding.ErrShortScalar
		}
		values = make([]int, p.Length)
		for i := range values {
			item, err := fatekPlannedItem(r, p.Offset+i, raw.random)
			if err != nil {
				return nil, nil, err
			}
			value, ok := item.(int)
			if !ok || value < 0 || uint64(value) > (uint64(1)<<uint(p.NativeBits))-1 {
				return nil, nil, fmt.Errorf("invalid FATEK register component")
			}
			values[i] = value
		}
	} else {
		if p.Offset > len(raw.words)-p.Length {
			return nil, nil, decoding.ErrShortScalar
		}
		values = raw.words[p.Offset : p.Offset+p.Length]
	}
	read := connector.ReadRequest{Address: p.Point.Address, Function: p.Point.Function, DataType: p.Point.DataType, Profile: profile, Count: p.Length}
	value, err := decodeFatekReadComponents(values, p.NativeBits, read)
	return value, intSliceToBytes(values), err
}
func fatekPlannedItem(r optimizer.AcquisitionRequest, index int, values map[string]any) (any, error) {
	if index < 0 || index >= len(r.Units) {
		return nil, decoding.ErrShortScalar
	}
	u := r.Units[index]
	key := fmt.Sprint(u.Area, u.Address)
	value, ok := values[key]
	if !ok {
		return nil, fmt.Errorf("missing FATEK random component %s", key)
	}
	return value, nil
}
