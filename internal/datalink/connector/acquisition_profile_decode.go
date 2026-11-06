package connector

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"go-gateway/internal/datalink/optimizer"
)

type acquisitionDeclaration struct {
	Version  int                  `json:"version"`
	Profile  string               `json:"profile"`
	Limits   *acquisitionLimits   `json:"limits,omitempty"`
	Coverage *acquisitionCoverage `json:"coverage,omitempty"`
}
type acquisitionLimits struct {
	NativeUnits   *int `json:"native_units,omitempty"`
	RequestBytes  *int `json:"request_bytes,omitempty"`
	ResponseBytes *int `json:"response_bytes,omitempty"`
	RandomPoints  *int `json:"random_points,omitempty"`
}
type acquisitionCoverage struct {
	Revision  string             `json:"revision"`
	Readable  []acquisitionRange `json:"readable,omitempty"`
	Forbidden []acquisitionRange `json:"forbidden,omitempty"`
	Excluded  []acquisitionRange `json:"excluded,omitempty"`
}
type acquisitionRange struct {
	Area  string `json:"area"`
	Start *int   `json:"start"`
	End   *int   `json:"end"`
}

// Preserve legacy top-level properties while rejecting an ambiguous repeated
// acquisition declaration in an exact installed configuration snapshot.
func acquisitionConfigFields(raw []byte) (map[string]json.RawMessage, error) {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	token, err := decoder.Token()
	if err != nil || token != json.Delim('{') {
		return nil, ErrAcquisitionConfig
	}
	fields := map[string]json.RawMessage{}
	for decoder.More() {
		token, err = decoder.Token()
		key, ok := token.(string)
		if err != nil || !ok {
			return nil, ErrAcquisitionConfig
		}
		if _, exists := fields[key]; key == "acquisition" && exists {
			return nil, ErrAcquisitionConfig
		}
		var value json.RawMessage
		if decoder.Decode(&value) != nil {
			return nil, ErrAcquisitionConfig
		}
		fields[key] = value
	}
	if token, err = decoder.Token(); err != nil || token != json.Delim('}') {
		return nil, ErrAcquisitionConfig
	}
	if _, err = decoder.Token(); err != io.EOF {
		return nil, ErrAcquisitionConfig
	}
	return fields, nil
}

func decodeAcquisitionDeclaration(raw []byte) (acquisitionDeclaration, error) {
	var result acquisitionDeclaration
	decoder := json.NewDecoder(bytes.NewReader(raw))
	if validateAcquisitionTokens(decoder) != nil {
		return result, ErrAcquisitionConfig
	}
	if _, err := decoder.Token(); err != io.EOF {
		return result, ErrAcquisitionConfig
	}
	fields, ok := exactAcquisitionObject(raw, "version", "profile", "limits", "coverage")
	if !ok {
		return result, ErrAcquisitionConfig
	}
	if nested, exists := fields["limits"]; exists {
		if _, ok := exactAcquisitionObject(nested, "native_units", "request_bytes", "response_bytes", "random_points"); !ok {
			return result, ErrAcquisitionConfig
		}
	}
	if nested, exists := fields["coverage"]; exists {
		coverage, ok := exactAcquisitionObject(nested, "revision", "readable", "forbidden", "excluded")
		if !ok {
			return result, ErrAcquisitionConfig
		}
		for _, key := range []string{"readable", "forbidden", "excluded"} {
			if ranges, exists := coverage[key]; exists {
				var items []json.RawMessage
				if json.Unmarshal(ranges, &items) != nil {
					return result, ErrAcquisitionConfig
				}
				for _, item := range items {
					if _, ok := exactAcquisitionObject(item, "area", "start", "end"); !ok {
						return result, ErrAcquisitionConfig
					}
				}
			}
		}
	}
	if json.Unmarshal(raw, &result) != nil || result.Version != 1 || result.Profile == "" {
		return result, ErrAcquisitionConfig
	}
	if result.Limits != nil {
		for _, value := range []*int{result.Limits.NativeUnits, result.Limits.RequestBytes, result.Limits.ResponseBytes, result.Limits.RandomPoints} {
			if value != nil && *value <= 0 {
				return result, ErrAcquisitionConfig
			}
		}
	}
	if result.Coverage != nil {
		if strings.TrimSpace(result.Coverage.Revision) == "" {
			return result, ErrAcquisitionConfig
		}
		for _, ranges := range [][]acquisitionRange{result.Coverage.Readable, result.Coverage.Forbidden, result.Coverage.Excluded} {
			for _, r := range ranges {
				if r.Start == nil || r.End == nil {
					return result, ErrAcquisitionConfig
				}
			}
		}
	}
	return result, nil
}
func exactAcquisitionObject(raw []byte, allowed ...string) (map[string]json.RawMessage, bool) {
	var fields map[string]json.RawMessage
	if json.Unmarshal(raw, &fields) != nil || fields == nil {
		return nil, false
	}
	for key := range fields {
		found := false
		for _, candidate := range allowed {
			if key == candidate {
				found = true
				break
			}
		}
		if !found {
			return nil, false
		}
	}
	return fields, true
}

// Validate raw tokens before decoding: encoding/json otherwise accepts null for
// scalar values and silently replaces duplicate object keys with the last one.
func validateAcquisitionTokens(decoder *json.Decoder) error {
	token, err := decoder.Token()
	if err != nil || token == nil {
		return ErrAcquisitionConfig
	}
	delimiter, ok := token.(json.Delim)
	if !ok {
		return nil
	}
	switch delimiter {
	case '{':
		seen := map[string]bool{}
		for decoder.More() {
			token, err := decoder.Token()
			key, ok := token.(string)
			if err != nil || !ok || seen[key] {
				return ErrAcquisitionConfig
			}
			seen[key] = true
			if validateAcquisitionTokens(decoder) != nil {
				return ErrAcquisitionConfig
			}
		}
	case '[':
		for decoder.More() {
			if validateAcquisitionTokens(decoder) != nil {
				return ErrAcquisitionConfig
			}
		}
	default:
		return ErrAcquisitionConfig
	}
	_, err = decoder.Token()
	return err
}
func positiveAcquisitionLimit(value *int) int {
	if value == nil {
		return 0
	}
	return *value
}
func (d acquisitionDeclaration) coverage() optimizer.ReadableCoverage {
	if d.Coverage == nil {
		return optimizer.ReadableCoverage{}
	}
	convert := func(ranges []acquisitionRange) []optimizer.NativeRange {
		result := make([]optimizer.NativeRange, len(ranges))
		for i, r := range ranges {
			result[i] = optimizer.NativeRange{Area: r.Area, Start: *r.Start, End: *r.End}
		}
		return result
	}
	return optimizer.ReadableCoverage{Revision: d.Coverage.Revision, Readable: convert(d.Coverage.Readable), Forbidden: convert(d.Coverage.Forbidden), Excluded: convert(d.Coverage.Excluded)}
}
