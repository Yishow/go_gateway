package snapshot

import (
	"testing"
	"time"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func mustRecordID(t *testing.T, workspace, group, revision, entity string, bucket time.Time) string {
	t.Helper()
	id, err := RecordID(workspace, group, revision, entity, bucket)
	require.NoError(t, err)
	require.NotEmpty(t, id)
	return id
}

func TestScopedRowIdentityAndLateArrivalRecordIDComponents(t *testing.T) {
	base := mustRecordID(t, "workspace-A", "group-G", "rev-1", "", fixtureStart)
	require.Equal(t, base, mustRecordID(t, "workspace-A", "group-G", "rev-1", "", fixtureStart), "deterministic across calls")
	require.Equal(t, base, mustRecordID(t, "workspace-A", "group-G", "rev-1", "", fixtureStart.In(time.FixedZone("UTC+8", 8*3600))), "time zone does not change the bucket identity")

	for name, other := range map[string]string{
		"workspace": mustRecordID(t, "workspace-B", "group-G", "rev-1", "", fixtureStart),
		"group":     mustRecordID(t, "workspace-A", "group-H", "rev-1", "", fixtureStart),
		"revision":  mustRecordID(t, "workspace-A", "group-G", "rev-2", "", fixtureStart),
		"entity":    mustRecordID(t, "workspace-A", "group-G", "rev-1", "plant-1", fixtureStart),
		"bucket":    mustRecordID(t, "workspace-A", "group-G", "rev-1", "", at(10)),
	} {
		require.NotEqual(t, base, other, name)
	}
	// The fixed group scope can never collide with a literal entity key.
	require.NotEqual(t,
		mustRecordID(t, "w", "g", "r", "", fixtureStart),
		mustRecordID(t, "w", "g", "r", "group-scope", fixtureStart))
	// Delimiter-like values cannot shift between fields.
	require.NotEqual(t,
		mustRecordID(t, "a", "b", "r", "c", fixtureStart),
		mustRecordID(t, "a", "bc", "r", "", fixtureStart))
}

func TestScopedRowIdentityAndLateArrivalRecordIDRequiresIdentity(t *testing.T) {
	_, err := RecordID("", "group-G", "rev-1", "", fixtureStart)
	require.ErrorIs(t, err, ErrInvalidConfig)
	_, err = RecordID("workspace-A", "", "rev-1", "", fixtureStart)
	require.ErrorIs(t, err, ErrInvalidConfig)
	_, err = RecordID("workspace-A", "group-G", "", "", fixtureStart)
	require.ErrorIs(t, err, ErrInvalidConfig)
	_, err = RecordID("workspace-A", "group-G", "rev-1", "", time.Time{})
	require.ErrorIs(t, err, ErrInvalidConfig)
}

func TestScopedRowIdentityAndLateArrivalEffectKeysAreDestinationScoped(t *testing.T) {
	record := mustRecordID(t, "workspace-A", "group-G", "rev-1", "", fixtureStart)
	scopeA, err := DestinationScope("connector-1", "crev-1", "db", "public", "readings")
	require.NoError(t, err)
	scopeB, err := DestinationScope("connector-1", "crev-1", "db", "public", "other_table")
	require.NoError(t, err)
	scopeC, err := DestinationScope("connector-1", "crev-2", "db", "public", "readings")
	require.NoError(t, err)
	require.Len(t, map[string]struct{}{scopeA: {}, scopeB: {}, scopeC: {}}, 3)

	keyA, err := EffectKey(scopeA, record)
	require.NoError(t, err)
	again, err := EffectKey(scopeA, record)
	require.NoError(t, err)
	require.Equal(t, keyA, again, "a finalized row resubmitted keeps its effect key")
	keyB, err := EffectKey(scopeB, record)
	require.NoError(t, err)
	require.NotEqual(t, keyA, keyB)
	require.NotEqual(t, keyA, record, "the effect key is namespaced, not the bare record id")

	_, err = EffectKey("", record)
	require.ErrorIs(t, err, ErrInvalidConfig)
	_, err = DestinationScope("", "crev-1", "db", "public", "readings")
	require.ErrorIs(t, err, ErrInvalidConfig)
	_, err = DestinationScope("connector-1", "crev-1", "db", "public", "")
	require.ErrorIs(t, err, ErrInvalidConfig)
}

func entityAssembler(t *testing.T) *Assembler {
	return newFixtureAssembler(t, func(c *Config) {
		c.Members = []Member{
			{Key: "plant-1/temp", EntityKey: "plant-1", Required: true},
			{Key: "plant-2/temp", EntityKey: "plant-2", Required: true},
		}
	})
}

func TestScopedRowIdentityAndLateArrivalSameTimestampDifferentEntityKeepsBothRows(t *testing.T) {
	a := entityAssembler(t)
	offerGood(t, a, "s1", "plant-1/temp", 5, 6)
	offerGood(t, a, "s2", "plant-2/temp", 5, 6)

	outcomes := a.Tick(at(10))
	require.Len(t, outcomes, 2)
	require.Equal(t, "plant-1", outcomes[0].EntityKey)
	require.Equal(t, "plant-2", outcomes[1].EntityKey)
	for _, outcome := range outcomes {
		require.Equal(t, OutcomeRow, outcome.Kind)
	}
	require.NotEqual(t, outcomes[0].RecordID, outcomes[1].RecordID)
	require.NotEqual(t, outcomes[0].EffectKey, outcomes[1].EffectKey)
	require.Equal(t, mustRecordID(t, "workspace-A", "group-G", "rev-1", "plant-1", fixtureStart), outcomes[0].RecordID)
	wantEffect, err := EffectKey("scope-X", outcomes[0].RecordID)
	require.NoError(t, err)
	require.Equal(t, wantEffect, outcomes[0].EffectKey)

	// A second assembler built from the same configuration (as after restart)
	// derives the same identities for the same bucket.
	b := entityAssembler(t)
	offerGood(t, b, "s1", "plant-1/temp", 5, 6)
	offerGood(t, b, "s2", "plant-2/temp", 5, 6)
	replayed := b.Tick(at(10))
	require.Equal(t, outcomes[0].RecordID, replayed[0].RecordID)
	require.Equal(t, outcomes[1].EffectKey, replayed[1].EffectKey)
}

func TestScopedRowIdentityAndLateArrivalOneEntitySilentStillReportsScopedNoData(t *testing.T) {
	a := entityAssembler(t)
	offerGood(t, a, "s1", "plant-1/temp", 5, 6)

	outcomes := a.Tick(at(10))
	require.Len(t, outcomes, 2)
	require.Equal(t, OutcomeRow, outcomes[0].Kind)
	require.Equal(t, OutcomeNoData, outcomes[1].Kind)
	require.Equal(t, "plant-2", outcomes[1].EntityKey)
	require.NotEmpty(t, outcomes[1].RecordID)
	require.Empty(t, outcomes[1].EffectKey, "only rows have destination effects")
}

func TestScopedRowIdentityAndLateArrivalClosedBucketRejectsLateRewrite(t *testing.T) {
	a := newFixtureAssembler(t, nil)
	offerGood(t, a, "t5", "temperature", 5, 6)
	offerGood(t, a, "p5", "pressure", 5, 6)
	closed := a.Tick(at(10))
	require.Len(t, closed, 1)
	frozen := closed[0]
	snapshotOfRow := frozen.Members[0].Sample.Value

	late := goodSample("t9-late", "temperature", 9, 999)
	result := a.Offer(late, at(11))
	require.Equal(t, OfferLate, result.Outcome)
	require.Equal(t, ReasonLate, result.Reason)
	require.Equal(t, uint64(1), a.Late())
	require.Zero(t, a.OpenBuckets(), "a late sample must not reopen or recreate the closed bucket")
	require.Empty(t, a.Tick(at(11)), "no second outcome for the closed bucket")
	require.True(t, snapshotOfRow.Equal(measurement.NewInt64(5)), "the emitted row is unchanged")

	// The same accepted sample re-sent after closure has no effect either.
	require.Equal(t, OfferLate, a.Offer(goodSample("t5", "temperature", 5, 5), at(11)).Outcome)
	require.Equal(t, uint64(2), a.Late())
	require.Empty(t, a.Tick(at(12)))
}

func TestScopedRowIdentityAndLateArrivalDuplicateBeforeCloseHasNoEffect(t *testing.T) {
	a := newFixtureAssembler(t, nil)
	offerGood(t, a, "t5", "temperature", 5, 6)
	require.Equal(t, OfferDuplicate, a.Offer(goodSample("t5", "temperature", 5, 5), at(7)).Outcome)
	conflict := goodSample("t5", "temperature", 5, 6)
	require.Equal(t, OfferConflict, a.Offer(conflict, at(7)).Outcome)
	offerGood(t, a, "p5", "pressure", 5, 7)

	row := a.Tick(at(10))[0]
	require.Equal(t, OutcomeRow, row.Kind)
	require.True(t, measurement.NewInt64(5).Equal(memberResult(t, row, "temperature").Sample.Value))
}

func TestScopedRowIdentityAndLateArrivalRevisionChangeDoesNotReinterpretOldRows(t *testing.T) {
	withRevisions := func(group, mapping string) func(*Config) {
		return func(c *Config) {
			c.GroupRevision = group
			c.Members[0].MappingRevision = mapping
		}
	}
	oldAssembler := newFixtureAssembler(t, withRevisions("rev-1", "scale-x1"))
	old := goodSample("t5", "temperature", 5, 50)
	old.MappingRevision = "scale-x1"
	require.Equal(t, OfferSelected, oldAssembler.Offer(old, at(6)).Outcome)
	offerGood(t, oldAssembler, "p5", "pressure", 5, 6)
	oldRow := oldAssembler.Tick(at(10))[0]
	require.Equal(t, OutcomeRow, oldRow.Kind)

	newAssembler := newFixtureAssembler(t, withRevisions("rev-2", "scale-x10"))
	stale := newAssembler.Offer(old, at(11))
	require.Equal(t, OfferRejected, stale.Outcome, "a value produced under the old scale is not reinterpreted")
	require.Equal(t, ReasonRevisionMismatch, stale.Reason)

	fresh := goodSample("t15", "temperature", 15, 500)
	fresh.MappingRevision = "scale-x10"
	require.Equal(t, OfferSelected, newAssembler.Offer(fresh, at(16)).Outcome)

	oldID := mustRecordID(t, "workspace-A", "group-G", "rev-1", "", fixtureStart)
	require.Equal(t, oldID, oldRow.RecordID, "accepted old rows keep their original identity")
	require.NotEqual(t, oldRow.RecordID, mustRecordID(t, "workspace-A", "group-G", "rev-2", "", fixtureStart))
}

func TestScopedRowIdentityAndLateArrivalConfigRequiresIdentityFields(t *testing.T) {
	for name, mutate := range map[string]func(*Config){
		"workspace": func(c *Config) { c.WorkspaceID = "" },
		"group":     func(c *Config) { c.GroupID = "" },
		"revision":  func(c *Config) { c.GroupRevision = "" },
		"scope":     func(c *Config) { c.DestinationScope = "" },
	} {
		config := fixtureConfig()
		config.Members = append([]Member(nil), config.Members...)
		mutate(&config)
		_, err := NewAssembler(config)
		require.ErrorIs(t, err, ErrInvalidConfig, name)
	}
}

func TestScopedRowIdentityAndLateArrivalRevisionIsPartOfDuplicateIdentity(t *testing.T) {
	selector := newFixtureSelector()
	first := goodSample("s8", "temperature", 8, 80)
	first.MappingRevision = "scale-x1"
	require.Equal(t, OfferSelected, selector.Offer(first, at(12)).Outcome)
	changed := first
	changed.MappingRevision = "scale-x10"
	require.Equal(t, OfferConflict, selector.Offer(changed, at(12)).Outcome)
	_ = schema.QualityGood
}
