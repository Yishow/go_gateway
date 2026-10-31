package runtime

import (
	"database/sql"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
	_ "modernc.org/sqlite"
	"path/filepath"
	"testing"
	"time"
)

func TestFormatRecoveryPendingCommitCannotPoll(t *testing.T) {
	f := newAcquisitionInstallFixture(t, false)
	p := f.projection.Points[0]
	p.FormatMode = schema.PointFormatInherit
	p.FormatVersion = 1
	p.RevisionID = "p"
	p.DataType = schema.DataTypeUint32
	f.projection.RuleLinks[0].Address = p.Address
	require.NoError(t, f.service.ApplyWorkspaceProjection(t.Context(), f.projection))
	require.True(t, f.service.SuspendPointRevision(p.ID, p.RevisionID))
	proof, ok := f.service.CaptureOrdinaryPoint(p.ID)
	require.True(t, ok)
	profile, err := f.service.ResolvePointReadProfile(t.Context(), p)
	require.NoError(t, err)
	path := filepath.Join(t.TempDir(), "accept.db")
	db, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	defer db.Close()
	db.SetMaxOpenConns(1)
	_, err = db.Exec(`PRAGMA journal_mode=DELETE;PRAGMA busy_timeout=2000;CREATE TABLE acceptance(id INTEGER PRIMARY KEY,value TEXT);INSERT INTO acceptance VALUES(1,'before')`)
	require.NoError(t, err)
	reader, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	defer reader.Close()
	held, err := reader.Begin()
	require.NoError(t, err)
	var old string
	require.NoError(t, held.QueryRow(`SELECT value FROM acceptance`).Scan(&old))
	tx, err := db.Begin()
	require.NoError(t, err)
	_, err = tx.Exec(`UPDATE acceptance SET value='accepted' WHERE id=1`)
	require.NoError(t, err)
	lease, ok := f.service.InstallPreparedFormatRecovery(p, proof, profile)
	require.True(t, ok)
	require.NotZero(t, lease)
	committed := make(chan error, 1)
	go func() { committed <- tx.Commit() }()
	defer func() { _ = held.Rollback(); <-committed }()
	select {
	case err := <-committed:
		committed <- err
		t.Fatalf("commit did not block: %v", err)
	case <-time.After(20 * time.Millisecond):
	}
	require.Empty(t, f.service.scheduler.PollNowContext(t.Context(), []string{p.ID}), "a staged recovery must not produce manual samples while durable commit is blocked")
	direct, scope, captured := f.service.CaptureDirectReadProof(p.ID)
	require.False(t, captured && f.service.AcceptsDirectReadProof(direct, scope), "direct fallback cannot bypass pending acceptance")
	require.False(t, f.service.scheduler.AcceptsInstallationSample(p.ID, lease, f.service.scheduler.InstallationScope()))
}
