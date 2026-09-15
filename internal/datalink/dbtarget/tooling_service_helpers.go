package dbtarget

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"go-gateway/internal/datalink/schema"
)

func resolveMappingColumnType(
	ctx context.Context,
	tagReader ConnectorTagReader,
	tagID string,
	kind schema.DatabaseConnectorKind,
) (string, error) {
	if tagReader == nil || strings.TrimSpace(tagID) == "" {
		return "TEXT", nil
	}

	tagEntity, err := tagReader.GetByID(ctx, tagID)
	if err != nil {
		return "", fmt.Errorf("取得標籤失敗: %w", err)
	}

	return columnTypeForDataType(kind, tagEntity.DataType), nil
}

func buildCreateTableStatement(
	kind schema.DatabaseConnectorKind,
	schemaName string,
	tableName string,
	valueColumn string,
	valueColumnType string,
	mappingRecord *schema.DatabaseTargetMapping,
) string {
	columns := []string{
		fmt.Sprintf("%s %s", quoteIdentifier(kind, valueColumn), valueColumnType),
	}
	if mappingRecord != nil && mappingRecord.WriteMode == schema.DatabaseWriteModeUpsert && mappingRecord.TimestampColumn != nil {
		ts := strings.TrimSpace(*mappingRecord.TimestampColumn)
		if ts != "" {
			columns = append(
				[]string{
					fmt.Sprintf("%s %s PRIMARY KEY", quoteIdentifier(kind, ts), timestampColumnType(kind)),
				},
				columns...,
			)
		}
	}

	return fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s (%s)",
		qualifiedTableName(kind, schemaName, tableName),
		strings.Join(columns, ", "),
	)
}

func buildAddColumnStatement(
	kind schema.DatabaseConnectorKind,
	schemaName string,
	tableName string,
	columnName string,
	columnType string,
) string {
	return fmt.Sprintf(
		"ALTER TABLE %s ADD COLUMN %s %s",
		qualifiedTableName(kind, schemaName, tableName),
		quoteIdentifier(kind, columnName),
		columnType,
	)
}

func buildEnsureUniqueIndexStatement(
	kind schema.DatabaseConnectorKind,
	schemaName string,
	tableName string,
	columnName string,
) (sqlStatement, normalizedIndexName string) {
	indexName := sanitizeIndexName(fmt.Sprintf("%s_%s_%s_uniq", schemaName, tableName, columnName))
	tableRef := qualifiedTableName(kind, schemaName, tableName)
	// MySQL 的 CREATE INDEX 不接受 IF NOT EXISTS；重複建立由欄位既有的主鍵 /
	// 唯一鍵資訊與同批次索引名稱去重擋下，不依賴語句層的守衛。
	existsGuard := "IF NOT EXISTS "
	if kind == schema.DatabaseConnectorKindMySQL {
		existsGuard = ""
	}
	statement := fmt.Sprintf(
		"CREATE UNIQUE INDEX %s%s ON %s (%s)",
		existsGuard,
		quoteIdentifier(kind, indexName),
		tableRef,
		quoteIdentifier(kind, columnName),
	)
	return statement, strings.ToLower(indexName)
}

func normalizeDryRunCandidateIDs(
	candidateIDs []string,
	mappings []*schema.DatabaseTargetMapping,
) (selectedIDs map[string]struct{}, unknownIDs []string) {
	if len(candidateIDs) == 0 {
		return nil, nil
	}

	existing := make(map[string]struct{}, len(mappings))
	for _, mappingRecord := range mappings {
		if mappingRecord == nil {
			continue
		}
		existing[mappingRecord.ID] = struct{}{}
	}

	selected := make(map[string]struct{}, len(candidateIDs))
	unknown := make([]string, 0)
	for _, candidateID := range candidateIDs {
		normalized := strings.TrimSpace(candidateID)
		if normalized == "" {
			continue
		}
		if _, seen := selected[normalized]; seen {
			continue
		}
		selected[normalized] = struct{}{}
		if _, found := existing[normalized]; !found {
			unknown = append(unknown, normalized)
		}
	}
	sort.Strings(unknown)
	return selected, unknown
}

func selectMappingsForDryRun(
	mappings []*schema.DatabaseTargetMapping,
	candidateIDs map[string]struct{},
) []*schema.DatabaseTargetMapping {
	selected := make([]*schema.DatabaseTargetMapping, 0, len(mappings))
	for _, mappingRecord := range mappings {
		if mappingRecord == nil {
			continue
		}
		if candidateIDs != nil {
			if _, include := candidateIDs[mappingRecord.ID]; !include {
				continue
			}
		}
		selected = append(selected, mappingRecord)
	}
	sort.Slice(selected, func(i, j int) bool {
		return selected[i].ID < selected[j].ID
	})
	return selected
}

func firstBlockingIssue(issues []ValidationIssue) (ValidationIssue, bool) {
	for _, issue := range issues {
		if strings.TrimSpace(issue.Severity) == validationSeverityError {
			return issue, true
		}
	}
	return ValidationIssue{}, false
}

func mapValidationIssueCodeForDryRun(code string) string {
	normalized := strings.TrimSpace(strings.ToLower(code))
	switch normalized {
	case validationTableMissingCode, validationColumnMissingCode, "timestamp_missing", validationTimestampColumnMissingCode, "timestamp_column_not_unique", "mapping_missing", validationTagMissingCode:
		return "schema_missing"
	case "column_type_mismatch":
		return "type_conflict"
	case "upsert_value_column_unique":
		return "conflict"
	default:
		return "conflict"
	}
}

func normalizeTableKey(schemaName, tableName string) tableKey {
	return tableKey{
		schema: strings.ToLower(strings.TrimSpace(schemaName)),
		name:   strings.ToLower(strings.TrimSpace(tableName)),
	}
}

func dedupeStatements(statements []string) []string {
	seen := make(map[string]struct{}, len(statements))
	result := make([]string, 0, len(statements))
	for _, statement := range statements {
		normalized := strings.TrimSpace(statement)
		if normalized == "" {
			continue
		}
		if _, exists := seen[normalized]; exists {
			continue
		}
		seen[normalized] = struct{}{}
		result = append(result, normalized)
	}
	return result
}

func sanitizeIndexName(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "idx_generated_unique"
	}
	var builder strings.Builder
	for _, char := range trimmed {
		switch {
		case char >= 'a' && char <= 'z':
			builder.WriteRune(char)
		case char >= 'A' && char <= 'Z':
			builder.WriteRune(char + ('a' - 'A'))
		case char >= '0' && char <= '9':
			builder.WriteRune(char)
		default:
			builder.WriteRune('_')
		}
	}
	return strings.Trim(builder.String(), "_")
}

func columnTypeForDataType(kind schema.DatabaseConnectorKind, dataType schema.DataType) string {
	switch dataType {
	case schema.DataTypeBool:
		if kind == schema.DatabaseConnectorKindPostgres {
			return "BOOLEAN"
		}
		return "INTEGER"
	case schema.DataTypeInt16, schema.DataTypeUint16, schema.DataTypeInt32, schema.DataTypeUint32, schema.DataTypeInt64, schema.DataTypeUint64:
		if kind == schema.DatabaseConnectorKindPostgres {
			return "BIGINT"
		}
		return "INTEGER"
	case schema.DataTypeFloat32, schema.DataTypeFloat64:
		if kind == schema.DatabaseConnectorKindPostgres {
			return "DOUBLE PRECISION"
		}
		return "REAL"
	default:
		return "TEXT"
	}
}

func timestampColumnType(kind schema.DatabaseConnectorKind) string {
	if kind == schema.DatabaseConnectorKindPostgres {
		return "TIMESTAMPTZ"
	}
	return "DATETIME"
}

// buildEnsureSchemaStatement 回傳建立資料表所屬容器的語句：PostgreSQL 為 schema，
// MySQL 為 database。SQLite 沒有對應概念，回傳空字串代表不需要語句。
func buildEnsureSchemaStatement(kind schema.DatabaseConnectorKind, schemaName string) string {
	if strings.TrimSpace(schemaName) == "" {
		return ""
	}
	switch kind {
	case schema.DatabaseConnectorKindPostgres:
		return fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", quoteIdentifier(kind, schemaName))
	case schema.DatabaseConnectorKindMySQL:
		return fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", quoteIdentifier(kind, schemaName))
	default:
		return ""
	}
}
