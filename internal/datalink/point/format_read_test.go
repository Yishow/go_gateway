package point

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
)

func TestReadFormatIsZeroWriteCanonicalSnapshot(t *testing.T) {
	db, svc, req := formatFixture(t)
	_, err := db.Exec(`UPDATE system_settings SET value=json_set(value,'$.database_setup_revision','not-source-revision') WHERE key='studio_v2_workspace'`)
	require.NoError(t, err)
	// Reject every configuration/receipt write, not just changed values.
	for _, table := range []string{"points", "system_settings", "source_rules", "source_rule_links"} {
		for _, op := range []string{"INSERT", "UPDATE", "DELETE"} {
			_, err = db.Exec(`CREATE TRIGGER no_read_write_` + table + op + ` BEFORE ` + op + ` ON ` + table + ` BEGIN SELECT RAISE(ABORT,'read wrote'); END`)
			require.NoError(t, err)
		}
	}
	got, err := svc.ReadFormat(t.Context(), FormatReadScope{WorkspaceID: req.WorkspaceID, DeviceID: req.DeviceID, PointID: req.PointID})
	require.NoError(t, err)
	require.Equal(t, "ws-old", got.WorkspaceRevision)
	require.Equal(t, "persisted", got.WorkspaceRevisionOrigin)
	require.Equal(t, "point-old", got.PointRevision)
	require.False(t, got.Editable, "missing live authority must not promise submit readiness")
	require.True(t, got.Resolution.Resolved)
	require.Equal(t, got.WorkspaceRevision, got.Resolution.WorkspaceRevision)
	require.NotEmpty(t, got.ConnectionRevision)
	raw, err := json.Marshal(got)
	require.NoError(t, err)
	require.NotContains(t, string(raw), "connection_config")
	require.NotContains(t, string(raw), "last_value")
}

func TestReadFormatInitialRevisionDoesNotInitialize(t *testing.T) {
	db, svc, req := formatFixture(t)
	_, err := db.Exec(`DELETE FROM system_settings WHERE key='modbus_share_workspace_revision:ws'; UPDATE points SET revision_id=''`)
	require.NoError(t, err)
	got, err := svc.ReadFormat(t.Context(), FormatReadScope{WorkspaceID: req.WorkspaceID, DeviceID: req.DeviceID, PointID: req.PointID})
	require.NoError(t, err)
	require.Equal(t, "rev-1", got.WorkspaceRevision)
	require.Equal(t, "initial", got.WorkspaceRevisionOrigin)
	require.Equal(t, "", got.PointRevision)
	var n int
	require.NoError(t, db.QueryRow(`SELECT count(*) FROM system_settings WHERE key LIKE 'point_format_receipt:%' OR key='modbus_share_workspace_revision:ws'`).Scan(&n))
	require.Zero(t, n)
}

func TestReadFormatCapturedVersionsSaveAndReset(t *testing.T) {
	db, svc, req := formatFixture(t)
	scope := FormatReadScope{WorkspaceID: req.WorkspaceID, DeviceID: req.DeviceID, PointID: req.PointID}
	before, err := svc.ReadFormat(t.Context(), scope)
	require.NoError(t, err)
	req.ExpectedWorkspaceRevision = before.WorkspaceRevision
	req.ExpectedPointRevision = &before.PointRevision
	saved, err := svc.MutateFormat(t.Context(), req)
	require.NoError(t, err)
	require.Equal(t, "CDAB", saved.Resolution.EffectiveFormat)
	dirty, err := svc.ReadFormat(t.Context(), scope)
	require.NoError(t, err)
	require.False(t, dirty.Editable)
	// SQL-only fixture: known successful invalidation is no longer unknown dirty;
	// live submit readiness is independently checked by the production gate.
	_, err = db.Exec(`UPDATE system_settings SET value=json_set(value,'$.dirty',json('false')) WHERE key='modbus_share_workspace_revision:ws'`)
	require.NoError(t, err)
	next, err := svc.ReadFormat(t.Context(), scope)
	require.NoError(t, err)
	req.ExpectedWorkspaceRevision = next.WorkspaceRevision
	req.ExpectedPointRevision = &next.PointRevision
	req.RequestID = "reset-2"
	req.FormatMode = schema.PointFormatInherit
	req.ConfiguredFormat = ""
	reset, err := svc.MutateFormat(t.Context(), req)
	require.NoError(t, err)
	require.Equal(t, schema.PointFormatInherit, reset.Resolution.Mode)
	req.RequestID = "old-edit"
	req.ExpectedWorkspaceRevision = before.WorkspaceRevision
	req.ExpectedPointRevision = &before.PointRevision
	_, err = svc.MutateFormat(t.Context(), req)
	require.ErrorIs(t, err, ErrFormatConflict)
}

func TestReadFormatScopeAndRuleDiscovery(t *testing.T) {
	db, svc, req := formatFixture(t)
	_, err := db.Exec(`INSERT INTO source_rules(id,device_id,start_address,count,data_type,naming_prefix,enabled,revision_id,data_format) VALUES('rule-owned',?,'40001',1,'uint32','r',1,'rule-v1','BADC')`, req.DeviceID)
	require.NoError(t, err)
	p, err := svc.GetByID(t.Context(), req.PointID)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO source_rule_links(id,rule_id,point_id,address) VALUES('link-owned','rule-owned',?,?)`, req.PointID, p.Address)
	require.NoError(t, err)
	scope := FormatReadScope{WorkspaceID: req.WorkspaceID, DeviceID: req.DeviceID, PointID: req.PointID}
	got, err := svc.ReadFormat(t.Context(), scope)
	require.NoError(t, err)
	require.Equal(t, "rule-owned", got.RuleID)
	require.Equal(t, "rule-v1", got.RuleRevision)
	for _, bad := range []FormatReadScope{
		{WorkspaceID: "foreign", DeviceID: req.DeviceID, PointID: req.PointID},
		{WorkspaceID: req.WorkspaceID, DeviceID: "foreign", PointID: req.PointID},
		{WorkspaceID: req.WorkspaceID, DeviceID: req.DeviceID, PointID: "missing"},
		{WorkspaceID: req.WorkspaceID, DeviceID: req.DeviceID, PointID: req.PointID, RuleID: "foreign"},
	} {
		out, err := svc.ReadFormat(t.Context(), bad)
		require.ErrorIs(t, err, ErrFormatScope)
		require.Nil(t, out)
	}
}

type formatReadLiveGate struct {
	formatTestInvalidator
	check func(FormatScope) error
}

func (g formatReadLiveGate) CheckPointFormatScope(_ context.Context, s FormatScope) error {
	return g.check(s)
}
func TestReadFormatUsesOriginalSnapshotAtLiveGate(t *testing.T) {
	for _, stale := range []bool{false, true} {
		db, svc, req := formatFixture(t)
		svc.SetFormatInvalidator(formatReadLiveGate{check: func(scope FormatScope) error {
			require.Equal(t, "ws-old", scope.ExpectedWorkspaceRevision)
			require.Equal(t, "point-old", *scope.ExpectedPointRevision)
			if stale {
				return ErrFormatConflict
			}
			return nil
		}})
		// Both SQL and live gate remain read-only, including initial revision reads.
		_, err := db.Exec(`CREATE TRIGGER forbid_read_settings_update BEFORE UPDATE ON system_settings BEGIN SELECT RAISE(ABORT,'read may not initialize'); END`)
		require.NoError(t, err)
		got, err := svc.ReadFormat(t.Context(), FormatReadScope{WorkspaceID: req.WorkspaceID, DeviceID: req.DeviceID, PointID: req.PointID})
		require.NoError(t, err)
		require.Equal(t, !stale, got.Editable)
		require.Equal(t, "ws-old", got.WorkspaceRevision)
		require.True(t, got.Resolution.Resolved)
	}
}

func TestReadFormatConcurrentAuthorityChangeDoesNotRebaseSnapshot(t *testing.T) {
	db, svc, req := formatFixture(t)
	svc.SetFormatInvalidator(formatReadLiveGate{check: func(scope FormatScope) error {
		// Deterministic external-writer window AFTER the SQL read snapshot closed,
		// before the same live gate used by Save checks the captured revision.
		_, err := db.Exec(`UPDATE system_settings SET value='{"revision":"ws-new","dirty":false}' WHERE key='modbus_share_workspace_revision:ws'`)
		require.NoError(t, err)
		var current string
		require.NoError(t, db.QueryRow(`SELECT json_extract(value,'$.revision') FROM system_settings WHERE key='modbus_share_workspace_revision:ws'`).Scan(&current))
		if scope.ExpectedWorkspaceRevision != current {
			return ErrFormatConflict
		}
		return nil
	}})
	got, err := svc.ReadFormat(t.Context(), FormatReadScope{WorkspaceID: req.WorkspaceID, DeviceID: req.DeviceID, PointID: req.PointID})
	require.NoError(t, err)
	require.False(t, got.Editable)
	require.Equal(t, "ws-old", got.WorkspaceRevision)
	require.Equal(t, "ws-old", got.Resolution.WorkspaceRevision)
}
