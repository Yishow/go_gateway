package datalink

import (
	"database/sql"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMigratorMemberCalculationAtomicExpansion(t *testing.T) {
	for _, partial := range []bool{false, true} {
		t.Run(map[bool]string{false: "empty", true: "partial"}[partial], func(t *testing.T) {
			db := legacyMixedStorageDB(t)
			var err error
			if partial {
				_, err = db.Exec(`ALTER TABLE write_groups ADD COLUMN result_shape TEXT`)
				require.NoError(t, err)
				require.NoError(t, ensureSQLiteMemberCalculationIntent(db))
				var count int
				require.NoError(t, db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('write_group_members') WHERE name='calculation'`).Scan(&count))
				require.Equal(t, 1, count, "restore missing nullable column without inferring intent")
				return
			}
			require.NoError(t, ensureSQLiteMemberCalculationIntent(db))
			var calculation, shape sql.NullString
			require.NoError(t, db.QueryRow(`SELECT calculation FROM write_group_members`).Scan(&calculation))
			require.NoError(t, db.QueryRow(`SELECT result_shape FROM write_groups`).Scan(&shape))
			require.False(t, calculation.Valid)
			require.False(t, shape.Valid)
			_, err = db.Exec(`UPDATE write_group_members SET calculation='mean'; UPDATE write_groups SET result_shape='member-calculation-v1'`)
			require.NoError(t, err)
			require.NoError(t, ensureSQLiteMemberCalculationIntent(db))
			require.NoError(t, db.QueryRow(`SELECT calculation FROM write_group_members`).Scan(&calculation))
			require.Equal(t, "mean", calculation.String)
			for _, bad := range []string{"", "average", "unknown"} {
				_, err = db.Exec(`UPDATE write_group_members SET calculation=?`, bad)
				require.Error(t, err)
			}
		})
	}
}
