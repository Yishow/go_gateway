// Package datalink 提供 Datalink 系統的資料庫連接管理。
//
// 本套件實作統一的資料庫連接工廠，支援 SQLite、PostgreSQL 和 MySQL。
package datalink

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	_ "github.com/lib/pq"              // PostgreSQL driver
	_ "modernc.org/sqlite"             // SQLite driver (pure Go)
)

// =============================================================================
// 資料庫類型定義
// =============================================================================

// DBType 資料庫類型
type DBType string

const (
	// DBTypeSQLite SQLite 資料庫
	DBTypeSQLite DBType = "sqlite"
	// DBTypePostgres PostgreSQL 資料庫
	DBTypePostgres DBType = "postgres"
	// DBTypeMySQL MySQL 資料庫
	DBTypeMySQL DBType = "mysql"
)

// =============================================================================
// 資料庫連接配置
// =============================================================================

// DBConfig 資料庫連接配置
type DBConfig struct {
	// Type 資料庫類型
	Type DBType `json:"type" yaml:"type"`
	// DSN 資料源名稱 (連接字串)
	DSN string `json:"dsn" yaml:"dsn"`
	// MaxOpenConns 最大開啟連接數
	MaxOpenConns int `json:"max_open_conns" yaml:"max_open_conns"`
	// MaxIdleConns 最大閒置連接數
	MaxIdleConns int `json:"max_idle_conns" yaml:"max_idle_conns"`
}

// DefaultEmbeddedSQLiteDSN 為內嵌閘道（如 cmd/test_ui）使用的 SQLite URI。
// - busy_timeout：鎖競爭時等待（毫秒），降低 SQLITE_BUSY。
// - journal_mode=WAL：讀寫並發較 DELETE 模式友善。
// - cache=shared、mode=rwc：與 database/sql 連線池共用頁面快取並允許讀寫。
const DefaultEmbeddedSQLiteDSN = "file:datalink.db?cache=shared&mode=rwc&_pragma=busy_timeout(15000)&_pragma=journal_mode(WAL)"

// ApplySQLitePoolDefaults 將 *sql.DB 設為單連線池設定，與單檔 SQLite 鎖行為相容。
// 在仍使用一個 *sql.DB 實例的前提下，可明顯減少並發連線造成的 database is locked。
func ApplySQLitePoolDefaults(db *sql.DB) {
	if db == nil {
		return
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)
}

// DefaultSQLiteConfig 預設 SQLite 配置
func DefaultSQLiteConfig() DBConfig {
	return DBConfig{
		Type:         DBTypeSQLite,
		DSN:          DefaultEmbeddedSQLiteDSN,
		MaxOpenConns: 1,
		MaxIdleConns: 1,
	}
}

// DefaultPostgresConfig 預設 PostgreSQL 配置
func DefaultPostgresConfig(host, port, user, password, dbname string) DBConfig {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	return DBConfig{
		Type:         DBTypePostgres,
		DSN:          dsn,
		MaxOpenConns: 10,
		MaxIdleConns: 5,
	}
}

// =============================================================================
// 資料庫連接管理
// =============================================================================

// DBManager 資料庫連接管理器
type DBManager struct {
	config DBConfig
	db     *sql.DB
	mu     sync.RWMutex
}

// NewDBManager 建立新的資料庫管理器
func NewDBManager(config DBConfig) *DBManager {
	return &DBManager{
		config: config,
	}
}

// Connect 建立資料庫連接
func (m *DBManager) Connect() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.db != nil {
		return nil // 已連接
	}

	var driverName string
	switch m.config.Type {
	case DBTypeSQLite:
		driverName = "sqlite"
	case DBTypePostgres:
		driverName = "postgres"
	case DBTypeMySQL:
		driverName = "mysql"
	default:
		return fmt.Errorf("不支援的資料庫類型: %s", m.config.Type)
	}

	db, err := sql.Open(driverName, m.config.DSN)
	if err != nil {
		return fmt.Errorf("開啟資料庫連接失敗: %w", err)
	}

	// 設定連接池
	if m.config.MaxOpenConns > 0 {
		db.SetMaxOpenConns(m.config.MaxOpenConns)
	}
	if m.config.MaxIdleConns > 0 {
		db.SetMaxIdleConns(m.config.MaxIdleConns)
	}

	// 測試連接
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()
	if err := db.PingContext(pingCtx); err != nil {
		db.Close()
		return fmt.Errorf("資料庫連接測試失敗: %w", err)
	}

	m.db = db
	return nil
}

// Close 關閉資料庫連接
func (m *DBManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.db == nil {
		return nil
	}

	err := m.db.Close()
	m.db = nil
	return err
}

// DB 取得資料庫連接
func (m *DBManager) DB() *sql.DB {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.db
}

// Type 取得資料庫類型
func (m *DBManager) Type() DBType {
	return m.config.Type
}

// =============================================================================
// Repository 工廠
// =============================================================================

// RepositoryFactory Repository 工廠介面
type RepositoryFactory interface {
	// DeviceRepository 取得 Device Repository
	DeviceRepository() interface{}
	// PointRepository 取得 Point Repository
	PointRepository() interface{}
	// TagRepository 取得 Tag Repository
	TagRepository() interface{}
	// MappingRepository 取得 Mapping Repository
	MappingRepository() interface{}
	// PollingGroupRepository 取得 PollingGroup Repository
	PollingGroupRepository() interface{}
	// SettingsRepository 取得 Settings Repository
	SettingsRepository() interface{}
}

// SQLRepositoryFactory SQL Repository 工廠
type SQLRepositoryFactory struct {
	dbManager *DBManager
}

// NewSQLRepositoryFactory 建立新的 SQL Repository 工廠
func NewSQLRepositoryFactory(dbManager *DBManager) *SQLRepositoryFactory {
	return &SQLRepositoryFactory{
		dbManager: dbManager,
	}
}

// DBManager 取得資料庫管理器
func (f *SQLRepositoryFactory) DBManager() *DBManager {
	return f.dbManager
}

// placeholder 方法，實際實現在各自的 sql_repo.go 中
func (f *SQLRepositoryFactory) DeviceRepository() interface{}       { return nil }
func (f *SQLRepositoryFactory) PointRepository() interface{}        { return nil }
func (f *SQLRepositoryFactory) TagRepository() interface{}          { return nil }
func (f *SQLRepositoryFactory) MappingRepository() interface{}      { return nil }
func (f *SQLRepositoryFactory) PollingGroupRepository() interface{} { return nil }
func (f *SQLRepositoryFactory) SettingsRepository() interface{}     { return nil }
