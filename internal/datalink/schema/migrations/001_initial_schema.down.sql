-- =============================================================================
-- Go Gateway - 資料收集管線資料庫遷移腳本 (PostgreSQL) - 回滾
-- 版本: 001_initial_schema
-- =============================================================================

-- 刪除觸發器
DROP TRIGGER IF EXISTS trigger_tags_key_lower ON tags;
DROP TRIGGER IF EXISTS trigger_system_settings_updated_at ON system_settings;
DROP TRIGGER IF EXISTS trigger_mappings_updated_at ON mappings;
DROP TRIGGER IF EXISTS trigger_tags_updated_at ON tags;
DROP TRIGGER IF EXISTS trigger_points_updated_at ON points;
DROP TRIGGER IF EXISTS trigger_polling_groups_updated_at ON polling_groups;
DROP TRIGGER IF EXISTS trigger_devices_updated_at ON devices;

-- 刪除函數
DROP FUNCTION IF EXISTS update_tag_key_lower();
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP FUNCTION IF EXISTS create_timeseries_partition_weekly(DATE);
DROP FUNCTION IF EXISTS create_timeseries_partition_daily(DATE);
DROP FUNCTION IF EXISTS create_timeseries_partition_monthly(DATE);

-- 刪除表 (按依賴順序)
DROP TABLE IF EXISTS system_settings;
DROP TABLE IF EXISTS timeseries;
DROP TABLE IF EXISTS mappings;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS points;
DROP TABLE IF EXISTS polling_groups;
DROP TABLE IF EXISTS devices;
