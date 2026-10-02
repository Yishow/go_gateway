-- Canonical Studio V2 write-group authority and its persisted source members.
CREATE TABLE IF NOT EXISTS write_groups (
    id                              TEXT PRIMARY KEY,
    workspace_id                    TEXT NOT NULL,
    revision                        TEXT NOT NULL,
    applied_revision                TEXT NOT NULL DEFAULT '',
    name                            TEXT NOT NULL,
    status                          TEXT NOT NULL DEFAULT 'draft',
    destination_connector_id        TEXT NOT NULL,
    destination_connector_revision  TEXT NOT NULL,
    destination_database            TEXT NOT NULL DEFAULT '',
    destination_table_schema        TEXT NOT NULL DEFAULT 'main',
    destination_table_name          TEXT NOT NULL,
    destination_storage_strategy    TEXT NOT NULL DEFAULT 'custom',
    destination_schema_revision     TEXT NOT NULL DEFAULT '',
    destination_schema_digest       TEXT NOT NULL DEFAULT '',
    row_policy                      TEXT NOT NULL DEFAULT '{}',
    write_policy                    TEXT NOT NULL DEFAULT '{}',
    migration                       TEXT NOT NULL DEFAULT '{}',
    created_at                      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_write_groups_workspace
    ON write_groups(workspace_id, created_at DESC, id DESC);

CREATE TABLE IF NOT EXISTS write_group_members (
    group_id             TEXT NOT NULL,
    member_index         INTEGER NOT NULL,
    device_id            TEXT NOT NULL,
    point_id             TEXT NOT NULL,
    tag_id               TEXT NOT NULL,
    entity_key           TEXT NOT NULL DEFAULT '',
    source_revision      TEXT NOT NULL DEFAULT '',
    mapping_revision     TEXT NOT NULL DEFAULT '',
    measurement_id       TEXT,
    target_column        TEXT NOT NULL,
    required             INTEGER NOT NULL DEFAULT 1,
    max_age_seconds      INTEGER,
    PRIMARY KEY (group_id, member_index),
    UNIQUE (group_id, device_id, point_id, tag_id),
    FOREIGN KEY (group_id) REFERENCES write_groups(id) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_write_group_members_source
    ON write_group_members(device_id, point_id, tag_id);

-- Immutable applied snapshots keep old bucket ownership available while a
-- newer draft is scheduled for the next UTC boundary.
CREATE TABLE IF NOT EXISTS write_group_versions (
    group_id       TEXT NOT NULL,
    group_revision TEXT NOT NULL,
    effective_at   DATETIME NOT NULL,
    payload        TEXT NOT NULL,
    created_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (group_id, group_revision),
    FOREIGN KEY (group_id) REFERENCES write_groups(id) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_write_group_versions_effective
    ON write_group_versions(group_id, effective_at DESC, group_revision DESC);

-- Durable identity and provenance for reviewed legacy single-mapping imports.
CREATE TABLE IF NOT EXISTS write_group_migration_maps (
    workspace_id     TEXT NOT NULL,
    source_kind      TEXT NOT NULL,
    source_id        TEXT NOT NULL,
    source_revision  TEXT NOT NULL,
    group_id         TEXT NOT NULL,
    before_intent    TEXT NOT NULL,
    review_digest    TEXT NOT NULL,
    adapter_version  TEXT NOT NULL,
    created_at       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (workspace_id, source_kind, source_id),
    FOREIGN KEY (group_id) REFERENCES write_groups(id) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_write_group_migration_maps_group
    ON write_group_migration_maps(group_id);
