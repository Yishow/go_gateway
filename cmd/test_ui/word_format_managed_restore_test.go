package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
)

func wordSQLiteBackupCopy(t *testing.T, source, destination string) {
	t.Helper()
	script := `import pathlib,sqlite3,sys
source=pathlib.Path(sys.argv[1]).resolve(strict=True)
with sqlite3.connect(source.as_uri()+"?mode=ro",uri=True) as src, sqlite3.connect(sys.argv[2]) as dst:
 src.backup(dst)
 assert dst.execute("PRAGMA integrity_check").fetchone()[0]=="ok"
`
	output, err := exec.CommandContext(t.Context(), "python3", "-c", script, source, destination).CombinedOutput()
	require.NoError(t, err, string(output))
}

func TestProductionWordFormatManagedBackupRestorePreservesIntentAndHistory(t *testing.T) {
	f, _ := newRecoveryFaultProductionFixture(t)
	req := recoveryFaultSave(t, f) // real HTTP explicit override + immutable receipt
	original, err := f.services.point.GetByID(t.Context(), req.PointID)
	require.NoError(t, err)
	require.Equal(t, schema.PointFormatOverride, original.FormatMode)
	require.Equal(t, "CDAB", original.FormatOverride)
	f.plc.set(0, []uint16{0x3344, 0x1122})
	historyPath := filepath.Join(t.TempDir(), "historical-output.db")
	history, err := sql.Open("sqlite", historyPath)
	require.NoError(t, err)
	_, err = history.Exec(`CREATE TABLE history(id INTEGER,value TEXT,source_revision TEXT); INSERT INTO history VALUES(1,'18446744073709551615','frozen-before-format-edit')`)
	require.NoError(t, err)
	require.NoError(t, history.Close())
	historicalBytes, err := os.ReadFile(historyPath)
	require.NoError(t, err)
	historicalHash := sha256.Sum256(historicalBytes)
	var seq int
	var name, source string
	require.NoError(t, f.db.QueryRow(`PRAGMA database_list`).Scan(&seq, &name, &source))
	backup := filepath.Join(t.TempDir(), "config-backup.db")
	wordSQLiteBackupCopy(t, source, backup)
	restoredPath := filepath.Join(t.TempDir(), "config-restored-copy.db")
	wordSQLiteBackupCopy(t, backup, restoredPath)
	// A backup containing modern intent must still refuse the legacy target.
	output, err := exec.CommandContext(t.Context(), "python3", "../../scripts/check_point_format_downgrade.py", "--sqlite-file", restoredPath, "--target-format-version", "0").CombinedOutput()
	var exit *exec.ExitError
	require.ErrorAs(t, err, &exit, string(output))
	require.Equal(t, 2, exit.ExitCode(), string(output))
	restoredDB, err := sql.Open("sqlite", restoredPath)
	require.NoError(t, err)
	restoredDB.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, restoredDB.Close()) })
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	t.Cleanup(func() { require.NoError(t, manager.CloseAll()) })
	restored, err := wireGatewayServices(restoredDB, manager)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, restored.runtime.Stop(context.Background())) })
	p, err := restored.point.GetByID(t.Context(), original.ID)
	require.NoError(t, err)
	require.Equal(t, original.FormatMode, p.FormatMode)
	require.Equal(t, original.FormatOverride, p.FormatOverride)
	require.Equal(t, original.FormatVersion, p.FormatVersion)
	require.Equal(t, original.RevisionID, p.RevisionID)
	require.NoError(t, restored.runtime.Start(t.Context()))
	values := restored.scheduler.PollNowContext(t.Context(), []string{p.ID})
	require.Len(t, values, 1)
	require.Equal(t, schema.QualityGood, values[0].Quality)
	require.Equal(t, uint32(0x11223344), values[0].Value)
	profile, ok := restored.scheduler.CapturePointReadProfile(p.ID, values[0].InstallationLease)
	require.True(t, ok)
	require.Equal(t, req.RuleID, profile.IdentityView().RuleID)
	require.Equal(t, "CDAB", profile.Resolution().EffectiveFormat)
	current, err := f.services.point.GetByID(t.Context(), p.ID)
	require.NoError(t, err)
	require.Equal(t, original.RevisionID, current.RevisionID, "restore never overwrites original settings")
	require.Equal(t, original.FormatOverride, current.FormatOverride)
	after, err := os.ReadFile(historyPath)
	require.NoError(t, err)
	require.Equal(t, historicalHash, sha256.Sum256(after), "historical SQL remains byte-for-byte frozen")
}
