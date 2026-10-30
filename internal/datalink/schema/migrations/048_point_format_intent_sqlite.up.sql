-- Add intent without inferring it from historical data_format values.
-- The Go migrator runs this atomically once and rejects partial schema state.
ALTER TABLE points ADD COLUMN format_mode TEXT
    CHECK (format_mode IS NULL OR format_mode IN ('', 'legacy', 'inherit', 'override'));
ALTER TABLE points ADD COLUMN format_override TEXT;
ALTER TABLE points ADD COLUMN format_version INTEGER NOT NULL DEFAULT 0;
ALTER TABLE points ADD COLUMN revision_id TEXT NOT NULL DEFAULT '';
