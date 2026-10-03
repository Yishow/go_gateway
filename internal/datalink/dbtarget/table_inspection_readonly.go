package dbtarget

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	"go-gateway/internal/datalink/schema"
)

// ReadOnlyTableInspector wraps ConnectorService with a SQLite inspection path
// that opens the saved file in mode=ro. PostgreSQL uses the existing catalog
// inspection queries against the same saved connector snapshot.
type ReadOnlyTableInspector struct {
	service *ConnectorService
}

// NewReadOnlyTableInspector creates a destination metadata inspector that
// refuses to create a missing SQLite database while inspecting it.
func NewReadOnlyTableInspector(service *ConnectorService) *ReadOnlyTableInspector {
	return &ReadOnlyTableInspector{service: service}
}

// InspectTable returns metadata for one saved connector destination without
// provisioning or altering the destination database.
func (i *ReadOnlyTableInspector) InspectTable(ctx context.Context, connectorID, schemaName, tableName string) (*TableInspection, error) {
	if i == nil || i.service == nil || i.service.repo == nil {
		return nil, fmt.Errorf("inspect destination table: %w", ErrConnectorNotFound)
	}
	connector, err := i.service.repo.GetByID(ctx, strings.TrimSpace(connectorID))
	if err != nil {
		return nil, fmt.Errorf("inspect destination connector: %w", err)
	}
	return inspectReadOnlySavedConnectorTable(ctx, connector, schemaName, tableName)
}

func inspectReadOnlySavedConnectorTable(ctx context.Context, connector *schema.DatabaseConnector, schemaName, tableName string) (*TableInspection, error) {
	if connector.Kind != schema.DatabaseConnectorKindSQLite {
		return inspectSavedConnectorTable(ctx, connector, schemaName, tableName)
	}
	config, err := parseConnectionConfig(connector.ConnectionConfig)
	if err != nil {
		return nil, err
	}
	dsn, err := sqliteReadOnlyDSN(config)
	if err != nil {
		return &TableInspection{
			Status: TableInspectionFailed, Schema: strings.TrimSpace(schemaName),
			Table: strings.TrimSpace(tableName), Columns: []ColumnInfo{},
			Reason: inspectionReasonInvalidConfiguration,
		}, nil
	}
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open read-only sqlite destination: %w", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	if err := db.PingContext(ctx); err != nil {
		return &TableInspection{
			Status: TableInspectionFailed, Schema: strings.TrimSpace(schemaName),
			Table: strings.TrimSpace(tableName), Columns: []ColumnInfo{},
			Reason: inspectionReasonConnectionFailed,
		}, nil
	}
	tables, err := inspectSQLiteTables(ctx, db)
	if err != nil {
		return (&TableInspection{
			Schema: strings.TrimSpace(schemaName), Table: strings.TrimSpace(tableName),
			Columns: []ColumnInfo{},
		}).classified(err), nil
	}
	result := &TableInspection{
		Schema: strings.TrimSpace(schemaName), Table: strings.TrimSpace(tableName),
		Columns: []ColumnInfo{},
	}
	result.Schema = defaultString(result.Schema, "main")
	return result.fromListing(tables, TableInspectionMissing, ""), nil
}

func inspectSavedConnectorTable(
	ctx context.Context,
	connector *schema.DatabaseConnector,
	schemaName, tableName string,
) (*TableInspection, error) {
	result := &TableInspection{
		Schema: strings.TrimSpace(schemaName), Table: strings.TrimSpace(tableName), Columns: []ColumnInfo{},
	}
	if connector.Kind != schema.DatabaseConnectorKindPostgres {
		return result.withStatus(TableInspectionFailed, inspectionReasonUnsupportedKind), nil
	}
	config, err := parseConnectionConfig(connector.ConnectionConfig)
	if err != nil {
		return result.withStatus(TableInspectionFailed, inspectionReasonInvalidConfiguration), nil
	}
	manager, err := openExternalDBManagerFunc(connector.Kind, config)
	if err != nil {
		if classifyConnectorError(err) == schema.DatabaseConnectorStatusAuthFailed {
			return result.withStatus(TableInspectionFailed, inspectionReasonAuthenticationFailed), nil
		}
		return result.withStatus(TableInspectionFailed, inspectionReasonConnectionFailed), nil
	}
	defer manager.Close()
	switch connector.Kind {
	case schema.DatabaseConnectorKindPostgres:
		result.Schema = defaultString(result.Schema, "public")
		exists, err := postgresCatalogHasTable(ctx, manager.DB(), result.Schema, result.Table)
		if err != nil {
			return result.classified(err), nil
		}
		if !exists {
			return result.withStatus(TableInspectionMissing, ""), nil
		}
		tables, err := inspectPostgresTables(ctx, manager.DB())
		if err != nil {
			return result.classified(err), nil
		}
		return result.fromListing(tables, TableInspectionForbidden, inspectionReasonPermissionDenied), nil
	default:
		return result.withStatus(TableInspectionFailed, inspectionReasonUnsupportedKind), nil
	}
}

func sqliteReadOnlyDSN(config ConnectionConfig) (string, error) {
	dsn := strings.TrimSpace(stringConfigValue(config, "dsn"))
	if dsn == "" {
		dsn = strings.TrimSpace(stringConfigValue(config, "path"))
	}
	if dsn == "" || strings.EqualFold(dsn, ":memory:") {
		return "", fmt.Errorf("sqlite read-only inspection requires a file dsn")
	}
	if strings.HasPrefix(strings.ToLower(dsn), "file:") {
		parsed, err := url.Parse(dsn)
		if err != nil || (parsed.Host != "" && !strings.EqualFold(parsed.Host, "localhost")) {
			return "", fmt.Errorf("sqlite read-only inspection requires a local file dsn")
		}
		if strings.EqualFold(parsed.Opaque, ":memory:") || strings.EqualFold(parsed.Path, ":memory:") {
			return "", fmt.Errorf("sqlite read-only inspection does not support memory databases")
		}
		query := parsed.Query()
		if strings.EqualFold(query.Get("mode"), "memory") {
			return "", fmt.Errorf("sqlite read-only inspection does not support memory databases")
		}
		query.Set("mode", "ro")
		query.Del("_pragma")
		parsed.RawQuery = query.Encode()
		return parsed.String(), nil
	}
	base, queryString, found := strings.Cut(dsn, "?")
	if found {
		if strings.EqualFold(strings.TrimSpace(base), ":memory:") || strings.EqualFold(strings.TrimSpace(base), "memory") {
			return "", fmt.Errorf("sqlite read-only inspection does not support memory databases")
		}
		query, err := url.ParseQuery(queryString)
		if err != nil {
			return "", fmt.Errorf("parse sqlite dsn query: %w", err)
		}
		if strings.EqualFold(query.Get("mode"), "memory") {
			return "", fmt.Errorf("sqlite read-only inspection does not support memory databases")
		}
		query.Set("mode", "ro")
		query.Del("_pragma")
		return "file:" + base + "?" + query.Encode(), nil
	}
	return "file:" + dsn + "?mode=ro", nil
}
