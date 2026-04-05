-- -----------------------------------------------------------------------------
-- 來源規則表新增最新 revision id 欄位 (SQLite 版本)
-- -----------------------------------------------------------------------------

ALTER TABLE source_rules ADD COLUMN revision_id TEXT NOT NULL DEFAULT '';

UPDATE source_rules
SET revision_id = id || ':legacy'
WHERE trim(coalesce(revision_id, '')) = '';
