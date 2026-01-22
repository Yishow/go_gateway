-- Revert: Remove collection stats and readiness columns from devices table

ALTER TABLE devices DROP COLUMN last_collected_at;
ALTER TABLE devices DROP COLUMN collection_count;
ALTER TABLE devices DROP COLUMN error_count;
ALTER TABLE devices DROP COLUMN readiness_status;
