package workspace

import (
	"database/sql"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func seedScaledMeasurement(t *testing.T, db *sql.DB, fixture writeGroupFixture, scale string, tagType schema.DataType) string {
	t.Helper()
	ctx := t.Context()
	_, err := db.ExecContext(ctx, `UPDATE points SET data_type = ? WHERE id = ?`, string(tagType), fixture.pointID)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `UPDATE tags SET data_type = ? WHERE id = ?`, string(tagType), fixture.tagID)
	require.NoError(t, err)

	repo := measurement.NewSQLRepository(db)
	measurementID := "meas-scaled-" + fixture.tagID
	def := &measurement.MeasurementDefinition{
		ID:                    measurementID,
		WorkspaceID:           fixture.workspaceID,
		DeviceID:              fixture.deviceID,
		PointID:               fixture.pointID,
		TagID:                 &fixture.tagID,
		EquipmentID:           fixture.deviceID,
		DefinitionRevision:    "def-rev-1",
		SourceBindingRevision: "bind-rev-1",
		SeriesEpoch:           "epoch-1",
		Name:                  "energy",
		Quantity:              "energy",
		Unit:                  "kWh",
		SemanticKind:          measurement.SemanticKindCounter,
		NumericEncoding:       string(tagType),
		ExactRecordingScale:   scale,
		CounterPolicy: &measurement.CounterPolicy{
			AllowNegative:    false,
			RolloverModulus:  "1000000",
			MaxRatePerSecond: "100",
		},
	}
	require.NoError(t, repo.Create(t.Context(), def))
	return measurementID
}

func makeManagedLongGroup(fixture writeGroupFixture, tagID, measurementID string) *WriteGroup {
	return &WriteGroup{
		WorkspaceID:   fixture.workspaceID,
		Name:          "scaled counter group",
		Status:        WriteGroupStatusDraft,
		LayoutVersion: WriteGroupLayoutManagedLongV1,
		StorageMode:   WriteGroupStorageModePlain,
		RecordingProfile: &WriteGroupRecordingProfile{
			Kind:            WriteGroupProfileCumulative,
			Interval:        10,
			Unit:            WriteGroupIntervalUnitMinute,
			SampleReduction: WriteGroupSampleReductionLatest,
		},
		Members: []WriteGroupMember{{
			DeviceID:      fixture.deviceID,
			PointID:       fixture.pointID,
			TagID:         tagID,
			CircuitID:     "circuit-1",
			CircuitName:   "Circuit 1",
			TargetColumn:  "val",
			MeasurementID: &measurementID,
			Required:      true,
		}},
		Destination: WriteGroupDestination{
			ConnectorID:       fixture.connectorID,
			ConnectorRevision: "connector-revision-1",
			Database:          "analytics",
			TableSchema:       "main",
			StorageStrategy:   WriteGroupStorageStrategyManaged,
		},
		RowPolicy: WriteGroupRowPolicy{IntervalSeconds: 600},
	}
}

func TestExactRecordingSource_StaleClientProof(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	fixture := seedWriteGroupFixture(ctx, t, db)
	measurementID := seedScaledMeasurement(t, db, fixture, "0.1", schema.DataTypeInt64)

	repo := NewSQLWriteGroupRepository(db)
	group := makeManagedLongGroup(fixture, fixture.tagID, measurementID)

	// Client sends stale proof with outdated definition revision and different scale
	group.Members[0].ExactRecordingSource = &ExactRecordingSource{
		MeasurementID:         measurementID,
		DefinitionRevision:    "stale-def-rev",
		SourceBindingRevision: "bind-rev-1",
		SeriesEpoch:           "epoch-1",
		OriginalDataType:      schema.DataTypeInt64,
		Scale:                 "0.05",
		Unit:                  "kWh",
		SemanticKind:          measurement.SemanticKindCounter,
	}

	err := repo.Create(ctx, group)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrWriteGroupValidation)
}

func TestExactRecordingSource_BackendDerivedDecimal(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	fixture := seedWriteGroupFixture(ctx, t, db)
	measurementID := seedScaledMeasurement(t, db, fixture, "0.1", schema.DataTypeInt64)

	repo := NewSQLWriteGroupRepository(db)
	group := makeManagedLongGroup(fixture, fixture.tagID, measurementID)
	// Proof omitted by caller
	group.Members[0].ExactRecordingSource = nil

	require.NoError(t, repo.Create(ctx, group))
	require.NotNil(t, group.Members[0].ExactRecordingSource)
	require.Equal(t, "0.1", group.Members[0].ExactRecordingSource.Scale)
	require.Equal(t, schema.DataTypeInt64, group.Members[0].SourceValueType)
	require.Equal(t, measurement.ExactDecimal, group.Members[0].ResultValueType)
	require.Equal(t, "def-rev-1", group.Members[0].ExactRecordingSource.DefinitionRevision)
	require.Equal(t, "bind-rev-1", group.Members[0].ExactRecordingSource.SourceBindingRevision)
	require.Equal(t, "epoch-1", group.Members[0].ExactRecordingSource.SeriesEpoch)
}

func TestExactRecordingSource_OriginalWrongType(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	fixture := seedWriteGroupFixture(ctx, t, db)

	// Update tag data type in database to float64
	_, err := db.ExecContext(ctx, `UPDATE tags SET data_type = 'float64' WHERE id = ?`, fixture.tagID)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `UPDATE points SET data_type = 'float64' WHERE id = ?`, fixture.pointID)
	require.NoError(t, err)

	measurementID := seedScaledMeasurement(t, db, fixture, "0.1", schema.DataTypeFloat64)

	repo := NewSQLWriteGroupRepository(db)
	group := makeManagedLongGroup(fixture, fixture.tagID, measurementID)

	err = repo.Create(ctx, group)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrWriteGroupValidation)
}

func TestExactRecordingSource_LegacyNULL(t *testing.T) {
	ctx := t.Context()
	dbPath := filepath.Join(t.TempDir(), "legacy.db")
	db := openWorkspaceTestDB(t, dbPath)
	fixture := seedWriteGroupFixture(ctx, t, db)

	repo := NewSQLWriteGroupRepository(db)
	group := &WriteGroup{
		WorkspaceID: fixture.workspaceID,
		Name:        "legacy unscaled group",
		Status:      WriteGroupStatusDraft,
		Members: []WriteGroupMember{{
			DeviceID:     fixture.deviceID,
			PointID:      fixture.pointID,
			TagID:        fixture.tagID,
			TargetColumn: "val",
			Required:     true,
		}},
		Destination: WriteGroupDestination{
			ConnectorID:       fixture.connectorID,
			ConnectorRevision: "connector-revision-1",
			Database:          "analytics",
			TableSchema:       "main",
			TableName:         "gw_legacy",
			StorageStrategy:   WriteGroupStorageStrategyCustom,
		},
	}
	require.NoError(t, repo.Create(ctx, group))

	got, err := repo.Get(ctx, fixture.workspaceID, group.ID)
	require.NoError(t, err)
	require.Nil(t, got.Members[0].ExactRecordingSource)
}

func TestExactRecordingSource_RepositoryRoundtrip(t *testing.T) {
	ctx := t.Context()
	dbPath := filepath.Join(t.TempDir(), "roundtrip.db")
	db := openWorkspaceTestDB(t, dbPath)
	fixture := seedWriteGroupFixture(ctx, t, db)
	measurementID := seedScaledMeasurement(t, db, fixture, "0.001", schema.DataTypeInt64)

	repo := NewSQLWriteGroupRepository(db)
	group := makeManagedLongGroup(fixture, fixture.tagID, measurementID)

	require.NoError(t, repo.Create(ctx, group))
	require.NotEmpty(t, group.ID)

	reopenedDB := openWorkspaceTestDB(t, dbPath)
	reopenedRepo := NewSQLWriteGroupRepository(reopenedDB)

	got, err := reopenedRepo.Get(ctx, fixture.workspaceID, group.ID)
	require.NoError(t, err)
	require.NotNil(t, got.Members[0].ExactRecordingSource)
	require.Equal(t, "0.001", got.Members[0].ExactRecordingSource.Scale)
	require.Equal(t, schema.DataTypeInt64, got.Members[0].ExactRecordingSource.OriginalDataType)
	require.Equal(t, "def-rev-1", got.Members[0].ExactRecordingSource.DefinitionRevision)
	require.NotNil(t, got.Members[0].ExactRecordingSource.CounterPolicy)
	require.Equal(t, "1000000", got.Members[0].ExactRecordingSource.CounterPolicy.RolloverModulus)
	require.Equal(t, "100", got.Members[0].ExactRecordingSource.CounterPolicy.MaxRatePerSecond)
	require.Equal(t, measurement.ExactDecimal, got.Members[0].ResultValueType)
}

func TestExactRecordingSource_SameAppliedSemantics_ScaleChangeFalse(t *testing.T) {
	fixture := writeGroupFixture{
		workspaceID: "ws-1",
		deviceID:    "dev-1",
		pointID:     "pt-1",
		tagID:       "tag-1",
		connectorID: "conn-1",
	}
	measID := "meas-1"
	groupA := makeManagedLongGroup(fixture, fixture.tagID, measID)
	groupA.Members[0].ExactRecordingSource = &ExactRecordingSource{
		MeasurementID:         measID,
		DefinitionRevision:    "def-rev-1",
		SourceBindingRevision: "bind-rev-1",
		SeriesEpoch:           "epoch-1",
		OriginalDataType:      schema.DataTypeInt64,
		Scale:                 "0.1",
		Unit:                  "kWh",
		SemanticKind:          measurement.SemanticKindCounter,
	}

	groupB := cloneWriteGroup(groupA)
	require.True(t, sameWriteGroupAppliedSemantics(groupA, groupB))

	// Change scale only
	groupB.Members[0].ExactRecordingSource.Scale = "0.2"
	require.False(t, sameWriteGroupAppliedSemantics(groupA, groupB))
}

func TestValidateExactRecordingSource(t *testing.T) {
	tagType := schema.DataTypeInt64
	measurementID := "meas-1"
	validMember := WriteGroupMember{
		MeasurementID: &measurementID, SourceValueType: tagType, ResultValueType: measurement.ExactDecimal,
		TagID:     "tag-1",
		ValueRole: WriteGroupValueRoleCounterRaw,
		ExactRecordingSource: &ExactRecordingSource{
			MeasurementID:         "meas-1",
			DefinitionRevision:    "def-1",
			SourceBindingRevision: "bind-1",
			SeriesEpoch:           "epoch-1",
			OriginalDataType:      tagType,
			Scale:                 "0.1",
			Unit:                  "kWh",
			SemanticKind:          measurement.SemanticKindCounter,
			CounterPolicy: &measurement.CounterPolicy{
				AllowNegative:    false,
				RolloverModulus:  "1000000",
				MaxRatePerSecond: "100",
			},
		},
	}
	require.NoError(t, ValidateExactRecordingSource(validMember, tagType))

	// Nil proof
	nilMember := validMember
	nilMember.ExactRecordingSource = nil
	require.ErrorIs(t, ValidateExactRecordingSource(nilMember, tagType), ErrWriteGroupValidation)

	// Gauge role rejected
	gaugeMember := validMember
	gaugeMember.ValueRole = WriteGroupValueRoleGauge
	require.ErrorIs(t, ValidateExactRecordingSource(gaugeMember, tagType), ErrWriteGroupValidation)

	// Float tag rejected
	floatMember := validMember
	floatMember.ExactRecordingSource = cloneExactRecordingSource(validMember.ExactRecordingSource)
	floatMember.ExactRecordingSource.OriginalDataType = schema.DataTypeFloat64
	require.ErrorIs(t, ValidateExactRecordingSource(floatMember, schema.DataTypeFloat64), ErrWriteGroupValidation)

	// Bad scale rejected
	badScaleMember := validMember
	badScaleMember.ExactRecordingSource = cloneExactRecordingSource(validMember.ExactRecordingSource)
	badScaleMember.ExactRecordingSource.Scale = "-0.1"
	require.ErrorIs(t, ValidateExactRecordingSource(badScaleMember, tagType), ErrWriteGroupValidation)
}

func TestLongMemberStoredTypeForMember(t *testing.T) {
	tagType := schema.DataTypeInt64
	measurementID := "meas-1"
	validMember := WriteGroupMember{
		MeasurementID: &measurementID, SourceValueType: tagType, ResultValueType: measurement.ExactDecimal,
		TagID:     "tag-1",
		ValueRole: WriteGroupValueRoleCounterRaw,
		ExactRecordingSource: &ExactRecordingSource{
			MeasurementID:         "meas-1",
			DefinitionRevision:    "def-1",
			SourceBindingRevision: "bind-1",
			SeriesEpoch:           "epoch-1",
			OriginalDataType:      tagType,
			Scale:                 "0.1",
			Unit:                  "kWh",
			SemanticKind:          measurement.SemanticKindCounter,
		},
	}
	// With valid proof, always Decimal
	stored, ok := LongMemberStoredTypeForMember(WriteGroupSampleReductionLatest, validMember, tagType)
	require.True(t, ok)
	require.Equal(t, measurement.ExactDecimal, stored)

	// Without proof, counter_raw keeps raw integer
	unscaledMember := validMember
	unscaledMember.ExactRecordingSource = nil
	stored, ok = LongMemberStoredTypeForMember(WriteGroupSampleReductionLatest, unscaledMember, tagType)
	require.True(t, ok)
	require.Equal(t, measurement.ExactInt64, stored)
}

func TestExactRecordingSource_ForwardedCASBoundary_Apply(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	service, fixture, _ := newReadinessGroup(ctx, t, db)
	measurementID := seedScaledMeasurement(t, db, fixture, "0.1", schema.DataTypeInt64)

	ws, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)

	group := makeManagedLongGroup(fixture, fixture.tagID, measurementID)
	created, err := service.Create(ctx, WriteGroupMutation{
		WorkspaceID:               fixture.workspaceID,
		ExpectedWorkspaceRevision: ws.DatabaseSetupRevision,
		ExpectedConnectorRevision: "connector-revision-1",
		Group:                     group,
	})
	require.NoError(t, err)
	require.NotNil(t, created.Group.Members[0].ExactRecordingSource)
	before, err := service.Readiness(ctx, created.Group.ID)
	require.NoError(t, err)
	require.True(t, before.ConfigReady, "the source definition is valid before it changes")

	// Mutate live definition in DB to create drift (bump revision and change scale)
	_, err = db.ExecContext(ctx, `
		UPDATE measurement_definitions
		SET definition_revision = 'def-rev-2', exact_recording_scale = '0.5'
		WHERE id = ?
	`, measurementID)
	require.NoError(t, err)

	after, err := service.Readiness(ctx, created.Group.ID)
	require.NoError(t, err)
	require.False(t, after.ConfigReady, "definition drift must invalidate the source configuration independently of schema readiness")
	foundSourceRefusal := false
	for _, issue := range after.Issues {
		if issue.Code == "group-invalid" {
			foundSourceRefusal = true
		}
	}
	require.True(t, foundSourceRefusal)

	ws, err = service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)

	// Apply must reject drifted proof and refuse to apply
	_, err = service.Apply(ctx, created.Group.ID, WriteGroupMutation{
		WorkspaceID:               fixture.workspaceID,
		ExpectedWorkspaceRevision: ws.DatabaseSetupRevision,
		ExpectedGroupRevision:     created.Group.Revision,
		ExpectedConnectorRevision: "connector-revision-1",
	})
	require.Error(t, err)
}
