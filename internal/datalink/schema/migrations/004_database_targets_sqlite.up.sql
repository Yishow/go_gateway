CREATE TABLE IF NOT EXISTS database_connectors (
    id                TEXT PRIMARY KEY,
    name              TEXT NOT NULL UNIQUE,
    kind              TEXT NOT NULL,
    connection_config TEXT NOT NULL,
    status            TEXT NOT NULL,
    last_check_at     DATETIME,
    last_check_error  TEXT NOT NULL DEFAULT '',
    enabled           BOOLEAN NOT NULL DEFAULT 1,
    created_at        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS database_target_mappings (
    id                TEXT PRIMARY KEY,
    tag_id            TEXT NOT NULL,
    connector_id      TEXT NOT NULL,
    table_schema      TEXT NOT NULL DEFAULT 'main',
    table_name        TEXT NOT NULL,
    column_name       TEXT NOT NULL,
    write_mode        TEXT NOT NULL DEFAULT 'insert',
    timestamp_column  TEXT,
    enabled           BOOLEAN NOT NULL DEFAULT 1,
    created_at        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(tag_id, connector_id)
);

CREATE INDEX IF NOT EXISTS idx_database_target_mappings_connector_id
    ON database_target_mappings(connector_id);

CREATE INDEX IF NOT EXISTS idx_database_target_mappings_tag_id
    ON database_target_mappings(tag_id);

INSERT OR IGNORE INTO system_settings (key, value, description) VALUES
    ('runtime_stream_enabled', 'true', '是否啟用 runtime SSE / live value 功能'),
    ('db_target_enabled', 'true', '是否啟用資料庫目標輸出功能');

CREATE TRIGGER IF NOT EXISTS trigger_database_connectors_updated_at
    AFTER UPDATE ON database_connectors
    FOR EACH ROW
BEGIN
    UPDATE database_connectors SET updated_at = datetime('now') WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS trigger_database_target_mappings_updated_at
    AFTER UPDATE ON database_target_mappings
    FOR EACH ROW
BEGIN
    UPDATE database_target_mappings SET updated_at = datetime('now') WHERE id = NEW.id;
END;
