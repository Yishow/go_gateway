package point

import (
	"testing"

	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestServicePointFormat_NewDefaultsInheritWithoutGuessingProjection(t *testing.T) {
	repo := NewMemoryRepository()
	service := NewService(repo, NewMemoryPollingGroupRepository())
	point, err := service.Create(t.Context(), CreatePointRequest{
		DeviceID: "device", Name: "new", Address: "HR100", DataType: schema.DataTypeUint32, DataFormat: "ABCD",
	})
	require.NoError(t, err)
	require.Equal(t, schema.PointFormatInherit, point.FormatMode)
	require.Empty(t, point.FormatOverride)
	require.Equal(t, schema.PointFormatVersion, point.FormatVersion)
	require.NotEmpty(t, point.RevisionID)
	require.Equal(t, "ABCD", point.DataFormat, "retain old projection, do not interpret it as intent")
	stored, err := repo.GetByID(t.Context(), point.ID)
	require.NoError(t, err)
	require.Equal(t, point, stored)
}

func TestServicePointFormat_UnrelatedUpdatePreservesLegacyAndOverride(t *testing.T) {
	for _, mode := range []schema.PointFormatMode{"", schema.PointFormatLegacy, schema.PointFormatOverride} {
		repo := NewMemoryRepository()
		service := NewService(repo, NewMemoryPollingGroupRepository())
		original := schema.Point{ID: "point", DeviceID: "device", Name: "before", Address: "D0", DataType: schema.DataTypeUint32,
			DataFormat: "ABCD", FormatMode: mode, RevisionID: "original"}
		if mode == schema.PointFormatOverride {
			original.FormatOverride = "ABCD"
			original.FormatVersion = 1
		}
		require.NoError(t, repo.Create(t.Context(), &original))
		defaults := decoding.FormatContext{Protocol: schema.ProtocolMC3E, ConnectionFormat: "CDAB"}
		before, err := decoding.ResolvePoint(original, defaults)
		require.NoError(t, err)
		name := "after"
		updated, err := service.Update(t.Context(), original.ID, UpdatePointRequest{Name: &name})
		require.NoError(t, err)
		require.Equal(t, original.FormatMode, updated.FormatMode)
		require.Equal(t, original.FormatOverride, updated.FormatOverride)
		require.Equal(t, original.FormatVersion, updated.FormatVersion)
		require.NotEqual(t, original.RevisionID, updated.RevisionID)
		after, err := decoding.ResolvePoint(*updated, defaults)
		require.NoError(t, err)
		require.Equal(t, before.EffectiveFormat, after.EffectiveFormat)
		require.Equal(t, before.Policy, after.Policy)
		require.NoError(t, service.Restore(t.Context(), &original))
		restored, err := service.GetByID(t.Context(), original.ID)
		require.NoError(t, err)
		require.Equal(t, original.FormatMode, restored.FormatMode)
		require.Equal(t, original.FormatOverride, restored.FormatOverride)
		require.Equal(t, original.RevisionID, restored.RevisionID)
	}
}
