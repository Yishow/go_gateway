package point

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFormatRecoveryAcceptanceLocksBeforeStaging(t *testing.T) {
	db, svc, mutation := formatFixture(t)
	saved, err := svc.MutateFormat(t.Context(), mutation)
	require.NoError(t, err)
	_, err = db.Exec(`UPDATE system_settings SET value='{"enabled":false,"settings_revision":"settings-1"}' WHERE key='modbus_share'`)
	require.NoError(t, err)
	req := FormatRecoveryRequest{FormatScope: mutation.FormatScope, ExpectedConnectionRevision: saved.Resolution.ConnectionRevision, ExpectedSettingsRevision: "settings-1"}
	req.ExpectedPointRevision = &saved.Point.RevisionID
	req.ExpectedWorkspaceRevision = saved.WorkspaceRevision
	evidence, err := svc.ReadFormatRecovery(t.Context(), req)
	require.NoError(t, err)
	_, err = db.Exec(`DROP TABLE system_settings`)
	require.NoError(t, err)
	staged := false
	err = svc.AcceptFormatRecovery(t.Context(), req, evidence, func() error { staged = true; return nil })
	require.Error(t, err, "acceptance must establish a trustworthy source snapshot before any publication")
	require.False(t, staged)
}
