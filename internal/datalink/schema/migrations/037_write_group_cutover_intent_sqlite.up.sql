-- The whole request a cutover answered, kept beside its review digest. An apply
-- whose answer was lost and is sent again is recognised by this intent, not by
-- the digest alone. One statement per line.
ALTER TABLE write_group_cutovers ADD COLUMN intent TEXT NOT NULL DEFAULT '';
