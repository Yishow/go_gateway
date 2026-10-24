package dbtarget

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strings"

	"go-gateway/internal/datalink/schema"
)

// IndexInfo is one explicit, non-primary-key index with its ordered key columns.
type IndexInfo struct {
	Name    string   `json:"name"`
	Columns []string `json:"columns"`
	Unique  bool     `json:"unique"`
}

// InspectTableIndexesAtRevision reads the explicit indexes of one managed table
// at the confirmed connector revision, through the same destination guard and
// read-only access as table inspection. The primary-key index is not listed (the
// key is checked from column metadata). Partial and expression indexes are
// omitted because they cannot satisfy a layout's declared index, and so is a
// PostgreSQL index that is not valid, ready and live: such an index is unsafe to
// query, does not take writes, or is being dropped, so it cannot be the
// confirmed index even when its name and columns match. A missing
// table or a SQLite file that does not exist yet has no indexes; it is never
// created.
func (i *ManagedTableInspector) InspectTableIndexesAtRevision(ctx context.Context, connectorID, expectedRevision, schemaName, tableName string) ([]IndexInfo, error) {
	if i == nil || i.service == nil || i.service.repo == nil {
		return nil, managedDestinationInvalid("saved connector service is unavailable")
	}
	connector, err := i.service.ResolveSavedTarget(ctx, strings.TrimSpace(connectorID), strings.TrimSpace(expectedRevision))
	if err != nil {
		return nil, err
	}
	missing, err := i.guardManagedConnector(ctx, connector, schemaName)
	if err != nil {
		return nil, err
	}
	if missing {
		return []IndexInfo{}, nil
	}
	config, err := parseConnectionConfig(connector.ConnectionConfig)
	if err != nil {
		return nil, managedDestinationInvalid("saved connector configuration is invalid")
	}
	if connector.Kind == schema.DatabaseConnectorKindSQLite {
		return inspectReadOnlySQLiteIndexes(ctx, config, strings.TrimSpace(tableName))
	}
	manager, err := openExternalDBManagerFunc(connector.Kind, config)
	if err != nil {
		return nil, fmt.Errorf("open PostgreSQL destination for index inspection: %w", err)
	}
	defer manager.Close()
	return inspectPostgresIndexes(ctx, manager.DB(), defaultString(strings.TrimSpace(schemaName), "public"), strings.TrimSpace(tableName))
}

func inspectReadOnlySQLiteIndexes(ctx context.Context, config ConnectionConfig, tableName string) ([]IndexInfo, error) {
	dsn, err := sqliteReadOnlyDSN(config)
	if err != nil {
		return nil, managedDestinationInvalid("saved SQLite connector must use a local file destination")
	}
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open read-only sqlite destination: %w", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	return inspectSQLiteIndexes(ctx, db, tableName)
}

// inspectSQLiteIndexes lists the explicit indexes of a table on an open handle.
func inspectSQLiteIndexes(ctx context.Context, db *sql.DB, tableName string) ([]IndexInfo, error) {
	type listed struct {
		name   string
		unique bool
	}
	var names []listed
	rows, err := db.QueryContext(ctx, fmt.Sprintf(`PRAGMA index_list('%s')`, escapeSQLiteStringLiteral(tableName)))
	if err != nil {
		return nil, fmt.Errorf("list sqlite indexes: %w", err)
	}
	for rows.Next() {
		var seq, unique, partial int
		var name, origin string
		if err := rows.Scan(&seq, &name, &unique, &origin, &partial); err != nil {
			_ = rows.Close()
			return nil, fmt.Errorf("scan sqlite index: %w", err)
		}
		// origin c is an explicit CREATE INDEX; pk and u are constraint indexes.
		if origin == "c" && partial == 0 {
			names = append(names, listed{name: name, unique: unique == 1})
		}
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return nil, fmt.Errorf("iterate sqlite indexes: %w", err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("close sqlite index list: %w", err)
	}

	indexes := make([]IndexInfo, 0, len(names))
	for _, index := range names {
		info, err := db.QueryContext(ctx, fmt.Sprintf(`PRAGMA index_info('%s')`, escapeSQLiteStringLiteral(index.name)))
		if err != nil {
			return nil, fmt.Errorf("read sqlite index columns: %w", err)
		}
		entry := IndexInfo{Name: index.name, Unique: index.unique, Columns: []string{}}
		expression := false
		for info.Next() {
			var seqno, cid int
			var column sql.NullString
			if err := info.Scan(&seqno, &cid, &column); err != nil {
				_ = info.Close()
				return nil, fmt.Errorf("scan sqlite index column: %w", err)
			}
			expression = expression || !column.Valid
			entry.Columns = append(entry.Columns, column.String)
		}
		if err := info.Close(); err != nil {
			return nil, fmt.Errorf("close sqlite index columns: %w", err)
		}
		if !expression {
			indexes = append(indexes, entry)
		}
	}
	slices.SortFunc(indexes, func(a, b IndexInfo) int { return strings.Compare(a.Name, b.Name) })
	return indexes, nil
}

func inspectPostgresIndexes(ctx context.Context, db *sql.DB, schemaName, tableName string) ([]IndexInfo, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT i.relname, ix.indisunique, string_agg(a.attname::text, ',' ORDER BY k.ord)
		FROM pg_catalog.pg_index ix
		JOIN pg_catalog.pg_class t ON t.oid = ix.indrelid
		JOIN pg_catalog.pg_namespace n ON n.oid = t.relnamespace
		JOIN pg_catalog.pg_class i ON i.oid = ix.indexrelid
		CROSS JOIN LATERAL unnest(ix.indkey::int2[]) WITH ORDINALITY AS k(attnum, ord)
		JOIN pg_catalog.pg_attribute a ON a.attrelid = t.oid AND a.attnum = k.attnum
		WHERE n.nspname = $1 AND t.relname = $2
		  AND NOT ix.indisprimary AND ix.indpred IS NULL AND ix.indexprs IS NULL
		  AND ix.indisvalid AND ix.indisready AND ix.indislive
		GROUP BY i.relname, ix.indisunique
		ORDER BY i.relname
	`, schemaName, tableName)
	if err != nil {
		return nil, fmt.Errorf("query postgres indexes: %w", err)
	}
	defer rows.Close()
	indexes := []IndexInfo{}
	for rows.Next() {
		var entry IndexInfo
		var columns string
		if err := rows.Scan(&entry.Name, &entry.Unique, &columns); err != nil {
			return nil, fmt.Errorf("scan postgres index: %w", err)
		}
		entry.Columns = strings.Split(columns, ",")
		indexes = append(indexes, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate postgres indexes: %w", err)
	}
	return indexes, nil
}

// guardManagedConnector applies the managed destination path guard shared by
// table and index inspection. missing is true only for a SQLite file that does
// not exist yet; nothing is opened or created in that case.
func (i *ManagedTableInspector) guardManagedConnector(ctx context.Context, connector *schema.DatabaseConnector, schemaName string) (missing bool, err error) {
	switch connector.Kind {
	case schema.DatabaseConnectorKindSQLite:
		if strings.TrimSpace(schemaName) != "" && strings.TrimSpace(schemaName) != managedSQLiteSchemaName {
			return false, managedDestinationInvalid("managed SQLite schema must be main")
		}
		return i.inspectSQLiteDestination(ctx, connector)
	case schema.DatabaseConnectorKindPostgres:
		return false, validateManagedPostgresConfig(connector)
	default:
		return false, managedDestinationInvalid("saved connector kind is unsupported")
	}
}
