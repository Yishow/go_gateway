package groupdelivery

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/measurement"
	"testing"
)

func TestMixedBaselineMustRetainOriginalSealedMeaning(t *testing.T) {
	store, _, _ := newTestStore(t)
	proof, sample, raw := journalProofFixture(t)
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), testKey, raw))
	cp := aggregation.CounterCheckpoint{Binding: aggregation.CounterBinding{WorkspaceID: testKey.WorkspaceID, GroupID: testKey.GroupID, ProfileRevision: testKey.GroupRevision, MemberKey: sample.MemberKey, MeasurementID: proof.MeasurementID, DefinitionRevision: proof.DefinitionRevision, SourceBindingRevision: proof.SourceBindingRevision, SeriesEpoch: proof.SeriesEpoch, Kind: aggregation.CounterKindMonotonic, Unit: proof.Unit}, WindowEnd: at(15), Reading: aggregation.CounterReading{SampleID: sample.SampleID, ObservedAt: sample.ObservedAt, Value: sample.Value}}
	check := func(c aggregation.CounterCheckpoint) error {
		p, e := aggregation.EncodeCounterCheckpoint(c)
		require.NoError(t, e)
		return verifyMixedBaselineProof(t.Context(), store.DB(), testKey, []CounterBaseline{{MemberKey: sample.MemberKey, Payload: p}})
	}
	require.NoError(t, check(cp))
	for name, mutate := range map[string]func(*aggregation.CounterCheckpoint){
		"definition":  func(c *aggregation.CounterCheckpoint) { c.Binding.DefinitionRevision = "new" },
		"measurement": func(c *aggregation.CounterCheckpoint) { c.Binding.MeasurementID = "other" },
		"epoch":       func(c *aggregation.CounterCheckpoint) { c.Binding.SeriesEpoch = "new" },
		"unit":        func(c *aggregation.CounterCheckpoint) { c.Binding.Unit = "MWh" },
		"type":        func(c *aggregation.CounterCheckpoint) { c.Reading.Value = measurement.NewUint64(7) },
		"policy": func(c *aggregation.CounterCheckpoint) {
			v := measurement.NewUint64(100)
			c.Binding.RolloverModulus = &v
		},
	} {
		t.Run(name, func(t *testing.T) { changed := cp; mutate(&changed); require.Error(t, check(changed)) })
	}
}
