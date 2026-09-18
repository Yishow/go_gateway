package dbtarget

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/schema"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
)

// TableInspectionStatus 描述單一資料表的實際查詢結果。
type TableInspectionStatus string

const (
	TableInspectionExists    TableInspectionStatus = "exists"
	TableInspectionMissing   TableInspectionStatus = "missing"
	TableInspectionForbidden TableInspectionStatus = "forbidden"
	TableInspectionFailed    TableInspectionStatus = "failed"
)

const (
	inspectionReasonPermissionDenied     = "permission_denied"
	inspectionReasonFailed               = "inspection_failed"
	inspectionReasonExistenceUnconfirmed = "existence_unconfirmed"
	inspectionReasonAuthenticationFailed = "authentication_failed"
	inspectionReasonConnectionFailed     = "connection_failed"
	inspectionReasonInvalidConfiguration = "invalid_configuration"
	inspectionReasonUnsupportedKind      = "unsupported_kind"
)

// TableInspection 是單一資料表的查詢結果；只有 exists 時 Columns 才是實際查得的欄位。
type TableInspection struct {
	Status  TableInspectionStatus `json:"inspection_status"`
	Schema  string                `json:"schema"`
	Table   string                `json:"table"`
	Columns []ColumnInfo          `json:"columns"`
	Reason  string                `json:"reason,omitempty"`
}

// InspectTable 查詢已存連線中的單一資料表。連線設定不存在時回錯誤；其餘結果一律
// 以狀態回報，連線失敗、權限不足或無法確認都不會被當成資料表不存在。
func (s *ConnectorService) InspectTable(ctx context.Context, connectorID, schemaName, tableName string) (*TableInspection, error) {
	connector, err := s.repo.GetByID(ctx, strings.TrimSpace(connectorID))
	if err != nil {
		return nil, fmt.Errorf("取得資料庫連接器失敗: %w", err)
	}
	result := &TableInspection{Schema: strings.TrimSpace(schemaName), Table: strings.TrimSpace(tableName), Columns: []ColumnInfo{}}
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

	db := manager.DB()
	switch connector.Kind {
	case schema.DatabaseConnectorKindSQLite:
		// sqlite_master 列出所有資料表，沒列出即為不存在。
		result.Schema = defaultString(result.Schema, "main")
		tables, err := inspectSQLiteTables(ctx, db)
		if err != nil {
			return result.classified(err), nil
		}
		return result.fromListing(tables, TableInspectionMissing, ""), nil
	case schema.DatabaseConnectorKindPostgres:
		result.Schema = defaultString(result.Schema, "public")
		exists, err := postgresCatalogHasTable(ctx, db, result.Schema, result.Table)
		if err != nil {
			return result.classified(err), nil
		}
		if !exists {
			return result.withStatus(TableInspectionMissing, ""), nil
		}
		tables, err := inspectPostgresTables(ctx, db)
		if err != nil {
			return result.classified(err), nil
		}
		// 系統目錄有這張表，權限過濾的欄位清單卻看不到：帳號無權查詢。
		return result.fromListing(tables, TableInspectionForbidden, inspectionReasonPermissionDenied), nil
	case schema.DatabaseConnectorKindMySQL:
		result.Schema = defaultString(result.Schema, defaultString(stringConfigValue(config, "database"), stringConfigValue(config, "dbname")))
		tables, err := inspectMySQLTables(ctx, db, []string{result.Schema})
		if err != nil {
			return result.classified(err), nil
		}
		// information_schema 只列出帳號有權限的表，沒列出無法確認資料表不存在。
		return result.fromListing(tables, TableInspectionFailed, inspectionReasonExistenceUnconfirmed), nil
	default:
		return result.withStatus(TableInspectionFailed, inspectionReasonUnsupportedKind), nil
	}
}

func (r *TableInspection) withStatus(status TableInspectionStatus, reason string) *TableInspection {
	r.Status = status
	r.Reason = reason
	r.Columns = []ColumnInfo{}
	return r
}

func (r *TableInspection) classified(err error) *TableInspection {
	status, reason := classifyInspectionError(err)
	return r.withStatus(status, reason)
}

func (r *TableInspection) fromListing(tables []TableInfo, notFound TableInspectionStatus, notFoundReason string) *TableInspection {
	for _, table := range tables {
		if table.Schema == r.Schema && table.Name == r.Table {
			r.Status = TableInspectionExists
			r.Reason = ""
			r.Columns = append([]ColumnInfo{}, table.Columns...)
			return r
		}
	}
	return r.withStatus(notFound, notFoundReason)
}

// postgresCatalogHasTable 以不受表權限過濾的 pg_catalog 判斷資料表是否存在。
func postgresCatalogHasTable(ctx context.Context, db *sql.DB, schemaName, tableName string) (bool, error) {
	var found bool
	err := db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM pg_catalog.pg_class c
			JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
			WHERE c.relkind IN ('r', 'p')
			  AND n.nspname = $1
			  AND c.relname = $2
		)
	`, schemaName, tableName).Scan(&found)
	if err != nil {
		return false, fmt.Errorf("查詢 postgres 系統目錄失敗: %w", err)
	}
	return found, nil
}

// classifyInspectionError 把權限不足與其他查詢失敗分開，兩者都不代表資料表不存在。
func classifyInspectionError(err error) (status TableInspectionStatus, reason string) {
	if isPermissionDeniedError(err) {
		return TableInspectionForbidden, inspectionReasonPermissionDenied
	}
	return TableInspectionFailed, inspectionReasonFailed
}

func isPermissionDeniedError(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "42501" {
		return true
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "42501" {
		return true
	}
	var mysqlErr *mysqldriver.MySQLError
	if errors.As(err, &mysqlErr) {
		switch mysqlErr.Number {
		case 1044, 1142, 1143, 1227:
			return true
		}
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "permission denied") || strings.Contains(message, "command denied")
}
