-- =============================================================================
-- Go Gateway - 資料收集管線資料庫遷移腳本 (SQLite)
-- 版本: 001_initial_schema
-- 描述: 建立設備、點位、標籤、映射、時序儲存的完整資料庫結構
-- 注意: SQLite 不支援表分區，使用單一表結構
-- =============================================================================

-- -----------------------------------------------------------------------------
-- 設備表 (devices)
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS devices (
    id              TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab', 1 + (abs(random()) % 4), 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
    name            TEXT NOT NULL,
    description     TEXT,
    protocol        TEXT NOT NULL CHECK (protocol IN ('modbus_tcp', 'modbus_rtu', 'modbus_udp', 'fatek_fbs', 'mc_3e', 'mqtt')),
    status          TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'disabled')),
    connection_config TEXT NOT NULL DEFAULT '{}',
    last_test_at    TEXT,
    last_test_success INTEGER,
    last_test_error TEXT,
    created_at      TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at      TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_devices_protocol ON devices(protocol);
CREATE INDEX IF NOT EXISTS idx_devices_status ON devices(status);
CREATE INDEX IF NOT EXISTS idx_devices_name ON devices(name);

-- -----------------------------------------------------------------------------
-- 輪詢群組表 (polling_groups)
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS polling_groups (
    id              TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab', 1 + (abs(random()) % 4), 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
    name            TEXT NOT NULL UNIQUE,
    description     TEXT,
    interval_ms     INTEGER NOT NULL DEFAULT 1000 CHECK (interval_ms >= 100),
    priority        INTEGER NOT NULL DEFAULT 100,
    enabled         INTEGER NOT NULL DEFAULT 1,
    created_at      TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at      TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_polling_groups_enabled ON polling_groups(enabled);
CREATE INDEX IF NOT EXISTS idx_polling_groups_priority ON polling_groups(priority);

-- -----------------------------------------------------------------------------
-- 點位表 (points)
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS points (
    id              TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab', 1 + (abs(random()) % 4), 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
    device_id       TEXT NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    description     TEXT,
    address         TEXT NOT NULL,
    function        TEXT NOT NULL DEFAULT '',
    data_type       TEXT NOT NULL CHECK (data_type IN ('bool', 'int16', 'uint16', 'int32', 'uint32', 'int64', 'uint64', 'float32', 'float64', 'string')),
    mode            TEXT NOT NULL DEFAULT 'read' CHECK (mode IN ('read', 'readwrite')),
    polling_group_id TEXT REFERENCES polling_groups(id) ON DELETE SET NULL,
    last_read_at    TEXT,
    last_value      TEXT,
    last_error      TEXT,
    enabled         INTEGER NOT NULL DEFAULT 1,
    created_at      TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at      TEXT NOT NULL DEFAULT (datetime('now')),
    UNIQUE (device_id, address, function)
);

CREATE INDEX IF NOT EXISTS idx_points_device_id ON points(device_id);
CREATE INDEX IF NOT EXISTS idx_points_polling_group_id ON points(polling_group_id);
CREATE INDEX IF NOT EXISTS idx_points_enabled ON points(enabled);
CREATE INDEX IF NOT EXISTS idx_points_data_type ON points(data_type);

-- -----------------------------------------------------------------------------
-- 標籤表 (tags)
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS tags (
    id              TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab', 1 + (abs(random()) % 4), 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
    key             TEXT NOT NULL,
    key_lower       TEXT NOT NULL UNIQUE,
    display_name    TEXT NOT NULL,
    description     TEXT,
    unit            TEXT,
    data_type       TEXT NOT NULL CHECK (data_type IN ('bool', 'int16', 'uint16', 'int32', 'uint32', 'int64', 'uint64', 'float32', 'float64', 'string')),
    status          TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'retired')),
    labels          TEXT DEFAULT '{}',
    created_at      TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at      TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_tags_key ON tags(key);
CREATE INDEX IF NOT EXISTS idx_tags_key_lower ON tags(key_lower);
CREATE INDEX IF NOT EXISTS idx_tags_status ON tags(status);
CREATE INDEX IF NOT EXISTS idx_tags_data_type ON tags(data_type);

-- -----------------------------------------------------------------------------
-- 映射表 (mappings)
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS mappings (
    id                  TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab', 1 + (abs(random()) % 4), 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
    point_id            TEXT NOT NULL REFERENCES points(id) ON DELETE CASCADE,
    tag_id              TEXT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    transform_pipeline  TEXT NOT NULL DEFAULT '[]',
    status              TEXT NOT NULL DEFAULT 'active',
    rule_candidate_id   TEXT,
    proposed_signature  TEXT,
    last_applied_signature TEXT,
    blocking_reason     TEXT,
    enabled             INTEGER NOT NULL DEFAULT 1,
    created_at          TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at          TEXT NOT NULL DEFAULT (datetime('now')),
    UNIQUE (point_id, tag_id)
);

CREATE INDEX IF NOT EXISTS idx_mappings_point_id ON mappings(point_id);
CREATE INDEX IF NOT EXISTS idx_mappings_tag_id ON mappings(tag_id);
CREATE INDEX IF NOT EXISTS idx_mappings_enabled ON mappings(enabled);
CREATE INDEX IF NOT EXISTS idx_mappings_rule_candidate_id ON mappings(rule_candidate_id);

-- -----------------------------------------------------------------------------
-- 時序資料表 (timeseries) - 非分區
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS timeseries (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    tag_id          TEXT NOT NULL,
    ts              TEXT NOT NULL,
    value_num       REAL,
    value_text      TEXT,
    value_bool      INTEGER,
    raw_value       TEXT,
    quality         TEXT NOT NULL DEFAULT 'good' CHECK (quality IN ('good', 'bad', 'uncertain'))
);

CREATE INDEX IF NOT EXISTS idx_timeseries_tag_ts ON timeseries(tag_id, ts DESC);
CREATE INDEX IF NOT EXISTS idx_timeseries_ts ON timeseries(ts DESC);
CREATE INDEX IF NOT EXISTS idx_timeseries_quality ON timeseries(quality) WHERE quality != 'good';

-- -----------------------------------------------------------------------------
-- 系統設定表 (system_settings)
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS system_settings (
    id              TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab', 1 + (abs(random()) % 4), 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
    key             TEXT NOT NULL UNIQUE,
    value           TEXT NOT NULL,
    description     TEXT,
    updated_at      TEXT NOT NULL DEFAULT (datetime('now'))
);

-- 插入預設系統設定
INSERT OR IGNORE INTO system_settings (key, value, description) VALUES
    ('write_precision', '"second"', '時序寫入時間精度 (second 或 millisecond)'),
    ('partition_interval', '"monthly"', '時序表分區間隔 (SQLite 不適用)'),
    ('batch_size', '1000', '批次寫入大小'),
    ('flush_interval', '5000', '批次刷新間隔 (毫秒)'),
    ('default_retry_count', '3', '預設重試次數'),
    ('default_retry_delay', '1000', '預設重試延遲 (毫秒)'),
    ('modbus_share', '{"enabled":false,"bind_address":"127.0.0.1","port":5020,"slave_id":1,"capacity_registers":32768,"settings_revision":"migration-018"}', 'Local Modbus Share listener settings');

-- -----------------------------------------------------------------------------
-- 自動更新 updated_at 觸發器
-- -----------------------------------------------------------------------------
CREATE TRIGGER IF NOT EXISTS trigger_devices_updated_at
    AFTER UPDATE ON devices
    FOR EACH ROW
BEGIN
    UPDATE devices SET updated_at = datetime('now') WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS trigger_polling_groups_updated_at
    AFTER UPDATE ON polling_groups
    FOR EACH ROW
BEGIN
    UPDATE polling_groups SET updated_at = datetime('now') WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS trigger_points_updated_at
    AFTER UPDATE ON points
    FOR EACH ROW
BEGIN
    UPDATE points SET updated_at = datetime('now') WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS trigger_tags_updated_at
    AFTER UPDATE ON tags
    FOR EACH ROW
BEGIN
    UPDATE tags SET updated_at = datetime('now') WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS trigger_mappings_updated_at
    AFTER UPDATE ON mappings
    FOR EACH ROW
BEGIN
    UPDATE mappings SET updated_at = datetime('now') WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS trigger_system_settings_updated_at
    AFTER UPDATE ON system_settings
    FOR EACH ROW
BEGIN
    UPDATE system_settings SET updated_at = datetime('now') WHERE id = NEW.id;
END;

-- -----------------------------------------------------------------------------
-- 標籤鍵小寫自動填充觸發器
-- -----------------------------------------------------------------------------
CREATE TRIGGER IF NOT EXISTS trigger_tags_key_lower_insert
    AFTER INSERT ON tags
    FOR EACH ROW
    WHEN NEW.key_lower IS NULL OR NEW.key_lower = ''
BEGIN
    UPDATE tags SET key_lower = lower(NEW.key) WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS trigger_tags_key_lower_update
    AFTER UPDATE OF key ON tags
    FOR EACH ROW
BEGIN
    UPDATE tags SET key_lower = lower(NEW.key) WHERE id = NEW.id;
END;
