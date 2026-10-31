package datalink

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func mixedDeliveryMigrationDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", t.TempDir()+"/migration.db")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	db.SetMaxOpenConns(1)
	require.NoError(t, NewMigrator().Migrate(db))
	_, err = db.Exec(`DROP INDEX IF EXISTS idx_wg_delivery_outbox_due_head; DROP TABLE wg_delivery_mixed_progress; DROP TABLE wg_delivery_mixed_outcomes; ALTER TABLE wg_delivery_outbox DROP COLUMN due_at; PRAGMA foreign_keys=ON`)
	require.NoError(t, err)
	return db
}

func mixedDeliverySchema(t *testing.T, db *sql.DB) map[string]string {
	t.Helper()
	rows, err := db.Query(`SELECT name,sql FROM sqlite_master WHERE sql IS NOT NULL ORDER BY name`)
	require.NoError(t, err)
	defer rows.Close()
	result := map[string]string{}
	for rows.Next() {
		var name, ddl string
		require.NoError(t, rows.Scan(&name, &ddl))
		result[name] = ddl
	}
	require.NoError(t, rows.Err())
	return result
}

func mixedDeliveryLegacyValues(t *testing.T, db *sql.DB) []any {
	t.Helper()
	rows, err := db.Query(`SELECT * FROM wg_delivery_outbox WHERE effect_key='old'`)
	require.NoError(t, err)
	defer rows.Close()
	cols, err := rows.Columns()
	require.NoError(t, err)
	require.True(t, rows.Next())
	values := make([]any, len(cols))
	dest := make([]any, len(cols))
	for i := range values {
		dest[i] = &values[i]
	}
	require.NoError(t, rows.Scan(dest...))
	if cols[len(cols)-1] == "due_at" {
		require.Nil(t, values[len(values)-1])
		values = values[:len(values)-1]
	}
	return values
}

func TestMixedDeliveryMigrationPreservesLegacyAndReentry(t *testing.T) {
	db := mixedDeliveryMigrationDB(t)
	_, err := db.Exec(`INSERT INTO wg_delivery_outbox(effect_key,record_id,workspace_id,group_id,group_revision,bucket_start,partition_key,destination_scope,connector_id,connector_revision,table_name,payload,payload_digest,next_retry_at,created_at,updated_at) VALUES('old','r','w','g','v','2026-10-31T00:00:00Z','part','target','c','cr','data',' {"legacy": 1} ','digest','next','created','updated');
 CREATE INDEX custom_outbox_record ON wg_delivery_outbox(record_id);
 CREATE TABLE migration_audit(value TEXT);
 CREATE TRIGGER custom_outbox_trigger AFTER UPDATE ON wg_delivery_outbox BEGIN INSERT INTO migration_audit VALUES(new.effect_key); END;
 CREATE VIEW custom_outbox_view AS SELECT effect_key,payload FROM wg_delivery_outbox;`)
	require.NoError(t, err)
	before := mixedDeliverySchema(t, db)
	values := mixedDeliveryLegacyValues(t, db)
	require.NoError(t, expandSQLiteMixedDelivery(db, nil))
	after := mixedDeliverySchema(t, db)
	for name, ddl := range before {
		if name != "wg_delivery_outbox" {
			require.Equal(t, ddl, after[name], name)
		}
	}
	require.Equal(t, values, mixedDeliveryLegacyValues(t, db))
	require.NoError(t, expandSQLiteMixedDelivery(db, nil))
	require.Equal(t, after, mixedDeliverySchema(t, db))
	// Full migration additionally installs 052 read-side indexes. Preserve every
	// 051 object exactly, and prove full migration is itself stable on reentry.
	require.NoError(t, NewMigrator().Migrate(db))
	full := mixedDeliverySchema(t, db)
	for name, ddl := range after {
		require.Equal(t, ddl, full[name], name)
	}
	require.Len(t, full, len(after)+2)
	require.NoError(t, NewMigrator().Migrate(db))
	require.Equal(t, full, mixedDeliverySchema(t, db))
	var fk int
	require.NoError(t, db.QueryRow(`PRAGMA foreign_keys`).Scan(&fk))
	require.Equal(t, 1, fk)
	_, err = db.Exec(`UPDATE wg_delivery_outbox SET state='retry' WHERE effect_key='old'`)
	require.NoError(t, err)
	var audited string
	require.NoError(t, db.QueryRow(`SELECT value FROM migration_audit`).Scan(&audited))
	require.Equal(t, "old", audited)
}

func TestMixedDeliveryMigrationAtomicFaultRollback(t *testing.T) {
	for _, stage := range []string{"wg_delivery_mixed_progress", "wg_delivery_mixed_outcomes", "outbox", "index", "commit"} {
		t.Run(stage, func(t *testing.T) {
			db := mixedDeliveryMigrationDB(t)
			before := mixedDeliverySchema(t, db)
			injected := errors.New("migration fault")
			err := expandSQLiteMixedDelivery(db, func(at string) error {
				if at == stage {
					return injected
				}
				return nil
			})
			require.ErrorIs(t, err, injected)
			require.Equal(t, before, mixedDeliverySchema(t, db))
			var fk int
			require.NoError(t, db.QueryRow(`PRAGMA foreign_keys`).Scan(&fk))
			require.Equal(t, 1, fk)
			require.NoError(t, expandSQLiteMixedDelivery(db, nil))
		})
	}
}

func TestMixedDeliveryMigrationRejectsUnknownSchema(t *testing.T) {
	for _, alter := range []string{
		`CREATE TABLE wg_delivery_mixed_progress(workspace_id TEXT)`,
		`CREATE TABLE wg_delivery_mixed_outcomes(record_id TEXT)`,
		`ALTER TABLE wg_delivery_outbox ADD COLUMN due_at INTEGER`,
		`ALTER TABLE wg_delivery_outbox ADD COLUMN due_at TEXT DEFAULT ''`,
		`ALTER TABLE wg_delivery_outbox ADD COLUMN mystery TEXT`,
		`CREATE INDEX idx_wg_delivery_outbox_due_head ON wg_delivery_outbox(partition_key,bucket_start,effect_key)`,
		`ALTER TABLE wg_delivery_outbox ADD COLUMN due_at TEXT; CREATE INDEX idx_wg_delivery_outbox_due_head ON wg_delivery_outbox(partition_key,COALESCE(bucket_start,due_at),effect_key)`,
		`ALTER TABLE wg_delivery_outbox ADD COLUMN due_at TEXT; CREATE INDEX idx_wg_delivery_outbox_due_head ON wg_delivery_outbox(partition_key,COALESCE(due_at,bucket_start),effect_key DESC)`,
	} {
		t.Run(alter, func(t *testing.T) {
			db := mixedDeliveryMigrationDB(t)
			_, err := db.Exec(alter)
			require.NoError(t, err)
			before := mixedDeliverySchema(t, db)
			require.ErrorContains(t, expandSQLiteMixedDelivery(db, nil), "unknown")
			require.Equal(t, before, mixedDeliverySchema(t, db))
		})
	}
}

func TestMixedDeliveryMigrationForeignKeysAndJSON(t *testing.T) {
	db := mixedDeliveryMigrationDB(t)
	require.NoError(t, expandSQLiteMixedDelivery(db, nil))
	_, err := db.Exec(`INSERT INTO wg_runtime_versions VALUES('w','g','r','{}')`)
	require.NoError(t, err)
	progress := `INSERT INTO wg_delivery_mixed_progress VALUES(?,?,?,?,?,?,?,?,?,?)`
	args := []any{"w", "g", "r", 1, "descriptor", "", "{}", "checkpoint", "closure", "now"}
	_, err = db.Exec(progress, args...)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO wg_delivery_mixed_outcomes VALUES('w','g','r','record','due','row','{}','digest','now')`)
	require.NoError(t, err)
	_, err = db.Exec(`DELETE FROM wg_runtime_versions WHERE workspace_id='w'`)
	require.Error(t, err, "no implicit cascade")
	for _, mutation := range []struct {
		index int
		value any
	}{{0, "foreign"}, {3, 2}, {6, "[]"}, {6, "bad-json"}} {
		_, err = db.Exec(`DELETE FROM wg_delivery_mixed_progress`)
		require.NoError(t, err)
		bad := append([]any(nil), args...)
		bad[mutation.index] = mutation.value
		_, err = db.Exec(progress, bad...)
		require.Error(t, err, fmt.Sprint(mutation))
	}
	for _, query := range []string{
		`INSERT INTO wg_delivery_mixed_outcomes VALUES('foreign','g','r','foreign','due','row','{}','digest','now')`,
		`INSERT INTO wg_delivery_mixed_outcomes VALUES('w','g','r','bad','due','row','[]','digest','now')`,
	} {
		_, err = db.Exec(query)
		require.Error(t, err)
	}
	require.NoError(t, expandSQLiteMixedDelivery(db, nil))
	// Existing corruption must fail closed even when schemas are already known.
	_, err = db.Exec(`PRAGMA foreign_keys=OFF; INSERT INTO wg_delivery_mixed_outcomes VALUES('orphan','g','r','bad','due','row','{}','digest','now'); PRAGMA foreign_keys=ON`)
	require.NoError(t, err)
	before := mixedDeliverySchema(t, db)
	require.ErrorContains(t, expandSQLiteMixedDelivery(db, nil), "foreign key check")
	require.True(t, reflect.DeepEqual(before, mixedDeliverySchema(t, db)))
}

func TestMixedDeliveryMigrationReopenPreservesMixedRows(t *testing.T) {
	db := mixedDeliveryMigrationDB(t)
	require.NoError(t, expandSQLiteMixedDelivery(db, nil))
	_, err := db.Exec(`INSERT INTO wg_runtime_versions VALUES('w','g','r','{"descriptor":true}');
 INSERT INTO wg_delivery_mixed_progress VALUES('w','g','r',1,'descriptor','','{"version":1}','checkpoint','closure','now');
 INSERT INTO wg_delivery_mixed_outcomes VALUES('w','g','r','record','due','skipped','{"reason":"baseline"}','digest','now');`)
	require.NoError(t, err)
	var path string
	require.NoError(t, db.QueryRow(`SELECT file FROM pragma_database_list WHERE name='main'`).Scan(&path))
	require.NoError(t, db.Close())
	reopened, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	defer reopened.Close()
	reopened.SetMaxOpenConns(1)
	require.NoError(t, NewMigrator().Migrate(reopened))
	var checkpoint, outcome string
	require.NoError(t, reopened.QueryRow(`SELECT checkpoint FROM wg_delivery_mixed_progress`).Scan(&checkpoint))
	require.NoError(t, reopened.QueryRow(`SELECT outcome FROM wg_delivery_mixed_outcomes`).Scan(&outcome))
	require.Equal(t, `{"version":1}`, checkpoint)
	require.Equal(t, `{"reason":"baseline"}`, outcome)
}

func TestMixedDeliveryMigrationHeadExpressionIndex(t *testing.T) {
	db := mixedDeliveryMigrationDB(t)
	require.NoError(t, expandSQLiteMixedDelivery(db, nil))
	var indexSQL string
	require.NoError(t, db.QueryRow(`SELECT sql FROM sqlite_master WHERE name='idx_wg_delivery_outbox_due_head'`).Scan(&indexSQL))
	const head = `SELECT o.partition_key,o.effect_key FROM wg_delivery_outbox o
 WHERE o.state IN ('pending','retrying') AND o.next_retry_at <= ?
 AND NOT EXISTS (SELECT 1 FROM wg_delivery_outbox e WHERE e.partition_key=o.partition_key
 AND (COALESCE(e.due_at,e.bucket_start)<COALESCE(o.due_at,o.bucket_start) OR (COALESCE(e.due_at,e.bucket_start)=COALESCE(o.due_at,o.bucket_start) AND e.effect_key<o.effect_key))
 AND e.state NOT IN ('sql_committed','operator_skipped')) ORDER BY COALESCE(o.due_at,o.bucket_start),o.effect_key`
	for _, row := range []struct {
		id, bucket string
		due        any
	}{{"legacy-early", "10", nil}, {"mixed-middle", "", "15"}, {"legacy-last", "20", nil}} {
		_, err := db.Exec(`INSERT INTO wg_delivery_outbox(effect_key,record_id,workspace_id,group_id,group_revision,bucket_start,due_at,partition_key,destination_scope,connector_id,connector_revision,table_name,payload,payload_digest,next_retry_at,created_at,updated_at) VALUES(?,'r','w','g','v',?,?,'part','target','c','cr','data','{}','digest','0','created','updated')`, row.id, row.bucket, row.due)
		require.NoError(t, err)
	}
	rows, err := db.Query(`EXPLAIN QUERY PLAN `+head, "z")
	require.NoError(t, err)
	var plans []string
	for rows.Next() {
		var id, parent, unused int
		var detail string
		require.NoError(t, rows.Scan(&id, &parent, &unused, &detail))
		plans = append(plans, detail)
	}
	require.NoError(t, rows.Err())
	require.NoError(t, rows.Close())
	require.Contains(t, strings.Join(plans, "\n"), "SEARCH e USING INDEX idx_wg_delivery_outbox_due_head", strings.Join(plans, "\n"))
	t.Logf("query plan:\n%s", strings.Join(plans, "\n"))
	for _, want := range []string{"legacy-early", "mixed-middle", "legacy-last"} {
		var partition, effect string
		require.NoError(t, db.QueryRow(head, "z").Scan(&partition, &effect))
		require.Equal(t, want, effect)
		_, err = db.Exec(`UPDATE wg_delivery_outbox SET state='sql_committed' WHERE effect_key=?`, effect)
		require.NoError(t, err)
	}
	before := mixedDeliverySchema(t, db)
	require.NoError(t, expandSQLiteMixedDelivery(db, nil))
	require.Equal(t, before, mixedDeliverySchema(t, db))
}
