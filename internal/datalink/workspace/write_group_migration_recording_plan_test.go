package workspace

import (
	"context"
	"database/sql"
	"testing"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type recordingPlanMigrationFixture struct {
	db            *sql.DB
	service       *WriteGroupService
	workspaceSvc  *Service
	base          writeGroupFixture
	plan          *recordingplan.RecordingPlan
	measurementID string
}

func newRecordingPlanMigrationFixture(ctx context.Context, t *testing.T) recordingPlanMigrationFixture {
	t.Helper()
	db := openWorkspaceTestDB(t, ":memory:") //nolint:contextcheck // test helper owns the migration setup context
	base := seedWriteGroupFixture(ctx, t, db)
	_, err := db.ExecContext(ctx, `UPDATE devices SET status = ? WHERE id = ?`, schema.DeviceStatusActive, base.deviceID)
	require.NoError(t, err)
	workspaceSvc := NewService(NewSQLRepository(db))
	workspace, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	_, err = workspaceSvc.UpdateDatabaseSetup(ctx, workspace.DatabaseSetupRevision, func(_ context.Context, _ *sql.Tx, record *Record) error {
		record.DatabaseConnectorID = base.connectorID
		return nil
	})
	require.NoError(t, err)
	measurementID := "measurement-A"
	tagID := base.tagID
	measurementDefinition := &measurement.MeasurementDefinition{
		ID:                    measurementID,
		WorkspaceID:           base.workspaceID,
		DeviceID:              base.deviceID,
		PointID:               base.pointID,
		TagID:                 &tagID,
		EquipmentID:           "equipment-A",
		DefinitionRevision:    "measurement-revision-1",
		SourceBindingRevision: "binding-revision-1",
		SeriesEpoch:           "epoch-1",
		Name:                  "Temperature",
		Quantity:              "temperature",
		Unit:                  "C",
		SemanticKind:          measurement.SemanticKindGauge,
		NumericEncoding:       "float64",
	}
	require.NoError(t, measurement.NewSQLRepository(db).Create(ctx, measurementDefinition))
	plan := &recordingplan.RecordingPlan{
		ID:          "legacy-plan-A",
		WorkspaceID: base.workspaceID,
		Revision:    "plan-revision-1",
		Name:        "Legacy plan A",
		Status:      recordingplan.PlanStatusDraft,
		Timezone:    "Asia/Taipei",
		Members: []recordingplan.PlanMember{{
			MemberID: "member-A", MeasurementID: measurementID, EquipmentID: "equipment-A", Name: "Temperature",
		}},
		Streams: []recordingplan.PlanStream{{
			StreamID: "stream-A", MeasurementID: measurementID, Mode: recordingplan.StreamModeRawHistory, RawPolicy: recordingplan.RawPolicyEverySample,
		}},
		Destinations: []recordingplan.PlanDestination{{
			DestinationID: "destination-A", ConnectorID: base.connectorID, ConnectorRevision: "connector-revision-1",
		}},
	}
	require.NoError(t, recordingplan.NewSQLRepository(db).CreatePlan(ctx, plan))
	return recordingPlanMigrationFixture{
		db: db, service: NewWriteGroupService(workspaceSvc, NewSQLWriteGroupRepository(db)), workspaceSvc: workspaceSvc,
		base: base, plan: plan, measurementID: measurementID,
	}
}

func (f recordingPlanMigrationFixture) close() {
	_ = f.db.Close()
}

func recordingPlanReviewRequest(preview *WriteGroupMigrationPreview) WriteGroupMigrationReviewRequest {
	return WriteGroupMigrationReviewRequest{
		WorkspaceID:               preview.WorkspaceID,
		ExpectedWorkspaceRevision: preview.WorkspaceRevision,
		ExpectedConnectorRevision: preview.ConnectorRevision,
		ReviewDigest:              preview.ReviewDigest,
		SourceIDs:                 []string{preview.Items[0].SourceID},
		ConfirmSnapshotConversion: true,
	}
}

func findingCodes(findings []WriteGroupMigrationFinding) map[string]bool {
	codes := make(map[string]bool, len(findings))
	for _, finding := range findings {
		codes[finding.Code] = true
	}
	return codes
}

func TestPreviewRecordingPlanMigrationUnknownPlanIsSafeNotFound(t *testing.T) {
	ctx := context.Background()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	fixture := seedWriteGroupFixture(ctx, t, db)
	workspaceSvc := NewService(NewSQLRepository(db))
	workspace, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	_, err = workspaceSvc.UpdateDatabaseSetup(ctx, workspace.DatabaseSetupRevision, func(_ context.Context, _ *sql.Tx, record *Record) error {
		record.DatabaseConnectorID = fixture.connectorID
		return nil
	})
	require.NoError(t, err)

	service := NewWriteGroupService(workspaceSvc, NewSQLWriteGroupRepository(db))
	_, err = service.PreviewRecordingPlanMigration(ctx, workspace.ID, []string{"missing-plan"})
	require.ErrorIs(t, err, ErrWriteGroupNotFound)
}

func TestPreviewRecordingPlanMigrationBlocksAndPreservesCompleteIntent(t *testing.T) {
	ctx := t.Context()
	fixture := newRecordingPlanMigrationFixture(ctx, t)
	defer fixture.close()

	preview, err := fixture.service.PreviewRecordingPlanMigration(ctx, fixture.base.workspaceID, []string{fixture.plan.ID})
	require.NoError(t, err)
	require.Len(t, preview.Items, 1)
	item := preview.Items[0]
	require.Equal(t, WriteGroupMigrationPreviewStatusBlocked, item.Status)
	require.Nil(t, item.CandidateGroup)
	require.Equal(t, "open_write_groups", item.RepairAction)
	require.NotEmpty(t, item.SourceRevision)
	require.NotNil(t, item.BeforeRecordingPlanIntent)
	require.Equal(t, fixture.plan.ID, item.BeforeRecordingPlanIntent.SourceID)
	require.Equal(t, fixture.plan.ID, item.BeforeRecordingPlanIntent.Plan.ID)
	require.Equal(t, recordingplan.RawPolicyEverySample, item.BeforeRecordingPlanIntent.Plan.Streams[0].RawPolicy)
	require.Len(t, item.BeforeRecordingPlanIntent.Sources, 1)
	require.Equal(t, recordingPlanSourceStatusResolved, item.BeforeRecordingPlanIntent.Sources[0].Status)
	codes := findingCodes(item.Issues)
	require.True(t, codes["recording-plan-column-binding-missing"])
	require.True(t, codes["recording-plan-row-identity-missing"])
	require.True(t, codes["recording-plan-every-sample-not-equivalent"])

	beforeWorkspace, err := fixture.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	var groups, maps int
	require.NoError(t, fixture.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM write_groups`).Scan(&groups))
	require.NoError(t, fixture.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM write_group_migration_maps`).Scan(&maps))
	_, err = fixture.service.ReviewRecordingPlanMigration(ctx, recordingPlanReviewRequest(preview))
	require.ErrorIs(t, err, ErrWriteGroupPlanMigrationBlocked)
	require.ErrorIs(t, err, ErrWriteGroupValidation)
	var afterGroups, afterMaps int
	require.NoError(t, fixture.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM write_groups`).Scan(&afterGroups))
	require.NoError(t, fixture.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM write_group_migration_maps`).Scan(&afterMaps))
	require.Equal(t, groups, afterGroups)
	require.Equal(t, maps, afterMaps)
	afterWorkspace, err := fixture.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	require.Equal(t, beforeWorkspace.DatabaseSetupRevision, afterWorkspace.DatabaseSetupRevision)
	stored, err := recordingplan.NewSQLRepository(fixture.db).GetPlanByWorkspace(ctx, fixture.plan.ID, fixture.base.workspaceID)
	require.NoError(t, err)
	require.Equal(t, fixture.plan.Revision, stored.Revision)
	require.Equal(t, fixture.plan.Status, stored.Status)
}

func TestPreviewRecordingPlanMigrationUsesEmptySourcesArrayForPlanWithoutReferences(t *testing.T) {
	ctx := t.Context()
	fixture := newRecordingPlanMigrationFixture(ctx, t)
	defer fixture.close()
	plan := *fixture.plan
	plan.ID = "legacy-plan-empty"
	plan.Name = "Empty plan"
	plan.Members = nil
	plan.Streams = nil
	plan.Destinations = nil
	require.NoError(t, recordingplan.NewSQLRepository(fixture.db).CreatePlan(ctx, &plan))

	preview, err := fixture.service.PreviewRecordingPlanMigration(ctx, fixture.base.workspaceID, []string{plan.ID})
	require.NoError(t, err)
	intent := preview.Items[0].BeforeRecordingPlanIntent
	require.NotNil(t, intent)
	require.NotNil(t, intent.Sources)
	require.Empty(t, intent.Sources)
	require.Nil(t, intent.Plan.Members)
	require.Nil(t, intent.Plan.Streams)
	require.Nil(t, intent.Plan.Destinations)
	require.True(t, findingCodes(preview.Items[0].Issues)["recording-plan-source-missing"])
}

func TestRecordingPlanMigrationDigestTracksUnresolvedMeasurementAndDestinationChanges(t *testing.T) {
	ctx := t.Context()
	fixture := newRecordingPlanMigrationFixture(ctx, t)
	defer fixture.close()
	_, err := fixture.db.ExecContext(ctx, `UPDATE mappings SET enabled = 0 WHERE point_id = ? AND tag_id = ?`, fixture.base.pointID, fixture.base.tagID)
	require.NoError(t, err)
	first, err := fixture.service.PreviewRecordingPlanMigration(ctx, fixture.base.workspaceID, []string{fixture.plan.ID})
	require.NoError(t, err)
	require.Equal(t, recordingPlanSourceStatusBlocked, first.Items[0].BeforeRecordingPlanIntent.Sources[0].Status)
	firstSourceRevision := first.Items[0].BeforeRecordingPlanIntent.Sources[0].SourceRevision
	measurementRepo := measurement.NewSQLRepository(fixture.db)
	definition, err := measurementRepo.GetByID(ctx, fixture.measurementID)
	require.NoError(t, err)
	definition.Unit = "kPa"
	require.NoError(t, measurementRepo.Update(ctx, definition))
	second, err := fixture.service.PreviewRecordingPlanMigration(ctx, fixture.base.workspaceID, []string{fixture.plan.ID})
	require.NoError(t, err)
	require.NotEqual(t, first.ReviewDigest, second.ReviewDigest)
	require.NotEqual(t, firstSourceRevision, second.Items[0].BeforeRecordingPlanIntent.Sources[0].SourceRevision)
	_, err = fixture.service.ReviewRecordingPlanMigration(ctx, recordingPlanReviewRequest(first))
	require.ErrorIs(t, err, ErrWriteGroupRevisionConflict)

	thirdBeforeConnector := second.ReviewDigest
	_, err = fixture.db.ExecContext(ctx, `UPDATE database_connectors SET connection_config = ? WHERE id = ?`, `{"database":"changed","schema":"main"}`, fixture.base.connectorID)
	require.NoError(t, err)
	third, err := fixture.service.PreviewRecordingPlanMigration(ctx, fixture.base.workspaceID, []string{fixture.plan.ID})
	require.NoError(t, err)
	require.NotEqual(t, thirdBeforeConnector, third.ReviewDigest)
}

func TestPreviewRecordingPlanMigrationBlocksAmbiguousSourceAcrossTag(t *testing.T) {
	ctx := t.Context()
	fixture := newRecordingPlanMigrationFixture(ctx, t)
	defer fixture.close()
	secondPointID := "point-second"
	_, err := fixture.db.ExecContext(ctx, `
		INSERT INTO points (id, device_id, name, address, function, data_type, mode, enabled)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, secondPointID, fixture.base.deviceID, "temperature-2", "40002", "FC03", schema.DataTypeFloat32, schema.PointModeReadOnly, true)
	require.NoError(t, err)
	_, err = fixture.db.ExecContext(ctx, `
		INSERT INTO mappings (id, point_id, tag_id, transform_pipeline, status, enabled)
		VALUES (?, ?, ?, ?, ?, ?)`, "mapping-second", secondPointID, fixture.base.tagID, `[]`, schema.MappingStatusActive, true)
	require.NoError(t, err)
	preview, err := fixture.service.PreviewRecordingPlanMigration(ctx, fixture.base.workspaceID, []string{fixture.plan.ID})
	require.NoError(t, err)
	codes := findingCodes(preview.Items[0].Issues)
	require.True(t, codes["multiple-enabled-source-mappings"])
	require.Equal(t, recordingPlanSourceStatusBlocked, preview.Items[0].BeforeRecordingPlanIntent.Sources[0].Status)
}

func TestReviewRecordingPlanMigrationRejectsConfirmationAndStaleEnvelope(t *testing.T) {
	ctx := t.Context()
	fixture := newRecordingPlanMigrationFixture(ctx, t)
	defer fixture.close()
	preview, err := fixture.service.PreviewRecordingPlanMigration(ctx, fixture.base.workspaceID, []string{fixture.plan.ID})
	require.NoError(t, err)
	request := recordingPlanReviewRequest(preview)
	request.ConfirmSnapshotConversion = false
	_, err = fixture.service.ReviewRecordingPlanMigration(ctx, request)
	require.ErrorIs(t, err, ErrWriteGroupValidation)
	request = recordingPlanReviewRequest(preview)
	request.ExpectedWorkspaceRevision = "stale-workspace-revision"
	_, err = fixture.service.ReviewRecordingPlanMigration(ctx, request)
	require.ErrorIs(t, err, ErrWriteGroupRevisionConflict)
	var groups, maps int
	require.NoError(t, fixture.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM write_groups`).Scan(&groups))
	require.NoError(t, fixture.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM write_group_migration_maps`).Scan(&maps))
	require.Zero(t, groups)
	require.Zero(t, maps)
}

func TestPreviewRecordingPlanMigrationReturnsSafeNotFoundForForeignPlan(t *testing.T) {
	ctx := t.Context()
	fixture := newRecordingPlanMigrationFixture(ctx, t)
	defer fixture.close()
	foreign := *fixture.plan
	foreign.ID = "foreign-plan"
	foreign.WorkspaceID = "foreign-workspace"
	require.NoError(t, recordingplan.NewSQLRepository(fixture.db).CreatePlan(ctx, &foreign))
	_, err := fixture.service.PreviewRecordingPlanMigration(ctx, fixture.base.workspaceID, []string{foreign.ID})
	require.ErrorIs(t, err, ErrWriteGroupNotFound)
}

func TestPreviewRecordingPlanMigrationBlocksUnsupportedModesAndTargets(t *testing.T) {
	ctx := t.Context()
	fixture := newRecordingPlanMigrationFixture(ctx, t)
	defer fixture.close()
	plan := *fixture.plan
	plan.ID = "legacy-plan-advanced"
	plan.Name = "Advanced plan"
	plan.Streams = []recordingplan.PlanStream{{
		StreamID: "stream-advanced", MeasurementID: fixture.measurementID, Mode: recordingplan.StreamModeLatestOnly,
		DestinationIDs: []string{"destination-A", "destination-B"},
	}}
	plan.Destinations = append(plan.Destinations, recordingplan.PlanDestination{
		DestinationID: "destination-B", ConnectorID: fixture.base.connectorID, ConnectorRevision: "connector-revision-1",
	})
	require.NoError(t, recordingplan.NewSQLRepository(fixture.db).CreatePlan(ctx, &plan))
	preview, err := fixture.service.PreviewRecordingPlanMigration(ctx, fixture.base.workspaceID, []string{plan.ID})
	require.NoError(t, err)
	codes := findingCodes(preview.Items[0].Issues)
	require.True(t, codes["recording-plan-stream-mode-unsupported"])
	require.True(t, codes["recording-plan-multiple-destinations"])
	require.True(t, codes["recording-plan-multiple-targets"])
	require.Nil(t, preview.Items[0].CandidateGroup)
}
