-- Adds collection stats and readiness columns to devices table

-- SQLite/PostgreSQL Compatible

ALTER TABLE devices ADD COLUMN last_collected_at TIMESTAMP;
ALTER TABLE devices ADD COLUMN collection_count BIGINT DEFAULT 0;
ALTER TABLE devices ADD COLUMN error_count BIGINT DEFAULT 0;
ALTER TABLE devices ADD COLUMN readiness_status TEXT;
