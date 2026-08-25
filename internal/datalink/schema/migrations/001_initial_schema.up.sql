-- =============================================================================
-- Go Gateway - 資料收集管線資料庫遷移腳本 (PostgreSQL)
-- 版本: 001_initial_schema
-- 描述: 建立設備、點位、標籤、映射、時序儲存的完整資料庫結構
-- =============================================================================

-- -----------------------------------------------------------------------------
-- 啟用必要擴充
-- -----------------------------------------------------------------------------
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- -----------------------------------------------------------------------------
-- 設備表 (devices)
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS devices (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    protocol        VARCHAR(50) NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'draft',
    connection_config JSONB NOT NULL DEFAULT '{}',
    last_test_at    TIMESTAMPTZ,
    last_test_success BOOLEAN,
    last_test_error TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT devices_protocol_check CHECK (
        protocol IN ('modbus_tcp', 'modbus_rtu', 'modbus_udp', 'fatek_fbs', 'mc_3e', 'mqtt')
    ),
    CONSTRAINT devices_status_check CHECK (
        status IN ('draft', 'active', 'disabled')
    )
);

CREATE INDEX idx_devices_protocol ON devices(protocol);
CREATE INDEX idx_devices_status ON devices(status);
CREATE INDEX idx_devices_name ON devices(name);

COMMENT ON TABLE devices IS '設備註冊表，儲存可連接的 PLC 和資料來源';
COMMENT ON COLUMN devices.protocol IS '通訊協議類型';
COMMENT ON COLUMN devices.connection_config IS '協議特定的連線配置 (JSON)';

-- -----------------------------------------------------------------------------
-- 輪詢群組表 (polling_groups)
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS polling_groups (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name            VARCHAR(255) NOT NULL UNIQUE,
    description     TEXT,
    interval_ms     INTEGER NOT NULL DEFAULT 1000,
    priority        INTEGER NOT NULL DEFAULT 100,
    enabled         BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT polling_groups_interval_check CHECK (interval_ms >= 100)
);

CREATE INDEX idx_polling_groups_enabled ON polling_groups(enabled);
CREATE INDEX idx_polling_groups_priority ON polling_groups(priority);

COMMENT ON TABLE polling_groups IS '輪詢群組，控制點位的採集間隔';
COMMENT ON COLUMN polling_groups.interval_ms IS '輪詢間隔 (毫秒)，最小 100ms';

-- -----------------------------------------------------------------------------
-- 點位表 (points)
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS points (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    device_id       UUID NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    address         VARCHAR(255) NOT NULL,
    function        VARCHAR(50) NOT NULL DEFAULT '',
    data_type       VARCHAR(20) NOT NULL,
    mode            VARCHAR(20) NOT NULL DEFAULT 'read',
    polling_group_id UUID REFERENCES polling_groups(id) ON DELETE SET NULL,
    last_read_at    TIMESTAMPTZ,
    last_value      JSONB,
    last_error      TEXT,
    enabled         BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT points_data_type_check CHECK (
        data_type IN ('bool', 'int16', 'uint16', 'int32', 'uint32', 
                      'int64', 'uint64', 'float32', 'float64', 'string')
    ),
    CONSTRAINT points_mode_check CHECK (
        mode IN ('read', 'readwrite')
    ),
    CONSTRAINT points_device_address_function_unique UNIQUE (device_id, address, function)
);

CREATE INDEX idx_points_device_id ON points(device_id);
CREATE INDEX idx_points_polling_group_id ON points(polling_group_id);
CREATE INDEX idx_points_enabled ON points(enabled);
CREATE INDEX idx_points_data_type ON points(data_type);

COMMENT ON TABLE points IS '點位目錄，設備上的可讀取/寫入位址';
COMMENT ON COLUMN points.address IS '協議特定的位址字串';
COMMENT ON COLUMN points.function IS '協議功能碼或存取方式';

-- -----------------------------------------------------------------------------
-- 標籤表 (tags)
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS tags (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    key             VARCHAR(128) NOT NULL,
    key_lower       VARCHAR(128) NOT NULL,
    display_name    VARCHAR(255) NOT NULL,
    description     TEXT,
    unit            VARCHAR(50),
    data_type       VARCHAR(20) NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'draft',
    labels          JSONB DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT tags_key_lower_unique UNIQUE (key_lower),
    CONSTRAINT tags_key_format_check CHECK (
        key ~ '^[a-zA-Z0-9_\-./]{1,128}$'
    ),
    CONSTRAINT tags_data_type_check CHECK (
        data_type IN ('bool', 'int16', 'uint16', 'int32', 'uint32', 
                      'int64', 'uint64', 'float32', 'float64', 'string')
    ),
    CONSTRAINT tags_status_check CHECK (
        status IN ('draft', 'active', 'retired')
    )
);

CREATE INDEX idx_tags_key ON tags(key);
CREATE INDEX idx_tags_key_lower ON tags(key_lower);
CREATE INDEX idx_tags_status ON tags(status);
CREATE INDEX idx_tags_data_type ON tags(data_type);
CREATE INDEX idx_tags_labels ON tags USING GIN (labels);

COMMENT ON TABLE tags IS '全域標籤字典';
COMMENT ON COLUMN tags.key IS '標籤鍵 (唯一，大小寫敏感顯示)';
COMMENT ON COLUMN tags.key_lower IS '標籤鍵小寫 (用於唯一性檢查)';
COMMENT ON COLUMN tags.labels IS '標籤屬性 (JSON)，用於分類和篩選';

-- -----------------------------------------------------------------------------
-- 映射表 (mappings)
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS mappings (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    point_id            UUID NOT NULL REFERENCES points(id) ON DELETE CASCADE,
    tag_id              UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    transform_pipeline  JSONB NOT NULL DEFAULT '[]',
    status              TEXT NOT NULL DEFAULT 'active',
    rule_candidate_id   TEXT,
    proposed_signature  TEXT,
    last_applied_signature TEXT,
    blocking_reason     TEXT,
    enabled             BOOLEAN NOT NULL DEFAULT true,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT mappings_point_tag_unique UNIQUE (point_id, tag_id)
);

CREATE INDEX idx_mappings_point_id ON mappings(point_id);
CREATE INDEX idx_mappings_tag_id ON mappings(tag_id);
CREATE INDEX idx_mappings_enabled ON mappings(enabled);
CREATE INDEX idx_mappings_rule_candidate_id ON mappings(rule_candidate_id);

COMMENT ON TABLE mappings IS '點位到標籤的映射關係';
COMMENT ON COLUMN mappings.transform_pipeline IS '轉換管線步驟 (JSON 陣列)';

-- -----------------------------------------------------------------------------
-- 時序資料表 (timeseries) - 分區父表
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS timeseries (
    id              BIGSERIAL,
    tag_id          UUID NOT NULL,
    ts              TIMESTAMPTZ NOT NULL,
    value_num       DOUBLE PRECISION,
    value_text      TEXT,
    value_bool      BOOLEAN,
    raw_value       JSONB,
    quality         VARCHAR(20) NOT NULL DEFAULT 'good',
    
    CONSTRAINT timeseries_quality_check CHECK (
        quality IN ('good', 'bad', 'uncertain')
    ),
    CONSTRAINT timeseries_value_check CHECK (
        (value_num IS NOT NULL)::int + 
        (value_text IS NOT NULL)::int + 
        (value_bool IS NOT NULL)::int <= 1
    ),
    PRIMARY KEY (id, ts)
) PARTITION BY RANGE (ts);

-- 時序表索引 (在分區上自動繼承)
CREATE INDEX idx_timeseries_tag_ts ON timeseries(tag_id, ts DESC);
CREATE INDEX idx_timeseries_ts ON timeseries(ts DESC);
CREATE INDEX idx_timeseries_quality ON timeseries(quality) WHERE quality != 'good';

COMMENT ON TABLE timeseries IS '時序資料儲存表 (分區)';
COMMENT ON COLUMN timeseries.ts IS '時間戳記';
COMMENT ON COLUMN timeseries.value_num IS '數值型態值';
COMMENT ON COLUMN timeseries.value_text IS '文字型態值';
COMMENT ON COLUMN timeseries.value_bool IS '布林型態值';
COMMENT ON COLUMN timeseries.raw_value IS '原始值 (JSON)，用於追溯';

-- -----------------------------------------------------------------------------
-- 系統設定表 (system_settings)
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS system_settings (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    key             VARCHAR(100) NOT NULL UNIQUE,
    value           JSONB NOT NULL,
    description     TEXT,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 插入預設系統設定
INSERT INTO system_settings (key, value, description) VALUES
    ('write_precision', '"second"', '時序寫入時間精度 (second 或 millisecond)'),
    ('partition_interval', '"monthly"', '時序表分區間隔 (daily, weekly, monthly)'),
    ('batch_size', '1000', '批次寫入大小'),
    ('flush_interval', '5000', '批次刷新間隔 (毫秒)'),
    ('default_retry_count', '3', '預設重試次數'),
    ('default_retry_delay', '1000', '預設重試延遲 (毫秒)'),
    ('modbus_share', '{"enabled": false, "bind_address": "127.0.0.1", "port": 5020, "slave_id": 1, "capacity_registers": 32768, "settings_revision": "migration-018"}', 'Local Modbus Share listener settings')
ON CONFLICT (key) DO NOTHING;

COMMENT ON TABLE system_settings IS '系統設定表';

-- -----------------------------------------------------------------------------
-- 分區管理函數
-- -----------------------------------------------------------------------------

-- 建立月度分區的函數
CREATE OR REPLACE FUNCTION create_timeseries_partition_monthly(
    partition_date DATE
) RETURNS TEXT AS $$
DECLARE
    partition_name TEXT;
    start_date DATE;
    end_date DATE;
BEGIN
    start_date := date_trunc('month', partition_date);
    end_date := start_date + INTERVAL '1 month';
    partition_name := 'timeseries_' || to_char(start_date, 'YYYY_MM');
    
    -- 檢查分區是否已存在
    IF NOT EXISTS (
        SELECT 1 FROM pg_class c
        JOIN pg_namespace n ON n.oid = c.relnamespace
        WHERE c.relname = partition_name
        AND n.nspname = current_schema()
    ) THEN
        EXECUTE format(
            'CREATE TABLE %I PARTITION OF timeseries
             FOR VALUES FROM (%L) TO (%L)',
            partition_name, start_date, end_date
        );
        RETURN partition_name;
    END IF;
    
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- 建立每日分區的函數
CREATE OR REPLACE FUNCTION create_timeseries_partition_daily(
    partition_date DATE
) RETURNS TEXT AS $$
DECLARE
    partition_name TEXT;
    start_date DATE;
    end_date DATE;
BEGIN
    start_date := partition_date;
    end_date := start_date + INTERVAL '1 day';
    partition_name := 'timeseries_' || to_char(start_date, 'YYYY_MM_DD');
    
    IF NOT EXISTS (
        SELECT 1 FROM pg_class c
        JOIN pg_namespace n ON n.oid = c.relnamespace
        WHERE c.relname = partition_name
        AND n.nspname = current_schema()
    ) THEN
        EXECUTE format(
            'CREATE TABLE %I PARTITION OF timeseries
             FOR VALUES FROM (%L) TO (%L)',
            partition_name, start_date, end_date
        );
        RETURN partition_name;
    END IF;
    
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- 建立每週分區的函數
CREATE OR REPLACE FUNCTION create_timeseries_partition_weekly(
    partition_date DATE
) RETURNS TEXT AS $$
DECLARE
    partition_name TEXT;
    start_date DATE;
    end_date DATE;
BEGIN
    start_date := date_trunc('week', partition_date);
    end_date := start_date + INTERVAL '1 week';
    partition_name := 'timeseries_' || to_char(start_date, 'YYYY_"W"IW');
    
    IF NOT EXISTS (
        SELECT 1 FROM pg_class c
        JOIN pg_namespace n ON n.oid = c.relnamespace
        WHERE c.relname = partition_name
        AND n.nspname = current_schema()
    ) THEN
        EXECUTE format(
            'CREATE TABLE %I PARTITION OF timeseries
             FOR VALUES FROM (%L) TO (%L)',
            partition_name, start_date, end_date
        );
        RETURN partition_name;
    END IF;
    
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- 確保當前月份分區存在
SELECT create_timeseries_partition_monthly(CURRENT_DATE);
-- 預建下個月的分區
SELECT create_timeseries_partition_monthly(CURRENT_DATE + INTERVAL '1 month');

COMMENT ON FUNCTION create_timeseries_partition_monthly IS '建立指定月份的時序表分區';
COMMENT ON FUNCTION create_timeseries_partition_daily IS '建立指定日期的時序表分區';
COMMENT ON FUNCTION create_timeseries_partition_weekly IS '建立指定週的時序表分區';

-- -----------------------------------------------------------------------------
-- 自動更新 updated_at 觸發器
-- -----------------------------------------------------------------------------
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_devices_updated_at
    BEFORE UPDATE ON devices
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_polling_groups_updated_at
    BEFORE UPDATE ON polling_groups
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_points_updated_at
    BEFORE UPDATE ON points
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_tags_updated_at
    BEFORE UPDATE ON tags
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_mappings_updated_at
    BEFORE UPDATE ON mappings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_system_settings_updated_at
    BEFORE UPDATE ON system_settings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- -----------------------------------------------------------------------------
-- 標籤鍵小寫自動填充觸發器
-- -----------------------------------------------------------------------------
CREATE OR REPLACE FUNCTION update_tag_key_lower()
RETURNS TRIGGER AS $$
BEGIN
    NEW.key_lower = LOWER(NEW.key);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_tags_key_lower
    BEFORE INSERT OR UPDATE ON tags
    FOR EACH ROW EXECUTE FUNCTION update_tag_key_lower();
