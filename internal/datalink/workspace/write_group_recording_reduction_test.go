package workspace

import (
	"database/sql"
	"encoding/json"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

// TestRecordingReductionCompatibility pins the profile-level recording reduction
// contract: latest is the semantics of an omitted reduction, an explicit mean
// survives every JSON path, unknown values and a mean cumulative profile are
// rejected, and omitted versus explicit latest stay one intent.
func TestRecordingReductionCompatibility(t *testing.T) {
	t.Run("omitted reduction stays absent on the wire", func(t *testing.T) {
		var profile WriteGroupRecordingProfile
		require.NoError(t, json.Unmarshal([]byte(`{"kind":"cumulative","interval":10,"unit":"minute"}`), &profile))
		encoded, err := json.Marshal(profile)
		require.NoError(t, err)
		require.NotContains(t, string(encoded), "sample_reduction",
			"a legacy profile must not grow a reduction key it never saved")
	})

	t.Run("omitted and explicit latest are one intent", func(t *testing.T) {
		var omitted, explicit WriteGroupRecordingProfile
		require.NoError(t, json.Unmarshal([]byte(`{"kind":"instantaneous","interval":15,"unit":"second"}`), &omitted))
		require.NoError(t, json.Unmarshal([]byte(`{"kind":"instantaneous","interval":15,"unit":"second","sample_reduction":"latest"}`), &explicit))
		require.True(t, sameRecordingProfile(&omitted, &explicit),
			"an omitted reduction and an explicit latest must replay the same saved group")
		require.True(t, sameRecordingProfile(&explicit, &omitted))
	})

	t.Run("mean is a different intent from latest", func(t *testing.T) {
		var latest, mean WriteGroupRecordingProfile
		require.NoError(t, json.Unmarshal([]byte(`{"kind":"instantaneous","interval":15,"unit":"second","sample_reduction":"latest"}`), &latest))
		require.NoError(t, json.Unmarshal([]byte(`{"kind":"instantaneous","interval":15,"unit":"second","sample_reduction":"mean"}`), &mean))
		require.False(t, sameRecordingProfile(&latest, &mean),
			"a mean request must conflict with a saved latest group instead of replaying it")
	})

	t.Run("unknown reduction is rejected at decode", func(t *testing.T) {
		var profile WriteGroupRecordingProfile
		err := json.Unmarshal([]byte(`{"kind":"instantaneous","interval":15,"unit":"second","sample_reduction":"average"}`), &profile)
		require.Error(t, err, "an unknown reduction must fail closed instead of decoding as latest")
		require.ErrorIs(t, err, ErrWriteGroupValidation)
	})

	t.Run("mean with a cumulative profile is rejected", func(t *testing.T) {
		var profile WriteGroupRecordingProfile
		require.NoError(t, json.Unmarshal([]byte(`{"kind":"cumulative","interval":10,"unit":"minute","sample_reduction":"mean"}`), &profile))
		err := profile.validate()
		require.Error(t, err, "a cumulative profile counts and never averages")
		require.ErrorIs(t, err, ErrWriteGroupValidation)
	})

	t.Run("an explicit empty, null or non-string reduction is rejected", func(t *testing.T) {
		for name, payload := range map[string]string{
			"empty string": `{"kind":"instantaneous","interval":15,"unit":"second","sample_reduction":""}`,
			"json null":    `{"kind":"instantaneous","interval":15,"unit":"second","sample_reduction":null}`,
			"number":       `{"kind":"instantaneous","interval":15,"unit":"second","sample_reduction":3}`,
			"object":       `{"kind":"instantaneous","interval":15,"unit":"second","sample_reduction":{"mode":"mean"}}`,
		} {
			var profile WriteGroupRecordingProfile
			err := json.Unmarshal([]byte(payload), &profile)
			require.Error(t, err, "%s must not silently degrade to latest", name)
			require.ErrorIs(t, err, ErrWriteGroupValidation, "%s must fail as a typed validation error", name)
		}
	})
}

// TestRecordingReductionCompatibility_SQLPersistence pins the durable half of
// the reduction contract: an explicit mean survives create, reload and update
// through the real SQL rows; an omitted reduction stays NULL and is never
// backfilled; ensure replays stay idempotent; and latest versus mean are one
// canonical identity with different intents.
func TestRecordingReductionCompatibility_SQLPersistence(t *testing.T) {
	meanProfile := func() *WriteGroupRecordingProfile {
		return &WriteGroupRecordingProfile{Kind: WriteGroupProfileInstantaneous, Interval: 15, Unit: WriteGroupIntervalUnitSecond, SampleReduction: WriteGroupSampleReductionMean}
	}
	latestProfile := func() *WriteGroupRecordingProfile {
		return &WriteGroupRecordingProfile{Kind: WriteGroupProfileInstantaneous, Interval: 15, Unit: WriteGroupIntervalUnitSecond, SampleReduction: WriteGroupSampleReductionLatest}
	}
	omittedProfile := func() *WriteGroupRecordingProfile {
		return &WriteGroupRecordingProfile{Kind: WriteGroupProfileInstantaneous, Interval: 15, Unit: WriteGroupIntervalUnitSecond}
	}

	t.Run("mean survives create, reload and update", func(t *testing.T) {
		h := newCircuitProfileHarness(t, ":memory:")
		draft := h.draft(WriteGroupProfileInstantaneous, "circuit-mean", "Mean", h.primary())
		draft.RecordingProfile = meanProfile()
		draft.Members[0] = memberWithJSON(t, draft.Members[0], `{"value_role":"gauge"}`)
		created := h.mustEnsure(draft)
		require.Equal(t, WriteGroupSampleReductionMean, created.RecordingProfile.SampleReduction)

		reloaded, err := h.service.Get(t.Context(), created.ID)
		require.NoError(t, err)
		require.Equal(t, WriteGroupSampleReductionMean, reloaded.Group.RecordingProfile.SampleReduction,
			"a saved mean must survive the SQL row, not degrade to latest")

		var stored string
		require.NoError(t, h.db.QueryRowContext(t.Context(),
			`SELECT sample_reduction FROM write_groups WHERE id = ?`, created.ID).Scan(&stored))
		require.Equal(t, "mean", stored, "the reduction must be persisted, not only projected")

		current := reloaded.Group
		current.Name = "Mean renamed"
		mutation := h.mutation(t.Context(), current)
		mutation.ExpectedGroupRevision = current.Revision
		updated, err := h.service.Update(t.Context(), current.ID, mutation)
		require.NoError(t, err)
		require.Equal(t, WriteGroupSampleReductionMean, updated.Group.RecordingProfile.SampleReduction,
			"the member-table rewrite of an update must keep the saved reduction")
	})

	t.Run("an omitted reduction stays NULL and is never backfilled", func(t *testing.T) {
		h := newCircuitProfileHarness(t, ":memory:")
		draft := h.draft(WriteGroupProfileInstantaneous, "circuit-legacy", "Legacy", h.primary())
		draft.RecordingProfile = omittedProfile()
		created := h.mustEnsure(draft)
		var stored sql.NullString
		require.NoError(t, h.db.QueryRowContext(t.Context(),
			`SELECT sample_reduction FROM write_groups WHERE id = ?`, created.ID).Scan(&stored))
		require.False(t, stored.Valid, "an omitted reduction must not be written as an explicit value")
		reloaded, err := h.service.Get(t.Context(), created.ID)
		require.NoError(t, err)
		require.Empty(t, reloaded.Group.RecordingProfile.SampleReduction)
		require.Equal(t, WriteGroupSampleReductionLatest, reloaded.Group.RecordingProfile.EffectiveSampleReduction())
	})

	t.Run("a resent mean ensure replays the same single group", func(t *testing.T) {
		h := newCircuitProfileHarness(t, ":memory:")
		draft := h.draft(WriteGroupProfileInstantaneous, "circuit-mean", "Mean", h.primary())
		draft.RecordingProfile = meanProfile()
		draft.Members[0] = memberWithJSON(t, draft.Members[0], `{"value_role":"gauge"}`)
		first := h.mustEnsure(draft)
		again := h.mustEnsure(draft)
		require.Equal(t, first.ID, again.ID, "a reload or resend must resolve to the saved group")
		require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM write_groups`))
		require.Equal(t, WriteGroupSampleReductionMean, again.RecordingProfile.SampleReduction)
	})

	t.Run("a mean request conflicts with a saved latest group", func(t *testing.T) {
		h := newCircuitProfileHarness(t, ":memory:")
		draft := h.draft(WriteGroupProfileInstantaneous, "circuit-mean", "Mean", h.primary())
		draft.RecordingProfile = latestProfile()
		saved := h.mustEnsure(draft)

		meanDraft := h.draft(WriteGroupProfileInstantaneous, "circuit-mean", "Mean", h.primary())
		meanDraft.RecordingProfile = meanProfile()
		meanDraft.Members[0] = memberWithJSON(t, meanDraft.Members[0], `{"value_role":"gauge"}`)
		_, err := h.ensure(meanDraft)
		require.ErrorIs(t, err, ErrWriteGroupCircuitProfileConflict,
			"mean is a different intent than the saved latest")
		require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM write_groups`))
		reloaded, err := h.service.Get(t.Context(), saved.ID)
		require.NoError(t, err)
		require.Equal(t, WriteGroupSampleReductionLatest, reloaded.Group.RecordingProfile.SampleReduction)
	})

	t.Run("an omitted reduction and an explicit latest replay the same group", func(t *testing.T) {
		h := newCircuitProfileHarness(t, ":memory:")
		draft := h.draft(WriteGroupProfileInstantaneous, "circuit-mean", "Mean", h.primary())
		draft.RecordingProfile = latestProfile()
		first := h.mustEnsure(draft)
		resent := h.draft(WriteGroupProfileInstantaneous, "circuit-mean", "Mean", h.primary())
		resent.RecordingProfile = omittedProfile()
		again := h.mustEnsure(resent)
		require.Equal(t, first.ID, again.ID)
		require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM write_groups`))
		var stored sql.NullString
		require.NoError(t, h.db.QueryRowContext(t.Context(),
			`SELECT sample_reduction FROM write_groups WHERE id = ?`, first.ID).Scan(&stored))
		require.Equal(t, "latest", stored.String, "the saved explicit value is never rewritten by a replay")
	})
}

// TestRecordingReductionCompatibility_Digest pins the managed source digest:
// an omitted reduction and an explicit latest marshal to the same canonical
// scope so a no-op resend never invalidates a confirmed long schema, while a
// mean is a different scope that forces a fresh confirmation. A wide group's
// digest never carried the profile and must not start.
func TestRecordingReductionCompatibility_Digest(t *testing.T) {
	tagTypes := map[string]schema.DataType{"a": schema.DataTypeFloat32}
	withReduction := func(reduction WriteGroupSampleReduction) *WriteGroup {
		group := longGroupFixture(WriteGroupProfileInstantaneous, WriteGroupStorageModePlain, false)
		group.RecordingProfile.SampleReduction = reduction
		return group
	}

	legacy, err := managedSourceDigest(withReduction(""), tagTypes)
	require.NoError(t, err)
	explicitLatest, err := managedSourceDigest(withReduction(WriteGroupSampleReductionLatest), tagTypes)
	require.NoError(t, err)
	mean, err := managedSourceDigest(withReduction(WriteGroupSampleReductionMean), tagTypes)
	require.NoError(t, err)

	require.Equal(t, legacy, explicitLatest,
		"an explicit latest must not invalidate a confirmed schema that omitted the reduction")
	require.NotEqual(t, legacy, mean,
		"a mean reduction is a different confirmed scope and must force a fresh confirmation")

	wide := longGroupFixture(WriteGroupProfileInstantaneous, WriteGroupStorageModePlain, false)
	wide.LayoutVersion = ""
	before, err := managedSourceDigest(wide, tagTypes)
	require.NoError(t, err)
	wide.RecordingProfile.SampleReduction = WriteGroupSampleReductionMean
	after, err := managedSourceDigest(wide, tagTypes)
	require.NoError(t, err)
	require.Equal(t, before, after, "a wide group's digest never carried the profile")
}

// TestRecordingReductionCompatibility_MeanIsBlockedUntilRuntimeProof pins the
// fail-closed runtime gate: a saved mean group stays savable but readiness
// reports it unwritable, so Apply and recording start cannot begin, while a
// legacy latest group keeps its ready path unchanged.
func TestRecordingReductionCompatibility_MeanIsBlockedUntilRuntimeProof(t *testing.T) {
	meanGroup := longGroupFixture(WriteGroupProfileInstantaneous, WriteGroupStorageModePlain, false)
	require.NoError(t, json.Unmarshal([]byte(`{"kind":"instantaneous","interval":15,"unit":"second","sample_reduction":"mean"}`), meanGroup.RecordingProfile))
	readiness := finalizedLong(meanGroup)
	require.False(t, readiness.Ready, "a mean group has no runtime proof yet and must not start a writer")
	require.Contains(t, issueCodesOf(readiness), "recording-reduction-unavailable")

	legacy := longGroupFixture(WriteGroupProfileInstantaneous, WriteGroupStorageModePlain, false)
	legacyReadiness := finalizedLong(legacy)
	require.True(t, legacyReadiness.Ready, "%v", issueCodesOf(legacyReadiness))
	require.Empty(t, legacyReadiness.Issues, "a group without a reduction keeps the existing ready path")
}

// TestRecordingReductionCompatibility_JSONMeanRetained pins the first recording
// reduction contract: an explicit sample_reduction=mean on a recording profile
// must survive a JSON decode/encode round trip. The profile is the only carrier
// of the shared reduction semantics, so a silently dropped mean would downgrade
// an operator's explicit selection to latest on every saved reload.
func TestRecordingReductionCompatibility_JSONMeanRetained(t *testing.T) {
	payload := []byte(`{"kind":"instantaneous","interval":15,"unit":"second","sample_reduction":"mean"}`)
	var profile WriteGroupRecordingProfile
	require.NoError(t, json.Unmarshal(payload, &profile), "an explicit mean reduction must decode")

	encoded, err := json.Marshal(profile)
	require.NoError(t, err)
	var wire map[string]any
	require.NoError(t, json.Unmarshal(encoded, &wire))
	require.Equal(t, "mean", wire["sample_reduction"],
		"an explicit sample_reduction=mean must survive the JSON round trip, not be silently dropped")
}
