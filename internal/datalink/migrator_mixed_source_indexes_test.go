package datalink

import (
	"database/sql"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMixedSourceLatestQueriesUseOrderedScopeIndexes(t *testing.T) {
	db, err := sql.Open("sqlite", t.TempDir()+"/evidence.db")
	require.NoError(t, err)
	defer db.Close()
	require.NoError(t, NewMigrator().Migrate(db))
	for _, tc := range []struct {
		name, filter, index string
		args                []any
	}{
		{"wide", "", "idx_wg_mixed_source_latest", []any{"w", "g", "r"}},
		{"long", ` AND json_extract(outcome,'$.EntityKey')=?`, "idx_wg_mixed_source_entity_latest", []any{"w", "g", "r", "member"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rows, err := db.Query(`EXPLAIN QUERY PLAN SELECT record_id,due_at,kind,outcome,outcome_digest FROM wg_delivery_mixed_outcomes WHERE workspace_id=? AND group_id=? AND group_revision=?`+tc.filter+` ORDER BY due_at DESC,record_id DESC LIMIT 1`, tc.args...)
			require.NoError(t, err)
			defer rows.Close()
			var details []string
			for rows.Next() {
				var a, b, c int
				var d string
				require.NoError(t, rows.Scan(&a, &b, &c, &d))
				details = append(details, d)
			}
			require.NoError(t, rows.Err())
			plan := strings.Join(details, "\n")
			require.Contains(t, plan, tc.index)
			require.NotContains(t, plan, "TEMP B-TREE")
		})
	}
}

func mixedSourceIndexFixture(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", t.TempDir()+"/indexes.db")
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	require.NoError(t, NewMigrator().Migrate(db))
	_, err = db.Exec(`DROP INDEX idx_wg_mixed_source_latest; DROP INDEX idx_wg_mixed_source_entity_latest;
 INSERT INTO wg_runtime_versions VALUES('w','g','r','{}');
 INSERT INTO wg_delivery_mixed_outcomes VALUES('w','g','r','record-a','2026-10-31T00:00:15Z','row',' {"EntityKey":"member", "legacySpacing": 1} ','hash-a','created');
 INSERT INTO wg_delivery_mixed_outcomes VALUES('w','g','r','record-b','2026-10-31T00:00:15Z','skipped','{"EntityKey":"member"}','hash-b','created');
 CREATE INDEX custom_mixed_kind ON wg_delivery_mixed_outcomes(kind);
 CREATE VIEW custom_mixed_view AS SELECT outcome FROM wg_delivery_mixed_outcomes;
 CREATE TABLE mixed_index_audit(value TEXT);
 CREATE TRIGGER custom_mixed_trigger AFTER UPDATE ON wg_delivery_mixed_outcomes BEGIN INSERT INTO mixed_index_audit VALUES(new.record_id); END;`)
	require.NoError(t, err)
	return db
}

func mixedSourceIndexRows(t *testing.T, db *sql.DB) [][]any {
	t.Helper()
	rows, err := db.Query(`SELECT * FROM wg_delivery_mixed_outcomes ORDER BY record_id`)
	require.NoError(t, err)
	defer rows.Close()
	cols, err := rows.Columns()
	require.NoError(t, err)
	var result [][]any
	for rows.Next() {
		v := make([]any, len(cols))
		p := make([]any, len(cols))
		for i := range v {
			p[i] = &v[i]
		}
		require.NoError(t, rows.Scan(p...))
		result = append(result, v)
	}
	require.NoError(t, rows.Err())
	return result
}

func TestMixedSourceIndexMigrationPreservesRowsDependenciesAndReentry(t *testing.T) {
	db := mixedSourceIndexFixture(t)
	before := mixedDeliverySchema(t, db)
	values := mixedSourceIndexRows(t, db)
	require.NoError(t, expandSQLiteMixedSourceIndexes(db, nil))
	after := mixedDeliverySchema(t, db)
	for name, ddl := range before {
		require.Equal(t, ddl, after[name], name)
	}
	require.Equal(t, values, mixedSourceIndexRows(t, db))
	require.NoError(t, NewMigrator().Migrate(db))
	require.Equal(t, after, mixedDeliverySchema(t, db))
	require.Equal(t, values, mixedSourceIndexRows(t, db))
	for _, filter := range []string{"", ` AND json_extract(outcome,'$.EntityKey')='member'`} {
		var id, due, kind, raw, digest string
		require.NoError(t, db.QueryRow(`SELECT record_id,due_at,kind,outcome,outcome_digest FROM wg_delivery_mixed_outcomes WHERE workspace_id='w' AND group_id='g' AND group_revision='r'`+filter+` ORDER BY due_at DESC,record_id DESC LIMIT 1`).Scan(&id, &due, &kind, &raw, &digest))
		require.Equal(t, "record-b", id)
	}
}

func TestMixedSourceIndexMigrationRollsBackEachStage(t *testing.T) {
	for _, stage := range []string{"idx_wg_mixed_source_latest", "idx_wg_mixed_source_entity_latest", "commit"} {
		t.Run(stage, func(t *testing.T) {
			db := mixedSourceIndexFixture(t)
			before := mixedDeliverySchema(t, db)
			values := mixedSourceIndexRows(t, db)
			err := expandSQLiteMixedSourceIndexes(db, func(at string) error {
				if at == stage {
					return errors.New("injected index fault")
				}
				return nil
			})
			require.ErrorContains(t, err, "injected index fault")
			require.Equal(t, before, mixedDeliverySchema(t, db))
			require.Equal(t, values, mixedSourceIndexRows(t, db))
			require.NoError(t, expandSQLiteMixedSourceIndexes(db, nil))
		})
	}
}

func TestMixedSourceIndexMigrationRefusesUnknownOrCorruptInput(t *testing.T) {
	for name, setup := range map[string]string{
		"unknown index":  `CREATE INDEX idx_wg_mixed_source_latest ON wg_delivery_mixed_outcomes(record_id)`,
		"name collision": `CREATE TABLE idx_wg_mixed_source_latest(value TEXT)`,
		"unknown table":  `ALTER TABLE wg_delivery_mixed_outcomes ADD COLUMN surprise TEXT`,
		"invalid json":   `PRAGMA ignore_check_constraints=ON; UPDATE wg_delivery_mixed_outcomes SET outcome='broken'; PRAGMA ignore_check_constraints=OFF`,
	} {
		t.Run(name, func(t *testing.T) {
			db := mixedSourceIndexFixture(t)
			_, err := db.Exec(setup)
			require.NoError(t, err)
			before := mixedDeliverySchema(t, db)
			values := mixedSourceIndexRows(t, db)
			require.Error(t, expandSQLiteMixedSourceIndexes(db, nil))
			require.Equal(t, before, mixedDeliverySchema(t, db))
			require.Equal(t, values, mixedSourceIndexRows(t, db))
		})
	}
}

func TestMixedSourceIndexReentryRejectsCorruptJSONEvenWithExistingIndexes(t *testing.T) {
	db := mixedSourceIndexFixture(t)
	require.NoError(t, expandSQLiteMixedSourceIndexes(db, nil))
	// Valid JSON non-object corruption is indexable but violates the retained table contract.
	_, err := db.Exec(`PRAGMA ignore_check_constraints=ON; UPDATE wg_delivery_mixed_outcomes SET outcome='[]'; PRAGMA ignore_check_constraints=OFF`)
	require.NoError(t, err)
	require.Error(t, expandSQLiteMixedSourceIndexes(db, nil))
}
