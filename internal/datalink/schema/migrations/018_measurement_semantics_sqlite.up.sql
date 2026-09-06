-- Migration: 018_measurement_semantics_sqlite.up.sql (SQLite)

CREATE TABLE IF NOT EXISTS measurement_definitions (
    id TEXT PRIMARY KEY,
    workspace_id TEXT NOT NULL,
    device_id TEXT NOT NULL,
    point_id TEXT NOT NULL,
    tag_id TEXT,
    equipment_id TEXT NOT NULL,
    definition_revision TEXT NOT NULL,
    source_binding_revision TEXT NOT NULL,
    series_epoch TEXT NOT NULL DEFAULT 'epoch-1',
    name TEXT NOT NULL,
    quantity TEXT NOT NULL,
    unit TEXT,
    semantic_kind TEXT NOT NULL,
    numeric_encoding TEXT,
    counter_policy TEXT,
    state_map TEXT,
    bitmask_labels TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_measurement_defs_ws_dev ON measurement_definitions(workspace_id, device_id);
CREATE INDEX IF NOT EXISTS idx_measurement_defs_point ON measurement_definitions(point_id);
