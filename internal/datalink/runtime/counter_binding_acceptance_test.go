package runtime

import (
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func counterAcceptanceFixture(t *testing.T) (*longFixture, *groupdelivery.Store, groupdelivery.GroupKey, []byte) {
	t.Helper()
	f := newCounterFixture(t)
	f.members = f.members[:1]
	f.config.Group.Members = f.members
	group, member := f.config.Group, f.members[0]
	proof, err := workspace.NewCounterBindingProof(workspace.CounterBindingFacts{
		Member:        workspace.CounterBindingMemberIdentity{WorkspaceID: group.WorkspaceID, GroupID: group.ID, AppliedRevision: group.AppliedRevision, DeviceID: member.DeviceID, PointID: member.PointID, TagID: member.TagID},
		MeasurementID: "meas-int", DefinitionRevision: "def-1", SourceBindingRevision: "bind-1", SeriesEpoch: "epoch-1", SemanticKind: measurement.SemanticKindCounter, Unit: "kWh",
		SourceRevision: "src-int", MappingRevision: "map-int", ScaleMode: workspace.CounterScaleUnscaled, ExactScale: "1", SourceValueType: schema.DataTypeInt64, ResultValueType: measurement.ExactInt64,
	})
	require.NoError(t, err)
	f.config.CounterBindings = workspace.CounterBindingProofs{proof.Member.Key(): proof}
	store := groupdelivery.NewStore(openDurableDB(t, filepath.Join(t.TempDir(), "counter-proof.db")))
	key := groupdelivery.GroupKey{WorkspaceID: group.WorkspaceID, GroupID: group.ID, GroupRevision: group.AppliedRevision}
	f.config.Ledger = store.Ledger(key)
	payload, err := json.Marshal(workspace.WriteGroupAppliedSnapshot{WorkspaceID: group.WorkspaceID, GroupID: group.ID, AppliedRevision: group.AppliedRevision, Group: group, RuntimeLayout: &workspace.WriteGroupRuntimeLayout{TagTypes: f.config.TagTypes}, CounterBindings: f.config.CounterBindings})
	require.NoError(t, err)
	return f, store, key, payload
}

func TestCounterProofAcceptanceFailureLeavesMemoryAndBaselineUntouched(t *testing.T) {
	for _, failure := range []string{"missing proof", "journal failure"} {
		t.Run(failure, func(t *testing.T) {
			f, store, key, payload := counterAcceptanceFixture(t)
			if failure == "journal failure" {
				require.NoError(t, store.SaveRuntimeVersion(t.Context(), key, payload))
				_, err := store.DB().ExecContext(t.Context(), `CREATE TRIGGER fail_counter BEFORE INSERT ON wg_delivery_samples BEGIN SELECT RAISE(ABORT, 'injected journal failure'); END`)
				require.NoError(t, err)
			}
			boundary := f.build(t)
			before := boundary.NextClose()
			envelope := f.envelope(0, "counter-1", at(time.Second), int64(100))
			err := boundary.AcceptSample(t.Context(), envelope)
			require.Error(t, err)
			if failure == "missing proof" {
				require.ErrorIs(t, err, counterproof.ErrCounterBindingProofMissing)
			}
			require.Zero(t, boundary.assembler.OpenBuckets())
			require.Empty(t, boundary.assembler.CounterCheckpoints())
			require.Equal(t, before, boundary.NextClose())
			require.Zero(t, countJournal(t, store.DB()))
			if failure == "journal failure" {
				_, err = store.DB().ExecContext(t.Context(), `DROP TRIGGER fail_counter`)
				require.NoError(t, err)
			} else {
				require.NoError(t, store.SaveRuntimeVersion(t.Context(), key, payload))
			}
			require.NoError(t, boundary.AcceptSample(t.Context(), envelope), "a failed attempt cannot become an unjournaled duplicate")
			require.Equal(t, 1, countJournal(t, store.DB()))
			restored, err := store.Restore(t.Context(), key)
			require.NoError(t, err)
			require.Len(t, restored.Samples, 1)
			require.False(t, restored.Samples[0].CounterProof.IsZero())
		})
	}
}

func TestCounterProofAcceptanceRejectsIncompleteConfiguredMaps(t *testing.T) {
	for name, change := range map[string]func(*GroupBoundaryConfig){
		"empty map": func(c *GroupBoundaryConfig) { clear(c.CounterBindings) },
		"no ledger": func(c *GroupBoundaryConfig) { c.Ledger = nil },
		"wrong map key": func(c *GroupBoundaryConfig) {
			for key, p := range c.CounterBindings {
				delete(c.CounterBindings, key)
				c.CounterBindings["wrong"] = p
			}
		},
		"foreign revision": func(c *GroupBoundaryConfig) { c.Group.AppliedRevision = "other"; c.Group.Revision = "other" },
		"wrong result type": func(c *GroupBoundaryConfig) {
			for key, p := range c.CounterBindings {
				facts := p.CounterBindingFacts
				facts.ScaleMode, facts.ExactScale, facts.ResultValueType = workspace.CounterScaleExact, "0.1", measurement.ExactDecimal
				var err error
				c.CounterBindings[key], err = workspace.NewCounterBindingProof(facts)
				require.NoError(t, err)
			}
		},
	} {
		t.Run(name, func(t *testing.T) {
			f, _, _, _ := counterAcceptanceFixture(t)
			change(&f.config)
			_, err := NewGroupBoundaryContext(t.Context(), f.config)
			var typed *counterproof.CounterBindingProofError
			require.ErrorAs(t, err, &typed)
		})
	}
}

func TestCounterProofAcceptanceScaledReferenceKeepsDecimalOnce(t *testing.T) {
	f, store, key, _ := counterAcceptanceFixture(t)
	member := &f.config.Group.Members[0]
	member.SourceValueType, member.ResultValueType = schema.DataTypeInt64, measurement.ExactDecimal
	member.ExactRecordingSource = &workspace.ExactRecordingSource{
		MeasurementID: "meas-int", DefinitionRevision: "def-1", SourceBindingRevision: "bind-1", SeriesEpoch: "epoch-1",
		Unit: "kWh", SemanticKind: measurement.SemanticKindCounter, OriginalDataType: schema.DataTypeInt64, Scale: "0.1",
	}
	for id, previous := range f.config.CounterBindings {
		proof, err := workspace.CounterBindingProofFromExactSource(previous.Member, *member, schema.DataTypeInt64)
		require.NoError(t, err)
		f.config.CounterBindings[id] = proof
	}
	payload, err := json.Marshal(workspace.WriteGroupAppliedSnapshot{WorkspaceID: key.WorkspaceID, GroupID: key.GroupID, AppliedRevision: key.GroupRevision, Group: f.config.Group, RuntimeLayout: &workspace.WriteGroupRuntimeLayout{TagTypes: f.config.TagTypes}, CounterBindings: f.config.CounterBindings})
	require.NoError(t, err)
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), key, payload))
	boundary := f.build(t)
	require.NoError(t, boundary.AcceptSample(t.Context(), f.envelope(0, "scaled-1", at(time.Second), int64(120004))))
	restored, err := store.Restore(t.Context(), key)
	require.NoError(t, err)
	require.Len(t, restored.Samples, 1)
	digits, ok := restored.Samples[0].Value.DecimalDigits()
	require.True(t, ok)
	require.Equal(t, "12000.4", digits)
	require.False(t, restored.Samples[0].CounterProof.IsZero())
}

func TestCounterProofAcceptanceRejectsCompetingCounterMeaning(t *testing.T) {
	for name, mutate := range map[string]func(*CounterSource){
		"definition":     func(c *CounterSource) { c.DefinitionRevision = "def-2" },
		"source binding": func(c *CounterSource) { c.SourceBindingRevision = "bind-2" },
		"epoch":          func(c *CounterSource) { c.SeriesEpoch = "epoch-2" },
		"unit":           func(c *CounterSource) { c.Unit = "MWh" },
		"kind":           func(c *CounterSource) { c.Kind = aggregation.CounterKindSigned },
		"rollover and rate": func(c *CounterSource) {
			modulus, err := measurement.NewDecimal("9007199254740993")
			require.NoError(t, err)
			rate, err := measurement.NewDecimal("0.125")
			require.NoError(t, err)
			c.RolloverModulus, c.MaxRatePerSecond = &modulus, &rate
		},
	} {
		t.Run(name, func(t *testing.T) {
			f, store, key, payload := counterAcceptanceFixture(t)
			require.NoError(t, store.SaveRuntimeVersion(t.Context(), key, payload))
			source := f.config.Counters["tag-int"]
			mutate(&source)
			f.config.Counters["tag-int"] = source
			boundary, err := NewGroupBoundaryContext(t.Context(), f.config)
			if err == nil {
				err = boundary.AcceptSample(t.Context(), f.envelope(0, "mixed-proof", at(time.Second), int64(100)))
			}
			require.ErrorIs(t, err, counterproof.ErrCounterBindingProofConflict)
			require.Zero(t, countJournal(t, store.DB()))
		})
	}
}
