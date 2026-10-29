package handlers

import (
	"encoding/json"
	"go-gateway/internal/datalink/schema"
)

// runtimeWorkspaceDatabaseConfig deliberately does not embed the authoring DTO.
// This is saved configuration, never proof of the currently effective destination.
type runtimeWorkspaceDatabaseConfig struct {
	ID                   string `json:"id"`
	WorkspaceID          string `json:"workspace_id"`
	Kind                 string `json:"kind"`
	Name                 string `json:"name"`
	Schema               string `json:"schema"`
	Table                string `json:"table"`
	WriteMode            string `json:"write_mode"`
	WriteIntervalSeconds int    `json:"write_interval_seconds"`
	TimestampColumn      string `json:"timestamp_column"`
	Status               string `json:"status"`
}

func buildRuntimeDatabaseConfig(workspaceID string, connector *schema.DatabaseConnector) runtimeWorkspaceDatabaseConfig {
	// Decode only named strings. Nested/unknown fields cannot enter this projection.
	var config struct {
		Schema          string `json:"schema"`
		Table           string `json:"table"`
		TimestampColumn string `json:"timestamp_column"`
		WriteMode       string `json:"write_mode"`
	}
	configErr := json.Unmarshal([]byte(connector.ConnectionConfig), &config)
	kind := "unavailable"
	switch connector.Kind {
	case "sqlite", "postgres", "mysql", "sqlserver":
		kind = string(connector.Kind)
	}
	mode := "unavailable"
	switch config.WriteMode {
	case "insert", "upsert":
		mode = config.WriteMode
	case "":
		mode = "insert"
	}
	if configErr != nil {
		config.Schema = ""
		config.Table = ""
		config.TimestampColumn = ""
		mode = "unavailable"
	}
	status := "unavailable"
	switch connector.Status {
	case "ready", "draft", "error", "disabled":
		status = string(connector.Status)
	}
	return runtimeWorkspaceDatabaseConfig{
		ID: runtimeSafeID(connector.ID), WorkspaceID: runtimeSafeID(workspaceID), Kind: kind,
		Name: runtimeSafeLabel(connector.Name), Schema: runtimeSafeLabel(config.Schema), Table: runtimeSafeLabel(config.Table),
		WriteMode: mode, WriteIntervalSeconds: connector.DefaultWriteIntervalSeconds,
		TimestampColumn: runtimeSafeLabel(config.TimestampColumn), Status: status,
	}
}

func runtimeSafeLabelPointer(value *string) *string {
	if value == nil {
		return nil
	}
	safe := runtimeSafeLabel(*value)
	if safe == "" {
		return nil
	}
	return &safe
}

func runtimeSafeDataType(value schema.DataType) schema.DataType {
	if value.IsValid() {
		return value
	}
	return ""
}

func runtimeSafeDataTypePointer(value *schema.DataType) *schema.DataType {
	if value == nil || !value.IsValid() {
		return nil
	}
	safe := *value
	return &safe
}

func runtimeSafeMappingStatus(value schema.MappingStatus) schema.MappingStatus {
	switch value {
	case "draft", "validated", "active", "out_of_sync", "error":
		return value
	}
	return ""
}

func runtimeSafeUnit(value string) string {
	if len(value) > 32 || !runtimeSafeSensorText(value) {
		return ""
	}
	for _, r := range value {
		if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '°' || r == '%' || r == '/' || r == '²' || r == '³' || r == 'µ' || r == 'Ω') {
			return ""
		}
	}
	return value
}
