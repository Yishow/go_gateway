package point

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
)

func recoveryBindingsFixture(t *testing.T) (*sql.DB, *schema.Point) {
	t.Helper()
	db, svc, req := formatFixture(t)
	for _, q := range []string{
		`INSERT INTO tags(id,key,key_lower,display_name,data_type,status,labels) VALUES('target','target','target','Target','uint32','active','{"owner":"initial"}')`,
		`INSERT INTO mappings(id,point_id,tag_id,transform_pipeline,enabled) SELECT 'binding',id,'target','[]',1 FROM points`,
	} {
		_, err := db.Exec(q)
		require.NoError(t, err)
	}
	p, err := svc.GetByID(t.Context(), req.PointID)
	require.NoError(t, err)
	return db, p
}
func readRecoveryBindings(t *testing.T, db *sql.DB, p *schema.Point) *FormatRecoveryTargetInputs {
	t.Helper()
	tx, err := db.Begin()
	require.NoError(t, err)
	defer tx.Rollback()
	got, err := ReadFormatRecoveryTargetInputs(t.Context(), tx, p)
	require.NoError(t, err)
	return got
}
func TestFormatRecoveryBindingsWitness(t *testing.T) {
	for _, kind := range []string{"unchanged", "telemetry-time", "mapping-transform", "mapping-disabled", "mapping-status", "mapping-add", "mapping-delete", "tag-type", "tag-labels", "tag-rename", "tag-keylower", "device-name", "device-connection"} {
		t.Run(kind, func(t *testing.T) {
			db, p := recoveryBindingsFixture(t)
			before := readRecoveryBindings(t, db, p)
			require.Len(t, before.Mappings, 1)
			require.Len(t, before.Tags, 1)
			require.Len(t, before.WitnessDigest, 64)
			query := ""
			switch kind {
			case "telemetry-time":
				query = `UPDATE points SET last_value='99',last_error='telemetry',last_read_at='2026-01-02'; UPDATE devices SET last_test_error='telemetry',last_test_success=0; UPDATE mappings SET updated_at='2020-01-01'; UPDATE tags SET updated_at='2020-01-01'`
			case "mapping-transform":
				query = `UPDATE mappings SET transform_pipeline='[{"type":"scale","params":{"factor":2}}]'`
			case "mapping-disabled":
				query = `UPDATE mappings SET enabled=0`
			case "mapping-status":
				query = `UPDATE mappings SET status='disabled'`
			case "mapping-add":
				query = `INSERT INTO tags(id,key,key_lower,display_name,data_type) VALUES('target2','target2','target2','Target2','uint32'); INSERT INTO mappings(id,point_id,tag_id,enabled) SELECT 'binding2',id,'target2',0 FROM points`
			case "mapping-delete":
				query = `DELETE FROM mappings`
			case "tag-type":
				query = `UPDATE tags SET data_type='float64'`
			case "tag-labels":
				query = `UPDATE tags SET labels='{"owner":"new"}'`
			case "tag-rename":
				query = `UPDATE tags SET display_name='New Target'`
			case "tag-keylower":
				query = `UPDATE tags SET key_lower='changed'`
			case "device-name":
				query = `UPDATE devices SET name='New Device'`
			case "device-connection":
				query = `UPDATE devices SET connection_config='{"password":"PRIVATE_CONNECTION","data_format":"CDAB"}'`
			}
			if query != "" {
				_, err := db.Exec(query)
				require.NoError(t, err)
			}
			after := readRecoveryBindings(t, db, p)
			if kind == "unchanged" || kind == "telemetry-time" {
				require.Equal(t, before.WitnessDigest, after.WitnessDigest)
			} else {
				require.NotEqual(t, before.WitnessDigest, after.WitnessDigest)
			}
			if kind == "mapping-add" {
				require.Len(t, after.Mappings, 2)
				require.False(t, after.Mappings[1].Enabled)
			}
			if kind == "mapping-delete" {
				require.NotNil(t, after.Mappings)
				require.Empty(t, after.Mappings)
				require.Empty(t, after.Tags)
			}
			raw, err := json.Marshal(after)
			require.NoError(t, err)
			require.JSONEq(t, `{}`, string(raw))
			require.NotContains(t, string(raw), "PRIVATE_CONNECTION")
		})
	}
}
func TestFormatRecoveryBindingsMissingTagFailsClosed(t *testing.T) {
	db, p := recoveryBindingsFixture(t)
	_, err := db.Exec(`PRAGMA foreign_keys=OFF; DELETE FROM tags`)
	require.NoError(t, err)
	tx, err := db.Begin()
	require.NoError(t, err)
	defer tx.Rollback()
	got, err := ReadFormatRecoveryTargetInputs(t.Context(), tx, p)
	require.ErrorIs(t, err, ErrFormatUnavailable)
	require.Nil(t, got)
}

func TestFormatRecoveryBindingsRejectsForeignDeviceEvidence(t *testing.T) {
	db, p := recoveryBindingsFixture(t)
	_, err := db.Exec(`INSERT INTO devices(id,name,protocol,status) VALUES('foreign','Foreign','modbus_tcp','active')`)
	require.NoError(t, err)
	p.DeviceID = "foreign"
	tx, err := db.Begin()
	require.NoError(t, err)
	defer tx.Rollback()
	got, err := ReadFormatRecoveryTargetInputs(t.Context(), tx, p)
	require.Error(t, err)
	require.Nil(t, got)
}

func TestFormatRecoveryAcceptanceRejectsChangedTargetBindings(t *testing.T) {
	for _, change := range []string{"insert", "transform", "tag", "delete"} {
		t.Run(change, func(t *testing.T) {
			db, req, evidence := recoveryAcceptanceFixture(t)
			if change != "insert" {
				_, err := db.Exec(`INSERT INTO tags(id,key,key_lower,display_name,data_type) VALUES('target','target','target','Target','uint32'); INSERT INTO mappings(id,point_id,tag_id,enabled) SELECT 'binding',id,'target',1 FROM points`)
				require.NoError(t, err)
				evidence, err = NewSQLRepository(db).ReadFormatRecovery(t.Context(), req)
				require.NoError(t, err)
			}
			query := map[string]string{
				"insert":    `INSERT INTO tags(id,key,key_lower,display_name,data_type) VALUES('target','target','target','Target','uint32'); INSERT INTO mappings(id,point_id,tag_id,enabled) SELECT 'binding',id,'target',1 FROM points`,
				"transform": `UPDATE mappings SET transform_pipeline='[{"type":"scale","params":{"factor":2}}]'`,
				"tag":       `UPDATE tags SET display_name='Changed'`,
				"delete":    `DELETE FROM mappings`,
			}[change]
			_, err := db.Exec(query)
			require.NoError(t, err)
			called := false
			err = NewSQLRepository(db).AcceptFormatRecovery(t.Context(), req, evidence, func() error { called = true; return nil })
			require.ErrorIs(t, err, ErrFormatConflict)
			require.False(t, called, "changed target authority must not publish")
		})
	}
}
