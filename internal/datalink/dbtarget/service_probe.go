package dbtarget

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
	"go-gateway/internal/datalink/schema"
)

var (
	openExternalDBManagerFunc           = openExternalDBManager
	ensurePostgresDatabaseIfMissingFunc = ensurePostgresDatabaseIfMissing
)

func probeConnector(ctx context.Context, kind schema.DatabaseConnectorKind, config ConnectionConfig) (schema.DatabaseConnectorStatus, *time.Time, string) {
	now := time.Now()
	manager, err := openExternalDBManagerFunc(kind, config)
	if err != nil && kind == schema.DatabaseConnectorKindPostgres {
		if createErr := ensurePostgresDatabaseIfMissingFunc(ctx, config, err); createErr == nil {
			manager, err = openExternalDBManagerFunc(kind, config)
		} else if isPostgresMissingDatabaseError(err) {
			err = createErr
		}
	}
	if err != nil {
		return classifyConnectorError(err), &now, connectorErrorMessage(kind, err)
	}
	defer manager.Close()

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := manager.DB().PingContext(pingCtx); err != nil {
		return classifyConnectorError(err), &now, connectorErrorMessage(kind, err)
	}

	return schema.DatabaseConnectorStatusReady, &now, ""
}

func ensurePostgresDatabaseIfMissing(ctx context.Context, config ConnectionConfig, connectErr error) error {
	if !isPostgresMissingDatabaseError(connectErr) {
		return connectErr
	}

	dbConfig, err := buildExternalDBConfig(schema.DatabaseConnectorKindPostgres, config)
	if err != nil {
		return err
	}

	connConfig, err := pgx.ParseConfig(dbConfig.DSN)
	if err != nil {
		return fmt.Errorf("解析 postgres 連接設定失敗: %w", err)
	}
	databaseName := strings.TrimSpace(connConfig.Database)
	if databaseName == "" {
		return fmt.Errorf("postgres 連接設定缺少 database")
	}

	adminConfig := connConfig.Copy()
	adminConfig.Database = postgresAdminDatabaseName(databaseName)
	connectCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	adminConn, err := pgx.ConnectConfig(connectCtx, adminConfig)
	if err != nil {
		return fmt.Errorf("建立 postgres 管理連線失敗: %w", err)
	}
	defer adminConn.Close(ctx)

	if _, err := adminConn.Exec(
		connectCtx,
		fmt.Sprintf("CREATE DATABASE %s", quoteIdentifier(schema.DatabaseConnectorKindPostgres, databaseName)),
	); err != nil {
		if isPostgresDuplicateDatabaseError(err) {
			return nil
		}
		return fmt.Errorf("建立 postgres 資料庫失敗: %w", err)
	}

	return nil
}

func postgresAdminDatabaseName(databaseName string) string {
	if strings.EqualFold(strings.TrimSpace(databaseName), "postgres") {
		return "template1"
	}
	return "postgres"
}

func isPostgresMissingDatabaseError(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return string(pqErr.Code) == "3D000"
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "3D000"
	}

	message := strings.ToLower(err.Error())
	return strings.Contains(message, "database") && strings.Contains(message, "does not exist")
}

func isPostgresDuplicateDatabaseError(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return string(pqErr.Code) == "42P04"
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "42P04"
	}

	message := strings.ToLower(err.Error())
	return strings.Contains(message, "database") && strings.Contains(message, "already exists")
}
