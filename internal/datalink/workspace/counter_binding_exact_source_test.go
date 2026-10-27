package workspace

import (
	"encoding/json"
	"testing"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func exactCounterBindingFixture() (CounterBindingMemberIdentity, WriteGroupMember) {
	id := "measurement-1"
	return CounterBindingMemberIdentity{WorkspaceID: "workspace-1", GroupID: "group-1", AppliedRevision: "applied-1", DeviceID: "device-1", PointID: "point-1", TagID: "tag-1"}, WriteGroupMember{
		DeviceID: "device-1", PointID: "point-1", TagID: "tag-1", MeasurementID: &id,
		SourceRevision: "source-1", MappingRevision: "mapping-1",
		SourceValueType: schema.DataTypeUint64, ResultValueType: measurement.ExactDecimal,
		ExactRecordingSource: &ExactRecordingSource{
			MeasurementID: id, DefinitionRevision: "def-1", SourceBindingRevision: "binding-1", SeriesEpoch: "epoch-1",
			Unit: "kWh", SemanticKind: measurement.SemanticKindCounter, OriginalDataType: schema.DataTypeUint64, Scale: "0.1",
			CounterPolicy: &measurement.CounterPolicy{RolloverModulus: "9007199254740993", MaxRatePerSecond: "0.125"},
		},
	}
}

func TestCounterBindingExactSourceGolden(t *testing.T) {
	identity, member := exactCounterBindingFixture()
	proof, err := CounterBindingProofFromExactSource(identity, member, schema.DataTypeUint64)
	require.NoError(t, err)
	wire, err := json.Marshal(proof)
	require.NoError(t, err)
	require.Equal(t, counterProofJSON, string(wire), "same literal complete v1 contract as independently supplied facts")
	member.ExactRecordingSource.Unit = "MWh"
	member.ExactRecordingSource.CounterPolicy.MaxRatePerSecond = "999"
	require.Equal(t, "kWh", proof.Unit)
	require.Equal(t, "0.125", proof.CounterPolicy.MaxRatePerSecond)
	require.NoError(t, proof.Validate(), "adapter output cannot alias mutable source policy")
}

func TestCounterBindingExactSourceRejectsContradictoryFacts(t *testing.T) {
	for name, mutate := range map[string]func(*CounterBindingMemberIdentity, *WriteGroupMember){
		"missing source": func(_ *CounterBindingMemberIdentity, m *WriteGroupMember) { m.ExactRecordingSource = nil },
		"foreign device": func(i *CounterBindingMemberIdentity, _ *WriteGroupMember) { i.DeviceID = "other" },
		"foreign point":  func(i *CounterBindingMemberIdentity, _ *WriteGroupMember) { i.PointID = "other" },
		"foreign tag":    func(i *CounterBindingMemberIdentity, _ *WriteGroupMember) { i.TagID = "other" },
		"gauge":          func(_ *CounterBindingMemberIdentity, m *WriteGroupMember) { m.ValueRole = WriteGroupValueRoleGauge },
		"state":          func(_ *CounterBindingMemberIdentity, m *WriteGroupMember) { m.ValueRole = WriteGroupValueRoleState },
		"wrong result": func(_ *CounterBindingMemberIdentity, m *WriteGroupMember) {
			m.ResultValueType = measurement.ExactUint64
		},
		"wrong scale": func(_ *CounterBindingMemberIdentity, m *WriteGroupMember) { m.ExactRecordingSource.Scale = "NaN" },
	} {
		t.Run(name, func(t *testing.T) {
			identity, member := exactCounterBindingFixture()
			mutate(&identity, &member)
			_, err := CounterBindingProofFromExactSource(identity, member, schema.DataTypeUint64)
			require.ErrorIs(t, err, ErrCounterBindingProofConflict)
		})
	}
}

func TestCounterBindingExactSourcePreservesScaleTextAndSignedKind(t *testing.T) {
	identity, member := exactCounterBindingFixture()
	member.ExactRecordingSource.Scale = "0.100"
	member.ExactRecordingSource.SemanticKind = measurement.SemanticKindSignedCounter
	member.ExactRecordingSource.CounterPolicy = &measurement.CounterPolicy{AllowNegative: true}
	proof, err := CounterBindingProofFromExactSource(identity, member, schema.DataTypeUint64)
	require.NoError(t, err)
	require.Equal(t, "0.100", proof.ExactScale)
	require.Equal(t, measurement.SemanticKindSignedCounter, proof.SemanticKind)
	require.Equal(t, measurement.CounterPolicy{AllowNegative: true}, proof.CounterPolicy)
	member.ExactRecordingSource.SemanticKind = measurement.SemanticKindCounter
	member.ExactRecordingSource.CounterPolicy = nil
	proof, err = CounterBindingProofFromExactSource(identity, member, schema.DataTypeUint64)
	require.NoError(t, err)
	require.Equal(t, measurement.CounterPolicy{}, proof.CounterPolicy)
	member.ExactRecordingSource.Scale = "1"
	proof, err = CounterBindingProofFromExactSource(identity, member, schema.DataTypeUint64)
	require.NoError(t, err)
	require.Equal(t, CounterScaleExact, proof.ScaleMode)
	require.Equal(t, measurement.ExactDecimal, proof.ResultValueType, "exact factor one is not the unscaled native-type contract")
}

func TestCounterBindingApplyPersistsExactSourceProof(t *testing.T) {
	for _, kind := range []WriteGroupProfileKind{WriteGroupProfileCumulative, WriteGroupProfileInstantaneous} {
		for _, tagType := range []schema.DataType{schema.DataTypeInt64, schema.DataTypeUint64} {
			t.Run(string(kind)+"/"+string(tagType), func(t *testing.T) {
				f := newCounterApplyFixtureWithScale(t, WriteGroupLayoutManagedLongV1, kind, tagType, "0.1")
				applied, err := f.service.Apply(t.Context(), f.group.ID, f.mutation)
				require.NoError(t, err)
				snapshot, before := readCounterApplySnapshot(t, f.db.QueryRowContext, f.group.ID)
				require.Len(t, snapshot.CounterBindings, 1, "scaled accepted revision must carry the same counter proof contract")
				member := snapshot.Group.Members[0]
				key := CounterBindingMemberIdentity{WorkspaceID: snapshot.WorkspaceID, GroupID: snapshot.GroupID, AppliedRevision: applied.Group.AppliedRevision, DeviceID: member.DeviceID, PointID: member.PointID, TagID: member.TagID}.Key()
				proof := snapshot.CounterBindings[key]
				require.Equal(t, "def-1", proof.DefinitionRevision)
				require.Equal(t, "binding-1", proof.SourceBindingRevision)
				require.Equal(t, "epoch-1", proof.SeriesEpoch)
				require.Equal(t, "kWh", proof.Unit)
				require.Equal(t, CounterScaleExact, proof.ScaleMode)
				require.Equal(t, "0.1", proof.ExactScale)
				require.Equal(t, tagType, proof.SourceValueType)
				require.Equal(t, measurement.ExactDecimal, proof.ResultValueType)
				require.Equal(t, "9007199254740993", proof.CounterPolicy.RolloverModulus)
				require.Equal(t, "0.125", proof.CounterPolicy.MaxRatePerSecond)
				_, err = f.db.ExecContext(t.Context(), `UPDATE measurement_definitions SET unit='MWh', exact_recording_scale='0.5', series_epoch='epoch-2'`)
				require.NoError(t, err)
				_, after := readCounterApplySnapshot(t, f.db.QueryRowContext, f.group.ID)
				require.Equal(t, before, after)
			})
		}
	}
}

func TestCounterBindingUnscaledSourceGolden(t *testing.T) {
	identity, member := exactCounterBindingFixture()
	def := member.ExactRecordingSource.Definition()
	def.ExactRecordingScale = ""
	member.ExactRecordingSource = nil
	member.ResultValueType = measurement.ExactUint64
	group := &WriteGroup{WorkspaceID: identity.WorkspaceID, ID: identity.GroupID, Revision: identity.AppliedRevision}
	proof, err := unscaledCounterBindingProof(group, member, schema.DataTypeUint64, &def)
	require.NoError(t, err)
	wire, err := json.Marshal(proof)
	require.NoError(t, err)
	require.Equal(t, `{"version":1,"member":{"workspace_id":"workspace-1","group_id":"group-1","applied_revision":"applied-1","device_id":"device-1","point_id":"point-1","tag_id":"tag-1"},"measurement_id":"measurement-1","definition_revision":"def-1","source_binding_revision":"binding-1","series_epoch":"epoch-1","semantic_kind":"counter","unit":"kWh","source_revision":"source-1","mapping_revision":"mapping-1","scale_mode":"unscaled","exact_scale":"1","source_value_type":"uint64","result_value_type":"uint64","counter_policy":{"allow_negative":false,"rollover_modulus":"9007199254740993","max_rate_per_second":"0.125"},"digest":"c8bdd515bf15596cb0fb8241688f13f4f72184f64d7af55cc6eb88c8ac1bb0af"}`, string(wire))
}

func TestCounterBindingExactSourceTypeCompatibility(t *testing.T) {
	for _, tagType := range []schema.DataType{schema.DataTypeInt16, schema.DataTypeInt32, schema.DataTypeInt64,
		schema.DataTypeUint16, schema.DataTypeUint32, schema.DataTypeUint64} {
		t.Run(string(tagType), func(t *testing.T) {
			identity, member := exactCounterBindingFixture()
			member.SourceValueType, member.ExactRecordingSource.OriginalDataType = tagType, tagType
			proof, err := CounterBindingProofFromExactSource(identity, member, tagType)
			require.NoError(t, err)
			require.Equal(t, tagType, proof.SourceValueType)
			require.Equal(t, measurement.ExactDecimal, proof.ResultValueType)
		})
	}
	for _, tagType := range []schema.DataType{schema.DataTypeFloat32, schema.DataTypeFloat64, schema.DataTypeBool, schema.DataTypeString, "decimal", "unknown"} {
		t.Run(string(tagType), func(t *testing.T) {
			identity, member := exactCounterBindingFixture()
			member.SourceValueType, member.ExactRecordingSource.OriginalDataType = tagType, tagType
			_, err := CounterBindingProofFromExactSource(identity, member, tagType)
			require.ErrorIs(t, err, ErrCounterBindingProofConflict)
		})
	}
	identity, member := exactCounterBindingFixture()
	_, err := CounterBindingProofFromExactSource(identity, member, schema.DataTypeInt64)
	require.ErrorIs(t, err, ErrCounterBindingProofConflict, "actual tag type must agree with saved source type")
}

func TestCounterBindingExactSourceKeepsUnsupportedLayoutsBlocked(t *testing.T) {
	for name, mutate := range map[string]func(*WriteGroup){
		"wide profile": func(g *WriteGroup) { g.LayoutVersion = WriteGroupLayoutManagedWideProfileV1; g.StorageMode = "" },
		"wide basic":   func(g *WriteGroup) { g.LayoutVersion = WriteGroupLayoutManagedWideV1; g.StorageMode = "" },
		"gauge":        func(g *WriteGroup) { g.Members[0].ValueRole = WriteGroupValueRoleGauge },
		"state":        func(g *WriteGroup) { g.Members[0].ValueRole = WriteGroupValueRoleState },
	} {
		t.Run(name, func(t *testing.T) {
			db := openWorkspaceTestDB(t, ":memory:")
			fixture := seedWriteGroupFixture(t.Context(), t, db)
			id := seedScaledMeasurement(t, db, fixture, "0.1", schema.DataTypeInt64)
			group := makeManagedLongGroup(fixture, fixture.tagID, id)
			mutate(group)
			if !group.isLongLayout() {
				group.Destination.TableName = "scaled_counter"
			}
			if group.Members[0].ValueRole == WriteGroupValueRoleGauge || group.Members[0].ValueRole == WriteGroupValueRoleState {
				group.RecordingProfile.Kind = WriteGroupProfileInstantaneous
				group.RecordingProfile.Interval, group.RecordingProfile.Unit = 15, WriteGroupIntervalUnitSecond
				group.RowPolicy.IntervalSeconds = 15
			}
			tx, err := db.BeginTx(t.Context(), nil)
			require.NoError(t, err)
			defer tx.Rollback()
			err = NewSQLWriteGroupRepository(db).validate(t.Context(), tx, group)
			require.ErrorIs(t, err, ErrWriteGroupValidation)
			require.ErrorContains(t, err, "exact recording scale requires", "refusal must come from the supported exact-source contract, not an invalid fixture")
		})
	}
}
