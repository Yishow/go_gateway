-- Forward cutovers and their rollbacks. A cutover moves selected sources of a
-- wide or custom group to long profile groups at one exact, immutable instant
-- (boundary): observations before it belong to the revision that ends, those at
-- or after it to the revision that starts. The row is the only authority for
-- that instant; nothing derives it from a clock after the fact.
CREATE TABLE IF NOT EXISTS write_group_cutovers (
    id               TEXT PRIMARY KEY,
    workspace_id     TEXT NOT NULL,
    kind             TEXT NOT NULL CHECK (kind IN ('forward', 'rollback')),
    source_group_id  TEXT NOT NULL,
    boundary         TEXT NOT NULL,
    prior_cutover_id TEXT NOT NULL DEFAULT '',
    review_digest    TEXT NOT NULL,
    detail           TEXT NOT NULL,
    created_at       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- A forward cutover can be rolled back once.
CREATE UNIQUE INDEX IF NOT EXISTS idx_write_group_cutovers_prior
    ON write_group_cutovers(prior_cutover_id) WHERE prior_cutover_id <> '';

CREATE INDEX IF NOT EXISTS idx_write_group_cutovers_source
    ON write_group_cutovers(workspace_id, source_group_id, boundary);

-- The instant an applied revision stops owning its sources. A revision has at
-- most one end, and the end is never moved.
CREATE TABLE IF NOT EXISTS write_group_revision_ends (
    workspace_id   TEXT NOT NULL,
    group_id       TEXT NOT NULL,
    group_revision TEXT NOT NULL,
    ends_at        TEXT NOT NULL,
    cutover_id     TEXT NOT NULL,
    PRIMARY KEY (group_id, group_revision)
);
