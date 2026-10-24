-- Server-owned identity of one shared managed_long_v1 table. The owner is the
-- workspace plus the confirmed destination (connector, database, schema) plus
-- storage mode and layout; no group, label or point count takes part in it.
-- table_name and owner_column are derived from the scope and frozen here with
-- the layout digest and storage generation, so later recovery can prove which
-- layout generation a table was confirmed for.
CREATE TABLE IF NOT EXISTS write_group_storage_scopes (
    workspace_id    TEXT NOT NULL,
    connector_id    TEXT NOT NULL,
    database_name   TEXT NOT NULL DEFAULT '',
    table_schema    TEXT NOT NULL DEFAULT '',
    storage_mode    TEXT NOT NULL,
    layout_version  TEXT NOT NULL,
    table_name      TEXT NOT NULL,
    owner_column    TEXT NOT NULL,
    layout_digest   TEXT NOT NULL,
    generation      INTEGER NOT NULL DEFAULT 1,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (workspace_id, connector_id, database_name, table_schema, storage_mode, layout_version)
);
