-- Additive, nullable recording reduction metadata for write groups. NULL
-- means a group saved before reductions existed: it keeps reading as latest
-- and is never backfilled. The member role and type columns of the same
-- recording contract are created beside the reduction so one migration
-- carries the semantics; they stay unwired until their own slice.
-- One statement per line; the migrator skips columns that already exist.
ALTER TABLE write_groups ADD COLUMN sample_reduction TEXT;
ALTER TABLE write_group_members ADD COLUMN value_role TEXT;
ALTER TABLE write_group_members ADD COLUMN source_value_type TEXT;
ALTER TABLE write_group_members ADD COLUMN result_value_type TEXT;
