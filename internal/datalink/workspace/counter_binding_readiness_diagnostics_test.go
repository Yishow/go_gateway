package workspace

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"testing"
)

// Validation deliberately leaves tagTypes empty after rejecting the original
// measurement. Missing type evidence must not mask the actionable cause.
func TestCounterBindingReadinessRetainsOriginalMeasurementDiagnosis(t *testing.T) {
	for _, tc := range []struct{ name, statement, reason string }{
		{"foreign point", `UPDATE measurement_definitions SET point_id='other-point'`, "counter-measurement-foreign"},
		{"foreign workspace", `UPDATE measurement_definitions SET workspace_id='other-workspace'`, "counter-measurement-foreign"},
		{"deleted", `DELETE FROM measurement_definitions`, "counter-measurement-unavailable"},
		{"unconfirmed epoch", `UPDATE measurement_definitions SET series_epoch=''`, "counter-measurement-unconfirmed"},
		{"not a counter", `UPDATE measurement_definitions SET semantic_kind='gauge'`, "counter-measurement-not-counter"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newCounterApplyFixture(t, WriteGroupLayoutManagedLongV1, WriteGroupProfileCumulative, schema.DataTypeUint64)
			_, err := f.db.Exec(tc.statement)
			require.NoError(t, err)
			ready, err := f.service.Readiness(t.Context(), f.group.ID)
			require.NoError(t, err)
			require.False(t, ready.Ready)
			require.Empty(t, ready.counterBindings)
			var codes []string
			for _, issue := range ready.Issues {
				codes = append(codes, issue.Code)
			}
			require.Contains(t, codes, tc.reason)
			// Failed role/measurement validation may leave no type map. Diagnose
			// the original saved definition while still forbidding any proof.
			tx, err := f.db.BeginTx(t.Context(), nil)
			require.NoError(t, err)
			proofs, issues := f.service.counterBindingEvidenceInTx(t.Context(), tx, f.group, nil, nil)
			require.NoError(t, tx.Rollback())
			require.Empty(t, proofs)
			codes = nil
			for _, issue := range issues {
				codes = append(codes, issue.Code)
			}
			require.Contains(t, codes, tc.reason)
			_, err = f.service.Apply(t.Context(), f.group.ID, f.mutation)
			require.ErrorIs(t, err, ErrWriteGroupApplyNotReady)
			f.assertUnapplied(t)
		})
	}
}

func TestCounterBindingReadinessMissingOrUnsupportedTypeCannotSealProof(t *testing.T) {
	f := newCounterApplyFixture(t, WriteGroupLayoutManagedLongV1, WriteGroupProfileCumulative, schema.DataTypeUint64)
	tx, err := f.db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	def, err := measurement.NewSQLRepository(f.db).GetByIDInTx(t.Context(), tx, "measurement-1")
	require.NoError(t, err)
	definitions := map[string]*measurement.MeasurementDefinition{def.ID: def}
	for _, kind := range []schema.DataType{"", "future-type", schema.DataTypeString, schema.DataTypeBool} {
		for _, withID := range []bool{false, true} {
			group := cloneWriteGroup(f.group)
			if !withID {
				group.Members[0].MeasurementID = nil
			}
			proofs, issues := f.service.counterBindingEvidenceInTx(t.Context(), tx, group, map[string]schema.DataType{group.Members[0].TagID: kind}, definitions)
			require.Empty(t, proofs, "type=%s measurement=%t", kind, withID)
			var codes []string
			for _, issue := range issues {
				codes = append(codes, issue.Code)
			}
			require.Contains(t, codes, "counter-type-unsupported", "type=%s measurement=%t", kind, withID)
		}
	}
}
