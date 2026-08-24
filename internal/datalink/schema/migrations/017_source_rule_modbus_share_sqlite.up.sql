ALTER TABLE source_rules ADD COLUMN share_enabled INTEGER NOT NULL DEFAULT 0;
ALTER TABLE source_rules ADD COLUMN share_start_register INTEGER;
ALTER TABLE source_rules ADD COLUMN share_stride INTEGER;
