-- Local lineage evidence of a shared managed_long_v1 table, kept on the scope
-- row 032 created. The confirmation of a generation freezes the layout the
-- table was verified against (including the extension schema of a hypertable),
-- the connector revision it was confirmed on and the operation that verified
-- it. A table rebuilt after it vanished starts the next generation and marks
-- that history stored before may be lost. repair_state is non-empty when an
-- automatic restore ended in a result nobody can vouch for; delivery to the
-- scope then waits for an explicit new confirmation. One statement per line:
-- the migrator skips an ADD COLUMN whose column already exists.
ALTER TABLE write_group_storage_scopes ADD COLUMN confirmed_layout TEXT NOT NULL DEFAULT '';
ALTER TABLE write_group_storage_scopes ADD COLUMN confirmed_connector_revision TEXT NOT NULL DEFAULT '';
ALTER TABLE write_group_storage_scopes ADD COLUMN confirmed_operation_id TEXT NOT NULL DEFAULT '';
ALTER TABLE write_group_storage_scopes ADD COLUMN confirmed_digest TEXT NOT NULL DEFAULT '';
ALTER TABLE write_group_storage_scopes ADD COLUMN history_may_be_lost INTEGER NOT NULL DEFAULT 0;
ALTER TABLE write_group_storage_scopes ADD COLUMN reconstruction_count INTEGER NOT NULL DEFAULT 0;
ALTER TABLE write_group_storage_scopes ADD COLUMN last_reconstructed_at TEXT NOT NULL DEFAULT '';
ALTER TABLE write_group_storage_scopes ADD COLUMN last_reconstruction_operation_id TEXT NOT NULL DEFAULT '';
ALTER TABLE write_group_storage_scopes ADD COLUMN repair_state TEXT NOT NULL DEFAULT '';
ALTER TABLE write_group_storage_scopes ADD COLUMN repair_reason TEXT NOT NULL DEFAULT '';
