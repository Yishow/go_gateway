-- Persistent identity for the one canonical basic managed group per device.
-- Keeping the key after a group is tombstoned prevents accidental reuse.
CREATE TABLE IF NOT EXISTS write_group_basic_keys (
    workspace_id    TEXT NOT NULL,
    device_id       TEXT NOT NULL,
    canonical_role  TEXT NOT NULL,
    group_id        TEXT NOT NULL,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (workspace_id, device_id, canonical_role),
    UNIQUE (group_id),
    FOREIGN KEY (group_id) REFERENCES write_groups(id) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_write_group_basic_keys_group
    ON write_group_basic_keys(group_id);
