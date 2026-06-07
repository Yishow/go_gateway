CREATE TABLE IF NOT EXISTS workspace_audit_history (
    id TEXT PRIMARY KEY,
    workspace_id TEXT NOT NULL,
    event_type TEXT NOT NULL,
    result TEXT NOT NULL,
    scope TEXT NOT NULL DEFAULT 'workspace',
    reference_id TEXT NOT NULL DEFAULT '',
    details TEXT NOT NULL DEFAULT '{}',
    occurred_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_workspace_audit_history_workspace_time
    ON workspace_audit_history(workspace_id, occurred_at DESC);
