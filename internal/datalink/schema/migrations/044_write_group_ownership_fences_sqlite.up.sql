-- Expand only: legacy cutovers have no fence and remain explicitly unknown.
-- The cutover ID is an opaque equality/CAS revision, never a sortable epoch.
-- Boundary authority stays in write_group_cutovers; these immutable rows bind
-- that intent to the exact ended/started revision membership and destination.
CREATE TABLE IF NOT EXISTS write_group_ownership_fences (
    fence_revision    TEXT NOT NULL,
    workspace_id      TEXT NOT NULL,
    direction         TEXT NOT NULL CHECK (direction IN ('ended', 'started')),
    group_id          TEXT NOT NULL,
    group_revision    TEXT NOT NULL,
    device_id         TEXT NOT NULL,
    point_id          TEXT NOT NULL,
    tag_id            TEXT NOT NULL,
    profile_kind      TEXT NOT NULL CHECK (profile_kind IN ('', 'instantaneous', 'cumulative')),
    connector_id      TEXT NOT NULL,
    connector_revision TEXT NOT NULL,
    database_name     TEXT NOT NULL,
    table_schema      TEXT NOT NULL,
    table_name        TEXT NOT NULL,
    PRIMARY KEY (fence_revision, direction, group_id, group_revision, device_id, point_id, tag_id),
    UNIQUE (fence_revision, direction, device_id, point_id, tag_id, profile_kind),
    FOREIGN KEY (fence_revision) REFERENCES write_group_cutovers(id) ON DELETE RESTRICT
);
CREATE INDEX IF NOT EXISTS idx_write_group_ownership_fences_revision
    ON write_group_ownership_fences(workspace_id, group_id, group_revision, direction);
CREATE INDEX IF NOT EXISTS idx_write_group_ownership_fences_source
    ON write_group_ownership_fences(workspace_id, device_id, point_id, tag_id, profile_kind);
