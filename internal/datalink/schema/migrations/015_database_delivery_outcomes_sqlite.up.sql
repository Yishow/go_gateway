ALTER TABLE database_connectors ADD COLUMN last_schema_ensure_at DATETIME;
ALTER TABLE database_connectors ADD COLUMN last_schema_ensure_status TEXT NOT NULL DEFAULT '';
ALTER TABLE database_connectors ADD COLUMN last_schema_ensure_error TEXT NOT NULL DEFAULT '';
ALTER TABLE database_connectors ADD COLUMN last_write_at DATETIME;
ALTER TABLE database_connectors ADD COLUMN last_write_status TEXT NOT NULL DEFAULT '';
ALTER TABLE database_connectors ADD COLUMN last_write_error TEXT NOT NULL DEFAULT '';
ALTER TABLE database_connectors ADD COLUMN last_flush_at DATETIME;
ALTER TABLE database_connectors ADD COLUMN last_flush_status TEXT NOT NULL DEFAULT '';
ALTER TABLE database_connectors ADD COLUMN last_flush_error TEXT NOT NULL DEFAULT '';
