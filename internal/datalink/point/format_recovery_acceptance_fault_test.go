package point

import (
	"database/sql"
	"errors"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/sqltx"
	"modernc.org/sqlite"
)

func recoveryAcceptanceFixture(t *testing.T) (*sql.DB, FormatRecoveryRequest, *FormatRecoveryEvidence) {
	t.Helper()
	db, svc, mutation := formatFixture(t)
	_, err := db.Exec(`INSERT INTO source_rules(id,device_id,start_address,count,data_type,naming_prefix,revision_id,data_format) VALUES('rule',?,'40001',1,'uint32','p','rule-old','ABCD')`, mutation.DeviceID)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO source_rule_links(id,rule_id,address,point_id) VALUES('link','rule','40001',?)`, mutation.PointID)
	require.NoError(t, err)
	mutation.RuleID = "rule"
	mutation.ExpectedRuleRevision = "rule-old"
	svc.SetFormatRuleGuard(&formatRuleTestGuard{})
	result, err := svc.MutateFormat(t.Context(), mutation)
	require.NoError(t, err)
	_, err = db.Exec(`UPDATE system_settings SET value='{"enabled":true,"settings_revision":"settings-1"}' WHERE key='modbus_share'`)
	require.NoError(t, err)
	req := FormatRecoveryRequest{FormatScope: mutation.FormatScope, ExpectedConnectionRevision: result.Resolution.ConnectionRevision, ExpectedSettingsRevision: "settings-1"}
	req.ExpectedPointRevision = &result.Point.RevisionID
	req.ExpectedWorkspaceRevision = result.WorkspaceRevision
	evidence, err := NewSQLRepository(db).ReadFormatRecovery(t.Context(), req)
	require.NoError(t, err)
	return db, req, evidence
}

// Capture authoritative values and timestamps. The zero-row reservation must
// never fire row triggers; publication callbacks never perform SQL.
func recoveryAcceptanceRows(t *testing.T, db *sql.DB) []string {
	t.Helper()
	var result []string
	for _, query := range []string{
		`SELECT id||'|'||revision_id||'|'||updated_at||'|'||format_mode||'|'||format_override FROM points ORDER BY id`,
		`SELECT key||'|'||value||'|'||COALESCE(updated_at,'') FROM system_settings ORDER BY key`,
		`SELECT id||'|'||revision_id||'|'||updated_at||'|'||data_format FROM source_rules ORDER BY id`,
		`SELECT id||'|'||status||'|'||connection_config||'|'||updated_at FROM devices ORDER BY id`,
	} {
		rows, err := db.Query(query)
		require.NoError(t, err)
		for rows.Next() {
			var raw string
			require.NoError(t, rows.Scan(&raw))
			result = append(result, raw)
		}
		require.NoError(t, rows.Err())
		require.NoError(t, rows.Close())
	}
	return result
}

func TestFormatRecoveryAcceptanceRejectsStaleEvidenceBeforeStage(t *testing.T) {
	for _, kind := range []string{"point", "workspace", "same-revision-dirty", "rule", "connection", "settings", "receipt", "point-config-same-revision", "device-inactive", "point-disabled"} {
		t.Run(kind, func(t *testing.T) {
			db, req, evidence := recoveryAcceptanceFixture(t)
			var err error
			switch kind {
			case "point":
				_, err = db.Exec(`UPDATE points SET revision_id='new'`)
			case "workspace":
				_, err = db.Exec(`UPDATE system_settings SET value=json_set(value,'$.revision','new') WHERE key='modbus_share_workspace_revision:ws'`)
			case "same-revision-dirty":
				_, err = db.Exec(`UPDATE system_settings SET value=json_set(value,'$.dirty',json('true')) WHERE key='modbus_share_workspace_revision:ws'`)
			case "rule":
				_, err = db.Exec(`UPDATE source_rules SET revision_id='new'`)
			case "connection":
				_, err = db.Exec(`UPDATE devices SET connection_config='{"data_format":"DCBA"}'`)
			case "settings":
				_, err = db.Exec(`UPDATE system_settings SET value=json_set(value,'$.settings_revision','new') WHERE key='modbus_share'`)
			case "receipt":
				_, err = db.Exec(`UPDATE system_settings SET value=json_set(value,'$.result.invalidation_status','pending') WHERE key=?`, formatReceiptKey(req.WorkspaceID, req.RequestID))
			case "point-config-same-revision":
				_, err = db.Exec(`UPDATE points SET name='new'`)
			case "device-inactive":
				_, err = db.Exec(`UPDATE devices SET status='draft'`)
			case "point-disabled":
				_, err = db.Exec(`UPDATE points SET enabled=0`)
			}
			require.NoError(t, err)
			before := recoveryAcceptanceRows(t, db)
			stages := 0
			err = NewSQLRepository(db).AcceptFormatRecovery(t.Context(), req, evidence, func() error { stages++; return nil })
			require.Error(t, err)
			require.Zero(t, stages)
			require.Equal(t, before, recoveryAcceptanceRows(t, db))
		})
	}
}

func TestFormatRecoveryAcceptancePreservesAuthorityRows(t *testing.T) {
	for _, failStage := range []bool{false, true} {
		t.Run(map[bool]string{false: "success", true: "stage-failure"}[failStage], func(t *testing.T) {
			db, req, evidence := recoveryAcceptanceFixture(t)
			before := recoveryAcceptanceRows(t, db)
			stageFailure := errors.New("stage rejected")
			stages := 0
			err := NewSQLRepository(db).AcceptFormatRecovery(t.Context(), req, evidence, func() error {
				stages++
				if failStage {
					return stageFailure
				}
				return nil
			})
			if failStage {
				require.ErrorIs(t, err, stageFailure)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, 1, stages)
			require.Equal(t, before, recoveryAcceptanceRows(t, db))
		})
	}
}

func TestFormatRecoveryAcceptanceCommitBusyCleansPool(t *testing.T) {
	memory, req, evidence := recoveryAcceptanceFixture(t)
	path := filepath.Join(t.TempDir(), "recovery-busy.db")
	_, err := memory.Exec(`VACUUM INTO ?`, path)
	require.NoError(t, err)
	db, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	defer db.Close()
	db.SetMaxOpenConns(1)
	_, err = db.Exec(`PRAGMA journal_mode=DELETE; PRAGMA busy_timeout=0`)
	require.NoError(t, err)
	reader, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	defer reader.Close()
	held, err := reader.Begin()
	require.NoError(t, err)
	defer held.Rollback()
	var count int
	require.NoError(t, held.QueryRow(`SELECT count(*) FROM points`).Scan(&count))
	before := recoveryAcceptanceRows(t, db)
	stages := 0
	err = NewSQLRepository(db).AcceptFormatRecovery(t.Context(), req, evidence, func() error { stages++; return nil })
	require.Error(t, err)
	require.ErrorIs(t, err, ErrFormatUnavailable)
	var sqliteErr *sqlite.Error
	require.ErrorAs(t, err, &sqliteErr)
	require.Equal(t, 5, sqliteErr.Code()&0xff, "must exercise native SQLITE_BUSY COMMIT")
	require.Equal(t, 1, stages, "failure must occur after the acceptance callback, at zero-write cleanup COMMIT")
	require.Equal(t, before, recoveryAcceptanceRows(t, db))
	require.NoError(t, held.Rollback())
	clean, err := db.Begin()
	require.NoError(t, err, "busy COMMIT must not leave a native transaction in the pool")
	require.NoError(t, clean.Rollback())
}

func TestFormatRecoveryAcceptanceLostCommitAckIsUnknown(t *testing.T) {
	memory, req, evidence := recoveryAcceptanceFixture(t)
	path := filepath.Join(t.TempDir(), "recovery-unknown.db")
	_, err := memory.Exec(`VACUUM INTO ?`, path)
	require.NoError(t, err)
	connector := &pointLostAckConnector{path: path, fail: true}
	db := sql.OpenDB(connector)
	defer db.Close()
	db.SetMaxOpenConns(1)
	before := recoveryAcceptanceRows(t, db)
	stages := 0
	err = NewSQLRepository(db).AcceptFormatRecovery(t.Context(), req, evidence, func() error { stages++; return nil })
	require.ErrorIs(t, err, sqltx.ErrCommitOutcomeUnknown)
	require.ErrorIs(t, err, ErrFormatUnavailable)
	require.NotContains(t, err.Error(), "private lost")
	require.Equal(t, 1, stages)
	require.Equal(t, before, recoveryAcceptanceRows(t, db))
	clean, err := db.Begin()
	require.NoError(t, err)
	require.NoError(t, clean.Rollback())
}
