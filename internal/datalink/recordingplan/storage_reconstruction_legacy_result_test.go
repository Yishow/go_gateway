package recordingplan

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPlainReconstructionPreservesAppliedResultOnEventFailure(t *testing.T) {
	db := lineageDB(t)
	seedLineage(t, db, lineageScope)
	_, err := db.ExecContext(t.Context(), `DROP TABLE write_group_storage_events`)
	require.NoError(t, err)
	changed, err := RecordStorageReconstruction(t.Context(), db, lineageScope, 1, "legacy-repair", "group", time.Now().UTC())
	require.Error(t, err, "the diagnostic event failed")
	require.True(t, changed, "the pre-existing plain path reports its already-applied generation update")
	lineage, readErr := ReadStorageLineage(t.Context(), db, lineageScope)
	require.NoError(t, readErr)
	require.Equal(t, 2, lineage.Generation)
}
