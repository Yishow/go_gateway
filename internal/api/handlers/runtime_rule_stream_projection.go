package handlers

import (
	"encoding/json"
	datalinkruntime "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
	"math"
)

// The original ValueEvent serializer preserves integer precision on the wire.
// Complex values are unsupported, not recursively stringified or partially displayed.
func safeRuntimeValueEvent(event datalinkruntime.ValueEvent) (json.RawMessage, error) {
	raw, rawOK := runtimeSafeSensorValue(event.RawValue)
	transformed, transformedOK := runtimeSafeSensorValue(event.TransformedValue)
	event.DeviceID = runtimeSafeID(event.DeviceID)
	event.PointID = runtimeSafeID(event.PointID)
	address := runtimeSafeAddress(event.Address)
	supported := rawOK && transformedOK && address != "" && address == event.Address && event.DeviceID != "" && event.PointID != ""
	event.Quality = schema.QualityFlag(runtimeSafeEnum(string(event.Quality), "good", "bad", "missing", "stale", "invalid", "uncertain"))
	event.Address = address
	event.RawValue = raw
	event.TransformedValue = transformed
	if !supported {
		event.Quality = schema.QualityInvalid
	}
	data, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}
	if supported {
		return data, nil
	}
	return append(data[:len(data)-1], []byte(`,"value_state":"unsupported"}`)...), nil
}

func runtimeSafeSensorValue(value any) (any, bool) {
	switch v := value.(type) {
	case nil, bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return v, true
	case float32:
		if !math.IsNaN(float64(v)) && !math.IsInf(float64(v), 0) {
			return v, true
		}
	case float64:
		if !math.IsNaN(v) && !math.IsInf(v, 0) {
			return v, true
		}
	case string:
		if runtimeSafeSensorText(v) {
			return v, true
		}
	}
	return nil, false
}
