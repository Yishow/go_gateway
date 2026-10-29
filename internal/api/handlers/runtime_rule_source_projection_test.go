package handlers

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
)

func TestRuntimeSourceEvidenceByteOrderScaleCastAndUnsupportedParams(t *testing.T) {
	for _, order := range []string{"ABCD", "BADC", "CDAB", "DCBA"} {
		t.Run(order, func(t *testing.T) {
			f := newRuntimeSourceFixture(t)
			f.point.DataFormat = order
			require.NoError(t, f.points.Update(t.Context(), f.point))
			s := f.source(t)
			require.Equal(t, order, *s.Saved.Transforms[0].Order)
		})
	}
	for name, pipeline := range map[string]string{
		"formula":           `[{"type":"formula","order":0,"params":{"expression":"password=fixture-secret"}}]`,
		"unknown":           `[{"type":"future","order":0,"params":{"endpoint":"postgres://fixture-secret@private/db"}}]`,
		"extra-scale-param": `[{"type":"scale","order":0,"params":{"scale":0.1,"offset":0,"secret":"fixture-secret"}}]`,
		"invalid-scale":     `[{"type":"scale","order":0,"params":{"scale":"fixture-secret","offset":0}}]`,
		"missing-offset":    `[{"type":"scale","order":0,"params":{"scale":0.1}}]`,
		"invalid-cast":      `[{"type":"cast","order":0,"params":{"target_type":"fixture-secret"}}]`,
		"nested":            `[{"type":"cast","order":0,"params":{"target_type":{"secret":"fixture-secret"}}}]`,
		"malformed":         `{"fixture-secret":true}`,
	} {
		t.Run(name, func(t *testing.T) {
			f := newRuntimeSourceFixture(t)
			f.mapping.TransformPipeline = pipeline
			s := f.source(t)
			require.Equal(t, "unsupported", s.Saved.Transforms[len(s.Saved.Transforms)-1].Kind)
			raw, err := json.Marshal(s)
			require.NoError(t, err)
			require.NotContains(t, string(raw), "fixture-secret")
			require.Nil(t, s.Effective)
		})
	}
}
func TestRuntimeSourceEvidenceMissingDecodeAndPollingStayUnavailable(t *testing.T) {
	f := newRuntimeSourceFixture(t)
	f.point.DataFormat = ""
	f.point.PollingGroupID = nil
	require.NoError(t, f.points.Update(t.Context(), f.point))
	s := f.source(t)
	require.Nil(t, s.Saved.Polling)
	for _, step := range s.Saved.Transforms {
		require.NotEqual(t, "byte-order", step.Kind)
	}
	f = newRuntimeSourceFixture(t)
	f.point.DataFormat = "postgres://fixture-secret@private/db"
	f.point.Address = "/private/fixture-secret"
	f.link.Address = f.point.Address
	f.point.Function = "password=fixture-secret"
	f.point.DataType = "future-type"
	require.NoError(t, f.points.Update(t.Context(), f.point))
	s = f.source(t)
	raw, err := json.Marshal(s)
	require.NoError(t, err)
	require.NotContains(t, string(raw), "fixture-secret")
	require.Nil(t, s.Saved.ReadType)
	require.Nil(t, s.Saved.Address)
	require.Nil(t, s.Saved.ProtocolFunction)
}
func TestRuntimeSourceEvidenceUnvalidatedMappingDoesNotInventTransforms(t *testing.T) {
	for _, status := range []schema.MappingStatus{schema.MappingStatusDraft, schema.MappingStatusOutOfSync, schema.MappingStatusError, "future"} {
		t.Run(string(status), func(t *testing.T) {
			f := newRuntimeSourceFixture(t)
			f.mapping.Status = status
			s := f.source(t)
			require.Nil(t, s.Saved.ResultType)
			require.Equal(t, "unsupported", s.Saved.Transforms[len(s.Saved.Transforms)-1].Kind)
		})
	}
}
func TestRuntimeSourceEvidenceUsesActualMappingAndOrderedPipeline(t *testing.T) {
	f := newRuntimeSourceFixture(t)
	desired := float64(99)
	f.rule.ScaleMultiplier = &desired
	f.mapping.TransformPipeline = `[{"type":"cast","order":2,"params":{"target_type":"float32"}},{"type":"scale","order":1,"params":{"scale":-0.1,"offset":2.5}}]`
	s := f.source(t)
	require.Equal(t, "scale", s.Saved.Transforms[1].Kind)
	require.Equal(t, -0.1, *s.Saved.Transforms[1].Multiplier)
	require.Equal(t, 2.5, *s.Saved.Transforms[1].Offset)
	require.Equal(t, "cast", s.Saved.Transforms[2].Kind)
	raw, _ := json.Marshal(s.Saved.Transforms)
	require.False(t, strings.Contains(string(raw), "99"))
}

func TestRuntimeSourceEvidenceEqualOrderRetainsSavedPipelineSequence(t *testing.T) {
	f := newRuntimeSourceFixture(t)
	f.mapping.TransformPipeline = `[{"type":"cast","order":0,"params":{"target_type":"float32"}},{"type":"scale","order":0,"params":{"scale":0.1,"offset":0}}]`
	s := f.source(t)
	require.Equal(t, "cast", s.Saved.Transforms[1].Kind)
	require.Equal(t, "scale", s.Saved.Transforms[2].Kind)
}

func TestRuntimeSourceEvidenceKnownEmptyPipelineIsNotUnsupported(t *testing.T) {
	for _, pipeline := range []string{"", "  ", "[]", " [ ] ", "null"} {
		t.Run(pipeline, func(t *testing.T) {
			f := newRuntimeSourceFixture(t)
			f.mapping.TransformPipeline = pipeline
			s := f.source(t)
			require.Len(t, s.Saved.Transforms, 1)
			require.Equal(t, "byte-order", s.Saved.Transforms[0].Kind)
		})
	}
}
