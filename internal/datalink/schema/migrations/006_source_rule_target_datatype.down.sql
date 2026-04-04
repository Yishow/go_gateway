-- -----------------------------------------------------------------------------
-- 回滾：移除來源規則表的目標資料型態與縮放欄位
-- -----------------------------------------------------------------------------

-- PostgreSQL 版本
ALTER TABLE source_rules DROP COLUMN IF EXISTS target_data_type;
ALTER TABLE source_rules DROP COLUMN IF EXISTS scale_multiplier;
ALTER TABLE source_rules DROP COLUMN IF EXISTS scale_offset;