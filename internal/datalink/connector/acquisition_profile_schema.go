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
	positive := func(description string) map[string]any {
		return map[string]any{"type": "integer", "minimum": 1, "description": description}
	}
	nativeRange := map[string]any{"type": "object", "additionalProperties": false, "required": []string{"area", "start", "end"}, "properties": map[string]any{
		"area":  map[string]any{"type": "string", "minLength": 1, "description": "Canonical protocol area; function namespaces stay separate."},
		"start": map[string]any{"type": "integer", "minimum": 0},
		"end":   map[string]any{"type": "integer", "minimum": 1, "description": "Exclusive physical address, not component count."},
	}}
	ranges := func() map[string]any { return map[string]any{"type": "array", "items": nativeRange} }
	declaration := map[string]any{"type": "object", "title": "Acquisition declaration", "additionalProperties": false, "required": []string{"version", "profile"}, "properties": map[string]any{
		"version": map[string]any{"type": "integer", "enum": []int{1}},
		"profile": map[string]any{"type": "string", "enum": profiles, "description": "Explicit supported family profile; omitted acquisition remains unknown. FX5 CPU profile is only local CPU built-in Ethernet, not routed or module endpoints."},
		"limits": map[string]any{"type": "object", "additionalProperties": false, "properties": map[string]any{
			"native_units":   positive("Optional stricter per-request native component count; protocol/client limits still apply."),
			"request_bytes":  positive("Optional stricter complete application-request frame bytes."),
			"response_bytes": positive("Optional stricter complete application-response frame bytes."),
		}},
		"coverage": map[string]any{"type": "object", "additionalProperties": false, "required": []string{"revision"}, "properties": map[string]any{
			"revision": map[string]any{"type": "string", "minLength": 1, "pattern": `\S`},
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
