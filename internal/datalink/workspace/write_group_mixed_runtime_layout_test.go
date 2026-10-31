package workspace

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMixedManagedSchemaLayoutDelegatesCanonicalBuilder(t *testing.T) {
	for _, version := range []WriteGroupLayoutVersion{WriteGroupLayoutManagedWideProfileV1, WriteGroupLayoutManagedLongV1} {
		h := newCircuitProfileHarness(t, ":memory:")
		g := mixedDraft(h, WriteGroupMemberCalculationMean, 15, "mixed_export")
		g.LayoutVersion = version
		if version == WriteGroupLayoutManagedLongV1 {
			g.StorageMode = WriteGroupStorageModePlain
			g.Destination.TableName = ""
		}
		saved, err := h.service.Create(t.Context(), h.mutation(t.Context(), g))
		require.NoError(t, err)
		ready, err := h.service.readinessSnapshot(t.Context(), saved.Group.ID, nil)
		require.NoError(t, err)
		for _, dialect := range []string{"sqlite", "postgres"} {
			expected, err := managedSchemaLayout(saved.Group, ready.tagTypes, dialect)
			require.NoError(t, err)
			got, err := MixedManagedSchemaLayout(saved.Group, ready.tagTypes, dialect)
			require.NoError(t, err)
			require.Equal(t, expected, got)
		}
		for _, dialect := range []string{"mysql", "sqlserver"} {
			_, err := MixedManagedSchemaLayout(saved.Group, ready.tagTypes, dialect)
			require.Error(t, err)
		}
		saved.Group.RecordingProfile.Kind = WriteGroupProfileInstantaneous
		_, err = MixedManagedSchemaLayout(saved.Group, ready.tagTypes, "sqlite")
		require.Error(t, err)
	}
}
