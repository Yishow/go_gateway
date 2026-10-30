package point

import (
	"database/sql"
	"testing"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/schema/migrations"

	"github.com/stretchr/testify/require"
)

func addPointFormatTestColumns(t *testing.T, db *sql.DB) {
	t.Helper()
	raw, err := migrations.FS.ReadFile("048_point_format_intent_sqlite.up.sql")
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), string(raw))
	require.NoError(t, err)
}

func TestSQLRepositoryPointFormatRoundtrip(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	deviceID, groupID := setupTestData(t, db)
	repo := NewSQLRepository(db)
	for _, mode := range []schema.PointFormatMode{"", schema.PointFormatLegacy, schema.PointFormatInherit, schema.PointFormatOverride} {
		t.Run(string(mode), func(t *testing.T) {
			p := createTestPoint(deviceID)
			p.ID, p.Name, p.Address = "point-"+string(mode), "point-"+string(mode), "D"+string(mode)
			p.PollingGroupID = &groupID
			p.DataFormat = "ABCD"
			p.FormatMode = mode
			if mode != "" {
				p.FormatVersion, p.RevisionID = schema.PointFormatVersion, "revision-"+string(mode)
			}
			if mode == schema.PointFormatOverride {
				p.FormatOverride = "ABCD"
			}
			require.NoError(t, repo.Create(t.Context(), p))
			got, err := repo.GetByID(t.Context(), p.ID)
			require.NoError(t, err)
			assertPointFormat(t, p, got)
			for _, list := range []func() ([]*schema.Point, error){
				func() ([]*schema.Point, error) { return repo.ListByDevice(t.Context(), deviceID) },
				func() ([]*schema.Point, error) { return repo.ListByPollingGroup(t.Context(), groupID) },
				func() ([]*schema.Point, error) { return repo.List(t.Context(), ListFilter{}) },
			} {
				points, err := list()
				require.NoError(t, err)
				found := false
				for _, item := range points {
					if item.ID == p.ID {
						assertPointFormat(t, p, item)
						found = true
					}
				}
				require.True(t, found)
			}
			p.Name = "unrelated edit"
			require.NoError(t, repo.Update(t.Context(), p))
			got, err = NewSQLRepository(db).GetByID(t.Context(), p.ID)
			require.NoError(t, err)
			assertPointFormat(t, p, got)
			// Repository snapshots support restore and do not infer intent from strings.
			p.FormatMode, p.FormatOverride, p.FormatVersion, p.RevisionID = schema.PointFormatInherit, "", schema.PointFormatVersion, "reset-revision"
			require.NoError(t, repo.Update(t.Context(), p))
			got, err = repo.GetByID(t.Context(), p.ID)
			require.NoError(t, err)
			assertPointFormat(t, p, got)
		})
	}
}

func assertPointFormat(t *testing.T, want, got *schema.Point) {
	t.Helper()
	require.Equal(t, want.DataFormat, got.DataFormat)
	require.Equal(t, want.FormatMode, got.FormatMode)
	require.Equal(t, want.FormatOverride, got.FormatOverride)
	require.Equal(t, want.FormatVersion, got.FormatVersion)
	require.Equal(t, want.RevisionID, got.RevisionID)
}

func TestSQLRepositoryPointFormatRejectsInvalidMode(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	deviceID, _ := setupTestData(t, db)
	repo := NewSQLRepository(db)
	p := createTestPoint(deviceID)
	p.FormatMode = schema.PointFormatInherit
	require.NoError(t, repo.Create(t.Context(), p))
	p.FormatMode = "invalid-mode"
	require.Error(t, repo.Update(t.Context(), p))
	got, err := repo.GetByID(t.Context(), p.ID)
	require.NoError(t, err)
	require.Equal(t, schema.PointFormatInherit, got.FormatMode)
	p.ID, p.Address = "invalid-new", "D2"
	require.Error(t, repo.Create(t.Context(), p))
	_, err = repo.GetByID(t.Context(), p.ID)
	require.ErrorIs(t, err, ErrPointNotFound)
}
