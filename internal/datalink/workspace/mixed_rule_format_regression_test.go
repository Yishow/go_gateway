package workspace

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMixedRuleItemFormatIsUsedForWriteGroupSourceRevision(t *testing.T) {
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	f := seedWriteGroupFixture(t.Context(), t, db)
	_, err := db.ExecContext(t.Context(), `UPDATE devices SET protocol='mc_3e',connection_config='{"data_format":"CDAB"}' WHERE id=?`, f.deviceID)
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), `UPDATE points SET address='D100',function='',data_type='uint32',data_format='ABCD',format_mode='inherit',format_version=1,revision_id='p1' WHERE id=?`, f.pointID)
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), `INSERT INTO source_rules(id,device_id,start_address,count,data_type,naming_prefix,revision_id,layout_mode,mixed_items) VALUES('r',?,'D100',1,'uint32','p','r1','mixed',?)`, f.deviceID, `[{"item_id":"value","address":"D100","data_type":"uint32","data_format":"ABCD"}]`)
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), `INSERT INTO source_rule_links(id,rule_id,address,point_id) VALUES('link','r','D100',?)`, f.pointID)
	require.NoError(t, err)
	repo := NewSQLWriteGroupRepository(db)
	tx, err := db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	sources, err := readMigrationPreviewSources(t.Context(), tx, repo, f.tagID)
	require.NoError(t, err)
	require.Len(t, sources, 1)
	require.Equal(t, "ABCD", sources[0].source.profile.Resolution().EffectiveFormat)
}
