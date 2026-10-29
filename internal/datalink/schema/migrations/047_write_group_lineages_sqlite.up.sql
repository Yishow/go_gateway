-- Canonical keys remain immutable genesis anchors. Only this registry owns heads.
CREATE TABLE write_group_lineages (
    lineage_id TEXT PRIMARY KEY NOT NULL CHECK (TRIM(lineage_id) <> ''),
    workspace_id TEXT NOT NULL CHECK (TRIM(workspace_id) <> ''),
    current_group_id TEXT NOT NULL UNIQUE REFERENCES write_groups(id) ON DELETE RESTRICT,
    reserved_group_id TEXT UNIQUE REFERENCES write_groups(id) ON DELETE RESTRICT,
    reservation_operation_id TEXT,
    reservation_intent_digest TEXT,
    revision INTEGER NOT NULL DEFAULT 1 CHECK (typeof(revision) = 'integer' AND revision >= 1),
    CHECK (
        (reserved_group_id IS NULL AND reservation_operation_id IS NULL AND reservation_intent_digest IS NULL)
        OR (reserved_group_id IS NOT NULL AND TRIM(reserved_group_id) <> ''
            AND reservation_operation_id IS NOT NULL AND TRIM(reservation_operation_id) <> ''
            AND reservation_intent_digest IS NOT NULL AND TRIM(reservation_intent_digest) <> '')
    ),
    CHECK (reserved_group_id IS NULL OR reserved_group_id <> current_group_id)
);

CREATE TABLE write_group_generations (
    group_id TEXT PRIMARY KEY NOT NULL REFERENCES write_groups(id) ON DELETE RESTRICT,
    lineage_id TEXT NOT NULL REFERENCES write_group_lineages(lineage_id) ON DELETE RESTRICT,
    generation INTEGER NOT NULL CHECK (typeof(generation) = 'integer' AND generation >= 1),
    predecessor_group_id TEXT REFERENCES write_group_generations(group_id) ON DELETE RESTRICT,
    UNIQUE (lineage_id, generation),
    CHECK (
        (generation = 1 AND predecessor_group_id IS NULL)
        OR (generation > 1 AND predecessor_group_id IS NOT NULL AND predecessor_group_id <> group_id)
    )
);

ALTER TABLE write_group_basic_keys ADD COLUMN lineage_id TEXT
    REFERENCES write_group_lineages(lineage_id) ON DELETE RESTRICT
    CHECK (lineage_id IS NULL OR TRIM(lineage_id) <> '');
ALTER TABLE write_group_profile_keys ADD COLUMN lineage_id TEXT
    REFERENCES write_group_lineages(lineage_id) ON DELETE RESTRICT
    CHECK (lineage_id IS NULL OR TRIM(lineage_id) <> '');
