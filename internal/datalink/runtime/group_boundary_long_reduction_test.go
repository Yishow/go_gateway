package runtime

import (
	"errors"
	"testing"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func blockedCode(t *testing.T, config GroupBoundaryConfig) string {
	t.Helper()
	_, err := NewGroupBoundary(config)
	if err == nil {
		return ""
	}
	var blocked *GroupBoundaryError
	require.True(t, errors.As(err, &blocked), "%v", err)
	return blocked.Code
}

// TestLongBoundaryRefusesAnUnknownSampleReduction pins that the production
// GroupBoundary entry itself fails closed on a sample reduction it does not
// know: it never falls back to the legacy latest rule, whether the members carry
// recording metadata or not, and for an instantaneous or a cumulative profile.
func TestLongBoundaryRefusesAnUnknownSampleReduction(t *testing.T) {
	future := workspace.WriteGroupSampleReduction("future")

	t.Run("legacy members without metadata", func(t *testing.T) {
		f := newLongFixture(t)
		f.config.Group.RecordingProfile.SampleReduction = future
		require.Equal(t, boundaryProfileUnsupported, blockedCode(t, f.config))
	})
	t.Run("members with complete mean metadata", func(t *testing.T) {
		f := counterRawFixture(t)
		f.config.Group.RecordingProfile.SampleReduction = future
		// Refused before the shape gate: the members' frozen result types contradict
		// a reduction that is not the mean they were frozen for.
		require.Equal(t, boundaryMemberMetadata, blockedCode(t, f.config))
	})
	t.Run("a confirmed cumulative profile", func(t *testing.T) {
		f := newCounterFixture(t)
		f.config.Group.RecordingProfile.SampleReduction = future
		require.Equal(t, boundaryProfileUnsupported, blockedCode(t, f.config))
	})

	t.Run("the legacy latest shapes still build", func(t *testing.T) {
		for name, reduction := range map[string]workspace.WriteGroupSampleReduction{
			"omitted": "", "pinned latest": workspace.WriteGroupSampleReductionLatest,
		} {
			f := newLongFixture(t)
			f.config.Group.RecordingProfile.SampleReduction = reduction
			require.Empty(t, blockedCode(t, f.config), name)
		}
		require.Empty(t, blockedCode(t, newCounterFixture(t).config), "a legacy cumulative profile")
		require.Empty(t, blockedCode(t, counterRawFixture(t).config), "a verified mean profile")
	})

	t.Run("a concrete tag error keeps its own code", func(t *testing.T) {
		f := counterRawFixture(t)
		delete(f.config.TagTypes, "tag-B")
		require.Equal(t, boundaryTagTypeUnsupport, blockedCode(t, f.config))
		f = counterRawFixture(t)
		f.config.TagTypes["tag-B"] = schema.DataType("blob")
		require.Equal(t, boundaryTagTypeUnsupport, blockedCode(t, f.config))
	})
}
