package point

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
	"os"
	"testing"
)

func formatFixture(t *testing.T) (*sql.DB, *Service, FormatMutationRequest) {
	t.Helper()
	db := setupTestDB(t)
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = db.Close() })
	for _, file := range []string{"005_source_rules_sqlite.up.sql", "009_source_rule_revision_sqlite.up.sql"} {
		raw, err := os.ReadFile("../schema/migrations/" + file)
		require.NoError(t, err)
		_, err = db.Exec(string(raw))
		require.NoError(t, err)
	}
	_, err := db.Exec(`ALTER TABLE source_rules ADD COLUMN data_format TEXT`)
	require.NoError(t, err)
	dev, _ := setupTestData(t, db)
	repo := NewSQLRepository(db)
	p := createTestPoint(dev)
	p.DataType = schema.DataTypeUint32
	p.RevisionID = "point-old"
	require.NoError(t, repo.Create(t.Context(), p))
	ws, _ := json.Marshal(map[string]any{"id": "ws", "ordered_device_ids": []string{dev}})
	_, err = db.Exec(`INSERT INTO system_settings(key,value) VALUES ('studio_v2_workspace',?),('modbus_share_workspace_revision:ws','{"revision":"ws-old","dirty":false}')`, string(ws))
	require.NoError(t, err)
	rev := p.RevisionID
	req := FormatMutationRequest{FormatScope: FormatScope{WorkspaceID: "ws", DeviceID: dev, PointID: p.ID, ExpectedWorkspaceRevision: "ws-old", ExpectedPointRevision: &rev, RequestID: "request-1"}, FormatMode: schema.PointFormatOverride, ConfiguredFormat: "CDAB"}
	svc := NewService(repo, nil)
	svc.SetFormatInvalidator(formatTestInvalidator{})
	return db, svc, req
}

func TestFormatMutationSQLCommitsAndReplaysExactly(t *testing.T) {
	db, svc, req := formatFixture(t)
	got, err := svc.MutateFormat(t.Context(), req)
	require.NoError(t, err)
	require.NotEqual(t, *req.ExpectedPointRevision, got.Point.RevisionID)
	require.Equal(t, "CDAB", got.Point.FormatOverride)
	require.Equal(t, "pending", got.RuntimeApplyStatus)
	restarted := NewService(NewSQLRepository(db), nil)
	restarted.SetFormatInvalidator(formatTestInvalidator{})
	replay, err := restarted.MutateFormat(t.Context(), req)
	require.NoError(t, err)
	require.Equal(t, got, replay)
	req.ConfiguredFormat = "ABCD"
	got, err = restarted.MutateFormat(t.Context(), req)
	require.ErrorIs(t, err, ErrFormatConflict)
	require.Nil(t, got)
}

func TestFormatMutationSQLScopeAndRevisionGuards(t *testing.T) {
	for _, tc := range []struct {
		name string
		edit func(*FormatMutationRequest)
		want error
	}{
		{"foreign-workspace", func(r *FormatMutationRequest) { r.WorkspaceID = "other" }, ErrFormatScope},
		{"foreign-device", func(r *FormatMutationRequest) { r.DeviceID = "other" }, ErrFormatScope},
		{"foreign-point", func(r *FormatMutationRequest) { r.PointID = "other" }, ErrFormatScope},
		{"stale-workspace", func(r *FormatMutationRequest) { r.ExpectedWorkspaceRevision = "stale" }, ErrFormatConflict},
		{"stale-point", func(r *FormatMutationRequest) { v := "stale"; r.ExpectedPointRevision = &v }, ErrFormatConflict},
		{"incompatible-type", func(r *FormatMutationRequest) {
			dtype := schema.DataTypeUint16
			name := "must-not-save"
			r.Patch = &UpdatePointRequest{DataType: &dtype, Name: &name}
		}, ErrFormatInvalid},
		{"unknown-format", func(r *FormatMutationRequest) { r.ConfiguredFormat = "secret://format" }, ErrFormatInvalid},
		{"missing-point-revision", func(r *FormatMutationRequest) { r.ExpectedPointRevision = nil }, ErrFormatInvalid},
		{"unsafe-request-key", func(r *FormatMutationRequest) { r.RequestID = "../../secret" }, ErrFormatInvalid},
	} {
		t.Run(tc.name, func(t *testing.T) {
			db, svc, req := formatFixture(t)
			original := req.PointID
			tc.edit(&req)
			got, err := svc.MutateFormat(t.Context(), req)
			require.ErrorIs(t, err, tc.want)
			require.Nil(t, got)
			require.NotContains(t, err.Error(), "secret")
			p, err := svc.GetByID(t.Context(), original)
			require.NoError(t, err)
			require.Equal(t, "point-old", p.RevisionID)
			var count int
			require.NoError(t, db.QueryRow(`SELECT count(*) FROM system_settings WHERE key LIKE 'point_format_receipt:%'`).Scan(&count))
			require.Zero(t, count)
		})
	}
}

func TestFormatMutationSQLRollbackFailureAndSanitizedReceipt(t *testing.T) {
	for _, target := range []string{"receipt", "workspace", "point"} {
		t.Run(target, func(t *testing.T) {
			db, svc, req := formatFixture(t)
			trigger := `CREATE TRIGGER fail_format BEFORE INSERT ON system_settings WHEN NEW.key LIKE 'point_format_receipt:%' BEGIN SELECT RAISE(ABORT,'secret://db-error'); END;`
			if target == "workspace" {
				trigger = `CREATE TRIGGER fail_format BEFORE UPDATE ON system_settings WHEN NEW.key LIKE 'modbus_share_workspace_revision:%' BEGIN SELECT RAISE(ABORT,'secret://db-error'); END;`
			}
			if target == "point" {
				trigger = `CREATE TRIGGER fail_format BEFORE UPDATE ON points BEGIN SELECT RAISE(ABORT,'secret://db-error'); END;`
			}
			_, err := db.Exec(trigger)
			require.NoError(t, err)
			got, err := svc.MutateFormat(t.Context(), req)
			require.ErrorIs(t, err, ErrFormatUnavailable)
			require.Nil(t, got)
			require.NotContains(t, err.Error(), "secret")
			var pointRevision, workspaceRevision string
			require.NoError(t, db.QueryRow(`SELECT revision_id FROM points`).Scan(&pointRevision))
			require.NoError(t, db.QueryRow(`SELECT json_extract(value,'$.revision') FROM system_settings WHERE key='modbus_share_workspace_revision:ws'`).Scan(&workspaceRevision))
			require.Equal(t, "point-old", pointRevision)
			require.Equal(t, "ws-old", workspaceRevision)
		})
	}
	db, svc, req := formatFixture(t)
	_, err := db.Exec(`UPDATE points SET last_error='secret://telemetry',last_value='private'`)
	require.NoError(t, err)
	got, err := svc.MutateFormat(t.Context(), req)
	require.NoError(t, err)
	require.Empty(t, got.Point.LastError)
	require.Nil(t, got.Point.LastValue)
	var receipt string
	require.NoError(t, db.QueryRow(`SELECT value FROM system_settings WHERE key LIKE 'point_format_receipt:%'`).Scan(&receipt))
	require.NotContains(t, receipt, "secret")
	require.NotContains(t, receipt, "private")
}

func TestFormatMutationSQLMixedLegacyPatchAndCAS(t *testing.T) {
	db, svc, req := formatFixture(t)
	before, err := svc.GetByID(t.Context(), req.PointID)
	require.NoError(t, err)
	name := "renamed"
	req.Patch = &UpdatePointRequest{Name: &name}
	got, err := svc.MutateFormat(t.Context(), req)
	require.NoError(t, err)
	require.Equal(t, name, got.Point.Name)
	before.Name = "stale overwrite"
	before.RevisionID = "other-new-revision"
	err = NewSQLRepository(db).UpdateRevision(t.Context(), before, *req.ExpectedPointRevision)
	require.ErrorIs(t, err, ErrFormatConflict)
	current, err := svc.GetByID(t.Context(), req.PointID)
	require.NoError(t, err)
	require.Equal(t, got.Point.RevisionID, current.RevisionID)
	require.Equal(t, "CDAB", current.FormatOverride)
	// A deliberate subsequent request after the existing hydration/reconcile clears dirty.
	_, err = db.Exec(`UPDATE system_settings SET value=json_set(value,'$.dirty',json('false')) WHERE key='modbus_share_workspace_revision:ws'`)
	require.NoError(t, err)
	req.RequestID = "reset"
	req.ExpectedWorkspaceRevision = got.WorkspaceRevision
	req.ExpectedPointRevision = &got.Point.RevisionID
	req.Patch = nil
	req.FormatMode = schema.PointFormatInherit
	req.ConfiguredFormat = ""
	reset, err := svc.MutateFormat(t.Context(), req)
	require.NoError(t, err)
	require.Equal(t, schema.PointFormatInherit, reset.Point.FormatMode)
	require.Empty(t, reset.Point.FormatOverride)
}

type formatTestInvalidator struct{}

func (formatTestInvalidator) InvalidatePointFormat(context.Context, FormatInvalidation) error {
	return nil
}
