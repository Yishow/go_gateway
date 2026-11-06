ALTER TABLE polling_groups ADD COLUMN owner_device_id TEXT REFERENCES devices(id) ON DELETE CASCADE CHECK (owner_device_id IS NULL OR (typeof(owner_device_id) = 'text' AND length(owner_device_id) > 0));
CREATE UNIQUE INDEX idx_polling_groups_owner_device ON polling_groups (owner_device_id) WHERE owner_device_id IS NOT NULL;
ALTER TABLE devices ADD COLUMN stale_threshold_ms INTEGER CHECK (stale_threshold_ms IS NULL OR (typeof(stale_threshold_ms) = 'integer' AND stale_threshold_ms > 0));
ALTER TABLE devices ADD COLUMN frozen_threshold_ms INTEGER NOT NULL DEFAULT 0 CHECK (typeof(frozen_threshold_ms) = 'integer' AND frozen_threshold_ms >= 0);
ALTER TABLE devices ADD COLUMN acquisition_auto_interval INTEGER NOT NULL DEFAULT 0 CHECK (typeof(acquisition_auto_interval) = 'integer' AND acquisition_auto_interval IN (0, 1));
