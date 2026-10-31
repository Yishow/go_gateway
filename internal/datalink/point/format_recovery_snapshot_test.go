package point

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
	"testing"
)

func TestReadFormatRecoveryEvidence(t *testing.T) {
	for _, kind := range []string{"valid", "missing-clean-marker", "valid-inherit", "changed-inherited-default", "pending", "corrupt-resolution", "missing-settings", "missing-connection-proof", "stale-result", "corrupt", "unknown", "stale-point", "stale-workspace", "stale-rule", "stale-connection", "stale-settings", "foreign", "disabled-point", "valid-draft", "dirty", "dirty-share-disabled"} {
		t.Run(kind, func(t *testing.T) {
			db, svc, mutation := formatFixture(t)
			_, err := db.Exec(`INSERT INTO source_rules(id,device_id,start_address,count,data_type,naming_prefix,revision_id,data_format) VALUES('rule',?,'40001',1,'uint32','p','rule-old','ABCD')`, mutation.DeviceID)
			require.NoError(t, err)
			_, err = db.Exec(`INSERT INTO source_rule_links(id,rule_id,address,point_id) VALUES('link','rule','40001',?)`, mutation.PointID)
			require.NoError(t, err)
			mutation.RuleID = "rule"
			mutation.ExpectedRuleRevision = "rule-old"
			svc.SetFormatRuleGuard(&formatRuleTestGuard{})
			if kind == "valid-inherit" || kind == "changed-inherited-default" {
				mutation.FormatMode = schema.PointFormatInherit
				mutation.ConfiguredFormat = ""
			}
			result, err := svc.MutateFormat(t.Context(), mutation)
			require.NoError(t, err)
			require.Equal(t, "invalidated", result.InvalidationStatus)
			_, err = db.Exec(`UPDATE system_settings SET value='{"enabled":true,"settings_revision":"settings-1"}' WHERE key='modbus_share'`)
			require.NoError(t, err)
			req := FormatRecoveryRequest{FormatScope: mutation.FormatScope, ExpectedConnectionRevision: result.Resolution.ConnectionRevision, ExpectedSettingsRevision: "settings-1"}
			req.ExpectedPointRevision = &result.Point.RevisionID
			req.ExpectedWorkspaceRevision = result.WorkspaceRevision
			switch kind {
			case "missing-clean-marker":
				_, err = db.Exec(`UPDATE system_settings SET value=json_remove(value,'$.dirty') WHERE key='modbus_share_workspace_revision:ws'`)

			case "changed-inherited-default":
				_, err = db.Exec(`UPDATE source_rules SET data_format='DCBA'`)
			case "corrupt-resolution":
				_, err = db.Exec(`UPDATE system_settings SET value=json_set(value,'$.result.resolution.resolved',json('false')) WHERE key=?`, formatReceiptKey(req.WorkspaceID, req.RequestID))
			case "missing-settings":
				_, err = db.Exec(`DELETE FROM system_settings WHERE key='modbus_share'`)
			case "missing-connection-proof":
				req.ExpectedConnectionRevision = ""
			case "stale-result":
				req.ExpectedWorkspaceRevision = "ws-old"

			case "pending":
				_, err = db.Exec(`UPDATE system_settings SET value=json_set(value,'$.result.invalidation_status','pending') WHERE key=?`, formatReceiptKey(req.WorkspaceID, req.RequestID))
			case "corrupt":
				_, err = db.Exec(`UPDATE system_settings SET value='{"password":"secret"}' WHERE key=?`, formatReceiptKey(req.WorkspaceID, req.RequestID))
			case "unknown":
				req.RequestID = "unknown"
			case "stale-point":
				_, err = db.Exec(`UPDATE points SET revision_id='new'`)
			case "stale-workspace":
				_, err = db.Exec(`UPDATE system_settings SET value=json_set(value,'$.revision','new') WHERE key='modbus_share_workspace_revision:ws'`)
			case "stale-rule":
				_, err = db.Exec(`UPDATE source_rules SET revision_id='new',data_format='DCBA'`)
			case "stale-connection":
				_, err = db.Exec(`UPDATE devices SET connection_config='{"data_format":"DCBA","password":"secret"}'`)
			case "stale-settings":
				req.ExpectedSettingsRevision = "old"
			case "foreign":
				req.DeviceID = "foreign"
			case "disabled-point":
				_, err = db.Exec(`UPDATE points SET enabled=0`)
			case "valid-draft":
				_, err = db.Exec(`UPDATE devices SET status=?`, schema.DeviceStatusDraft)
			case "dirty", "dirty-share-disabled":
				_, err = db.Exec(`UPDATE system_settings SET value=json_set(value,'$.dirty',json('true')) WHERE key='modbus_share_workspace_revision:ws'`)
				if kind == "dirty-share-disabled" {
					_, err = db.Exec(`UPDATE system_settings SET value=json_set(value,'$.enabled',json('false')) WHERE key='modbus_share'`)
				}
			}
			require.NoError(t, err)
			var before int
			require.NoError(t, db.QueryRow(`SELECT total_changes()`).Scan(&before))
			evidence, err := NewSQLRepository(db).ReadFormatRecovery(t.Context(), req)
			var after int
			require.NoError(t, db.QueryRow(`SELECT total_changes()`).Scan(&after))
			require.Equal(t, before, after, "read must perform zero writes")
			if kind != "valid" && kind != "valid-inherit" && kind != "valid-draft" {
				require.Error(t, err)
				require.Nil(t, evidence)
				require.NotContains(t, err.Error(), "secret")
				return
			}
			require.NoError(t, err)
			require.Equal(t, result.Point.RevisionID, evidence.Point.RevisionID)
			require.Equal(t, result.Resolution, evidence.Resolution)
			require.Equal(t, req.RuleID, evidence.Profile.IdentityView().RuleID)
			require.Len(t, evidence.ReceiptFingerprint, 64)
			require.Contains(t, evidence.WorkspaceWitness, result.WorkspaceRevision)
			raw, err := json.Marshal(evidence)
			require.NoError(t, err)
			require.NotContains(t, string(raw), "WorkspaceWitness")
			require.NotContains(t, string(raw), "connection_config")
		})
	}
}
