-- Additive, nullable exact recording source proof. NULL keeps every existing
-- member on its old unscaled behavior; the column is never backfilled or inferred.
ALTER TABLE write_group_members ADD COLUMN exact_recording_source TEXT;
