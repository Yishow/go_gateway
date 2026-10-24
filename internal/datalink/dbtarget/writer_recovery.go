package dbtarget

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"

	"go-gateway/internal/datalink/schema"

	"modernc.org/sqlite"
)

func writerConnectionConfig(connector *schema.DatabaseConnector) (ConnectionConfig, error) {
	config, err := parseConnectionConfig(connector.ConnectionConfig)
	if err != nil {
		return nil, err
	}
	if connector.Kind == schema.DatabaseConnectorKindSQLite {
		return prepareSQLiteDeliveryConfig(config)
	}
	return config, nil
}

// executeWithTableRecovery repairs only a missing legacy output table. The
// failed statement did not write a row, so it is safe to retry exactly once.
// Canonical write groups keep their own schema/ownership and delivery gates.
func (w *Writer) executeWithTableRecovery(ctx context.Context, db *sql.DB, connector *schema.DatabaseConnector, schemaName, tableName, query string, args []any) error {
	_, writeErr := db.ExecContext(ctx, query, args...)
	if writeErr == nil || !missingWriterTable(connector.Kind, writeErr) {
		return writeErr
	}
	if err := w.recoverWriterTable(ctx, db, connector, schemaName, tableName); err != nil {
		return errors.Join(writeErr, fmt.Errorf("自動重建資料表失敗: %w", err))
	}
	_, err := db.ExecContext(ctx, query, args...)
	return err
}

func missingWriterTable(kind schema.DatabaseConnectorKind, err error) bool {
	switch kind {
	case schema.DatabaseConnectorKindPostgres:
		var state interface{ SQLState() string }
		return errors.As(err, &state) && state.SQLState() == "42P01"
	case schema.DatabaseConnectorKindSQLite:
		var sqliteErr *sqlite.Error
		return errors.As(err, &sqliteErr) && sqliteErr.Code()&0xff == sqliteError &&
			strings.HasPrefix(sqliteErr.Error(), "SQL logic error: no such table:")
	default:
		return false
	}
}

func (w *Writer) recoverWriterTable(ctx context.Context, db *sql.DB, connector *schema.DatabaseConnector, schemaName, tableName string) error {
	// Concurrent samples must inspect again after the first recovery commits.
	w.recoveryMu.Lock()
	defer w.recoveryMu.Unlock()

	if schemaName == "" {
		schemaName = defaultSchemaForConnector(connector)
		if connector.Kind == schema.DatabaseConnectorKindPostgres {
			if err := db.QueryRowContext(ctx, `SELECT current_schema()`).Scan(&schemaName); err != nil {
				return fmt.Errorf("resolve write schema: %w", err)
			}
		}
	}
	exists, err := writerTableExists(ctx, db, connector.Kind, schemaName, tableName)
	if err != nil || exists {
		return err
	}

	err = w.createMissingWriterTable(ctx, db, connector, schemaName, tableName)
	status, summary := deliveryOutcomeSuccess, ""
	if err != nil {
		status, summary = deliveryOutcomeFailed, err.Error()
	}
	if outcomeErr := w.connectorRepo.UpdateSchemaEnsureOutcome(ctx, connector.ID, w.now().UTC(), status, summary); outcomeErr != nil {
		return errors.Join(err, fmt.Errorf("記錄自動建表結果失敗: %w", outcomeErr))
	}
	return err
}

func writerTableExists(ctx context.Context, db *sql.DB, kind schema.DatabaseConnectorKind, schemaName, tableName string) (bool, error) {
	var exists bool
	if kind == schema.DatabaseConnectorKindPostgres {
		err := db.QueryRowContext(ctx, `SELECT to_regclass($1) IS NOT NULL`, qualifiedTableName(kind, schemaName, tableName)).Scan(&exists)
		return exists, err
	}
	err := db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM sqlite_master WHERE name = ? COLLATE NOCASE)`, tableName).Scan(&exists)
	return exists, err
}

func (w *Writer) createMissingWriterTable(ctx context.Context, db *sql.DB, connector *schema.DatabaseConnector, schemaName, tableName string) error {
	current, err := w.connectorRepo.GetByID(ctx, connector.ID)
	if err != nil {
		return err
	}
	if !current.Enabled || current.IdentityRevision != connector.IdentityRevision {
		return ErrConnectorRevisionConflict
	}
	projection, err := listLiveTargetProjection(ctx, w.mappingRepo, TargetMappingListFilter{ConnectorID: &connector.ID}, w.tagReader)
	if err != nil {
		return err
	}
	mappings := make([]*schema.DatabaseTargetMapping, 0)
	for _, mapping := range projection.Mappings {
		if mapping.TableName != tableName || (mapping.TableSchema != "" && mapping.TableSchema != schemaName) ||
			(w.suppress != nil && w.suppress(connector.ID, mapping.TagID)) {
			continue
		}
		// Work from persisted members only, including columns whose samples have
		// not arrived yet. Never infer SQL types from the current sample value.
		mappings = append(mappings, mapping)
	}
	if len(mappings) == 0 {
		return fmt.Errorf("找不到此資料表的有效寫入映射")
	}
	statement, err := w.missingWriterTableStatement(ctx, connector.Kind, schemaName, tableName, mappings)
	if err != nil {
		return err
	}
	statements := []string{}
	if ensureSchema := buildEnsureSchemaStatement(connector.Kind, schemaName); ensureSchema != "" {
		statements = append(statements, ensureSchema)
	}
	statements = append(statements, statement)
	_, err = executeSchemaBatch(ctx, db, connector.Kind, schemaName, statements)
	return err
}

func (w *Writer) missingWriterTableStatement(ctx context.Context, kind schema.DatabaseConnectorKind, schemaName, tableName string, mappings []*schema.DatabaseTargetMapping) (string, error) {
	columns := map[string]string{}
	unique := map[string]bool{}
	addColumn := func(name, columnType string) error {
		if existing, ok := columns[name]; name == "" || (ok && existing != columnType) {
			return fmt.Errorf("資料表映射欄位缺失或型別衝突")
		}
		columns[name] = columnType
		return nil
	}
	for _, mapping := range mappings {
		columnType, err := resolveMappingColumnType(ctx, w.tagReader, mapping.TagID, kind)
		if err != nil {
			return "", err
		}
		if err := addColumn(mapping.ColumnName, columnType); err != nil {
			return "", err
		}
		if mapping.WriteMode == schema.DatabaseWriteModeUpsert {
			name := derefOptionalString(mapping.TimestampColumn)
			if err := addColumn(name, timestampColumnType(kind)); err != nil {
				return "", err
			}
			unique[name] = true
		}
	}
	definitions := make([]string, 0, len(columns))
	for _, name := range slices.Sorted(maps.Keys(columns)) {
		definition := quoteIdentifier(kind, name) + " " + columns[name]
		if unique[name] {
			definition += " UNIQUE"
		}
		definitions = append(definitions, definition)
	}
	// One CREATE avoids ALTERing a table another owner creates after our check.
	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", qualifiedTableName(kind, schemaName, tableName), strings.Join(definitions, ", ")), nil
}
