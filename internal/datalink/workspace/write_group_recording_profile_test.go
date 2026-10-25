package workspace

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRecordingProfileIntervalBounds(t *testing.T) {
	type profileCase struct {
		kind        WriteGroupProfileKind
		interval    int
		unit        WriteGroupIntervalUnit
		wantSeconds int
	}
	accepted := []profileCase{
		{WriteGroupProfileInstantaneous, 1, WriteGroupIntervalUnitSecond, 1},
		{WriteGroupProfileInstantaneous, 15, WriteGroupIntervalUnitSecond, 15},
		{WriteGroupProfileInstantaneous, 300, WriteGroupIntervalUnitSecond, 300},
		{WriteGroupProfileCumulative, 1, WriteGroupIntervalUnitMinute, 60},
		{WriteGroupProfileCumulative, 10, WriteGroupIntervalUnitMinute, 600},
		{WriteGroupProfileCumulative, 600, WriteGroupIntervalUnitMinute, 36000},
	}
	for _, tc := range accepted {
		t.Run(fmt.Sprintf("accepts %s %d %s", tc.kind, tc.interval, tc.unit), func(t *testing.T) {
			h := newCircuitProfileHarness(t, ":memory:")
			circuit := "circuit-a"
			draft := h.draft(tc.kind, circuit, "A", h.primary())
			draft.RecordingProfile = &WriteGroupRecordingProfile{Kind: tc.kind, Interval: tc.interval, Unit: tc.unit}
			// A contradictory caller value must not survive: the profile is the authority.
			draft.RowPolicy.IntervalSeconds = 999

			created := h.mustEnsure(draft)
			require.Equal(t, tc.wantSeconds, created.RowPolicy.IntervalSeconds)
			require.Equal(t, tc.interval, created.RecordingProfile.Interval, "the public interval keeps its own unit")
			require.Equal(t, tc.unit, created.RecordingProfile.Unit)
			require.Equal(t, tc.wantSeconds, created.RecordingProfile.IntervalSeconds())

			reloaded, err := h.service.Get(t.Context(), created.ID)
			require.NoError(t, err)
			require.Equal(t, tc.wantSeconds, reloaded.Group.RowPolicy.IntervalSeconds)
			require.Equal(t, *created.RecordingProfile, *reloaded.Group.RecordingProfile)
		})
	}

	rejected := []profileCase{
		{WriteGroupProfileInstantaneous, 0, WriteGroupIntervalUnitSecond, 0},
		{WriteGroupProfileInstantaneous, 301, WriteGroupIntervalUnitSecond, 0},
		{WriteGroupProfileInstantaneous, -1, WriteGroupIntervalUnitSecond, 0},
		{WriteGroupProfileInstantaneous, 15, WriteGroupIntervalUnitMinute, 0},
		{WriteGroupProfileInstantaneous, 15, "", 0},
		{WriteGroupProfileInstantaneous, 15, "hour", 0},
		{WriteGroupProfileCumulative, 0, WriteGroupIntervalUnitMinute, 0},
		{WriteGroupProfileCumulative, 601, WriteGroupIntervalUnitMinute, 0},
		{WriteGroupProfileCumulative, 10, WriteGroupIntervalUnitSecond, 0},
		{WriteGroupProfileCumulative, 10, "", 0},
		{WriteGroupProfileCumulative, 600 * 60, WriteGroupIntervalUnitSecond, 0},
	}
	for _, tc := range rejected {
		t.Run(fmt.Sprintf("rejects %s %d %q", tc.kind, tc.interval, tc.unit), func(t *testing.T) {
			h := newCircuitProfileHarness(t, ":memory:")
			draft := h.draft(tc.kind, "circuit-a", "A", h.primary())
			draft.RecordingProfile = &WriteGroupRecordingProfile{Kind: tc.kind, Interval: tc.interval, Unit: tc.unit}

			_, err := h.ensure(draft)
			require.ErrorIs(t, err, ErrWriteGroupProfileIntervalInvalid)
			require.ErrorIs(t, err, ErrWriteGroupValidation)
			require.Zero(t, h.countRows(`SELECT COUNT(*) FROM write_groups`))

			// The generic create path enforces the same boundary.
			draft.Destination.StorageStrategy = WriteGroupStorageStrategyManaged
			_, err = h.service.Create(t.Context(), h.mutation(t.Context(), draft))
			require.ErrorIs(t, err, ErrWriteGroupProfileIntervalInvalid)
			require.Zero(t, h.countRows(`SELECT COUNT(*) FROM write_groups`))
		})
	}
}

func TestRecordingProfileDefaults(t *testing.T) {
	instantaneous, err := DefaultWriteGroupRecordingProfile(WriteGroupProfileInstantaneous)
	require.NoError(t, err)
	require.Equal(t, WriteGroupRecordingProfile{
		Kind: WriteGroupProfileInstantaneous, Interval: 15, Unit: WriteGroupIntervalUnitSecond,
		SampleReduction: WriteGroupSampleReductionMean,
	}, instantaneous, "a new instantaneous recording defaults to a mean with the 15-second interval")
	cumulative, err := DefaultWriteGroupRecordingProfile(WriteGroupProfileCumulative)
	require.NoError(t, err)
	require.Equal(t, WriteGroupRecordingProfile{
		Kind: WriteGroupProfileCumulative, Interval: 10, Unit: WriteGroupIntervalUnitMinute,
		SampleReduction: WriteGroupSampleReductionLatest,
	}, cumulative, "a cumulative profile counts and never averages")
	_, err = DefaultWriteGroupRecordingProfile("average")
	require.ErrorIs(t, err, ErrWriteGroupValidation)
}

func TestRecordingProfileUpdateKeepsBoundsAndInvalidatesEvidence(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	created := h.mustEnsure(h.draft(WriteGroupProfileInstantaneous, "circuit-a", "A", h.primary()))
	require.Equal(t, 15, created.RowPolicy.IntervalSeconds)
	_, err := h.db.ExecContext(t.Context(), `UPDATE write_groups
		SET destination_schema_revision = 'confirmed-revision', destination_schema_digest = 'confirmed-digest' WHERE id = ?`, created.ID)
	require.NoError(t, err)
	saved, err := h.service.Get(t.Context(), created.ID)
	require.NoError(t, err)
	current := saved.Group

	update := func(mutate func(*WriteGroup)) (*WriteGroupSaveResult, error) {
		next := cloneWriteGroup(current)
		mutate(next)
		mutation := h.mutation(t.Context(), next)
		mutation.ExpectedGroupRevision = current.Revision
		return h.service.Update(t.Context(), current.ID, mutation)
	}

	t.Run("a rename keeps the confirmed schema evidence", func(t *testing.T) {
		renamed, err := update(func(g *WriteGroup) { g.Name = "renamed" })
		require.NoError(t, err)
		require.Equal(t, "confirmed-revision", renamed.Group.Destination.SchemaRevision)
		require.Equal(t, "confirmed-digest", renamed.Group.Destination.SchemaDigest)
		current = renamed.Group
	})

	t.Run("an out-of-range edit is rejected and changes nothing", func(t *testing.T) {
		for _, interval := range []int{0, 301} {
			_, err := update(func(g *WriteGroup) { g.RecordingProfile.Interval = interval })
			require.ErrorIs(t, err, ErrWriteGroupProfileIntervalInvalid)
		}
		_, err := update(func(g *WriteGroup) { g.RecordingProfile.Unit = WriteGroupIntervalUnitMinute })
		require.ErrorIs(t, err, ErrWriteGroupProfileIntervalInvalid)
		reloaded, err := h.service.Get(t.Context(), current.ID)
		require.NoError(t, err)
		require.Equal(t, current.Revision, reloaded.Group.Revision)
		require.Equal(t, 15, reloaded.Group.RowPolicy.IntervalSeconds)
	})

	t.Run("a valid interval edit re-derives seconds and drops stale evidence", func(t *testing.T) {
		edited, err := update(func(g *WriteGroup) {
			g.RecordingProfile.Interval = 300
			g.RowPolicy.IntervalSeconds = 15 // stale client value
		})
		require.NoError(t, err)
		require.Equal(t, 300, edited.Group.RowPolicy.IntervalSeconds)
		require.Empty(t, edited.Group.Destination.SchemaRevision)
		require.Empty(t, edited.Group.Destination.SchemaDigest)
		reloaded, err := h.service.Get(t.Context(), current.ID)
		require.NoError(t, err)
		require.Equal(t, 300, reloaded.Group.RowPolicy.IntervalSeconds)
		require.Equal(t, 300, reloaded.Group.RecordingProfile.Interval)
	})
}

func TestRecordingProfileLegacyIntervalIsNotRewritten(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	legacy := newBasicManagedDraft(h.fixture)
	legacy.Destination.StorageStrategy = WriteGroupStorageStrategyManaged
	legacy.RowPolicy.IntervalSeconds = 60
	created, err := h.service.Create(t.Context(), h.mutation(t.Context(), legacy))
	require.NoError(t, err)

	edit := cloneWriteGroup(created.Group)
	edit.Name = "legacy renamed"
	mutation := h.mutation(t.Context(), edit)
	mutation.ExpectedGroupRevision = created.Group.Revision
	updated, err := h.service.Update(t.Context(), created.Group.ID, mutation)
	require.NoError(t, err)
	require.Equal(t, 60, updated.Group.RowPolicy.IntervalSeconds)
	require.Nil(t, updated.Group.RecordingProfile)

	// A legacy interval outside the profile bounds stays legal for a plain group.
	odd := cloneWriteGroup(updated.Group)
	odd.RowPolicy.IntervalSeconds = 7200
	mutation = h.mutation(t.Context(), odd)
	mutation.ExpectedGroupRevision = updated.Group.Revision
	saved, err := h.service.Update(t.Context(), updated.Group.ID, mutation)
	require.NoError(t, err)
	require.Equal(t, 7200, saved.Group.RowPolicy.IntervalSeconds)
}

func TestRecordingProfileJSONRejectsNonIntegerIntervals(t *testing.T) {
	decode := func(interval string) error {
		var profile WriteGroupRecordingProfile
		return json.Unmarshal(fmt.Appendf(nil, `{"kind":"instantaneous","unit":"second"%s}`, interval), &profile)
	}
	for name, field := range map[string]string{
		"fraction":       `,"interval":15.5`,
		"integral float": `,"interval":15.0`,
		"exponent":       `,"interval":1e1`,
		"quoted number":  `,"interval":"15"`,
		"null":           `,"interval":null`,
		"missing":        ``,
		"boolean":        `,"interval":true`,
		"overflow":       `,"interval":99999999999999999999`,
		"array":          `,"interval":[15]`,
	} {
		t.Run(name, func(t *testing.T) {
			require.ErrorIs(t, decode(field), ErrWriteGroupProfileIntervalInvalid)
		})
	}

	t.Run("integers decode, including values the service later bounds", func(t *testing.T) {
		for _, interval := range []string{"1", "15", "300", "0", "-4", "301"} {
			require.NoError(t, decode(`,"interval":`+interval), interval)
		}
	})

	t.Run("a group payload surfaces the same error", func(t *testing.T) {
		var group WriteGroup
		err := json.Unmarshal([]byte(`{"name":"g","recording_profile":{"kind":"cumulative","interval":2.5,"unit":"minute"}}`), &group)
		require.ErrorIs(t, err, ErrWriteGroupProfileIntervalInvalid)
	})

	t.Run("a whole group keeps its profile and circuit through a snapshot encode", func(t *testing.T) {
		original := WriteGroup{
			ID: "g", LayoutVersion: WriteGroupLayoutManagedWideV1,
			RecordingProfile: &WriteGroupRecordingProfile{Kind: WriteGroupProfileCumulative, Interval: 600, Unit: WriteGroupIntervalUnitMinute},
			Members:          []WriteGroupMember{{DeviceID: "d", PointID: "p", TagID: "t", CircuitID: "c", CircuitName: "Circuit"}},
		}
		encoded, err := json.Marshal(&original)
		require.NoError(t, err)
		var decoded WriteGroup
		require.NoError(t, json.Unmarshal(encoded, &decoded))
		require.Equal(t, original.RecordingProfile, decoded.RecordingProfile)
		require.Equal(t, original.LayoutVersion, decoded.LayoutVersion)
		require.Equal(t, original.Members, decoded.Members)
	})

	t.Run("a profile round-trips with its unit", func(t *testing.T) {
		original := WriteGroupRecordingProfile{Kind: WriteGroupProfileCumulative, Interval: 600, Unit: WriteGroupIntervalUnitMinute}
		encoded, err := json.Marshal(original)
		require.NoError(t, err)
		require.JSONEq(t, `{"kind":"cumulative","interval":600,"unit":"minute"}`, string(encoded))
		var decoded WriteGroupRecordingProfile
		require.NoError(t, json.Unmarshal(encoded, &decoded))
		require.Equal(t, original, decoded)
	})
}
