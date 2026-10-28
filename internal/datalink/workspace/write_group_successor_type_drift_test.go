package workspace

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
)

func TestRecreateRepositoryAppliedCumulativeTypeDriftCannotRelabelSource(t *testing.T) {
	for _, mode := range []string{"saved_revision", "refresh_legacy_omission", "refresh_explicit_float"} {
		t.Run(mode, func(t *testing.T) {
			f := newCounterApplyFixture(t, WriteGroupLayoutManagedLongV1, WriteGroupProfileCumulative, schema.DataTypeUint64)
			f.db.SetMaxOpenConns(1)
			applied, err := f.service.Apply(t.Context(), f.group.ID, f.mutation)
			require.NoError(t, err)
			frozen, originalPayload := readCounterApplySnapshot(t, f.db.QueryRowContext, f.group.ID)
			member := frozen.Group.Members[0]
			require.Equal(t, schema.DataTypeUint64, frozen.RuntimeLayout.TagTypes[member.TagID])
			require.Empty(t, member.ValueRole)
			require.Empty(t, member.SourceValueType)
			require.Empty(t, member.ResultValueType)
			require.Nil(t, member.ExactRecordingSource, "exercise cumulative long's unscaled legacy type omission")
			h := &circuitProfileHarness{t: t, db: f.db, service: f.service, workspace: f.service.workspaceSvc}
			lineage := requireLineageGenesis(t, f.db, applied.Group)
			_, err = f.db.Exec("UPDATE write_groups SET status='deleted' WHERE id=?", applied.Group.ID)
			require.NoError(t, err)
			// Only the live tag type changes; source tuple, mapping, and measurement
			// identity/binding/epoch remain exactly those in the applied uint64 source.
			_, err = f.db.Exec("UPDATE tags SET data_type='float64' WHERE id=?", member.TagID)
			require.NoError(t, err)
			var encoding, definition, binding, epoch string
			require.NoError(t, f.db.QueryRow("SELECT numeric_encoding,definition_revision,source_binding_revision,series_epoch FROM measurement_definitions WHERE id=?", *member.MeasurementID).Scan(&encoding, &definition, &binding, &epoch))
			require.Equal(t, []string{"uint64", "def-1", "binding-1", "epoch-1"}, []string{encoding, definition, binding, epoch})
			target := cloneWriteGroup(applied.Group)
			if mode != "saved_revision" {
				target.Members[0].SourceRevision = ""
			}
			if mode == "refresh_explicit_float" {
				target.Members[0].SourceValueType = schema.DataTypeFloat64
				target.Members[0].ResultValueType = measurement.ExactFloat64
			}
			req := writeGroupRecreateReservationRequest{WorkspaceID: applied.Group.WorkspaceID, ExpectedWorkspaceRevision: applied.WorkspaceRevision,
				LineageID: lineage, ExpectedLineageRevision: 1, PredecessorGroupID: applied.Group.ID, ExpectedPredecessorRevision: applied.Group.Revision,
				ExpectedPredecessorConnectorRevision: applied.Group.Destination.ConnectorRevision, ExpectedTargetConnectorRevision: applied.Group.Destination.ConnectorRevision,
				OperationID: "O2-type-drift", Group: target}
			before := recreateRepositorySnapshot(t, h)
			_, err = reserveRepositoryTest(t, h, req)
			if mode == "saved_revision" {
				require.ErrorIs(t, err, ErrWriteGroupSourceRevisionConflict)
			} else {
				require.ErrorIs(t, err, ErrWriteGroupValidation)
				var rejected *WriteGroupTransitionError
				require.ErrorAs(t, err, &rejected)
				require.Equal(t, "recording-meaning-change", rejected.ReasonCode)
			}
			require.Equal(t, before, recreateRepositorySnapshot(t, h))
			_, payload := readCounterApplySnapshot(t, f.db.QueryRowContext, f.group.ID)
			require.Equal(t, originalPayload, payload)
		})
	}
}
