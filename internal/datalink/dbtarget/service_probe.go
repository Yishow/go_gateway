package dbtarget

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"
)

var (
	openExternalDBManagerFunc           = openExternalDBManager
	ensurePostgresDatabaseIfMissingFunc = ensurePostgresDatabaseIfMissing
	ensureMySQLDatabaseIfMissingFunc    = ensureMySQLDatabaseIfMissing
)

// openExternalDBManagerWithMySQLDatabaseEnsure opens a target database and, for
// MySQL only, creates a missing database before retrying the open once.
func openExternalDBManagerWithMySQLDatabaseEnsure(
	ctx context.Context,
	kind schema.DatabaseConnectorKind,
	config ConnectionConfig,
) (*datalinkbase.DBManager, error) {
	manager, err := openExternalDBManagerFunc(kind, config)
	if err == nil || kind != schema.DatabaseConnectorKindMySQL || !isMySQLMissingDatabaseError(err) {
		return manager, err
	}

	if createErr := ensureMySQLDatabaseIfMissingFunc(ctx, config, err); createErr != nil {
		return nil, createErr
	}
	return openExternalDBManagerFunc(kind, config)
}

func probeConnector(ctx context.Context, kind schema.DatabaseConnectorKind, config ConnectionConfig) (schema.DatabaseConnectorStatus, *time.Time, string) {
	now := time.Now()
	var manager *datalinkbase.DBManager
	var err error
	if kind == schema.DatabaseConnectorKindMySQL {
		manager, err = openExternalDBManagerWithMySQLDatabaseEnsure(ctx, kind, config)
	} else {
		manager, err = openExternalDBManagerFunc(kind, config)
	}
	// MySQL 的建庫重試已由 openExternalDBManagerWithMySQLDatabaseEnsure 內建處理，
	// 這裡只補上 PostgreSQL 尚未包裝的同等流程。
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
	if err == nil {
		return false
	}
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
	if err == nil {
		return false
	}
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

func ensureMySQLDatabaseIfMissing(ctx context.Context, config ConnectionConfig, connectErr error) error {
	if !isMySQLMissingDatabaseError(connectErr) {
		return connectErr
	}

	user := defaultString(
		stringConfigValue(config, "user"),
		stringConfigValue(config, "username"),
	)
	databaseName := defaultString(
		stringConfigValue(config, "database"),
		stringConfigValue(config, "dbname"),
	)
	if user == "" || databaseName == "" {
		// 以 dsn 設定的連接器沒有這些離散鍵，這裡不是設定缺漏，而是自動建庫
		// 無法進行；回傳原始連線錯誤，避免用「缺少 user/database」蓋掉真正原因。
		return connectErr
	}

	mysqlConfig := newMySQLDriverConfig(config, "")

	if timeout := strings.TrimSpace(stringConfigValue(config, "timeout")); timeout != "" {
		if duration, err := time.ParseDuration(timeout); err == nil {
			mysqlConfig.Timeout = duration
			mysqlConfig.ReadTimeout = duration
			mysqlConfig.WriteTimeout = duration
		}
	}

	adminDB, err := sql.Open("mysql", mysqlConfig.FormatDSN())
	if err != nil {
		return fmt.Errorf("建立 mysql 管理連線失敗: %w", err)
	}
	defer adminDB.Close()

	connectCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	createSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", strings.ReplaceAll(databaseName, "`", "``"))
	if _, err := adminDB.ExecContext(connectCtx, createSQL); err != nil {
		return fmt.Errorf("建立 mysql 資料庫失敗: %w", err)
	}

	return nil
}

func isMySQLMissingDatabaseError(err error) bool {
	if err == nil {
		return false
	}
	var myErr *mysqldriver.MySQLError
	if errors.As(err, &myErr) {
		return myErr.Number == 1049
	}

	message := strings.ToLower(err.Error())
	return strings.Contains(message, "error 1049") || (strings.Contains(message, "unknown database") && strings.Contains(message, "42000"))
}

// ensureExternalDatabaseExists 是明確的佈建進入點：在需要時建立缺漏的 MySQL
// 目標資料庫後立即關閉連線。唯讀檢查與執行期寫入路徑不得呼叫它。
func ensureExternalDatabaseExists(
	ctx context.Context,
	kind schema.DatabaseConnectorKind,
	config ConnectionConfig,
) error {
	if kind != schema.DatabaseConnectorKindMySQL {
		// 其他資料庫類型沒有對應的自動建庫流程，連線錯誤留給後續步驟回報。
		return nil
	}
	manager, err := openExternalDBManagerWithMySQLDatabaseEnsure(ctx, kind, config)
	if err != nil {
		return err
	}
	return manager.Close()
}
