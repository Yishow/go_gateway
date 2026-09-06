package datalink

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// legacyMySQLTargetSchema 是 MySQL 連接器在支援 MySQL 檢查之前被寫入的預設
// table_schema。當時該值來自 SQLite 的 "main"，在 MySQL 上是不存在的資料庫名稱。
const legacyMySQLTargetSchema = "main"

// ensureMySQLTargetMappingSchema 把 MySQL 連接器底下殘留的 "main" table_schema
// 正規化為連線設定實際指向的資料庫名稱。沒有這一步，舊資料列會在檢查時對不上
// 資料表清單，並產生指向 `main` 的語句而失敗於 Unknown database。
func ensureMySQLTargetMappingSchema(db *sql.DB) error {
	const migrationName = "018_database_target_mysql_schema_normalization"

	rows, err := db.QueryContext(
		context.Background(),
		`SELECT c.id, c.connection_config
		 FROM database_connectors c
		 WHERE c.kind = 'mysql'
		   AND EXISTS (
		       SELECT 1 FROM database_target_mappings m
		       WHERE m.connector_id = c.id AND m.table_schema = ?
		   )`,
		legacyMySQLTargetSchema,
	)
	if err != nil {
		// 尚未建立資料表時（首次啟動）無須正規化。
		if strings.Contains(strings.ToLower(err.Error()), "no such table") {
			return nil
		}
		return fmt.Errorf("failed to inspect mysql connectors for migration %s: %w", migrationName, err)
	}
	defer rows.Close()

	type connectorSchema struct {
		id           string
		databaseName string
	}
	pending := make([]connectorSchema, 0)
	for rows.Next() {
		var id string
		var connectionConfig string
		if err := rows.Scan(&id, &connectionConfig); err != nil {
			return fmt.Errorf("failed to scan mysql connector for migration %s: %w", migrationName, err)
		}
		databaseName := mysqlDatabaseNameFromConnectionConfig(connectionConfig)
		if databaseName == "" || databaseName == legacyMySQLTargetSchema {
			// 找不到可靠的資料庫名稱時保持原值，讓錯誤在使用時明確浮現，
			// 而不是猜一個名字改寫使用者資料。
			continue
		}
		pending = append(pending, connectorSchema{id: id, databaseName: databaseName})
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("failed to iterate mysql connectors for migration %s: %w", migrationName, err)
	}

	for _, item := range pending {
		log.Printf("Executing SQLite migration: %s (connector %s -> %s)", migrationName, item.id, item.databaseName)
		if _, err := db.ExecContext(
			context.Background(),
			`UPDATE database_target_mappings SET table_schema = ? WHERE connector_id = ? AND table_schema = ?`,
			item.databaseName,
			item.id,
			legacyMySQLTargetSchema,
		); err != nil {
			return fmt.Errorf("failed to execute migration %s for connector %s: %w", migrationName, item.id, err)
		}
	}
	return nil
}

func mysqlDatabaseNameFromConnectionConfig(connectionConfig string) string {
	trimmed := strings.TrimSpace(connectionConfig)
	if trimmed == "" {
		return ""
	}
	config := map[string]any{}
	if err := json.Unmarshal([]byte(trimmed), &config); err != nil {
		return ""
	}
	for _, key := range []string{"database", "dbname"} {
		value, ok := config[key]
		if !ok || value == nil {
			continue
		}
		if name := strings.TrimSpace(fmt.Sprintf("%v", value)); name != "" {
			return name
		}
	}
	return ""
}
