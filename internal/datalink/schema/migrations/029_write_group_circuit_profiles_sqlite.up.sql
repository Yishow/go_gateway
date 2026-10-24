-- Additive, nullable circuit recording profile metadata for write groups.
-- NULL means a group saved before profiles existed: it keeps reading as
-- managed_wide_v1 with no profile, and existing intervals are untouched.
-- The migrator executes this file one statement per line and skips any
-- ADD COLUMN whose column already exists, so keep one statement per line.
ALTER TABLE write_groups ADD COLUMN layout_version TEXT;
ALTER TABLE write_groups ADD COLUMN profile_kind TEXT CHECK (profile_kind IS NULL OR profile_kind IN ('instantaneous', 'cumulative'));
ALTER TABLE write_groups ADD COLUMN profile_interval INTEGER;
ALTER TABLE write_groups ADD COLUMN profile_unit TEXT CHECK (profile_unit IS NULL OR profile_unit IN ('second', 'minute'));
ALTER TABLE write_group_members ADD COLUMN circuit_id TEXT;
ALTER TABLE write_group_members ADD COLUMN circuit_name TEXT;
