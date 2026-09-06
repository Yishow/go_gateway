-- Migration: 018_measurement_semantics.up.sql (PostgreSQL)

CREATE TABLE IF NOT EXISTS measurement_definitions (
    id VARCHAR(64) PRIMARY KEY,
    workspace_id VARCHAR(64) NOT NULL,
    device_id VARCHAR(64) NOT NULL,
    point_id VARCHAR(64) NOT NULL,
    tag_id VARCHAR(64),
    equipment_id VARCHAR(64) NOT NULL,
    definition_revision VARCHAR(64) NOT NULL,
    source_binding_revision VARCHAR(64) NOT NULL,
    series_epoch VARCHAR(64) NOT NULL DEFAULT 'epoch-1',
    name VARCHAR(255) NOT NULL,
    quantity VARCHAR(64) NOT NULL,
    unit VARCHAR(32),
    semantic_kind VARCHAR(32) NOT NULL,
    numeric_encoding VARCHAR(32),
    counter_policy TEXT,
    state_map TEXT,
    bitmask_labels TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_measurement_defs_ws_dev ON measurement_definitions(workspace_id, device_id);
CREATE INDEX IF NOT EXISTS idx_measurement_defs_point ON measurement_definitions(point_id);

ALTER TABLE source_rules ADD COLUMN IF NOT EXISTS layout_mode VARCHAR(32) NOT NULL DEFAULT 'homogeneous';
ALTER TABLE source_rules ADD COLUMN IF NOT EXISTS mixed_items TEXT;
