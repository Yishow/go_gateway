package dbtarget

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
)

const connectorSelectColumns = `
	id, name, kind, connection_config, status, COALESCE(CAST(last_check_at AS TEXT), ''),
	last_check_error, enabled, default_write_interval_seconds,
	COALESCE(CAST(last_schema_ensure_at AS TEXT), ''), last_schema_ensure_status, last_schema_ensure_error,
	COALESCE(CAST(last_write_at AS TEXT), ''), last_write_status, last_write_error,
	COALESCE(CAST(last_flush_at AS TEXT), ''), last_flush_status, last_flush_error,
	CAST(created_at AS TEXT), CAST(updated_at AS TEXT)`

func scanConnectorRow(row *sql.Row) (*schema.DatabaseConnector, error) {
	var connector schema.DatabaseConnector
	lastCheckAt, lastSchemaAt, lastWriteAt, lastFlushAt, createdAt, updatedAt := "", "", "", "", "", ""

	err := scanConnector(row.Scan, &connector, &lastCheckAt, &lastSchemaAt, &lastWriteAt, &lastFlushAt, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("資料庫連接器不存在")
	}
	if err != nil {
		return nil, fmt.Errorf("掃描資料庫連接器失敗: %w", err)
	}
	if err := parseConnectorTimes(&connector, lastCheckAt, lastSchemaAt, lastWriteAt, lastFlushAt, createdAt, updatedAt); err != nil {
		return nil, err
	}
	return &connector, nil
}

func scanConnectorRows(rows *sql.Rows) ([]*schema.DatabaseConnector, error) {
	var connectors []*schema.DatabaseConnector

	for rows.Next() {
		var connector schema.DatabaseConnector
		lastCheckAt, lastSchemaAt, lastWriteAt, lastFlushAt, createdAt, updatedAt := "", "", "", "", "", ""
		if err := scanConnector(rows.Scan, &connector, &lastCheckAt, &lastSchemaAt, &lastWriteAt, &lastFlushAt, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("掃描資料庫連接器失敗: %w", err)
		}
		if err := parseConnectorTimes(&connector, lastCheckAt, lastSchemaAt, lastWriteAt, lastFlushAt, createdAt, updatedAt); err != nil {
			return nil, err
		}
		connectors = append(connectors, &connector)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍歷資料庫連接器失敗: %w", err)
	}
	return connectors, nil
}

func scanConnector(
	scan func(dest ...any) error,
	connector *schema.DatabaseConnector,
	lastCheckAt, lastSchemaAt, lastWriteAt, lastFlushAt, createdAt, updatedAt *string,
) error {
	return scan(
		&connector.ID,
		&connector.Name,
		&connector.Kind,
		&connector.ConnectionConfig,
		&connector.Status,
		lastCheckAt,
		&connector.LastCheckError,
		&connector.Enabled,
		&connector.DefaultWriteIntervalSeconds,
		lastSchemaAt,
		&connector.LastSchemaEnsureStatus,
		&connector.LastSchemaEnsureError,
		lastWriteAt,
		&connector.LastWriteStatus,
		&connector.LastWriteError,
		lastFlushAt,
		&connector.LastFlushStatus,
		&connector.LastFlushError,
		createdAt,
		updatedAt,
	)
}

func parseConnectorTimes(
	connector *schema.DatabaseConnector,
	lastCheckAt, lastSchemaAt, lastWriteAt, lastFlushAt, createdAt, updatedAt string,
) error {
	parsedCreatedAt, err := common.ParseTimeString(createdAt)
	if err != nil {
		return fmt.Errorf("解析資料庫連接器建立時間失敗: %w", err)
	}
	parsedUpdatedAt, err := common.ParseTimeString(updatedAt)
	if err != nil {
		return fmt.Errorf("解析資料庫連接器更新時間失敗: %w", err)
	}
	connector.CreatedAt = parsedCreatedAt
	connector.UpdatedAt = parsedUpdatedAt

	if err := setOptionalConnectorTime(lastCheckAt, &connector.LastCheckAt, "解析資料庫連接器檢查時間失敗"); err != nil {
		return err
	}
	if err := setOptionalConnectorTime(lastSchemaAt, &connector.LastSchemaEnsureAt, "解析資料庫連接器 schema ensure 時間失敗"); err != nil {
		return err
	}
	if err := setOptionalConnectorTime(lastWriteAt, &connector.LastWriteAt, "解析資料庫連接器寫入時間失敗"); err != nil {
		return err
	}
	return setOptionalConnectorTime(lastFlushAt, &connector.LastFlushAt, "解析資料庫連接器 flush 時間失敗")
}

func setOptionalConnectorTime(raw string, target **time.Time, message string) error {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parsed, err := common.ParseTimeString(raw)
	if err != nil {
		return fmt.Errorf("%s: %w", message, err)
	}
	*target = &parsed
	return nil
}
