package workspace

import (
	"database/sql"
	"testing"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func exactSourceMember() WriteGroupMember {
	id := "measurement"
	return WriteGroupMember{
		MeasurementID: &id, SourceValueType: schema.DataTypeInt64, ResultValueType: measurement.ExactDecimal,
		ExactRecordingSource: &ExactRecordingSource{
			MeasurementID: id, DefinitionRevision: "def-1", SourceBindingRevision: "bind-1", SeriesEpoch: "epoch-1",
			OriginalDataType: schema.DataTypeInt64, Scale: "0.1", Unit: "kWh", SemanticKind: measurement.SemanticKindCounter,
		},
	}
}

func TestExactRecordingSourceStrictProof(t *testing.T) {
	for name, mutate := range map[string]func(*WriteGroupMember){
		"missing measurement": func(m *WriteGroupMember) { m.MeasurementID = nil },
		"missing source type": func(m *WriteGroupMember) { m.SourceValueType = "" },
		"missing result type": func(m *WriteGroupMember) { m.ResultValueType = "" },
	} {
		t.Run(name, func(t *testing.T) {
			member := exactSourceMember()
			mutate(&member)
			require.ErrorIs(t, ValidateExactRecordingSource(member, schema.DataTypeInt64), ErrWriteGroupValidation)
		})
	}
	member := exactSourceMember()
	member.ExactRecordingSource.Scale = "NaN"
	_, supported := LongMemberStoredTypeForMember(WriteGroupSampleReductionLatest, member, schema.DataTypeInt64)
	require.False(t, supported, "invalid frozen proof cannot downgrade to a raw integer")
	for _, raw := range []string{"null", "{}", `{"scale":"0.1"}`} {
		_, err := storedExactRecordingSource(sql.NullString{String: raw, Valid: true})
		require.ErrorIs(t, err, ErrWriteGroupValidation, raw)
	}
}

func TestExactRecordingSourceCloneOwnsPolicyPointers(t *testing.T) {
	member := exactSourceMember()
	threshold, rollover := 1.0, 65536.0
	member.ExactRecordingSource.CounterPolicy = &measurement.CounterPolicy{ResetThreshold: &threshold, MaxRollOver: &rollover}
	cloned := cloneWriteGroup(&WriteGroup{Members: []WriteGroupMember{member}})
	*cloned.Members[0].ExactRecordingSource.CounterPolicy.ResetThreshold = 2
	require.Equal(t, 1.0, threshold)
	definition := member.ExactRecordingSource.Definition()
	*definition.CounterPolicy.MaxRollOver = 12
	require.Equal(t, 65536.0, rollover)
}

func TestExactRecordingSourceDerivesBeforeMetadataValidation(t *testing.T) {
	db := openWorkspaceTestDB(t, ":memory:")
	t.Cleanup(func() { _ = db.Close() })
	fixture := seedWriteGroupFixture(t.Context(), t, db)
	id := seedScaledMeasurement(t, db, fixture, "0.1", schema.DataTypeInt64)
	group := makeManagedLongGroup(fixture, fixture.tagID, id)
	group.Members[0].SourceValueType, group.Members[0].ResultValueType = schema.DataTypeInt64, measurement.ExactDecimal
	require.NoError(t, NewSQLWriteGroupRepository(db).Create(t.Context(), group))
	require.NotNil(t, group.Members[0].ExactRecordingSource)
}

func TestExactRecordingSourceEncodingMustMatchTag(t *testing.T) {
	db := openWorkspaceTestDB(t, ":memory:")
	t.Cleanup(func() { _ = db.Close() })
	fixture := seedWriteGroupFixture(t.Context(), t, db)
	id := seedScaledMeasurement(t, db, fixture, "0.1", schema.DataTypeInt64)
	_, err := db.ExecContext(t.Context(), `UPDATE measurement_definitions SET numeric_encoding='uint64' WHERE id=?`, id)
	require.NoError(t, err)
	group := makeManagedLongGroup(fixture, fixture.tagID, id)
	require.ErrorIs(t, NewSQLWriteGroupRepository(db).Create(t.Context(), group), ErrWriteGroupValidation)
}
