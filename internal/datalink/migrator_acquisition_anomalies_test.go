package datalink

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAcquisitionAnomaliesMigrationIsIdempotentAndEnforcesOneActive(t *testing.T) {
	db := openDeviceAcquisitionDB(t)
	require.NoError(t, NewMigrator().Migrate(db), "re-running the migration is a no-op")
	seedDeviceAcquisition(t, db,
		`INSERT INTO devices (id, name, protocol, status, connection_config) VALUES ('dev', 'dev', 'modbus_tcp', 'draft', '{}')`,
		`INSERT INTO acquisition_anomalies (device_id, point_id, kind, started_at) VALUES ('dev', 'p', 'stale', 1)`,
		`INSERT INTO acquisition_anomalies (device_id, kind, started_at) VALUES ('dev', 'disconnected', 1)`,
	)
	_, err := db.ExecContext(t.Context(), `INSERT INTO acquisition_anomalies (device_id, point_id, kind, started_at) VALUES ('dev', 'p', 'stale', 2)`)
	require.ErrorContains(t, err, "UNIQUE", "one active event per device, point and kind")
	_, err = db.ExecContext(t.Context(), `INSERT INTO acquisition_anomalies (device_id, kind, started_at) VALUES ('dev', 'disconnected', 2)`)
	require.ErrorContains(t, err, "UNIQUE", "device-level events use the same uniqueness")
	_, err = db.ExecContext(t.Context(), `INSERT INTO acquisition_anomalies (device_id, kind, started_at) VALUES ('dev', 'stale', 2)`)
	require.ErrorContains(t, err, "CHECK", "point events require a point")
}

func TestAcquisitionAnomaliesMigrationRefusesUnknownSchema(t *testing.T) {
	db := openDeviceAcquisitionDB(t)
	seedDeviceAcquisition(t, db,
		`DROP TABLE acquisition_anomalies`,
		`CREATE TABLE acquisition_anomalies (id INTEGER PRIMARY KEY, device_id TEXT)`,
	)
	require.ErrorContains(t, NewMigrator().Migrate(db), "unknown acquisition_anomalies schema")
}
