package recordingplan

import (
	"encoding/json"
	"fmt"
)

type planJSONFields struct {
	members, streams, destinations, retention, limits []byte
}

func encodeRecordingPlan(plan *RecordingPlan) (planJSONFields, error) {
	var encoded planJSONFields
	for _, field := range []struct {
		name   string
		value  any
		target *[]byte
	}{
		{"members", plan.Members, &encoded.members},
		{"streams", plan.Streams, &encoded.streams},
		{"destinations", plan.Destinations, &encoded.destinations},
		{"retention", plan.Retention, &encoded.retention},
		{"limits", plan.Limits, &encoded.limits},
	} {
		value, err := json.Marshal(field.value)
		if err != nil {
			return planJSONFields{}, fmt.Errorf("encode recordingplan %s: %w", field.name, err)
		}
		*field.target = value
	}
	return encoded, nil
}
