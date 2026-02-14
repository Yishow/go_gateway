// Package point 提供點位管理功能的 SQL Repository 實現。
package point

import "database/sql"

// =============================================================================
// SQL Repository 實現
// =============================================================================

// SQLRepository SQL 點位儲存庫
type SQLRepository struct {
	db *sql.DB
}

// NewSQLRepository 建立新的 SQL 儲存庫
func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}
