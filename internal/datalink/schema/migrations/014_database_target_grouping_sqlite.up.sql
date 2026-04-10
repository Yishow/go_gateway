ALTER TABLE database_connectors
    ADD COLUMN default_write_interval_seconds INTEGER NOT NULL DEFAULT 15;

ALTER TABLE database_target_mappings
    ADD COLUMN group_key TEXT;

ALTER TABLE database_target_mappings
    ADD COLUMN write_interval_seconds INTEGER;
