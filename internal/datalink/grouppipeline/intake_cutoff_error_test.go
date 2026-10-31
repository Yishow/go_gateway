package grouppipeline

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIntakeCutoffFailureIsReportedWithoutSuccessObservation(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		t.Run(map[bool]string{false: "without observations", true: "with observations"}[enabled], func(t *testing.T) {
			f := newFixture(t)
			p := f.pipeline()
			require.NoError(t, p.Reconcile(t.Context()))
			group := f.groups.groups[0]
			m := p.boundaries[boundaryKey(group.ID, group.AppliedRevision)]
			require.NotNil(t, m)
			var out []IntakeQuotaObservation
			reads := 0
			if enabled {
				intakeCapture(p, &out, &reads)
			}
			before := m.boundary.ReplayState()
			expected := errors.New("cutoff-invalid")
			var reported []error
			p.config.OnError = func(err error) { reported = append(reported, err) }
			p.observedIntakeCutoff(t.Context(), m, func() error { return expected }, true)
			require.Equal(t, []error{expected}, reported)
			require.Empty(t, out, "failure must not claim retiring or paused")
			require.Equal(t, before, m.boundary.ReplayState())
		})
	}
}
