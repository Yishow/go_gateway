ALTER TABLE database_connectors
    ADD COLUMN IF NOT EXISTS default_write_interval_seconds INTEGER NOT NULL DEFAULT 15;

ALTER TABLE database_target_mappings
    ADD COLUMN IF NOT EXISTS group_key TEXT;

ALTER TABLE database_target_mappings
    ADD COLUMN IF NOT EXISTS write_interval_seconds INTEGER;
