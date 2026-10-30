package workspace

import (
	"database/sql"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func TestCalculationStorageLegacyAbsenceSurvivesRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "configuration.db")
	h := newCircuitProfileHarness(t, path)
	draft := h.draft(WriteGroupProfileInstantaneous, "legacy", "Legacy", h.primary())
	draft.RecordingProfile.SampleReduction = ""
	saved := h.mustEnsure(draft)
	var calculation, shape sql.NullString
	require.NoError(t, h.db.QueryRow(`SELECT calculation FROM write_group_members WHERE group_id=?`, saved.ID).Scan(&calculation))
	require.NoError(t, h.db.QueryRow(`SELECT result_shape FROM write_groups WHERE id=?`, saved.ID).Scan(&shape))
	require.False(t, calculation.Valid)
	require.False(t, shape.Valid)
	require.NoError(t, h.db.Close())
	db := openWorkspaceTestDB(t, path)
	restarted := newCircuitProfileHarnessOn(t, db, h)
	got, err := restarted.service.Get(t.Context(), saved.ID)
	require.NoError(t, err)
	require.Equal(t, saved, got.Group)
	require.Empty(t, got.Group.Members[0].Calculation)
	require.Empty(t, got.Group.RecordingProfile.ResultShape)
}

// The low-level writer is not a public save/activation API. This pins expansion
// storage separately from the fail-closed mixed configuration capability.
func TestCalculationStorageExplicitSnapshotRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "snapshot.db")
	h := newCircuitProfileHarness(t, path)
	legacy := h.mustEnsure(h.draft(WriteGroupProfileInstantaneous, "legacy", "Legacy", h.primary()))
	candidate := cloneWriteGroup(legacy)
	candidate.ID = "internal-mixed-snapshot"
	candidate.RecordingProfile = &WriteGroupRecordingProfile{Kind: WriteGroupProfileMixed, ResultShape: WriteGroupResultShapeMemberCalculationV1}
	candidate.Members[0].Calculation = WriteGroupMemberCalculationMean
	repo := NewSQLWriteGroupRepository(h.db)
	tx, err := h.db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	require.NoError(t, repo.insertGroup(t.Context(), tx, candidate))
	require.NoError(t, tx.Commit())
	got, err := repo.Get(t.Context(), h.fixture.workspaceID, candidate.ID)
	require.NoError(t, err)
	require.Equal(t, candidate.RecordingProfile, got.RecordingProfile)
	require.Equal(t, candidate.Members, got.Members)
	candidate.Members[0].Calculation = WriteGroupMemberCalculationLatest
	tx, err = h.db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	require.NoError(t, repo.updateGroup(t.Context(), tx, candidate))
	require.NoError(t, tx.Commit())
	require.NoError(t, h.db.Close())
	db := openWorkspaceTestDB(t, path)
	t.Cleanup(func() { _ = db.Close() })
	repo = NewSQLWriteGroupRepository(db)
	got, err = repo.Get(t.Context(), h.fixture.workspaceID, candidate.ID)
	require.NoError(t, err)
	require.Equal(t, candidate.RecordingProfile, got.RecordingProfile)
	require.Equal(t, candidate.Members, got.Members)
	require.True(t, memberMetadataUnsupported(&got.Members[0]))
	require.ErrorIs(t, normalizeWriteGroupProfile(got), ErrWriteGroupValidation, "persisted intent cannot activate normal mixed save")

	_, err = repo.Get(t.Context(), "foreign-workspace", candidate.ID)
	require.ErrorIs(t, err, ErrWriteGroupNotFound)
	// Legacy reads cannot acquire a method from an unrelated stored record.
	untouched, err := repo.Get(t.Context(), h.fixture.workspaceID, legacy.ID)
	require.NoError(t, err)
	require.Empty(t, untouched.Members[0].Calculation)
}

func TestCalculationStorageCorruptCrossShapeRejected(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	saved := h.mustEnsure(h.draft(WriteGroupProfileInstantaneous, "legacy", "Legacy", h.primary()))
	repo := NewSQLWriteGroupRepository(h.db)
	_, err := h.db.Exec(`UPDATE write_group_members SET calculation='mean' WHERE group_id=?`, saved.ID)
	require.NoError(t, err)
	_, err = repo.Get(t.Context(), h.fixture.workspaceID, saved.ID)
	require.ErrorIs(t, err, ErrWriteGroupValidation)
	_, err = h.db.Exec(`UPDATE write_group_members SET calculation=NULL WHERE group_id=?`, saved.ID)
	require.NoError(t, err)
	_, err = h.db.Exec(`UPDATE write_groups SET result_shape='member-calculation-v1' WHERE id=?`, saved.ID)
	require.NoError(t, err)
	_, err = repo.Get(t.Context(), h.fixture.workspaceID, saved.ID)
	require.ErrorIs(t, err, ErrWriteGroupValidation)
}

func TestCalculationStorageRejectsInvalidSnapshotBeforeMutation(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	saved := h.mustEnsure(h.draft(WriteGroupProfileInstantaneous, "legacy", "Legacy", h.primary()))
	repo := NewSQLWriteGroupRepository(h.db)
	for _, mode := range []string{"missing method", "conflicting interval"} {
		t.Run(mode, func(t *testing.T) {
			candidate := cloneWriteGroup(saved)
			candidate.RecordingProfile = &WriteGroupRecordingProfile{Kind: WriteGroupProfileMixed, ResultShape: WriteGroupResultShapeMemberCalculationV1}
			if mode == "conflicting interval" {
				candidate.Members[0].Calculation = WriteGroupMemberCalculationMean
				candidate.RecordingProfile.Interval = 15
			}
			tx, err := h.db.BeginTx(t.Context(), nil)
			require.NoError(t, err)
			require.ErrorIs(t, repo.updateGroup(t.Context(), tx, candidate), ErrWriteGroupValidation)
			require.NoError(t, tx.Rollback())
			got, err := repo.Get(t.Context(), h.fixture.workspaceID, saved.ID)
			require.NoError(t, err)
			require.Equal(t, saved, got)
			candidate.ID = "must-not-save"
			tx, err = h.db.BeginTx(t.Context(), nil)
			require.NoError(t, err)
			require.ErrorIs(t, repo.insertGroup(t.Context(), tx, candidate), ErrWriteGroupValidation)
			require.NoError(t, tx.Rollback())
			_, err = repo.Get(t.Context(), h.fixture.workspaceID, candidate.ID)
			require.ErrorIs(t, err, ErrWriteGroupNotFound)
		})
	}
}
