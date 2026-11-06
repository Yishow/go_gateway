package connector

import (
	"encoding/json"

	"go-gateway/internal/datalink/schema"
)

// AcquisitionConfigSchema describes the same optional declaration resolved by
// ResolveAcquisitionConfig. It never advertises random support for an undeclared
// device, and a profile selection remains explicit operator configuration.
func AcquisitionConfigSchema(protocol schema.ProtocolType) json.RawMessage {
	profiles := acquisitionProfiles(protocol)
	if len(profiles) == 0 {
		return nil
	}
	const (
		additionalPropertiesKey = "additionalProperties"
		descriptionKey          = "description"
		integerType             = "integer"
		minimumKey              = "minimum"
		objectType              = "object"
		propertiesKey           = "properties"
		requiredKey             = "required"
		stringType              = "string"
		typeKey                 = "type"
	)
	positive := func(description string) map[string]any {
		return map[string]any{typeKey: integerType, minimumKey: 1, descriptionKey: description}
	}
	nativeRange := map[string]any{typeKey: objectType, additionalPropertiesKey: false, requiredKey: []string{"area", "start", "end"}, propertiesKey: map[string]any{
		"area":  map[string]any{typeKey: stringType, "minLength": 1, descriptionKey: "Canonical protocol area; function namespaces stay separate."},
		"start": map[string]any{typeKey: integerType, minimumKey: 0},
		"end":   map[string]any{typeKey: integerType, minimumKey: 1, descriptionKey: "Exclusive physical address, not component count."},
	}}
	ranges := func() map[string]any { return map[string]any{"type": "array", "items": nativeRange} }
	declaration := map[string]any{typeKey: objectType, "title": "Acquisition declaration", additionalPropertiesKey: false, requiredKey: []string{"version", "profile"}, propertiesKey: map[string]any{
		"version": map[string]any{typeKey: integerType, "enum": []int{1}},
		"profile": map[string]any{typeKey: stringType, "enum": profiles, descriptionKey: "Explicit supported family profile; omitted acquisition remains unknown. FX5 CPU profile is only local CPU built-in Ethernet, not routed or module endpoints."},
		"limits": map[string]any{typeKey: objectType, additionalPropertiesKey: false, propertiesKey: map[string]any{
			"native_units":   positive("Optional stricter per-request native component count; protocol/client limits still apply."),
			"request_bytes":  positive("Optional stricter complete application-request frame bytes."),
			"response_bytes": positive("Optional stricter complete application-response frame bytes."),
			"random_points":  positive("Optional lower per-request random point limit; must not exceed the profile limit (MC 192, FATEK 64)."),
		}},
		"coverage": map[string]any{typeKey: objectType, additionalPropertiesKey: false, requiredKey: []string{"revision"}, propertiesKey: map[string]any{
			"revision": map[string]any{typeKey: stringType, "minLength": 1, "pattern": `\S`},
			"readable": ranges(), "forbidden": ranges(), "excluded": ranges(),
		}},
	}}
	raw, err := json.Marshal(declaration)
	if err != nil {
		return nil
	}
	return raw
}

func withAcquisitionConfigSchema(info ProtocolInfo) ProtocolInfo {
	declaration := AcquisitionConfigSchema(info.Type)
	if len(declaration) == 0 {
		return info
	}
	var config map[string]json.RawMessage
	var properties map[string]json.RawMessage
	if json.Unmarshal(info.ConfigSchema, &config) != nil || json.Unmarshal(config["properties"], &properties) != nil {
		return info
	}
	properties["acquisition"] = declaration
	rawProperties, err := json.Marshal(properties)
	if err != nil {
		return info
	}
	config["properties"] = rawProperties
	rawConfig, err := json.Marshal(config)
	if err != nil {
		return info
	}
	info.ConfigSchema = rawConfig
	return info
}
