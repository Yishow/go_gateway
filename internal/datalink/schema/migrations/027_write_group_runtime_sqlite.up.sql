-- Frozen, verified applied descriptors for offline intake and journal closure.
-- This is recovery metadata; accepted samples and delivery remain in 025.
CREATE TABLE IF NOT EXISTS wg_runtime_versions (
    workspace_id TEXT NOT NULL,
    group_id TEXT NOT NULL,
    group_revision TEXT NOT NULL,
    payload TEXT NOT NULL,
    PRIMARY KEY (workspace_id, group_id, group_revision)
);
