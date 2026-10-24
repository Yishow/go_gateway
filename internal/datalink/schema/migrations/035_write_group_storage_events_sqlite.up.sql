-- Durable, value-free diagnostics of managed long storage: confirmations,
-- rebuilt tables and recoveries that were refused. A repeated event of the same
-- kind and reason for the same group and generation is counted, not stacked.
CREATE TABLE IF NOT EXISTS write_group_storage_events (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    workspace_id    TEXT NOT NULL,
    connector_id    TEXT NOT NULL,
    database_name   TEXT NOT NULL DEFAULT '',
    table_schema    TEXT NOT NULL DEFAULT '',
    storage_mode    TEXT NOT NULL,
    layout_version  TEXT NOT NULL,
    generation      INTEGER NOT NULL,
    group_id        TEXT NOT NULL DEFAULT '',
    kind            TEXT NOT NULL,
    code            TEXT NOT NULL DEFAULT '',
    reason          TEXT NOT NULL DEFAULT '',
    operation_id    TEXT NOT NULL DEFAULT '',
    occurrences     INTEGER NOT NULL DEFAULT 1,
    occurred_at     TEXT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_write_group_storage_events_identity
    ON write_group_storage_events(workspace_id, connector_id, database_name, table_schema, storage_mode, layout_version, generation, group_id, kind, reason);

CREATE INDEX IF NOT EXISTS idx_write_group_storage_events_group
    ON write_group_storage_events(workspace_id, group_id, id);
