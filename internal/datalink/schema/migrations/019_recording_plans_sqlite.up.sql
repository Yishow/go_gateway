-- Migration: 019_recording_plans_sqlite.up.sql (SQLite)

CREATE TABLE IF NOT EXISTS recording_plans (
    id TEXT PRIMARY KEY,
    workspace_id TEXT NOT NULL,
    revision TEXT NOT NULL DEFAULT 'rev-1',
    applied_revision TEXT NOT NULL DEFAULT '',
    name TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'draft',
    timezone TEXT NOT NULL DEFAULT 'Asia/Taipei',
    members TEXT NOT NULL,
    streams TEXT NOT NULL,
    destinations TEXT NOT NULL,
    retention TEXT NOT NULL,
    limits TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_recording_plans_ws ON recording_plans(workspace_id);

CREATE TABLE IF NOT EXISTS managed_schema_preview_tokens (
    token TEXT PRIMARY KEY,
    workspace_id TEXT NOT NULL,
    plan_id TEXT NOT NULL,
    plan_revision TEXT NOT NULL,
    connector_id TEXT NOT NULL,
    table_prefix TEXT NOT NULL,
    statements TEXT NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
