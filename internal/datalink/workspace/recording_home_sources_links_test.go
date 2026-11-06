package workspace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRecordingHomeSourcesLinkEvidenceIsNotPresentationDigest(t *testing.T) {
	input := homeFixtures(t)[1].Before
	mappingID, tagID := "M1", "T1"
	input.Links = []RecordingHomeSourceLink{{ID: "L1", RuleID: "R1", PointID: "P1", Address: "40001", MappingID: &mappingID, TagID: &tagID}}
	observe := func() RecordingHomeSourcesObservation {
		got, err := ReadRecordingHomeSources(t.Context(), "W1", func(context.Context, string) (RecordingHomeSourceInput, error) { return input, nil })
		require.NoError(t, err)
		require.Equal(t, "available", got.State)
		return got
	}
	first := observe()
	input.Links[0].ID = "L2"
	second := observe()
	require.Equal(t, first.SourceDigest, second.SourceDigest)
	calls := 0
	got, err := ReadRecordingHomeSources(t.Context(), "W1", func(context.Context, string) (RecordingHomeSourceInput, error) {
		calls++
		input.Links[0].ID = "L1"
		if calls == 2 {
			input.Links[0].ID = "L2"
		}
		return input, nil
	})
	require.NoError(t, err)
	require.Equal(t, "unknown", got.State)
	require.Equal(t, "source-changed", got.ReasonCode)
	require.Nil(t, got.RecordingHomeSourceSummary)
}
