-- 建表操作帳本：每次預覽確認對應一筆 operation。operation_id 與 token 各自唯一，
-- 同一目標範圍同時只允許一筆 pending/running，多程序或重啟後都不會重複取得執行權。
CREATE TABLE IF NOT EXISTS managed_schema_operations (
    operation_id TEXT PRIMARY KEY,
    token TEXT NOT NULL,
    workspace_id TEXT NOT NULL,
    scope_key TEXT NOT NULL,
    owner TEXT NOT NULL,
    action TEXT NOT NULL,
    status TEXT NOT NULL,
    executed_statements INTEGER NOT NULL DEFAULT 0,
    verified_digest TEXT NOT NULL DEFAULT '',
    reason TEXT NOT NULL DEFAULT '',
    next_action TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    completed_at DATETIME
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_managed_schema_operations_token
    ON managed_schema_operations(token);

CREATE UNIQUE INDEX IF NOT EXISTS idx_managed_schema_operations_active_scope
    ON managed_schema_operations(scope_key) WHERE status IN ('pending', 'running');

CREATE INDEX IF NOT EXISTS idx_managed_schema_operations_workspace
    ON managed_schema_operations(workspace_id, created_at);
