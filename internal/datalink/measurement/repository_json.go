package measurement

import (
	"encoding/json"
	"fmt"
)

type measurementJSONFields struct {
	counterPolicy, stateMap, bitmaskLabels []byte
}

func encodeMeasurementDefinition(def *MeasurementDefinition) (measurementJSONFields, error) {
	var encoded measurementJSONFields
	for _, field := range []struct {
		name   string
		value  any
		target *[]byte
	}{
		{"counterPolicy", def.CounterPolicy, &encoded.counterPolicy},
		{"stateMap", def.StateMap, &encoded.stateMap},
		{"bitmaskLabels", def.BitmaskLabels, &encoded.bitmaskLabels},
	} {
		value, err := json.Marshal(field.value)
		if err != nil {
			return measurementJSONFields{}, fmt.Errorf("encode measurement %s: %w", field.name, err)
		}
		*field.target = value
	}
	return encoded, nil
}
