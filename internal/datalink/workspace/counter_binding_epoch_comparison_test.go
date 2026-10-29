package workspace

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestCounterBindingExplicitApplyComparesEveryReviewedMeaningFact(t *testing.T) {
	for name, statement := range map[string]string{
		"definition": `UPDATE measurement_definitions SET definition_revision='def-2'`,
		"binding":    `UPDATE measurement_definitions SET source_binding_revision='binding-2'`,
		"epoch":      `UPDATE measurement_definitions SET series_epoch='epoch-2'`,
		"unit":       `UPDATE measurement_definitions SET unit='MWh'`,
		"rollover":   `UPDATE measurement_definitions SET counter_policy='{"rollover_modulus":"9007199254740994","max_rate_per_second":"0.125"}'`,
		"rate":       `UPDATE measurement_definitions SET counter_policy='{"rollover_modulus":"9007199254740993","max_rate_per_second":"0.25"}'`,
		"type":       `UPDATE measurement_definitions SET numeric_encoding='int64'; UPDATE tags SET data_type='int64'; UPDATE points SET data_type='int64'`,
	} {
		t.Run(name, func(t *testing.T) {
			f := newCounterApplyFixture(t, WriteGroupLayoutManagedLongV1, WriteGroupProfileCumulative, schema.DataTypeUint64)
			first, err := f.service.Apply(t.Context(), f.group.ID, f.mutation)
			require.NoError(t, err)
			old, original := counterEpochRuntime(t, f, first.Group.AppliedRevision)
			_, err = f.db.Exec(statement)
			require.NoError(t, err)
			draft := cloneWriteGroup(first.Group)
			draft.Name = "explicit reviewed facts"
			if name == "type" {
				draft.Members[0].SourceRevision = ""
			} // explicit resave resolves the changed point/type source revision
			mutation := counterEpochMutation(t, f, first.Group)
			mutation.Group = draft
			saved, err := f.service.Update(t.Context(), f.group.ID, mutation)
			require.NoError(t, err)
			saveManagedReadinessProof(t, f.db, f.group.ID)
			f.service.WithClock(func() time.Time { return old.EffectiveAt.Add(time.Second) })
			next, err := f.service.Apply(t.Context(), f.group.ID, counterEpochMutation(t, f, saved.Group))
			require.NoError(t, err)
			require.Equal(t, saved.Group.Revision, next.Group.AppliedRevision)
			require.NotEqual(t, first.Group.AppliedRevision, next.Group.AppliedRevision, "each independently changed fact is semantic")
			_, retained := counterEpochRuntime(t, f, first.Group.AppliedRevision)
			require.Equal(t, original, retained)
		})
	}
}

func TestCounterBindingDisplayRenameRetainsOriginalAppliedProof(t *testing.T) {
	for _, tc := range []struct {
		name   string
		layout WriteGroupLayoutVersion
		kind   WriteGroupProfileKind
		scale  string
	}{
		{"cumulative", WriteGroupLayoutManagedLongV1, WriteGroupProfileCumulative, ""},
		{"exact-cumulative", WriteGroupLayoutManagedLongV1, WriteGroupProfileCumulative, "0.1"},
		{"raw", WriteGroupLayoutManagedLongV1, WriteGroupProfileInstantaneous, ""},
		{"exact-raw", WriteGroupLayoutManagedLongV1, WriteGroupProfileInstantaneous, "0.1"},
		{"wide-raw", WriteGroupLayoutManagedWideProfileV1, WriteGroupProfileInstantaneous, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newCounterApplyFixtureWithScale(t, tc.layout, tc.kind, schema.DataTypeUint64, tc.scale)
			first, err := f.service.Apply(t.Context(), f.group.ID, f.mutation)
			require.NoError(t, err)
			_, original := counterEpochRuntime(t, f, first.Group.AppliedRevision)
			draft := cloneWriteGroup(first.Group)
			draft.Name = "only a display rename"
			mutation := counterEpochMutation(t, f, first.Group)
			mutation.Group = draft
			saved, err := f.service.Update(t.Context(), f.group.ID, mutation)
			require.NoError(t, err)
			require.NotEqual(t, first.Group.Revision, saved.Group.Revision)
			renamed, err := f.service.Apply(t.Context(), f.group.ID, counterEpochMutation(t, f, saved.Group))
			require.NoError(t, err)
			require.Equal(t, first.Group.AppliedRevision, renamed.Group.AppliedRevision)
			_, retained := counterEpochRuntime(t, f, first.Group.AppliedRevision)
			require.Equal(t, original, retained)
			var count int
			require.NoError(t, f.db.QueryRow(`SELECT COUNT(*) FROM wg_runtime_versions WHERE group_id=?`, f.group.ID).Scan(&count))
			require.Equal(t, 1, count)
		})
	}
}

func TestCounterBindingFreshRevisionDoesNotRepairUnprovableOldEvidence(t *testing.T) {
	for _, damage := range []string{"missing", "conflict", "unsupported"} {
		t.Run(damage, func(t *testing.T) {
			f := newCounterApplyFixture(t, WriteGroupLayoutManagedLongV1, WriteGroupProfileCumulative, schema.DataTypeUint64)
			first, err := f.service.Apply(t.Context(), f.group.ID, f.mutation)
			require.NoError(t, err)
			old, original := counterEpochRuntime(t, f, first.Group.AppliedRevision)
			broken := original
			switch damage {
			case "missing":
				old.CounterBindings = nil
				bytes, err := json.Marshal(old)
				require.NoError(t, err)
				broken = string(bytes)
			case "conflict":
				broken = strings.Replace(original, `"unit":"kWh"`, `"unit":"MWh"`, 1)
			case "unsupported":
				broken = strings.Replace(original, `"version":1`, `"version":99`, 1)
			}
			require.NotEqual(t, original, broken)
			_, err = f.db.Exec(`UPDATE wg_runtime_versions SET payload=? WHERE group_id=?`, broken, f.group.ID)
			require.NoError(t, err)
			draft := cloneWriteGroup(first.Group)
			draft.Name = "explicit new owner, never repair old evidence"
			mutation := counterEpochMutation(t, f, first.Group)
			mutation.Group = draft
			saved, err := f.service.Update(t.Context(), f.group.ID, mutation)
			require.NoError(t, err)
			f.service.WithClock(func() time.Time { return old.EffectiveAt.Add(time.Second) })
			next, err := f.service.Apply(t.Context(), f.group.ID, counterEpochMutation(t, f, saved.Group))
			require.NoError(t, err)
			require.Equal(t, saved.Group.Revision, next.Group.AppliedRevision)
			require.NotEqual(t, first.Group.AppliedRevision, next.Group.AppliedRevision)
			var retained string
			require.NoError(t, f.db.QueryRow(`SELECT payload FROM wg_runtime_versions WHERE group_id=? AND group_revision=?`, f.group.ID, first.Group.AppliedRevision).Scan(&retained))
			require.Equal(t, broken, retained)
			proven, _ := counterEpochRuntime(t, f, next.Group.AppliedRevision)
			require.Len(t, proven.CounterBindings, 1)
		})
	}
}

// Exact-source already made scale part of group semantics before this change;
// keep that path protected rather than claiming it was the unscaled RED.
func TestCounterBindingExplicitExactScaleChangeCreatesNewOwner(t *testing.T) {
	f := newCounterApplyFixtureWithScale(t, WriteGroupLayoutManagedLongV1, WriteGroupProfileCumulative, schema.DataTypeUint64, "0.1")
	first, err := f.service.Apply(t.Context(), f.group.ID, f.mutation)
	require.NoError(t, err)
	old, original := counterEpochRuntime(t, f, first.Group.AppliedRevision)
	_, err = f.db.Exec(`UPDATE measurement_definitions SET exact_recording_scale='0.01'`)
	require.NoError(t, err)
	draft := cloneWriteGroup(first.Group)
	draft.Members[0].ExactRecordingSource.Scale = "0.01"
	mutation := counterEpochMutation(t, f, first.Group)
	mutation.Group = draft
	saved, err := f.service.Update(t.Context(), f.group.ID, mutation)
	require.NoError(t, err)
	saveManagedReadinessProof(t, f.db, f.group.ID)
	f.service.WithClock(func() time.Time { return old.EffectiveAt.Add(time.Second) })
	next, err := f.service.Apply(t.Context(), f.group.ID, counterEpochMutation(t, f, saved.Group))
	require.NoError(t, err)
	require.Equal(t, saved.Group.Revision, next.Group.AppliedRevision)
	require.NotEqual(t, first.Group.AppliedRevision, next.Group.AppliedRevision)
	frozen, _ := counterEpochRuntime(t, f, next.Group.AppliedRevision)
	for _, proof := range frozen.CounterBindings {
		require.Equal(t, "0.01", proof.ExactScale)
	}
	_, retained := counterEpochRuntime(t, f, first.Group.AppliedRevision)
	require.Equal(t, original, retained)
}

func TestCounterBindingMeaningComparisonDoesNotConvertInfrastructureErrorToReuse(t *testing.T) {
	f := newCounterApplyFixture(t, WriteGroupLayoutManagedLongV1, WriteGroupProfileCumulative, schema.DataTypeUint64)
	first, err := f.service.Apply(t.Context(), f.group.ID, f.mutation)
	require.NoError(t, err)
	old, original := counterEpochRuntime(t, f, first.Group.AppliedRevision)
	tx, err := f.db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	same, err := f.service.sameAppliedCounterMeaningInTx(ctx, tx, &old, old.CounterBindings)
	require.ErrorIs(t, err, context.Canceled)
	require.False(t, same)
	require.NoError(t, tx.Rollback())
	_, retained := counterEpochRuntime(t, f, first.Group.AppliedRevision)
	require.Equal(t, original, retained)
}
