package workspace

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

const (
	writeGroupConnectorKindSQLite   = "sqlite"
	writeGroupConnectorKindPostgres = "postgres"
)

// resolveWriteGroupDestinationScope derives the saved database scope from the
// connector identity. A client cannot use the destination payload to redirect
// a group to another database; an opaque connector DSN simply leaves the
// database name unresolved for a later server-side scope resolution.
func resolveWriteGroupDestinationScope(kind, rawConfig string, destination *WriteGroupDestination) error {
	config := make(map[string]json.RawMessage)
	if strings.TrimSpace(rawConfig) != "" {
		if err := json.Unmarshal([]byte(rawConfig), &config); err != nil {
			return fmt.Errorf("%w: connector scope configuration is invalid", ErrWriteGroupValidation)
		}
	}

	database := connectorScopeString(config, "database", "dbname")
	schema := connectorScopeString(config, "schema", "table_schema")
	if database == "" {
		database = connectorDatabaseFromDSN(kind, connectorScopeString(config, "dsn"))
	}
	if database == "" && kind == writeGroupConnectorKindSQLite {
		database = connectorScopeString(config, "path")
	}

	if schema == "" {
		schema = strings.TrimSpace(destination.TableSchema)
	}
	if schema == "" {
		schema = "main"
		if kind == writeGroupConnectorKindPostgres {
			schema = "public"
		}
	}

	destination.Database = database
	destination.TableSchema = schema
	return nil
}

func connectorScopeString(config map[string]json.RawMessage, keys ...string) string {
	for _, key := range keys {
		raw, ok := config[key]
		if !ok {
			continue
		}
		var value string
		if err := json.Unmarshal(raw, &value); err == nil && strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func connectorDatabaseFromDSN(kind, dsn string) string {
	dsn = strings.TrimSpace(dsn)
	if dsn == "" {
		return ""
	}
	if kind == writeGroupConnectorKindSQLite {
		return strings.TrimSpace(strings.SplitN(dsn, "?", 2)[0])
	}
	if kind != writeGroupConnectorKindPostgres {
		return ""
	}
	if parsed, err := url.Parse(dsn); err == nil && parsed.Scheme != "" {
		if database := strings.Trim(strings.TrimSpace(parsed.Path), "/"); database != "" {
			return database
		}
	}
	for _, field := range strings.Fields(dsn) {
		key, value, found := strings.Cut(field, "=")
		if found && strings.EqualFold(strings.TrimSpace(key), "dbname") {
			return strings.Trim(strings.TrimSpace(value), `"'`)
		}
	}
	return ""
}
