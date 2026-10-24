package dbtarget

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestWriter_PostgresMissingTableRecovery(t *testing.T) {
	dsn, fields := postgresKeywordDSN(t)
	admin := postgresAdminDB(t, dsn)
	for _, grouped := range []bool{false, true} {
		t.Run(fmt.Sprintf("grouped=%t", grouped), func(t *testing.T) {
			schemaName := postgresScopedSchema(t, admin)
			f := newWriterRecoveryFixture(t, schema.DatabaseConnectorKindPostgres,
				postgresConnectorConfig(fields, schemaName, fields["user"], fields["password"]), schemaName, grouped)
			table := pgIdentifier(t, schemaName) + `.sensor_readings`
			at := time.Date(2026, 10, 24, 0, 55, 0, 0, time.UTC)
			for cycle := range 2 {
				if cycle > 0 {
					_, err := admin.ExecContext(t.Context(), `DROP TABLE `+table)
					require.NoError(t, err)
				}
				for i, mapping := range f.mappings {
					require.NoError(t, f.writer.WriteTagValue(t.Context(), mapping.TagID, 21.5+float64(i), at))
				}
				if grouped {
					require.NoError(t, f.writer.flushBuckets(t.Context(), at.Add(time.Minute), false))
				}
				var temperature, pressure float64
				var savedAt time.Time
				require.NoError(t, admin.QueryRowContext(t.Context(), `SELECT temperature, pressure, ts FROM `+table).Scan(&temperature, &pressure, &savedAt))
				require.Equal(t, 21.5, temperature)
				require.Equal(t, 22.5, pressure)
				require.True(t, at.Equal(savedAt))
				var count int
				require.NoError(t, admin.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM `+table).Scan(&count))
				require.Equal(t, 1, count, "a recovery retry must not duplicate the sample")
			}
			updated, err := f.connectors.GetByID(t.Context(), f.connector.ID)
			require.NoError(t, err)
			require.Equal(t, "success", updated.LastSchemaEnsureStatus)
		})
	}
}

func TestWriter_PostgresMissingTableRecoveryPermissionDenied(t *testing.T) {
	dsn, fields := postgresKeywordDSN(t)
	admin := postgresAdminDB(t, dsn)
	schemaName := postgresScopedSchema(t, admin)
	roleName := fmt.Sprintf("gw_writer_%d", time.Now().UnixNano())
	// This password is an isolated test fixture, never a saved user credential.
	const password = "gw_writer_recovery_test"
	_, err := admin.ExecContext(t.Context(), `CREATE ROLE `+pgIdentifier(t, roleName)+` LOGIN PASSWORD '`+password+`'`)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = admin.ExecContext(context.Background(), `DROP OWNED BY `+pgIdentifier(t, roleName))
		_, _ = admin.ExecContext(context.Background(), `DROP ROLE `+pgIdentifier(t, roleName))
	})
	_, err = admin.ExecContext(t.Context(), `GRANT USAGE ON SCHEMA `+pgIdentifier(t, schemaName)+` TO `+pgIdentifier(t, roleName))
	require.NoError(t, err)
	f := newWriterRecoveryFixture(t, schema.DatabaseConnectorKindPostgres,
		postgresConnectorConfig(fields, schemaName, roleName, password), schemaName, false)
	err = f.writer.WriteTagValue(t.Context(), f.mappings[0].TagID, 42.0, time.Now())
	require.Error(t, err)
	require.Contains(t, err.Error(), "自動重建資料表失敗")
	var exists bool
	require.NoError(t, admin.QueryRowContext(t.Context(), `SELECT to_regclass($1) IS NOT NULL`, schemaName+`.sensor_readings`).Scan(&exists))
	require.False(t, exists)
	updated, err := f.connectors.GetByID(t.Context(), f.connector.ID)
	require.NoError(t, err)
	require.Equal(t, "failed", updated.LastSchemaEnsureStatus)
	require.Equal(t, "failed", updated.LastWriteStatus)
}
