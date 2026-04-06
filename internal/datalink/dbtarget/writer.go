package dbtarget

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go-gateway/internal/datalink/schema"
)

type Writer struct {
	connectorRepo ConnectorRepository
	mappingRepo   TargetMappingRepository
}

func NewWriter(connectorRepo ConnectorRepository, mappingRepo TargetMappingRepository) *Writer {
	return &Writer{
		connectorRepo: connectorRepo,
		mappingRepo:   mappingRepo,
	}
}

func (w *Writer) WriteTagValue(ctx context.Context, tagID string, value any, observedAt time.Time) error {
	if strings.TrimSpace(tagID) == "" {
		return fmt.Errorf("tagID 不可為空")
	}

	enabled := true
	mappings, err := w.mappingRepo.List(ctx, TargetMappingListFilter{
		TagID:   &tagID,
		Enabled: &enabled,
	})
	if err != nil {
		return fmt.Errorf("查詢資料庫目標映射失敗: %w", err)
	}
	if len(mappings) == 0 {
		return nil
	}

	var failures []string
	for _, mapping := range mappings {
		connector, err := w.connectorRepo.GetByID(ctx, mapping.ConnectorID)
		if err != nil {
			recordWriteHistory(mapping.ConnectorID, "failed", 0, err.Error())
			failures = append(failures, fmt.Sprintf("mapping %s 取得連接器失敗: %v", mapping.ID, err))
			continue
		}
		if !connector.Enabled {
			continue
		}

		if err := w.writeMapping(ctx, connector, mapping, value, observedAt); err != nil {
			recordWriteHistory(connector.ID, "failed", 0, err.Error())
			failures = append(failures, fmt.Sprintf("mapping %s 寫入失敗: %v", mapping.ID, err))
			continue
		}
		recordWriteHistory(connector.ID, "success", 1, "")
	}

	if len(failures) > 0 {
		return fmt.Errorf("%s", strings.Join(failures, "; "))
	}

	return nil
}

func (w *Writer) writeMapping(
	ctx context.Context,
	connector *schema.DatabaseConnector,
	mapping *schema.DatabaseTargetMapping,
	value any,
	observedAt time.Time,
) error {
	connectionConfig, err := parseConnectionConfig(connector.ConnectionConfig)
	if err != nil {
		return err
	}

	manager, err := openExternalDBManager(connector.Kind, connectionConfig)
	if err != nil {
		return err
	}
	defer manager.Close()

	query, args, err := buildWriteStatement(connector.Kind, mapping, value, observedAt)
	if err != nil {
		return err
	}

	if _, err := manager.DB().ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("執行資料庫寫入失敗: %w", err)
	}

	return nil
}

func buildWriteStatement(
	kind schema.DatabaseConnectorKind,
	mapping *schema.DatabaseTargetMapping,
	value any,
	observedAt time.Time,
) (string, []any, error) {
	if mapping.TableName == "" || mapping.ColumnName == "" {
		return "", nil, fmt.Errorf("資料庫目標映射缺少資料表或欄位資訊")
	}

	columns := []string{quoteIdentifier(kind, mapping.ColumnName)}
	args := []any{normalizeDBValue(value)}
	if mapping.WriteMode == schema.DatabaseWriteModeUpsert && mapping.TimestampColumn != nil {
		columns = append(columns, quoteIdentifier(kind, *mapping.TimestampColumn))
		args = append(args, observedAt.UTC())
	}

	placeholders := buildPlaceholders(kind, len(args))
	tableRef := qualifiedTableName(kind, mapping.TableSchema, mapping.TableName)
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		tableRef,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	if mapping.WriteMode == schema.DatabaseWriteModeUpsert {
		if mapping.TimestampColumn == nil {
			return "", nil, fmt.Errorf("upsert 模式缺少 timestamp_column")
		}

		valueColumn := quoteIdentifier(kind, mapping.ColumnName)
		timestampColumn := quoteIdentifier(kind, *mapping.TimestampColumn)
		query += fmt.Sprintf(
			" ON CONFLICT (%s) DO UPDATE SET %s = EXCLUDED.%s",
			timestampColumn,
			valueColumn,
			valueColumn,
		)
	}

	return query, args, nil
}

func buildPlaceholders(kind schema.DatabaseConnectorKind, count int) []string {
	placeholders := make([]string, 0, count)
	for i := 0; i < count; i++ {
		switch kind {
		case schema.DatabaseConnectorKindPostgres:
			placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		default:
			placeholders = append(placeholders, "?")
		}
	}
	return placeholders
}

func quoteIdentifier(kind schema.DatabaseConnectorKind, name string) string {
	escaped := strings.ReplaceAll(strings.TrimSpace(name), `"`, `""`)
	return `"` + escaped + `"`
}

func qualifiedTableName(kind schema.DatabaseConnectorKind, schemaName string, tableName string) string {
	quotedTable := quoteIdentifier(kind, tableName)
	if strings.TrimSpace(schemaName) == "" || kind == schema.DatabaseConnectorKindSQLite {
		return quotedTable
	}
	return quoteIdentifier(kind, schemaName) + "." + quotedTable
}

func normalizeDBValue(value any) any {
	switch typed := value.(type) {
	case nil:
		return nil
	case map[string]any, []any:
		data, err := json.Marshal(typed)
		if err != nil {
			return fmt.Sprintf("%v", typed)
		}
		return string(data)
	case time.Time:
		return typed.UTC()
	default:
		return value
	}
}

var _ interface {
	WriteTagValue(context.Context, string, any, time.Time) error
} = (*Writer)(nil)
