package workspace

import (
	"context"
	"database/sql"
	"testing"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestCounterBindingApplyPersistsUnscaledProofBeforeOwnership(t *testing.T) {
	for _, layout := range []WriteGroupLayoutVersion{WriteGroupLayoutManagedLongV1, WriteGroupLayoutManagedWideProfileV1} {
		for _, kind := range []WriteGroupProfileKind{WriteGroupProfileCumulative, WriteGroupProfileInstantaneous} {
			for _, dataType := range []schema.DataType{schema.DataTypeInt64, schema.DataTypeFloat32} {
				t.Run(string(layout)+"/"+string(kind)+"/"+string(dataType), func(t *testing.T) {
					f := newCounterApplyFixture(t, layout, kind, dataType)
					var prepared WriteGroupAppliedSnapshot
					f.service.WithWriterOwnershipActivationBarrier(writerOwnershipActivationFunc(func(_ context.Context, tx *sql.Tx, _ WriterOwnershipActivationRequest) error {
						prepared, _ = readCounterApplySnapshot(t, tx.QueryRowContext, f.group.ID)
						return nil
					}))
					applied, err := f.service.Apply(t.Context(), f.group.ID, f.mutation)
					require.NoError(t, err)
					require.Len(t, prepared.CounterBindings, 1, "the ownership barrier must already see the complete proof in its transaction")
					member := f.group.Members[0]
					identity := CounterBindingMemberIdentity{WorkspaceID: f.group.WorkspaceID, GroupID: f.group.ID, AppliedRevision: applied.Group.AppliedRevision, DeviceID: member.DeviceID, PointID: member.PointID, TagID: member.TagID}
					proof, ok := prepared.CounterBindings[identity.Key()]
					require.True(t, ok)
					require.Equal(t, "measurement-1", proof.MeasurementID)
					require.Equal(t, "def-1", proof.DefinitionRevision)
					require.Equal(t, "binding-1", proof.SourceBindingRevision)
					require.Equal(t, "epoch-1", proof.SeriesEpoch)
					require.Equal(t, "kWh", proof.Unit)
					require.Equal(t, measurement.SemanticKindCounter, proof.SemanticKind)
					require.Equal(t, CounterScaleUnscaled, proof.ScaleMode)
					require.Equal(t, "1", proof.ExactScale)
					require.Equal(t, dataType, proof.SourceValueType)
					wantResult := measurement.ExactInt64
					if dataType == schema.DataTypeFloat32 {
						wantResult = measurement.ExactFloat64
						require.Empty(t, proof.CounterPolicy.RolloverModulus)
						require.Empty(t, proof.CounterPolicy.MaxRatePerSecond)
					} else {
						require.Equal(t, "9007199254740993", proof.CounterPolicy.RolloverModulus)
						require.Equal(t, "0.125", proof.CounterPolicy.MaxRatePerSecond)
					}
					require.Equal(t, wantResult, proof.ResultValueType)
					require.Equal(t, member.SourceRevision, proof.SourceRevision)
					require.Equal(t, member.MappingRevision, proof.MappingRevision)
					require.NoError(t, proof.Validate())
					persisted, before := readCounterApplySnapshot(t, f.db.QueryRowContext, f.group.ID)
					require.Equal(t, prepared.CounterBindings, persisted.CounterBindings)
					_, err = f.db.ExecContext(t.Context(), `UPDATE measurement_definitions SET definition_revision='def-2', series_epoch='epoch-2', unit='MWh' WHERE id='measurement-1'`)
					require.NoError(t, err)
					_, after := readCounterApplySnapshot(t, f.db.QueryRowContext, f.group.ID)
					require.Equal(t, before, after, "editing current definitions cannot rewrite persisted proof bytes")
				})
			}
		}
	}
}

func TestCounterBindingApplyRejectsDefinitionEditDuringInspection(t *testing.T) {
	f := newCounterApplyFixture(t, WriteGroupLayoutManagedLongV1, WriteGroupProfileCumulative, schema.DataTypeInt64)
	f.inspector.mutate = func() {
		_, err := f.db.ExecContext(t.Context(), `UPDATE measurement_definitions SET definition_revision='def-2', series_epoch='epoch-2', unit='MWh' WHERE id='measurement-1'`)
		require.NoError(t, err)
	}
	_, err := f.service.Apply(t.Context(), f.group.ID, f.mutation)
	require.ErrorIs(t, err, ErrWriteGroupApplyNotReady)
	f.assertUnapplied(t)
}

func TestCounterBindingApplyCASAfterReview(t *testing.T) {
	for name, change := range map[string]string{
		"definition revision":     `definition_revision='def-2'`,
		"source binding revision": `source_binding_revision='binding-2'`,
		"semantic facts":          `unit='MWh', series_epoch='epoch-2'`,
	} {
		t.Run(name, func(t *testing.T) {
			f := newCounterApplyFixture(t, WriteGroupLayoutManagedLongV1, WriteGroupProfileCumulative, schema.DataTypeInt64)
			store := &counterApplyRaceRepository{SQLRepository: NewSQLRepository(f.db), beforeApply: func() {
				_, err := f.db.ExecContext(t.Context(), `UPDATE measurement_definitions SET `+change+` WHERE id='measurement-1'`)
				require.NoError(t, err)
			}}
			f.service.workspaceSvc.repo = store
			_, err := f.service.Apply(t.Context(), f.group.ID, f.mutation)
			require.Equal(t, 3, store.begins)
			require.ErrorIs(t, err, ErrWriteGroupRevisionConflict)
			f.assertUnapplied(t)
		})
	}
}

func TestCounterBindingApplyKeepsIndependentGroupFactsIsolated(t *testing.T) {
	f := newCounterApplyFixture(t, WriteGroupLayoutManagedLongV1, WriteGroupProfileCumulative, schema.DataTypeInt64)
	other := seedAdditionalWriteGroupDevice(t.Context(), t, f.db, f.group.WorkspaceID, f.group.Destination.ConnectorID)
	measurementID := seedScaledMeasurement(t, f.db, other, "", schema.DataTypeInt64)
	_, err := f.db.ExecContext(t.Context(), `UPDATE measurement_definitions SET definition_revision='def-2', unit='MWh', series_epoch='epoch-2' WHERE id=?`, measurementID)
	require.NoError(t, err)
	current, err := f.service.workspaceSvc.GetOrCreate(t.Context())
	require.NoError(t, err)
	otherDraft := makeManagedLongGroup(other, other.tagID, measurementID)
	otherDraft.Members[0].CircuitID = "circuit-2"
	created, err := f.service.Create(t.Context(), WriteGroupMutation{WorkspaceID: other.workspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedConnectorRevision: "connector-revision-1", Group: otherDraft})
	require.NoError(t, err)
	saveManagedReadinessProof(t, f.db, created.Group.ID)
	for _, tc := range []struct {
		group       *WriteGroup
		unit, epoch string
	}{
		{f.group, "kWh", "epoch-1"}, {created.Group, "MWh", "epoch-2"},
	} {
		current, err := f.service.workspaceSvc.GetOrCreate(t.Context())
		require.NoError(t, err)
		_, err = f.service.Apply(t.Context(), tc.group.ID, WriteGroupMutation{WorkspaceID: other.workspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
			ExpectedGroupRevision: tc.group.Revision, ExpectedConnectorRevision: "connector-revision-1"})
		require.NoError(t, err)
		snapshot, _ := readCounterApplySnapshot(t, f.db.QueryRowContext, tc.group.ID)
		require.Len(t, snapshot.CounterBindings, 1)
		for _, proof := range snapshot.CounterBindings {
			require.Equal(t, tc.group.ID, proof.Member.GroupID)
			require.Equal(t, tc.unit, proof.Unit)
			require.Equal(t, tc.epoch, proof.SeriesEpoch)
		}
	}
	first, _ := readCounterApplySnapshot(t, f.db.QueryRowContext, f.group.ID)
	for _, proof := range first.CounterBindings {
		require.Equal(t, "kWh", proof.Unit)
		require.Equal(t, "def-1", proof.DefinitionRevision)
	}
}
