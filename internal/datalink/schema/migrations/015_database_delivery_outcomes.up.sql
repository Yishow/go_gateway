ALTER TABLE database_connectors
    ADD COLUMN IF NOT EXISTS last_schema_ensure_at TIMESTAMPTZ;

ALTER TABLE database_connectors
    ADD COLUMN IF NOT EXISTS last_schema_ensure_status TEXT NOT NULL DEFAULT '';

ALTER TABLE database_connectors
    ADD COLUMN IF NOT EXISTS last_schema_ensure_error TEXT NOT NULL DEFAULT '';

ALTER TABLE database_connectors
    ADD COLUMN IF NOT EXISTS last_write_at TIMESTAMPTZ;

ALTER TABLE database_connectors
    ADD COLUMN IF NOT EXISTS last_write_status TEXT NOT NULL DEFAULT '';

ALTER TABLE database_connectors
    ADD COLUMN IF NOT EXISTS last_write_error TEXT NOT NULL DEFAULT '';

ALTER TABLE database_connectors
    ADD COLUMN IF NOT EXISTS last_flush_at TIMESTAMPTZ;

ALTER TABLE database_connectors
    ADD COLUMN IF NOT EXISTS last_flush_status TEXT NOT NULL DEFAULT '';

ALTER TABLE database_connectors
    ADD COLUMN IF NOT EXISTS last_flush_error TEXT NOT NULL DEFAULT '';
