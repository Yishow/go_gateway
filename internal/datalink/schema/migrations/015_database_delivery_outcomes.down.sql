ALTER TABLE database_connectors
    DROP COLUMN IF EXISTS last_flush_error;

ALTER TABLE database_connectors
    DROP COLUMN IF EXISTS last_flush_status;

ALTER TABLE database_connectors
    DROP COLUMN IF EXISTS last_flush_at;

ALTER TABLE database_connectors
    DROP COLUMN IF EXISTS last_write_error;

ALTER TABLE database_connectors
    DROP COLUMN IF EXISTS last_write_status;

ALTER TABLE database_connectors
    DROP COLUMN IF EXISTS last_write_at;

ALTER TABLE database_connectors
    DROP COLUMN IF EXISTS last_schema_ensure_error;

ALTER TABLE database_connectors
    DROP COLUMN IF EXISTS last_schema_ensure_status;

ALTER TABLE database_connectors
    DROP COLUMN IF EXISTS last_schema_ensure_at;
