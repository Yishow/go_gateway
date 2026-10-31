package workspace

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestModernSourceRevisionMatchesMigrationSQLReadback(t *testing.T) {
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	f := seedWriteGroupFixture(t.Context(), t, db)
	_, err := db.ExecContext(t.Context(), `UPDATE devices SET protocol='mc_3e',connection_config='{"data_format":"CDAB"}' WHERE id=?`, f.deviceID)
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), `UPDATE points SET address='D0',function='',data_type='uint32',data_format='ABCD',format_mode='override',format_override='ABCD',format_version=1,revision_id='p1' WHERE id=?`, f.pointID)
	require.NoError(t, err)
	repo := NewSQLWriteGroupRepository(db)
	g := newBasicWriteGroup(f)
	require.NoError(t, repo.Create(t.Context(), g))
	tx, err := db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	sources, err := readMigrationPreviewSources(t.Context(), tx, repo, f.tagID)
	require.NoError(t, err)
	require.Len(t, sources, 1)
	revision, err := sources[0].source.revision()
	require.NoError(t, err)
	require.Equal(t, g.Members[0].SourceRevision, revision)
	require.NoError(t, tx.Rollback())
	// The previously identical legacy DataFormat does not authorize reusing an
	// old source binding after switching to explicit per-point policy.
	_, err = db.ExecContext(t.Context(), `UPDATE points SET format_mode='legacy',format_override='',format_version=0 WHERE id=?`, f.pointID)
	require.NoError(t, err)
	legacy := newBasicWriteGroup(f)
	require.NoError(t, repo.Create(t.Context(), legacy))
	require.NotEqual(t, g.Members[0].SourceRevision, legacy.Members[0].SourceRevision)
}
