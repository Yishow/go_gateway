package handlers

import (
	"encoding/json"
	"go-gateway/internal/datalink/schema"
	"regexp"
	"strings"
)

var runtimePointFunctionPattern = regexp.MustCompile(`^(?:|[A-Za-z_][A-Za-z0-9_]{0,31}|0[1-9]|1[0-6])$`)

type runtimePointReadResponse struct {
	*schema.Point
	ValueState    string `json:"value_state,omitempty"`
	MetadataState string `json:"metadata_state,omitempty"`
}

// Only read responses are projected. The repository retains the original diagnostic/value.
func safeRuntimePointRead(record *schema.Point) runtimePointReadResponse {
	if record == nil {
		return runtimePointReadResponse{}
	}
	copyRecord := *record
	copyRecord.LastError = runtimeSafeDiagnostic(record.LastError)
	response := runtimePointReadResponse{Point: &copyRecord}
	if !safeRuntimePointMetadata(&copyRecord) {
		response.MetadataState = "unsupported"
	}
	if record.LastValue == nil {
		return response
	}
	raw := *record.LastValue
	if record.DataType.IsValid() && safeRuntimeStoredScalar(raw) {
		return response
	}
	copyRecord.LastValue = nil
	response.ValueState = "unsupported"
	return response
}

func safeRuntimeStoredScalar(raw string) bool {
	var value any
	decoder := json.NewDecoder(strings.NewReader(raw))
	decoder.UseNumber()
	if json.Valid([]byte(raw)) && decoder.Decode(&value) == nil {
		if _, ok := value.(json.Number); ok {
			return true
		}
		_, ok := runtimeSafeSensorValue(value)
		return ok
	}
	// Legacy plain text is valid sensor input, not an arbitrary JSON container.
	if strings.HasPrefix(strings.TrimSpace(raw), "{") || strings.HasPrefix(strings.TrimSpace(raw), "[") {
		return false
	}
	_, ok := runtimeSafeSensorValue(raw)
	return ok
}

func safeRuntimePointMetadata(p *schema.Point) bool {
	supported := true
	project := func(value *string, safe string) {
		if *value != safe {
			supported = false
			*value = safe
		}
	}
	project(&p.ID, runtimeSafeID(p.ID))
	project(&p.DeviceID, runtimeSafeID(p.DeviceID))
	for _, value := range []*string{&p.Name, &p.Description} {
		if !runtimeSafeSensorText(*value) {
			project(value, "")
		}
	}
	project(&p.Address, runtimeSafeAddress(p.Address))
	if !runtimePointFunctionPattern.MatchString(p.Function) {
		project(&p.Function, "")
	}
	project(&p.DataFormat, runtimeSafeEnum(p.DataFormat, "", "ABCD", "BADC", "CDAB", "DCBA"))
	mode := schema.PointMode(runtimeSafeEnum(string(p.Mode), "", "read", "readwrite"))
	if mode != p.Mode {
		supported = false
		p.Mode = mode
	}
	dataType := runtimeSafeDataType(p.DataType)
	if dataType != p.DataType {
		supported = false
		p.DataType = dataType
	}
	if p.PollingGroupID != nil && runtimeSafeID(*p.PollingGroupID) != *p.PollingGroupID {
		supported = false
		p.PollingGroupID = nil
	}
	return supported
}
