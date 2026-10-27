-- Additive, nullable exact recording scale. NULL/empty keeps every existing
-- definition on its old behavior; the column is never backfilled or inferred.
ALTER TABLE measurement_definitions ADD COLUMN exact_recording_scale TEXT;
