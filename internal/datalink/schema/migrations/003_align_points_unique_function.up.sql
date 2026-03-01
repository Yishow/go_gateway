-- Align points uniqueness from (device_id, address) to (device_id, address, function)

ALTER TABLE points ALTER COLUMN function SET DEFAULT '';
UPDATE points SET function = '' WHERE function IS NULL;
ALTER TABLE points ALTER COLUMN function SET NOT NULL;

ALTER TABLE points DROP CONSTRAINT IF EXISTS points_device_address_unique;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'points_device_address_function_unique'
          AND conrelid = 'points'::regclass
    ) THEN
        ALTER TABLE points
            ADD CONSTRAINT points_device_address_function_unique UNIQUE (device_id, address, function);
    END IF;
END
$$;
