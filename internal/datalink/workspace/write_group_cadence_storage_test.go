package workspace

import (
	"database/sql"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemberCadenceStorageLeavesLegacyNullAndJSONUnchanged(t *testing.T) {
	path := filepath.Join(t.TempDir(), "legacy-cadence.db")
	h := newCircuitProfileHarness(t, path)
	saved := h.mustEnsure(h.draft(WriteGroupProfileInstantaneous, "circuit-legacy", "Legacy", h.primary()))
	before, err := json.Marshal(saved)
	require.NoError(t, err)
	var cadence sql.NullString
	require.NoError(t, h.db.QueryRowContext(t.Context(), `SELECT cadence FROM write_group_members WHERE group_id=?`, saved.ID).Scan(&cadence))
	require.False(t, cadence.Valid)
	require.NoError(t, h.db.Close())
	db := openWorkspaceTestDB(t, path)
	restarted := newCircuitProfileHarnessOn(t, db, h)
	got, err := restarted.service.Get(t.Context(), saved.ID)
	require.NoError(t, err)
	after, err := json.Marshal(got.Group)
	require.NoError(t, err)
	require.Equal(t, string(before), string(after))
	require.Nil(t, got.Group.Members[0].Cadence)
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT cadence FROM write_group_members WHERE group_id=?`, saved.ID).Scan(&cadence))
	require.False(t, cadence.Valid)
}
