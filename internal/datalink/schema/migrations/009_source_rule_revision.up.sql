-- -----------------------------------------------------------------------------
-- 來源規則表新增最新 revision id 欄位
-- -----------------------------------------------------------------------------

ALTER TABLE source_rules ADD COLUMN revision_id TEXT;

UPDATE source_rules
SET revision_id = id || ':legacy'
WHERE revision_id IS NULL OR btrim(revision_id) = '';

ALTER TABLE source_rules ALTER COLUMN revision_id SET NOT NULL;
