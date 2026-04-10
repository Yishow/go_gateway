ALTER TABLE database_target_mappings
    DROP COLUMN IF EXISTS write_interval_seconds;

ALTER TABLE database_target_mappings
    DROP COLUMN IF EXISTS group_key;

ALTER TABLE database_connectors
    DROP COLUMN IF EXISTS default_write_interval_seconds;
