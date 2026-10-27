package workspace

import (
	"encoding/json"
	"testing"
	"time"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"

	"github.com/stretchr/testify/require"
)

func counterEpochRuntime(t *testing.T, f counterApplyFixture, revision string) (WriteGroupAppliedSnapshot, string) {
	t.Helper()
	var raw string
	require.NoError(t, f.db.QueryRow(`SELECT payload FROM wg_runtime_versions WHERE group_id=? AND group_revision=?`, f.group.ID, revision).Scan(&raw))
	var frozen WriteGroupAppliedSnapshot
	require.NoError(t, json.Unmarshal([]byte(raw), &frozen))
	return frozen, raw
}

func counterEpochMutation(t *testing.T, f counterApplyFixture, group *WriteGroup) WriteGroupMutation {
	t.Helper()
	current, err := f.service.workspaceSvc.GetOrCreate(t.Context())
	require.NoError(t, err)
	return WriteGroupMutation{WorkspaceID: group.WorkspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision, ExpectedGroupRevision: group.Revision, ExpectedConnectorRevision: group.Destination.ConnectorRevision}
}

func TestCounterBindingExplicitEpochApplyCreatesNewProofAndRetainsAcceptedOwner(t *testing.T) {
	for _, scale := range []string{"", "0.1"} {
		t.Run("scale="+scale, func(t *testing.T) {
			f := newCounterApplyFixtureWithScale(t, WriteGroupLayoutManagedLongV1, WriteGroupProfileCumulative, schema.DataTypeUint64, scale)
			first, err := f.service.Apply(t.Context(), f.group.ID, f.mutation)
			require.NoError(t, err)
			original, originalBytes := counterEpochRuntime(t, f, first.Group.AppliedRevision)
			var proof CounterBindingProof
			for _, value := range original.CounterBindings {
				proof = value
			}
			ref, err := proof.Reference()
			require.NoError(t, err)
			value := measurement.NewUint64(100)
			if scale != "" {
				value, err = measurement.NewDecimal("100.0")
				require.NoError(t, err)
			}
			store := groupdelivery.NewStore(f.db)
			key := groupdelivery.GroupKey{WorkspaceID: f.group.WorkspaceID, GroupID: f.group.ID, GroupRevision: first.Group.AppliedRevision}
			require.NoError(t, store.AppendSample(t.Context(), key, original.EffectiveAt, snapshot.Sample{SampleID: "old-accepted", MemberKey: proof.Member.SampleMemberKey(), ObservedAt: original.EffectiveAt.Add(time.Second), Value: value, Quality: schema.QualityGood, SourceRevision: proof.SourceRevision, MappingRevision: proof.MappingRevision, CounterProof: ref}))
			var payload, digest string
			require.NoError(t, f.db.QueryRow(`SELECT payload,payload_digest FROM wg_delivery_samples WHERE sample_id='old-accepted'`).Scan(&payload, &digest))
			_, err = f.db.Exec(`UPDATE measurement_definitions SET definition_revision='def-2',source_binding_revision='binding-2',series_epoch='epoch-2',unit='MWh' WHERE id='measurement-1'`)
			require.NoError(t, err)
			draft := cloneWriteGroup(first.Group)
			draft.Name = "explicit epoch confirmation"
			if source := draft.Members[0].ExactRecordingSource; source != nil {
				source.DefinitionRevision, source.SourceBindingRevision, source.SeriesEpoch, source.Unit = "def-2", "binding-2", "epoch-2", "MWh"
			}
			mutation := counterEpochMutation(t, f, first.Group)
			mutation.Group = draft
			saved, err := f.service.Update(t.Context(), f.group.ID, mutation)
			require.NoError(t, err)
			// Exact source changes invalidate the schema selection confirmation; this
			// fixture supplies the newly reviewed confirmation before explicit Apply.
			saveManagedReadinessProof(t, f.db, f.group.ID)
			f.service.WithClock(func() time.Time { return original.EffectiveAt.Add(time.Second) })
			second, err := f.service.Apply(t.Context(), f.group.ID, counterEpochMutation(t, f, saved.Group))
			require.NoError(t, err)
			require.Equal(t, saved.Group.Revision, second.Group.AppliedRevision, "new reviewed meaning is not a display-only rename")
			require.NotEqual(t, first.Group.AppliedRevision, second.Group.AppliedRevision)
			current, _ := counterEpochRuntime(t, f, second.Group.AppliedRevision)
			require.Len(t, current.CounterBindings, 1)
			for _, p := range current.CounterBindings {
				require.Equal(t, []string{"def-2", "binding-2", "epoch-2", "MWh"}, []string{p.DefinitionRevision, p.SourceBindingRevision, p.SeriesEpoch, p.Unit})
				require.Equal(t, second.Group.AppliedRevision, p.Member.AppliedRevision)
			}
			_, afterOriginal := counterEpochRuntime(t, f, first.Group.AppliedRevision)
			require.Equal(t, originalBytes, afterOriginal)
			var afterPayload, afterDigest string
			var consumed, versions int
			require.NoError(t, f.db.QueryRow(`SELECT payload,payload_digest,consumed FROM wg_delivery_samples WHERE sample_id='old-accepted'`).Scan(&afterPayload, &afterDigest, &consumed))
			require.Equal(t, payload, afterPayload)
			require.Equal(t, digest, afterDigest)
			require.Zero(t, consumed)
			again, err := f.service.Apply(t.Context(), f.group.ID, counterEpochMutation(t, f, second.Group))
			require.NoError(t, err)
			require.Equal(t, second.Group.AppliedRevision, again.Group.AppliedRevision)
			require.NoError(t, f.db.QueryRow(`SELECT COUNT(*) FROM wg_runtime_versions WHERE group_id=?`, f.group.ID).Scan(&versions))
			require.Equal(t, 2, versions, "retry creates no third owner or snapshot")
		})
	}
}

func TestCounterBindingChangedMeaningCannotReuseTheSameSavedRevision(t *testing.T) {
	f := newCounterApplyFixture(t, WriteGroupLayoutManagedLongV1, WriteGroupProfileCumulative, schema.DataTypeUint64)
	first, err := f.service.Apply(t.Context(), f.group.ID, f.mutation)
	require.NoError(t, err)
	original, raw := counterEpochRuntime(t, f, first.Group.AppliedRevision)
	f.service.WithClock(func() time.Time { return original.EffectiveAt.Add(time.Second) })
	_, err = f.db.Exec(`UPDATE measurement_definitions SET definition_revision='def-2',series_epoch='epoch-2',unit='MWh' WHERE id='measurement-1'`)
	require.NoError(t, err)
	mutation := counterEpochMutation(t, f, first.Group)
	_, err = f.service.Apply(t.Context(), f.group.ID, mutation)
	require.ErrorIs(t, err, ErrWriteGroupRevisionConflict, "changed meaning needs an explicit fresh saved revision")
	_, after := counterEpochRuntime(t, f, first.Group.AppliedRevision)
	require.Equal(t, raw, after)
	current, err := f.service.workspaceSvc.GetOrCreate(t.Context())
	require.NoError(t, err)
	require.Equal(t, mutation.ExpectedWorkspaceRevision, current.DatabaseSetupRevision)
}
