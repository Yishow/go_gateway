package runtime

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/groupdelivery"
	"path/filepath"
	"testing"
)

func TestMixedDurableDeltaSourceProjectionSeparatesBaselineAndResult(t *testing.T) {
	for _, scaled := range []bool{false, true} {
		t.Run(map[bool]string{false: "unscaled", true: "scaled"}[scaled], func(t *testing.T) {
			f := mixedCounterBoundaryFixture(t, scaled)
			db := openDurableDB(t, filepath.Join(t.TempDir(), "journal.db"))
			store := mixedDurableFixture(t, f, db)
			b, err := NewGroupBoundaryContext(t.Context(), f.config)
			require.NoError(t, err)
			g := f.config.Group
			m := g.Members[0]
			scope := groupdelivery.SourceScope{WorkspaceID: g.WorkspaceID, GroupID: g.ID, GroupRevision: g.AppliedRevision, ProfileKind: "mixed", ResultShape: string(g.RecordingProfile.ResultShape), RecordIdentity: g.RecordID, ConnectorID: g.Destination.ConnectorID, ConnectorRevision: g.Destination.ConnectorRevision}
			ref := groupdelivery.SourceRef{DeviceID: m.DeviceID, PointID: m.PointID, TagID: m.TagID, Column: m.TargetColumn, Calculation: string(m.Calculation), CadenceSeconds: 15, Alignment: m.Cadence.Alignment, Required: m.Required, MeasurementID: *m.MeasurementID, SourceType: string(m.SourceValueType), ResultType: string(m.ResultValueType)}
			f.accept(t, b, f.envelope(0, "first", boundaryAt(14), uint64(9007199254740993)))
			require.NoError(t, b.Tick(t.Context(), boundaryAt(15)))
			first, err := store.SourceRecords(t.Context(), scope, ref)
			require.NoError(t, err)
			require.NotNil(t, first.Latest)
			require.True(t, first.Latest.Available, "cause=%s", first.Latest.Cause)
			require.Equal(t, groupdelivery.SourceDeliveryWindowSkipped, first.Latest.Delivery)
			require.Empty(t, first.Latest.EffectKey)
			require.Empty(t, first.Latest.PayloadDigest)
			require.True(t, first.Latest.Calculation.RawPresent)
			require.False(t, first.Latest.Calculation.ResultPresent)
			require.False(t, first.Latest.Calculation.RowEmitted)
			require.NotNil(t, first.Latest.DeltaKnown)
			require.False(t, *first.Latest.DeltaKnown)
			require.Nil(t, first.LastSQLCommitted)
			f.accept(t, b, f.envelope(0, "second", boundaryAt(29), uint64(9007199254741001)))
			require.NoError(t, b.Tick(t.Context(), boundaryAt(30)))
			next, err := store.SourceRecords(t.Context(), scope, ref)
			require.NoError(t, err)
			require.NotNil(t, next.Latest)
			require.True(t, next.Latest.Available, "cause=%s", next.Latest.Cause)
			c := next.Latest.Calculation
			require.Equal(t, "delta", c.Method)
			require.Equal(t, "uint64", c.SourceType)
			require.Equal(t, "decimal", c.ResultType)
			require.True(t, c.ResultPresent)
			require.True(t, c.RawPresent)
			require.True(t, c.RowEmitted)
			require.Equal(t, groupdelivery.SourceDeliveryQueued, next.Latest.Delivery)
			require.Nil(t, next.LastSQLCommitted, "calculation result is not destination commit proof")
			require.NotEmpty(t, next.Latest.PayloadDigest)
			require.NotEmpty(t, c.OutcomeDigest)
		})
	}
}
