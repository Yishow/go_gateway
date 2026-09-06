-- Migration: 019_recording_plans.up.sql (PostgreSQL)

CREATE TABLE IF NOT EXISTS recording_plans (
    id VARCHAR(64) PRIMARY KEY,
    workspace_id VARCHAR(64) NOT NULL,
    revision VARCHAR(64) NOT NULL DEFAULT 'rev-1',
    applied_revision VARCHAR(64) NOT NULL DEFAULT '',
    name VARCHAR(255) NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'draft',
    timezone VARCHAR(64) NOT NULL DEFAULT 'Asia/Taipei',
    members TEXT NOT NULL,
    streams TEXT NOT NULL,
    destinations TEXT NOT NULL,
    retention TEXT NOT NULL,
    limits TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_recording_plans_ws ON recording_plans(workspace_id);

CREATE TABLE IF NOT EXISTS managed_schema_preview_tokens (
    token VARCHAR(64) PRIMARY KEY,
    workspace_id VARCHAR(64) NOT NULL,
    plan_id VARCHAR(64) NOT NULL,
    plan_revision VARCHAR(64) NOT NULL,
    connector_id VARCHAR(64) NOT NULL,
    table_prefix VARCHAR(64) NOT NULL,
    statements TEXT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
