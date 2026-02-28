-- Rollback points uniqueness from (device_id, address, function) to (device_id, address)
-- Precheck: rollback only works when each device+address still has at most one row.

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM points
        GROUP BY device_id, address
        HAVING COUNT(*) > 1
    ) THEN
        RAISE EXCEPTION 'rollback blocked: duplicate device_id+address exists in points';
    END IF;
END
$$;

ALTER TABLE points DROP CONSTRAINT IF EXISTS points_device_address_function_unique;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'points_device_address_unique'
          AND conrelid = 'points'::regclass
    ) THEN
        ALTER TABLE points
            ADD CONSTRAINT points_device_address_unique UNIQUE (device_id, address);
    END IF;
END
$$;

UPDATE points SET function = NULL WHERE function = '';
ALTER TABLE points ALTER COLUMN function DROP NOT NULL;
ALTER TABLE points ALTER COLUMN function DROP DEFAULT;
