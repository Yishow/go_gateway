-- Persistent create-once identity for one canonical group per circuit profile.
-- An assigned circuit is keyed by circuit_id; an unassigned raw selection is
-- keyed by its device through selection_scope. The two scopes live in separate
-- columns so an opaque circuit_id can never collide with a device id, and the
-- CHECK admits exactly one scope per role. Keeping the key after a group is
-- tombstoned prevents accidental reuse.
CREATE TABLE IF NOT EXISTS write_group_profile_keys (
    workspace_id    TEXT NOT NULL,
    circuit_id      TEXT NOT NULL DEFAULT '',
    selection_scope TEXT NOT NULL DEFAULT '',
    profile_kind    TEXT NOT NULL CHECK (profile_kind IN ('instantaneous', 'cumulative')),
    canonical_role  TEXT NOT NULL CHECK (canonical_role IN ('circuit', 'unassigned')),
    group_id        TEXT NOT NULL,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (workspace_id, circuit_id, selection_scope, profile_kind, canonical_role),
    UNIQUE (group_id),
    CHECK (
        (canonical_role = 'circuit' AND circuit_id <> '' AND selection_scope = '')
        OR (canonical_role = 'unassigned' AND circuit_id = '' AND selection_scope <> '' AND profile_kind = 'instantaneous')
    ),
    FOREIGN KEY (group_id) REFERENCES write_groups(id) ON DELETE RESTRICT
);
