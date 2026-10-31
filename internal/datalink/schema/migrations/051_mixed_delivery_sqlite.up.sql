CREATE TABLE wg_delivery_mixed_progress (
    workspace_id TEXT NOT NULL,
    group_id TEXT NOT NULL,
    group_revision TEXT NOT NULL,
    version INTEGER NOT NULL CHECK (version = 1),
    descriptor_digest TEXT NOT NULL,
    closed_through TEXT NOT NULL,
    checkpoint TEXT NOT NULL CHECK (json_valid(checkpoint) AND json_type(checkpoint) = 'object'),
    checkpoint_digest TEXT NOT NULL,
    last_closure_digest TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    PRIMARY KEY (workspace_id, group_id, group_revision),
    FOREIGN KEY (workspace_id, group_id, group_revision)
        REFERENCES wg_runtime_versions (workspace_id, group_id, group_revision)
);
CREATE TABLE wg_delivery_mixed_outcomes (
    workspace_id TEXT NOT NULL,
    group_id TEXT NOT NULL,
    group_revision TEXT NOT NULL,
    record_id TEXT NOT NULL,
    due_at TEXT NOT NULL,
    kind TEXT NOT NULL,
    outcome TEXT NOT NULL CHECK (json_valid(outcome) AND json_type(outcome) = 'object'),
    outcome_digest TEXT NOT NULL,
    created_at TEXT NOT NULL,
    PRIMARY KEY (workspace_id, group_id, group_revision, record_id),
    FOREIGN KEY (workspace_id, group_id, group_revision)
        REFERENCES wg_runtime_versions (workspace_id, group_id, group_revision)
);
ALTER TABLE wg_delivery_outbox ADD COLUMN due_at TEXT;
CREATE INDEX idx_wg_delivery_outbox_due_head ON wg_delivery_outbox (partition_key, COALESCE(due_at,bucket_start), effect_key);
