ALTER TABLE source_rules
    ADD COLUMN IF NOT EXISTS share_enabled BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE source_rules
    ADD COLUMN IF NOT EXISTS share_start_register INTEGER;

ALTER TABLE source_rules
    ADD COLUMN IF NOT EXISTS share_stride INTEGER;
