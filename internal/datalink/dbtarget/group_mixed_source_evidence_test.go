package dbtarget

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/snapshot"
	"testing"
)

func TestMixedSourceEvidenceUsesActualEncoderSelection(t *testing.T) {
	o := mixedEncodeFixture(t, snapshot.ModeMixedWide)[0]
	spec := mixedWideSpec(SQLDialectSQLite)
	for i, m := range spec.Members {
		evidence, err := MixedSourceEvidence(m, o, o.Members[i], true)
		require.NoError(t, err)
		expected, err := json.Marshal(mixedProvenance(m.mixedIntent(), o.Members[i], o.DueAt))
		require.NoError(t, err)
		require.JSONEq(t, string(expected), string(evidence))
		v, err := MixedSelectedValue(m.Calculation, o.Members[i])
		require.NoError(t, err)
		want, err := mixedMemberValue(m.Calculation, o.Members[i])
		require.NoError(t, err)
		require.Equal(t, want, v)
	}
}
func TestMixedSourceEvidenceRequiredFailureIsNotARow(t *testing.T) {
	o := mixedEncodeFixture(t, snapshot.ModeMixedWide)[0]
	m := mixedWideSpec(SQLDialectSQLite).Members[0]
	r := o.Members[0]
	r.Status = snapshot.MemberMissing
	r.Reason = snapshot.ReasonNoSamples
	r.Sample = nil
	r.Reduction = nil
	r.Evidence = nil
	_, err := MixedSourceEvidence(m, o, r, true)
	require.Error(t, err)
	o.Kind = snapshot.OutcomeSkipped
	_, err = MixedSourceEvidence(m, o, r, true)
	require.NoError(t, err)
	_, err = MixedSelectedValue(m.Calculation, r)
	require.Error(t, err)
}
