-- Mutable authorized SQL Server native-partition projection is separate from
-- the immutable confirmed layout and from the MySQL-family state.
CREATE TABLE IF NOT EXISTS write_group_sqlserver_partition_maintenance (
    workspace_id TEXT NOT NULL,
    connector_id TEXT NOT NULL,
    database_name TEXT NOT NULL,
    table_schema TEXT NOT NULL,
    storage_mode TEXT NOT NULL,
    layout_version TEXT NOT NULL,
    base_layout_digest TEXT NOT NULL,
    connector_revision TEXT NOT NULL,
    generation INTEGER NOT NULL,
    current_evidence TEXT NOT NULL,
    revision INTEGER NOT NULL,
    last_operation_id TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL,
    reason TEXT NOT NULL DEFAULT '',
    updated_at TEXT NOT NULL,
    PRIMARY KEY (workspace_id, connector_id, database_name, table_schema, storage_mode, layout_version)
);
