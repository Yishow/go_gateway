-- 為建表預覽 token 補上完整範圍與保護欄位。既有 token 一律補空字串：
-- 套用前的重驗會把缺少保護欄位的 token 視為舊版並要求重新預覽，不補猜有效性。
ALTER TABLE managed_schema_preview_tokens ADD COLUMN operation_id TEXT NOT NULL DEFAULT '';
ALTER TABLE managed_schema_preview_tokens ADD COLUMN action TEXT NOT NULL DEFAULT '';
ALTER TABLE managed_schema_preview_tokens ADD COLUMN workspace_revision TEXT NOT NULL DEFAULT '';
ALTER TABLE managed_schema_preview_tokens ADD COLUMN connector_revision TEXT NOT NULL DEFAULT '';
ALTER TABLE managed_schema_preview_tokens ADD COLUMN dialect TEXT NOT NULL DEFAULT '';
ALTER TABLE managed_schema_preview_tokens ADD COLUMN database_name TEXT NOT NULL DEFAULT '';
ALTER TABLE managed_schema_preview_tokens ADD COLUMN schema_name TEXT NOT NULL DEFAULT '';
ALTER TABLE managed_schema_preview_tokens ADD COLUMN no_change_reason TEXT NOT NULL DEFAULT '';
ALTER TABLE managed_schema_preview_tokens ADD COLUMN tables TEXT NOT NULL DEFAULT '[]';
ALTER TABLE managed_schema_preview_tokens ADD COLUMN digest TEXT NOT NULL DEFAULT '';
