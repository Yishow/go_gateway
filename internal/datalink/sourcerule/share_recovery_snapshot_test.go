package sourcerule

import (
	"database/sql"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/modbusshare"
	"testing"
)

func shareRecoveryFixture(t *testing.T) (*sql.DB, modbusshare.Settings) {
	t.Helper()
	db := setupSQLRepoDB(t)
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { db.Close() })
	for _, q := range []string{
		`INSERT INTO devices(id,name,protocol,status) VALUES('d','d','modbus_tcp','active')`,
		`INSERT INTO points(id,device_id,name,address,data_type) VALUES('p','d','p','40001','uint16')`,
		`INSERT INTO tags(id,key,key_lower,display_name,data_type,status) VALUES('t','tag','tag','Tag','uint16','active')`,
		`INSERT INTO mappings(id,point_id,tag_id) VALUES('m','p','t')`,
		`INSERT INTO source_rules(id,device_id,start_address,count,data_type,naming_prefix,revision_id,share_enabled,share_stride) VALUES('r','d','40001',1,'uint16','p','rev',1,1)`,
		`INSERT INTO source_rule_links(id,rule_id,address,point_id,tag_id,mapping_id) VALUES('l','r','40001','p','t','m')`,
		`INSERT INTO source_rule_candidate_snapshots(source_rule_id,revision_id,candidate_type,payload,status,reason,generated_at) VALUES('r','rev','local_modbus_outputs','{"candidates":[{"point_id":"p","tag_id":"t","mapping_id":"m","register":0,"data_type":"uint16","status":"ready"}]}','ready','','2026-01-01T00:00:00Z')`,
		`INSERT INTO system_settings(key,value) VALUES('studio_v2_workspace','{"id":"w","ordered_device_ids":["d"]}')`,
	} {
		_, err := db.Exec(q)
		require.NoError(t, err)
	}
	settings := modbusshare.DefaultSettings()
	settings.Enabled = true
	settings.SettingsRevision = "set-1"
	raw, err := json.Marshal(settings)
	require.NoError(t, err)
	_, err = db.Exec(`UPDATE system_settings SET value=? WHERE key='modbus_share'`, string(raw))
	require.NoError(t, err)
	return db, settings
}

func TestShareRecoverySnapshotInputWitness(t *testing.T) {
	for _, kind := range []string{"unchanged", "candidate", "link", "tag", "mapping", "membership", "settings", "new-rule", "missing-candidate", "foreign-point-device"} {
		t.Run(kind, func(t *testing.T) {
			db, settings := shareRecoveryFixture(t)
			tx, err := db.Begin()
			require.NoError(t, err)
			initial, err := ReadShareRecoverySnapshot(t.Context(), tx, "w", settings)
			require.NoError(t, err)
			require.Len(t, initial.Desired, 1)
			require.Len(t, initial.WitnessDigest, 64)
			require.NoError(t, tx.Rollback())
			var query string
			switch kind {
			case "candidate":
				query = `UPDATE source_rule_candidate_snapshots SET reason='changed'`
			case "link":
				query = `UPDATE source_rule_links SET address='40002'`
			case "tag":
				query = `UPDATE tags SET display_name='changed'`
			case "mapping":
				query = `UPDATE mappings SET transform_pipeline='[{"type":"scale","params":{"factor":2}}]'`
			case "membership":
				query = `UPDATE system_settings SET value='{"id":"w","ordered_device_ids":[]}' WHERE key='studio_v2_workspace'`
			case "settings":
				query = `UPDATE system_settings SET value=json_set(value,'$.port',5029) WHERE key='modbus_share'`
			case "new-rule":
				query = `INSERT INTO source_rules(id,device_id,start_address,count,data_type,naming_prefix,revision_id,share_enabled) VALUES('new','d','40002',1,'uint16','p','newrev',0)`
			case "missing-candidate":
				query = `DELETE FROM source_rule_candidate_snapshots`
			case "foreign-point-device":
				_, err = db.Exec(`INSERT INTO devices(id,name,protocol,status) VALUES('foreign','foreign','modbus_tcp','active')`)
				require.NoError(t, err)
				query = `UPDATE points SET device_id='foreign'`
			}
			if query != "" {
				_, err = db.Exec(query)
				require.NoError(t, err)
			}
			tx, err = db.Begin()
			require.NoError(t, err)
			defer tx.Rollback()
			err = ValidateShareRecoverySnapshot(t.Context(), tx, "w", settings, initial.WitnessDigest)
			if kind == "unchanged" {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, ErrShareRecoverySnapshot)
			}
		})
	}
}

func TestShareRecoverySnapshotRejectsIncompleteCandidate(t *testing.T) {
	for _, payload := range []string{`{}`, `{"candidates":null}`, `{"candidates":[{"point_id":"p","data_type":"uint16","status":"ready"}]}`, `{"candidates":[{"point_id":"p","tag_id":"t","register":0,"data_type":"uint16","status":"out_of_sync"}]}`, `{"candidates":[{"point_id":"p","tag_id":"t","register":0,"data_type":"uint16","status":"deferred"}]}`} {
		t.Run(payload, func(t *testing.T) {
			db, settings := shareRecoveryFixture(t)
			_, err := db.Exec(`UPDATE source_rule_candidate_snapshots SET payload=?`, payload)
			require.NoError(t, err)
			tx, err := db.Begin()
			require.NoError(t, err)
			defer tx.Rollback()
			_, err = ReadShareRecoverySnapshot(t.Context(), tx, "w", settings)
			require.ErrorIs(t, err, ErrShareRecoverySnapshot)
		})
	}
}

func TestShareRecoverySnapshotAllowsExplicitEmptyCandidates(t *testing.T) {
	db, settings := shareRecoveryFixture(t)
	_, err := db.Exec(`UPDATE source_rule_candidate_snapshots SET payload='{"candidates":[]}'`)
	require.NoError(t, err)
	tx, err := db.Begin()
	require.NoError(t, err)
	defer tx.Rollback()
	evidence, err := ReadShareRecoverySnapshot(t.Context(), tx, "w", settings)
	require.NoError(t, err)
	require.Empty(t, evidence.Desired)
}

func TestShareRecoverySnapshotRejectsStaleSettingsInput(t *testing.T) {
	db, settings := shareRecoveryFixture(t)
	settings.CapacityRegisters--
	tx, err := db.Begin()
	require.NoError(t, err)
	defer tx.Rollback()
	_, err = ReadShareRecoverySnapshot(t.Context(), tx, "w", settings)
	require.ErrorIs(t, err, ErrShareRecoverySnapshot)
}

func TestShareRecoverySnapshotWitnessesAbsentTag(t *testing.T) {
	db, settings := shareRecoveryFixture(t)
	settings.Enabled = false
	raw, err := json.Marshal(settings)
	require.NoError(t, err)
	_, err = db.Exec(`UPDATE system_settings SET value=? WHERE key='modbus_share'`, string(raw))
	require.NoError(t, err)
	_, err = db.Exec(`UPDATE source_rule_candidate_snapshots SET payload='{"candidates":[{"point_id":"p","tag_id":"absent","data_type":"uint16","status":"deferred"}]}'`)
	require.NoError(t, err)
	tx, err := db.Begin()
	require.NoError(t, err)
	evidence, err := ReadShareRecoverySnapshot(t.Context(), tx, "w", settings)
	require.NoError(t, err)
	require.Empty(t, evidence.Desired)
	require.NoError(t, tx.Rollback())
	_, err = db.Exec(`INSERT INTO tags(id,key,key_lower,display_name,data_type) VALUES('absent','absent','absent','Absent','uint16')`)
	require.NoError(t, err)
	tx, err = db.Begin()
	require.NoError(t, err)
	defer tx.Rollback()
	require.ErrorIs(t, ValidateShareRecoverySnapshot(t.Context(), tx, "w", settings, evidence.WitnessDigest), ErrShareRecoverySnapshot)
}

func TestShareRecoverySnapshotRejectsBlockedSnapshot(t *testing.T) {
	db, settings := shareRecoveryFixture(t)
	_, err := db.Exec(`UPDATE source_rule_candidate_snapshots SET status='blocked',reason='repair-required'`)
	require.NoError(t, err)
	tx, err := db.Begin()
	require.NoError(t, err)
	defer tx.Rollback()
	_, err = ReadShareRecoverySnapshot(t.Context(), tx, "w", settings)
	require.ErrorIs(t, err, ErrShareRecoverySnapshot)
}
