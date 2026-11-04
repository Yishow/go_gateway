package dbtarget

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strings"

	"go-gateway/internal/datalink/schema"
)

// CatalogSchema is one schema visible to the saved connector. Tables is kept
// as an empty (rather than nil) slice when the schema has no user tables so a
// caller can represent an accessible, empty schema without inventing a row.
type CatalogSchema struct {
	Name   string   `json:"name"`
	Tables []string `json:"tables"`
}

// DatabaseCatalog is a read-only projection of the saved connector catalog.
// It intentionally contains no connection configuration or server details.
type DatabaseCatalog struct {
	Kind          schema.DatabaseConnectorKind `json:"kind"`
	DefaultSchema string                       `json:"default_schema"`
	Schemas       []CatalogSchema              `json:"schemas"`
}

// ReadCatalog resolves the enabled connector at an exact identity revision,
// then reads its catalog without running a probe, DDL, or schema creation.
func (s *ConnectorService) ReadCatalog(ctx context.Context, connectorID, expectedRevision string) (*DatabaseCatalog, error) {
	connector, err := s.ResolveSavedTarget(ctx, connectorID, expectedRevision)
	if err != nil {
		return nil, err
	}
	return readSavedConnectorCatalog(ctx, connector)
}

func readSavedConnectorCatalog(ctx context.Context, connector *schema.DatabaseConnector) (*DatabaseCatalog, error) {
	if connector == nil {
		return nil, fmt.Errorf("read database catalog: %w", ErrConnectorNotFound)
	}
	config, err := parseConnectionConfig(connector.ConnectionConfig)
	if err != nil {
		return nil, err
	}

	var catalog *DatabaseCatalog
	if connector.Kind == schema.DatabaseConnectorKindSQLite {
		catalog, err = readSQLiteCatalog(ctx, config, connector.Kind)
	} else {
		manager, openErr := openExternalDBManagerFunc(connector.Kind, config)
		if openErr != nil {
			return nil, openErr
		}
		if manager == nil || manager.DB() == nil {
			if manager != nil {
				_ = manager.Close()
			}
			return nil, fmt.Errorf("read database catalog: empty database connection")
		}
		defer manager.Close()
		catalog, err = readCatalogByKind(ctx, manager.DB(), connector.Kind, config)
	}
	if err != nil {
		return nil, err
	}
	return catalog, nil
}

func readSQLiteCatalog(ctx context.Context, config ConnectionConfig, kind schema.DatabaseConnectorKind) (*DatabaseCatalog, error) {
	dsn, err := sqliteReadOnlyDSN(config)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open(string(SQLDialectSQLite), dsn)
	if err != nil {
		return nil, fmt.Errorf("open read-only sqlite catalog: %w", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("read sqlite catalog: %w", err)
	}
	return readCatalogByKind(ctx, db, kind, config)
}

func readCatalogByKind(ctx context.Context, db *sql.DB, kind schema.DatabaseConnectorKind, config ConnectionConfig) (*DatabaseCatalog, error) {
	if db == nil {
		return nil, fmt.Errorf("read database catalog: database connection is unavailable")
	}

	var (
		schemas []catalogSchemaRow
		tables  []catalogTableRow
		err     error
		filter  catalogFilter
	)
	switch kind {
	case schema.DatabaseConnectorKindSQLite:
		tables, err = readSQLiteCatalogTables(ctx, db)
		schemas = []catalogSchemaRow{{Name: managedSQLiteSchemaName}}
		filter = catalogFilterSQLite
	case schema.DatabaseConnectorKindPostgres:
		schemas, err = readCatalogSchemas(ctx, db, `
			SELECT schema_name
			FROM information_schema.schemata
			ORDER BY schema_name
		`)
		if err == nil {
			tables, err = readCatalogTables(ctx, db, `
				SELECT table_schema, table_name
				FROM information_schema.tables
				WHERE table_type = 'BASE TABLE'
				ORDER BY table_schema, table_name
			`)
		}
		filter = catalogFilterPostgres
	case schema.DatabaseConnectorKindMySQL:
		schemas, err = readCatalogSchemas(ctx, db, `
			SELECT schema_name
			FROM information_schema.schemata
			ORDER BY schema_name
		`)
		if err == nil {
			tables, err = readCatalogTables(ctx, db, `
				SELECT table_schema, table_name
				FROM information_schema.tables
				WHERE table_type = 'BASE TABLE'
				ORDER BY table_schema, table_name
			`)
		}
		filter = catalogFilterMySQL
	case schema.DatabaseConnectorKindSQLServer:
		schemas, err = readCatalogSchemas(ctx, db, `
			SELECT name
			FROM sys.schemas
			ORDER BY name
		`)
		if err == nil {
			tables, err = readCatalogTables(ctx, db, `
				SELECT SCHEMA_NAME(schema_id), name
				FROM sys.objects
				WHERE type = 'U'
				  AND ISNULL(is_ms_shipped, 0) = 0
				ORDER BY SCHEMA_NAME(schema_id), name
			`)
		}
		filter = catalogFilterSQLServer
	default:
		return nil, fmt.Errorf("unsupported database catalog kind: %s", kind)
	}
	if err != nil {
		return nil, err
	}

	defaultSchema := defaultSchemaForCatalogConnector(kind, config)
	if kind == schema.DatabaseConnectorKindSQLServer && defaultSchema == "" {
		if err := db.QueryRowContext(ctx, `SELECT SCHEMA_NAME()`).Scan(&defaultSchema); err != nil {
			return nil, fmt.Errorf("query sql server default schema: %w", err)
		}
		if strings.TrimSpace(defaultSchema) == "" {
			return nil, fmt.Errorf("query sql server default schema: empty schema")
		}
	}

	return &DatabaseCatalog{
		Kind:          kind,
		DefaultSchema: defaultSchema,
		Schemas:       catalogSchemasFromRows(schemas, tables, filter),
	}, nil
}

func readCatalogSchemas(ctx context.Context, db *sql.DB, query string) ([]catalogSchemaRow, error) {
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query database schemas: %w", err)
	}
	defer rows.Close()

	result := make([]catalogSchemaRow, 0)
	for rows.Next() {
		var row catalogSchemaRow
		if err := rows.Scan(&row.Name); err != nil {
			return nil, fmt.Errorf("scan database schema: %w", err)
		}
		result = append(result, row)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate database schemas: %w", err)
	}
	return result, nil
}

func readCatalogTables(ctx context.Context, db *sql.DB, query string) ([]catalogTableRow, error) {
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query database tables: %w", err)
	}
	defer rows.Close()

	result := make([]catalogTableRow, 0)
	for rows.Next() {
		var row catalogTableRow
		if err := rows.Scan(&row.Schema, &row.Name); err != nil {
			return nil, fmt.Errorf("scan database table: %w", err)
		}
		result = append(result, row)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate database tables: %w", err)
	}
	return result, nil
}

func readSQLiteCatalogTables(ctx context.Context, db *sql.DB) ([]catalogTableRow, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT name
		FROM sqlite_master
		WHERE type = 'table'
		ORDER BY name
	`)
	if err != nil {
		return nil, fmt.Errorf("query sqlite tables: %w", err)
	}
	defer rows.Close()

	result := make([]catalogTableRow, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("scan sqlite table: %w", err)
		}
		result = append(result, catalogTableRow{Schema: managedSQLiteSchemaName, Name: name})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate sqlite tables: %w", err)
	}
	return result, nil
}

type catalogSchemaRow struct{ Name string }

type catalogTableRow struct {
	Schema string
	Name   string
}

type catalogFilter uint8

const (
	catalogFilterSQLite catalogFilter = iota + 1
	catalogFilterPostgres
	catalogFilterMySQL
	catalogFilterSQLServer
)

const catalogInformationSchemaName = "information_schema"

func catalogSchemasFromRows(schemaRows []catalogSchemaRow, tableRows []catalogTableRow, filter catalogFilter) []CatalogSchema {
	bySchema := make(map[string]*CatalogSchema, len(schemaRows)+len(tableRows))
	for _, row := range schemaRows {
		name := row.Name
		if strings.TrimSpace(name) == "" || !catalogSchemaAllowed(name, filter) {
			continue
		}
		if _, exists := bySchema[name]; !exists {
			bySchema[name] = &CatalogSchema{Name: name, Tables: make([]string, 0)}
		}
	}
	for _, row := range tableRows {
		schemaName := row.Schema
		tableName := row.Name
		if strings.TrimSpace(schemaName) == "" || strings.TrimSpace(tableName) == "" || !catalogSchemaAllowed(schemaName, filter) || !catalogTableAllowed(tableName, filter) {
			continue
		}
		entry, exists := bySchema[schemaName]
		if !exists {
			entry = &CatalogSchema{Name: schemaName, Tables: make([]string, 0)}
			bySchema[schemaName] = entry
		}
		entry.Tables = append(entry.Tables, tableName)
	}

	result := make([]CatalogSchema, 0, len(bySchema))
	for _, entry := range bySchema {
		slices.Sort(entry.Tables)
		entry.Tables = slices.Compact(entry.Tables)
		result = append(result, *entry)
	}
	slices.SortFunc(result, func(a, b CatalogSchema) int { return strings.Compare(a.Name, b.Name) })
	return result
}

func catalogSchemaAllowed(name string, filter catalogFilter) bool {
	switch filter {
	case catalogFilterSQLite:
		return name == managedSQLiteSchemaName
	case catalogFilterPostgres:
		return name != "pg_catalog" && name != catalogInformationSchemaName && !strings.HasPrefix(name, "pg_")
	case catalogFilterMySQL:
		return name != "mysql" && name != catalogInformationSchemaName && name != "performance_schema" && name != "sys"
	case catalogFilterSQLServer:
		return name != "sys" && name != "INFORMATION_SCHEMA"
	default:
		return false
	}
}

func catalogTableAllowed(name string, filter catalogFilter) bool {
	if filter == catalogFilterSQLite {
		return !strings.HasPrefix(name, "sqlite_")
	}
	return true
}

func defaultSchemaForCatalogConnector(kind schema.DatabaseConnectorKind, config ConnectionConfig) string {
	if schemaName := stringConfigValue(config, "schema"); kind == schema.DatabaseConnectorKindPostgres && schemaName != "" {
		return schemaName
	}
	if kind == schema.DatabaseConnectorKindMySQL {
		return defaultString(stringConfigValue(config, "database"), stringConfigValue(config, "dbname"))
	}
	if kind == schema.DatabaseConnectorKindSQLServer {
		return stringConfigValue(config, "schema")
	}
	return defaultSchemaForKind(kind)
}
