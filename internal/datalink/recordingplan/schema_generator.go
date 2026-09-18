package recordingplan

import (
	"errors"
	"fmt"
	"strings"
)

// ErrUnsupportedSchemaDialect marks a target kind without verified managed DDL.
var ErrUnsupportedSchemaDialect = errors.New("unsupported database dialect")

// Dialect spellings accepted from connector kinds.
const (
	dialectSQLite     = "sqlite"
	dialectSQLite3    = "sqlite3"
	dialectPostgres   = "postgres"
	dialectPostgreSQL = "postgresql"
	dialectPgx        = "pgx"
)

// GenerateManagedSchemaDDL 根據目標資料庫種類與表名前綴產生完整的 managed 長表與索引 DDL。
func GenerateManagedSchemaDDL(dialect, prefix string) ([]string, error) {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		prefix = defaultManagedTablePrefix
	}

	d := strings.ToLower(strings.TrimSpace(dialect))
	switch d {
	case dialectSQLite, dialectSQLite3:
		return generateSQLiteDDL(prefix), nil
	case dialectPostgres, dialectPostgreSQL, dialectPgx:
		return generatePostgreSQLDDL(prefix), nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedSchemaDialect, dialect)
	}
}

func generateSQLiteDDL(p string) []string {
	return []string{
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %ssamples (
			workspace_id TEXT NOT NULL,
			plan_id TEXT NOT NULL,
			stream_id TEXT NOT NULL,
			record_id TEXT NOT NULL PRIMARY KEY,
			measurement_id TEXT NOT NULL,
			series_epoch TEXT NOT NULL,
			observed_at DATETIME NOT NULL,
			received_at DATETIME NOT NULL,
			quality TEXT NOT NULL,
			quality_reason TEXT,
			value_type TEXT NOT NULL,
			val_num REAL,
			val_dec TEXT,
			val_str TEXT,
			val_bool INTEGER,
			is_test INTEGER NOT NULL DEFAULT 0
		);`, p),
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS idx_%ssamples_plan_time ON %ssamples(workspace_id, plan_id, observed_at);`, p, p),
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS idx_%ssamples_meas_epoch ON %ssamples(measurement_id, series_epoch, observed_at);`, p, p),

		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %sintervals (
			workspace_id TEXT NOT NULL,
			plan_id TEXT NOT NULL,
			stream_id TEXT NOT NULL,
			record_id TEXT NOT NULL PRIMARY KEY,
			measurement_id TEXT NOT NULL,
			series_epoch TEXT NOT NULL,
			interval_start DATETIME NOT NULL,
			interval_end DATETIME NOT NULL,
			calculation_revision TEXT NOT NULL,
			sample_count INTEGER NOT NULL,
			mean_val REAL,
			min_val REAL,
			max_val REAL,
			quantity_delta REAL,
			known_subtotal REAL,
			is_estimated INTEGER NOT NULL DEFAULT 0,
			is_complete INTEGER NOT NULL DEFAULT 1,
			is_test INTEGER NOT NULL DEFAULT 0
		);`, p),
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS idx_%sintervals_time ON %sintervals(workspace_id, plan_id, interval_start, interval_end);`, p, p),

		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %sevents (
			workspace_id TEXT NOT NULL,
			plan_id TEXT NOT NULL,
			stream_id TEXT NOT NULL,
			record_id TEXT NOT NULL PRIMARY KEY,
			measurement_id TEXT NOT NULL,
			series_epoch TEXT NOT NULL,
			observed_at DATETIME NOT NULL,
			event_type TEXT NOT NULL,
			state_from TEXT,
			state_to TEXT,
			duration_ms INTEGER,
			message TEXT,
			is_test INTEGER NOT NULL DEFAULT 0
		);`, p),

		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %ssnapshots (
			workspace_id TEXT NOT NULL,
			plan_id TEXT NOT NULL,
			stream_id TEXT NOT NULL,
			record_id TEXT NOT NULL PRIMARY KEY,
			batch_id TEXT NOT NULL,
			trigger_id TEXT,
			observed_at DATETIME NOT NULL,
			completeness TEXT NOT NULL,
			payload_json TEXT NOT NULL,
			is_test INTEGER NOT NULL DEFAULT 0
		);`, p),

		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %sdefinitions (
			workspace_id TEXT NOT NULL,
			plan_id TEXT NOT NULL,
			measurement_id TEXT NOT NULL,
			definition_revision TEXT NOT NULL,
			series_epoch TEXT NOT NULL,
			quantity TEXT NOT NULL,
			unit TEXT,
			semantic_kind TEXT NOT NULL,
			schema_json TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (workspace_id, plan_id, measurement_id, definition_revision)
		);`, p),

		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %sreceipts (
			workspace_id TEXT NOT NULL,
			plan_id TEXT NOT NULL,
			destination_id TEXT NOT NULL,
			batch_id TEXT NOT NULL,
			last_record_id TEXT NOT NULL,
			delivered_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			record_count INTEGER NOT NULL,
			PRIMARY KEY (workspace_id, plan_id, destination_id, batch_id)
		);`, p),
	}
}

func generatePostgreSQLDDL(p string) []string {
	return []string{
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %ssamples (
			workspace_id VARCHAR(64) NOT NULL,
			plan_id VARCHAR(64) NOT NULL,
			stream_id VARCHAR(64) NOT NULL,
			record_id VARCHAR(64) NOT NULL PRIMARY KEY,
			measurement_id VARCHAR(64) NOT NULL,
			series_epoch VARCHAR(64) NOT NULL,
			observed_at TIMESTAMP WITH TIME ZONE NOT NULL,
			received_at TIMESTAMP WITH TIME ZONE NOT NULL,
			quality VARCHAR(32) NOT NULL,
			quality_reason TEXT,
			value_type VARCHAR(32) NOT NULL,
			val_num DOUBLE PRECISION,
			val_dec NUMERIC,
			val_str TEXT,
			val_bool BOOLEAN,
			is_test BOOLEAN NOT NULL DEFAULT FALSE
		);`, p),
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS idx_%ssamples_plan_time ON %ssamples(workspace_id, plan_id, observed_at);`, p, p),
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS idx_%ssamples_meas_epoch ON %ssamples(measurement_id, series_epoch, observed_at);`, p, p),

		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %sintervals (
			workspace_id VARCHAR(64) NOT NULL,
			plan_id VARCHAR(64) NOT NULL,
			stream_id VARCHAR(64) NOT NULL,
			record_id VARCHAR(64) NOT NULL PRIMARY KEY,
			measurement_id VARCHAR(64) NOT NULL,
			series_epoch VARCHAR(64) NOT NULL,
			interval_start TIMESTAMP WITH TIME ZONE NOT NULL,
			interval_end TIMESTAMP WITH TIME ZONE NOT NULL,
			calculation_revision VARCHAR(64) NOT NULL,
			sample_count INTEGER NOT NULL,
			mean_val DOUBLE PRECISION,
			min_val DOUBLE PRECISION,
			max_val DOUBLE PRECISION,
			quantity_delta DOUBLE PRECISION,
			known_subtotal DOUBLE PRECISION,
			is_estimated BOOLEAN NOT NULL DEFAULT FALSE,
			is_complete BOOLEAN NOT NULL DEFAULT TRUE,
			is_test BOOLEAN NOT NULL DEFAULT FALSE
		);`, p),
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS idx_%sintervals_time ON %sintervals(workspace_id, plan_id, interval_start, interval_end);`, p, p),

		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %sevents (
			workspace_id VARCHAR(64) NOT NULL,
			plan_id VARCHAR(64) NOT NULL,
			stream_id VARCHAR(64) NOT NULL,
			record_id VARCHAR(64) NOT NULL PRIMARY KEY,
			measurement_id VARCHAR(64) NOT NULL,
			series_epoch VARCHAR(64) NOT NULL,
			observed_at TIMESTAMP WITH TIME ZONE NOT NULL,
			event_type VARCHAR(64) NOT NULL,
			state_from VARCHAR(64),
			state_to VARCHAR(64),
			duration_ms BIGINT,
			message TEXT,
			is_test BOOLEAN NOT NULL DEFAULT FALSE
		);`, p),

		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %ssnapshots (
			workspace_id VARCHAR(64) NOT NULL,
			plan_id VARCHAR(64) NOT NULL,
			stream_id VARCHAR(64) NOT NULL,
			record_id VARCHAR(64) NOT NULL PRIMARY KEY,
			batch_id VARCHAR(64) NOT NULL,
			trigger_id VARCHAR(64),
			observed_at TIMESTAMP WITH TIME ZONE NOT NULL,
			completeness VARCHAR(32) NOT NULL,
			payload_json JSONB NOT NULL,
			is_test BOOLEAN NOT NULL DEFAULT FALSE
		);`, p),

		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %sdefinitions (
			workspace_id VARCHAR(64) NOT NULL,
			plan_id VARCHAR(64) NOT NULL,
			measurement_id VARCHAR(64) NOT NULL,
			definition_revision VARCHAR(64) NOT NULL,
			series_epoch VARCHAR(64) NOT NULL,
			quantity VARCHAR(64) NOT NULL,
			unit VARCHAR(32),
			semantic_kind VARCHAR(32) NOT NULL,
			schema_json JSONB NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (workspace_id, plan_id, measurement_id, definition_revision)
		);`, p),

		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %sreceipts (
			workspace_id VARCHAR(64) NOT NULL,
			plan_id VARCHAR(64) NOT NULL,
			destination_id VARCHAR(64) NOT NULL,
			batch_id VARCHAR(64) NOT NULL,
			last_record_id VARCHAR(64) NOT NULL,
			delivered_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
			record_count INTEGER NOT NULL,
			PRIMARY KEY (workspace_id, plan_id, destination_id, batch_id)
		);`, p),
	}
}
