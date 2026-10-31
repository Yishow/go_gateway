package datalink

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema/migrations"
)

func mixedFenceLegacyDB(t *testing.T) *sql.DB {
	db := mixedDeliveryMigrationDB(t)
	_, err := db.Exec(`DROP TABLE write_group_ownership_fences`)
	require.NoError(t, err)
	raw, err := migrations.FS.ReadFile("044_write_group_ownership_fences_sqlite.up.sql")
	require.NoError(t, err)
	_, err = db.Exec(string(raw))
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO write_group_cutovers(id,workspace_id,kind,source_group_id,boundary,review_digest,detail) VALUES('cut','w','forward','g','boundary','review','{}')`)
	require.NoError(t, err)
	for i, profile := range []string{"", "instantaneous", "cumulative"} {
		_, err = db.Exec(`INSERT INTO write_group_ownership_fences VALUES('cut','w','ended','g','r','d',?,'tag',?,'c','cr','db','schema','table trailing  ')`, fmt.Sprint(i), profile)
		require.NoError(t, err)
	}
	return db
}
func mixedFenceRowBytes(t *testing.T, db *sql.DB) string {
	t.Helper()
	columns := strings.Fields("fence_revision workspace_id direction group_id group_revision device_id point_id tag_id profile_kind connector_id connector_revision database_name table_schema table_name")
	var fields []string
	for _, col := range columns {
		fields = append(fields, col, "typeof("+col+")")
	}
	var raw string
	require.NoError(t, db.QueryRow("SELECT json_group_array(json_array("+strings.Join(fields, ",")+")) FROM (SELECT * FROM write_group_ownership_fences ORDER BY point_id)").Scan(&raw))
	return raw
}
func TestMixedFenceExpansionPreservesRowsDependenciesAndConstraints(t *testing.T) {
	db := mixedFenceLegacyDB(t)
	_, err := db.Exec(`CREATE INDEX custom_fence_index ON write_group_ownership_fences(table_name);
 CREATE TABLE fence_audit(value TEXT);
 CREATE TRIGGER custom_fence_trigger AFTER UPDATE ON write_group_ownership_fences BEGIN INSERT INTO fence_audit VALUES(new.table_name); END;
 CREATE VIEW fence_view AS SELECT * FROM write_group_ownership_fences;
 CREATE TABLE fence_ref(f TEXT,direction TEXT,g TEXT,r TEXT,d TEXT,p TEXT,t TEXT,FOREIGN KEY(f,direction,g,r,d,p,t) REFERENCES write_group_ownership_fences(fence_revision,direction,group_id,group_revision,device_id,point_id,tag_id));
 INSERT INTO fence_ref VALUES('cut','ended','g','r','d','0','tag');`)
	require.NoError(t, err)
	before := mixedDeliverySchema(t, db)
	values := mixedFenceRowBytes(t, db)
	require.NoError(t, expandSQLiteMixedDelivery(db, nil))
	after := mixedDeliverySchema(t, db)
	for name, ddl := range before {
		if name != "write_group_ownership_fences" && name != "wg_delivery_outbox" {
			require.Equal(t, ddl, after[name], name)
		}
	}
	require.Equal(t, values, mixedFenceRowBytes(t, db))
	require.NoError(t, expandSQLiteMixedDelivery(db, nil))
	require.Equal(t, after, mixedDeliverySchema(t, db))
	_, err = db.Exec(`INSERT INTO write_group_ownership_fences VALUES('cut','w','ended','g','r','d','mixed','tag','mixed','c','cr','db','schema','mixed-table')`)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO write_group_ownership_fences VALUES('cut','w','ended','other','r','d','mixed','tag','mixed','c','cr','db','schema','mixed-table')`)
	require.Error(t, err, "existing UNIQUE authority preserved")
	_, err = db.Exec(`DELETE FROM write_group_cutovers WHERE id='cut'`)
	require.Error(t, err, "existing RESTRICT FK preserved")
	_, err = db.Exec(`UPDATE write_group_ownership_fences SET table_name='new' WHERE point_id='mixed'`)
	require.NoError(t, err)
	var audit string
	require.NoError(t, db.QueryRow(`SELECT value FROM fence_audit`).Scan(&audit))
	require.Equal(t, "new", audit)
	var count int
	require.NoError(t, db.QueryRow(`SELECT count(*) FROM fence_view`).Scan(&count))
	require.Equal(t, 4, count)
}
func TestMixedFenceExpansionRollsBackEveryStage(t *testing.T) {
	for _, stage := range []string{"fence_copy", "fence_replace", "fence_indexes", "fence_foreign_keys", "fence_commit"} {
		t.Run(stage, func(t *testing.T) {
			db := mixedFenceLegacyDB(t)
			before := mixedDeliverySchema(t, db)
			values := mixedFenceRowBytes(t, db)
			injected := errors.New("fence-fault")
			require.ErrorIs(t, expandSQLiteMixedDelivery(db, func(at string) error {
				if at == stage {
					return injected
				}
				return nil
			}), injected)
			require.Equal(t, before, mixedDeliverySchema(t, db))
			require.Equal(t, values, mixedFenceRowBytes(t, db))
			var fk int
			require.NoError(t, db.QueryRow(`PRAGMA foreign_keys`).Scan(&fk))
			require.Equal(t, 1, fk)
			require.NoError(t, expandSQLiteMixedDelivery(db, nil))
		})
	}
}
func TestMixedFenceExpansionRejectsUnknownAndForeignCorruption(t *testing.T) {
	for _, alter := range []string{`ALTER TABLE write_group_ownership_fences ADD COLUMN unknown TEXT`, `CREATE TABLE write_group_ownership_fences_mixed_expand(value TEXT)`, `PRAGMA foreign_keys=OFF; UPDATE write_group_ownership_fences SET fence_revision='missing'; PRAGMA foreign_keys=ON`} {
		t.Run(alter, func(t *testing.T) {
			db := mixedFenceLegacyDB(t)
			_, err := db.Exec(alter)
			require.NoError(t, err)
			before := mixedDeliverySchema(t, db)
			values := mixedFenceRowBytes(t, db)
			require.Error(t, expandSQLiteMixedDelivery(db, nil))
			require.Equal(t, before, mixedDeliverySchema(t, db))
			require.Equal(t, values, mixedFenceRowBytes(t, db))
		})
	}
}
