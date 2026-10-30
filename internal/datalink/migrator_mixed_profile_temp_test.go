package datalink

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMixedProfileMigrationRestoresConnectionIdentity(t *testing.T) {
	for _, store := range []int{0, 1, 2} {
		for _, preexisting := range []bool{false, true} {
			for _, fail := range []bool{false, true} {
				t.Run(fmt.Sprintf("store%d/temp%t/fail%t", store, preexisting, fail), func(t *testing.T) {
					db := legacyMixedStorageDB(t)
					_, err := db.Exec(fmt.Sprintf("PRAGMA temp_store=%d", store))
					require.NoError(t, err)
					if preexisting {
						_, err = db.Exec(`CREATE TEMP TABLE keep_my_data(value TEXT); INSERT INTO keep_my_data VALUES('unchanged')`)
						require.NoError(t, err)
					}
					injected := errors.New("rollback after rename")
					err = expandSQLiteMixedProfile(db, func(at string) error {
						if fail && at == "indexes" {
							return injected
						}
						return nil
					})
					if fail {
						require.ErrorIs(t, err, injected)
					} else {
						require.NoError(t, err)
					}
					conn, err := db.Conn(t.Context())
					require.NoError(t, err)
					defer conn.Close()
					hasTemp, err := sqliteMigrationHasTemp(t.Context(), conn)
					require.NoError(t, err)
					require.Equal(t, preexisting, hasTemp)
					var actual, fk int
					require.NoError(t, conn.QueryRowContext(t.Context(), "PRAGMA temp_store").Scan(&actual))
					require.Equal(t, store, actual)
					require.NoError(t, conn.QueryRowContext(t.Context(), "PRAGMA foreign_keys").Scan(&fk))
					require.Equal(t, 1, fk)
					if preexisting {
						var value string
						require.NoError(t, conn.QueryRowContext(t.Context(), "SELECT value FROM temp.keep_my_data").Scan(&value))
						require.Equal(t, "unchanged", value)
					}
				})
			}
		}
	}
}

func TestMixedProfileTempCleanupRefusesUnexpectedObjects(t *testing.T) {
	db := legacyMixedStorageDB(t)
	conn, err := db.Conn(t.Context())
	require.NoError(t, err)
	defer conn.Close()
	_, err = conn.ExecContext(t.Context(), `CREATE TEMP TABLE unexpected(value TEXT); INSERT INTO unexpected VALUES('must not delete')`)
	require.NoError(t, err)
	require.ErrorContains(t, restoreSQLiteMigrationTemp(t.Context(), conn, false, 0), "unexpectedly created temporary objects")
	var value string
	require.NoError(t, conn.QueryRowContext(t.Context(), "SELECT value FROM temp.unexpected").Scan(&value))
	require.Equal(t, "must not delete", value)
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	require.Error(t, restoreSQLiteMigrationTemp(ctx, conn, false, 0), "uncertain state is an error, not successful cleanup")
	require.NoError(t, conn.QueryRowContext(t.Context(), "SELECT value FROM temp.unexpected").Scan(&value))
}
