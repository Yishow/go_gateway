ALTER TABLE write_group_members ADD COLUMN cadence TEXT CHECK (cadence IS NULL OR (typeof(cadence) = 'text' AND json_valid(cadence) AND json_type(cadence) = 'object'));
